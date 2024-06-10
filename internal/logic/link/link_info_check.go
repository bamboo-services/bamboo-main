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
	"github.com/gogf/gf/v2/os/gtime"
	"xiaoMain/internal/dao"
	"xiaoMain/internal/model/do"
	"xiaoMain/internal/model/entity"
)

// CheckLinkName
//
// # 检查地址信息
//
// 用于检查用户添加的链接名是否已经存在，如果存在则返回错误，否则返回成功
//
// # 参数:
//   - ctx: 上下文对象，用于传递和控制请求的生命周期。
//   - linkName: 用户尝试添加的链接名。
//
// # 返回:
//   - err: 如果链接名已存在，返回错误；否则返回 nil.
func (s *sLink) CheckLinkName(ctx context.Context, linkName string) (err error) {
	g.Log().Notice(ctx, "[LOGIC] Link:CheckLinkName | 检查地址信息")
	// 检查链接名是否重复
	var linkInfo *do.LinkList
	err = dao.LinkList.Ctx(ctx).Where(do.LinkList{SiteName: linkName}).Scan(&linkInfo)
	if err != nil {
		g.Log().Errorf(ctx, "[LOGIC] 数据库查询错误，错误原因： %s", err.Error())
		return berror.NewErrorHasError(bcode.ServerInternalError, err)
	}
	if linkInfo != nil {
		g.Log().Errorf(ctx, "[LOGIC] 网站名已存在，网站名：%s", linkName)
		return berror.NewError(bcode.AlreadyExists, "网站名已存在")
	} else {
		return nil
	}
}

// CheckLinkURL
//
// # 检查链接地址
//
// 用于检查用户添加的链接地址是否已经存在，如果存在则返回错误，否则返回成功
//
// # 参数:
//   - ctx: 上下文对象，用于传递和控制请求的生命周期。
//   - siteURL: 用户尝试添加的链接地址。
//
// # 返回:
//   - err: 如果链接地址已存在，返回错误；否则返回 nil.
func (s *sLink) CheckLinkURL(ctx context.Context, siteURL string) (err error) {
	g.Log().Notice(ctx, "[LOGIC] Link:CheckLinkUrl | 检查链接地址")
	// 检查链接名是否重复
	var linkInfo *do.LinkList
	err = dao.LinkList.Ctx(ctx).Where(do.LinkList{SiteUrl: siteURL}).Scan(&linkInfo)
	if err != nil {
		g.Log().Errorf(ctx, "[LOGIC] 数据库查询错误，错误原因： %s", err.Error())
		return berror.NewErrorHasError(bcode.ServerInternalError, err)
	}
	if linkInfo != nil {
		g.Log().Errorf(ctx, "[LOGIC] 网站链接已存在，网站链接：%s", siteURL)
		return berror.NewError(bcode.AlreadyExists, "网站链接已存在")
	} else {
		return nil
	}
}

// CheckLinkHasConnect
//
// # 检查链接是否可以连接
//
// 用于检查用户添加的链接地址是否可以连接，如果可以则返回 nil，否则返回错误
//
// # 参数:
//   - ctx: 上下文对象，用于传递和控制请求的生命周期。
//   - linkID: 用户尝试添加的链接ID。
//
// # 返回:
//   - delay: 如果链接可以连接，返回延迟时间；否则返回错误.
//   - err: 如果链接不存在，返回错误；否则返回 nil.
func (s *sLink) CheckLinkHasConnect(ctx context.Context, linkID string) (delay *int64, err error) {
	g.Log().Notice(ctx, "[LOGIC] Link:CheckLinkHasConnect | 检查链接是否可以连接")
	// 获取链接信息
	var getLink *entity.LinkList
	err = dao.LinkList.Ctx(ctx).Where(do.LinkList{Id: linkID}).Scan(&getLink)
	if err != nil {
		g.Log().Errorf(ctx, "[LOGIC] 数据库查询错误，错误原因： %s", err.Error())
		return nil, berror.NewErrorHasError(bcode.ServerInternalError, err)
	}
	if getLink == nil {
		g.Log().Warningf(ctx, "[LOGIC] 链接不存在，链接ID[%s]", linkID)
		return nil, berror.NewError(bcode.NotExist, "链接不存在")
	}
	// 检查链接是否可以连接
	getNowTimestamp := gtime.Now().TimestampMicro()
	err = s.CheckLinkCanAccess(ctx, getLink.SiteUrl)
	if err != nil {
		return nil, err
	}
	// 获取延迟信息
	*delay = gtime.Now().TimestampMicro() - getNowTimestamp
	return delay, nil
}

