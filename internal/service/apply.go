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
	IApply interface {
		// NewApplyLink
		//
		// 创建一个新的友情链接申请并存储到数据库中。
		// 如果存储过程中出现错误，返回对应的错误码。
		// 此方法还会生成唯一的UUID用于标识新申请，并清理全局缓存数据。
		NewApplyLink(ctx context.Context, newApplyLink base.LinkFriendDTO) *berror.ErrorCode
	}
)

var (
	localApply IApply
)

func Apply() IApply {
	if localApply == nil {
		panic("implement not found for interface IApply, forgot register?")
	}
	return localApply
}

func RegisterApply(i IApply) {
	localApply = i
}
