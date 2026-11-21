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
	"github.com/bamboo-services/bamboo-base-go/application/v2"
	"github.com/bamboo-services/bamboo-base-go/application/v2/builtin"
	"github.com/bamboo-services/bamboo-base-go/foundation/v2/apperr"
	"github.com/bamboo-services/bamboo-base-go/foundation/v2/logging"
	"github.com/bamboo-services/bamboo-base-go/plugins/cache/v2"
	"github.com/bamboo-services/bamboo-base-go/plugins/database/v2"
	"github.com/bamboo-services/bamboo-base-go/plugins/web-gin/v2/httpresult"
	"github.com/bamboo-services/bamboo-main/internal/repository"
	"github.com/bamboo-services/bamboo-main/pkg/constants"
	ctxUtil "github.com/bamboo-services/bamboo-main/pkg/util/ctx"
	"github.com/gin-gonic/gin"
)

// RequireAdmin 要求当前登录用户为系统唯一管理员。
func RequireAdmin(c *gin.Context) {
	userID, exists := ctxUtil.GetUserID(c)
	if !exists {
		_ = c.Error(apperr.New(httpresult.Unauthorized, "未认证的用户"))
		c.Abort()
		return
	}

	reqCtx := c.Request.Context()
	db, dbErr := database.Get(reqCtx)
	if dbErr != nil || db == nil {
		_ = c.Error(apperr.New(database.DatabaseError, "数据库连接异常", dbErr))
		c.Abort()
		return
	}
	m := cache.MustGet(reqCtx)
	logger, _ := application.Get[*logging.Logger](reqCtx, builtin.LoggerKey)

	userRepo := repository.NewSystemUserRepo(db, m, logger)
	user, found, xErr := userRepo.GetByID(reqCtx, userID)
	if xErr != nil {
		_ = c.Error(apperr.New(database.DatabaseError, "用户信息查询失败", xErr))
		c.Abort()
		return
	}
	if !found || user.Status != constants.StatusActive {
		_ = c.Error(apperr.New(httpresult.NotFound, "用户不存在或已被禁用"))
		c.Abort()
		return
	}

	systemRepo := repository.NewSystemRepo(db, m, logger)
	isAdmin, xErr := systemRepo.IsAdmin(reqCtx, userID)
	if xErr != nil {
		_ = c.Error(apperr.New(database.DatabaseError, "管理员身份校验失败", xErr))
		c.Abort()
		return
	}
	if !isAdmin {
		_ = c.Error(apperr.New(httpresult.PermissionDenied, "权限不足"))
		c.Abort()
		return
	}

	c.Next()
}
