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

import "github.com/gogf/gf/v2/os/gtime"

// LinkGroupDTO
//
// 是友链分组数据传输对象
// 它包含了友链分组的相关信息，如分组名称、唯一标识符、描述、排序以及创建和更新时间等。
type LinkGroupDTO struct {
	GroupName      string      `json:"group_name"       description:"分组名称"`    // 分组名称
	GroupUuid      string      `json:"group_uuid"       description:"分组唯一标识符"` // 分组唯一标识符
	GroupDesc      string      `json:"group_desc"       description:"分组描述"`    // 分组描述
	GroupOrder     int         `json:"group_order"      description:"分组排序"`    // 分组排序
	GroupCreatedAt *gtime.Time `json:"group_created_at" description:"分组创建时间"`  // 分组创建时间
	GroupUpdatedAt *gtime.Time `json:"group_updated_at" description:"分组更新时间"`  // 分组更新时间
}
