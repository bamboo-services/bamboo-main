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

	xError "github.com/bamboo-services/bamboo-base-go/common/error"
	xLog "github.com/bamboo-services/bamboo-base-go/common/log"
	xSnowflake "github.com/bamboo-services/bamboo-base-go/common/snowflake"
	xCache "github.com/bamboo-services/bamboo-base-go/major/cache"
	"github.com/bamboo-services/bamboo-main/internal/entity"
	"github.com/bamboo-services/bamboo-main/pkg/constants"
	"gorm.io/gorm"
)

// LinkRepo 友情链接数据访问层
//
// 收口友情链接实体的 CRUD、分页查询与缓存读写，缓存经 xCache.Manager 统一失效。
type LinkRepo struct {
	db  *gorm.DB
	kc  xCache.KeyCache[string, entity.LinkFriend]
	log *xLog.LogNamedLogger
}

// FriendQuery 友情链接分页查询条件
//
// repository 层自有查询模型，由 logic 层从 transport DTO（api/link）转换而来，
// 避免数据访问层直接耦合传输层契约。
type FriendQuery struct {
	Page        int
	PageSize    int
	LinkName    string
	LinkStatus  *int
	LinkFail    *int
	LinkGroupID xSnowflake.SnowflakeID
	UserID      xSnowflake.SnowflakeID
	SortBy      string
	SortOrder   string
}

// NewLinkRepo 创建 LinkRepo 实例
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

// Create 创建友情链接
func (r *LinkRepo) Create(ctx context.Context, link *entity.LinkFriend, tx *gorm.DB) (*entity.LinkFriend, *xError.Error) {
	r.log.Info(ctx, "Create - 创建友情链接")

	err := r.pickDB(tx).WithContext(ctx).Create(link).Error
	if err != nil {
		return nil, xError.NewError(ctx, xError.DatabaseError, "创建友情链接失败", false, err)
	}

	if cacheErr := r.kc.Set(ctx, constants.RedisLinkFriend.Get(link.ID).String(), link, xCache.WithTTL(15*time.Minute)); cacheErr != nil {
		r.log.Warn(ctx, cacheErr.Error())
	}

	return link, nil
}

// BindUserByEmail 将指定邮箱名下尚未归属的友链绑定到用户
//
// 用于游客提交友链后注册/登录时的归属关联：把所有 user_id 为空且联系邮箱匹配的友链归属到该用户。
func (r *LinkRepo) BindUserByEmail(ctx context.Context, userID xSnowflake.SnowflakeID, email string) *xError.Error {
	r.log.Info(ctx, "BindUserByEmail - 按邮箱绑定友链归属")

	email = strings.ToLower(strings.TrimSpace(email))
	if email == "" {
		return nil
	}

	result := r.db.WithContext(ctx).Model(&entity.LinkFriend{}).
		Where("user_id IS NULL AND lower(email) = ?", email).
		Update("user_id", userID)
	if result.Error != nil {
		return xError.NewError(ctx, xError.DatabaseError, "绑定友链归属失败", false, result.Error)
	}
	return nil
}

// Save 保存友情链接（存在则更新）
func (r *LinkRepo) Save(ctx context.Context, link *entity.LinkFriend, tx *gorm.DB) (*entity.LinkFriend, *xError.Error) {
	r.log.Info(ctx, "Save - 保存友情链接")

	err := r.pickDB(tx).WithContext(ctx).Save(link).Error
	if err != nil {
		return nil, xError.NewError(ctx, xError.DatabaseError, "保存友情链接失败", false, err)
	}

	if cacheErr := r.kc.Delete(ctx, constants.RedisLinkFriend.Get(link.ID).String()); cacheErr != nil {
		r.log.Warn(ctx, cacheErr.Error())
	}

	return link, nil
}

