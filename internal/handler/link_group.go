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
	apiLinkGroup "github.com/bamboo-services/bamboo-main/api/linkgroup"
	logic "github.com/bamboo-services/bamboo-main/internal/logic"

	xResult "github.com/bamboo-services/bamboo-base-go/result"
	xValid "github.com/bamboo-services/bamboo-base-go/validator"
	"github.com/gin-gonic/gin"
)

// LinkGroupHandler 友链分组处理器
type LinkGroupHandler struct {
	groupLogic *logic.LinkGroupLogic
}

// NewLinkGroupHandler 创建友链分组处理器
func NewLinkGroupHandler() *LinkGroupHandler {
	return &LinkGroupHandler{
		groupLogic: logic.NewLinkGroupLogic(),
	}
}

// Add 添加友链分组
// @Summary 添加友链分组
// @Description 创建新的友链分组
// @Tags 友链分组管理
// @Accept json
// @Produce json
// @Security Bearer
// @Param request body apiLinkGroup.AddRequest true "添加友链分组请求"
// @Success 200 {object} apiLinkGroup.AddResponse "添加成功"
// @Failure 400 {object} map[string]interface{} "请求参数错误"
// @Failure 401 {object} map[string]interface{} "未认证"
// @Failure 500 {object} map[string]interface{} "服务器内部错误"
// @Router /api/v1/admin/groups [post]
func (h *LinkGroupHandler) Add(c *gin.Context) {
	var req apiLinkGroup.AddRequest

	// 绑定请求数据
	bindErr := c.ShouldBindJSON(&req)
	if bindErr != nil {
		xValid.HandleValidationError(c, bindErr)
		return
	}

	// 调用服务层
	group, err := h.groupLogic.Add(c, &req)
	if err != nil {
		_ = c.Error(err)
		return
	}

	// 返回成功响应
	resp := apiLinkGroup.AddResponse{LinkGroupDetailDTO: *group}
	xResult.SuccessHasData(c, "友链分组添加成功", resp)
}

// Update 更新友链分组
// @Summary 更新友链分组
// @Description 更新指定友链分组的名称和描述
// @Tags 友链分组管理
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path int64 true "友链分组ID"
// @Param request body apiLinkGroup.UpdateRequest true "更新友链分组请求"
// @Success 200 {object} apiLinkGroup.UpdateResponse "更新成功"
// @Failure 400 {object} map[string]interface{} "请求参数错误"
// @Failure 401 {object} map[string]interface{} "未认证"
// @Failure 404 {object} map[string]interface{} "友链分组不存在"
// @Failure 500 {object} map[string]interface{} "服务器内部错误"
// @Router /api/v1/admin/groups/{id} [put]
func (h *LinkGroupHandler) Update(c *gin.Context) {
	groupIDStr := c.Param("id")
	var req apiLinkGroup.UpdateRequest

	// 绑定请求数据
	bindErr := c.ShouldBindJSON(&req)
	if bindErr != nil {
		xValid.HandleValidationError(c, bindErr)
		return
	}

	// 调用服务层
	group, err := h.groupLogic.Update(c, groupIDStr, &req)
	if err != nil {
		_ = c.Error(err)
		return
	}

	// 返回成功响应
	resp := apiLinkGroup.UpdateResponse{LinkGroupDetailDTO: *group}
	xResult.SuccessHasData(c, "友链分组更新成功", resp)
}

