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

	"github.com/bamboo-services/bamboo-base-go/application/v2"
	"github.com/bamboo-services/bamboo-base-go/application/v2/builtin"
	"github.com/bamboo-services/bamboo-base-go/foundation/v2/apperr"
	"github.com/bamboo-services/bamboo-base-go/foundation/v2/logging"
	"github.com/bamboo-services/bamboo-base-go/foundation/v2/snowflake"
	"github.com/bamboo-services/bamboo-base-go/plugins/cache/v2"
	"github.com/bamboo-services/bamboo-base-go/plugins/database/v2"
	"github.com/bamboo-services/bamboo-base-go/plugins/web-gin/v2/httpresult"
	apiLinkGroup "github.com/bamboo-services/bamboo-main/api/link"
	"github.com/bamboo-services/bamboo-main/internal/entity"
	"github.com/bamboo-services/bamboo-main/internal/models/base"
	"github.com/bamboo-services/bamboo-main/internal/repository"
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
	db := database.MustGet(ctx)
	m := cache.MustGet(ctx)
	logger := application.MustGet[*logging.Logger](ctx, builtin.LoggerKey)

	return &LinkGroupLogic{
		logic: logic{
			db:     db,
			cache:  m,
			logger: logger,
			log:    logger.WithName(logging.NamedLOGC, "LinkGroupLogic"),
		},
		repo: linkGroupRepo{
			group:  repository.NewLinkGroupRepo(db, m, logger),
			link:   repository.NewLinkRepo(db, m, logger),
			system: repository.NewSystemRepo(db, m, logger),
		},
	}
}

// Add 添加友链分组，按当前最大排序值自增生成新分组。
func (l *LinkGroupLogic) Add(ctx context.Context, req *apiLinkGroup.GroupAddRequest) (*entity.LinkGroup, *apperr.Error) {
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
		return nil, apperr.New(httpresult.NotFound, "友链分组不存在")
	}

	return reloaded, nil
}

// Update 更新友链分组，按请求字段覆盖名称、描述、排序与状态。
func (l *LinkGroupLogic) Update(ctx context.Context, groupID snowflake.SnowflakeID, req *apiLinkGroup.GroupUpdateRequest) (*entity.LinkGroup, *apperr.Error) {
	// 内置「已失效」分组为系统语义分组（不落库），拒绝编辑
	if entity.IsBuiltinGroupID(groupID) {
		return nil, apperr.New(httpresult.BadRequest, "内置已失效分组请使用「已失效配置」接口修改")
	}

	group, found, xErr := l.repo.group.GetByID(ctx, groupID, false, nil)
	if xErr != nil {
		return nil, xErr
	}
	if !found {
		return nil, apperr.New(httpresult.NotFound, "友链分组不存在")
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
		return nil, apperr.New(httpresult.NotFound, "友链分组不存在")
	}

	return reloaded, nil
}

// UpdateSort 更新友链分组排序，在事务内按 ID 列表批量重置排序值。
func (l *LinkGroupLogic) UpdateSort(ctx context.Context, req *apiLinkGroup.GroupSortRequest) *apperr.Error {
	startSort := 0
	if req.SortOrder != nil && *req.SortOrder > 0 {
		startSort = *req.SortOrder
	}

	// 内置分组不参与排序持久化：剔除保留 ID，真实分组排序值顺延内置占位（0、1 预留给内置）
	validIDs := make([]snowflake.SnowflakeID, 0, len(req.GroupIDs))
	for _, id := range req.GroupIDs {
		if !entity.IsBuiltinGroupID(id) {
			validIDs = append(validIDs, id)
		}
	}
	if len(validIDs) == 0 {
		return nil
	}
	startSort += len(req.GroupIDs) - len(validIDs)

	return l.withTx(ctx, func(tx *gorm.DB) *apperr.Error {
		return l.repo.group.UpdateSortByIDs(ctx, validIDs, startSort, tx)
	})
}

// UpdateStatus 更新友链分组启用状态。
func (l *LinkGroupLogic) UpdateStatus(ctx context.Context, groupID snowflake.SnowflakeID, req *apiLinkGroup.GroupStatusRequest) *apperr.Error {
	// 内置「已失效」分组恒为启用状态（不落库），拒绝切换
	if entity.IsBuiltinGroupID(groupID) {
		return apperr.New(httpresult.BadRequest, "内置已失效分组恒为启用，不支持切换")
	}

	ok, xErr := l.repo.group.UpdateStatusByID(ctx, groupID, req.Status, nil)
	if xErr != nil {
		return xErr
	}
	if !ok {
		return apperr.New(httpresult.NotFound, "友链分组不存在")
	}

	return nil
}

