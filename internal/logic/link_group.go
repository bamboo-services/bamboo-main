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

package logic

import (
	"context"
	"strings"

	xError "github.com/bamboo-services/bamboo-base-go/common/error"
	xLog "github.com/bamboo-services/bamboo-base-go/common/log"
	xSnowflake "github.com/bamboo-services/bamboo-base-go/common/snowflake"
	xCtxUtil "github.com/bamboo-services/bamboo-base-go/major/utility/context"
	apiLinkGroup "github.com/bamboo-services/bamboo-main/api/link"
	"github.com/bamboo-services/bamboo-main/internal/entity"
	"github.com/bamboo-services/bamboo-main/internal/models/base"
	"github.com/bamboo-services/bamboo-main/internal/repository"
	"gorm.io/gorm"
)

type linkGroupRepo struct {
	group *repository.LinkGroupRepo
	link  *repository.LinkRepo
}

// LinkGroupLogic 友链分组业务逻辑
type LinkGroupLogic struct {
	logic
	repo linkGroupRepo
}

// NewLinkGroupLogic 创建 LinkGroupLogic 实例，从上下文获取数据库与缓存并初始化分组与友链仓储依赖。
func NewLinkGroupLogic(ctx context.Context) *LinkGroupLogic {
	db := xCtxUtil.MustGetDB(ctx)
	m := xCtxUtil.MustGetCacheManager(ctx)

	return &LinkGroupLogic{
		logic: logic{
			db:    db,
			cache: m,
			log:   xLog.WithName(xLog.NamedLOGC, "LinkGroupLogic"),
		},
		repo: linkGroupRepo{
			group: repository.NewLinkGroupRepo(db, m),
			link:  repository.NewLinkRepo(db, m),
		},
	}
}

// Add 添加友链分组，按当前最大排序值自增生成新分组。
func (l *LinkGroupLogic) Add(ctx context.Context, req *apiLinkGroup.GroupAddRequest) (*entity.LinkGroup, *xError.Error) {
	maxSort, xErr := l.repo.group.GetMaxSortOrder(ctx, nil)
	if xErr != nil {
		return nil, xErr
	}

	sortOrder := maxSort + 1
	// 排序位 0、1 预留给内置分组（首页/友链页），自定义分组排序值从内置数量起
	if sortOrder < len(entity.NewBuiltinGroups()) {
		sortOrder = len(entity.NewBuiltinGroups())
	}

	group := &entity.LinkGroup{
		Name:        req.GroupName,
		Description: &req.GroupDesc,
		SortOrder:   sortOrder,
		Status:      true,
	}

	_, xErr = l.repo.group.Create(ctx, group, nil)
	if xErr != nil {
		return nil, xErr
	}

	reloaded, found, xErr := l.repo.group.GetByID(ctx, group.ID, true, nil)
	if xErr != nil {
		return nil, xErr
	}
	if !found {
		return nil, xError.NewError(ctx, xError.NotFound, "友链分组不存在", false)
	}

	return reloaded, nil
}

// Update 更新友链分组，按请求字段覆盖名称、描述、排序与状态。
func (l *LinkGroupLogic) Update(ctx context.Context, groupID xSnowflake.SnowflakeID, req *apiLinkGroup.GroupUpdateRequest) (*entity.LinkGroup, *xError.Error) {
	// 内置分组为系统预设位置（不落库），拒绝编辑
	if entity.IsBuiltinGroupID(groupID) {
		return nil, xError.NewError(ctx, xError.BadRequest, "内置分组为系统预设位置，不支持编辑", false)
	}

	group, found, xErr := l.repo.group.GetByID(ctx, groupID, false, nil)
	if xErr != nil {
		return nil, xErr
	}
	if !found {
		return nil, xError.NewError(ctx, xError.NotFound, "友链分组不存在", false)
	}

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

	_, xErr = l.repo.group.Save(ctx, group, nil)
	if xErr != nil {
		return nil, xErr
	}

	reloaded, found, xErr := l.repo.group.GetByID(ctx, groupID, true, nil)
	if xErr != nil {
		return nil, xErr
	}
	if !found {
		return nil, xError.NewError(ctx, xError.NotFound, "友链分组不存在", false)
	}

	return reloaded, nil
}

// UpdateSort 更新友链分组排序，在事务内按 ID 列表批量重置排序值。
func (l *LinkGroupLogic) UpdateSort(ctx context.Context, req *apiLinkGroup.GroupSortRequest) *xError.Error {
	startSort := 0
	if req.SortOrder != nil && *req.SortOrder > 0 {
		startSort = *req.SortOrder
	}

	// 内置分组不参与排序持久化：剔除保留 ID，真实分组排序值顺延内置占位（0、1 预留给内置）
	validIDs := make([]xSnowflake.SnowflakeID, 0, len(req.GroupIDs))
	for _, id := range req.GroupIDs {
		if !entity.IsBuiltinGroupID(id) {
			validIDs = append(validIDs, id)
		}
	}
	if len(validIDs) == 0 {
		return nil
	}
	startSort += len(req.GroupIDs) - len(validIDs)

	return l.withTx(ctx, func(tx *gorm.DB) *xError.Error {
		return l.repo.group.UpdateSortByIDs(ctx, validIDs, startSort, tx)
	})
}

// UpdateStatus 更新友链分组启用状态。
func (l *LinkGroupLogic) UpdateStatus(ctx context.Context, groupID xSnowflake.SnowflakeID, req *apiLinkGroup.GroupStatusRequest) *xError.Error {
	// 内置分组恒为启用状态（不落库），拒绝切换
	if entity.IsBuiltinGroupID(groupID) {
		return xError.NewError(ctx, xError.BadRequest, "内置分组为系统预设位置，恒为启用", false)
	}

	ok, xErr := l.repo.group.UpdateStatusByID(ctx, groupID, req.Status, nil)
	if xErr != nil {
		return xErr
	}
	if !ok {
		return xError.NewError(ctx, xError.NotFound, "友链分组不存在", false)
	}

	return nil
}

