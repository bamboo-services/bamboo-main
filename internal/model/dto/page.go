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

// Page
//
// 分页数据结构，包含当前页数据、页码、每页数据量、总数据量和总页数。
// T 是一个泛型类型参数，表示当前页数据的类型。
type Page[T interface{}] struct {
	Record  []*T `json:"record" description:"当前页数据" v:"required"`
	Current int  `json:"current" description:"当前页码" v:"required"`
	Size    int  `json:"size" description:"每页数据量" v:"required|between:1,100"`
	Total   int  `json:"total" description:"总数据量" v:"required"`
	Pages   int  `json:"pages" description:"总页数" v:"required"`
}
