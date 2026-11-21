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

package repository

import (
	"io"
	"log/slog"
	"sync"

	"github.com/bamboo-services/bamboo-base-go/foundation/v2/logging"
)

var (
	repoFallbackLoggerOnce sync.Once
	repoFallbackLogger     *logging.Logger
)

// resolveRepoLogger 从可选注入的 *logging.Logger 派生仓储命名日志器；
// 未传或为 nil 时回退至静默日志器（兼容测试替身与独立调用）。
func resolveRepoLogger(optionName string, logger ...*logging.Logger) *logging.LogNamedLogger {
	if len(logger) > 0 && logger[0] != nil {
		return logger[0].WithName(logging.NamedREPO, optionName)
	}
	repoFallbackLoggerOnce.Do(func() {
		l, _ := logging.New(logging.NewLogHandler(logging.HandlerConfig{
			Console: io.Discard,
			Level:   slog.LevelInfo,
		}))
		repoFallbackLogger = l
	})
	return repoFallbackLogger.WithName(logging.NamedREPO, optionName)
}
