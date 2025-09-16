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
	"bamboo-main/internal/model/base"
	"bamboo-main/internal/model/dto"
	"bamboo-main/internal/model/entity"
	"bamboo-main/internal/model/request"

	xError "github.com/bamboo-services/bamboo-base-go/error"
	xUtil "github.com/bamboo-services/bamboo-base-go/utility"
	xCtxUtil "github.com/bamboo-services/bamboo-base-go/utility/ctxutil"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// LinkGroupLogic 友链分组业务逻辑
type LinkGroupLogic struct {
}

// Add 添加友链分组
func (l *LinkGroupLogic) Add(ctx *gin.Context, req *request.LinkGroupAddReq) (*dto.LinkGroupDetailDTO, *xError.Error) {
	// 获取数据库连接
	db := xCtxUtil.GetDB(ctx)

	// 创建友链分组实体
	group := &entity.LinkGroup{
		Name:        req.GroupName,
		Description: &req.GroupDesc,
		Status:      true, // 默认启用
	}

	// 设置排序值：查询当前最大排序值并+1
	var maxSort int
	db.WithContext(ctx).Model(&entity.LinkGroup{}).Select("COALESCE(MAX(sort_order), 0)").Scan(&maxSort)
	group.SortOrder = maxSort + 1

	// 保存到数据库
	err := db.WithContext(ctx).Create(group).Error
	if err != nil {
		return nil, xError.NewError(ctx, xError.DatabaseError, "创建友链分组失败", false, err)
	}

	// 预加载关联数据并统计友链数量
	err = db.WithContext(ctx).Preload("LinksFKey").First(group, "uuid = ?", group.UUID).Error
	if err != nil {
		return nil, xError.NewError(ctx, xError.DatabaseError, "查询友链分组失败", false, err)
	}

	return convertLinkGroupToDetailDTO(group), nil
}

// Update 更新友链分组（名称和描述）
func (l *LinkGroupLogic) Update(ctx *gin.Context, groupUUID string, req *request.LinkGroupUpdateReq) (*dto.LinkGroupDetailDTO, *xError.Error) {
	// 获取数据库连接
	db := xCtxUtil.GetDB(ctx)

	// 解析UUID
	groupID, err := uuid.Parse(groupUUID)
	if err != nil {
		return nil, xError.NewError(ctx, xError.BadRequest, "无效的分组UUID", false)
	}

	// 查找友链分组
	var group entity.LinkGroup
	err = db.WithContext(ctx).First(&group, "uuid = ?", groupID).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, xError.NewError(ctx, xError.NotFound, "友链分组不存在", false)
		}
		return nil, xError.NewError(ctx, xError.DatabaseError, "查询友链分组失败", false, err)
	}

	// 更新字段
	if req.GroupName != "" {
		group.Name = req.GroupName
	}
	if req.GroupDesc != "" {
		group.Description = &req.GroupDesc
	}
	if req.GroupOrder != nil {
		group.SortOrder = *req.GroupOrder
	}
	if req.GroupStatus != nil {
		group.Status = *req.GroupStatus == 1
	}

	// 保存更新
	err = db.WithContext(ctx).Updates(&group).Error
	if err != nil {
		return nil, xError.NewError(ctx, xError.DatabaseError, "更新友链分组失败", false, err)
	}

	// 预加载关联数据
	err = db.WithContext(ctx).Preload("LinksFKey").First(&group, "uuid = ?", group.UUID).Error
	if err != nil {
		return nil, xError.NewError(ctx, xError.DatabaseError, "查询友链分组失败", false, err)
	}

	return convertLinkGroupToDetailDTO(&group), nil
}

