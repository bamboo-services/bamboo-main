package auth

import (
	"context"
	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/os/glog"
	"xiaoMain/api/auth/v1"
	"xiaoMain/internal/service"
)

// UserLogin 用户登录控制器
func (c *ControllerV1) UserLogin(ctx context.Context, req *v1.UserLoginReq) (res *v1.UserLoginRes, err error) {
	glog.Info(ctx, "[CONTROL] 控制层 UserLogin 接口")
	// 检查用户登录是否有效
	if service.AuthLogic().IsUserLogin(ctx) {

	} else {

	}
	return nil, gerror.NewCode(gcode.CodeNotImplemented)
}
