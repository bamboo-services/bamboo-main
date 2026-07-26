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

	xEnv "github.com/bamboo-services/bamboo-base-go/defined/env"
	xMiddle "github.com/bamboo-services/bamboo-base-go/major/middleware"
	xRoute "github.com/bamboo-services/bamboo-base-go/major/route"
	"github.com/gin-gonic/gin"
	bSdkRoute "github.com/phalanx-labs/beacon-sso-sdk/route"
)

type route struct {
	engine  *gin.Engine
	context context.Context
}

// NewRoute 路由注册器，装配中间件链与各领域子路由
func NewRoute(ctx context.Context, serve *gin.Engine) {
	r := &route{
		engine:  serve,
		context: ctx,
	}

	r.engine.NoMethod(xRoute.NoMethod)
	r.engine.NoRoute(noRouteHandler)

	r.engine.Use(xMiddle.ResponseMiddleware)
	r.engine.Use(xMiddle.ReleaseAllCors)
	r.engine.Use(xMiddle.AllowOptionRequest)

	if xEnv.GetEnvBool(xEnv.Debug, false) {
		swaggerRegister(r.engine)
	}

	oauthRoute := bSdkRoute.NewRoute(r.context)

	{
		apiRouter := r.engine.Group("/api/v1")

		oauthRoute.OAuthRouter(apiRouter)

		r.publicRouter(apiRouter)
		r.authRouter(apiRouter)
		r.linkRouter(apiRouter)
		r.infoRouter(apiRouter)
		r.sponsorRouter(apiRouter)
		r.adminRouter(apiRouter)
	}
}
