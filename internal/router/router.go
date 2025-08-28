package router

import (
	"bamboo-main/internal/model/base"

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
	// TODO: 添加中间件支持

	// 初始化路由表
	route := New(getGroup)

	// 注册路由
	route.PublicRouter()
	route.AuthRouter()
	route.AdminRouter()

	// 未匹配路由处理
	engine.NoRoute(xRoute.NoRoute)
}
