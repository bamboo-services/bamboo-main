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
	logcHelper "github.com/bamboo-services/bamboo-main/internal/logic/helper"
	"github.com/bamboo-services/bamboo-main/pkg/constants"

	xError "github.com/bamboo-services/bamboo-base-go/common/error"
	xHttp "github.com/bamboo-services/bamboo-base-go/defined/http"
	xCtxUtil "github.com/bamboo-services/bamboo-base-go/major/utility/context"
	"github.com/gin-gonic/gin"
)

// AuthMiddleware 认证中间件，基于本地 Redis 会话校验访问令牌并注入用户上下文。
//
// 令牌来源无关：密码登录与 SSO 登录均在登录时写入本地会话，此处统一读取会话，
// 不再每请求回查 SSO Userinfo。
func AuthMiddleware(c *gin.Context) {
	accessToken := xHttp.GetToken(c, xHttp.HeaderAuthorization)
	if accessToken == "" {
		_ = c.Error(xError.NewError(c, xError.Unauthorized, "未检测到访问令牌", false))
		c.Abort()
		return
	}

	sessionLogic := logcHelper.NewSessionLogic(xCtxUtil.MustGetCacheManager(c))
	session, found, xErr := sessionLogic.GetUserSession(c.Request.Context(), accessToken)
	if xErr != nil {
		_ = c.Error(xErr)
		c.Abort()
		return
	}
	if !found {
		_ = c.Error(xError.NewError(c, xError.Unauthorized, "登录状态无效或已过期", false))
		c.Abort()
		return
	}

	c.Set(constants.ContextKeyUserID, session.UserID)
	c.Set(constants.ContextKeyToken, accessToken)

	c.Next()
}
