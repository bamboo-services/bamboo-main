// ================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// You can delete these comments if you wish manually maintain this interface file.
// ================================================================================

package service

import (
	"context"
)

type (
	IAuthLogic interface {
		// IsUserLogin 检查用户是否已经登录
		IsUserLogin(ctx context.Context) bool
	}
)

var (
	localAuthLogic IAuthLogic
)

func AuthLogic() IAuthLogic {
	if localAuthLogic == nil {
		panic("implement not found for interface IAuthLogic, forgot register?")
	}
	return localAuthLogic
}

func RegisterAuthLogic(i IAuthLogic) {
	localAuthLogic = i
}
