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
	xError "github.com/bamboo-services/bamboo-base-go/common/error"
	xResult "github.com/bamboo-services/bamboo-base-go/major/result"
	xUtil "github.com/bamboo-services/bamboo-base-go/major/utility"
	xValid "github.com/bamboo-services/bamboo-base-go/major/validator"
	apiSponsorRecord "github.com/bamboo-services/bamboo-main/api/sponsor"
	"github.com/bamboo-services/bamboo-main/pkg/util/ctx"
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
// @Success 200 {object} xBase.BaseResponse{data=apiSponsorRecord.RecordAdminPageResponse} "获取成功"
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
	resp := apiSponsorRecord.RecordAdminPageResponse{
		PaginationResponse: result.PaginationResponse,
		PendingCount:       result.PendingCount,
	}
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

// ApplySponsor 访客自助申请赞助展示
//
// @Summary [公开] 申请赞助展示
// @Description 游客或登录用户提交赞助展示申请，需提供联系邮箱用于归属确认，提交后进入待审核状态
// @Tags 赞助记录接口
// @Accept json
// @Produce json
// @Param request body apiSponsorRecord.SponsorApplyRequest true "赞助展示申请请求"
// @Success 200 {object} xBase.BaseResponse{data=apiSponsorRecord.RecordAddResponse} "申请提交成功"
// @Failure 400 {object} xBase.BaseResponse "请求参数错误"
// @Failure 404 {object} xBase.BaseResponse "赞助渠道不存在"
// @Failure 500 {object} xBase.BaseResponse "服务器内部错误"
// @Router /api/v1/sponsors/apply [POST]
func (h *SponsorRecordHandler) ApplySponsor(c *gin.Context) {
	var req apiSponsorRecord.SponsorApplyRequest

	// 绑定请求数据
	bindErr := c.ShouldBindJSON(&req)
	if bindErr != nil {
		xValid.HandleValidationError(c, bindErr)
		return
	}

	// 调用服务层
	record, err := h.service.sponsorRecordLogic.Apply(c.Request.Context(), &req)
	if err != nil {
		_ = c.Error(err)
		return
	}

	// 返回成功响应
	resp := apiSponsorRecord.RecordAddResponse{RecordEntityResponse: *record}
	xResult.SuccessHasData(c, "赞助申请已提交，请等待管理员审核", resp)
}

// UpdateStatus 更新赞助记录审核状态
//
// @Summary [管理] 审核赞助记录
// @Description 审核赞助展示申请，可置为通过（1）或拒绝（2），并填写审核备注反馈给申请者
// @Tags 赞助记录接口
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path int true "赞助记录ID"
// @Param request body apiSponsorRecord.SponsorStatusRequest true "审核赞助记录请求"
// @Success 200 {object} xBase.BaseResponse "状态更新成功"
// @Failure 400 {object} xBase.BaseResponse "请求参数错误"
// @Failure 401 {object} xBase.BaseResponse "未认证"
// @Failure 404 {object} xBase.BaseResponse "赞助记录不存在"
// @Failure 500 {object} xBase.BaseResponse "服务器内部错误"
// @Router /api/v1/admin/sponsors/records/{id}/status [PUT]
func (h *SponsorRecordHandler) UpdateStatus(c *gin.Context) {
	uri := xUtil.Bind(c, &apiSponsorRecord.RecordIDRequest{}).URI()
	if uri == nil {
		return
	}

	var req apiSponsorRecord.SponsorStatusRequest

	// 绑定请求数据
	bindErr := c.ShouldBindJSON(&req)
	if bindErr != nil {
		xValid.HandleValidationError(c, bindErr)
		return
	}

	// 调用服务层
	err := h.service.sponsorRecordLogic.UpdateStatus(c.Request.Context(), uri.ID, &req)
	if err != nil {
		_ = c.Error(err)
		return
	}

	// 返回成功响应
	xResult.Success(c, "状态更新成功")
}

