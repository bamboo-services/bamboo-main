// ================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// You can delete these comments if you wish manually maintain this interface file.
// ================================================================================

package service

import (
	"context"
	v1 "xiaoMain/api/auth/v1"
	"xiaoMain/internal/model/entity"
)

type (
	IAuthLogic interface {
		// IsUserLogin
		// 检查用户是否已经登录，若用户已经完成登录的相关操作，若用户的此次登录有效则返回 true 否则返回 false
		IsUserLogin(ctx context.Context) bool
		// CheckUserLogin
		// 对用户的登录进行检查。主要用于对用户输入的信息与数据库的内容进行校验，当用户名与用户校验通过后 isCorrect 返回正确值，否则返回错误的内容
		// 并且当用户正常登录后，将会返回用户的 UUID 作为下一步的登录操作
		CheckUserLogin(ctx context.Context, getData *v1.UserLoginReq) (userUuid *string, isCorrect bool)
		// RegisteredUserLogin
		// 对用户的登录内容进行登记，将用户的 UUID 传入后存入 token 数据表中，作为用户登录的登录依据。在检查用户是否登录时候，此数据表的内容作为登录
		// 依据。
		//
		// 依据 index 数据表字段 key 中的 auth_limit 所对应的 value 的大小作为允许登录节点数的限制内容
		RegisteredUserLogin(ctx context.Context, userUuid string, remember bool) (userToken *entity.XfToken, err error)
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
