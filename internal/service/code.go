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
	ICode interface {
		// GenerateEmailCode
		//
		// 生成邮箱验证码；
		// 如果验证码已存在，则不重新生成；
		// 如果验证码不存在，则随机生成一个6位数的验证码；
		// 如果存储验证码到缓存失败，则返回错误码。
		GenerateEmailCode(ctx context.Context, email string, code *string) *berror.ErrorCode
		// VerifyEmailCode
		//
		// 验证邮箱验证码；
		// 如果验证码验证成功，则删除缓存中的验证码；
		// 如果验证码不存在，则返回错误码 cerror.ErrMailCodeNotExist；
		// 如果验证码不匹配，则返回错误码 cerror.ErrMailCodeInvalid；
		// 如果获取验证码失败，则返回错误码 berror.ErrInternalServer。
		// 如果删除验证码缓存失败，则返回错误码 berror.ErrInternalServer。
		// 如果验证成功，则返回 nil。
		VerifyEmailCode(ctx context.Context, email string, code string) *berror.ErrorCode
	}
)

var (
	localCode ICode
)

func Code() ICode {
	if localCode == nil {
		panic("implement not found for interface ICode, forgot register?")
	}
	return localCode
}

func RegisterCode(i ICode) {
	localCode = i
}
