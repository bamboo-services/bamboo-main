/*
 * --------------------------------------------------------------------------------
 * Copyright (c) 2016-NOW(至今) 筱锋
 * Author: 筱锋「xiao_lfeng」(https://www.x-lf.com)
 * --------------------------------------------------------------------------------
 * 许可证声明：版权所有 (c) 2016-2025 筱锋。保留所有权利。
 * 有关MIT许可证的更多信息，请查看项目根目录下的LICENSE文件或访问：
 * https://opensource.org/licenses/MIT
 * --------------------------------------------------------------------------------
 */

package startup

import (
	"bamboo-main/internal/model/entity"
	"errors"
	"strconv"

	xConsts "github.com/bamboo-services/bamboo-base-go/constants"
	xUtil "github.com/bamboo-services/bamboo-base-go/utility"
	"gorm.io/gorm"
)

func (r *Reg) DatabaseUserInit() {
	log := r.Serv.Logger.Named(xConsts.LogINIT)
	log.Info("初始化默认用户")

	var hasUser = true

	var getValue *entity.System
	result := r.DB.Model(&entity.System{}).Where(
		"key = ?",
		"system.admin.id",
	).First(&getValue)
	if result.Error != nil {
		if !errors.Is(result.Error, gorm.ErrRecordNotFound) {
			panic(result.Error)
		} else {
			hasUser = false
		}
	}

	if getValue.Value == nil {
		hasUser = false
	}

	// 如果没有用户，则创建默认用户
	if !hasUser {
		passwordString, _ := xUtil.EncryptPasswordString("xiao_lfeng")
		var user = &entity.SystemUser{
			ID:       r.SnowflakeNode.Generate().Int64(), // 手动生成 Snowflake ID
			Username: "xiao_lfeng",
			Password: passwordString,
			Email:    "gm@x-lf.cn",
			Nickname: xUtil.Ptr("筱锋"),
			Avatar:   xUtil.Ptr("https://i-cdn.akass.cn/2024/05/664870a814c0d.png!wp60"),
			Role:     "admin",
			Status:   1,
		}
		err := r.DB.Create(user).Error
		if err != nil {
			panic(err)
		}

		err = r.DB.Create(&entity.System{
			ID:    r.SnowflakeNode.Generate().Int64(), // 手动生成 Snowflake ID
			Key:   "system.admin.id",
			Value: xUtil.Ptr(strconv.FormatInt(user.ID, 10)),
		}).Error
		if err != nil {
			panic(err)
		}
	}
}
