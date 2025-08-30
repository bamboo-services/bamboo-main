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
	"net/http"

	ctxUtil "bamboo-main/pkg/util/ctx"

	"github.com/gin-gonic/gin"
)

// RequireRole 要求特定角色的中间件
func RequireRole(roles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		user, exists := ctxUtil.GetUserFromContext(c)
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code":    401,
				"message": "未认证的用户",
				"data":    nil,
			})
			c.Abort()
			return
		}

		// 检查用户角色
		hasRole := false
		for _, role := range roles {
			if user.Role == role {
				hasRole = true
				break
			}
		}

		if !hasRole {
			c.JSON(http.StatusForbidden, gin.H{
				"code":    403,
				"message": "权限不足",
				"data":    nil,
			})
			c.Abort()
			return
		}

		c.Next()
	}
}
