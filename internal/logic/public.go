/*
 * --------------------------------------------------------------------------------
 * Copyright (c) 2016-NOW(至今) 筱锋
 * Author: 筱锋「xiao_lfeng」(https://www.x-lf.com)
 * --------------------------------------------------------------------------------
 * 许可证声明：版权所有 (c) 2016-2025 筱锋。保留所有权利。
 * 有关MIT许可证的更多信息，请查看项目根目录下的LICENSE文件或访问：
 * https://opensource.org/licenses/MIT
 * --------------------------------------------------------------------------------
 */

package logic

import (
	"context"
	"runtime"
	"time"

	apiPublic "github.com/bamboo-services/bamboo-main/api/public"

	xError "github.com/bamboo-services/bamboo-base-go/common/error"
	xLog "github.com/bamboo-services/bamboo-base-go/common/log"
	xCtxUtil "github.com/bamboo-services/bamboo-base-go/common/utility/context"
	"github.com/gin-gonic/gin"
)

// PublicLogic 公开接口业务逻辑
type PublicLogic struct {
	logic
}

func NewPublicLogic(ctx context.Context) *PublicLogic {
	db := xCtxUtil.MustGetDB(ctx)
	rdb := xCtxUtil.MustGetRDB(ctx)

	return &PublicLogic{
		logic: logic{
			db:  db,
			rdb: rdb,
			log: xLog.WithName(xLog.NamedLOGC, "PublicLogic"),
		},
	}
}

// HealthCheck 健康检查
func (p *PublicLogic) HealthCheck(ctx *gin.Context) (*apiPublic.HealthResponse, *xError.Error) {
	// 检查数据库连接
	db := p.db
	if db == nil {
		return nil, xError.NewError(ctx, xError.DatabaseError, "数据库连接失败", false)
	}

	// 执行简单查询测试数据库连接
	sqlDB, err := db.DB()
	if err != nil {
		return nil, xError.NewError(ctx, xError.DatabaseError, "获取数据库连接失败", false, err)
	}

	if err := sqlDB.Ping(); err != nil {
		return nil, xError.NewError(ctx, xError.DatabaseError, "数据库连接测试失败", false, err)
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
func (p *PublicLogic) Ping(ctx *gin.Context) (*apiPublic.PingResponse, *xError.Error) {
	pingResponse := &apiPublic.PingResponse{
		Message:   "pong",
		Timestamp: time.Now(),
	}

	return pingResponse, nil
}
