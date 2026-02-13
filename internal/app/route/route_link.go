package route

import (
	"github.com/bamboo-services/bamboo-main/internal/handler"
	"github.com/gin-gonic/gin"
)

func (r *route) linkRouter(route gin.IRouter) {
	linkHandler := handler.NewHandler[handler.LinkHandler](r.context, "LinkHandler")
	linkGroup := route.Group("/links")
	{
		linkGroup.GET("", linkHandler.GetPublicLinks)
	}
}
