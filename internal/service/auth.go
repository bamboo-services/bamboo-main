package service

import (
	"bamboo-main/internal/logic"
	"bamboo-main/pkg/startup"
)

// AuthService 认证服务类型别名
type AuthService = logic.AuthLogic

// NewAuthService 创建认证服务实例
func NewAuthService(reg *startup.Reg) *AuthService {
	return logic.NewAuthLogic(reg)
}