// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// XfLinkList is the golang structure of table xf_link_list for DAO operations like Where/Data.
type XfLinkList struct {
	g.Meta          `orm:"table:xf_link_list, do:true"`
	Id              interface{} // 主键
	WebmasterEmail  interface{} // 站长邮箱
	ServiceProvider interface{} // 服务提供商
	SiteName        interface{} // 站点名字
	SiteUrl         interface{} // 站点地址
	SiteLogo        interface{} // 站点 logo
	SiteDescription interface{} // 站点描述
	SiteRssUrl      interface{} // 站点订阅地址
	HasAdv          interface{} // 是否有广告
	DesiredLocation interface{} // 理想位置
	Location        interface{} // 所在位置
	DesiredColor    interface{} // 理想颜色
	Color           interface{} // 颜色
	WebmasterRemark interface{} // 站长留言
	Remark          interface{} // 我的留言
	Status          interface{} // 0 待审核，1 通过，-1 审核拒绝
	CreatedAt       *gtime.Time // 创建时间
	UpdatedAt       *gtime.Time // 修改时间
	DeletedAt       *gtime.Time // 删除时间
}
