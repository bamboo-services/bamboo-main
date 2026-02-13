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

package logic

import (
	"context"
	"errors"
	"strconv"

	xError "github.com/bamboo-services/bamboo-base-go/error"
	xLog "github.com/bamboo-services/bamboo-base-go/log"
	xUtil "github.com/bamboo-services/bamboo-base-go/utility"
	xCtxUtil "github.com/bamboo-services/bamboo-base-go/utility/ctxutil"
	apiLinkColor "github.com/bamboo-services/bamboo-main/api/link"
	entity2 "github.com/bamboo-services/bamboo-main/internal/entity"
	"github.com/bamboo-services/bamboo-main/internal/models/base"
	"github.com/bamboo-services/bamboo-main/internal/models/dto"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// LinkColorLogic 友链颜色业务逻辑
type LinkColorLogic struct {
	logic
}

func NewLinkColorLogic(ctx context.Context) *LinkColorLogic {
	db := xCtxUtil.MustGetDB(ctx)
	rdb := xCtxUtil.MustGetRDB(ctx)

	return &LinkColorLogic{
		logic: logic{
			db:  db,
			rdb: rdb,
			log: xLog.WithName(xLog.NamedLOGC, "LinkColorLogic"),
		},
	}
}

// Add 添加友链颜色
func (l *LinkColorLogic) Add(ctx *gin.Context, req *apiLinkColor.ColorAddRequest) (*dto.LinkColorDetailDTO, *xError.Error) {
	db := xCtxUtil.MustGetDB(ctx)

	// 业务规则校验：type=0 时三个颜色字段必填
	if req.ColorType == 0 {
		if req.PrimaryColor == nil || req.SubColor == nil || req.HoverColor == nil {
			return nil, xError.NewError(ctx, xError.BadRequest, "普通颜色类型需要设置主颜色、副颜色和悬停颜色", false)
		}
		if *req.PrimaryColor == "" || *req.SubColor == "" || *req.HoverColor == "" {
			return nil, xError.NewError(ctx, xError.BadRequest, "普通颜色类型的颜色值不能为空", false)
		}
	}

	// 创建友链颜色实体
	color := &entity2.LinkColor{
		Name:   req.ColorName,
		Type:   req.ColorType,
		Status: true, // 默认启用
	}

	// 根据类型设置颜色值
	if req.ColorType == 0 {
		color.PrimaryColor = req.PrimaryColor
		color.SubColor = req.SubColor
		color.HoverColor = req.HoverColor
	}
	// type=1 时颜色字段保持 nil

	// 设置排序值：查询当前最大排序值并+1
	var maxSort int
	db.Model(&entity2.LinkColor{}).Select("COALESCE(MAX(sort_order), 0)").Scan(&maxSort)
	color.SortOrder = maxSort + 1

	// 保存到数据库
	if err := db.Create(color).Error; err != nil {
		return nil, xError.NewError(ctx, xError.DatabaseError, "创建友链颜色失败", false, err)
	}

	// 预加载关联数据
	if err := db.Preload("LinksFKey").First(color, "id = ?", color.ID).Error; err != nil {
		return nil, xError.NewError(ctx, xError.DatabaseError, "查询友链颜色失败", false, err)
	}

	return convertLinkColorToDetailDTO(color), nil
}

// Update 更新友链颜色
func (l *LinkColorLogic) Update(ctx *gin.Context, colorIDStr string, req *apiLinkColor.ColorUpdateRequest) (*dto.LinkColorDetailDTO, *xError.Error) {
	db := xCtxUtil.MustGetDB(ctx)

	// 解析ID
	colorID, err := strconv.ParseInt(colorIDStr, 10, 64)
	if err != nil {
		return nil, xError.NewError(ctx, xError.BadRequest, "无效的颜色ID", false)
	}

	// 查找友链颜色
	var color entity2.LinkColor
	if err := db.First(&color, "id = ?", colorID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, xError.NewError(ctx, xError.NotFound, "友链颜色不存在", false)
		}
		return nil, xError.NewError(ctx, xError.DatabaseError, "查询友链颜色失败", false, err)
	}

	// 更新字段
	if req.ColorName != nil {
		color.Name = *req.ColorName
	}
	if req.ColorType != nil {
		color.Type = *req.ColorType
	}
	if req.ColorOrder != nil {
		color.SortOrder = *req.ColorOrder
	}

	// 处理颜色字段更新
	if req.PrimaryColor != nil {
		if *req.PrimaryColor == "" {
			color.PrimaryColor = nil
		} else {
			color.PrimaryColor = req.PrimaryColor
		}
	}
	if req.SubColor != nil {
		if *req.SubColor == "" {
			color.SubColor = nil
		} else {
			color.SubColor = req.SubColor
		}
	}
	if req.HoverColor != nil {
		if *req.HoverColor == "" {
			color.HoverColor = nil
		} else {
			color.HoverColor = req.HoverColor
		}
	}

	// 业务规则校验：如果最终 type=0，则三个颜色字段必须有值
	if color.Type == 0 {
		if color.PrimaryColor == nil || color.SubColor == nil || color.HoverColor == nil {
			return nil, xError.NewError(ctx, xError.BadRequest, "普通颜色类型需要设置主颜色、副颜色和悬停颜色", false)
		}
	}

	// 保存更新
	if err := db.Updates(&color).Error; err != nil {
		return nil, xError.NewError(ctx, xError.DatabaseError, "更新友链颜色失败", false, err)
	}

	// 预加载关联数据
	if err := db.Preload("LinksFKey").First(&color, "id = ?", color.ID).Error; err != nil {
		return nil, xError.NewError(ctx, xError.DatabaseError, "查询友链颜色失败", false, err)
	}

	return convertLinkColorToDetailDTO(&color), nil
}

