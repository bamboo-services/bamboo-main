package router

import (
	"bamboo-main/internal/handler"
)

// PublicRouter 注册公开相关的路由
// 公开路由无需认证，完全对外开放，无业务逻辑
func (r *Route) PublicRouter() {
	publicHandler := handler.NewPublicHandler()
	public := r.router.Group("/public")
	{
		// 健康检查
		public.GET("/health", publicHandler.HealthCheck)

		// ping 接口
		public.GET("/ping", publicHandler.Ping)
	}
}

