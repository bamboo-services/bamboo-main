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
	"os"
	"path/filepath"
	"sync"
	"testing"
	"time"

	"github.com/bamboo-services/bamboo-main/internal/entity"
	"github.com/bamboo-services/bamboo-main/pkg/constants"

	xError "github.com/bamboo-services/bamboo-base-go/common/error"
	xSnowflake "github.com/bamboo-services/bamboo-base-go/common/snowflake"
	xModels "github.com/bamboo-services/bamboo-base-go/major/models"
	"gorm.io/gorm"
)

// fakeRepo 内存替身仓储，实现 LinkRepo 接口
type fakeRepo struct {
	mu      sync.Mutex
	links   map[xSnowflake.SnowflakeID]entity.LinkFriend
	updates []xSnowflake.SnowflakeID
	update  chan xSnowflake.SnowflakeID
	getErr  *xError.Error
}

func newFakeRepo() *fakeRepo {
	return &fakeRepo{
		links:  make(map[xSnowflake.SnowflakeID]entity.LinkFriend),
		update: make(chan xSnowflake.SnowflakeID, 64),
	}
}

func (f *fakeRepo) ListApprovedForScreenshot(_ context.Context, _ *gorm.DB) ([]entity.LinkFriend, *xError.Error) {
	f.mu.Lock()
	defer f.mu.Unlock()

	var out []entity.LinkFriend
	for _, link := range f.links {
		if link.Status == constants.LinkStatusApproved && link.IsFailure == constants.LinkFailNormal {
			out = append(out, link)
		}
	}
	return out, nil
}

func (f *fakeRepo) GetByID(_ context.Context, id xSnowflake.SnowflakeID, _ bool, _ *gorm.DB) (*entity.LinkFriend, bool, *xError.Error) {
	f.mu.Lock()
	defer f.mu.Unlock()

	if f.getErr != nil {
		return nil, false, f.getErr
	}
	link, ok := f.links[id]
	if !ok {
		return nil, false, nil
	}
	return &link, true, nil
}

func (f *fakeRepo) UpdateScreenshot(_ context.Context, id xSnowflake.SnowflakeID, _ string, _ time.Time, _ *gorm.DB) *xError.Error {
	f.mu.Lock()
	defer f.mu.Unlock()

	f.updates = append(f.updates, id)
	f.update <- id
	return nil
}

func (f *fakeRepo) waitUpdate(t *testing.T) xSnowflake.SnowflakeID {
	t.Helper()
	select {
	case id := <-f.update:
		return id
	case <-time.After(5 * time.Second):
		t.Fatal("等待截图处理结果超时")
		return 0
	}
}

// recordCapture 记录截图调用顺序的替身截图函数
func recordCapture(order *[]string, mu *sync.Mutex, fail map[string]error) CaptureFunc {
	return func(_ context.Context, url string) ([]byte, error) {
		mu.Lock()
		*order = append(*order, url)
		mu.Unlock()
		if err := fail[url]; err != nil {
			return nil, err
		}
		return []byte("fake-png"), nil
	}
}

func testConfig(t *testing.T) Config {
	t.Helper()
	return Config{
		Enabled: true,
		Dir:     t.TempDir(),
		Timeout: 3 * time.Second,
	}
}

func approvedLink(id xSnowflake.SnowflakeID, url string) entity.LinkFriend {
	return entity.LinkFriend{
		BaseEntity: xModels.BaseEntity{ID: id},
		URL:        url,
		Status:     constants.LinkStatusApproved,
		IsFailure:  constants.LinkFailNormal,
	}
}

// TestManager_EnqueueDedup 验证同一友链重复入队只处理一次
func TestManager_EnqueueDedup(t *testing.T) {
	repo := newFakeRepo()
	const id = xSnowflake.SnowflakeID(1001)
	repo.links[id] = approvedLink(id, "https://a.example")

	var (
		mu    sync.Mutex
		order []string
	)
	m := NewManager(testConfig(t), repo, recordCapture(&order, &mu, nil), nil)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// worker 启动前连续入队，去重应只保留一条
	m.Enqueue(id)
	m.Enqueue(id)
	m.Enqueue(id)
	go m.Run(ctx)

	if got := repo.waitUpdate(t); got != id {
		t.Fatalf("期望更新友链 %d，实际 %d", id, got)
	}
	mu.Lock()
	defer mu.Unlock()
	if len(order) != 1 || order[0] != "https://a.example" {
		t.Fatalf("期望仅截图一次 https://a.example，实际 %v", order)
	}
}

// TestManager_FIFO 验证队列按入队顺序串行处理
func TestManager_FIFO(t *testing.T) {
	repo := newFakeRepo()
	ids := []xSnowflake.SnowflakeID{2001, 2002, 2003}
	urls := map[xSnowflake.SnowflakeID]string{
		ids[0]: "https://one.example",
		ids[1]: "https://two.example",
		ids[2]: "https://three.example",
	}
	for _, id := range ids {
		repo.links[id] = approvedLink(id, urls[id])
	}

	var (
		mu    sync.Mutex
		order []string
	)
	m := NewManager(testConfig(t), repo, recordCapture(&order, &mu, nil), nil)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go m.Run(ctx)

	for _, id := range ids {
		m.Enqueue(id)
	}
	for i := 0; i < len(ids); i++ {
		_ = repo.waitUpdate(t)
	}

	mu.Lock()
	defer mu.Unlock()
	want := []string{"https://one.example", "https://two.example", "https://three.example"}
	if len(order) != len(want) {
		t.Fatalf("期望处理 %d 次，实际 %d（%v）", len(want), len(order), order)
	}
	for i := range want {
		if order[i] != want[i] {
			t.Fatalf("FIFO 顺序不符：期望 %v，实际 %v", want, order)
		}
	}
}

