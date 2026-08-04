/*
 * --------------------------------------------------------------------------------
 * Copyright (c) 2016-NOW(至今) 筱锋
 * Author: 筱锋「xiao_lfeng」(https://www.x-lf.com)
 * --------------------------------------------------------------------------------
 * 许可证声明：版权所有 (c) 2016-2026 筱锋。保留所有权利。
 * 有关MIT许可证的更多信息，请查看项目根目录下的LICENSE文件或访问：
 * https://opensource.org/licenses/MIT
 * --------------------------------------------------------------------------------
 */

package logic

import (
	"context"

	xError "github.com/bamboo-services/bamboo-base-go/common/error"
	xLog "github.com/bamboo-services/bamboo-base-go/common/log"
	xCache "github.com/bamboo-services/bamboo-base-go/major/cache"
	"gorm.io/gorm"
)

type logic struct {
	db    *gorm.DB
	cache *xCache.Manager
	log   *xLog.LogNamedLogger
}

// withTx 在独立事务中执行 fn：panic 时回滚，业务错误回滚后透传，成功后提交。
//
// 收敛各 logic 手写的 Begin + defer recover + Commit 样板，保证回滚路径一致。
func (l *logic) withTx(ctx context.Context, fn func(tx *gorm.DB) *xError.Error) *xError.Error {
	tx := l.db.WithContext(ctx).Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if xErr := fn(tx); xErr != nil {
		tx.Rollback()
		return xErr
	}

	if err := tx.Commit().Error; err != nil {
		return xError.NewError(ctx, xError.DatabaseError, "提交事务失败", false, err)
	}

	return nil
}
