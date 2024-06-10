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

package link

import (
	"context"
	"github.com/bamboo-services/bamboo-utils/bcode"
	"github.com/bamboo-services/bamboo-utils/berror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"xiaoMain/internal/model/vo"
	"xiaoMain/internal/service"
	"xiaoMain/utility/result"

	"xiaoMain/api/link/v1"
)

// LinkGetLocation
//
// # 获取期望位置信息
//
// 获取期望位置信息, 用于用户选择链接的期望位置。
//
// # 参数
//   - ctx: 请求的上下文，用于管理超时和取消信号。
//   - req: 用户的请求，包含获取期望位置信息的详细信息。
//
// # 返回
//   - res: 发送给用户的响应。如果获取期望位置信息成功，它将返回成功的消息。
//   - err: 在获取期望位置信息过程中发生的任何错误。
func (c *ControllerV1) LinkGetLocation(
	ctx context.Context,
	_ *v1.LinkGetLocationReq,
) (res *v1.LinkGetLocationRes, err error) {
	g.Log().Notice(ctx, "[CONTROL] 控制层 LinkGetLocation 接口")
	getRequest := ghttp.RequestFromCtx(ctx)
	// 获取期望位置信息
	getLocation, err := service.Link().GetLocation(ctx)
	if err != nil {
		return nil, err
	}
	if getLocation == nil {
		return nil, berror.NewError(bcode.NotExist, "没有找到期望位置信息")
	} else {
		g.Log().Debugf(ctx, "[CONTROL] 获取期望位置信息成功, 数量[%d]", len(getLocation))
		getLocationList := make([]vo.LinkLocationVO, len(getLocation))
		for i, location := range getLocation {
			getLocationList[i] = vo.LinkLocationVO{
				ID:          location.Id,
				DisplayName: location.DisplayName,
			}
		}
		result.Success("获取期望位置信息成功", getLocationList).Response(getRequest)
	}
	return nil, nil
}
