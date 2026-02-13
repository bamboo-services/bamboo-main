package route

import (
	"github.com/bamboo-services/bamboo-main/internal/handler"
	"github.com/bamboo-services/bamboo-main/internal/middleware"
	"github.com/gin-gonic/gin"
)

func (r *route) authRouter(route gin.IRouter) {
	authGroup := route.Group("/auth")
	authHandler := handler.NewAuthHandler()
	{
		authGroup.POST("/login", authHandler.Login)
		authGroup.POST("/register", authHandler.Register)
		authGroup.PATCH("/password/reset", authHandler.ResetPassword)
		authGroup.GET("/verify-email", authHandler.VerifyEmail)
		authGroup.GET("/reset-password", authHandler.VerifyResetToken)
		authGroup.POST("/reset-password", authHandler.ConfirmResetPassword)

		authRequiredGroup := authGroup.Group("")
		authRequiredGroup.Use(middleware.AuthMiddleware)
		{
			authRequiredGroup.PATCH("/logout", authHandler.Logout)
			authRequiredGroup.GET("/user", authHandler.GetUserInfo)
			authRequiredGroup.PUT("/password/change", authHandler.ChangePassword)
		}
	}
}
