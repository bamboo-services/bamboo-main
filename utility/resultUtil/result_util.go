package resultUtil

import (
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/os/gtime"
	"golang.org/x/net/context"
)

// SuccessNoData
// 成功代码输出
func SuccessNoData(ctx context.Context, message *string) {
	if message == nil {
		*message = "操作成功"
	}
	getHttp := ghttp.RequestFromCtx(ctx)
	getHttp.Response.WriteHeader(200)
	getHttp.Response.WriteJson(BaseResponse{
		Output:  "Success",
		Code:    200,
		Message: *message,
		Time:    gtime.Datetime(),
		Data:    nil,
	})
}

func SuccessData(ctx context.Context, data interface{}, message *string) {
	if message == nil {
		*message = "操作成功"
	}
	getHttp := ghttp.RequestFromCtx(ctx)
	getHttp.Response.WriteHeader(200)
	getHttp.Response.WriteJson(BaseResponse{
		Output:  "Success",
		Code:    200,
		Message: *message,
		Time:    gtime.Datetime(),
		Data:    data,
	})
}

// ErrorNoData
// 错误代码输出
func ErrorNoData(ctx context.Context, code ErrorCode, errorMessage string) {
	var getHttp = ghttp.RequestFromCtx(ctx)
	getHttp.Response.WriteHeader(int(code.code / 100))
	getHttp.Response.WriteJson(ErrorBaseResponse{
		Output:  code.output,
		Code:    code.code,
		Message: code.message,
		Time:    gtime.Datetime(),
		Data: errorData{
			ErrorMessage: errorMessage,
			ErrorData:    nil,
		},
	})
}

// ErrorData
// 错误代码输出，包含数据
func ErrorData(ctx context.Context, code ErrorCode, errorMessage string, data interface{}) {
	getHttp := ghttp.RequestFromCtx(ctx)
	getHttp.Response.WriteHeader(int(code.code / 100))
	getHttp.Response.WriteJson(ErrorBaseResponse{
		Output:  code.output,
		Code:    code.code,
		Message: code.message,
		Time:    gtime.Datetime(),
		Data: errorData{
			ErrorMessage: errorMessage,
			ErrorData:    data,
		},
	})
}
