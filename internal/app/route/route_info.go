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

	"github.com/gin-gonic/gin"
)

func (r *route) infoRouter(route gin.IRouter) {
	infoHandler := handler.NewHandler[handler.InfoHandler](r.context, "InfoHandler")
	infoGroup := route.Group("/info")
	{
		infoGroup.GET("/site", infoHandler.GetSiteInfo)
		infoGroup.GET("/about", infoHandler.GetAbout)
	}
}

func (r *route) infoAdminRouter(route gin.IRouter) {
	infoHandler := handler.NewHandler[handler.InfoHandler](r.context, "InfoHandler")
	infoGroup := route.Group("/info")
	{
		infoGroup.PUT("/site", infoHandler.UpdateSiteInfo)
		infoGroup.PUT("/about", infoHandler.UpdateAbout)
	}
}