// UpdateSort 批量更新友链颜色排序
func (l *LinkColorLogic) UpdateSort(ctx *gin.Context, req *apiLinkColor.ColorSortRequest) *xError.Error {
	db := xCtxUtil.MustGetDB(ctx)
	colorIDs := req.ColorIDs

	tx := db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	startSort := 0
	if req.SortOrder != nil && *req.SortOrder > 0 {
		startSort = *req.SortOrder
	}

	for i, colorID := range colorIDs {
		result := tx.Model(&entity2.LinkColor{}).
			Where("id = ?", colorID).
			Update("sort_order", startSort+i)

		if result.Error != nil {
			tx.Rollback()
			return xError.NewError(ctx, xError.DatabaseError, "更新颜色排序失败", false, result.Error)
		}

		if result.RowsAffected == 0 {
			tx.Rollback()
			return xError.NewError(ctx, xError.NotFound, "颜色不存在", false)
		}
	}

	if err := tx.Commit().Error; err != nil {
		return xError.NewError(ctx, xError.DatabaseError, "提交排序更新失败", false, err)
	}

	return nil
}

// UpdateStatus 更新友链颜色状态
func (l *LinkColorLogic) UpdateStatus(ctx *gin.Context, colorIDStr string, req *apiLinkColor.ColorStatusRequest) *xError.Error {
	db := xCtxUtil.MustGetDB(ctx)

	colorID, err := strconv.ParseInt(colorIDStr, 10, 64)
	if err != nil {
		return xError.NewError(ctx, xError.BadRequest, "无效的颜色ID", false)
	}

	result := db.Model(&entity2.LinkColor{}).
		Where("id = ?", colorID).
		Update("status", req.Status)

	if result.Error != nil {
		return xError.NewError(ctx, xError.DatabaseError, "更新颜色状态失败", false, result.Error)
	}

	if result.RowsAffected == 0 {
		return xError.NewError(ctx, xError.NotFound, "友链颜色不存在", false)
	}

	return nil
}

