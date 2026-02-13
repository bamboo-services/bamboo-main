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
	apiAuth "github.com/bamboo-services/bamboo-main/api/auth"
	logic "github.com/bamboo-services/bamboo-main/internal/logic"
	ctxUtil "github.com/bamboo-services/bamboo-main/pkg/util/ctx"

	xError "github.com/bamboo-services/bamboo-base-go/error"
	xResult "github.com/bamboo-services/bamboo-base-go/result"
	xValid "github.com/bamboo-services/bamboo-base-go/validator"
	"github.com/gin-gonic/gin"
)

// AuthHandler 认证处理器
type AuthHandler struct {
	authLogic *logic.AuthLogic
}

// NewAuthHandler 创建认证处理器
func NewAuthHandler() *AuthHandler {
	return &AuthHandler{
		authLogic: logic.NewAuthLogic(),
	}
}

// Login 用户登录
// @Summary 用户登录
// @Description 管理员用户登录，返回用户信息、访问令牌及Token时间信息
// @Tags 认证管理
// @Accept json
// @Produce json
// @Param request body apiAuth.LoginRequest true "登录请求"
// @Success 200 {object} apiAuth.LoginResponse "登录成功，包含用户信息、Token及时间信息"
// @Failure 400 {object} map[string]interface{} "请求参数错误"
// @Failure 401 {object} map[string]interface{} "认证失败"
// @Failure 500 {object} map[string]interface{} "服务器内部错误"
// @Router /api/v1/auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	var req apiAuth.LoginRequest

	// 绑定请求数据
	bindErr := c.ShouldBindJSON(&req)
	if bindErr != nil {
		xValid.HandleValidationError(c, bindErr)
		return
	}

	// 调用服务层
	user, token, createdAt, expiredAt, err := h.authLogic.Login(c, &req)
	if err != nil {
		_ = c.Error(err)
		return
	}

	// 返回成功响应
	resp := apiAuth.LoginResponse{
		User:      *user,
		Token:     token,
		CreatedAt: *createdAt,
		ExpiredAt: *expiredAt,
	}
	xResult.SuccessHasData(c, "登录成功", resp)
}

// Register 用户注册
// @Summary 用户注册
// @Description 注册新用户账户，注册成功后自动登录并返回访问令牌
// @Tags 认证管理
// @Accept json
// @Produce json
// @Param request body apiAuth.RegisterRequest true "注册请求"
// @Success 200 {object} apiAuth.RegisterResponse "注册成功，包含用户信息、Token及时间信息"
// @Failure 400 {object} map[string]interface{} "请求参数错误（用户名或邮箱已存在）"
// @Failure 500 {object} map[string]interface{} "服务器内部错误"
// @Router /api/v1/auth/register [post]
func (h *AuthHandler) Register(c *gin.Context) {
	var req apiAuth.RegisterRequest

	// 绑定请求数据
	bindErr := c.ShouldBindJSON(&req)
	if bindErr != nil {
		xValid.HandleValidationError(c, bindErr)
		return
	}

	// 调用服务层
	user, token, createdAt, expiredAt, err := h.authLogic.Register(c, &req)
	if err != nil {
		_ = c.Error(err)
		return
	}

	// 返回成功响应
	resp := apiAuth.RegisterResponse{
		User:      *user,
		Token:     token,
		CreatedAt: *createdAt,
		ExpiredAt: *expiredAt,
	}
	xResult.SuccessHasData(c, "注册成功", resp)
}

