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
	"bamboo-main/internal/model/base"

	xMiddle "github.com/bamboo-services/bamboo-base-go/middleware"
	xRoute "github.com/bamboo-services/bamboo-base-go/route"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// Route 表示一个路由类型，用于定义和管理应用的路由逻辑
type Route struct {
	router *gin.RouterGroup
}

// New 创建并返回一个新的 Route 实例
func New(router *gin.RouterGroup) *Route {
	return &Route{
		router: router,
	}
}

// Init 初始化路由配置
func Init(engine *gin.Engine, config *base.BambooConfig) {
	// Swagger 文档注册（只在调试模式下开放）
	if config.Xlf.Debug {
		engine.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	}

	// 创建路由组
	getGroup := engine.Group("/api/v1")
	getGroup.Use(xMiddle.ResponseMiddleware)

	// 初始化路由表
	route := New(getGroup)

	// 注册路由
	route.PublicRouter() // 纯公开接口 (health, ping)
	route.CommonRouter() // 业务接口但无需认证 (links)
	route.AuthRouter()   // 认证相关
	route.AdminRouter()  // 管理员专用

	// 未匹配路由处理
	engine.NoRoute(xRoute.NoRoute)
}
