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
	"github.com/XiaoLFeng/bamboo-utils/blog"
)

// GetUserSimple
//
// 获取用户的简单信息；
// 如果获取成功，则返回用户的简单信息；
// 如果获取失败，则会产生恐慌（一般情况下都可以正常输出）
func (s *sUser) GetUserSimple(ctx context.Context) *base.UserSimpleDTO {
	blog.ServiceInfo(ctx, "GetUserSimple", "获取用户的简单信息")

	var userSimpleData = &base.UserSimpleDTO{
		Username: dao.System.MustGetSystemValue(ctx, consts.SystemUserUsernameKey),
		Email:    dao.System.MustGetSystemValue(ctx, consts.SystemUserEmailKey),
		Phone:    dao.System.MustGetSystemValue(ctx, consts.SystemUserPhoneKey),
	}
	return userSimpleData
}

// GetUserDetail
//
// 获取用户的详细信息；
// 如果获取成功，则返回用户的详细信息；
// 如果获取失败，则会产生恐慌（一般情况下都可以正常输出）
func (s *sUser) GetUserDetail(ctx context.Context) *base.UserDetailDTO {
	blog.ServiceInfo(ctx, "GetUserDetail", "获取用户的详细信息")

	simpleData := s.GetUserSimple(ctx)

	var userDetailData = &base.UserDetailDTO{
		UserSimpleDTO: *simpleData,
		Nickname:      dao.System.MustGetSystemValue(ctx, consts.SystemUserNicknameKey),
		AvatarType:    dao.System.MustGetSystemValue(ctx, consts.SystemUserAvatarTypeKey),
		AvatarBase64:  dao.System.MustGetSystemValue(ctx, consts.SystemUserAvatarBase64Key),
		AvatarURL:     dao.System.MustGetSystemValue(ctx, consts.SystemUserAvatarURLKey),
	}
	return userDetailData
}
