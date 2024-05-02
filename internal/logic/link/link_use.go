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
	"errors"
	"github.com/gogf/gf/v2/os/glog"
	v1 "xiaoMain/api/link/v1"
	"xiaoMain/internal/dao"
	"xiaoMain/internal/model/do"
	"xiaoMain/internal/model/entity"
)

// AddLink 添加链接
// 用于添加链接，如果添加成功则返回 nil，否则返回错误。
//
// 参数：
// ctx: 请求的上下文，用于管理超时和取消信号。
// req: 用户的请求，包含添加链接的详细信息。
//
// 返回：
// err: 如果添加链接成功，返回 nil；否则返回错误。其中主要包含的错误数据库错误以及对应内容不存在，返回内容均为自定义描述值
func (s *sLinkLogic) AddLink(ctx context.Context, req v1.LinkAddReq) (err error) {
	glog.Info(ctx, "[LOGIC] 执行 LinkLogic:AddLink 服务层")
	var getLocation *do.XfLocation
	var getColor *do.XfColor
	// 对颜色进行数据库获取，获取指定的颜色是否存在
	if dao.XfLocation.Ctx(ctx).Where(do.XfLocation{Name: req.DesiredLocation}).Scan(&getLocation) != nil {
		glog.Info(ctx, "[LOGIC] 数据库错误<期望位置不存在>")
		return errors.New("期望位置不存在")
	}
	if getLocation == nil {
		glog.Infof(ctx, "[LOGIC] 期望位置不存在[%s]", req.DesiredLocation)
		return errors.New("期望位置不存在")
	}
	// 对颜色进行数据库获取，获取指定的颜色是否存在
	if dao.XfColor.Ctx(ctx).Where(do.XfColor{Name: req.DesiredColor}).Scan(&getColor) != nil {
		glog.Info(ctx, "[LOGIC] 数据库错误<期望位置不存在>")
		return errors.New("期望颜色不存在")
	}
	if getColor == nil {
		glog.Infof(ctx, "[LOGIC] 期望颜色不存在[%s]", req.DesiredColor)
		return errors.New("期望颜色不存在")
	}
	// 对数据进行插入
	if _, err = dao.XfLinkList.Ctx(ctx).Insert(do.XfLinkList{
		WebmasterEmail:  req.WebmasterEmail,
		ServiceProvider: req.ServiceProvider,
		SiteName:        req.SiteName,
		SiteUrl:         req.SiteURL,
		SiteLogo:        req.SiteLogo,
		SiteDescription: req.SiteDescription,
		SiteRssUrl:      req.SiteRssURL,
		HasAdv:          req.HasAdv,
		DesiredLocation: getLocation.Id,
		DesiredColor:    getColor.Id,
		WebmasterRemark: req.Remark,
	}); err != nil {
		glog.Infof(ctx, "[LOGIC] 添加链接失败[%s]", err.Error())
		return errors.New("添加链接失败")
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
func (s *sLinkLogic) GetColor(ctx context.Context) (getColors []*entity.XfColor, err error) {
	glog.Info(ctx, "[LOGIC] 执行 LinkLogic:GetColor 服务层")
	err = dao.XfColor.Ctx(ctx).Where(do.XfColor{HasSelect: true}).Scan(&getColors)
	if err != nil {
		glog.Error(ctx, "[LOGIC] 数据库错误<获取期望颜色失败>")
		return nil, err
	}
	return getColors, nil
}

// GetLocation 获取期望位置信息
// 用于获取期望位置信息, 如果成功则返回期望位置信息，否则返回错误。
//
// 参数：
// ctx: 请求的上下文，用于管理超时和取消信号。
//
// 返回：
// getLocation: 如果获取期望位置信息成功，返回期望位置信息；否则返回错误。
// err: 如果获取期望位置信息成功，返回 nil；否则返回错误。其中错误的返回信息在此函数中主要包含内容为数据库错误。
func (s *sLinkLogic) GetLocation(ctx context.Context) (getLocation []*entity.XfLocation, err error) {
	glog.Info(ctx, "[LOGIC] 执行 LinkLogic:GetLocation 服务层")
	err = dao.XfLocation.Ctx(ctx).Where(do.XfLocation{Reveal: true}).Scan(&getLocation)
	if err != nil {
		glog.Error(ctx, "[LOGIC] 数据库错误<获取期望位置失败>")
		return nil, err
	}
	return getLocation, nil
}

// GetLocationAllInformation 获取所有的位置信息
// 用于获取所有的位置信息, 如果成功则返回期望位置信息，否则返回错误。
//
// 参数：
// ctx: 请求的上下文，用于管理超时和取消信号。
//
// 返回：
// getLocation: 如果获取期望位置信息成功，返回期望位置信息；否则返回错误。
// err: 如果获取期望位置信息成功，返回 nil；否则返回错误。其中错误的返回信息在此函数中主要包含内容为数据库错误。
func (s *sLinkLogic) GetLocationAllInformation(ctx context.Context) (getLocation []*entity.XfLocation, err error) {
	glog.Info(ctx, "[LOGIC] 执行 LinkLogic:GetLocationAllInformation 服务层")
	err = dao.XfLocation.Ctx(ctx).Scan(&getLocation)
	if err != nil {
		glog.Error(ctx, "[LOGIC] 数据库错误<获取期望位置失败>")
		return nil, err
	}
	return getLocation, nil
}

func (s *sLinkLogic) AddLocation(
	ctx context.Context,
	name string,
	displayName string,
	description string,
	reveal bool,
	sort int,
) (err error) {
	glog.Info(ctx, "[LOGIC] 执行 LinkLogic:AddLocation 服务层")
	if _, err = dao.XfLocation.Ctx(ctx).Insert(do.XfLocation{
		Name:        name,
		DisplayName: displayName,
		Description: description,
		Reveal:      reveal,
		Sort:        sort,
	}); err != nil {
		glog.Infof(ctx, "[LOGIC] 添加链接位置失败[%s]", err.Error())
		return errors.New("[LOGIC] 数据库错误<添加链接位置失败>")
	}
	return nil
}
