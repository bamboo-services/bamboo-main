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

package utility

import (
	"bamboo-main/internal/model/dto"
	"github.com/gogf/gf/v2/util/gconv"
)

// MakePageToTarget
// 创建一个分页对象，包含当前页数据、页码、每页数据量、总数据量和总页数。
// 该函数通过 record 切片长度自动计算总数据量。
//
// 参数:
// - record: 完整的原始数据切片。
// - current: 当前请求的页码。
// - size: 每页显示的数据量。
// - _targetElement: 泛型参数，仅用于类型推断。
//
// 请注意：在使用该函数时，确保 T 和 E 类型之间的兼容性，否则可能产生恐慌
func MakePageToTarget[T, E interface{}](record []*T, current, size int, _targetElement E) *dto.Page[E] {
	// 自动从 record 切片计算总数据量
	totalCount := 0
	if record != nil {
		totalCount = len(record)
	}

	// 处理页码和每页大小的默认值及边界情况
	if current <= 0 {
		current = 1
	}
	if size <= 0 {
		size = 10
	}

	// 计算总页数
	totalPages := (totalCount + size - 1) / size
	if totalCount == 0 {
		totalPages = 0
	}

	// 处理极端情况：record 为 nil 或总数据量为 0
	if record == nil || totalCount == 0 {
		return &dto.Page[E]{
			Record:  make([]*E, 0),
			Current: current,
			Size:    size,
			Total:   0,
			Pages:   totalPages,
		}
	}

	// 确保请求的当前页码在有效范围内
	if current > totalPages && totalPages > 0 {
		current = totalPages
	}
	if totalPages == 0 && current != 1 {
		current = 1
	}

	// 计算当前页数据的起始和结束索引
	startIndex := (current - 1) * size
	endIndex := startIndex + size

	// 安全地从原始 record 切片中获取当前页的数据
	if startIndex >= len(record) { // 使用 len(record) 而不是 totalCount 来判断切片边界
		return &dto.Page[E]{
			Record:  make([]*E, 0),
			Current: current,
			Size:    size,
			Total:   totalCount, // 这里 Total 仍然是实际总数
			Pages:   totalPages,
		}
	}
	if endIndex > len(record) { // 使用 len(record) 而不是 totalCount 来判断切片边界
		endIndex = len(record)
	}

	cutRecord := record[startIndex:endIndex]

	// 将获取到的当前页数据切片转换为目标类型 E
	var targetRecord []*E
	operateErr := gconv.Struct(cutRecord, &targetRecord)
	if operateErr != nil {
		panic("数据转换失败，请检查 T 和 E 类型的兼容性或数据结构: " + operateErr.Error())
	}

	// 构建并返回完整的分页数据结构对象
	return &dto.Page[E]{
		Record:  targetRecord,
		Current: current,
		Size:    size,
		Total:   totalCount,
		Pages:   totalPages,
	}
}

// MakePage
// 创建一个分页对象，包含当前页数据、页码、每页数据量、总数据量和总页数。
// 该函数直接对原始数据切片进行分页，不进行额外的类型转换，并通过 record 切片长度自动计算总数据量。
//
// 参数:
// - record: 完整的原始数据切片。
// - current: 当前请求的页码。
// - size: 每页显示的数据量。
func MakePage[T interface{}](record []*T, current, size int) *dto.Page[T] {
	// 自动从 record 切片计算总数据量
	totalCount := 0
	if record != nil {
		totalCount = len(record)
	}

	// 处理页码和每页大小的默认值及边界情况
	if current <= 0 {
		current = 1
	}
	if size <= 0 {
		size = 10
	}

	// 计算总页数
	totalPages := (totalCount + size - 1) / size
	if totalCount == 0 {
		totalPages = 0
	}

	// 处理极端情况：record 为 nil 或总数据量为 0
	if record == nil || totalCount == 0 {
		return &dto.Page[T]{
			Record:  make([]*T, 0),
			Current: current,
			Size:    size,
			Total:   0,
			Pages:   totalPages,
		}
	}

	// 确保请求的当前页码在有效范围内
	if current > totalPages && totalPages > 0 {
		current = totalPages
	}
	if totalPages == 0 && current != 1 {
		current = 1
	}

	// 计算当前页数据的起始和结束索引
	startIndex := (current - 1) * size
	endIndex := startIndex + size

	// 安全地从原始 record 切片中获取当前页的数据
	if startIndex >= len(record) {
		return &dto.Page[T]{
			Record:  make([]*T, 0),
			Current: current,
			Size:    size,
			Total:   totalCount,
			Pages:   totalPages,
		}
	}
	if endIndex > len(record) {
		endIndex = len(record)
	}

	cutRecord := record[startIndex:endIndex]

	// 直接返回分页数据结构对象，不进行类型转换
	return &dto.Page[T]{
		Record:  cutRecord,
		Current: current,
		Size:    size,
		Total:   totalCount,
		Pages:   totalPages,
	}
}
