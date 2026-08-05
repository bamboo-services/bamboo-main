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

package logic

import (
	"context"
	"strings"
	"time"

	apiInfo "github.com/bamboo-services/bamboo-main/api/info"
	"github.com/bamboo-services/bamboo-main/internal/entity"
	"github.com/bamboo-services/bamboo-main/internal/repository"
	bConst "github.com/bamboo-services/bamboo-main/pkg/constants"

	xError "github.com/bamboo-services/bamboo-base-go/common/error"
	xLog "github.com/bamboo-services/bamboo-base-go/common/log"
	xCtxUtil "github.com/bamboo-services/bamboo-base-go/major/utility/context"
)

// 站点信息键名常量
const (
	KeySiteName         = "site.name"
	KeySiteDescription  = "site.description"
	KeySiteIntroduction = "site.introduction"
	KeyProfileAbout     = "profile.about"

	// 申请站点展示键名（供「交换友链」场景读取：站点名字/描述/地址/图片/订阅/邮箱，
	// 展示于 operate/apply 申请页，语义独立于站点信息与博主信息）
	KeyApplySiteName  = "site.apply.name"
	KeyApplySiteDesc  = "site.apply.description"
	KeyApplySiteUrl   = "site.apply.url"
	KeyApplySiteImg   = "site.apply.image"
	KeyApplySiteRss   = "site.apply.rss"
	KeyApplySiteEmail = "site.apply.email"

	// 博主信息键名（供「关于我」名士帖个人展示：昵称/个人简介/博客链接/头像）
	KeyBloggerNick    = "blogger.nick"
	KeyBloggerDesc    = "blogger.description"
	KeyBloggerBlogUrl = "blogger.blog_url"
	KeyBloggerAvatar  = "blogger.avatar"
)

// infoRepo 站点信息仓储依赖集合
type infoRepo struct {
	system *repository.SystemRepo
}

// InfoLogic 站点信息业务逻辑
type InfoLogic struct {
	logic
	repo infoRepo
}

// NewInfoLogic 创建 InfoLogic 实例，从上下文获取数据库与缓存并初始化系统配置仓储依赖。
func NewInfoLogic(ctx context.Context) *InfoLogic {
	db := xCtxUtil.MustGetDB(ctx)
	m := xCtxUtil.MustGetCacheManager(ctx)

	return &InfoLogic{
		logic: logic{
			db:    db,
			cache: m,
			log:   xLog.WithName(xLog.NamedLOGC, "InfoLogic"),
		},
		repo: infoRepo{
			system: repository.NewSystemRepo(db, m),
		},
	}
}

// GetSiteInfo 获取站点信息
func (l *InfoLogic) GetSiteInfo(ctx context.Context) (*apiInfo.SiteResponse, *xError.Error) {
	// 批量查询站点相关配置
	keys := []string{KeySiteName, KeySiteIntroduction}
	configs, xErr := l.repo.system.ListByKeys(ctx, keys)
	if xErr != nil {
		return nil, xError.NewError(ctx, xError.DatabaseError, "获取站点信息失败", false, xErr)
	}

	// 转换为 map 便于访问
	configMap := make(map[string]*entity.System)
	for i := range configs {
		configMap[configs[i].Key] = &configs[i]
	}

	result := &apiInfo.SiteResponse{
		SiteName:     getConfigValue(configMap, KeySiteName),
		Introduction: getConfigValue(configMap, KeySiteIntroduction),
		UpdatedAt:    getLatestUpdateTime(configMap, keys),
	}

	return result, nil
}

// UpdateSiteInfo 更新站点信息
func (l *InfoLogic) UpdateSiteInfo(ctx context.Context, req *apiInfo.SiteUpdateRequest) (*apiInfo.SiteResponse, *xError.Error) {
	// 收集需要更新的字段（仅更新非 nil 的字段）
	updates := make(map[string]*string)
	if req.SiteName != nil {
		updates[KeySiteName] = req.SiteName
	}
	if req.Introduction != nil {
		updates[KeySiteIntroduction] = req.Introduction
	}

	// 如果没有任何更新字段，直接返回当前值
	if len(updates) == 0 {
		return l.GetSiteInfo(ctx)
	}

	// 执行更新
	for key, value := range updates {
		xErr := l.repo.system.UpdateValueByKey(ctx, key, value)
		if xErr != nil {
			return nil, xError.NewError(ctx, xError.DatabaseError, "更新站点信息失败", false, xErr)
		}
	}

	return l.GetSiteInfo(ctx)
}

