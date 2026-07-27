// --------------------------------------------------------------------------------
// Copyright (c) 2016-NOW(至今) 筱锋
// Author: 筱锋「xiao_lfeng」(https://www.x-lf.com)
// --------------------------------------------------------------------------------
// 许可证声明：版权所有 (c) 2016-2026 筱锋。保留所有权利。
// 有关MIT许可证的更多信息，请查看项目根目录下的LICENSE文件或访问：
// https://opensource.org/licenses/MIT
// --------------------------------------------------------------------------------

package logic

import (
	"context"
	"fmt"
	"runtime"
	"time"

	apiPublic "github.com/bamboo-services/bamboo-main/api/public"
	"github.com/bamboo-services/bamboo-main/internal/metrics"
	"github.com/bamboo-services/bamboo-main/internal/repository"

	xError "github.com/bamboo-services/bamboo-base-go/common/error"
	xLog "github.com/bamboo-services/bamboo-base-go/common/log"
	xCtxUtil "github.com/bamboo-services/bamboo-base-go/major/utility/context"
)

type publicRepo struct {
	health *repository.HealthRepo
}

// PublicLogic 公开接口业务逻辑
type PublicLogic struct {
	logic
	repo publicRepo
}

// NewPublicLogic 创建 PublicLogic 实例，从上下文获取数据库与缓存并初始化健康检查仓储依赖。
func NewPublicLogic(ctx context.Context) *PublicLogic {
	db := xCtxUtil.MustGetDB(ctx)
	m := xCtxUtil.MustGetCacheManager(ctx)

	return &PublicLogic{
		logic: logic{
			db:    db,
			cache: m,
			log:   xLog.WithName(xLog.NamedLOGC, "PublicLogic"),
		},
		repo: publicRepo{
			health: repository.NewHealthRepo(db, m),
		},
	}
}

// HealthCheck 健康检查
func (p *PublicLogic) HealthCheck(ctx context.Context) (*apiPublic.HealthResponse, *xError.Error) {
	// 检查数据库连接
	if _, xErr := p.repo.health.DatabaseReady(ctx); xErr != nil {
		return nil, xErr
	}

	// 采集真实运行时指标（运行时长 / 协程数 / 内存 / CPU）
	snap := metrics.Snapshot()

	// 构建健康检查响应
	healthResponse := &apiPublic.HealthResponse{
		Status:    "healthy",
		Timestamp: time.Now(),
		System: apiPublic.SystemInfo{
			Version:     "v1.0.0",
			Environment: "development", // 可以从配置获取
			Platform:    runtime.GOOS,
			GoVersion:   runtime.Version(),
		},
		Runtime: apiPublic.RuntimeInfo{
			Uptime:      formatDuration(snap.Uptime),
			Goroutines:  snap.Goroutines,
			MemoryUsage: formatBytes(snap.MemoryBytes),
			CPUUsage:    fmt.Sprintf("%.2f%%", snap.CPUPercent),
		},
	}

	return healthResponse, nil
}

// formatDuration 将运行时长格式化为紧凑形式（如 1h2m3s）
func formatDuration(d time.Duration) string {
	d = d.Round(time.Second)
	h := int(d.Hours())
	m := int(d.Minutes()) % 60
	s := int(d.Seconds()) % 60
	switch {
	case h > 0:
		return fmt.Sprintf("%dh%dm%ds", h, m, s)
	case m > 0:
		return fmt.Sprintf("%dm%ds", m, s)
	default:
		return fmt.Sprintf("%ds", s)
	}
}

// formatBytes 将字节数格式化为人类可读形式（B/KB/MB/GB）
func formatBytes(bytes uint64) string {
	const (
		kb = 1024
		mb = 1024 * kb
		gb = 1024 * mb
	)
	switch {
	case bytes >= gb:
		return fmt.Sprintf("%.2f GB", float64(bytes)/float64(gb))
	case bytes >= mb:
		return fmt.Sprintf("%.2f MB", float64(bytes)/float64(mb))
	case bytes >= kb:
		return fmt.Sprintf("%.2f KB", float64(bytes)/float64(kb))
	default:
		return fmt.Sprintf("%d B", bytes)
	}
}

// Ping 简单连通性测试
func (p *PublicLogic) Ping(ctx context.Context) (*apiPublic.PingResponse, *xError.Error) {
	pingResponse := &apiPublic.PingResponse{
		Message:   "pong",
		Timestamp: time.Now(),
	}

	return pingResponse, nil
}
