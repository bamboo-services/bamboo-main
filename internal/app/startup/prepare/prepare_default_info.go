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

	xLog "github.com/bamboo-services/bamboo-base-go/common/log"
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
	// 博主信息（供「交换友链」场景，访客申请前需先在自站添加博主友链）
	{Key: "blogger.site_name", Value: "凌中的锋雨"},
	{Key: "blogger.site_description", Value: "不为如何，只为在茫茫人海中有自己的一片天空~"},
	{Key: "blogger.site_url", Value: "https://www.x-lf.com"},
	{Key: "blogger.site_image", Value: "https://i-cdn.akass.cn/2024/05/664870a814c0d.png!wp60"},
	{Key: "blogger.rss", Value: "https://blog.x-lf.com/atom.xml"},
	{Key: "blogger.email", Value: "gm@x-lf.cn"},
	// 博主个人展示信息（供「关于我」名士帖：昵称/个人简介/博客链接/头像）
	{Key: "blogger.nick", Value: "筱锋"},
	{Key: "blogger.description", Value: "一个热爱技术的开发者，喜欢折腾各种有趣的东西。"},
	{Key: "blogger.blog_url", Value: "https://blog.x-lf.com"},
	{Key: "blogger.avatar", Value: "https://i-cdn.akass.cn/2024/05/664870a814c0d.png!wp60"},
	// 内置「已失效」分组配置（供友链失效自动归集与公开「已失效」章节展示）
	{Key: "group.builtin.invalid.name", Value: "已失效"},
	{Key: "group.builtin.invalid.description", Value: ""},
}

// DefaultInfo 初始化默认站点信息配置的 xOptionDB.PrepareFunc。
//
// 幂等：按 key 跳过已存在的条目，仅写入缺失项。失败时返回 error。
func DefaultInfo(ctx context.Context, db *gorm.DB) error {
	log := xLog.WithName(xLog.NamedINIT)

	for _, item := range defaultInfoConfigs {
		var existing entity.System
		err := db.WithContext(ctx).Where("key = ?", item.Key).First(&existing).Error
		if err == nil {
			continue
		}
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}

		value := item.Value
		if err = db.WithContext(ctx).Create(&entity.System{Key: item.Key, Value: &value}).Error; err != nil {
			return err
		}
	}

	log.Info(ctx, "默认站点信息配置初始化完成")
	return nil
}