// Logout 用户登出
// @Summary 用户登出
// @Description 注销当前登录会话
// @Tags 认证管理
// @Accept json
// @Produce json
// @Security Bearer
// @Success 200 {object} apiAuth.LogoutResponse "登出成功"
// @Failure 401 {object} map[string]interface{} "未认证"
// @Failure 500 {object} map[string]interface{} "服务器内部错误"
// @Router /api/v1/auth/logout [patch]
func (h *AuthHandler) Logout(c *gin.Context) {
	token, exists := c.Get("token")
	if !exists {
		_ = c.Error(xError.NewError(c, xError.Unauthorized, "未找到认证令牌", false))
		return
	}

	// 调用服务层
	err := h.authLogic.Logout(c, token.(string))
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
// @Success 200 {object} apiAuth.UserInfoResponse "用户信息"
// @Failure 401 {object} map[string]interface{} "未认证"
// @Failure 500 {object} map[string]interface{} "服务器内部错误"
// @Router /api/v1/auth/user [get]
func (h *AuthHandler) GetUserInfo(c *gin.Context) {
	userUUID, exists := ctxUtil.GetUserID(c)
	if !exists {
		_ = c.Error(xError.NewError(c, xError.Unauthorized, "用户信息获取失败", false))
		return
	}

	// 调用服务层
	userInfo, err := h.authLogic.GetUserInfo(c, userUUID)
	if err != nil {
		_ = c.Error(err)
		return
	}

	// 返回成功响应
	resp := apiAuth.UserInfoResponse{SystemUserDetailDTO: *userInfo}
	xResult.SuccessHasData(c, "获取成功", resp)
}

// ChangePassword 修改密码
// @Summary 修改密码
// @Description 修改当前用户的登录密码
// @Tags 认证管理
// @Accept json
// @Produce json
// @Security Bearer
// @Param request body apiAuth.PasswordChangeRequest true "修改密码请求"
// @Success 200 {object} apiAuth.PasswordChangeResponse "修改成功"
// @Failure 400 {object} map[string]interface{} "请求参数错误"
// @Failure 401 {object} map[string]interface{} "未认证或旧密码错误"
// @Failure 500 {object} map[string]interface{} "服务器内部错误"
// @Router /api/v1/auth/password/change [put]
func (h *AuthHandler) ChangePassword(c *gin.Context) {
	var req apiAuth.PasswordChangeRequest

	// 绑定请求数据
	bindErr := c.ShouldBindJSON(&req)
	if bindErr != nil {
		xValid.HandleValidationError(c, bindErr)
		return
	}

	userUUID, exists := ctxUtil.GetUserID(c)
	if !exists {
		_ = c.Error(xError.NewError(c, xError.Unauthorized, "用户信息获取失败", false))
		return
	}

	// 调用服务层
	err := h.authLogic.ChangePassword(c, userUUID, &req)
	if err != nil {
		_ = c.Error(err)
		return
	}

	// 返回成功响应
	xResult.Success(c, "密码修改成功")
}

// ResetPassword 重置密码
// @Summary 重置密码
// @Description 通过邮箱重置用户密码，发送重置链接到邮箱
// @Tags 认证管理
// @Accept json
// @Produce json
// @Param request body apiAuth.PasswordResetRequest true "重置密码请求"
// @Success 200 {object} apiAuth.PasswordResetResponse "重置链接已发送"
// @Failure 400 {object} map[string]interface{} "请求参数错误"
// @Failure 404 {object} map[string]interface{} "邮箱不存在"
// @Failure 500 {object} map[string]interface{} "服务器内部错误"
// @Router /api/v1/auth/password/reset [patch]
func (h *AuthHandler) ResetPassword(c *gin.Context) {
	var req apiAuth.PasswordResetRequest

	// 绑定请求数据
	bindErr := c.ShouldBindJSON(&req)
	if bindErr != nil {
		xValid.HandleValidationError(c, bindErr)
		return
	}

	// 调用服务层
	err := h.authLogic.ResetPassword(c, &req)
	if err != nil {
		_ = c.Error(err)
		return
	}

	// 返回成功响应
	xResult.Success(c, "重置链接已发送到您的邮箱，请在1小时内完成密码重置")
}

// VerifyEmail 验证邮箱
// @Summary 验证邮箱
// @Description 通过邮箱中的验证链接验证用户邮箱
// @Tags 认证管理
// @Accept json
// @Produce json
// @Param token query string true "验证Token"
// @Success 200 {object} map[string]interface{} "验证成功"
// @Failure 400 {object} map[string]interface{} "Token无效或已过期"
// @Failure 500 {object} map[string]interface{} "服务器内部错误"
// @Router /api/v1/auth/verify-email [get]
func (h *AuthHandler) VerifyEmail(c *gin.Context) {
	var req apiAuth.VerifyEmailRequest

	// 绑定请求数据
	bindErr := c.ShouldBindQuery(&req)
	if bindErr != nil {
		xValid.HandleValidationError(c, bindErr)
		return
	}

	// 调用服务层
	err := h.authLogic.VerifyEmail(c, &req)
	if err != nil {
		_ = c.Error(err)
		return
	}

	// 返回成功响应
	xResult.Success(c, "邮箱验证成功")
}

// VerifyResetToken 验证重置密码Token
// @Summary 验证重置密码Token
// @Description 检查密码重置链接是否有效
// @Tags 认证管理
// @Accept json
// @Produce json
// @Param token query string true "重置Token"
// @Success 200 {object} map[string]interface{} "Token有效"
// @Failure 400 {object} map[string]interface{} "Token无效或已过期"
// @Failure 500 {object} map[string]interface{} "服务器内部错误"
// @Router /api/v1/auth/reset-password [get]
func (h *AuthHandler) VerifyResetToken(c *gin.Context) {
	var req apiAuth.VerifyResetTokenRequest

	// 绑定请求数据
	bindErr := c.ShouldBindQuery(&req)
	if bindErr != nil {
		xValid.HandleValidationError(c, bindErr)
		return
	}

	// 调用服务层
	valid, err := h.authLogic.VerifyResetToken(c, &req)
	if err != nil {
		_ = c.Error(err)
		return
	}

	if !valid {
		_ = c.Error(xError.NewError(c, xError.BadRequest, "重置链接无效或已过期", false))
		return
	}

	// 返回成功响应
	xResult.Success(c, "重置链接有效，请设置新密码")
}

// ConfirmResetPassword 确认重置密码
// @Summary 确认重置密码
// @Description 通过重置Token设置新密码
// @Tags 认证管理
// @Accept json
// @Produce json
// @Param request body apiAuth.ConfirmResetPasswordRequest true "确认重置密码请求"
// @Success 200 {object} map[string]interface{} "密码重置成功"
// @Failure 400 {object} map[string]interface{} "Token无效或已过期"
// @Failure 500 {object} map[string]interface{} "服务器内部错误"
// @Router /api/v1/auth/reset-password [post]
func (h *AuthHandler) ConfirmResetPassword(c *gin.Context) {
	var req apiAuth.ConfirmResetPasswordRequest

	// 绑定请求数据
	bindErr := c.ShouldBindJSON(&req)
	if bindErr != nil {
		xValid.HandleValidationError(c, bindErr)
		return
	}

	// 调用服务层
	err := h.authLogic.ConfirmResetPassword(c, &req)
	if err != nil {
		_ = c.Error(err)
		return
	}

	// 返回成功响应
	xResult.Success(c, "密码重置成功，请使用新密码登录")
}
