/*
 * --------------------------------------------------------------------------------
 * Copyright (c) 2016-NOW(至今) 筱锋
 * Author: 筱锋(https://www.x-lf.com)
 *
 * 本文件包含 XiaoMain 的源代码，该项目的所有源代码均遵循MIT开源许可证协议。
 * --------------------------------------------------------------------------------
 * 许可证声明：
 *
 * 版权所有 (c) 2016-2024 筱锋。保留所有权利。
 *
 * 本软件是“按原样”提供的，没有任何形式的明示或暗示的保证，包括但不限于
 * 对适销性、特定用途的适用性和非侵权性的暗示保证。在任何情况下，
 * 作者或版权持有人均不承担因软件或软件的使用或其他交易而产生的、
 * 由此引起的或以任何方式与此软件有关的任何索赔、损害或其他责任。
 *
 * 使用本软件即表示您了解此声明并同意其条款。
 *
 * 有关MIT许可证的更多信息，请查看项目根目录下的LICENSE文件或访问：
 * https://opensource.org/licenses/MIT
 * --------------------------------------------------------------------------------
 * 免责声明：
 *
 * 使用本软件的风险由用户自担。作者或版权持有人在法律允许的最大范围内，
 * 对因使用本软件内容而导致的任何直接或间接的损失不承担任何责任。
 * --------------------------------------------------------------------------------
 */

package cmd

import (
	"context"
	"github.com/bamboo-services/bamboo-utils/bmiddle"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/os/gcmd"
	"xiaoMain/internal/config/middleware"
	"xiaoMain/internal/config/startup"
	"xiaoMain/internal/config/task"
	"xiaoMain/internal/controller/auth"
	"xiaoMain/internal/controller/info"
	"xiaoMain/internal/controller/link"
	"xiaoMain/internal/controller/mail"
	"xiaoMain/internal/controller/rss"
)

var (
	Main = gcmd.Command{
		Name:  "main",
		Usage: "main",
		Brief: "start http server",
		Func: func(ctx context.Context, parser *gcmd.Parser) (err error) {
			// 服务器启动
			s := g.Server()

			// 数据进行初始化
			startup.Initial(ctx)
			// 定时任务
			task.Task(ctx)

			// 关闭路由映射输出
			s.SetDumpRouterMap(false)

			// 路由注册
			s.Group("/", func(group *ghttp.RouterGroup) {
				group.Middleware(middleware.MiddleTimeHandler) // 接口时间统计接口

				// 前端部分
				group.Group("", func(group *ghttp.RouterGroup) {

				})

				// 后端部分
				group.Group("/api/v1", func(group *ghttp.RouterGroup) {
					group.Middleware(middleware.MiddleAccessUserHandler) // 访问处理中间件
					group.Middleware(bmiddle.BambooMiddleHandler)        // 错误集中处理中间件

					// 用户操作
					group.Bind(
						auth.NewV1().AuthLogin,
						auth.NewV1().AuthResetPassword,
					)

					// 链接操作
					group.Group("/link", func(group *ghttp.RouterGroup) {
						group.Bind(
							link.NewV1().LinkAdd,
							link.NewV1().CheckLinkURLHasConnect,
							link.NewV1().CheckLogoURLHasConnect,
							link.NewV1().CheckRssURLHasConnect,
							link.NewV1().LinkGetColor,
							link.NewV1().LinkGetLocation,
						)
						group.Middleware(middleware.MiddleAuthHandler).Bind(
							link.NewV1().LinkColorAdd,
							link.NewV1().LinkLocationAdd,
							link.NewV1().LinkGetColorFull,
							link.NewV1().LinkGetLocationFull,
							link.NewV1().LinkAddAdmin,
							link.NewV1().CheckLinkIDHasConnect,
						)
					})

					// 系统信息
					group.Group("/info", func(group *ghttp.RouterGroup) {
						group.Bind(
							info.NewV1().GetWebInfo,
							info.NewV1().EditWebInfo,
						)
					})

					// 邮件操作
					group.Group("/mail", func(group *ghttp.RouterGroup) {
						group.Bind(
							mail.NewV1().MailSend,
						)
					})

					// Rss订阅消息
					group.Group("/rss", func(group *ghttp.RouterGroup) {
						group.Bind(
							rss.NewV1().GetLinkRssInfo,
						)
					})

					group.Middleware(middleware.MiddleAuthHandler).Bind(
						auth.NewV1().AuthChangePassword,
					)
				})
			})
			s.Run()
			return nil
		},
	}
)
