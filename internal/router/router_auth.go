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

package router

import (
	"bamboo-main/internal/handler"
	"bamboo-main/internal/middleware"
)

// AuthRouter 注册认证相关的路由
// 包含登录、登出、密码管理等用户认证相关功能
func (r *Route) AuthRouter() {
	auth := r.router.Group("/auth")
	authHandler := handler.NewAuthHandler()
	{
		// 无需认证的路由
		auth.POST("/login", authHandler.Login)                   // 用户登录
		auth.PATCH("/password/reset", authHandler.ResetPassword) // 重置密码

		// 需要认证的路由
		authRequired := auth.Group("")
		authRequired.Use(middleware.AuthMiddleware)
		{
			authRequired.PATCH("/logout", authHandler.Logout)                // 用户登出
			authRequired.GET("/user", authHandler.GetUserInfo)               // 获取用户信息
			authRequired.PUT("/password/change", authHandler.ChangePassword) // 修改密码
			// TODO: 实现更新个人资料和头像上传功能
			// authRequired.PUT("/profile", authHandler.UpdateProfile)
			// authRequired.POST("/avatar/upload", authHandler.UploadAvatar)
		}
	}
}
