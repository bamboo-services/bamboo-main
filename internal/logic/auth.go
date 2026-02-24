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

package logic

import (
	"context"
	"fmt"
	"strconv"
	"time"

	xUtil "github.com/bamboo-services/bamboo-base-go/common/utility"
	apiAuth "github.com/bamboo-services/bamboo-main/api/auth"
	"github.com/bamboo-services/bamboo-main/internal/entity"
	logcHelper "github.com/bamboo-services/bamboo-main/internal/logic/helper"
	"github.com/bamboo-services/bamboo-main/internal/repository"
	"github.com/bamboo-services/bamboo-main/pkg/constants"
	ctxUtil "github.com/bamboo-services/bamboo-main/pkg/util/ctx"

	"crypto/rand"

	xError "github.com/bamboo-services/bamboo-base-go/common/error"
	xLog "github.com/bamboo-services/bamboo-base-go/common/log"
	xCtxUtil "github.com/bamboo-services/bamboo-base-go/common/utility/context"
	"github.com/gin-gonic/gin"
	bSdkLogic "github.com/phalanx-labs/beacon-sso-sdk/logic"
)

type authRepo struct {
	user *repository.SystemUserRepo
}

// AuthLogic 认证业务逻辑
type AuthLogic struct {
	logic
	SessionService *logcHelper.SessionLogic
	repo           authRepo
}

func NewAuthLogic(ctx context.Context) *AuthLogic {
	db := xCtxUtil.MustGetDB(ctx)
	rdb := xCtxUtil.MustGetRDB(ctx)

	return &AuthLogic{
		logic: logic{
			db:  db,
			rdb: rdb,
			log: xLog.WithName(xLog.NamedLOGC, "AuthLogic"),
		},
		SessionService: &logcHelper.SessionLogic{},
		repo: authRepo{
			user: repository.NewSystemUserRepo(db, rdb),
		},
	}
}

// Login 用户登录
func (a *AuthLogic) Login(ctx *gin.Context, req *apiAuth.LoginRequest) (*entity.SystemUser, string, *time.Time, *time.Time, *xError.Error) {
	user, found, xErr := a.repo.user.GetByUsernameOrEmail(ctx, req.Username)
	if xErr != nil {
		return nil, "", nil, nil, xErr
	}
	if !found {
		return nil, "", nil, nil, xError.NewError(ctx, xError.LoginFailed, "用户名或密码错误", false)
	}

	// 检查用户状态
	if user.Status == 0 {
		return nil, "", nil, nil, xError.NewError(ctx, xError.Forbidden, "用户已被禁用", false)
	}

	// 验证密码
	if !xUtil.Password().IsValid(req.Password, user.Password) {
		return nil, "", nil, nil, xError.NewError(ctx, xError.LoginFailed, "用户名或密码错误", false)
	}

	// 生成 token
	token := xUtil.Security().GenerateKey()

	// 记录时间信息
	now := time.Now()
	expireAt := now.Add(24 * time.Hour) // 24小时过期

	// 创建用户会话
	err := a.SessionService.CreateUserSession(ctx, user, token)
	if err != nil {
		return nil, "", nil, nil, xError.NewError(ctx, xError.ServerInternalError, "创建用户会话失败", false, err)
	}

	// 更新最后登录时间
	xErr = a.repo.user.UpdateLastLoginByID(ctx, user.ID, &now)
	if xErr != nil {
		// 记录错误但不影响登录
		xLog.WithName(xLog.NamedLOGC, "AUTH").Error(ctx, fmt.Sprintf("更新最后登录时间失败: %v", xErr))
	}
	return user, token, &now, &expireAt, nil
}

