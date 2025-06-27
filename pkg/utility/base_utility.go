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
	"bamboo-main/internal/model/dto/base"
	"context"
)

// GetCtxVarToUserEntity
//
// 从上下文中获取用户实体信息；
// 如果上下文中没有用户实体信息，则返回一个空的 UserSimpleDTO 实例。
func GetCtxVarToUserEntity(ctx context.Context) *base.UserSimpleDTO {
	var userEntity = ctx.Value("UserEntity").(*base.UserSimpleDTO)
	if userEntity == nil {
		return &base.UserSimpleDTO{}
	}
	return userEntity
}

// GetOrDefaultString
//
// 获取字符串指针的值，如果值为 nil 或空字符串，则返回默认值。
// 参数 value 是一个指向字符串的指针，def 是默认值。
// 如果 value 为 nil 或指向的字符串为空，则返回 def；否则返回指向的字符串值。
// 用于处理可能为 nil 的字符串指针，避免空指针解引用错误。
func GetOrDefaultString(value *string, def string) string {
	if value == nil || *value == "" {
		return def
	}
	return *value
}

// GetOrDefault
//
// 获取泛型值，如果 value 为 nil 或指向的值为空，则返回默认值。
// 参数 value 是一个指向 T 的指针，def 是默认值。
func GetOrDefault[E interface{}](value *E, def E) E {
	if &value == nil {
		return def
	}
	return *value
}

// GetOrDefaultArray
//
// 获取指向泛型数组的指针的值，如果值为 nil 或数组为空，则返回默认值。
func GetOrDefaultArray[E interface{}](value *[]*E, def []*E) []*E {
	if value == nil || len(*value) == 0 {
		return def
	}
	return *value

}
