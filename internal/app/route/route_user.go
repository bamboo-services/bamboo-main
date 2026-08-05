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

// userRouter 注册用户自助路由组
//
// 面向已登录的普通用户：管理自己名下的友链（查看/编辑/申请下架）与个人资料。
// 全部端点经 AuthMiddleware 鉴权，友链归属校验由各 logic 方法内部完成。
func (r *route) userRouter(route gin.IRouter) {
	linkHandler := handler.NewHandler[handler.LinkHandler](r.context)
	authHandler := handler.NewHandler[handler.AuthHandler](r.context)
	sponsorRecordHandler := handler.NewHandler[handler.SponsorRecordHandler](r.context)

	userGroup := route.Group("/user")
	userGroup.Use(middleware.AuthMiddleware)
	{
		userGroup.GET("/links", linkHandler.ListMyLinks)
		userGroup.GET("/links/:id", linkHandler.GetMyLink)
		userGroup.PUT("/links/:id", linkHandler.UpdateMyLink)
		userGroup.PUT("/links/:id/takedown", linkHandler.RequestTakedown)
		userGroup.GET("/sponsors", sponsorRecordHandler.ListMyRecords)
		userGroup.GET("/sponsors/:id", sponsorRecordHandler.GetMyRecord)
		userGroup.PUT("/sponsors/:id", sponsorRecordHandler.UpdateMyRecord)
		userGroup.PUT("/profile", authHandler.UpdateProfile)
	}
}
