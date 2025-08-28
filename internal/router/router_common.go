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
}
