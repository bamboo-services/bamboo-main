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

package dtoBase

// PaginationInfo 分页信息
type PaginationInfo struct {
	Page       int  `json:"page"`        // 当前页码
	PageSize   int  `json:"page_size"`   // 每页数量
	Total      int  `json:"total"`       // 总记录数
	TotalPages int  `json:"total_pages"` // 总页数
	HasNext    bool `json:"has_next"`    // 是否有下一页
	HasPrev    bool `json:"has_prev"`    // 是否有上一页
}
