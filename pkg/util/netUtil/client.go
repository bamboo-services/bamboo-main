/*
 * --------------------------------------------------------------------------------
 * Copyright (c) 2016-NOW(至今) 筱锋
 * Author: 筱锋「xiao_lfeng」(https://www.x-lf.com)
 * --------------------------------------------------------------------------------
 * 许可证声明：版权所有 (c) 2016-2026 筱锋。保留所有权利。
 * 有关MIT许可证的更多信息，请查看项目根目录下的LICENSE文件或访问：
 * https://opensource.org/licenses/MIT
 * --------------------------------------------------------------------------------
 */

package netUtil

import (
	"net"
	"strings"

	"github.com/bamboo-services/bamboo-base-go/foundation/v2/toolkit"
	"github.com/gin-gonic/gin"
)

// GetClientIP 获取客户端真实IP地址
func GetClientIP(c *gin.Context) string {
	headers := []string{
		"X-Original-Forwarded-For",
		"X-Forwarded-For",
		"X-Real-IP",
		"X-Client-IP",
		"CF-Connecting-IP",
		"True-Client-IP",
	}

	for _, header := range headers {
		ip := c.GetHeader(header)
		if ip != "" {
			if strings.Contains(ip, ",") {
				ip = strings.TrimSpace(strings.Split(ip, ",")[0])
			}
			if toolkit.IsIP(ip) {
				return ip
			}
		}
	}

	ip, _, err := net.SplitHostPort(c.Request.RemoteAddr)
	if err != nil {
		return c.Request.RemoteAddr
	}
	return ip
}
