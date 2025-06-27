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

package link_color

import (
	"bamboo-main/internal/service"
	"context"
	"github.com/XiaoLFeng/bamboo-utils/blog"
	"github.com/XiaoLFeng/bamboo-utils/bresult"

	"bamboo-main/api/link_color/v1"
)

// LinkColorPage
//
// 获取友链颜色列表的控制器方法；接收一个 LinkColorPageReq 请求参数，返回一个 LinkColorPageRes 响应结果。
func (c *ControllerV1) LinkColorPage(ctx context.Context, req *v1.LinkColorPageReq) (res *v1.LinkColorPageRes, err error) {
	// 日志记录
	blog.ControllerInfo(ctx, "LinkColorPage", "获取友链颜色列表")

	// 调用逻辑层处理
	iColor := service.Color()
	pagedResult, errorCode := iColor.GetList(ctx, req.Search, req.Page, req.Size)
	if errorCode != nil {
		return nil, errorCode
	}

	return &v1.LinkColorPageRes{
		ResponseDTO: bresult.SuccessHasData(ctx, "获取成功", *pagedResult),
	}, nil
}
