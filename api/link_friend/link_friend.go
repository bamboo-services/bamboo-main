// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package link_friend

import (
	"context"

	"bamboo-main/api/link_friend/v1"
)

type ILinkFriendV1 interface {
	LinkFriendAdd(ctx context.Context, req *v1.LinkFriendAddReq) (res *v1.LinkFriendAddRes, err error)
	LinkFriendEdit(ctx context.Context, req *v1.LinkFriendEditReq) (res *v1.LinkFriendEditRes, err error)
	LinkFriendFail(ctx context.Context, req *v1.LinkFriendFailReq) (res *v1.LinkFriendFailRes, err error)
	LinkFriendStatus(ctx context.Context, req *v1.LinkFriendStatusReq) (res *v1.LinkFriendStatusRes, err error)
}
