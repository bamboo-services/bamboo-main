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
	"errors"
	"strings"
	"time"

	xError "github.com/bamboo-services/bamboo-base-go/common/error"
	xLog "github.com/bamboo-services/bamboo-base-go/common/log"
	xSnowflake "github.com/bamboo-services/bamboo-base-go/common/snowflake"
	xCache "github.com/bamboo-services/bamboo-base-go/major/cache"
	apiLink "github.com/bamboo-services/bamboo-main/api/link"
	"github.com/bamboo-services/bamboo-main/internal/entity"
	"github.com/bamboo-services/bamboo-main/pkg/constants"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type LinkRepo struct {
	db  *gorm.DB
	kc  xCache.KeyCache[string, entity.LinkFriend]
	log *xLog.LogNamedLogger
}

func NewLinkRepo(db *gorm.DB, m *xCache.Manager) *LinkRepo {
	return &LinkRepo{
		db:  db,
		kc:  xCache.KeyCacheOf[string, entity.LinkFriend](m),
		log: xLog.WithName(xLog.NamedREPO, "LinkRepo"),
	}
}

func (r *LinkRepo) pickDB(tx *gorm.DB) *gorm.DB {
	if tx != nil {
		return tx
	}
	return r.db
}

func (r *LinkRepo) Create(ctx *gin.Context, link *entity.LinkFriend, tx *gorm.DB) (*entity.LinkFriend, *xError.Error) {
	r.log.Info(ctx, "Create - 创建友情链接")

	err := r.pickDB(tx).WithContext(ctx).Create(link).Error
	if err != nil {
		return nil, xError.NewError(ctx, xError.DatabaseError, "创建友情链接失败", false, err)
	}

	if cacheErr := r.kc.Set(ctx.Request.Context(), constants.RedisLinkFriend.Get(link.ID).String(), link, xCache.WithTTL(15*time.Minute)); cacheErr != nil {
		r.log.Warn(ctx, cacheErr.Error())
	}

	return link, nil
}

func (r *LinkRepo) Save(ctx *gin.Context, link *entity.LinkFriend, tx *gorm.DB) (*entity.LinkFriend, *xError.Error) {
	r.log.Info(ctx, "Save - 保存友情链接")

	err := r.pickDB(tx).WithContext(ctx).Save(link).Error
	if err != nil {
		return nil, xError.NewError(ctx, xError.DatabaseError, "保存友情链接失败", false, err)
	}

	if cacheErr := r.kc.Delete(ctx.Request.Context(), constants.RedisLinkFriend.Get(link.ID).String()); cacheErr != nil {
		r.log.Warn(ctx, cacheErr.Error())
	}

	return link, nil
}

func (r *LinkRepo) GetByID(ctx *gin.Context, id xSnowflake.SnowflakeID, withAssociations bool, tx *gorm.DB) (*entity.LinkFriend, bool, *xError.Error) {
	r.log.Info(ctx, "GetByID - 获取友情链接")

	if !withAssociations {
		if link, ok, err := r.kc.Get(ctx.Request.Context(), constants.RedisLinkFriend.Get(id).String()); err != nil {
			return nil, false, xError.NewError(ctx, xError.CacheError, "获取友情链接缓存失败", true, err)
		} else if ok {
			return link, true, nil
		}
	}

	query := r.pickDB(tx).WithContext(ctx)
	if withAssociations {
		query = query.Preload("GroupFKey").Preload("ColorFKey")
	}

	var link entity.LinkFriend
	err := query.Where("id = ?", id).First(&link).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, false, nil
		}
		return nil, false, xError.NewError(ctx, xError.DatabaseError, "查询友情链接失败", false, err)
	}

	if cacheErr := r.kc.Set(ctx.Request.Context(), constants.RedisLinkFriend.Get(link.ID).String(), &link, xCache.WithTTL(15*time.Minute)); cacheErr != nil {
		r.log.Warn(ctx, cacheErr.Error())
	}

	return &link, true, nil
}

