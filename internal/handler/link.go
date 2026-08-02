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
	apiLink "github.com/bamboo-services/bamboo-main/api/link"
	ctxUtil "github.com/bamboo-services/bamboo-main/pkg/util/ctx"

	xError "github.com/bamboo-services/bamboo-base-go/common/error"
	xResult "github.com/bamboo-services/bamboo-base-go/major/result"
	xUtil "github.com/bamboo-services/bamboo-base-go/major/utility"
	xValid "github.com/bamboo-services/bamboo-base-go/major/validator"
	"github.com/gin-gonic/gin"
)

// Add 添加友情链接
//
// @Summary [管理] 添加友情链接
// @Description 添加新的友情链接申请
// @Tags 友情链接接口
// @Accept json
// @Produce json
// @Security Bearer
// @Param request body apiLink.FriendAddRequest true "添加友情链接请求"
// @Success 200 {object} xBase.BaseResponse{data=apiLink.FriendAddResponse} "添加成功"
// @Failure 400 {object} xBase.BaseResponse "请求参数错误"
// @Failure 401 {object} xBase.BaseResponse "未认证"
// @Failure 500 {object} xBase.BaseResponse "服务器内部错误"
// @Router /api/v1/admin/links [POST]
func (h *LinkHandler) Add(c *gin.Context) {
	var req apiLink.FriendAddRequest

	// 绑定请求数据
	bindErr := c.ShouldBindJSON(&req)
	if bindErr != nil {
		xValid.HandleValidationError(c, bindErr)
		return
	}

	// 调用服务层
	link, err := h.service.linkLogic.Add(c.Request.Context(), &req)
	if err != nil {
		_ = c.Error(err)
		return
	}

	// 返回成功响应
	resp := apiLink.FriendAddResponse{LinkFriend: *link}
	xResult.SuccessHasData(c, "友情链接添加成功", resp)
}

// Update 更新友情链接
//
// @Summary [管理] 更新友情链接
// @Description 更新指定的友情链接信息
// @Tags 友情链接接口
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path int true "友情链接ID"
// @Param request body apiLink.FriendUpdateRequest true "更新友情链接请求"
// @Success 200 {object} xBase.BaseResponse{data=apiLink.FriendUpdateResponse} "更新成功"
// @Failure 400 {object} xBase.BaseResponse "请求参数错误"
// @Failure 401 {object} xBase.BaseResponse "未认证"
// @Failure 404 {object} xBase.BaseResponse "友情链接不存在"
// @Failure 500 {object} xBase.BaseResponse "服务器内部错误"
// @Router /api/v1/admin/links/{id} [PUT]
func (h *LinkHandler) Update(c *gin.Context) {
	uri := xUtil.Bind(c, &apiLink.LinkIDRequest{}).URI()
	if uri == nil {
		return
	}

	var req apiLink.FriendUpdateRequest
	// 绑定请求数据
	bindErr := c.ShouldBindJSON(&req)
	if bindErr != nil {
		xValid.HandleValidationError(c, bindErr)
		return
	}

	// 调用服务层
	link, err := h.service.linkLogic.Update(c.Request.Context(), uri.ID, &req)
	if err != nil {
		_ = c.Error(err)
		return
	}

	// 返回成功响应
	resp := apiLink.FriendUpdateResponse{LinkFriend: *link}
	xResult.SuccessHasData(c, "友情链接更新成功", resp)
}

// Delete 删除友情链接
//
// @Summary [管理] 删除友情链接
// @Description 删除指定的友情链接
// @Tags 友情链接接口
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path int true "友情链接ID"
// @Success 200 {object} xBase.BaseResponse "删除成功"
// @Failure 400 {object} xBase.BaseResponse "请求参数错误"
// @Failure 401 {object} xBase.BaseResponse "未认证"
// @Failure 404 {object} xBase.BaseResponse "友情链接不存在"
// @Failure 500 {object} xBase.BaseResponse "服务器内部错误"
// @Router /api/v1/admin/links/{id} [DELETE]
func (h *LinkHandler) Delete(c *gin.Context) {
	uri := xUtil.Bind(c, &apiLink.LinkIDRequest{}).URI()
	if uri == nil {
		return
	}

	// 调用服务层
	err := h.service.linkLogic.Delete(c.Request.Context(), uri.ID)
	if err != nil {
		_ = c.Error(err)
		return
	}

	// 返回成功响应
	xResult.Success(c, "友情链接删除成功")
}

