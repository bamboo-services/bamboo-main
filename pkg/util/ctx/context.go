/*
 * --------------------------------------------------------------------------------
 * Copyright (c) 2016-NOW(至今) 筱锋
 * Author: 筱锋「xiao_lfeng」(https://www.x-lf.com)
 * --------------------------------------------------------------------------------
 * 许可证声明：版权所有 (c) 2016-2026 筱锋。保留所有权利。
 * 有关MIT许可证的更多信息，请查看项目根目录下的LICENSE文件或访问：
 * https://opensource.org/licenses/MIT
 * --------------------------------------------------------------------------------
 */

package ctxUtil

import (
	"context"

	"github.com/bamboo-services/bamboo-base-go/application/v2"
	"github.com/bamboo-services/bamboo-base-go/foundation/v2/ctxkey"
	"github.com/bamboo-services/bamboo-base-go/foundation/v2/snowflake"
	"github.com/bamboo-services/bamboo-main/internal/models/base"
	"github.com/bamboo-services/bamboo-main/pkg/constants"
	"github.com/gin-gonic/gin"
)

// GetConfig 从上下文中获取 `BambooConfig` 配置实例。
func GetConfig(ctx context.Context) *base.BambooConfig {
	value, err := application.Get[*base.BambooConfig](ctx, constants.ContextCustomConfig)
	if err == nil {
		return value
	}
	return nil
}

// GetSnowflake 从上下文中获取 Snowflake 节点实例。
func GetSnowflake(ctx context.Context) *snowflake.Node {
	node, err := application.Get[*snowflake.Node](ctx, ctxkey.SnowflakeNodeKey)
	if err == nil {
		return node
	}
	return nil
}

// GetUserID 从上下文获取当前认证用户的 Snowflake ID。
func GetUserID(c *gin.Context) (snowflake.SnowflakeID, bool) {
	userID, exists := c.Get(constants.ContextKeyUserID)
	if !exists {
		return 0, false
	}

	userIDSnowflake, ok := userID.(snowflake.SnowflakeID)
	if !ok {
		return 0, false
	}

	return userIDSnowflake, true
}
