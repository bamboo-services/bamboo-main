package prepare

import (
	"errors"

	"github.com/bamboo-services/bamboo-main/internal/entity"
	"gorm.io/gorm"
)

var defaultInfoConfigs = []struct {
	Key   string
	Value string
}{
	{Key: "site.name", Value: "Bamboo Links"},
	{Key: "site.description", Value: "一个简洁优雅的友情链接管理系统"},
	{Key: "site.introduction", Value: "欢迎来到我的友链主页！这里收录了我的好朋友们的博客链接。"},
	{Key: "profile.about", Value: "# 关于我\n\n这里是自我介绍，支持 **Markdown** 格式。"},
}

func (p *Prepare) prepareDefaultInfo() error {
	for _, item := range defaultInfoConfigs {
		var existing entity.System
		err := p.db.WithContext(p.ctx).Where("key = ?", item.Key).First(&existing).Error
		if err == nil {
			continue
		}
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}

		value := item.Value
		if err = p.db.WithContext(p.ctx).Create(&entity.System{Key: item.Key, Value: &value}).Error; err != nil {
			return err
		}
	}

	return nil
}