// ListMyRecords 获取当前用户的赞助记录列表
//
// @Summary [用户] 获取我的赞助记录列表
// @Description 分页查询当前登录用户名下的赞助记录
// @Tags 用户赞助接口
// @Accept json
// @Produce json
// @Security Bearer
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页数量" default(10)
// @Param sponsor_status query int false "赞助状态 0:待审核 1:已通过 2:已拒绝"
// @Success 200 {object} xBase.BaseResponse{data=apiSponsorRecord.RecordPageResponse} "获取成功"
// @Failure 401 {object} xBase.BaseResponse "未认证"
// @Failure 500 {object} xBase.BaseResponse "服务器内部错误"
// @Router /api/v1/user/sponsors [GET]
func (h *SponsorRecordHandler) ListMyRecords(c *gin.Context) {
	userID, exists := ctxUtil.GetUserID(c)
	if !exists {
		_ = c.Error(xError.NewError(c, xError.Unauthorized, "用户信息获取失败", false))
		return
	}

	var req apiSponsorRecord.SponsorUserQueryRequest
	// 绑定查询参数
	bindErr := c.ShouldBindQuery(&req)
	if bindErr != nil {
		xValid.HandleValidationError(c, bindErr)
		return
	}

	// 调用服务层
	result, err := h.service.sponsorRecordLogic.ListMine(c.Request.Context(), userID, &req)
	if err != nil {
		_ = c.Error(err)
		return
	}

	// 返回成功响应
	resp := apiSponsorRecord.RecordPageResponse{PaginationResponse: *result}
	xResult.SuccessHasData(c, "获取成功", resp)
}

// GetMyRecord 获取当前用户的赞助记录详情
//
// @Summary [用户] 获取我的赞助记录详情
// @Description 获取当前登录用户名下指定赞助记录的详细信息
// @Tags 用户赞助接口
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path int true "赞助记录ID"
// @Success 200 {object} xBase.BaseResponse{data=apiSponsorRecord.RecordDetailResponse} "获取成功"
// @Failure 401 {object} xBase.BaseResponse "未认证"
// @Failure 404 {object} xBase.BaseResponse "赞助记录不存在"
// @Failure 500 {object} xBase.BaseResponse "服务器内部错误"
// @Router /api/v1/user/sponsors/{id} [GET]
func (h *SponsorRecordHandler) GetMyRecord(c *gin.Context) {
	userID, exists := ctxUtil.GetUserID(c)
	if !exists {
		_ = c.Error(xError.NewError(c, xError.Unauthorized, "用户信息获取失败", false))
		return
	}

	uri := xUtil.Bind(c, &apiSponsorRecord.RecordIDRequest{}).URI()
	if uri == nil {
		return
	}

	// 调用服务层
	record, err := h.service.sponsorRecordLogic.GetMine(c.Request.Context(), userID, uri.ID)
	if err != nil {
		_ = c.Error(err)
		return
	}

	// 返回成功响应
	resp := apiSponsorRecord.RecordDetailResponse{RecordEntityResponse: *record}
	xResult.SuccessHasData(c, "获取成功", resp)
}

// UpdateMyRecord 更新当前用户的赞助记录
//
// @Summary [用户] 更新我的赞助记录
// @Description 更新当前登录用户名下指定赞助记录的展示信息（金额与渠道不可修改）
// @Tags 用户赞助接口
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path int true "赞助记录ID"
// @Param request body apiSponsorRecord.SponsorUserUpdateRequest true "更新赞助记录请求"
// @Success 200 {object} xBase.BaseResponse{data=apiSponsorRecord.RecordUpdateResponse} "更新成功"
// @Failure 400 {object} xBase.BaseResponse "请求参数错误"
// @Failure 401 {object} xBase.BaseResponse "未认证"
// @Failure 404 {object} xBase.BaseResponse "赞助记录不存在"
// @Failure 500 {object} xBase.BaseResponse "服务器内部错误"
// @Router /api/v1/user/sponsors/{id} [PUT]
func (h *SponsorRecordHandler) UpdateMyRecord(c *gin.Context) {
	userID, exists := ctxUtil.GetUserID(c)
	if !exists {
		_ = c.Error(xError.NewError(c, xError.Unauthorized, "用户信息获取失败", false))
		return
	}

	uri := xUtil.Bind(c, &apiSponsorRecord.RecordIDRequest{}).URI()
	if uri == nil {
		return
	}

	var req apiSponsorRecord.SponsorUserUpdateRequest
	// 绑定请求数据
	bindErr := c.ShouldBindJSON(&req)
	if bindErr != nil {
		xValid.HandleValidationError(c, bindErr)
		return
	}

	// 调用服务层
	record, err := h.service.sponsorRecordLogic.UpdateMine(c.Request.Context(), userID, uri.ID, &req)
	if err != nil {
		_ = c.Error(err)
		return
	}

	// 返回成功响应
	resp := apiSponsorRecord.RecordUpdateResponse{RecordEntityResponse: *record}
	xResult.SuccessHasData(c, "赞助记录更新成功", resp)
}
