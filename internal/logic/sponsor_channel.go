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
	"errors"
	"strconv"

	xError "github.com/bamboo-services/bamboo-base-go/error"
	xCtxUtil "github.com/bamboo-services/bamboo-base-go/utility/ctxutil"
	apiSponsorChannel "github.com/bamboo-services/bamboo-main/api/sponsorchannel"
	entity2 "github.com/bamboo-services/bamboo-main/internal/entity"
	"github.com/bamboo-services/bamboo-main/internal/model/base"
	"github.com/bamboo-services/bamboo-main/internal/model/dto"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// SponsorChannelLogic 赞助渠道业务逻辑
type SponsorChannelLogic struct {
}

func NewSponsorChannelLogic() *SponsorChannelLogic {
	return &SponsorChannelLogic{}
}

// Add 添加赞助渠道
func (l *SponsorChannelLogic) Add(ctx *gin.Context, req *apiSponsorChannel.AddRequest) (*dto.SponsorChannelDetailDTO, *xError.Error) {
	// 获取数据库连接 - 注意：不要再次调用 WithContext，已包含 Snowflake 节点
	db := xCtxUtil.MustGetDB(ctx)

	// 创建赞助渠道实体
	channel := &entity2.SponsorChannel{
		Name:        req.Name,
		Icon:        req.Icon,
		Description: req.Description,
		SortOrder:   req.SortOrder,
		Status:      true, // 默认启用
	}

	// 保存到数据库
	err := db.Create(channel).Error
	if err != nil {
		return nil, xError.NewError(ctx, xError.DatabaseError, "创建赞助渠道失败", false, err)
	}

	// 查询刚创建的渠道详情
	err = db.First(channel, "id = ?", channel.ID).Error
	if err != nil {
		return nil, xError.NewError(ctx, xError.DatabaseError, "查询赞助渠道失败", false, err)
	}

	return convertChannelToDetailDTO(channel, 0), nil
}

// Update 更新赞助渠道
func (l *SponsorChannelLogic) Update(ctx *gin.Context, idStr string, req *apiSponsorChannel.UpdateRequest) (*dto.SponsorChannelDetailDTO, *xError.Error) {
	// 获取数据库连接
	db := xCtxUtil.MustGetDB(ctx)

	// 解析ID
	channelID, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return nil, xError.NewError(ctx, xError.BadRequest, "无效的渠道ID", false)
	}

	// 查找赞助渠道
	var channel entity2.SponsorChannel
	err = db.First(&channel, "id = ?", channelID).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, xError.NewError(ctx, xError.NotFound, "赞助渠道不存在", false)
		}
		return nil, xError.NewError(ctx, xError.DatabaseError, "查询赞助渠道失败", false, err)
	}

	// 更新字段 - 只更新提供的字段
	if req.Name != nil {
		channel.Name = *req.Name
	}
	if req.Icon != nil {
		channel.Icon = req.Icon
	}
	if req.Description != nil {
		channel.Description = req.Description
	}
	if req.SortOrder != nil {
		channel.SortOrder = *req.SortOrder
	}

	// 保存更新
	err = db.Updates(&channel).Error
	if err != nil {
		return nil, xError.NewError(ctx, xError.DatabaseError, "更新赞助渠道失败", false, err)
	}

	// 查询赞助记录数量
	var sponsorCount int64
	err = db.Model(&entity2.SponsorRecord{}).Where("channel_id = ?", channelID).Count(&sponsorCount).Error
	if err != nil {
		return nil, xError.NewError(ctx, xError.DatabaseError, "查询赞助数量失败", false, err)
	}

	return convertChannelToDetailDTO(&channel, int(sponsorCount)), nil
}

// UpdateStatus 更新赞助渠道状态
func (l *SponsorChannelLogic) UpdateStatus(ctx *gin.Context, idStr string, req *apiSponsorChannel.StatusRequest) (bool, *xError.Error) {
	// 获取数据库连接
	db := xCtxUtil.MustGetDB(ctx)

	// 解析ID
	channelID, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return false, xError.NewError(ctx, xError.BadRequest, "无效的渠道ID", false)
	}

	// 更新状态
	result := db.Model(&entity2.SponsorChannel{}).
		Where("id = ?", channelID).
		Update("status", req.Status)

	if result.Error != nil {
		return false, xError.NewError(ctx, xError.DatabaseError, "更新渠道状态失败", false, result.Error)
	}

	if result.RowsAffected == 0 {
		return false, xError.NewError(ctx, xError.NotFound, "赞助渠道不存在", false)
	}

	return req.Status, nil
}

