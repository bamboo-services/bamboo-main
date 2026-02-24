package repository

import (
	"errors"
	"time"

	xError "github.com/bamboo-services/bamboo-base-go/common/error"
	xLog "github.com/bamboo-services/bamboo-base-go/common/log"
	"github.com/bamboo-services/bamboo-main/internal/entity"
	"github.com/bamboo-services/bamboo-main/internal/repository/cache"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type SponsorChannelRepo struct {
	db    *gorm.DB
	cache *cache.SponsorChannelCache
	log   *xLog.LogNamedLogger
}

type SponsorChannelListQuery struct {
	Status      *bool
	OnlyEnabled bool
	Name        string
	OrderBy     string
	Order       string
}

type SponsorChannelPageQuery struct {
	Page     int
	PageSize int
	Status   *bool
	Name     string
	OrderBy  string
	Order    string
}

func NewSponsorChannelRepo(db *gorm.DB, rdb *redis.Client) *SponsorChannelRepo {
	return &SponsorChannelRepo{
		db: db,
		cache: &cache.SponsorChannelCache{
			RDB: rdb,
			TTL: time.Minute * 10,
		},
		log: xLog.WithName(xLog.NamedREPO, "SponsorChannelRepo"),
	}
}

func (r *SponsorChannelRepo) Create(ctx *gin.Context, channel *entity.SponsorChannel) (*entity.SponsorChannel, *xError.Error) {
	r.log.Info(ctx, "Create - 创建赞助渠道")

	err := r.db.WithContext(ctx).Create(channel).Error
	if err != nil {
		return nil, xError.NewError(ctx, xError.DatabaseError, "创建赞助渠道失败", true, err)
	}

	if cacheErr := r.cache.Set(ctx.Request.Context(), channel); cacheErr != nil {
		r.log.Warn(ctx, cacheErr.Error())
	}
	return channel, nil
}

func (r *SponsorChannelRepo) GetByID(ctx *gin.Context, id int64) (*entity.SponsorChannel, bool, *xError.Error) {
	r.log.Info(ctx, "GetByID - 获取赞助渠道")

	if channel, err := r.cache.Get(ctx.Request.Context(), id); err != nil {
		return nil, false, xError.NewError(ctx, xError.CacheError, "获取赞助渠道缓存失败", true, err)
	} else if channel != nil {
		return channel, true, nil
	}

	var channel entity.SponsorChannel
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&channel).Error
	if err == nil {
		if cacheErr := r.cache.Set(ctx.Request.Context(), &channel); cacheErr != nil {
			r.log.Warn(ctx, cacheErr.Error())
		}
		return &channel, true, nil
	}
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, false, nil
	}
	return nil, false, xError.NewError(ctx, xError.DatabaseError, "查询赞助渠道失败", true, err)
}

func (r *SponsorChannelRepo) UpdateByID(ctx *gin.Context, id int64, updates map[string]any) (*entity.SponsorChannel, bool, *xError.Error) {
	r.log.Info(ctx, "UpdateByID - 更新赞助渠道")

	result := r.db.WithContext(ctx).Model(&entity.SponsorChannel{}).Where("id = ?", id).Updates(updates)
	if result.Error != nil {
		return nil, false, xError.NewError(ctx, xError.DatabaseError, "更新赞助渠道失败", true, result.Error)
	}
	if result.RowsAffected == 0 {
		return nil, false, nil
	}

	if cacheErr := r.cache.Delete(ctx.Request.Context(), id); cacheErr != nil {
		r.log.Warn(ctx, cacheErr.Error())
	}

	channel, found, xErr := r.GetByID(ctx, id)
	if xErr != nil {
		return nil, false, xErr
	}
	return channel, found, nil
}

func (r *SponsorChannelRepo) HardDeleteByID(ctx *gin.Context, id int64) (bool, *xError.Error) {
	r.log.Info(ctx, "HardDeleteByID - 删除赞助渠道")

	result := r.db.WithContext(ctx).Unscoped().Where("id = ?", id).Delete(&entity.SponsorChannel{})
	if result.Error != nil {
		return false, xError.NewError(ctx, xError.DatabaseError, "删除赞助渠道失败", true, result.Error)
	}
	if result.RowsAffected == 0 {
		return false, nil
	}

	if cacheErr := r.cache.Delete(ctx.Request.Context(), id); cacheErr != nil {
		r.log.Warn(ctx, cacheErr.Error())
	}
	return true, nil
}

