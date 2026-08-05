// --------------------------------------------------------------------------------
// Copyright (c) 2016-NOW(至今) 筱锋
// Author: 筱锋「xiao_lfeng」(https://www.x-lf.com)
// --------------------------------------------------------------------------------
// 许可证声明：版权所有 (c) 2016-2026 筱锋。保留所有权利。
// 有关MIT许可证的更多信息，请查看项目根目录下的LICENSE文件或访问：
// https://opensource.org/licenses/MIT
// --------------------------------------------------------------------------------

package entity

import (
	"time"

	xSnowflake "github.com/bamboo-services/bamboo-base-go/common/snowflake"
	xModels "github.com/bamboo-services/bamboo-base-go/major/models"

	bConst "github.com/bamboo-services/bamboo-main/pkg/constants"
)

// SponsorRecord 表示一条赞助记录实体，用于记录每一笔赞助的详细信息。
//
// 该类型包含赞助者昵称、跳转地址、赞助金额、赞助渠道、留言等信息。
// 支持匿名展示和隐藏功能，便于灵活管理赞助记录的前台展示。
// 同时记录审核状态与归属信息：游客可自助申请展示，经管理员审核后前台可见。
//
// 注意: 金额以"分"为单位存储，避免浮点精度问题。展示时由前端转换为元。
type SponsorRecord struct {
	xModels.BaseEntity                         // 嵌入基础实体（SnowflakeID 主键 + timestamptz 时间戳）
	Nickname           string                  `json:"nickname" gorm:"type:varchar(100);not null;comment:赞助者昵称"`        // 赞助者展示名称
	RedirectURL        *string                 `json:"redirect_url,omitempty" gorm:"type:varchar(500);comment:赞助者跳转地址"` // 点击昵称时的跳转地址
	Amount             int64                   `json:"amount" gorm:"type:bigint;not null;comment:赞助金额（单位：分）"`           // 金额，单位为分（如 666 表示 6.66 元）
	ChannelID          *xSnowflake.SnowflakeID `json:"channel_id,omitempty" gorm:"comment:赞助渠道ID;index"`                // 关联的赞助渠道ID
	Message            *string                 `json:"message,omitempty" gorm:"type:text;comment:赞助留言"`                 // 赞助者的留言内容
	SponsorAt          *time.Time              `json:"sponsor_at,omitempty" gorm:"type:timestamptz;comment:赞助发生时间"`     // 实际赞助发生的时间
	SortOrder          int                     `json:"sort_order" gorm:"type:int;default:0;comment:显示排序"`               // 排序值，用于前台赞助墙排序
	IsAnonymous        bool                    `json:"is_anonymous" gorm:"type:boolean;default:false;comment:是否匿名展示"`     // 为 true 时前台显示"匿名用户"
	IsHidden           bool                    `json:"is_hidden" gorm:"type:boolean;default:false;comment:是否在前台隐藏"`       // 为 true 时前台不展示该记录
	Status             int                     `json:"status" gorm:"type:int;not null;default:1;comment:赞助状态（0: 待审核, 1: 已通过, 2: 已拒绝）"` // 赞助状态（默认已通过，兼容历史手动录入记录）
	Email              *string                 `json:"email,omitempty" gorm:"type:varchar(100);comment:联系邮箱"`              // 联系邮箱（用于确认归属与审核结果通知）
	UserID             *xSnowflake.SnowflakeID `json:"user_id,omitempty" gorm:"comment:归属用户ID;index"`                     // 归属用户ID（按邮箱确认归属，为空表示游客提交尚未关联）
	ApplyRemark        *string                 `json:"apply_remark,omitempty" gorm:"type:text;comment:申请者备注"`              // 申请者备注
	ReviewRemark       *string                 `json:"review_remark,omitempty" gorm:"type:text;comment:审核备注"`              // 审核备注

	// 关联关系
	ChannelFKey *SponsorChannel `json:"channel_f_key,omitempty" gorm:"foreignKey:ChannelID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;comment:赞助渠道外键"` // 赞助渠道外键，关联 SponsorChannel 类型
	UserFKey    *SystemUser     `json:"user_f_key,omitempty" gorm:"foreignKey:UserID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;comment:归属用户外键"`      // 归属用户外键，关联 SystemUser 类型
}

// GetGene 返回 xSnowflake.Gene，用于标识该实体在 ID 生成时使用的基因类型。
func (_ *SponsorRecord) GetGene() xSnowflake.Gene {
	return bConst.GeneSponsorRecord
}
