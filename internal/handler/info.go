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
// @Description 获取站点名称、主页介绍等公开信息
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
// @Description 管理员更新站点名称、主页介绍
// @Tags 站点信息接口
// @Accept json
// @Produce json
// @Security Bearer
// @Param request body apiInfo.SiteUpdateRequest true "站点信息更新请求"
// @Success 200 {object} xBase.BaseResponse{data=apiInfo.SiteResponse} "更新成功"
// @Failure 400 {object} xBase.BaseResponse "请求参数错误"
// @Failure 401 {object} xBase.BaseResponse "未认证"
// @Failure 500 {object} xBase.BaseResponse "服务器内部错误"
// @Router /api/v1/info/admin/site [PUT]
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

// GetArchiveInfo 获取站点档案
//
// @Summary [用户] 获取站点档案
// @Description 获取站点描述与自我介绍（均 Markdown），供 about/me 展示
// @Tags 站点信息接口
// @Accept json
// @Produce json
// @Success 200 {object} xBase.BaseResponse{data=apiInfo.ArchiveResponse} "获取成功"
// @Failure 500 {object} xBase.BaseResponse "服务器内部错误"
// @Router /api/v1/info/archive [GET]
func (h *InfoHandler) GetArchiveInfo(c *gin.Context) {
	result, err := h.service.infoLogic.GetArchiveInfo(c.Request.Context())
	if err != nil {
		_ = c.Error(err)
		return
	}

	xResult.SuccessHasData(c, "获取站点档案成功", result)
}

// GetBuiltinInvalidGroup 获取内置「已失效」分组配置
//
// @Summary [用户] 获取内置已失效分组配置
// @Description 返回内置「已失效」分组的名称与描述（经 bm_system 热修改）
// @Tags 站点信息接口
// @Accept json
// @Produce json
// @Success 200 {object} xBase.BaseResponse{data=apiInfo.BuiltinInvalidGroupResponse} "获取成功"
// @Failure 500 {object} xBase.BaseResponse "服务器内部错误"
// @Router /api/v1/info/builtin-invalid-group [GET]
func (h *InfoHandler) GetBuiltinInvalidGroup(c *gin.Context) {
	group, err := h.service.infoLogic.GetBuiltinInvalidGroup(c.Request.Context())
	if err != nil {
		_ = c.Error(err)
		return
	}

	resp := apiInfo.BuiltinInvalidGroupResponse{LinkGroup: *group}
	xResult.SuccessHasData(c, "获取内置已失效分组配置成功", resp)
}

// GetApplySiteInfo 获取申请站点展示
//
// @Summary [用户] 获取申请站点展示
// @Description 获取博主站点资料（站点名字/描述/地址/图片/订阅/邮箱），供 operate/apply 申请页交换友链时复制添加
// @Tags 站点信息接口
// @Accept json
// @Produce json
// @Success 200 {object} xBase.BaseResponse{data=apiInfo.ApplySiteResponse} "获取成功"
// @Failure 500 {object} xBase.BaseResponse "服务器内部错误"
// @Router /api/v1/info/apply-site [GET]
func (h *InfoHandler) GetApplySiteInfo(c *gin.Context) {
	result, err := h.service.infoLogic.GetApplySiteInfo(c.Request.Context())
	if err != nil {
		_ = c.Error(err)
		return
	}

	xResult.SuccessHasData(c, "获取申请站点展示成功", result)
}

// GetBloggerInfo 获取博主信息
//
// @Summary [用户] 获取博主信息
// @Description 获取博主个人展示信息（昵称/简介/博客链接/头像），供「关于我」名士帖展示
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

// UpdateApplySiteInfo 更新申请站点展示
//
// @Summary [管理] 更新申请站点展示
// @Description 管理员更新博主站点资料（站点名字/描述/地址/图片/订阅/邮箱），供 operate/apply 申请页交换友链时复制
// @Tags 站点信息接口
// @Accept json
// @Produce json
// @Security Bearer
// @Param request body apiInfo.ApplySiteUpdateRequest true "申请站点展示更新请求"
// @Success 200 {object} xBase.BaseResponse{data=apiInfo.ApplySiteResponse} "更新成功"
// @Failure 400 {object} xBase.BaseResponse "请求参数错误"
// @Failure 401 {object} xBase.BaseResponse "未认证"
// @Failure 500 {object} xBase.BaseResponse "服务器内部错误"
// @Router /api/v1/info/admin/apply-site [PUT]
func (h *InfoHandler) UpdateApplySiteInfo(c *gin.Context) {
	var req apiInfo.ApplySiteUpdateRequest

	if bindErr := c.ShouldBindJSON(&req); bindErr != nil {
		xValid.HandleValidationError(c, bindErr)
		return
	}

	result, err := h.service.infoLogic.UpdateApplySiteInfo(c.Request.Context(), &req)
	if err != nil {
		_ = c.Error(err)
		return
	}

	xResult.SuccessHasData(c, "申请站点展示更新成功", result)
}

