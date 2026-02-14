package repository

import (
	"errors"
	"time"

	xError "github.com/bamboo-services/bamboo-base-go/error"
	xLog "github.com/bamboo-services/bamboo-base-go/log"
	"github.com/bamboo-services/bamboo-main/internal/entity"
	"github.com/bamboo-services/bamboo-main/internal/repository/cache"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type SponsorRecordRepo struct {
	db    *gorm.DB
	cache *cache.SponsorRecordCache
	log   *xLog.LogNamedLogger
}

type SponsorRecordPageQuery struct {
	Page        int
	PageSize    int
	ChannelID   *int64
	Nickname    string
	IsAnonymous *bool
	IsHidden    *bool
	OrderBy     string
	Order       string
}

type SponsorRecordPublicPageQuery struct {
	Page      int
	PageSize  int
	ChannelID *int64
	OrderBy   string
	Order     string
}

func NewSponsorRecordRepo(db *gorm.DB, rdb *redis.Client) *SponsorRecordRepo {
	return &SponsorRecordRepo{
		db: db,
		cache: &cache.SponsorRecordCache{
			RDB: rdb,
			TTL: time.Minute * 10,
		},
		log: xLog.WithName(xLog.NamedREPO, "SponsorRecordRepo"),
	}
}

func (r *SponsorRecordRepo) Create(ctx *gin.Context, record *entity.SponsorRecord) (*entity.SponsorRecord, *xError.Error) {
	r.log.Info(ctx, "Create - 创建赞助记录")

	err := r.db.WithContext(ctx).Create(record).Error
	if err != nil {
		return nil, xError.NewError(ctx, xError.DatabaseError, "创建赞助记录失败", true, err)
	}

	if cacheErr := r.cache.Set(ctx.Request.Context(), record); cacheErr != nil {
		r.log.Warn(ctx, cacheErr.Error())
	}

	detail, found, xErr := r.GetDetailByID(ctx, record.ID)
	if xErr != nil {
		return nil, xErr
	}
	if !found {
		return nil, xError.NewError(ctx, xError.NotFound, "赞助记录不存在", false)
	}
	return detail, nil
}

func (r *SponsorRecordRepo) GetByID(ctx *gin.Context, id int64) (*entity.SponsorRecord, bool, *xError.Error) {
	r.log.Info(ctx, "GetByID - 获取赞助记录")

	if record, err := r.cache.Get(ctx.Request.Context(), id); err != nil {
		return nil, false, xError.NewError(ctx, xError.CacheError, "获取赞助记录缓存失败", true, err)
	} else if record != nil {
		return record, true, nil
	}

	var record entity.SponsorRecord
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&record).Error
	if err == nil {
		if cacheErr := r.cache.Set(ctx.Request.Context(), &record); cacheErr != nil {
			r.log.Warn(ctx, cacheErr.Error())
		}
		return &record, true, nil
	}
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, false, nil
	}
	return nil, false, xError.NewError(ctx, xError.DatabaseError, "查询赞助记录失败", true, err)
}

func (r *SponsorRecordRepo) GetDetailByID(ctx *gin.Context, id int64) (*entity.SponsorRecord, bool, *xError.Error) {
	r.log.Info(ctx, "GetDetailByID - 获取赞助记录详情")

	var record entity.SponsorRecord
	err := r.db.WithContext(ctx).Preload("ChannelFKey").Where("id = ?", id).First(&record).Error
	if err == nil {
		if cacheErr := r.cache.Set(ctx.Request.Context(), &record); cacheErr != nil {
			r.log.Warn(ctx, cacheErr.Error())
		}
		return &record, true, nil
	}
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, false, nil
	}
	return nil, false, xError.NewError(ctx, xError.DatabaseError, "查询赞助记录失败", true, err)
}

