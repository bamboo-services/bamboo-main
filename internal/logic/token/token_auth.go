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
		Token:        uuid.New(),
		RefreshToken: uuid.New(),
		UserAgent:    request.UserAgent(),
		UserIP:       request.GetClientIp(),
		CreatedAt:    gtime.Now(),
		ExpiresAt:    gtime.Now().Add(12 * time.Hour),
		RefreshAt:    gtime.Now().Add(24 * time.Hour),
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
	_, cacheErr = g.Redis().HSetNX(ctx, consts.UserTokenRedisKey, newUserToken.Token.String(), gjson.MustEncodeString(newUserToken))
	if cacheErr != nil {
		blog.ServiceError(ctx, "GenerateUserToken", "生成用户令牌失败 %v", cacheErr)
		return nil, berror.ErrorAddData(&berror.ErrInternalServer, "生成用户令牌失败: "+cacheErr.Error())
	}

	blog.ServiceNotice(ctx, "GenerateUserToken", "生成用户令牌成功 %s", newUserToken.Token.String())
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
