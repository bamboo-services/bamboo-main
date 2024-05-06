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

package task

import (
	"context"
	"encoding/json"
	"github.com/gogf/gf/v2/os/glog"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/os/gtimer"
	"time"
	"xiaoMain/internal/dao"
	"xiaoMain/internal/model/do"
	"xiaoMain/internal/model/dto"
	"xiaoMain/internal/model/entity"
	"xiaoMain/internal/service"
)

// RssObtain 是一个定时任务，用于获取 RSS 订阅信息。
// 它会在每 10 分钟执行一次。
// 本接口会根据已有的链接信息对 RSS 信息进行获取，若获取失败返回失败信息，若成功返回成功信息
// 如果成功则返回 nil，否则返回错误
func RssObtain(ctx context.Context) {
	gtimer.Add(ctx, time.Minute*10, func(_ context.Context) {
		getNowTimestamp := time.Now().UnixMilli()
		glog.Notice(ctx, "[TASK] 开始操作 RSS 订阅")
		// 获取所有链接的Rss信息
		var getLinkHasRss []*entity.XfLinkList
		err := dao.XfLinkList.Ctx(ctx).
			Where(do.XfLinkList{Status: 1}).WhereNull("deleted_at").WhereNotNull("site_rss_url").
			Scan(&getLinkHasRss)
		if err != nil {
			glog.Warning(ctx, "[TASK] 获取链接信息失败，等待下次重试")
			glog.Warning(ctx, err)
			return // 如果获取失败则直接返回
		}
		// 遍历所有需要刷新的链接
		if len(getLinkHasRss) > 0 {
			for _, rssLink := range getLinkHasRss {
				var getRssInfo *entity.XfLinkRss
				err := dao.XfLinkRss.Ctx(ctx).Where(do.XfLinkRss{LinkId: rssLink.Id}).Scan(&getRssInfo)
				if err != nil {
					continue
				}
				if getRssInfo != nil && getRssInfo.CheckTime.Before(gtime.Now()) {
					continue
				}

				// 对 Rss 进行处理操作
				// RSS Hexo Feed
				link, this := service.RssLogic().RssWithHexoFeed(ctx, rssLink.SiteRssUrl)
				if this && saveData(ctx, link, rssLink) {
					continue
				}
				// RSS Hugo Feed
				link, this = service.RssLogic().RssWithHugoFeed(ctx, rssLink.SiteRssUrl)
				if this && saveData(ctx, link, rssLink) {
					continue
				}
				// RSS WordPress Feed
				link, this = service.RssLogic().RssWithWordpressFeed(ctx, rssLink.SiteRssUrl)
				if this && saveData(ctx, link, rssLink) {
					continue
				}
				// 无法正常获取
				glog.Warningf(ctx, "无法正常获取 %v 的 RSS 订阅信息", rssLink.SiteRssUrl)
			}
		}
		glog.Noticef(ctx, "[TASK] 操作 RSS 订阅完成, 耗时: %vms", time.Now().UnixMilli()-getNowTimestamp)
	})
}

// saveData 保存数据
// 处理输入的数据并保存到数据库
// 主要处理从 RSS 获取的数据并保存到数据库
func saveData(ctx context.Context, link *[]dto.RssLinkDTO, rssLink *entity.XfLinkList) bool {
	// 更新数据表
	marshal, err := json.Marshal(link)
	if err != nil {
		return false // 如果序列化失败则直接返回
	}
	_, err = dao.XfLinkRss.Ctx(ctx).Data(do.XfLinkRss{
		LinkId:    rssLink.Id,
		RssJson:   string(marshal),
		CheckTime: gtime.NewFromTime(time.Now().Add(time.Hour * 1)),
	}).OnConflict("link_id").Save()
	if err != nil {
		return false // 如果保存失败则直接返回
	}
	return true
}
