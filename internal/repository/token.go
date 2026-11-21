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
	"github.com/bamboo-services/bamboo-base-go/foundation/v2/snowflake"
	"github.com/bamboo-services/bamboo-base-go/plugins/cache/v2"
	"github.com/bamboo-services/bamboo-base-go/plugins/web-gin/v2/httpresult"
	"github.com/bamboo-services/bamboo-main/pkg/constants"
)

// TokenRepo 令牌仓储
//
// 收口邮箱验证、密码重置两类短期令牌的 Redis 读写，使 logic 层不再直接拼接
// Redis 命令与键名。键名构造仍走 constants.RedisX 统一前缀，TTL 由调用方按业务语义传入。
type TokenRepo struct {
	kc   cache.KeyCache[string, snowflake.SnowflakeID]
	code cache.KeyCache[string, string]
	log  *logging.LogNamedLogger
}

// NewTokenRepo 创建 TokenRepo 实例
func NewTokenRepo(m *cache.Manager, logger ...*logging.Logger) *TokenRepo {
	return &TokenRepo{
		kc:   cache.KeyCacheOf[string, snowflake.SnowflakeID](m),
		code: cache.KeyCacheOf[string, string](m),
		log:  resolveRepoLogger("TokenRepo", logger...),
	}
}

// SaveEmailVerifyToken 保存邮箱验证令牌（24h 过期）
func (r *TokenRepo) SaveEmailVerifyToken(ctx context.Context, token string, userID snowflake.SnowflakeID) *apperr.Error {
	key := constants.RedisEmailVerify.Get(token).String()
	if err := r.kc.Set(ctx, key, &userID, cache.WithTTL(24*time.Hour)); err != nil {
		return apperr.New(httpresult.InternalError, "保存邮箱验证Token失败", err)
	}
	return nil
}

// ConsumeEmailVerifyToken 读取并删除邮箱验证令牌；token 不存在时 found=false（非错误）
func (r *TokenRepo) ConsumeEmailVerifyToken(ctx context.Context, token string) (snowflake.SnowflakeID, bool, *apperr.Error) {
	key := constants.RedisEmailVerify.Get(token).String()
	userID, found, err := r.kc.Get(ctx, key)
	if err != nil {
		return 0, false, apperr.New(httpresult.InternalError, "读取邮箱验证Token失败", err)
	}
	if !found {
		return 0, false, nil
	}
	_ = r.kc.Delete(ctx, key)
	return *userID, true, nil
}

// SavePasswordResetToken 保存密码重置令牌（1h 过期）
func (r *TokenRepo) SavePasswordResetToken(ctx context.Context, token string, userID snowflake.SnowflakeID) *apperr.Error {
	key := constants.RedisPasswordReset.Get(token).String()
	if err := r.kc.Set(ctx, key, &userID, cache.WithTTL(time.Hour)); err != nil {
		return apperr.New(httpresult.InternalError, "保存重置Token失败", err)
	}
	return nil
}

// ExistsPasswordResetToken 检查密码重置令牌是否存在
func (r *TokenRepo) ExistsPasswordResetToken(ctx context.Context, token string) (bool, *apperr.Error) {
	key := constants.RedisPasswordReset.Get(token).String()
	exists, err := r.kc.Exists(ctx, key)
	if err != nil {
		return false, apperr.New(httpresult.InternalError, "验证Token失败", err)
	}
	return exists, nil
}

// ConsumePasswordResetToken 读取并删除密码重置令牌；token 不存在时 found=false
func (r *TokenRepo) ConsumePasswordResetToken(ctx context.Context, token string) (snowflake.SnowflakeID, bool, *apperr.Error) {
	key := constants.RedisPasswordReset.Get(token).String()
	userID, found, err := r.kc.Get(ctx, key)
	if err != nil {
		return 0, false, apperr.New(httpresult.InternalError, "读取重置Token失败", err)
	}
	if !found {
		return 0, false, nil
	}
	_ = r.kc.Delete(ctx, key)
	return *userID, true, nil
}

// SaveRegisterCode 保存注册邮箱验证码（10min 过期）
func (r *TokenRepo) SaveRegisterCode(ctx context.Context, email, code string) *apperr.Error {
	key := constants.RedisRegisterCode.Get(email).String()
	if err := r.code.Set(ctx, key, &code, cache.WithTTL(10*time.Minute)); err != nil {
		return apperr.New(httpresult.InternalError, "保存注册验证码失败", err)
	}
	return nil
}

// ConsumeRegisterCode 校验并消费注册验证码；验证码匹配时删除并返回 true，不匹配或不存在返回 false
func (r *TokenRepo) ConsumeRegisterCode(ctx context.Context, email, code string) (bool, *apperr.Error) {
	key := constants.RedisRegisterCode.Get(email).String()
	stored, found, err := r.code.Get(ctx, key)
	if err != nil {
		return false, apperr.New(httpresult.InternalError, "读取注册验证码失败", err)
	}
	if !found || stored == nil || *stored != code {
		return false, nil
	}
	_ = r.code.Delete(ctx, key)
	return true, nil
}

// SetRegisterCodeLimit 设置注册验证码发送频率限制（60s 过期）
func (r *TokenRepo) SetRegisterCodeLimit(ctx context.Context, email string) *apperr.Error {
	key := constants.RedisRegisterCodeLimit.Get(email).String()
	marker := "1"
	if err := r.code.Set(ctx, key, &marker, cache.WithTTL(time.Minute)); err != nil {
		return apperr.New(httpresult.InternalError, "设置发送频率限制失败", err)
	}
	return nil
}

// ExistsRegisterCodeLimit 检查注册验证码发送频率限制是否存在
func (r *TokenRepo) ExistsRegisterCodeLimit(ctx context.Context, email string) (bool, *apperr.Error) {
	key := constants.RedisRegisterCodeLimit.Get(email).String()
	exists, err := r.code.Exists(ctx, key)
	if err != nil {
		return false, apperr.New(httpresult.InternalError, "检查发送频率限制失败", err)
	}
	return exists, nil
}
