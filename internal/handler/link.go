package handler

import (
	"bamboo-main/internal/model/dto/response"
	"bamboo-main/internal/model/request"
	"bamboo-main/internal/service"

	xResult "github.com/bamboo-services/bamboo-base-go/result"
	xValid "github.com/bamboo-services/bamboo-base-go/validator"
	"github.com/gin-gonic/gin"
)

// LinkHandler 友情链接处理器
type LinkHandler struct {
	linkService service.ILinkService
}

// NewLinkHandler 创建友情链接处理器
func NewLinkHandler() *LinkHandler {
	return &LinkHandler{
		linkService: service.NewLinkService(),
	}
}

// Add 添加友情链接
// @Summary 添加友情链接
// @Description 添加新的友情链接申请
// @Tags 友情链接管理
// @Accept json
// @Produce json
// @Security Bearer
// @Param request body request.LinkFriendAddReq true "添加友情链接请求"
// @Success 200 {object} response.LinkAddResponse "添加成功"
// @Failure 400 {object} map[string]interface{} "请求参数错误"
// @Failure 401 {object} map[string]interface{} "未认证"
// @Failure 500 {object} map[string]interface{} "服务器内部错误"
// @Router /api/v1/admin/links [post]
func (h *LinkHandler) Add(c *gin.Context) {
	var req request.LinkFriendAddReq

	// 绑定请求数据
	bindErr := c.ShouldBindJSON(&req)
	if bindErr != nil {
		xValid.HandleValidationError(c, bindErr)
		return
	}

	// 调用服务层
	link, err := h.linkService.Add(c, &req)
	if err != nil {
		_ = c.Error(err)
		return
	}

	// 返回成功响应
	resp := response.LinkAddResponse{LinkFriendDTO: *link}
	xResult.SuccessHasData(c, "友情链接添加成功", resp)
}

// Update 更新友情链接
// @Summary 更新友情链接
// @Description 更新指定的友情链接信息
// @Tags 友情链接管理
// @Accept json
// @Produce json
// @Security Bearer
// @Param link_uuid path string true "友情链接UUID"
// @Param request body request.LinkFriendUpdateReq true "更新友情链接请求"
// @Success 200 {object} response.LinkUpdateResponse "更新成功"
// @Failure 400 {object} map[string]interface{} "请求参数错误"
// @Failure 401 {object} map[string]interface{} "未认证"
// @Failure 404 {object} map[string]interface{} "友情链接不存在"
// @Failure 500 {object} map[string]interface{} "服务器内部错误"
// @Router /api/v1/admin/links/{link_uuid} [put]
func (h *LinkHandler) Update(c *gin.Context) {
	linkUUID := c.Param("link_uuid")
	var req request.LinkFriendUpdateReq

	// 绑定请求数据
	bindErr := c.ShouldBindJSON(&req)
	if bindErr != nil {
		xValid.HandleValidationError(c, bindErr)
		return
	}

	// 调用服务层
	link, err := h.linkService.Update(c, linkUUID, &req)
	if err != nil {
		_ = c.Error(err)
		return
	}

	// 返回成功响应
	resp := response.LinkUpdateResponse{LinkFriendDTO: *link}
	xResult.SuccessHasData(c, "友情链接更新成功", resp)
}

// Delete 删除友情链接
// @Summary 删除友情链接
// @Description 删除指定的友情链接
// @Tags 友情链接管理
// @Accept json
// @Produce json
// @Security Bearer
// @Param link_uuid path string true "友情链接UUID"
// @Success 200 {object} response.MessageResponse "删除成功"
// @Failure 400 {object} map[string]interface{} "请求参数错误"
// @Failure 401 {object} map[string]interface{} "未认证"
// @Failure 404 {object} map[string]interface{} "友情链接不存在"
// @Failure 500 {object} map[string]interface{} "服务器内部错误"
// @Router /api/v1/admin/links/{link_uuid} [delete]
func (h *LinkHandler) Delete(c *gin.Context) {
	linkUUID := c.Param("link_uuid")

	// 调用服务层
	err := h.linkService.Delete(c, linkUUID)
	if err != nil {
		_ = c.Error(err)
		return
	}

	// 返回成功响应
	xResult.Success(c, "友情链接删除成功")
}

// Get 获取友情链接详情
// @Summary 获取友情链接详情
// @Description 获取指定友情链接的详细信息
// @Tags 友情链接管理
// @Accept json
// @Produce json
// @Security Bearer
// @Param link_uuid path string true "友情链接UUID"
// @Success 200 {object} response.LinkDetailResponse "获取成功"
// @Failure 400 {object} map[string]interface{} "请求参数错误"
// @Failure 401 {object} map[string]interface{} "未认证"
// @Failure 404 {object} map[string]interface{} "友情链接不存在"
// @Failure 500 {object} map[string]interface{} "服务器内部错误"
// @Router /api/v1/admin/links/{link_uuid} [get]
func (h *LinkHandler) Get(c *gin.Context) {
	linkUUID := c.Param("link_uuid")

	// 调用服务层
	link, err := h.linkService.Get(c, linkUUID)
	if err != nil {
		_ = c.Error(err)
		return
	}

	// 返回成功响应
	resp := response.LinkDetailResponse{LinkFriendDTO: *link}
	xResult.SuccessHasData(c, "获取成功", resp)
}

