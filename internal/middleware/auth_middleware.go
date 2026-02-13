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
	"github.com/bamboo-services/bamboo-main/internal/logic"
	"github.com/bamboo-services/bamboo-main/pkg/constants"
	bSdkLogic "github.com/phalanx/beacon-sso-sdk/logic"
	bSdkUtil "github.com/phalanx/beacon-sso-sdk/utility"

	xError "github.com/bamboo-services/bamboo-base-go/error"
	"github.com/gin-gonic/gin"
)

func AuthMiddleware(c *gin.Context) {
	accessToken := bSdkUtil.GetAuthorization(c)
	if accessToken == "" {
		_ = c.Error(xError.NewError(c, xError.Unauthorized, "未检测到 OAuth2 访问令牌", false))
		c.Abort()
		return
	}

	oauthLogic := bSdkLogic.NewBusiness(c)
	userinfo, xErr := oauthLogic.Userinfo(c, accessToken)
	if xErr != nil {
		_ = c.Error(xErr)
		c.Abort()
		return
	}

	authLogic := logic.NewAuthLogic(c)
	localUser, xErr := authLogic.SyncOAuthUser(c, userinfo)
	if xErr != nil {
		_ = c.Error(xErr)
		c.Abort()
		return
	}

	c.Set(constants.ContextKeyUserID, localUser.ID)
	c.Set(constants.ContextKeyToken, accessToken)

	c.Next()
}
