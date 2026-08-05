// --------------------------------------------------------------------------------
// Copyright (c) 2016-NOW(至今) 筱锋
// Author: 筱锋「xiao_lfeng」(https://www.x-lf.com)
// --------------------------------------------------------------------------------
// 许可证声明：版权所有 (c) 2016-2026 筱锋。保留所有权利。
// 有关MIT许可证的更多信息，请查看项目根目录下的LICENSE文件或访问：
// https://opensource.org/licenses/MIT
// --------------------------------------------------------------------------------

package screenshot

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/bamboo-services/bamboo-main/internal/entity"
	"github.com/bamboo-services/bamboo-main/pkg/constants"

	xError "github.com/bamboo-services/bamboo-base-go/common/error"
	xLog "github.com/bamboo-services/bamboo-base-go/common/log"
	xSnowflake "github.com/bamboo-services/bamboo-base-go/common/snowflake"
	xCtxUtil "github.com/bamboo-services/bamboo-base-go/major/utility/context"
	"gorm.io/gorm"
)

// LinkRepo 截图服务所需的友链仓储能力（由 repository.LinkRepo 实现，测试可注入替身）
type LinkRepo interface {
	// ListApprovedForScreenshot 查询全部已通过且未失效的友链（按队列顺序）
	ListApprovedForScreenshot(ctx context.Context, tx *gorm.DB) ([]entity.LinkFriend, *xError.Error)
	// GetByID 查询友链详情
	GetByID(ctx context.Context, id xSnowflake.SnowflakeID, withAssociations bool, tx *gorm.DB) (*entity.LinkFriend, bool, *xError.Error)
	// UpdateScreenshot 更新友链截图信息
	UpdateScreenshot(ctx context.Context, id xSnowflake.SnowflakeID, url string, at time.Time, tx *gorm.DB) *xError.Error
}

// Manager 友链站点截图任务管理器。
//
// 内存 FIFO 队列 + 单 worker 串行消费，天然满足「按队列顺序自动处理」；
// 任务入队时经 inflight map 去重，避免同一友链被重复排队。
// worker 常驻运行并随 Runner 生命周期退出；进程重启丢失的任务由每日全量入队兜底。
type Manager struct {
	cfg     Config
	repo    LinkRepo
	capture CaptureFunc
	log     *xLog.LogNamedLogger

	mu       sync.Mutex
	queue    []xSnowflake.SnowflakeID
	inflight map[xSnowflake.SnowflakeID]struct{}
	wake     chan struct{}
	closed   bool
}

// NewManager 创建截图任务管理器；capture 为 nil 时使用 rod 默认实现
func NewManager(cfg Config, repo LinkRepo, capture CaptureFunc, log *xLog.LogNamedLogger) *Manager {
	if capture == nil {
		capture = NewRodCapture(cfg)
	}
	if log == nil {
		log = xLog.WithName(xLog.NamedTASK, "Screenshot")
	}
	return &Manager{
		cfg:      cfg,
		repo:     repo,
		capture:  capture,
		log:      log,
		inflight: make(map[xSnowflake.SnowflakeID]struct{}),
		wake:     make(chan struct{}, 1),
	}
}

// Enqueue 入队单个友链截图任务（非阻塞；已在队列或处理中则忽略）
func (m *Manager) Enqueue(id xSnowflake.SnowflakeID) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.closed {
		return
	}
	if _, ok := m.inflight[id]; ok {
		return
	}
	m.inflight[id] = struct{}{}
	m.queue = append(m.queue, id)
	m.wakeUp()
}

// EnqueueAll 全量入队：查询全部已通过且未失效的友链，排队依次截图
func (m *Manager) EnqueueAll(ctx context.Context) {
	links, xErr := m.repo.ListApprovedForScreenshot(ctx, nil)
	if xErr != nil {
		m.log.SugarError(ctx, "全量截图入队失败：查询友链列表出错", "error", xErr.Error())
		return
	}
	for _, link := range links {
		m.Enqueue(link.ID)
	}
	m.log.SugarInfo(ctx, "全量截图任务已入队", "count", len(links))
}

