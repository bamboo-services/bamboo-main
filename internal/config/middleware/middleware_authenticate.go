/*
 * --------------------------------------------------------------------------------
 * Copyright (c) 2016-NOW(至今) 筱锋
 * Author: 筱锋(https://www.x-lf.com)
 *
 * 本文件包含 XiaoMain 的源代码，该项目的所有源代码均遵循MIT开源许可证协议。
 * --------------------------------------------------------------------------------
 * 许可证声明：
 *
 * 版权所有 (c) 2016-2024 筱锋。保留所有权利。
 *
 * 本软件是“按原样”提供的，没有任何形式的明示或暗示的保证，包括但不限于
 * 对适销性、特定用途的适用性和非侵权性的暗示保证。在任何情况下，
 * 作者或版权持有人均不承担因软件或软件的使用或其他交易而产生的、
 * 由此引起的或以任何方式与此软件有关的任何索赔、损害或其他责任。
 *
 * 使用本软件即表示您了解此声明并同意其条款。
 *
 * 有关MIT许可证的更多信息，请查看项目根目录下的LICENSE文件或访问：
 * https://opensource.org/licenses/MIT
 * --------------------------------------------------------------------------------
 * 免责声明：
 *
 * 使用本软件的风险由用户自担。作者或版权持有人在法律允许的最大范围内，
 * 对因使用本软件内容而导致的任何直接或间接的损失不承担任何责任。
 * --------------------------------------------------------------------------------
 */

package middleware

import (
	"github.com/bamboo-services/bamboo-utils/bcode"
	"github.com/bamboo-services/bamboo-utils/berror"
	"github.com/bamboo-services/bamboo-utils/bresult"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/os/gtime"
	"xiaoMain/internal/dao"
	"xiaoMain/internal/model/do"
	"xiaoMain/internal/model/entity"
	"xiaoMain/utility"
)

// MiddleAuthenticate
//
// # 身份验证器
//
// 获取路由信息，检查路由是否允许匿名访问，如果不允许匿名访问，则检查是否有登录信息
// 如果没有登录信息，则返回未授权信息
func MiddleAuthenticate(r *ghttp.Request) {
	// 获取用户的 IP 地址 以及 User-Agent
	userIP := r.GetClientIp()
	userAgent := r.GetHeader("User-Agent")
	// 两者内容不能为空
	if userIP == "" || userAgent == "" {
		g.Log().Error(r.Context(), "[REQU] 用户访问异常")
		if userIP == "" {
			g.Log().Error(r.Context(), "[REQU] 用户 IP 为空")
		}
		if userAgent == "" {
			g.Log().Error(r.Context(), "[REQU] 用户 User-Agent 为空")
		}
		bresult.WriteResponse(r.GetCtx(), bcode.OperationFailed, "用户访问异常", nil)
	} else {
		g.Log().Noticef(r.Context(), "[REQU] User-IP [%s] ", userIP)
		g.Log().Debugf(r.Context(), "[REQU] User-Agent [%s]", userAgent)
		// 获取路由信息
		route := r.Router.Uri
		method := r.Method
		// 检查路由是否权限鉴别
		if isAuthenticate(route, method) {
			// 检查是否有登录信息
			if authHandler(r) {
				r.Middleware.Next()
			} else {
				r.SetError(berror.NewError(bcode.Unauthorized, "用户未授权"))
			}
		} else {
			r.Middleware.Next()
		}
	}
}

// isAuthenticate
//
// # 是否需要权限鉴别
//
// 检查路由是否需要权限鉴别
//
// # 参数
//   - route: 路由路由
//   - method: 请求方法
//
// # 返回
//   - bool: 是否需要权限鉴别
func isAuthenticate(route, method string) bool {
	// 需要权限鉴别访问的路由
	routeList := []map[string]string{
		{"method": "PUT", "route": "/api/v1/user/change-password"},
		{"method": "PUT", "route": "/api/v1/user/reset-password"},
		{"method": "PUT", "route": "/api/v1/info"},
		{"method": "POST", "route": "/api/v1/link/add/color"},
		{"method": "POST", "route": "/api/v1/link/add/location"},
		{"method": "POST", "route": "/api/v1/link/admin"},
	}
	for _, getRoute := range routeList {
		if getRoute["route"] == route && getRoute["method"] == method {
			return true
		}
	}
	return false
}

// authHandler
//
// # 用户授权处理
//
// 获取授权信息，检查授权信息是否有效
// 如果授权信息有效，则返回 true，否则返回 false
//
// # 参数
//   - r: 请求信息
//
// # 返回
//   - bool: 是否有效
func authHandler(r *ghttp.Request) bool {
	getAuthorize, err := utility.GetAuthorizationFromHeader(r.Context())
	if err != nil {
		g.Log().Warning(r.Context(), "[REQU] 用户授权异常|获取授权错误|用户未登录")
		r.SetError(berror.NewError(bcode.UnknownError, "获取授权错误"))
		return false
	}
	if getAuthorize == "" {
		g.Log().Warning(r.Context(), "[REQU] 用户授权异常|无授权头|用户未登录")
		r.SetError(berror.NewError(bcode.UnknownError, "无授权头"))
		return false
	}
	// 数据库检查
	var tokenInfo *entity.Token
	err = dao.Token.Ctx(r.Context()).Where(do.Token{
		Token: getAuthorize,
	}).Limit(1).OrderDesc("expired_at").Scan(&tokenInfo)
	if err != nil {
		g.Log().Error(r.Context(), "[REQU] 数据库查询错误", err.Error())
		r.SetError(berror.NewError(bcode.UnknownError, "数据库查询错误"))
		return false
	}
	// 对数据库进行有效性检查
	if tokenInfo != nil {
		if gtime.Now().Before(tokenInfo.ExpiredAt) {
			g.Log().Noticef(r.Context(), "[REQU] 用户授权有效|用户UUID[%s]", tokenInfo.UserUuid)
			return true
		} else {
			g.Log().Warning(r.Context(), "[REQU] 用户授权异常|授权已过期|用户未登录")
			r.SetError(berror.NewError(bcode.UnknownError, "授权已过期"))
			// 删除数据库中的授权信息
			_, _ = dao.Token.Ctx(r.Context()).Where(do.Token{Id: tokenInfo.Id}).Delete()
		}
	} else {
		g.Log().Warning(r.Context(), "[REQU] 用户授权异常|授权不存在|用户未登录")
	}
	return false
}
