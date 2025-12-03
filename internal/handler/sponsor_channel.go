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
	"bamboo-main/internal/model/request"
	"bamboo-main/internal/service"

	xResult "github.com/bamboo-services/bamboo-base-go/result"
	xValid "github.com/bamboo-services/bamboo-base-go/validator"
	"github.com/gin-gonic/gin"
)

// SponsorChannelHandler 赞助渠道处理器
type SponsorChannelHandler struct {
	channelService service.ISponsorChannelService
}

// NewSponsorChannelHandler 创建赞助渠道处理器
func NewSponsorChannelHandler() *SponsorChannelHandler {
	return &SponsorChannelHandler{
		channelService: service.NewSponsorChannelService(),
	}
}

// Add 添加赞助渠道
// @Summary 添加赞助渠道
// @Description 创建新的赞助渠道
// @Tags 赞助渠道管理
// @Accept json
// @Produce json
// @Security Bearer
// @Param request body request.SponsorChannelAddReq true "添加赞助渠道请求"
// @Success 200 {object} dto.SponsorChannelDetailDTO "添加成功"
// @Failure 400 {object} map[string]interface{} "请求参数错误"
// @Failure 401 {object} map[string]interface{} "未认证"
// @Failure 500 {object} map[string]interface{} "服务器内部错误"
// @Router /api/v1/admin/sponsors/channels [post]
func (h *SponsorChannelHandler) Add(c *gin.Context) {
	var req request.SponsorChannelAddReq

	// 绑定请求数据
	bindErr := c.ShouldBindJSON(&req)
	if bindErr != nil {
		xValid.HandleValidationError(c, bindErr)
		return
	}

	// 调用服务层
	channel, err := h.channelService.Add(c, &req)
	if err != nil {
		_ = c.Error(err)
		return
	}

	// 返回成功响应
	xResult.SuccessHasData(c, "赞助渠道添加成功", channel)
}

// Update 更新赞助渠道
// @Summary 更新赞助渠道
// @Description 更新指定赞助渠道的信息
// @Tags 赞助渠道管理
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path int64 true "赞助渠道ID"
// @Param request body request.SponsorChannelUpdateReq true "更新赞助渠道请求"
// @Success 200 {object} dto.SponsorChannelDetailDTO "更新成功"
// @Failure 400 {object} map[string]interface{} "请求参数错误"
// @Failure 401 {object} map[string]interface{} "未认证"
// @Failure 404 {object} map[string]interface{} "赞助渠道不存在"
// @Failure 500 {object} map[string]interface{} "服务器内部错误"
// @Router /api/v1/admin/sponsors/channels/{id} [put]
func (h *SponsorChannelHandler) Update(c *gin.Context) {
	channelIDStr := c.Param("id")
	var req request.SponsorChannelUpdateReq

	// 绑定请求数据
	bindErr := c.ShouldBindJSON(&req)
	if bindErr != nil {
		xValid.HandleValidationError(c, bindErr)
		return
	}

	// 调用服务层
	channel, err := h.channelService.Update(c, channelIDStr, &req)
	if err != nil {
		_ = c.Error(err)
		return
	}

	// 返回成功响应
	xResult.SuccessHasData(c, "赞助渠道更新成功", channel)
}

// UpdateStatus 更新赞助渠道状态
// @Summary 更新赞助渠道状态
// @Description 切换指定赞助渠道的启用/禁用状态
// @Tags 赞助渠道管理
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path int64 true "赞助渠道ID"
// @Param request body request.SponsorChannelStatusReq true "渠道状态请求"
// @Success 200 {object} map[string]interface{} "状态更新成功"
// @Failure 400 {object} map[string]interface{} "请求参数错误"
// @Failure 401 {object} map[string]interface{} "未认证"
// @Failure 404 {object} map[string]interface{} "赞助渠道不存在"
// @Failure 500 {object} map[string]interface{} "服务器内部错误"
// @Router /api/v1/admin/sponsors/channels/{id}/status [patch]
func (h *SponsorChannelHandler) UpdateStatus(c *gin.Context) {
	channelIDStr := c.Param("id")
	var req request.SponsorChannelStatusReq

	// 绑定请求数据
	bindErr := c.ShouldBindJSON(&req)
	if bindErr != nil {
		xValid.HandleValidationError(c, bindErr)
		return
	}

	// 调用服务层
	status, err := h.channelService.UpdateStatus(c, channelIDStr, &req)
	if err != nil {
		_ = c.Error(err)
		return
	}

	// 返回成功响应
	statusText := "禁用"
	if status {
		statusText = "启用"
	}
	xResult.SuccessHasData(c, "渠道已"+statusText, gin.H{
		"status": status,
	})
}

// Delete 删除赞助渠道
// @Summary 删除赞助渠道
// @Description 删除指定的赞助渠道
// @Tags 赞助渠道管理
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path int64 true "赞助渠道ID"
// @Success 200 {object} map[string]interface{} "删除成功"
// @Failure 400 {object} map[string]interface{} "请求参数错误"
// @Failure 401 {object} map[string]interface{} "未认证"
// @Failure 404 {object} map[string]interface{} "赞助渠道不存在"
// @Failure 500 {object} map[string]interface{} "服务器内部错误"
// @Router /api/v1/admin/sponsors/channels/{id} [delete]
func (h *SponsorChannelHandler) Delete(c *gin.Context) {
	channelIDStr := c.Param("id")

	// 调用服务层
	err := h.channelService.Delete(c, channelIDStr)
	if err != nil {
		_ = c.Error(err)
		return
	}

	// 返回成功响应
	xResult.Success(c, "赞助渠道删除成功")
}

