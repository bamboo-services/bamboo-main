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
	v1 "xiaoMain/api/link/v1"
	"xiaoMain/internal/dao"
	"xiaoMain/internal/model/do"
	"xiaoMain/internal/model/entity"
)

// AddLink
//
// # 添加链接
//
// 用于添加链接，如果添加成功则返回 nil，否则返回错误。
//
// # 参数:
//   - ctx: 上下文对象，用于传递和控制请求的生命周期。
//   - req: 添加链接请求(v1.LinkAddReq)
//
// # 返回:
//   - err: 如果添加链接成功，返回 nil；否则返回错误.
func (s *sLink) AddLink(ctx context.Context, req v1.LinkAddReq) (err error) {
	g.Log().Notice(ctx, "[LOGIC] Link:AddLink | 添加链接")
	var getLocation *do.Location
	var getColor *do.Color
	// 对颜色进行数据库获取，获取指定的颜色是否存在
	err = dao.Location.Ctx(ctx).Where(do.Location{Id: req.DesiredLocation}).Scan(&getLocation)
	if err != nil {
		g.Log().Warning(ctx, "[LOGIC] 数据库错误<期望位置不存在>")
		return berror.NewErrorHasError(bcode.ServerInternalError, err)
	}
	if getLocation == nil {
		g.Log().Warningf(ctx, "[LOGIC] 期望位置不存在[%v]", req.DesiredLocation)
		return berror.NewError(bcode.NotExist, "期望位置不存在")
	}
	// 对颜色进行数据库获取，获取指定的颜色是否存在
	if dao.Color.Ctx(ctx).Where(do.Color{Id: req.DesiredColor}).Scan(&getColor) != nil {
		g.Log().Warning(ctx, "[LOGIC] 数据库错误<期望位置不存在>")
		return berror.NewErrorHasError(bcode.ServerInternalError, err)
	}
	if getColor == nil {
		g.Log().Warningf(ctx, "[LOGIC] 期望颜色不存在[%v]", req.DesiredColor)
		return berror.NewError(bcode.NotExist, "期望颜色不存在")
	}
	// 对数据进行插入
	if _, err = dao.LinkList.Ctx(ctx).Insert(do.LinkList{
		WebmasterEmail:  req.WebmasterEmail,
		ServiceProvider: req.ServiceProvider,
		SiteName:        req.SiteName,
		SiteUrl:         req.SiteURL,
		SiteLogo:        req.SiteLogo,
		SiteDescription: req.SiteDescription,
		SiteRssUrl:      req.SiteRssURL,
		HasAdv:          false,
		DesiredLocation: getLocation.Id,
		DesiredColor:    getColor.Id,
		WebmasterRemark: req.Remark,
	}); err != nil {
		g.Log().Warningf(ctx, "[LOGIC] 添加链接失败[%s]", err.Error())
		return berror.NewErrorHasError(bcode.ServerInternalError, err)
	}
	return nil
}

// GetColor 获取期望颜色信息
// 用于获取期望颜色信息, 如果成功则返回期望颜色信息，否则返回错误。
//
// 参数：
// ctx: 请求的上下文，用于管理超时和取消信号。
//
// 返回：
// getColors: 如果获取期望颜色信息成功，返回期望颜色信息；否则返回错误。
// err: 如果获取期望颜色信息成功，返回 nil；否则返回错误。其中错误的返回信息在此函数中主要包含内容为数据库错误。
func (s *sLink) GetColor(ctx context.Context) (getColors []*entity.Color, err error) {
	g.Log().Notice(ctx, "[LOGIC] 执行 Link:GetColor 服务层")
	err = dao.Color.Ctx(ctx).Where(do.Color{HasSelect: true}).Scan(&getColors)
	if err != nil {
		g.Log().Error(ctx, "[LOGIC] 数据库错误<获取期望颜色失败>")
		return nil, err
	}
	return getColors, nil
}

