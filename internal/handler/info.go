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

package handler

import (
	apiInfo "github.com/bamboo-services/bamboo-main/api/info"

	xResult "github.com/bamboo-services/bamboo-base-go/major/result"
	xValid "github.com/bamboo-services/bamboo-base-go/major/validator"
	"github.com/gin-gonic/gin"
)

// GetSiteInfo 获取站点信息
//
// @Summary [用户] 获取站点信息
// @Description 获取站点名称、描述、主页介绍等公开信息
// @Tags 站点信息接口
// @Accept json
// @Produce json
// @Success 200 {object} xBase.BaseResponse{data=apiInfo.SiteResponse} "获取成功"
// @Failure 500 {object} xBase.BaseResponse "服务器内部错误"
// @Router /api/v1/info/site [GET]
func (h *InfoHandler) GetSiteInfo(c *gin.Context) {
	result, err := h.service.infoLogic.GetSiteInfo(c.Request.Context())
	if err != nil {
		_ = c.Error(err)
		return
	}

	xResult.SuccessHasData(c, "获取站点信息成功", result)
}

// UpdateSiteInfo 更新站点信息
//
// @Summary [管理] 更新站点信息
// @Description 管理员更新站点名称、描述、主页介绍
// @Tags 站点信息接口
// @Accept json
// @Produce json
// @Security Bearer
// @Param request body apiInfo.SiteUpdateRequest true "站点信息更新请求"
// @Success 200 {object} xBase.BaseResponse{data=apiInfo.SiteResponse} "更新成功"
// @Failure 400 {object} xBase.BaseResponse "请求参数错误"
// @Failure 401 {object} xBase.BaseResponse "未认证"
// @Failure 500 {object} xBase.BaseResponse "服务器内部错误"
// @Router /api/v1/admin/info/site [PUT]
func (h *InfoHandler) UpdateSiteInfo(c *gin.Context) {
	var req apiInfo.SiteUpdateRequest

	if bindErr := c.ShouldBindJSON(&req); bindErr != nil {
		xValid.HandleValidationError(c, bindErr)
		return
	}

	result, err := h.service.infoLogic.UpdateSiteInfo(c.Request.Context(), &req)
	if err != nil {
		_ = c.Error(err)
		return
	}

	xResult.SuccessHasData(c, "站点信息更新成功", result)
}

// GetAbout 获取自我介绍
//
// @Summary [用户] 获取自我介绍
// @Description 获取 Markdown 格式的自我介绍内容
// @Tags 站点信息接口
// @Accept json
// @Produce json
// @Success 200 {object} xBase.BaseResponse{data=apiInfo.AboutResponse} "获取成功"
// @Failure 500 {object} xBase.BaseResponse "服务器内部错误"
// @Router /api/v1/info/about [GET]
func (h *InfoHandler) GetAbout(c *gin.Context) {
	result, err := h.service.infoLogic.GetAbout(c.Request.Context())
	if err != nil {
		_ = c.Error(err)
		return
	}

	xResult.SuccessHasData(c, "获取自我介绍成功", result)
}

// GetBloggerInfo 获取博主信息
//
// @Summary [用户] 获取博主信息
// @Description 获取博主站点资料（站点名字/描述/地址/图片/订阅/邮箱）与个人展示信息（昵称/简介/博客链接/头像），供交换友链复制及「关于我」名士帖展示
// @Tags 站点信息接口
// @Accept json
// @Produce json
// @Success 200 {object} xBase.BaseResponse{data=apiInfo.BloggerResponse} "获取成功"
// @Failure 500 {object} xBase.BaseResponse "服务器内部错误"
// @Router /api/v1/info/blogger [GET]
func (h *InfoHandler) GetBloggerInfo(c *gin.Context) {
	result, err := h.service.infoLogic.GetBloggerInfo(c.Request.Context())
	if err != nil {
		_ = c.Error(err)
		return
	}

	xResult.SuccessHasData(c, "获取博主信息成功", result)
}

// UpdateBloggerInfo 更新博主信息
//
// @Summary [管理] 更新博主信息
// @Description 管理员更新博主站点资料（站点名字/描述/地址/图片/订阅/邮箱）与个人展示信息（昵称/简介/博客链接/头像）
// @Tags 站点信息接口
// @Accept json
// @Produce json
// @Security Bearer
// @Param request body apiInfo.BloggerUpdateRequest true "博主信息更新请求"
// @Success 200 {object} xBase.BaseResponse{data=apiInfo.BloggerResponse} "更新成功"
// @Failure 400 {object} xBase.BaseResponse "请求参数错误"
// @Failure 401 {object} xBase.BaseResponse "未认证"
// @Failure 500 {object} xBase.BaseResponse "服务器内部错误"
// @Router /api/v1/admin/info/blogger [PUT]
func (h *InfoHandler) UpdateBloggerInfo(c *gin.Context) {
	var req apiInfo.BloggerUpdateRequest

	if bindErr := c.ShouldBindJSON(&req); bindErr != nil {
		xValid.HandleValidationError(c, bindErr)
		return
	}

	result, err := h.service.infoLogic.UpdateBloggerInfo(c.Request.Context(), &req)
	if err != nil {
		_ = c.Error(err)
		return
	}

	xResult.SuccessHasData(c, "博主信息更新成功", result)
}

// UpdateAbout 更新自我介绍
//
// @Summary [管理] 更新自我介绍
// @Description 管理员更新 Markdown 格式的自我介绍
// @Tags 站点信息接口
// @Accept json
// @Produce json
// @Security Bearer
// @Param request body apiInfo.AboutUpdateRequest true "自我介绍更新请求"
// @Success 200 {object} xBase.BaseResponse{data=apiInfo.AboutResponse} "更新成功"
// @Failure 400 {object} xBase.BaseResponse "请求参数错误"
// @Failure 401 {object} xBase.BaseResponse "未认证"
// @Failure 500 {object} xBase.BaseResponse "服务器内部错误"
// @Router /api/v1/admin/info/about [PUT]
func (h *InfoHandler) UpdateAbout(c *gin.Context) {
	var req apiInfo.AboutUpdateRequest

	if bindErr := c.ShouldBindJSON(&req); bindErr != nil {
		xValid.HandleValidationError(c, bindErr)
		return
	}

	result, err := h.service.infoLogic.UpdateAbout(c.Request.Context(), &req)
	if err != nil {
		_ = c.Error(err)
		return
	}

	xResult.SuccessHasData(c, "自我介绍更新成功", result)
}
