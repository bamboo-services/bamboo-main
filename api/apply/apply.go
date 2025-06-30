// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package apply

import (
	"context"

	"bamboo-main/api/apply/v1"
)

type IApplyV1 interface {
	ApplyLink(ctx context.Context, req *v1.ApplyLinkReq) (res *v1.ApplyLinkRes, err error)
}
