package route

import (
	"github.com/bamboo-services/bamboo-main/internal/handler"
	"github.com/gin-gonic/gin"
)

func (r *route) publicRouter(route gin.IRouter) {
	publicHandler := handler.NewHandler[handler.PublicHandler](r.context, "PublicHandler")
	publicGroup := route.Group("/public")
	{
		publicGroup.GET("/health", publicHandler.HealthCheck)
		publicGroup.GET("/ping", publicHandler.Ping)
	}
}