// Register 用户注册
func (a *AuthLogic) Register(ctx *gin.Context, req *apiAuth.RegisterRequest) (*entity.SystemUser, string, *time.Time, *time.Time, *xError.Error) {
	exists, xErr := a.repo.user.ExistsByUsername(ctx, req.Username)
	if xErr != nil {
		return nil, "", nil, nil, xErr
	}
	if exists {
		return nil, "", nil, nil, xError.NewError(ctx, xError.ParameterError, "用户名已存在", false)
	}

	exists, xErr = a.repo.user.ExistsByEmailExceptID(ctx, req.Email, 0)
	if xErr != nil {
		return nil, "", nil, nil, xErr
	}
	if exists {
		return nil, "", nil, nil, xError.NewError(ctx, xError.ParameterError, "邮箱已被注册", false)
	}

	hashedPassword, err := xUtil.Password().EncryptString(req.Password)
	if err != nil {
		return nil, "", nil, nil, xError.NewError(ctx, xError.ServerInternalError, "密码加密失败", false, err)
	}

	newUser := entity.SystemUser{
		Username:    req.Username,
		Password:    hashedPassword,
		Email:       req.Email,
		Nickname:    req.Nickname,
		Role:        "user", // 新用户角色为 user
		Status:      1,      // 默认启用
		EmailVerify: false,  // 默认未验证邮箱
	}

	if _, xErr = a.repo.user.Create(ctx, &newUser); xErr != nil {
		return nil, "", nil, nil, xErr
	}

	token := xUtil.Security().GenerateKey()

	now := time.Now()
	expireAt := now.Add(24 * time.Hour) // 24小时过期

	err = a.SessionService.CreateUserSession(ctx, &newUser, token)
	if err != nil {
		return nil, "", nil, nil, xError.NewError(ctx, xError.ServerInternalError, "创建用户会话失败", false, err)
	}

	go a.sendEmailVerification(ctx, &newUser)
	return &newUser, token, &now, &expireAt, nil
}

// Logout 用户登出
func (a *AuthLogic) Logout(ctx *gin.Context, token string) *xError.Error {
	if token == "" {
		return xError.NewError(ctx, xError.ParameterEmpty, "访问令牌不能为空", false)
	}

	oauthLogic := bSdkLogic.NewOAuth(ctx)
	xErr := oauthLogic.Logout(ctx, "access_token", token)
	if xErr != nil {
		return xErr
	}

	return nil
}

// ChangePassword 修改密码
func (a *AuthLogic) ChangePassword(ctx *gin.Context, userID int64, req *apiAuth.PasswordChangeRequest) *xError.Error {
	user, found, xErr := a.repo.user.GetByID(ctx, userID)
	if xErr != nil {
		return xErr
	}
	if !found {
		return xError.NewError(ctx, xError.NotFound, "用户不存在", false)
	}

	// 验证旧密码
	if !xUtil.Password().IsValid(req.OldPassword, user.Password) {
		return xError.NewError(ctx, xError.ParameterError, "旧密码错误", false)
	}

	// 加密新密码
	hashedPassword, err := xUtil.Password().EncryptString(req.NewPassword)
	if err != nil {
		return xError.NewError(ctx, xError.ServerInternalError, "密码加密失败", false, err)
	}

	// 更新密码
	xErr = a.repo.user.UpdatePasswordByID(ctx, userID, hashedPassword)
	if xErr != nil {
		return xErr
	}

	return nil
}

// ResetPassword 重置密码（发送重置链接）
func (a *AuthLogic) ResetPassword(ctx *gin.Context, req *apiAuth.PasswordResetRequest) *xError.Error {
	user, found, xErr := a.repo.user.GetByEmail(ctx, req.Email)
	if xErr != nil {
		return xErr
	}
	if !found {
		return xError.NewError(ctx, xError.NotFound, "邮箱不存在", false)
	}

	// 生成重置 Token（32位随机字符串）
	resetToken := generateRandomString(32)

	// 存储到 Redis（1小时过期）
	rdb := ctxUtil.GetRedisClient(ctx)
	if rdb == nil {
		return xError.NewError(ctx, xError.ServerInternalError, "Redis 客户端不可用", false)
	}

	redisKey := constants.RedisPasswordReset.Get(resetToken).String()
	err := rdb.Set(ctx.Request.Context(), redisKey, user.ID, time.Hour).Err()
	if err != nil {
		return xError.NewError(ctx, xError.ServerInternalError, "保存重置Token失败", false, err)
	}

	// 发送重置密码邮件（异步，不阻断主流程）
	go a.sendPasswordResetEmail(ctx, user, resetToken)

	return nil
}

// GetUserInfo 获取用户信息
func (a *AuthLogic) GetUserInfo(ctx *gin.Context, userID int64) (*entity.SystemUser, *xError.Error) {
	user, found, xErr := a.repo.user.GetByID(ctx, userID)
	if xErr != nil {
		return nil, xErr
	}
	if !found {
		return nil, xError.NewError(ctx, xError.NotFound, "用户不存在", false)
	}
	return user, nil
}

