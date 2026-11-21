// --------------------------------------------------------------------------------
// Copyright (c) 2016-NOW(至今) 筱锋
// Author: 筱锋「xiao_lfeng」(https://www.x-lf.com)
// --------------------------------------------------------------------------------
// 许可证声明：版权所有 (c) 2016-2026 筱锋。保留所有权利。
// 有关MIT许可证的更多信息，请查看项目根目录下的LICENSE文件或访问：
// https://opensource.org/licenses/MIT
// --------------------------------------------------------------------------------

package constants

import (
	"github.com/bamboo-services/bamboo-base-go/foundation/v2/snowflake"
)

// BuiltinGroupInvalidID 内置「已失效」分组的保留 ID。
//
// 已失效为系统内置的语义分组（不落库）：友链标记失效时自动归入该保留 ID，
// 名字与描述由 bm_system 配置（group.builtin.invalid.name / .description）热修改，
// 查询返回时再由后端注入对应的虚拟分组对象。
// 该值在雪花 ID 空间之外，数据库真实分组记录不可能占用。
const BuiltinGroupInvalidID snowflake.SnowflakeID = 1