// IsColorExistByName
//
// # 获取颜色信息
//
// 用于获取颜色信息，如果成功则返回颜色信息，否则返回错误。
//
// # 参数:
//   - ctx: 请求的上下文，用于管理超时和取消信号。
//   - getColor: 用户尝试获取的颜色名称。
//
// # 返回:
//   - err: 如果颜色存在，返回错误；否则返回 nil.
func (s *sLink) IsColorExistByName(ctx context.Context, getColor string) error {
	g.Log().Notice(ctx, "[LOGIC] Link:IsColorExistByName | 获取颜色信息")
	// 查询指定的颜色是否存在
	var getColorInfo *entity.Color
	err := dao.Color.Ctx(ctx).Where(do.Color{Name: getColor}).Scan(&getColorInfo)
	if err != nil {
		g.Log().Errorf(ctx, "[LOGIC] 数据库查询错误，错误原因： %s", err.Error())
		return berror.NewErrorHasError(bcode.ServerInternalError, err)
	}
	if getColorInfo != nil {
		return nil
	} else {
		return berror.NewError(bcode.NotExist, "颜色不存在")
	}
}

// IsColorExistByColorID
//
// # 获取颜色信息
//
// 用于获取颜色信息，如果成功则返回颜色信息，否则返回错误。
//
// # 参数:
//   - ctx: 请求的上下文，用于管理超时和取消信号。
//   - getColor: 用户尝试获取的颜色名称。
//
// # 返回:
//   - err: 如果颜色存在，返回错误；否则返回 nil.
func (s *sLink) IsColorExistByColorID(ctx context.Context, getColor string) error {
	g.Log().Notice(ctx, "[LOGIC] Link:IsColorExistByColorID | 获取颜色信息")
	// 查询指定的颜色是否存在
	var getColorInfo *entity.Color
	err := dao.Color.Ctx(ctx).Where(do.Color{Color: getColor}).Scan(&getColorInfo)
	if err != nil {
		g.Log().Errorf(ctx, "[LOGIC] 数据库查询错误，错误原因： %s", err.Error())
		return berror.NewErrorHasError(bcode.ServerInternalError, err)
	}
	if getColorInfo != nil {
		return nil
	} else {
		return berror.NewError(bcode.NotExist, "颜色不存在")
	}
}

// IsLocationExist
//
// # 获取位置信息
//
// 用于获取位置信息，如果成功则返回位置信息，否则返回错误。
//
// # 参数:
//   - ctx: 请求的上下文，用于管理超时和取消信号。
//   - name: 用户尝试获取的位置名称。
//
// # 返回:
//   - err: 如果位置存在，返回错误；否则返回 nil.
func (s *sLink) IsLocationExist(ctx context.Context, name string) (err error) {
	g.Log().Notice(ctx, "[LOGIC] 执行 Link:IsLocationExist 服务层")
	// 查询指定的位置是否存在
	var getLocationInfo *entity.Location
	err = dao.Location.Ctx(ctx).Where(do.Location{Name: name}).Scan(&getLocationInfo)
	if err != nil {
		g.Log().Errorf(ctx, "[LOGIC] 数据库查询错误，错误原因： %s", err.Error())
		return berror.NewErrorHasError(bcode.ServerInternalError, err)
	}
	if getLocationInfo != nil {
		return nil
	} else {
		return berror.NewError(bcode.NotExist, "位置不存在")
	}
}
