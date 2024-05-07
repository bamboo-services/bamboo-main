// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package rss

import (
	"context"

	"xiaoMain/api/rss/v1"
)

type IRssV1 interface {
	GetLinkRssInfo(ctx context.Context, req *v1.GetLinkRssInfoReq) (res *v1.GetLinkRssInfoRes, err error)
}
