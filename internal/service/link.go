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

// ================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// You can delete these comments if you wish manually maintain this interface file.
// ================================================================================

package service

import (
	"context"
	v1 "xiaoMain/api/link/v1"
	"xiaoMain/internal/model/entity"
)

type (
	ILinkLogic interface {
		// CheckLinkCanAccess 检查链接是否可以访问
		// 用于检查用户添加的链接是否可以访问，如果可以访问则返回 nil，否则返回错误
		//
		// 参数：
		// ctx: 请求的上下文，用于管理超时和取消信号。
		// siteURL: 用户尝试添加的链接地址。
		//
		// 返回：
		// err: 如果链接可以访问，返回 nil；否则返回错误。
		CheckLinkCanAccess(ctx context.Context, siteURL string) (err error)
		// CheckLogoCanAccess 检查 Logo 是否可以访问
		// 用于检查用户添加的 Logo 是否可以访问，如果可以访问则返回 nil，否则返回错误
		// 并且检查获取的状态是否是图片，如果不是图片则返回错误，否则返回 nil
		//
		// 参数：
		// ctx: 请求的上下文，用于管理超时和取消信号。
		// siteLogo: 用户尝试添加的 Logo 地址。
		//
		// 返回：
		// err: 如果 Logo 可以访问，返回 nil；否则返回错误。
		CheckLogoCanAccess(ctx context.Context, siteLogo string) (err error)
		// CheckRSSURL 检查链接是否可以访问
		// 用于检查用户添加的链接是否可以访问，如果可以访问则返回 nil，否则返回错误
		// 还需要检查是否为 RSS URL 是否为 XML 格式
		//
		// 参数：
		// ctx: 请求的上下文，用于管理超时和取消信号。
		// siteRSS: 用户尝试添加的 RSS 地址。
		//
		// 返回：
		// err: 如果链接可以访问，返回 nil；否则返回错误。
		CheckRSSCanAccess(ctx context.Context, siteRSS string) (err error)
		// CheckLinkName 检查链接名是否重复
		// 用于检查用户添加的链接名是否已经存在，如果存在则返回错误，否则返回成功
		//
		// 参数：
		// ctx: 请求的上下文，用于管理超时和取消信号。
		// linkName: 用户尝试添加的链接名。
		//
		// 返回：
		// err: 如果链接名已存在，返回错误；否则返回 nil。
		CheckLinkName(ctx context.Context, linkName string) (err error)
		// CheckLinkURL 检查链接URL是否重复
		// 用于检查用户添加的链接URL是否已经存在，如果存在则返回错误，否则返回成功
		//
		// 参数：
		// ctx: 请求的上下文，用于管理超时和取消信号。
		// siteURL: 用户尝试添加的链接URL。
		//
		// 返回：
		// err: 如果链接URL已存在，返回错误；否则返回 nil。
		CheckLinkURL(ctx context.Context, siteURL string) (err error)
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
		CheckLinkHasConnect(ctx context.Context, linkID string) (delay *int64, err error)
		// AddLink 添加链接
		// 用于添加链接，如果添加成功则返回 nil，否则返回错误。
		//
		// 参数：
		// ctx: 请求的上下文，用于管理超时和取消信号。
		// req: 用户的请求，包含添加链接的详细信息。
		//
		// 返回：
		// err: 如果添加链接成功，返回 nil；否则返回错误。
		AddLink(ctx context.Context, req v1.LinkAddReq) (err error)
		// GetColor 获取期望颜色信息
		// 用于获取期望颜色信息, 如果成功则返回期望颜色信息，否则返回错误。
		//
		// 参数：
		// ctx: 请求的上下文，用于管理超时和取消信号。
		//
		// 返回：
		// getColors: 如果获取期望颜色信息成功，返回期望颜色信息；否则返回错误。
		// err: 如果获取期望颜色信息成功，返回 nil；否则返回错误。
		GetColor(ctx context.Context) (getColors []*entity.XfColor, err error)
		// GetLocation 获取期望位置信息
		// 用于获取期望位置信息, 如果成功则返回期望位置信息，否则返回错误。
		//
		// 参数：
		// ctx: 请求的上下文，用于管理超时和取消信号。
		//
		// 返回：
		// getLocation: 如果获取期望位置信息成功，返回期望位置信息；否则返回错误。
		// err: 如果获取期望位置信息成功，返回 nil；否则返回错误。
		GetLocation(ctx context.Context) (getLocation []*entity.XfLocation, err error)
	}
)

var (
	localLinkLogic ILinkLogic
)

func LinkLogic() ILinkLogic {
	if localLinkLogic == nil {
		panic("implement not found for interface ILinkLogic, forgot register?")
	}
	return localLinkLogic
}

func RegisterLinkLogic(i ILinkLogic) {
	localLinkLogic = i
}
