// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package link

import (
	"context"

	"xiaoMain/api/link/v1"
)

type ILinkV1 interface {
	LinkAddColor(ctx context.Context, req *v1.LinkAddColorReq) (res *v1.LinkAddColorRes, err error)
	LinkAdd(ctx context.Context, req *v1.LinkAddReq) (res *v1.LinkAddRes, err error)
	LinkPlateAdd(ctx context.Context, req *v1.LinkPlateAddReq) (res *v1.LinkPlateAddRes, err error)
	LinkGetColor(ctx context.Context, req *v1.LinkGetColorReq) (res *v1.LinkGetColorRes, err error)
	LinkGetPlate(ctx context.Context, req *v1.LinkGetPlateReq) (res *v1.LinkGetPlateRes, err error)
}
