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
	"bamboo-main/internal/model/dto/response"
	"bamboo-main/internal/model/request"
	"bamboo-main/internal/service"
	ctxUtil "bamboo-main/pkg/util/ctx"
	"errors"

	xResult "github.com/bamboo-services/bamboo-base-go/result"
	xValid "github.com/bamboo-services/bamboo-base-go/validator"
	"github.com/gin-gonic/gin"
)

// AuthHandler 认证处理器
type AuthHandler struct {
	authService service.IAuthService
}

// NewAuthHandler 创建认证处理器
func NewAuthHandler() *AuthHandler {
	return &AuthHandler{
		authService: service.NewAuthService(),
	}
}

// Login 用户登录
// @Summary 用户登录
// @Description 管理员用户登录，返回用户信息、访问令牌及Token时间信息
// @Tags 认证管理
// @Accept json
// @Produce json
// @Param request body request.AuthLoginReq true "登录请求"
// @Success 200 {object} response.AuthLoginResponse "登录成功，包含用户信息、Token及时间信息"
// @Failure 400 {object} map[string]interface{} "请求参数错误"
// @Failure 401 {object} map[string]interface{} "认证失败"
// @Failure 500 {object} map[string]interface{} "服务器内部错误"
// @Router /api/v1/auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	var req request.AuthLoginReq

	// 绑定请求数据
	bindErr := c.ShouldBindJSON(&req)
	if bindErr != nil {
		xValid.HandleValidationError(c, bindErr)
		return
	}

	// 调用服务层
	user, token, createdAt, expiredAt, err := h.authService.Login(c, &req)
	if err != nil {
		_ = c.Error(err)
		return
	}

	// 返回成功响应
	resp := response.AuthLoginResponse{
		User:      *user,
		Token:     token,
		CreatedAt: *createdAt,
		ExpiredAt: *expiredAt,
	}
	xResult.SuccessHasData(c, "登录成功", resp)
}

// Logout 用户登出
// @Summary 用户登出
// @Description 注销当前登录会话
// @Tags 认证管理
// @Accept json
// @Produce json
// @Security Bearer
// @Success 200 {object} response.AuthLogoutResponse "登出成功"
// @Failure 401 {object} map[string]interface{} "未认证"
// @Failure 500 {object} map[string]interface{} "服务器内部错误"
// @Router /api/v1/auth/logout [post]
func (h *AuthHandler) Logout(c *gin.Context) {
	token, exists := c.Get("token")
	if !exists {
		_ = c.Error(errors.New("未找到认证令牌"))
		return
	}

	// 调用服务层
	err := h.authService.Logout(c, token.(string))
	if err != nil {
		_ = c.Error(err)
		return
	}

	// 返回成功响应
	xResult.Success(c, "登出成功")
}

// GetUserInfo 获取当前用户信息
// @Summary 获取用户信息
// @Description 获取当前登录用户的详细信息
// @Tags 认证管理
// @Accept json
// @Produce json
// @Security Bearer
// @Success 200 {object} response.AuthUserInfoResponse "用户信息"
// @Failure 401 {object} map[string]interface{} "未认证"
// @Failure 500 {object} map[string]interface{} "服务器内部错误"
// @Router /api/v1/auth/user [get]
func (h *AuthHandler) GetUserInfo(c *gin.Context) {
	user, exists := ctxUtil.GetUserFromContext(c)
	if !exists {
		_ = c.Error(errors.New("用户信息获取失败"))
		return
	}

	// 调用服务层
	userInfo, err := h.authService.GetUserInfo(c, user.UserUUID)
	if err != nil {
		_ = c.Error(err)
		return
	}

	// 返回成功响应
	resp := response.AuthUserInfoResponse{SystemUserDetailDTO: *userInfo}
	xResult.SuccessHasData(c, "获取成功", resp)
}

// ChangePassword 修改密码
// @Summary 修改密码
// @Description 修改当前用户的登录密码
// @Tags 认证管理
// @Accept json
// @Produce json
// @Security Bearer
// @Param request body request.AuthPasswordChangeReq true "修改密码请求"
// @Success 200 {object} response.AuthPasswordChangeResponse "修改成功"
// @Failure 400 {object} map[string]interface{} "请求参数错误"
// @Failure 401 {object} map[string]interface{} "未认证或旧密码错误"
// @Failure 500 {object} map[string]interface{} "服务器内部错误"
// @Router /api/v1/auth/password/change [post]
func (h *AuthHandler) ChangePassword(c *gin.Context) {
	var req request.AuthPasswordChangeReq

	// 绑定请求数据
	bindErr := c.ShouldBindJSON(&req)
	if bindErr != nil {
		xValid.HandleValidationError(c, bindErr)
		return
	}

	user, exists := ctxUtil.GetUserFromContext(c)
	if !exists {
		_ = c.Error(errors.New("用户信息获取失败"))
		return
	}

	// 调用服务层
	err := h.authService.ChangePassword(c, user.UserUUID, &req)
	if err != nil {
		_ = c.Error(err)
		return
	}

	// 返回成功响应
	xResult.Success(c, "密码修改成功")
}

// ResetPassword 重置密码
// @Summary 重置密码
// @Description 通过邮箱重置用户密码
// @Tags 认证管理
// @Accept json
// @Produce json
// @Param request body request.AuthPasswordResetReq true "重置密码请求"
// @Success 200 {object} response.AuthPasswordResetResponse "重置成功"
// @Failure 400 {object} map[string]interface{} "请求参数错误"
// @Failure 404 {object} map[string]interface{} "邮箱不存在"
// @Failure 500 {object} map[string]interface{} "服务器内部错误"
// @Router /api/v1/auth/password/reset [post]
func (h *AuthHandler) ResetPassword(c *gin.Context) {
	var req request.AuthPasswordResetReq

	// 绑定请求数据
	bindErr := c.ShouldBindJSON(&req)
	if bindErr != nil {
		xValid.HandleValidationError(c, bindErr)
		return
	}

	// 调用服务层
	err := h.authService.ResetPassword(c, &req)
	if err != nil {
		_ = c.Error(err)
		return
	}

	// 返回成功响应
	xResult.Success(c, "密码重置成功，新密码已发送到邮箱")
}
