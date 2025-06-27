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

type LinkGroupEditReq struct {
	g.Meta      `path:"/link/group" method:"PUT" tags:"链接控制器" summary:"编辑链接分组" dc:"编辑友情链接分组接口，允许用户修改现有的友情链接分组信息，包括名称、描述和排序等属性"`
	UUID        string `json:"uuid" v:"required|regex:^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$#请输入有效的UUID格式" dc:"分组唯一标识符，不能为空，用于标识要编辑的分组"`
	Name        string `json:"name" v:"required#请输入分组名称" dc:"分组名称，不能为空，表示分组的名称"`
	Description string `json:"description" v:"required#请输入分组描述" dc:"分组描述，不能为空，表示分组的详细描述信息"`
	Order       int    `json:"order" default:"0" v:"required|between:0,100#请输入分组排序|分组顺序只允许 0 到 100" dc:"分组排序，数字越小越靠前，默认为0，表示分组的显示顺序"`
}

type LinkGroupEditRes struct {
	g.Meta `mime:"application/json; charset=utf-8" dc:"编辑链接分组响应"`
	*bmodels.ResponseDTO[types.Nil]
}
