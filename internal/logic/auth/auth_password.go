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

package auth

import (
	"bamboo-main/internal/consts"
	"bamboo-main/internal/dao"
	"bamboo-main/internal/model/do"
	"bamboo-main/internal/model/entity"
	"bamboo-main/pkg/utility"
	"context"
	"fmt"
	"github.com/XiaoLFeng/bamboo-utils/berror"
	"github.com/XiaoLFeng/bamboo-utils/blog"
	"github.com/XiaoLFeng/bamboo-utils/butil"
	"github.com/gogf/gf/v2/database/gdb"
)

// getUserPassword
//
// 获取用户密码；
// 如果获取失败，则返回错误码；
// 如果获取的密码出现错误或获取不到密码，则返回内部服务器错误。
func getUserPassword(ctx context.Context) (string, *berror.ErrorCode) {
	getPassword, errorCode := dao.System.GetSystemValue(ctx, consts.SystemUserPasswordKey)
	if errorCode != nil {
		return "", errorCode
	}
	if getPassword == "" {
		return "", berror.ErrorAddData(&berror.ErrInternalServer, "系统错误 getUserPassword 函数出现意外错误")
	}
	return getPassword, nil
}

// VerifyPassword
//
// 验证用户的登录密码；
// 如果验证成功，则返回 nil；
// 如果验证失败，则返回错误码；
// 如果获取用户密码失败，则返回内部服务器错误。
func (s *sAuth) VerifyPassword(ctx context.Context, password string) *berror.ErrorCode {
	blog.ServiceInfo(ctx, "VerifyPassword", "验证用户的登录密码")

	// 获取用户密码
	getUserPassword, errorCode := getUserPassword(ctx)
	if errorCode != nil {
		return errorCode
	}

	// 验证密码
	if !butil.PasswordVerify(password, getUserPassword) {
		blog.ServiceNotice(ctx, "VerifyPassword", "密码错误")
		return berror.ErrorAddData(&berror.ErrUnauthorized, "密码错误")
	}
	blog.ServiceInfo(ctx, "VerifyPassword", "密码验证成功")
	return nil
}

// ChangePassword
//
// 修改用户密码；
// 如果修改成功，则返回 nil；
// 如果修改失败，则返回错误码；
// 如果新密码为空，则返回错误码。
func (s *sAuth) ChangePassword(ctx context.Context, newPassword string) *berror.ErrorCode {
	blog.ServiceInfo(ctx, "ChangePassword", "修改 %s 用户密码", utility.GetCtxVarToUserEntity(ctx).Username)

	// 验证新密码是否为空
	if newPassword == "" {
		return berror.ErrorAddData(&berror.ErrBadRequest, "新密码不能为空")
	}

	// 对新密码进行加密
	encodedPassword := butil.PasswordEncode(newPassword)

	// 更新系统值中的密码
	_, daoErr := dao.System.Ctx(ctx).Cache(gdb.CacheOption{
		Duration: -1,
		Name:     fmt.Sprintf(consts.SystemFieldsRedisKey, consts.SystemUserPasswordKey),
	}).Where(&do.System{SystemName: consts.SystemUserPasswordKey}).OmitEmpty().Update(&entity.System{SystemValue: encodedPassword})
	if daoErr != nil {
		blog.ServiceError(ctx, "ChangePassword", "更新密码失败: %v", daoErr)
		return berror.ErrorAddData(&berror.ErrInternalServer, "更新密码失败: "+daoErr.Error())
	}

	blog.ServiceNotice(ctx, "ChangePassword", "用户 %s 密码修改成功", utility.GetCtxVarToUserEntity(ctx).Username)
	return nil
}
