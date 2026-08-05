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

// BuiltinFancyColorID 内置炫彩颜色的保留 ID。
//
// 炫彩为系统内置的特殊颜色，不落库：颜色列表接口会以该 ID 注入一条虚拟炫彩记录，
// 友链引用炫彩时 color_id 保存该保留值，查询返回时再由后端注入对应的虚拟颜色对象。
// 该值在雪花 ID 空间之外，数据库真实颜色记录不可能占用。
const BuiltinFancyColorID xSnowflake.SnowflakeID = 1