func (r *LinkRepo) DeleteByID(ctx *gin.Context, id xSnowflake.SnowflakeID, tx *gorm.DB) (bool, *xError.Error) {
	r.log.Info(ctx, "DeleteByID - 删除友情链接")

	result := r.pickDB(tx).WithContext(ctx).Where("id = ?", id).Delete(&entity.LinkFriend{})
	if result.Error != nil {
		return false, xError.NewError(ctx, xError.DatabaseError, "删除友情链接失败", false, result.Error)
	}
	if result.RowsAffected == 0 {
		return false, nil
	}

	if cacheErr := r.kc.Delete(ctx.Request.Context(), constants.RedisLinkFriend.Get(id).String()); cacheErr != nil {
		r.log.Warn(ctx, cacheErr.Error())
	}

	return true, nil
}

func (r *LinkRepo) List(ctx *gin.Context, req *apiLink.FriendQueryRequest, tx *gorm.DB) ([]entity.LinkFriend, int64, *xError.Error) {
	r.log.Info(ctx, "List - 查询友情链接分页列表")

	query := r.pickDB(tx).WithContext(ctx).Model(&entity.LinkFriend{})

	if req.LinkName != "" {
		query = query.Where("name ILIKE ?", "%"+req.LinkName+"%")
	}
	if req.LinkStatus != nil {
		query = query.Where("status = ?", *req.LinkStatus)
	}
	if req.LinkFail != nil {
		query = query.Where("is_failure = ?", *req.LinkFail)
	}
	if req.LinkGroupID != 0 {
		query = query.Where("group_id = ?", req.LinkGroupID)
	}

	orderBy := "created_at"
	switch req.SortBy {
	case "created_at":
		orderBy = "created_at"
	case "updated_at":
		orderBy = "updated_at"
	case "link_order":
		orderBy = "sort_order"
	case "link_name":
		orderBy = "name"
	}

	sortOrder := "DESC"
	if strings.EqualFold(req.SortOrder, "asc") {
		sortOrder = "ASC"
	}

	query = query.Order(orderBy + " " + sortOrder)

	var total int64
	err := query.Count(&total).Error
	if err != nil {
		return nil, 0, xError.NewError(ctx, xError.DatabaseError, "统计友情链接数量失败", true, err)
	}

	offset := (req.Page - 1) * req.PageSize
	var links []entity.LinkFriend
	err = query.Preload("GroupFKey").Preload("ColorFKey").Offset(offset).Limit(req.PageSize).Find(&links).Error
	if err != nil {
		return nil, 0, xError.NewError(ctx, xError.DatabaseError, "查询友情链接列表失败", false, err)
	}

	return links, total, nil
}

func (r *LinkRepo) UpdateStatusByID(ctx *gin.Context, id xSnowflake.SnowflakeID, status int, reviewRemark string, tx *gorm.DB) (bool, *xError.Error) {
	r.log.Info(ctx, "UpdateStatusByID - 更新友情链接状态")

	updates := map[string]any{
		"status":        status,
		"review_remark": reviewRemark,
	}

	result := r.pickDB(tx).WithContext(ctx).Model(&entity.LinkFriend{}).Where("id = ?", id).Updates(updates)
	if result.Error != nil {
		return false, xError.NewError(ctx, xError.DatabaseError, "更新友情链接状态失败", false, result.Error)
	}
	if result.RowsAffected == 0 {
		return false, nil
	}

	if cacheErr := r.kc.Delete(ctx.Request.Context(), constants.RedisLinkFriend.Get(id).String()); cacheErr != nil {
		r.log.Warn(ctx, cacheErr.Error())
	}

	return true, nil
}

func (r *LinkRepo) UpdateFailureByID(ctx *gin.Context, id xSnowflake.SnowflakeID, isFailure int, failReason string, tx *gorm.DB) (bool, *xError.Error) {
	r.log.Info(ctx, "UpdateFailureByID - 更新友情链接失效状态")

	updates := map[string]any{
		"is_failure":  isFailure,
		"fail_reason": failReason,
	}

	result := r.pickDB(tx).WithContext(ctx).Model(&entity.LinkFriend{}).Where("id = ?", id).Updates(updates)
	if result.Error != nil {
		return false, xError.NewError(ctx, xError.DatabaseError, "更新友情链接失效状态失败", false, result.Error)
	}
	if result.RowsAffected == 0 {
		return false, nil
	}

	if cacheErr := r.kc.Delete(ctx.Request.Context(), constants.RedisLinkFriend.Get(id).String()); cacheErr != nil {
		r.log.Warn(ctx, cacheErr.Error())
	}

	return true, nil
}

