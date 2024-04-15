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
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/os/glog"
	"xiaoMain/internal/service"
	"xiaoMain/utility/result"

	"xiaoMain/api/auth/v1"
)

func (c *ControllerV1) AuthChangePassword(
	ctx context.Context,
	req *v1.AuthChangePasswordReq,
) (res *v1.AuthChangePasswordRes, err error) {
	glog.Info(ctx, "[CONTROL] 控制层 UserChangePassword 接口")
	// 获取 Request
	getRequest := ghttp.RequestFromCtx(ctx)
	// 检查用户登录是否有效
	hasLogin, message := service.AuthLogic().IsUserLogin(ctx)
	if !hasLogin {
		result.NotLoggedIn.SetErrorMessage(message).Response(getRequest)
		return nil, nil
	}

	// 检查用户邮箱是否正确
	hasCheck, info := service.UserMailLogic().CheckUserMail(ctx, req.Email)
	if !hasCheck {
		result.RequestBodyValidationError.SetErrorMessage(info).Response(getRequest)
		return nil, nil
	}

	// 检查验证码是否正确
	isCorrect, info := service.MailUserLogic().
		VerificationCodeHasCorrect(ctx, req.Email, req.EmailCode, "ChangePassword")
	if !isCorrect {
		result.VerificationFailed.SetErrorMessage(info).Response(getRequest)
		return nil, nil
	}

	// 对密码进行修改
	err = service.AuthLogic().ChangeUserPassword(ctx, req.NewPassword)
	if err != nil {
		if err.Error() == "密码与原密码相同" {
			result.VerificationFailed.SetErrorMessage(err.Error()).Response(getRequest)
		} else {
			result.ServerInternalError.Response(getRequest)
		}
	}
	result.Success("密码修改成功", nil)
	return nil, nil
}