// Delete 删除友链颜色
func (l *LinkColorLogic) Delete(ctx *gin.Context, colorIDStr string, req *apiLinkColor.ColorDeleteRequest) ([]dto.LinkColorDeleteConflictDTO, *xError.Error) {
	db := xCtxUtil.MustGetDB(ctx)

	colorID, err := strconv.ParseInt(colorIDStr, 10, 64)
	if err != nil {
		return nil, xError.NewError(ctx, xError.BadRequest, "无效的颜色ID", false)
	}

	// 检查颜色是否存在
	var color entity2.LinkColor
	if err := db.First(&color, "id = ?", colorID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, xError.NewError(ctx, xError.NotFound, "友链颜色不存在", false)
		}
		return nil, xError.NewError(ctx, xError.DatabaseError, "查询友链颜色失败", false, err)
	}

	// 查询关联的友链
	var linkCount int64
	if err := db.Model(&entity2.LinkFriend{}).Where("color_id = ?", colorID).Count(&linkCount).Error; err != nil {
		return nil, xError.NewError(ctx, xError.DatabaseError, "查询关联友链失败", false, err)
	}

	// 如果有关联友链且不是强制删除，返回冲突信息
	if linkCount > 0 && !req.Force {
		var conflictLinks []entity2.LinkFriend
		if err := db.Where("color_id = ?", colorID).Limit(10).Find(&conflictLinks).Error; err != nil {
			return nil, xError.NewError(ctx, xError.DatabaseError, "查询冲突友链失败", false, err)
		}

		conflictDTOs := make([]dto.LinkColorDeleteConflictDTO, len(conflictLinks))
		for i, link := range conflictLinks {
			conflictDTOs[i] = dto.LinkColorDeleteConflictDTO{
				ID:   link.ID,
				Name: link.Name,
				URL:  link.URL,
			}
		}

		return conflictDTOs, xError.NewError(ctx, xError.BadRequest, "颜色下存在友链，无法删除", false)
	}

	tx := db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 如果是强制删除，先清空关联友链的 color_id
	if req.Force && linkCount > 0 {
		if result := tx.Model(&entity2.LinkFriend{}).Where("color_id = ?", colorID).Update("color_id", nil); result.Error != nil {
			tx.Rollback()
			return nil, xError.NewError(ctx, xError.DatabaseError, "清空友链颜色关联失败", false, result.Error)
		}
	}

	// 删除颜色（硬删除）
	if result := tx.Unscoped().Where("id = ?", colorID).Delete(&entity2.LinkColor{}); result.Error != nil {
		tx.Rollback()
		return nil, xError.NewError(ctx, xError.DatabaseError, "删除友链颜色失败", false, result.Error)
	}

	if err := tx.Commit().Error; err != nil {
		return nil, xError.NewError(ctx, xError.DatabaseError, "提交删除操作失败", false, err)
	}

	return nil, nil
}

// Get 获取友链颜色详情
func (l *LinkColorLogic) Get(ctx *gin.Context, colorIDStr string) (*dto.LinkColorDetailDTO, *xError.Error) {
	db := xCtxUtil.MustGetDB(ctx)

	colorID, err := strconv.ParseInt(colorIDStr, 10, 64)
	if err != nil {
		return nil, xError.NewError(ctx, xError.BadRequest, "无效的颜色ID", false)
	}

	var color entity2.LinkColor
	if err := db.Preload("LinksFKey").First(&color, "id = ?", colorID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, xError.NewError(ctx, xError.NotFound, "友链颜色不存在", false)
		}
		return nil, xError.NewError(ctx, xError.DatabaseError, "查询友链颜色失败", false, err)
	}

	return convertLinkColorToDetailDTO(&color), nil
}

