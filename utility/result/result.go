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

package result

import (
	"encoding/json"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/os/gtime"
)

type BaseResponse struct {
	Output  string      `json:"output"`
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Time    string      `json:"time"`
	Data    interface{} `json:"data,omitempty"`
}

type BaseResponseNoData struct {
	Output  string `json:"output"`
	Code    int    `json:"code"`
	Message string `json:"message"`
	Time    string `json:"time"`
}

type ErrorBaseResponse struct {
	Output  string    `json:"output"`
	Code    int       `json:"code"`
	Message string    `json:"message"`
	Time    string    `json:"time"`
	Data    errorData `json:"data"`
}

type errorData struct {
	ErrorMessage string      `json:"errorMessage"`
	ErrorData    interface{} `json:"errorData,omitempty"`
}

/*
 * 正确返回
 */

// GetCode 返回当前的状态码
func (base BaseResponse) GetCode() int {
	return base.Code
}

// GetOutput 返回当前错误代码的整数。
func (base BaseResponse) GetOutput() string {
	return base.Output
}

// GetMessage 返回当前错误代码的简短消息。
func (base BaseResponse) GetMessage() string {
	return base.Message
}

// GetTime 返回错误代码的时间
func (base BaseResponse) GetTime() string {
	return base.Time
}

// GetData 获取错误代码的数据
// 仅在填写之后才有对应的数据
func (base BaseResponse) GetData() interface{} {
	return base.Data
}

// String 返回当前样式的 String 类型
func (base BaseResponse) String() string {
	base.Time = gtime.Datetime()
	jsonData, err := json.Marshal(base)
	if err != nil {
		return ""
	}
	return string(jsonData)
}

// Response 返回当前样式的 Response 类型
func (base BaseResponse) Response(r *ghttp.Request) {
	base.Time = gtime.Datetime()
	r.Response.Status = 200
	r.Response.WriteJson(base)
}

// Get

/*
 * 错误相关
 */

// GetErrorOutput 返回当前错误代码的整数。
func (baseError ErrorBaseResponse) GetErrorOutput() string {
	return baseError.Output
}

// GetErrorCode 返回当前错误代码的整数。
func (baseError ErrorBaseResponse) GetErrorCode() int {
	return baseError.Code
}

// GetErrorMessage 返回当前错误代码的简短消息。
func (baseError ErrorBaseResponse) GetErrorMessage() string {
	return baseError.Message
}

// GetErrorTime 返回错误代码的时间
func (baseError ErrorBaseResponse) GetErrorTime() string {
	return baseError.Time
}

// GetErrorDetail 返回错误代码的详细消息
func (baseError ErrorBaseResponse) GetErrorDetail() string {
	return baseError.Data.ErrorMessage
}

// GetErrorData 获取错误代码的数据
// 仅在填写之后才有对应的数据
func (baseError ErrorBaseResponse) GetErrorData() interface{} {
	return baseError.Data.ErrorData
}

// String 返回当前错误代码的字符串表示形式
func (baseError ErrorBaseResponse) String() string {
	baseError.Time = gtime.Datetime()
	jsonData, err := json.Marshal(baseError)
	if err != nil {
		panic(jsonData)
	}
	return string(jsonData)
}

// Response 返回当前错误代码的字符串表示形式
func (baseError ErrorBaseResponse) Response(r *ghttp.Request) {
	baseError.Time = gtime.Datetime()
	// 检查 baseError 是否存在 errorMessage
	if baseError.Data.ErrorMessage == "" {
		baseError.Data.ErrorMessage = baseError.Message
	}
	r.Response.Status = baseError.Code / 100
	r.Response.WriteJson(baseError)
}

// SetErrorMessage 设置错误消息
func (baseError ErrorBaseResponse) SetErrorMessage(msg string) ErrorBaseResponse {
	baseError.Data.ErrorMessage = msg
	return baseError
}

// SetErrorData 设置错误数据
func (baseError ErrorBaseResponse) SetErrorData(data any) ErrorBaseResponse {
	baseError.Data.ErrorData = data
	return baseError
}
