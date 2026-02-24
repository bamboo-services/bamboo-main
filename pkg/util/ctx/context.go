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
	"github.com/bamboo-services/bamboo-main/internal/models/base"
	"github.com/bamboo-services/bamboo-main/pkg/constants"

	xSnowflake "github.com/bamboo-services/bamboo-base-go/common/snowflake"
	xCtxUtil "github.com/bamboo-services/bamboo-base-go/common/utility/context"
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
	value, err := xCtxUtil.Get[*base.BambooConfig](c, constants.ContextCustomConfig)
	if err == nil {
		return value
	}
	return nil
}

// GetSnowflake 从 gin.Context 中获取 Snowflake 节点实例。
// 如果上下文中不存在对应的 Snowflake 节点，则返回 nil。
// 参数 c 表示请求上下文 gin.Context。
// 返回值为指向 Snowflake 节点的指针或 nil。
func GetSnowflake(c *gin.Context) *xSnowflake.Node {
	return xCtxUtil.GetSnowflakeNode(c)
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
	value, err := xCtxUtil.GetRDB(c)
	if err == nil {
		return value
	}
	return nil
}

// GetUserID 从上下文获取用户ID
//
// 从Gin上下文中获取当前认证用户的ID(Snowflake ID)。
//
// 参数说明:
//   - c: Gin 上下文指针，用于存储和传递请求相关的数据。
//
// 返回值:
//   - 用户ID(int64)，如果用户未认证则为0。
//   - 布尔值，表示是否成功获取到用户ID。
func GetUserID(c *gin.Context) (int64, bool) {
	userID, exists := c.Get(constants.ContextKeyUserID)
	if !exists {
		return 0, false
	}

	userIDInt64, ok := userID.(int64)
	if !ok {
		return 0, false
	}

	return userIDInt64, true
}
