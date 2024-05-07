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
	"errors"
	"github.com/gogf/gf/v2/os/glog"
	"github.com/gogf/gf/v2/os/gtime"
	"io"
	"regexp"
	time2 "time"
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
	getBody, err := s.rssLinkAccess(ctx, rssURL)
	if err != nil {
		return nil, false
	}
	// 解析 XML
	getFeed := new(dto.HexoFeedDTO)
	err = xml.Unmarshal(getBody, &getFeed)
	if err != nil {
		glog.Warningf(ctx, "解析 XML 失败: %v", err.Error())
		return nil, false
	}
	// 处理数据
	if getFeed != nil {
		if getFeed.Generator == "Hexo" {
			rssLink = new([]dto.RssLinkDTO)
			for _, item := range getFeed.Entry {
				// 处理 Category
				categories := new([]string)
				for _, category := range item.Category {
					*categories = append(*categories, category.Term)
				}
				description := item.Summary
				if len(description) > 100 {
					description = description[:100]
				}
				// 添加数据
				*rssLink = append(*rssLink, dto.RssLinkDTO{
					Title:    item.Title,
					Link:     item.ID,
					Summary:  description,
					Category: *categories,
					Timer:    gtime.NewFromStr(item.Published).TimestampMilliStr(),
				})
			}
			return rssLink, true
		}
	} else {
		glog.Warning(ctx, "解析 XML 成功，但数据为空")
	}
	return nil, false
}

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
func (s *sRssLogic) RssWithHugoFeed(ctx context.Context, rssURL string) (rssLink *[]dto.RssLinkDTO, hasThis bool) {
	glog.Noticef(ctx, "[LOGIC] 尝试获取 Hugo 中 原生 的 feed 内容")
	getBody, err := s.rssLinkAccess(ctx, rssURL)
	if err != nil {
		return nil, false
	}
	// 解析 XML
	getFeed := new(dto.HugoFeedDTO)
	err = xml.Unmarshal(getBody, &getFeed)
	if err != nil {
		glog.Warningf(ctx, "解析 XML 失败: %v", err.Error())
		return nil, false
	}
	// 处理数据
	if getFeed != nil {
		if getFeed.Generator == "Hugo -- gohugo.io" {
			rssLink = new([]dto.RssLinkDTO)
			for _, item := range getFeed.Items {
				description := item.Description
				if len(description) > 100 {
					description = description[:100]
				}
				// 时间处理
				parse, _ := time2.Parse("Mon, 02 Jan 2006 15:04:05 -0700", item.PubDate)
				// 添加数据
				*rssLink = append(*rssLink, dto.RssLinkDTO{
					Title:   item.Title,
					Link:    item.Link,
					Summary: description,
					Timer:   gtime.NewFromTime(parse).TimestampMilliStr(),
				})
			}
			return rssLink, true
		}
	} else {
		glog.Warning(ctx, "解析 XML 成功，但数据为空")
	}
	return nil, false
}

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
func (s *sRssLogic) RssWithWordpressFeed(ctx context.Context, rssURL string) (rssLink *[]dto.RssLinkDTO, hasThis bool) {
	glog.Noticef(ctx, "[LOGIC] 尝试获取 WordPress 中 原生 的 feed 内容")
	getBody, err := s.rssLinkAccess(ctx, rssURL)
	if err != nil {
		return nil, false
	}
	// 解析 XML
	getFeed := new(dto.WordPressFeedDTO)
	err = xml.Unmarshal(getBody, &getFeed)
	if err != nil {
		glog.Warningf(ctx, "解析 XML 失败: %v", err.Error())
		return nil, false
	}
	// 处理数据
	if getFeed != nil {
		if match, _ := regexp.MatchString(`^https?://\S*wordpress\.org\S*$`, getFeed.Channel.Generator); match {
			rssLink = new([]dto.RssLinkDTO)
			for _, item := range getFeed.Channel.Items {
				description := item.Description
				if len(description) > 100 {
					description = description[:100]
				}
				// 时间处理
				parse, _ := time2.Parse("Mon, 02 Jan 2006 15:04:05 -0700", item.PubDate)
				// 添加数据
				*rssLink = append(*rssLink, dto.RssLinkDTO{
					Title:   item.Title,
					Link:    item.Link,
					Summary: description,
					Timer:   gtime.NewFromTime(parse).TimestampMilliStr(),
				})
			}
			return rssLink, true
		}
	} else {
		glog.Warning(ctx, "解析 XML 成功，但数据为空")
	}
	return nil, false
}

// rssLinkAccess 获取 RSS 链接信息
// 用于获取 RSS 链接信息
// 如果成功则返回 RSS 链接信息，否则返回错误
//
// 参数：
// ctx: 请求的上下文，用于管理超时和取消信号。
// rssURL: 站点的 rss 订阅地址。
//
// 返回：
// getBody: 如果获取 RSS 链接信息成功，返回 RSS 链接信息；否则返回错误。
// err: 如果获取 RSS 链接信息成功，返回 nil；否则返回错误。
func (s *sRssLogic) rssLinkAccess(ctx context.Context, rssURL string) ([]byte, error) {
	getResp, err := lutil.LinkAccess(rssURL)
	defer func() {
		if getResp != nil {
			_ = getResp.Body.Close()
		}
	}()
	// 处理错误检查信息
	if err != nil {
		glog.Warningf(ctx, "获取链接信息失败: %v", err)
		return nil, errors.New("获取链接信息失败")
	}
	// 获取 Body 并解析 XML
	if err != nil {
		glog.Warningf(ctx, "UTF8解析链接信息失败: %v", err)
		return nil, errors.New("UTF8解析链接信息失败")
	}
	getBody, err := io.ReadAll(getResp.Body)
	if err != nil {
		glog.Warningf(ctx, "获取链接信息失败: %v", err)
		return nil, errors.New("获取链接信息失败")
	}
	return getBody, nil
}
