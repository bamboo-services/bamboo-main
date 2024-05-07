// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package link

import (
	"context"

	"xiaoMain/api/link/v1"
)

type ILinkV1 interface {
	LinkColorAdd(ctx context.Context, req *v1.LinkColorAddReq) (res *v1.LinkColorAddRes, err error)
	LinkAdd(ctx context.Context, req *v1.LinkAddReq) (res *v1.LinkAddRes, err error)
	LinkAddAdmin(ctx context.Context, req *v1.LinkAddAdminReq) (res *v1.LinkAddAdminRes, err error)
	LinkLocationAdd(ctx context.Context, req *v1.LinkLocationAddReq) (res *v1.LinkLocationAddRes, err error)
	CheckLinkURLHasConnect(ctx context.Context, req *v1.CheckLinkURLHasConnectReq) (res *v1.CheckLinkURLHasConnectRes, err error)
	CheckLinkIDHasConnect(ctx context.Context, req *v1.CheckLinkIDHasConnectReq) (res *v1.CheckLinkIDHasConnectRes, err error)
	CheckLogoURLHasConnect(ctx context.Context, req *v1.CheckLogoURLHasConnectReq) (res *v1.CheckLogoURLHasConnectRes, err error)
	CheckRssURLHasConnect(ctx context.Context, req *v1.CheckRssURLHasConnectReq) (res *v1.CheckRssURLHasConnectRes, err error)
	LinkGetColor(ctx context.Context, req *v1.LinkGetColorReq) (res *v1.LinkGetColorRes, err error)
	LinkGetColorFull(ctx context.Context, req *v1.LinkGetColorFullReq) (res *v1.LinkGetColorFullRes, err error)
	LinkGetLocation(ctx context.Context, req *v1.LinkGetLocationReq) (res *v1.LinkGetLocationRes, err error)
	LinkGetLocationFull(ctx context.Context, req *v1.LinkGetLocationFullReq) (res *v1.LinkGetLocationFullRes, err error)
}
