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
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/os/gtime"
	"regexp"
	"xiaoMain/internal/dao"
	"xiaoMain/internal/model/do"
	"xiaoMain/internal/model/entity"
	"xiaoMain/utility"
)

// MiddleAuthHandler 是用于处理用户授权的中间件。
// 它检查用户的授权信息是否有效。
//
// 参数:
// r: 请求的上下文，用于管理请求的信息。
//
// 返回:
// 无
func MiddleAuthHandler(r *ghttp.Request) {
	getAuthorize, err := utility.GetAuthorizationFromHeader(r.Context())
	if err != nil {
		g.Log().Warning(r.Context(), "[MIDDLE] 用户授权异常|获取授权错误|用户未登录")
		r.SetError(berror.NewError(bcode.UnknownError, "获取授权错误"))
		return

	}
	if getAuthorize == "" {
		g.Log().Warning(r.Context(), "[MIDDLE] 用户授权异常|无授权头|用户未登录")
		r.SetError(berror.NewError(bcode.UnknownError, "无授权头"))
		return
	}
	// 对获取数据进行正则表达式校验
	if hasMatch, _ := regexp.Match(
		`^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-4[0-9a-fA-F]{3}-[89abAB][0-9a-fA-F]{3}-[0-9a-fA-F]{12}$`,
		[]byte(getAuthorize),
	); hasMatch != true {
		g.Log().Warningf(r.Context(), "[MIDDLE] 用户授权异常|授权格式错误|用户未登录 %s", getAuthorize)
		r.SetError(berror.NewError(bcode.UnknownError, "授权格式错误"))
		return
	}
	// 数据库检查
	var tokenInfo *entity.Token
	err = dao.Token.Ctx(r.Context()).Where(do.Token{
		Token: getAuthorize,
	}).Limit(1).OrderDesc("expired_at").Scan(&tokenInfo)
	if err != nil {
		g.Log().Error(r.Context(), "[MIDDLE] 数据库查询错误", err.Error())
		r.SetError(berror.NewError(bcode.UnknownError, "数据库查询错误"))
		return
	}
	// 对数据库进行有效性检查
	if tokenInfo != nil {
		if gtime.Now().Before(tokenInfo.ExpiredAt) {
			g.Log().Noticef(r.Context(), "[MIDDLE] 用户授权有效|用户UUID[%s]", tokenInfo.UserUuid)
			r.Middleware.Next()
		} else {
			g.Log().Warning(r.Context(), "[MIDDLE] 用户授权异常|授权已过期|用户未登录")
			r.SetError(berror.NewError(bcode.UnknownError, "授权已过期"))
			// 删除数据库中的授权信息
			_, _ = dao.Token.Ctx(r.Context()).Where(do.Token{Id: tokenInfo.Id}).Delete()
		}
	} else {
		g.Log().Warning(r.Context(), "[MIDDLE] 用户授权异常|授权不存在|用户未登录")
	}
}
