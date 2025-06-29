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

package base

import "github.com/gogf/gf/v2/os/gtime"

// LinkFriendDTO
//
// 表示一个友链的数据信息传输对象。
// 包含友链的唯一标识符、名称、URL地址、头像URL、描述、联系邮箱、所属分组ID、颜色ID、排序、状态、申请者备注、审核备注以及创建和更新时间等字段。
type LinkFriendDTO struct {
	LinkUuid         string      `json:"link_uuid"          description:"友链唯一标识符" v:"required|regex:^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$"` // 友链唯一标识符
	LinkName         string      `json:"link_name"          description:"友链名称" v:"required|max-length:100"`                                                                         // 友链名称
	LinkUrl          string      `json:"link_url"           description:"友链URL地址" v:"required|url"`                                                                                 // 友链URL地址
	LinkAvatar       string      `json:"link_avatar"        description:"友链头像URL"`                                                                                                  // 友链头像URL
	LinkDesc         string      `json:"link_desc"          description:"友链描述"`                                                                                                     // 友链描述
	LinkEmail        string      `json:"link_email"         description:"友链联系邮箱" v:"email"`                                                                                         // 友链联系邮箱
	LinkGroupUuid    string      `json:"link_group_uuid"    description:"所属分组ID" v:"regex:^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$"`           // 所属分组ID
	LinkColorUuid    string      `json:"link_color_uuid"    description:"友链颜色ID" v:"regex:^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$"`           // 友链颜色ID
	LinkOrder        int         `json:"link_order"         description:"友链排序" v:"between:0,100"`                                                                                   // 友链排序
	LinkStatus       int         `json:"link_status"        description:"友链状态（0: 待审核, 1: 已通过, 2: 已拒绝）" v:"required|regex:^[012]$"`                                                  // 友链状态（0: 待审核, 1: 已通过, 2: 已拒绝）
	LinkApplyRemark  string      `json:"link_apply_remark"  description:"申请者备注" v:"max-length:10240"`                                                                               // 申请者备注
	LinkReviewRemark string      `json:"link_review_remark" description:"审核备注" v:"max-length:10240"`                                                                                // 审核备注
	LinkCreatedAt    *gtime.Time `json:"link_created_at"    description:"友链创建时间"`                                                                                                   // 友链创建时间
	LinkUpdatedAt    *gtime.Time `json:"link_updated_at"    description:"友链更新时间"`                                                                                                   // 友链更新时间
}
