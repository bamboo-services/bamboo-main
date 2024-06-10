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

package utility

import (
	"context"
	"github.com/bamboo-services/bamboo-utils/bcode"
	"github.com/bamboo-services/bamboo-utils/berror"
	"github.com/bamboo-services/bamboo-utils/butil"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/google/uuid"
	"xiaoMain/internal/constants"
)

// GetAuthorizationFromHeader
//
// # 获取请求头中的 Authorization
//
// 用于获取请求头中的 Authorization 字段，如果存在则返回去除 Bearer 后的值，否则返回空字符串。
//
// # 参数:
//   - ctx: 上下文对象，用于传递和控制请求的生命周期。
func GetAuthorizationFromHeader(ctx context.Context) (string, error) {
	getAuthorization := g.RequestFromCtx(ctx).GetHeader("Authorization")
	if getAuthorization != "" {
		authorizationCode := butil.TokenRemoveBearer(getAuthorization)
		getUUID, err := uuid.Parse(authorizationCode)
		if err != nil {
			g.Log().Warning(ctx, "[UTIL] 获取用户授权异常 | UUID错误")
			return "", berror.NewError(bcode.OperationFailed, "UUID错误")
		} else {
			return getUUID.String(), nil
		}
	} else {
		return "", berror.NewError(bcode.NotExist, "无授权头")
	}
}

// GetMailTemplateByScene
//
// # 根据场景获取邮件模板
//
// 用于根据场景获取邮件模板，如果存在则返回模板内容，否则返回空字符串。
//
// # 参数:
//   - scene: 场景(constants.Scene)
//
// # 返回:
//   - string: 如果存在则返回模板内容，否则返回空字符串。
func GetMailTemplateByScene(scene constants.Scene) string {
	// 从 Scene 获取模板
	for _, template := range constants.MailTemplate {
		if template.Name == string(scene) {
			return template.Data
		}
	}
	return ""
}

// GetMailSendPort
//
// # 获取邮件发送端口
//
// 用于获取邮件发送端口，如果是 SSL 则返回 constants.SMTPPortSSL，否则返回 constants.SMTPPortTLS。
//
// # 参数:
//   - ctx: 上下文对象，用于传递和控制请求的生命周期。
//
// # 返回:
//   - int: 如果是 SSL 则返回 465，否则返回 587。
func GetMailSendPort(ctx context.Context) int {
	if butil.GetVisitorProtocol(ctx) == "SSL" {
		return constants.SMTPPortSSL
	} else {
		return constants.SMTPPortTLS
	}
}
