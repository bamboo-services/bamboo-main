// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package info

import (
	"context"

	"xiaoMain/api/info/v1"
)

type IInfoV1 interface {
	InfoUser(ctx context.Context, req *v1.InfoUserReq) (res *v1.InfoUserRes, err error)
	GetWebInfo(ctx context.Context, req *v1.GetWebInfoReq) (res *v1.GetWebInfoRes, err error)
	EditWebInfo(ctx context.Context, req *v1.EditWebInfoReq) (res *v1.EditWebInfoRes, err error)
}
