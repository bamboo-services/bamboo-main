/*
 * --------------------------------------------------------------------------------
 * Copyright (c) 2016-NOW(至今) 筱锋
 * Author: 筱锋「xiao_lfeng」(https://www.x-lf.com)
 * --------------------------------------------------------------------------------
 * 许可证声明:版权所有 (c) 2016-2026 筱锋。保留所有权利。
 * 有关MIT许可证的更多信息,请查看项目根目录下的LICENSE文件或访问:
 * https://opensource.org/licenses/MIT
 * --------------------------------------------------------------------------------
 */

package handler

import (
	xResult "github.com/bamboo-services/bamboo-base-go/major/result"
	xUtil "github.com/bamboo-services/bamboo-base-go/major/utility"
	xValid "github.com/bamboo-services/bamboo-base-go/major/validator"
	apiSponsorRecord "github.com/bamboo-services/bamboo-main/api/sponsor"
	"github.com/gin-gonic/gin"
)

// Add 添加赞助记录
//
// @Summary [管理] 添加赞助记录
// @Description 创建新的赞助记录
// @Tags 赞助记录接口
// @Accept json
// @Produce json
// @Security Bearer
// @Param request body apiSponsorRecord.RecordAddRequest true "添加赞助记录请求"
// @Success 200 {object} xBase.BaseResponse{data=apiSponsorRecord.RecordAddResponse} "添加成功"
// @Failure 400 {object} xBase.BaseResponse "请求参数错误"
// @Failure 401 {object} xBase.BaseResponse "未认证"
// @Failure 500 {object} xBase.BaseResponse "服务器内部错误"
// @Router /api/v1/admin/sponsors/records [POST]
func (h *SponsorRecordHandler) Add(c *gin.Context) {
	var req apiSponsorRecord.RecordAddRequest

	// 绑定请求数据
	bindErr := c.ShouldBindJSON(&req)
	if bindErr != nil {
		xValid.HandleValidationError(c, bindErr)
		return
	}

	// 调用服务层
	record, err := h.service.sponsorRecordLogic.Add(c.Request.Context(), &req)
	if err != nil {
		_ = c.Error(err)
		return
	}

	// 返回成功响应
	resp := apiSponsorRecord.RecordAddResponse{RecordEntityResponse: *record}
	xResult.SuccessHasData(c, "赞助记录添加成功", resp)
}

// Update 更新赞助记录
//
// @Summary [管理] 更新赞助记录
// @Description 更新指定赞助记录的信息
// @Tags 赞助记录接口
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path int true "赞助记录ID"
// @Param request body apiSponsorRecord.RecordUpdateRequest true "更新赞助记录请求"
// @Success 200 {object} xBase.BaseResponse{data=apiSponsorRecord.RecordUpdateResponse} "更新成功"
// @Failure 400 {object} xBase.BaseResponse "请求参数错误"
// @Failure 401 {object} xBase.BaseResponse "未认证"
// @Failure 404 {object} xBase.BaseResponse "赞助记录不存在"
// @Failure 500 {object} xBase.BaseResponse "服务器内部错误"
// @Router /api/v1/admin/sponsors/records/{id} [PUT]
func (h *SponsorRecordHandler) Update(c *gin.Context) {
	uri := xUtil.Bind(c, &apiSponsorRecord.RecordIDRequest{}).URI()
	if uri == nil {
		return
	}

	var req apiSponsorRecord.RecordUpdateRequest

	// 绑定请求数据
	bindErr := c.ShouldBindJSON(&req)
	if bindErr != nil {
		xValid.HandleValidationError(c, bindErr)
		return
	}

	// 调用服务层
	record, err := h.service.sponsorRecordLogic.Update(c.Request.Context(), uri.ID, &req)
	if err != nil {
		_ = c.Error(err)
		return
	}

	// 返回成功响应
	resp := apiSponsorRecord.RecordUpdateResponse{RecordEntityResponse: *record}
	xResult.SuccessHasData(c, "赞助记录更新成功", resp)
}

// Delete 删除赞助记录
//
// @Summary [管理] 删除赞助记录
// @Description 删除指定的赞助记录
// @Tags 赞助记录接口
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path int true "赞助记录ID"
// @Success 200 {object} xBase.BaseResponse "删除成功"
// @Failure 400 {object} xBase.BaseResponse "请求参数错误"
// @Failure 401 {object} xBase.BaseResponse "未认证"
// @Failure 404 {object} xBase.BaseResponse "赞助记录不存在"
// @Failure 500 {object} xBase.BaseResponse "服务器内部错误"
// @Router /api/v1/admin/sponsors/records/{id} [DELETE]
func (h *SponsorRecordHandler) Delete(c *gin.Context) {
	uri := xUtil.Bind(c, &apiSponsorRecord.RecordIDRequest{}).URI()
	if uri == nil {
		return
	}

	// 调用服务层
	err := h.service.sponsorRecordLogic.Delete(c.Request.Context(), uri.ID)
	if err != nil {
		_ = c.Error(err)
		return
	}

	// 返回成功响应
	xResult.Success(c, "赞助记录删除成功")
}