// GetList 获取友链颜色列表（不分页）
func (l *LinkColorLogic) GetList(ctx *gin.Context, req *apiLinkColor.ColorListRequest) ([]dto.LinkColorListDTO, *xError.Error) {
	db := xCtxUtil.MustGetDB(ctx)
	query := db.Model(&entity2.LinkColor{})

	// 应用过滤条件
	if req.Status != nil {
		status := *req.Status == 1
		query = query.Where("status = ?", status)
	}
	if req.OnlyEnabled != nil && *req.OnlyEnabled {
		query = query.Where("status = ?", true)
	}
	if req.Type != nil {
		query = query.Where("type = ?", *req.Type)
	}
	if req.Name != nil && *req.Name != "" {
		query = query.Where("name ILIKE ?", "%"+*req.Name+"%")
	}

	// 设置排序
	orderBy := "sort_order"
	if req.OrderBy != nil && *req.OrderBy != "" {
		orderBy = *req.OrderBy
	}
	order := "asc"
	if req.Order != nil && *req.Order != "" {
		order = *req.Order
	}
	query = query.Order(orderBy + " " + order)

	var colors []entity2.LinkColor
	if err := query.Find(&colors).Error; err != nil {
		return nil, xError.NewError(ctx, xError.DatabaseError, "查询友链颜色列表失败", false, err)
	}

	// 批量查询友链数量
	colorIDs := make([]int64, len(colors))
	for i, color := range colors {
		colorIDs[i] = color.ID
	}

	linkCounts := make(map[int64]int64)
	if len(colorIDs) > 0 {
		var countResults []struct {
			ColorID int64 `gorm:"column:color_id"`
			Count   int64 `gorm:"column:count"`
		}

		if err := db.Model(&entity2.LinkFriend{}).
			Select("color_id, COUNT(*) as count").
			Where("color_id IN ?", colorIDs).
			Group("color_id").
			Find(&countResults).Error; err != nil {
			return nil, xError.NewError(ctx, xError.DatabaseError, "查询友链数量失败", false, err)
		}

		for _, result := range countResults {
			linkCounts[result.ColorID] = result.Count
		}
	}

	// 转换为DTO
	colorDTOs := make([]dto.LinkColorListDTO, len(colors))
	for i, color := range colors {
		colorDTOs[i] = convertLinkColorToListDTO(&color, linkCounts[color.ID])
	}

	return colorDTOs, nil
}

// GetPage 获取友链颜色分页列表
func (l *LinkColorLogic) GetPage(ctx *gin.Context, req *apiLinkColor.ColorPageRequest) (*base.PaginationResponse[dto.LinkColorNormalDTO], *xError.Error) {
	db := xCtxUtil.MustGetDB(ctx)

	// 设置默认值
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PageSize <= 0 || req.PageSize > 100 {
		req.PageSize = 10
	}

	query := db.Model(&entity2.LinkColor{})

	// 应用过滤条件
	if req.Status != nil {
		status := *req.Status == 1
		query = query.Where("status = ?", status)
	}
	if req.Type != nil {
		query = query.Where("type = ?", *req.Type)
	}
	if req.Name != nil && *req.Name != "" {
		query = query.Where("name ILIKE ?", "%"+*req.Name+"%")
	}

	// 统计总数
	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, xError.NewError(ctx, xError.DatabaseError, "统计友链颜色数量失败", false, err)
	}

	// 设置排序
	orderBy := "sort_order"
	if req.OrderBy != nil && *req.OrderBy != "" {
		orderBy = *req.OrderBy
	}
	order := "asc"
	if req.Order != nil && *req.Order != "" {
		order = *req.Order
	}
	query = query.Order(orderBy + " " + order)

	// 分页查询
	var colors []entity2.LinkColor
	offset := (req.Page - 1) * req.PageSize
	if err := query.Offset(offset).Limit(req.PageSize).Find(&colors).Error; err != nil {
		return nil, xError.NewError(ctx, xError.DatabaseError, "查询友链颜色列表失败", false, err)
	}

	// 批量查询友链数量
	colorIDs := make([]int64, len(colors))
	for i, color := range colors {
		colorIDs[i] = color.ID
	}

	linkCounts := make(map[int64]int64)
	if len(colorIDs) > 0 {
		var countResults []struct {
			ColorID int64 `gorm:"column:color_id"`
			Count   int64 `gorm:"column:count"`
		}

		if err := db.Model(&entity2.LinkFriend{}).
			Select("color_id, COUNT(*) as count").
			Where("color_id IN ?", colorIDs).
			Group("color_id").
			Find(&countResults).Error; err != nil {
			return nil, xError.NewError(ctx, xError.DatabaseError, "查询友链数量失败", false, err)
		}

		for _, result := range countResults {
			linkCounts[result.ColorID] = result.Count
		}
	}

	// 转换为DTO
	colorDTOs := make([]dto.LinkColorNormalDTO, len(colors))
	for i, color := range colors {
		colorDTOs[i] = convertLinkColorToNormalDTO(&color, linkCounts[color.ID])
	}

	return base.NewPaginationResponse(colorDTOs, req.Page, req.PageSize, total), nil
}

