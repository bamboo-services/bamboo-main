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

// ILinkService 友情链接服务接口
type ILinkService interface {
	Add(ctx *gin.Context, req *request.LinkFriendAddReq) (*dto.LinkFriendDetailDTO, *xError.Error)
	Update(ctx *gin.Context, linkIDStr string, req *request.LinkFriendUpdateReq) (*dto.LinkFriendDetailDTO, *xError.Error)
	Delete(ctx *gin.Context, linkIDStr string) *xError.Error
	Get(ctx *gin.Context, linkIDStr string) (*dto.LinkFriendDetailDTO, *xError.Error)
	List(ctx *gin.Context, req *request.LinkFriendQueryReq) (*base.PaginationResponse[dto.LinkFriendDetailDTO], *xError.Error)
	UpdateStatus(ctx *gin.Context, linkIDStr string, req *request.LinkFriendStatusReq) *xError.Error
	UpdateFailStatus(ctx *gin.Context, linkIDStr string, req *request.LinkFriendFailReq) *xError.Error
	GetPublicLinks(ctx *gin.Context, groupIDStr string) ([]dto.LinkFriendDetailDTO, *xError.Error)
}

// NewLinkService 创建友情链接服务实例
func NewLinkService() ILinkService {
	return &logic.LinkLogic{}
}
