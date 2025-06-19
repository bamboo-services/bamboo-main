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

package setup

import (
	"bamboo-main/internal/consts"
	"context"
	"fmt"
)

// projectSetup
//
// 用于项目的初始化设置
type projectSetup struct {
	ctx context.Context
}

// NewSetup
//
// 创建一个新的 projectSetup 实例
func NewSetup(ctx context.Context) {
	var newSetup = &projectSetup{ctx: ctx}

	newSetup.ProjectSetupForSQL()
	newSetup.ProjectSetupSystemValueInit()

	/*
	 * 系统初始化完成
	 */
	fmt.Print(`
` + "\033[1;32m" + `   _  __ _             ` + "\033[1;34m" + `__  ___      _     
` + "\033[1;32m" + `  | |/ /(_)___ _____  ` + "\033[1;34m" + `/  |/  /___ _(_)___ 
` + "\033[1;32m" + `  |   // / __ ` + "`" + `/ __ \` + "\033[1;34m" + `/ /|_/ / __ ` + "`" + `/ / __ \
` + "\033[1;32m" + ` /   |/ / /_/ / /_/ ` + "\033[1;34m" + `/ /  / / /_/ / / / / /
` + "\033[1;32m" + `/_/|_/_/\__,_/\____` + "\033[1;34m" + `/_/  /_/\__,_/_/_/ /_/
`)
	fmt.Println("\033[1;33m   ::: XiaoMain :::	::: " + consts.SystemVersion + " :::")
	fmt.Println("\033[0m")
}