// ============ 辅助函数 ============

// getColorTypeText 获取颜色类型文本
func getColorTypeText(colorType int) string {
	switch colorType {
	case 0:
		return "普通"
	case 1:
		return "炫彩"
	default:
		return "未知"
	}
}

// convertLinkColorToDetailDTO 将友链颜色实体转换为详细DTO
func convertLinkColorToDetailDTO(color *entity2.LinkColor) *dto.LinkColorDetailDTO {
	if color == nil {
		return nil
	}

	links := make([]dto.LinkFriendSimpleDTO, len(color.LinksFKey))
	for i, link := range color.LinksFKey {
		if link != nil {
			links[i] = dto.LinkFriendSimpleDTO{
				ID:     link.ID,
				Name:   link.Name,
				URL:    link.URL,
				Avatar: xUtil.Val(link.Avatar),
			}
		}
	}

	status := 0
	if color.Status {
		status = 1
	}

	return &dto.LinkColorDetailDTO{
		ID:           color.ID,
		Name:         color.Name,
		Type:         color.Type,
		TypeText:     getColorTypeText(color.Type),
		PrimaryColor: color.PrimaryColor,
		SubColor:     color.SubColor,
		HoverColor:   color.HoverColor,
		SortOrder:    color.SortOrder,
		Status:       status,
		LinkCount:    len(links),
		CreatedAt:    color.CreatedAt,
		UpdatedAt:    color.UpdatedAt,
		Links:        links,
	}
}

// convertLinkColorToNormalDTO 将友链颜色实体转换为标准DTO
func convertLinkColorToNormalDTO(color *entity2.LinkColor, linkCount int64) dto.LinkColorNormalDTO {
	status := 0
	if color.Status {
		status = 1
	}

	return dto.LinkColorNormalDTO{
		ID:           color.ID,
		Name:         color.Name,
		Type:         color.Type,
		TypeText:     getColorTypeText(color.Type),
		PrimaryColor: color.PrimaryColor,
		SubColor:     color.SubColor,
		HoverColor:   color.HoverColor,
		SortOrder:    color.SortOrder,
		Status:       status,
		LinkCount:    int(linkCount),
		CreatedAt:    color.CreatedAt,
		UpdatedAt:    color.UpdatedAt,
	}
}

// convertLinkColorToListDTO 将友链颜色实体转换为列表DTO
func convertLinkColorToListDTO(color *entity2.LinkColor, linkCount int64) dto.LinkColorListDTO {
	status := 0
	if color.Status {
		status = 1
	}

	return dto.LinkColorListDTO{
		ID:           color.ID,
		Name:         color.Name,
		Type:         color.Type,
		PrimaryColor: color.PrimaryColor,
		SubColor:     color.SubColor,
		HoverColor:   color.HoverColor,
		SortOrder:    color.SortOrder,
		Status:       status,
		LinkCount:    int(linkCount),
	}
}
