package startup

import (
	xConsts "github.com/bamboo-services/bamboo-base-go/constants"
	"github.com/gin-gonic/gin"
)

type handler struct {
	*Reg
}

// SystemContextInit 初始化系统上下文。
//
// 此方法用于设置全局系统上下文，包括数据库连接、Redis连接、配置信息等的全局访问。
// 这些上下文信息将在整个应用程序生命周期中可用，供各个组件使用。
//
// 功能包括:
//   - 设置全局数据库连接上下文
//   - 设置全局Redis连接上下文
//   - 设置全局配置信息上下文
//   - 初始化其他系统级别的上下文信息
//
// 注意: 此方法应在所有其他初始化方法之后调用，确保所有必要的组件都已正确初始化。
func (r *Reg) SystemContextInit() {
	r.Serv.Logger.Named(xConsts.LogINIT).Info("初始化系统上下文")

	// 创建处理器实例
	handler := &handler{r}

	// 注册系统上下文处理函数
	r.Serv.Serve.Use(handler.systemContextHandlerFunc)
}

// systemContextHandlerFunc 系统上下文中间件处理函数。
//
// 为当前请求的 Gin 上下文设置必要的信息。
// 包括数据库实例和 Redis 客户端，这些信息通过 `Context` 存储，
// 供整个请求生命周期中的其他处理函数使用。
//
// 参数说明:
//   - c: Gin 上下文对象，代表当前请求的上下文。
//
// 此方法设置完上下文后会调用 `c.Next()` 放行请求，以执行后续中间件或路由处理函数。
func (h *handler) systemContextHandlerFunc(c *gin.Context) {
	// 设置 Context
	c.Set(xConsts.ContextDatabase.String(), h.DB)
	c.Set(xConsts.ContextRedisClient.String(), h.Rdb)
	c.Set(xConsts.ContextCustomConfig.String(), h.Config)

	// 放行内容
	c.Next()
}
