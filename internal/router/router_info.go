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

// registerInfoPublicRouter 注册公开的站点信息路由
// 无需认证，用于前台展示
func (r *Route) registerInfoPublicRouter() {
	infoHandler := handler.NewInfoHandler()
	info := r.router.Group("/info")
	{
		info.GET("/site", infoHandler.GetSiteInfo) // 获取站点信息
		info.GET("/about", infoHandler.GetAbout)   // 获取自我介绍
	}
}

// registerInfoAdminRouter 注册管理员站点信息路由
// 需要管理员权限
func (r *Route) registerInfoAdminRouter(admin *gin.RouterGroup) {
	infoHandler := handler.NewInfoHandler()
	info := admin.Group("/info")
	{
		info.PUT("/site", infoHandler.UpdateSiteInfo) // 更新站点信息
		info.PUT("/about", infoHandler.UpdateAbout)   // 更新自我介绍
	}
}
