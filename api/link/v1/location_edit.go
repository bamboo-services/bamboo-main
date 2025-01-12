/*
 * --------------------------------------------------------------------------------
 * Copyright (c) 2016-NOW(至今) 筱锋
 * Author: 筱锋(https://www.x-lf.com)
 *
 * 本文件包含 XiaoMain 的源代码，该项目的所有源代码均遵循MIT开源许可证协议。
 * --------------------------------------------------------------------------------
 * 许可证声明：
 *
 * 版权所有 (c) 2016-2024 筱锋。保留所有权利。
 *
 * 本软件是“按原样”提供的，没有任何形式的明示或暗示的保证，包括但不限于
 * 对适销性、特定用途的适用性和非侵权性的暗示保证。在任何情况下，
 * 作者或版权持有人均不承担因软件或软件的使用或其他交易而产生的、
 * 由此引起的或以任何方式与此软件有关的任何索赔、损害或其他责任。
 *
 * 使用本软件即表示您了解此声明并同意其条款。
 *
 * 有关MIT许可证的更多信息，请查看项目根目录下的LICENSE文件或访问：
 * https://opensource.org/licenses/MIT
 * --------------------------------------------------------------------------------
 * 免责声明：
 *
 * 使用本软件的风险由用户自担。作者或版权持有人在法律允许的最大范围内，
 * 对因使用本软件内容而导致的任何直接或间接的损失不承担任何责任。
 * --------------------------------------------------------------------------------
 */

package v1

import (
	"github.com/gogf/gf/v2/frame/g"
)

// EditLocationReq
//
// # 编辑位置
//
// 编辑位置(板块)信息。
//
// # 参数
//   - ID: 位置ID
//   - Sort: 排序
//   - Name: 名称
//   - DisplayName: 显示名称
//   - Description: 描述
//   - Reveal: 是否显示
type EditLocationReq struct {
	g.Meta      `path:"/link/location" method:"Put" tags:"链接控制器" summary:"编辑位置(板块)信息"`
	ID          int64  `json:"id" v:"required|regex:^[0-9]+$#请输入ID|ID为数字" dc:"ID"`
	Sort        int64  `json:"sort" v:"required|regex:^[0-9]+$#请输入排序|排序为数字" dc:"排序"`
	Name        string `json:"name" v:"required|length:1,40|regex:^[0-9A-Za-z]+$#请输入名称|名称长度为:1-40位|名称只允许0-9,A-Z,a-z" dc:"名称"` //nolint:lll
	DisplayName string `json:"display_name" v:"required|length:3,40#请输入显示名称|显示名称长度为:3-40位" dc:"显示名称"`
	Description string `json:"description" v:"required#请输入描述" dc:"描述"`
	Reveal      bool   `json:"reveal" v:"required|boolean#请输入是否显示|是否显示为布尔值" dc:"是否显示"`
	UserToken   string `json:"Authorization" in:"header"`
}

// EditLocationRes
//
// # 编辑位置返回
//
// 编辑位置(板块)信息返回。
type EditLocationRes struct {
	g.Meta `mime:"application/json"`
}
