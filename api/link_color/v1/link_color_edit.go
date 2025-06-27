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

package v1

import (
	"github.com/XiaoLFeng/bamboo-utils/bmodels"
	"github.com/gogf/gf/v2/frame/g"
	"go/types"
)

type LinkColorEditReq struct {
	g.Meta      `path:"/link/color" method:"PUT" tags:"链接颜色控制器" summary:"编辑链接颜色" dc:"编辑友情链接颜色接口，允许用户修改现有的友情链接颜色信息，包括名称、颜色值和描述等属性"`
	UUID        string `json:"uuid" v:"required|regex:^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$#请输入有效的UUID格式|请输入有效的UUID格式" dc:"颜色唯一标识符，不能为空，用于标识要编辑的颜色"`
	Name        string `json:"name" v:"required#请输入颜色名称" dc:"颜色名称，不能为空，表示颜色的名称"`
	Value       string `json:"value" v:"required|regex:^#[0-9A-Fa-f]{6}$#请输入颜色值|请输入有效的HEX颜色值" dc:"颜色值，如HEX值：#FFFFFF，不能为空"`
	Description string `json:"description" dc:"颜色描述，表示颜色的详细描述信息"`
}

type LinkColorEditRes struct {
	g.Meta `mime:"application/json; charset=utf-8" dc:"编辑链接颜色响应"`
	*bmodels.ResponseDTO[types.Nil]
}
