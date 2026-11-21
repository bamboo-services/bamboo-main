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

package route

import (
	"context"

	"github.com/bamboo-services/bamboo-base-go/foundation/v2/apperr"
	"github.com/bamboo-services/bamboo-base-go/foundation/v2/env"
	"github.com/bamboo-services/bamboo-base-go/plugins/web-gin/v2/httpresult"
	"github.com/bamboo-services/bamboo-main/internal/metrics"
	"github.com/bamboo-services/bamboo-main/internal/service/screenshot"
	"github.com/gin-gonic/gin"
	bSdkRoute "github.com/phalanx-labs/beacon-sso-sdk/v2/route"
)

type route struct {
	engine  *gin.Engine
	context context.Context
}

// NewRoute 路由注册器，装配业务中间件链与各领域子路由（Recovery/Response/CORS/Binder 已由 webgin 插件挂载）
func NewRoute(ctx context.Context, serve *gin.Engine) {
	r := &route{
		engine:  serve,
		context: ctx,
	}

	r.engine.NoMethod(func(c *gin.Context) {
		httpresult.Error(c, apperr.New(httpresult.MethodNotAllowed, "请求方法不被允许"))
	})
	r.engine.NoRoute(noRouteHandler)

	r.engine.Use(metrics.Middleware())

	if debugVal, ok := env.Lookup(env.Debug); ok && !debugVal.IsEmpty() {
		if debug, err := debugVal.Bool(); err == nil && debug {
			swaggerRegister(r.engine)
		}
	}

	// Prometheus 指标端点（根路径，公开无鉴权；promhttp 直写，不经 envelope.Response）
	r.engine.GET("/metrics", metrics.Handler())

	// 友链站点截图静态资源：运行时生成于 SCREENSHOT_DIR，独立于内嵌前端资源
	screenshotDir := "data/screenshots"
	if dirVal, ok := env.Lookup(screenshot.EnvScreenshotDir); ok && !dirVal.IsEmpty() {
		screenshotDir = dirVal.String()
	}
	r.engine.Static("/screenshots", screenshotDir)

	oauthRoute := bSdkRoute.NewRoute(r.context)

	{
		apiRouter := r.engine.Group("/api/v1")

		oauthRoute.OAuthRouter(apiRouter)

		r.publicRouter(apiRouter)
		r.authRouter(apiRouter)
		r.linkRouter(apiRouter)
		r.userRouter(apiRouter)
		r.infoRouter(apiRouter)
		r.sponsorRouter(apiRouter)
		r.adminRouter(apiRouter)
	}
}
