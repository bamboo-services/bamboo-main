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
	"bamboo-main/internal/model/base"
	"bamboo-main/internal/model/dto"
	"bamboo-main/internal/model/request"

	xError "github.com/bamboo-services/bamboo-base-go/error"
	"github.com/gin-gonic/gin"
)

// ISponsorChannelService 赞助渠道服务接口
type ISponsorChannelService interface {
	Add(ctx *gin.Context, req *request.SponsorChannelAddReq) (*dto.SponsorChannelDetailDTO, *xError.Error)
	Update(ctx *gin.Context, idStr string, req *request.SponsorChannelUpdateReq) (*dto.SponsorChannelDetailDTO, *xError.Error)
	UpdateStatus(ctx *gin.Context, idStr string, req *request.SponsorChannelStatusReq) (bool, *xError.Error)
	Delete(ctx *gin.Context, idStr string) *xError.Error
	Get(ctx *gin.Context, idStr string) (*dto.SponsorChannelDetailDTO, *xError.Error)
	GetList(ctx *gin.Context, req *request.SponsorChannelListReq) ([]dto.SponsorChannelListDTO, *xError.Error)
	GetPage(ctx *gin.Context, req *request.SponsorChannelPageReq) (*base.PaginationResponse[dto.SponsorChannelNormalDTO], *xError.Error)
	GetPublicList(ctx *gin.Context) ([]dto.SponsorChannelListDTO, *xError.Error)
}

// NewSponsorChannelService 创建赞助渠道服务实例
func NewSponsorChannelService() ISponsorChannelService {
	return &logic.SponsorChannelLogic{}
}
