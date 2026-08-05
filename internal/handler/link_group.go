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
	xResult "github.com/bamboo-services/bamboo-base-go/major/result"
	xUtil "github.com/bamboo-services/bamboo-base-go/major/utility"
	xValid "github.com/bamboo-services/bamboo-base-go/major/validator"
	apiLinkGroup "github.com/bamboo-services/bamboo-main/api/link"
	"github.com/gin-gonic/gin"
)

// Add 添加友链分组
//
// @Summary [管理] 添加友链分组
// @Description 创建新的友链分组
// @Tags 友链分组接口
// @Accept json
// @Produce json
// @Security Bearer
// @Param request body apiLinkGroup.GroupAddRequest true "添加友链分组请求"
// @Success 200 {object} xBase.BaseResponse{data=apiLinkGroup.GroupAddResponse} "添加成功"
// @Failure 400 {object} xBase.BaseResponse "请求参数错误"
// @Failure 401 {object} xBase.BaseResponse "未认证"
// @Failure 500 {object} xBase.BaseResponse "服务器内部错误"
// @Router /api/v1/admin/groups [POST]
func (h *LinkGroupHandler) Add(c *gin.Context) {
	var req apiLinkGroup.GroupAddRequest

	// 绑定请求数据
	bindErr := c.ShouldBindJSON(&req)
	if bindErr != nil {
		xValid.HandleValidationError(c, bindErr)
		return
	}

	// 调用服务层
	group, err := h.service.linkGroupLogic.Add(c.Request.Context(), &req)
	if err != nil {
		_ = c.Error(err)
		return
	}

	// 返回成功响应
	resp := apiLinkGroup.GroupAddResponse{LinkGroup: *group}
	xResult.SuccessHasData(c, "友链分组添加成功", resp)
}

// Update 更新友链分组
//
// @Summary [管理] 更新友链分组
// @Description 更新指定友链分组的名称和描述
// @Tags 友链分组接口
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path int true "友链分组ID"
// @Param request body apiLinkGroup.GroupUpdateRequest true "更新友链分组请求"
// @Success 200 {object} xBase.BaseResponse{data=apiLinkGroup.GroupUpdateResponse} "更新成功"
// @Failure 400 {object} xBase.BaseResponse "请求参数错误"
// @Failure 401 {object} xBase.BaseResponse "未认证"
// @Failure 404 {object} xBase.BaseResponse "友链分组不存在"
// @Failure 500 {object} xBase.BaseResponse "服务器内部错误"
// @Router /api/v1/admin/groups/{id} [PUT]
func (h *LinkGroupHandler) Update(c *gin.Context) {
	uri := xUtil.Bind(c, &apiLinkGroup.GroupIDRequest{}).URI()
	if uri == nil {
		return
	}

	var req apiLinkGroup.GroupUpdateRequest

	// 绑定请求数据
	bindErr := c.ShouldBindJSON(&req)
	if bindErr != nil {
		xValid.HandleValidationError(c, bindErr)
		return
	}

	// 调用服务层
	group, err := h.service.linkGroupLogic.Update(c.Request.Context(), uri.ID, &req)
	if err != nil {
		_ = c.Error(err)
		return
	}

	// 返回成功响应
	resp := apiLinkGroup.GroupUpdateResponse{LinkGroup: *group}
	xResult.SuccessHasData(c, "友链分组更新成功", resp)
}

// UpdateSort 批量更新友链分组排序
//
// @Summary [管理] 批量更新友链分组排序
// @Description 按照传入的UUID数组顺序重新设置分组排序
// @Tags 友链分组接口
// @Accept json
// @Produce json
// @Security Bearer
// @Param request body apiLinkGroup.GroupSortRequest true "分组排序请求"
// @Success 200 {object} xBase.BaseResponse{data=apiLinkGroup.GroupSortResponse} "排序更新成功"
// @Failure 400 {object} xBase.BaseResponse "请求参数错误"
// @Failure 401 {object} xBase.BaseResponse "未认证"
// @Failure 404 {object} xBase.BaseResponse "分组不存在"
// @Failure 500 {object} xBase.BaseResponse "服务器内部错误"
// @Router /api/v1/admin/groups/sort [PATCH]
func (h *LinkGroupHandler) UpdateSort(c *gin.Context) {
	var req apiLinkGroup.GroupSortRequest

	// 绑定请求数据
	bindErr := c.ShouldBindJSON(&req)
	if bindErr != nil {
		xValid.HandleValidationError(c, bindErr)
		return
	}

	// 调用服务层
	err := h.service.linkGroupLogic.UpdateSort(c.Request.Context(), &req)
	if err != nil {
		_ = c.Error(err)
		return
	}

	// 返回成功响应
	resp := apiLinkGroup.GroupSortResponse{
		Count: len(req.GroupIDs),
	}
	xResult.SuccessHasData(c, "分组排序更新成功", resp)
}

