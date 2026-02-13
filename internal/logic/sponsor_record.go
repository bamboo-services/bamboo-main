/*
 * --------------------------------------------------------------------------------
 * Copyright (c) 2016-NOW(至今) 筱锋
 * Author: 筱锋「xiao_lfeng」(https://www.x-lf.com)
 * --------------------------------------------------------------------------------
 * 许可证声明:版权所有 (c) 2016-2025 筱锋。保留所有权利。
 * 有关MIT许可证的更多信息,请查看项目根目录下的LICENSE文件或访问:
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
	xCtxUtil "github.com/bamboo-services/bamboo-base-go/utility/ctxutil"
	apiSponsorRecord "github.com/bamboo-services/bamboo-main/api/sponsor"
	entity2 "github.com/bamboo-services/bamboo-main/internal/entity"
	"github.com/bamboo-services/bamboo-main/internal/models/base"
	"github.com/bamboo-services/bamboo-main/internal/models/dto"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// SponsorRecordLogic 赞助记录业务逻辑
type SponsorRecordLogic struct {
	logic
}

func NewSponsorRecordLogic(ctx context.Context) *SponsorRecordLogic {
	db := xCtxUtil.MustGetDB(ctx)
	rdb := xCtxUtil.MustGetRDB(ctx)

	return &SponsorRecordLogic{
		logic: logic{
			db:  db,
			rdb: rdb,
			log: xLog.WithName(xLog.NamedLOGC, "SponsorRecordLogic"),
		},
	}
}

// Add 添加赞助记录
func (l *SponsorRecordLogic) Add(ctx *gin.Context, req *apiSponsorRecord.RecordAddRequest) (*dto.SponsorRecordDetailDTO, *xError.Error) {
	// 获取数据库连接
	db := xCtxUtil.MustGetDB(ctx)

	// 如果提供了渠道ID,需要验证渠道是否存在
	if req.ChannelID != nil {
		var channel entity2.SponsorChannel
		err := db.First(&channel, "id = ?", *req.ChannelID).Error
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, xError.NewError(ctx, xError.NotFound, "赞助渠道不存在", false)
			}
			return nil, xError.NewError(ctx, xError.DatabaseError, "查询赞助渠道失败", false, err)
		}
	}

	// 创建赞助记录实体
	record := &entity2.SponsorRecord{
		Nickname:    req.Nickname,
		RedirectURL: req.RedirectURL,
		Amount:      req.Amount,
		ChannelID:   req.ChannelID,
		Message:     req.Message,
		SponsorAt:   req.SponsorAt,
		SortOrder:   req.SortOrder,
		IsAnonymous: req.IsAnonymous,
		IsHidden:    req.IsHidden,
	}

	// 保存到数据库
	err := db.Create(record).Error
	if err != nil {
		return nil, xError.NewError(ctx, xError.DatabaseError, "创建赞助记录失败", false, err)
	}

	// 预加载关联渠道数据
	err = db.Preload("ChannelFKey").First(record, "id = ?", record.ID).Error
	if err != nil {
		return nil, xError.NewError(ctx, xError.DatabaseError, "查询赞助记录失败", false, err)
	}

	return convertRecordToDetailDTO(record), nil
}

// Update 更新赞助记录
func (l *SponsorRecordLogic) Update(ctx *gin.Context, idStr string, req *apiSponsorRecord.RecordUpdateRequest) (*dto.SponsorRecordDetailDTO, *xError.Error) {
	// 获取数据库连接
	db := xCtxUtil.MustGetDB(ctx)

	// 解析ID
	recordID, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return nil, xError.NewError(ctx, xError.BadRequest, "无效的记录ID", false)
	}

	// 查找赞助记录
	var record entity2.SponsorRecord
	err = db.First(&record, "id = ?", recordID).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, xError.NewError(ctx, xError.NotFound, "赞助记录不存在", false)
		}
		return nil, xError.NewError(ctx, xError.DatabaseError, "查询赞助记录失败", false, err)
	}

	// 如果更新了渠道ID,需要验证渠道是否存在
	if req.ChannelID != nil {
		var channel entity2.SponsorChannel
		err = db.First(&channel, "id = ?", *req.ChannelID).Error
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, xError.NewError(ctx, xError.NotFound, "赞助渠道不存在", false)
			}
			return nil, xError.NewError(ctx, xError.DatabaseError, "查询赞助渠道失败", false, err)
		}
	}

	// 更新字段
	if req.Nickname != nil {
		record.Nickname = *req.Nickname
	}
	if req.RedirectURL != nil {
		record.RedirectURL = req.RedirectURL
	}
	if req.Amount != nil {
		record.Amount = *req.Amount
	}
	if req.ChannelID != nil {
		record.ChannelID = req.ChannelID
	}
	if req.Message != nil {
		record.Message = req.Message
	}
	if req.SponsorAt != nil {
		record.SponsorAt = req.SponsorAt
	}
	if req.SortOrder != nil {
		record.SortOrder = *req.SortOrder
	}
	if req.IsAnonymous != nil {
		record.IsAnonymous = *req.IsAnonymous
	}
	if req.IsHidden != nil {
		record.IsHidden = *req.IsHidden
	}

	// 保存更新
	err = db.Updates(&record).Error
	if err != nil {
		return nil, xError.NewError(ctx, xError.DatabaseError, "更新赞助记录失败", false, err)
	}

	// 预加载关联数据
	err = db.Preload("ChannelFKey").First(&record, "id = ?", record.ID).Error
	if err != nil {
		return nil, xError.NewError(ctx, xError.DatabaseError, "查询赞助记录失败", false, err)
	}

	return convertRecordToDetailDTO(&record), nil
}

// Delete 删除赞助记录
func (l *SponsorRecordLogic) Delete(ctx *gin.Context, idStr string) *xError.Error {
	// 获取数据库连接
	db := xCtxUtil.MustGetDB(ctx)

	// 解析ID
	recordID, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return xError.NewError(ctx, xError.BadRequest, "无效的记录ID", false)
	}

	// 检查记录是否存在
	var record entity2.SponsorRecord
	err = db.First(&record, "id = ?", recordID).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return xError.NewError(ctx, xError.NotFound, "赞助记录不存在", false)
		}
		return xError.NewError(ctx, xError.DatabaseError, "查询赞助记录失败", false, err)
	}

	// 删除记录（硬删除）
	result := db.Unscoped().Where("id = ?", recordID).Delete(&entity2.SponsorRecord{})
	if result.Error != nil {
		return xError.NewError(ctx, xError.DatabaseError, "删除赞助记录失败", false, result.Error)
	}

	return nil
}

// Get 获取赞助记录详情
func (l *SponsorRecordLogic) Get(ctx *gin.Context, idStr string) (*dto.SponsorRecordDetailDTO, *xError.Error) {
	// 获取数据库连接
	db := xCtxUtil.MustGetDB(ctx)

	// 解析ID
	recordID, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return nil, xError.NewError(ctx, xError.BadRequest, "无效的记录ID", false)
	}

	// 查询赞助记录（预加载关联渠道）
	var record entity2.SponsorRecord
	err = db.Preload("ChannelFKey").First(&record, "id = ?", recordID).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, xError.NewError(ctx, xError.NotFound, "赞助记录不存在", false)
		}
		return nil, xError.NewError(ctx, xError.DatabaseError, "查询赞助记录失败", false, err)
	}

	return convertRecordToDetailDTO(&record), nil
}

// GetPage 获取赞助记录分页列表（后台）
func (l *SponsorRecordLogic) GetPage(ctx *gin.Context, req *apiSponsorRecord.RecordPageRequest) (*base.PaginationResponse[dto.SponsorRecordNormalDTO], *xError.Error) {
	// 获取数据库连接
	db := xCtxUtil.MustGetDB(ctx)

	// 设置默认值
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PageSize <= 0 || req.PageSize > 100 {
		req.PageSize = 10
	}

	// 构建查询
	query := db.Model(&entity2.SponsorRecord{})

	// 应用过滤条件
	if req.ChannelID != nil {
		query = query.Where("channel_id = ?", *req.ChannelID)
	}

	if req.Nickname != nil && *req.Nickname != "" {
		query = query.Where("nickname ILIKE ?", "%"+*req.Nickname+"%")
	}

	if req.IsAnonymous != nil {
		query = query.Where("is_anonymous = ?", *req.IsAnonymous)
	}

	if req.IsHidden != nil {
		query = query.Where("is_hidden = ?", *req.IsHidden)
	}

	// 统计总数
	var total int64
	err := query.Count(&total).Error
	if err != nil {
		return nil, xError.NewError(ctx, xError.DatabaseError, "统计赞助记录数量失败", false, err)
	}

	// 设置排序
	orderBy := "sort_order"
	if req.OrderBy != nil && *req.OrderBy != "" {
		orderBy = *req.OrderBy
	}
	order := "desc"
	if req.Order != nil && *req.Order != "" {
		order = *req.Order
	}
	query = query.Order(orderBy + " " + order)

	// 分页查询
	var records []entity2.SponsorRecord
	offset := (req.Page - 1) * req.PageSize
	err = query.Preload("ChannelFKey").Offset(offset).Limit(req.PageSize).Find(&records).Error
	if err != nil {
		return nil, xError.NewError(ctx, xError.DatabaseError, "查询赞助记录列表失败", false, err)
	}

	// 转换为DTO
	recordDTOs := make([]dto.SponsorRecordNormalDTO, len(records))
	for i, record := range records {
		recordDTOs[i] = convertRecordToNormalDTO(&record)
	}

	return base.NewPaginationResponse(recordDTOs, req.Page, req.PageSize, total), nil
}

// GetPublicPage 获取赞助记录公开分页列表（前台）
func (l *SponsorRecordLogic) GetPublicPage(ctx *gin.Context, req *apiSponsorRecord.RecordPublicPageRequest) (*base.PaginationResponse[dto.SponsorRecordSimpleDTO], *xError.Error) {
	// 获取数据库连接
	db := xCtxUtil.MustGetDB(ctx)

	// 设置默认值
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PageSize <= 0 || req.PageSize > 50 {
		req.PageSize = 20
	}

	// 构建查询（只查询未隐藏的记录）
	query := db.Model(&entity2.SponsorRecord{}).Where("is_hidden = ?", false)

	// 应用过滤条件
	if req.ChannelID != nil {
		query = query.Where("channel_id = ?", *req.ChannelID)
	}

	// 统计总数
	var total int64
	err := query.Count(&total).Error
	if err != nil {
		return nil, xError.NewError(ctx, xError.DatabaseError, "统计赞助记录数量失败", false, err)
	}

	// 设置排序
	orderBy := "sort_order"
	if req.OrderBy != nil && *req.OrderBy != "" {
		orderBy = *req.OrderBy
	}
	order := "desc"
	if req.Order != nil && *req.Order != "" {
		order = *req.Order
	}
	query = query.Order(orderBy + " " + order)

	// 分页查询
	var records []entity2.SponsorRecord
	offset := (req.Page - 1) * req.PageSize
	err = query.Preload("ChannelFKey").Offset(offset).Limit(req.PageSize).Find(&records).Error
	if err != nil {
		return nil, xError.NewError(ctx, xError.DatabaseError, "查询赞助记录列表失败", false, err)
	}

	// 转换为DTO（匿名记录特殊处理）
	recordDTOs := make([]dto.SponsorRecordSimpleDTO, len(records))
	for i, record := range records {
		recordDTOs[i] = convertRecordToSimpleDTO(&record)
	}

	return base.NewPaginationResponse(recordDTOs, req.Page, req.PageSize, total), nil
}

// 辅助函数：将赞助记录实体转换为详细DTO
func convertRecordToDetailDTO(record *entity2.SponsorRecord) *dto.SponsorRecordDetailDTO {
	if record == nil {
		return nil
	}

	return &dto.SponsorRecordDetailDTO{
		SponsorRecordNormalDTO: convertRecordToNormalDTO(record),
	}
}

// 辅助函数：将赞助记录实体转换为标准DTO
func convertRecordToNormalDTO(record *entity2.SponsorRecord) dto.SponsorRecordNormalDTO {
	normalDTO := dto.SponsorRecordNormalDTO{
		ID:          record.ID,
		Nickname:    record.Nickname,
		RedirectURL: record.RedirectURL,
		Amount:      record.Amount,
		ChannelID:   record.ChannelID,
		Message:     record.Message,
		SponsorAt:   record.SponsorAt,
		SortOrder:   record.SortOrder,
		IsAnonymous: record.IsAnonymous,
		IsHidden:    record.IsHidden,
		CreatedAt:   record.CreatedAt,
		UpdatedAt:   record.UpdatedAt,
	}

	// 转换关联的渠道信息
	if record.ChannelFKey != nil {
		normalDTO.Channel = &dto.SponsorChannelSimpleDTO{
			ID:   record.ChannelFKey.ID,
			Name: record.ChannelFKey.Name,
			Icon: record.ChannelFKey.Icon,
		}
	}

	return normalDTO
}

// 辅助函数：将赞助记录实体转换为简单DTO（公开接口用）
func convertRecordToSimpleDTO(record *entity2.SponsorRecord) dto.SponsorRecordSimpleDTO {
	// 匿名记录特殊处理
	nickname := record.Nickname
	redirectURL := record.RedirectURL
	if record.IsAnonymous {
		nickname = "匿名用户"
		redirectURL = nil
	}

	simpleDTO := dto.SponsorRecordSimpleDTO{
		ID:          record.ID,
		Nickname:    nickname,
		RedirectURL: redirectURL,
		Amount:      record.Amount,
		Message:     record.Message,
		SponsorAt:   record.SponsorAt,
	}

	// 转换关联的渠道信息
	if record.ChannelFKey != nil {
		simpleDTO.Channel = &dto.SponsorChannelSimpleDTO{
			ID:   record.ChannelFKey.ID,
			Name: record.ChannelFKey.Name,
			Icon: record.ChannelFKey.Icon,
		}
	}

	return simpleDTO
}
