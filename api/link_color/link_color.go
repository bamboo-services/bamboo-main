// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package link_color

import (
	"context"

	"bamboo-main/api/link_color/v1"
)

type ILinkColorV1 interface {
	LinkColorAdd(ctx context.Context, req *v1.LinkColorAddReq) (res *v1.LinkColorAddRes, err error)
	LinkColorDelete(ctx context.Context, req *v1.LinkColorDeleteReq) (res *v1.LinkColorDeleteRes, err error)
	LinkColorEdit(ctx context.Context, req *v1.LinkColorEditReq) (res *v1.LinkColorEditRes, err error)
	LinkColorGet(ctx context.Context, req *v1.LinkColorGetReq) (res *v1.LinkColorGetRes, err error)
	LinkColorPage(ctx context.Context, req *v1.LinkColorPageReq) (res *v1.LinkColorPageRes, err error)
}
