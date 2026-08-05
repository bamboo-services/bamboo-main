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
	apiLinkColor "github.com/bamboo-services/bamboo-main/api/link"
	"github.com/bamboo-services/bamboo-main/internal/entity"
	"github.com/bamboo-services/bamboo-main/internal/models/base"
	"github.com/bamboo-services/bamboo-main/internal/repository"
	bConst "github.com/bamboo-services/bamboo-main/pkg/constants"
	"gorm.io/gorm"
)

type linkColorRepo struct {
	color *repository.LinkColorRepo
	link  *repository.LinkRepo
}

// LinkColorLogic 友链颜色业务逻辑
type LinkColorLogic struct {
	logic
	repo linkColorRepo
}

// NewLinkColorLogic 创建 LinkColorLogic 实例，从上下文获取数据库与缓存并初始化颜色与友链仓储依赖。
func NewLinkColorLogic(ctx context.Context) *LinkColorLogic {
	db := xCtxUtil.MustGetDB(ctx)
	m := xCtxUtil.MustGetCacheManager(ctx)

	return &LinkColorLogic{
		logic: logic{
			db:    db,
			cache: m,
			log:   xLog.WithName(xLog.NamedLOGC, "LinkColorLogic"),
		},
		repo: linkColorRepo{
			color: repository.NewLinkColorRepo(db, m),
			link:  repository.NewLinkRepo(db, m),
		},
	}
}

// Add 添加友链颜色，按类型校验配色并自增生成新颜色。
//
// 普通配色（Type=0）仅需主色必填，副色/悬停色可空（渲染取单色）；
// 高级配色（Type=1）需主色、副色、悬停色三者齐备（渲染取三色渐变）。
func (l *LinkColorLogic) Add(ctx context.Context, req *apiLinkColor.ColorAddRequest) (*entity.LinkColor, *xError.Error) {
	if req.PrimaryColor == nil || *req.PrimaryColor == "" {
		return nil, xError.NewError(ctx, xError.BadRequest, "需要设置主颜色", false)
	}
	if req.Type == bConst.ColorTypePremium && (req.SubColor == nil || *req.SubColor == "" || req.HoverColor == nil || *req.HoverColor == "") {
		return nil, xError.NewError(ctx, xError.BadRequest, "高级配色需设置主颜色、副颜色和悬停颜色", false)
	}

	maxSort, xErr := l.repo.color.GetMaxSortOrder(ctx, nil)
	if xErr != nil {
		return nil, xErr
	}

	color := &entity.LinkColor{
		Type:         req.Type,
		Name:         req.ColorName,
		PrimaryColor: req.PrimaryColor,
		SubColor:     req.SubColor,
		HoverColor:   req.HoverColor,
		SortOrder:    maxSort + 1,
		Status:       true,
	}

	_, xErr = l.repo.color.Create(ctx, color, nil)
	if xErr != nil {
		return nil, xErr
	}

	reloaded, found, xErr := l.repo.color.GetByID(ctx, color.ID, true, nil)
	if xErr != nil {
		return nil, xErr
	}
	if !found {
		return nil, xError.NewError(ctx, xError.NotFound, "友链颜色不存在", false)
	}

	return reloaded, nil
}

// Update 更新友链颜色，按请求字段覆盖颜色属性并按结果类型校验配色。
func (l *LinkColorLogic) Update(ctx context.Context, colorID xSnowflake.SnowflakeID, req *apiLinkColor.ColorUpdateRequest) (*entity.LinkColor, *xError.Error) {
	color, found, xErr := l.repo.color.GetByID(ctx, colorID, false, nil)
	if xErr != nil {
		return nil, xErr
	}
	if !found {
		return nil, xError.NewError(ctx, xError.NotFound, "友链颜色不存在", false)
	}

	if req.Type != nil {
		color.Type = *req.Type
	}
	if req.ColorName != nil {
		color.Name = *req.ColorName
	}
	if req.ColorOrder != nil {
		color.SortOrder = *req.ColorOrder
	}

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

	// 按结果类型校验：普通配色仅主色必填；高级配色需三色齐备
	if color.PrimaryColor == nil || *color.PrimaryColor == "" {
		return nil, xError.NewError(ctx, xError.BadRequest, "需要设置主颜色", false)
	}
	if color.Type == bConst.ColorTypePremium && (color.SubColor == nil || *color.SubColor == "" || color.HoverColor == nil || *color.HoverColor == "") {
		return nil, xError.NewError(ctx, xError.BadRequest, "高级配色需设置主颜色、副颜色和悬停颜色", false)
	}

	_, xErr = l.repo.color.Save(ctx, color, nil)
	if xErr != nil {
		return nil, xErr
	}

	reloaded, found, xErr := l.repo.color.GetByID(ctx, colorID, true, nil)
	if xErr != nil {
		return nil, xErr
	}
	if !found {
		return nil, xError.NewError(ctx, xError.NotFound, "友链颜色不存在", false)
	}

	return reloaded, nil
}

// UpdateSort 更新友链颜色排序，在事务内按 ID 列表批量重置排序值。
func (l *LinkColorLogic) UpdateSort(ctx context.Context, req *apiLinkColor.ColorSortRequest) *xError.Error {
	startSort := 0
	if req.SortOrder != nil && *req.SortOrder > 0 {
		startSort = *req.SortOrder
	}

	return l.withTx(ctx, func(tx *gorm.DB) *xError.Error {
		return l.repo.color.UpdateSortByIDs(ctx, req.ColorIDs, startSort, tx)
	})
}

