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

package sponsor

import (
	"context"
	"xiaoMain/internal/service"

	"xiaoMain/api/sponsor/v1"
)

// SponsorTypeGet 获取赞助类型
// 用于获取赞助类型，如果成功则返回 nil，否则返回错误。
// 本接口会根据已有的赞助类型信息对赞助类型进行获取，若获取失败返回失败信息，若成功返回成功信息
// 当用户执行接口的时候，可以获取可以选择的所有赞助类型
//
// 参数：
// ctx: 请求的上下文，用于管理超时和取消信号。
// req: 用户的请求，包含获取赞助类型的详细信息。
//
// 返回：
// res: 如果获取赞助类型成功，返回 nil；否则返回错误。
func (c *ControllerV1) SponsorTypeGet(ctx context.Context, _ *v1.SponsorTypeGetReq) (res *v1.SponsorTypeGetRes, err error) {
	sponsorType, err := service.Sponsor().GetSponsorType(ctx)
	if err == nil {
		return &v1.SponsorTypeGetRes{
			Link: sponsorType,
		}, nil
	}
	return nil, err
}