// Get 获取友情链接详情
//
// @Summary [管理] 获取友情链接详情
// @Description 获取指定友情链接的详细信息
// @Tags 友情链接接口
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path int true "友情链接ID"
// @Success 200 {object} xBase.BaseResponse{data=apiLink.FriendDetailResponse} "获取成功"
// @Failure 400 {object} xBase.BaseResponse "请求参数错误"
// @Failure 401 {object} xBase.BaseResponse "未认证"
// @Failure 404 {object} xBase.BaseResponse "友情链接不存在"
// @Failure 500 {object} xBase.BaseResponse "服务器内部错误"
// @Router /api/v1/admin/links/{id} [GET]
func (h *LinkHandler) Get(c *gin.Context) {
	uri := xUtil.Bind(c, &apiLink.LinkIDRequest{}).URI()
	if uri == nil {
		return
	}

	// 调用服务层
	link, err := h.service.linkLogic.Get(c.Request.Context(), uri.ID)
	if err != nil {
		_ = c.Error(err)
		return
	}

	// 返回成功响应
	resp := apiLink.FriendDetailResponse{LinkFriend: *link}
	xResult.SuccessHasData(c, "获取成功", resp)
}

// List 获取友情链接列表
//
// @Summary [管理] 获取友情链接列表
// @Description 分页查询友情链接列表
// @Tags 友情链接接口
// @Accept json
// @Produce json
// @Security Bearer
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页数量" default(10)
// @Param link_name query string false "友情链接名称"
// @Param link_status query int false "友情链接状态 0:待审核 1:已通过 2:已拒绝 3:下架待审核 4:已下架"
// @Param link_fail query int false "失效状态 0:正常 1:失效"
// @Param link_anomaly query bool false "异常过滤：status 非 0/1 或已失效（true）"
// @Param link_group_id query int64 false "分组ID"
// @Param sort_by query string false "排序字段" Enums(created_at, updated_at, link_order, link_name)
// @Param sort_order query string false "排序方式" Enums(asc, desc)
// @Success 200 {object} xBase.BaseResponse{data=apiLink.FriendListResponse} "获取成功"
// @Failure 400 {object} xBase.BaseResponse "请求参数错误"
// @Failure 401 {object} xBase.BaseResponse "未认证"
// @Failure 500 {object} xBase.BaseResponse "服务器内部错误"
// @Router /api/v1/admin/links [GET]
func (h *LinkHandler) List(c *gin.Context) {
	var req apiLink.FriendQueryRequest

	// 绑定查询参数
	bindErr := c.ShouldBindQuery(&req)
	if bindErr != nil {
		xValid.HandleValidationError(c, bindErr)
		return
	}

	// 调用服务层
	resp, err := h.service.linkLogic.List(c.Request.Context(), &req)
	if err != nil {
		_ = c.Error(err)
		return
	}

	// 返回成功响应
	xResult.SuccessHasData(c, "获取成功", resp)
}

// UpdateStatus 更新友情链接状态
//
// @Summary [管理] 更新友情链接状态
// @Description 审核友情链接，更新状态为通过或拒绝
// @Tags 友情链接接口
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path int true "友情链接ID"
// @Param request body apiLink.FriendStatusRequest true "更新状态请求"
// @Success 200 {object} xBase.BaseResponse "更新成功"
// @Failure 400 {object} xBase.BaseResponse "请求参数错误"
// @Failure 401 {object} xBase.BaseResponse "未认证"
// @Failure 404 {object} xBase.BaseResponse "友情链接不存在"
// @Failure 500 {object} xBase.BaseResponse "服务器内部错误"
// @Router /api/v1/admin/links/{id}/status [PUT]
func (h *LinkHandler) UpdateStatus(c *gin.Context) {
	uri := xUtil.Bind(c, &apiLink.LinkIDRequest{}).URI()
	if uri == nil {
		return
	}

	var req apiLink.FriendStatusRequest
	// 绑定请求数据
	bindErr := c.ShouldBindJSON(&req)
	if bindErr != nil {
		xValid.HandleValidationError(c, bindErr)
		return
	}

	// 调用服务层
	err := h.service.linkLogic.UpdateStatus(c.Request.Context(), uri.ID, &req)
	if err != nil {
		_ = c.Error(err)
		return
	}

	// 返回成功响应
	xResult.Success(c, "状态更新成功")
}