// Run 启动 worker 常驻消费队列，ctx 取消时退出
func (m *Manager) Run(ctx context.Context) {
	if !m.cfg.Enabled {
		m.log.SugarInfo(ctx, "截图服务已禁用（SCREENSHOT_ENABLED=false），worker 不启动")
		m.mu.Lock()
		m.closed = true
		m.mu.Unlock()
		return
	}
	if err := os.MkdirAll(m.cfg.Dir, 0o755); err != nil {
		m.log.SugarWarn(ctx, "截图存储目录创建失败", "dir", m.cfg.Dir, "error", err.Error())
	}
	m.log.SugarInfo(ctx, "截图 worker 已启动", "dir", m.cfg.Dir)

	for {
		m.mu.Lock()
		if len(m.queue) == 0 {
			m.mu.Unlock()
			select {
			case <-m.wake:
				continue
			case <-ctx.Done():
				m.mu.Lock()
				m.closed = true
				m.mu.Unlock()
				return
			}
		}
		id := m.queue[0]
		m.queue = m.queue[1:]
		m.mu.Unlock()

		m.process(ctx, id)

		m.mu.Lock()
		delete(m.inflight, id)
		m.mu.Unlock()
	}
}

// process 处理单个友链截图任务
func (m *Manager) process(ctx context.Context, id xSnowflake.SnowflakeID) {
	linkID := fmt.Sprintf("%d", id)

	link, found, xErr := m.repo.GetByID(ctx, id, false, nil)
	if xErr != nil {
		m.log.SugarWarn(ctx, "截图任务跳过：查询友链失败", "link_id", linkID, "error", xErr.Error())
		return
	}
	if !found {
		m.log.SugarWarn(ctx, "截图任务跳过：友链不存在", "link_id", linkID)
		return
	}
	if link.Status != constants.LinkStatusApproved {
		m.log.SugarInfo(ctx, "截图任务跳过：友链非已通过状态", "link_id", linkID, "status", link.Status)
		return
	}

	data, err := m.capture(ctx, link.URL)
	if err != nil {
		// 截图失败保留旧截图，等待下次（每日全量或再次触发）自然重试
		m.log.SugarWarn(ctx, "站点截图失败，保留旧截图", "link_id", linkID, "url", link.URL, "error", err.Error())
		return
	}
	if err = m.saveFile(id, data); err != nil {
		m.log.SugarWarn(ctx, "截图文件写入失败", "link_id", linkID, "error", err.Error())
		return
	}

	path := fmt.Sprintf("/screenshots/%d.png", id)
	if xErr = m.repo.UpdateScreenshot(ctx, id, path, time.Now(), nil); xErr != nil {
		m.log.SugarWarn(ctx, "截图信息更新失败", "link_id", linkID, "error", xErr.Error())
		return
	}
	m.log.SugarInfo(ctx, "截图成功",
		"link_id", linkID,
		"url", link.URL,
		"path", path,
		"size", len(data),
	)
}

// saveFile 原子写入截图文件：先写临时文件再 rename 替换，避免半写文件被读取
func (m *Manager) saveFile(id xSnowflake.SnowflakeID, data []byte) error {
	if err := os.MkdirAll(m.cfg.Dir, 0o755); err != nil {
		return err
	}
	tmpPath := filepath.Join(m.cfg.Dir, fmt.Sprintf("%d.tmp", id))
	finalPath := filepath.Join(m.cfg.Dir, fmt.Sprintf("%d.png", id))
	if err := os.WriteFile(tmpPath, data, 0o644); err != nil {
		return err
	}
	return os.Rename(tmpPath, finalPath)
}

// wakeUp 唤醒 worker（有缓冲的 channel，不阻塞）
func (m *Manager) wakeUp() {
	select {
	case m.wake <- struct{}{}:
	default:
	}
}

// GetManager 从上下文获取截图任务管理器实例
func GetManager(ctx context.Context) *Manager {
	manager, err := xCtxUtil.Get[*Manager](ctx, constants.ContextScreenshotManager)
	if err != nil {
		return nil
	}
	return manager
}
