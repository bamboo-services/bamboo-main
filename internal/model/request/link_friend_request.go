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

package request

// LinkFriendAddReq 添加友情链接请求
type LinkFriendAddReq struct {
	LinkName        string `json:"link_name" binding:"required,min=1,max=100" example:"示例网站"`
	LinkURL         string `json:"link_url" binding:"required,url,max=500" example:"https://example.com"`
	LinkAvatar      string `json:"link_avatar" binding:"omitempty,url,max=500" example:"https://example.com/avatar.jpg"`
	LinkRSS         string `json:"link_rss" binding:"omitempty,url,max=500" example:"https://example.com/rss.xml"`
	LinkDesc        string `json:"link_desc" binding:"omitempty,max=500" example:"这是一个示例网站"`
	LinkEmail       string `json:"link_email" binding:"omitempty,email,max=100" example:"admin@example.com"`
	LinkGroupID     int64  `json:"link_group_id" binding:"omitempty,number" example:"1"`
	LinkColorID     int64  `json:"link_color_id" binding:"omitempty,number" example:"1"`
	LinkOrder       int    `json:"link_order" binding:"omitempty,min=0" example:"0"`
	LinkApplyRemark string `json:"link_apply_remark" binding:"omitempty,max=500" example:"申请友链"`
}

// LinkFriendUpdateReq 更新友情链接请求
type LinkFriendUpdateReq struct {
	LinkName        string `json:"link_name" binding:"omitempty,min=1,max=100" example:"示例网站"`
	LinkURL         string `json:"link_url" binding:"omitempty,url,max=500" example:"https://example.com"`
	LinkAvatar      string `json:"link_avatar" binding:"omitempty,url,max=500" example:"https://example.com/avatar.jpg"`
	LinkRSS         string `json:"link_rss" binding:"omitempty,url,max=500" example:"https://example.com/rss.xml"`
	LinkDesc        string `json:"link_desc" binding:"omitempty,max=500" example:"这是一个示例网站"`
	LinkEmail       string `json:"link_email" binding:"omitempty,email,max=100" example:"admin@example.com"`
	LinkGroupID     int64  `json:"link_group_id" binding:"omitempty,number" example:"1"`
	LinkColorID     int64  `json:"link_color_id" binding:"omitempty,number" example:"1"`
	LinkOrder       *int   `json:"link_order" binding:"omitempty,min=0" example:"0"`
	LinkApplyRemark string `json:"link_apply_remark" binding:"omitempty,max=500" example:"申请友链"`
}

// LinkFriendQueryReq 查询友情链接请求
type LinkFriendQueryReq struct {
	Page        int    `form:"page" binding:"omitempty,min=1" example:"1"`
	PageSize    int    `form:"page_size" binding:"omitempty,min=1,max=100" example:"10"`
	LinkName    string `form:"link_name" binding:"omitempty,max=100" example:"示例"`
	LinkStatus  *int   `form:"link_status" binding:"omitempty,oneof=0 1 2" example:"1"`
	LinkFail    *int   `form:"link_fail" binding:"omitempty,oneof=0 1" example:"0"`
	LinkGroupID int64  `form:"link_group_id" binding:"omitempty,number" example:"1"`
	SortBy      string `form:"sort_by" binding:"omitempty,oneof=created_at updated_at link_order link_name" example:"created_at"`
	SortOrder   string `form:"sort_order" binding:"omitempty,oneof=asc desc" example:"desc"`
}

// LinkFriendStatusReq 更新友情链接状态请求
type LinkFriendStatusReq struct {
	LinkStatus       int    `json:"link_status" binding:"required,oneof=0 1 2" example:"1"`
	LinkReviewRemark string `json:"link_review_remark" binding:"omitempty,max=500" example:"审核通过"`
}

// LinkFriendFailReq 更新友情链接失效状态请求
type LinkFriendFailReq struct {
	LinkFail       int    `json:"link_fail" binding:"required,oneof=0 1" example:"1"`
	LinkFailReason string `json:"link_fail_reason" binding:"omitempty,max=500" example:"网站无法访问"`
}
