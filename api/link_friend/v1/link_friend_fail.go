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

type LinkFriendFailReq struct {
	g.Meta   `path:"/link/friend/fail" method:"PATCH" tags:"链接控制器" summary:"更新友情链接失败状态" dc:"更新友情链接失败状态的接口，允许管理员将友情链接标记为失败状态，通常用于处理审核未通过的链接请求"`
	LinkUUID string `json:"link_uuid" v:"required|regex:^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$#请输入有效的UUID格式|请输入有效的UUID格式" dc:"友情链接唯一标识符，不能为空，用于指定要更新的友情链接"`
	IsFail   bool   `json:"is_fail" v:"required#请输入是否失败状态" dc:"是否失败状态，不能为空，表示友情链接是否被标记为失败状态，通常用于审核未通过的情况"`
	Reason   string `json:"reason" v:"required|max-length:1024#请输入失败原因|失败原因长度不能超过1024个字符" dc:"失败原因，不能为空，表示友情链接审核未通过的具体原因，长度不能超过1024个字符"`
}

type LinkFriendFailRes struct {
	g.Meta `mime:"application/json; charset=utf-8" dc:"更新友情链接失败状态响应"`
	*bmodels.ResponseDTO[types.Nil]
}