// GetApplySiteInfo 获取申请站点展示
//
// 申请站点展示用于「交换友链」场景：站点名字/描述/地址/图片/订阅/邮箱，
// 供访客申请友链前在自站添加博主友链时复制，展示于 operate/apply 申请页。
// 与站点信息及博主个人展示语义独立，单独取数。
func (l *InfoLogic) GetApplySiteInfo(ctx context.Context) (*apiInfo.ApplySiteResponse, *xError.Error) {
	keys := []string{
		KeyApplySiteName, KeyApplySiteDesc, KeyApplySiteUrl,
		KeyApplySiteImg, KeyApplySiteRss, KeyApplySiteEmail,
	}
	configs, xErr := l.repo.system.ListByKeys(ctx, keys)
	if xErr != nil {
		return nil, xError.NewError(ctx, xError.DatabaseError, "获取申请站点展示失败", false, xErr)
	}

	configMap := make(map[string]*entity.System)
	for i := range configs {
		configMap[configs[i].Key] = &configs[i]
	}

	result := &apiInfo.ApplySiteResponse{
		SiteName:        getConfigValue(configMap, KeyApplySiteName),
		SiteDescription: getConfigValue(configMap, KeyApplySiteDesc),
		SiteUrl:         getConfigValue(configMap, KeyApplySiteUrl),
		SiteImage:       getConfigValue(configMap, KeyApplySiteImg),
		Rss:             getConfigValue(configMap, KeyApplySiteRss),
		Email:           getConfigValue(configMap, KeyApplySiteEmail),
		UpdatedAt:       getLatestUpdateTime(configMap, keys),
	}

	return result, nil
}

// UpdateApplySiteInfo 更新申请站点展示
func (l *InfoLogic) UpdateApplySiteInfo(ctx context.Context, req *apiInfo.ApplySiteUpdateRequest) (*apiInfo.ApplySiteResponse, *xError.Error) {
	updates := make(map[string]*string)
	if req.SiteName != nil {
		updates[KeyApplySiteName] = req.SiteName
	}
	if req.SiteDescription != nil {
		updates[KeyApplySiteDesc] = req.SiteDescription
	}
	if req.SiteUrl != nil {
		updates[KeyApplySiteUrl] = req.SiteUrl
	}
	if req.SiteImage != nil {
		updates[KeyApplySiteImg] = req.SiteImage
	}
	if req.Rss != nil {
		updates[KeyApplySiteRss] = req.Rss
	}
	if req.Email != nil {
		updates[KeyApplySiteEmail] = req.Email
	}

	if len(updates) == 0 {
		return l.GetApplySiteInfo(ctx)
	}

	for key, value := range updates {
		xErr := l.repo.system.UpdateValueByKey(ctx, key, value)
		if xErr != nil {
			return nil, xError.NewError(ctx, xError.DatabaseError, "更新申请站点展示失败", false, xErr)
		}
	}

	return l.GetApplySiteInfo(ctx)
}

// GetBloggerInfo 获取博主信息
//
// 博主信息用于「关于我」名士帖个人展示：昵称/个人简介/博客链接/头像。
// 与站点信息及申请站点展示语义独立，单独取数。
func (l *InfoLogic) GetBloggerInfo(ctx context.Context) (*apiInfo.BloggerResponse, *xError.Error) {
	keys := []string{
		KeyBloggerNick, KeyBloggerDesc, KeyBloggerBlogUrl, KeyBloggerAvatar,
	}
	configs, xErr := l.repo.system.ListByKeys(ctx, keys)
	if xErr != nil {
		return nil, xError.NewError(ctx, xError.DatabaseError, "获取博主信息失败", false, xErr)
	}

	configMap := make(map[string]*entity.System)
	for i := range configs {
		configMap[configs[i].Key] = &configs[i]
	}

	result := &apiInfo.BloggerResponse{
		Nick:        getConfigValue(configMap, KeyBloggerNick),
		Description: getConfigValue(configMap, KeyBloggerDesc),
		BlogUrl:     getConfigValue(configMap, KeyBloggerBlogUrl),
		Avatar:      getConfigValue(configMap, KeyBloggerAvatar),
		UpdatedAt:   getLatestUpdateTime(configMap, keys),
	}

	return result, nil
}

// UpdateBloggerInfo 更新博主信息
func (l *InfoLogic) UpdateBloggerInfo(ctx context.Context, req *apiInfo.BloggerUpdateRequest) (*apiInfo.BloggerResponse, *xError.Error) {
	updates := make(map[string]*string)
	if req.Nick != nil {
		updates[KeyBloggerNick] = req.Nick
	}
	if req.Description != nil {
		updates[KeyBloggerDesc] = req.Description
	}
	if req.BlogUrl != nil {
		updates[KeyBloggerBlogUrl] = req.BlogUrl
	}
	if req.Avatar != nil {
		updates[KeyBloggerAvatar] = req.Avatar
	}

	if len(updates) == 0 {
		return l.GetBloggerInfo(ctx)
	}

	for key, value := range updates {
		xErr := l.repo.system.UpdateValueByKey(ctx, key, value)
		if xErr != nil {
			return nil, xError.NewError(ctx, xError.DatabaseError, "更新博主信息失败", false, xErr)
		}
	}

	return l.GetBloggerInfo(ctx)
}

