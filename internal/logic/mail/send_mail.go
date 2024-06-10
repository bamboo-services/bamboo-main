/*
 * --------------------------------------------------------------------------------
 * Copyright (c) 2016-NOW(至今) 筱锋
 * Author: 筱锋(https://www.x-lf.com)
 *
 * 本文件包含 XiaoMain 的源代码，该项目的所有源代码均遵循MIT开源许可证协议。
 * --------------------------------------------------------------------------------
 * 许可证声明：
 *
 * 版权所有 (c) 2016-2024 筱锋。保留所有权利。
 *
 * 本软件是“按原样”提供的，没有任何形式的明示或暗示的保证，包括但不限于
 * 对适销性、特定用途的适用性和非侵权性的暗示保证。在任何情况下，
 * 作者或版权持有人均不承担因软件或软件的使用或其他交易而产生的、
 * 由此引起的或以任何方式与此软件有关的任何索赔、损害或其他责任。
 *
 * 使用本软件即表示您了解此声明并同意其条款。
 *
 * 有关MIT许可证的更多信息，请查看项目根目录下的LICENSE文件或访问：
 * https://opensource.org/licenses/MIT
 * --------------------------------------------------------------------------------
 * 免责声明：
 *
 * 使用本软件的风险由用户自担。作者或版权持有人在法律允许的最大范围内，
 * 对因使用本软件内容而导致的任何直接或间接的损失不承担任何责任。
 * --------------------------------------------------------------------------------
 */

package mail

import (
	"context"
	"github.com/bamboo-services/bamboo-utils/bcode"
	"github.com/bamboo-services/bamboo-utils/berror"
	"github.com/gogf/gf/v2/frame/g"
	"gopkg.in/gomail.v2"
	"strings"
	"xiaoMain/internal/constants"
	"xiaoMain/internal/dao"
	"xiaoMain/internal/model/do"
	"xiaoMain/internal/model/entity"
	"xiaoMain/internal/model/vo"
	"xiaoMain/utility"
)

// sendMail
//
// # 发送邮件
//
// 用于发送邮件，如果发送成功则返回 nil，否则返回错误信息。会根据传入的模板进行邮件的发送。
//
// # 参数:
//   - ctx: 上下文对象，用于传递和控制请求的生命周期。
//   - mail: 邮箱地址(string)
//   - template: 邮件模板(constants.Scene)
//   - data: 邮件发送数据(vo.MailSendData)
//
// # 返回:
//   - err: 如果发送过程中发生错误，返回错误信息。否则返回 nil.
func (s *sMail) sendMail(
	ctx context.Context,
	mail string,
	template constants.Scene,
	data vo.MailSendData,
) (err error) {
	g.Log().Notice(ctx, "[LOGIC] Mail:sendMail | 发送邮件")

	// 数据库邮件模板
	var getMailTemplateTitle entity.Index
	var getMailTemplate entity.Index
	err = dao.Index.Ctx(ctx).
		Where(do.Index{Key: "mail_template_" + utility.GetMailTemplateByScene(template) + "_title"}).
		Limit(1).Scan(&getMailTemplateTitle)
	if err != nil {
		return berror.NewError(bcode.ServerInternalError, "未查询到邮件模板标题")
	}
	err = dao.Index.Ctx(ctx).
		Where(do.Index{Key: "mail_template_" + utility.GetMailTemplateByScene(template)}).
		Limit(1).Scan(&getMailTemplate)
	if err != nil {
		return berror.NewError(bcode.ServerInternalError, "未查询到邮件模板标题")
	}

	// 对邮件模板的替换
	mailTemplate := getMailTemplate.Value
	mailTemplate = strings.ReplaceAll(mailTemplate, "%XiaoMain%", data.XiaoMain)
	mailTemplate = strings.ReplaceAll(mailTemplate, "%Email%", data.Email)
	mailTemplate = strings.ReplaceAll(mailTemplate, "%Code%", data.Code)
	mailTemplate = strings.ReplaceAll(mailTemplate, "%DateTime%", data.DateTime.String())
	mailTemplate = strings.ReplaceAll(mailTemplate, "%Copyright%", data.Copyright)

	// 创建一个新的消息
	sendMail := gomail.NewMessage()
	// 设置需要发送的人
	sendMail.SetHeader("To", mail)
	sendMail.SetHeader("From", constants.SMTPUser)
	sendMail.SetHeader("Subject", getMailTemplateTitle.Value)
	sendMail.SetBody("text/html", mailTemplate)

	// 配置邮件服务器
	dialer := gomail.NewDialer(constants.SMTPHost, utility.GetMailSendPort(ctx), constants.SMTPUser, constants.SMTPPass)
	// 发送邮件
	return dialer.DialAndSend(sendMail)
}