// Delete 删除友链分组，存在关联友链时按 Force 决定阻断或事务清空外键后删除。
func (l *LinkGroupLogic) Delete(ctx context.Context, groupID xSnowflake.SnowflakeID, req *apiLinkGroup.GroupDeleteRequest) ([]entity.LinkFriend, *xError.Error) {
	// 内置分组为系统预设位置（不落库），拒绝删除
	if entity.IsBuiltinGroupID(groupID) {
		return nil, xError.NewError(ctx, xError.BadRequest, "内置分组为系统预设位置，不可删除", false)
	}

	_, found, xErr := l.repo.group.GetByID(ctx, groupID, false, nil)
	if xErr != nil {
		return nil, xErr
	}
	if !found {
		return nil, xError.NewError(ctx, xError.NotFound, "友链分组不存在", false)
	}

	linkCount, xErr := l.repo.link.CountByGroupID(ctx, groupID, nil)
	if xErr != nil {
		return nil, xErr
	}

	if linkCount > 0 && !req.Force {
		conflictLinks, xErr := l.repo.link.ListByGroupID(ctx, groupID, 10, nil)
		if xErr != nil {
			return nil, xErr
		}
		return conflictLinks, xError.NewError(ctx, xError.BadRequest, "分组下存在友链，无法删除", false)
	}

	tx := l.db.WithContext(ctx).Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if req.Force && linkCount > 0 {
		if xErr = l.repo.link.ClearGroupID(ctx, groupID, tx); xErr != nil {
			tx.Rollback()
			return nil, xErr
		}
	}

	ok, xErr := l.repo.group.DeleteByID(ctx, groupID, tx)
	if xErr != nil {
		tx.Rollback()
		return nil, xErr
	}
	if !ok {
		tx.Rollback()
		return nil, xError.NewError(ctx, xError.NotFound, "友链分组不存在", false)
	}

	if err := tx.Commit().Error; err != nil {
		return nil, xError.NewError(ctx, xError.DatabaseError, "提交删除操作失败", false, err)
	}

	return nil, nil
}

// Get 获取友链分组详情。
func (l *LinkGroupLogic) Get(ctx context.Context, groupID xSnowflake.SnowflakeID) (*entity.LinkGroup, *xError.Error) {
	group, found, xErr := l.repo.group.GetByID(ctx, groupID, true, nil)
	if xErr != nil {
		return nil, xErr
	}
	if !found {
		return nil, xError.NewError(ctx, xError.NotFound, "友链分组不存在", false)
	}

	return group, nil
}

// GetList 获取友链分组列表。
func (l *LinkGroupLogic) GetList(ctx context.Context, req *apiLinkGroup.GroupListRequest) ([]entity.LinkGroup, *xError.Error) {
	groups, xErr := l.repo.group.List(ctx, &repository.GroupListQuery{
		Status:      req.Status,
		Name:        req.Name,
		WithLinks:   req.WithLinks,
		OnlyEnabled: req.OnlyEnabled,
		OrderBy:     req.OrderBy,
		Order:       req.Order,
	}, nil)
	if xErr != nil {
		return nil, xErr
	}

	// 内置分组（不落库）：满足过滤条件时置顶注入虚拟记录，供前端选择器使用
	if builtinGroupsVisible(req.Status, req.Name, req.OnlyEnabled) {
		groups = append(builtinGroupValues(), groups...)
	}

	return groups, nil
}

// GetPage 分页获取友链分组。
func (l *LinkGroupLogic) GetPage(ctx context.Context, req *apiLinkGroup.GroupPageRequest) (*base.PaginationResponse[entity.LinkGroup], *xError.Error) {
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PageSize <= 0 || req.PageSize > 100 {
		req.PageSize = 10
	}

	groups, total, xErr := l.repo.group.Page(ctx, &repository.GroupPageQuery{
		Page:     req.Page,
		PageSize: req.PageSize,
		Status:   req.Status,
		Name:     req.Name,
		OrderBy:  req.OrderBy,
		Order:    req.Order,
	}, nil)
	if xErr != nil {
		return nil, xErr
	}

	// 内置分组置顶注入（仅第 1 页；total 计入内置条数，保证分页一致）
	if req.Page == 1 && builtinGroupsVisible(req.Status, req.Name, nil) {
		values := builtinGroupValues()
		groups = append(values, groups...)
		total += int64(len(values))
	}

	return base.NewPaginationResponse(groups, req.Page, req.PageSize, total), nil
}

// builtinGroupValues 将内置分组对象列表转为值切片（保持固定顺序：首页 → 友链页）。
func builtinGroupValues() []entity.LinkGroup {
	builtins := entity.NewBuiltinGroups()
	values := make([]entity.LinkGroup, 0, len(builtins))
	for _, group := range builtins {
		values = append(values, *group)
	}
	return values
}

// builtinGroupsVisible 判断内置分组是否满足当前过滤条件。
//
// 内置分组恒为启用状态，仅当状态/名称/启用过滤未排除它们时才注入。
func builtinGroupsVisible(status *int, name *string, onlyEnabled *bool) bool {
	if status != nil && *status != 1 {
		return false
	}
	if onlyEnabled != nil && !*onlyEnabled {
		return false
	}
	if name != nil && *name != "" {
		for _, group := range entity.NewBuiltinGroups() {
			if strings.Contains(group.Name, *name) {
				return true
			}
		}
		return false
	}
	return true
}
