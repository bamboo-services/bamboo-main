package result

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// Copyright 2016-2024(ToDate) XResult xiao_lfeng Author(https://blog.x-lf.com). All Rights Reserved.
//
// This Source GetErrorCode Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at.

type NormalCode interface {
	GetOutput() string
	GetCode() int
	GetMessage() string
	GetTime() string
	GetData() interface{}
}

type ErrorCode interface {
	GetErrorOutput() string
	GetErrorCode() int
	GetErrorMessage() string
	GetErrorTime() string
	GetErrorDetail() string
	GetErrorData() interface{}
}

// ================================================================================================================
// 常见错误代码定义
// 框架保留了内部错误代码：代码 < xxx30
// ================================================================================================================

var (
	RequestBodyValidationError = ErrorBaseResponse{Output: "RequestBodyValidationError", Code: 40301, Message: "请求正文验证错误"}
	ServerInternalError        = ErrorBaseResponse{Output: "ServerInternalError", Code: 50000, Message: "服务器内部错误", Data: errorData{ErrorMessage: "服务器内部错误"}}
)

// NewNormalCode 创建新的常规返回
func NewNormalCode(output string, code int, message string, data interface{}) NormalCode {
	return BaseResponse{
		Output:  output,
		Code:    code,
		Message: message,
		Time:    gtime.Datetime(),
		Data:    data,
	}
}

// NewErrorCode 创建新的错误返回
func NewErrorCode(output string, code int, message string, errorMessage string, data interface{}) ErrorCode {
	return ErrorBaseResponse{
		Output:  output,
		Code:    code,
		Message: message,
		Time:    gtime.Datetime(),
		Data: errorData{
			ErrorMessage: errorMessage,
			ErrorData:    data,
		},
	}
}
