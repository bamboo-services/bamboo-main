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

package cmd

import (
	"bamboo-main/internal/controller/auth"
	"bamboo-main/internal/middleware/handler"
	"context"
	"github.com/XiaoLFeng/bamboo-utils/bhandler/bhook"
	"github.com/XiaoLFeng/bamboo-utils/bhandler/bmiddle"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/os/gcmd"
	"github.com/gogf/gf/v2/os/glog"
)

var (
	Main = gcmd.Command{
		Name:  "main",
		Usage: "main",
		Brief: "start http server",
		Func: func(ctx context.Context, parser *gcmd.Parser) (err error) {
			s := g.Server()

			// Hook
			if glog.GetLevel() == glog.LEVEL_DEV {
				s.BindHookHandler("/**", ghttp.HookBeforeServe, bhook.BambooHookDefaultCors)
				s.BindHookHandler("/api/v1/**", ghttp.HookBeforeServe, bhook.BambooHookRequestInfo)
			}

			// RESTful API 路由
			s.Group("/api/v1", func(group *ghttp.RouterGroup) {
				group.Middleware(bmiddle.BambooHandlerResponse)

				// 非认证
				group.Bind(
					auth.NewV1().AuthLogin,
				)

				// 认证
				group.Middleware(handler.AuthenticationHandler)
				group.Bind(
					auth.NewV1().AuthPasswordChange,
				)
			})

			s.Run()
			return nil
		},
	}
)
