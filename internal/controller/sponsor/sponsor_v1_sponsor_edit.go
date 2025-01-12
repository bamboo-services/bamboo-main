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

// SponsorEdit 编辑赞助
// 用于编辑赞助，如果成功则返回 nil，否则返回错误。
// 本接口会根据用户提供的赞助信息进行编辑，若编辑失败返回失败信息，若成功返回成功信息
//
// 参数：
// ctx: 请求的上下文，用于管理超时和取消信号。
// req: 用户的请求，包含编辑赞助的详细信息。
//
// 返回：
// res: 如果编辑赞助成功，返回 nil；否则返回错误。
func (c *ControllerV1) SponsorEdit(ctx context.Context, req *v1.SponsorEditReq) (res *v1.SponsorEditRes, err error) {
	sponsor, err := service.Sponsor().GetSponsorByUUID(ctx, req.Uuid)
	if err == nil {
		if sponsor == nil {
			return nil, berror.NewError(bcode.NotExist, "没有找到赞助信息")
		}
		err := service.Sponsor().EditSponsor(ctx, req.Uuid, req)
		if err == nil {
			return &v1.SponsorEditRes{}, nil
		} else {
			return nil, err
		}
	}
	return nil, err
}