// Delete 删除友链分组，存在关联友链时按 target_group_id / force 决定迁移、清空或阻断。
//
// 处理分组下关联友链的三种结局：target_group_id 指定时迁移至目标分组、force 时清空为未分组、
// 均未指定且存在友链则阻断。target_group_id 与 force 互斥，目标分组必须存在且为启用状态。
func (l *LinkGroupLogic) Delete(ctx context.Context, groupID snowflake.SnowflakeID, req *apiLinkGroup.GroupDeleteRequest) ([]entity.LinkFriend, *apperr.Error) {
	// 内置「已失效」分组为系统语义分组（不落库），拒绝删除
	if entity.IsBuiltinGroupID(groupID) {
		return nil, apperr.New(httpresult.BadRequest, "内置已失效分组不可删除")
	}

	// force 与 target_group_id 互斥，避免语义歧义
	if req.Force && req.TargetGroupID != nil {
		return nil, apperr.New(httpresult.BadRequest, "force 与 target_group_id 不能同时指定")
	}

	_, found, xErr := l.repo.group.GetByID(ctx, groupID, false, nil)
	if xErr != nil {
		return nil, xErr
	}
	if !found {
		return nil, apperr.New(httpresult.NotFound, "友链分组不存在")
	}

	// 目标分组前置校验：非内置、非自身、存在且启用
	if req.TargetGroupID != nil {
		if entity.IsBuiltinGroupID(*req.TargetGroupID) {
			return nil, apperr.New(httpresult.BadRequest, "内置已失效分组不可作为迁移目标")
		}
		if *req.TargetGroupID == groupID {
			return nil, apperr.New(httpresult.BadRequest, "不可迁移到分组自身")
		}
		target, targetFound, xErr := l.repo.group.GetByID(ctx, *req.TargetGroupID, false, nil)
		if xErr != nil {
			return nil, xErr
		}
		if !targetFound {
			return nil, apperr.New(httpresult.NotFound, "目标友链分组不存在")
		}
		if !target.Status {
			return nil, apperr.New(httpresult.BadRequest, "目标分组已禁用，无法迁移")
		}
	}

	linkCount, xErr := l.repo.link.CountByGroupID(ctx, groupID, nil)
	if xErr != nil {
		return nil, xErr
	}

	// 无目标分组且非强制且存在友链：阻断并返回冲突友链
	if linkCount > 0 && !req.Force && req.TargetGroupID == nil {
		conflictLinks, xErr := l.repo.link.ListByGroupID(ctx, groupID, 10, nil)
		if xErr != nil {
			return nil, xErr
		}
		return conflictLinks, apperr.New(httpresult.BadRequest, "分组下存在友链，无法删除")
	}

	// 事务内：迁移或清空关联后物理删除分组
	if xErr := l.withTx(ctx, func(tx *gorm.DB) *apperr.Error {
		if req.TargetGroupID != nil && linkCount > 0 {
			if xErr := l.repo.link.MoveGroupID(ctx, groupID, *req.TargetGroupID, tx); xErr != nil {
				return xErr
			}
		} else if req.Force && linkCount > 0 {
			if xErr := l.repo.link.ClearGroupID(ctx, groupID, tx); xErr != nil {
				return xErr
			}
		}

		ok, xErr := l.repo.group.DeleteByID(ctx, groupID, tx)
		if xErr != nil {
			return xErr
		}
		if !ok {
			return apperr.New(httpresult.NotFound, "友链分组不存在")
		}
		return nil
	}); xErr != nil {
		return nil, xErr
	}

	return nil, nil
}

// Get 获取友链分组详情。
//
// 内置「已失效」分组不落库，按保留 ID 请求时返回 bm_system 配置构造的虚拟分组。
func (l *LinkGroupLogic) Get(ctx context.Context, groupID snowflake.SnowflakeID) (*entity.LinkGroup, *apperr.Error) {
	if entity.IsBuiltinGroupID(groupID) {
		return l.repo.system.BuildBuiltinInvalidGroup(ctx)
	}

	group, found, xErr := l.repo.group.GetByID(ctx, groupID, true, nil)
	if xErr != nil {
		return nil, xErr
	}
	if !found {
		return nil, apperr.New(httpresult.NotFound, "友链分组不存在")
	}

	return group, nil
}

// GetList 获取友链分组列表。
//
// 内置「已失效」分组为系统语义分组（非可选展示位置），不参与列表注入。
func (l *LinkGroupLogic) GetList(ctx context.Context, req *apiLinkGroup.GroupListRequest) ([]entity.LinkGroup, *apperr.Error) {
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
func (l *LinkGroupLogic) GetPage(ctx context.Context, req *apiLinkGroup.GroupPageRequest) (*base.PaginationResponse[entity.LinkGroup], *apperr.Error) {
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