// UpdateSort 批量更新友链分组排序
// @Summary 批量更新友链分组排序
// @Description 按照传入的UUID数组顺序重新设置分组排序
// @Tags 友链分组管理
// @Accept json
// @Produce json
// @Security Bearer
// @Param request body apiLinkGroup.SortRequest true "分组排序请求"
// @Success 200 {object} apiLinkGroup.SortResponse "排序更新成功"
// @Failure 400 {object} map[string]interface{} "请求参数错误"
// @Failure 401 {object} map[string]interface{} "未认证"
// @Failure 404 {object} map[string]interface{} "分组不存在"
// @Failure 500 {object} map[string]interface{} "服务器内部错误"
// @Router /api/v1/admin/groups/sort [patch]
func (h *LinkGroupHandler) UpdateSort(c *gin.Context) {
	var req apiLinkGroup.SortRequest

	// 绑定请求数据
	bindErr := c.ShouldBindJSON(&req)
	if bindErr != nil {
		xValid.HandleValidationError(c, bindErr)
		return
	}

	// 调用服务层
	err := h.groupLogic.UpdateSort(c, &req)
	if err != nil {
		_ = c.Error(err)
		return
	}

	// 返回成功响应
	resp := apiLinkGroup.SortResponse{
		Message: "分组排序更新成功",
		Count:   len(req.GroupIDs),
	}
	xResult.SuccessHasData(c, "分组排序更新成功", resp)
}

// UpdateStatus 更新友链分组状态
// @Summary 更新友链分组状态
// @Description 切换指定友链分组的启用/禁用状态
// @Tags 友链分组管理
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path int64 true "友链分组ID"
// @Param request body apiLinkGroup.StatusRequest true "分组状态请求"
// @Success 200 {object} apiLinkGroup.StatusResponse "状态更新成功"
// @Failure 400 {object} map[string]interface{} "请求参数错误"
// @Failure 401 {object} map[string]interface{} "未认证"
// @Failure 404 {object} map[string]interface{} "友链分组不存在"
// @Failure 500 {object} map[string]interface{} "服务器内部错误"
// @Router /api/v1/admin/groups/{id}/status [patch]
func (h *LinkGroupHandler) UpdateStatus(c *gin.Context) {
	groupIDStr := c.Param("id")
	var req apiLinkGroup.StatusRequest

	// 绑定请求数据
	bindErr := c.ShouldBindJSON(&req)
	if bindErr != nil {
		xValid.HandleValidationError(c, bindErr)
		return
	}

	// 调用服务层
	err := h.groupLogic.UpdateStatus(c, groupIDStr, &req)
	if err != nil {
		_ = c.Error(err)
		return
	}

	// 返回成功响应
	statusText := "禁用"
	if req.Status {
		statusText = "启用"
	}
	resp := apiLinkGroup.StatusResponse{
		Message: "分组状态更新成功",
		Status:  req.Status,
	}
	xResult.SuccessHasData(c, "分组已"+statusText, resp)
}

// Delete 删除友链分组
// @Summary 删除友链分组
// @Description 删除指定的友链分组，支持强制删除模式
// @Tags 友链分组管理
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path int64 true "友链分组ID"
// @Param force query bool false "是否强制删除（默认false）"
// @Success 200 {object} apiLinkGroup.DeleteResponse "删除成功"
// @Failure 400 {object} map[string]interface{} "请求参数错误"
// @Failure 401 {object} map[string]interface{} "未认证"
// @Failure 404 {object} map[string]interface{} "友链分组不存在"
// @Failure 409 {object} apiLinkGroup.DeleteConflictResponse "存在关联数据冲突"
// @Failure 500 {object} map[string]interface{} "服务器内部错误"
// @Router /api/v1/admin/groups/{id} [delete]
func (h *LinkGroupHandler) Delete(c *gin.Context) {
	groupIDStr := c.Param("id")

	// 获取force参数
	var req apiLinkGroup.DeleteRequest
	req.Force = c.Query("force") == "true"

	// 调用服务层
	_, err := h.groupLogic.Delete(c, groupIDStr, &req)
	if err != nil {
		_ = c.Error(err)
		return
	}

	// 返回成功响应
	resp := apiLinkGroup.DeleteResponse{
		Message: "友链分组删除成功",
	}
	xResult.SuccessHasData(c, "友链分组删除成功", resp)
}

