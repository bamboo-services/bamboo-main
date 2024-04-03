package result

// Copyright 2016-2024(ToDate) XResult xiao_lfeng Author(https://blog.x-lf.com). All Rights Reserved.
//
// This Source GetErrorCode Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at.

// ================================================================================================================
// 常见错误代码定义
// 框架保留了内部错误代码：代码 < xxx30
// ================================================================================================================

var (
	RequestBodyValidationError = ErrorBaseResponse{Output: "RequestBodyValidationError", Code: 40301, Message: "请求正文验证错误"}
	VerificationFailed         = ErrorBaseResponse{Output: "VerificationFailed", Code: 40302, Message: "校验失败"}
	ServerInternalError        = ErrorBaseResponse{Output: "ServerInternalError", Code: 50000, Message: "服务器内部错误", Data: errorData{ErrorMessage: "服务器内部错误"}}
	DatabaseError              = ErrorBaseResponse{Output: "DatabaseError", Code: 50001, Message: "服务器内部错误", Data: errorData{ErrorMessage: "数据库错误"}}
)
