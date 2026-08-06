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

func (r *route) adminRouter(route gin.IRouter) {
	adminGroup := route.Group("/admin")
	adminGroup.Use(middleware.AuthMiddleware)
	adminGroup.Use(middleware.RequireAdmin)
	{
		r.linkAdminRouter(adminGroup)
		r.linkGroupAdminRouter(adminGroup)
		r.linkColorAdminRouter(adminGroup)
		r.dashboardAdminRouter(adminGroup)
		r.systemUserAdminRouter(adminGroup)
		r.systemLogRouter(adminGroup)
		r.sponsorAdminRouter(adminGroup)
	}
}

func (r *route) dashboardAdminRouter(route gin.IRouter) {
	dashboardHandler := handler.NewHandler[handler.DashboardHandler](r.context)
	dashboardRouter := route.Group("/dashboard")
	{
		dashboardRouter.GET("/stats", dashboardHandler.Stats)
	}
}

func (r *route) linkAdminRouter(route gin.IRouter) {
	linkHandler := handler.NewHandler[handler.LinkHandler](r.context)
	linkGroup := route.Group("/links")
	{
		linkGroup.POST("", linkHandler.Add)
		linkGroup.GET("", linkHandler.List)
		linkGroup.PATCH("/sort", linkHandler.UpdateSort)
		linkGroup.GET("/:id", linkHandler.Get)
		linkGroup.PUT("/:id", linkHandler.Update)
		linkGroup.DELETE("/:id", linkHandler.Delete)
		linkGroup.PUT("/:id/status", linkHandler.UpdateStatus)
		linkGroup.PUT("/:id/fail", linkHandler.UpdateFailStatus)
		linkGroup.POST("/:id/screenshot", linkHandler.ReScreenshot)
		linkGroup.POST("/:id/edit-request/approve", linkHandler.ApproveEditRequest)
		linkGroup.POST("/:id/edit-request/reject", linkHandler.RejectEditRequest)
	}
}

func (r *route) linkGroupAdminRouter(route gin.IRouter) {
	groupHandler := handler.NewHandler[handler.LinkGroupHandler](r.context)
	groupRouter := route.Group("/groups")
	{
		groupRouter.POST("", groupHandler.Add)
		groupRouter.GET("", groupHandler.GetPage)
		groupRouter.GET("/all", groupHandler.GetList)
		groupRouter.PATCH("/sort", groupHandler.UpdateSort)
		groupRouter.GET("/:id", groupHandler.Get)
		groupRouter.PUT("/:id", groupHandler.Update)
		groupRouter.PATCH("/:id/status", groupHandler.UpdateStatus)
		groupRouter.DELETE("/:id", groupHandler.Delete)
	}
}

func (r *route) linkColorAdminRouter(route gin.IRouter) {
	colorHandler := handler.NewHandler[handler.LinkColorHandler](r.context)
	colorRouter := route.Group("/colors")
	{
		colorRouter.POST("", colorHandler.Add)
		colorRouter.GET("", colorHandler.GetPage)
		colorRouter.GET("/all", colorHandler.GetList)
		colorRouter.PATCH("/sort", colorHandler.UpdateSort)
		colorRouter.GET("/:id", colorHandler.Get)
		colorRouter.PUT("/:id", colorHandler.Update)
		colorRouter.PATCH("/:id/status", colorHandler.UpdateStatus)
		colorRouter.DELETE("/:id", colorHandler.Delete)
	}
}

func (r *route) systemUserAdminRouter(route gin.IRouter) {
	_ = route
}

func (r *route) systemLogRouter(route gin.IRouter) {
	_ = route
}
