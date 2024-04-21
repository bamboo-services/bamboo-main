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
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/os/glog"
	"sync"
	"xiaoMain/api/link/v1"
	"xiaoMain/internal/logic/link"
	"xiaoMain/internal/service"
	"xiaoMain/utility/result"
)

// LinkAdd 是 ControllerV1 结构体的一个方法。
// 它处理用户尝试添加链接的过程。
//
// 参数:
// ctx: 请求的上下文，用于管理超时和取消信号。
// req: 用户的请求，包含添加链接的详细信息。
//
// 返回:
// res: 发送给用户的响应。如果添加链接成功，它将返回成功的消息。
func (c *ControllerV1) LinkAdd(ctx context.Context, req *v1.LinkAddReq) (res *v1.LinkAddRes, err error) {
	glog.Info(ctx, "[CONTROL] 控制层 LinkAdd 接口")
	// 获取 Request
	getRequest := ghttp.RequestFromCtx(ctx)
	service.RegisterLinkLogic(link.New())
	// 异步操作
	getError := error(nil)
	wg := sync.WaitGroup{}
	wg.Add(1)
	// 检查网站名是否重复
	go func(request *v1.LinkAddReq) {
		if getError = service.LinkLogic().CheckLinkName(ctx, request.SiteName); getError != nil {
			result.AddLinkFailed.SetErrorMessage(getError.Error()).Response(getRequest)
		}
		wg.Done()
	}(req)
	// 检查网站链接是否重复
	wg.Add(1)
	go func(request *v1.LinkAddReq) {
		if getError = service.LinkLogic().CheckLinkURL(ctx, request.SiteURL); getError != nil {
			result.AddLinkFailed.SetErrorMessage(getError.Error()).Response(getRequest)
		}
		wg.Done()
	}(req)
	// 检查链接是否可以访问
	wg.Add(1)
	go func(request *v1.LinkAddReq) {
		if getError = service.LinkLogic().CheckLinkCanAccess(ctx, request.SiteURL); getError != nil {
			result.AddLinkFailed.SetErrorMessage(getError.Error()).Response(getRequest)
		}
		wg.Done()
	}(req)
	// 检查 Logo 是否可以访问
	wg.Add(1)
	go func(request *v1.LinkAddReq) {
		if getError = service.LinkLogic().CheckLogoCanAccess(ctx, request.SiteLogo); getError != nil {
			result.AddLinkFailed.SetErrorMessage(getError.Error()).Response(getRequest)
		}
		wg.Done()
	}(req)
	// 检查 RSS URL 是否合法
	wg.Add(1)
	go func(request *v1.LinkAddReq) {
		if getError = service.LinkLogic().CheckRSSURL(ctx, request.SiteRssURL); getError != nil {
			result.AddLinkFailed.SetErrorMessage(getError.Error()).Response(getRequest)
		}
		wg.Done()
	}(req)
	// 等待异步操作完成
	wg.Wait()
	// 对内容进行插入
	if getError = service.LinkLogic().AddLink(ctx, *req); getError != nil {
		result.AddLinkFailed.SetErrorMessage(getError.Error()).Response(getRequest)
	}
	return nil, nil
}
