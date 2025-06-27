package link_group

import (
	"bamboo-main/internal/model/dto/base"
	"bamboo-main/internal/service"
	"context"
	"github.com/XiaoLFeng/bamboo-utils/berror"
	"github.com/XiaoLFeng/bamboo-utils/blog"
	"github.com/XiaoLFeng/bamboo-utils/bresult"
	"github.com/gogf/gf/v2/util/gconv"

	"bamboo-main/api/link_group/v1"
)

// LinkGroupGet
//
// 获取友链分组详情的控制器方法；接收一个 LinkGroupGetReq 请求参数，返回一个 LinkGroupGetRes 响应结果。
func (c *ControllerV1) LinkGroupGet(ctx context.Context, req *v1.LinkGroupGetReq) (res *v1.LinkGroupGetRes, err error) {
	// 日志记录
	blog.ControllerInfo(ctx, "LinkGroupGet", "获取友链分组详情")

	// 调用逻辑层处理
	iGroup := service.Group()
	linkGroup, errorCode := iGroup.GetOneByUUID(ctx, req.GroupUUID)
	if errorCode != nil {
		return nil, errorCode
	}

	// 格式转换
	var linkGroupDTO base.LinkGroupDTO
	operateErr := gconv.Struct(linkGroup, &linkGroupDTO)
	if operateErr != nil {
		return nil, &berror.ErrInternalServer
	}

	return &v1.LinkGroupGetRes{
		ResponseDTO: bresult.SuccessHasData(ctx, "获取成功", linkGroupDTO),
	}, nil
}
