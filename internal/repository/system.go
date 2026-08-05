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

	xError "github.com/bamboo-services/bamboo-base-go/common/error"
	xLog "github.com/bamboo-services/bamboo-base-go/common/log"
	xCache "github.com/bamboo-services/bamboo-base-go/major/cache"
	"github.com/bamboo-services/bamboo-main/internal/entity"
	bConst "github.com/bamboo-services/bamboo-main/pkg/constants"
	"gorm.io/gorm"
)

// SystemRepo 系统配置数据访问层
//
// 收口系统配置实体的按键查询与按键更新，供 logic 层读取和修改全局配置项。
type SystemRepo struct {
	db  *gorm.DB
	log *xLog.LogNamedLogger
}

// NewSystemRepo 创建 SystemRepo 实例
func NewSystemRepo(db *gorm.DB, _ *xCache.Manager) *SystemRepo {
	return &SystemRepo{
		db:  db,
		log: xLog.WithName(xLog.NamedREPO, "SystemRepo"),
	}
}

// ListByKeys 根据键集合批量查询系统配置
func (r *SystemRepo) ListByKeys(ctx context.Context, keys []string) ([]entity.System, *xError.Error) {
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

// GetByKey 根据键获取单条系统配置
func (r *SystemRepo) GetByKey(ctx context.Context, key string) (*entity.System, bool, *xError.Error) {
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

// UpdateValueByKey 根据键更新系统配置的值
func (r *SystemRepo) UpdateValueByKey(ctx context.Context, key string, value *string) *xError.Error {
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

// GetValuesByKeys 根据键集合批量查询系统配置，返回 key → 值映射。
//
// 未命中或值为空串的键不会出现在结果中，便于调用方按缺省值兜底。
func (r *SystemRepo) GetValuesByKeys(ctx context.Context, keys []string) (map[string]string, *xError.Error) {
	configs, xErr := r.ListByKeys(ctx, keys)
	if xErr != nil {
		return nil, xErr
	}

	values := make(map[string]string, len(configs))
	for i := range configs {
		if configs[i].Value != nil && *configs[i].Value != "" {
			values[configs[i].Key] = *configs[i].Value
		}
	}

	return values, nil
}

// BuildBuiltinInvalidGroup 读取内置「已失效」分组配置并构造虚拟分组对象。
//
// 名称缺省回退「已失效」，描述缺省为 nil；DB 读取失败返回错误，由调用方决定降级策略
// （公开接口建议以 entity.NewDefaultBuiltinGroup 兜底，不阻断查询）。
func (r *SystemRepo) BuildBuiltinInvalidGroup(ctx context.Context) (*entity.LinkGroup, *xError.Error) {
	values, xErr := r.GetValuesByKeys(ctx, []string{
		bConst.KeyBuiltinInvalidGroupName,
		bConst.KeyBuiltinInvalidGroupDesc,
	})
	if xErr != nil {
		return nil, xErr
	}

	name := entity.NewDefaultBuiltinGroup().Name
	if v, ok := values[bConst.KeyBuiltinInvalidGroupName]; ok {
		name = v
	}

	var desc *string
	if v, ok := values[bConst.KeyBuiltinInvalidGroupDesc]; ok {
		desc = &v
	}

	return entity.NewBuiltinGroup(name, desc), nil
}
