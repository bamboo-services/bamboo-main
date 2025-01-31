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
	"github.com/gogf/gf/v2/frame/g"
	v1 "xiaoMain/api/link/v1"
	"xiaoMain/internal/dao"
	"xiaoMain/internal/model/do"
	"xiaoMain/internal/model/entity"
)

// GetSingleLink
//
// # 获取单个链接
//
// 获取单个链接，用于获取单个链接之后进行展示使用。接口将会查询 ID 的链接信息是否存在，若存在将数据返回；
// 若不存在则返回错误信息。
// 若产生数据库等非业务类型报错，将会执行 bcode.ServerInternalError 错误码。
//
// # 参数
//   - ctx: 请求的上下文，用于管理超时和取消信号。
//   - linkID: 链接的 ID。
//
// # 返回
//   - linkInfo: 链接的信息。
//   - error: 返回的错误码信息
func (s *sLink) GetSingleLink(ctx context.Context, linkID int64) (linkInfo *entity.LinkList, err error) {
	g.Log().Notice(ctx, "[LOGIC] Link:GetSingleLink | 获取单个链接")
	// 查询单个链接
	var link *entity.LinkList
	err = dao.LinkList.Ctx(ctx).Where(do.LinkList{Id: linkID}).Scan(&link)
	if err != nil {
		return nil, berror.NewErrorHasError(bcode.ServerInternalError, err)
	}
	// 检查链接是否为空
	if link == nil {
		return nil, berror.NewError(bcode.NotExist, "链接不存在")
	}
	return link, nil
}

// EditLink
//
// # 编辑链接
//
// 获取链接的 ID 随后，对链接的信息进行获取操作，用于编辑链接的信息。当完成修改操作后将返回 nil；
// 若操作的过程中产生问题，将不会存入数据库，并返回 error
//
// # 参数
//   - ctx: 请求的上下文，用于管理超时和取消信号。
//   - data: 用户的请求，包含编辑链接的详细信息。
//
// # 返回
//   - error: 在编辑链接过程中发生的任何错误。
func (s *sLink) EditLink(ctx context.Context, data v1.LinkEditReq) (err error) {
	g.Log().Notice(ctx, "[LOGIC] Link:EditLink | 编辑链接")
	// 查询链接
	info, err := s.GetSingleLink(ctx, data.ID)
	if err != nil {
		return err
	}
	// 检查错误参数
	if data.Location == 0 {
		return berror.NewError(bcode.OperationFailed, "位置不能为空")
	}
	if *(data.Color) != 0 {
		info.Color = *(data.Color)
	} else {
		data.Color = nil
	}
	// 更新链接信息
	_, err = dao.LinkList.Ctx(ctx).Data(data).Where(do.LinkList{Id: data.ID}).Update()
	if err != nil {
		return berror.NewErrorHasError(bcode.ServerInternalError, err)
	}
	return nil
}

// AddLinkAdmin
//
// # 添加链接
//
// 添加链接, 由管理员直接进行操作；添加的链接可以直接在用户界面进行查看，若创建出现错误则会抛出错误，否则将会返回 nil
//
// # 参数
//   - ctx: 请求的上下文，用于管理超时和取消信号。
//   - data: 用户的请求，包含添加链接的详细信息。
//
// # 返回
//   - error: 在添加链接过程中发生的任何错误。
func (s *sLink) AddLinkAdmin(ctx context.Context, data v1.LinkAddAdminReq) (err error) {
	g.Log().Notice(ctx, "[LOGIC] Link:AddLinkAdmin | 添加链接")
	// 对数据进行获取，检查位置，颜色以及所填写信息是否已存在
	if data.Location == 0 {
		return berror.NewError(bcode.OperationFailed, "位置不能为空")
	}
	// 检查数据是已存在
	count, err := dao.LinkList.Ctx(ctx).
		Where(do.LinkList{SiteName: data.SiteName}).
		WhereOr(do.LinkList{SiteUrl: data.SiteURL}).
		Count()
	if err != nil {
		return berror.NewErrorHasError(bcode.ServerInternalError, err)
	}
	if count > 0 {
		return berror.NewError(bcode.OperationFailed, "链接已存在")
	}
	// 插入数据
	_, err = dao.LinkList.Ctx(ctx).Data(do.LinkList{
		WebmasterEmail:  data.WebmasterEmail,
		ServiceProvider: data.ServiceProvider,
		SiteName:        data.SiteName,
		SiteUrl:         data.SiteURL,
		SiteLogo:        data.SiteLogo,
		SiteDescription: data.SiteDescription,
		SiteRssUrl:      data.SiteRssURL,
		HasAdv:          data.HasAdv,
		DesiredLocation: data.Location,
		Location:        data.Location,
		DesiredColor:    data.Color,
		Color:           data.Color,
		Remark:          data.Remark,
		Status:          1,
		AbleConnect:     true,
	}).Insert()
	if err != nil {
		return berror.NewErrorHasError(bcode.ServerInternalError, err)
	}
	return nil
}

// Verify
//
// # 审核链接
//
// 审核链接
//
// # 参数
//   - ctx: 请求的上下文，用于管理超时和取消信号。
//   - req: 用户的请求，包含审核链接的详细信息。
//
// # 返回
//   - error: 在审核链接过程中发生的任何错误。
func (s *sLink) Verify(ctx context.Context, req *v1.LinkVerifyReq) (err error) {
	g.Log().Notice(ctx, "[LOGIC] Link:Verify | 审核链接")
	// 根据 id 获取链接信息
	info, err := s.GetSingleLink(ctx, req.Id)
	if req.Status {
		info.Status = 1
		info.Location = req.DesiredLocation
		info.Color = req.DesiredColor
		_, err := dao.LinkList.Ctx(ctx).Data(info).Where(do.LinkList{Id: req.Id}).Update()
		if err != nil {
			return berror.NewErrorHasError(bcode.ServerInternalError, err)
		}
	} else {
		info.Status = 2
		_, err := dao.LinkList.Ctx(ctx).Data(info).Where(do.LinkList{Id: req.Id}).Update()
		if err != nil {
			return berror.NewErrorHasError(bcode.ServerInternalError, err)
		}
	}
	return nil
}
