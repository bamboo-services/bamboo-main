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

package auth

import (
	"bamboo-main/internal/consts"
	"bamboo-main/internal/dao"
	"context"
	"github.com/XiaoLFeng/bamboo-utils/berror"
	"github.com/XiaoLFeng/bamboo-utils/blog"
)

// VerifyUserByUsername
//
// 验证用户的登录信息；
// 如果验证成功，则返回 nil；
// 如果验证失败，则返回错误码；
// 如果获取用户信息或密码失败，则返回内部服务器错误。
func (s *sAuth) VerifyUserByUsername(ctx context.Context, username, password string) *berror.ErrorCode {
	blog.ServiceInfo(ctx, "VerifyUserByUsername", "验证用户 %s 的登录信息", username)

	// 获取用户信息
	getUsername, errorCode := dao.System.GetSystemValue(ctx, consts.SystemUserUsernameKey)
	if errorCode != nil {
		return errorCode
	}
	if getUsername == "" {
		return berror.ErrorAddData(&berror.ErrInternalServer, "系统错误 VerifyUserByUsername 函数出现意外错误")
	}

	// 验证用户名
	if getUsername != username {
		blog.ServiceNotice(ctx, "VerifyUserByUsername", "用户名错误")
		return berror.ErrorAddData(&berror.ErrUnauthorized, "用户名或密码错误")
	}

	// 验证密码
	if errorCode := s.VerifyPassword(ctx, password); errorCode != nil {
		return errorCode
	}

	blog.ServiceInfo(ctx, "VerifyUserByUsername", "用户 %s 验证成功", username)
	return nil
}

// VerifyUserByEmail
//
// 通过邮箱验证用户的登录信息；
// 如果验证成功，则返回 nil；
// 如果验证失败，则返回错误码；
// 如果获取用户信息或密码失败，则返回内部服务器错误。
func (s *sAuth) VerifyUserByEmail(ctx context.Context, email, password string) *berror.ErrorCode {
	blog.ServiceInfo(ctx, "VerifyUserByEmail", "通过邮箱验证用户 %s 的登录信息", email)

	// 获取用户邮箱
	getUserEmail, errorCode := dao.System.GetSystemValue(ctx, consts.SystemUserEmailKey)
	if errorCode != nil {
		return errorCode
	}
	if getUserEmail == "" {
		return berror.ErrorAddData(&berror.ErrInternalServer, "系统错误 VerifyUserByEmail 函数出现意外错误")
	}

	// 验证邮箱
	if getUserEmail != email {
		blog.ServiceNotice(ctx, "VerifyUserByEmail", "邮箱错误")
		return berror.ErrorAddData(&berror.ErrUnauthorized, "邮箱或密码错误")
	}

	// 验证密码
	if errorCode := s.VerifyPassword(ctx, password); errorCode != nil {
		return errorCode
	}

	blog.ServiceInfo(ctx, "VerifyUserByEmail", "邮箱 %s 验证成功", email)
	return nil
}

// VerifyUserByPhone
//
// 通过手机号验证用户的登录信息；
// 如果验证成功，则返回 nil；
// 如果验证失败，则返回错误码；
// 如果��取用户信息或密码失败，则返回内部服务器错误。
func (s *sAuth) VerifyUserByPhone(ctx context.Context, phone, password string) *berror.ErrorCode {
	blog.ServiceInfo(ctx, "VerifyUserByPhone", "通过手机号验证用户 %s 的登录信息", phone)

	// 获取用户手机号
	getUserPhone, errorCode := dao.System.GetSystemValue(ctx, consts.SystemUserPhoneKey)
	if errorCode != nil {
		return errorCode
	}
	if getUserPhone == "" {
		return berror.ErrorAddData(&berror.ErrInternalServer, "系统错误 VerifyUserByPhone 函数出现意外错误")
	}

	// 验证手机号
	if getUserPhone != phone {
		blog.ServiceNotice(ctx, "VerifyUserByPhone", "手机号错误")
		return berror.ErrorAddData(&berror.ErrUnauthorized, "手机号或密码错误")
	}

	// 验证密码
	if errorCode := s.VerifyPassword(ctx, password); errorCode != nil {
		return errorCode
	}

	blog.ServiceInfo(ctx, "VerifyUserByPhone", "手机号 %s 验证成功", phone)
	return nil
}