// UpdateLastLogin 更新最后登录时间
func (a *AuthLogic) UpdateLastLogin(ctx *gin.Context, userID int64) *xError.Error {
	now := time.Now()
	return a.repo.user.UpdateLastLoginByID(ctx, userID, &now)
}

// ValidateToken 验证令牌
func (a *AuthLogic) ValidateToken(ctx *gin.Context, token string) (*entity.SystemUser, *xError.Error) {
	// 这个方法主要通过中间件来处理，这里提供一个备用实现
	// 实际项目中可以根据需要实现更复杂的验证逻辑
	return nil, xError.NewError(ctx, xError.OperationInvalid, "请通过认证中间件验证令牌", false)
}

// generateRandomString 生成随机字符串
func generateRandomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

	b := make([]byte, length)
	rand.Read(b)

	result := make([]byte, length)
	for i := range b {
		result[i] = charset[b[i]%byte(len(charset))]
	}

	return string(result)
}

func parseUserID(value string) int64 {
	userID, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		return 0
	}
	return userID
}

// VerifyEmail 验证邮箱
func (a *AuthLogic) VerifyEmail(ctx *gin.Context, req *apiAuth.VerifyEmailRequest) *xError.Error {
	logger := xLog.WithName(xLog.NamedLOGC, "AUTH")

	// 获取 Redis 客户端
	rdb := ctxUtil.GetRedisClient(ctx)
	if rdb == nil {
		return xError.NewError(ctx, xError.ServerInternalError, "Redis 客户端不可用", false)
	}

	// 从 Redis 获取用户 ID
	redisKey := constants.RedisEmailVerify.Get(req.Token).String()
	userIDStr, err := rdb.Get(ctx.Request.Context(), redisKey).Result()
	if err != nil {
		logger.Warn(ctx, fmt.Sprintf("邮箱验证Token无效或已过期: %s", req.Token))
		return xError.NewError(ctx, xError.BadRequest, "验证链接无效或已过期", false)
	}

	// 删除已使用的 Token
	rdb.Del(ctx.Request.Context(), redisKey)

	// 更新用户邮箱验证状态
	_, xErr := a.repo.user.UpdateFieldsByID(ctx, parseUserID(userIDStr), map[string]any{"email_verify": true})
	if xErr != nil {
		return xErr
	}

	logger.Info(ctx, fmt.Sprintf("用户 %s 邮箱验证成功", userIDStr))
	return nil
}

// VerifyResetToken 验证重置密码Token
func (a *AuthLogic) VerifyResetToken(ctx *gin.Context, req *apiAuth.VerifyResetTokenRequest) (bool, *xError.Error) {
	// 获取 Redis 客户端
	rdb := ctxUtil.GetRedisClient(ctx)
	if rdb == nil {
		return false, xError.NewError(ctx, xError.ServerInternalError, "Redis 客户端不可用", false)
	}

	// 检查 Token 是否存在
	redisKey := constants.RedisPasswordReset.Get(req.Token).String()
	exists, err := rdb.Exists(ctx.Request.Context(), redisKey).Result()
	if err != nil {
		return false, xError.NewError(ctx, xError.ServerInternalError, "验证Token失败", false, err)
	}

	return exists > 0, nil
}

// ConfirmResetPassword 确认重置密码
func (a *AuthLogic) ConfirmResetPassword(ctx *gin.Context, req *apiAuth.ConfirmResetPasswordRequest) *xError.Error {
	logger := xLog.WithName(xLog.NamedLOGC, "AUTH")

	// 获取 Redis 客户端
	rdb := ctxUtil.GetRedisClient(ctx)
	if rdb == nil {
		return xError.NewError(ctx, xError.ServerInternalError, "Redis 客户端不可用", false)
	}

	// 从 Redis 获取用户 ID
	redisKey := constants.RedisPasswordReset.Get(req.Token).String()
	userIDStr, err := rdb.Get(ctx.Request.Context(), redisKey).Result()
	if err != nil {
		logger.Warn(ctx, fmt.Sprintf("密码重置Token无效或已过期: %s", req.Token))
		return xError.NewError(ctx, xError.BadRequest, "重置链接无效或已过期", false)
	}

	// 删除已使用的 Token
	rdb.Del(ctx.Request.Context(), redisKey)

	// 加密新密码
	hashedPassword, err := xUtil.Password().EncryptString(req.NewPassword)
	if err != nil {
		return xError.NewError(ctx, xError.ServerInternalError, "密码加密失败", false, err)
	}

	// 更新用户密码
	xErr := a.repo.user.UpdatePasswordByID(ctx, parseUserID(userIDStr), hashedPassword)
	if xErr != nil {
		return xErr
	}

	logger.Info(ctx, fmt.Sprintf("用户 %s 密码重置成功", userIDStr))
	return nil
}