// UpdateFailStatus 更新友情链接失效状态
//
// @Summary [管理] 更新友情链接失效状态
// @Description 标记友情链接为失效或恢复正常
// @Tags 友情链接接口
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path int true "友情链接ID"
// @Param request body apiLink.FriendFailRequest true "更新失效状态请求"
// @Success 200 {object} xBase.BaseResponse "更新成功"
// @Failure 400 {object} xBase.BaseResponse "请求参数错误"
// @Failure 401 {object} xBase.BaseResponse "未认证"
// @Failure 404 {object} xBase.BaseResponse "友情链接不存在"
// @Failure 500 {object} xBase.BaseResponse "服务器内部错误"
// @Router /api/v1/admin/links/{id}/fail [PUT]
func (h *LinkHandler) UpdateFailStatus(c *gin.Context) {
	uri := xUtil.Bind(c, &apiLink.LinkIDRequest{}).URI()
	if uri == nil {
		return
	}

	var req apiLink.FriendFailRequest
	// 绑定请求数据
	bindErr := c.ShouldBindJSON(&req)
	if bindErr != nil {
		xValid.HandleValidationError(c, bindErr)
		return
	}

	// 调用服务层
	err := h.service.linkLogic.UpdateFailStatus(c.Request.Context(), uri.ID, &req)
	if err != nil {
		_ = c.Error(err)
		return
	}

	// 返回成功响应
	xResult.Success(c, "失效状态更新成功")
}

// ReScreenshot 重新生成友链站点截图
//
// @Summary [管理] 重新生成友链站点截图
// @Description 手动触发友链站点截图重新生成，任务进入截图队列串行处理
// @Tags 友情链接接口
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path int true "友情链接ID"
// @Success 200 {object} xBase.BaseResponse "已加入截图队列"
// @Failure 400 {object} xBase.BaseResponse "请求参数错误"
// @Failure 401 {object} xBase.BaseResponse "未认证"
// @Failure 404 {object} xBase.BaseResponse "友情链接不存在"
// @Failure 500 {object} xBase.BaseResponse "服务器内部错误"
// @Router /api/v1/admin/links/{id}/screenshot [POST]
func (h *LinkHandler) ReScreenshot(c *gin.Context) {
	uri := xUtil.Bind(c, &apiLink.LinkIDRequest{}).URI()
	if uri == nil {
		return
	}

	err := h.service.linkLogic.ReScreenshot(c.Request.Context(), uri.ID)
	if err != nil {
		_ = c.Error(err)
		return
	}

	xResult.Success(c, "已加入截图队列")
}

// GetPublicLinks 获取公开的友情链接
//
// @Summary [用户] 获取公开友情链接
// @Description 获取已通过审核且正常的友情链接列表，用于前台展示
// @Tags 公开接口
// @Accept json
// @Produce json
// @Param group_id query int64 false "分组ID，不传则获取所有"
// @Success 200 {object} xBase.BaseResponse{data=apiLink.FriendPublicResponse} "获取成功"
// @Failure 500 {object} xBase.BaseResponse "服务器内部错误"
// @Router /api/v1/public/links [GET]
func (h *LinkHandler) GetPublicLinks(c *gin.Context) {
	groupIDStr := c.Query("group_id")

	// 调用服务层
	links, err := h.service.linkLogic.GetPublicLinks(c.Request.Context(), groupIDStr)
	if err != nil {
		_ = c.Error(err)
		return
	}

	// 返回成功响应
	resp := apiLink.FriendPublicResponse{Links: links}
	xResult.SuccessHasData(c, "获取成功", resp)
}

