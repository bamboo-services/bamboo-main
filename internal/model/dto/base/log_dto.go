/*
 * --------------------------------------------------------------------------------
 * Copyright (c) 2016-NOW(至今) 筱锋
 * Author: 筱锋「xiao_lfeng」(https://www.x-lf.com)
 * --------------------------------------------------------------------------------
 * 许可证声明：版权所有 (c) 2016-2025 筱锋。保留所有权利。
 * 有关MIT许可证的更多信息，请查看项目根目录下的LICENSE文件或访问：
 * https://opensource.org/licenses/MIT
 * --------------------------------------------------------------------------------
 */

package base

// LogDTO
//
// 定义了日志数据传输对象（DTO）的结构体，用于记录请求的相关信息。
// 该结构体包含了请求的路径、方法、请求体、请求头、
// IP地址、用户代理、引用页面以及日志消息等字段。
type LogDTO struct {
	Path    string `json:"path"`    // 请求路径
	Method  string `json:"method"`  // 请求方法
	Body    string `json:"body"`    // 请求体
	IP      string `json:"ip"`      // 请求IP地址
	UA      string `json:"ua"`      // 用户代理
	Referer string `json:"referer"` // 引用页面
	Msg     string `json:"msg"`     // 日志消息
}
