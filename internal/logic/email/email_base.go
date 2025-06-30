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
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/util/gconv"
	"time"
)

type sEmail struct {
}

func init() {
	service.RegisterEmail(New())
}

func New() *sEmail {
	return &sEmail{}
}

// getBaseData
//
// 获取邮件发送的基础数据；
// 包括系统名称、版本和描述等信息；
// 这些数据将用于邮件模板中，提供给用户关于系统的基本信息。
func (s *sEmail) getBaseData(ctx context.Context) g.Map {
	baseData := g.Map{
		"system_name":    dao.System.MustGetSystemValue(ctx, consts.SystemNameKey),
		"system_version": dao.System.MustGetSystemValue(ctx, consts.SystemVersionKey),
		"system_desc":    dao.System.MustGetSystemValue(ctx, consts.SystemDescriptionKey),
		"system_since":   dao.System.MustGetSystemValue(ctx, consts.SystemSinceYearKey),
		"system_website": dao.System.MustGetSystemValue(ctx, consts.SystemWebsiteKey),
	}
	baseData["system_copyright"] = fmt.Sprintf("Copyright &copy; %s-%d %s. All Rights Reserved.", baseData["system_since"], gtime.Now().Year(), baseData["system_name"])
	return baseData
}

// setNextSendTime
//
// 设置下次发送邮件的时间；
// 这里设置为当前时间加上1分钟，防止短时间内多次发送邮件；
// 如果设置失败，则记录错误日志并抛出异常。
func (s *sEmail) setNextSendTime(ctx context.Context, email string) *berror.ErrorCode {
	cacheErr := g.Redis().SetEX(
		ctx,
		fmt.Sprintf(consts.NextSendEmailTimeRedisKey, email),
		gtime.Now().Add(time.Minute),
		gconv.Int64(time.Minute.Seconds()),
	)
	if cacheErr != nil {
		blog.ServiceError(ctx, "SendMailByPasswordReset", "存储验证码到缓存失败，收件人：%s, 错误：%v", email, cacheErr)
		return berror.ErrorAddData(cerror.ErrMailSend, cacheErr.Error())
	}
	return nil
}
