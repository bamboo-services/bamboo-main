package apiInfo

import "github.com/bamboo-services/bamboo-main/internal/models/dto"

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
	dto.SiteInfoDTO
}

// AboutResponse 自我介绍响应
type AboutResponse struct {
	dto.AboutDTO
}
