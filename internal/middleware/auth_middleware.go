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
	"github.com/bamboo-services/bamboo-base-go/plugins/web-gin/v2/header"
	"github.com/bamboo-services/bamboo-base-go/plugins/web-gin/v2/httpresult"
	logcHelper "github.com/bamboo-services/bamboo-main/internal/logic/helper"
	"github.com/bamboo-services/bamboo-main/pkg/constants"
	"github.com/gin-gonic/gin"
)

// AuthMiddleware 认证中间件，基于本地 Redis 会话校验访问令牌并注入用户上下文。
func AuthMiddleware(c *gin.Context) {
	accessToken := header.GetToken(c, header.HeaderAuthorization)
	if accessToken == "" {
		_ = c.Error(apperr.New(httpresult.Unauthorized, "未检测到访问令牌"))
		c.Abort()
		return
	}

	reqCtx := c.Request.Context()
	m := cache.MustGet(reqCtx)
	logger, _ := application.Get[*logging.Logger](reqCtx, builtin.LoggerKey)

	sessionLogic := logcHelper.NewSessionLogic(m, logger)
	session, found, xErr := sessionLogic.GetUserSession(reqCtx, accessToken)
	if xErr != nil {
		_ = c.Error(xErr)
		c.Abort()
		return
	}
	if !found {
		_ = c.Error(apperr.New(httpresult.Unauthorized, "登录状态无效或已过期"))
		c.Abort()
		return
	}

	c.Set(constants.ContextKeyUserID, session.UserID)
	c.Set(constants.ContextKeyToken, accessToken)

	c.Next()
}
