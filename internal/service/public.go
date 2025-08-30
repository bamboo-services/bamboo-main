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
