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

package logic

import (
	"context"
	"time"

	apiInfo "github.com/bamboo-services/bamboo-main/api/info"
	"github.com/bamboo-services/bamboo-main/internal/entity"
	"github.com/bamboo-services/bamboo-main/internal/repository"

	xError "github.com/bamboo-services/bamboo-base-go/error"
	xLog "github.com/bamboo-services/bamboo-base-go/log"
	xCtxUtil "github.com/bamboo-services/bamboo-base-go/utility/ctxutil"
	"github.com/gin-gonic/gin"
)

// 站点信息键名常量
const (
	KeySiteName         = "site.name"
	KeySiteDescription  = "site.description"
	KeySiteIntroduction = "site.introduction"
	KeyProfileAbout     = "profile.about"
)

// InfoLogic 站点信息业务逻辑
type infoRepo struct {
	system *repository.SystemRepo
}

type InfoLogic struct {
	logic
	repo infoRepo
}

func NewInfoLogic(ctx context.Context) *InfoLogic {
	db := xCtxUtil.MustGetDB(ctx)
	rdb := xCtxUtil.MustGetRDB(ctx)

	return &InfoLogic{
		logic: logic{
			db:  db,
			rdb: rdb,
			log: xLog.WithName(xLog.NamedLOGC, "InfoLogic"),
		},
		repo: infoRepo{
			system: repository.NewSystemRepo(db, rdb),
		},
	}
}

// GetSiteInfo 获取站点信息
func (l *InfoLogic) GetSiteInfo(ctx *gin.Context) (*apiInfo.SiteResponse, *xError.Error) {
	// 批量查询站点相关配置
	keys := []string{KeySiteName, KeySiteDescription, KeySiteIntroduction}
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
		SiteName:        getConfigValue(configMap, KeySiteName),
		SiteDescription: getConfigValue(configMap, KeySiteDescription),
		Introduction:    getConfigValue(configMap, KeySiteIntroduction),
		UpdatedAt:       getLatestUpdateTime(configMap, keys),
	}

	return result, nil
}

// UpdateSiteInfo 更新站点信息
func (l *InfoLogic) UpdateSiteInfo(ctx *gin.Context, req *apiInfo.SiteUpdateRequest) (*apiInfo.SiteResponse, *xError.Error) {
	// 收集需要更新的字段（仅更新非 nil 的字段）
	updates := make(map[string]*string)
	if req.SiteName != nil {
		updates[KeySiteName] = req.SiteName
	}
	if req.SiteDescription != nil {
		updates[KeySiteDescription] = req.SiteDescription
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

// GetAbout 获取自我介绍
func (l *InfoLogic) GetAbout(ctx *gin.Context) (*apiInfo.AboutResponse, *xError.Error) {
	config, found, xErr := l.repo.system.GetByKey(ctx, KeyProfileAbout)
	if xErr != nil {
		return nil, xError.NewError(ctx, xError.DatabaseError, "获取自我介绍失败", false, xErr)
	}
	if !found {
		return nil, xError.NewError(ctx, xError.DatabaseError, "获取自我介绍失败", false)
	}

	content := ""
	if config.Value != nil {
		content = *config.Value
	}

	return &apiInfo.AboutResponse{
		Content:   content,
		UpdatedAt: config.UpdatedAt,
	}, nil
}

// UpdateAbout 更新自我介绍
func (l *InfoLogic) UpdateAbout(ctx *gin.Context, req *apiInfo.AboutUpdateRequest) (*apiInfo.AboutResponse, *xError.Error) {
	content := req.Content
	xErr := l.repo.system.UpdateValueByKey(ctx, KeyProfileAbout, &content)
	if xErr != nil {
		return nil, xError.NewError(ctx, xError.DatabaseError, "更新自我介绍失败", false, xErr)
	}

	return l.GetAbout(ctx)
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