// ApplyLink 访客自助申请友情链接
//
// @Summary [公开] 申请友情链接
// @Description 访客提交友情链接申请，进入待审核状态；联系邮箱用于确认友链归属
// @Tags 公开接口
// @Accept json
// @Produce json
// @Param request body apiLink.FriendApplyRequest true "申请友情链接请求"
// @Success 200 {object} xBase.BaseResponse{data=apiLink.FriendAddResponse} "申请成功"
// @Failure 400 {object} xBase.BaseResponse "请求参数错误"
// @Failure 500 {object} xBase.BaseResponse "服务器内部错误"
// @Router /api/v1/links/apply [POST]
func (h *LinkHandler) ApplyLink(c *gin.Context) {
	var req apiLink.FriendApplyRequest

	// 绑定请求数据
	bindErr := c.ShouldBindJSON(&req)
	if bindErr != nil {
		xValid.HandleValidationError(c, bindErr)
		return
	}

	// 调用服务层
	link, err := h.service.linkLogic.Apply(c.Request.Context(), &req)
	if err != nil {
		_ = c.Error(err)
		return
	}

	// 返回成功响应
	resp := apiLink.FriendAddResponse{LinkFriend: *link}
	xResult.SuccessHasData(c, "友链申请已提交，请等待管理员审核", resp)
}

// ListMyLinks 获取当前用户的友情链接列表
//
// @Summary [用户] 获取我的友情链接列表
// @Description 分页查询当前登录用户名下的友情链接
// @Tags 用户友链接口
// @Accept json
// @Produce json
// @Security Bearer
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页数量" default(10)
// @Param link_status query int false "友情链接状态 0:待审核 1:已通过 2:已拒绝 3:下架待审核 4:已下架"
// @Success 200 {object} xBase.BaseResponse{data=apiLink.FriendListResponse} "获取成功"
// @Failure 401 {object} xBase.BaseResponse "未认证"
// @Failure 500 {object} xBase.BaseResponse "服务器内部错误"
// @Router /api/v1/user/links [GET]
func (h *LinkHandler) ListMyLinks(c *gin.Context) {
	userID, exists := ctxUtil.GetUserID(c)
	if !exists {
		_ = c.Error(xError.NewError(c, xError.Unauthorized, "用户信息获取失败", false))
		return
	}

	var req apiLink.FriendUserQueryRequest
	// 绑定查询参数
	bindErr := c.ShouldBindQuery(&req)
	if bindErr != nil {
		xValid.HandleValidationError(c, bindErr)
		return
	}

	// 调用服务层
	result, err := h.service.linkLogic.ListMine(c.Request.Context(), userID, &req)
	if err != nil {
		_ = c.Error(err)
		return
	}

	// 返回成功响应
	resp := apiLink.FriendListResponse{PaginationResponse: *result}
	xResult.SuccessHasData(c, "获取成功", resp)
}

// GetMyLink 获取当前用户的友情链接详情
//
// @Summary [用户] 获取我的友情链接详情
// @Description 获取当前登录用户名下指定友情链接的详细信息
// @Tags 用户友链接口
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path int true "友情链接ID"
// @Success 200 {object} xBase.BaseResponse{data=apiLink.FriendDetailResponse} "获取成功"
// @Failure 401 {object} xBase.BaseResponse "未认证"
// @Failure 404 {object} xBase.BaseResponse "友情链接不存在"
// @Failure 500 {object} xBase.BaseResponse "服务器内部错误"
// @Router /api/v1/user/links/{id} [GET]
func (h *LinkHandler) GetMyLink(c *gin.Context) {
	userID, exists := ctxUtil.GetUserID(c)
	if !exists {
		_ = c.Error(xError.NewError(c, xError.Unauthorized, "用户信息获取失败", false))
		return
	}

	uri := xUtil.Bind(c, &apiLink.LinkIDRequest{}).URI()
	if uri == nil {
		return
	}

	// 调用服务层
	link, err := h.service.linkLogic.GetMine(c.Request.Context(), userID, uri.ID)
	if err != nil {
		_ = c.Error(err)
		return
	}

	// 返回成功响应
	resp := apiLink.FriendDetailResponse{LinkFriend: *link}
	xResult.SuccessHasData(c, "获取成功", resp)
}

