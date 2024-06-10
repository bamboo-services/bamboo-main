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
	"xiaoMain/internal/service"
	"xiaoMain/utility"

	"xiaoMain/api/mail/v1"
)

func (c *ControllerV1) MailSend(ctx context.Context, req *v1.MailSendReq) (res *v1.MailSendRes, err error) {
	g.Log().Notice(ctx, "[CONTROL] MailSend | 邮件发送接口")
	// 检查邮箱是否正确
	err = service.User().IsMailHasConsoleUser(ctx, req.To)
	if err != nil {
		return nil, err
	}
	// 检查 Scene 是否正确
	if utility.GetMailTemplateByScene(req.Scene) == "" {
		return nil, berror.NewError(bcode.OperationFailed, "场景错误")
	}
	// 对场景内容进行邮件发送
	err = service.Mail().SendMail(ctx, req.To, req.Scene, req.GetVariables)
	if err != nil {
		return nil, err
	}
	return nil, nil
}
