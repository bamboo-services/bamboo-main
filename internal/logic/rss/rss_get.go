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
	"encoding/json"
	"errors"
	"github.com/gogf/gf/v2/os/glog"
	"github.com/gogf/gf/v2/os/gtime"
	"xiaoMain/internal/dao"
	"xiaoMain/internal/model/dto"
	"xiaoMain/internal/model/entity"
)

// GetAllLinkRssInfo 获取所有链接的Rss信息
// 用于获取所有链接的Rss信息
// 如果成功则返回 nil，否则返回错误
// 本接口会根据已有的链接信息对Rss信息进行获取，若获取失败返回失败信息，若成功返回成功信息
//
// 参数：
// ctx: 请求的上下文，用于管理超时和取消信号。
//
// 返回：
// err: 如果获取Rss信息成功，返回 nil；否则返回错误。
func (s *sRssLogic) GetAllLinkRssInfo(ctx context.Context) (getRss *[]*dto.RssLinkDTO, err error) {
	glog.Noticef(ctx, "[LOGIC] 执行 RssLogic:GetAllLinkRssInfo 服务层")
	var getRssInfo []*entity.XfLinkRss
	err = dao.XfLinkRss.Ctx(ctx).Scan(&getRssInfo)
	if err != nil {
		glog.Errorf(ctx, "[LOGIC] 获取 Rss 信息失败: %v", err.Error())
		return nil, errors.New("数据库获取失败<Rss信息提取失败>")
	}
	// 对数据进行并且进行插入
	getRss = new([]*dto.RssLinkDTO)
	// 从数据库遍历数据
	for _, sqlLinkRss := range getRssInfo {
		// 获取单个遍历数据的信息
		var getSQLSingleRssInfo *[]*dto.RssLinkDTO
		err := json.Unmarshal([]byte(sqlLinkRss.RssJson), &getSQLSingleRssInfo)
		if err != nil {
			glog.Warningf(ctx, "[LOGIC] 解析 RssJson 失败: %v", err.Error())
			return nil, errors.New("数据解析失败<RssJson解析失败>")
		}
		// 检查数据是否为空
		*getRss = append(*getRss, *getSQLSingleRssInfo...)
	}
	if len(*getRss) > 0 {
		for i, linkDTO := range *getRss {
			// 按照时间戳的大小进行排序
			for j, reLinkDTO := range *getRss {
				thisTime := gtime.NewFromStr(linkDTO.Timer)
				nextTime := gtime.NewFromStr(reLinkDTO.Timer)
				if nextTime.Before(thisTime) {
					(*getRss)[i], (*getRss)[j] = (*getRss)[j], (*getRss)[i]
				}
			}
		}
	}
	// 限制输出前 20 个
	if len(*getRss) > 20 {
		*getRss = (*getRss)[:20]
	}
	// 时间戳转为时间
	if len(*getRss) > 0 {
		for _, linkDTO := range *getRss {
			linkDTO.Timer = gtime.NewFromStr(linkDTO.Timer).Format("Y年m月d日 H点i分")
		}
	}
	return getRss, nil
}
