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

package service

import (
	"bamboo-main/internal/logic"
	servHelper "bamboo-main/internal/service/helper"
	"time"

	"bamboo-main/internal/model/dto"
	"bamboo-main/internal/model/request"

	xError "github.com/bamboo-services/bamboo-base-go/error"
	"github.com/gin-gonic/gin"
)

// IAuthService 认证服务接口
type IAuthService interface {
	// Login 用户登录
	Login(ctx *gin.Context, req *request.AuthLoginReq) (*dto.SystemUserDTO, string, *time.Time, *time.Time, *xError.Error)

	// Register 用户注册
	Register(ctx *gin.Context, r *request.AuthRegisterReq) (*dto.SystemUserDTO, string, *time.Time, *time.Time, *xError.Error)

	// Logout 用户登出
	Logout(ctx *gin.Context, token string) *xError.Error

	// ChangePassword 修改密码
	ChangePassword(ctx *gin.Context, userID int64, req *request.AuthPasswordChangeReq) *xError.Error

	// ResetPassword 重置密码
	ResetPassword(ctx *gin.Context, req *request.AuthPasswordResetReq) *xError.Error

	// GetUserInfo 获取用户信息
	GetUserInfo(ctx *gin.Context, userID int64) (*dto.SystemUserDTO, *xError.Error)

	// UpdateLastLogin 更新最后登录时间
	UpdateLastLogin(ctx *gin.Context, userID int64) *xError.Error

	// ValidateToken 验证令牌
	ValidateToken(ctx *gin.Context, token string) (*dto.SystemUserDTO, *xError.Error)
}

// NewAuthService 创建认证服务实例
func NewAuthService() IAuthService {
	return &logic.AuthLogic{
		SessionService: servHelper.NewSessionService(),
	}
}