// sendEmailVerification 发送邮箱验证邮件
//
// 此函数应在 goroutine 中异步调用，不会阻断主流程
func (a *AuthLogic) sendEmailVerification(ctx *gin.Context, user *entity.SystemUser) {
	logger := xLog.WithName(xLog.NamedLOGC, "MAIL")

	// 获取配置
	config := ctxUtil.GetConfig(ctx)
	if config == nil {
		logger.Warn(ctx, "无法获取配置，跳过发送邮箱验证邮件")
		return
	}

	// 获取 Redis 客户端
	rdb := ctxUtil.GetRedisClient(ctx)
	if rdb == nil {
		logger.Warn(ctx, "Redis 客户端不可用，跳过发送邮箱验证邮件")
		return
	}

	// 生成验证 Token
	verifyToken := generateRandomString(32)

	// 存储到 Redis（24小时过期）
	redisKey := constants.RedisEmailVerify.Get(verifyToken).String()
	err := rdb.Set(ctx.Request.Context(), redisKey, user.ID, 24*time.Hour).Err()
	if err != nil {
		logger.Warn(ctx, fmt.Sprintf("保存验证Token失败: %v", err))
		return
	}

	// 构建验证链接（TODO: 从配置读取域名前缀）
	verifyLink := fmt.Sprintf("https://localhost/api/v1/auth/verify-email?token=%s", verifyToken)

	// 获取用户昵称
	username := user.Username
	if user.Nickname != nil && *user.Nickname != "" {
		username = *user.Nickname
	}

	// 构建模板变量
	variables := map[string]string{
		"Username":   username,
		"VerifyLink": verifyLink,
		"ExpireTime": "24小时",
		"FromName":   config.Email.FromName,
	}

	// 发送邮件
	mailLogic := &MailLogic{TemplateService: &logcHelper.MailTemplateLogic{}, MaxRetry: 3}
	mailErr := mailLogic.SendWithTemplate(
		ctx,
		"email_verify",
		[]string{user.Email},
		"请验证您的邮箱地址",
		variables,
	)
	if mailErr != nil {
		logger.Warn(ctx, fmt.Sprintf("发送邮箱验证邮件失败: %v", mailErr))
	} else {
		logger.Info(ctx, fmt.Sprintf("已发送邮箱验证邮件到: %s", user.Email))
	}
}

// sendPasswordResetEmail 发送密码重置邮件
//
// 此函数应在 goroutine 中异步调用，不会阻断主流程
func (a *AuthLogic) sendPasswordResetEmail(ctx *gin.Context, user *entity.SystemUser, resetToken string) {
	logger := xLog.WithName(xLog.NamedLOGC, "MAIL")

	// 获取配置
	config := ctxUtil.GetConfig(ctx)
	if config == nil {
		logger.Warn(ctx, "无法获取配置，跳过发送密码重置邮件")
		return
	}

	// 构建重置链接（TODO: 从配置读取域名前缀）
	resetLink := fmt.Sprintf("https://localhost/api/v1/auth/reset-password?token=%s", resetToken)

	// 获取用户昵称
	username := user.Username
	if user.Nickname != nil && *user.Nickname != "" {
		username = *user.Nickname
	}

	// 构建模板变量
	variables := map[string]string{
		"Username":   username,
		"ResetLink":  resetLink,
		"ExpireTime": "1小时",
		"FromName":   config.Email.FromName,
	}

	// 发送邮件
	mailLogic := &MailLogic{TemplateService: &logcHelper.MailTemplateLogic{}, MaxRetry: 3}
	mailErr := mailLogic.SendWithTemplate(
		ctx,
		"password_reset",
		[]string{user.Email},
		"密码重置请求",
		variables,
	)
	if mailErr != nil {
		logger.Warn(ctx, fmt.Sprintf("发送密码重置邮件失败: %v", mailErr))
	} else {
		logger.Info(ctx, fmt.Sprintf("已发送密码重置邮件到: %s", user.Email))
	}
}
