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
	"errors"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/util/gvalid"
	"xiaoMain/utility/result"
)

func MiddleErrorHandler(r *ghttp.Request) {
	r.Middleware.Next()

	ctx := r.GetCtx()
	getError := r.GetError()
	if getError != nil {
		// 对 Error 进行转换
		errorCode := gerror.Code(getError).Code()
		r.Response.ClearBuffer()
		switch errorCode {
		// 校验错误处理
		case 51:
			var validError gvalid.Error
			errors.As(getError, &validError)
			// 获取错误信息
			valid, _ := validError.FirstRule()
			// 解析错误
			resultError := result.RequestBodyValidationError
			resultError.Data.ErrorMessage = "请求体内容出现错误，具体内容请见 Data 值"
			resultError.Data.ErrorData = g.Map{
				"content":     validError.String(),
				"requirement": valid,
			}
			r.Response.WriteJson(resultError.String())
			g.Log().Warningf(ctx, "[Exception] 用户参数信息输入有误：%s", getError.Error())
			break
		case 52:
			result.ServerInternalError.
				SetErrorMessage("数据库内部错误，请询问管理员").
				SetErrorData(g.Map{"userCtx": ctx}).
				Response(r)
			g.Log().Warningf(ctx, "[Exception] 数据库错误：%s", getError.Error())
			break
		case 58:
			g.Log().Errorf(ctx, "[Exception] 未处理错误=> [%v]: %s", errorCode, getError.Error())
			result.ServerInternalError.SetErrorMessage("页面已定义未实现").Response(r)
			break
		default:
			g.Log().Warningf(ctx, "[Exception] 未定义的系统错误=> [%v]: %s", errorCode, getError.Error())
			result.ServerInternalError.SetErrorMessage("未定义的系统错误").Response(r)
		}

		// 结束错误
		r.SetError(nil)
	}

	// 判断是否有自定义抛出
	if r.Response.BufferLength() == 0 {
		result.Success("操作成功", r.GetHandlerResponse()).Response(r)
		g.Log().Notice(ctx, "[RETURN] 执行操作成功，返回结果")
	}
}
