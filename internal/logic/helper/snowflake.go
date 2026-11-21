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

	"github.com/bamboo-services/bamboo-base-go/foundation/v2/snowflake"
)

// ParseSnowflakeID 将字符串解析为 snowflake.SnowflakeID。
func ParseSnowflakeID(s string) (snowflake.SnowflakeID, error) {
	id, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return 0, err
	}
	return snowflake.SnowflakeID(id), nil
}
