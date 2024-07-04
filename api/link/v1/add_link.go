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

package v1

import "github.com/gogf/gf/v2/frame/g"

// LinkAddReq
//
// # 添加新的链接请求参数
//
// 添加新的链接
//
// 参数：
// WebmasterEmail: 站长邮箱
// ServiceProvider: 服务商
// SiteName: 站点名称
// SiteURL: 站点URL
// SiteLogo: 站点Logo
// SiteDescription: 站点描述
// SiteRssURL: 站点RSS URL
// DesiredLocation: 期望位置
// DesiredColor: 期望颜色
// HasAdv: 是否有广告
// Remark: 备注
type LinkAddReq struct {
	g.Meta          `path:"/link" method:"post" tags:"链接控制器" summary:"添加新的链接" `
	WebmasterEmail  string `json:"webmaster_email" v:"required|email#请输入站长邮箱|站长邮箱格式不正确" dc:"站长邮箱"`
	ServiceProvider string `json:"service_provider" v:"required#请输入服务商" dc:"服务商"`
	SiteName        string `json:"site_name" v:"required#请输入站点名称" dc:"站点名称"`
	SiteURL         string `json:"site_url" v:"required|url#请输入站点URL|站点URL格式不正确" dc:"站点URL"`
	SiteLogo        string `json:"site_logo" v:"required|url#请输入站点Logo|站点Logo格式不正确" dc:"站点Logo"`
	SiteDescription string `json:"site_description" v:"required#请输入站点描述" dc:"站点描述"`
	SiteRssURL      string `json:"site_rss_url" v:"url#站点RSS格式不正确" dc:"站点RSS URL"`
	DesiredLocation string `json:"desired_location" v:"required#请输入期望位置" dc:"期望位置"`
	DesiredColor    string `json:"desired_color" v:"required#请输入期望颜色" dc:"期望颜色"`
	HasAdv          bool   `json:"has_adv" v:"required#请输入是否有广告" dc:"是否有广告"`
	Remark          string `json:"remark" v:"required#请输入备注" dc:"备注"`
}

// LinkAddRes 添加新的链接响应参数
// 添加新的链接
type LinkAddRes struct {
	g.Meta `mime:"application/json"`
}
