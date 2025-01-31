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
	"github.com/gogf/gf/v2/os/gtime"
	"xiaoMain/internal/service"

	"xiaoMain/api/link/v1"
)

// LinkGetAdmin
//
// # 获取链接
//
// 获取链接，包括友情链接、社交链接等。
//
// # 参数:
//   - ctx: 上下文对象，用于传递和控制请求的生命周期。
//   - req: 请求参数，包括获取链接的请求参数。
//
// # 返回:
//   - res: 响应参数，返回获取链接的响应结果。
//   - err: 如果处理过程中发生错误，返回错误信息。
func (c *ControllerV1) LinkGetAdmin(ctx context.Context, req *v1.LinkGetAdminReq) (res *v1.LinkGetAdminRes, err error) {
	// 检查用户是否登录
	err = service.Auth().IsUserLogin(ctx)
	if err != nil {
		return nil, err
	}
	// 获取链接
	link, total, err := service.Link().GetLinkAdmin(ctx)
	if err != nil {
		return nil, err
	}
	// 待审核的链接
	var (
		reviewedNumber   uint64 = 0
		recentlyAdded    uint64 = 0
		recentlyModified uint64 = 0
	)
	for _, list := range link {
		if list.Location == 0 {
			reviewedNumber++
		}
		// 查找创建时间，找出最近七天
		if gtime.Now().Sub(list.CreatedAt).Hours() < 24*7 && list.Location != 0 {
			recentlyAdded++
		}
		// 查找修改时间，找出最近七天
		if gtime.Now().Sub(list.UpdatedAt).Hours() < 24*7 && list.Location != 0 {
			recentlyModified++
		}
	}
	return &v1.LinkGetAdminRes{
		Links:            link,
		Total:            total,
		Reviewed:         reviewedNumber,
		RecentlyAdded:    recentlyAdded,
		RecentlyModified: recentlyModified,
	}, nil
}
