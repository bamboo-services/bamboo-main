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

package startup

import (
	"context"
	"github.com/bamboo-services/bamboo-utils/butil"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/util/gconv"
	"xiaoMain/internal/constants"
)

// InitStruct
//
// # 初始化结构体
//
// 该结构体为初始化启动准备的结构体。
type InitStruct struct {
	CTX context.Context
}

// New
//
// # 新建初始化结构体
//
// 该结构体为初始化启动准备的结构体。
func New(ctx context.Context) *InitStruct {
	return &InitStruct{
		CTX: ctx,
	}
}

// InitialDatabase
//
// # 初始化数据库
//
// 是一个初始化数据库的函数。
// 它会检查数据库表是否完整，并插入必要的数据。
// 这个函数在应用程序的启动过程中被调用。
func (is *InitStruct) InitialDatabase() {
	/*
	 * 检查数据表是否完善
	 */
	// 记录日志，开始初始化数据表
	g.Log().Notice(is.CTX, "[BOOT] 数据表初始化中")
	// 初始化信息表
	is.initialSQL(is.CTX, "xf_index")
	// 初始化登录信息表
	is.initialSQL(is.CTX, "xf_token")
	// 初始化验证码表
	is.initialSQL(is.CTX, "xf_verification_code")
	// 初始化日志表
	is.initialSQL(is.CTX, "xf_logs")
	// 初始化位置表
	is.initialSQL(is.CTX, "xf_location")
	// 初始化颜色表
	is.initialSQL(is.CTX, "xf_color")
	// 初始化链接表
	is.initialSQL(is.CTX, "xf_link_list")
	// 初始化 RSS 订阅表
	is.initialSQL(is.CTX, "xf_link_rss")
}

// InitialIndexTable
//
// # 初始化主要数据表
//
// 是一个初始化主要数据表的函数。
// 它会插入主要数据表的数据。
// 这个函数在应用程序的启动过程中被调用。
func (is *InitStruct) InitialIndexTable() {
	/**
	 * 插入主要数据表 index
	 */
	// 记录日志，开始初始化数据库表信息
	g.Log().Notice(is.CTX, "[BOOT] 数据库表信息初始化中")
	g.Log().Notice(is.CTX, "[BOOT] 初始化信息表")
	// 插入软件版本信息
	is.insertIndexData(is.CTX, "version", constants.XiaoMainVersion)
	// 插入软件作者信息
	is.insertIndexData(is.CTX, "author", constants.XiaoMainAuthor)
	// 插入站点信息
	is.insertIndexData(is.CTX, "site_name", "XiaoMain")
	// 插入站点描述
	is.insertIndexData(is.CTX, "description", "一个由Go开发的开源项目，用于快速搭建个人页面、介绍的网站。")
	// 插入站点关键字
	is.insertIndexData(is.CTX, "keywords", "XiaoMain,筱锋,开源项目,Go,个人网站,介绍")

	// 生成并插入用户的唯一 UUID
	is.insertIndexData(is.CTX, "uuid", butil.GenerateRandUUID().String())
	// 新建默认用户
	is.insertIndexData(is.CTX, "user", "admin")
	// 设置默认用户密码
	is.insertIndexData(is.CTX, "password", butil.PasswordEncode("admin-admin"))
	// 设置默认用户邮箱
	is.insertIndexData(is.CTX, "email", "admin@xiaoMain.com")
	// 站长昵称
	is.insertIndexData(is.CTX, "blogger_name", "xiao_lfeng")
	// 站长展示名字
	is.insertIndexData(is.CTX, "blogger_nick", "筱锋")
	// 站长座右铭（一句话介绍）
	is.insertIndexData(is.CTX, "blogger_description", "愿你的人生如璀璨星辰，勇敢梦想，坚韧追寻，发现真我之光。在每个转角，以坚定意志和无尽创意，拥抱挑战，种下希望。遇见激励之人，共成长；面对风雨，保持乐观，让每一步都踏出意义深远的足迹。在追梦途中，听从内心之声，珍惜遇见，让生活不仅是冒险，更是自我发现的诗篇。") //nolint:lll

	// 设置允许登录的节点数（设备数）
	is.insertIndexData(is.CTX, "auth_limit", "3")

	// SMTP 邮件服务器配置
	is.insertIndexData(is.CTX, "smtp_host", "smtp.x-lf.com")
	// SMTP 邮件服务器端口(默认)
	is.insertIndexData(is.CTX, "smtp_port_tls", "25")
	// SMTP 邮件服务器端口(SSL)
	is.insertIndexData(is.CTX, "smtp_port_ssl", "465")
	// SMTP 邮件服务器用户名
	is.insertIndexData(is.CTX, "smtp_user", "noreplay@xiaoMain.com")
	// SMTP 邮件服务器密码
	is.insertIndexData(is.CTX, "smtp_pass", "password")

	// 初始化邮件模板标题(user-change-password)
	is.insertIndexData(is.CTX, "mail_template_user_change_password_title", "修改密码")
	// 初始化邮件模板(user-change-password)
	is.insertIndexData(is.CTX, "mail_template_user_change_password", is.getMailTemplate("user_change_password"))
	// 初始化邮件模板标题(user-forget-password)
	is.insertIndexData(is.CTX, "mail_template_user_reset_password_title", "重置密码")
	// 初始化邮件模板(user-forget-password)
	is.insertIndexData(is.CTX, "mail_template_user_reset_password", is.getMailTemplate("user_reset_password"))
}

