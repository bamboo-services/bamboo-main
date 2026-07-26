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

package handler

import (
	apiPublic "github.com/bamboo-services/bamboo-main/api/public"

	xResult "github.com/bamboo-services/bamboo-base-go/major/result"
	"github.com/gin-gonic/gin"
)

// HealthCheck 健康检查接口
//
// @Summary [用户] 系统健康检查
// @Description 检查系统、数据库、Redis连接状态
// @Tags 公开接口
// @Accept json
// @Produce json
// @Success 200 {object} xBase.BaseResponse{data=apiPublic.HealthResponse} "健康检查成功"
// @Failure 500 {object} xBase.BaseResponse "服务器内部错误"
// @Router /api/v1/public/health [GET]
func (h *PublicHandler) HealthCheck(c *gin.Context) {
	// 调用服务层
	result, err := h.service.publicLogic.HealthCheck(c.Request.Context())
	if err != nil {
		_ = c.Error(err)
		return
	}

	// 返回成功响应
	resp := apiPublic.HealthResponse{}
	if result != nil {
		resp = *result
	}
	xResult.SuccessHasData(c, "系统状态正常", resp)
}

// Ping 简单连通性测试接口
//
// @Summary [用户] 连通性测试
// @Description 简单的服务连通性测试
// @Tags 公开接口
// @Accept json
// @Produce json
// @Success 200 {object} xBase.BaseResponse{data=apiPublic.PingResponse} "连通测试成功"
// @Failure 500 {object} xBase.BaseResponse "服务器内部错误"
// @Router /api/v1/public/ping [GET]
func (h *PublicHandler) Ping(c *gin.Context) {
	// 调用服务层
	result, err := h.service.publicLogic.Ping(c.Request.Context())
	if err != nil {
		_ = c.Error(err)
		return
	}

	// 返回成功响应
	resp := apiPublic.PingResponse{}
	if result != nil {
		resp = *result
	}
	xResult.SuccessHasData(c, "pong", resp)
}
