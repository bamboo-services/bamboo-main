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
	"errors"
	"fmt"
	"time"

	"bamboo-main/internal/model/dto"
	"bamboo-main/internal/model/entity"
	"bamboo-main/internal/model/request"
	servHelper "bamboo-main/internal/service/helper"
	"bamboo-main/pkg/constants"
	ctxUtil "bamboo-main/pkg/util/ctx"

	"crypto/rand"

	xError "github.com/bamboo-services/bamboo-base-go/error"
	xUtil "github.com/bamboo-services/bamboo-base-go/utility"
	xCtxUtil "github.com/bamboo-services/bamboo-base-go/utility/ctxutil"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// AuthLogic 认证业务逻辑
type AuthLogic struct {
	SessionService servHelper.ISessionService
}

// Login 用户登录
func (a *AuthLogic) Login(ctx *gin.Context, req *request.AuthLoginReq) (*dto.SystemUserDetailDTO, string, *time.Time, *time.Time, *xError.Error) {
	// 获取数据库连接
	db := xCtxUtil.GetDB(ctx)

	// 查找用户
	var user entity.SystemUser
	err := db.Where("username = ? OR email = ?", req.Username, req.Username).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, "", nil, nil, xError.NewError(ctx, xError.LoginFailed, "用户名或密码错误", false)
		}
		return nil, "", nil, nil, xError.NewError(ctx, xError.DatabaseError, "查询用户失败", false, err)
	}

	// 检查用户状态
	if user.Status == 0 {
		return nil, "", nil, nil, xError.NewError(ctx, xError.Forbidden, "用户已被禁用", false)
	}

	// 验证密码
	if !xUtil.IsPasswordValid(req.Password, user.Password) {
		return nil, "", nil, nil, xError.NewError(ctx, xError.LoginFailed, "用户名或密码错误", false)
	}

	// 生成 token
	token := xUtil.GenerateSecurityKey()

	// 记录时间信息
	now := time.Now()
	expireAt := now.Add(24 * time.Hour) // 24小时过期

	// 创建用户会话
	err = a.SessionService.CreateUserSession(ctx, &user, token)
	if err != nil {
		return nil, "", nil, nil, xError.NewError(ctx, xError.ServerInternalError, "创建用户会话失败", false, err)
	}

	// 更新最后登录时间
	err = db.Model(&user).Update("last_login_at", &now).Error
	if err != nil {
		// 记录错误但不影响登录
		xCtxUtil.GetSugarLogger(ctx, "").Errorf("更新最后登录时间失败: %v", err)
	}

	userDTO := &dto.SystemUserDetailDTO{
		ID:          user.ID,
		Username:    user.Username,
		Email:       user.Email,
		Nickname:    user.Nickname,
		Avatar:      user.Avatar,
		Role:        user.Role,
		Status:      user.Status,
		LastLoginAt: user.LastLoginAt,
		CreatedAt:   user.CreatedAt,
		UpdatedAt:   user.UpdatedAt,
	}
	return userDTO, token, &now, &expireAt, nil
}

