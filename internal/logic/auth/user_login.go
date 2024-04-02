package auth

import (
	"context"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/os/glog"
	"github.com/gogf/gf/v2/os/gtime"
	"xiaoMain/internal/dao"
	"xiaoMain/internal/model/entity"
	"xiaoMain/internal/service"
)

type sAuthLogic struct {
}

func init() {
	service.RegisterAuthLogic(New())
}

func New() *sAuthLogic {
	return &sAuthLogic{}
}

// IsUserLogin 检查用户是否已经登录
func (aL *sAuthLogic) IsUserLogin(ctx context.Context) bool {
	glog.Info(ctx, "[LOGIC] 执行 IsUserLogin 服务层")
	// 根据 ctx 获取 Request 信息
	getRequest := ghttp.RequestFromCtx(ctx)
	// 获取用户的 UUID(UID) 以及 认证密钥
	getUserUid := getRequest.Header.Get("X-User-Uid")
	getUserAuthorize := getRequest.Header.Get("Authorization")
	// 对内容进行校验
	if &getUserUid != nil && &getUserAuthorize != nil {
		var getTokenDO entity.XfToken
		err := dao.XfToken.Ctx(ctx).Where("").Limit(1).Scan(&getTokenDO)
		if err != nil {
			glog.Error(ctx, "[LOGIC] 获取数据库出错")
			return false
		}
		// 检查是否过期
		if gtime.Timestamp() < getTokenDO.ExpiredAt.Timestamp() {
			// 验证登录有效
			if getUserAuthorize == getTokenDO.UserToken {
				glog.Infof(ctx, "[LOGIC] 用户UID %s 任然登录状态", getTokenDO.UserUuid)
				return true
			} else {
				glog.Warning(ctx, "[LOGIC] 用户登录已失效")
				return false
			}
		}
	}
	glog.Warning(ctx, "[LOGIC] 用户未登录")
	return false
}
