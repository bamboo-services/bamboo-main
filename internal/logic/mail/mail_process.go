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
	"github.com/bamboo-services/bamboo-utils/butil"
	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
	"sync"
	"xiaoMain/internal/constants"
	"xiaoMain/internal/dao"
	"xiaoMain/internal/model/do"
)

// sendEmailVerificationCode
//
// # 发送邮件验证码
//
// 用于发送邮件验证码，如果发送成功则返回 nil，否则返回具体的报错信息。会根据传入的场景进行邮件的发送。
//
// # 参数:
//   - ctx: 上下文对象，用于传递和控制请求的生命周期。
//   - mail: 邮箱地址(string)
//   - scenes: 场景(constants.Scene)
//
// # 返回:
//   - err: 如果发送过程中发生错误，返回错误信息。否则返回 nil.
func (s *sMail) sendEmailVerificationCode(ctx context.Context, mail string, scenes constants.Scene) (err error) {
	g.Log().Notice(ctx, "[LOGIC] Mail:sendEmailVerificationCode | 发送邮件验证码")
	wg := sync.WaitGroup{}
	// 验证码存入数据库
	err = dao.VerificationCode.Ctx(ctx).Transaction(ctx, func(_ context.Context, tx gdb.TX) error {
		sendMailData := g.Map{
			"Code":      butil.RandomString(6),
			"Email":     "gm@x-lf.cn",
			"XiaoMain":  "XiaoMain",
			"Copyright": "CopyRight (C) 2016-2024 筱锋 All Rights Reserved.",
			"DateTime":  gtime.Now().String(),
		}
		// 异步创建验证码的发送
		wg.Add(1)
		go func(ctx context.Context, mail string, scene constants.Scene, data g.Map) {
			// 创建验证码并发送
			err := s.sendMail(ctx, mail, scene, data)
			if err != nil {
				g.Log().Warningf(ctx, "[LOGIC] 发送邮件时发生错误，不进行数据库插入，错误原因：%s", err.Error())
			}
			wg.Done()
		}(ctx, mail, scenes, sendMailData)
		// 存入验证码
		_, err = tx.Insert(dao.VerificationCode.Table(), do.VerificationCode{
			Type:      true,
			Contact:   mail,
			Code:      sendMailData["Code"].(string),
			Scenes:    string(scenes),
			ExpiredAt: gtime.NewFromTimeStamp(gtime.TimestampMicro() + 300000),
		})
		if err != nil {
			g.Log().Warningf(ctx, "[LOGIC] 发送邮件时发生错误，回滚 %s 数据表中发送给 %s 的数据，错误原因：%s",
				dao.VerificationCode.Table(),
				mail,
				err.Error(),
			)
			return err
		}
		return nil
	})
	if err != nil {
		return berror.NewErrorHasError(bcode.ServerInternalError, err)
	}
	// 发送成功
	wg.Wait()
	return nil
}