// UpdateSort 批量更新友链分组排序
func (l *LinkGroupLogic) UpdateSort(ctx *gin.Context, req *request.LinkGroupSortReq) *xError.Error {
	// 获取数据库连接
	db := xCtxUtil.GetDB(ctx)

	// 解析UUID数组
	groupIDs := make([]uuid.UUID, len(req.GroupUUIDs))
	for i, uuidStr := range req.GroupUUIDs {
		groupID, err := uuid.Parse(uuidStr)
		if err != nil {
			return xError.NewError(ctx, xError.BadRequest, "无效的分组UUID", false)
		}
		groupIDs[i] = groupID
	}

	// 开始事务
	tx := db.WithContext(ctx).Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 计算起始排序值
	startSort := 0
	if req.SortOrder != nil && *req.SortOrder > 0 {
		startSort = *req.SortOrder
	}

	// 按顺序更新每个分组的排序值
	for i, groupID := range groupIDs {
		result := tx.Model(&entity.LinkGroup{}).
			Where("uuid = ?", groupID).
			Update("sort_order", startSort+i)

		if result.Error != nil {
			tx.Rollback()
			return xError.NewError(ctx, xError.DatabaseError, "更新分组排序失败", false, result.Error)
		}

		if result.RowsAffected == 0 {
			tx.Rollback()
			return xError.NewError(ctx, xError.NotFound, "分组不存在", false)
		}
	}

	// 提交事务
	if err := tx.Commit().Error; err != nil {
		return xError.NewError(ctx, xError.DatabaseError, "提交排序更新失败", false, err)
	}

	return nil
}

// UpdateStatus 更新友链分组状态
func (l *LinkGroupLogic) UpdateStatus(ctx *gin.Context, groupUUID string, req *request.LinkGroupStatusReq) *xError.Error {
	// 获取数据库连接
	db := xCtxUtil.GetDB(ctx)

	// 解析UUID
	groupID, err := uuid.Parse(groupUUID)
	if err != nil {
		return xError.NewError(ctx, xError.BadRequest, "无效的分组UUID", false)
	}

	// 更新状态
	result := db.WithContext(ctx).Model(&entity.LinkGroup{}).
		Where("uuid = ?", groupID).
		Update("status", req.Status)

	if result.Error != nil {
		return xError.NewError(ctx, xError.DatabaseError, "更新分组状态失败", false, result.Error)
	}

	if result.RowsAffected == 0 {
		return xError.NewError(ctx, xError.NotFound, "友链分组不存在", false)
	}

	return nil
}

// Delete 删除友链分组
func (l *LinkGroupLogic) Delete(ctx *gin.Context, groupUUID string, req *request.LinkGroupDeleteReq) ([]dto.LinkGroupDeleteConflictDTO, *xError.Error) {
	// 获取数据库连接
	db := xCtxUtil.GetDB(ctx)

	// 解析UUID
	groupID, err := uuid.Parse(groupUUID)
	if err != nil {
		return nil, xError.NewError(ctx, xError.BadRequest, "无效的分组UUID", false)
	}

	// 检查分组是否存在
	var group entity.LinkGroup
	err = db.WithContext(ctx).First(&group, "uuid = ?", groupID).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, xError.NewError(ctx, xError.NotFound, "友链分组不存在", false)
		}
		return nil, xError.NewError(ctx, xError.DatabaseError, "查询友链分组失败", false, err)
	}

	// 查询关联的友链
	var linkCount int64
	err = db.WithContext(ctx).Model(&entity.LinkFriend{}).
		Where("group_uuid = ?", groupID).
		Count(&linkCount).Error
	if err != nil {
		return nil, xError.NewError(ctx, xError.DatabaseError, "查询关联友链失败", false, err)
	}

	// 如果有关联友链且不是强制删除，返回冲突信息
	if linkCount > 0 && !req.Force {
		var conflictLinks []entity.LinkFriend
		err = db.WithContext(ctx).
			Where("group_uuid = ?", groupID).
			Limit(10).
			Find(&conflictLinks).Error
		if err != nil {
			return nil, xError.NewError(ctx, xError.DatabaseError, "查询冲突友链失败", false, err)
		}

		// 转换为冲突DTO
		conflictDTOs := make([]dto.LinkGroupDeleteConflictDTO, len(conflictLinks))
		for i, link := range conflictLinks {
			conflictDTOs[i] = dto.LinkGroupDeleteConflictDTO{
				UUID: link.UUID.String(),
				Name: link.Name,
				URL:  link.URL,
			}
		}

		return conflictDTOs, xError.NewError(ctx, xError.BadRequest, "分组下存在友链，无法删除", false)
	}

	// 开始事务
	tx := db.WithContext(ctx).Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 如果是强制删除，先清空关联友链的group_uuid
	if req.Force && linkCount > 0 {
		result := tx.Model(&entity.LinkFriend{}).
			Where("group_uuid = ?", groupID).
			Update("group_uuid", nil)
		if result.Error != nil {
			tx.Rollback()
			return nil, xError.NewError(ctx, xError.DatabaseError, "清空友链分组关联失败", false, result.Error)
		}
	}

	// 删除分组（硬删除）
	result := tx.Unscoped().Where("uuid = ?", groupID).Delete(&entity.LinkGroup{})
	if result.Error != nil {
		tx.Rollback()
		return nil, xError.NewError(ctx, xError.DatabaseError, "删除友链分组失败", false, result.Error)
	}

	// 提交事务
	if err := tx.Commit().Error; err != nil {
		return nil, xError.NewError(ctx, xError.DatabaseError, "提交删除操作失败", false, err)
	}

	return nil, nil
}