// GetArchiveInfo 获取站点档案
//
// 站点档案聚合「站点描述」（site.description）与「自我介绍」（profile.about），
// 均以 Markdown 书写，供设置页「站点档案」板块统一编辑与公开页 about/me 展示。
func (l *InfoLogic) GetArchiveInfo(ctx context.Context) (*apiInfo.ArchiveResponse, *xError.Error) {
	keys := []string{KeySiteDescription, KeyProfileAbout}
	configs, xErr := l.repo.system.ListByKeys(ctx, keys)
	if xErr != nil {
		return nil, xError.NewError(ctx, xError.DatabaseError, "获取站点档案失败", false, xErr)
	}

	configMap := make(map[string]*entity.System)
	for i := range configs {
		configMap[configs[i].Key] = &configs[i]
	}

	result := &apiInfo.ArchiveResponse{
		SiteDescription: getConfigValue(configMap, KeySiteDescription),
		About:           getConfigValue(configMap, KeyProfileAbout),
		UpdatedAt:       getLatestUpdateTime(configMap, keys),
	}

	return result, nil
}

// UpdateArchiveInfo 更新站点档案
func (l *InfoLogic) UpdateArchiveInfo(ctx context.Context, req *apiInfo.ArchiveUpdateRequest) (*apiInfo.ArchiveResponse, *xError.Error) {
	updates := make(map[string]*string)
	if req.SiteDescription != nil {
		updates[KeySiteDescription] = req.SiteDescription
	}
	if req.About != nil {
		updates[KeyProfileAbout] = req.About
	}

	if len(updates) == 0 {
		return l.GetArchiveInfo(ctx)
	}

	for key, value := range updates {
		xErr := l.repo.system.UpdateValueByKey(ctx, key, value)
		if xErr != nil {
			return nil, xError.NewError(ctx, xError.DatabaseError, "更新站点档案失败", false, xErr)
		}
	}

	return l.GetArchiveInfo(ctx)
}

// GetBuiltinInvalidGroup 获取内置「已失效」分组配置
//
// 名称/描述存于 bm_system（group.builtin.invalid.*），经 BuildBuiltinInvalidGroup 组装虚拟分组。
func (l *InfoLogic) GetBuiltinInvalidGroup(ctx context.Context) (*entity.LinkGroup, *xError.Error) {
	return l.repo.system.BuildBuiltinInvalidGroup(ctx)
}

// UpdateBuiltinInvalidGroup 更新内置「已失效」分组配置
func (l *InfoLogic) UpdateBuiltinInvalidGroup(ctx context.Context, req *apiInfo.BuiltinInvalidGroupUpdateRequest) (*entity.LinkGroup, *xError.Error) {
	updates := make(map[string]*string)

	if req.Name != nil {
		if strings.TrimSpace(*req.Name) == "" {
			return nil, xError.NewError(ctx, xError.BadRequest, "已失效分组名称不能为空", false)
		}
		updates[bConst.KeyBuiltinInvalidGroupName] = req.Name
	}
	if req.Description != nil {
		updates[bConst.KeyBuiltinInvalidGroupDesc] = req.Description
	}
	if len(updates) == 0 {
		return l.GetBuiltinInvalidGroup(ctx)
	}

	for key, value := range updates {
		if xErr := l.repo.system.UpdateValueByKey(ctx, key, value); xErr != nil {
			return nil, xErr
		}
	}

	return l.GetBuiltinInvalidGroup(ctx)
}

// ============ 辅助函数 ============

// getConfigValue 从配置 map 中获取值
func getConfigValue(configMap map[string]*entity.System, key string) string {
	if config, ok := configMap[key]; ok && config.Value != nil {
		return *config.Value
	}
	return ""
}

// getLatestUpdateTime 获取最新的更新时间
func getLatestUpdateTime(configMap map[string]*entity.System, keys []string) time.Time {
	var latest time.Time
	for _, key := range keys {
		if config, ok := configMap[key]; ok {
			if config.UpdatedAt.After(latest) {
				latest = config.UpdatedAt
			}
		}
	}
	if latest.IsZero() {
		latest = time.Now()
	}
	return latest
}
