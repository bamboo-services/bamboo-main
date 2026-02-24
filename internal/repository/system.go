package repository

import (
	"errors"
	"strings"

	xError "github.com/bamboo-services/bamboo-base-go/common/error"
	xLog "github.com/bamboo-services/bamboo-base-go/common/log"
	"github.com/bamboo-services/bamboo-main/internal/entity"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type SystemRepo struct {
	db  *gorm.DB
	log *xLog.LogNamedLogger
}

func NewSystemRepo(db *gorm.DB, _ *redis.Client) *SystemRepo {
	return &SystemRepo{
		db:  db,
		log: xLog.WithName(xLog.NamedREPO, "SystemRepo"),
	}
}

func (r *SystemRepo) ListByKeys(ctx *gin.Context, keys []string) ([]entity.System, *xError.Error) {
	r.log.Info(ctx, "ListByKeys - 查询系统配置")

	if len(keys) == 0 {
		return make([]entity.System, 0), nil
	}

	var configs []entity.System
	err := r.db.WithContext(ctx).Where("key IN ?", keys).Find(&configs).Error
	if err != nil {
		return nil, xError.NewError(ctx, xError.DatabaseError, "查询系统配置失败", true, err)
	}

	return configs, nil
}

func (r *SystemRepo) GetByKey(ctx *gin.Context, key string) (*entity.System, bool, *xError.Error) {
	r.log.Info(ctx, "GetByKey - 查询系统配置")

	key = strings.TrimSpace(key)
	if key == "" {
		return nil, false, nil
	}

	var config entity.System
	err := r.db.WithContext(ctx).Where("key = ?", key).First(&config).Error
	if err == nil {
		return &config, true, nil
	}
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, false, nil
	}

	return nil, false, xError.NewError(ctx, xError.DatabaseError, "查询系统配置失败", true, err)
}

func (r *SystemRepo) UpdateValueByKey(ctx *gin.Context, key string, value *string) *xError.Error {
	r.log.Info(ctx, "UpdateValueByKey - 更新系统配置")

	key = strings.TrimSpace(key)
	if key == "" {
		return nil
	}

	err := r.db.WithContext(ctx).Model(&entity.System{}).Where("key = ?", key).Update("value", value).Error
	if err != nil {
		return xError.NewError(ctx, xError.DatabaseError, "更新系统配置失败", true, err)
	}

	return nil
}
