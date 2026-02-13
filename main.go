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
	"context"

	_ "github.com/bamboo-services/bamboo-main/docs"
	"github.com/bamboo-services/bamboo-main/internal/app/route"
	"github.com/bamboo-services/bamboo-main/internal/app/startup"

	xLog "github.com/bamboo-services/bamboo-base-go/log"
	xMain "github.com/bamboo-services/bamboo-base-go/main"
	xReg "github.com/bamboo-services/bamboo-base-go/register"
)

// @title BambooMain
// @version v1.0.0
// @description 友情链接管理系统
// @termsOfService https://www.aiawaken.top/
// @contact.name 筱锋 xiao_lfeng
// @contact.url https://www.x-lf.com/
// @contact.email gm@x-lf.cn
// @host localhost:23333
// @BasePath /api/v1
func main() {
	reg := xReg.Register(startup.Init())
	log := xLog.WithName(xLog.NamedMAIN)

	xMain.Runner(reg, log, route.NewRoute, func(ctx context.Context, _ ...any) {
		if worker := startup.GetMailWorker(ctx); worker != nil {
			worker.Stop()
		}
	})
}
