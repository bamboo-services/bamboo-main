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

type LinkColorDeleteReq struct {
	g.Meta    `path:"/link/color" method:"DELETE" tags:"链接颜色控制器" summary:"删除链接颜色" dc:"删除友链颜色的接口，允许用户删除指定的链接颜色"`
	ColorUuid string `json:"color_uuid" in:"query" v:"required#请输入颜色唯一标识符|regex:^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$#请输入有效的UUID格式|请输入有效的UUID格式" dc:"颜色唯一标识符，不能为空，用于指定要删除的颜色"`
}

type LinkColorDeleteRes struct {
	g.Meta `mime:"application/json; charset=utf-8" dc:"删除链接颜色响应"`
	*bmodels.ResponseDTO[types.Nil]
}
