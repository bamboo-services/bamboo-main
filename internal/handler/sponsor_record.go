/*
 * --------------------------------------------------------------------------------
 * Copyright (c) 2016-NOW(至今) 筱锋
 * Author: 筱锋「xiao_lfeng」(https://www.x-lf.com)
 * --------------------------------------------------------------------------------
 * 许可证声明:版权所有 (c) 2016-2025 筱锋。保留所有权利。
 * 有关MIT许可证的更多信息,请查看项目根目录下的LICENSE文件或访问:
 * https://opensource.org/licenses/MIT
 * --------------------------------------------------------------------------------
 */

package handler

import (
	xValid "github.com/bamboo-services/bamboo-base-go/common/validator"
	xResult "github.com/bamboo-services/bamboo-base-go/major/result"
	apiSponsorRecord "github.com/bamboo-services/bamboo-main/api/sponsor"
	"github.com/gin-gonic/gin"
)

// Add 添加赞助记录
// @Summary 添加赞助记录
// @Description 创建新的赞助记录
// @Tags 赞助记录管理
// @Accept json
// @Produce json
// @Security Bearer
// @Param request body apiSponsorRecord.RecordAddRequest true "添加赞助记录请求"
// @Success 200 {object} apiSponsorRecord.RecordAddResponse "添加成功"
// @Failure 400 {object} map[string]interface{} "请求参数错误"
// @Failure 401 {object} map[string]interface{} "未认证"
// @Failure 500 {object} map[string]interface{} "服务器内部错误"
// @Router /api/v1/admin/sponsors/records [post]
func (h *SponsorRecordHandler) Add(c *gin.Context) {
	var req apiSponsorRecord.RecordAddRequest

	// 绑定请求数据
	bindErr := c.ShouldBindJSON(&req)
	if bindErr != nil {
		xValid.HandleValidationError(c, bindErr)
		return
	}

	// 调用服务层
	record, err := h.service.sponsorRecordLogic.Add(c, &req)
	if err != nil {
		_ = c.Error(err)
		return
	}

	// 返回成功响应
	resp := apiSponsorRecord.RecordAddResponse{RecordEntityResponse: *record}
	xResult.SuccessHasData(c, "赞助记录添加成功", resp)
}

// Update 更新赞助记录
// @Summary 更新赞助记录
// @Description 更新指定赞助记录的信息
// @Tags 赞助记录管理
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path int64 true "赞助记录ID"
// @Param request body apiSponsorRecord.RecordUpdateRequest true "更新赞助记录请求"
// @Success 200 {object} apiSponsorRecord.RecordUpdateResponse "更新成功"
// @Failure 400 {object} map[string]interface{} "请求参数错误"
// @Failure 401 {object} map[string]interface{} "未认证"
// @Failure 404 {object} map[string]interface{} "赞助记录不存在"
// @Failure 500 {object} map[string]interface{} "服务器内部错误"
// @Router /api/v1/admin/sponsors/records/{id} [put]
func (h *SponsorRecordHandler) Update(c *gin.Context) {
	recordIDStr := c.Param("id")
	var req apiSponsorRecord.RecordUpdateRequest

	// 绑定请求数据
	bindErr := c.ShouldBindJSON(&req)
	if bindErr != nil {
		xValid.HandleValidationError(c, bindErr)
		return
	}

	// 调用服务层
	record, err := h.service.sponsorRecordLogic.Update(c, recordIDStr, &req)
	if err != nil {
		_ = c.Error(err)
		return
	}

	// 返回成功响应
	resp := apiSponsorRecord.RecordUpdateResponse{RecordEntityResponse: *record}
	xResult.SuccessHasData(c, "赞助记录更新成功", resp)
}

// Delete 删除赞助记录
// @Summary 删除赞助记录
// @Description 删除指定的赞助记录
// @Tags 赞助记录管理
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path int64 true "赞助记录ID"
// @Success 200 {object} apiSponsorRecord.RecordDeleteResponse "删除成功"
// @Failure 400 {object} map[string]interface{} "请求参数错误"
// @Failure 401 {object} map[string]interface{} "未认证"
// @Failure 404 {object} map[string]interface{} "赞助记录不存在"
// @Failure 500 {object} map[string]interface{} "服务器内部错误"
// @Router /api/v1/admin/sponsors/records/{id} [delete]
func (h *SponsorRecordHandler) Delete(c *gin.Context) {
	recordIDStr := c.Param("id")

	// 调用服务层
	err := h.service.sponsorRecordLogic.Delete(c, recordIDStr)
	if err != nil {
		_ = c.Error(err)
		return
	}

	// 返回成功响应
	resp := apiSponsorRecord.RecordDeleteResponse{
		Message: "赞助记录删除成功",
	}
	xResult.SuccessHasData(c, "赞助记录删除成功", resp)
}

