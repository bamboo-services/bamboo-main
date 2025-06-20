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
		// GetUserToken
		//
		// 获取用户令牌；
		// 如果获取成功，则返回用户令牌；
		// 如果获取失败，则返回错误码；
		// 如果获取的用户令牌不存在，则返回未找到错误。
		GetUserToken(ctx context.Context, token string) (*base.UserTokenDTO, *berror.ErrorCode)
		// VerifyAndRefreshUserToken
		//
		// 验证并刷新用户令牌；
		// 如果验证成功，则返回用户令牌；
		// 如果验证失败，则返回错误码；
		// 如果验证的用户令牌不存在，则返回未找到错误；
		// 如果刷新令牌不匹配或已过期，则返回未授权错误。
		// 令牌为严格检查，必须提供 UserAgent，并且与当前请求的 UserAgent 匹配。
		VerifyAndRefreshUserToken(ctx context.Context, token string, refreshToken *string) (*base.UserTokenDTO, *berror.ErrorCode)
		// RemoveUserAllToken
		//
		// 删除所有用户令牌；
		// 如果删除成功，则返回 nil；
		// 如果删除失败，则返回错误码。
		RemoveUserAllToken(ctx context.Context) *berror.ErrorCode
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
