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

	xError "github.com/bamboo-services/bamboo-base-go/common/error"
	xLog "github.com/bamboo-services/bamboo-base-go/common/log"
	xSnowflake "github.com/bamboo-services/bamboo-base-go/common/snowflake"
	xCache "github.com/bamboo-services/bamboo-base-go/major/cache"
	"github.com/bamboo-services/bamboo-main/internal/entity"
	"github.com/bamboo-services/bamboo-main/pkg/constants"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// LinkColorRepo 友链颜色数据访问层
//
// 收口友链颜色实体的 CRUD、列表与分页查询，缓存经 xCache.Manager 统一失效。
type LinkColorRepo struct {
	db  *gorm.DB
	kc  xCache.KeyCache[string, entity.LinkColor]
	log *xLog.LogNamedLogger
}

// ColorListQuery 友链颜色列表查询条件
//
// repository 层自有查询模型，由 logic 层从 transport DTO（api/link）转换而来。
type ColorListQuery struct {
	Status      *int
	Name        *string
	OnlyEnabled *bool
	OrderBy     *string
	Order       *string
}

// ColorPageQuery 友链颜色分页查询条件
type ColorPageQuery struct {
	Page     int
	PageSize int
	Status   *int
	Name     *string
	OrderBy  *string
	Order    *string
}

// NewLinkColorRepo 创建 LinkColorRepo 实例
func NewLinkColorRepo(db *gorm.DB, m *xCache.Manager) *LinkColorRepo {
	return &LinkColorRepo{
		db:  db,
		kc:  xCache.KeyCacheOf[string, entity.LinkColor](m),
		log: xLog.WithName(xLog.NamedREPO, "LinkColorRepo"),
	}
}

func (r *LinkColorRepo) pickDB(tx *gorm.DB) *gorm.DB {
	if tx != nil {
		return tx
	}
	return r.db
}

// GetMaxSortOrder 获取友链颜色当前最大排序值
func (r *LinkColorRepo) GetMaxSortOrder(ctx context.Context, tx *gorm.DB) (int, *xError.Error) {
	var maxSort int
	err := r.pickDB(tx).WithContext(ctx).Model(&entity.LinkColor{}).Select("COALESCE(MAX(sort_order), 0)").Scan(&maxSort).Error
	if err != nil {
		return 0, xError.NewError(ctx, xError.DatabaseError, "查询颜色最大排序失败", false, err)
	}

	return maxSort, nil
}

// Create 创建友链颜色
func (r *LinkColorRepo) Create(ctx context.Context, color *entity.LinkColor, tx *gorm.DB) (*entity.LinkColor, *xError.Error) {
	err := r.pickDB(tx).WithContext(ctx).Create(color).Error
	if err != nil {
		return nil, xError.NewError(ctx, xError.DatabaseError, "创建友链颜色失败", false, err)
	}

	if cacheErr := r.kc.Set(ctx, constants.RedisLinkColor.Get(color.ID).String(), color, xCache.WithTTL(15*time.Minute)); cacheErr != nil {
		r.log.Warn(ctx, cacheErr.Error())
	}

	return color, nil
}

// Save 保存友链颜色（存在则更新）
//
// 经 Omit(clause.Associations) 收敛关联写入：实体可能来自缓存快照（含 LinksFKey），
// 保留关联会让 GORM Save 用友链旧快照覆盖友链记录。关联引用不参与保存，交由外键字段落库。
func (r *LinkColorRepo) Save(ctx context.Context, color *entity.LinkColor, tx *gorm.DB) (*entity.LinkColor, *xError.Error) {
	err := r.pickDB(tx).WithContext(ctx).Omit(clause.Associations).Save(color).Error
	if err != nil {
		return nil, xError.NewError(ctx, xError.DatabaseError, "保存友链颜色失败", false, err)
	}

	if cacheErr := r.kc.Delete(ctx, constants.RedisLinkColor.Get(color.ID).String()); cacheErr != nil {
		r.log.Warn(ctx, cacheErr.Error())
	}

	return color, nil
}

// GetByID 根据ID获取友链颜色
func (r *LinkColorRepo) GetByID(ctx context.Context, id xSnowflake.SnowflakeID, withLinks bool, tx *gorm.DB) (*entity.LinkColor, bool, *xError.Error) {
	if !withLinks {
		if color, ok, err := r.kc.Get(ctx, constants.RedisLinkColor.Get(id).String()); err != nil {
			return nil, false, xError.NewError(ctx, xError.CacheError, "获取友链颜色缓存失败", true, err)
		} else if ok {
			return color, true, nil
		}
	}

	query := r.pickDB(tx).WithContext(ctx)
	if withLinks {
		query = query.Preload("LinksFKey")
	}

	var color entity.LinkColor
	err := query.Where("id = ?", id).First(&color).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, false, nil
		}
		return nil, false, xError.NewError(ctx, xError.DatabaseError, "查询友链颜色失败", false, err)
	}

	if cacheErr := r.kc.Set(ctx, constants.RedisLinkColor.Get(color.ID).String(), &color, xCache.WithTTL(15*time.Minute)); cacheErr != nil {
		r.log.Warn(ctx, cacheErr.Error())
	}

	return &color, true, nil
}

