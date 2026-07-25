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

package worker

import (
	"context"
	"fmt"

	"github.com/bamboo-services/bamboo-main/internal/models/base"
	"github.com/bamboo-services/bamboo-main/internal/task"
	"github.com/bamboo-services/bamboo-main/pkg/constants"

	xLog "github.com/bamboo-services/bamboo-base-go/common/log"
	xCtxUtil "github.com/bamboo-services/bamboo-base-go/major/utility/context"
	"go.uber.org/zap"
)

// MailWorkerRunner 启动邮件 worker 并阻塞直至上下文取消，随后执行优雅停止。
func MailWorkerRunner(ctx context.Context, _ ...any) {
	log := xLog.WithName(xLog.NamedMAIN, "MailWorker")
	config, err := xCtxUtil.Get[*base.BambooConfig](ctx, constants.ContextCustomConfig)
	if err != nil {
		panic(fmt.Sprintf("获取配置失败: %v", err))
	}

	rdb := xCtxUtil.MustGetRDB(ctx)
	worker := task.NewMailWorker(rdb, &config.Email, zap.NewNop().Sugar())

	worker.Start()

	<-ctx.Done()
	worker.Stop()
	log.Info(ctx, "邮件 worker 已优雅停止")
}