// Register 用户注册
func (a *AuthLogic) Register(ctx *gin.Context, req *request.AuthRegisterReq) (*dto.SystemUserDetailDTO, string, *time.Time, *time.Time, *xError.Error) {
	// 获取数据库连接
	db := xCtxUtil.GetDB(ctx)

	// 1. 检查用户名是否已存在
	var existingUser entity.SystemUser
	err := db.Where("username = ?", req.Username).First(&existingUser).Error
	if err == nil {
		return nil, "", nil, nil, xError.NewError(ctx, xError.ParameterError, "用户名已存在", false)
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, "", nil, nil, xError.NewError(ctx, xError.DatabaseError, "查询用户失败", false, err)
	}

	// 2. 检查邮箱是否已存在
	err = db.Where("email = ?", req.Email).First(&existingUser).Error
	if err == nil {
		return nil, "", nil, nil, xError.NewError(ctx, xError.ParameterError, "邮箱已被注册", false)
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, "", nil, nil, xError.NewError(ctx, xError.DatabaseError, "查询邮箱失败", false, err)
	}

	// 3. 加密密码
	hashedPassword, err := xUtil.EncryptPasswordString(req.Password)
	if err != nil {
		return nil, "", nil, nil, xError.NewError(ctx, xError.ServerInternalError, "密码加密失败", false, err)
	}

	// 4. 构建用户实体
	newUser := entity.SystemUser{
		Username:    req.Username,
		Password:    hashedPassword,
		Email:       req.Email,
		Nickname:    req.Nickname,
		Role:        "user", // 新用户角色为 user
		Status:      1,      // 默认启用
		EmailVerify: false,  // 默认未验证邮箱
	}

	// 5. 创建用户（BeforeCreate Hook 会自动生成 ID 和时间戳）
	err = db.Create(&newUser).Error
	if err != nil {
		return nil, "", nil, nil, xError.NewError(ctx, xError.DatabaseError, "创建用户失败", false, err)
	}

	// 6. 生成 Token
	token := xUtil.GenerateSecurityKey()

	// 7. 记录时间信息
	now := time.Now()
	expireAt := now.Add(24 * time.Hour) // 24小时过期

	// 8. 创建用户会话
	err = a.SessionService.CreateUserSession(ctx, &newUser, token)
	if err != nil {
		return nil, "", nil, nil, xError.NewError(ctx, xError.ServerInternalError, "创建用户会话失败", false, err)
	}

	// 9. 发送邮箱验证邮件（异步，不阻断主流程）
	go a.sendEmailVerification(ctx, &newUser)

	// 10. 构建返回 DTO（注册不更新 last_login_at）
	userDTO := &dto.SystemUserDetailDTO{
		ID:          newUser.ID,
		Username:    newUser.Username,
		Email:       newUser.Email,
		Nickname:    newUser.Nickname,
		Avatar:      newUser.Avatar,
		Role:        newUser.Role,
		Status:      newUser.Status,
		EmailVerify: newUser.EmailVerify,
		LastLoginAt: nil, // 注册时不设置最后登录时间
		CreatedAt:   newUser.CreatedAt,
		UpdatedAt:   newUser.UpdatedAt,
	}

	return userDTO, token, &now, &expireAt, nil
}

// Logout 用户登出
func (a *AuthLogic) Logout(ctx *gin.Context, token string) *xError.Error {
	err := a.SessionService.DeleteUserSession(ctx, token)
	if err != nil {
		return xError.NewError(ctx, xError.ServerInternalError, "删除用户会话失败", false, err)
	}
	return nil
}

// ChangePassword 修改密码
func (a *AuthLogic) ChangePassword(ctx *gin.Context, userID int64, req *request.AuthPasswordChangeReq) *xError.Error {
	// 获取数据库连接
	db := xCtxUtil.GetDB(ctx)

	// 查找用户
	var user entity.SystemUser
	err := db.First(&user, "id = ?", userID).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return xError.NewError(ctx, xError.NotFound, "用户不存在", false)
		}
		return xError.NewError(ctx, xError.DatabaseError, "查询用户失败", false, err)
	}

	// 验证旧密码
	if !xUtil.IsPasswordValid(req.OldPassword, user.Password) {
		return xError.NewError(ctx, xError.ParameterError, "旧密码错误", false)
	}

	// 加密新密码
	hashedPassword, err := xUtil.EncryptPasswordString(req.NewPassword)
	if err != nil {
		return xError.NewError(ctx, xError.ServerInternalError, "密码加密失败", false, err)
	}

	// 更新密码
	err = db.Model(&user).Update("password", hashedPassword).Error
	if err != nil {
		return xError.NewError(ctx, xError.DatabaseError, "更新密码失败", false, err)
	}

	return nil
}

// ResetPassword 重置密码（发送重置链接）
func (a *AuthLogic) ResetPassword(ctx *gin.Context, req *request.AuthPasswordResetReq) *xError.Error {
	// 获取数据库连接
	db := xCtxUtil.GetDB(ctx)

	// 查找用户
	var user entity.SystemUser
	err := db.First(&user, "email = ?", req.Email).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return xError.NewError(ctx, xError.NotFound, "邮箱不存在", false)
		}
		return xError.NewError(ctx, xError.DatabaseError, "查询用户失败", false, err)
	}

	// 生成重置 Token（32位随机字符串）
	resetToken := generateRandomString(32)

	// 存储到 Redis（1小时过期）
	rdb := ctxUtil.GetRedisClient(ctx)
	if rdb == nil {
		return xError.NewError(ctx, xError.ServerInternalError, "Redis 客户端不可用", false)
	}

	redisKey := fmt.Sprintf(constants.PasswordResetTokenPrefix, resetToken)
	err = rdb.Set(ctx.Request.Context(), redisKey, user.ID, time.Hour).Err()
	if err != nil {
		return xError.NewError(ctx, xError.ServerInternalError, "保存重置Token失败", false, err)
	}

	// 发送重置密码邮件（异步，不阻断主流程）
	go a.sendPasswordResetEmail(ctx, &user, resetToken)

	return nil
}

