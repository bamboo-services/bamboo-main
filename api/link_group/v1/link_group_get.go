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

type LinkGroupGetReq struct {
	g.Meta    `path:"/link/group" method:"GET" tags:"链接分组控制器" summary:"获取链接分组" dc:"获取友情链接分组的接口，允许用户查询现有的友情链接分组信息，包括名称、描述和排序等属性"`
	GroupUUID string `json:"group_uuid" in:"query" v:"required|regex:^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$#请输入有效的UUID格式|请输入有效的UUID格式" dc:"分组唯一标识符，不能为空，用于指定要查询的分组"`
}

type LinkGroupGetRes struct {
	g.Meta `mime:"application/json; charset=utf-8" dc:"获取链接分组响应"`
	*bmodels.ResponseDTO[base.LinkGroupDTO]
}
