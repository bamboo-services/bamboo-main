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

package middleware

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	dtoRedis "bamboo-main/internal/model/dto/redis"
	"bamboo-main/pkg/constants"
	ctxUtil "bamboo-main/pkg/util/ctx"

	xError "github.com/bamboo-services/bamboo-base-go/error"
	xResult "github.com/bamboo-services/bamboo-base-go/result"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

// AuthMiddleware 认证中间件
func AuthMiddleware(c *gin.Context) {
	// 获取 Authorization 头部
	authHeader := c.GetHeader(constants.HeaderAuthorization)
	if authHeader == "" {
		xResult.Error(c, xError.HeaderMissing, "认证令牌缺失", nil)
		c.Abort()
		return
	}

	// 检查 Bearer 前缀
	if !strings.HasPrefix(authHeader, constants.TokenPrefix) {
		xResult.Error(c, xError.HeaderIllegal, "认证令牌格式错误", nil)
		c.Abort()
		return
	}

	// 提取 token
	token := strings.TrimPrefix(authHeader, constants.TokenPrefix)
	if token == "" {
		xResult.Error(c, xError.HeaderMissing, "认证令牌不能为空", nil)
		c.Abort()
		return
	}

	// 获取 Redis 客户端
	rdb := ctxUtil.GetRedisClient(c)
	if rdb == nil {
		xResult.Error(c, xError.DatabaseError, "Redis 链接异常", nil)
		c.Abort()
		return
	}

	// 从 Redis 获取用户会话
	redisKey := fmt.Sprintf(constants.AuthTokenPrefix, token)
	sessionData, err := rdb.Get(c.Request.Context(), redisKey).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			xResult.Error(c, xError.Unauthorized, "认证令牌已过期或无效", nil)
		} else {
			xResult.Error(c, xError.DatabaseError, "认证服务异常", nil)
		}
		c.Abort()
		return
	}

	// 解析用户会话数据
	var session dtoRedis.TokenDTO
	err = json.Unmarshal([]byte(sessionData), &session)
	if err != nil {
		xResult.Error(c, xError.ServerInternalError, "会话数据解析失败", nil)
		c.Abort()
		return
	}

	// 检查会话是否过期
	if time.Now().After(session.ExpiredAt) {
		rdb.Del(c.Request.Context(), redisKey)
		xResult.Error(c, xError.Unauthorized, "认证令牌已过期", nil)
		c.Abort()
		return
	}

	// 将用户信息存储到上下文中
	c.Set(constants.ContextKeyUserID, session.UserUUID)
	c.Set(constants.ContextKeyToken, token)

	c.Next()
}
