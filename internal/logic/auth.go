package logic

import (
	"context"
	"fmt"
	"time"

	"bamboo-main/internal/middleware"
	"bamboo-main/internal/model/dto"
	"bamboo-main/internal/model/entity"
	"bamboo-main/internal/model/request"
	"bamboo-main/pkg/startup"

	"crypto/rand"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// AuthLogic 认证业务逻辑
type AuthLogic struct {
	reg *startup.Reg
}

// NewAuthLogic 创建认证业务逻辑实例
func NewAuthLogic(reg *startup.Reg) *AuthLogic {
	return &AuthLogic{
		reg: reg,
	}
}

// Login 用户登录
func (a *AuthLogic) Login(ctx *gin.Context, req *request.AuthLoginReq) (*dto.SystemUserDTO, string, error) {
	// 查找用户
	var user entity.SystemUser
	err := a.reg.DB.WithContext(ctx.Request.Context()).Where("username = ? OR email = ?", req.Username, req.Username).First(&user).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, "", fmt.Errorf("用户名或密码错误")
		}
		return nil, "", fmt.Errorf("查询用户失败: %w", err)
	}

	// 检查用户状态
	if user.Status == 0 {
		return nil, "", fmt.Errorf("用户已被禁用")
	}

	// 验证密码
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		return nil, "", fmt.Errorf("用户名或密码错误")
	}

	// 生成 token
	token := generateSecurityKey()

	// 创建用户会话
	err = middleware.CreateUserSession(a.reg, &user, token)
	if err != nil {
		return nil, "", fmt.Errorf("创建用户会话失败: %w", err)
	}

	// 更新最后登录时间
	now := time.Now()
	err = a.reg.DB.WithContext(ctx.Request.Context()).Model(&user).Update("last_login_at", &now).Error
	if err != nil {
		// 记录错误但不影响登录
		a.reg.Serv.Logger.Sugar().Errorf("更新最后登录时间失败: %v", err)
	}

	userDTO := &dto.SystemUserDTO{
		UUID:        user.UUID.String(),
		Username:    user.Username,
		Email:       user.Email,
		Nickname:    safeStringValue(user.Nickname),
		Avatar:      safeStringValue(user.Avatar),
		Role:        user.Role,
		Status:      user.Status,
		LastLoginAt: user.LastLoginAt,
		CreatedAt:   user.CreatedAt,
		UpdatedAt:   user.UpdatedAt,
	}
	return userDTO, token, nil
}

// Logout 用户登出
func (a *AuthLogic) Logout(ctx *gin.Context, token string) error {
	return middleware.DeleteUserSession(a.reg, token)
}

// ChangePassword 修改密码
func (a *AuthLogic) ChangePassword(ctx *gin.Context, userUUID string, req *request.AuthPasswordChangeReq) error {
	// 查找用户
	var user entity.SystemUser
	err := a.reg.DB.WithContext(ctx.Request.Context()).First(&user, "uuid = ?", userUUID).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return fmt.Errorf("用户不存在")
		}
		return fmt.Errorf("查询用户失败: %w", err)
	}

	// 验证旧密码
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.OldPassword))
	if err != nil {
		return fmt.Errorf("旧密码错误")
	}

	// 加密新密码
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("密码加密失败: %w", err)
	}

	// 更新密码
	err = a.reg.DB.WithContext(ctx.Request.Context()).Model(&user).Update("password", string(hashedPassword)).Error
	if err != nil {
		return fmt.Errorf("更新密码失败: %w", err)
	}

	return nil
}

// ResetPassword 重置密码
func (a *AuthLogic) ResetPassword(ctx *gin.Context, req *request.AuthPasswordResetReq) error {
	// 查找用户
	var user entity.SystemUser
	err := a.reg.DB.WithContext(ctx.Request.Context()).First(&user, "email = ?", req.Email).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return fmt.Errorf("邮箱不存在")
		}
		return fmt.Errorf("查询用户失败: %w", err)
	}

	// 生成临时密码
	tempPassword := generateRandomString(12)

	// 加密密码
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(tempPassword), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("密码加密失败: %w", err)
	}

	// 更新密码
	err = a.reg.DB.WithContext(ctx.Request.Context()).Model(&user).Update("password", string(hashedPassword)).Error
	if err != nil {
		return fmt.Errorf("重置密码失败: %w", err)
	}

	// TODO: 发送邮件通知新密码
	a.reg.Serv.Logger.Sugar().Infof("用户 %s 的临时密码为: %s", user.Email, tempPassword)

	return nil
}

// GetUserInfo 获取用户信息
func (a *AuthLogic) GetUserInfo(ctx *gin.Context, userUUID string) (*dto.SystemUserDTO, error) {
	var user entity.SystemUser
	err := a.reg.DB.WithContext(ctx.Request.Context()).First(&user, "uuid = ?", userUUID).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("用户不存在")
		}
		return nil, fmt.Errorf("查询用户失败: %w", err)
	}

	userDTO := &dto.SystemUserDTO{
		UUID:        user.UUID.String(),
		Username:    user.Username,
		Email:       user.Email,
		Nickname:    safeStringValue(user.Nickname),
		Avatar:      safeStringValue(user.Avatar),
		Role:        user.Role,
		Status:      user.Status,
		LastLoginAt: user.LastLoginAt,
		CreatedAt:   user.CreatedAt,
		UpdatedAt:   user.UpdatedAt,
	}
	return userDTO, nil
}

// UpdateLastLogin 更新最后登录时间
func (a *AuthLogic) UpdateLastLogin(ctx context.Context, userUUID string) error {
	now := time.Now()
	err := a.reg.DB.WithContext(ctx).Model(&entity.SystemUser{}).Where("uuid = ?", userUUID).Update("last_login_at", &now).Error
	if err != nil {
		return fmt.Errorf("更新最后登录时间失败: %w", err)
	}
	return nil
}

// ValidateToken 验证令牌
func (a *AuthLogic) ValidateToken(ctx context.Context, token string) (*dto.SystemUserDTO, error) {
	// 这个方法主要通过中间件来处理，这里提供一个备用实现
	// 实际项目中可以根据需要实现更复杂的验证逻辑
	return nil, fmt.Errorf("请通过认证中间件验证令牌")
}

// generateSecurityKey 生成安全密钥，格式：cs_ + 32字符 + 32字符
func generateSecurityKey() string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	const keyLength = 64
	
	b := make([]byte, keyLength)
	rand.Read(b)
	
	result := make([]byte, keyLength)
	for i := range b {
		result[i] = charset[b[i]%byte(len(charset))]
	}
	
	return "cs_" + string(result)
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

// safeStringValue 安全转换指针字符串为字符串
func safeStringValue(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}