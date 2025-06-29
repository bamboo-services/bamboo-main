/*
 * --------------------------------------------------------------------------------
 * Copyright (c) 2016-NOW(至今) 筱锋
 * Author: 筱锋「xiao_lfeng」(https://www.x-lf.com)
 * --------------------------------------------------------------------------------
 * 许可证声明：版权所有 (c) 2016-2025 筱锋。保留所有权利。
 * 有关MIT许可证的更多信息，请查看项目根目录下的LICENSE文件或访问：
 * https://opensource.org/licenses/MIT
 * --------------------------------------------------------------------------------
 */

package v1

import (
	"bamboo-main/internal/model/dto"
	"bamboo-main/internal/model/dto/base"
	"github.com/XiaoLFeng/bamboo-utils/bmodels"
	"github.com/gogf/gf/v2/frame/g"
)

type LinkFriendPageReq struct {
	g.Meta `path:"/link/friends" method:"GET" tags:"链接控制器" summary:"获取友情链接列表" dc:"获取友情链接列表的接口，支持分页和筛选功能"`
	Search string `json:"search" in:"query" v:"max-length:100#搜索内容长度不能超过100个字符" dc:"搜索内容，表示要搜索的关键词，可以为空，长度不能超过100个字符"`
	Page   int    `json:"page" in:"query" default:"1" v:"required#请输入页码" dc:"页码，不能为空，用于指定要获取的页数"`
	Size   int    `json:"size" in:"query" default:"20" v:"required|between:1,100#请输入每页数量|每页数量必须在1到100之间" dc:"每页数量，不能为空，表示每页显示的友情链接数量，必须在1到100之间"`
	IsAll  bool   `json:"is_all" in:"query" default:" false" dc:"是否获取全部数据，默认为 false，表示只获取有效的链接，如果为 true，则获取所有链接，包括已失效的链接"`
}

type LinkFriendPageRes struct {
	g.Meta `mime:"application/json; charset=utf-8" dc:"获取友情链接列表响应"`
	*bmodels.ResponseDTO[dto.Page[base.LinkFriendDTO]]
}
