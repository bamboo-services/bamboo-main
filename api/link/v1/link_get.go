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

package v1

import (
	"github.com/gogf/gf/v2/frame/g"
	"xiaoMain/internal/model/dto/flow"
	"xiaoMain/internal/model/entity"
)

// LinkGetReq
//
// # 获取链接
//
// 获取链接，包括友情链接、社交链接等。
type LinkGetReq struct {
	g.Meta `path:"/link" method:"Get" tags:"链接控制器" summary:"获取链接" dc:"获取链接，包括友情链接、社交链接等。用于获取链接之后进行展示使用"`
}

// LinkGetRes
//
// # 获取链接
//
// 获取链接的响应。
type LinkGetRes struct {
	g.Meta        `mime:"application/json"`
	Links         []flow.LinkGetDTO `json:"locations" dc:"链接列表"`
	LocationTotal uint64            `json:"location_total" dc:"位置总数"`
}

// LinkGetAdminReq
//
// # 获取链接
//
// 获取链接，包括友情链接、社交链接等。
//
// # 参数:
//   - UserToken: 用户Token
type LinkGetAdminReq struct {
	g.Meta    `path:"/link/admin" method:"Get" tags:"链接控制器" summary:"管理员获取链接" dc:"获取链接，包括友情链接、社交链接等。"`
	UserToken string `json:"Authorization" in:"header" required:"true" dc:"用户Token"`
}

// LinkGetAdminRes
//
// # 获取链接
//
// 获取链接的响应。
type LinkGetAdminRes struct {
	g.Meta           `mime:"application/json"`
	Links            []entity.LinkList `json:"links" dc:"链接列表"`
	Total            uint64            `json:"total" dc:"链接总数"`
	Reviewed         uint64            `json:"reviewed" dc:"待审核链接总数"`
	DeleteTheTotal   uint64            `json:"deleted" dc:"删除链接总数"`
	RecentlyAdded    uint64            `json:"recently_added" dc:"最近添加链接总数"`
	RecentlyModified uint64            `json:"recently_modified" dc:"最近修改链接总数"`
}
