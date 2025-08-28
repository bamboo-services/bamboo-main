package middleware

import (
	"encoding/json"
	"fmt"
	"time"

	"bamboo-main/internal/helper"
	"bamboo-main/internal/model/entity"
	"bamboo-main/pkg/constants"
	ctxUtil "bamboo-main/pkg/util/ctx"

	"github.com/gin-gonic/gin"
)

// CreateUserSession 创建用户会话
func CreateUserSession(c *gin.Context, user *entity.SystemUser, token string) error {
	// 获取 Redis 客户端
	rdb := ctxUtil.GetRedisClient(c)
	if rdb == nil {
		return fmt.Errorf("Redis 客户端不可用")
	}

	session := helper.UserSession{
		UserUUID: user.UUID.String(),
		Username: user.Username,
		Email:    user.Email,
		Role:     user.Role,
		LoginAt:  time.Now(),
		ExpireAt: time.Now().Add(24 * time.Hour), // 24小时过期
	}

	// 序列化会话数据
	sessionData, err := json.Marshal(session)
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
func DeleteUserSession(c *gin.Context, token string) error {
	// 获取 Redis 客户端
	rdb := ctxUtil.GetRedisClient(c)
	if rdb == nil {
		return fmt.Errorf("Redis 客户端不可用")
	}

	redisKey := fmt.Sprintf(constants.AuthTokenPrefix, token)
	return rdb.Del(c.Request.Context(), redisKey).Err()
}