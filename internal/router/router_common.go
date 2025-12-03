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
)

// CommonRouter 注册通用相关的路由
// 通用路由无需认证，但包含业务功能的接口
func (r *Route) CommonRouter() {
	// 友情链接接口 (无需认证，但有业务功能)
	linkHandler := handler.NewLinkHandler()

	r.router.GET("/links", linkHandler.GetPublicLinks)
	r.router.GET("/links/group/:group_uuid", linkHandler.GetPublicLinks) // 按分组获取

	// 站点信息接口 (无需认证)
	r.registerInfoPublicRouter()

	// 赞助公开接口 (无需认证，用于前台赞助墙展示)
	r.registerSponsorPublicRouter()
}
