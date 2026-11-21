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

// SystemUserRepo 系统用户数据访问层
//
// 收口系统用户实体的查询、创建、保存与字段更新，缓存经 cache.Manager 统一读写与失效。
type SystemUserRepo struct {
	db  *gorm.DB
	kc  cache.KeyCache[string, entity.SystemUser]
	log *logging.LogNamedLogger
}

// NewSystemUserRepo 创建 SystemUserRepo 实例
func NewSystemUserRepo(db *gorm.DB, m *cache.Manager, logger ...*logging.Logger) *SystemUserRepo {
	return &SystemUserRepo{
		db:  db,
		kc:  cache.KeyCacheOf[string, entity.SystemUser](m),
		log: resolveRepoLogger("SystemUserRepo", logger...),
	}
}

// UserQuery 系统用户分页查询条件
type UserQuery struct {
	Page      int    // 页码，从 1 开始
	PageSize  int    // 每页数量
	Keyword   string // 搜索关键词（用户名/邮箱/昵称模糊匹配）
	Status    *int   // 状态筛选（0: 禁用, 1: 启用）；指针类型以区分「未筛选」与「筛选禁用」
	SortBy    string // 排序字段（created_at/updated_at/username/email/last_login_at）
	SortOrder string // 排序方式（asc/desc）
}

// GetByID 根据ID获取用户
func (r *SystemUserRepo) GetByID(ctx context.Context, id snowflake.SnowflakeID) (*entity.SystemUser, bool, *apperr.Error) {
	r.log.Info(ctx, "GetByID - 获取用户信息")

	if user, ok, err := r.kc.Get(ctx, constants.RedisSystemUser.Get(id).String()); err != nil {
		return nil, false, apperr.New(cache.CacheError, "获取用户缓存失败", err)
	} else if ok {
		return user, true, nil
	}

	var user entity.SystemUser
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&user).Error
	if err == nil {
		if cacheErr := r.kc.Set(ctx, constants.RedisSystemUser.Get(user.ID).String(), &user, cache.WithTTL(15*time.Minute)); cacheErr != nil {
			r.log.Warn(ctx, cacheErr.Error())
		}
		return &user, true, nil
	}
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, false, nil
	}
	return nil, false, apperr.New(database.DatabaseError, "查询用户失败", err)
}

// GetPasswordByID 直接查询用户密码哈希（绕过缓存）
//
// 用户实体经缓存序列化（encoding/json）时，Password 字段因 json:"-" 被丢弃，
// 命中缓存读回的实体密码字段恒为空。修改密码等需要校验旧密码哈希的场景
// 必须直查数据库，避免以空哈希验证导致校验恒失败。
func (r *SystemUserRepo) GetPasswordByID(ctx context.Context, userID snowflake.SnowflakeID) (string, bool, *apperr.Error) {
	r.log.Info(ctx, "GetPasswordByID - 查询用户密码")

	var user entity.SystemUser
	err := r.db.WithContext(ctx).Where("id = ?", userID).Select("id", "password").First(&user).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return "", false, nil
	}
	if err != nil {
		return "", false, apperr.New(database.DatabaseError, "查询用户密码失败", err)
	}
	return user.Password, true, nil
}

// GetByUsernameOrEmail 根据用户名或邮箱获取用户
func (r *SystemUserRepo) GetByUsernameOrEmail(ctx context.Context, keyword string) (*entity.SystemUser, bool, *apperr.Error) {
	r.log.Info(ctx, "GetByUsernameOrEmail - 查询用户")

	var user entity.SystemUser
	err := r.db.WithContext(ctx).Where("username = ? OR email = ?", keyword, keyword).First(&user).Error
	if err == nil {
		if cacheErr := r.kc.Set(ctx, constants.RedisSystemUser.Get(user.ID).String(), &user, cache.WithTTL(15*time.Minute)); cacheErr != nil {
			r.log.Warn(ctx, cacheErr.Error())
		}
		return &user, true, nil
	}
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, false, nil
	}
	return nil, false, apperr.New(database.DatabaseError, "查询用户失败", err)
}

