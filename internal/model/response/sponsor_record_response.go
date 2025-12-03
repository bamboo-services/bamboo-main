/*
 * --------------------------------------------------------------------------------
 * Copyright (c) 2016-NOW(至今) 筱锋
 * Author: 筱锋「xiao_lfeng」(https://www.x-lf.com)
 * --------------------------------------------------------------------------------
 * 许可证声明:版权所有 (c) 2016-2025 筱锋。保留所有权利。
 * 有关MIT许可证的更多信息,请查看项目根目录下的LICENSE文件或访问:
 * https://opensource.org/licenses/MIT
 * --------------------------------------------------------------------------------
 */

package response

import (
	"bamboo-main/internal/model/base"
	"bamboo-main/internal/model/dto"
)

// SponsorRecordAddResponse 添加记录响应
type SponsorRecordAddResponse struct {
	dto.SponsorRecordDetailDTO
}

// SponsorRecordUpdateResponse 更新记录响应
type SponsorRecordUpdateResponse struct {
	dto.SponsorRecordDetailDTO
}

// SponsorRecordDetailResponse 详情响应
type SponsorRecordDetailResponse struct {
	dto.SponsorRecordDetailDTO
}

// SponsorRecordPageResponse 分页响应（后台）
type SponsorRecordPageResponse struct {
	base.PaginationResponse[dto.SponsorRecordNormalDTO]
}

// SponsorRecordPublicPageResponse 公开分页响应（前台）
type SponsorRecordPublicPageResponse struct {
	base.PaginationResponse[dto.SponsorRecordSimpleDTO]
}

// SponsorRecordDeleteResponse 删除响应
type SponsorRecordDeleteResponse struct {
	Message string `json:"message"`
}
