// --------------------------------------------------------------------------------
// Copyright (c) 2016-NOW(至今) 筱锋
// Author: 筱锋「xiao_lfeng」(https://www.x-lf.com)
// --------------------------------------------------------------------------------
// 许可证声明：版权所有 (c) 2016-2026 筱锋。保留所有权利。
// 有关MIT许可证的更多信息，请查看项目根目录下的LICENSE文件或访问：
// https://opensource.org/licenses/MIT
// --------------------------------------------------------------------------------

package repository

import (
	"context"
	"time"

	xError "github.com/bamboo-services/bamboo-base-go/common/error"
	xLog "github.com/bamboo-services/bamboo-base-go/common/log"
	xCache "github.com/bamboo-services/bamboo-base-go/major/cache"
	rediscache "github.com/bamboo-services/bamboo-main/internal/models/cache/redis"
	"github.com/bamboo-services/bamboo-main/pkg/constants"
)

// SessionRepo 会话仓储
//
// 收口登录会话（TokenSession）的序列化与 Redis 读写，使 logic 层不再直接拼接
// Redis 命令与键名。会话键名经 constants.RedisAuthToken 统一前缀构造。
type SessionRepo struct {
	kc  xCache.KeyCache[string, rediscache.TokenSession]
	log *xLog.LogNamedLogger
}

// NewSessionRepo 创建 SessionRepo 实例
func NewSessionRepo(m *xCache.Manager) *SessionRepo {
	return &SessionRepo{
		kc:  xCache.KeyCacheOf[string, rediscache.TokenSession](m),
		log: xLog.WithName(xLog.NamedREPO, "SessionRepo"),
	}
}

// SaveSession 按 TTL 将会话数据写入缓存
func (r *SessionRepo) SaveSession(ctx context.Context, token string, session *rediscache.TokenSession, ttl time.Duration) *xError.Error {
	key := constants.RedisAuthToken.Get(token).String()
	if err := r.kc.Set(ctx, key, session, xCache.WithTTL(ttl)); err != nil {
		return xError.NewError(ctx, xError.ServerInternalError, "保存用户会话失败", false, err)
	}
	return nil
}

// GetSession 读取指定令牌对应的会话；会话不存在时 found=false（非错误）
func (r *SessionRepo) GetSession(ctx context.Context, token string) (*rediscache.TokenSession, bool, *xError.Error) {
	key := constants.RedisAuthToken.Get(token).String()
	session, found, err := r.kc.Get(ctx, key)
	if err != nil {
		return nil, false, xError.NewError(ctx, xError.ServerInternalError, "读取用户会话失败", false, err)
	}
	if !found {
		return nil, false, nil
	}
	return session, true, nil
}

// DeleteSession 删除指定令牌对应的会话
func (r *SessionRepo) DeleteSession(ctx context.Context, token string) *xError.Error {
	key := constants.RedisAuthToken.Get(token).String()
	if err := r.kc.Delete(ctx, key); err != nil {
		return xError.NewError(ctx, xError.ServerInternalError, "删除用户会话失败", false, err)
	}
	return nil
}
