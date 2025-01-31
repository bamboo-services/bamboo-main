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
	"github.com/bamboo-services/bamboo-utils/bcode"
	"github.com/bamboo-services/bamboo-utils/berror"
	"github.com/young2j/gocopy"
	"xiaoMain/internal/dao"
	"xiaoMain/internal/model/do"
	"xiaoMain/internal/model/dto/flow"
	"xiaoMain/internal/model/entity"
)

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
func (s *sLink) GetLink(ctx context.Context) (getLink []flow.LinkGetDTO, total uint64, err error) {
	// 获取所有链接
	linkList := make([]entity.LinkList, 0)
	linkTotal := 0
	err = dao.LinkList.Ctx(ctx).Where(do.LinkList{Status: 1}).OrderDesc(dao.LinkList.Columns().CreatedAt).
		ScanAndCount(&linkList, &linkTotal, true)
	if err != nil {
		return nil, 0, berror.NewErrorHasError(bcode.ServerInternalError, err)
	}
	// 获取所在位置
	locationList := make([]entity.Location, 0)
	locationTotal := 0
	err = dao.Location.Ctx(ctx).Where(do.Location{Reveal: true}).OrderAsc(dao.Location.Columns().Sort).
		ScanAndCount(&locationList, &locationTotal, true)
	if err != nil {
		return nil, 0, berror.NewErrorHasError(bcode.ServerInternalError, err)
	}
	// 对链接进行序列化
	getLinkDTO := make([]flow.LinkGetDTO, 0)
	for _, location := range locationList {
		locationLinkDTO := make([]flow.LinkInfoDTO, 0)
		locationLinkTotal := 0
		for _, list := range linkList {
			if list.Location == location.Id {
				newLink := new(flow.LinkInfoDTO)
				gocopy.Copy(newLink, list)
				locationLinkDTO = append(locationLinkDTO, *newLink)
				locationLinkTotal++
			}
		}
		getLinkDTO = append(getLinkDTO, flow.LinkGetDTO{
			Location: location,
			Total:    locationLinkTotal,
			Link:     locationLinkDTO,
		})
	}
	return getLinkDTO, uint64(locationTotal), nil
}

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
func (s *sLink) GetLinkAdmin(ctx context.Context) (getAllLink []entity.LinkList, total uint64, err error) {
	// 获取所有的链接
	linkList := make([]entity.LinkList, 0)
	linkTotal := 0
	err = dao.LinkList.Ctx(ctx).OrderDesc(dao.LinkList.Columns().CreatedAt).ScanAndCount(&linkList, &linkTotal, true)
	if err != nil {
		return nil, 0, berror.NewErrorHasError(bcode.ServerInternalError, err)
	}
	return linkList, uint64(linkTotal), nil
}
