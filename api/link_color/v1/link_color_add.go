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
	"bamboo-main/internal/model/dto/base"
	"github.com/XiaoLFeng/bamboo-utils/bmodels"
	"github.com/gogf/gf/v2/frame/g"
)

type LinkColorAddReq struct {
	g.Meta      `path:"/link/color" method:"POST" tags:"链接颜色控制器" summary:"添加链接颜色" dc:"添加友链颜色的接口，允许用户创建新的链接颜色"`
	Name        string `json:"name" v:"required#请输入颜色名称" dc:"颜色名称，不能为空"`
	Value       string `json:"value" v:"required|regex:^#[0-9A-Fa-f]{6}$#请输入颜色值|请输入有效的HEX颜色值" dc:"颜色值，如HEX值：#FFFFFF，不能为空"`
	Description string `json:"description" dc:"颜色描述，表示颜色的详细描述信息"`
}

type LinkColorAddRes struct {
	g.Meta `mime:"application/json; charset=utf-8" dc:"添加链接颜色响应"`
	*bmodels.ResponseDTO[base.LinkColorDTO]
}