// TestManager_FailureNotBlocking 验证单个截图失败不阻塞后续任务
func TestManager_FailureNotBlocking(t *testing.T) {
	repo := newFakeRepo()
	const (
		failID = xSnowflake.SnowflakeID(3001)
		okID   = xSnowflake.SnowflakeID(3002)
	)
	repo.links[failID] = approvedLink(failID, "https://fail.example")
	repo.links[okID] = approvedLink(okID, "https://ok.example")

	var (
		mu    sync.Mutex
		order []string
	)
	m := NewManager(testConfig(t), repo, recordCapture(&order, &mu, map[string]error{
		"https://fail.example": xError.NewError(context.Background(), xError.DatabaseError, "mock 截图失败", false),
	}), nil)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go m.Run(ctx)

	m.Enqueue(failID)
	m.Enqueue(okID)

	if got := repo.waitUpdate(t); got != okID {
		t.Fatalf("期望只有友链 %d 更新成功，实际 %d", okID, got)
	}
	// 失败任务已被消费（不阻塞队列），成功任务完成更新
	waitFor(t, func() bool {
		mu.Lock()
		defer mu.Unlock()
		return len(order) == 2
	})

	mu.Lock()
	defer mu.Unlock()
	if len(order) != 2 || order[0] != "https://fail.example" || order[1] != "https://ok.example" {
		t.Fatalf("期望按序处理两个任务，实际 %v", order)
	}
	repo.mu.Lock()
	defer repo.mu.Unlock()
	if len(repo.updates) != 1 || repo.updates[0] != okID {
		t.Fatalf("期望仅成功任务更新截图，实际 %v", repo.updates)
	}
}

// TestManager_SkipNonApproved 验证非已通过状态友链跳过截图
func TestManager_SkipNonApproved(t *testing.T) {
	repo := newFakeRepo()
	const id = xSnowflake.SnowflakeID(4001)
	link := approvedLink(id, "https://pending.example")
	link.Status = constants.LinkStatusPending // 待审核
	repo.links[id] = link

	var (
		mu    sync.Mutex
		order []string
	)
	m := NewManager(testConfig(t), repo, recordCapture(&order, &mu, nil), nil)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go m.Run(ctx)

	m.Enqueue(id)
	waitFor(t, func() bool {
		m.mu.Lock()
		defer m.mu.Unlock()
		return len(m.queue) == 0 && len(m.inflight) == 0
	})

	mu.Lock()
	defer mu.Unlock()
	if len(order) != 0 {
		t.Fatalf("非已通过友链不应截图，实际 %v", order)
	}
	repo.mu.Lock()
	defer repo.mu.Unlock()
	if len(repo.updates) != 0 {
		t.Fatalf("非已通过友链不应更新截图，实际 %v", repo.updates)
	}
}

// TestManager_Disabled 验证关闭截图服务后入队与消费全部失效
func TestManager_Disabled(t *testing.T) {
	repo := newFakeRepo()
	const id = xSnowflake.SnowflakeID(5001)
	repo.links[id] = approvedLink(id, "https://disabled.example")

	var (
		mu    sync.Mutex
		order []string
	)
	cfg := testConfig(t)
	cfg.Enabled = false
	m := NewManager(cfg, repo, recordCapture(&order, &mu, nil), nil)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	m.Run(ctx) // 同步返回：disabled 时不启动 worker
	m.Enqueue(id)

	if len(m.queue) != 0 || len(m.inflight) != 0 {
		t.Fatalf("disabled 后队列不应增长，queue=%v inflight=%v", m.queue, m.inflight)
	}
	if len(order) != 0 {
		t.Fatalf("disabled 后不应截图，实际 %v", order)
	}
}

// TestManager_SaveFile 验证原子写入：文件存在、内容一致、无临时文件残留
func TestManager_SaveFile(t *testing.T) {
	const id = xSnowflake.SnowflakeID(6001)
	m := NewManager(testConfig(t), newFakeRepo(), nil, nil)
	data := []byte("png-bytes")

	if err := m.saveFile(id, data); err != nil {
		t.Fatalf("saveFile 失败: %v", err)
	}

	finalPath := filepath.Join(m.cfg.Dir, "6001.png")
	got, err := os.ReadFile(finalPath)
	if err != nil {
		t.Fatalf("读取截图文件失败: %v", err)
	}
	if string(got) != string(data) {
		t.Fatal("截图文件内容与写入不一致")
	}
	if _, err := os.Stat(filepath.Join(m.cfg.Dir, "6001.tmp")); !os.IsNotExist(err) {
		t.Fatal("临时文件未清理，存在残留")
	}
}

// waitFor 轮询等待条件满足（测试辅助）
func waitFor(t *testing.T, cond func() bool) {
	t.Helper()
	deadline := time.Now().Add(5 * time.Second)
	for time.Now().Before(deadline) {
		if cond() {
			return
		}
		time.Sleep(10 * time.Millisecond)
	}
	t.Fatal("等待条件满足超时")
}
