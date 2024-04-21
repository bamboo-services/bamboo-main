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
	CdnLogoUrl      interface{} // 镜像站点 logo
	SiteDescription interface{} // 站点描述
	SiteRssUrl      interface{} // 站点订阅地址
	HasAdv          interface{} // 是否有广告
	DesiredLocation interface{} // 理想位置
	Location        interface{} // 所在位置
	DesiredColor    interface{} // 理想颜色
	Color           interface{} // 颜色
	WebmasterRemark interface{} // 站长留言
	Remark          interface{} // 我的留言
	Status          interface{} // 0 待审核，1 通过，-1 审核拒绝，2 回收站
	AbleConnect     interface{} // 能否连接
	CreatedAt       *gtime.Time // 创建时间
	UpdatedAt       *gtime.Time // 修改时间
	DeletedAt       *gtime.Time // 删除时间
}
