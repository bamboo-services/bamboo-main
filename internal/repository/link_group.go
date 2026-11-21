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

package repository

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/bamboo-services/bamboo-base-go/foundation/v2/apperr"
	"github.com/bamboo-services/bamboo-base-go/foundation/v2/logging"
	"github.com/bamboo-services/bamboo-base-go/foundation/v2/snowflake"
	"github.com/bamboo-services/bamboo-base-go/plugins/cache/v2"
	"github.com/bamboo-services/bamboo-base-go/plugins/database/v2"
	"github.com/bamboo-services/bamboo-base-go/plugins/web-gin/v2/httpresult"
	"github.com/bamboo-services/bamboo-main/internal/entity"
	"github.com/bamboo-services/bamboo-main/pkg/constants"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// LinkGroupRepo 友链分组数据访问层
//
// 收口友链分组实体的 CRUD、列表与分页查询，缓存经 cache.Manager 统一失效。
type LinkGroupRepo struct {
	db  *gorm.DB
	kc  cache.KeyCache[string, entity.LinkGroup]
	log *logging.LogNamedLogger
}

// GroupListQuery 友链分组列表查询条件
//
// repository 层自有查询模型，由 logic 层从 transport DTO（api/link）转换而来。
type GroupListQuery struct {
	Status      *int
	Name        *string
	WithLinks   *bool
	OnlyEnabled *bool
	OrderBy     *string
	Order       *string
}

// GroupPageQuery 友链分组分页查询条件
type GroupPageQuery struct {
	Page     int
	PageSize int
	Status   *int
	Name     *string
	OrderBy  *string
	Order    *string
}

// NewLinkGroupRepo 创建 LinkGroupRepo 实例
func NewLinkGroupRepo(db *gorm.DB, m *cache.Manager, logger ...*logging.Logger) *LinkGroupRepo {
	return &LinkGroupRepo{
		db:  db,
		kc:  cache.KeyCacheOf[string, entity.LinkGroup](m),
		log: resolveRepoLogger("LinkGroupRepo", logger...),
	}
}

func (r *LinkGroupRepo) pickDB(tx *gorm.DB) *gorm.DB {
	if tx != nil {
		return tx
	}
	return r.db
}

// GetMaxSortOrder 获取友链分组当前最大排序值
func (r *LinkGroupRepo) GetMaxSortOrder(ctx context.Context, tx *gorm.DB) (int, *apperr.Error) {
	var maxSort int
	err := r.pickDB(tx).WithContext(ctx).Model(&entity.LinkGroup{}).Select("COALESCE(MAX(sort_order), 0)").Scan(&maxSort).Error
	if err != nil {
		return 0, apperr.New(database.DatabaseError, "查询分组最大排序失败", err)
	}

	return maxSort, nil
}

// Create 创建友链分组
func (r *LinkGroupRepo) Create(ctx context.Context, group *entity.LinkGroup, tx *gorm.DB) (*entity.LinkGroup, *apperr.Error) {
	err := r.pickDB(tx).WithContext(ctx).Create(group).Error
	if err != nil {
		return nil, apperr.New(database.DatabaseError, "创建友链分组失败", err)
	}

	if cacheErr := r.kc.Set(ctx, constants.RedisLinkGroup.Get(group.ID).String(), group, cache.WithTTL(15*time.Minute)); cacheErr != nil {
		r.log.Warn(ctx, cacheErr.Error())
	}

	return group, nil
}

// Save 保存友链分组（存在则更新）
//
// 经 Omit(clause.Associations) 收敛关联写入：实体可能来自缓存快照（含 LinksFKey），
// 保留关联会让 GORM Save 用友链旧快照覆盖友链记录。关联引用不参与保存，交由外键字段落库。
func (r *LinkGroupRepo) Save(ctx context.Context, group *entity.LinkGroup, tx *gorm.DB) (*entity.LinkGroup, *apperr.Error) {
	err := r.pickDB(tx).WithContext(ctx).Omit(clause.Associations).Save(group).Error
	if err != nil {
		return nil, apperr.New(database.DatabaseError, "保存友链分组失败", err)
	}

	if cacheErr := r.kc.Delete(ctx, constants.RedisLinkGroup.Get(group.ID).String()); cacheErr != nil {
		r.log.Warn(ctx, cacheErr.Error())
	}

	return group, nil
}

