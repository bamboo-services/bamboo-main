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

package apiLink

import (
	"encoding/json"

	xSnowflake "github.com/bamboo-services/bamboo-base-go/common/snowflake"
	"github.com/bamboo-services/bamboo-main/internal/entity"
	"github.com/bamboo-services/bamboo-main/internal/models/base"
)

// LinkIDRequest 友情链接路径参数
type LinkIDRequest struct {
	ID xSnowflake.SnowflakeID `uri:"id" binding:"required"`
}

// NullableSnowflakeID 可空雪花 ID 请求字段（三态）：
//   - JSON 省略：未提供，业务按「保持原值」处理
//   - JSON null：显式清空（如取消颜色/分组选择）
//   - JSON 值：设置新值
//
// 用于区分「不传」与「传 null」，避免取消选择被零值哨兵吞掉。
type NullableSnowflakeID struct {
	present bool
	value   *xSnowflake.SnowflakeID
}

// UnmarshalJSON 解析 JSON 值；null 与省略通过 Provided 区分
func (n *NullableSnowflakeID) UnmarshalJSON(data []byte) error {
	n.present = true
	if string(data) == "null" {
		n.value = nil
		return nil
	}
	var id xSnowflake.SnowflakeID
	if err := json.Unmarshal(data, &id); err != nil {
		return err
	}
	n.value = &id
	return nil
}

// Provided 字段是否在请求中显式出现（null 与值均为 true，省略为 false）
func (n NullableSnowflakeID) Provided() bool { return n.present }

// Value 返回解析值；显式 null 时为 nil
func (n NullableSnowflakeID) Value() *xSnowflake.SnowflakeID { return n.value }

// FriendAddRequest 添加友情链接请求
type FriendAddRequest struct {
	LinkName        string              `json:"link_name" binding:"required,min=1,max=100" example:"示例网站"`
	LinkURL         string              `json:"link_url" binding:"required,url,max=500" example:"https://example.com"`
	LinkAvatar      string              `json:"link_avatar" binding:"omitempty,url,max=500" example:"https://example.com/avatar.jpg"`
	LinkRSS         string              `json:"link_rss" binding:"omitempty,url,max=500" example:"https://example.com/rss.xml"`
	LinkDesc        string              `json:"link_desc" binding:"omitempty,max=500" example:"这是一个示例网站"`
	LinkEmail       string              `json:"link_email" binding:"omitempty,email,max=100" example:"admin@example.com"`
	LinkGroupID     NullableSnowflakeID `json:"link_group_id" binding:"omitempty" example:"1"`
	LinkColorID     NullableSnowflakeID `json:"link_color_id" binding:"omitempty" example:"1"`
	LinkOrder       int                 `json:"link_order" binding:"omitempty,min=0" example:"0"`
	LinkLevel       int                 `json:"link_level" binding:"omitempty,oneof=0 1 2 3" example:"0"`
	LinkApplyRemark string              `json:"link_apply_remark" binding:"omitempty,max=500" example:"申请友链"`
}

// FriendUpdateRequest 更新友情链接请求
type FriendUpdateRequest struct {
	LinkName        string              `json:"link_name" binding:"omitempty,min=1,max=100" example:"示例网站"`
	LinkURL         string              `json:"link_url" binding:"omitempty,url,max=500" example:"https://example.com"`
	LinkAvatar      string              `json:"link_avatar" binding:"omitempty,url,max=500" example:"https://example.com/avatar.jpg"`
	LinkRSS         string              `json:"link_rss" binding:"omitempty,url,max=500" example:"https://example.com/rss.xml"`
	LinkDesc        string              `json:"link_desc" binding:"omitempty,max=500" example:"这是一个示例网站"`
	LinkEmail       string              `json:"link_email" binding:"omitempty,email,max=100" example:"admin@example.com"`
	LinkGroupID     NullableSnowflakeID `json:"link_group_id" binding:"omitempty" example:"1"`
	LinkColorID     NullableSnowflakeID `json:"link_color_id" binding:"omitempty" example:"1"`
	LinkOrder       *int                `json:"link_order" binding:"omitempty,min=0" example:"0"`
	LinkLevel       *int                `json:"link_level" binding:"omitempty,oneof=0 1 2 3" example:"0"`
	LinkApplyRemark string              `json:"link_apply_remark" binding:"omitempty,max=500" example:"申请友链"`
}

// FriendQueryRequest 查询友情链接请求
type FriendQueryRequest struct {
	Page        int                    `form:"page" binding:"omitempty,min=1" example:"1"`
	PageSize    int                    `form:"page_size" binding:"omitempty,min=1,max=100" example:"10"`
	LinkName    string                 `form:"link_name" binding:"omitempty,max=100" example:"示例"`
	LinkStatus  *int                   `form:"link_status" binding:"omitempty,oneof=0 1 2 3 4" example:"1"`
	LinkFail    *int                   `form:"link_fail" binding:"omitempty,oneof=0 1" example:"0"`
	LinkAnomaly *bool                  `form:"link_anomaly" binding:"omitempty" example:"true"` // 异常过滤：status 非 0/1 或已失效（true）
	LinkGroupID xSnowflake.SnowflakeID `form:"link_group_id" binding:"omitempty,number" example:"1"`
	SortBy      string                 `form:"sort_by" binding:"omitempty,oneof=created_at updated_at link_order link_name" example:"created_at"`
	SortOrder   string                 `form:"sort_order" binding:"omitempty,oneof=asc desc" example:"desc"`
}

