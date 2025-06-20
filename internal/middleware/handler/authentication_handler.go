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

package handler

import (
	"bamboo-main/internal/service"
	"github.com/XiaoLFeng/bamboo-utils/butil"
	"github.com/gogf/gf/v2/net/ghttp"
)

// AuthenticationHandler
//
// 中间件处理函数，用于处理请求的认证逻辑；
// 该函数会在请求处理链中被调用，确保请求在继续处理之前经过认证检查。
func AuthenticationHandler(r *ghttp.Request) {
	getToken := butil.TokenRemoveBearer(r.GetHeader("Authorization"))
	getRefreshToken := r.GetHeader("X-Refresh-Token")

	// 验证用户令牌
	newUserToken, errorCode := service.Token().VerifyAndRefreshUserToken(r.GetCtx(), getToken, butil.Ptr(getRefreshToken))
	if errorCode != nil {
		r.SetError(errorCode)
		return
	}

	// 设置新的用户令牌到请求上下文
	getUserEntity, errorCode := service.User().GetUserSimple(r.GetCtx())
	if errorCode != nil {
		r.SetError(errorCode)
		return
	}
	r.SetCtxVar("UserEntity", getUserEntity)

	// 验证通过
	r.Middleware.Next()

	r.Response.Header().Set("X-New-Refresh-Token", newUserToken.RefreshToken)
}
