// --------------------------------------------------------------------------------
// Copyright (c) 2016-NOW(至今) 筱锋
// Author: 筱锋「xiao_lfeng」(https://www.x-lf.com)
// --------------------------------------------------------------------------------
// 许可证声明：版权所有 (c) 2016-2026 筱锋。保留所有权利。
// 有关MIT许可证的更多信息，请查看项目根目录下的LICENSE文件或访问：
// https://opensource.org/licenses/MIT
// --------------------------------------------------------------------------------

package logic

import (
	"context"
	"fmt"
	"time"

	xSnowflake "github.com/bamboo-services/bamboo-base-go/common/snowflake"
	xUtil "github.com/bamboo-services/bamboo-base-go/common/utility"
	apiAuth "github.com/bamboo-services/bamboo-main/api/auth"
	"github.com/bamboo-services/bamboo-main/internal/entity"
	logcHelper "github.com/bamboo-services/bamboo-main/internal/logic/helper"
	"github.com/bamboo-services/bamboo-main/internal/repository"

	"crypto/rand"

	xError "github.com/bamboo-services/bamboo-base-go/common/error"
	xLog "github.com/bamboo-services/bamboo-base-go/common/log"
	xCtxUtil "github.com/bamboo-services/bamboo-base-go/major/utility/context"
	xAsync "github.com/bamboo-services/bamboo-base-go/plugins/async"
	bSdkLogic "github.com/phalanx-labs/beacon-sso-sdk/logic"
)

type authRepo struct {
	user  *repository.SystemUserRepo
	token *repository.TokenRepo
}

// AuthLogic 认证业务逻辑
type AuthLogic struct {
	logic
	SessionService *logcHelper.SessionLogic
	repo           authRepo
}

// NewAuthLogic 创建 AuthLogic 实例，从上下文获取数据库与缓存管理器并初始化会话与认证仓储依赖。
func NewAuthLogic(ctx context.Context) *AuthLogic {
	db := xCtxUtil.MustGetDB(ctx)
	m := xCtxUtil.MustGetCacheManager(ctx)

	return &AuthLogic{
		logic: logic{
			db:    db,
			cache: m,
			log:   xLog.WithName(xLog.NamedLOGC, "AuthLogic"),
		},
		SessionService: logcHelper.NewSessionLogic(m),
		repo: authRepo{
			user:  repository.NewSystemUserRepo(db, m),
			token: repository.NewTokenRepo(m),
		},
	}
}

// Login 用户登录
func (a *AuthLogic) Login(ctx context.Context, req *apiAuth.LoginRequest, meta logcHelper.SessionMeta) (*entity.SystemUser, string, *time.Time, *time.Time, *xError.Error) {
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
	if xErr := a.SessionService.CreateUserSession(ctx, user, token, meta); xErr != nil {
		return nil, "", nil, nil, xErr
	}

	// 更新最后登录时间
	xErr = a.repo.user.UpdateLastLoginByID(ctx, user.ID, &now)
	if xErr != nil {
		// 记录错误但不影响登录
		a.log.Error(ctx, fmt.Sprintf("更新最后登录时间失败: %v", xErr))
	}
	return user, token, &now, &expireAt, nil
}

// Register 用户注册
func (a *AuthLogic) Register(ctx context.Context, req *apiAuth.RegisterRequest, meta logcHelper.SessionMeta) (*entity.SystemUser, string, *time.Time, *time.Time, *xError.Error) {
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

	if xErr := a.SessionService.CreateUserSession(ctx, &newUser, token, meta); xErr != nil {
		return nil, "", nil, nil, xErr
	}

	// 异步发送邮箱验证邮件（xAsync 解耦请求上下文，不阻断主流程）
	xAsync.Async(ctx, func(asyncCtx context.Context) {
		a.sendEmailVerification(asyncCtx, &newUser)
	}, xAsync.WithName("MAIL"))
	return &newUser, token, &now, &expireAt, nil
}

