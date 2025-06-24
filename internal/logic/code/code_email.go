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

package code

import (
	"bamboo-main/internal/consts"
	"bamboo-main/internal/dao"
	"bamboo-main/pkg/cerror"
	"context"
	"fmt"
	"github.com/XiaoLFeng/bamboo-utils/berror"
	"github.com/XiaoLFeng/bamboo-utils/blog"
	"github.com/XiaoLFeng/bamboo-utils/butil"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/util/gconv"
)

// GenerateEmailCode
//
// 生成邮箱验证码；
// 如果验证码已存在，则不重新生成；
// 如果验证码不存在，则随机生成一个6位数的验证码；
// 如果存储验证码到缓存失败，则返回错误码。
func (s *sCode) GenerateEmailCode(ctx context.Context, email string, code *string) *berror.ErrorCode {
	blog.ServiceInfo(ctx, "GenerateEmailCode", "生成 %s 的验证码", email)

	// 检查验证码是否已存在
	if code == nil || *code == "" {
		// 验证码为空，随机生成一个验证码
		*code = butil.RandomString(6)
	}

	// 对验证码进行存储
	cacheErr := g.Redis().SetEX(
		ctx,
		fmt.Sprintf(consts.SendEmailSmsRedisKey, email),
		code,
		gconv.Int64(dao.System.MustGetSystemValue(ctx, consts.SystemEmailVerifyCodeExpireTimeKey)),
	)
	if cacheErr != nil {
		blog.ServiceError(ctx, "SendMailByPasswordReset", "存储验证码到缓存失败，收件人：%s, 错误：%v", email, cacheErr)
		return cerror.ErrMailSend
	}
	return nil
}

// VerifyEmailCode
//
// 验证邮箱验证码；
// 如果验证码验证成功，则删除缓存中的验证码；
// 如果验证码不存在，则返回错误码 cerror.ErrMailCodeNotExist；
// 如果验证码不匹配，则返回错误码 cerror.ErrMailCodeInvalid；
// 如果获取验证码失败，则返回错误码 berror.ErrInternalServer。
// 如果删除验证码缓存失败，则返回错误码 berror.ErrInternalServer。
// 如果验证成功，则返回 nil。
func (s *sCode) VerifyEmailCode(ctx context.Context, email, code string) *berror.ErrorCode {
	blog.ServiceInfo(ctx, "VerifyEmailCode", "验证 %s 的验证码", email)

	// 从缓存中获取验证码
	cacheValue, cacheErr := g.Redis().Get(ctx, fmt.Sprintf(consts.SendEmailSmsRedisKey, email))
	if cacheErr != nil {
		blog.ServiceError(ctx, "VerifyEmailCode", "获取验证码失败，收件人：%s, 错误：%v", email, cacheErr)
		return berror.ErrorAddData(&berror.ErrInternalServer, cacheErr.Error())
	}

	if cacheValue.IsEmpty() {
		blog.ServiceNotice(ctx, "VerifyEmailCode", "验证码不存在，收件人：%s", email)
		return cerror.ErrMailCodeNotExist
	}

	if code != cacheValue.String() {
		blog.ServiceNotice(ctx, "VerifyEmailCode", "验证码不匹配，收件人：%s, 提交的验证码：%s, 缓存中的验证码：%s", email, code, cacheValue.String())
		return cerror.ErrMailCodeInvalid
	}

	// 验证码验证成功后，删除缓存中的验证码
	_, cacheErr = g.Redis().Del(ctx, fmt.Sprintf(consts.SendEmailSmsRedisKey, email))
	if cacheErr != nil {
		blog.ServiceError(ctx, "VerifyEmailCode", "删除验证码缓存失败，收件人：%s, 错误：%v", email, cacheErr)
		return berror.ErrorAddData(&berror.ErrInternalServer, cacheErr.Error())
	}

	return nil
}
