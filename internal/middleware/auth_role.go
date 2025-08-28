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