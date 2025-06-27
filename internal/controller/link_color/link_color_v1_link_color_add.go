package link_color

import (
	v1 "bamboo-main/api/link_color/v1"
	"bamboo-main/internal/model/dto/base"
	"bamboo-main/internal/service"
	"context"
	"github.com/XiaoLFeng/bamboo-utils/berror"
	"github.com/XiaoLFeng/bamboo-utils/blog"
	"github.com/XiaoLFeng/bamboo-utils/bresult"
	"github.com/gogf/gf/v2/util/gconv"
)

// LinkColorAdd
//
// 添加友链颜色的控制器方法；
// 该方法接收一个 LinkColorAddReq 请求参数，包含颜色名称、颜色值和描述等信息；
// 返回一个 LinkColorAddRes 响应结果，包含创建成功的颜色信息。
func (c *ControllerV1) LinkColorAdd(ctx context.Context, req *v1.LinkColorAddReq) (res *v1.LinkColorAddRes, err error) {
	blog.ControllerInfo(ctx, "LinkColorAdd", "添加链接颜色")

	// 调用逻辑层处理
	iColor := service.Color()
	linkColor, errorCode := iColor.Create(ctx, req.Name, req.Value, req.Description)
	if errorCode != nil {
		return nil, errorCode
	}

	// 结构处理
	var linkColorDTO base.LinkColorDTO
	operateErr := gconv.Struct(linkColor, &linkColorDTO)
	if operateErr != nil {
		blog.ControllerError(ctx, "LinkColorAdd", "转换链接颜色数据失败: %v", operateErr)
		return nil, &berror.ErrInternalServer
	}

	return &v1.LinkColorAddRes{
		ResponseDTO: bresult.SuccessHasData(ctx, "创建成功", linkColorDTO),
	}, nil
}
