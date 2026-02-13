package route

import (
	"github.com/bamboo-services/bamboo-main/internal/handler"

	"github.com/gin-gonic/gin"
)

func (r *route) infoRouter(route gin.IRouter) {
	infoHandler := handler.NewInfoHandler()
	infoGroup := route.Group("/info")
	{
		infoGroup.GET("/site", infoHandler.GetSiteInfo)
		infoGroup.GET("/about", infoHandler.GetAbout)
	}
}

func (r *route) infoAdminRouter(route gin.IRouter) {
	infoHandler := handler.NewInfoHandler()
	infoGroup := route.Group("/info")
	{
		infoGroup.PUT("/site", infoHandler.UpdateSiteInfo)
		infoGroup.PUT("/about", infoHandler.UpdateAbout)
	}
}
