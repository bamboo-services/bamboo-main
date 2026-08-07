// --------------------------------------------------------------------------------
// Copyright (c) 2016-NOW(至今) 筱锋
// Author: 筱锋「xiao_lfeng」(https://www.x-lf.com)
// --------------------------------------------------------------------------------
// 许可证声明：版权所有 (c) 2016-2026 筱锋。保留所有权利。
// 有关MIT许可证的更多信息，请查看项目根目录下的LICENSE文件或访问：
// https://opensource.org/licenses/MIT
// --------------------------------------------------------------------------------

// Package templates 将邮件模板目录（templates/mail）以 go:embed 形式内嵌进二进制，
// 供 xEmail 插件经 AddTemplateFS 加载渲染。
//
// 模板内嵌后单二进制部署无需运行时文件系统目录，Docker 镜像也无需 COPY templates。
package templates

import (
	"embed"
	"io/fs"
)

//go:embed all:mail
var mailFS embed.FS

// Mail 返回内嵌的邮件模板文件系统。
//
// FS 根目录等价于 templates/，访问单个模板使用路径 "mail/xxx.html"，
// 配合 xEmail 客户端 AddTemplateFS(fs, "mail/*.html") 使用。
func Mail() fs.FS {
	return mailFS
}