// Get 获取赞助记录详情
// @Summary 获取赞助记录详情
// @Description 根据ID获取指定赞助记录的详细信息
// @Tags 赞助记录管理
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path int64 true "赞助记录ID"
// @Success 200 {object} apiSponsorRecord.RecordDetailResponse "获取成功"
// @Failure 400 {object} map[string]interface{} "请求参数错误"
// @Failure 401 {object} map[string]interface{} "未认证"
// @Failure 404 {object} map[string]interface{} "赞助记录不存在"
// @Failure 500 {object} map[string]interface{} "服务器内部错误"
// @Router /api/v1/admin/sponsors/records/{id} [get]
func (h *SponsorRecordHandler) Get(c *gin.Context) {
	recordIDStr := c.Param("id")

	// 调用服务层
	record, err := h.service.sponsorRecordLogic.Get(c, recordIDStr)
	if err != nil {
		_ = c.Error(err)
		return
	}

	// 返回成功响应
	resp := apiSponsorRecord.RecordDetailResponse{RecordEntityResponse: *record}
	xResult.SuccessHasData(c, "获取赞助记录详情成功", resp)
}

// GetPage 获取赞助记录分页列表
// @Summary 获取赞助记录分页列表
// @Description 分页获取赞助记录列表，支持过滤和排序
// @Tags 赞助记录管理
// @Accept json
// @Produce json
// @Security Bearer
// @Param page query int false "页码（默认1）"
// @Param page_size query int false "每页数量（默认10，最大100）"
// @Param channel_id query int64 false "渠道ID过滤"
// @Param nickname query string false "昵称模糊搜索"
// @Param is_anonymous query bool false "是否匿名过滤"
// @Param is_hidden query bool false "是否隐藏过滤"
// @Param order_by query string false "排序字段（nickname, amount, sponsor_at, sort_order, created_at）"
// @Param order query string false "排序方向（asc, desc）"
// @Success 200 {object} apiSponsorRecord.RecordPageResponse "获取成功"
// @Failure 400 {object} map[string]interface{} "请求参数错误"
// @Failure 401 {object} map[string]interface{} "未认证"
// @Failure 500 {object} map[string]interface{} "服务器内部错误"
// @Router /api/v1/admin/sponsors/records [get]
func (h *SponsorRecordHandler) GetPage(c *gin.Context) {
	var req apiSponsorRecord.RecordPageRequest

	// 绑定查询参数
	bindErr := c.ShouldBindQuery(&req)
	if bindErr != nil {
		xValid.HandleValidationError(c, bindErr)
		return
	}

	// 调用服务层
	result, err := h.service.sponsorRecordLogic.GetPage(c, &req)
	if err != nil {
		_ = c.Error(err)
		return
	}

	// 返回成功响应
	resp := apiSponsorRecord.RecordPageResponse{PaginationResponse: *result}
	xResult.SuccessHasData(c, "获取赞助记录分页列表成功", resp)
}

// GetPublicPage 获取赞助记录公开分页列表
// @Summary 获取赞助记录公开分页列表
// @Description 分页获取前台赞助墙展示的记录列表，只返回未隐藏的记录，匿名记录显示为"匿名用户"
// @Tags 赞助记录公开接口
// @Accept json
// @Produce json
// @Param page query int false "页码（默认1）"
// @Param page_size query int false "每页数量（默认20，最大50）"
// @Param channel_id query int64 false "渠道ID过滤"
// @Param order_by query string false "排序字段（amount, sponsor_at, sort_order）"
// @Param order query string false "排序方向（asc, desc）"
// @Success 200 {object} apiSponsorRecord.RecordPublicPageResponse "获取成功"
// @Failure 400 {object} map[string]interface{} "请求参数错误"
// @Failure 500 {object} map[string]interface{} "服务器内部错误"
// @Router /api/v1/sponsors/records [get]
func (h *SponsorRecordHandler) GetPublicPage(c *gin.Context) {
	var req apiSponsorRecord.RecordPublicPageRequest

	// 绑定查询参数
	bindErr := c.ShouldBindQuery(&req)
	if bindErr != nil {
		xValid.HandleValidationError(c, bindErr)
		return
	}

	// 调用服务层
	result, err := h.service.sponsorRecordLogic.GetPublicPage(c, &req)
	if err != nil {
		_ = c.Error(err)
		return
	}

	// 返回成功响应
	resp := apiSponsorRecord.RecordPublicPageResponse{PaginationResponse: *result}
	xResult.SuccessHasData(c, "获取公开赞助记录列表成功", resp)
}