// Delete 删除赞助渠道
func (l *SponsorChannelLogic) Delete(ctx *gin.Context, idStr string) *xError.Error {
	// 获取数据库连接
	db := xCtxUtil.MustGetDB(ctx)

	// 解析ID
	channelID, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return xError.NewError(ctx, xError.BadRequest, "无效的渠道ID", false)
	}

	// 检查渠道是否存在
	var channel entity2.SponsorChannel
	err = db.First(&channel, "id = ?", channelID).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return xError.NewError(ctx, xError.NotFound, "赞助渠道不存在", false)
		}
		return xError.NewError(ctx, xError.DatabaseError, "查询赞助渠道失败", false, err)
	}

	// 查询关联的赞助记录数量
	var sponsorCount int64
	err = db.Model(&entity2.SponsorRecord{}).
		Where("channel_id = ?", channelID).
		Count(&sponsorCount).Error
	if err != nil {
		return xError.NewError(ctx, xError.DatabaseError, "查询关联赞助记录失败", false, err)
	}

	// 如果有关联记录，不允许删除
	if sponsorCount > 0 {
		return xError.NewError(ctx, xError.BadRequest, "该渠道下存在赞助记录，无法删除", false)
	}

	// 删除渠道（硬删除）
	result := db.Unscoped().Where("id = ?", channelID).Delete(&entity2.SponsorChannel{})
	if result.Error != nil {
		return xError.NewError(ctx, xError.DatabaseError, "删除赞助渠道失败", false, result.Error)
	}

	return nil
}

// Get 获取赞助渠道详情
func (l *SponsorChannelLogic) Get(ctx *gin.Context, idStr string) (*dto.SponsorChannelDetailDTO, *xError.Error) {
	// 获取数据库连接
	db := xCtxUtil.MustGetDB(ctx)

	// 解析ID
	channelID, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return nil, xError.NewError(ctx, xError.BadRequest, "无效的渠道ID", false)
	}

	// 查询赞助渠道
	var channel entity2.SponsorChannel
	err = db.First(&channel, "id = ?", channelID).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, xError.NewError(ctx, xError.NotFound, "赞助渠道不存在", false)
		}
		return nil, xError.NewError(ctx, xError.DatabaseError, "查询赞助渠道失败", false, err)
	}

	// 查询赞助记录数量
	var sponsorCount int64
	err = db.Model(&entity2.SponsorRecord{}).Where("channel_id = ?", channelID).Count(&sponsorCount).Error
	if err != nil {
		return nil, xError.NewError(ctx, xError.DatabaseError, "查询赞助数量失败", false, err)
	}

	return convertChannelToDetailDTO(&channel, int(sponsorCount)), nil
}

