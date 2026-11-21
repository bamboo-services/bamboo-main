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

	"github.com/bamboo-services/bamboo-base-go/foundation/v2/apperr"
	"github.com/bamboo-services/bamboo-base-go/foundation/v2/logging"
	"github.com/bamboo-services/bamboo-base-go/plugins/cache/v2"
	"github.com/bamboo-services/bamboo-base-go/plugins/web-gin/v2/httpresult"
	rediscache "github.com/bamboo-services/bamboo-main/internal/models/cache/redis"
	"github.com/bamboo-services/bamboo-main/pkg/constants"
)

// SessionRepo 会话仓储
//
// 收口登录会话（TokenSession）的序列化与 Redis 读写，使 logic 层不再直接拼接
// Redis 命令与键名。会话键名经 constants.RedisAuthToken 统一前缀构造。
type SessionRepo struct {
	kc  cache.KeyCache[string, rediscache.TokenSession]
	log *logging.LogNamedLogger
}

// NewSessionRepo 创建 SessionRepo 实例
func NewSessionRepo(m *cache.Manager, logger ...*logging.Logger) *SessionRepo {
	return &SessionRepo{
		kc:  cache.KeyCacheOf[string, rediscache.TokenSession](m),
		log: resolveRepoLogger("SessionRepo", logger...),
	}
}

// SaveSession 按 TTL 将会话数据写入缓存
func (r *SessionRepo) SaveSession(ctx context.Context, token string, session *rediscache.TokenSession, ttl time.Duration) *apperr.Error {
	key := constants.RedisAuthToken.Get(token).String()
	if err := r.kc.Set(ctx, key, session, cache.WithTTL(ttl)); err != nil {
		return apperr.New(httpresult.InternalError, "保存用户会话失败", err)
	}
	return nil
}

// GetSession 读取指定令牌对应的会话；会话不存在时 found=false（非错误）
func (r *SessionRepo) GetSession(ctx context.Context, token string) (*rediscache.TokenSession, bool, *apperr.Error) {
	key := constants.RedisAuthToken.Get(token).String()
	session, found, err := r.kc.Get(ctx, key)
	if err != nil {
		return nil, false, apperr.New(httpresult.InternalError, "读取用户会话失败", err)
	}
	if !found {
		return nil, false, nil
	}
	return session, true, nil
}

// DeleteSession 删除指定令牌对应的会话
func (r *SessionRepo) DeleteSession(ctx context.Context, token string) *apperr.Error {
	key := constants.RedisAuthToken.Get(token).String()
	if err := r.kc.Delete(ctx, key); err != nil {
		return apperr.New(httpresult.InternalError, "删除用户会话失败", err)
	}
	return nil
}
