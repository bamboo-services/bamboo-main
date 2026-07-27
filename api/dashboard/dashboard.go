// --------------------------------------------------------------------------------
// Copyright (c) 2016-NOW(至今) 筱锋
// Author: 筱锋「xiao_lfeng」(https://www.x-lf.com)
// --------------------------------------------------------------------------------
// 许可证声明：版权所有 (c) 2016-2026 筱锋。保留所有权利。
// 有关MIT许可证的更多信息，请查看项目根目录下的LICENSE文件或访问：
// https://opensource.org/licenses/MIT
// --------------------------------------------------------------------------------

package apiDashboard

import (
	"time"

	xSnowflake "github.com/bamboo-services/bamboo-base-go/common/snowflake"
)

// RecentApplicationItem 最近友链申请条目
type RecentApplicationItem struct {
	ID        xSnowflake.SnowflakeID `json:"id"`         // 友链ID
	Name      string                 `json:"name"`       // 友链名称
	Avatar    string                 `json:"avatar"`     // 友链头像URL
	URL       string                 `json:"url"`        // 友链URL地址
	CreatedAt time.Time              `json:"created_at"` // 申请时间
}

// StatsResponse 仪表盘统计响应
type StatsResponse struct {
	TotalLinks         int64                   `json:"total_links"`         // 友链总数
	PendingLinks       int64                   `json:"pending_links"`       // 待审核友链数
	ApprovedLinks      int64                   `json:"approved_links"`      // 已通过友链数
	RecentApplications []RecentApplicationItem `json:"recent_applications"` // 最近友链申请列表
}
