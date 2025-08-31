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

package servHelper

import (
	logcHelper "bamboo-main/internal/logic/helper"
	"bamboo-main/internal/model/entity"

	"github.com/gin-gonic/gin"
)

// ISessionService 会话管理服务接口
type ISessionService interface {
	// CreateUserSession 创建用户会话
	CreateUserSession(c *gin.Context, user *entity.SystemUser, token string) error

	// DeleteUserSession 删除用户会话
	DeleteUserSession(c *gin.Context, token string) error
}

// NewSessionService 创建会话管理服务实例
func NewSessionService() ISessionService {
	return &logcHelper.SessionLogic{}
}
