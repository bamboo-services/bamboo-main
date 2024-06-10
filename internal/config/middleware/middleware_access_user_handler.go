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
	"github.com/gogf/gf/v2/net/ghttp"
	"xiaoMain/utility/result"
)

// MiddleAccessUserHandler 是用于处理用户访问的中间件。
// 它检查用户的 IP 地址和 User-Agent 是否为空。
//
// 参数:
// r: 请求的上下文，用于管理请求的信息。
//
// 返回:
// 无
func MiddleAccessUserHandler(r *ghttp.Request) {
	// 继续执行后续中间件
	ctx := r.GetCtx()
	// 获取用户的 IP 地址 以及 User-Agent
	userIP := r.GetClientIp()
	userAgent := r.GetHeader("User-Agent")
	// 两者内容不能为空
	if userIP == "" || userAgent == "" {
		g.Log().Error(ctx, "[MIDDLE] 用户访问异常")
		if userIP == "" {
			g.Log().Error(ctx, "[MIDDLE] 用户 IP 为空")
		}
		if userAgent == "" {
			g.Log().Error(ctx, "[MIDDLE] 用户 User-Agent 为空")
		}
		result.AccessError.Response(r)
	} else {
		g.Log().Noticef(ctx, "[MIDDLE] 访问者 [%s] ", userIP)
		g.Log().Debugf(ctx, "[MIDDLE] User-Agent [%s]", userAgent)
		r.Middleware.Next()
	}
}