func (r *SponsorRecordRepo) UpdateByID(ctx *gin.Context, id int64, updates map[string]any) (*entity.SponsorRecord, bool, *xError.Error) {
	r.log.Info(ctx, "UpdateByID - 更新赞助记录")

	result := r.db.WithContext(ctx).Model(&entity.SponsorRecord{}).Where("id = ?", id).Updates(updates)
	if result.Error != nil {
		return nil, false, xError.NewError(ctx, xError.DatabaseError, "更新赞助记录失败", true, result.Error)
	}
	if result.RowsAffected == 0 {
		return nil, false, nil
	}

	if cacheErr := r.cache.Delete(ctx.Request.Context(), id); cacheErr != nil {
		r.log.Warn(ctx, cacheErr.Error())
	}

	detail, found, xErr := r.GetDetailByID(ctx, id)
	if xErr != nil {
		return nil, false, xErr
	}
	return detail, found, nil
}

func (r *SponsorRecordRepo) HardDeleteByID(ctx *gin.Context, id int64) (bool, *xError.Error) {
	r.log.Info(ctx, "HardDeleteByID - 删除赞助记录")

	result := r.db.WithContext(ctx).Unscoped().Where("id = ?", id).Delete(&entity.SponsorRecord{})
	if result.Error != nil {
		return false, xError.NewError(ctx, xError.DatabaseError, "删除赞助记录失败", true, result.Error)
	}
	if result.RowsAffected == 0 {
		return false, nil
	}

	if cacheErr := r.cache.Delete(ctx.Request.Context(), id); cacheErr != nil {
		r.log.Warn(ctx, cacheErr.Error())
	}
	return true, nil
}

func (r *SponsorRecordRepo) Page(ctx *gin.Context, query SponsorRecordPageQuery) ([]entity.SponsorRecord, int64, *xError.Error) {
	r.log.Info(ctx, "Page - 分页查询赞助记录")

	gormQuery := r.db.WithContext(ctx).Model(&entity.SponsorRecord{})
	gormQuery = r.applyAdminFilters(gormQuery, query.ChannelID, query.Nickname, query.IsAnonymous, query.IsHidden)

	var total int64
	err := gormQuery.Count(&total).Error
	if err != nil {
		return nil, 0, xError.NewError(ctx, xError.DatabaseError, "统计赞助记录数量失败", true, err)
	}

	offset := (query.Page - 1) * query.PageSize
	gormQuery = gormQuery.Order(buildOrder(query.OrderBy, query.Order, "sort_order", "desc")).Offset(offset).Limit(query.PageSize)

	var records []entity.SponsorRecord
	err = gormQuery.Preload("ChannelFKey").Find(&records).Error
	if err != nil {
		return nil, 0, xError.NewError(ctx, xError.DatabaseError, "查询赞助记录列表失败", true, err)
	}
	return records, total, nil
}

func (r *SponsorRecordRepo) PublicPage(ctx *gin.Context, query SponsorRecordPublicPageQuery) ([]entity.SponsorRecord, int64, *xError.Error) {
	r.log.Info(ctx, "PublicPage - 分页查询公开赞助记录")

	gormQuery := r.db.WithContext(ctx).Model(&entity.SponsorRecord{}).Where("is_hidden = ?", false)
	if query.ChannelID != nil {
		gormQuery = gormQuery.Where("channel_id = ?", *query.ChannelID)
	}

	var total int64
	err := gormQuery.Count(&total).Error
	if err != nil {
		return nil, 0, xError.NewError(ctx, xError.DatabaseError, "统计赞助记录数量失败", true, err)
	}

	offset := (query.Page - 1) * query.PageSize
	gormQuery = gormQuery.Order(buildOrder(query.OrderBy, query.Order, "sort_order", "desc")).Offset(offset).Limit(query.PageSize)

	var records []entity.SponsorRecord
	err = gormQuery.Preload("ChannelFKey").Find(&records).Error
	if err != nil {
		return nil, 0, xError.NewError(ctx, xError.DatabaseError, "查询赞助记录列表失败", true, err)
	}
	return records, total, nil
}

func (r *SponsorRecordRepo) applyAdminFilters(query *gorm.DB, channelID *int64, nickname string, isAnonymous *bool, isHidden *bool) *gorm.DB {
	if channelID != nil {
		query = query.Where("channel_id = ?", *channelID)
	}
	if nickname != "" {
		query = query.Where("nickname ILIKE ?", "%"+nickname+"%")
	}
	if isAnonymous != nil {
		query = query.Where("is_anonymous = ?", *isAnonymous)
	}
	if isHidden != nil {
		query = query.Where("is_hidden = ?", *isHidden)
	}
	return query
}