// GetUserInfo 获取用户信息
func (a *AuthLogic) GetUserInfo(ctx *gin.Context, userID int64) (*dto.SystemUserDetailDTO, *xError.Error) {
	// 获取数据库连接
	db := xCtxUtil.GetDB(ctx)

	var user entity.SystemUser
	err := db.First(&user, "id = ?", userID).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, xError.NewError(ctx, xError.NotFound, "用户不存在", false)
		}
		return nil, xError.NewError(ctx, xError.DatabaseError, "查询用户失败", false, err)
	}

	userDTO := &dto.SystemUserDetailDTO{
		ID:          user.ID,
		Username:    user.Username,
		Email:       user.Email,
		Nickname:    user.Nickname,
		Avatar:      user.Avatar,
		Role:        user.Role,
		Status:      user.Status,
		LastLoginAt: user.LastLoginAt,
		CreatedAt:   user.CreatedAt,
		UpdatedAt:   user.UpdatedAt,
	}
	return userDTO, nil
}

// UpdateLastLogin 更新最后登录时间
func (a *AuthLogic) UpdateLastLogin(ctx *gin.Context, userID int64) *xError.Error {
	// 获取数据库连接
	db := xCtxUtil.GetDB(ctx)

	now := time.Now()
	err := db.Model(&entity.SystemUser{}).Where("id = ?", userID).Update("last_login_at", &now).Error
	if err != nil {
		return xError.NewError(ctx, xError.DatabaseError, "更新最后登录时间失败", false, err)
	}
	return nil
}