// Get 获取赞助记录详情
//
// @Summary [管理] 获取赞助记录详情
// @Description 根据ID获取指定赞助记录的详细信息
// @Tags 赞助记录接口
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path int true "赞助记录ID"
// @Success 200 {object} xBase.BaseResponse{data=apiSponsorRecord.RecordDetailResponse} "获取成功"
// @Failure 400 {object} xBase.BaseResponse "请求参数错误"
// @Failure 401 {object} xBase.BaseResponse "未认证"
// @Failure 404 {object} xBase.BaseResponse "赞助记录不存在"
// @Failure 500 {object} xBase.BaseResponse "服务器内部错误"
// @Router /api/v1/admin/sponsors/records/{id} [GET]
func (h *SponsorRecordHandler) Get(c *gin.Context) {
	uri := xUtil.Bind(c, &apiSponsorRecord.RecordIDRequest{}).URI()
	if uri == nil {
		return
	}

	// 调用服务层
	record, err := h.service.sponsorRecordLogic.Get(c.Request.Context(), uri.ID)
	if err != nil {
		_ = c.Error(err)
		return
	}

	// 返回成功响应
	resp := apiSponsorRecord.RecordDetailResponse{RecordEntityResponse: *record}
	xResult.SuccessHasData(c, "获取赞助记录详情成功", resp)
}

// GetPage 获取赞助记录分页列表
//
// @Summary [管理] 获取赞助记录分页列表
// @Description 分页获取赞助记录列表，支持过滤和排序
// @Tags 赞助记录接口
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
// @Success 200 {object} xBase.BaseResponse{data=apiSponsorRecord.RecordPageResponse} "获取成功"
// @Failure 400 {object} xBase.BaseResponse "请求参数错误"
// @Failure 401 {object} xBase.BaseResponse "未认证"
// @Failure 500 {object} xBase.BaseResponse "服务器内部错误"
// @Router /api/v1/admin/sponsors/records [GET]
func (h *SponsorRecordHandler) GetPage(c *gin.Context) {
	var req apiSponsorRecord.RecordPageRequest

	// 绑定查询参数
	bindErr := c.ShouldBindQuery(&req)
	if bindErr != nil {
		xValid.HandleValidationError(c, bindErr)
		return
	}

	// 调用服务层
	result, err := h.service.sponsorRecordLogic.GetPage(c.Request.Context(), &req)
	if err != nil {
		_ = c.Error(err)
		return
	}

	// 返回成功响应
	resp := apiSponsorRecord.RecordPageResponse{PaginationResponse: *result}
	xResult.SuccessHasData(c, "获取赞助记录分页列表成功", resp)
}

// GetPublicPage 获取赞助记录公开分页列表
//
// @Summary [用户] 获取赞助记录公开分页列表
// @Description 分页获取前台赞助墙展示的记录列表，只返回未隐藏的记录，匿名记录显示为"匿名用户"
// @Tags 赞助记录接口
// @Accept json
// @Produce json
// @Param page query int false "页码（默认1）"
// @Param page_size query int false "每页数量（默认20，最大50）"
// @Param channel_id query int64 false "渠道ID过滤"
// @Param order_by query string false "排序字段（amount, sponsor_at, sort_order）"
// @Param order query string false "排序方向（asc, desc）"
// @Success 200 {object} xBase.BaseResponse{data=apiSponsorRecord.RecordPublicPageResponse} "获取成功"
// @Failure 400 {object} xBase.BaseResponse "请求参数错误"
// @Failure 500 {object} xBase.BaseResponse "服务器内部错误"
// @Router /api/v1/sponsors/records [GET]
func (h *SponsorRecordHandler) GetPublicPage(c *gin.Context) {
	var req apiSponsorRecord.RecordPublicPageRequest

	// 绑定查询参数
	bindErr := c.ShouldBindQuery(&req)
	if bindErr != nil {
		xValid.HandleValidationError(c, bindErr)
		return
	}

	// 调用服务层
	result, err := h.service.sponsorRecordLogic.GetPublicPage(c.Request.Context(), &req)
	if err != nil {
		_ = c.Error(err)
		return
	}

	// 返回成功响应
	resp := apiSponsorRecord.RecordPublicPageResponse{PaginationResponse: *result}
	xResult.SuccessHasData(c, "获取公开赞助记录列表成功", resp)
}
