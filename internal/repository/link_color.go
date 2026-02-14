package repository

import (
	"errors"
	"strings"
	"time"

	xError "github.com/bamboo-services/bamboo-base-go/error"
	xLog "github.com/bamboo-services/bamboo-base-go/log"
	apiLink "github.com/bamboo-services/bamboo-main/api/link"
	"github.com/bamboo-services/bamboo-main/internal/entity"
	"github.com/bamboo-services/bamboo-main/internal/repository/cache"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type LinkColorRepo struct {
	db    *gorm.DB
	cache *cache.LinkColorCache
	log   *xLog.LogNamedLogger
}

func NewLinkColorRepo(db *gorm.DB, rdb *redis.Client) *LinkColorRepo {
	return &LinkColorRepo{
		db: db,
		cache: &cache.LinkColorCache{
			RDB: rdb,
			TTL: time.Minute * 15,
		},
		log: xLog.WithName(xLog.NamedREPO, "LinkColorRepo"),
	}
}

func (r *LinkColorRepo) pickDB(tx *gorm.DB) *gorm.DB {
	if tx != nil {
		return tx
	}
	return r.db
}

func (r *LinkColorRepo) GetMaxSortOrder(ctx *gin.Context, tx *gorm.DB) (int, *xError.Error) {
	var maxSort int
	err := r.pickDB(tx).WithContext(ctx).Model(&entity.LinkColor{}).Select("COALESCE(MAX(sort_order), 0)").Scan(&maxSort).Error
	if err != nil {
		return 0, xError.NewError(ctx, xError.DatabaseError, "查询颜色最大排序失败", false, err)
	}

	return maxSort, nil
}

func (r *LinkColorRepo) Create(ctx *gin.Context, color *entity.LinkColor, tx *gorm.DB) (*entity.LinkColor, *xError.Error) {
	err := r.pickDB(tx).WithContext(ctx).Create(color).Error
	if err != nil {
		return nil, xError.NewError(ctx, xError.DatabaseError, "创建友链颜色失败", false, err)
	}

	if cacheErr := r.cache.Set(ctx.Request.Context(), color); cacheErr != nil {
		r.log.Warn(ctx, cacheErr.Error())
	}

	return color, nil
}

func (r *LinkColorRepo) Save(ctx *gin.Context, color *entity.LinkColor, tx *gorm.DB) (*entity.LinkColor, *xError.Error) {
	err := r.pickDB(tx).WithContext(ctx).Save(color).Error
	if err != nil {
		return nil, xError.NewError(ctx, xError.DatabaseError, "保存友链颜色失败", false, err)
	}

	if cacheErr := r.cache.Delete(ctx.Request.Context(), color.ID); cacheErr != nil {
		r.log.Warn(ctx, cacheErr.Error())
	}

	return color, nil
}

func (r *LinkColorRepo) GetByID(ctx *gin.Context, id int64, withLinks bool, tx *gorm.DB) (*entity.LinkColor, bool, *xError.Error) {
	if !withLinks {
		if color, err := r.cache.Get(ctx.Request.Context(), id); err != nil {
			return nil, false, xError.NewError(ctx, xError.CacheError, "获取友链颜色缓存失败", true, err)
		} else if color != nil {
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

	if cacheErr := r.cache.Set(ctx.Request.Context(), &color); cacheErr != nil {
		r.log.Warn(ctx, cacheErr.Error())
	}

	return &color, true, nil
}

func (r *LinkColorRepo) UpdateStatusByID(ctx *gin.Context, id int64, status bool, tx *gorm.DB) (bool, *xError.Error) {
	result := r.pickDB(tx).WithContext(ctx).Model(&entity.LinkColor{}).Where("id = ?", id).Update("status", status)
	if result.Error != nil {
		return false, xError.NewError(ctx, xError.DatabaseError, "更新颜色状态失败", false, result.Error)
	}
	if result.RowsAffected == 0 {
		return false, nil
	}

	if cacheErr := r.cache.Delete(ctx.Request.Context(), id); cacheErr != nil {
		r.log.Warn(ctx, cacheErr.Error())
	}

	return true, nil
}

func (r *LinkColorRepo) UpdateSortByIDs(ctx *gin.Context, ids []int64, startSort int, tx *gorm.DB) *xError.Error {
	db := r.pickDB(tx).WithContext(ctx)
	for i, id := range ids {
		result := db.Model(&entity.LinkColor{}).Where("id = ?", id).Update("sort_order", startSort+i)
		if result.Error != nil {
			return xError.NewError(ctx, xError.DatabaseError, "更新颜色排序失败", false, result.Error)
		}
		if result.RowsAffected == 0 {
			return xError.NewError(ctx, xError.NotFound, "颜色不存在", false)
		}
		if cacheErr := r.cache.Delete(ctx.Request.Context(), id); cacheErr != nil {
			r.log.Warn(ctx, cacheErr.Error())
		}
	}

	return nil
}

func (r *LinkColorRepo) DeleteByID(ctx *gin.Context, id int64, tx *gorm.DB) (bool, *xError.Error) {
	result := r.pickDB(tx).WithContext(ctx).Unscoped().Where("id = ?", id).Delete(&entity.LinkColor{})
	if result.Error != nil {
		return false, xError.NewError(ctx, xError.DatabaseError, "删除友链颜色失败", false, result.Error)
	}
	if result.RowsAffected == 0 {
		return false, nil
	}

	if cacheErr := r.cache.Delete(ctx.Request.Context(), id); cacheErr != nil {
		r.log.Warn(ctx, cacheErr.Error())
	}

	return true, nil
}

func (r *LinkColorRepo) List(ctx *gin.Context, req *apiLink.ColorListRequest, tx *gorm.DB) ([]entity.LinkColor, *xError.Error) {
	query := r.pickDB(tx).WithContext(ctx).Model(&entity.LinkColor{})

	if req.Status != nil {
		query = query.Where("status = ?", *req.Status == 1)
	}
	if req.OnlyEnabled != nil && *req.OnlyEnabled {
		query = query.Where("status = ?", true)
	}
	if req.Type != nil {
		query = query.Where("type = ?", *req.Type)
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

func (r *LinkColorRepo) Page(ctx *gin.Context, req *apiLink.ColorPageRequest, tx *gorm.DB) ([]entity.LinkColor, int64, *xError.Error) {
	query := r.pickDB(tx).WithContext(ctx).Model(&entity.LinkColor{})

	if req.Status != nil {
		query = query.Where("status = ?", *req.Status == 1)
	}
	if req.Type != nil {
		query = query.Where("type = ?", *req.Type)
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
