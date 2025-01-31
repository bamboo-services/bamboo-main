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
	"github.com/gogf/gf/v2/os/gtime"
	"xiaoMain/api/link/v1"
	"xiaoMain/internal/dao"
	"xiaoMain/internal/model/do"
	"xiaoMain/internal/model/entity"
)

func (s *sLink) EditLocation(ctx context.Context, req v1.EditLocationReq) (err error) {
	// 对数据进行获取
	var getLocation *entity.Location
	err = dao.Location.Ctx(ctx).Where(do.Location{Id: req.ID}).Scan(&getLocation)
	if err != nil {
		return berror.NewErrorHasError(bcode.ServerInternalError, err)
	}
	// 检查数据是否存在
	if getLocation == nil {
		return berror.NewError(bcode.NotExist, "位置不存在")
	}
	// 检查数据是否有重复
	count, err := dao.Location.Ctx(ctx).
		Where(do.Location{Name: req.Name}).
		WhereOr(do.Location{DisplayName: req.DisplayName}).
		Count()
	if err != nil {
		return berror.NewErrorHasError(bcode.ServerInternalError, err)
	}
	// 检查数据是否已经存在
	if count > 1 {
		return berror.NewError(bcode.AlreadyExists, "位置名称或展示名称已经存在")
	}
	// 数据不存在进行插入操作
	getLocation.Sort = req.Sort
	getLocation.Name = req.Name
	getLocation.DisplayName = req.DisplayName
	getLocation.Description = req.Description
	getLocation.Reveal = req.Reveal
	getLocation.UpdatedAt = gtime.Now()
	_, err = dao.Location.Ctx(ctx).Data(getLocation).Where(do.Location{Id: getLocation.Id}).Update()
	if err != nil {
		return berror.NewErrorHasError(bcode.ServerInternalError, err)
	}
	return nil
}

func (s *sLink) LocationMove(ctx context.Context, id, moveID int64) (err error) {
	// 检查原有 id 是否存在
	var getOriginalLocation *entity.Location
	err = dao.Location.Ctx(ctx).Where(do.Location{Id: id}).Scan(&getOriginalLocation)
	if err != nil {
		return berror.NewErrorHasError(bcode.ServerInternalError, err)
	}
	if getOriginalLocation == nil {
		return berror.NewError(bcode.NotExist, "原有链接不存在")
	}
	// 根据原有 id 获取所有匹配的友链
	var links []entity.LinkList
	err = dao.LinkList.Ctx(ctx).Where(do.LinkList{Location: id}).Scan(&links)
	if err != nil {
		return berror.NewErrorHasError(bcode.ServerInternalError, err)
	}
	// 检查修改后的链接是否存在
	if moveID != 0 {
		count, err := dao.Location.Ctx(ctx).Where(do.Location{Id: moveID}).Count()
		if err != nil {
			return berror.NewErrorHasError(bcode.ServerInternalError, err)
		}
		if count > 0 {
			// 检查链接是否不为空
			if len(links) > 0 {
				for _, link := range links {
					// 对链接进行操作切换 LocationID 位置
					link.Location = moveID
					_, err := dao.LinkList.Ctx(ctx).Data(link).Where(do.LinkList{Id: link.Id}).Update()
					if err != nil {
						return berror.NewErrorHasError(bcode.ServerInternalError, err)
					}
				}
			}
		} else {
			return berror.NewError(bcode.NotExist, "链接不存在")
		}
	} else {
		// 获取优先级最高
		var getLastLocation []entity.Location
		err := dao.Location.Ctx(ctx).OrderAsc(dao.Location.Columns().Sort).Limit(2).Scan(&getLastLocation)
		if err != nil {
			return berror.NewErrorHasError(bcode.ServerInternalError, err)
		}
		if len(getLastLocation) == 1 {
			return berror.NewError(bcode.OperationFailed, "当前仅有一个位置，不可删除")
		}
		// 检查链接是否不为空
		if len(links) > 0 {
			for _, link := range links {
				// 对链接进行操作切换 LocationID 位置
				link.Location = moveID
				var err error
				if getLastLocation[0].Id == id {
					_, err = dao.LinkList.Ctx(ctx).Data(link).Where(do.LinkList{
						Id:              getLastLocation[1].Id,
						DesiredLocation: getLastLocation[1].Id,
					}).Update()
				} else {
					_, err = dao.LinkList.Ctx(ctx).Data(link).Where(do.LinkList{
						Id:              getLastLocation[0].Id,
						DesiredLocation: getLastLocation[0].Id,
					}).Update()
				}
				if err != nil {
					return berror.NewErrorHasError(bcode.ServerInternalError, err)
				}
			}
		}
	}
	return nil
}

func (s *sLink) DelLocation(ctx context.Context, id int64) (err error) {
	// 检查原有 id 是否存在
	var getOriginalLocation *entity.Location
	err = dao.Location.Ctx(ctx).Where(do.Location{Id: id}).Scan(&getOriginalLocation)
	if err != nil {
		return berror.NewErrorHasError(bcode.ServerInternalError, err)
	}
	if getOriginalLocation == nil {
		return berror.NewError(bcode.NotExist, "原有链接不存在")
	}
	// 对 id 进行删除
	_, err = dao.Location.Ctx(ctx).Where(do.Location{Id: getOriginalLocation.Id}).Delete()
	if err != nil {
		return berror.NewErrorHasError(bcode.ServerInternalError, err)
	}
	return nil
}
