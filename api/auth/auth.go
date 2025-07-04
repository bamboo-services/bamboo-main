// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package auth

import (
	"context"

	"bamboo-main/api/auth/v1"
)

type IAuthV1 interface {
	AuthLogin(ctx context.Context, req *v1.AuthLoginReq) (res *v1.AuthLoginRes, err error)
	AuthPasswordChange(ctx context.Context, req *v1.AuthPasswordChangeReq) (res *v1.AuthPasswordChangeRes, err error)
	AuthPasswordReset(ctx context.Context, req *v1.AuthPasswordResetReq) (res *v1.AuthPasswordResetRes, err error)
}
