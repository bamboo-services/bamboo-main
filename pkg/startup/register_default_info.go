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

	xConsts "github.com/bamboo-services/bamboo-base-go/constants"
	xUtil "github.com/bamboo-services/bamboo-base-go/utility"
	"gorm.io/gorm"
)

// defaultInfoConfigs 站点信息默认配置
var defaultInfoConfigs = []struct {
	Key   string
	Value string
}{
	{Key: "site.name", Value: "Bamboo Links"},
	{Key: "site.description", Value: "一个简洁优雅的友情链接管理系统"},
	{Key: "site.introduction", Value: "欢迎来到我的友链主页！这里收录了我的好朋友们的博客链接。"},
	{Key: "profile.about", Value: "# 关于我\n\n这里是自我介绍，支持 **Markdown** 格式。"},
}

// DatabaseInfoInit 初始化站点信息默认配置
// 检查每个配置键是否存在，不存在则创建，已存在则跳过（不覆盖用户数据）
func (r *Reg) DatabaseInfoInit() {
	log := r.Serv.Logger.Named(xConsts.LogINIT).Sugar()
	log.Info("初始化站点信息默认配置")

	for _, config := range defaultInfoConfigs {
		var existing entity.System
		result := r.DB.Where("key = ?", config.Key).First(&existing)

		if result.Error != nil {
			if errors.Is(result.Error, gorm.ErrRecordNotFound) {
				// 不存在则创建
				newConfig := &entity.System{
					ID:    r.SnowflakeNode.Generate().Int64(),
					Key:   config.Key,
					Value: xUtil.Ptr(config.Value),
				}
				if err := r.DB.Create(newConfig).Error; err != nil {
					log.Errorf("创建默认配置 [%s] 失败: %v", config.Key, err)
				} else {
					log.Infof("创建默认配置 [%s] 成功", config.Key)
				}
			} else {
				log.Errorf("查询配置 [%s] 失败: %v", config.Key, result.Error)
			}
		}
		// 如果已存在则跳过，不覆盖用户设置的值
	}
}
