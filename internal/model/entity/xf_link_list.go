// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// XfLinkList is the golang structure for table xf_link_list.
type XfLinkList struct {
	Id              int64       `json:"id"              orm:"id"               ` // 主键
	WebmasterEmail  string      `json:"webmasterEmail"  orm:"webmaster_email"  ` // 站长邮箱
	ServiceProvider string      `json:"serviceProvider" orm:"service_provider" ` // 服务提供商
	SiteName        string      `json:"siteName"        orm:"site_name"        ` // 站点名字
	SiteUrl         string      `json:"siteUrl"         orm:"site_url"         ` // 站点地址
	SiteLogo        string      `json:"siteLogo"        orm:"site_logo"        ` // 站点 logo
	SiteDescription string      `json:"siteDescription" orm:"site_description" ` // 站点描述
	SiteRssUrl      string      `json:"siteRssUrl"      orm:"site_rss_url"     ` // 站点订阅地址
	HasAdv          bool        `json:"hasAdv"          orm:"has_adv"          ` // 是否有广告
	DesiredLocation int         `json:"desiredLocation" orm:"desired_location" ` // 理想位置
	Location        int         `json:"location"        orm:"location"         ` // 所在位置
	DesiredColor    int         `json:"desiredColor"    orm:"desired_color"    ` // 理想颜色
	Color           int         `json:"color"           orm:"color"            ` // 颜色
	WebmasterRemark string      `json:"webmasterRemark" orm:"webmaster_remark" ` // 站长留言
	Remark          string      `json:"remark"          orm:"remark"           ` // 我的留言
	Status          int         `json:"status"          orm:"status"           ` // 0 待审核，1 通过，-1 审核拒绝
	CreatedAt       *gtime.Time `json:"createdAt"       orm:"created_at"       ` // 创建时间
	UpdatedAt       *gtime.Time `json:"updatedAt"       orm:"updated_at"       ` // 修改时间
	DeletedAt       *gtime.Time `json:"deletedAt"       orm:"deleted_at"       ` // 删除时间
}
