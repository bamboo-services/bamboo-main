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

package service

import (
	"bamboo-main/internal/logic"
	"bamboo-main/internal/model/dto"
	"bamboo-main/internal/model/request"

	xError "github.com/bamboo-services/bamboo-base-go/error"
	"github.com/gin-gonic/gin"
)

// IInfoService 站点信息服务接口
type IInfoService interface {
	// GetSiteInfo 获取站点信息
	GetSiteInfo(ctx *gin.Context) (*dto.SiteInfoDTO, *xError.Error)

	// UpdateSiteInfo 更新站点信息
	UpdateSiteInfo(ctx *gin.Context, req *request.SiteInfoUpdateReq) (*dto.SiteInfoDTO, *xError.Error)

	// GetAbout 获取自我介绍
	GetAbout(ctx *gin.Context) (*dto.AboutDTO, *xError.Error)

	// UpdateAbout 更新自我介绍
	UpdateAbout(ctx *gin.Context, req *request.AboutUpdateReq) (*dto.AboutDTO, *xError.Error)
}

// NewInfoService 创建站点信息服务实例
func NewInfoService() IInfoService {
	return &logic.InfoLogic{}
}
