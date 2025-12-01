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

package dto

import "time"

// SiteInfoDTO 站点信息DTO
// 用于返回站点基本信息，包括站点名称、描述和主页介绍
type SiteInfoDTO struct {
	SiteName        string    `json:"site_name"`        // 站点名字
	SiteDescription string    `json:"site_description"` // 站点描述
	Introduction    string    `json:"introduction"`     // 主页介绍
	UpdatedAt       time.Time `json:"updated_at"`       // 最后更新时间
}

// AboutDTO 自我介绍DTO
// 用于返回Markdown格式的自我介绍内容
type AboutDTO struct {
	Content   string    `json:"content"`    // 自我介绍内容（Markdown格式）
	UpdatedAt time.Time `json:"updated_at"` // 最后更新时间
}