// GetByID 根据ID获取友链分组
func (r *LinkGroupRepo) GetByID(ctx context.Context, id snowflake.SnowflakeID, withLinks bool, tx *gorm.DB) (*entity.LinkGroup, bool, *apperr.Error) {
	if !withLinks {
		if group, ok, err := r.kc.Get(ctx, constants.RedisLinkGroup.Get(id).String()); err != nil {
			return nil, false, apperr.New(cache.CacheError, "获取友链分组缓存失败", err)
		} else if ok {
			return group, true, nil
		}
	}

	query := r.pickDB(tx).WithContext(ctx)
	if withLinks {
		query = query.Preload("LinksFKey")
	}

	var group entity.LinkGroup
	err := query.Where("id = ?", id).First(&group).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, false, nil
		}
		return nil, false, apperr.New(database.DatabaseError, "查询友链分组失败", err)
	}

	if cacheErr := r.kc.Set(ctx, constants.RedisLinkGroup.Get(group.ID).String(), &group, cache.WithTTL(15*time.Minute)); cacheErr != nil {
		r.log.Warn(ctx, cacheErr.Error())
	}

	return &group, true, nil
}

// UpdateStatusByID 更新友链分组的启用状态
func (r *LinkGroupRepo) UpdateStatusByID(ctx context.Context, id snowflake.SnowflakeID, status bool, tx *gorm.DB) (bool, *apperr.Error) {
	result := r.pickDB(tx).WithContext(ctx).Model(&entity.LinkGroup{}).Where("id = ?", id).Update("status", status)
	if result.Error != nil {
		return false, apperr.New(database.DatabaseError, "更新分组状态失败", result.Error)
	}
	if result.RowsAffected == 0 {
		return false, nil
	}

	if cacheErr := r.kc.Delete(ctx, constants.RedisLinkGroup.Get(id).String()); cacheErr != nil {
		r.log.Warn(ctx, cacheErr.Error())
	}

	return true, nil
}

// groupSortRowTemplate 友链分组重排 VALUES 行模板（id, sort_order）
//
// 占位符逐项显式 cast：id 为 bigint，sort_order 为 int；避免 PostgreSQL 将纯参数列
// 默认推断为 text 触发 SQLSTATE 42883。占位符数量由 TestSortRowTemplate 兜底。
const groupSortRowTemplate = "(?::bigint, ?::int)"

// UpdateSortByIDs 批量更新友链分组的排序值
//
// 单条 UPDATE ... FROM (VALUES ...) 一次落库全部条目，消除逐行 N+1 往返；
// 排序视为内容更新，刷新 updated_at；RowsAffected==len(ids) 依赖 PostgreSQL
// 「匹配行」语义作原子存在性兜底；VALUES 语法与 ::cast 为 PostgreSQL 专属。
// 写后逐条失效单条缓存。
func (r *LinkGroupRepo) UpdateSortByIDs(ctx context.Context, ids []snowflake.SnowflakeID, startSort int, tx *gorm.DB) *apperr.Error {
	if len(ids) == 0 {
		return nil
	}

	db := r.pickDB(tx).WithContext(ctx)
	values := make([]string, 0, len(ids))
	args := make([]any, 0, len(ids)*2)
	for i, id := range ids {
		values = append(values, groupSortRowTemplate)
		args = append(args, id, startSort+i)
	}

	result := db.Exec(
		fmt.Sprintf(`UPDATE bm_link_group
			SET sort_order = v.sort_order, updated_at = now()
			FROM (VALUES %s) AS v(id, sort_order)
			WHERE bm_link_group.id = v.id`, strings.Join(values, ",")),
		args...,
	)
	if result.Error != nil {
		return apperr.New(database.DatabaseError, "更新分组排序失败", result.Error)
	}
	if result.RowsAffected != int64(len(ids)) {
		return apperr.New(httpresult.NotFound, "分组不存在")
	}

	for _, id := range ids {
		if cacheErr := r.kc.Delete(ctx, constants.RedisLinkGroup.Get(id).String()); cacheErr != nil {
			r.log.Warn(ctx, cacheErr.Error())
		}
	}

	return nil
}

