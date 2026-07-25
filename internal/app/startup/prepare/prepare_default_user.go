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
	"errors"

	"github.com/bamboo-services/bamboo-main/internal/entity"
	"github.com/bamboo-services/bamboo-main/pkg/constants"

	xUtil "github.com/bamboo-services/bamboo-base-go/common/utility"
	"gorm.io/gorm"
)

func (p *Prepare) prepareDefaultUser() error {
	var adminSystem entity.System
	err := p.db.WithContext(p.ctx).Model(&entity.System{}).Where("key = ?", "system.admin.id").First(&adminSystem).Error
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
	if err = p.db.WithContext(p.ctx).Create(user).Error; err != nil {
		return err
	}

	adminID := user.ID.String()
	return p.db.WithContext(p.ctx).Create(&entity.System{
		Key:   "system.admin.id",
		Value: &adminID,
	}).Error
}
