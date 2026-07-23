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

package main

import (
	_ "github.com/bamboo-services/bamboo-main/docs"
	"github.com/bamboo-services/bamboo-main/internal/app/route"
	"github.com/bamboo-services/bamboo-main/internal/app/startup"
	"github.com/bamboo-services/bamboo-main/internal/app/startup/prepare"
	"github.com/bamboo-services/bamboo-main/internal/app/startup/worker"
	"github.com/bamboo-services/bamboo-main/internal/entity"

	xLog "github.com/bamboo-services/bamboo-base-go/common/log"
	xMain "github.com/bamboo-services/bamboo-base-go/major/main"
	xOption "github.com/bamboo-services/bamboo-base-go/major/option"
	xOptionCache "github.com/bamboo-services/bamboo-base-go/major/option/cache"
	xOptionDB "github.com/bamboo-services/bamboo-base-go/major/option/database"
	xReg "github.com/bamboo-services/bamboo-base-go/major/register"
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
			xOptionDB.WithPrepare(prepare.DefaultData),
		),
		xOption.WithCache(xOptionCache.FromEnv()),
		xOption.WithRoute(route.NewRoute),
	}

	reg := xReg.Register(ctx, nodeList, opts...)
	log := xLog.WithName(xLog.NamedMAIN)

	xMain.Runner(reg, log, worker.MailWorkerRunner)
}