// GetByIDs 按 ID 列表批量查询友链分组（无预加载、不走缓存，供排序前存在性与启用状态校验）
func (r *LinkGroupRepo) GetByIDs(ctx context.Context, ids []snowflake.SnowflakeID, tx *gorm.DB) ([]entity.LinkGroup, *apperr.Error) {
	r.log.Info(ctx, "GetByIDs - 批量查询友链分组")

	var groups []entity.LinkGroup
	err := r.pickDB(tx).WithContext(ctx).Where("id IN ?", ids).Find(&groups).Error
	if err != nil {
		return nil, apperr.New(database.DatabaseError, "批量查询友链分组失败", err)
	}

	return groups, nil
}

// DeleteByID 物理删除友链分组
func (r *LinkGroupRepo) DeleteByID(ctx context.Context, id snowflake.SnowflakeID, tx *gorm.DB) (bool, *apperr.Error) {
	result := r.pickDB(tx).WithContext(ctx).Unscoped().Where("id = ?", id).Delete(&entity.LinkGroup{})
	if result.Error != nil {
		return false, apperr.New(database.DatabaseError, "删除友链分组失败", result.Error)
	}
	if result.RowsAffected == 0 {
		return false, nil
	}

	if cacheErr := r.kc.Delete(ctx, constants.RedisLinkGroup.Get(id).String()); cacheErr != nil {
		r.log.Warn(ctx, cacheErr.Error())
	}

	return true, nil
}

// List 查询友链分组列表
func (r *LinkGroupRepo) List(ctx context.Context, req *GroupListQuery, tx *gorm.DB) ([]entity.LinkGroup, *apperr.Error) {
	query := r.pickDB(tx).WithContext(ctx).Model(&entity.LinkGroup{})

	if req.Status != nil {
		query = query.Where("status = ?", *req.Status == 1)
	}
	if req.OnlyEnabled != nil && *req.OnlyEnabled {
		query = query.Where("status = ?", true)
	}
	if req.Name != nil && *req.Name != "" {
		query = query.Where("name ILIKE ?", "%"+*req.Name+"%")
	}
	if req.WithLinks != nil && *req.WithLinks {
		query = query.Preload("LinksFKey")
	}

	orderBy := "sort_order"
	if req.OrderBy != nil && *req.OrderBy != "" {
		switch *req.OrderBy {
		case "name", "sort_order", "created_at":
			orderBy = *req.OrderBy
		}
	}

	order := "ASC"
	if req.Order != nil && strings.EqualFold(*req.Order, "desc") {
		order = "DESC"
	}

	var groups []entity.LinkGroup
	err := query.Order(orderBy + " " + order).Find(&groups).Error
	if err != nil {
		return nil, apperr.New(database.DatabaseError, "查询友链分组列表失败", err)
	}

	return groups, nil
}

// Page 分页查询友链分组
func (r *LinkGroupRepo) Page(ctx context.Context, req *GroupPageQuery, tx *gorm.DB) ([]entity.LinkGroup, int64, *apperr.Error) {
	query := r.pickDB(tx).WithContext(ctx).Model(&entity.LinkGroup{})

	if req.Status != nil {
		query = query.Where("status = ?", *req.Status == 1)
	}
	if req.Name != nil && *req.Name != "" {
		query = query.Where("name ILIKE ?", "%"+*req.Name+"%")
	}

	var total int64
	err := query.Count(&total).Error
	if err != nil {
		return nil, 0, apperr.New(database.DatabaseError, "统计友链分组数量失败", err)
	}

	orderBy := "sort_order"
	if req.OrderBy != nil && *req.OrderBy != "" {
		switch *req.OrderBy {
		case "name", "sort_order", "created_at":
			orderBy = *req.OrderBy
		}
	}

	order := "ASC"
	if req.Order != nil && strings.EqualFold(*req.Order, "desc") {
		order = "DESC"
	}

	offset := (req.Page - 1) * req.PageSize
	var groups []entity.LinkGroup
	err = query.Order(orderBy + " " + order).Offset(offset).Limit(req.PageSize).Find(&groups).Error
	if err != nil {
		return nil, 0, apperr.New(database.DatabaseError, "查询友链分组列表失败", err)
	}

	return groups, total, nil
}
