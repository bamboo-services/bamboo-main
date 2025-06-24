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

package cerror

import "github.com/XiaoLFeng/bamboo-utils/berror"

var (
	ErrMailSend                   = berror.NewErrorCode(40304, "邮件发送失败", nil)     // ErrMailSend 邮件发送失败错误
	ErrMailNextSendTimeNotReached = berror.NewErrorCode(40305, "下次发送邮件时间未到", nil) // ErrMailNextSendTimeNotReached 下次发送邮件时间未到错误
	ErrMailCodeNotExist           = berror.NewErrorCode(40306, "邮件验证码不存在", nil)   // ErrMailCodeNotExist 邮件验证码不存在错误
	ErrMailCodeExpired            = berror.NewErrorCode(40307, "邮件验证码已过期", nil)   // ErrMailCodeExpired 邮件验证码已过期错误
	ErrMailCodeInvalid            = berror.NewErrorCode(40308, "邮件验证码无效", nil)    // ErrMailCodeInvalid 邮件验证码无效错误
)
