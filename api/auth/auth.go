// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package auth

import (
	"context"

	"develop/api/auth/v1"
)

type IAuthV1 interface {
	UserChangePassword(ctx context.Context, req *v1.UserChangePasswordReq) (res *v1.UserChangePasswordRes, err error)
	UserLogin(ctx context.Context, req *v1.UserLoginReq) (res *v1.UserLoginRes, err error)
	UserResetPassword(ctx context.Context, req *v1.UserResetPasswordReq) (res *v1.UserResetPasswordRes, err error)
}
