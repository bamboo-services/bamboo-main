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
	"github.com/gogf/gf/v2/os/glog"
	"xiaoMain/internal/dao"
	"xiaoMain/internal/model/do"
	"xiaoMain/internal/model/entity"
)

// CheckMailHasConsoleUser 这个函数的主要作用是检查用户输入的邮箱地址是否与数据库中存储的邮箱地址匹配
// 在具体实现中，函数首先从数据库中获取指定的邮箱信息，然后将用户输入的邮箱地址与数据库中的邮箱地址进行比较。如果两者匹配，函数将返回 true 和
// "邮箱匹配" 的信息；如果两者不匹配，函数将返回 false 和 "邮箱不匹配" 的信息。如果在查询数据库过程中出现错误，函数将返回 false 和 "未查询
// 到邮箱" 的信息。
//
// 参数:
// ctx: 请求的上下文，用于管理超时和取消信号。
// email: 用户输入的邮箱地址。
//
// 返回:
// checkMail: 如果用户输入的邮箱与数据库中的邮箱匹配，则返回 true，否则返回 false。
// info: 返回的信息，如果邮箱匹配，返回 "邮箱匹配"，如果不匹配，返回 "邮箱不匹配"，如果查询过程中出现错误，返回 "未查询到邮箱"。
func (s *sUserLogic) CheckMailHasConsoleUser(ctx context.Context, email string) (checkMail bool, info string) {
	glog.Info(ctx, "[LOGIC] 执行 UserLogic:CheckMailHasConsoleUser 服务层")
	// 从数据库获取指定信息
	var getAdminEmail entity.XfIndex
	if dao.XfIndex.Ctx(ctx).Where(do.XfIndex{Key: "email"}).Scan(&getAdminEmail) != nil {
		return false, "未查询到邮箱"
	}
	// 对邮箱进行匹配
	if getAdminEmail.Value == email {
		// 返回正确信息
		return true, "邮箱匹配"
	} else {
		return false, "邮箱不匹配"
	}
}
