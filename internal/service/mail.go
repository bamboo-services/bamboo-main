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
)

type (
	IMailUserLogic interface {
		// VerificationCodeHasCorrect
		// 验证验证码是否正确，若验证码正确将会返回 true，否则返回 false；
		// 若返回错误的内容将会返回具体的错误原因，不会抛出 Error
		VerificationCodeHasCorrect(ctx context.Context, email string, code string, scenes string) (isCorrect bool, info string)
		// SendEmailVerificationCode
		// 根据输入的场景进行邮箱的发送，需要保证场景的合法性，场景的合法性参考 consts.Scenes 的参考值
		// 若邮件发送的过程中出现错误将会终止发件并且返回 error 信息，发件成功返回 nil
		SendEmailVerificationCode(ctx context.Context, mail string, scenes string) (err error)
	}
)

var (
	localMailUserLogic IMailUserLogic
)

func MailUserLogic() IMailUserLogic {
	if localMailUserLogic == nil {
		panic("implement not found for interface IMailUserLogic, forgot register?")
	}
	return localMailUserLogic
}

func RegisterMailUserLogic(i IMailUserLogic) {
	localMailUserLogic = i
}
