/*
 * --------------------------------------------------------------------------------
 * Copyright (c) 2016-NOW(至今) 筱锋
 * Author: 筱锋「xiao_lfeng」(https://www.x-lf.com)
 * --------------------------------------------------------------------------------
 * 许可证声明：版权所有 (c) 2016-2025 筱锋。保留所有权利。
 * 有关MIT许可证的更多信息，请查看项目根目录下的LICENSE文件或访问：
 * https://opensource.org/licenses/MIT
 * --------------------------------------------------------------------------------
 */

package request

// SiteInfoUpdateReq 站点信息更新请求
type SiteInfoUpdateReq struct {
	SiteName        *string `json:"site_name" binding:"omitempty,min=1,max=100" example:"筱锋的小站"`
	SiteDescription *string `json:"site_description" binding:"omitempty,max=500" example:"一个有趣的个人博客"`
	Introduction    *string `json:"introduction" binding:"omitempty,max=2000" example:"欢迎来到我的主页！"`
}

// AboutUpdateReq 自我介绍更新请求
type AboutUpdateReq struct {
	Content string `json:"content" binding:"required,min=1,max=10000" example:"# 关于我\n我是筱锋..."`
}