// UpdateStatus 更新友链分组状态
//
// @Summary [管理] 更新友链分组状态
// @Description 切换指定友链分组的启用/禁用状态
// @Tags 友链分组接口
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path int true "友链分组ID"
// @Param request body apiLinkGroup.GroupStatusRequest true "分组状态请求"
// @Success 200 {object} xBase.BaseResponse{data=apiLinkGroup.GroupStatusResponse} "状态更新成功"
// @Failure 400 {object} xBase.BaseResponse "请求参数错误"
// @Failure 401 {object} xBase.BaseResponse "未认证"
// @Failure 404 {object} xBase.BaseResponse "友链分组不存在"
// @Failure 500 {object} xBase.BaseResponse "服务器内部错误"
// @Router /api/v1/admin/groups/{id}/status [PATCH]
func (h *LinkGroupHandler) UpdateStatus(c *gin.Context) {
	uri := xUtil.Bind(c, &apiLinkGroup.GroupIDRequest{}).URI()
	if uri == nil {
		return
	}

	var req apiLinkGroup.GroupStatusRequest

	// 绑定请求数据
	bindErr := c.ShouldBindJSON(&req)
	if bindErr != nil {
		xValid.HandleValidationError(c, bindErr)
		return
	}

	// 调用服务层
	err := h.service.linkGroupLogic.UpdateStatus(c.Request.Context(), uri.ID, &req)
	if err != nil {
		_ = c.Error(err)
		return
	}

	// 返回成功响应
	statusText := "禁用"
	if req.Status {
		statusText = "启用"
	}
	resp := apiLinkGroup.GroupStatusResponse{
		Status: req.Status,
	}
	xResult.SuccessHasData(c, "分组已"+statusText, resp)
}

// Delete 删除友链分组
//
// @Summary [管理] 删除友链分组
// @Description 删除指定的友链分组，支持强制删除模式
// @Tags 友链分组接口
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path int true "友链分组ID"
// @Param force query bool false "是否强制删除（默认false）"
// @Success 200 {object} xBase.BaseResponse "删除成功"
// @Failure 400 {object} xBase.BaseResponse "请求参数错误"
// @Failure 401 {object} xBase.BaseResponse "未认证"
// @Failure 404 {object} xBase.BaseResponse "友链分组不存在"
// @Failure 409 {object} apiLinkGroup.GroupDeleteConflictResponse "存在关联数据冲突"
// @Failure 500 {object} xBase.BaseResponse "服务器内部错误"
// @Router /api/v1/admin/groups/{id} [DELETE]
func (h *LinkGroupHandler) Delete(c *gin.Context) {
	uri := xUtil.Bind(c, &apiLinkGroup.GroupIDRequest{}).URI()
	if uri == nil {
		return
	}

	// 获取force参数
	var req apiLinkGroup.GroupDeleteRequest
	req.Force = c.Query("force") == "true"

	// 调用服务层
	_, err := h.service.linkGroupLogic.Delete(c.Request.Context(), uri.ID, &req)
	if err != nil {
		_ = c.Error(err)
		return
	}

	// 返回成功响应
	xResult.Success(c, "友链分组删除成功")
}