// List 获取友情链接列表
// @Summary 获取友情链接列表
// @Description 分页查询友情链接列表
// @Tags 友情链接管理
// @Accept json
// @Produce json
// @Security Bearer
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页数量" default(10)
// @Param link_name query string false "友情链接名称"
// @Param link_status query int false "友情链接状态 0:待审核 1:已通过 2:已拒绝"
// @Param link_fail query int false "失效状态 0:正常 1:失效"
// @Param link_group_uuid query string false "分组UUID"
// @Param sort_by query string false "排序字段" Enums(created_at, updated_at, link_order, link_name)
// @Param sort_order query string false "排序方式" Enums(asc, desc)
// @Success 200 {object} response.LinkListResponse "获取成功"
// @Failure 400 {object} map[string]interface{} "请求参数错误"
// @Failure 401 {object} map[string]interface{} "未认证"
// @Failure 500 {object} map[string]interface{} "服务器内部错误"
// @Router /api/v1/admin/links [get]
func (h *LinkHandler) List(c *gin.Context) {
	var req request.LinkFriendQueryReq

	// 绑定查询参数
	bindErr := c.ShouldBindQuery(&req)
	if bindErr != nil {
		xValid.HandleValidationError(c, bindErr)
		return
	}

	// 调用服务层
	result, err := h.linkService.List(c, &req)
	if err != nil {
		_ = c.Error(err)
		return
	}

	// 返回成功响应
	resp := response.LinkListResponse{PaginationResponse: *result}
	xResult.SuccessHasData(c, "获取成功", resp)
}

// UpdateStatus 更新友情链接状态
// @Summary 更新友情链接状态
// @Description 审核友情链接，更新状态为通过或拒绝
// @Tags 友情链接管理
// @Accept json
// @Produce json
// @Security Bearer
// @Param link_uuid path string true "友情链接UUID"
// @Param request body request.LinkFriendStatusReq true "更新状态请求"
// @Success 200 {object} response.MessageResponse "更新成功"
// @Failure 400 {object} map[string]interface{} "请求参数错误"
// @Failure 401 {object} map[string]interface{} "未认证"
// @Failure 404 {object} map[string]interface{} "友情链接不存在"
// @Failure 500 {object} map[string]interface{} "服务器内部错误"
// @Router /api/v1/admin/links/{link_uuid}/status [put]
func (h *LinkHandler) UpdateStatus(c *gin.Context) {
	linkUUID := c.Param("link_uuid")
	var req request.LinkFriendStatusReq

	// 绑定请求数据
	bindErr := c.ShouldBindJSON(&req)
	if bindErr != nil {
		xValid.HandleValidationError(c, bindErr)
		return
	}

	// 调用服务层
	err := h.linkService.UpdateStatus(c, linkUUID, &req)
	if err != nil {
		_ = c.Error(err)
		return
	}

	// 返回成功响应
	xResult.Success(c, "状态更新成功")
}

// UpdateFailStatus 更新友情链接失效状态
// @Summary 更新友情链接失效状态
// @Description 标记友情链接为失效或恢复正常
// @Tags 友情链接管理
// @Accept json
// @Produce json
// @Security Bearer
// @Param link_uuid path string true "友情链接UUID"
// @Param request body request.LinkFriendFailReq true "更新失效状态请求"
// @Success 200 {object} response.MessageResponse "更新成功"
// @Failure 400 {object} map[string]interface{} "请求参数错误"
// @Failure 401 {object} map[string]interface{} "未认证"
// @Failure 404 {object} map[string]interface{} "友情链接不存在"
// @Failure 500 {object} map[string]interface{} "服务器内部错误"
// @Router /api/v1/admin/links/{link_uuid}/fail [put]
func (h *LinkHandler) UpdateFailStatus(c *gin.Context) {
	linkUUID := c.Param("link_uuid")
	var req request.LinkFriendFailReq

	// 绑定请求数据
	bindErr := c.ShouldBindJSON(&req)
	if bindErr != nil {
		xValid.HandleValidationError(c, bindErr)
		return
	}

	// 调用服务层
	err := h.linkService.UpdateFailStatus(c, linkUUID, &req)
	if err != nil {
		_ = c.Error(err)
		return
	}

	// 返回成功响应
	xResult.Success(c, "失效状态更新成功")
}

// GetPublicLinks 获取公开的友情链接
// @Summary 获取公开友情链接
// @Description 获取已通过审核且正常的友情链接列表，用于前台展示
// @Tags 公开接口
// @Accept json
// @Produce json
// @Param group_uuid query string false "分组UUID，不传则获取所有"
// @Success 200 {object} response.LinkPublicResponse "获取成功"
// @Failure 500 {object} map[string]interface{} "服务器内部错误"
// @Router /api/v1/public/links [get]
func (h *LinkHandler) GetPublicLinks(c *gin.Context) {
	groupUUID := c.Query("group_uuid")

	// 调用服务层
	links, err := h.linkService.GetPublicLinks(c, groupUUID)
	if err != nil {
		_ = c.Error(err)
		return
	}

	// 返回成功响应
	resp := response.LinkPublicResponse{Links: links}
	xResult.SuccessHasData(c, "获取成功", resp)
}
