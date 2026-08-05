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
	bConst "github.com/bamboo-services/bamboo-main/pkg/constants"
	"gorm.io/gorm"
)

type linkGroupRepo struct {
	group  *repository.LinkGroupRepo
	link   *repository.LinkRepo
	system *repository.SystemRepo
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
			group:  repository.NewLinkGroupRepo(db, m),
			link:   repository.NewLinkRepo(db, m),
			system: repository.NewSystemRepo(db, m),
		},
	}
}

// Add 添加友链分组，按当前最大排序值自增生成新分组。
func (l *LinkGroupLogic) Add(ctx context.Context, req *apiLinkGroup.GroupAddRequest) (*entity.LinkGroup, *xError.Error) {
	maxSort, xErr := l.repo.group.GetMaxSortOrder(ctx, nil)
	if xErr != nil {
		return nil, xErr
	}

	// 内置「已失效」分组不参与排序持久化，自定义分组排序值自最大排序 +1 顺延
	sortOrder := maxSort + 1

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
	// 内置「已失效」分组为系统语义分组（不落库），拒绝编辑
	if entity.IsBuiltinGroupID(groupID) {
		return nil, xError.NewError(ctx, xError.BadRequest, "内置已失效分组请使用「已失效配置」接口修改", false)
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
	// 内置「已失效」分组恒为启用状态（不落库），拒绝切换
	if entity.IsBuiltinGroupID(groupID) {
		return xError.NewError(ctx, xError.BadRequest, "内置已失效分组恒为启用，不支持切换", false)
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
	// 内置「已失效」分组为系统语义分组（不落库），拒绝删除
	if entity.IsBuiltinGroupID(groupID) {
		return nil, xError.NewError(ctx, xError.BadRequest, "内置已失效分组不可删除", false)
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
//
// 内置「已失效」分组不落库，按保留 ID 请求时返回 bm_system 配置构造的虚拟分组。
func (l *LinkGroupLogic) Get(ctx context.Context, groupID xSnowflake.SnowflakeID) (*entity.LinkGroup, *xError.Error) {
	if entity.IsBuiltinGroupID(groupID) {
		return l.GetBuiltinInvalidGroup(ctx)
	}

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
//
// 内置「已失效」分组为系统语义分组（非可选展示位置），不参与列表注入。
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

	return groups, nil
}

// GetPage 分页获取友链分组。
//
// 内置「已失效」分组不参与分页注入。
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

	return base.NewPaginationResponse(groups, req.Page, req.PageSize, total), nil
}

// GetBuiltinInvalidGroup 获取内置「已失效」分组配置（名称与描述，经 bm_system 热修改）。
func (l *LinkGroupLogic) GetBuiltinInvalidGroup(ctx context.Context) (*entity.LinkGroup, *xError.Error) {
	return l.repo.system.BuildBuiltinInvalidGroup(ctx)
}

// UpdateBuiltinInvalidGroup 更新内置「已失效」分组配置。
//
// 按请求字段逐 key 写入 bm_system（PATCH 语义：仅更新非 nil 字段；描述传空串即清空），
// 更新后回读最新配置返回。
func (l *LinkGroupLogic) UpdateBuiltinInvalidGroup(ctx context.Context, req *apiLinkGroup.BuiltinInvalidGroupUpdateRequest) (*entity.LinkGroup, *xError.Error) {
	updates := make(map[string]*string)

	if req.Name != nil {
		if strings.TrimSpace(*req.Name) == "" {
			return nil, xError.NewError(ctx, xError.BadRequest, "已失效分组名称不能为空", false)
		}
		updates[bConst.KeyBuiltinInvalidGroupName] = req.Name
	}
	if req.Description != nil {
		updates[bConst.KeyBuiltinInvalidGroupDesc] = req.Description
	}
	if len(updates) == 0 {
		return l.GetBuiltinInvalidGroup(ctx)
	}

	for key, value := range updates {
		if xErr := l.repo.system.UpdateValueByKey(ctx, key, value); xErr != nil {
			return nil, xErr
		}
	}

	return l.GetBuiltinInvalidGroup(ctx)
}
