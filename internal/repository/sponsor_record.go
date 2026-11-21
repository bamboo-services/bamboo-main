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
)

// SponsorRecordRepo 赞助记录数据访问层
//
// 收口赞助记录实体的创建、查询、更新与分页查询，缓存经 cache.Manager 统一失效。
type SponsorRecordRepo struct {
	db  *gorm.DB
	kc  cache.KeyCache[string, entity.SponsorRecord]
	log *logging.LogNamedLogger
}

// SponsorRecordPageQuery 赞助记录分页查询条件
//
// repository 层自有查询模型，由 logic 层从 transport DTO 转换而来。
type SponsorRecordPageQuery struct {
	Page        int
	PageSize    int
	ChannelID   *snowflake.SnowflakeID
	Nickname    string
	IsAnonymous *bool
	IsHidden    *bool
	Status      *int
	UserID      snowflake.SnowflakeID
	OrderBy     string
	Order       string
}

// SponsorRecordPublicPageQuery 赞助记录公开分页查询条件
//
// repository 层自有查询模型，由 logic 层从 transport DTO 转换而来。
type SponsorRecordPublicPageQuery struct {
	Page      int
	PageSize  int
	ChannelID *snowflake.SnowflakeID
	OrderBy   string
	Order     string
}

// NewSponsorRecordRepo 创建 SponsorRecordRepo 实例
func NewSponsorRecordRepo(db *gorm.DB, m *cache.Manager, logger ...*logging.Logger) *SponsorRecordRepo {
	return &SponsorRecordRepo{
		db:  db,
		kc:  cache.KeyCacheOf[string, entity.SponsorRecord](m),
		log: resolveRepoLogger("SponsorRecordRepo", logger...),
	}
}

// pickDB 解析目标数据库：事务存在时优先使用事务连接，否则使用默认连接
func (r *SponsorRecordRepo) pickDB(tx *gorm.DB) *gorm.DB {
	if tx != nil {
		return tx
	}
	return r.db
}

// Create 创建赞助记录
func (r *SponsorRecordRepo) Create(ctx context.Context, record *entity.SponsorRecord, tx *gorm.DB) (*entity.SponsorRecord, *apperr.Error) {
	r.log.Info(ctx, "Create - 创建赞助记录")

	err := r.pickDB(tx).WithContext(ctx).Create(record).Error
	if err != nil {
		return nil, apperr.New(database.DatabaseError, "创建赞助记录失败", err)
	}

	if cacheErr := r.kc.Set(ctx, constants.RedisSponsorRecord.Get(record.ID).String(), record, cache.WithTTL(10*time.Minute)); cacheErr != nil {
		r.log.Warn(ctx, cacheErr.Error())
	}

	detail, found, xErr := r.GetDetailByID(ctx, record.ID)
	if xErr != nil {
		return nil, xErr
	}
	if !found {
		return nil, apperr.New(httpresult.NotFound, "赞助记录不存在")
	}
	return detail, nil
}

// GetByID 根据ID获取赞助记录
func (r *SponsorRecordRepo) GetByID(ctx context.Context, id snowflake.SnowflakeID) (*entity.SponsorRecord, bool, *apperr.Error) {
	r.log.Info(ctx, "GetByID - 获取赞助记录")

	if record, ok, err := r.kc.Get(ctx, constants.RedisSponsorRecord.Get(id).String()); err != nil {
		return nil, false, apperr.New(cache.CacheError, "获取赞助记录缓存失败", err)
	} else if ok {
		return record, true, nil
	}

	var record entity.SponsorRecord
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&record).Error
	if err == nil {
		if cacheErr := r.kc.Set(ctx, constants.RedisSponsorRecord.Get(record.ID).String(), &record, cache.WithTTL(10*time.Minute)); cacheErr != nil {
			r.log.Warn(ctx, cacheErr.Error())
		}
		return &record, true, nil
	}
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, false, nil
	}
	return nil, false, apperr.New(database.DatabaseError, "查询赞助记录失败", err)
}