// GetByEmail 根据邮箱获取用户
func (r *SystemUserRepo) GetByEmail(ctx context.Context, email string) (*entity.SystemUser, bool, *apperr.Error) {
	r.log.Info(ctx, "GetByEmail - 查询用户")

	email = strings.TrimSpace(email)
	if email == "" {
		return nil, false, nil
	}

	var user entity.SystemUser
	err := r.db.WithContext(ctx).Where("email = ?", email).First(&user).Error
	if err == nil {
		if cacheErr := r.kc.Set(ctx, constants.RedisSystemUser.Get(user.ID).String(), &user, cache.WithTTL(15*time.Minute)); cacheErr != nil {
			r.log.Warn(ctx, cacheErr.Error())
		}
		return &user, true, nil
	}
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, false, nil
	}
	return nil, false, apperr.New(database.DatabaseError, "查询用户失败", err)
}

// GetByOAuthUserID 根据OAuth用户标识获取用户
func (r *SystemUserRepo) GetByOAuthUserID(ctx context.Context, oauthUserID string) (*entity.SystemUser, bool, *apperr.Error) {
	r.log.Info(ctx, "GetByOAuthUserID - 查询用户")

	oauthUserID = strings.TrimSpace(oauthUserID)
	if oauthUserID == "" {
		return nil, false, nil
	}

	var user entity.SystemUser
	err := r.db.WithContext(ctx).Where("oauth_user_id = ?", oauthUserID).First(&user).Error
	if err == nil {
		if cacheErr := r.kc.Set(ctx, constants.RedisSystemUser.Get(user.ID).String(), &user, cache.WithTTL(15*time.Minute)); cacheErr != nil {
			r.log.Warn(ctx, cacheErr.Error())
		}
		return &user, true, nil
	}
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, false, nil
	}
	return nil, false, apperr.New(database.DatabaseError, "查询用户失败", err)
}

// ExistsByUsername 判断用户名是否已存在
func (r *SystemUserRepo) ExistsByUsername(ctx context.Context, username string) (bool, *apperr.Error) {
	r.log.Info(ctx, "ExistsByUsername - 检查用户名")

	var count int64
	err := r.db.WithContext(ctx).Model(&entity.SystemUser{}).Where("username = ?", username).Count(&count).Error
	if err != nil {
		return false, apperr.New(database.DatabaseError, "检查用户名失败", err)
	}
	return count > 0, nil
}

// ExistsByUsernameExceptID 判断用户名是否已被其他用户占用
func (r *SystemUserRepo) ExistsByUsernameExceptID(ctx context.Context, username string, exceptID snowflake.SnowflakeID) (bool, *apperr.Error) {
	r.log.Info(ctx, "ExistsByUsernameExceptID - 检查用户名")

	var count int64
	query := r.db.WithContext(ctx).Model(&entity.SystemUser{}).Where("username = ?", username)
	if exceptID > 0 {
		query = query.Where("id <> ?", exceptID)
	}
	err := query.Count(&count).Error
	if err != nil {
		return false, apperr.New(database.DatabaseError, "检查用户名失败", err)
	}
	return count > 0, nil
}

// ExistsByEmailExceptID 判断邮箱是否已被其他用户占用
func (r *SystemUserRepo) ExistsByEmailExceptID(ctx context.Context, email string, exceptID snowflake.SnowflakeID) (bool, *apperr.Error) {
	r.log.Info(ctx, "ExistsByEmailExceptID - 检查邮箱")

	var count int64
	query := r.db.WithContext(ctx).Model(&entity.SystemUser{}).Where("email = ?", email)
	if exceptID > 0 {
		query = query.Where("id <> ?", exceptID)
	}
	err := query.Count(&count).Error
	if err != nil {
		return false, apperr.New(database.DatabaseError, "检查邮箱失败", err)
	}
	return count > 0, nil
}

// Create 创建用户
func (r *SystemUserRepo) Create(ctx context.Context, user *entity.SystemUser) (*entity.SystemUser, *apperr.Error) {
	r.log.Info(ctx, "Create - 创建用户")

	err := r.db.WithContext(ctx).Create(user).Error
	if err != nil {
		return nil, apperr.New(database.DatabaseError, "创建用户失败", err)
	}

	if cacheErr := r.kc.Set(ctx, constants.RedisSystemUser.Get(user.ID).String(), user, cache.WithTTL(15*time.Minute)); cacheErr != nil {
		r.log.Warn(ctx, cacheErr.Error())
	}

	return user, nil
}

