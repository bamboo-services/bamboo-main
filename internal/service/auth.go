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
		// IsUserLogin
		//
		// 检查用户是否已经登录，若用户已经完成登录的相关操作，若用户的此次登录有效则返回 true 否则返回 false
		IsUserLogin(ctx context.Context) (hasLogin bool, message string)
		// CheckUserLogin
		//
		// 对用户的登录进行检查。主要用于对用户输入的信息与数据库的内容进行校验，当用户名与用户校验通过后 isCorrect 返回正确值，否则返回错误的内容
		// 并且当用户正常登录后，将会返回用户的 UUID 作为下一步的登录操作
		CheckUserLogin(ctx context.Context, getData *v1.UserLoginReq) (userUUID *string, isCorrect bool)
		// RegisteredUserLogin
		//
		// 对用户的登录内容进行登记，将用户的 UUID 传入后存入 token 数据表中，作为用户登录的登录依据。在检查用户是否登录时候，此数据表的内容作为登录
		// 依据。
		//
		// 依据 index 数据表字段 key 中的 auth_limit 所对应的 value 的大小作为允许登录节点数的限制
		RegisteredUserLogin(ctx context.Context, userUUID string, remember bool) (userToken *entity.XfToken, err error)
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
