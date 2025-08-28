package middleware

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"bamboo-main/internal/helper"
	"bamboo-main/pkg/constants"
	ctxUtil "bamboo-main/pkg/util/ctx"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

// AuthMiddleware 认证中间件
func AuthMiddleware(c *gin.Context) {
	// 获取 Authorization 头部
	authHeader := c.GetHeader(constants.HeaderAuthorization)
	if authHeader == "" {
		c.JSON(http.StatusUnauthorized, gin.H{
			"code":    401,
			"message": "缺少认证令牌",
			"data":    nil,
		})
		c.Abort()
		return
	}

	// 检查 Bearer 前缀
	if !strings.HasPrefix(authHeader, constants.TokenPrefix) {
		c.JSON(http.StatusUnauthorized, gin.H{
			"code":    401,
			"message": "认证令牌格式错误",
			"data":    nil,
		})
		c.Abort()
		return
	}

	// 提取 token
	token := strings.TrimPrefix(authHeader, constants.TokenPrefix)
	if token == "" {
		c.JSON(http.StatusUnauthorized, gin.H{
			"code":    401,
			"message": "认证令牌不能为空",
			"data":    nil,
		})
		c.Abort()
		return
	}

	// 获取 Redis 客户端
	rdb := ctxUtil.GetRedisClient(c)
	if rdb == nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "Redis 连接异常",
			"data":    nil,
		})
		c.Abort()
		return
	}

	// 从 Redis 获取用户会话
	redisKey := fmt.Sprintf(constants.AuthTokenPrefix, token)
	sessionData, err := rdb.Get(c.Request.Context(), redisKey).Result()
	if err != nil {
		if err == redis.Nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code":    401,
				"message": "认证令牌已过期或无效",
				"data":    nil,
			})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{
				"code":    500,
				"message": "认证服务异常",
				"data":    nil,
			})
		}
		c.Abort()
		return
	}

	// 解析用户会话数据
	var session helper.UserSession
	err = json.Unmarshal([]byte(sessionData), &session)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "会话数据解析失败",
			"data":    nil,
		})
		c.Abort()
		return
	}

	// 检查会话是否过期
	if time.Now().After(session.ExpireAt) {
		// 删除过期的 token
		rdb.Del(c.Request.Context(), redisKey)
		c.JSON(http.StatusUnauthorized, gin.H{
			"code":    401,
			"message": "认证令牌已过期",
			"data":    nil,
		})
		c.Abort()
		return
	}

	// 将用户信息存储到上下文中
	c.Set(constants.ContextKeyUser, session)
	c.Set(constants.ContextKeyUserID, session.UserUUID)
	c.Set(constants.ContextKeyToken, token)

	c.Next()
}