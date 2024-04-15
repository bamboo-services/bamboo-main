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
	"errors"
	"github.com/gogf/gf/v2/os/glog"
	"xiaoMain/internal/dao"
	"xiaoMain/internal/model/do"
	"xiaoMain/internal/model/entity"
	"xiaoMain/utility"
)

// ChangeUserPassword
// 修改用户密码，若密码修改完毕后。将会清理掉用户的登录状态，需要用户重新进行登录；
// 若用户的密码修改失败，将会返回错误信息，若正确将返回 nil。
func (s *sAuthLogic) ChangeUserPassword(ctx context.Context, password string) (err error) {
	glog.Info(ctx, "[LOGIC] 执行 AuthLogic:ChangeUserPassword 服务层")
	// 检查用户的密码是否与前密码一致
	var getUserPassword entity.XfIndex
	if dao.XfIndex.Ctx(ctx).Where(do.XfIndex{Key: "password"}).Scan(&getUserPassword) != nil {
		glog.Error(ctx, "[LOGIC] 数据库读取错误")
		return errors.New("数据库错误")
	}
	// 对密码进行检查
	if utility.PasswordVerify(getUserPassword.Value, password) {
		glog.Info(ctx, "[LOGIC] 用户修改的密码与原密码相同")
		return errors.New("密码与原密码相同")
	}
	// 对密码进行修改
	getHashPassword := utility.PasswordEncode(password)
	_, err = dao.XfIndex.Ctx(ctx).Data(do.XfIndex{Value: getHashPassword}).Where(do.XfIndex{Key: "password"}).Update()
	if err != nil {
		glog.Error(ctx, "[LOGIC] 数据库写入错误", err)
		return err
	}
	return nil
}