// Save 保存用户（存在则更新）
func (r *SystemUserRepo) Save(ctx context.Context, user *entity.SystemUser) (*entity.SystemUser, *apperr.Error) {
	r.log.Info(ctx, "Save - 保存用户")

	err := r.db.WithContext(ctx).Save(user).Error
	if err != nil {
		return nil, apperr.New(database.DatabaseError, "保存用户失败", err)
	}

	if cacheErr := r.kc.Set(ctx, constants.RedisSystemUser.Get(user.ID).String(), user, cache.WithTTL(15*time.Minute)); cacheErr != nil {
		r.log.Warn(ctx, cacheErr.Error())
	}

	return user, nil
}

// UpdateFieldsByID 根据ID更新用户的指定字段
func (r *SystemUserRepo) UpdateFieldsByID(ctx context.Context, userID snowflake.SnowflakeID, updates map[string]any) (*entity.SystemUser, *apperr.Error) {
	r.log.Info(ctx, "UpdateFieldsByID - 更新用户字段")

	err := r.db.WithContext(ctx).Model(&entity.SystemUser{}).Where("id = ?", userID).Updates(updates).Error
	if err != nil {
		return nil, apperr.New(database.DatabaseError, "更新用户失败", err)
	}

	if cacheErr := r.kc.Delete(ctx, constants.RedisSystemUser.Get(userID).String()); cacheErr != nil {
		r.log.Warn(ctx, cacheErr.Error())
	}

	user, found, xErr := r.GetByID(ctx, userID)
	if xErr != nil {
		return nil, xErr
	}
	if !found {
		return nil, apperr.New(httpresult.NotFound, "用户不存在")
	}

	return user, nil
}

// UpdatePasswordByID 根据ID更新用户密码
func (r *SystemUserRepo) UpdatePasswordByID(ctx context.Context, userID snowflake.SnowflakeID, hashedPassword string) *apperr.Error {
	_, xErr := r.UpdateFieldsByID(ctx, userID, map[string]any{"password": hashedPassword})
	return xErr
}

// UpdateLastLoginByID 根据ID更新用户最近登录时间
func (r *SystemUserRepo) UpdateLastLoginByID(ctx context.Context, userID snowflake.SnowflakeID, loginAt *time.Time) *apperr.Error {
	updates := map[string]any{"last_login_at": loginAt}
	_, xErr := r.UpdateFieldsByID(ctx, userID, updates)
	return xErr
}

// List 分页查询系统用户
//
// 管理端用户列表查询：支持用户名/邮箱/昵称模糊匹配、状态筛选与白名单排序。
// 分页列表不过缓存（cache-aside 仅针对单实体），保证列表结果实时新鲜。
func (r *SystemUserRepo) List(ctx context.Context, req *UserQuery) ([]entity.SystemUser, int64, *apperr.Error) {
	r.log.Info(ctx, "List - 查询系统用户分页列表")

	query := r.db.WithContext(ctx).Model(&entity.SystemUser{})

	if req.Keyword != "" {
		kw := "%" + req.Keyword + "%"
		query = query.Where("(username ILIKE ? OR email ILIKE ? OR nickname ILIKE ?)", kw, kw, kw)
	}
	if req.Status != nil {
		query = query.Where("status = ?", *req.Status)
	}

	// 白名单排序，杜绝任意列名注入
	orderBy := "created_at"
	switch req.SortBy {
	case "updated_at":
		orderBy = "updated_at"
	case "username":
		orderBy = "username"
	case "email":
		orderBy = "email"
	case "last_login_at":
		orderBy = "last_login_at"
	}

	sortOrder := "DESC"
	if strings.EqualFold(req.SortOrder, "asc") {
		sortOrder = "ASC"
	}

	query = query.Order(orderBy + " " + sortOrder)

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, apperr.New(database.DatabaseError, "统计系统用户数量失败", err)
	}

	offset := (req.Page - 1) * req.PageSize
	var users []entity.SystemUser
	if err := query.Offset(offset).Limit(req.PageSize).Find(&users).Error; err != nil {
		return nil, 0, apperr.New(database.DatabaseError, "查询系统用户列表失败", err)
	}

	return users, total, nil
}
