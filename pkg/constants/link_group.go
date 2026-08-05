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
	xSnowflake "github.com/bamboo-services/bamboo-base-go/common/snowflake"
)

// BuiltinGroupHomepageID 内置「首页」分组的保留 ID。
//
// 首页为系统内置的预设位置（不落库）：分组列表接口会以该 ID 注入一条虚拟分组记录，
// 友链引用内置位置时 group_id 保存该保留值，查询返回时再由后端注入对应的虚拟分组对象。
// 该值在雪花 ID 空间之外，数据库真实分组记录不可能占用。
const BuiltinGroupHomepageID xSnowflake.SnowflakeID = 1

// BuiltinGroupFriendsID 内置「友链页」分组的保留 ID（语义同 BuiltinGroupHomepageID）。
const BuiltinGroupFriendsID xSnowflake.SnowflakeID = 2
