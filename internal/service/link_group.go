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

package service

import (
	"bamboo-main/internal/logic"
	"bamboo-main/internal/model/base"
	"bamboo-main/internal/model/dto"
	"bamboo-main/internal/model/request"

	xError "github.com/bamboo-services/bamboo-base-go/error"
	"github.com/gin-gonic/gin"
)

// ILinkGroupService 友链分组服务接口
type ILinkGroupService interface {
	// Add 添加友链分组
	Add(ctx *gin.Context, req *request.LinkGroupAddReq) (*dto.LinkGroupDetailDTO, *xError.Error)

	// Update 更新友链分组（名称和描述）
	Update(ctx *gin.Context, groupUUID string, req *request.LinkGroupUpdateReq) (*dto.LinkGroupDetailDTO, *xError.Error)

	// UpdateSort 批量更新友链分组排序
	UpdateSort(ctx *gin.Context, req *request.LinkGroupSortReq) *xError.Error

	// UpdateStatus 更新友链分组状态（启用/禁用）
	UpdateStatus(ctx *gin.Context, groupUUID string, req *request.LinkGroupStatusReq) *xError.Error

	// Delete 删除友链分组
	// 返回值：删除冲突的友链列表（如果force=false且存在关联），错误信息
	Delete(ctx *gin.Context, groupUUID string, req *request.LinkGroupDeleteReq) ([]dto.LinkGroupDeleteConflictDTO, *xError.Error)

	// Get 获取友链分组详情
	Get(ctx *gin.Context, groupUUID string) (*dto.LinkGroupDetailDTO, *xError.Error)

	// GetList 获取友链分组列表（不分页）
	GetList(ctx *gin.Context, req *request.LinkGroupListReq) ([]dto.LinkGroupListDTO, *xError.Error)

	// GetPage 获取友链分组分页列表
	GetPage(ctx *gin.Context, req *request.LinkGroupPageReq) (*base.PaginationResponse[dto.LinkGroupNormalDTO], *xError.Error)
}

// NewLinkGroupService 创建友链分组服务实例
func NewLinkGroupService() ILinkGroupService {
	return &logic.LinkGroupLogic{}
}
