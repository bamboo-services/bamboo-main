/*
 * --------------------------------------------------------------------------------
 * Copyright (c) 2016-NOW(至今) 筱锋
 * Author: 筱锋「xiao_lfeng」(https://www.x-lf.com)
 * --------------------------------------------------------------------------------
 * 许可证声明：版权所有 (c) 2016-2026 筱锋。保留所有权利。
 * 有关MIT许可证的更多信息，请查看项目根目录下的LICENSE文件或访问：
 * https://opensource.org/licenses/MIT
 * --------------------------------------------------------------------------------
 */

package apiInfo

import "time"

// SiteUpdateRequest 站点信息更新请求
type SiteUpdateRequest struct {
	SiteName        *string `json:"site_name" binding:"omitempty,min=1,max=100" example:"筱锋的小站"`
	SiteDescription *string `json:"site_description" binding:"omitempty,max=500" example:"一个有趣的个人博客"`
	Introduction    *string `json:"introduction" binding:"omitempty,max=2000" example:"欢迎来到我的主页！"`
}

// AboutUpdateRequest 自我介绍更新请求
type AboutUpdateRequest struct {
	Content string `json:"content" binding:"required,min=1,max=10000" example:"# 关于我\n我是筱锋..."`
}

// SiteResponse 站点信息响应
type SiteResponse struct {
	SiteName        string    `json:"site_name"`
	SiteDescription string    `json:"site_description"`
	Introduction    string    `json:"introduction"`
	UpdatedAt       time.Time `json:"updated_at"`
}

// AboutResponse 自我介绍响应
type AboutResponse struct {
	Content   string    `json:"content"`
	UpdatedAt time.Time `json:"updated_at"`
}

// BloggerUpdateRequest 博主信息更新请求
//
// 博主信息用于「关于我」名士帖个人展示：昵称/个人简介/博客链接/头像。
// 与 SiteUpdateRequest（系统站点信息）及 ApplySiteUpdateRequest（申请站点展示）语义独立，故单独建模。
type BloggerUpdateRequest struct {
	Nick        *string `json:"nick" binding:"omitempty,max=100" example:"筱锋"`
	Description *string `json:"description" binding:"omitempty,max=500" example:"一个热爱技术的开发者，喜欢折腾各种有趣的东西。"`
	BlogUrl     *string `json:"blog_url" binding:"omitempty,max=500" example:"https://blog.x-lf.com"`
	Avatar      *string `json:"avatar" binding:"omitempty,max=500" example:"https://example.com/avatar.png"`
}

// BloggerResponse 博主信息响应
type BloggerResponse struct {
	Nick        string    `json:"nick"`
	Description string    `json:"description"`
	BlogUrl     string    `json:"blog_url"`
	Avatar      string    `json:"avatar"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// ApplySiteUpdateRequest 申请站点展示更新请求
//
// 申请站点展示用于「交换友链」场景：访客申请友链前需先在自站添加博主友链，
// 此处的站点名字/描述/地址/图片/订阅/邮箱即供访客复制添加的博主站点资料，
// 展示于 operate/apply 申请页。与 SiteUpdateRequest（系统站点信息）及
// BloggerUpdateRequest（博主个人展示）语义独立，故单独建模。
type ApplySiteUpdateRequest struct {
	SiteName        *string `json:"site_name" binding:"omitempty,max=100" example:"凌中的锋雨"`
	SiteDescription *string `json:"site_description" binding:"omitempty,max=500" example:"不为如何，只为在茫茫人海中有自己的一片天空~"`
	SiteUrl         *string `json:"site_url" binding:"omitempty,max=500" example:"https://www.x-lf.com"`
	SiteImage       *string `json:"site_image" binding:"omitempty,max=500" example:"https://example.com/logo.png"`
	Rss             *string `json:"rss" binding:"omitempty,max=500" example:"https://blog.example.com/atom.xml"`
	Email           *string `json:"email" binding:"omitempty,max=200" example:"gm@x-lf.cn"`
}

// ApplySiteResponse 申请站点展示响应
type ApplySiteResponse struct {
	SiteName        string    `json:"site_name"`
	SiteDescription string    `json:"site_description"`
	SiteUrl         string    `json:"site_url"`
	SiteImage       string    `json:"site_image"`
	Rss             string    `json:"rss"`
	Email           string    `json:"email"`
	UpdatedAt       time.Time `json:"updated_at"`
}
