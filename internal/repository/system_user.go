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

// SystemUserRepo 系统用户数据访问层
//
// 收口系统用户实体的查询、创建、保存与字段更新，缓存经 xCache.Manager 统一读写与失效。
type SystemUserRepo struct {
	db  *gorm.DB
	kc  xCache.KeyCache[string, entity.SystemUser]
	log *xLog.LogNamedLogger
}

// NewSystemUserRepo 创建 SystemUserRepo 实例
func NewSystemUserRepo(db *gorm.DB, m *xCache.Manager) *SystemUserRepo {
	return &SystemUserRepo{
		db:  db,
		kc:  xCache.KeyCacheOf[string, entity.SystemUser](m),
		log: xLog.WithName(xLog.NamedREPO, "SystemUserRepo"),
	}
}

// GetByID 根据ID获取用户
func (r *SystemUserRepo) GetByID(ctx context.Context, id xSnowflake.SnowflakeID) (*entity.SystemUser, bool, *xError.Error) {
	r.log.Info(ctx, "GetByID - 获取用户信息")

	if user, ok, err := r.kc.Get(ctx, constants.RedisSystemUser.Get(id).String()); err != nil {
		return nil, false, xError.NewError(ctx, xError.CacheError, "获取用户缓存失败", true, err)
	} else if ok {
		return user, true, nil
	}

	var user entity.SystemUser
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&user).Error
	if err == nil {
		if cacheErr := r.kc.Set(ctx, constants.RedisSystemUser.Get(user.ID).String(), &user, xCache.WithTTL(15*time.Minute)); cacheErr != nil {
			r.log.Warn(ctx, cacheErr.Error())
		}
		return &user, true, nil
	}
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, false, nil
	}
	return nil, false, xError.NewError(ctx, xError.DatabaseError, "查询用户失败", true, err)
}

// GetPasswordByID 直接查询用户密码哈希（绕过缓存）
//
// 用户实体经缓存序列化（encoding/json）时，Password 字段因 json:"-" 被丢弃，
// 命中缓存读回的实体密码字段恒为空。修改密码等需要校验旧密码哈希的场景
// 必须直查数据库，避免以空哈希验证导致校验恒失败。
func (r *SystemUserRepo) GetPasswordByID(ctx context.Context, userID xSnowflake.SnowflakeID) (string, bool, *xError.Error) {
	r.log.Info(ctx, "GetPasswordByID - 查询用户密码")

	var user entity.SystemUser
	err := r.db.WithContext(ctx).Where("id = ?", userID).Select("id", "password").First(&user).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return "", false, nil
	}
	if err != nil {
		return "", false, xError.NewError(ctx, xError.DatabaseError, "查询用户密码失败", true, err)
	}
	return user.Password, true, nil
}

// GetByUsernameOrEmail 根据用户名或邮箱获取用户
func (r *SystemUserRepo) GetByUsernameOrEmail(ctx context.Context, keyword string) (*entity.SystemUser, bool, *xError.Error) {
	r.log.Info(ctx, "GetByUsernameOrEmail - 查询用户")

	var user entity.SystemUser
	err := r.db.WithContext(ctx).Where("username = ? OR email = ?", keyword, keyword).First(&user).Error
	if err == nil {
		if cacheErr := r.kc.Set(ctx, constants.RedisSystemUser.Get(user.ID).String(), &user, xCache.WithTTL(15*time.Minute)); cacheErr != nil {
			r.log.Warn(ctx, cacheErr.Error())
		}
		return &user, true, nil
	}
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, false, nil
	}
	return nil, false, xError.NewError(ctx, xError.DatabaseError, "查询用户失败", true, err)
}

// GetByEmail 根据邮箱获取用户
func (r *SystemUserRepo) GetByEmail(ctx context.Context, email string) (*entity.SystemUser, bool, *xError.Error) {
	r.log.Info(ctx, "GetByEmail - 查询用户")

	email = strings.TrimSpace(email)
	if email == "" {
		return nil, false, nil
	}

	var user entity.SystemUser
	err := r.db.WithContext(ctx).Where("email = ?", email).First(&user).Error
	if err == nil {
		if cacheErr := r.kc.Set(ctx, constants.RedisSystemUser.Get(user.ID).String(), &user, xCache.WithTTL(15*time.Minute)); cacheErr != nil {
			r.log.Warn(ctx, cacheErr.Error())
		}
		return &user, true, nil
	}
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, false, nil
	}
	return nil, false, xError.NewError(ctx, xError.DatabaseError, "查询用户失败", true, err)
}

// GetByOAuthUserID 根据OAuth用户标识获取用户
func (r *SystemUserRepo) GetByOAuthUserID(ctx context.Context, oauthUserID string) (*entity.SystemUser, bool, *xError.Error) {
	r.log.Info(ctx, "GetByOAuthUserID - 查询用户")

	oauthUserID = strings.TrimSpace(oauthUserID)
	if oauthUserID == "" {
		return nil, false, nil
	}

	var user entity.SystemUser
	err := r.db.WithContext(ctx).Where("oauth_user_id = ?", oauthUserID).First(&user).Error
	if err == nil {
		if cacheErr := r.kc.Set(ctx, constants.RedisSystemUser.Get(user.ID).String(), &user, xCache.WithTTL(15*time.Minute)); cacheErr != nil {
			r.log.Warn(ctx, cacheErr.Error())
		}
		return &user, true, nil
	}
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, false, nil
	}
	return nil, false, xError.NewError(ctx, xError.DatabaseError, "查询用户失败", true, err)
}

