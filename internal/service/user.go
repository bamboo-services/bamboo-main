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
	IUser interface {
		// GetUserCurrent
		//
		// # 获取当前用户信息
		//
		// 获取当前用户信息，需要用户UUID。根据用户UUID获取当前用户的信息。
		//
		// # 参数
		//   - ctx: 请求的上下文，用于管理超时和取消信号。
		//   - userUUID: 用户UUID
		//
		// # 返回
		//   - userCurrent: 当前用户信息
		//   - err: 在获取过程中发生的任何错误。
		GetUserCurrent(ctx context.Context) (userCurrent *flow.UserCurrentDTO, err error)
		// IsMailHasConsoleUser
		//
		// # 检查邮箱是否为站长邮箱
		//
		// 用于检查该邮箱是否为站长的邮箱，如果是则返回 nil，否则返回错误信息。
		//
		// # 参数:
		//   - ctx: 上下文对象，用于传递和控制请求的生命周期。
		//   - email: 邮箱地址(string)
		//
		// # 返回:
		//   - err: 如果检查过程中发生错误，返回错误信息。否则返回 nil.
		IsMailHasConsoleUser(ctx context.Context, email string) (err error)
		// IsUserHasConsoleUser
		//
		// # 检查用户是否为站长用户
		//
		// 用于检查该用户是否为站长用户，如果是则返回 nil，否则返回错误信息。
		//
		// # 参数:
		//   - ctx: 上下文对象，用于传递和控制请求的生命周期。
		//   - user: 用户名(string)
		//
		// # 返回:
		//   - err: 如果检查过程中发生错误，返回错误信息。否则返回 nil.
		IsUserHasConsoleUser(ctx context.Context, user string) (err error)
	}
)

var (
	localUser IUser
)

func User() IUser {
	if localUser == nil {
		panic("implement not found for interface IUser, forgot register?")
	}
	return localUser
}

func RegisterUser(i IUser) {
	localUser = i
}