// InitialLocationTable
//
// # 初始化位置表
//
// 是一个初始化位置表的函数。
// 它会插入位置表的数据。
// 这个函数在应用程序的启动过程中被调用。
func (is *InitStruct) InitialLocationTable() {
	// 记录日志，开始初始化数据表
	g.Log().Notice(is.CTX, "[BOOT] 初始化位置表")

	// 初始化期望位置表
	is.insertLocationData(is.CTX, 1, "favorite", "最喜欢", "这是我最喜欢的东西，我当然要置顶啦", true)
	is.insertLocationData(is.CTX, 100, "fail", "失效的", "这是失效的友链，希望你快回来嗷", false)
}

// InitialColorTable
//
// # 初始化颜色表
//
// 是一个初始化颜色表的函数。
// 它会插入颜色表的数据。
// 这个函数在应用程序的启动过程中被调用。
func (is *InitStruct) InitialColorTable() {
	g.Log().Notice(is.CTX, "[BOOT] 初始化颜色表")

	// 初始化期望颜色表
	is.insertColorData(is.CTX, "red", "红色", "FF0000")
	is.insertColorData(is.CTX, "orange", "橙色", "FFA500")
	is.insertColorData(is.CTX, "yellow", "黄色", "FFFF00")
	is.insertColorData(is.CTX, "green", "绿色", "008000")
	is.insertColorData(is.CTX, "cyan", "青色", "00FFFF")
	is.insertColorData(is.CTX, "blue", "蓝色", "0000FF")
	is.insertColorData(is.CTX, "purple", "紫色", "800080")
	is.insertColorData(is.CTX, "pink", "粉色", "FFC0CB")
	is.insertColorData(is.CTX, "black", "黑色", "000000")
	is.insertColorData(is.CTX, "white", "白色", "FFFFFF")
	is.insertColorData(is.CTX, "gray", "灰色", "808080")
	is.insertColorData(is.CTX, "brown", "棕色", "A52A2A")
	is.insertColorData(is.CTX, "gold", "金色", "FFD700")
	is.insertColorData(is.CTX, "silver", "银色", "C0C0C0")
	is.insertColorData(is.CTX, "bronze", "铜色", "CD7F32")
	is.insertColorData(is.CTX, "rose", "玫瑰金", "FFC0CB")
	is.insertColorData(is.CTX, "champagne", "香槟金", "F7E7CE")
	is.insertColorData(is.CTX, "peach", "桃红", "FFDAB9")
	is.insertColorData(is.CTX, "apricot", "杏色", "FBCEB1")
	is.insertColorData(is.CTX, "coral", "珊瑚红", "FF7F50")
	is.insertColorData(is.CTX, "salmon", "鲑鱼红", "FA8072")
	is.insertColorData(is.CTX, "tomato", "番茄红", "FF6347")
	is.insertColorData(is.CTX, "maroon", "栗色", "800000")
	is.insertColorData(is.CTX, "burgundy", "酒红", "800020")
	is.insertColorData(is.CTX, "ruby", "红宝石", "E0115F")
	is.insertColorData(is.CTX, "sapphire", "蓝宝石", "0F52BA")
	is.insertColorData(is.CTX, "emerald", "翡翠绿", "50C878")
	is.insertColorData(is.CTX, "amethyst", "紫水晶", "9966CC")
	is.insertColorData(is.CTX, "topaz", "黄玉", "FFC87C")
	is.insertColorData(is.CTX, "turquoise", "绿松石", "40E0D0")
	is.insertColorData(is.CTX, "aquamarine", "海蓝宝石", "7FFFD4")
	is.insertColorData(is.CTX, "peridot", "橄榄石", "E6E200")
	is.insertColorData(is.CTX, "opal", "蛋白石", "A8C3BC")
	is.insertColorData(is.CTX, "pearl", "珍珠", "F0EAD6")
	is.insertColorData(is.CTX, "moonstone", "月光石", "E3E4FA")
	is.insertColorData(is.CTX, "diamond", "钻石", "B9F2FF")
}

// InitCommonData
//
// # 初始化常量数据
//
// 是一个初始化常量数据的函数。
// 将会从数据库 Index 中读取常量数据。
// 这个函数在应用程序的启动过程中被调用。
func (is *InitStruct) InitCommonData() {
	g.Log().Notice(is.CTX, "[BOOT] 初始化常量数据")

	// 记录日志，开始初始化常用数据
	constants.SMTPHost = gconv.String(is.prepareCommonData("smtp_host"))
	constants.SMTPPortTLS = gconv.Int(is.prepareCommonData("smtp_port_tls"))
	constants.SMTPPortSSL = gconv.Int(is.prepareCommonData("smtp_port_ssl"))
	constants.SMTPUser = gconv.String(is.prepareCommonData("smtp_user"))
	constants.SMTPPass = gconv.String(is.prepareCommonData("smtp_pass"))
}
