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
	"io"
	"log/slog"
	"sync"

	"github.com/bamboo-services/bamboo-base-go/application/v2"
	"github.com/bamboo-services/bamboo-base-go/application/v2/builtin"
	"github.com/bamboo-services/bamboo-base-go/foundation/v2/apperr"
	"github.com/bamboo-services/bamboo-base-go/foundation/v2/logging"
	"github.com/bamboo-services/bamboo-base-go/plugins/cache/v2"
	"github.com/bamboo-services/bamboo-base-go/plugins/database/v2"
	"gorm.io/gorm"
)

var (
	logicFallbackLoggerOnce sync.Once
	logicFallbackLogger     *logging.Logger
)

type logic struct {
	db     *gorm.DB
	cache  *cache.Manager
	logger *logging.Logger
	log    *logging.LogNamedLogger
}

func resolveLogicNamedLogger(ctx context.Context, optionName ...string) *logging.LogNamedLogger {
	if ctx != nil {
		if l, err := application.Get[*logging.Logger](ctx, builtin.LoggerKey); err == nil && l != nil {
			return l.WithName(logging.NamedLOGC, optionName...)
		}
	}
	logicFallbackLoggerOnce.Do(func() {
		l, _ := logging.New(logging.NewLogHandler(logging.HandlerConfig{
			Console: io.Discard,
			Level:   slog.LevelInfo,
		}))
		logicFallbackLogger = l
	})
	return logicFallbackLogger.WithName(logging.NamedLOGC, optionName...)
}

// withTx 在独立事务中执行 fn：panic 时回滚，业务错误回滚后透传，成功后提交。
func (l *logic) withTx(ctx context.Context, fn func(tx *gorm.DB) *apperr.Error) *apperr.Error {
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
		return apperr.New(database.TransactionFailed, "提交事务失败", err)
	}

	return nil
}
