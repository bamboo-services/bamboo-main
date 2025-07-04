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

package main

import (
	"bamboo-main/internal/config/cmd"
	"bamboo-main/internal/config/setup"
	"github.com/gogf/gf/v2/os/gctx"

	_ "bamboo-main/internal/logic"
	_ "bamboo-main/internal/packed"
	_ "github.com/gogf/gf/contrib/drivers/pgsql/v2"
	_ "github.com/gogf/gf/contrib/nosql/redis/v2"
)

// main
//
// 系统服务启动
func main() {
	// 创建新的系统上下文
	getSystemContext := gctx.GetInitCtx()

	// 系统初始化服务
	setup.NewSetup(getSystemContext)

	// 启动服务
	cmd.Main.Run(getSystemContext)
}