// GetList 获取赞助渠道列表（不分页）
func (l *SponsorChannelLogic) GetList(ctx *gin.Context, req *apiSponsorChannel.ListRequest) ([]dto.SponsorChannelListDTO, *xError.Error) {
	// 获取数据库连接
	db := xCtxUtil.MustGetDB(ctx)

	// 构建查询
	query := db.Model(&entity2.SponsorChannel{})

	// 应用过滤条件
	if req.Status != nil {
		query = query.Where("status = ?", *req.Status)
	}

	if req.OnlyEnabled != nil && *req.OnlyEnabled {
		query = query.Where("status = ?", true)
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

	// 执行查询
	var channels []entity2.SponsorChannel
	err := query.Find(&channels).Error
	if err != nil {
		return nil, xError.NewError(ctx, xError.DatabaseError, "查询赞助渠道列表失败", false, err)
	}

	// 批量查询赞助记录数量
	channelIDs := make([]int64, len(channels))
	for i, channel := range channels {
		channelIDs[i] = channel.ID
	}

	sponsorCounts := make(map[int64]int64)
	if len(channelIDs) > 0 {
		var countResults []struct {
			ChannelID int64 `gorm:"column:channel_id"`
			Count     int64 `gorm:"column:count"`
		}

		err = db.
			Model(&entity2.SponsorRecord{}).
			Select("channel_id, COUNT(*) as count").
			Where("channel_id IN ?", channelIDs).
			Group("channel_id").
			Find(&countResults).Error
		if err != nil {
			return nil, xError.NewError(ctx, xError.DatabaseError, "查询赞助数量失败", false, err)
		}

		for _, result := range countResults {
			sponsorCounts[result.ChannelID] = result.Count
		}
	}

	// 转换为DTO
	channelDTOs := make([]dto.SponsorChannelListDTO, len(channels))
	for i, channel := range channels {
		channelDTOs[i] = convertChannelToListDTO(&channel, int(sponsorCounts[channel.ID]))
	}

	return channelDTOs, nil
}

// GetPage 获取赞助渠道分页列表
func (l *SponsorChannelLogic) GetPage(ctx *gin.Context, req *apiSponsorChannel.PageRequest) (*base.PaginationResponse[dto.SponsorChannelNormalDTO], *xError.Error) {
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
	query := db.Model(&entity2.SponsorChannel{})

	// 应用过滤条件
	if req.Status != nil {
		query = query.Where("status = ?", *req.Status)
	}

	if req.Name != nil && *req.Name != "" {
		query = query.Where("name ILIKE ?", "%"+*req.Name+"%")
	}

	// 统计总数
	var total int64
	err := query.Count(&total).Error
	if err != nil {
		return nil, xError.NewError(ctx, xError.DatabaseError, "统计赞助渠道数量失败", false, err)
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
	var channels []entity2.SponsorChannel
	offset := (req.Page - 1) * req.PageSize
	err = query.Offset(offset).Limit(req.PageSize).Find(&channels).Error
	if err != nil {
		return nil, xError.NewError(ctx, xError.DatabaseError, "查询赞助渠道列表失败", false, err)
	}

	// 批量查询赞助记录数量
	channelIDs := make([]int64, len(channels))
	for i, channel := range channels {
		channelIDs[i] = channel.ID
	}

	sponsorCounts := make(map[int64]int64)
	if len(channelIDs) > 0 {
		var countResults []struct {
			ChannelID int64 `gorm:"column:channel_id"`
			Count     int64 `gorm:"column:count"`
		}

		err = db.
			Model(&entity2.SponsorRecord{}).
			Select("channel_id, COUNT(*) as count").
			Where("channel_id IN ?", channelIDs).
			Group("channel_id").
			Find(&countResults).Error
		if err != nil {
			return nil, xError.NewError(ctx, xError.DatabaseError, "查询赞助数量失败", false, err)
		}

		for _, result := range countResults {
			sponsorCounts[result.ChannelID] = result.Count
		}
	}

	// 转换为DTO
	channelDTOs := make([]dto.SponsorChannelNormalDTO, len(channels))
	for i, channel := range channels {
		channelDTOs[i] = convertChannelToNormalDTO(&channel, int(sponsorCounts[channel.ID]))
	}

	return base.NewPaginationResponse(channelDTOs, req.Page, req.PageSize, total), nil
}

// GetPublicList 获取公开的赞助渠道列表（仅返回启用状态的渠道）
func (l *SponsorChannelLogic) GetPublicList(ctx *gin.Context) ([]dto.SponsorChannelListDTO, *xError.Error) {
	// 获取数据库连接
	db := xCtxUtil.MustGetDB(ctx)

	// 构建查询：只查询启用状态的渠道
	query := db.Model(&entity2.SponsorChannel{}).
		Where("status = ?", true).
		Order("sort_order asc")

	// 执行查询
	var channels []entity2.SponsorChannel
	err := query.Find(&channels).Error
	if err != nil {
		return nil, xError.NewError(ctx, xError.DatabaseError, "查询公开渠道列表失败", false, err)
	}

	// 批量查询赞助记录数量
	channelIDs := make([]int64, len(channels))
	for i, channel := range channels {
		channelIDs[i] = channel.ID
	}

	sponsorCounts := make(map[int64]int64)
	if len(channelIDs) > 0 {
		var countResults []struct {
			ChannelID int64 `gorm:"column:channel_id"`
			Count     int64 `gorm:"column:count"`
		}

		err = db.
			Model(&entity2.SponsorRecord{}).
			Select("channel_id, COUNT(*) as count").
			Where("channel_id IN ?", channelIDs).
			Group("channel_id").
			Find(&countResults).Error
		if err != nil {
			return nil, xError.NewError(ctx, xError.DatabaseError, "查询赞助数量失败", false, err)
		}

		for _, result := range countResults {
			sponsorCounts[result.ChannelID] = result.Count
		}
	}

	// 转换为DTO
	channelDTOs := make([]dto.SponsorChannelListDTO, len(channels))
	for i, channel := range channels {
		channelDTOs[i] = convertChannelToListDTO(&channel, int(sponsorCounts[channel.ID]))
	}

	return channelDTOs, nil
}

// 辅助函数：将赞助渠道实体转换为详细DTO
func convertChannelToDetailDTO(channel *entity2.SponsorChannel, sponsorCount int) *dto.SponsorChannelDetailDTO {
	if channel == nil {
		return nil
	}

	return &dto.SponsorChannelDetailDTO{
		SponsorChannelNormalDTO: dto.SponsorChannelNormalDTO{
			ID:           channel.ID,
			Name:         channel.Name,
			Icon:         channel.Icon,
			Description:  channel.Description,
			SortOrder:    channel.SortOrder,
			Status:       channel.Status,
			SponsorCount: sponsorCount,
			CreatedAt:    channel.CreatedAt,
			UpdatedAt:    channel.UpdatedAt,
		},
	}
}

// 辅助函数：将赞助渠道实体转换为标准DTO
func convertChannelToNormalDTO(channel *entity2.SponsorChannel, sponsorCount int) dto.SponsorChannelNormalDTO {
	return dto.SponsorChannelNormalDTO{
		ID:           channel.ID,
		Name:         channel.Name,
		Icon:         channel.Icon,
		Description:  channel.Description,
		SortOrder:    channel.SortOrder,
		Status:       channel.Status,
		SponsorCount: sponsorCount,
		CreatedAt:    channel.CreatedAt,
		UpdatedAt:    channel.UpdatedAt,
	}
}

// 辅助函数：将赞助渠道实体转换为列表DTO
func convertChannelToListDTO(channel *entity2.SponsorChannel, sponsorCount int) dto.SponsorChannelListDTO {
	return dto.SponsorChannelListDTO{
		ID:           channel.ID,
		Name:         channel.Name,
		Icon:         channel.Icon,
		SortOrder:    channel.SortOrder,
		Status:       channel.Status,
		SponsorCount: sponsorCount,
	}
}
