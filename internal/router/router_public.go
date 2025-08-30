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
