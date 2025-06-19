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
	IUser interface {
		// GetUserSimple
		//
		// 获取用户的简单信息；
		// 如果获取成功，则返回用户的简单信息；
		// 如果获取失败，则返回错误码；
		// 如果获取的用户信息出现错误或获取不到用户信息，则返回内部服务器错误。
		GetUserSimple(ctx context.Context) (*base.UserSimpleDTO, *berror.ErrorCode)
		// GetUserDetail
		//
		// 获取用户的详细信息；
		// 如果获取成功，则返回用户的详细信息；
		// 如果获取失败，则返回错误码；
		// 如果获取的用户信息出现错误或获取不到用户信息，则返回内部服务器错误。
		GetUserDetail(ctx context.Context) (*base.UserDetailDTO, *berror.ErrorCode)
	}
)

var (
	localUser IUser
)

func User() IUser {
	if localUser == nil {
		panic("implement not found for interface IUser, forgot register?")
	}
	return localUser
}

func RegisterUser(i IUser) {
	localUser = i
}
