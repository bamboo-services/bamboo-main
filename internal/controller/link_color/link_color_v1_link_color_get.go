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
	"bamboo-main/internal/model/dto/base"
	"bamboo-main/internal/service"
	"context"
	"github.com/XiaoLFeng/bamboo-utils/berror"
	"github.com/XiaoLFeng/bamboo-utils/blog"
	"github.com/XiaoLFeng/bamboo-utils/bresult"
	"github.com/gogf/gf/v2/util/gconv"

	"bamboo-main/api/link_color/v1"
)

// LinkColorGet
//
// 获取友链颜色详情的控制器方法；接收一个 LinkColorGetReq 请求参数，返回一个 LinkColorGetRes 响应结果。
func (c *ControllerV1) LinkColorGet(ctx context.Context, req *v1.LinkColorGetReq) (res *v1.LinkColorGetRes, err error) {
	// 日志记录
	blog.ControllerInfo(ctx, "LinkColorGet", "获取友链颜色详情")

	// 调用逻辑层处理
	iColor := service.Color()
	linkColor, errorCode := iColor.GetOneByUUID(ctx, req.ColorUUID)
	if errorCode != nil {
		return nil, errorCode
	}

	// 格式转换
	var linkColorDTO base.LinkColorDTO
	operateErr := gconv.Struct(linkColor, &linkColorDTO)
	if operateErr != nil {
		return nil, &berror.ErrInternalServer
	}

	return &v1.LinkColorGetRes{
		ResponseDTO: bresult.SuccessHasData(ctx, "获取成功", linkColorDTO),
	}, nil
}
