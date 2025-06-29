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

type LinkFriendAddReq struct {
	g.Meta       `path:"/link/friend" method:"POST" tags:"链接控制器" summary:"添加友情链接" dc:"添加友情链接接口，允许用户提交新的友情链接请求"`
	Name         string `json:"name" v:"required|max-length:100#请输入友链名称|友链名称长度不能超过100个字符" dc:"友链名称，不能为空，表示友情链接的名称"`
	URL          string `json:"url" v:"required|url#请输入友链地址|请输入有效的URL地址" dc:"友链地址，不能为空，表示友情链接的URL地址，必须是有效的URL格式"`
	Avatar       string `json:"avatar" v:"url#请输入有效的头像URL地址" dc:"友链头像地址，表示友情链接的头像图片URL地址，必须是有效的URL格式"`
	Email        string `json:"email" v:"email#请输入有效的电子邮箱地址" dc:"联系邮箱，表示友情链接的联系邮箱，必须是有效的电子邮箱格式"`
	Description  string `json:"description" v:"max-length:2048#友链描述长度不能超过2048个字符" dc:"友链描述，表示友情链接的详细描述信息，长度不能超过2048个字符"`
	Group        string `json:"group" v:"required#请输入分组名称" dc:"分组名称，不能为空，表示友情链接所属的分组UUID"`
	Color        string `json:"color" v:"required#请输入颜色名称" dc:"颜色名称，不能为空，表示友情链接的颜色名称，必须是已存在的颜色UUID"`
	Order        int    `json:"order" default:"50" v:"required|between:0,100#请输入友链排序|友链顺序只允许 0 到 100" dc:"友链排序，数字越小越靠前，默认为0，表示友情链接的显示顺序"`
	ReviewRemark string `json:"review_remark" v:"max-length:10240#审核备注长度不能超过10240个字符" dc:"审核备注，表示管理员对友情链接的审核意见，长度不能超过10240个字符"`
}

type LinkFriendAddRes struct {
	g.Meta `mime:"application/json; charset=utf-8" dc:"添加友情链接响应"`
	*bmodels.ResponseDTO[types.Nil]
}
