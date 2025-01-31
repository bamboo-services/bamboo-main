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

package auth

import (
	"context"
	"github.com/bamboo-services/bamboo-utils/bcode"
	"github.com/bamboo-services/bamboo-utils/berror"
	"github.com/bamboo-services/bamboo-utils/butil"
	"github.com/gogf/gf/v2/frame/g"
	"xiaoMain/internal/dao"
	"xiaoMain/internal/model/do"
	"xiaoMain/internal/model/entity"
)

// ChangeUserPassword 用于修改用户密码。如果密码修改成功，将会清理用户的登录状态，需要用户重新进行登录。
// 如果用户的密码修改失败，将会返回错误信息。如果修改成功，将返回 nil。
//
// 参数:
// ctx: 上下文对象，用于传递和控制请求的生命周期。
// password: 用户新的密码字符串。
//
// 返回值:
// err: 如果密码修改成功，返回 nil。否则返回错误信息。
func (s *sAuth) ChangeUserPassword(ctx context.Context, password string) (err error) {
	g.Log().Notice(ctx, "[LOGIC] 执行 AuthLogic:ChangeUserPassword 服务层")
	// 检查用户的密码是否与前密码一致
	var getUserPassword entity.Index
	err = dao.Index.Ctx(ctx).Where(do.Index{Key: "password"}).Scan(&getUserPassword)
	if err != nil {
		return berror.NewErrorHasError(bcode.ServerInternalError, err)
	}
	// 对密码进行检查
	if butil.PasswordVerify(password, getUserPassword.Value) {
		g.Log().Warning(ctx, "[LOGIC] 用户修改的密码与原密码相同")
		return berror.NewError(bcode.OperationFailed, "新密码与原密码相同")
	}
	// 对密码进行修改
	getHashPassword := butil.PasswordEncode(password)
	_, err = dao.Index.Ctx(ctx).Data(do.Index{Value: getHashPassword}).Where(do.Index{Key: "password"}).Update()
	if err != nil {
		return berror.NewErrorHasError(bcode.ServerInternalError, err)
	}
	return nil
}
