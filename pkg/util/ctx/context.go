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

package ctxUtil

import (
	"bamboo-main/internal/helper"
	"bamboo-main/internal/model/base"
	"bamboo-main/pkg/constants"

	xConsts "github.com/bamboo-services/bamboo-base-go/constants"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

// GetConfig 从 Gin 的上下文中获取 `BambooConfig` 配置实例。
//
// 如果上下文中存在对应配置，则返回其指针；如果不存在，则返回 `nil`。
// 
// 注意: 这是项目特定的配置获取函数，与 bamboo-base-go 的 GetConfig 不同
//
// 参数说明:
//   - c: Gin 上下文指针，用于存储和传递请求相关的数据。
//
// 返回值:
//   - 配置实例指针 (`*base.BambooConfig`)，如果配置不存在则为 `nil`。
func GetConfig(c *gin.Context) *base.BambooConfig {
	value, exists := c.Get(xConsts.ContextCustomConfig.String())
	if exists {
		return value.(*base.BambooConfig)
	}
	return nil
}

// GetRedisClient 从 Gin 的上下文中获取 Redis 客户端实例。
//
// 如果上下文中存在 Redis 客户端，则返回其指针；如果不存在，则返回 `nil`。
//
// 参数说明:
//   - c: Gin 上下文指针，用于存储和传递请求相关的数据。
//
// 返回值:
//   - Redis 客户端实例指针 (`*redis.Client`)，如果客户端不存在则为 `nil`。
func GetRedisClient(c *gin.Context) *redis.Client {
	value, exists := c.Get(xConsts.ContextRedisClient.String())
	if exists {
		return value.(*redis.Client)
	}
	return nil
}

// GetUserFromContext 从上下文获取用户信息
//
// 从Gin上下文中获取当前认证用户的会话信息。
//
// 参数说明:
//   - c: Gin 上下文指针，用于存储和传递请求相关的数据。
//
// 返回值:
//   - 用户会话信息指针 (`*helper.UserSession`)，如果用户未认证则为 `nil`。
//   - 布尔值，表示是否成功获取到用户信息。
func GetUserFromContext(c *gin.Context) (*helper.UserSession, bool) {
	user, exists := c.Get(constants.ContextKeyUser)
	if !exists {
		return nil, false
	}

	userSession, ok := user.(helper.UserSession)
	if !ok {
		return nil, false
	}

	return &userSession, true
}
