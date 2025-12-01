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
	"bamboo-main/internal/model/response"
	"bamboo-main/internal/service"

	xResult "github.com/bamboo-services/bamboo-base-go/result"
	xValid "github.com/bamboo-services/bamboo-base-go/validator"
	"github.com/gin-gonic/gin"
)

// LinkColorHandler 友链颜色处理器
type LinkColorHandler struct {
	colorService service.ILinkColorService
}

// NewLinkColorHandler 创建友链颜色处理器
func NewLinkColorHandler() *LinkColorHandler {
	return &LinkColorHandler{
		colorService: service.NewLinkColorService(),
	}
}

// Add 添加友链颜色
// @Summary 添加友链颜色
// @Description 创建新的友链颜色，支持普通颜色和炫彩颜色两种类型
// @Tags 友链颜色管理
// @Accept json
// @Produce json
// @Security Bearer
// @Param request body request.LinkColorAddReq true "添加友链颜色请求"
// @Success 200 {object} response.LinkColorAddResponse "添加成功"
// @Failure 400 {object} map[string]interface{} "请求参数错误"
// @Failure 401 {object} map[string]interface{} "未认证"
// @Failure 500 {object} map[string]interface{} "服务器内部错误"
// @Router /api/v1/admin/colors [post]
func (h *LinkColorHandler) Add(c *gin.Context) {
	var req request.LinkColorAddReq

	if bindErr := c.ShouldBindJSON(&req); bindErr != nil {
		xValid.HandleValidationError(c, bindErr)
		return
	}

	color, err := h.colorService.Add(c, &req)
	if err != nil {
		_ = c.Error(err)
		return
	}

	resp := response.LinkColorAddResponse{LinkColorDetailDTO: *color}
	xResult.SuccessHasData(c, "友链颜色添加成功", resp)
}

// Update 更新友链颜色
// @Summary 更新友链颜色
// @Description 更新指定友链颜色的信息
// @Tags 友链颜色管理
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path int64 true "友链颜色ID"
// @Param request body request.LinkColorUpdateReq true "更新友链颜色请求"
// @Success 200 {object} response.LinkColorUpdateResponse "更新成功"
// @Failure 400 {object} map[string]interface{} "请求参数错误"
// @Failure 401 {object} map[string]interface{} "未认证"
// @Failure 404 {object} map[string]interface{} "友链颜色不存在"
// @Failure 500 {object} map[string]interface{} "服务器内部错误"
// @Router /api/v1/admin/colors/{id} [put]
func (h *LinkColorHandler) Update(c *gin.Context) {
	colorIDStr := c.Param("id")
	var req request.LinkColorUpdateReq

	if bindErr := c.ShouldBindJSON(&req); bindErr != nil {
		xValid.HandleValidationError(c, bindErr)
		return
	}

	color, err := h.colorService.Update(c, colorIDStr, &req)
	if err != nil {
		_ = c.Error(err)
		return
	}

	resp := response.LinkColorUpdateResponse{LinkColorDetailDTO: *color}
	xResult.SuccessHasData(c, "友链颜色更新成功", resp)
}

// UpdateSort 批量更新友链颜色排序
// @Summary 批量更新友链颜色排序
// @Description 按照传入的ID数组顺序重新设置颜色排序
// @Tags 友链颜色管理
// @Accept json
// @Produce json
// @Security Bearer
// @Param request body request.LinkColorSortReq true "颜色排序请求"
// @Success 200 {object} response.LinkColorSortResponse "排序更新成功"
// @Failure 400 {object} map[string]interface{} "请求参数错误"
// @Failure 401 {object} map[string]interface{} "未认证"
// @Failure 404 {object} map[string]interface{} "颜色不存在"
// @Failure 500 {object} map[string]interface{} "服务器内部错误"
// @Router /api/v1/admin/colors/sort [patch]
func (h *LinkColorHandler) UpdateSort(c *gin.Context) {
	var req request.LinkColorSortReq

	if bindErr := c.ShouldBindJSON(&req); bindErr != nil {
		xValid.HandleValidationError(c, bindErr)
		return
	}

	if err := h.colorService.UpdateSort(c, &req); err != nil {
		_ = c.Error(err)
		return
	}

	resp := response.LinkColorSortResponse{
		Message: "颜色排序更新成功",
		Count:   len(req.ColorIDs),
	}
	xResult.SuccessHasData(c, "颜色排序更新成功", resp)
}

// UpdateStatus 更新友链颜色状态
// @Summary 更新友链颜色状态
// @Description 切换指定友链颜色的启用/禁用状态
// @Tags 友链颜色管理
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path int64 true "友链颜色ID"
// @Param request body request.LinkColorStatusReq true "颜色状态请求"
// @Success 200 {object} response.LinkColorStatusResponse "状态更新成功"
// @Failure 400 {object} map[string]interface{} "请求参数错误"
// @Failure 401 {object} map[string]interface{} "未认证"
// @Failure 404 {object} map[string]interface{} "友链颜色不存在"
// @Failure 500 {object} map[string]interface{} "服务器内部错误"
// @Router /api/v1/admin/colors/{id}/status [patch]
func (h *LinkColorHandler) UpdateStatus(c *gin.Context) {
	colorIDStr := c.Param("id")
	var req request.LinkColorStatusReq

	if bindErr := c.ShouldBindJSON(&req); bindErr != nil {
		xValid.HandleValidationError(c, bindErr)
		return
	}

	if err := h.colorService.UpdateStatus(c, colorIDStr, &req); err != nil {
		_ = c.Error(err)
		return
	}

	statusText := "禁用"
	if req.Status {
		statusText = "启用"
	}
	resp := response.LinkColorStatusResponse{
		Message: "颜色状态更新成功",
		Status:  req.Status,
	}
	xResult.SuccessHasData(c, "颜色已"+statusText, resp)
}