// Get 获取友链分组详情
func (l *LinkGroupLogic) Get(ctx *gin.Context, groupUUID string) (*dto.LinkGroupDetailDTO, *xError.Error) {
	// 获取数据库连接
	db := xCtxUtil.GetDB(ctx)

	// 解析UUID
	groupID, err := uuid.Parse(groupUUID)
	if err != nil {
		return nil, xError.NewError(ctx, xError.BadRequest, "无效的分组UUID", false)
	}

	// 查询友链分组（预加载关联友链）
	var group entity.LinkGroup
	err = db.WithContext(ctx).Preload("LinksFKey").First(&group, "uuid = ?", groupID).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, xError.NewError(ctx, xError.NotFound, "友链分组不存在", false)
		}
		return nil, xError.NewError(ctx, xError.DatabaseError, "查询友链分组失败", false, err)
	}

	return convertLinkGroupToDetailDTO(&group), nil
}

// GetList 获取友链分组列表（不分页）
func (l *LinkGroupLogic) GetList(ctx *gin.Context, req *request.LinkGroupListReq) ([]dto.LinkGroupListDTO, *xError.Error) {
	// 获取数据库连接
	db := xCtxUtil.GetDB(ctx)

	// 构建查询
	query := db.WithContext(ctx).Model(&entity.LinkGroup{})

	// 应用过滤条件
	if req.Status != nil {
		status := *req.Status == 1
		query = query.Where("status = ?", status)
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
	var groups []entity.LinkGroup
	err := query.Find(&groups).Error
	if err != nil {
		return nil, xError.NewError(ctx, xError.DatabaseError, "查询友链分组列表失败", false, err)
	}

	// 批量查询友链数量
	groupUUIDs := make([]uuid.UUID, len(groups))
	for i, group := range groups {
		groupUUIDs[i] = group.UUID
	}

	linkCounts := make(map[uuid.UUID]int64)
	if len(groupUUIDs) > 0 {
		var countResults []struct {
			GroupUUID uuid.UUID `gorm:"column:group_uuid"`
			Count     int64     `gorm:"column:count"`
		}

		err = db.WithContext(ctx).
			Model(&entity.LinkFriend{}).
			Select("group_uuid, COUNT(*) as count").
			Where("group_uuid IN ?", groupUUIDs).
			Group("group_uuid").
			Find(&countResults).Error
		if err != nil {
			return nil, xError.NewError(ctx, xError.DatabaseError, "查询友链数量失败", false, err)
		}

		for _, result := range countResults {
			linkCounts[result.GroupUUID] = result.Count
		}
	}

	// 转换为DTO
	groupDTOs := make([]dto.LinkGroupListDTO, len(groups))
	for i, group := range groups {
		groupDTOs[i] = convertLinkGroupToListDTO(&group, linkCounts[group.UUID])
	}

	return groupDTOs, nil
}

// GetPage 获取友链分组分页列表
func (l *LinkGroupLogic) GetPage(ctx *gin.Context, req *request.LinkGroupPageReq) (*base.PaginationResponse[dto.LinkGroupNormalDTO], *xError.Error) {
	// 获取数据库连接
	db := xCtxUtil.GetDB(ctx)

	// 设置默认值
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PageSize <= 0 || req.PageSize > 100 {
		req.PageSize = 10
	}

	// 构建查询
	query := db.WithContext(ctx).Model(&entity.LinkGroup{})

	// 应用过滤条件
	if req.Status != nil {
		status := *req.Status == 1
		query = query.Where("status = ?", status)
	}

	if req.Name != nil && *req.Name != "" {
		query = query.Where("name ILIKE ?", "%"+*req.Name+"%")
	}

	// 统计总数
	var total int64
	err := query.Count(&total).Error
	if err != nil {
		return nil, xError.NewError(ctx, xError.DatabaseError, "统计友链分组数量失败", false, err)
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
	var groups []entity.LinkGroup
	offset := (req.Page - 1) * req.PageSize
	err = query.Offset(offset).Limit(req.PageSize).Find(&groups).Error
	if err != nil {
		return nil, xError.NewError(ctx, xError.DatabaseError, "查询友链分组列表失败", false, err)
	}

	// 批量查询友链数量
	groupUUIDs := make([]uuid.UUID, len(groups))
	for i, group := range groups {
		groupUUIDs[i] = group.UUID
	}

	linkCounts := make(map[uuid.UUID]int64)
	if len(groupUUIDs) > 0 {
		var countResults []struct {
			GroupUUID uuid.UUID `gorm:"column:group_uuid"`
			Count     int64     `gorm:"column:count"`
		}

		err = db.WithContext(ctx).
			Model(&entity.LinkFriend{}).
			Select("group_uuid, COUNT(*) as count").
			Where("group_uuid IN ?", groupUUIDs).
			Group("group_uuid").
			Find(&countResults).Error
		if err != nil {
			return nil, xError.NewError(ctx, xError.DatabaseError, "查询友链数量失败", false, err)
		}

		for _, result := range countResults {
			linkCounts[result.GroupUUID] = result.Count
		}
	}

	// 转换为DTO
	groupDTOs := make([]dto.LinkGroupNormalDTO, len(groups))
	for i, group := range groups {
		groupDTOs[i] = convertLinkGroupToNormalDTO(&group, linkCounts[group.UUID])
	}

	return base.NewPaginationResponse(groupDTOs, req.Page, req.PageSize, total), nil
}

// 辅助函数：将友链分组实体转换为详细DTO
func convertLinkGroupToDetailDTO(group *entity.LinkGroup) *dto.LinkGroupDetailDTO {
	if group == nil {
		return nil
	}

	// 转换关联的友链信息
	links := make([]dto.LinkFriendBasicInfo, len(group.LinksFKey))
	for i, link := range group.LinksFKey {
		if link != nil {
			links[i] = dto.LinkFriendBasicInfo{
				UUID:   link.UUID.String(),
				Name:   link.Name,
				URL:    link.URL,
				Avatar: xUtil.Val(link.Avatar),
			}
		}
	}

	// 构建状态值（true->1, false->0）
	status := 0
	if group.Status {
		status = 1
	}

	return &dto.LinkGroupDetailDTO{
		UUID:        group.UUID.String(),
		Name:        group.Name,
		Description: xUtil.Val(group.Description),
		SortOrder:   group.SortOrder,
		Status:      status,
		LinkCount:   len(links),
		CreatedAt:   group.CreatedAt,
		UpdatedAt:   group.UpdatedAt,
		Links:       links,
	}
}

// 辅助函数：将友链分组实体转换为标准DTO
func convertLinkGroupToNormalDTO(group *entity.LinkGroup, linkCount int64) dto.LinkGroupNormalDTO {
	// 构建状态值（true->1, false->0）
	status := 0
	if group.Status {
		status = 1
	}

	return dto.LinkGroupNormalDTO{
		UUID:        group.UUID.String(),
		Name:        group.Name,
		Description: xUtil.Val(group.Description),
		SortOrder:   group.SortOrder,
		Status:      status,
		LinkCount:   int(linkCount),
		CreatedAt:   group.CreatedAt,
		UpdatedAt:   group.UpdatedAt,
	}
}

// 辅助函数：将友链分组实体转换为列表DTO
func convertLinkGroupToListDTO(group *entity.LinkGroup, linkCount int64) dto.LinkGroupListDTO {
	// 构建状态值（true->1, false->0）
	status := 0
	if group.Status {
		status = 1
	}

	return dto.LinkGroupListDTO{
		UUID:      group.UUID.String(),
		Name:      group.Name,
		SortOrder: group.SortOrder,
		Status:    status,
		LinkCount: int(linkCount),
	}
}
