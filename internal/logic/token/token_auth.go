/*
 * --------------------------------------------------------------------------------
 * Copyright (c) 2016-NOW(至今) 筱锋
 * Author: 筱锋「xiao_lfeng」(https://www.x-lf.com)
 * --------------------------------------------------------------------------------
 * 许可证声明：版权所有 (c) 2016-2025 筱锋。保留所有权利。
 * 有关MIT许可证的更多信息，请查看项目根目录下的LICENSE文件或访问：
 * https://opensource.org/licenses/MIT
 * --------------------------------------------------------------------------------
 */

package token

import (
	"bamboo-main/internal/consts"
	"bamboo-main/internal/model/dto/base"
	"context"
	"github.com/XiaoLFeng/bamboo-utils/berror"
	"github.com/XiaoLFeng/bamboo-utils/blog"
	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/util/gconv"
	"github.com/google/uuid"
	"time"
)

// GenerateUserToken
//
// 生成用户令牌；
// 如果生成成功，则返回用户令牌；
// 如果生成失败，则返回错误码。
func (s *sToken) GenerateUserToken(ctx context.Context) (*base.UserTokenDTO, *berror.ErrorCode) {
	blog.ServiceInfo(ctx, "GenerateUserToken", "生成用户令牌")

	request := ghttp.RequestFromCtx(ctx)

	// 制作新的用户令牌
	var newUserToken = &base.UserTokenDTO{
		Token:        uuid.New().String(),
		RefreshToken: uuid.New().String(),
		UserAgent:    request.UserAgent(),
		UserIP:       request.GetClientIp(),
		CreatedAt:    gtime.Now(),
		ExpiresAt:    gtime.Now().Add(12 * time.Hour),
		RefreshAt:    gtime.Now().Add(7 * 24 * time.Hour),
	}
	getTokenUUIDList, cacheErr := g.Redis().HGetAll(ctx, consts.UserTokenRedisKey)
	if cacheErr != nil {
		blog.ServiceError(ctx, "GenerateUserToken", "获取用户令牌缓存失败 %v", cacheErr)
		return nil, berror.ErrorAddData(&berror.ErrInternalServer, "获取用户令牌缓存失败: "+cacheErr.Error())
	}
	if len(getTokenUUIDList.Map()) > 4 {
		// 删除最早的一份用户令牌
		var oldestToken string
		var oldestTime *gtime.Time
		for token, tokenEntityVar := range getTokenUUIDList.Map() {
			var tokenEntity *base.UserTokenDTO
			operateErr := gconv.Struct(tokenEntityVar, &tokenEntity)
			if operateErr != nil {
				blog.ServiceError(ctx, "GenerateUserToken", "解析用户令牌缓存失败 %v", operateErr)
				return nil, berror.ErrorAddData(&berror.ErrInternalServer, "解析用户令牌缓存失败: "+operateErr.Error())
			}
			if oldestTime == nil || tokenEntity.CreatedAt.Before(oldestTime) {
				oldestTime = tokenEntity.CreatedAt
				oldestToken = token
			}
		}
		if oldestToken != "" {
			_, cacheErr := g.Redis().HDel(ctx, consts.UserTokenRedisKey, oldestToken)
			if cacheErr != nil {
				blog.ServiceError(ctx, "GenerateUserToken", "删除最早的用户令牌失败 %v", cacheErr)
				return nil, berror.ErrorAddData(&berror.ErrInternalServer, "删除最早的用户令牌失败: "+cacheErr.Error())
			}
		}
	}
	_, cacheErr = g.Redis().HSetNX(ctx, consts.UserTokenRedisKey, newUserToken.Token, gjson.MustEncodeString(newUserToken))
	if cacheErr != nil {
		blog.ServiceError(ctx, "GenerateUserToken", "生成用户令牌失败 %v", cacheErr)
		return nil, berror.ErrorAddData(&berror.ErrInternalServer, "生成用户令牌失败: "+cacheErr.Error())
	}

	blog.ServiceNotice(ctx, "GenerateUserToken", "生成用户令牌成功 %s", newUserToken.Token)
	return newUserToken, nil
}