// UpdateBloggerInfo 更新博主信息
//
// @Summary [管理] 更新博主信息
// @Description 管理员更新博主个人展示信息（昵称/简介/博客链接/头像），供「关于我」名士帖展示
// @Tags 站点信息接口
// @Accept json
// @Produce json
// @Security Bearer
// @Param request body apiInfo.BloggerUpdateRequest true "博主信息更新请求"
// @Success 200 {object} xBase.BaseResponse{data=apiInfo.BloggerResponse} "更新成功"
// @Failure 400 {object} xBase.BaseResponse "请求参数错误"
// @Failure 401 {object} xBase.BaseResponse "未认证"
// @Failure 500 {object} xBase.BaseResponse "服务器内部错误"
// @Router /api/v1/info/admin/blogger [PUT]
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

// UpdateArchiveInfo 更新站点档案
//
// @Summary [管理] 更新站点档案
// @Description 管理员更新站点描述与自我介绍（均 Markdown），一次保存
// @Tags 站点信息接口
// @Accept json
// @Produce json
// @Security Bearer
// @Param request body apiInfo.ArchiveUpdateRequest true "站点档案更新请求"
// @Success 200 {object} xBase.BaseResponse{data=apiInfo.ArchiveResponse} "更新成功"
// @Failure 400 {object} xBase.BaseResponse "请求参数错误"
// @Failure 401 {object} xBase.BaseResponse "未认证"
// @Failure 500 {object} xBase.BaseResponse "服务器内部错误"
// @Router /api/v1/info/admin/archive [PUT]
func (h *InfoHandler) UpdateArchiveInfo(c *gin.Context) {
	var req apiInfo.ArchiveUpdateRequest

	if bindErr := c.ShouldBindJSON(&req); bindErr != nil {
		xValid.HandleValidationError(c, bindErr)
		return
	}

	result, err := h.service.infoLogic.UpdateArchiveInfo(c.Request.Context(), &req)
	if err != nil {
		_ = c.Error(err)
		return
	}

	xResult.SuccessHasData(c, "站点档案更新成功", result)
}

// UpdateBuiltinInvalidGroup 更新内置「已失效」分组配置
//
// @Summary [管理] 更新内置已失效分组配置
// @Description 更新内置「已失效」分组的名称与描述（写入 bm_system，即时生效）
// @Tags 站点信息接口
// @Accept json
// @Produce json
// @Security Bearer
// @Param request body apiInfo.BuiltinInvalidGroupUpdateRequest true "更新内置已失效分组配置请求"
// @Success 200 {object} xBase.BaseResponse{data=apiInfo.BuiltinInvalidGroupResponse} "更新成功"
// @Failure 400 {object} xBase.BaseResponse "请求参数错误"
// @Failure 401 {object} xBase.BaseResponse "未认证"
// @Failure 500 {object} xBase.BaseResponse "服务器内部错误"
// @Router /api/v1/info/admin/builtin-invalid-group [PUT]
func (h *InfoHandler) UpdateBuiltinInvalidGroup(c *gin.Context) {
	var req apiInfo.BuiltinInvalidGroupUpdateRequest

	if bindErr := c.ShouldBindJSON(&req); bindErr != nil {
		xValid.HandleValidationError(c, bindErr)
		return
	}

	group, err := h.service.infoLogic.UpdateBuiltinInvalidGroup(c.Request.Context(), &req)
	if err != nil {
		_ = c.Error(err)
		return
	}

	resp := apiInfo.BuiltinInvalidGroupResponse{LinkGroup: *group}
	xResult.SuccessHasData(c, "内置已失效分组配置更新成功", resp)
}

// GetColorMode 获取高级配色模式
//
// @Summary [用户] 获取高级配色模式
// @Description 获取站点当前高级配色模式（normal=普通, premium=高级），决定颜色选择器可见范围
// @Tags 站点信息接口
// @Accept json
// @Produce json
// @Success 200 {object} xBase.BaseResponse{data=apiInfo.ColorModeResponse} "获取成功"
// @Failure 500 {object} xBase.BaseResponse "服务器内部错误"
// @Router /api/v1/info/color-mode [GET]
func (h *InfoHandler) GetColorMode(c *gin.Context) {
	result, err := h.service.infoLogic.GetColorMode(c.Request.Context())
	if err != nil {
		_ = c.Error(err)
		return
	}

	xResult.SuccessHasData(c, "获取高级配色模式成功", result)
}

// UpdateColorMode 更新高级配色模式
//
// @Summary [管理] 更新高级配色模式
// @Description 管理员切换站点高级配色模式（normal=普通, premium=高级），开启后高级配色可选并渐变渲染
// @Tags 站点信息接口
// @Accept json
// @Produce json
// @Security Bearer
// @Param request body apiInfo.ColorModeUpdateRequest true "高级配色模式更新请求"
// @Success 200 {object} xBase.BaseResponse{data=apiInfo.ColorModeResponse} "更新成功"
// @Failure 400 {object} xBase.BaseResponse "请求参数错误"
// @Failure 401 {object} xBase.BaseResponse "未认证"
// @Failure 500 {object} xBase.BaseResponse "服务器内部错误"
// @Router /api/v1/info/admin/color-mode [PUT]
func (h *InfoHandler) UpdateColorMode(c *gin.Context) {
	var req apiInfo.ColorModeUpdateRequest

	if bindErr := c.ShouldBindJSON(&req); bindErr != nil {
		xValid.HandleValidationError(c, bindErr)
		return
	}

	result, err := h.service.infoLogic.UpdateColorMode(c.Request.Context(), &req)
	if err != nil {
		_ = c.Error(err)
		return
	}

	xResult.SuccessHasData(c, "高级配色模式更新成功", result)
}
