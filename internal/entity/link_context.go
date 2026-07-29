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
	xSnowflake "github.com/bamboo-services/bamboo-base-go/common/snowflake"
	xModels "github.com/bamboo-services/bamboo-base-go/major/models"

	bConst "github.com/bamboo-services/bamboo-main/pkg/constants"
)

// LinkFriend 表示一个友情链接实体，用于描述友情链接及其属性。
//
// 该类型包含友情链接的唯一标识符、名称、URL、头像、RSS、描述等信息。
// 同时记录该友情链接的创建时间和更新时间，便于数据管理和审计。
//
// 注意: 友情链接通过外键关联到分组和颜色表。
type LinkFriend struct {
	xModels.BaseEntity                         // 嵌入基础实体（SnowflakeID 主键 + timestamptz 时间戳）
	Name               string                  `json:"name" gorm:"type:varchar(100);not null;comment:友链名称"`                                     // 友链名称
	URL                string                  `json:"url" gorm:"type:varchar(500);not null;comment:友链URL地址"`                                   // 友链URL地址
	Avatar             *string                 `json:"avatar,omitempty" gorm:"type:varchar(500);comment:友链头像URL"`                               // 友链头像URL
	RSS                *string                 `json:"rss,omitempty" gorm:"type:varchar(500);comment:友链RSS地址"`                                  // 友链RSS地址
	Description        *string                 `json:"description,omitempty" gorm:"type:text;comment:友链描述"`                                     // 友链描述
	Email              *string                 `json:"email,omitempty" gorm:"type:varchar(100);comment:友链联系邮箱"`                                 // 友链联系邮箱
	UserID             *xSnowflake.SnowflakeID `json:"user_id,omitempty" gorm:"comment:归属用户ID;index"`                                           // 归属用户ID（按邮箱确认归属，为空表示游客提交尚未关联）
	GroupID            *xSnowflake.SnowflakeID `json:"group_id,omitempty" gorm:"comment:所属分组ID"`                                                // 所属分组ID
	ColorID            *xSnowflake.SnowflakeID `json:"color_id,omitempty" gorm:"comment:友链颜色ID"`                                                // 友链颜色ID
	SortOrder          int                     `json:"sort_order" gorm:"type:int;default:0;comment:友链排序"`                                       // 友链排序
	Status             int                     `json:"status" gorm:"type:int;default:0;comment:友链状态（0: 待审核, 1: 已通过, 2: 已拒绝, 3: 下架待审核, 4: 已下架）"` // 友链状态
	IsFailure          int                     `json:"is_failure" gorm:"type:int;default:0;comment:友链失效标志（0: 正常, 1: 失效）"`                       // 友链失效标志
	Level              int                     `json:"level" gorm:"type:int;default:0;comment:友链级别（0: 一般, 1: 好友, 2: 高级, 3: 广告）"`                // 友链级别
	FailReason         *string                 `json:"fail_reason,omitempty" gorm:"type:text;comment:友链失效原因"`                                   // 友链失效原因
	ApplyRemark        *string                 `json:"apply_remark,omitempty" gorm:"type:text;comment:申请者备注"`                                   // 申请者备注
	ReviewRemark       *string                 `json:"review_remark,omitempty" gorm:"type:text;comment:审核备注"`                                   // 审核备注

	// 关联关系
	GroupFKey *LinkGroup  `json:"group_f_key,omitempty" gorm:"foreignKey:GroupID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;comment:友链分组外键"` // 友链分组外键
	ColorFKey *LinkColor  `json:"color_f_key,omitempty" gorm:"foreignKey:ColorID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;comment:友链颜色外键"` // 友链颜色外键
	UserFKey  *SystemUser `json:"user_f_key,omitempty" gorm:"foreignKey:UserID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;comment:归属用户外键"`   // 归属用户外键
}

// GetGene 返回 xSnowflake.Gene，用于标识该实体在 ID 生成时使用的基因类型。
func (_ *LinkFriend) GetGene() xSnowflake.Gene {
	return bConst.GeneLink
}
