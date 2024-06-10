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

package boot

import (
	"context"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/util/gconv"
	"xiaoMain/internal/constants"
	"xiaoMain/internal/dao"
	"xiaoMain/internal/model/do"
	"xiaoMain/internal/model/entity"
)

// InitCommonData 是一个初始化常用数据的函数。
// 它从数据库中读取邮件服务器的相关信息，并将这些信息存储在内存中，以便后续使用。
// 这个函数在应用程序的启动过程中被调用。
//
// 参数:
// ctx: context.Context 对象，用于管理 go 协程和其他与上下文相关的任务。
//
// 返回: 无
func InitCommonData(ctx context.Context) {
	// 从数据库读取邮件准备信息进入内存
	var getSMTPHost entity.Index
	err := dao.Index.Ctx(ctx).Where(do.Index{Key: "smtp_host"}).Scan(&getSMTPHost)
	if err != nil {
		g.Log().Panic(ctx, "[INIT] 获取邮件服务器地址失败")
	}
	var getSMTPPortTLS entity.Index
	err = dao.Index.Ctx(ctx).Where(do.Index{Key: "smtp_port_tls"}).Scan(&getSMTPPortTLS)
	if err != nil {
		g.Log().Panic(ctx, "[INIT] 获取邮件服务器端口失败")
	}
	var getSMTPPortSSL entity.Index
	err = dao.Index.Ctx(ctx).Where(do.Index{Key: "smtp_port_ssl"}).Scan(&getSMTPPortSSL)
	if err != nil {
		g.Log().Panic(ctx, "[INIT] 获取邮件服务器端口失败")
	}
	var getSMTPUser entity.Index
	err = dao.Index.Ctx(ctx).Where(do.Index{Key: "smtp_user"}).Scan(&getSMTPUser)
	if err != nil {
		g.Log().Panic(ctx, "[INIT] 获取邮件服务器用户名失败")
	}
	var getSMTPPass entity.Index
	err = dao.Index.Ctx(ctx).Where(do.Index{Key: "smtp_pass"}).Scan(&getSMTPPass)
	if err != nil {
		g.Log().Panic(ctx, "[INIT] 获取邮件服务器密码失败")
	}

	// 数据写入 const
	constants.SMTPHost = getSMTPHost.Value
	constants.SMTPPortTLS = gconv.Int(getSMTPPortTLS.Value)
	constants.SMTPPortSSL = gconv.Int(getSMTPPortSSL.Value)
	constants.SMTPUser = getSMTPUser.Value
	constants.SMTPPass = getSMTPPass.Value
}
