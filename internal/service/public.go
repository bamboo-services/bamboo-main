package service

import (
	"bamboo-main/internal/logic"
	"bamboo-main/internal/model/dto/response"
	xError "github.com/bamboo-services/bamboo-base-go/error"
	"github.com/gin-gonic/gin"
)

// IPublicService 公开接口服务接口
type IPublicService interface {
	HealthCheck(ctx *gin.Context) (*response.HealthResponse, *xError.Error)
	Ping(ctx *gin.Context) (*response.PingResponse, *xError.Error)
}

// NewPublicService 创建公开接口服务实例，返回Logic实现
func NewPublicService() *logic.PublicLogic {
	return &logic.PublicLogic{}
}