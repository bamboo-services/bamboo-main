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
	Add(ctx *gin.Context, req *request.LinkFriendAddReq) (*dto.LinkFriendDTO, *xError.Error)
	Update(ctx *gin.Context, linkUUID string, req *request.LinkFriendUpdateReq) (*dto.LinkFriendDTO, *xError.Error)
	Delete(ctx *gin.Context, linkUUID string) *xError.Error
	Get(ctx *gin.Context, linkUUID string) (*dto.LinkFriendDTO, *xError.Error)
	List(ctx *gin.Context, req *request.LinkFriendQueryReq) (*base.PaginationResponse[dto.LinkFriendDTO], *xError.Error)
	UpdateStatus(ctx *gin.Context, linkUUID string, req *request.LinkFriendStatusReq) *xError.Error
	UpdateFailStatus(ctx *gin.Context, linkUUID string, req *request.LinkFriendFailReq) *xError.Error
	GetPublicLinks(ctx *gin.Context, groupUUID string) ([]dto.LinkFriendDTO, *xError.Error)
}

// NewLinkService 创建友情链接服务实例，返回Logic实现
func NewLinkService() *logic.LinkLogic {
	return &logic.LinkLogic{}
}
