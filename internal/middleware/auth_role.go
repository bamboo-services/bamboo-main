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

package middleware

import (
	"github.com/bamboo-services/bamboo-main/internal/repository"
	"github.com/bamboo-services/bamboo-main/pkg/constants"
	ctxUtil "github.com/bamboo-services/bamboo-main/pkg/util/ctx"

	xError "github.com/bamboo-services/bamboo-base-go/common/error"
	xCtxUtil "github.com/bamboo-services/bamboo-base-go/major/utility/context"
	"github.com/gin-gonic/gin"
)

// RequireAdmin 要求当前登录用户为系统唯一管理员。
//
// 管理员身份由 bm_system 配置表的 system.admin.id 唯一标记：先校验用户存在且启用，
// 再经 SystemRepo.IsAdmin 比对当前用户 ID，非管理员一律拒绝（403）。
func RequireAdmin(c *gin.Context) {
	// 获取当前认证用户 ID
	userID, exists := ctxUtil.GetUserID(c)
	if !exists {
		_ = c.Error(xError.NewError(c, xError.Unauthorized, "未认证的用户", false))
		c.Abort()
		return
	}

	// 从上下文获取数据库与缓存管理器
	db := xCtxUtil.MustGetDB(c)
	m := xCtxUtil.MustGetCacheManager(c)
	if db == nil {
		_ = c.Error(xError.NewError(c, xError.DatabaseError, "数据库连接异常", false))
		c.Abort()
		return
	}

	// 校验用户存在且处于启用状态
	userRepo := repository.NewSystemUserRepo(db, m)
	user, found, xErr := userRepo.GetByID(c, userID)
	if xErr != nil {
		_ = c.Error(xError.NewError(c, xError.DatabaseError, "用户信息查询失败", false, xErr))
		c.Abort()
		return
	}
	if !found || user.Status != constants.StatusActive {
		_ = c.Error(xError.NewError(c, xError.NotFound, "用户不存在或已被禁用", false))
		c.Abort()
		return
	}

	// 校验是否为系统唯一管理员
	systemRepo := repository.NewSystemRepo(db, m)
	isAdmin, xErr := systemRepo.IsAdmin(c, userID)
	if xErr != nil {
		_ = c.Error(xError.NewError(c, xError.DatabaseError, "管理员身份校验失败", false, xErr))
		c.Abort()
		return
	}
	if !isAdmin {
		_ = c.Error(xError.NewError(c, xError.PermissionDenied, "权限不足", false))
		c.Abort()
		return
	}

	c.Next()
}
