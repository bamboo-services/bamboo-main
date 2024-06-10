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

package flow

// WebInfoDTO 站点信息返回
// 返回一些基本的站点信息内容
//
// 参数:
// Main 主要配置信息
// SiteBloggerDTO 站长的一些个人信息
type WebInfoDTO struct {
	Main    *SiteMainDTO    `json:"site" dc:"主要配置信息"`
	Blogger *SiteBloggerDTO `json:"blogger" dc:"站长个人信息"`
}

// SiteMainDTO 主要配置信息
// 包含了一些主要的配置信息进行返回
//
// 参数:
// SiteName 站点名字
// Author 软件作者
// Version 软件版本
// Description 站点描述
// Keywords 站点关键字
type SiteMainDTO struct {
	SiteName    string `json:"site_name" dc:"站点名称"`
	Author      string `json:"author" dc:"软件作者"`
	Version     string `json:"version" dc:"软件版本"`
	Description string `json:"description" dc:"站点描述"`
	Keywords    string `json:"keywords" dc:"站点关键词"`
}

// SiteBloggerDTO 站长信息
// 包含了站长一些基础信息需要进行站的内容
//
// 参数:
// Name 站长名字
// Nick 站长昵称
// Email 站长邮箱
// Description 站长的一些描述（或者说座右铭）
type SiteBloggerDTO struct {
	Name        string `json:"name" dc:"站长名字"`
	Nick        string `json:"nick" dc:"站长昵称"`
	Email       string `json:"email" dc:"站长邮箱"`
	Description string `json:"description" dc:"站长描述（座右铭）"`
}
