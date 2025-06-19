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

package user

import (
	"bamboo-main/internal/consts"
	"bamboo-main/internal/dao"
	"bamboo-main/internal/model/dto/base"
	"context"
	"github.com/XiaoLFeng/bamboo-utils/berror"
	"github.com/XiaoLFeng/bamboo-utils/blog"
)

// GetUserSimple
//
// 获取用户的简单信息；
// 如果获取成功，则返回用户的简单信息；
// 如果获取失败，则返回错误码；
// 如果获取的用户信息出现错误或获取不到用户信息，则返回内部服务器错误。
func (s *sUser) GetUserSimple(ctx context.Context) (*base.UserSimpleDTO, *berror.ErrorCode) {
	blog.ServiceInfo(ctx, "GetUserSimple", "获取用户的简单信息")

	username, errorCode := dao.System.GetSystemValue(ctx, consts.SystemUserUsernameKey)
	if errorCode != nil {
		return nil, errorCode
	}
	email, errorCode := dao.System.GetSystemValue(ctx, consts.SystemUserEmailKey)
	if errorCode != nil {
		return nil, errorCode
	}
	phone, errorCode := dao.System.GetSystemValue(ctx, consts.SystemUserPhoneKey)
	if errorCode != nil {
		return nil, errorCode
	}

	var userSimpleData = &base.UserSimpleDTO{
		Username: username,
		Email:    email,
		Phone:    phone,
	}
	return userSimpleData, nil
}

// GetUserDetail
//
// 获取用户的详细信息；
// 如果获取成功，则返回用户的详细信息；
// 如果获取失败，则返回错误码；
// 如果获取的用户信息出现错误或获取不到用户信息，则返回内部服务器错误。
func (s *sUser) GetUserDetail(ctx context.Context) (*base.UserDetailDTO, *berror.ErrorCode) {
	blog.ServiceInfo(ctx, "GetUserDetail", "获取用户的详细信息")

	simpleData, errorCode := s.GetUserSimple(ctx)
	if errorCode != nil {
		return nil, errorCode
	}

	nickname, errorCode := dao.System.GetSystemValue(ctx, consts.SystemUserNicknameKey)
	if errorCode != nil {
		return nil, errorCode
	}
	avatarType, errorCode := dao.System.GetSystemValue(ctx, consts.SystemUserAvatarTypeKey)
	if errorCode != nil {
		return nil, errorCode
	}
	avatarBase64, errorCode := dao.System.GetSystemValue(ctx, consts.SystemUserAvatarBase64Key)
	if errorCode != nil {
		return nil, errorCode
	}
	avatarURL, errorCode := dao.System.GetSystemValue(ctx, consts.SystemUserAvatarURLKey)
	if errorCode != nil {
		return nil, errorCode
	}

	var userDetailData = &base.UserDetailDTO{
		UserSimpleDTO: *simpleData,
		Nickname:      nickname,
		AvatarType:    avatarType,
		AvatarBase64:  avatarBase64,
		AvatarURL:     avatarURL,
	}
	return userDetailData, nil
}
