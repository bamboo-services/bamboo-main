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
	"github.com/XiaoLFeng/bamboo-utils/bmodels"
	"github.com/gogf/gf/v2/frame/g"
	"go/types"
)

type ApplyLinkReq struct {
	g.Meta      `path:"/apply/link" method:"POST" tags:"申请控制器" sm:"链接申请" dc:"用于申请链接，申请站内链接，在站内显示使用"`
	Name        string `json:"name" v:"required|max-length:100#请输入友链名称|友链名称长度不能超过100个字符" dc:"友链名称，不能为空，表示友情链接的名称"`
	URL         string `json:"url" v:"required|url#请输入友链地址|请输入有效的URL地址" dc:"友链地址，不能为空，表示友情链接的URL地址，必须是有效的URL格式"`
	Avatar      string `json:"avatar" v:"url#请输入有效的头像URL地址" dc:"友链头像地址，表示友情链接的头像图片URL地址，必须是有效的URL格式"`
	Rss         string `json:"rss" v:"url#请输入有效的RSS地址" dc:"友链RSS地址，表示友情链接的RSS订阅地址，必须是有效的URL格式"`
	Email       string `json:"email" v:"required|email#请输入有效的电子邮箱地址|请输入有效的电子邮箱地址" dc:"联系邮箱，必须输入，表示友情链接的联系邮箱，必须是有效的电子邮箱格式"`
	Description string `json:"description" v:"max-length:2048#友链描述长度不能超过2048个字符" dc:"友链描述，表示友情链接的详细描述信息，长度不能超过2048个字符"`
	Group       string `json:"group" v:"required#请输入分组UUID" dc:"分组名称，不能为空，表示友情链接所属的分组UUID"`
	Color       string `json:"color" v:"required#请输入颜色UUID" dc:"颜色名称，不能为空，表示友情链接的颜色名称，必须是已存在的颜色UUID"`
	ApplyRemark string `json:"link_apply_remark" dc:"申请者备注，用于申请者额外填写的一些信息作为备注给管理员看使用"`
}

type ApplyLinkRes struct {
	g.Meta `mime:"application/json; charset=utf-8" dc:"申请链接响应"`
	*bmodels.ResponseDTO[types.Nil]
}