// ExistsByUsername 判断用户名是否已存在
func (r *SystemUserRepo) ExistsByUsername(ctx context.Context, username string) (bool, *xError.Error) {
	r.log.Info(ctx, "ExistsByUsername - 检查用户名")

	var count int64
	err := r.db.WithContext(ctx).Model(&entity.SystemUser{}).Where("username = ?", username).Count(&count).Error
	if err != nil {
		return false, xError.NewError(ctx, xError.DatabaseError, "检查用户名失败", true, err)
	}
	return count > 0, nil
}

// ExistsByUsernameExceptID 判断用户名是否已被其他用户占用
func (r *SystemUserRepo) ExistsByUsernameExceptID(ctx context.Context, username string, exceptID xSnowflake.SnowflakeID) (bool, *xError.Error) {
	r.log.Info(ctx, "ExistsByUsernameExceptID - 检查用户名")

	var count int64
	query := r.db.WithContext(ctx).Model(&entity.SystemUser{}).Where("username = ?", username)
	if exceptID > 0 {
		query = query.Where("id <> ?", exceptID)
	}
	err := query.Count(&count).Error
	if err != nil {
		return false, xError.NewError(ctx, xError.DatabaseError, "检查用户名失败", true, err)
	}
	return count > 0, nil
}

// ExistsByEmailExceptID 判断邮箱是否已被其他用户占用
func (r *SystemUserRepo) ExistsByEmailExceptID(ctx context.Context, email string, exceptID xSnowflake.SnowflakeID) (bool, *xError.Error) {
	r.log.Info(ctx, "ExistsByEmailExceptID - 检查邮箱")

	var count int64
	query := r.db.WithContext(ctx).Model(&entity.SystemUser{}).Where("email = ?", email)
	if exceptID > 0 {
		query = query.Where("id <> ?", exceptID)
	}
	err := query.Count(&count).Error
	if err != nil {
		return false, xError.NewError(ctx, xError.DatabaseError, "检查邮箱失败", true, err)
	}
	return count > 0, nil
}

// Create 创建用户
func (r *SystemUserRepo) Create(ctx context.Context, user *entity.SystemUser) (*entity.SystemUser, *xError.Error) {
	r.log.Info(ctx, "Create - 创建用户")

	err := r.db.WithContext(ctx).Create(user).Error
	if err != nil {
		return nil, xError.NewError(ctx, xError.DatabaseError, "创建用户失败", true, err)
	}

	if cacheErr := r.kc.Set(ctx, constants.RedisSystemUser.Get(user.ID).String(), user, xCache.WithTTL(15*time.Minute)); cacheErr != nil {
		r.log.Warn(ctx, cacheErr.Error())
	}

	return user, nil
}

// Save 保存用户（存在则更新）
func (r *SystemUserRepo) Save(ctx context.Context, user *entity.SystemUser) (*entity.SystemUser, *xError.Error) {
	r.log.Info(ctx, "Save - 保存用户")

	err := r.db.WithContext(ctx).Save(user).Error
	if err != nil {
		return nil, xError.NewError(ctx, xError.DatabaseError, "保存用户失败", true, err)
	}

	if cacheErr := r.kc.Set(ctx, constants.RedisSystemUser.Get(user.ID).String(), user, xCache.WithTTL(15*time.Minute)); cacheErr != nil {
		r.log.Warn(ctx, cacheErr.Error())
	}

	return user, nil
}

// UpdateFieldsByID 根据ID更新用户的指定字段
func (r *SystemUserRepo) UpdateFieldsByID(ctx context.Context, userID xSnowflake.SnowflakeID, updates map[string]any) (*entity.SystemUser, *xError.Error) {
	r.log.Info(ctx, "UpdateFieldsByID - 更新用户字段")

	err := r.db.WithContext(ctx).Model(&entity.SystemUser{}).Where("id = ?", userID).Updates(updates).Error
	if err != nil {
		return nil, xError.NewError(ctx, xError.DatabaseError, "更新用户失败", true, err)
	}

	if cacheErr := r.kc.Delete(ctx, constants.RedisSystemUser.Get(userID).String()); cacheErr != nil {
		r.log.Warn(ctx, cacheErr.Error())
	}

	user, found, xErr := r.GetByID(ctx, userID)
	if xErr != nil {
		return nil, xErr
	}
	if !found {
		return nil, xError.NewError(ctx, xError.NotFound, "用户不存在", false)
	}

	return user, nil
}

// UpdatePasswordByID 根据ID更新用户密码
func (r *SystemUserRepo) UpdatePasswordByID(ctx context.Context, userID xSnowflake.SnowflakeID, hashedPassword string) *xError.Error {
	_, xErr := r.UpdateFieldsByID(ctx, userID, map[string]any{"password": hashedPassword})
	return xErr
}

// UpdateLastLoginByID 根据ID更新用户最近登录时间
func (r *SystemUserRepo) UpdateLastLoginByID(ctx context.Context, userID xSnowflake.SnowflakeID, loginAt *time.Time) *xError.Error {
	updates := map[string]any{"last_login_at": loginAt}
	_, xErr := r.UpdateFieldsByID(ctx, userID, updates)
	return xErr
}
