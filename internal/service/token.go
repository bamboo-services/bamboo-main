// ================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// You can delete these comments if you wish manually maintain this interface file.
// ================================================================================

package service

import (
	"bamboo-main/internal/model/dto/base"
	"context"

	"github.com/XiaoLFeng/bamboo-utils/berror"
)

type (
	IToken interface {
		// GenerateUserToken
		//
		// 生成用户令牌；
		// 如果生成成功，则返回用户令牌；
		// 如果生成失败，则返回错误码。
		GenerateUserToken(ctx context.Context) (*base.UserTokenDTO, *berror.ErrorCode)
		// RemoveUserToken
		//
		// 删除用户令牌；
		// 如果删除成功，则返回 nil；
		// 如果删除失败，则返回错误码。
		RemoveUserToken(ctx context.Context, tokenUUID string) *berror.ErrorCode
	}
)

var (
	localToken IToken
)

func Token() IToken {
	if localToken == nil {
		panic("implement not found for interface IToken, forgot register?")
	}
	return localToken
}

func RegisterToken(i IToken) {
	localToken = i
}
