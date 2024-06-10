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
	"github.com/bamboo-services/bamboo-utils/bcode"
	"github.com/bamboo-services/bamboo-utils/berror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
	"xiaoMain/internal/dao"
	"xiaoMain/internal/model/do"
	"xiaoMain/internal/model/dto"
	"xiaoMain/internal/model/entity"
)

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
func (s *sRss) GetAllLinkRssInfo(ctx context.Context) (getRss *[]*dto.RssLinkDTO, err error) {
	g.Log().Notice(ctx, "[LOGIC] Rss:GetAllLinkRssInfo | 获取所有链接的Rss信息")
	var getRssInfo []*entity.LinkRss
	err = dao.LinkRss.Ctx(ctx).Scan(&getRssInfo)
	if err != nil {
		g.Log().Errorf(ctx, "[LOGIC] 获取 Rss 信息失败: %v", err.Error())
		return nil, berror.NewErrorHasError(bcode.ServerInternalError, err)
	}
	// 对数据进行并且进行插入
	getRss = new([]*dto.RssLinkDTO)
	// 从数据库遍历数据
	if len(getRssInfo) > 0 {
		for _, sqlLinkRss := range getRssInfo {
			// 获取单个遍历数据的信息
			var getSQLSingleRssInfo *[]*dto.RssLinkDTO
			err := json.Unmarshal([]byte(sqlLinkRss.RssJson), &getSQLSingleRssInfo)
			if err != nil {
				g.Log().Warningf(ctx, "[LOGIC] 解析 RssJson 失败: %v", err.Error())
				return nil, berror.NewError(bcode.OperationFailed, "数据解析失败<RssJson解析失败>")
			}
			// 检查数据是否为空
			*getRss = append(*getRss, *getSQLSingleRssInfo...)
		}
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
	// 时间戳转为时间
	if len(*getRss) > 0 {
		for _, linkDTO := range *getRss {
			linkDTO.Timer = gtime.NewFromStr(linkDTO.Timer).Format("Y年m月d日 H点i分")
		}
	}
	return getRss, nil
}

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
func (s *sRss) GetLinkRssInfoWithLinkID(
	ctx context.Context,
	linkID int64,
) (getRss *[]*dto.RssLinkDTO, err error) {
	g.Log().Notice(ctx, "[LOGIC] Rss:GetLinkRssInfoWithLinkID | 获取链接的Rss信息")
	var getLinkInfo *entity.LinkList
	err = dao.LinkList.Ctx(ctx).Where(do.LinkList{Id: linkID, Status: 1}).
		WhereNotNull(dao.LinkList.Columns().SiteRssUrl).Limit(1).Scan(&getLinkInfo)
	if err != nil {
		g.Log().Errorf(ctx, "[LOGIC] 获取链接信息失败: %v", err.Error())
		return nil, berror.NewErrorHasError(bcode.ServerInternalError, err)
	}
	return s.rssLinkToDTO(ctx, *getLinkInfo)
}

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
func (s *sRss) GetLinkRssWithLinkName(
	ctx context.Context,
	linkName string,
) (getRss *[]*dto.RssLinkDTO, err error) {
	g.Log().Notice(ctx, "[LOGIC] Rss:GetLinkRssWithLinkName | 获取链接的Rss信息")
	var getLinkInfo *entity.LinkList
	err = dao.LinkList.Ctx(ctx).Where(do.LinkList{SiteName: linkName, Status: 1}).
		WhereNotNull(dao.LinkList.Columns().SiteRssUrl).Limit(1).Scan(&getLinkInfo)
	if err != nil {
		g.Log().Errorf(ctx, "[LOGIC] 获取链接信息失败: %v", err.Error())
		return nil, berror.NewErrorHasError(bcode.ServerInternalError, err)
	}
	return s.rssLinkToDTO(ctx, *getLinkInfo)
}

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
func (s *sRss) GetLinkRssWithLinkLocation(
	ctx context.Context,
	linkLocation int64,
) (getRss *[]*dto.RssLinkDTO, err error) {
	g.Log().Notice(ctx, "[LOGIC] Rss:GetLinkRssWithLinkName | 获取链接的Rss信息")
	var getLinkInfo []*entity.LinkList
	err = dao.LinkList.Ctx(ctx).Where(do.LinkList{Location: linkLocation, Status: 1}).
		WhereNotNull(dao.LinkList.Columns().SiteRssUrl).Scan(&getLinkInfo)
	if err != nil {
		g.Log().Errorf(ctx, "[LOGIC] 获取链接信息失败: %v", err.Error())
		return nil, berror.NewErrorHasError(bcode.ServerInternalError, err)
	}
	getRss = new([]*dto.RssLinkDTO)
	if len(getLinkInfo) > 0 {
		for _, linkRss := range getLinkInfo {
			var getLinkRssInfo *entity.LinkRss
			err := dao.LinkRss.Ctx(ctx).Where(do.LinkRss{LinkId: linkRss.Id}).Scan(&getLinkRssInfo)
			if err != nil || getLinkRssInfo == nil {
				continue
			}
			// 获取单个遍历数据的信息
			var getSQLSingleRssInfo *[]*dto.RssLinkDTO
			err = json.Unmarshal([]byte(getLinkRssInfo.RssJson), &getSQLSingleRssInfo)
			if err != nil {
				g.Log().Warningf(ctx, "[LOGIC] 解析 RssJson 失败: %v", err.Error())
				return nil, berror.NewError(bcode.OperationFailed, "数据解析失败<RssJson解析失败>")
			}
			// 检查数据是否为空
			*getRss = append(*getRss, *getSQLSingleRssInfo...)
		}
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
	// 时间戳转为时间
	if len(*getRss) > 0 {
		for _, linkDTO := range *getRss {
			linkDTO.Timer = gtime.NewFromStr(linkDTO.Timer).Format("Y年m月d日 H点i分")
		}
	}
	return getRss, nil
}

// rssLinkToDTO
//
// # Rss 链接转 DTO
//
// 用于将 Rss 链接转换为 DTO，如果成功则返回 nil，否则返回错误。
//
// # 参数:
//   - ctx: 请求的上下文，用于管理超时和取消信号。
//   - getLinkInfo: 链接的信息。
//
// # 返回:
//   - getRss: 如果获取 Rss 信息成功，返回 nil；否则返回错误。
//   - err: 如果获取 Rss 信息成功，返回 nil；否则返回错误。
func (s *sRss) rssLinkToDTO(
	ctx context.Context,
	getLinkInfo entity.LinkList,
) (getRss *[]*dto.RssLinkDTO, err error) {
	g.Log().Notice(ctx, "[LOGIC] Rss:rssLinkToDTO | Rss 链接转 DTO")
	var getRssInfo *entity.LinkRss
	err = dao.LinkRss.Ctx(ctx).Where(do.LinkRss{LinkId: getLinkInfo.Id}).Scan(&getRssInfo)
	if err != nil {
		g.Log().Errorf(ctx, "[LOGIC] 获取 Rss 信息失败: %v", err.Error())
		return nil, errors.New("数据库获取失败<Rss信息提取失败>")
	}
	var getSQLRssInfo *[]*dto.RssLinkDTO
	err = json.Unmarshal([]byte(getRssInfo.RssJson), &getSQLRssInfo)
	if err != nil {
		g.Log().Warningf(ctx, "[LOGIC] 解析 RssJson 失败: %v", err.Error())
		return nil, errors.New("数据解析失败<RssJson解析失败>")
	}
	// 时间戳转换
	if len(*getSQLRssInfo) > 0 {
		for _, linkDTO := range *getSQLRssInfo {
			linkDTO.Timer = gtime.NewFromStr(linkDTO.Timer).Format("Y年m月d日 H点i分")
		}
	}
	return getSQLRssInfo, nil
}
