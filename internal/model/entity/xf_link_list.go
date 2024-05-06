/*
 * --------------------------------------------------------------------------------
 * Copyright (c) 2016-NOW(至今) 筱锋
 * Author: 筱锋(https://www.x-lf.com)
 *
 * 本文件包含 XiaoMain 的源代码，该项目的所有源代码均遵循MIT开源许可证协议。
 * --------------------------------------------------------------------------------
 * 许可证声明：
 *
 * 版权所有 (c) 2016-2024 筱锋。保留所有权利。
 *
 * 本软件是“按原样”提供的，没有任何形式的明示或暗示的保证，包括但不限于
 * 对适销性、特定用途的适用性和非侵权性的暗示保证。在任何情况下，
 * 作者或版权持有人均不承担因软件或软件的使用或其他交易而产生的、
 * 由此引起的或以任何方式与此软件有关的任何索赔、损害或其他责任。
 *
 * 使用本软件即表示您了解此声明并同意其条款。
 *
 * 有关MIT许可证的更多信息，请查看项目根目录下的LICENSE文件或访问：
 * https://opensource.org/licenses/MIT
 * --------------------------------------------------------------------------------
 * 免责声明：
 *
 * 使用本软件的风险由用户自担。作者或版权持有人在法律允许的最大范围内，
 * 对因使用本软件内容而导致的任何直接或间接的损失不承担任何责任。
 * --------------------------------------------------------------------------------
 */

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
	CdnLogoUrl      string      `json:"cdnLogoUrl"      orm:"cdn_logo_url"     ` // 镜像站点 logo
	SiteDescription string      `json:"siteDescription" orm:"site_description" ` // 站点描述
	SiteRssUrl      string      `json:"siteRssUrl"      orm:"site_rss_url"     ` // 站点订阅地址
	HasAdv          bool        `json:"hasAdv"          orm:"has_adv"          ` // 是否有广告
	DesiredLocation int64       `json:"desiredLocation" orm:"desired_location" ` // 理想位置
	Location        int64       `json:"location"        orm:"location"         ` // 所在位置
	DesiredColor    int64       `json:"desiredColor"    orm:"desired_color"    ` // 理想颜色
	Color           int64       `json:"color"           orm:"color"            ` // 颜色
	WebmasterRemark string      `json:"webmasterRemark" orm:"webmaster_remark" ` // 站长留言
	Remark          string      `json:"remark"          orm:"remark"           ` // 我的留言
	Status          int         `json:"status"          orm:"status"           ` // 0 待审核，1 通过，-1 审核拒绝，2 回收站
	AbleConnect     bool        `json:"ableConnect"     orm:"able_connect"     ` // 能否连接
	CreatedAt       *gtime.Time `json:"createdAt"       orm:"created_at"       ` // 创建时间
	UpdatedAt       *gtime.Time `json:"updatedAt"       orm:"updated_at"       ` // 修改时间
	DeletedAt       *gtime.Time `json:"deletedAt"       orm:"deleted_at"       ` // 删除时间
}
