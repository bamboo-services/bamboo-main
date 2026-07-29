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
	xSnowflake "github.com/bamboo-services/bamboo-base-go/common/snowflake"
	xCache "github.com/bamboo-services/bamboo-base-go/major/cache"
	"github.com/bamboo-services/bamboo-main/pkg/constants"
)

// TokenRepo 令牌仓储
//
// 收口邮箱验证、密码重置两类短期令牌的 Redis 读写，使 logic 层不再直接拼接
// Redis 命令与键名。键名构造仍走 constants.RedisX 统一前缀，TTL 由调用方按业务语义传入。
type TokenRepo struct {
	kc   xCache.KeyCache[string, xSnowflake.SnowflakeID]
	code xCache.KeyCache[string, string]
	log  *xLog.LogNamedLogger
}

// NewTokenRepo 创建 TokenRepo 实例
func NewTokenRepo(m *xCache.Manager) *TokenRepo {
	return &TokenRepo{
		kc:   xCache.KeyCacheOf[string, xSnowflake.SnowflakeID](m),
		code: xCache.KeyCacheOf[string, string](m),
		log:  xLog.WithName(xLog.NamedREPO, "TokenRepo"),
	}
}

// SaveEmailVerifyToken 保存邮箱验证令牌（24h 过期）
func (r *TokenRepo) SaveEmailVerifyToken(ctx context.Context, token string, userID xSnowflake.SnowflakeID) *xError.Error {
	key := constants.RedisEmailVerify.Get(token).String()
	if err := r.kc.Set(ctx, key, &userID, xCache.WithTTL(24*time.Hour)); err != nil {
		return xError.NewError(ctx, xError.ServerInternalError, "保存邮箱验证Token失败", false, err)
	}
	return nil
}

// ConsumeEmailVerifyToken 读取并删除邮箱验证令牌；token 不存在时 found=false（非错误）
func (r *TokenRepo) ConsumeEmailVerifyToken(ctx context.Context, token string) (xSnowflake.SnowflakeID, bool, *xError.Error) {
	key := constants.RedisEmailVerify.Get(token).String()
	userID, found, err := r.kc.Get(ctx, key)
	if err != nil {
		return 0, false, xError.NewError(ctx, xError.ServerInternalError, "读取邮箱验证Token失败", false, err)
	}
	if !found {
		return 0, false, nil
	}
	_ = r.kc.Delete(ctx, key)
	return *userID, true, nil
}

// SavePasswordResetToken 保存密码重置令牌（1h 过期）
func (r *TokenRepo) SavePasswordResetToken(ctx context.Context, token string, userID xSnowflake.SnowflakeID) *xError.Error {
	key := constants.RedisPasswordReset.Get(token).String()
	if err := r.kc.Set(ctx, key, &userID, xCache.WithTTL(time.Hour)); err != nil {
		return xError.NewError(ctx, xError.ServerInternalError, "保存重置Token失败", false, err)
	}
	return nil
}

// ExistsPasswordResetToken 检查密码重置令牌是否存在
func (r *TokenRepo) ExistsPasswordResetToken(ctx context.Context, token string) (bool, *xError.Error) {
	key := constants.RedisPasswordReset.Get(token).String()
	exists, err := r.kc.Exists(ctx, key)
	if err != nil {
		return false, xError.NewError(ctx, xError.ServerInternalError, "验证Token失败", false, err)
	}
	return exists, nil
}

// ConsumePasswordResetToken 读取并删除密码重置令牌；token 不存在时 found=false
func (r *TokenRepo) ConsumePasswordResetToken(ctx context.Context, token string) (xSnowflake.SnowflakeID, bool, *xError.Error) {
	key := constants.RedisPasswordReset.Get(token).String()
	userID, found, err := r.kc.Get(ctx, key)
	if err != nil {
		return 0, false, xError.NewError(ctx, xError.ServerInternalError, "读取重置Token失败", false, err)
	}
	if !found {
		return 0, false, nil
	}
	_ = r.kc.Delete(ctx, key)
	return *userID, true, nil
}

// SaveRegisterCode 保存注册邮箱验证码（10min 过期）
func (r *TokenRepo) SaveRegisterCode(ctx context.Context, email, code string) *xError.Error {
	key := constants.RedisRegisterCode.Get(email).String()
	if err := r.code.Set(ctx, key, &code, xCache.WithTTL(10*time.Minute)); err != nil {
		return xError.NewError(ctx, xError.ServerInternalError, "保存注册验证码失败", false, err)
	}
	return nil
}

// ConsumeRegisterCode 校验并消费注册验证码；验证码匹配时删除并返回 true，不匹配或不存在返回 false
func (r *TokenRepo) ConsumeRegisterCode(ctx context.Context, email, code string) (bool, *xError.Error) {
	key := constants.RedisRegisterCode.Get(email).String()
	stored, found, err := r.code.Get(ctx, key)
	if err != nil {
		return false, xError.NewError(ctx, xError.ServerInternalError, "读取注册验证码失败", false, err)
	}
	if !found || stored == nil || *stored != code {
		return false, nil
	}
	_ = r.code.Delete(ctx, key)
	return true, nil
}

// SetRegisterCodeLimit 设置注册验证码发送频率限制（60s 过期）
func (r *TokenRepo) SetRegisterCodeLimit(ctx context.Context, email string) *xError.Error {
	key := constants.RedisRegisterCodeLimit.Get(email).String()
	marker := "1"
	if err := r.code.Set(ctx, key, &marker, xCache.WithTTL(time.Minute)); err != nil {
		return xError.NewError(ctx, xError.ServerInternalError, "设置发送频率限制失败", false, err)
	}
	return nil
}

// ExistsRegisterCodeLimit 检查注册验证码发送频率限制是否存在
func (r *TokenRepo) ExistsRegisterCodeLimit(ctx context.Context, email string) (bool, *xError.Error) {
	key := constants.RedisRegisterCodeLimit.Get(email).String()
	exists, err := r.code.Exists(ctx, key)
	if err != nil {
		return false, xError.NewError(ctx, xError.ServerInternalError, "检查发送频率限制失败", false, err)
	}
	return exists, nil
}