// GetLocation
//
// # 获取位置信息
//
// 用于获取一些位置信息，这些位置信息的查询受限于是否公开。如果成功则返回期望位置信息，否则返回错误。若查询的位置信息是公开的，则返回所有位置信息。
// 若查询的位置信息是不公开的，则返回所有不公开的位置信息。公开数据表参数 Reveal 为 true，不公开数据表参数 Reveal 为 false。
//
// # 参数:
//   - ctx: 请求的上下文，用于管理超时和取消信号。
//
// # 返回:
//   - getLocation: 如果获取期望位置信息成功，返回期望位置信息；否则返回错误.
//   - err: 如果获取期望位置信息成功，返回 nil；否则返回错误.其中错误的返回信息在此函数中主要包含内容为数据库错误。
func (s *sLink) GetLocation(ctx context.Context) (getLocation []*entity.Location, err error) {
	g.Log().Notice(ctx, "[LOGIC] Link:GetLocation | 获取期望位置信息")
	err = dao.Location.Ctx(ctx).Where(do.Location{Reveal: true}).Scan(&getLocation)
	if err != nil {
		g.Log().Error(ctx, "[LOGIC] 数据库错误<获取期望位置失败>")
		return nil, berror.NewErrorHasError(bcode.ServerInternalError, err)
	}
	return getLocation, nil
}

// GetLocationNoReveal
//
// # 获取所有位置信息
//
// 用于获取所有位置信息, 如果成功则返回所有位置信息，否则返回错误。
//
// # 参数:
//   - ctx: 请求的上下文，用于管理超时和取消信号。
//
// # 返回:
//   - getLocation: 如果获取所有位置信息成功，返回所有位置信息；否则返回错误.
//   - err: 如果获取所有位置信息成功，返回 nil；否则返回错误.其中错误的返回信息在此函数中主要包含内容为数据库错误。
func (s *sLink) GetLocationNoReveal(ctx context.Context) (getLocation []*entity.Location, err error) {
	g.Log().Notice(ctx, "[LOGIC] Link:GetLocationNoReveal | 获取所有位置信息")
	err = dao.Location.Ctx(ctx).OrderAsc(dao.Location.Columns().Sort).Scan(&getLocation)
	if err != nil {
		g.Log().Error(ctx, "[LOGIC] 数据库错误<获取期望位置失败>")
		return nil, berror.NewErrorHasError(bcode.ServerInternalError, err)
	}
	return getLocation, nil
}

// AddLocation
//
// # 添加链接位置
//
// 用于添加链接位置，如果成功则返回 nil，否则返回错误。
//
// # 参数:
//   - ctx: 请求的上下文，用于管理超时和取消信号。
//   - name: 用户尝试添加的位置名称。
//   - displayName: 用户尝试添加的位置显示名称。
//   - description: 用户尝试添加的位置描述。
//   - reveal: 用户尝试添加的位置是否公开。
//   - sort: 用户尝试添加的位置排序。
//
// # 返回:
//   - err: 如果添加链接位置成功，返回 nil；否则返回错误.
func (s *sLink) AddLocation(
	ctx context.Context,
	name string,
	displayName string,
	description string,
	reveal bool,
	sort int,
) (err error) {
	g.Log().Notice(ctx, "[LOGIC] Link:AddLocation | 添加链接位置")
	_, err = dao.Location.Ctx(ctx).Data(do.Location{
		Name:        name,
		DisplayName: displayName,
		Description: description,
		Reveal:      reveal,
		Sort:        sort,
	}).Insert()
	if err != nil {
		g.Log().Errorf(ctx, "[LOGIC] 添加链接位置失败[%s]", err.Error())
		return berror.NewErrorHasError(bcode.ServerInternalError, err)
	}
	return nil
}

// AddColor
//
// # 添加链接颜色
//
// 用于添加链接颜色，如果成功则返回 nil，否则返回错误。
//
// # 参数:
//   - ctx: 请求的上下文，用于管理超时和取消信号。
//   - name: 用户尝试添加的颜色名称。
//   - displayName: 用户尝试添加的颜色显示名称。
//   - color: 用户尝试添加的颜色。
//   - hasSelect: 用户尝试添加的颜色是否可选。
//
// # 返回:
//   - err: 如果添加链接颜色成功，返回 nil；否则返回错误.
func (s *sLink) AddColor(
	ctx context.Context,
	name string,
	displayName string,
	color string,
	hasSelect bool,
) (err error) {
	g.Log().Notice(ctx, "[LOGIC] Link:AddColor | 添加链接颜色")
	_, err = dao.Color.Ctx(ctx).Data(do.Color{
		Name:        name,
		DisplayName: displayName,
		Color:       color,
		HasSelect:   hasSelect,
	}).Insert()
	if err != nil {
		g.Log().Noticef(ctx, "[LOGIC] 添加链接颜色失败[%s]", err.Error())
		return berror.NewErrorHasError(bcode.ServerInternalError, err)
	}
	return nil
}