// Get 获取友链分组详情
// @Summary 获取友链分组详情
// @Description 根据ID获取指定友链分组的详细信息
// @Tags 友链分组管理
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path int64 true "友链分组ID"
// @Success 200 {object} apiLinkGroup.DetailResponse "获取成功"
// @Failure 400 {object} map[string]interface{} "请求参数错误"
// @Failure 401 {object} map[string]interface{} "未认证"
// @Failure 404 {object} map[string]interface{} "友链分组不存在"
// @Failure 500 {object} map[string]interface{} "服务器内部错误"
// @Router /api/v1/admin/groups/{id} [get]
func (h *LinkGroupHandler) Get(c *gin.Context) {
	groupIDStr := c.Param("id")

	// 调用服务层
	group, err := h.groupLogic.Get(c, groupIDStr)
	if err != nil {
		_ = c.Error(err)
		return
	}

	// 返回成功响应
	resp := apiLinkGroup.DetailResponse{LinkGroupDetailDTO: *group}
	xResult.SuccessHasData(c, "获取友链分组详情成功", resp)
}

// GetList 获取友链分组列表
// @Summary 获取友链分组列表
// @Description 获取友链分组列表（不分页），支持过滤和排序
// @Tags 友链分组管理
// @Accept json
// @Produce json
// @Security Bearer
// @Param status query int false "状态过滤（0=禁用，1=启用）"
// @Param name query string false "名称模糊搜索"
// @Param with_links query bool false "是否包含友链列表"
// @Param only_enabled query bool false "仅查询启用的分组"
// @Param order_by query string false "排序字段（name, sort_order, created_at）"
// @Param order query string false "排序方向（asc, desc）"
// @Success 200 {object} []dto.LinkGroupListDTO "获取成功"
// @Failure 400 {object} map[string]interface{} "请求参数错误"
// @Failure 401 {object} map[string]interface{} "未认证"
// @Failure 500 {object} map[string]interface{} "服务器内部错误"
// @Router /api/v1/admin/groups/all [get]
func (h *LinkGroupHandler) GetList(c *gin.Context) {
	var req apiLinkGroup.ListRequest

	// 绑定查询参数
	bindErr := c.ShouldBindQuery(&req)
	if bindErr != nil {
		xValid.HandleValidationError(c, bindErr)
		return
	}

	// 调用服务层
	groups, err := h.groupLogic.GetList(c, &req)
	if err != nil {
		_ = c.Error(err)
		return
	}

	// 返回成功响应
	xResult.SuccessHasData(c, "获取友链分组列表成功", groups)
}

// GetPage 获取友链分组分页列表
// @Summary 获取友链分组分页列表
// @Description 分页获取友链分组列表，支持过滤和排序
// @Tags 友链分组管理
// @Accept json
// @Produce json
// @Security Bearer
// @Param page query int false "页码（默认1）"
// @Param page_size query int false "每页数量（默认10，最大100）"
// @Param status query int false "状态过滤（0=禁用，1=启用）"
// @Param name query string false "名称模糊搜索"
// @Param order_by query string false "排序字段（name, sort_order, created_at）"
// @Param order query string false "排序方向（asc, desc）"
// @Success 200 {object} apiLinkGroup.PageResponse "获取成功"
// @Failure 400 {object} map[string]interface{} "请求参数错误"
// @Failure 401 {object} map[string]interface{} "未认证"
// @Failure 500 {object} map[string]interface{} "服务器内部错误"
// @Router /api/v1/admin/groups [get]
func (h *LinkGroupHandler) GetPage(c *gin.Context) {
	var req apiLinkGroup.PageRequest

	// 绑定查询参数
	bindErr := c.ShouldBindQuery(&req)
	if bindErr != nil {
		xValid.HandleValidationError(c, bindErr)
		return
	}

	// 调用服务层
	result, err := h.groupLogic.GetPage(c, &req)
	if err != nil {
		_ = c.Error(err)
		return
	}

	// 返回成功响应
	resp := apiLinkGroup.PageResponse{PaginationResponse: *result}
	xResult.SuccessHasData(c, "获取友链分组分页列表成功", resp)
}
