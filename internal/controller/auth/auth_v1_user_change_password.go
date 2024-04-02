package auth

import (
	"context"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"

	"xiaoMain/api/auth/v1"
)

func (c *ControllerV1) UserChangePassword(ctx context.Context, req *v1.UserChangePasswordReq) (res *v1.UserChangePasswordRes, err error) {
	return nil, gerror.NewCode(gcode.CodeNotImplemented)
}
