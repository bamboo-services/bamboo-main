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

// LinkColorEdit
//
// 编辑链接颜色的控制器方法；接收编辑请求参数，执行颜色更新操作并返回处理结果。
func (c *ControllerV1) LinkColorEdit(ctx context.Context, req *v1.LinkColorEditReq) (res *v1.LinkColorEditRes, err error) {
	blog.ControllerInfo(ctx, "LinkColorEdit", "编辑链接颜色")

	// 调用逻辑层处理
	iColor := service.Color()
	errorCode := iColor.Update(ctx, req.UUID, req.Name, req.Value, req.Description)
	if errorCode != nil {
		return nil, errorCode
	}

	return &v1.LinkColorEditRes{
		ResponseDTO: bresult.Success(ctx, "编辑成功"),
	}, nil
}
