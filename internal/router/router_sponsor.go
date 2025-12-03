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

	"github.com/gin-gonic/gin"
)

// registerSponsorAdminRouter 注册赞助管理路由（后台管理）
// 这些路由需要认证且需要管理员权限
func (r *Route) registerSponsorAdminRouter(admin *gin.RouterGroup) {
	channelHandler := handler.NewSponsorChannelHandler()
	recordHandler := handler.NewSponsorRecordHandler()

	sponsors := admin.Group("/sponsors")
	{
		// 渠道管理 /api/v1/admin/sponsors/channels
		channels := sponsors.Group("/channels")
		{
			channels.POST("", channelHandler.Add)                      // 添加赞助渠道
			channels.GET("", channelHandler.GetPage)                   // 获取赞助渠道分页列表
			channels.GET("/all", channelHandler.GetList)               // 获取所有赞助渠道（用于下拉选择）
			channels.GET("/:id", channelHandler.Get)                   // 获取赞助渠道详情
			channels.PUT("/:id", channelHandler.Update)                // 更新赞助渠道
			channels.PATCH("/:id/status", channelHandler.UpdateStatus) // 更新赞助渠道状态
			channels.DELETE("/:id", channelHandler.Delete)             // 删除赞助渠道
		}

		// 记录管理 /api/v1/admin/sponsors/records
		records := sponsors.Group("/records")
		{
			records.POST("", recordHandler.Add)          // 添加赞助记录
			records.GET("", recordHandler.GetPage)       // 获取赞助记录分页列表
			records.GET("/:id", recordHandler.Get)       // 获取赞助记录详情
			records.PUT("/:id", recordHandler.Update)    // 更新赞助记录
			records.DELETE("/:id", recordHandler.Delete) // 删除赞助记录
		}
	}
}

// registerSponsorPublicRouter 注册赞助公开路由（前台展示）
// 这些路由不需要认证，用于前台赞助墙展示
func (r *Route) registerSponsorPublicRouter() {
	channelHandler := handler.NewSponsorChannelHandler()
	recordHandler := handler.NewSponsorRecordHandler()

	sponsors := r.router.Group("/sponsors")
	{
		sponsors.GET("/channels", channelHandler.GetPublicList) // 获取启用的赞助渠道列表
		sponsors.GET("/records", recordHandler.GetPublicPage)   // 获取公开的赞助记录（分页）
	}
}
