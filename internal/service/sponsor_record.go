/*
 * --------------------------------------------------------------------------------
 * Copyright (c) 2016-NOW(至今) 筱锋
 * Author: 筱锋「xiao_lfeng」(https://www.x-lf.com)
 * --------------------------------------------------------------------------------
 * 许可证声明:版权所有 (c) 2016-2025 筱锋。保留所有权利。
 * 有关MIT许可证的更多信息,请查看项目根目录下的LICENSE文件或访问:
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

// ISponsorRecordService 赞助记录服务接口
type ISponsorRecordService interface {
	Add(ctx *gin.Context, req *request.SponsorRecordAddReq) (*dto.SponsorRecordDetailDTO, *xError.Error)
	Update(ctx *gin.Context, idStr string, req *request.SponsorRecordUpdateReq) (*dto.SponsorRecordDetailDTO, *xError.Error)
	Delete(ctx *gin.Context, idStr string) *xError.Error
	Get(ctx *gin.Context, idStr string) (*dto.SponsorRecordDetailDTO, *xError.Error)
	GetPage(ctx *gin.Context, req *request.SponsorRecordPageReq) (*base.PaginationResponse[dto.SponsorRecordNormalDTO], *xError.Error)
	GetPublicPage(ctx *gin.Context, req *request.SponsorRecordPublicPageReq) (*base.PaginationResponse[dto.SponsorRecordSimpleDTO], *xError.Error)
}

// NewSponsorRecordService 创建赞助记录服务实例
func NewSponsorRecordService() ISponsorRecordService {
	return &logic.SponsorRecordLogic{}
}
