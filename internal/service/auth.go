package service

import (
	"bamboo-main/internal/logic"
)

// AuthService 认证服务类型别名
type AuthService = logic.AuthLogic

// NewAuthService 创建认证服务实例
func NewAuthService() *AuthService {
	return logic.NewAuthLogic()
}