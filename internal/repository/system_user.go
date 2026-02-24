package repository

import (
	"errors"
	"strings"
	"time"

	xError "github.com/bamboo-services/bamboo-base-go/common/error"
	xLog "github.com/bamboo-services/bamboo-base-go/common/log"
	"github.com/bamboo-services/bamboo-main/internal/entity"
	"github.com/bamboo-services/bamboo-main/internal/repository/cache"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type SystemUserRepo struct {
	db    *gorm.DB
	cache *cache.SystemUserCache
	log   *xLog.LogNamedLogger
}

func NewSystemUserRepo(db *gorm.DB, rdb *redis.Client) *SystemUserRepo {
	return &SystemUserRepo{
		db: db,
		cache: &cache.SystemUserCache{
			RDB: rdb,
			TTL: time.Minute * 15,
		},
		log: xLog.WithName(xLog.NamedREPO, "SystemUserRepo"),
	}
}

func (r *SystemUserRepo) GetByID(ctx *gin.Context, id int64) (*entity.SystemUser, bool, *xError.Error) {
	r.log.Info(ctx, "GetByID - 获取用户信息")

	if user, err := r.cache.Get(ctx.Request.Context(), id); err != nil {
		return nil, false, xError.NewError(ctx, xError.CacheError, "获取用户缓存失败", true, err)
	} else if user != nil {
		return user, true, nil
	}

	var user entity.SystemUser
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&user).Error
	if err == nil {
		if cacheErr := r.cache.Set(ctx.Request.Context(), &user); cacheErr != nil {
			r.log.Warn(ctx, cacheErr.Error())
		}
		return &user, true, nil
	}
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, false, nil
	}
	return nil, false, xError.NewError(ctx, xError.DatabaseError, "查询用户失败", true, err)
}

func (r *SystemUserRepo) GetByUsernameOrEmail(ctx *gin.Context, keyword string) (*entity.SystemUser, bool, *xError.Error) {
	r.log.Info(ctx, "GetByUsernameOrEmail - 查询用户")

	var user entity.SystemUser
	err := r.db.WithContext(ctx).Where("username = ? OR email = ?", keyword, keyword).First(&user).Error
	if err == nil {
		if cacheErr := r.cache.Set(ctx.Request.Context(), &user); cacheErr != nil {
			r.log.Warn(ctx, cacheErr.Error())
		}
		return &user, true, nil
	}
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, false, nil
	}
	return nil, false, xError.NewError(ctx, xError.DatabaseError, "查询用户失败", true, err)
}

func (r *SystemUserRepo) GetByEmail(ctx *gin.Context, email string) (*entity.SystemUser, bool, *xError.Error) {
	r.log.Info(ctx, "GetByEmail - 查询用户")

	email = strings.TrimSpace(email)
	if email == "" {
		return nil, false, nil
	}

	var user entity.SystemUser
	err := r.db.WithContext(ctx).Where("email = ?", email).First(&user).Error
	if err == nil {
		if cacheErr := r.cache.Set(ctx.Request.Context(), &user); cacheErr != nil {
			r.log.Warn(ctx, cacheErr.Error())
		}
		return &user, true, nil
	}
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, false, nil
	}
	return nil, false, xError.NewError(ctx, xError.DatabaseError, "查询用户失败", true, err)
}

func (r *SystemUserRepo) GetByOAuthUserID(ctx *gin.Context, oauthUserID string) (*entity.SystemUser, bool, *xError.Error) {
	r.log.Info(ctx, "GetByOAuthUserID - 查询用户")

	oauthUserID = strings.TrimSpace(oauthUserID)
	if oauthUserID == "" {
		return nil, false, nil
	}

	var user entity.SystemUser
	err := r.db.WithContext(ctx).Where("oauth_user_id = ?", oauthUserID).First(&user).Error
	if err == nil {
		if cacheErr := r.cache.Set(ctx.Request.Context(), &user); cacheErr != nil {
			r.log.Warn(ctx, cacheErr.Error())
		}
		return &user, true, nil
	}
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, false, nil
	}
	return nil, false, xError.NewError(ctx, xError.DatabaseError, "查询用户失败", true, err)
}

func (r *SystemUserRepo) ExistsByUsername(ctx *gin.Context, username string) (bool, *xError.Error) {
	r.log.Info(ctx, "ExistsByUsername - 检查用户名")

	var count int64
	err := r.db.WithContext(ctx).Model(&entity.SystemUser{}).Where("username = ?", username).Count(&count).Error
	if err != nil {
		return false, xError.NewError(ctx, xError.DatabaseError, "检查用户名失败", true, err)
	}
	return count > 0, nil
}

func (r *SystemUserRepo) ExistsByEmailExceptID(ctx *gin.Context, email string, exceptID int64) (bool, *xError.Error) {
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

func (r *SystemUserRepo) Create(ctx *gin.Context, user *entity.SystemUser) (*entity.SystemUser, *xError.Error) {
	r.log.Info(ctx, "Create - 创建用户")

	err := r.db.WithContext(ctx).Create(user).Error
	if err != nil {
		return nil, xError.NewError(ctx, xError.DatabaseError, "创建用户失败", true, err)
	}

	if cacheErr := r.cache.Set(ctx.Request.Context(), user); cacheErr != nil {
		r.log.Warn(ctx, cacheErr.Error())
	}

	return user, nil
}

func (r *SystemUserRepo) Save(ctx *gin.Context, user *entity.SystemUser) (*entity.SystemUser, *xError.Error) {
	r.log.Info(ctx, "Save - 保存用户")

	err := r.db.WithContext(ctx).Save(user).Error
	if err != nil {
		return nil, xError.NewError(ctx, xError.DatabaseError, "保存用户失败", true, err)
	}

	if cacheErr := r.cache.Set(ctx.Request.Context(), user); cacheErr != nil {
		r.log.Warn(ctx, cacheErr.Error())
	}

	return user, nil
}

func (r *SystemUserRepo) UpdateFieldsByID(ctx *gin.Context, userID int64, updates map[string]any) (*entity.SystemUser, *xError.Error) {
	r.log.Info(ctx, "UpdateFieldsByID - 更新用户字段")

	err := r.db.WithContext(ctx).Model(&entity.SystemUser{}).Where("id = ?", userID).Updates(updates).Error
	if err != nil {
		return nil, xError.NewError(ctx, xError.DatabaseError, "更新用户失败", true, err)
	}

	if cacheErr := r.cache.Delete(ctx.Request.Context(), userID); cacheErr != nil {
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

func (r *SystemUserRepo) UpdatePasswordByID(ctx *gin.Context, userID int64, hashedPassword string) *xError.Error {
	_, xErr := r.UpdateFieldsByID(ctx, userID, map[string]any{"password": hashedPassword})
	return xErr
}

func (r *SystemUserRepo) UpdateLastLoginByID(ctx *gin.Context, userID int64, loginAt *time.Time) *xError.Error {
	updates := map[string]any{"last_login_at": loginAt}
	_, xErr := r.UpdateFieldsByID(ctx, userID, updates)
	return xErr
}
