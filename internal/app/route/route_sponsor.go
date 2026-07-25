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

	"github.com/gin-gonic/gin"
)

func (r *route) sponsorAdminRouter(route gin.IRouter) {
	channelHandler := handler.NewHandler[handler.SponsorChannelHandler](r.context)
	recordHandler := handler.NewHandler[handler.SponsorRecordHandler](r.context)

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
	channelHandler := handler.NewHandler[handler.SponsorChannelHandler](r.context)
	recordHandler := handler.NewHandler[handler.SponsorRecordHandler](r.context)

	sponsorGroup := route.Group("/sponsors")
	{
		sponsorGroup.GET("/channels", channelHandler.GetPublicList)
		sponsorGroup.GET("/records", recordHandler.GetPublicPage)
	}
}
