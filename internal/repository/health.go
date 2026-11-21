// --------------------------------------------------------------------------------
// Copyright (c) 2016-NOW(至今) 筱锋
// Author: 筱锋「xiao_lfeng」(https://www.x-lf.com)
// --------------------------------------------------------------------------------
// 许可证声明：版权所有 (c) 2016-2026 筱锋。保留所有权利。
// 有关MIT许可证的更多信息，请查看项目根目录下的LICENSE文件或访问：
// https://opensource.org/licenses/MIT
// --------------------------------------------------------------------------------

package repository

import (
	"context"

	"github.com/bamboo-services/bamboo-base-go/foundation/v2/apperr"
	"github.com/bamboo-services/bamboo-base-go/foundation/v2/logging"
	"github.com/bamboo-services/bamboo-base-go/plugins/cache/v2"
	"github.com/bamboo-services/bamboo-base-go/plugins/database/v2"
	"gorm.io/gorm"
)

// HealthRepo 健康检查仓储
//
// 收口数据库连通性探测，使 logic 层不再直接持有 *gorm.DB 并执行 Ping。
type HealthRepo struct {
	db  *gorm.DB
	log *logging.LogNamedLogger
}

// NewHealthRepo 创建 HealthRepo 实例
func NewHealthRepo(db *gorm.DB, _ *cache.Manager, logger ...*logging.Logger) *HealthRepo {
	return &HealthRepo{
		db:  db,
		log: resolveRepoLogger("HealthRepo", logger...),
	}
}

// DatabaseReady 探测数据库连接是否可用
func (r *HealthRepo) DatabaseReady(ctx context.Context) (bool, *apperr.Error) {
	if r.db == nil {
		return false, apperr.New(database.DatabaseError, "数据库连接失败")
	}

	sqlDB, err := r.db.WithContext(ctx).DB()
	if err != nil {
		return false, apperr.New(database.DatabaseError, "获取数据库连接失败", err)
	}

	if err := sqlDB.PingContext(ctx); err != nil {
		return false, apperr.New(database.DatabaseError, "数据库连接测试失败", err)
	}

	return true, nil
}
