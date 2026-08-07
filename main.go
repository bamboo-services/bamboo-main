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

package main

import (
	"context"
	"time"

	_ "github.com/bamboo-services/bamboo-main/docs"
	"github.com/bamboo-services/bamboo-main/internal/app/route"
	"github.com/bamboo-services/bamboo-main/internal/app/startup"
	"github.com/bamboo-services/bamboo-main/internal/app/startup/prepare"
	"github.com/bamboo-services/bamboo-main/internal/entity"
	"github.com/bamboo-services/bamboo-main/internal/service/screenshot"

	xLog "github.com/bamboo-services/bamboo-base-go/common/log"
	xMain "github.com/bamboo-services/bamboo-base-go/major/main"
	xOption "github.com/bamboo-services/bamboo-base-go/major/option"
	xOptionCache "github.com/bamboo-services/bamboo-base-go/major/option/cache"
	xOptionDB "github.com/bamboo-services/bamboo-base-go/major/option/database"
	xReg "github.com/bamboo-services/bamboo-base-go/major/register"
	xCron "github.com/bamboo-services/bamboo-base-go/plugins/cron"
	xCronRunner "github.com/bamboo-services/bamboo-base-go/plugins/cron/runner"

	// v1.2.0 起数据库驱动插件化：DATABASE_DRIVER=postgres 需空白导入对应驱动插件
	_ "github.com/bamboo-services/bamboo-base-go/plugins/database/postgres"
)

func main() {
	ctx, nodeList := startup.Init()

	opts := []xOption.Option{
		xOption.WithDatabase(
			xOptionDB.FromEnv(),
			xOptionDB.WithTablePrefix("bm_"),
			xOptionDB.WithAutoMigrate(
				&entity.SystemUser{},
				&entity.LinkGroup{},
				&entity.LinkColor{},
				&entity.LinkFriend{},
				&entity.SystemLog{},
				&entity.System{},
				&entity.SponsorChannel{},
				&entity.SponsorRecord{},
			),
			xOptionDB.WithPrepare(
				prepare.DefaultUser,
				prepare.DefaultInfo,
			),
		),
		xOption.WithCache(xOptionCache.FromEnv()),
		xOption.WithRoute(route.NewRoute),
	}

	// 定时任务时区：固定 Asia/Shanghai（每日 0 点全量更新站点截图）
	location, err := time.LoadLocation("Asia/Shanghai")
	if err != nil {
		location = time.FixedZone("Asia/Shanghai", 8*3600)
	}

	xMain.Runner(
		xReg.Register(ctx, nodeList, opts...),
		xLog.WithName(xLog.NamedMAIN),
		// 每日 0 点全量刷新友链站点截图（cron 插件挂载到 Runner 附加协程）
		xCronRunner.New(
			xCronRunner.WithRegister(
				xCron.NewJob("0 0 * * *", func(ctx context.Context) {
					if manager := screenshot.GetManager(ctx); manager != nil {
						manager.EnqueueAll(ctx)
					}
				}),
			),
			xCronRunner.WithLocation(location),
		),
		// 截图任务队列 worker 常驻协程，随 Runner 生命周期启停
		func(ctx context.Context, _ ...any) {
			if manager := screenshot.GetManager(ctx); manager != nil {
				manager.Run(ctx)
			}
		},
	)
}
