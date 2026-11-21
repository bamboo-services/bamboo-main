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

package prepare

import (
	"context"
	"errors"

	"github.com/bamboo-services/bamboo-base-go/foundation/v2/logging"
	"github.com/bamboo-services/bamboo-base-go/foundation/v2/toolkit"
	"github.com/bamboo-services/bamboo-main/internal/entity"
	"github.com/bamboo-services/bamboo-main/pkg/constants"
	"gorm.io/gorm"
)

// DefaultUser 初始化默认管理员用户（基于 system.admin.id 幂等）。
func DefaultUser(ctx context.Context, db *gorm.DB, logger *logging.Logger) error {
	var adminSystem entity.System
	err := db.WithContext(ctx).Model(&entity.System{}).Where("key = ?", constants.KeySystemAdminID).First(&adminSystem).Error
	if err == nil && adminSystem.Value != nil && *adminSystem.Value != "" {
		return nil
	}
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}

	passwordString, appErr := toolkit.EncryptPasswordString("123456")
	if appErr != nil {
		return appErr
	}

	user := &entity.SystemUser{
		Username:    "super_admin",
		Password:    passwordString,
		Email:       "gm@x-lf.cn",
		Nickname:    toolkit.Ptr("筱锋"),
		Avatar:      toolkit.Ptr("https://i-cdn.akass.cn/2024/05/664870a814c0d.png!wp60"),
		Status:      constants.StatusActive,
		EmailVerify: true,
	}
	if err = db.WithContext(ctx).Create(user).Error; err != nil {
		return err
	}

	adminID := user.ID.String()
	if err = db.WithContext(ctx).Create(&entity.System{
		Key:   constants.KeySystemAdminID,
		Value: &adminID,
	}).Error; err != nil {
		return err
	}

	if logger != nil {
		logger.WithName(logging.NamedINIT).Info(ctx, "默认管理员用户初始化完成")
	}
	return nil
}
