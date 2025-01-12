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
	"github.com/bamboo-services/bamboo-utils/butil"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/google/uuid"
	v1 "xiaoMain/api/sponsor/v1"
	"xiaoMain/internal/dao"
	"xiaoMain/internal/model/do"
	"xiaoMain/internal/model/entity"
)

// AddSponsor
//
// # 赞助添加
//
// 用于添加赞助，添加一个赞助信息
//
// # 接口
//   - ctx: 上下文
//   - req: 请求参数
//
// # 返回
//   - err: 错误信息
func (s *sSponsor) AddSponsor(ctx context.Context, req *v1.SponsorAddReq) (err error) {
	g.Log().Notice(ctx, "[LOGIC] SponsorLogic:AddSponsor | 赞助添加")
	// 赞助添加
	if _, err = dao.Sponsor.Ctx(ctx).Data(do.Sponsor{
		SponsorUuid: butil.GenerateRandUUID().String(),
		Name:        req.Name,
		Url:         req.Url,
		Type:        req.Type,
		Money:       req.Money,
		Time:        req.Time,
	}).Insert(); err != nil {
		return berror.NewError(bcode.ServerInternalError, err.Error())
	}
	return nil
}

// GetSponsorByUUID
//
// # 获取赞助信息
//
// 用于获取赞助信息，从数据库中获取赞助的数据；
// 获取的数据直接输出出来；
//
// # 参数
//   - ctx: 上下文
//   - uuid: 赞助UUID
//
// # 返回
//   - err: 错误信息
func (s *sSponsor) GetSponsorByUUID(ctx context.Context, uuid uuid.UUID) (*entity.Sponsor, error) {
	g.Log().Notice(ctx, "[LOGIC] SponsorLogic:GetSponsorByUUID | 获取赞助信息")
	var sponsor *entity.Sponsor
	err := dao.Sponsor.Ctx(ctx).Where("sponsor_uuid", uuid).Scan(&sponsor)
	if err != nil {
		return nil, berror.NewError(bcode.ServerInternalError, err.Error())
	}
	return sponsor, nil
}

// EditSponsor
//
// # 赞助编辑
//
// 用于编辑赞助，编辑一个赞助信息
//
// # 接口
//   - ctx: 上下文
//   - uuid: 赞助UUID
//   - req: 请求参数
//
// # 返回
//   - err: 错误信息
func (s *sSponsor) EditSponsor(ctx context.Context, uuid uuid.UUID, req *v1.SponsorEditReq) (err error) {
	g.Log().Notice(ctx, "[LOGIC] SponsorLogic:EditSponsor | 赞助编辑")
	// 赞助编辑
	if _, err = dao.Sponsor.Ctx(ctx).Data(do.Sponsor{
		Name:  req.Name,
		Url:   req.Url,
		Type:  req.Type,
		Money: req.Money,
		Time:  req.Time,
	}).Where("sponsor_uuid", uuid).Update(); err != nil {
		return berror.NewError(bcode.ServerInternalError, err.Error())
	}
	return nil
}

// DelSponsor
//
// # 赞助删除
//
// 用于删除赞助，删除一个赞助信息
//
// # 接口
//   - ctx: 上下文
//   - uuid: 赞助UUID
//
// # 返回
//   - err: 错误信息
func (s *sSponsor) DelSponsor(ctx context.Context, uuid uuid.UUID) (err error) {
	g.Log().Notice(ctx, "[LOGIC] SponsorLogic:DelSponsor | 赞助删除")
	// 赞助删除
	if _, err = dao.Sponsor.Ctx(ctx).Where("sponsor_uuid", uuid).Delete(); err != nil {
		return berror.NewError(bcode.ServerInternalError, err.Error())
	}
	return nil
}
