package repository

import (
	"errors"
	"strings"
	"time"

	xError "github.com/bamboo-services/bamboo-base-go/common/error"
	xLog "github.com/bamboo-services/bamboo-base-go/common/log"
	apiLink "github.com/bamboo-services/bamboo-main/api/link"
	"github.com/bamboo-services/bamboo-main/internal/entity"
	"github.com/bamboo-services/bamboo-main/internal/repository/cache"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type LinkGroupRepo struct {
	db    *gorm.DB
	cache *cache.LinkGroupCache
	log   *xLog.LogNamedLogger
}

func NewLinkGroupRepo(db *gorm.DB, rdb *redis.Client) *LinkGroupRepo {
	return &LinkGroupRepo{
		db: db,
		cache: &cache.LinkGroupCache{
			RDB: rdb,
			TTL: time.Minute * 15,
		},
		log: xLog.WithName(xLog.NamedREPO, "LinkGroupRepo"),
	}
}

func (r *LinkGroupRepo) pickDB(tx *gorm.DB) *gorm.DB {
	if tx != nil {
		return tx
	}
	return r.db
}

func (r *LinkGroupRepo) GetMaxSortOrder(ctx *gin.Context, tx *gorm.DB) (int, *xError.Error) {
	var maxSort int
	err := r.pickDB(tx).WithContext(ctx).Model(&entity.LinkGroup{}).Select("COALESCE(MAX(sort_order), 0)").Scan(&maxSort).Error
	if err != nil {
		return 0, xError.NewError(ctx, xError.DatabaseError, "查询分组最大排序失败", false, err)
	}

	return maxSort, nil
}

func (r *LinkGroupRepo) Create(ctx *gin.Context, group *entity.LinkGroup, tx *gorm.DB) (*entity.LinkGroup, *xError.Error) {
	err := r.pickDB(tx).WithContext(ctx).Create(group).Error
	if err != nil {
		return nil, xError.NewError(ctx, xError.DatabaseError, "创建友链分组失败", false, err)
	}

	if cacheErr := r.cache.Set(ctx.Request.Context(), group); cacheErr != nil {
		r.log.Warn(ctx, cacheErr.Error())
	}

	return group, nil
}

func (r *LinkGroupRepo) Save(ctx *gin.Context, group *entity.LinkGroup, tx *gorm.DB) (*entity.LinkGroup, *xError.Error) {
	err := r.pickDB(tx).WithContext(ctx).Save(group).Error
	if err != nil {
		return nil, xError.NewError(ctx, xError.DatabaseError, "保存友链分组失败", false, err)
	}

	if cacheErr := r.cache.Delete(ctx.Request.Context(), group.ID); cacheErr != nil {
		r.log.Warn(ctx, cacheErr.Error())
	}

	return group, nil
}

func (r *LinkGroupRepo) GetByID(ctx *gin.Context, id int64, withLinks bool, tx *gorm.DB) (*entity.LinkGroup, bool, *xError.Error) {
	if !withLinks {
		if group, err := r.cache.Get(ctx.Request.Context(), id); err != nil {
			return nil, false, xError.NewError(ctx, xError.CacheError, "获取友链分组缓存失败", true, err)
		} else if group != nil {
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
		return nil, false, xError.NewError(ctx, xError.DatabaseError, "查询友链分组失败", false, err)
	}

	if cacheErr := r.cache.Set(ctx.Request.Context(), &group); cacheErr != nil {
		r.log.Warn(ctx, cacheErr.Error())
	}

	return &group, true, nil
}

func (r *LinkGroupRepo) UpdateStatusByID(ctx *gin.Context, id int64, status bool, tx *gorm.DB) (bool, *xError.Error) {
	result := r.pickDB(tx).WithContext(ctx).Model(&entity.LinkGroup{}).Where("id = ?", id).Update("status", status)
	if result.Error != nil {
		return false, xError.NewError(ctx, xError.DatabaseError, "更新分组状态失败", false, result.Error)
	}
	if result.RowsAffected == 0 {
		return false, nil
	}

	if cacheErr := r.cache.Delete(ctx.Request.Context(), id); cacheErr != nil {
		r.log.Warn(ctx, cacheErr.Error())
	}

	return true, nil
}

func (r *LinkGroupRepo) UpdateSortByIDs(ctx *gin.Context, ids []int64, startSort int, tx *gorm.DB) *xError.Error {
	db := r.pickDB(tx).WithContext(ctx)
	for i, id := range ids {
		result := db.Model(&entity.LinkGroup{}).Where("id = ?", id).Update("sort_order", startSort+i)
		if result.Error != nil {
			return xError.NewError(ctx, xError.DatabaseError, "更新分组排序失败", false, result.Error)
		}
		if result.RowsAffected == 0 {
			return xError.NewError(ctx, xError.NotFound, "分组不存在", false)
		}
		if cacheErr := r.cache.Delete(ctx.Request.Context(), id); cacheErr != nil {
			r.log.Warn(ctx, cacheErr.Error())
		}
	}

	return nil
}

func (r *LinkGroupRepo) DeleteByID(ctx *gin.Context, id int64, tx *gorm.DB) (bool, *xError.Error) {
	result := r.pickDB(tx).WithContext(ctx).Unscoped().Where("id = ?", id).Delete(&entity.LinkGroup{})
	if result.Error != nil {
		return false, xError.NewError(ctx, xError.DatabaseError, "删除友链分组失败", false, result.Error)
	}
	if result.RowsAffected == 0 {
		return false, nil
	}

	if cacheErr := r.cache.Delete(ctx.Request.Context(), id); cacheErr != nil {
		r.log.Warn(ctx, cacheErr.Error())
	}

	return true, nil
}

func (r *LinkGroupRepo) List(ctx *gin.Context, req *apiLink.GroupListRequest, tx *gorm.DB) ([]entity.LinkGroup, *xError.Error) {
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
		return nil, xError.NewError(ctx, xError.DatabaseError, "查询友链分组列表失败", false, err)
	}

	return groups, nil
}

func (r *LinkGroupRepo) Page(ctx *gin.Context, req *apiLink.GroupPageRequest, tx *gorm.DB) ([]entity.LinkGroup, int64, *xError.Error) {
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
		return nil, 0, xError.NewError(ctx, xError.DatabaseError, "统计友链分组数量失败", false, err)
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
		return nil, 0, xError.NewError(ctx, xError.DatabaseError, "查询友链分组列表失败", false, err)
	}

	return groups, total, nil
}
