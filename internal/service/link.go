// ================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// You can delete these comments if you wish manually maintain this interface file.
// ================================================================================

package service

import (
	"context"
	v1 "xiaoMain/api/link/v1"
	"xiaoMain/internal/model/dto/flow"
	"xiaoMain/internal/model/entity"
)

type (
	ILink interface {
		// CheckLinkCanAccess
		//
		// # 检查链接是否可以访问
		//
		// 用于检查用户添加的链接是否可以访问，如果可以访问则返回 nil，否则返回错误。
		//
		// # 参数:
		//   - ctx: 上下文对象，用于传递和控制请求的生命周期。
		//   - siteURL: 网站地址(string)
		//
		// # 返回:
		//   - err: 如果链接可以访问，返回 nil。否则返回错误信息。
		CheckLinkCanAccess(ctx context.Context, siteURL string) (err error)
		// CheckLogoCanAccess
		//
		// # 检查 Logo 是否可以访问
		//
		// 用于检查用户添加的 Logo 是否可以访问，如果可以访问则返回 nil，否则返回错误。
		//
		// # 参数:
		//   - ctx: 上下文对象，用于传递和控制请求的生命周期。
		//   - siteLogo: 网站 Logo 地址(string)
		//
		// # 返回:
		//   - err: 如果 Logo 可以访问，返回 nil。否则返回错误信息。
		CheckLogoCanAccess(ctx context.Context, siteLogo string) (err error)
		// CheckRSSCanAccess
		//
		// # 检查 RSS 是否可以访问
		//
		// 用于检查用户添加的 RSS 是否可以访问，如果可以访问则返回 nil，否则返回错误。
		//
		// # 参数:
		//   - ctx: 上下文对象，用于传递和控制请求的生命周期。
		//   - siteRSS: 网站 RSS 地址(string)
		//
		// # 返回:
		//   - err: 如果 RSS 可以访问，返回 nil。否则返回错误信息。
		CheckRSSCanAccess(ctx context.Context, siteRSS string) (err error)
		// GetLink
		//
		// # 获取链接
		//
		// 获取链接，包括友情链接、社交链接等。用于获取链接之后进行展示使用
		//
		// # 参数:
		//   - ctx: 上下文对象，用于传递和控制请求的生命周期。
		//
		// # 返回:
		//   - getLink: 获取到的链接列表
		//   - total: 获取到的链接总数
		//   - err: 如果处理过程中发生错误，返回错误信息。
		GetLink(ctx context.Context) (getLink []flow.LinkGetDTO, total uint64, err error)
		// GetLinkAdmin
		//
		// # 获取链接
		//
		// 获取链接，包括友情链接、社交链接等。用于获取链接之后进行展示使用。该链接获只有管理员可以对内容进行获取。
		//
		// # 参数:
		//   - ctx: 上下文对象，用于传递和控制请求的生命周期。
		//
		// # 返回:
		//   - getAllLink: 获取到的链接列表
		//   - total: 获取到的链接总数
		//   - err: 如果处理过程中发生错误，返回错误信息。
		GetLinkAdmin(ctx context.Context) (getAllLink []entity.LinkList, total uint64, err error)
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
		CheckLinkName(ctx context.Context, linkName string) (err error)
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
		CheckLinkURL(ctx context.Context, siteURL string) (err error)
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
		IsColorExistByName(ctx context.Context, getColor string) error
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
		IsColorExistByColorID(ctx context.Context, getColor string) error
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
		IsLocationExist(ctx context.Context, name string) (err error)
		EditLocation(ctx context.Context, req v1.EditLocationReq) (err error)
		LocationMove(ctx context.Context, id, moveID int64) (err error)
		DelLocation(ctx context.Context, id int64) (err error)
		// GetLinkByID
		//
		// # 通过 ID 获取链接
		//
		// 通过 ID 获取链接, 需要用户提供链接的 ID。
		//
		// # 参数
		//   - ctx: 请求的上下文，用于管理超时和取消信号。
		//   - getLinkID: 用户的请求，包含获取链接的详细信息。
		//
		// # 返回
		//   - err: 在获取链接过程中发生的任何错误。
		GetLinkByID(ctx context.Context, getLinkID int64) (link *entity.LinkList, err error)
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
		AddLink(ctx context.Context, req v1.LinkAddReq) (err error)
		// GetColor 获取期望颜色信息
		// 用于获取期望颜色信息, 如果成功则返回期望颜色信息，否则返回错误。
		//
		// 参数：
		// ctx: 请求的上下文，用于管理超时和取消信号。
		//
		// 返回：
		// getColors: 如果获取期望颜色信息成功，返回期望颜色信息；否则返回错误。
		// err: 如果获取期望颜色信息成功，返回 nil；否则返回错误。其中错误的返回信息在此函数中主要包含内容为数据库错误。
		GetColor(ctx context.Context) (getColors []*entity.Color, err error)
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
		GetLocation(ctx context.Context) (getLocation []*entity.Location, err error)
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
		GetLocationNoReveal(ctx context.Context) (getLocation []*entity.Location, err error)
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
		AddLocation(ctx context.Context, name string, displayName string, description string, reveal bool, sort int) (err error)
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
		AddColor(ctx context.Context, name string, displayName string, color string, hasSelect bool) (err error)
	}
)

var (
	localLink ILink
)

func Link() ILink {
	if localLink == nil {
		panic("implement not found for interface ILink, forgot register?")
	}
	return localLink
}

func RegisterLink(i ILink) {
	localLink = i
}