// FriendStatusRequest 更新友情链接状态请求
type FriendStatusRequest struct {
	LinkStatus       int    `json:"link_status" binding:"required,oneof=0 1 2 3 4" example:"1"`
	LinkReviewRemark string `json:"link_review_remark" binding:"omitempty,max=500" example:"审核通过"`
}

// FriendFailRequest 更新友情链接失效状态请求
type FriendFailRequest struct {
	LinkFail       int    `json:"link_fail" binding:"required,oneof=0 1" example:"1"`
	LinkFailReason string `json:"link_fail_reason" binding:"omitempty,max=500" example:"网站无法访问"`
}

// FriendAddResponse 添加友情链接响应
type FriendAddResponse struct {
	entity.LinkFriend
}

// FriendUpdateResponse 更新友情链接响应
type FriendUpdateResponse struct {
	entity.LinkFriend
}

// FriendDetailResponse 友情链接详情响应
type FriendDetailResponse struct {
	entity.LinkFriend
}

// FriendListResponse 友情链接列表响应（附带待审核/异常计数，供管理端入口徽章展示）
type FriendListResponse struct {
	base.PaginationResponse[entity.LinkFriend]
	PendingCount int64 `json:"pending_count"` // 待审核友链数量
	AnomalyCount int64 `json:"anomaly_count"` // 异常友链数量（status 非 0/1 或已失效）
}

// FriendPublicResponse 公开友情链接响应
type FriendPublicResponse struct {
	Links []entity.LinkFriend `json:"links"`
}

// FriendSortItem 友链排序条目
//
// items 数组顺序 = 目标全局展示顺序（后端按此序重写全局 sort_order 为 0..N-1）。
// GroupID 三态（复用 NullableSnowflakeID）：省略=保持原分组，null=移入未分组，值=移入该分组。
type FriendSortItem struct {
	ID      xSnowflake.SnowflakeID `json:"id" binding:"required"`        // 友链ID
	GroupID NullableSnowflakeID    `json:"group_id" binding:"omitempty"` // 目标分组ID（三态）
}

// FriendSortRequest 友链批量排序/位置请求
type FriendSortRequest struct {
	Items []FriendSortItem `json:"items" binding:"required,min=1,max=500,dive"` // 排序条目列表，顺序即目标全局展示顺序
}

// FriendSortResponse 友链批量排序/位置响应
type FriendSortResponse struct {
	Count int `json:"count"` // 更新的友链数量
}

// FriendApplyRequest 访客自助申请友情链接请求
//
// 面向游客与登录用户的公开申请入口：仅需站点基础信息，联系邮箱必填（用于确认友链归属），
// 分组/颜色/级别/排序等管理员专属字段不在此开放，由管理员审核时分配。
type FriendApplyRequest struct {
	LinkName        string              `json:"link_name" binding:"required,min=1,max=100" example:"示例网站"`
	LinkURL         string              `json:"link_url" binding:"required,url,max=500" example:"https://example.com"`
	LinkAvatar      string              `json:"link_avatar" binding:"omitempty,url,max=500" example:"https://example.com/avatar.jpg"`
	LinkRSS         string              `json:"link_rss" binding:"omitempty,url,max=500" example:"https://example.com/rss.xml"`
	LinkDesc        string              `json:"link_desc" binding:"omitempty,max=500" example:"这是一个示例网站"`
	LinkEmail       string              `json:"link_email" binding:"required,email,max=100" example:"admin@example.com"`
	LinkGroupID     NullableSnowflakeID `json:"link_group_id" binding:"omitempty" example:"1"`
	LinkColorID     NullableSnowflakeID `json:"link_color_id" binding:"omitempty" example:"1"`
	LinkApplyRemark string              `json:"link_apply_remark" binding:"omitempty,max=500" example:"申请友链"`
}

// FriendUserUpdateRequest 用户更新自己友情链接请求
//
// 仅允许更新站点基础信息字段，分组/颜色/级别/排序等管理员专属字段不可改，审核状态保持不变。
type FriendUserUpdateRequest struct {
	LinkName        string `json:"link_name" binding:"omitempty,min=1,max=100" example:"示例网站"`
	LinkURL         string `json:"link_url" binding:"omitempty,url,max=500" example:"https://example.com"`
	LinkAvatar      string `json:"link_avatar" binding:"omitempty,url,max=500" example:"https://example.com/avatar.jpg"`
	LinkRSS         string `json:"link_rss" binding:"omitempty,url,max=500" example:"https://example.com/rss.xml"`
	LinkDesc        string `json:"link_desc" binding:"omitempty,max=500" example:"这是一个示例网站"`
	LinkEmail       string `json:"link_email" binding:"omitempty,email,max=100" example:"admin@example.com"`
	LinkApplyRemark string `json:"link_apply_remark" binding:"omitempty,max=500" example:"申请友链"`
}

// FriendUserQueryRequest 用户查询自己友情链接请求
type FriendUserQueryRequest struct {
	Page       int  `form:"page" binding:"omitempty,min=1" example:"1"`
	PageSize   int  `form:"page_size" binding:"omitempty,min=1,max=100" example:"10"`
	LinkStatus *int `form:"link_status" binding:"omitempty,oneof=0 1 2 3 4" example:"1"`
}
