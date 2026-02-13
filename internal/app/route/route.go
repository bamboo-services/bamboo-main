package route

import (
	"context"

	xEnv "github.com/bamboo-services/bamboo-base-go/env"
	xMiddle "github.com/bamboo-services/bamboo-base-go/middleware"
	xReg "github.com/bamboo-services/bamboo-base-go/register"
	xRoute "github.com/bamboo-services/bamboo-base-go/route"
	"github.com/gin-gonic/gin"
)

type route struct {
	engine  *gin.Engine
	context context.Context
}

func NewRoute(reg *xReg.Reg) {
	r := &route{
		engine:  reg.Serve,
		context: reg.Init.Ctx,
	}

	r.engine.NoMethod(xRoute.NoMethod)
	r.engine.NoRoute(xRoute.NoRoute)

	r.engine.Use(xMiddle.ResponseMiddleware)
	r.engine.Use(xMiddle.ReleaseAllCors)
	r.engine.Use(xMiddle.AllowOption)

	if xEnv.GetEnvBool(xEnv.Debug, false) {
		swaggerRegister(r.engine)
	}

	{
		apiRouter := r.engine.Group("/api/v1")

		r.publicRouter(apiRouter)
		r.authRouter(apiRouter)
		r.linkRouter(apiRouter)
		r.infoRouter(apiRouter)
		r.sponsorRouter(apiRouter)
		r.adminRouter(apiRouter)
	}
}
