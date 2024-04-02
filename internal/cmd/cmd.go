package cmd

import (
	"context"
	"develop/internal/controller/auth"
	"develop/internal/middleware"
	"develop/manifest/boot"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/os/gcmd"
)

var (
	Main = gcmd.Command{
		Name:  "main",
		Usage: "main",
		Brief: "start http server",
		Func: func(ctx context.Context, parser *gcmd.Parser) (err error) {
			// 数据进行初始化
			boot.InitialDatabase(ctx)
			// 服务器启动
			s := g.Server()
			s.Group("/", func(group *ghttp.RouterGroup) {
				group.Middleware(middleware.MiddleErrorHandler)
				//group.Middleware(ghttp.MiddlewareHandlerResponse)
				group.Bind(
					auth.NewV1(),
				)
			})
			s.Run()
			return nil
		},
	}
)