// Get 获取赞助渠道详情
// @Summary 获取赞助渠道详情
// @Description 根据ID获取指定赞助渠道的详细信息
// @Tags 赞助渠道管理
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path int64 true "赞助渠道ID"
// @Success 200 {object} dto.SponsorChannelDetailDTO "获取成功"
// @Failure 400 {object} map[string]interface{} "请求参数错误"
// @Failure 401 {object} map[string]interface{} "未认证"
// @Failure 404 {object} map[string]interface{} "赞助渠道不存在"
// @Failure 500 {object} map[string]interface{} "服务器内部错误"
// @Router /api/v1/admin/sponsors/channels/{id} [get]
func (h *SponsorChannelHandler) Get(c *gin.Context) {
	channelIDStr := c.Param("id")

	// 调用服务层
	channel, err := h.channelService.Get(c, channelIDStr)
	if err != nil {
		_ = c.Error(err)
		return
	}

	// 返回成功响应
	xResult.SuccessHasData(c, "获取赞助渠道详情成功", channel)
}

// GetList 获取赞助渠道列表
// @Summary 获取赞助渠道列表
// @Description 获取赞助渠道列表（不分页），支持过滤和排序
// @Tags 赞助渠道管理
// @Accept json
// @Produce json
// @Security Bearer
// @Param status query bool false "状态过滤（true=启用，false=禁用）"
// @Param name query string false "名称模糊搜索"
// @Param only_enabled query bool false "仅查询启用的渠道"
// @Param order_by query string false "排序字段（name, sort_order, created_at）"
// @Param order query string false "排序方向（asc, desc）"
// @Success 200 {array} dto.SponsorChannelListDTO "获取成功"
// @Failure 400 {object} map[string]interface{} "请求参数错误"
// @Failure 401 {object} map[string]interface{} "未认证"
// @Failure 500 {object} map[string]interface{} "服务器内部错误"
// @Router /api/v1/admin/sponsors/channels/all [get]
func (h *SponsorChannelHandler) GetList(c *gin.Context) {
	var req request.SponsorChannelListReq

	// 绑定查询参数
	bindErr := c.ShouldBindQuery(&req)
	if bindErr != nil {
		xValid.HandleValidationError(c, bindErr)
		return
	}

	// 调用服务层
	channels, err := h.channelService.GetList(c, &req)
	if err != nil {
		_ = c.Error(err)
		return
	}

	// 返回成功响应
	xResult.SuccessHasData(c, "获取赞助渠道列表成功", channels)
}

// GetPage 获取赞助渠道分页列表
// @Summary 获取赞助渠道分页列表
// @Description 分页获取赞助渠道列表，支持过滤和排序
// @Tags 赞助渠道管理
// @Accept json
// @Produce json
// @Security Bearer
// @Param page query int false "页码（默认1）"
// @Param page_size query int false "每页数量（默认10，最大100）"
// @Param status query bool false "状态过滤（true=启用，false=禁用）"
// @Param name query string false "名称模糊搜索"
// @Param order_by query string false "排序字段（name, sort_order, created_at）"
// @Param order query string false "排序方向（asc, desc）"
// @Success 200 {object} response.SponsorChannelPageResponse "获取成功"
// @Failure 400 {object} map[string]interface{} "请求参数错误"
// @Failure 401 {object} map[string]interface{} "未认证"
// @Failure 500 {object} map[string]interface{} "服务器内部错误"
// @Router /api/v1/admin/sponsors/channels [get]
func (h *SponsorChannelHandler) GetPage(c *gin.Context) {
	var req request.SponsorChannelPageReq

	// 绑定查询参数
	bindErr := c.ShouldBindQuery(&req)
	if bindErr != nil {
		xValid.HandleValidationError(c, bindErr)
		return
	}

	// 调用服务层
	result, err := h.channelService.GetPage(c, &req)
	if err != nil {
		_ = c.Error(err)
		return
	}

	// 返回成功响应
	xResult.SuccessHasData(c, "获取赞助渠道分页列表成功", result)
}

// GetPublicList 获取公开的赞助渠道列表
// @Summary 获取公开的赞助渠道列表
// @Description 获取所有启用状态的赞助渠道（不分页，公开接口）
// @Tags 赞助渠道
// @Accept json
// @Produce json
// @Success 200 {array} dto.SponsorChannelListDTO "获取成功"
// @Failure 500 {object} map[string]interface{} "服务器内部错误"
// @Router /api/v1/sponsors/channels [get]
func (h *SponsorChannelHandler) GetPublicList(c *gin.Context) {
	// 调用服务层
	channels, err := h.channelService.GetPublicList(c)
	if err != nil {
		_ = c.Error(err)
		return
	}

	// 返回成功响应
	xResult.SuccessHasData(c, "获取公开渠道列表成功", channels)
}
