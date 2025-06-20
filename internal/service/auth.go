// ================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// You can delete these comments if you wish manually maintain this interface file.
// ================================================================================

package service

import (
	"context"

	"github.com/XiaoLFeng/bamboo-utils/berror"
)

type (
	IAuth interface {
		// VerifyPassword
		//
		// 验证用户的登录密码；
		// 如果验证成功，则返回 nil；
		// 如果验证失败，则返回错误码；
		// 如果获取用户密码失败，则返回内部服务器错误。
		VerifyPassword(ctx context.Context, password string) *berror.ErrorCode
		// ChangePassword
		//
		// 修改用户密码；
		// 如果修改成功，则返回 nil；
		// 如果修改失败，则返回错误码；
		// 如果新密码为空，则返回错误码。
		ChangePassword(ctx context.Context, newPassword string) *berror.ErrorCode
		// VerifyUserByUsername
		//
		// 验证用户的登录信息；
		// 如果验证成功，则返回 nil；
		// 如果验证失败，则返回错误码；
		// 如果获取用户信息或密码失败，则返回内部服务器错误。
		VerifyUserByUsername(ctx context.Context, username string, password string) *berror.ErrorCode
		// VerifyUserByEmail
		//
		// 通过邮箱验证用户的登录信息；
		// 如果验证成功，则返回 nil；
		// 如果验证失败，则返回错误码；
		// 如果获取用户信息或密码失败，则返回内部服务器错误。
		VerifyUserByEmail(ctx context.Context, email string, password string) *berror.ErrorCode
		// VerifyUserByPhone
		//
		// 通过手机号验证用户的登录信息；
		// 如果验证成功，则返回 nil；
		// 如果验证失败，则返回错误码；
		// 如果��取用户信息或密码失败，则返回内部服务器错误。
		VerifyUserByPhone(ctx context.Context, phone string, password string) *berror.ErrorCode
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
