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
	apiAuth "github.com/bamboo-services/bamboo-main/api/auth"

	xResult "github.com/bamboo-services/bamboo-base-go/major/result"
	xUtil "github.com/bamboo-services/bamboo-base-go/major/utility"
	xValid "github.com/bamboo-services/bamboo-base-go/major/validator"
	"github.com/gin-gonic/gin"
)

// List 获取系统用户列表
//
// @Summary [管理] 获取系统用户列表
// @Description 分页查询系统用户列表，支持关键词搜索、状态筛选与排序
// @Tags 系统用户管理接口
// @Accept json
// @Produce json
// @Security Bearer
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页数量" default(10)
// @Param keyword query string false "搜索关键词（用户名/邮箱/昵称模糊匹配）"
// @Param status query int false "用户状态 0:禁用 1:启用" Enums(0,1)
// @Param sort_by query string false "排序字段" Enums(created_at, updated_at, username, email, last_login_at)
// @Param sort_order query string false "排序方式" Enums(asc, desc)
// @Success 200 {object} xBase.BaseResponse{data=apiAuth.UserListResponse} "获取成功"
// @Failure 400 {object} xBase.BaseResponse "请求参数错误"
// @Failure 401 {object} xBase.BaseResponse "未认证"
// @Failure 500 {object} xBase.BaseResponse "服务器内部错误"
// @Router /api/v1/admin/users [GET]
func (h *SystemUserHandler) List(c *gin.Context) {
	var req apiAuth.UserQueryRequest

	// 绑定查询参数
	bindErr := c.ShouldBindQuery(&req)
	if bindErr != nil {
		xValid.HandleValidationError(c, bindErr)
		return
	}

	// 调用服务层
	resp, err := h.service.systemUserLogic.ListUsers(c.Request.Context(), &req)
	if err != nil {
		_ = c.Error(err)
		return
	}

	// 返回成功响应
	xResult.SuccessHasData(c, "获取成功", resp)
}

// UpdateStatus 启用/禁用用户
//
// @Summary [管理] 启用/禁用用户
// @Description 启用或禁用指定用户账号；系统唯一管理员不可被禁用
// @Tags 系统用户管理接口
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path int true "用户ID"
// @Param request body apiAuth.UserStatusRequest true "更新用户状态请求"
// @Success 200 {object} xBase.BaseResponse "更新成功"
// @Failure 400 {object} xBase.BaseResponse "请求参数错误"
// @Failure 401 {object} xBase.BaseResponse "未认证"
// @Failure 404 {object} xBase.BaseResponse "用户不存在"
// @Failure 500 {object} xBase.BaseResponse "服务器内部错误"
// @Router /api/v1/admin/users/{id}/status [PATCH]
func (h *SystemUserHandler) UpdateStatus(c *gin.Context) {
	uri := xUtil.Bind(c, &apiAuth.UserIDRequest{}).URI()
	if uri == nil {
		return
	}

	var req apiAuth.UserStatusRequest
	bindErr := c.ShouldBindJSON(&req)
	if bindErr != nil {
		xValid.HandleValidationError(c, bindErr)
		return
	}

	// 调用服务层
	if err := h.service.systemUserLogic.UpdateUserStatus(c.Request.Context(), uri.ID, &req); err != nil {
		_ = c.Error(err)
		return
	}

	// 返回成功响应
	xResult.Success(c, "状态更新成功")
}
