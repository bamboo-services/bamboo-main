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

// ================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// You can delete these comments if you wish manually maintain this interface file.
// ================================================================================

package service

import (
	"context"
	"xiaoMain/internal/model/dto/flow"
)

type (
	IInfo interface {
		// GetMainInfo
		//
		// # 获取主要信息
		//
		// 用于获取主要信息，如果返回成功则返回具体的信息，若某些情况下无法获取则获取的内容为空
		//
		// # 参数:
		//   - ctx: 上下文对象，用于传递和控制请求的生命周期。
		//
		// # 返回:
		//   - *vo.SiteMainDTO: 如果获取成功，返回具体的信息；否则返回空值。
		GetMainInfo(ctx context.Context) *flow.SiteMainDTO
		// GetBloggerInfo
		//
		// # 获取站长信息
		//
		// 用于获取站长信息，如果返回成功则返回具体的信息，若某些情况下无法获取则获取的内容为空
		//
		// # 参数:
		//   - ctx: 上下文对象，用于传递和控制请求的生命周期。
		//
		// # 返回:
		//   - *vo.SiteBloggerDTO: 如果获取成功，返回具体的信息；否则返回空值。
		GetBloggerInfo(ctx context.Context) *flow.SiteBloggerDTO
		// GetIndexTableData
		//
		// # 获取 Index 数据库中的信息
		//
		// 用于获取 Index 数据库中的信息，如果成功则返回具体的信息，否则返回空值
		//
		// # 参数:
		//   - ctx: 上下文对象，用于传递和控制请求的生命周期。
		//   - name: 需要获取的信息名称(string)
		//
		// # 返回:
		//   - string: 如果获取成功，返回具体的信息；否则返回空值。
		GetIndexTableData(ctx context.Context, name string) string
	}
)

var (
	localInfo IInfo
)

func Info() IInfo {
	if localInfo == nil {
		panic("implement not found for interface IInfo, forgot register?")
	}
	return localInfo
}

func RegisterInfo(i IInfo) {
	localInfo = i
}
