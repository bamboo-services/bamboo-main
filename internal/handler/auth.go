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
	logcHelper "github.com/bamboo-services/bamboo-main/internal/logic/helper"
	ctxUtil "github.com/bamboo-services/bamboo-main/pkg/util/ctx"
	"github.com/bamboo-services/bamboo-main/pkg/util/netUtil"

	xError "github.com/bamboo-services/bamboo-base-go/common/error"
	xResult "github.com/bamboo-services/bamboo-base-go/major/result"
	xValid "github.com/bamboo-services/bamboo-base-go/major/validator"
	"github.com/gin-gonic/gin"
	bSdkUtil "github.com/phalanx-labs/beacon-sso-sdk/utility"
)

// Login 用户登录（账号密码）
//
// @Summary [公开] 用户登录
// @Description 使用用户名/邮箱 + 密码登录，返回用户信息、访问令牌及Token时间信息
// @Tags 认证接口
// @Accept json
// @Produce json
// @Param request body apiAuth.LoginRequest true "登录请求"
// @Success 200 {object} xBase.BaseResponse{data=apiAuth.LoginResponse} "登录成功，包含用户信息、Token及时间信息"
// @Failure 400 {object} xBase.BaseResponse "请求参数错误"
// @Failure 401 {object} xBase.BaseResponse "用户名或密码错误"
// @Failure 500 {object} xBase.BaseResponse "服务器内部错误"
// @Router /api/v1/auth/login [POST]
func (h *AuthHandler) Login(c *gin.Context) {
	var req apiAuth.LoginRequest

	// 绑定请求数据
	bindErr := c.ShouldBindJSON(&req)
	if bindErr != nil {
		xValid.HandleValidationError(c, bindErr)
		return
	}

	// 调用服务层
	meta := logcHelper.SessionMeta{
		ClientIP:  netUtil.GetClientIP(c),
		UserAgent: c.GetHeader("User-Agent"),
	}
	user, token, createdAt, expiredAt, err := h.service.authLogic.Login(c.Request.Context(), &req, meta)
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

// OAuthLogin 用户登录（SSO OAuth）
//
// @Summary [用户] SSO OAuth 登录
// @Description 携带 SSO 访问令牌换取本地会话，返回用户信息、访问令牌及Token时间信息
// @Tags 认证接口
// @Accept json
// @Produce json
// @Security Bearer
// @Success 200 {object} xBase.BaseResponse{data=apiAuth.LoginResponse} "登录成功，包含用户信息、Token及时间信息"
// @Failure 401 {object} xBase.BaseResponse "认证失败"
// @Failure 500 {object} xBase.BaseResponse "服务器内部错误"
// @Router /api/v1/auth/oauth/login [POST]
func (h *AuthHandler) OAuthLogin(c *gin.Context) {
	accessToken := bSdkUtil.GetAuthorization(c)
	userinfo, xErr := h.service.oauthLogic.Userinfo(c.Request.Context(), accessToken)
	if xErr != nil {
		_ = c.Error(xErr)
		return
	}

	meta := logcHelper.SessionMeta{
		ClientIP:  netUtil.GetClientIP(c),
		UserAgent: c.GetHeader("User-Agent"),
	}
	user, token, createdAt, expiredAt, err := h.service.authLogic.LoginByOAuth(c.Request.Context(), userinfo, accessToken, meta)
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
//
// @Summary [用户] 用户注册
// @Description 注册新用户账户，注册成功后自动登录并返回访问令牌
// @Tags 认证接口
// @Accept json
// @Produce json
// @Param request body apiAuth.RegisterRequest true "注册请求"
// @Success 200 {object} xBase.BaseResponse{data=apiAuth.RegisterResponse} "注册成功，包含用户信息、Token及时间信息"
// @Failure 400 {object} xBase.BaseResponse "请求参数错误（用户名或邮箱已存在）"
// @Failure 500 {object} xBase.BaseResponse "服务器内部错误"
// @Router /api/v1/auth/register [POST]
func (h *AuthHandler) Register(c *gin.Context) {
	var req apiAuth.RegisterRequest

	// 绑定请求数据
	bindErr := c.ShouldBindJSON(&req)
	if bindErr != nil {
		xValid.HandleValidationError(c, bindErr)
		return
	}

	// 调用服务层
	meta := logcHelper.SessionMeta{
		ClientIP:  netUtil.GetClientIP(c),
		UserAgent: c.GetHeader("User-Agent"),
	}
	user, token, createdAt, expiredAt, err := h.service.authLogic.Register(c.Request.Context(), &req, meta)
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

// SendRegisterCode 发送注册邮箱验证码
//
// @Summary [公开] 发送注册验证码
// @Description 向指定邮箱发送 6 位注册验证码，验证码 10 分钟内有效，60 秒内不可重复发送
// @Tags 认证接口
// @Accept json
// @Produce json
// @Param request body apiAuth.RegisterCodeRequest true "发送验证码请求"
// @Success 200 {object} xBase.BaseResponse "验证码已发送"
// @Failure 400 {object} xBase.BaseResponse "请求参数错误（邮箱已被注册）"
// @Failure 500 {object} xBase.BaseResponse "服务器内部错误或发送过于频繁"
// @Router /api/v1/auth/register/code [POST]
func (h *AuthHandler) SendRegisterCode(c *gin.Context) {
	var req apiAuth.RegisterCodeRequest

	// 绑定请求数据
	bindErr := c.ShouldBindJSON(&req)
	if bindErr != nil {
		xValid.HandleValidationError(c, bindErr)
		return
	}

	// 调用服务层
	err := h.service.authLogic.SendRegisterCode(c.Request.Context(), &req)
	if err != nil {
		_ = c.Error(err)
		return
	}

	// 返回成功响应
	xResult.Success(c, "验证码已发送到您的邮箱，请在10分钟内完成注册")
}

// Logout 用户登出
//
// @Summary [用户] 用户登出
// @Description 注销当前登录会话
// @Tags 认证接口
// @Accept json
// @Produce json
// @Security Bearer
// @Success 200 {object} xBase.BaseResponse "登出成功"
// @Failure 401 {object} xBase.BaseResponse "未认证"
// @Failure 500 {object} xBase.BaseResponse "服务器内部错误"
// @Router /api/v1/auth/logout [PATCH]
func (h *AuthHandler) Logout(c *gin.Context) {
	token := bSdkUtil.GetAuthorization(c)
	if token == "" {
		if ctxToken, exists := c.Get("token"); exists {
			token, _ = ctxToken.(string)
		}
	}

	if token == "" {
		_ = c.Error(xError.NewError(c, xError.Unauthorized, "未找到认证令牌", false))
		return
	}

	err := h.service.authLogic.Logout(c.Request.Context(), token)
	if err != nil {
		_ = c.Error(err)
		return
	}

	// 返回成功响应
	xResult.Success(c, "登出成功")
}

// GetUserInfo 获取当前用户信息
//
// @Summary [用户] 获取用户信息
// @Description 获取当前登录用户的详细信息
// @Tags 认证接口
// @Accept json
// @Produce json
// @Security Bearer
// @Success 200 {object} xBase.BaseResponse{data=apiAuth.UserInfoResponse} "用户信息"
// @Failure 401 {object} xBase.BaseResponse "未认证"
// @Failure 500 {object} xBase.BaseResponse "服务器内部错误"
// @Router /api/v1/auth/user [GET]
func (h *AuthHandler) GetUserInfo(c *gin.Context) {
	userUUID, exists := ctxUtil.GetUserID(c)
	if !exists {
		_ = c.Error(xError.NewError(c, xError.Unauthorized, "用户信息获取失败", false))
		return
	}

	// 调用服务层
	userInfo, err := h.service.authLogic.GetUserInfo(c.Request.Context(), userUUID)
	if err != nil {
		_ = c.Error(err)
		return
	}

	// 返回成功响应
	resp := apiAuth.UserInfoResponse{User: *userInfo}
	xResult.SuccessHasData(c, "获取成功", resp)
}

// UpdateProfile 更新用户资料
//
// @Summary [用户] 更新用户资料
// @Description 更新当前用户的用户名、昵称与头像
// @Tags 认证接口
// @Accept json
// @Produce json
// @Security Bearer
// @Param request body apiAuth.UpdateProfileRequest true "更新资料请求"
// @Success 200 {object} xBase.BaseResponse{data=apiAuth.UserInfoResponse} "更新成功"
// @Failure 400 {object} xBase.BaseResponse "请求参数错误"
// @Failure 401 {object} xBase.BaseResponse "未认证"
// @Failure 500 {object} xBase.BaseResponse "服务器内部错误"
// @Router /api/v1/user/profile [PUT]
func (h *AuthHandler) UpdateProfile(c *gin.Context) {
	var req apiAuth.UpdateProfileRequest

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
	user, err := h.service.authLogic.UpdateProfile(c.Request.Context(), userUUID, &req)
	if err != nil {
		_ = c.Error(err)
		return
	}

	// 返回成功响应
	resp := apiAuth.UserInfoResponse{User: *user}
	xResult.SuccessHasData(c, "资料更新成功", resp)
}

// ChangePassword 修改密码
//
// @Summary [用户] 修改密码
// @Description 修改当前用户的登录密码
// @Tags 认证接口
// @Accept json
// @Produce json
// @Security Bearer
// @Param request body apiAuth.PasswordChangeRequest true "修改密码请求"
// @Success 200 {object} xBase.BaseResponse "修改成功"
// @Failure 400 {object} xBase.BaseResponse "请求参数错误"
// @Failure 401 {object} xBase.BaseResponse "未认证或旧密码错误"
// @Failure 500 {object} xBase.BaseResponse "服务器内部错误"
// @Router /api/v1/auth/password/change [PUT]
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
	err := h.service.authLogic.ChangePassword(c.Request.Context(), userUUID, &req)
	if err != nil {
		_ = c.Error(err)
		return
	}

	// 返回成功响应
	xResult.Success(c, "密码修改成功")
}

// ResetPassword 重置密码
//
// @Summary [用户] 重置密码
// @Description 通过邮箱重置用户密码，发送重置链接到邮箱
// @Tags 认证接口
// @Accept json
// @Produce json
// @Param request body apiAuth.PasswordResetRequest true "重置密码请求"
// @Success 200 {object} xBase.BaseResponse "重置链接已发送"
// @Failure 400 {object} xBase.BaseResponse "请求参数错误"
// @Failure 404 {object} xBase.BaseResponse "邮箱不存在"
// @Failure 500 {object} xBase.BaseResponse "服务器内部错误"
// @Router /api/v1/auth/password/reset [PATCH]
func (h *AuthHandler) ResetPassword(c *gin.Context) {
	var req apiAuth.PasswordResetRequest

	// 绑定请求数据
	bindErr := c.ShouldBindJSON(&req)
	if bindErr != nil {
		xValid.HandleValidationError(c, bindErr)
		return
	}

	// 调用服务层
	err := h.service.authLogic.ResetPassword(c.Request.Context(), &req)
	if err != nil {
		_ = c.Error(err)
		return
	}

	// 返回成功响应
	xResult.Success(c, "重置链接已发送到您的邮箱，请在1小时内完成密码重置")
}

// VerifyEmail 验证邮箱
//
// @Summary [用户] 验证邮箱
// @Description 通过邮箱中的验证链接验证用户邮箱
// @Tags 认证接口
// @Accept json
// @Produce json
// @Param token query string true "验证Token"
// @Success 200 {object} xBase.BaseResponse "验证成功"
// @Failure 400 {object} xBase.BaseResponse "Token无效或已过期"
// @Failure 500 {object} xBase.BaseResponse "服务器内部错误"
// @Router /api/v1/auth/verify-email [GET]
func (h *AuthHandler) VerifyEmail(c *gin.Context) {
	var req apiAuth.VerifyEmailRequest

	// 绑定请求数据
	bindErr := c.ShouldBindQuery(&req)
	if bindErr != nil {
		xValid.HandleValidationError(c, bindErr)
		return
	}

	// 调用服务层
	err := h.service.authLogic.VerifyEmail(c.Request.Context(), &req)
	if err != nil {
		_ = c.Error(err)
		return
	}

	// 返回成功响应
	xResult.Success(c, "邮箱验证成功")
}

// VerifyResetToken 验证重置密码Token
//
// @Summary [用户] 验证重置密码Token
// @Description 检查密码重置链接是否有效
// @Tags 认证接口
// @Accept json
// @Produce json
// @Param token query string true "重置Token"
// @Success 200 {object} xBase.BaseResponse "Token有效"
// @Failure 400 {object} xBase.BaseResponse "Token无效或已过期"
// @Failure 500 {object} xBase.BaseResponse "服务器内部错误"
// @Router /api/v1/auth/reset-password [GET]
func (h *AuthHandler) VerifyResetToken(c *gin.Context) {
	var req apiAuth.VerifyResetTokenRequest

	// 绑定请求数据
	bindErr := c.ShouldBindQuery(&req)
	if bindErr != nil {
		xValid.HandleValidationError(c, bindErr)
		return
	}

	// 调用服务层（令牌无效或已过期时 logic 直接返回业务错误）
	err := h.service.authLogic.VerifyResetToken(c.Request.Context(), &req)
	if err != nil {
		_ = c.Error(err)
		return
	}

	// 返回成功响应
	xResult.Success(c, "重置链接有效，请设置新密码")
}

// ConfirmResetPassword 确认重置密码
//
// @Summary [用户] 确认重置密码
// @Description 通过重置Token设置新密码
// @Tags 认证接口
// @Accept json
// @Produce json
// @Param request body apiAuth.ConfirmResetPasswordRequest true "确认重置密码请求"
// @Success 200 {object} xBase.BaseResponse "密码重置成功"
// @Failure 400 {object} xBase.BaseResponse "Token无效或已过期"
// @Failure 500 {object} xBase.BaseResponse "服务器内部错误"
// @Router /api/v1/auth/reset-password [POST]
func (h *AuthHandler) ConfirmResetPassword(c *gin.Context) {
	var req apiAuth.ConfirmResetPasswordRequest

	// 绑定请求数据
	bindErr := c.ShouldBindJSON(&req)
	if bindErr != nil {
		xValid.HandleValidationError(c, bindErr)
		return
	}

	// 调用服务层
	err := h.service.authLogic.ConfirmResetPassword(c.Request.Context(), &req)
	if err != nil {
		_ = c.Error(err)
		return
	}

	// 返回成功响应
	xResult.Success(c, "密码重置成功，请使用新密码登录")
}
