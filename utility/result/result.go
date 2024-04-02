// Copyright 2016-2024(ToDate) XResult xiao_lfeng Author(https://blog.x-lf.com). All Rights Reserved.
//
// This Source GetErrorCode Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at.

package result

import (
	"encoding/json"
	"github.com/gogf/gf/v2/os/gtime"
)

type BaseResponse struct {
	Output  string      `json:"output"`
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Time    string      `json:"time"`
	Data    interface{} `json:"data,omitempty"`
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

func (base BaseResponse) String() string {
	base.Time = gtime.Datetime()
	jsonData, err := json.Marshal(base)
	if err != nil {
		return ""
	}
	return string(jsonData)
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

func (baseError ErrorBaseResponse) String() string {
	baseError.Time = gtime.Datetime()
	jsonData, err := json.Marshal(baseError)
	if err != nil {
		panic(jsonData)
		return ""
	}
	return string(jsonData)
}
