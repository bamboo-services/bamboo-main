// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// XfLinkList is the golang structure for table xf_link_list.
type XfLinkList struct {
	Id              int64       `json:"id"              ` // 主键
	WebmasterEmail  string      `json:"webmasterEmail"  ` // 站长邮箱
	ServiceProvider string      `json:"serviceProvider" ` // 服务提供商
	SiteName        string      `json:"siteName"        ` // 站点名字
	SiteUrl         string      `json:"siteUrl"         ` // 站点地址
	SiteLogo        string      `json:"siteLogo"        ` // 站点 logo
	SiteDescription string      `json:"siteDescription" ` // 站点描述
	SiteRssUrl      string      `json:"siteRssUrl"      ` // 站点订阅地址
	HasAdv          bool        `json:"hasAdv"          ` // 是否有广告
	DesiredLocation int         `json:"desiredLocation" ` // 理想位置
	Location        int         `json:"location"        ` // 所在位置
	DesiredColor    int         `json:"desiredColor"    ` // 理想颜色
	Color           int         `json:"color"           ` // 颜色
	WebmasterRemark string      `json:"webmasterRemark" ` // 站长留言
	Remark          string      `json:"remark"          ` // 我的留言
	Status          int         `json:"status"          ` // 0 待审核，1 通过，-1 审核拒绝
	CreatedAt       *gtime.Time `json:"createdAt"       ` // 创建时间
	UpdatedAt       *gtime.Time `json:"updatedAt"       ` // 修改时间
	DeletedAt       *gtime.Time `json:"deletedAt"       ` // 删除时间
}
