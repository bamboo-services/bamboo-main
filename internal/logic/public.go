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
	"runtime"
	"time"

	apiPublic "github.com/bamboo-services/bamboo-main/api/public"
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

	// 检查 Redis 连接（可选，如果 Redis 不可用不影响基本功能）
	// 注意：暂时注释Redis检查，等待Redis相关工具函数实现

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
			Uptime:      "0m", // 可以计算实际运行时间
			Goroutines:  runtime.NumGoroutine(),
			MemoryUsage: "N/A", // 可以计算实际内存使用
			CPUUsage:    "N/A", // 可以计算实际CPU使用率
		},
	}

	return healthResponse, nil
}

// Ping 简单连通性测试
func (p *PublicLogic) Ping(ctx context.Context) (*apiPublic.PingResponse, *xError.Error) {
	pingResponse := &apiPublic.PingResponse{
		Message:   "pong",
		Timestamp: time.Now(),
	}

	return pingResponse, nil
}
