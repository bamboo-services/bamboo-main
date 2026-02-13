package route

import (
	"github.com/bamboo-services/bamboo-main/internal/handler"
	"github.com/bamboo-services/bamboo-main/internal/middleware"

	"github.com/gin-gonic/gin"
)

func (r *route) adminRouter(route gin.IRouter) {
	adminGroup := route.Group("/admin")
	adminGroup.Use(middleware.AuthMiddleware)
	adminGroup.Use(middleware.RequireRole("admin"))
	{
		r.linkAdminRouter(adminGroup)
		r.linkGroupAdminRouter(adminGroup)
		r.linkColorAdminRouter(adminGroup)
		r.infoAdminRouter(adminGroup)
		r.systemUserAdminRouter(adminGroup)
		r.systemLogRouter(adminGroup)
		r.sponsorAdminRouter(adminGroup)
	}
}

func (r *route) linkAdminRouter(route gin.IRouter) {
	linkHandler := handler.NewLinkHandler()
	linkGroup := route.Group("/links")
	{
		linkGroup.POST("", linkHandler.Add)
		linkGroup.GET("", linkHandler.List)
		linkGroup.GET("/:id", linkHandler.Get)
		linkGroup.PUT("/:id", linkHandler.Update)
		linkGroup.DELETE("/:id", linkHandler.Delete)
		linkGroup.PUT("/:id/status", linkHandler.UpdateStatus)
		linkGroup.PUT("/:id/fail", linkHandler.UpdateFailStatus)
	}
}

func (r *route) linkGroupAdminRouter(route gin.IRouter) {
	groupHandler := handler.NewLinkGroupHandler()
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
	colorHandler := handler.NewLinkColorHandler()
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
