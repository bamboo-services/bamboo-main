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
	v1 "xiaoMain/api/auth/v1"
)

type (
	IAuth interface {
		// IsUserLogin
		//
		// # 用户是否已登录
		//
		// 用于检查用户是否登录，如果登录则返回 nil, 否则返回错误信息。
		//
		// # 参数:
		//   - ctx: 上下文对象，用于传递和控制请求的生命周期。
		//
		// # 返回:
		//   - err: 如果验证过程中发生错误，返回错误信息。否则返回 nil.
		IsUserLogin(ctx context.Context) (err error)
		// UserLogin
		//
		// # 进行用户登录检查
		//
		// 用于检查用户的登录信息是否正确，如果正确则返回 nil, 否则返回错误信息。
		//
		// # 参数:
		//   - ctx: 上下文对象，用于传递和控制请求的生命周期。
		//   - getData: 用户登录信息(v1.AuthLoginReq)
		//
		// # 返回:
		//   - userUUID: 如果用户登录成功，返回用户的 UUID。否则返回 nil.
		//   - isCorrect: 如果用户登录成功，返回 true。否则返回 false.
		UserLogin(ctx context.Context, getData *v1.AuthLoginReq) (userUUID *string, isCorrect bool, err error)
		// RegisteredUserLogin
		//
		// # 注册用户登录
		//
		// 用于注册用户登录，如果注册成功则返回用户的 Token 信息，否则返回错误信息。
		//
		// # 参数:
		//   - ctx: 上下文对象，用于传递和控制请求的生命周期。
		//   - userUUID: 用户的 UUID。
		//   - remember: 是否记住密码。
		//
		// # 返回:
		//   - userToken: 如果注册成功，返回用户的 Token 信息。否则返回 nil.
		//   - err: 如果注册过程中发生错误，返回错误信息。否则返回 nil.
		RegisteredUserLogin(ctx context.Context, userUUID string, remember bool) (userToken string, err error)
		// ChangeUserPassword 用于修改用户密码。如果密码修改成功，将会清理用户的登录状态，需要用户重新进行登录。
		// 如果用户的密码修改失败，将会返回错误信息。如果修改成功，将返回 nil。
		//
		// 参数:
		// ctx: 上下文对象，用于传递和控制请求的生命周期。
		// password: 用户新的密码字符串。
		//
		// 返回值:
		// err: 如果密码修改成功，返回 nil。否则返回错误信息。
		ChangeUserPassword(ctx context.Context, password string) (err error)
	}
)

var (
	localAuth IAuth
)

func Auth() IAuth {
	if localAuth == nil {
		panic("implement not found for interface IAuth, forgot register?")
	}
	return localAuth
}

func RegisterAuth(i IAuth) {
	localAuth = i
}