// GetDetailByID 根据ID获取赞助记录详情（含渠道关联）
func (r *SponsorRecordRepo) GetDetailByID(ctx context.Context, id snowflake.SnowflakeID) (*entity.SponsorRecord, bool, *apperr.Error) {
	r.log.Info(ctx, "GetDetailByID - 获取赞助记录详情")

	var record entity.SponsorRecord
	err := r.db.WithContext(ctx).Preload("ChannelFKey").Where("id = ?", id).First(&record).Error
	if err == nil {
		if cacheErr := r.kc.Set(ctx, constants.RedisSponsorRecord.Get(record.ID).String(), &record, cache.WithTTL(10*time.Minute)); cacheErr != nil {
			r.log.Warn(ctx, cacheErr.Error())
		}
		return &record, true, nil
	}
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, false, nil
	}
	return nil, false, apperr.New(database.DatabaseError, "查询赞助记录失败", err)
}

// UpdateByID 根据ID更新赞助记录的指定字段
func (r *SponsorRecordRepo) UpdateByID(ctx context.Context, id snowflake.SnowflakeID, updates map[string]any, tx *gorm.DB) (*entity.SponsorRecord, bool, *apperr.Error) {
	r.log.Info(ctx, "UpdateByID - 更新赞助记录")

	result := r.pickDB(tx).WithContext(ctx).Model(&entity.SponsorRecord{}).Where("id = ?", id).Updates(updates)
	if result.Error != nil {
		return nil, false, apperr.New(database.DatabaseError, "更新赞助记录失败", result.Error)
	}
	if result.RowsAffected == 0 {
		return nil, false, nil
	}

	if cacheErr := r.kc.Delete(ctx, constants.RedisSponsorRecord.Get(id).String()); cacheErr != nil {
		r.log.Warn(ctx, cacheErr.Error())
	}

	detail, found, xErr := r.GetDetailByID(ctx, id)
	if xErr != nil {
		return nil, false, xErr
	}
	return detail, found, nil
}

// HardDeleteByID 物理删除赞助记录
func (r *SponsorRecordRepo) HardDeleteByID(ctx context.Context, id snowflake.SnowflakeID, tx *gorm.DB) (bool, *apperr.Error) {
	r.log.Info(ctx, "HardDeleteByID - 删除赞助记录")

	result := r.pickDB(tx).WithContext(ctx).Unscoped().Where("id = ?", id).Delete(&entity.SponsorRecord{})
	if result.Error != nil {
		return false, apperr.New(database.DatabaseError, "删除赞助记录失败", result.Error)
	}
	if result.RowsAffected == 0 {
		return false, nil
	}

	if cacheErr := r.kc.Delete(ctx, constants.RedisSponsorRecord.Get(id).String()); cacheErr != nil {
		r.log.Warn(ctx, cacheErr.Error())
	}
	return true, nil
}

// Page 分页查询赞助记录
func (r *SponsorRecordRepo) Page(ctx context.Context, query SponsorRecordPageQuery) ([]entity.SponsorRecord, int64, *apperr.Error) {
	r.log.Info(ctx, "Page - 分页查询赞助记录")

	gormQuery := r.db.WithContext(ctx).Model(&entity.SponsorRecord{})
	gormQuery = r.applyAdminFilters(gormQuery, query.ChannelID, query.Nickname, query.IsAnonymous, query.IsHidden)
	// 审核状态过滤（管理端按状态筛选、用户中心按归属查询）
	if query.Status != nil {
		gormQuery = gormQuery.Where("status = ?", *query.Status)
	}
	if query.UserID != 0 {
		gormQuery = gormQuery.Where("user_id = ?", query.UserID)
	}

	var total int64
	err := gormQuery.Count(&total).Error
	if err != nil {
		return nil, 0, apperr.New(database.DatabaseError, "统计赞助记录数量失败", err)
	}

	offset := (query.Page - 1) * query.PageSize
	gormQuery = gormQuery.Order(buildOrder(query.OrderBy, query.Order, "sort_order", "desc")).Offset(offset).Limit(query.PageSize)

	var records []entity.SponsorRecord
	err = gormQuery.Preload("ChannelFKey").Find(&records).Error
	if err != nil {
		return nil, 0, apperr.New(database.DatabaseError, "查询赞助记录列表失败", err)
	}
	return records, total, nil
}

