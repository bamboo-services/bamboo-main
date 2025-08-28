package base

// PaginationInfo 分页信息结构
// 通用的分页信息，适用于所有需要分页的查询
type PaginationInfo struct {
	Page       int   `json:"page"`        // 当前页码
	PageSize   int   `json:"page_size"`   // 每页数量
	Total      int64 `json:"total"`       // 总记录数
	TotalPages int   `json:"total_pages"` // 总页数
	HasNext    bool  `json:"has_next"`    // 是否有下一页
	HasPrev    bool  `json:"has_prev"`    // 是否有上一页
}

// PaginationResponse 分页响应结构
// 泛型结构，可以用于任何类型的分页响应
type PaginationResponse[T any] struct {
	Data       []T            `json:"data"`       // 数据列表
	Pagination PaginationInfo `json:"pagination"` // 分页信息
}

// NewPaginationInfo 创建分页信息
// 根据当前页码、每页大小和总记录数创建分页信息
func NewPaginationInfo(page, pageSize int, total int64) PaginationInfo {
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	totalPages := int((total + int64(pageSize) - 1) / int64(pageSize))
	if totalPages == 0 {
		totalPages = 1
	}

	return PaginationInfo{
		Page:       page,
		PageSize:   pageSize,
		Total:      total,
		TotalPages: totalPages,
		HasNext:    page < totalPages,
		HasPrev:    page > 1,
	}
}

// NewPaginationResponse 创建分页响应
// 泛型函数，用于创建包含数据和分页信息的分页响应
func NewPaginationResponse[T any](data []T, page, pageSize int, total int64) *PaginationResponse[T] {
	if data == nil {
		data = make([]T, 0)
	}

	return &PaginationResponse[T]{
		Data:       data,
		Pagination: NewPaginationInfo(page, pageSize, total),
	}
}

// PaginationQuery 基础分页查询结构
// 提供通用的分页查询参数
type PaginationQuery struct {
	Page      int     `json:"page" form:"page" binding:"omitempty,min=1" default:"1"`                    // 页码，默认为1
	PageSize  int     `json:"page_size" form:"page_size" binding:"omitempty,min=1,max=100" default:"20"` // 每页数量，默认为20，最大100
	SortBy    *string `json:"sort_by" form:"sort_by"`                                                    // 排序字段
	SortOrder *string `json:"sort_order" form:"sort_order" binding:"omitempty" default:"desc"`           // 排序方向，asc或desc，默认为desc
}

// GetPage 获取页码，默认为1
func (q *PaginationQuery) GetPage() int {
	if q.Page <= 0 {
		return 1
	}
	return q.Page
}

// GetPageSize 获取每页数量，默认为10
func (q *PaginationQuery) GetPageSize() int {
	if q.PageSize <= 0 {
		return 20
	}
	if q.PageSize > 100 {
		return 100
	}
	return q.PageSize
}

// GetSortOrder 获取排序方向，默认为desc
func (q *PaginationQuery) GetSortOrder() string {
	if q.SortOrder == nil || *q.SortOrder == "" {
		return "desc"
	}
	return *q.SortOrder
}