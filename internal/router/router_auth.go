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
		auth.POST("/login", authHandler.Login)                  // 用户登录
		auth.POST("/password/reset", authHandler.ResetPassword) // 重置密码

		// 需要认证的路由
		authRequired := auth.Group("")
		authRequired.Use(middleware.AuthMiddleware)
		{
			authRequired.POST("/logout", authHandler.Logout)                  // 用户登出
			authRequired.GET("/user", authHandler.GetUserInfo)                // 获取用户信息
			authRequired.POST("/password/change", authHandler.ChangePassword) // 修改密码
			// TODO: 实现更新个人资料和头像上传功能
			// authRequired.PUT("/profile", authHandler.UpdateProfile)
			// authRequired.POST("/avatar/upload", authHandler.UploadAvatar)
		}
	}
}
