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
)

// AuthService 认证服务类型别名
type AuthService = logic.AuthLogic

// NewAuthService 创建认证服务实例
func NewAuthService() *AuthService {
	return logic.NewAuthLogic()
}
