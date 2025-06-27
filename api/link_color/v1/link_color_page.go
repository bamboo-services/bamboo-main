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

type LinkColorPageReq struct {
	g.Meta `path:"/link/colors" method:"GET" tags:"链接颜色控制器" summary:"获取链接颜色列表" dc:"获取友情链接颜色列表接口，返回所有可用的友情链接颜色信息"`
	Search string `json:"search" in:"query" v:"max-length:100#搜索内容长度不能超过100个字符" dc:"搜索内容，允许用户通过颜色名称或描述进行模糊搜索"`
	Page   int    `json:"page" default:"1" in:"query" v:"required|min:1#请输入页码|页码不能小于 1" dc:"页码，当前页码，默认为1"`
	Size   int    `json:"size" default:"20" in:"query" v:"required|between:1,100#请输入每页数量|每页数量必须在1到100之间" dc:"每页数量，每页显示的颜色数量，默认为20"`
}

type LinkColorPageRes struct {
	g.Meta `mime:"application/json; charset=utf-8" dc:"获取友情链接颜色列表响应"`
	*bmodels.ResponseDTO[dto.Page[base.LinkColorDTO]]
}