// Logout 用户登出
func (a *AuthLogic) Logout(ctx context.Context, token string) *xError.Error {
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
func (a *AuthLogic) ChangePassword(ctx context.Context, userID xSnowflake.SnowflakeID, req *apiAuth.PasswordChangeRequest) *xError.Error {
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
func (a *AuthLogic) ResetPassword(ctx context.Context, req *apiAuth.PasswordResetRequest) *xError.Error {
	user, found, xErr := a.repo.user.GetByEmail(ctx, req.Email)
	if xErr != nil {
		return xErr
	}
	if !found {
		return xError.NewError(ctx, xError.NotFound, "邮箱不存在", false)
	}

	// 生成重置 Token（32位随机字符串）并存储（1小时过期）
	resetToken := generateRandomString(32)
	if xErr := a.repo.token.SavePasswordResetToken(ctx, resetToken, user.ID); xErr != nil {
		return xErr
	}

	// 发送重置密码邮件（xAsync 解耦请求上下文，不阻断主流程）
	xAsync.Async(ctx, func(asyncCtx context.Context) {
		a.sendPasswordResetEmail(asyncCtx, user, resetToken)
	}, xAsync.WithName("MAIL"))

	return nil
}

// GetUserInfo 获取用户信息
func (a *AuthLogic) GetUserInfo(ctx context.Context, userID xSnowflake.SnowflakeID) (*entity.SystemUser, *xError.Error) {
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
func (a *AuthLogic) UpdateLastLogin(ctx context.Context, userID xSnowflake.SnowflakeID) *xError.Error {
	now := time.Now()
	return a.repo.user.UpdateLastLoginByID(ctx, userID, &now)
}

// ValidateToken 验证令牌
func (a *AuthLogic) ValidateToken(ctx context.Context, token string) (*entity.SystemUser, *xError.Error) {
	// 这个方法主要通过中间件来处理，这里提供一个备用实现
	// 实际项目中可以根据需要实现更复杂的验证逻辑
	return nil, xError.NewError(ctx, xError.OperationInvalid, "请通过认证中间件验证令牌", false)
}

// generateRandomString 生成随机字符串
func generateRandomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

	b := make([]byte, length)
	if _, err := rand.Read(b); err != nil {
		// crypto/rand.Read 在正常环境下几乎不会失败；失败时无法安全继续
		panic(fmt.Sprintf("generateRandomString: rand.Read failed: %v", err))
	}

	result := make([]byte, length)
	for i := range b {
		result[i] = charset[b[i]%byte(len(charset))]
	}

	return string(result)
}

// VerifyEmail 验证邮箱
func (a *AuthLogic) VerifyEmail(ctx context.Context, req *apiAuth.VerifyEmailRequest) *xError.Error {
	userID, found, xErr := a.repo.token.ConsumeEmailVerifyToken(ctx, req.Token)
	if xErr != nil {
		return xErr
	}
	if !found {
		a.log.Warn(ctx, fmt.Sprintf("邮箱验证Token无效或已过期: %s", req.Token))
		return xError.NewError(ctx, xError.BadRequest, "验证链接无效或已过期", false)
	}

	// 更新用户邮箱验证状态
	if _, xErr := a.repo.user.UpdateFieldsByID(ctx, userID, map[string]any{"email_verify": true}); xErr != nil {
		return xErr
	}

	a.log.Info(ctx, fmt.Sprintf("用户 %d 邮箱验证成功", userID))
	return nil
}

// VerifyResetToken 验证重置密码Token
//
// 令牌无效或已过期时直接返回业务错误，handler 无需再做有效性判断。
func (a *AuthLogic) VerifyResetToken(ctx context.Context, req *apiAuth.VerifyResetTokenRequest) *xError.Error {
	exists, xErr := a.repo.token.ExistsPasswordResetToken(ctx, req.Token)
	if xErr != nil {
		return xErr
	}
	if !exists {
		return xError.NewError(ctx, xError.BadRequest, "重置链接无效或已过期", false)
	}
	return nil
}

// ConfirmResetPassword 确认重置密码
func (a *AuthLogic) ConfirmResetPassword(ctx context.Context, req *apiAuth.ConfirmResetPasswordRequest) *xError.Error {
	userID, found, xErr := a.repo.token.ConsumePasswordResetToken(ctx, req.Token)
	if xErr != nil {
		return xErr
	}
	if !found {
		a.log.Warn(ctx, fmt.Sprintf("密码重置Token无效或已过期: %s", req.Token))
		return xError.NewError(ctx, xError.BadRequest, "重置链接无效或已过期", false)
	}

	// 加密新密码
	hashedPassword, err := xUtil.Password().EncryptString(req.NewPassword)
	if err != nil {
		return xError.NewError(ctx, xError.ServerInternalError, "密码加密失败", false, err)
	}

	// 更新用户密码
	if xErr := a.repo.user.UpdatePasswordByID(ctx, userID, hashedPassword); xErr != nil {
		return xErr
	}

	a.log.Info(ctx, fmt.Sprintf("用户 %d 密码重置成功", userID))
	return nil
}

// sendEmailVerification 发送邮箱验证邮件
//
// 此函数应在 xAsync 异步任务中调用，ctx 为解耦后的独立上下文，不会阻断主流程
func (a *AuthLogic) sendEmailVerification(ctx context.Context, user *entity.SystemUser) {
	logger := xLog.WithName(xLog.NamedLOGC, "MAIL")

	// 生成验证 Token 并存储（24小时过期）
	verifyToken := generateRandomString(32)
	if xErr := a.repo.token.SaveEmailVerifyToken(ctx, verifyToken, user.ID); xErr != nil {
		logger.Warn(ctx, fmt.Sprintf("保存验证Token失败: %v", xErr))
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
	}

	// 发送邮件
	mailLogic := NewMailLogic()
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
// 此函数应在 xAsync 异步任务中调用，ctx 为解耦后的独立上下文，不会阻断主流程
func (a *AuthLogic) sendPasswordResetEmail(ctx context.Context, user *entity.SystemUser, resetToken string) {
	logger := xLog.WithName(xLog.NamedLOGC, "MAIL")

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
	}

	// 发送邮件
	mailLogic := NewMailLogic()
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