// Get 获取友链分组详情
//
// @Summary [管理] 获取友链分组详情
// @Description 根据ID获取指定友链分组的详细信息
// @Tags 友链分组接口
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path int true "友链分组ID"
// @Success 200 {object} xBase.BaseResponse{data=apiLinkGroup.GroupDetailResponse} "获取成功"
// @Failure 400 {object} xBase.BaseResponse "请求参数错误"
// @Failure 401 {object} xBase.BaseResponse "未认证"
// @Failure 404 {object} xBase.BaseResponse "友链分组不存在"
// @Failure 500 {object} xBase.BaseResponse "服务器内部错误"
// @Router /api/v1/admin/groups/{id} [GET]
func (h *LinkGroupHandler) Get(c *gin.Context) {
	uri := xUtil.Bind(c, &apiLinkGroup.GroupIDRequest{}).URI()
	if uri == nil {
		return
	}

	// 调用服务层
	group, err := h.service.linkGroupLogic.Get(c.Request.Context(), uri.ID)
	if err != nil {
		_ = c.Error(err)
		return
	}

	// 返回成功响应
	resp := apiLinkGroup.GroupDetailResponse{LinkGroup: *group}
	xResult.SuccessHasData(c, "获取友链分组详情成功", resp)
}

// GetList 获取友链分组列表
//
// @Summary [管理] 获取友链分组列表
// @Description 获取友链分组列表（不分页），支持过滤和排序
// @Tags 友链分组接口
// @Accept json
// @Produce json
// @Security Bearer
// @Param status query int false "状态过滤（0=禁用，1=启用）"
// @Param name query string false "名称模糊搜索"
// @Param with_links query bool false "是否包含友链列表"
// @Param only_enabled query bool false "仅查询启用的分组"
// @Param order_by query string false "排序字段（name, sort_order, created_at）"
// @Param order query string false "排序方向（asc, desc）"
// @Success 200 {object} xBase.BaseResponse{data=[]entity.LinkGroup} "获取成功"
// @Failure 400 {object} xBase.BaseResponse "请求参数错误"
// @Failure 401 {object} xBase.BaseResponse "未认证"
// @Failure 500 {object} xBase.BaseResponse "服务器内部错误"
// @Router /api/v1/admin/groups/all [GET]
func (h *LinkGroupHandler) GetList(c *gin.Context) {
	var req apiLinkGroup.GroupListRequest

	// 绑定查询参数
	bindErr := c.ShouldBindQuery(&req)
	if bindErr != nil {
		xValid.HandleValidationError(c, bindErr)
		return
	}

	// 调用服务层
	groups, err := h.service.linkGroupLogic.GetList(c.Request.Context(), &req)
	if err != nil {
		_ = c.Error(err)
		return
	}

	// 返回成功响应
	xResult.SuccessHasData(c, "获取友链分组列表成功", groups)
}

// GetPage 获取友链分组分页列表
//
// @Summary [管理] 获取友链分组分页列表
// @Description 分页获取友链分组列表，支持过滤和排序
// @Tags 友链分组接口
// @Accept json
// @Produce json
// @Security Bearer
// @Param page query int false "页码（默认1）"
// @Param page_size query int false "每页数量（默认10，最大100）"
// @Param status query int false "状态过滤（0=禁用，1=启用）"
// @Param name query string false "名称模糊搜索"
// @Param order_by query string false "排序字段（name, sort_order, created_at）"
// @Param order query string false "排序方向（asc, desc）"
// @Success 200 {object} xBase.BaseResponse{data=apiLinkGroup.GroupPageResponse} "获取成功"
// @Failure 400 {object} xBase.BaseResponse "请求参数错误"
// @Failure 401 {object} xBase.BaseResponse "未认证"
// @Failure 500 {object} xBase.BaseResponse "服务器内部错误"
// @Router /api/v1/admin/groups [GET]
func (h *LinkGroupHandler) GetPage(c *gin.Context) {
	var req apiLinkGroup.GroupPageRequest

	// 绑定查询参数
	bindErr := c.ShouldBindQuery(&req)
	if bindErr != nil {
		xValid.HandleValidationError(c, bindErr)
		return
	}

	// 调用服务层
	result, err := h.service.linkGroupLogic.GetPage(c.Request.Context(), &req)
	if err != nil {
		_ = c.Error(err)
		return
	}

	// 返回成功响应
	resp := apiLinkGroup.GroupPageResponse{PaginationResponse: *result}
	xResult.SuccessHasData(c, "获取友链分组分页列表成功", resp)
}
