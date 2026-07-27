// --------------------------------------------------------------------------------
// Copyright (c) 2016-NOW(至今) 筱锋
// Author: 筱锋「xiao_lfeng」(https://www.x-lf.com)
// --------------------------------------------------------------------------------
// 许可证声明：版权所有 (c) 2016-2026 筱锋。保留所有权利。
// 有关MIT许可证的更多信息，请查看项目根目录下的LICENSE文件或访问：
// https://opensource.org/licenses/MIT
// --------------------------------------------------------------------------------

package handler

import (
	xResult "github.com/bamboo-services/bamboo-base-go/major/result"
	"github.com/gin-gonic/gin"
)

// Stats 获取仪表盘统计数据
//
// @Summary [管理] 获取仪表盘统计数据
// @Description 管理员获取友链总数、待审核、已通过计数及最近友链申请列表
// @Tags 仪表盘接口
// @Accept json
// @Produce json
// @Security Bearer
// @Success 200 {object} xBase.BaseResponse{data=apiDashboard.StatsResponse} "获取成功"
// @Failure 401 {object} xBase.BaseResponse "未认证"
// @Failure 500 {object} xBase.BaseResponse "服务器内部错误"
// @Router /api/v1/admin/dashboard/stats [GET]
func (h *DashboardHandler) Stats(c *gin.Context) {
	result, err := h.service.dashboardLogic.Stats(c.Request.Context())
	if err != nil {
		_ = c.Error(err)
		return
	}

	xResult.SuccessHasData(c, "获取仪表盘统计成功", result)
}
