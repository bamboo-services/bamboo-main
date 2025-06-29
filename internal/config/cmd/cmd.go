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
	"bamboo-main/internal/controller/link_color"
	"bamboo-main/internal/controller/link_friend"
	"bamboo-main/internal/controller/link_group"
	"bamboo-main/internal/handler/hook"
	"bamboo-main/internal/handler/middleware"
	"context"
	"github.com/XiaoLFeng/bamboo-utils/bhandler/bhook"
	"github.com/XiaoLFeng/bamboo-utils/bhandler/bmiddle"
	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/os/gcmd"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/gogf/gf/v2/os/glog"
	"github.com/gogf/gf/v2/os/gtime"
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
				s.BindHookHandler("/api/*", ghttp.HookBeforeServe, bhook.BambooHookDefaultCors)
				s.BindHookHandler("/api/v1/*", ghttp.HookBeforeServe, bhook.BambooHookRequestInfo)
			}
			s.BindHookHandler("/api/v1/*", ghttp.HookAfterOutput, hook.LogHook)

			// RESTful API 路由
			s.Group("/api/v1", func(group *ghttp.RouterGroup) {
				group.Middleware(bmiddle.BambooHandlerResponse)

				// 非认证
				group.Bind(
					// 授权
					auth.NewV1().AuthLogin,
					auth.NewV1().AuthPasswordReset,
					// 友链颜色
					link_color.NewV1().LinkColorGet,
					// 友链分组
					link_group.NewV1().LinkGroupGet,
				)

				// 认证
				group.Middleware(middleware.AuthenticationMiddleware)
				group.Bind(
					// 授权
					auth.NewV1().AuthPasswordChange,
					// 友链分组
					link_group.NewV1().LinkGroupAdd,
					link_group.NewV1().LinkGroupEdit,
					link_group.NewV1().LinkGroupDelete,
					link_group.NewV1().LinkGroupPage,
					// 友链颜色
					link_color.NewV1().LinkColorAdd,
					link_color.NewV1().LinkColorEdit,
					link_color.NewV1().LinkColorDelete,
					link_color.NewV1().LinkColorPage,
					// 友链
					link_friend.NewV1().LinkFriendAdd,
					link_friend.NewV1().LinkFriendEdit,
					link_friend.NewV1().LinkFriendStatus,
					link_friend.NewV1().LinkFriendFail,
					link_friend.NewV1().LinkFriendDelete,
				)
			})

			// 404 Not Found
			s.BindHandler("/api/v1/*", func(r *ghttp.Request) {
				getValue := g.Map{
					"context": gctx.CtxId(r.GetCtx()),
					"code":    40401,
					"message": "页面不存在",
					"time":    gtime.Now().TimestampMilli(),
				}
				r.Response.Header().Set("Content-Type", "application/json; charset=utf-8")
				if glog.GetLevel() == glog.LEVEL_DEV {
					getValue["overhead"] = gtime.Now().Sub(r.EnterTime).Milliseconds()
				}
				jsonEncodeData := gjson.MustEncodeString(getValue)
				r.Response.WriteStatus(404, jsonEncodeData)
			})

			s.Run()
			return nil
		},
	}
)