// PublicPage 分页查询公开的赞助记录
func (r *SponsorRecordRepo) PublicPage(ctx context.Context, query SponsorRecordPublicPageQuery) ([]entity.SponsorRecord, int64, *apperr.Error) {
	r.log.Info(ctx, "PublicPage - 分页查询公开赞助记录")

	// 仅展示已通过审核且未被隐藏的记录
	gormQuery := r.db.WithContext(ctx).Model(&entity.SponsorRecord{}).
		Where("status = ? AND is_hidden = ?", constants.SponsorStatusApproved, false)
	if query.ChannelID != nil {
		gormQuery = gormQuery.Where("channel_id = ?", *query.ChannelID)
	}

	var total int64
	err := gormQuery.Count(&total).Error
	if err != nil {
		return nil, 0, apperr.New(database.DatabaseError, "统计赞助记录数量失败", err)
	}

	offset := (query.Page - 1) * query.PageSize
	gormQuery = gormQuery.Order(buildOrder(query.OrderBy, query.Order, "sort_order", "desc")).Offset(offset).Limit(query.PageSize)

	var records []entity.SponsorRecord
	err = gormQuery.Preload("ChannelFKey").Find(&records).Error
	if err != nil {
		return nil, 0, apperr.New(database.DatabaseError, "查询赞助记录列表失败", err)
	}
	return records, total, nil
}

func (r *SponsorRecordRepo) applyAdminFilters(query *gorm.DB, channelID *snowflake.SnowflakeID, nickname string, isAnonymous *bool, isHidden *bool) *gorm.DB {
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

// UpdateStatusByID 更新赞助记录的审核状态
func (r *SponsorRecordRepo) UpdateStatusByID(ctx context.Context, id snowflake.SnowflakeID, status int, reviewRemark string, tx *gorm.DB) (bool, *apperr.Error) {
	r.log.Info(ctx, "UpdateStatusByID - 更新赞助记录状态")

	updates := map[string]any{
		"status":        status,
		"review_remark": reviewRemark,
	}

	result := r.pickDB(tx).WithContext(ctx).Model(&entity.SponsorRecord{}).Where("id = ?", id).Updates(updates)
	if result.Error != nil {
		return false, apperr.New(database.DatabaseError, "更新赞助记录状态失败", result.Error)
	}
	if result.RowsAffected == 0 {
		return false, nil
	}

	if cacheErr := r.kc.Delete(ctx, constants.RedisSponsorRecord.Get(id).String()); cacheErr != nil {
		r.log.Warn(ctx, cacheErr.Error())
	}

	return true, nil
}

// BindUserByEmail 将指定邮箱名下尚未归属的赞助记录绑定到用户
//
// 用于游客申请赞助后注册/登录时的归属关联：把所有 user_id 为空且联系邮箱匹配的赞助记录归属到该用户。
func (r *SponsorRecordRepo) BindUserByEmail(ctx context.Context, userID snowflake.SnowflakeID, email string, tx *gorm.DB) *apperr.Error {
	r.log.Info(ctx, "BindUserByEmail - 按邮箱绑定赞助记录归属")

	email = strings.ToLower(strings.TrimSpace(email))
	if email == "" {
		return nil
	}

	result := r.pickDB(tx).WithContext(ctx).Model(&entity.SponsorRecord{}).
		Where("user_id IS NULL AND lower(email) = ?", email).
		Update("user_id", userID)
	if result.Error != nil {
		return apperr.New(database.DatabaseError, "绑定赞助记录归属失败", result.Error)
	}
	return nil
}

// CountByStatus 统计指定审核状态的赞助记录数量
func (r *SponsorRecordRepo) CountByStatus(ctx context.Context, status int, tx *gorm.DB) (int64, *apperr.Error) {
	r.log.Info(ctx, "CountByStatus - 统计赞助记录状态数量")

	var count int64
	err := r.pickDB(tx).WithContext(ctx).Model(&entity.SponsorRecord{}).Where("status = ?", status).Count(&count).Error
	if err != nil {
		return 0, apperr.New(database.DatabaseError, "统计赞助记录数量失败", err)
	}
	return count, nil
}
