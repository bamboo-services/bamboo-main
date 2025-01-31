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
	"github.com/bamboo-services/bamboo-utils/bcode"
	"github.com/bamboo-services/bamboo-utils/berror"
	"xiaoMain/internal/service"

	"xiaoMain/api/sponsor/v1"
)

// SponsorTypeAdd 添加赞助商类型
// 用于添加赞助商类型，如果赞助商类型已存在则返回错误
// 用于添加其他额外自定义的赞助类型信息
//
// # 接口
//   - ctx: 上下文
//   - req: 请求参数
//
// # 返回
//   - res: 返回结果
//   - err: 错误信息
func (c *ControllerV1) SponsorTypeAdd(ctx context.Context, req *v1.SponsorTypeAddReq) (res *v1.SponsorTypeAddRes, err error) {
	err = service.Auth().IsUserLogin(ctx)
	if err != nil {
		return nil, err
	}
	sponsorType, err := service.Sponsor().GetSingleSponsorTypeByName(ctx, req.Name)
	if err == nil {
		if sponsorType == nil {
			err = service.Sponsor().AddSponsorType(ctx, req)
			if err == nil {
				return &v1.SponsorTypeAddRes{}, nil
			} else {
				return nil, err
			}
		} else {
			return nil, berror.NewError(bcode.AlreadyExists, "赞助商类型已存在")
		}
	}
	return nil, err
}
