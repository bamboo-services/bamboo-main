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
