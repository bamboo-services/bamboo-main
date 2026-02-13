package route

import (
	"github.com/bamboo-services/bamboo-main/internal/handler"

	"github.com/gin-gonic/gin"
)

func (r *route) sponsorAdminRouter(route gin.IRouter) {
	channelHandler := handler.NewHandler[handler.SponsorChannelHandler](r.context, "SponsorChannelHandler")
	recordHandler := handler.NewHandler[handler.SponsorRecordHandler](r.context, "SponsorRecordHandler")

	sponsorGroup := route.Group("/sponsors")
	{
		channelGroup := sponsorGroup.Group("/channels")
		{
			channelGroup.POST("", channelHandler.Add)
			channelGroup.GET("", channelHandler.GetPage)
			channelGroup.GET("/all", channelHandler.GetList)
			channelGroup.GET("/:id", channelHandler.Get)
			channelGroup.PUT("/:id", channelHandler.Update)
			channelGroup.PATCH("/:id/status", channelHandler.UpdateStatus)
			channelGroup.DELETE("/:id", channelHandler.Delete)
		}

		recordGroup := sponsorGroup.Group("/records")
		{
			recordGroup.POST("", recordHandler.Add)
			recordGroup.GET("", recordHandler.GetPage)
			recordGroup.GET("/:id", recordHandler.Get)
			recordGroup.PUT("/:id", recordHandler.Update)
			recordGroup.DELETE("/:id", recordHandler.Delete)
		}
	}
}

func (r *route) sponsorRouter(route gin.IRouter) {
	channelHandler := handler.NewHandler[handler.SponsorChannelHandler](r.context, "SponsorChannelHandler")
	recordHandler := handler.NewHandler[handler.SponsorRecordHandler](r.context, "SponsorRecordHandler")

	sponsorGroup := route.Group("/sponsors")
	{
		sponsorGroup.GET("/channels", channelHandler.GetPublicList)
		sponsorGroup.GET("/records", recordHandler.GetPublicPage)
	}
}
