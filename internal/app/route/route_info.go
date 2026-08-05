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
)

// infoRouter 站点信息路由：GET 公开、PUT 统一挂 /info/admin 子组并鉴权管理员。
func (r *route) infoRouter(route gin.IRouter) {
	infoHandler := handler.NewHandler[handler.InfoHandler](r.context)
	infoGroup := route.Group("/info")
	{
		infoGroup.GET("/site", infoHandler.GetSiteInfo)
		infoGroup.GET("/archive", infoHandler.GetArchiveInfo)
		infoGroup.GET("/apply-site", infoHandler.GetApplySiteInfo)
		infoGroup.GET("/blogger", infoHandler.GetBloggerInfo)
		infoGroup.GET("/builtin-invalid-group", infoHandler.GetBuiltinInvalidGroup)
		infoGroup.GET("/color-mode", infoHandler.GetColorMode)

		// 管理端写操作：统一 /info/admin 前置路径，鉴权管理员
		adminGroup := infoGroup.Group("/admin")
		adminGroup.Use(middleware.AuthMiddleware)
		adminGroup.Use(middleware.RequireRole("admin"))
		{
			adminGroup.PUT("/site", infoHandler.UpdateSiteInfo)
			adminGroup.PUT("/archive", infoHandler.UpdateArchiveInfo)
			adminGroup.PUT("/apply-site", infoHandler.UpdateApplySiteInfo)
			adminGroup.PUT("/blogger", infoHandler.UpdateBloggerInfo)
			adminGroup.PUT("/builtin-invalid-group", infoHandler.UpdateBuiltinInvalidGroup)
			adminGroup.PUT("/color-mode", infoHandler.UpdateColorMode)
		}
	}
}
