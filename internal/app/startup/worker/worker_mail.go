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

package worker

import (
	"context"
	"fmt"

	"github.com/bamboo-services/bamboo-main/internal/models/base"
	"github.com/bamboo-services/bamboo-main/internal/task"
	"github.com/bamboo-services/bamboo-main/pkg/constants"

	xLog "github.com/bamboo-services/bamboo-base-go/log"
	xCtxUtil "github.com/bamboo-services/bamboo-base-go/utility/ctxutil"
	"go.uber.org/zap"
)

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
	return
}
