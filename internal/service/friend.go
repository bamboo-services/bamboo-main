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
	IFriend interface {
		AddFriend(ctx context.Context, friend base.LinkFriendDTO) (err *berror.ErrorCode)
	}
)

var (
	localFriend IFriend
)

func Friend() IFriend {
	if localFriend == nil {
		panic("implement not found for interface IFriend, forgot register?")
	}
	return localFriend
}

func RegisterFriend(i IFriend) {
	localFriend = i
}
