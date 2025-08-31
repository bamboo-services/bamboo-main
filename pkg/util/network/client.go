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

package network

import (
	"net"
	"strings"

	xUtil "github.com/bamboo-services/bamboo-base-go/utility"
	"github.com/gin-gonic/gin"
)

// GetRealClientIP 获取客户端真实IP地址
func GetRealClientIP(c *gin.Context) string {
	// 按优先级检查各种可能的IP头部
	headers := []string{
		"X-Original-Forwarded-For",
		"X-Forwarded-For",
		"X-Real-IP",
		"X-Client-IP",
		"CF-Connecting-IP", // Cloudflare
		"True-Client-IP",   // Akamai and Cloudflare
	}

	for _, header := range headers {
		ip := c.GetHeader(header)
		if ip != "" {
			// X-Forwarded-For 可能包含多个IP，取第一个
			if strings.Contains(ip, ",") {
				ip = strings.TrimSpace(strings.Split(ip, ",")[0])
			}
			// 使用 bamboo-base-go 的工具验证IP是否有效
			if xUtil.IsValidIP(ip) {
				return ip
			}
		}
	}

	// 如果没有找到有效的IP头部，使用连接的远程地址
	ip, _, err := net.SplitHostPort(c.Request.RemoteAddr)
	if err != nil {
		return c.Request.RemoteAddr
	}
	return ip
}