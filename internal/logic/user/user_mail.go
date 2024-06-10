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

package user

import (
	"context"
	"github.com/bamboo-services/bamboo-utils/bcode"
	"github.com/bamboo-services/bamboo-utils/berror"
	"github.com/gogf/gf/v2/frame/g"
	"xiaoMain/internal/dao"
	"xiaoMain/internal/model/do"
	"xiaoMain/internal/model/entity"
)

// CheckMailHasConsoleUser
//
// # 检查邮箱是否有控制台用户
//
// 用于检查邮箱是否有控制台用户，如果有则返回 true，否则返回 false.
//
// # 参数:
//   - ctx: 上下文对象，用于传递和控制请求的生命周期。
//   - email: 邮箱地址(string)
//
// # 返回:
//   - checkMail: 如果邮箱有控制台用户，返回 true。否则返回 false.
//   - err: 如果检查过程中发生错误，返回错误信息。否则返回 nil.
func (s *sUser) CheckMailHasConsoleUser(ctx context.Context, email string) (checkMail bool, err error) {
	g.Log().Notice(ctx, "[LOGIC] UserLogic:CheckMailHasConsoleUser | 检查邮箱是否有控制台用户")
	// 从数据库获取指定信息
	var getAdminEmail *entity.Index
	if dao.Index.Ctx(ctx).Where(do.Index{Key: "email"}).Scan(&getAdminEmail) != nil {
		return false, berror.NewErrorHasError(bcode.ServerInternalError, err)
	}
	// 对邮箱进行匹配
	if getAdminEmail != nil {
		if getAdminEmail.Value == email {
			return true, nil
		} else {
			return false, berror.NewError(bcode.ServerInternalError, "邮箱不匹配")
		}
	} else {
		return false, berror.NewError(bcode.ServerInternalError, "未查询到邮箱")
	}
}
