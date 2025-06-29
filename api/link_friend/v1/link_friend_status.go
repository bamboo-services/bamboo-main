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

type LinkFriendStatusReq struct {
	g.Meta   `path:"/link/friend/status" method:"PATCH" tags:"链接控制器" summary:"更新友情链接状态" dc:"更新友情链接状态的接口，允许管理员修改友情链接的审核状态和显示状态"`
	LinkUUID string `json:"link_uuid" v:"required|regex:^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$#请输入有效的UUID格式|请输入有效的UUID格式" dc:"友情链接唯一标识符，不能为空，用于指定要更新的友情链接"`
	Status   int    `json:"status" v:"required|in:0,1,2#请输入友情链接状态|友情链接状态只允许 0（待审核）、1（已通过）、2（已拒绝）" dc:"友情链接状态，0表示待审核，1表示已通过，2表示已拒绝，不能为空"`
}

type LinkFriendStatusRes struct {
	g.Meta `mime:"application/json; charset=utf-8" dc:"更新友情链接状态响应"`
	*bmodels.ResponseDTO[types.Nil]
}
