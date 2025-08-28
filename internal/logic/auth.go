package logic

import (
	"time"

	"bamboo-main/internal/middleware"
	"bamboo-main/internal/model/dto"
	"bamboo-main/internal/model/entity"
	"bamboo-main/internal/model/request"

	"crypto/rand"

	xError "github.com/bamboo-services/bamboo-base-go/error"
	xCtxUtil "github.com/bamboo-services/bamboo-base-go/utility/ctxutil"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// AuthLogic 认证业务逻辑
type AuthLogic struct {
}

// NewAuthLogic 创建认证业务逻辑实例
func NewAuthLogic() *AuthLogic {
	return &AuthLogic{}
}

// Login 用户登录
func (a *AuthLogic) Login(ctx *gin.Context, req *request.AuthLoginReq) (*dto.SystemUserDTO, string, *xError.Error) {
	// 获取数据库连接
	db := xCtxUtil.GetDB(ctx)

	// 查找用户
	var user entity.SystemUser
	err := db.WithContext(ctx.Request.Context()).Where("username = ? OR email = ?", req.Username, req.Username).First(&user).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, "", xError.NewError(ctx, xError.LoginFailed, "用户名或密码错误", false)
		}
		return nil, "", xError.NewError(ctx, xError.DatabaseError, "查询用户失败", false, err)
	}

	// 检查用户状态
	if user.Status == 0 {
		return nil, "", xError.NewError(ctx, xError.Forbidden, "用户已被禁用", false)
	}

	// 验证密码
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		return nil, "", xError.NewError(ctx, xError.LoginFailed, "用户名或密码错误", false)
	}

	// 生成 token
	token := generateSecurityKey()

	// 创建用户会话
	err = middleware.CreateUserSession(ctx, &user, token)
	if err != nil {
		return nil, "", xError.NewError(ctx, xError.ServerInternalError, "创建用户会话失败", false, err)
	}

	// 更新最后登录时间
	now := time.Now()
	err = db.WithContext(ctx.Request.Context()).Model(&user).Update("last_login_at", &now).Error
	if err != nil {
		// 记录错误但不影响登录
		xCtxUtil.GetSugarLogger(ctx).Errorf("更新最后登录时间失败: %v", err)
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
func (a *AuthLogic) Logout(ctx *gin.Context, token string) *xError.Error {
	err := middleware.DeleteUserSession(ctx, token)
	if err != nil {
		return xError.NewError(ctx, xError.ServerInternalError, "删除用户会话失败", false, err)
	}
	return nil
}

// ChangePassword 修改密码
func (a *AuthLogic) ChangePassword(ctx *gin.Context, userUUID string, req *request.AuthPasswordChangeReq) *xError.Error {
	// 获取数据库连接
	db := xCtxUtil.GetDB(ctx)

	// 查找用户
	var user entity.SystemUser
	err := db.WithContext(ctx.Request.Context()).First(&user, "uuid = ?", userUUID).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return xError.NewError(ctx, xError.NotFound, "用户不存在", false)
		}
		return xError.NewError(ctx, xError.DatabaseError, "查询用户失败", false, err)
	}

	// 验证旧密码
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.OldPassword))
	if err != nil {
		return xError.NewError(ctx, xError.ParameterError, "旧密码错误", false)
	}

	// 加密新密码
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		return xError.NewError(ctx, xError.ServerInternalError, "密码加密失败", false, err)
	}

	// 更新密码
	err = db.WithContext(ctx.Request.Context()).Model(&user).Update("password", string(hashedPassword)).Error
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
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(tempPassword), bcrypt.DefaultCost)
	if err != nil {
		return xError.NewError(ctx, xError.ServerInternalError, "密码加密失败", false, err)
	}

	// 更新密码
	err = db.WithContext(ctx.Request.Context()).Model(&user).Update("password", string(hashedPassword)).Error
	if err != nil {
		return xError.NewError(ctx, xError.DatabaseError, "重置密码失败", false, err)
	}

	// TODO: 发送邮件通知新密码
	xCtxUtil.GetSugarLogger(ctx).Infof("用户 %s 的临时密码为: %s", user.Email, tempPassword)

	return nil
}

// GetUserInfo 获取用户信息
func (a *AuthLogic) GetUserInfo(ctx *gin.Context, userUUID string) (*dto.SystemUserDTO, *xError.Error) {
	// 获取数据库连接
	db := xCtxUtil.GetDB(ctx)

	var user entity.SystemUser
	err := db.WithContext(ctx.Request.Context()).First(&user, "uuid = ?", userUUID).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, xError.NewError(ctx, xError.NotFound, "用户不存在", false)
		}
		return nil, xError.NewError(ctx, xError.DatabaseError, "查询用户失败", false, err)
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
func (a *AuthLogic) UpdateLastLogin(ctx *gin.Context, userUUID string) *xError.Error {
	// 获取数据库连接
	db := xCtxUtil.GetDB(ctx)

	now := time.Now()
	err := db.WithContext(ctx.Request.Context()).Model(&entity.SystemUser{}).Where("uuid = ?", userUUID).Update("last_login_at", &now).Error
	if err != nil {
		return xError.NewError(ctx, xError.DatabaseError, "更新最后登录时间失败", false, err)
	}
	return nil
}

// ValidateToken 验证令牌
func (a *AuthLogic) ValidateToken(ctx *gin.Context, token string) (*dto.SystemUserDTO, *xError.Error) {
	// 这个方法主要通过中间件来处理，这里提供一个备用实现
	// 实际项目中可以根据需要实现更复杂的验证逻辑
	return nil, xError.NewError(ctx, xError.OperationNotSupported, "请通过认证中间件验证令牌", false)
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
