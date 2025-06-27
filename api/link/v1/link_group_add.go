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

type LinkGroupAddReq struct {
	g.Meta      `path:"/link/group" method:"POST" tags:"链接控制器" summary:"添加链接分组" dc:"添加友链分组的接口，允许用户创建新的链接分组"`
	Name        string `json:"name" v:"required#请输入分组名称" dc:"分组名称，不能为空"`
	Description string `json:"description" v:"required#请输入分组描述" dc:"分组描述，不能为空"`
	Order       int    `json:"order" default:"0" v:"required|between:0,100#请输入分组排序|分组顺序只允许 0 到 100" dc:"分组排序，数字越小越靠前，默认为0"`
}

type LinkGroupAddRes struct {
	g.Meta `mime:"application/json; charset=utf-8" dc:"添加链接分组响应"`
	*bmodels.ResponseDTO[base.LinkGroupDTO]
}
