// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// LinkContext is the golang structure of table xf_link_context for DAO operations like Where/Data.
type LinkContext struct {
	g.Meta           `orm:"table:xf_link_context, do:true"`
	LinkUuid         interface{} // 友链唯一标识符
	LinkName         interface{} // 友链名称
	LinkUrl          interface{} // 友链URL地址
	LinkAvatar       interface{} // 友链头像URL
	LinkRss          interface{} // 友链RSS地址
	LinkDesc         interface{} // 友链描述
	LinkEmail        interface{} // 友链联系邮箱
	LinkGroupUuid    interface{} // 所属分组ID
	LinkColorUuid    interface{} // 友链颜色ID
	LinkOrder        interface{} // 友链排序
	LinkStatus       interface{} // 友链状态（0: 待审核, 1: 已通过, 2: 已拒绝）
	LinkFail         interface{} // 友链失效标志（0: 正常, 1: 失效）
	LinkFailReason   interface{} // 友链失效原因
	LinkApplyRemark  interface{} // 申请者备注
	LinkReviewRemark interface{} // 审核备注
	LinkCreatedAt    *gtime.Time // 友链创建时间
	LinkUpdatedAt    *gtime.Time // 友链更新时间
}
