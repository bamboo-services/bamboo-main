package middleware

import (
	"errors"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/os/glog"
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
			glog.Warningf(ctx, "[Exception] 用户参数信息输入有误：%s", getError.Error())
		case 58:
			glog.Errorf(ctx, "[Exception] 未处理错误=> [%v]: %s", errorCode, getError.Error())
			r.Response.WriteJson(result.ServerInternalError.String())
		default:
			glog.Warningf(ctx, "[Exception] 未定义的系统错误=> [%v]: %s", errorCode, getError.Error())
		}

		// 结束错误
		r.SetError(nil)
	}
}
