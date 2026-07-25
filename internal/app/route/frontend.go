// --------------------------------------------------------------------------------
// Copyright (c) 2016-NOW(至今) 筱锋
// Author: 筱锋「xiao_lfeng」(https://www.x-lf.com)
// --------------------------------------------------------------------------------
// 许可证声明：版权所有 (c) 2016-2026 筱锋。保留所有权利。
// 有关MIT许可证的更多信息，请查看项目根目录下的LICENSE文件或访问：
// https://opensource.org/licenses/MIT
// --------------------------------------------------------------------------------

package route

import (
	"net/http"
	"strings"

	xRoute "github.com/bamboo-services/bamboo-base-go/major/route"
	"github.com/bamboo-services/bamboo-main/resources/frontend"
	"github.com/gin-gonic/gin"
)

// frontendFileServer 基于内嵌前端产物构建的静态文件服务器，用于 NoRoute 阶段托管。
var frontendFileServer = http.FileServer(http.FS(frontend.FS()))

// noRouteHandler 统一接管未匹配路由的请求，按路径前缀分流：
//   - /api/* 与 /swagger/*：交由 xRoute.NoRoute 返回标准化 JSON 404；
//   - 其余路径：优先匹配内嵌静态资源，命中则直接返回；
//   - 未命中静态资源时回退到 index.html，以承载前端 SPA 路由（history mode）。
//
// 设计目的：单二进制即可同时提供后端 API 与前端页面，无需独立部署前端服务。
func noRouteHandler(ctx *gin.Context) {
	path := ctx.Request.URL.Path

	// 接口与文档类路径仍走 JSON 404，避免污染 API 错误语义
	if strings.HasPrefix(path, "/api/") || strings.HasPrefix(path, "/swagger/") {
		xRoute.NoRoute(ctx)
		return
	}

	// 根路径直接由 FileServer 返回 index.html
	if path == "/" {
		frontendFileServer.ServeHTTP(ctx.Writer, ctx.Request)
		return
	}

	// 命中静态资源则原样返回（含 hashed assets、favicon 等）
	if file, err := frontend.FS().Open(strings.TrimPrefix(path, "/")); err == nil {
		_ = file.Close()
		frontendFileServer.ServeHTTP(ctx.Writer, ctx.Request)
		return
	}

	// SPA fallback：重写路径到根，由 FileServer 返回 index.html
	ctx.Request.URL.Path = "/"
	frontendFileServer.ServeHTTP(ctx.Writer, ctx.Request)
}
