package router

import (
	"fmt"

	"bamboo-main/pkg/startup"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// Route 表示一个路由类型，用于定义和管理应用的路由逻辑
type Route struct {
	router *gin.RouterGroup
	reg    *startup.Reg
}

// New 创建并返回一个新的 Route 实例
func New(router *gin.RouterGroup, reg *startup.Reg) *Route {
	return &Route{
		router: router,
		reg:    reg,
	}
}

// Init 初始化路由配置
func Init(engine *gin.Engine, reg *startup.Reg) {
	// Swagger 文档注册（只在调试模式下开放）
	if reg.Config.Xlf.Debug {
		engine.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	}

	// 创建路由组
	getGroup := engine.Group("/api/v1")
	// TODO: 添加中间件支持

	// 初始化路由表
	route := New(getGroup, reg)

	// 注册路由
	route.PublicRouter()
	route.AuthRouter()
	route.AdminRouter()

	// 未匹配路由处理
	engine.NoRoute(route.NoRoute)
}

// NoRoute 处理未定义路由的请求
func (r *Route) NoRoute(ctx *gin.Context) {
	r.reg.Serv.Logger.Sugar().Warnf("未找到路由: %s", ctx.Request.URL.Path)
	ctx.JSON(404, gin.H{
		"code":    404,
		"message": fmt.Sprintf("页面 [%s] 不存在，请检查路由是否正确", ctx.Request.URL.Path),
		"data":    nil,
	})
}
