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

package handler

import (
	apiInfo "github.com/bamboo-services/bamboo-main/api/info"

	xResult "github.com/bamboo-services/bamboo-base-go/result"
	xValid "github.com/bamboo-services/bamboo-base-go/validator"
	"github.com/gin-gonic/gin"
)

// GetSiteInfo 获取站点信息
// @Summary 获取站点信息
// @Description 获取站点名称、描述、主页介绍等公开信息
// @Tags 站点信息
// @Accept json
// @Produce json
// @Success 200 {object} apiInfo.SiteResponse "获取成功"
// @Failure 500 {object} map[string]interface{} "服务器内部错误"
// @Router /api/v1/info/site [get]
func (h *InfoHandler) GetSiteInfo(c *gin.Context) {
	result, err := h.service.infoLogic.GetSiteInfo(c)
	if err != nil {
		_ = c.Error(err)
		return
	}

	resp := apiInfo.SiteResponse{SiteInfoDTO: *result}
	xResult.SuccessHasData(c, "获取站点信息成功", resp)
}

// UpdateSiteInfo 更新站点信息
// @Summary 更新站点信息
// @Description 管理员更新站点名称、描述、主页介绍
// @Tags 站点信息管理
// @Accept json
// @Produce json
// @Security Bearer
// @Param request body apiInfo.SiteUpdateRequest true "站点信息更新请求"
// @Success 200 {object} apiInfo.SiteResponse "更新成功"
// @Failure 400 {object} map[string]interface{} "请求参数错误"
// @Failure 401 {object} map[string]interface{} "未认证"
// @Failure 500 {object} map[string]interface{} "服务器内部错误"
// @Router /api/v1/admin/info/site [put]
func (h *InfoHandler) UpdateSiteInfo(c *gin.Context) {
	var req apiInfo.SiteUpdateRequest

	if bindErr := c.ShouldBindJSON(&req); bindErr != nil {
		xValid.HandleValidationError(c, bindErr)
		return
	}

	result, err := h.service.infoLogic.UpdateSiteInfo(c, &req)
	if err != nil {
		_ = c.Error(err)
		return
	}

	resp := apiInfo.SiteResponse{SiteInfoDTO: *result}
	xResult.SuccessHasData(c, "站点信息更新成功", resp)
}

// GetAbout 获取自我介绍
// @Summary 获取自我介绍
// @Description 获取 Markdown 格式的自我介绍内容
// @Tags 站点信息
// @Accept json
// @Produce json
// @Success 200 {object} apiInfo.AboutResponse "获取成功"
// @Failure 500 {object} map[string]interface{} "服务器内部错误"
// @Router /api/v1/info/about [get]
func (h *InfoHandler) GetAbout(c *gin.Context) {
	result, err := h.service.infoLogic.GetAbout(c)
	if err != nil {
		_ = c.Error(err)
		return
	}

	resp := apiInfo.AboutResponse{AboutDTO: *result}
	xResult.SuccessHasData(c, "获取自我介绍成功", resp)
}

// UpdateAbout 更新自我介绍
// @Summary 更新自我介绍
// @Description 管理员更新 Markdown 格式的自我介绍
// @Tags 站点信息管理
// @Accept json
// @Produce json
// @Security Bearer
// @Param request body apiInfo.AboutUpdateRequest true "自我介绍更新请求"
// @Success 200 {object} apiInfo.AboutResponse "更新成功"
// @Failure 400 {object} map[string]interface{} "请求参数错误"
// @Failure 401 {object} map[string]interface{} "未认证"
// @Failure 500 {object} map[string]interface{} "服务器内部错误"
// @Router /api/v1/admin/info/about [put]
func (h *InfoHandler) UpdateAbout(c *gin.Context) {
	var req apiInfo.AboutUpdateRequest

	if bindErr := c.ShouldBindJSON(&req); bindErr != nil {
		xValid.HandleValidationError(c, bindErr)
		return
	}

	result, err := h.service.infoLogic.UpdateAbout(c, &req)
	if err != nil {
		_ = c.Error(err)
		return
	}

	resp := apiInfo.AboutResponse{AboutDTO: *result}
	xResult.SuccessHasData(c, "自我介绍更新成功", resp)
}
