package router

import (
	"bamboo-main/internal/handler"

	"github.com/gin-gonic/gin"
)

// PublicRouter 注册公开相关的路由
// 公开路由无需认证，任何人都可以访问
func (r *Route) PublicRouter() {
	public := r.router.Group("/public")
	{
		// 健康检查
		public.GET("/health", r.healthCheck)

		// ping 接口
		public.GET("/ping", r.ping)

		// 公开的友情链接接口
		linkFriendHandler := handler.NewLinkFriendHandler(r.reg)
		public.GET("/links", linkFriendHandler.GetPublicLinks)
		public.GET("/links/group/:group_uuid", linkFriendHandler.GetPublicLinks) // 按分组获取
	}
}

// healthCheck 健康检查接口
func (r *Route) healthCheck(c *gin.Context) {
	// TODO: 实现完整的健康检查逻辑，包括数据库连接、Redis连接等
	c.JSON(200, gin.H{
		"code":    200,
		"message": "服务正常",
		"data": gin.H{
			"service": "bamboo-main",
			"version": "v1.0.0",
			"status":  "healthy",
		},
	})
}

// ping 简单的连通性测试接口
func (r *Route) ping(c *gin.Context) {
	c.JSON(200, gin.H{
		"code":    200,
		"message": "pong",
		"data":    nil,
	})
}