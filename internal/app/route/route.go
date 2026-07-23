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

func NewRoute(ctx context.Context, serve *gin.Engine) {
	r := &route{
		engine:  serve,
		context: ctx,
	}

	r.engine.NoMethod(xRoute.NoMethod)
	r.engine.NoRoute(xRoute.NoRoute)

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
