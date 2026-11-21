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
	"sync"

	"github.com/bamboo-services/bamboo-base-go/application/v2"
	"github.com/bamboo-services/bamboo-base-go/application/v2/builtin"
	"github.com/bamboo-services/bamboo-base-go/foundation/v2/apperr"
	"github.com/bamboo-services/bamboo-base-go/foundation/v2/ctxkey"
	"github.com/bamboo-services/bamboo-base-go/foundation/v2/logging"
	"github.com/bamboo-services/bamboo-base-go/plugins/cache/v2"
	"github.com/bamboo-services/bamboo-base-go/plugins/database/v2"
	"github.com/bamboo-services/bamboo-main/internal/repository"
	"github.com/bamboo-services/bamboo-main/pkg/constants"
)

// Worker 将截图任务管理器封装为受 application.App 托管的 ModeWorker 插件。
type Worker struct {
	mu      sync.Mutex
	manager *Manager
}

// NewWorker 构造截图服务常驻 Worker 插件信息与工厂闭包。
func NewWorker() (application.PluginInfo, func() application.RunnablePlugin) {
	info := application.PluginInfo{
		Key:      constants.PluginScreenshotKey,
		Name:     "bamboo-screenshot-worker",
		Mode:     application.ModeWorker,
		Requires: []ctxkey.Key{ctxkey.DatabaseKey, ctxkey.CacheManagerKey, builtin.LoggerKey, constants.ContextDatabaseMigrated},
		Provides: []ctxkey.Key{constants.ContextScreenshotManager},
	}
	factory := func() application.RunnablePlugin {
		return &Worker{}
	}
	return info, factory
}

// Start 在受控启动上下文中构造 Manager 实例并准备暴露到共享容器。
func (w *Worker) Start(ctx context.Context) *apperr.Error {
	db, appErr := database.Get(ctx)
	if appErr != nil {
		return appErr
	}
	cacheMgr, appErr := cache.Get(ctx)
	if appErr != nil {
		return appErr
	}
	logger, appErr := application.Get[*logging.Logger](ctx, builtin.LoggerKey)
	if appErr != nil {
		return appErr
	}

	cfg := LoadConfig()
	mgr := NewManager(cfg, repository.NewLinkRepo(db, cacheMgr, logger), nil, logger.WithName(logging.NamedTASK, "Screenshot"))

	w.mu.Lock()
	w.manager = mgr
	w.mu.Unlock()
	return nil
}

// Run 阻塞运行截图消费队列，直至应用停机取消。
func (w *Worker) Run(ctx context.Context) *apperr.Error {
	w.mu.Lock()
	mgr := w.manager
	w.mu.Unlock()
	if mgr != nil {
		mgr.Run(ctx)
	}
	return nil
}

// Stop 在逆拓扑停机阶段关闭队列准入并唤醒退出。
func (w *Worker) Stop(_ context.Context) *apperr.Error {
	w.mu.Lock()
	mgr := w.manager
	w.mu.Unlock()
	if mgr != nil {
		mgr.Close()
	}
	return nil
}

// Expose 向应用共享容器发布 *Manager 实例。
func (w *Worker) Expose() map[ctxkey.Key]any {
	w.mu.Lock()
	defer w.mu.Unlock()
	return map[ctxkey.Key]any{
		constants.ContextScreenshotManager: w.manager,
	}
}

// RefreshComponent 为每日定时触发的全量截图入队激发态组件（OverlapSkip 策略）。
type RefreshComponent struct {
	manager *Manager
}

// NewRefreshComponent 构造定时全量截图激发态 Component 插件信息与工厂闭包。
func NewRefreshComponent() (application.PluginInfo, func() application.RunnablePlugin) {
	info := application.PluginInfo{
		Key:      constants.ComponentScreenshotRefreshKey,
		Name:     "bamboo-screenshot-refresh",
		Mode:     application.ModeComponent,
		Requires: []ctxkey.Key{constants.ContextScreenshotManager},
		Overlap:  application.OverlapPolicy{Kind: application.OverlapSkip},
	}
	factory := func() application.RunnablePlugin {
		return &RefreshComponent{}
	}
	return info, factory
}

// Start 提取已启动的截图管理器句柄。
func (c *RefreshComponent) Start(ctx context.Context) *apperr.Error {
	mgr, appErr := application.Get[*Manager](ctx, constants.ContextScreenshotManager)
	if appErr != nil {
		return appErr
	}
	c.manager = mgr
	return nil
}

// Run 执行一次全量已通过友链截图入队后即结束本次激发态生命周期。
func (c *RefreshComponent) Run(ctx context.Context) *apperr.Error {
	if c.manager != nil {
		c.manager.EnqueueAll(ctx)
	}
	return nil
}

// Stop 清理组件引用。
func (c *RefreshComponent) Stop(_ context.Context) *apperr.Error {
	return nil
}
