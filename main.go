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
	"fmt"
	"os"
	"time"

	_ "github.com/bamboo-services/bamboo-main/docs"
	"github.com/bamboo-services/bamboo-main/internal/app/route"
	"github.com/bamboo-services/bamboo-main/internal/app/startup"
	"github.com/bamboo-services/bamboo-main/internal/entity"
	"github.com/bamboo-services/bamboo-main/internal/service/screenshot"
	"github.com/bamboo-services/bamboo-main/pkg/constants"
	"github.com/bamboo-services/bamboo-main/templates"

	"github.com/bamboo-services/bamboo-base-go/application/v2"
	"github.com/bamboo-services/bamboo-base-go/application/v2/builtin"
	"github.com/bamboo-services/bamboo-base-go/foundation/v2/ctxkey"
	"github.com/bamboo-services/bamboo-base-go/foundation/v2/env"
	"github.com/bamboo-services/bamboo-base-go/plugins/async/v2"
	"github.com/bamboo-services/bamboo-base-go/plugins/cache/v2"
	"github.com/bamboo-services/bamboo-base-go/plugins/cron/v2"
	postgres "github.com/bamboo-services/bamboo-base-go/plugins/database/adapter/postgres/v2"
	"github.com/bamboo-services/bamboo-base-go/plugins/database/v2"
	"github.com/bamboo-services/bamboo-base-go/plugins/email/v2"
	webgin "github.com/bamboo-services/bamboo-base-go/plugins/web-gin/v2"
	"github.com/joho/godotenv"
	bSdkConst "github.com/phalanx-labs/beacon-sso-sdk/v2/constant"
	bSdkStartup "github.com/phalanx-labs/beacon-sso-sdk/v2/startup"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

func main() {
	// v2 底座不自动加载 .env，入口显式加载（不存在时忽略，兼容容器环境注入）
	_ = godotenv.Load()

	app := application.New()
	regTool := mustRegisterTool(app)
	regWorker := mustRegisterWorker(app)
	regComponent := mustRegisterComponent(app)

	// 1. 注册内置核心工具：受管日志与雪花发号器
	regTool(builtin.NewLogger(startup.BuildLoggerConfig()))
	dcID, nodeID := startup.ParseSnowflakeIDs()
	regTool(builtin.NewSnowflake(dcID, nodeID))

	// 2. 注册基础设施工具：配置、数据库（PostgreSQL）、缓存、邮件（内嵌模板）、异步池、SSO
	regTool(startup.NewConfigTool())

	tablePrefix := "bm_"
	if val, ok := env.Lookup(env.DatabasePrefix); ok && !val.IsEmpty() {
		tablePrefix = val.String()
	}
	regTool(database.New(
		database.WithAdapter(postgres.FromEnv()),
		database.WithGormConfig(&gorm.Config{
			NamingStrategy: schema.NamingStrategy{
				TablePrefix: tablePrefix,
			},
		}),
	))

	regTool(cache.New(cache.FromEnv()))
	regTool(email.New(
		email.FromEnv(),
		email.WithTemplateFS(templates.Mail(), "mail/*.html"),
	))
	regTool(async.New())
	regTool(bSdkStartup.New(bSdkStartup.WithRequires(builtin.LoggerKey)))

	// 3. 注册数据库迁移与默认种子工具
	regTool(startup.NewMigrationTool(
		&entity.SystemUser{},
		&entity.LinkGroup{},
		&entity.LinkColor{},
		&entity.LinkFriend{},
		&entity.SystemLog{},
		&entity.System{},
		&entity.SponsorChannel{},
		&entity.SponsorRecord{},
	))

	// 4. 注册激发态组件：每日全量友链截图任务
	regComponent(screenshot.NewRefreshComponent())

	// 5. 注册常驻 Worker：截图消费队列、定时调度器、Gin Web 服务
	regWorker(screenshot.NewWorker())

	location, err := time.LoadLocation("Asia/Shanghai")
	if err != nil {
		location = time.FixedZone("Asia/Shanghai", 8*3600)
	}
	regWorker(cron.New(
		cron.WithActivator(app.Activator()),
		cron.WithLocation(location),
		cron.WithCronJob("0 0 * * *", constants.ComponentScreenshotRefreshKey),
	))

	regWorker(webgin.New(
		webgin.WithRequires(
			builtin.LoggerKey,
			ctxkey.SnowflakeNodeKey,
			ctxkey.DatabaseKey,
			ctxkey.CacheManagerKey,
			ctxkey.EmailClientKey,
			async.PoolKey,
			constants.ContextCustomConfig,
			constants.ContextDatabaseMigrated,
			constants.ContextScreenshotManager,
			bSdkConst.CtxOAuthConfig,
			bSdkConst.CtxOAuthUserinfoURI,
			bSdkConst.CtxSsoClient,
		),
		webgin.WithRoute(route.NewRoute),
	))

	if runErr := app.Run(); runErr != nil {
		fmt.Fprintf(os.Stderr, "bamboo-main 退出异常: %v\n", runErr)
		os.Exit(1)
	}
}

func mustRegisterTool(app *application.App) func(application.PluginInfo, func() application.Plugin) {
	return func(info application.PluginInfo, factory func() application.Plugin) {
		if err := app.RegisterTool(info, factory); err != nil {
			fmt.Fprintf(os.Stderr, "注册 Tool [%s] 失败: %v\n", info.Key, err)
			os.Exit(1)
		}
	}
}

func mustRegisterWorker(app *application.App) func(application.PluginInfo, func() application.RunnablePlugin) {
	return func(info application.PluginInfo, factory func() application.RunnablePlugin) {
		if err := app.RegisterWorker(info, factory); err != nil {
			fmt.Fprintf(os.Stderr, "注册 Worker [%s] 失败: %v\n", info.Key, err)
			os.Exit(1)
		}
	}
}

func mustRegisterComponent(app *application.App) func(application.PluginInfo, func() application.RunnablePlugin) {
	return func(info application.PluginInfo, factory func() application.RunnablePlugin) {
		if err := app.RegisterComponent(info, factory); err != nil {
			fmt.Fprintf(os.Stderr, "注册 Component [%s] 失败: %v\n", info.Key, err)
			os.Exit(1)
		}
	}
}