// UpdateMyLink 更新当前用户的友情链接
//
// @Summary [用户] 更新我的友情链接
// @Description 更新当前登录用户名下指定友情链接的站点基础信息
// @Tags 用户友链接口
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path int true "友情链接ID"
// @Param request body apiLink.FriendUserUpdateRequest true "更新友情链接请求"
// @Success 200 {object} xBase.BaseResponse{data=apiLink.FriendUpdateResponse} "更新成功"
// @Failure 400 {object} xBase.BaseResponse "请求参数错误"
// @Failure 401 {object} xBase.BaseResponse "未认证"
// @Failure 404 {object} xBase.BaseResponse "友情链接不存在"
// @Failure 500 {object} xBase.BaseResponse "服务器内部错误"
// @Router /api/v1/user/links/{id} [PUT]
func (h *LinkHandler) UpdateMyLink(c *gin.Context) {
	userID, exists := ctxUtil.GetUserID(c)
	if !exists {
		_ = c.Error(xError.NewError(c, xError.Unauthorized, "用户信息获取失败", false))
		return
	}

	uri := xUtil.Bind(c, &apiLink.LinkIDRequest{}).URI()
	if uri == nil {
		return
	}

	var req apiLink.FriendUserUpdateRequest
	// 绑定请求数据
	bindErr := c.ShouldBindJSON(&req)
	if bindErr != nil {
		xValid.HandleValidationError(c, bindErr)
		return
	}

	// 调用服务层
	link, err := h.service.linkLogic.UpdateMine(c.Request.Context(), userID, uri.ID, &req)
	if err != nil {
		_ = c.Error(err)
		return
	}

	// 返回成功响应
	resp := apiLink.FriendUpdateResponse{LinkFriend: *link}
	xResult.SuccessHasData(c, "友情链接更新成功", resp)
}

// RequestTakedown 当前用户申请下架自己的友情链接
//
// @Summary [用户] 申请下架我的友情链接
// @Description 对已通过的友链发起下架申请，进入下架待审核状态
// @Tags 用户友链接口
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path int true "友情链接ID"
// @Success 200 {object} xBase.BaseResponse "申请成功"
// @Failure 400 {object} xBase.BaseResponse "状态不允许下架"
// @Failure 401 {object} xBase.BaseResponse "未认证"
// @Failure 404 {object} xBase.BaseResponse "友情链接不存在"
// @Failure 500 {object} xBase.BaseResponse "服务器内部错误"
// @Router /api/v1/user/links/{id}/takedown [PUT]
func (h *LinkHandler) RequestTakedown(c *gin.Context) {
	userID, exists := ctxUtil.GetUserID(c)
	if !exists {
		_ = c.Error(xError.NewError(c, xError.Unauthorized, "用户信息获取失败", false))
		return
	}

	uri := xUtil.Bind(c, &apiLink.LinkIDRequest{}).URI()
	if uri == nil {
		return
	}

	// 调用服务层
	err := h.service.linkLogic.RequestTakedown(c.Request.Context(), userID, uri.ID)
	if err != nil {
		_ = c.Error(err)
		return
	}

	// 返回成功响应
	xResult.Success(c, "下架申请已提交，请等待管理员审核")
}

// GetPublicGroups 获取启用的友链分组列表（公开接口，供申请表单选择器使用）
//
// @Summary [公开] 获取友链分组列表
// @Description 返回所有启用的友链分组，供访客申请时预选展示位置
// @Tags 友情链接接口
// @Produce json
// @Success 200 {object} xBase.BaseResponse{data=[]entity.LinkGroup} "获取成功"
// @Router /api/v1/links/groups [GET]
func (h *LinkHandler) GetPublicGroups(c *gin.Context) {
	onlyEnabled := true
	req := apiLink.GroupListRequest{
		OnlyEnabled: &onlyEnabled,
	}
	groups, err := h.service.linkGroupLogic.GetList(c.Request.Context(), &req)
	if err != nil {
		_ = c.Error(err)
		return
	}
	xResult.SuccessHasData(c, "获取友链分组列表成功", groups)
}

// GetPublicColors 获取启用的友链颜色列表（公开接口，供申请表单选择器使用）
//
// @Summary [公开] 获取友链颜色列表
// @Description 返回所有启用的友链颜色，供访客申请时预选展示颜色
// @Tags 友情链接接口
// @Produce json
// @Success 200 {object} xBase.BaseResponse{data=[]entity.LinkColor} "获取成功"
// @Router /api/v1/links/colors [GET]
func (h *LinkHandler) GetPublicColors(c *gin.Context) {
	onlyEnabled := true
	req := apiLink.ColorListRequest{
		OnlyEnabled: &onlyEnabled,
	}
	colors, err := h.service.linkColorLogic.GetList(c.Request.Context(), &req)
	if err != nil {
		_ = c.Error(err)
		return
	}
	xResult.SuccessHasData(c, "获取友链颜色列表成功", colors)
}