// RemoveUserToken
//
// 删除用户令牌；
// 如果删除成功，则返回 nil；
// 如果删除失败，则返回错误码。
func (s *sToken) RemoveUserToken(ctx context.Context, tokenUUID string) *berror.ErrorCode {
	blog.ServiceInfo(ctx, "RemoveUserToken", "删除用户令牌 %s", tokenUUID)

	// 删除用户令牌
	_, cacheErr := g.Redis().HDel(ctx, consts.UserTokenRedisKey, tokenUUID)
	if cacheErr != nil {
		blog.ServiceError(ctx, "RemoveUserToken", "删除用户令牌失败 %v", cacheErr)
		return berror.ErrorAddData(&berror.ErrInternalServer, "删除用户令牌失败: "+cacheErr.Error())
	}

	return nil
}

// GetUserToken
//
// 获取用户令牌；
// 如果获取成功，则返回用户令牌；
// 如果获取失败，则返回错误码；
// 如果获取的用户令牌不存在，则返回未找到错误。
func (s *sToken) GetUserToken(ctx context.Context, token string) (*base.UserTokenDTO, *berror.ErrorCode) {
	blog.ServiceInfo(ctx, "GetUserToken", "获取用户令牌 %s", token)

	// 获取用户令牌
	getTokenEntity, cacheErr := g.Redis().HGet(ctx, consts.UserTokenRedisKey, token)
	if cacheErr != nil {
		blog.ServiceError(ctx, "GetUserToken", "获取用户令牌失败 %v", cacheErr)
		return nil, berror.ErrorAddData(&berror.ErrInternalServer, "获取用户令牌失败: "+cacheErr.Error())
	}
	if getTokenEntity.IsEmpty() {
		blog.ServiceNotice(ctx, "GetUserToken", "用户令牌 %s 不存在", token)
		return nil, berror.ErrorAddData(&berror.ErrUnauthorized, "用户令牌不存在")
	}

	var userToken base.UserTokenDTO
	err := gconv.Struct(getTokenEntity, &userToken)
	if err != nil {
		blog.ServiceError(ctx, "GetUserToken", "解析用户令牌失败 %v", err)
		return nil, berror.ErrorAddData(&berror.ErrInternalServer, "解析用户令牌失败: "+err.Error())
	}

	return &userToken, nil
}

