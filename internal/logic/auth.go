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
	"time"

	"bamboo-main/internal/model/dto"
	"bamboo-main/internal/model/entity"
	"bamboo-main/internal/model/request"
	servHelper "bamboo-main/internal/service/helper"

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
	err := db.WithContext(ctx.Request.Context()).Where("username = ? OR email = ?", req.Username, req.Username).First(&user).Error
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
	err = db.WithContext(ctx.Request.Context()).Model(&user).Update("last_login_at", &now).Error
	if err != nil {
		// 记录错误但不影响登录
		xCtxUtil.GetSugarLogger(ctx, "").Errorf("更新最后登录时间失败: %v", err)
	}

	userDTO := &dto.SystemUserDetailDTO{
		ID:          user.ID,
		Username:    user.Username,
		Email:       user.Email,
		Nickname:    user.Nickname, // 直接赋值指针 *string → *string
		Avatar:      user.Avatar,   // 直接赋值指针 *string → *string
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
	err := db.WithContext(ctx.Request.Context()).Where("username = ?", req.Username).First(&existingUser).Error
	if err == nil {
		return nil, "", nil, nil, xError.NewError(ctx, xError.ParameterError, "用户名已存在", false)
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, "", nil, nil, xError.NewError(ctx, xError.DatabaseError, "查询用户失败", false, err)
	}

	// 2. 检查邮箱是否已存在
	err = db.WithContext(ctx.Request.Context()).Where("email = ?", req.Email).First(&existingUser).Error
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
	err = db.WithContext(ctx.Request.Context()).Create(&newUser).Error
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

	// 9. TODO: 发送邮箱验证邮件
	// 功能说明：
	//   - 生成邮箱验证令牌（建议使用 JWT 或随机字符串）
	//   - 将令牌存储到 Redis，设置过期时间（如 24 小时）
	//   - 发送包含验证链接的邮件到用户邮箱
	//   - 验证链接格式：https://域名/api/v1/auth/verify-email?token=xxx
	//   - 用户点击链接后，验证 Token 并更新 email_verify 字段为 true
	// 相关文件：
	//   - 邮件发送服务：待实现（参考 pkg/constants/redis.go 的 EmailLimitPrefix）
	//   - 验证接口：待添加到 router_auth.go
	xCtxUtil.GetSugarLogger(ctx, "").Infof("用户 %s 注册成功，邮箱 %s 待验证", newUser.Username, newUser.Email)

	// 10. 构建返回 DTO（注册不更新 last_login_at）
	userDTO := &dto.SystemUserDetailDTO{
		ID:          newUser.ID,
		Username:    newUser.Username,
		Email:       newUser.Email,
		Nickname:    newUser.Nickname, // 直接赋值指针 *string → *string
		Avatar:      newUser.Avatar,   // 直接赋值指针 *string → *string
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
	err := db.WithContext(ctx.Request.Context()).First(&user, "id = ?", userID).Error
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
	err = db.WithContext(ctx.Request.Context()).Model(&user).Update("password", hashedPassword).Error
	if err != nil {
		return xError.NewError(ctx, xError.DatabaseError, "更新密码失败", false, err)
	}

	return nil
}

// ResetPassword 重置密码
func (a *AuthLogic) ResetPassword(ctx *gin.Context, req *request.AuthPasswordResetReq) *xError.Error {
	// 获取数据库连接
	db := xCtxUtil.GetDB(ctx)

	// 查找用户
	var user entity.SystemUser
	err := db.WithContext(ctx.Request.Context()).First(&user, "email = ?", req.Email).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return xError.NewError(ctx, xError.NotFound, "邮箱不存在", false)
		}
		return xError.NewError(ctx, xError.DatabaseError, "查询用户失败", false, err)
	}

	// 生成临时密码
	tempPassword := generateRandomString(12)

	// 加密密码
	hashedPassword, err := xUtil.EncryptPasswordString(tempPassword)
	if err != nil {
		return xError.NewError(ctx, xError.ServerInternalError, "密码加密失败", false, err)
	}

	// 更新密码
	err = db.WithContext(ctx.Request.Context()).Model(&user).Update("password", hashedPassword).Error
	if err != nil {
		return xError.NewError(ctx, xError.DatabaseError, "重置密码失败", false, err)
	}

	// TODO: 发送邮件通知新密码
	xCtxUtil.GetSugarLogger(ctx, "").Infof("用户 %s 的临时密码为: %s", user.Email, tempPassword)

	return nil
}

// GetUserInfo 获取用户信息
func (a *AuthLogic) GetUserInfo(ctx *gin.Context, userID int64) (*dto.SystemUserDetailDTO, *xError.Error) {
	// 获取数据库连接
	db := xCtxUtil.GetDB(ctx)

	var user entity.SystemUser
	err := db.WithContext(ctx.Request.Context()).First(&user, "id = ?", userID).Error
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
		Nickname:    user.Nickname, // 直接赋值指针 *string → *string
		Avatar:      user.Avatar,   // 直接赋值指针 *string → *string
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
	err := db.WithContext(ctx.Request.Context()).Model(&entity.SystemUser{}).Where("id = ?", userID).Update("last_login_at", &now).Error
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
