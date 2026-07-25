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

	"github.com/bamboo-services/bamboo-main/internal/entity"
	"github.com/bamboo-services/bamboo-main/pkg/constants"

	xLog "github.com/bamboo-services/bamboo-base-go/common/log"
	xUtil "github.com/bamboo-services/bamboo-base-go/common/utility"
	"gorm.io/gorm"
)

// DefaultUser 初始化默认管理员用户的 xOptionDB.PrepareFunc。
//
// 幂等：基于 system.admin.id 标记判定是否已初始化，已存在则直接跳过。
// 失败时返回 error，由 xOption 启动流程中断并向上传递。
func DefaultUser(ctx context.Context, db *gorm.DB) error {
	log := xLog.WithName(xLog.NamedINIT)

	var adminSystem entity.System
	err := db.WithContext(ctx).Model(&entity.System{}).Where("key = ?", "system.admin.id").First(&adminSystem).Error
	if err == nil && adminSystem.Value != nil && *adminSystem.Value != "" {
		return nil
	}
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}

	passwordString, err := xUtil.Password().EncryptString("xiao_lfeng")
	if err != nil {
		return err
	}

	user := &entity.SystemUser{
		Username: "xiao_lfeng",
		Password: passwordString,
		Email:    "gm@x-lf.cn",
		Nickname: xUtil.Ptr("筱锋"),
		Avatar:   xUtil.Ptr("https://i-cdn.akass.cn/2024/05/664870a814c0d.png!wp60"),
		Role:     constants.RoleAdmin,
		Status:   constants.StatusActive,
	}
	if err = db.WithContext(ctx).Create(user).Error; err != nil {
		return err
	}

	adminID := user.ID.String()
	if err = db.WithContext(ctx).Create(&entity.System{
		Key:   "system.admin.id",
		Value: &adminID,
	}).Error; err != nil {
		return err
	}

	log.Info(ctx, "默认管理员用户初始化完成")
	return nil
}
