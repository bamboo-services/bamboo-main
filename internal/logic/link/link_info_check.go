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
	"github.com/gogf/gf/v2/os/gtime"
	"xiaoMain/internal/dao"
	"xiaoMain/internal/model/do"
	"xiaoMain/internal/model/entity"
)

// CheckLinkName 检查链接名是否重复
// 用于检查用户添加的链接名是否已经存在，如果存在则返回错误，否则返回成功
//
// 参数：
// ctx: 请求的上下文，用于管理超时和取消信号。
// linkName: 用户尝试添加的链接名。
//
// 返回：
// err: 如果链接名已存在，返回错误；否则返回 nil。
func (s *sLinkLogic) CheckLinkName(ctx context.Context, linkName string) (err error) {
	glog.Info(ctx, "[LOGIC] 执行 LinkLogic:CheckLinkName 服务层")
	// 检查链接名是否重复
	var linkInfo *do.XfLinkList
	err = dao.XfLinkList.Ctx(ctx).Where(do.XfLinkList{SiteName: linkName}).Scan(&linkInfo)
	if err != nil {
		glog.Errorf(ctx, "[LOGIC] 数据库查询错误，错误原因： %s", err.Error())
		return errors.New("数据库查询错误")
	}
	if linkInfo != nil {
		glog.Errorf(ctx, "[LOGIC] 网站名已存在，网站名：%s", linkName)
		return errors.New("网站名已存在")
	} else {
		return nil
	}
}

// CheckLinkURL 检查链接URL是否重复
// 用于检查用户添加的链接URL是否已经存在，如果存在则返回错误，否则返回成功
//
// 参数：
// ctx: 请求的上下文，用于管理超时和取消信号。
// siteURL: 用户尝试添加的链接URL。
//
// 返回：
// err: 如果链接URL已存在，返回错误；否则返回 nil。
func (s *sLinkLogic) CheckLinkURL(ctx context.Context, siteURL string) (err error) {
	glog.Info(ctx, "[LOGIC] 执行 LinkLogic:CheckLinkUrl 服务层")
	// 检查链接名是否重复
	var linkInfo *do.XfLinkList
	err = dao.XfLinkList.Ctx(ctx).Where(do.XfLinkList{SiteUrl: siteURL}).Scan(&linkInfo)
	if err != nil {
		glog.Errorf(ctx, "[LOGIC] 数据库查询错误，错误原因： %s", err.Error())
		return errors.New("数据库查询错误")
	}
	if linkInfo != nil {
		glog.Errorf(ctx, "[LOGIC] 网站链接已存在，网站链接：%s", siteURL)
		return errors.New("网站链接已存在")
	} else {
		return nil
	}
}

// CheckLinkHasConnect 检查链接是否已经连接
// 用于检查链接是否已经连接，如果成功则返回 nil，否则返回错误。
// 本接口会根据已有的链接信息对链接进行链接检查是否可以连接，若连接失败返回失败信息，若成功返回成功信息
//
// 参数：
// ctx: 请求的上下文，用于管理超时和取消信号。
// linkID: 用户尝试添加的链接ID。
//
// 返回：
// err: 如果链接已连接，返回错误；否则返回 nil。
func (s *sLinkLogic) CheckLinkHasConnect(ctx context.Context, linkID string) (delay *int64, err error) {
	glog.Info(ctx, "[LOGIC] 执行 LinkLogic:CheckLinkHasConnect 服务层")
	// 获取链接信息
	var getLink *entity.XfLinkList
	err = dao.XfLinkList.Ctx(ctx).Where(do.XfLinkList{Id: linkID}).Scan(&getLink)
	if err != nil {
		glog.Errorf(ctx, "[LOGIC] 数据库查询错误，错误原因： %s", err.Error())
		return nil, errors.New("数据库查询错误")
	}
	if getLink == nil {
		glog.Warningf(ctx, "[LOGIC] 链接不存在，链接ID[%s]", linkID)
		return nil, errors.New("链接不存在")
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

// HasColorByName 获取颜色信息
// 用于获取颜色信息，如果成功则返回颜色信息，否则返回错误。
// 本接口会根据已有的颜色信息对颜色进行查询，若查询失败返回失败信息，若成功返回成功信息
//
// 参数：
// ctx: 请求的上下文，用于管理超时和取消信号。
// getColor: 用户尝试获取的颜色名称。
//
// 返回：
// bool: 如果颜色存在，返回 false；否则返回 true。
func (s *sLinkLogic) HasColorByName(ctx context.Context, getColor string) bool {
	glog.Info(ctx, "[LOGIC] 执行 LinkLogic:HasColorByName 服务层")
	// 查询指定的颜色是否存在
	var getColorInfo *entity.XfColor
	err := dao.XfColor.Ctx(ctx).Where(do.XfColor{Name: getColor}).Scan(&getColorInfo)
	if err != nil {
		glog.Errorf(ctx, "[LOGIC] 数据库查询错误，错误原因： %s", err.Error())
		return false
	}
	if getColorInfo != nil {
		glog.Errorf(ctx, "[LOGIC] 颜色已存在，颜色：%s", getColor)
		return false
	} else {
		return true
	}
}

// HasColorByColor 获取颜色信息
// 用于获取颜色信息，如果成功则返回颜色信息，否则返回错误。
// 本接口会根据已有的颜色信息对颜色进行查询，若查询失败返回失败信息，若成功返回成功信息
//
// 参数：
// ctx: 请求的上下文，用于管理超时和取消信号。
// getColor: 用户尝试获取的颜色名称。
//
// 返回：
// bool: 如果颜色存在，返回 false；否则返回 true。
func (s *sLinkLogic) HasColorByColor(ctx context.Context, getColor string) bool {
	glog.Info(ctx, "[LOGIC] 执行 LinkLogic:HasColorByColor 服务层")
	// 查询指定的颜色是否存在
	var getColorInfo *entity.XfColor
	err := dao.XfColor.Ctx(ctx).Where(do.XfColor{Color: getColor}).Scan(&getColorInfo)
	if err != nil {
		glog.Errorf(ctx, "[LOGIC] 数据库查询错误，错误原因： %s", err.Error())
		return false
	}
	if getColorInfo != nil {
		glog.Errorf(ctx, "[LOGIC] 颜色已存在，颜色：%s", getColor)
		return false
	} else {
		return true
	}
}

// CheckLocationExist 检查位置是否存在
// 用于检查位置是否存在，如果成功则返回 nil，否则返回错误。
// 本接口会根据已有的位置信息对位置进行查询，若查询失败返回失败信息，若成功返回成功信息
//
// 参数：
// ctx: 请求的上下文，用于管理超时和取消信号。
// name: 用户尝试添加的位置名称。
//
// 返回：
// err: 如果位置存在，返回错误；否则返回 nil。
func (s *sLinkLogic) CheckLocationExist(ctx context.Context, name string) (err error) {
	glog.Info(ctx, "[LOGIC] 执行 LinkLogic:CheckLocationExist 服务层")
	// 查询指定的位置是否存在
	var getLocationInfo *entity.XfLocation
	err = dao.XfLocation.Ctx(ctx).Where(do.XfLocation{Name: name}).Scan(&getLocationInfo)
	if err != nil {
		glog.Errorf(ctx, "[LOGIC] 数据库查询错误，错误原因： %s", err.Error())
		return errors.New("数据库查询错误")
	}
	if getLocationInfo != nil {
		glog.Errorf(ctx, "[LOGIC] 位置已存在，位置：%s", name)
		return errors.New("位置已存在")
	} else {
		return nil
	}
}
