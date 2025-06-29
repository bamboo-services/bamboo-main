// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// LinkContext is the golang structure for table link_context.
type LinkContext struct {
	LinkUuid         string      `json:"link_uuid"          orm:"link_uuid"          description:"友链唯一标识符"`                      // 友链唯一标识符
	LinkName         string      `json:"link_name"          orm:"link_name"          description:"友链名称"`                         // 友链名称
	LinkUrl          string      `json:"link_url"           orm:"link_url"           description:"友链URL地址"`                      // 友链URL地址
	LinkAvatar       string      `json:"link_avatar"        orm:"link_avatar"        description:"友链头像URL"`                      // 友链头像URL
	LinkRss          string      `json:"link_rss"           orm:"link_rss"           description:"友链RSS地址"`                      // 友链RSS地址
	LinkDesc         string      `json:"link_desc"          orm:"link_desc"          description:"友链描述"`                         // 友链描述
	LinkEmail        string      `json:"link_email"         orm:"link_email"         description:"友链联系邮箱"`                       // 友链联系邮箱
	LinkGroupUuid    string      `json:"link_group_uuid"    orm:"link_group_uuid"    description:"所属分组ID"`                       // 所属分组ID
	LinkColorUuid    string      `json:"link_color_uuid"    orm:"link_color_uuid"    description:"友链颜色ID"`                       // 友链颜色ID
	LinkOrder        int         `json:"link_order"         orm:"link_order"         description:"友链排序"`                         // 友链排序
	LinkStatus       int         `json:"link_status"        orm:"link_status"        description:"友链状态（0: 待审核, 1: 已通过, 2: 已拒绝）"` // 友链状态（0: 待审核, 1: 已通过, 2: 已拒绝）
	LinkFail         int         `json:"link_fail"          orm:"link_fail"          description:"友链失效标志（0: 正常, 1: 失效）"`         // 友链失效标志（0: 正常, 1: 失效）
	LinkFailReason   string      `json:"link_fail_reason"   orm:"link_fail_reason"   description:"友链失效原因"`                       // 友链失效原因
	LinkApplyRemark  string      `json:"link_apply_remark"  orm:"link_apply_remark"  description:"申请者备注"`                        // 申请者备注
	LinkReviewRemark string      `json:"link_review_remark" orm:"link_review_remark" description:"审核备注"`                         // 审核备注
	LinkCreatedAt    *gtime.Time `json:"link_created_at"    orm:"link_created_at"    description:"友链创建时间"`                       // 友链创建时间
	LinkUpdatedAt    *gtime.Time `json:"link_updated_at"    orm:"link_updated_at"    description:"友链更新时间"`                       // 友链更新时间
}
