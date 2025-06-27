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

// LinkGroupAdd
//
// 添加友链分组的控制器方法；
// 该方法接收一个 LinkGroupAddReq 请求参数，包含分组名称、描述和排序等信息；
// 返回一个 LinkGroupAddRes 响应结果，包含创建成功的分组信息。
func (c *ControllerV1) LinkGroupAdd(ctx context.Context, req *v1.LinkGroupAddReq) (res *v1.LinkGroupAddRes, err error) {
	blog.ControllerInfo(ctx, "LinkGroupAdd", "添加链接分组")

	// 调用逻辑层处理
	iGroup := service.Group()
	linkGroup, errorCode := iGroup.Create(ctx, req.Name, req.Description, req.Order)
	if errorCode != nil {
		return nil, errorCode
	}

	// 结构处理
	var linkGroupDTO base.LinkGroupDTO
	operateErr := gconv.Struct(linkGroup, &linkGroupDTO)
	if operateErr != nil {
		blog.ControllerError(ctx, "LinkGroupAdd", "转换链接分组数据失败: %v", operateErr)
		return nil, &berror.ErrInternalServer
	}

	return &v1.LinkGroupAddRes{
		ResponseDTO: bresult.SuccessHasData(ctx, "创建成功", linkGroupDTO),
	}, nil
}