func (r *LinkRepo) ListPublic(ctx *gin.Context, groupID *xSnowflake.SnowflakeID, approvedStatus int, normalFail int, tx *gorm.DB) ([]entity.LinkFriend, *xError.Error) {
	r.log.Info(ctx, "ListPublic - 查询公开友情链接")

	query := r.pickDB(tx).WithContext(ctx).
		Where("status = ? AND is_failure = ?", approvedStatus, normalFail)

	if groupID != nil {
		query = query.Where("group_id = ?", *groupID)
	}

	var links []entity.LinkFriend
	err := query.Preload("GroupFKey").Preload("ColorFKey").Order("sort_order ASC, created_at DESC").Find(&links).Error
	if err != nil {
		return nil, xError.NewError(ctx, xError.DatabaseError, "查询公开友情链接失败", false, err)
	}

	return links, nil
}

func (r *LinkRepo) CountByGroupID(ctx *gin.Context, groupID xSnowflake.SnowflakeID, tx *gorm.DB) (int64, *xError.Error) {
	var count int64
	err := r.pickDB(tx).WithContext(ctx).Model(&entity.LinkFriend{}).Where("group_id = ?", groupID).Count(&count).Error
	if err != nil {
		return 0, xError.NewError(ctx, xError.DatabaseError, "查询关联友链失败", false, err)
	}
	return count, nil
}

func (r *LinkRepo) CountByColorID(ctx *gin.Context, colorID xSnowflake.SnowflakeID, tx *gorm.DB) (int64, *xError.Error) {
	var count int64
	err := r.pickDB(tx).WithContext(ctx).Model(&entity.LinkFriend{}).Where("color_id = ?", colorID).Count(&count).Error
	if err != nil {
		return 0, xError.NewError(ctx, xError.DatabaseError, "查询关联友链失败", false, err)
	}
	return count, nil
}

func (r *LinkRepo) ListByGroupID(ctx *gin.Context, groupID xSnowflake.SnowflakeID, limit int, tx *gorm.DB) ([]entity.LinkFriend, *xError.Error) {
	query := r.pickDB(tx).WithContext(ctx).Where("group_id = ?", groupID)
	if limit > 0 {
		query = query.Limit(limit)
	}

	var links []entity.LinkFriend
	err := query.Find(&links).Error
	if err != nil {
		return nil, xError.NewError(ctx, xError.DatabaseError, "查询冲突友链失败", false, err)
	}

	return links, nil
}

func (r *LinkRepo) ListByColorID(ctx *gin.Context, colorID xSnowflake.SnowflakeID, limit int, tx *gorm.DB) ([]entity.LinkFriend, *xError.Error) {
	query := r.pickDB(tx).WithContext(ctx).Where("color_id = ?", colorID)
	if limit > 0 {
		query = query.Limit(limit)
	}

	var links []entity.LinkFriend
	err := query.Find(&links).Error
	if err != nil {
		return nil, xError.NewError(ctx, xError.DatabaseError, "查询冲突友链失败", false, err)
	}

	return links, nil
}

func (r *LinkRepo) ClearGroupID(ctx *gin.Context, groupID xSnowflake.SnowflakeID, tx *gorm.DB) *xError.Error {
	result := r.pickDB(tx).WithContext(ctx).Model(&entity.LinkFriend{}).Where("group_id = ?", groupID).Update("group_id", nil)
	if result.Error != nil {
		return xError.NewError(ctx, xError.DatabaseError, "清空友链分组关联失败", false, result.Error)
	}
	return nil
}

func (r *LinkRepo) ClearColorID(ctx *gin.Context, colorID xSnowflake.SnowflakeID, tx *gorm.DB) *xError.Error {
	result := r.pickDB(tx).WithContext(ctx).Model(&entity.LinkFriend{}).Where("color_id = ?", colorID).Update("color_id", nil)
	if result.Error != nil {
		return xError.NewError(ctx, xError.DatabaseError, "清空友链颜色关联失败", false, result.Error)
	}
	return nil
}
