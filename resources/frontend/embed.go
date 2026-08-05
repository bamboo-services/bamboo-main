// --------------------------------------------------------------------------------
// Copyright (c) 2016-NOW(至今) 筱锋
// Author: 筱锋「xiao_lfeng」(https://www.x-lf.com)
// --------------------------------------------------------------------------------
// 许可证声明：版权所有 (c) 2016-2026 筱锋。保留所有权利。
// 有关MIT许可证的更多信息，请查看项目根目录下的LICENSE文件或访问：
// https://opensource.org/licenses/MIT
// --------------------------------------------------------------------------------

// Package frontend 将 Vite 构建产物（resources/frontend/dist）以 go:embed 形式内嵌进二进制，
// 供 route 包在 NoRoute 阶段提供静态资源托管与 SPA fallback 能力。
//
// 注意：dist 目录下至少保留一个文件（当前为 .gitkeep 空占位），否则 go:embed 将在编译期失败。
// 构建产物不入库，开发阶段若未执行 `make build-frontend`，静态资源不可用（FileServer 返回 404）。
package frontend

import (
	"embed"
	"io/fs"
)

// distFS 内嵌的前端构建产物，根目录等价于 resources/frontend/dist。
//
//go:embed all:dist
var distFS embed.FS

// FS 返回以 dist 为根的只读文件系统。访问 index.html 时使用路径 "index.html"（无前导斜杠）。
func FS() fs.FS {
	sub, err := fs.Sub(distFS, "dist")
	if err != nil {
		// 仅在 go:embed 声明与目录结构不一致时发生，属于编译期可发现的静态错误
		panic("frontend: failed to subtree dist: " + err.Error())
	}
	return sub
}
