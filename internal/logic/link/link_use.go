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
)

// AddLink 添加链接
// 用于添加链接，如果添加成功则返回 nil，否则返回错误。
//
// 参数：
// ctx: 请求的上下文，用于管理超时和取消信号。
// req: 用户的请求，包含添加链接的详细信息。
//
// 返回：
// err: 如果添加链接成功，返回 nil；否则返回错误。
func (s *sLinkLogic) AddLink(ctx context.Context, req v1.LinkAddReq) (err error) {
	glog.Info(ctx, "[LOGIC] 执行 LinkLogic:AddLink 服务层")
	var getLocation *do.XfDesiredLocation
	var getColor *do.XfDesiredColor
	// 对颜色进行数据库获取，获取指定的颜色是否存在
	if dao.XfDesiredLocation.Ctx(ctx).Where(do.XfDesiredLocation{Name: req.DesiredLocation}).Scan(&getLocation) != nil {
		glog.Info(ctx, "[LOGIC] 数据库错误<期望位置不存在>")
		return errors.New("期望位置不存在")
	}
	if getLocation == nil {
		glog.Infof(ctx, "[LOGIC] 期望位置不存在[%s]", req.DesiredLocation)
		return errors.New("期望位置不存在")
	}
	// 对颜色进行数据库获取，获取指定的颜色是否存在
	if dao.XfDesiredColor.Ctx(ctx).Where(do.XfDesiredColor{Name: req.DesiredColor}).Scan(&getColor) != nil {
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
