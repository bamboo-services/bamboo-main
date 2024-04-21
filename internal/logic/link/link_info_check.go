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
	"xiaoMain/internal/dao"
	"xiaoMain/internal/model/do"
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
