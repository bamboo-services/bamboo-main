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
	"github.com/bamboo-services/bamboo-main/internal/repository"
	"github.com/bamboo-services/bamboo-main/pkg/constants"
	ctxUtil "github.com/bamboo-services/bamboo-main/pkg/util/ctx"

	xError "github.com/bamboo-services/bamboo-base-go/error"
	xCtxUtil "github.com/bamboo-services/bamboo-base-go/utility/ctxutil"
	"github.com/gin-gonic/gin"
)

// RequireRole 要求特定角色的中间件
func RequireRole(roles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 获取用户UUID
		userID, exists := ctxUtil.GetUserID(c)
		if !exists {
			_ = c.Error(xError.NewError(c, xError.Unauthorized, "未认证的用户", false))
			return
		}

		// 从数据库获取用户信息
		db := xCtxUtil.MustGetDB(c)
		rdb := xCtxUtil.MustGetRDB(c)
		if db == nil {
			_ = c.Error(xError.NewError(c, xError.DatabaseError, "数据库连接异常", false))
			return
		}

		repo := repository.NewSystemUserRepo(db, rdb)
		user, found, xErr := repo.GetByID(c, userID)
		if xErr != nil {
			_ = c.Error(xError.NewError(c, xError.DatabaseError, "用户信息查询失败", false, xErr))
			return
		}
		if !found || user.Status != constants.StatusActive {
			_ = c.Error(xError.NewError(c, xError.NotFound, "用户不存在或已被禁用", false))
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
			_ = c.Error(xError.NewError(c, xError.PermissionDenied, "权限不足", false))
			return
		}

		c.Next()
	}
}
