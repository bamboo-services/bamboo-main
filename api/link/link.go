// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package link

import (
	"context"

	"bamboo-main/api/link/v1"
)

type ILinkV1 interface {
	LinkFriendAdd(ctx context.Context, req *v1.LinkFriendAddReq) (res *v1.LinkFriendAddRes, err error)
	LinkGroupAdd(ctx context.Context, req *v1.LinkGroupAddReq) (res *v1.LinkGroupAddRes, err error)
	LinkGroupDelete(ctx context.Context, req *v1.LinkGroupDeleteReq) (res *v1.LinkGroupDeleteRes, err error)
	LinkGroupEdit(ctx context.Context, req *v1.LinkGroupEditReq) (res *v1.LinkGroupEditRes, err error)
	LinkGroupList(ctx context.Context, req *v1.LinkGroupListReq) (res *v1.LinkGroupListRes, err error)
}
