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
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
	"xiaoMain/api/link/v1"
	"xiaoMain/internal/model/vo"
	"xiaoMain/internal/service"
)

// CheckRssURLHasConnect
//
// # 检查链接是否已经连接
//
// 检查链接是否已经连接, 需要用户提供链接的URL。
//
// # 参数
//   - ctx: 请求的上下文，用于管理超时和取消信号。
//   - req: 用户的请求，包含检查链接是否已经连接的详细信息。
//
// # 返回
//   - res: 发送给用户的响应。如果检查链接是否已经连接成功，它将返回成功的消息。
//   - err: 在检查链接是否已经连接过程中发生的任何错误。
func (c *ControllerV1) CheckRssURLHasConnect(
	ctx context.Context,
	req *v1.CheckRssURLHasConnectReq,
) (res *v1.CheckRssURLHasConnectRes, err error) {
	g.Log().Notice(ctx, "[CONTROL] 控制层 CheckRssURLHasConnect 接口")
	// 获取 RSS 链接是否已经连接
	getNowTimestamp := gtime.TimestampMilli()
	err = service.Link().CheckRSSCanAccess(ctx, req.RssURL)
	if err != nil {
		return nil, err
	}
	getTime := gtime.TimestampMilli() - getNowTimestamp
	return &v1.CheckRssURLHasConnectRes{
		LinkConnectRes: vo.LinkConnectRes{
			Message: "站点读取成功",
			Delay:   &getTime,
		},
	}, nil
}