// Delete 删除友链颜色
// @Summary 删除友链颜色
// @Description 删除指定的友链颜色，支持强制删除模式
// @Tags 友链颜色管理
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path int64 true "友链颜色ID"
// @Param force query bool false "是否强制删除（默认false）"
// @Success 200 {object} response.LinkColorDeleteResponse "删除成功"
// @Failure 400 {object} map[string]interface{} "请求参数错误"
// @Failure 401 {object} map[string]interface{} "未认证"
// @Failure 404 {object} map[string]interface{} "友链颜色不存在"
// @Failure 409 {object} response.LinkColorDeleteConflictResponse "存在关联数据冲突"
// @Failure 500 {object} map[string]interface{} "服务器内部错误"
// @Router /api/v1/admin/colors/{id} [delete]
func (h *LinkColorHandler) Delete(c *gin.Context) {
	colorIDStr := c.Param("id")

	var req request.LinkColorDeleteReq
	req.Force = c.Query("force") == "true"

	if _, err := h.colorService.Delete(c, colorIDStr, &req); err != nil {
		_ = c.Error(err)
		return
	}

	resp := response.LinkColorDeleteResponse{
		Message: "友链颜色删除成功",
	}
	xResult.SuccessHasData(c, "友链颜色删除成功", resp)
}

// Get 获取友链颜色详情
// @Summary 获取友链颜色详情
// @Description 根据ID获取指定友链颜色的详细信息
// @Tags 友链颜色管理
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path int64 true "友链颜色ID"
// @Success 200 {object} response.LinkColorDetailResponse "获取成功"
// @Failure 400 {object} map[string]interface{} "请求参数错误"
// @Failure 401 {object} map[string]interface{} "未认证"
// @Failure 404 {object} map[string]interface{} "友链颜色不存在"
// @Failure 500 {object} map[string]interface{} "服务器内部错误"
// @Router /api/v1/admin/colors/{id} [get]
func (h *LinkColorHandler) Get(c *gin.Context) {
	colorIDStr := c.Param("id")

	color, err := h.colorService.Get(c, colorIDStr)
	if err != nil {
		_ = c.Error(err)
		return
	}

	resp := response.LinkColorDetailResponse{LinkColorDetailDTO: *color}
	xResult.SuccessHasData(c, "获取友链颜色详情成功", resp)
}

// GetList 获取友链颜色列表
// @Summary 获取友链颜色列表
// @Description 获取友链颜色列表（不分页），支持过滤和排序
// @Tags 友链颜色管理
// @Accept json
// @Produce json
// @Security Bearer
// @Param status query int false "状态过滤（0=禁用，1=启用）"
// @Param type query int false "类型过滤（0=普通，1=炫彩）"
// @Param name query string false "名称模糊搜索"
// @Param only_enabled query bool false "仅查询启用的颜色"
// @Param order_by query string false "排序字段（name, sort_order, created_at）"
// @Param order query string false "排序方向（asc, desc）"
// @Success 200 {object} []dto.LinkColorListDTO "获取成功"
// @Failure 400 {object} map[string]interface{} "请求参数错误"
// @Failure 401 {object} map[string]interface{} "未认证"
// @Failure 500 {object} map[string]interface{} "服务器内部错误"
// @Router /api/v1/admin/colors/all [get]
func (h *LinkColorHandler) GetList(c *gin.Context) {
	var req request.LinkColorListReq

	if bindErr := c.ShouldBindQuery(&req); bindErr != nil {
		xValid.HandleValidationError(c, bindErr)
		return
	}

	colors, err := h.colorService.GetList(c, &req)
	if err != nil {
		_ = c.Error(err)
		return
	}

	xResult.SuccessHasData(c, "获取友链颜色列表成功", colors)
}

// GetPage 获取友链颜色分页列表
// @Summary 获取友链颜色分页列表
// @Description 分页获取友链颜色列表，支持过滤和排序
// @Tags 友链颜色管理
// @Accept json
// @Produce json
// @Security Bearer
// @Param page query int false "页码（默认1）"
// @Param page_size query int false "每页数量（默认10，最大100）"
// @Param status query int false "状态过滤（0=禁用，1=启用）"
// @Param type query int false "类型过滤（0=普通，1=炫彩）"
// @Param name query string false "名称模糊搜索"
// @Param order_by query string false "排序字段（name, sort_order, created_at）"
// @Param order query string false "排序方向（asc, desc）"
// @Success 200 {object} response.LinkColorPageResponse "获取成功"
// @Failure 400 {object} map[string]interface{} "请求参数错误"
// @Failure 401 {object} map[string]interface{} "未认证"
// @Failure 500 {object} map[string]interface{} "服务器内部错误"
// @Router /api/v1/admin/colors [get]
func (h *LinkColorHandler) GetPage(c *gin.Context) {
	var req request.LinkColorPageReq

	if bindErr := c.ShouldBindQuery(&req); bindErr != nil {
		xValid.HandleValidationError(c, bindErr)
		return
	}

	result, err := h.colorService.GetPage(c, &req)
	if err != nil {
		_ = c.Error(err)
		return
	}

	resp := response.LinkColorPageResponse{PaginationResponse: *result}
	xResult.SuccessHasData(c, "获取友链颜色分页列表成功", resp)
}
