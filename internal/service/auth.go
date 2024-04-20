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
	"xiaoMain/internal/model/entity"
)

type (
	IAuthLogic interface {
		// IsUserLogin 是一个用于检查用户是否已经登录的函数。
		// 它主要用于获取用户的 UUID 和认证密钥，并对这些信息进行校验。如果用户已经完成登录的相关操作，并且此次登录有效，则返回 true 和空字符串。否则，返回 false 和错误信息。
		//
		// 参数:
		// ctx: 上下文对象，用于传递和控制请求的生命周期。
		//
		// 返回值:
		// hasLogin: 如果用户已经登录并且此次登录有效，返回 true。否则返回 false。
		// message: 如果用户未登录或登录已失效，返回错误信息。否则返回空字符串。
		IsUserLogin(ctx context.Context) (hasLogin bool, message string)
		// CheckUserLogin 是一个用于检查用户登录的函数。
		// 它主要用于对用户输入的信息与数据库的内容进行校验，当用户名与用户校验通过后 isCorrect 返回正确值，否则返回错误的内容。
		// 并且当用户正常登录后，将会返回用户的 UUID 作为下一步的登录操作。
		//
		// 参数:
		// ctx: 上下文对象，用于传递和控制请求的生命周期。
		// getData: 用户的登录请求数据，包含了用户的用户名和密码。
		//
		// 返回值:
		// userUUID: 如果用户登录成功，返回用户的 UUID 字符串。
		// isCorrect: 如果用户登录成功，返回 true。否则返回 false。
		// errMessage: 如果用户登录失败，返回错误信息。否则返回空字符串。
		CheckUserLogin(ctx context.Context, getData *v1.AuthLoginReq) (userUUID *string, isCorrect bool, errMessage string)
		// RegisteredUserLogin 用于登记用户的登录信息。当用户完成登录操作后，该方法会将用户的 UUID 存入 token 数据表中，作为用户登录的依据。
		// 在检查用户是否登录时，此数据表的内容作为登录依据。依据 index 数据表字段 key 中的 auth_limit 所对应的 value 的大小作为允许登录节点数的限制。
		//
		// 参数:
		// ctx: 上下文对象，用于传递和控制请求的生命周期。
		// userUUID: 用户的 UUID 字符串。
		// remember: 用户是否选择记住登录状态的布尔值。
		//
		// 返回值:
		// userToken: 用户的 token 信息，包含了用户的 UUID、token、IP、验证信息、User-Agent 和过期时间等信息。
		// err: 如果登录登记成功，返回 nil。否则返回错误信息。
		RegisteredUserLogin(ctx context.Context, userUUID string, remember bool) (userToken *entity.XfToken, err error)
		// CheckUserHasConsoleUser 是一个用于检查用户是否存在于控制台用户列表中的函数。
		// 它主要用于获取用户的用户名，并与数据库中的内容进行比对。如果用户名存在于数据库中，则返回 nil。否则，返回错误信息。
		//
		// 参数:
		// ctx: 上下文对象，用于传递和控制请求的生命周期。
		// username: 需要检查的用户名字符串。
		//
		// 返回值:
		// err: 如果用户名存在于数据库中，返回 nil。否则返回错误信息。
		CheckUserHasConsoleUser(ctx context.Context, username string) (err error)
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
	localAuthLogic IAuthLogic
)

func AuthLogic() IAuthLogic {
	if localAuthLogic == nil {
		panic("implement not found for interface IAuthLogic, forgot register?")
	}
	return localAuthLogic
}

func RegisterAuthLogic(i IAuthLogic) {
	localAuthLogic = i
}
