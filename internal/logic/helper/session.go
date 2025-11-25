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

package logcHelper

import (
	"encoding/json"
	"fmt"
	"time"

	dtoRedis "bamboo-main/internal/model/dto/redis"
	"bamboo-main/internal/model/entity"
	"bamboo-main/pkg/constants"
	ctxUtil "bamboo-main/pkg/util/ctx"
	"bamboo-main/pkg/util/netUtil"

	"github.com/gin-gonic/gin"
)

// SessionLogic 会话管理服务实现
type SessionLogic struct{}

// CreateUserSession 创建用户会话
func (s *SessionLogic) CreateUserSession(c *gin.Context, user *entity.SystemUser, token string) error {
	// 获取 Redis 客户端
	rdb := ctxUtil.GetRedisClient(c)
	if rdb == nil {
		return fmt.Errorf("redis 客户端不可用")
	}

	// 获取客户端真实IP地址
	clientIP := netUtil.GetClientIP(c)

	// 获取用户代理
	userAgent := c.GetHeader("User-Agent")

	// 创建Token会话数据
	now := time.Now()
	tokenDTO := dtoRedis.TokenDTO{
		UserID:    user.ID,
		Username:  user.Username,
		Email:     user.Email,
		Role:      user.Role,
		LoginIP:   clientIP,
		UserAgent: userAgent,
		CreatedAt: now,
		ExpiredAt: now.Add(24 * time.Hour), // 24小时过期
	}

	// 序列化会话数据
	sessionData, err := json.Marshal(tokenDTO)
	if err != nil {
		return err
	}

	// 存储到 Redis
	redisKey := fmt.Sprintf(constants.AuthTokenPrefix, token)
	err = rdb.Set(c.Request.Context(), redisKey, sessionData, 24*time.Hour).Err()
	if err != nil {
		return err
	}

	return nil
}

// DeleteUserSession 删除用户会话
func (s *SessionLogic) DeleteUserSession(c *gin.Context, token string) error {
	// 获取 Redis 客户端
	rdb := ctxUtil.GetRedisClient(c)
	if rdb == nil {
		return fmt.Errorf("redis 客户端不可用")
	}

	redisKey := fmt.Sprintf(constants.AuthTokenPrefix, token)
	return rdb.Del(c.Request.Context(), redisKey).Err()
}
