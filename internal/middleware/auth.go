package middleware

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"bamboo-main/internal/model/entity"
	"bamboo-main/pkg/constants"
	"bamboo-main/pkg/startup"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

// UserSession 用户会话结构
type UserSession struct {
	UserUUID  string    `json:"user_uuid"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	Role      string    `json:"role"`
	LoginAt   time.Time `json:"login_at"`
	ExpireAt  time.Time `json:"expire_at"`
}

// AuthMiddleware 认证中间件
func AuthMiddleware(reg *startup.Reg) gin.HandlerFunc {
	return func(c *gin.Context) {
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

		// 从 Redis 获取用户会话
		redisKey := fmt.Sprintf(constants.AuthTokenPrefix, token)
		sessionData, err := reg.Rdb.Get(c.Request.Context(), redisKey).Result()
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
		var session UserSession
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
			reg.Rdb.Del(c.Request.Context(), redisKey)
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
}

// GetUserFromContext 从上下文获取用户信息
func GetUserFromContext(c *gin.Context) (*UserSession, bool) {
	user, exists := c.Get(constants.ContextKeyUser)
	if !exists {
		return nil, false
	}
	
	userSession, ok := user.(UserSession)
	if !ok {
		return nil, false
	}
	
	return &userSession, true
}

// RequireRole 要求特定角色的中间件
func RequireRole(roles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		user, exists := GetUserFromContext(c)
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

// CreateUserSession 创建用户会话
func CreateUserSession(reg *startup.Reg, user *entity.SystemUser, token string) error {
	session := UserSession{
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
	err = reg.Rdb.Set(reg.Serv.Context, redisKey, sessionData, 24*time.Hour).Err()
	if err != nil {
		return err
	}

	return nil
}

// DeleteUserSession 删除用户会话
func DeleteUserSession(reg *startup.Reg, token string) error {
	redisKey := fmt.Sprintf(constants.AuthTokenPrefix, token)
	return reg.Rdb.Del(reg.Serv.Context, redisKey).Err()
}