// ================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// You can delete these comments if you wish manually maintain this interface file.
// ================================================================================

package service

import (
	"bamboo-main/internal/model/dto/base"
	"context"
)

type (
	IUser interface {
		// GetUserSimple
		//
		// 获取用户的简单信息；
		// 如果获取成功，则返回用户的简单信息；
		// 如果获取失败，则会产生恐慌（一般情况下都可以正常输出）
		GetUserSimple(ctx context.Context) *base.UserSimpleDTO
		// GetUserDetail
		//
		// 获取用户的详细信息；
		// 如果获取成功，则返回用户的详细信息；
		// 如果获取失败，则会产生恐慌（一般情况下都可以正常输出）
		GetUserDetail(ctx context.Context) *base.UserDetailDTO
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