// ValidateToken 验证令牌
func (a *AuthLogic) ValidateToken(ctx *gin.Context, token string) (*dto.SystemUserDetailDTO, *xError.Error) {
	// 这个方法主要通过中间件来处理，这里提供一个备用实现
	// 实际项目中可以根据需要实现更复杂的验证逻辑
	return nil, xError.NewError(ctx, xError.OperationNotSupported, "请通过认证中间件验证令牌", false)
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

// VerifyEmail 验证邮箱
func (a *AuthLogic) VerifyEmail(ctx *gin.Context, req *request.AuthVerifyEmailReq) *xError.Error {
	logger := xCtxUtil.GetSugarLogger(ctx, "AUTH")

	// 获取 Redis 客户端
	rdb := ctxUtil.GetRedisClient(ctx)
	if rdb == nil {
		return xError.NewError(ctx, xError.ServerInternalError, "Redis 客户端不可用", false)
	}

	// 从 Redis 获取用户 ID
	redisKey := fmt.Sprintf(constants.EmailVerifyTokenPrefix, req.Token)
	userIDStr, err := rdb.Get(ctx.Request.Context(), redisKey).Result()
	if err != nil {
		logger.Warnf("邮箱验证Token无效或已过期: %s", req.Token)
		return xError.NewError(ctx, xError.BadRequest, "验证链接无效或已过期", false)
	}

	// 删除已使用的 Token
	rdb.Del(ctx.Request.Context(), redisKey)

	// 更新用户邮箱验证状态
	db := xCtxUtil.GetDB(ctx)
	err = db.Model(&entity.SystemUser{}).Where("id = ?", userIDStr).Update("email_verify", true).Error
	if err != nil {
		return xError.NewError(ctx, xError.DatabaseError, "更新邮箱验证状态失败", false, err)
	}

	logger.Infof("用户 %s 邮箱验证成功", userIDStr)
	return nil
}

// VerifyResetToken 验证重置密码Token
func (a *AuthLogic) VerifyResetToken(ctx *gin.Context, req *request.AuthVerifyResetTokenReq) (bool, *xError.Error) {
	// 获取 Redis 客户端
	rdb := ctxUtil.GetRedisClient(ctx)
	if rdb == nil {
		return false, xError.NewError(ctx, xError.ServerInternalError, "Redis 客户端不可用", false)
	}

	// 检查 Token 是否存在
	redisKey := fmt.Sprintf(constants.PasswordResetTokenPrefix, req.Token)
	exists, err := rdb.Exists(ctx.Request.Context(), redisKey).Result()
	if err != nil {
		return false, xError.NewError(ctx, xError.ServerInternalError, "验证Token失败", false, err)
	}

	return exists > 0, nil
}

// ConfirmResetPassword 确认重置密码
func (a *AuthLogic) ConfirmResetPassword(ctx *gin.Context, req *request.AuthConfirmResetPasswordReq) *xError.Error {
	logger := xCtxUtil.GetSugarLogger(ctx, "AUTH")

	// 获取 Redis 客户端
	rdb := ctxUtil.GetRedisClient(ctx)
	if rdb == nil {
		return xError.NewError(ctx, xError.ServerInternalError, "Redis 客户端不可用", false)
	}

	// 从 Redis 获取用户 ID
	redisKey := fmt.Sprintf(constants.PasswordResetTokenPrefix, req.Token)
	userIDStr, err := rdb.Get(ctx.Request.Context(), redisKey).Result()
	if err != nil {
		logger.Warnf("密码重置Token无效或已过期: %s", req.Token)
		return xError.NewError(ctx, xError.BadRequest, "重置链接无效或已过期", false)
	}

	// 删除已使用的 Token
	rdb.Del(ctx.Request.Context(), redisKey)

	// 加密新密码
	hashedPassword, err := xUtil.EncryptPasswordString(req.NewPassword)
	if err != nil {
		return xError.NewError(ctx, xError.ServerInternalError, "密码加密失败", false, err)
	}

	// 更新用户密码
	db := xCtxUtil.GetDB(ctx)
	err = db.Model(&entity.SystemUser{}).Where("id = ?", userIDStr).Update("password", hashedPassword).Error
	if err != nil {
		return xError.NewError(ctx, xError.DatabaseError, "重置密码失败", false, err)
	}

	logger.Infof("用户 %s 密码重置成功", userIDStr)
	return nil
}

// sendEmailVerification 发送邮箱验证邮件
//
// 此函数应在 goroutine 中异步调用，不会阻断主流程
func (a *AuthLogic) sendEmailVerification(ctx *gin.Context, user *entity.SystemUser) {
	logger := xCtxUtil.GetSugarLogger(ctx, "MAIL")

	// 获取配置
	config := ctxUtil.GetConfig(ctx)
	if config == nil {
		logger.Warn("无法获取配置，跳过发送邮箱验证邮件")
		return
	}

	// 获取 Redis 客户端
	rdb := ctxUtil.GetRedisClient(ctx)
	if rdb == nil {
		logger.Warn("Redis 客户端不可用，跳过发送邮箱验证邮件")
		return
	}

	// 生成验证 Token
	verifyToken := generateRandomString(32)

	// 存储到 Redis（24小时过期）
	redisKey := fmt.Sprintf(constants.EmailVerifyTokenPrefix, verifyToken)
	err := rdb.Set(ctx.Request.Context(), redisKey, user.ID, 24*time.Hour).Err()
	if err != nil {
		logger.Warnf("保存验证Token失败: %v", err)
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
	mailLogic := &MailLogic{TemplateService: servHelper.NewMailTemplateService(), MaxRetry: 3}
	mailErr := mailLogic.SendWithTemplate(
		ctx,
		"email_verify",
		[]string{user.Email},
		"请验证您的邮箱地址",
		variables,
	)
	if mailErr != nil {
		logger.Warnf("发送邮箱验证邮件失败: %v", mailErr)
	} else {
		logger.Infof("已发送邮箱验证邮件到: %s", user.Email)
	}
}

// sendPasswordResetEmail 发送密码重置邮件
//
// 此函数应在 goroutine 中异步调用，不会阻断主流程
func (a *AuthLogic) sendPasswordResetEmail(ctx *gin.Context, user *entity.SystemUser, resetToken string) {
	logger := xCtxUtil.GetSugarLogger(ctx, "MAIL")

	// 获取配置
	config := ctxUtil.GetConfig(ctx)
	if config == nil {
		logger.Warn("无法获取配置，跳过发送密码重置邮件")
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
	mailLogic := &MailLogic{TemplateService: servHelper.NewMailTemplateService(), MaxRetry: 3}
	mailErr := mailLogic.SendWithTemplate(
		ctx,
		"password_reset",
		[]string{user.Email},
		"密码重置请求",
		variables,
	)
	if mailErr != nil {
		logger.Warnf("发送密码重置邮件失败: %v", mailErr)
	} else {
		logger.Infof("已发送密码重置邮件到: %s", user.Email)
	}
}
