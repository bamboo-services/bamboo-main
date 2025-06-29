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

package utility

// GetBaseDomain
//
// 从给定的URL中提取基础域名，支持处理带或不带协议头的URL，以及以斜杠结尾的情况。
// 如果URL长度小于8，则直接返回原始URL。
// 如果URL以 "https://" 开头，则返回去掉协议头后的域名部分；
// 如果URL以 "http://" 开头，则同样返回去掉协议头后的域名部分；
// 如果URL以斜杠结尾，则返回去掉斜杠后的域名部分；
// 否则，直接返回去掉协议头后的完整URL。
func GetBaseDomain(url string) string {
	if len(url) < 8 {
		return url
	}
	if url[:8] == "https://" {
		if lastChar := len(url) - 1; lastChar >= 0 && url[lastChar] == '/' {
			return url[8:lastChar]
		}
		return url[8:]
	} else if url[:7] == "http://" {
		if lastChar := len(url) - 1; lastChar >= 0 && url[lastChar] == '/' {
			return url[7:lastChar]
		}
		return url[7:]
	}
	if lastChar := len(url) - 1; lastChar >= 0 && url[lastChar] == '/' {
		return url[:lastChar]
	}
	return url
}
