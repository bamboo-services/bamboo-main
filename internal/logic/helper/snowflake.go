// --------------------------------------------------------------------------------
// Copyright (c) 2016-NOW(至今) 筱锋
// Author: 筱锋「xiao_lfeng」(https://www.x-lf.com)
// --------------------------------------------------------------------------------
// 许可证声明：版权所有 (c) 2016-2026 筱锋。保留所有权利。
// 有关MIT许可证的更多信息，请查看项目根目录下的LICENSE文件或访问：
// https://opensource.org/licenses/MIT
// --------------------------------------------------------------------------------

package logcHelper

import (
	"strconv"

	xSnowflake "github.com/bamboo-services/bamboo-base-go/common/snowflake"
)

// ParseSnowflakeID 将字符串解析为 xSnowflake.SnowflakeID。
//
// 用于 handler 层路径参数到实体主键类型的统一转换，
// 替代直接 strconv.ParseInt 返回 int64 后再做类型转换的写法。
func ParseSnowflakeID(s string) (xSnowflake.SnowflakeID, error) {
	id, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return 0, err
	}
	return xSnowflake.SnowflakeID(id), nil
}
