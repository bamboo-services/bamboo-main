package prepare

import (
	"errors"
	"strconv"

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

	adminID := strconv.FormatInt(user.ID, 10)
	return p.db.WithContext(p.ctx).Create(&entity.System{
		Key:   "system.admin.id",
		Value: &adminID,
	}).Error
}
