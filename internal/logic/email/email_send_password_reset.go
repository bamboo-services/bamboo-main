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

package email

import (
	"bamboo-main/internal/consts"
	"bamboo-main/internal/dao"
	"bamboo-main/internal/service"
	"bamboo-main/pkg/cerror"
	"context"
	"fmt"
	"github.com/XiaoLFeng/bamboo-utils/berror"
	"github.com/XiaoLFeng/bamboo-utils/blog"
	"github.com/XiaoLFeng/bamboo-utils/butil"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/util/gconv"
	"sync"
	"time"
)

// SendMailByPasswordReset
//
// 发送重置密码的邮件；
// 如果发送成功，则返回 nil；
// 如果发送失败，则返回错误码；
// 如果获取邮件模板失败，则返回错误码。
func (s *sEmail) SendMailByPasswordReset(ctx context.Context, toMail string) *berror.ErrorCode {
	blog.ServiceInfo(ctx, "SendMailByPasswordReset", "发送重置密码邮件，收件人：%s", toMail)
	// 检查是否可再次发送验证码
	getValue, cacheErr := g.Redis().Get(ctx, fmt.Sprintf(consts.NextSendEmailTimeRedisKey, toMail))
	if cacheErr != nil {
		blog.ServiceError(ctx, "SendMailByPasswordReset", "获取下次发送邮件时间失败，收件人：%s, 错误：%v", toMail, cacheErr)
		return cerror.ErrMailSend
	}
	if !getValue.IsEmpty() {
		nextSendTime, operateErr := gtime.StrToTime(getValue.String())
		if operateErr != nil {
			blog.ServiceError(ctx, "SendMailByPasswordReset", "解析下次发送邮件时间失败，收件人：%s, 错误：%v", toMail, operateErr)
			return cerror.ErrMailSend
		}
		if nextSendTime.After(gtime.Now()) {
			blog.ServiceNotice(ctx, "SendMailByPasswordReset", "下次发送邮件时间未到，收件人：%s, 下次发送时间：%s", toMail, nextSendTime)
			return berror.ErrorAddData(cerror.ErrMailNextSendTimeNotReached, g.Map{
				"next_send_time":     nextSendTime,
				"surplus_milli_time": nextSendTime.Sub(gtime.Now()).Milliseconds(),
				"mail_to":            toMail,
			})
		}
	}

	// 数据准备
	getData := s.getBaseData(ctx)
	getData["user_name"] = toMail
	getData["verification_code"] = butil.RandomString(6)
	getData["expire_at"] = gtime.Now().Add(gconv.Duration(dao.System.MustGetSystemValue(ctx, consts.SystemEmailVerifyCodeExpireTimeKey)) * time.Second).Format("Y-m-d H:i:s")
	getData["expire_time"] = dao.System.MustGetSystemValue(ctx, consts.SystemEmailVerifyCodeExpireTimeKey)

	// 异步
	wg := sync.WaitGroup{}
	wg.Add(2)

	// 发送邮件
	go func(ctx context.Context, data g.Map) {
		_ = s.SendMail(ctx, toMail, "重置密码", "auth_user_password_reset", true, data)
		defer wg.Done()
	}(ctx, getData)
	// 验证码存入缓存
	go func(ctx context.Context, mail string, data g.Map) {
		smsWg := &sync.WaitGroup{}
		smsWg.Add(2)
		// 存入下一次可再次发送的时间
		go func(ctx context.Context, mail string) {
			errorCode := s.setNextSendTime(ctx, mail)
			if errorCode != nil {
				panic(errorCode)
			}
			defer smsWg.Done()
		}(ctx, toMail)
		// 存入验证码
		go func(ctx context.Context, mail, code string) {
			errorCode := service.Code().GenerateEmailCode(ctx, mail, &code)
			if errorCode != nil {
				panic(errorCode)
			}
			defer smsWg.Done()
		}(ctx, toMail, gconv.String(data["verification_code"]))

		smsWg.Wait()
		defer wg.Done()
	}(ctx, toMail, getData)

	// 线程等待
	wg.Wait()
	return nil
}