// UpdateStatusByID 更新友链颜色的启用状态
func (r *LinkColorRepo) UpdateStatusByID(ctx context.Context, id xSnowflake.SnowflakeID, status bool, tx *gorm.DB) (bool, *xError.Error) {
	result := r.pickDB(tx).WithContext(ctx).Model(&entity.LinkColor{}).Where("id = ?", id).Update("status", status)
	if result.Error != nil {
		return false, xError.NewError(ctx, xError.DatabaseError, "更新颜色状态失败", false, result.Error)
	}
	if result.RowsAffected == 0 {
		return false, nil
	}

	if cacheErr := r.kc.Delete(ctx, constants.RedisLinkColor.Get(id).String()); cacheErr != nil {
		r.log.Warn(ctx, cacheErr.Error())
	}

	return true, nil
}

// colorSortRowTemplate 友链颜色重排 VALUES 行模板（id, sort_order）
//
// 占位符逐项显式 cast：id 为 bigint，sort_order 为 int；避免 PostgreSQL 将纯参数列
// 默认推断为 text 触发 SQLSTATE 42883。占位符数量由 TestSortRowTemplate 兜底。
const colorSortRowTemplate = "(?::bigint, ?::int)"

// UpdateSortByIDs 批量更新友链颜色的排序值
//
// 单条 UPDATE ... FROM (VALUES ...) 一次落库全部条目，消除逐行 N+1 往返；
// 排序视为内容更新，刷新 updated_at；RowsAffected==len(ids) 依赖 PostgreSQL
// 「匹配行」语义作原子存在性兜底；VALUES 语法与 ::cast 为 PostgreSQL 专属。
// 写后逐条失效单条缓存。
func (r *LinkColorRepo) UpdateSortByIDs(ctx context.Context, ids []xSnowflake.SnowflakeID, startSort int, tx *gorm.DB) *xError.Error {
	if len(ids) == 0 {
		return nil
	}

	db := r.pickDB(tx).WithContext(ctx)
	values := make([]string, 0, len(ids))
	args := make([]any, 0, len(ids)*2)
	for i, id := range ids {
		values = append(values, colorSortRowTemplate)
		args = append(args, id, startSort+i)
	}

	result := db.Exec(
		fmt.Sprintf(`UPDATE bm_link_color
			SET sort_order = v.sort_order, updated_at = now()
			FROM (VALUES %s) AS v(id, sort_order)
			WHERE bm_link_color.id = v.id`, strings.Join(values, ",")),
		args...,
	)
	if result.Error != nil {
		return xError.NewError(ctx, xError.DatabaseError, "更新颜色排序失败", false, result.Error)
	}
	if result.RowsAffected != int64(len(ids)) {
		return xError.NewError(ctx, xError.NotFound, "颜色不存在", false)
	}

	for _, id := range ids {
		if cacheErr := r.kc.Delete(ctx, constants.RedisLinkColor.Get(id).String()); cacheErr != nil {
			r.log.Warn(ctx, cacheErr.Error())
		}
	}

	return nil
}

// DeleteByID 物理删除友链颜色
func (r *LinkColorRepo) DeleteByID(ctx context.Context, id xSnowflake.SnowflakeID, tx *gorm.DB) (bool, *xError.Error) {
	result := r.pickDB(tx).WithContext(ctx).Unscoped().Where("id = ?", id).Delete(&entity.LinkColor{})
	if result.Error != nil {
		return false, xError.NewError(ctx, xError.DatabaseError, "删除友链颜色失败", false, result.Error)
	}
	if result.RowsAffected == 0 {
		return false, nil
	}

	if cacheErr := r.kc.Delete(ctx, constants.RedisLinkColor.Get(id).String()); cacheErr != nil {
		r.log.Warn(ctx, cacheErr.Error())
	}

	return true, nil
}

// List 查询友链颜色列表
func (r *LinkColorRepo) List(ctx context.Context, req *ColorListQuery, tx *gorm.DB) ([]entity.LinkColor, *xError.Error) {
	query := r.pickDB(tx).WithContext(ctx).Model(&entity.LinkColor{})

	if req.Status != nil {
		query = query.Where("status = ?", *req.Status == 1)
	}
	if req.OnlyEnabled != nil && *req.OnlyEnabled {
		query = query.Where("status = ?", true)
	}
	if req.Name != nil && *req.Name != "" {
		query = query.Where("name ILIKE ?", "%"+*req.Name+"%")
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

	var colors []entity.LinkColor
	err := query.Order(orderBy + " " + order).Find(&colors).Error
	if err != nil {
		return nil, xError.NewError(ctx, xError.DatabaseError, "查询友链颜色列表失败", false, err)
	}

	return colors, nil
}

// Page 分页查询友链颜色
func (r *LinkColorRepo) Page(ctx context.Context, req *ColorPageQuery, tx *gorm.DB) ([]entity.LinkColor, int64, *xError.Error) {
	query := r.pickDB(tx).WithContext(ctx).Model(&entity.LinkColor{})

	if req.Status != nil {
		query = query.Where("status = ?", *req.Status == 1)
	}
	if req.Name != nil && *req.Name != "" {
		query = query.Where("name ILIKE ?", "%"+*req.Name+"%")
	}

	var total int64
	err := query.Count(&total).Error
	if err != nil {
		return nil, 0, xError.NewError(ctx, xError.DatabaseError, "统计友链颜色数量失败", false, err)
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
	var colors []entity.LinkColor
	err = query.Order(orderBy + " " + order).Offset(offset).Limit(req.PageSize).Find(&colors).Error
	if err != nil {
		return nil, 0, xError.NewError(ctx, xError.DatabaseError, "查询友链颜色列表失败", false, err)
	}

	return colors, total, nil
}
