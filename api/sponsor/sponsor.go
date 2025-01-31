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

// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package sponsor

import (
	"context"

	"xiaoMain/api/sponsor/v1"
)

type ISponsorV1 interface {
	SponsorAdd(ctx context.Context, req *v1.SponsorAddReq) (res *v1.SponsorAddRes, err error)
	SponsorDel(ctx context.Context, req *v1.SponsorDelReq) (res *v1.SponsorDelRes, err error)
	SponsorEdit(ctx context.Context, req *v1.SponsorEditReq) (res *v1.SponsorEditRes, err error)
	SponsorGet(ctx context.Context, req *v1.SponsorGetReq) (res *v1.SponsorGetRes, err error)
	SponsorTypeAdd(ctx context.Context, req *v1.SponsorTypeAddReq) (res *v1.SponsorTypeAddRes, err error)
	SponsorTypeDel(ctx context.Context, req *v1.SponsorTypeDelReq) (res *v1.SponsorTypeDelRes, err error)
	SponsorTypeEdit(ctx context.Context, req *v1.SponsorTypeEditReq) (res *v1.SponsorTypeEditRes, err error)
	SponsorTypeGet(ctx context.Context, req *v1.SponsorTypeGetReq) (res *v1.SponsorTypeGetRes, err error)
}
