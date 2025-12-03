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
	"bamboo-main/internal/middleware"

	"github.com/gin-gonic/gin"
)

// AdminRouter 注册管理员相关的路由
// 管理员路由需要认证且需要管理员权限
func (r *Route) AdminRouter() {
	admin := r.router.Group("/admin")
	admin.Use(middleware.AuthMiddleware)
	admin.Use(middleware.RequireRole("admin"))
	{
		r.registerLinkAdminRouter(admin)       // 友情链接管理路由
		r.registerLinkGroupAdminRouter(admin)  // 友链分组管理路由
		r.registerLinkColorAdminRouter(admin)  // 友链颜色管理路由
		r.registerInfoAdminRouter(admin)       // 站点信息管理路由
		r.registerSystemUserAdminRouter(admin) // 系统用户管理路由
		r.registerSystemLogRouter(admin)       // 系统日志路由
		r.registerSponsorAdminRouter(admin)    // 赞助管理路由
	}
}

// registerLinkAdminRouter 注册友情链接管理路由
func (r *Route) registerLinkAdminRouter(admin *gin.RouterGroup) {
	linkHandler := handler.NewLinkHandler()
	links := admin.Group("/links")
	{
		links.POST("", linkHandler.Add)                             // 添加友情链接
		links.GET("", linkHandler.List)                             // 获取友情链接分页列表
		links.GET("/:link_uuid", linkHandler.Get)                   // 获取友情链接详情
		links.PUT("/:link_uuid", linkHandler.Update)                // 更新友情链接
		links.DELETE("/:link_uuid", linkHandler.Delete)             // 删除友情链接
		links.PUT("/:link_uuid/status", linkHandler.UpdateStatus)   // 更新友情链接状态
		links.PUT("/:link_uuid/fail", linkHandler.UpdateFailStatus) // 更新友情链接失效状态
	}
}

// registerLinkGroupAdminRouter 注册友链分组管理路由
func (r *Route) registerLinkGroupAdminRouter(admin *gin.RouterGroup) {
	groupHandler := handler.NewLinkGroupHandler()
	groups := admin.Group("/groups")
	{
		groups.POST("", groupHandler.Add)                              // 添加友链分组
		groups.GET("", groupHandler.GetPage)                           // 获取友链分组分页列表
		groups.GET("/all", groupHandler.GetList)                       // 获取所有友链分组（用于下拉选择）
		groups.PATCH("/sort", groupHandler.UpdateSort)                 // 批量更新分组排序
		groups.GET("/:group_uuid", groupHandler.Get)                   // 获取友链分组详情
		groups.PUT("/:group_uuid", groupHandler.Update)                // 更新友链分组
		groups.PATCH("/:group_uuid/status", groupHandler.UpdateStatus) // 更新友链分组状态
		groups.DELETE("/:group_uuid", groupHandler.Delete)             // 删除友链分组
	}
}

// registerLinkColorAdminRouter 注册友链颜色管理路由
func (r *Route) registerLinkColorAdminRouter(admin *gin.RouterGroup) {
	colorHandler := handler.NewLinkColorHandler()
	colors := admin.Group("/colors")
	{
		colors.POST("", colorHandler.Add)                      // 添加友链颜色
		colors.GET("", colorHandler.GetPage)                   // 获取友链颜色分页列表
		colors.GET("/all", colorHandler.GetList)               // 获取所有友链颜色（用于下拉选择）
		colors.PATCH("/sort", colorHandler.UpdateSort)         // 批量更新颜色排序
		colors.GET("/:id", colorHandler.Get)                   // 获取友链颜色详情
		colors.PUT("/:id", colorHandler.Update)                // 更新友链颜色
		colors.PATCH("/:id/status", colorHandler.UpdateStatus) // 更新友链颜色状态
		colors.DELETE("/:id", colorHandler.Delete)             // 删除友链颜色
	}
}

// registerSystemUserAdminRouter 注册系统用户管理路由
func (r *Route) registerSystemUserAdminRouter(admin *gin.RouterGroup) {
	// TODO: 创建 SystemUserHandler 或扩展 AuthHandler
	users := admin.Group("/users")
	{
		// 这些路由需要等handler实现后再添加
		users.POST("", nil)                           // 添加系统用户
		users.GET("", nil)                            // 获取系统用户分页列表
		users.GET("/:user_uuid", nil)                 // 获取系统用户详情
		users.PUT("/:user_uuid", nil)                 // 更新系统用户
		users.DELETE("/:user_uuid", nil)              // 删除系统用户
		users.PUT("/:user_uuid/status", nil)          // 更新用户状态
		users.POST("/:user_uuid/password/reset", nil) // 重置用户密码
	}
}

// registerSystemLogRouter 注册系统日志路由
func (r *Route) registerSystemLogRouter(admin *gin.RouterGroup) {
	// TODO: 创建 SystemLogHandler
	logs := admin.Group("/logs")
	{
		// 这些路由需要等handler实现后再添加
		logs.GET("", nil)              // 获取系统日志分页列表
		logs.GET("/:log_uuid", nil)    // 获取系统日志详情
		logs.DELETE("/:log_uuid", nil) // 删除系统日志
		logs.POST("/cleanup", nil)     // 清理旧日志
	}
}
