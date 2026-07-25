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
	"github.com/bamboo-services/bamboo-main/internal/handler"
	"github.com/bamboo-services/bamboo-main/internal/middleware"
	"github.com/gin-gonic/gin"
	bSdkMiddle "github.com/phalanx-labs/beacon-sso-sdk/middleware"
)

func (r *route) authRouter(route gin.IRouter) {
	authGroup := route.Group("/auth")
	authHandler := handler.NewHandler[handler.AuthHandler](r.context, "AuthHandler")
	{
		authGroup.POST("/login", bSdkMiddle.CheckAuth(r.context), authHandler.Login)
		authGroup.POST("/register", authHandler.Register)
		authGroup.PATCH("/password/reset", authHandler.ResetPassword)
		authGroup.GET("/verify-email", authHandler.VerifyEmail)
		authGroup.GET("/reset-password", authHandler.VerifyResetToken)
		authGroup.POST("/reset-password", authHandler.ConfirmResetPassword)

		authRequiredGroup := authGroup.Group("")
		authRequiredGroup.Use(bSdkMiddle.CheckAuth(r.context))
		authRequiredGroup.Use(middleware.AuthMiddleware)
		{
			authRequiredGroup.PATCH("/logout", authHandler.Logout)
			authRequiredGroup.GET("/user", authHandler.GetUserInfo)
			authRequiredGroup.PUT("/password/change", authHandler.ChangePassword)
		}
	}
}
