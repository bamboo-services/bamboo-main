// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package link_group

import (
	"context"

	"bamboo-main/api/link_group/v1"
)

type ILinkGroupV1 interface {
	LinkGroupAdd(ctx context.Context, req *v1.LinkGroupAddReq) (res *v1.LinkGroupAddRes, err error)
	LinkGroupDelete(ctx context.Context, req *v1.LinkGroupDeleteReq) (res *v1.LinkGroupDeleteRes, err error)
	LinkGroupEdit(ctx context.Context, req *v1.LinkGroupEditReq) (res *v1.LinkGroupEditRes, err error)
	LinkGroupGet(ctx context.Context, req *v1.LinkGroupGetReq) (res *v1.LinkGroupGetRes, err error)
	LinkGroupList(ctx context.Context, req *v1.LinkGroupListReq) (res *v1.LinkGroupListRes, err error)
}