// VerifyAndRefreshUserToken
//
// 验证并刷新用户令牌；
// 如果验证成功，则返回用户令牌；
// 如果验证失败，则返回错误码；
// 如果验证的用户令牌不存在，则返回未找到错误；
// 如果刷新令牌不匹配或已过期，则返回未授权错误。
// 令牌为严格检查，必须提供 UserAgent，并且与当前请求的 UserAgent 匹配。
func (s *sToken) VerifyAndRefreshUserToken(ctx context.Context, token string, refreshToken *string) (*base.UserTokenDTO, *berror.ErrorCode) {
	blog.ServiceInfo(ctx, "VerifyAndRefreshUserToken", "验证并刷新用户令牌 %s", token)

	// 获取用户令牌
	userToken, errorCode := s.GetUserToken(ctx, token)
	if errorCode != nil {
		return nil, errorCode
	}

	// 验证用户令牌
	if userToken == nil {
		blog.ServiceNotice(ctx, "VerifyAndRefreshUserToken", "用户令牌 %s 不存在", token)
		return nil, berror.ErrorAddData(&berror.ErrUnauthorized, "用户令牌不存在")
	}

	// 验证用户令牌的 UserAgent
	request := ghttp.RequestFromCtx(ctx)
	if userToken.UserAgent != request.UserAgent() {
		blog.ServiceNotice(ctx, "VerifyAndRefreshUserToken", "令牌 %s 非法使用，当前用户端 UserAgent: %s | IP: %s", token, request.UserAgent(), request.GetClientIp())
		return nil, berror.ErrorAddData(&berror.ErrUnauthorized, "令牌非法使用")
	}

	// 检查是否存在刷新令牌
	if refreshToken != nil && *refreshToken != "" {
		// 检查刷新令牌是否匹配
		if userToken.RefreshToken == *refreshToken {
			// 检查刷新令牌是否过期
			if userToken.RefreshAt.Before(gtime.Now()) {
				blog.ServiceNotice(ctx, "VerifyAndRefreshUserToken", "刷新令牌 %s 已经过期", *refreshToken)
				_, cacheErr := g.Redis().HDel(ctx, consts.UserTokenRedisKey, token)
				if cacheErr != nil {
					blog.ServiceError(ctx, "VerifyAndRefreshUserToken", "删除过期用户令牌失败 %v", cacheErr)
					return nil, berror.ErrorAddData(&berror.ErrInternalServer, "删除过期用户令牌失败: "+cacheErr.Error())
				}
				return nil, berror.ErrorAddData(&berror.ErrUnauthorized, "刷新令牌已过期")
			}
		} else {
			blog.ServiceNotice(ctx, "VerifyAndRefreshUserToken", "刷新令牌不匹配，直接验证用户令牌", token)
			if err := s.checkTokenExpiration(ctx, token, userToken); err != nil {
				return nil, err
			}
		}
	} else {
		blog.ServiceNotice(ctx, "VerifyAndRefreshUserToken", "没有提供刷新令牌，直接验证用户令牌 %s", token)
		if err := s.checkTokenExpiration(ctx, token, userToken); err != nil {
			return nil, err
		}
	}

	// 刷新用户令牌
	if refreshToken != nil && *refreshToken != "" && userToken.RefreshToken == *refreshToken {
		// 更新用户令牌的过期时间
		userToken.RefreshToken = uuid.New().String()
		userToken.ExpiresAt = gtime.Now().Add(12 * time.Hour)
		userToken.RefreshAt = gtime.Now().Add(7 * 24 * time.Hour)
	}

	// 更新用户令牌到缓存
	_, cacheErr := g.Redis().HDel(ctx, consts.UserTokenRedisKey, token)
	if cacheErr != nil {
		blog.ServiceError(ctx, "VerifyAndRefreshUserToken", "删除用户令牌缓存失败 %v", cacheErr)
		return nil, berror.ErrorAddData(&berror.ErrInternalServer, "删除用户令牌缓存失败: "+cacheErr.Error())
	}
	_, cacheErr = g.Redis().HSetNX(ctx, consts.UserTokenRedisKey, userToken.Token, gjson.MustEncodeString(userToken))
	if cacheErr != nil {
		blog.ServiceError(ctx, "VerifyAndRefreshUserToken", "更新用户令牌缓存失败 %v", cacheErr)
		return nil, berror.ErrorAddData(&berror.ErrInternalServer, "更新用户令牌缓存失败: "+cacheErr.Error())
	}

	blog.ServiceDebug(ctx, "VerifyAndRefreshUserToken", "用户令牌 %s 验证并刷新成功", token)
	return userToken, nil
}

// RemoveUserAllToken
//
// 删除所有用户令牌；
// 如果删除成功，则返回 nil；
// 如果删除失败，则返回错误码。
func (s *sToken) RemoveUserAllToken(ctx context.Context) *berror.ErrorCode {
	blog.ServiceInfo(ctx, "RemoveUserAllToken", "删除所有用户令牌")

	// 删除所有用户令牌
	_, cacheErr := g.Redis().Del(ctx, consts.UserTokenRedisKey)
	if cacheErr != nil {
		blog.ServiceError(ctx, "RemoveUserAllToken", "删除所有用户令牌失败 %v", cacheErr)
		return berror.ErrorAddData(&berror.ErrInternalServer, "删除所有用户令牌失败: "+cacheErr.Error())
	}

	blog.ServiceNotice(ctx, "RemoveUserAllToken", "所有用户令牌已被删除")
	return nil
}

// checkTokenExpiration
//
// 检查令牌是否过期；
// 如果令牌过期，则删除令牌并返回错误；
// 如果令牌未过期，则返回 nil
func (s *sToken) checkTokenExpiration(ctx context.Context, token string, userToken *base.UserTokenDTO) *berror.ErrorCode {
	if userToken.ExpiresAt.Before(gtime.Now()) {
		blog.ServiceNotice(ctx, "checkTokenExpiration", "用户令牌 %s 已经过期", token)
		_, cacheErr := g.Redis().HDel(ctx, consts.UserTokenRedisKey, token)
		if cacheErr != nil {
			blog.ServiceError(ctx, "checkTokenExpiration", "删除过期用户令牌失败 %v", cacheErr)
			return berror.ErrorAddData(&berror.ErrInternalServer, "删除过期用户令牌失败: "+cacheErr.Error())
		}
		return berror.ErrorAddData(&berror.ErrUnauthorized, "用户令牌已过期")
	}
	return nil
}
