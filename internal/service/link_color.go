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

// ILinkColorService 友链颜色服务接口
type ILinkColorService interface {
	// Add 添加友链颜色
	Add(ctx *gin.Context, req *request.LinkColorAddReq) (*dto.LinkColorDetailDTO, *xError.Error)

	// Update 更新友链颜色
	Update(ctx *gin.Context, colorIDStr string, req *request.LinkColorUpdateReq) (*dto.LinkColorDetailDTO, *xError.Error)

	// UpdateSort 批量更新友链颜色排序
	UpdateSort(ctx *gin.Context, req *request.LinkColorSortReq) *xError.Error

	// UpdateStatus 更新友链颜色状态（启用/禁用）
	UpdateStatus(ctx *gin.Context, colorIDStr string, req *request.LinkColorStatusReq) *xError.Error

	// Delete 删除友链颜色
	// 返回值：删除冲突的友链列表（如果force=false且存在关联），错误信息
	Delete(ctx *gin.Context, colorIDStr string, req *request.LinkColorDeleteReq) ([]dto.LinkColorDeleteConflictDTO, *xError.Error)

	// Get 获取友链颜色详情
	Get(ctx *gin.Context, colorIDStr string) (*dto.LinkColorDetailDTO, *xError.Error)

	// GetList 获取友链颜色列表（不分页）
	GetList(ctx *gin.Context, req *request.LinkColorListReq) ([]dto.LinkColorListDTO, *xError.Error)

	// GetPage 获取友链颜色分页列表
	GetPage(ctx *gin.Context, req *request.LinkColorPageReq) (*base.PaginationResponse[dto.LinkColorNormalDTO], *xError.Error)
}

// NewLinkColorService 创建友链颜色服务实例
func NewLinkColorService() ILinkColorService {
	return &logic.LinkColorLogic{}
}