func (r *SponsorChannelRepo) CountSponsorsByChannelID(ctx *gin.Context, channelID int64) (int64, *xError.Error) {
	r.log.Info(ctx, "CountSponsorsByChannelID - 统计赞助记录数量")

	var total int64
	err := r.db.WithContext(ctx).Model(&entity.SponsorRecord{}).Where("channel_id = ?", channelID).Count(&total).Error
	if err != nil {
		return 0, xError.NewError(ctx, xError.DatabaseError, "查询赞助数量失败", true, err)
	}
	return total, nil
}

func (r *SponsorChannelRepo) CountSponsorsByChannelIDs(ctx *gin.Context, channelIDs []int64) (map[int64]int64, *xError.Error) {
	r.log.Info(ctx, "CountSponsorsByChannelIDs - 批量统计赞助记录数量")

	sponsorCounts := make(map[int64]int64)
	if len(channelIDs) == 0 {
		return sponsorCounts, nil
	}

	var countResults []struct {
		ChannelID int64 `gorm:"column:channel_id"`
		Count     int64 `gorm:"column:count"`
	}

	err := r.db.WithContext(ctx).
		Model(&entity.SponsorRecord{}).
		Select("channel_id, COUNT(*) as count").
		Where("channel_id IN ?", channelIDs).
		Group("channel_id").
		Find(&countResults).Error
	if err != nil {
		return nil, xError.NewError(ctx, xError.DatabaseError, "查询赞助数量失败", true, err)
	}

	for _, result := range countResults {
		sponsorCounts[result.ChannelID] = result.Count
	}
	return sponsorCounts, nil
}

func (r *SponsorChannelRepo) List(ctx *gin.Context, query SponsorChannelListQuery) ([]entity.SponsorChannel, *xError.Error) {
	r.log.Info(ctx, "List - 查询赞助渠道列表")

	gormQuery := r.applyListFilters(r.db.WithContext(ctx).Model(&entity.SponsorChannel{}), query.Status, query.OnlyEnabled, query.Name)
	gormQuery = gormQuery.Order(buildOrder(query.OrderBy, query.Order, "sort_order", "asc"))

	var channels []entity.SponsorChannel
	err := gormQuery.Find(&channels).Error
	if err != nil {
		return nil, xError.NewError(ctx, xError.DatabaseError, "查询赞助渠道列表失败", true, err)
	}
	return channels, nil
}

func (r *SponsorChannelRepo) Page(ctx *gin.Context, query SponsorChannelPageQuery) ([]entity.SponsorChannel, int64, *xError.Error) {
	r.log.Info(ctx, "Page - 分页查询赞助渠道")

	gormQuery := r.applyListFilters(r.db.WithContext(ctx).Model(&entity.SponsorChannel{}), query.Status, false, query.Name)

	var total int64
	err := gormQuery.Count(&total).Error
	if err != nil {
		return nil, 0, xError.NewError(ctx, xError.DatabaseError, "统计赞助渠道数量失败", true, err)
	}

	offset := (query.Page - 1) * query.PageSize
	gormQuery = gormQuery.Order(buildOrder(query.OrderBy, query.Order, "sort_order", "asc")).Offset(offset).Limit(query.PageSize)

	var channels []entity.SponsorChannel
	err = gormQuery.Find(&channels).Error
	if err != nil {
		return nil, 0, xError.NewError(ctx, xError.DatabaseError, "查询赞助渠道列表失败", true, err)
	}
	return channels, total, nil
}

func (r *SponsorChannelRepo) PublicList(ctx *gin.Context) ([]entity.SponsorChannel, *xError.Error) {
	r.log.Info(ctx, "PublicList - 查询公开赞助渠道列表")

	var channels []entity.SponsorChannel
	err := r.db.WithContext(ctx).Model(&entity.SponsorChannel{}).Where("status = ?", true).Order("sort_order asc").Find(&channels).Error
	if err != nil {
		return nil, xError.NewError(ctx, xError.DatabaseError, "查询公开渠道列表失败", true, err)
	}
	return channels, nil
}

func (r *SponsorChannelRepo) applyListFilters(query *gorm.DB, status *bool, onlyEnabled bool, name string) *gorm.DB {
	if status != nil {
		query = query.Where("status = ?", *status)
	}
	if onlyEnabled {
		query = query.Where("status = ?", true)
	}
	if name != "" {
		query = query.Where("name ILIKE ?", "%"+name+"%")
	}
	return query
}

func buildOrder(orderBy string, order string, defaultBy string, defaultOrder string) string {
	if orderBy == "" {
		orderBy = defaultBy
	}
	if order == "" {
		order = defaultOrder
	}
	return orderBy + " " + order
}
