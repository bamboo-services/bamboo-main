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
	"xiaoMain/internal/model/dto/dmiddle"
)

type (
	IRss interface {
		// GetAllLinkRssInfo
		//
		// # 获取所有链接的Rss信息
		//
		// 用于获取所有链接的Rss信息，如果成功则返回 nil，否则返回错误。
		//
		// # 参数:
		//   - ctx: 上下文对象，用于传递和控制请求的生命周期。
		//
		// # 返回:
		//   - getRss: 如果获取Rss信息成功，返回 nil；否则返回错误。
		//   - err: 如果获取Rss信息成功，返回 nil；否则返回错误。
		GetAllLinkRssInfo(ctx context.Context) (getRss *[]*dmiddle.RssLinkDTO, err error)
		// GetLinkRssInfoWithLinkID
		//
		// # 获取链接的Rss信息
		//
		// 用于获取链接的Rss信息，如果成功则返回 nil，否则返回错误。
		// 本接口会根据已有的链接信息对Rss信息进行获取，若获取失败返回失败信息，若成功返回成功信息。
		//
		// # 参数:
		//   - ctx: 请求的上下文，用于管理超时和取消信号。
		//   - linkID: 链接的ID。
		//
		// # 返回:
		//   - getRss: 如果获取Rss信息成功，返回 nil；否则返回错误。
		//   - err: 如果获取Rss信息成功，返回 nil；否则返回错误。
		GetLinkRssInfoWithLinkID(ctx context.Context, linkID int64) (getRss *[]*dmiddle.RssLinkDTO, err error)
		// GetLinkRssWithLinkName
		//
		// # 获取链接的Rss信息
		//
		// 用于获取链接的Rss信息，如果成功则返回 nil，否则返回错误。
		// 本接口会根据已有的链接信息对Rss信息进行获取，若获取失败返回失败信息，若成功返回成功信息。
		//
		// # 参数:
		//   - ctx: 请求的上下文，用于管理超时和取消信号。
		//   - linkName: 链接的名称。
		//
		// # 返回:
		//   - getRss: 如果获取Rss信息成功，返回 nil；否则返回错误。
		//   - err: 如果获取Rss信息成功，返回 nil；否则返回错误。
		GetLinkRssWithLinkName(ctx context.Context, linkName string) (getRss *[]*dmiddle.RssLinkDTO, err error)
		// GetLinkRssWithLinkLocation
		//
		// # 获取链接的Rss信息
		//
		// 用于获取链接的Rss信息，如果成功则返回 nil，否则返回错误。
		// 本接口会根据已有的链接信息对Rss信息进行获取，若获取失败返回失败信息，若成功返回成功信息。
		//
		// # 参数:
		//   - ctx: 请求的上下文，用于管理超时和取消信号。
		//   - linkLocation: 链接的位置。
		//
		// # 返回:
		//   - getRss: 如果获取Rss信息成功，返回 nil；否则返回错误。
		//   - err: 如果获取Rss信息成功，返回 nil；否则返回错误。
		GetLinkRssWithLinkLocation(ctx context.Context, linkLocation int64) (getRss *[]*dmiddle.RssLinkDTO, err error)
		// RssWithHexoFeed 通过Hexo的Rss信息获取Rss信息
		// 用于获取Hexo的Rss内容（插件：hexo-generator-feed ｜ 插件内容 https://github.com/hexojs/hexo-generator-feed）
		// 如果成功则返回 RssLinkDTO，否则返回 false
		//
		// 参数：
		// ctx: 请求的上下文，用于管理超时和取消信号。
		// rssURL: 站点的 rss 订阅地址。
		//
		// 返回：
		// rssLink: 如果获取Rss信息成功，返回 RssLinkDTO；否则返回 nil。
		// hasThis: 如果获取Rss信息成功，返回 true；否则返回 false。
		RssWithHexoFeed(ctx context.Context, rssURL string) (rssLink *[]dmiddle.RssLinkDTO, hasThis bool)
		// RssWithHugoFeed 通过Hugo的Rss信息获取Rss信息
		// 用于获取Hugo的Rss内容（插件：Hugo ｜ 插件内容 https://gohugo.io/）
		// 如果成功则返回 nil，否则返回错误
		//
		// 参数：
		// ctx: 请求的上下文，用于管理超时和取消信号。
		// rssURL: 站点的 rss 订阅地址。
		//
		// 返回：
		// rssLink: 如果获取Rss信息成功，返回 RssLinkDTO；否则返回 nil。
		// hasThis: 如果获取Rss信息成功，返回 true；否则返回 false。
		RssWithHugoFeed(ctx context.Context, rssURL string) (rssLink *[]dmiddle.RssLinkDTO, hasThis bool)
		// RssWithWordpressFeed 通过WordPress的Rss信息获取Rss信息
		// 用于获取WordPress的Rss内容（插件：WordPress ｜ 插件内容 https://wordpress.org/）
		// 如果成功则返回 nil，否则返回错误
		//
		// 参数：
		// ctx: 请求的上下文，用于管理超时和取消信号。
		// rssURL: 站点的 rss 订阅地址。
		//
		// 返回：
		// rssLink: 如果获取Rss信息成功，返回 RssLinkDTO；否则返回 nil。
		// hasThis: 如果获取Rss信息成功，返回 true；否则返回 false。
		RssWithWordpressFeed(ctx context.Context, rssURL string) (rssLink *[]dmiddle.RssLinkDTO, hasThis bool)
	}
)

var (
	localRss IRss
)

func Rss() IRss {
	if localRss == nil {
		panic("implement not found for interface IRss, forgot register?")
	}
	return localRss
}

func RegisterRss(i IRss) {
	localRss = i
}