// GetByID 根据ID获取友情链接
func (r *LinkRepo) GetByID(ctx context.Context, id xSnowflake.SnowflakeID, withAssociations bool, tx *gorm.DB) (*entity.LinkFriend, bool, *xError.Error) {
	r.log.Info(ctx, "GetByID - 获取友情链接")

	if !withAssociations {
		if link, ok, err := r.kc.Get(ctx, constants.RedisLinkFriend.Get(id).String()); err != nil {
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

	if cacheErr := r.kc.Set(ctx, constants.RedisLinkFriend.Get(link.ID).String(), &link, xCache.WithTTL(15*time.Minute)); cacheErr != nil {
		r.log.Warn(ctx, cacheErr.Error())
	}

	return &link, true, nil
}

// DeleteByID 删除友情链接
func (r *LinkRepo) DeleteByID(ctx context.Context, id xSnowflake.SnowflakeID, tx *gorm.DB) (bool, *xError.Error) {
	r.log.Info(ctx, "DeleteByID - 删除友情链接")

	result := r.pickDB(tx).WithContext(ctx).Where("id = ?", id).Delete(&entity.LinkFriend{})
	if result.Error != nil {
		return false, xError.NewError(ctx, xError.DatabaseError, "删除友情链接失败", false, result.Error)
	}
	if result.RowsAffected == 0 {
		return false, nil
	}

	if cacheErr := r.kc.Delete(ctx, constants.RedisLinkFriend.Get(id).String()); cacheErr != nil {
		r.log.Warn(ctx, cacheErr.Error())
	}

	return true, nil
}

// List 分页查询友情链接
func (r *LinkRepo) List(ctx context.Context, req *FriendQuery, tx *gorm.DB) ([]entity.LinkFriend, int64, *xError.Error) {
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
	if req.UserID != 0 {
		query = query.Where("user_id = ?", req.UserID)
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

// UpdateStatusByID 更新友情链接的审核状态
func (r *LinkRepo) UpdateStatusByID(ctx context.Context, id xSnowflake.SnowflakeID, status int, reviewRemark string, tx *gorm.DB) (bool, *xError.Error) {
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

	if cacheErr := r.kc.Delete(ctx, constants.RedisLinkFriend.Get(id).String()); cacheErr != nil {
		r.log.Warn(ctx, cacheErr.Error())
	}

	return true, nil
}

// UpdateFailureByID 更新友情链接的失效状态
func (r *LinkRepo) UpdateFailureByID(ctx context.Context, id xSnowflake.SnowflakeID, isFailure int, failReason string, tx *gorm.DB) (bool, *xError.Error) {
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

	if cacheErr := r.kc.Delete(ctx, constants.RedisLinkFriend.Get(id).String()); cacheErr != nil {
		r.log.Warn(ctx, cacheErr.Error())
	}

	return true, nil
}

// ListPublic 查询公开的友情链接列表
func (r *LinkRepo) ListPublic(ctx context.Context, groupID *xSnowflake.SnowflakeID, approvedStatus int, normalFail int, tx *gorm.DB) ([]entity.LinkFriend, *xError.Error) {
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

// ListApprovedForScreenshot 查询全部已通过且未失效的友链（按排序与创建时间，供截图全量入队）
func (r *LinkRepo) ListApprovedForScreenshot(ctx context.Context, tx *gorm.DB) ([]entity.LinkFriend, *xError.Error) {
	r.log.Info(ctx, "ListApprovedForScreenshot - 查询待截图友链")

	var links []entity.LinkFriend
	err := r.pickDB(tx).WithContext(ctx).
		Where("status = ? AND is_failure = ?", constants.LinkStatusApproved, constants.LinkFailNormal).
		Order("sort_order ASC, created_at DESC").
		Find(&links).Error
	if err != nil {
		return nil, xError.NewError(ctx, xError.DatabaseError, "查询待截图友链失败", false, err)
	}

	return links, nil
}

// UpdateScreenshot 更新友链站点截图信息（URL 与最近截图时间）
func (r *LinkRepo) UpdateScreenshot(ctx context.Context, id xSnowflake.SnowflakeID, url string, at time.Time, tx *gorm.DB) *xError.Error {
	r.log.Info(ctx, "UpdateScreenshot - 更新友链截图信息")

	result := r.pickDB(tx).WithContext(ctx).Model(&entity.LinkFriend{}).
		Where("id = ?", id).
		Updates(map[string]any{
			"screenshot_url": url,
			"screenshot_at":  at,
		})
	if result.Error != nil {
		return xError.NewError(ctx, xError.DatabaseError, "更新友链截图信息失败", false, result.Error)
	}
	if result.RowsAffected == 0 {
		return xError.NewError(ctx, xError.NotFound, "友情链接不存在", false)
	}

	if cacheErr := r.kc.Delete(ctx, constants.RedisLinkFriend.Get(id).String()); cacheErr != nil {
		r.log.Warn(ctx, cacheErr.Error())
	}
	return nil
}

// CountByStatus 按审核状态统计友情链接数量；status 为负数时统计全部
func (r *LinkRepo) CountByStatus(ctx context.Context, status int, tx *gorm.DB) (int64, *xError.Error) {
	query := r.pickDB(tx).WithContext(ctx).Model(&entity.LinkFriend{})
	if status >= 0 {
		query = query.Where("status = ?", status)
	}

	var count int64
	if err := query.Count(&count).Error; err != nil {
		return 0, xError.NewError(ctx, xError.DatabaseError, "统计友情链接数量失败", true, err)
	}
	return count, nil
}

// CountByGroupID 统计指定分组下友情链接数量
func (r *LinkRepo) CountByGroupID(ctx context.Context, groupID xSnowflake.SnowflakeID, tx *gorm.DB) (int64, *xError.Error) {
	var count int64
	err := r.pickDB(tx).WithContext(ctx).Model(&entity.LinkFriend{}).Where("group_id = ?", groupID).Count(&count).Error
	if err != nil {
		return 0, xError.NewError(ctx, xError.DatabaseError, "查询关联友链失败", false, err)
	}
	return count, nil
}

// CountByColorID 统计指定颜色下友情链接数量
func (r *LinkRepo) CountByColorID(ctx context.Context, colorID xSnowflake.SnowflakeID, tx *gorm.DB) (int64, *xError.Error) {
	var count int64
	err := r.pickDB(tx).WithContext(ctx).Model(&entity.LinkFriend{}).Where("color_id = ?", colorID).Count(&count).Error
	if err != nil {
		return 0, xError.NewError(ctx, xError.DatabaseError, "查询关联友链失败", false, err)
	}
	return count, nil
}

// ListByGroupID 根据分组ID查询友情链接列表
func (r *LinkRepo) ListByGroupID(ctx context.Context, groupID xSnowflake.SnowflakeID, limit int, tx *gorm.DB) ([]entity.LinkFriend, *xError.Error) {
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

// ListByColorID 根据颜色ID查询友情链接列表
func (r *LinkRepo) ListByColorID(ctx context.Context, colorID xSnowflake.SnowflakeID, limit int, tx *gorm.DB) ([]entity.LinkFriend, *xError.Error) {
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

// ClearGroupID 清空指定分组下友情链接的分组关联
func (r *LinkRepo) ClearGroupID(ctx context.Context, groupID xSnowflake.SnowflakeID, tx *gorm.DB) *xError.Error {
	result := r.pickDB(tx).WithContext(ctx).Model(&entity.LinkFriend{}).Where("group_id = ?", groupID).Update("group_id", nil)
	if result.Error != nil {
		return xError.NewError(ctx, xError.DatabaseError, "清空友链分组关联失败", false, result.Error)
	}
	return nil
}

// ClearColorID 清空指定颜色下友情链接的颜色关联
func (r *LinkRepo) ClearColorID(ctx context.Context, colorID xSnowflake.SnowflakeID, tx *gorm.DB) *xError.Error {
	result := r.pickDB(tx).WithContext(ctx).Model(&entity.LinkFriend{}).Where("color_id = ?", colorID).Update("color_id", nil)
	if result.Error != nil {
		return xError.NewError(ctx, xError.DatabaseError, "清空友链颜色关联失败", false, result.Error)
	}
	return nil
}
