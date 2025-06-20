package auth

import (
	"bamboo-main/api/auth/v1"
	"bamboo-main/internal/service"
	"bamboo-main/pkg/utility"
	"context"
	"github.com/XiaoLFeng/bamboo-utils/blog"
	"github.com/XiaoLFeng/bamboo-utils/bresult"
)

// AuthPasswordChange
//
// 用户修改密码，需要用户提供原始密码和新密码。
func (c *ControllerV1) AuthPasswordChange(ctx context.Context, req *v1.AuthPasswordChangeReq) (res *v1.AuthPasswordChangeRes, err error) {
	blog.ControllerInfo(ctx, "AuthPasswordChange", "用户 %s 修改密码", utility.GetCtxVarToUserEntity(ctx).Username)

	// 检查原始密码是否正确
	iAuth := service.Auth()
	errorCode := iAuth.VerifyPassword(ctx, req.OriginalPassword)
	if errorCode != nil {
		blog.ControllerError(ctx, "AuthPasswordChange", "用户 %s 修改密码失败，原始密码错误", utility.GetCtxVarToUserEntity(ctx).Username)
		return nil, errorCode
	}

	// 修改密码
	if errorCode := iAuth.ChangePassword(ctx, req.NewPassword); errorCode != nil {
		blog.ControllerError(ctx, "AuthPasswordChange", "用户 %s 修改密码失败: %v", utility.GetCtxVarToUserEntity(ctx).Username, err)
		return nil, errorCode
	}

	// 是否需要刷新登录
	if req.NeedAllDeviceRefresh {
		iToken := service.Token()
		errorCode = iToken.RemoveUserAllToken(ctx)
		if errorCode != nil {
			blog.ControllerError(ctx, "AuthPasswordChange", "用户 %s 修改密码后刷新登录失败: %v", utility.GetCtxVarToUserEntity(ctx).Username, errorCode)
			return nil, errorCode
		}
	}

	return &v1.AuthPasswordChangeRes{
		ResponseDTO: bresult.Success(ctx, "用户密码修改成功"),
	}, nil
}
