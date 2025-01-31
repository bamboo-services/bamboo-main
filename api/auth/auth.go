// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package auth

import (
	"context"

	"xiaoMain/api/auth/v1"
)

type IAuthV1 interface {
	AuthChangePassword(ctx context.Context, req *v1.AuthChangePasswordReq) (res *v1.AuthChangePasswordRes, err error)
	AuthLogin(ctx context.Context, req *v1.AuthLoginReq) (res *v1.AuthLoginRes, err error)
	AuthResetPassword(ctx context.Context, req *v1.AuthResetPasswordReq) (res *v1.AuthResetPasswordRes, err error)
}
