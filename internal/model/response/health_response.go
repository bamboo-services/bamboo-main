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

package response

import (
	"time"
)

// HealthResponse 健康检查响应结构
type HealthResponse struct {
	Status    string      `json:"status"`    // 系统状态
	Timestamp time.Time   `json:"timestamp"` // 检查时间
	System    SystemInfo  `json:"system"`    // 系统信息
	Runtime   RuntimeInfo `json:"runtime"`   // 运行时信息
}

// SystemInfo 系统信息
type SystemInfo struct {
	Version     string `json:"version"`     // 应用版本
	Environment string `json:"environment"` // 运行环境
	Platform    string `json:"platform"`    // 运行平台
	GoVersion   string `json:"go_version"`  // Go版本
}

// RuntimeInfo 运行时信息
type RuntimeInfo struct {
	Uptime      string `json:"uptime"`       // 运行时间
	Goroutines  int    `json:"goroutines"`   // 协程数量
	MemoryUsage string `json:"memory_usage"` // 内存使用
	CPUUsage    string `json:"cpu_usage"`    // CPU使用率
}

// PingResponse Ping响应
type PingResponse struct {
	Message   string    `json:"message"`   // 响应消息
	Timestamp time.Time `json:"timestamp"` // 响应时间
}
