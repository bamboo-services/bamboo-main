/*
 * --------------------------------------------------------------------------------
 * Copyright (c) 2016-NOW(至今) 筱锋
 * Author: 筱锋「xiao_lfeng」(https://www.x-lf.com)
 * --------------------------------------------------------------------------------
 * 许可证声明：版权所有 (c) 2016-2026 筱锋。保留所有权利。
 * 有关MIT许可证的更多信息，请查看项目根目录下的LICENSE文件或访问：
 * https://opensource.org/licenses/MIT
 * --------------------------------------------------------------------------------
 */

package apiAuth

import (
	xSnowflake "github.com/bamboo-services/bamboo-base-go/common/snowflake"
	"github.com/bamboo-services/bamboo-main/internal/entity"
	"github.com/bamboo-services/bamboo-main/internal/models/base"
)

// UserIDRequest 系统用户路径参数
type UserIDRequest struct {
	ID xSnowflake.SnowflakeID `uri:"id" binding:"required"`
}

// UserQueryRequest 系统用户分页查询请求
//
// Status 使用指针类型：status=0（禁用）时 form 解码写入指针指向的 0，
// omitempty 仅对 nil 指针放行，oneof=0 1 校验 0 通过，避免禁用筛选被零值吞掉。
type UserQueryRequest struct {
	Page      int    `form:"page" binding:"omitempty,min=1" example:"1"`                            // 页码，默认1
	PageSize  int    `form:"page_size" binding:"omitempty,min=1,max=100" example:"10"`              // 每页数量，默认10
	Keyword   string `form:"keyword" binding:"omitempty,max=100" example:"admin"`                   // 搜索关键词（用户名/邮箱/昵称模糊匹配）
	Status    *int   `form:"status" binding:"omitempty,oneof=0 1" example:"1"`                      // 用户状态（0: 禁用, 1: 启用）
	SortBy    string `form:"sort_by" binding:"omitempty,oneof=created_at updated_at username email last_login_at" example:"created_at"` // 排序字段
	SortOrder string `form:"sort_order" binding:"omitempty,oneof=asc desc" example:"desc"`           // 排序方式
}

// UserListResponse 系统用户分页列表响应
//
// Password 与 OAuthUserID 因 json:"-" 在序列化时自动脱敏，不会随列表返回。
type UserListResponse struct {
	base.PaginationResponse[entity.SystemUser]
}

// UserStatusRequest 更新系统用户状态请求
type UserStatusRequest struct {
	Status int `json:"status" binding:"required,oneof=0 1" example:"1"` // 目标状态（0: 禁用, 1: 启用）
}