// UpdateStatus 更新友链颜色启用状态。
func (l *LinkColorLogic) UpdateStatus(ctx context.Context, colorID xSnowflake.SnowflakeID, req *apiLinkColor.ColorStatusRequest) *xError.Error {
	ok, xErr := l.repo.color.UpdateStatusByID(ctx, colorID, req.Status, nil)
	if xErr != nil {
		return xErr
	}
	if !ok {
		return xError.NewError(ctx, xError.NotFound, "友链颜色不存在", false)
	}

	return nil
}

// Delete 删除友链颜色，存在关联友链时按 Force 决定阻断或事务清空外键后删除。
func (l *LinkColorLogic) Delete(ctx context.Context, colorID xSnowflake.SnowflakeID, req *apiLinkColor.ColorDeleteRequest) ([]entity.LinkFriend, *xError.Error) {
	_, found, xErr := l.repo.color.GetByID(ctx, colorID, false, nil)
	if xErr != nil {
		return nil, xErr
	}
	if !found {
		return nil, xError.NewError(ctx, xError.NotFound, "友链颜色不存在", false)
	}

	linkCount, xErr := l.repo.link.CountByColorID(ctx, colorID, nil)
	if xErr != nil {
		return nil, xErr
	}

	if linkCount > 0 && !req.Force {
		conflictLinks, xErr := l.repo.link.ListByColorID(ctx, colorID, 10, nil)
		if xErr != nil {
			return nil, xErr
		}
		return conflictLinks, xError.NewError(ctx, xError.BadRequest, "颜色下存在友链，无法删除", false)
	}

	tx := l.db.WithContext(ctx).Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if req.Force && linkCount > 0 {
		if xErr = l.repo.link.ClearColorID(ctx, colorID, tx); xErr != nil {
			tx.Rollback()
			return nil, xErr
		}
	}

	ok, xErr := l.repo.color.DeleteByID(ctx, colorID, tx)
	if xErr != nil {
		tx.Rollback()
		return nil, xErr
	}
	if !ok {
		tx.Rollback()
		return nil, xError.NewError(ctx, xError.NotFound, "友链颜色不存在", false)
	}

	if err := tx.Commit().Error; err != nil {
		return nil, xError.NewError(ctx, xError.DatabaseError, "提交删除操作失败", false, err)
	}

	return nil, nil
}

// Get 获取友链颜色详情。
func (l *LinkColorLogic) Get(ctx context.Context, colorID xSnowflake.SnowflakeID) (*entity.LinkColor, *xError.Error) {
	color, found, xErr := l.repo.color.GetByID(ctx, colorID, true, nil)
	if xErr != nil {
		return nil, xErr
	}
	if !found {
		return nil, xError.NewError(ctx, xError.NotFound, "友链颜色不存在", false)
	}

	return color, nil
}

// GetList 获取友链颜色列表。
//
// 颜色按类型分级（type=0 普通 / type=1 高级），未显式传 type 时返回全部启用配色；
// 显式传 type 参数可按级别过滤。内置炫彩按过滤条件置顶注入虚拟记录（不落库）。
func (l *LinkColorLogic) GetList(ctx context.Context, req *apiLinkColor.ColorListRequest) ([]entity.LinkColor, *xError.Error) {
	colors, xErr := l.repo.color.List(ctx, &repository.ColorListQuery{
		Status:      req.Status,
		Name:        req.Name,
		OnlyEnabled: req.OnlyEnabled,
		Type:        req.Type,
		OrderBy:     req.OrderBy,
		Order:       req.Order,
	}, nil)
	if xErr != nil {
		return nil, xErr
	}

	// 内置炫彩（不落库）：满足过滤条件时置顶注入虚拟记录，供前端选择器使用
	if builtinFancyVisible(req, req.Type) {
		colors = append([]entity.LinkColor{*entity.NewFancyColor()}, colors...)
	}

	return colors, nil
}

// builtinFancyVisible 判断内置炫彩是否满足当前列表过滤条件。
//
// 炫彩恒为启用状态，仅当名称/状态过滤未排除它、且未被按类型过滤（type=1 高级色）剔除时才注入。
func builtinFancyVisible(req *apiLinkColor.ColorListRequest, typeFilter *int) bool {
	if req.Status != nil && *req.Status != 1 {
		return false
	}
	if req.OnlyEnabled != nil && !*req.OnlyEnabled {
		return false
	}
	if req.Name != nil && *req.Name != "" && !strings.Contains("炫彩", *req.Name) {
		return false
	}
	if typeFilter != nil && *typeFilter != bConst.ColorTypeNormal {
		return false
	}
	return true
}

// GetPage 分页获取友链颜色。
func (l *LinkColorLogic) GetPage(ctx context.Context, req *apiLinkColor.ColorPageRequest) (*base.PaginationResponse[entity.LinkColor], *xError.Error) {
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PageSize <= 0 || req.PageSize > 100 {
		req.PageSize = 10
	}

	colors, total, xErr := l.repo.color.Page(ctx, &repository.ColorPageQuery{
		Page:     req.Page,
		PageSize: req.PageSize,
		Status:   req.Status,
		Name:     req.Name,
		Type:     req.Type,
		OrderBy:  req.OrderBy,
		Order:    req.Order,
	}, nil)
	if xErr != nil {
		return nil, xErr
	}

	return base.NewPaginationResponse(colors, req.Page, req.PageSize, total), nil
}
