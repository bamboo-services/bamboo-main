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

package rss

import (
	"context"
	"encoding/xml"
	"github.com/gogf/gf/v2/os/glog"
	"io"
	"xiaoMain/internal/lutil"
	"xiaoMain/internal/model/dto"
)

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
func (s *sRssLogic) RssWithHexoFeed(ctx context.Context, rssURL string) (rssLink *[]dto.RssLinkDTO, hasThis bool) {
	glog.Noticef(ctx, "[LOGIC] 尝试获取 Hexo 中 hexo-generator-feed 插件的 feed 内容")
	getResp, err := lutil.LinkAccess(rssURL)
	defer func() {
		if getResp != nil {
			_ = getResp.Body.Close()
		}
	}()
	// 处理错误检查信息
	if err != nil {
		glog.Warningf(ctx, "获取链接信息失败: %v", err)
	}
	// 获取 Body 并解析 XML
	getBody, err := io.ReadAll(getResp.Body)
	if err != nil {
		glog.Warningf(ctx, "获取链接信息失败: %v", err)
	}
	// 解析 XML
	getFeed := new(dto.FeedDTO)
	err = xml.Unmarshal(getBody, &getFeed)
	if err != nil {
		glog.Warningf(ctx, "解析 XML 失败: %v", err.Error())
		return nil, false
	}
	// 处理数据
	if getFeed.Generator == "Hexo" {
		if getFeed != nil {
			rssLink = new([]dto.RssLinkDTO)
			for _, item := range getFeed.Entry {
				// 处理 Category
				categories := new([]string)
				for _, category := range item.Category {
					*categories = append(*categories, category.Term)
				}
				// 添加数据
				*rssLink = append(*rssLink, dto.RssLinkDTO{
					Title:    item.Title,
					Link:     item.ID,
					Summary:  item.Summary,
					Category: *categories,
					Timer:    item.Published,
				})
			}
			return rssLink, true
		} else {
			glog.Warning(ctx, "解析 XML 成功，但数据为空")
			return nil, false
		}
	} else {
		return nil, false
	}
}
