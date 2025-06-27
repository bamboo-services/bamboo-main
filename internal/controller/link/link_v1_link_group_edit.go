package link

import (
	"bamboo-main/internal/service"
	"context"
	"github.com/XiaoLFeng/bamboo-utils/blog"
	"github.com/XiaoLFeng/bamboo-utils/bresult"

	"bamboo-main/api/link/v1"
)

// LinkGroupEdit
//
// 编辑链接分组的控制器方法；接收编辑请求参数，执行分组更新操作并返回处理结果。
func (c *ControllerV1) LinkGroupEdit(ctx context.Context, req *v1.LinkGroupEditReq) (res *v1.LinkGroupEditRes, err error) {
	blog.ControllerInfo(ctx, "LinkGroupEdit", "编辑链接分组")

	// 调用逻辑层处理
	iGroup := service.Group()
	errorCode := iGroup.Update(ctx, req.UUID, req.Name, req.Description, req.Order)
	if errorCode != nil {
		return nil, errorCode
	}

	return &v1.LinkGroupEditRes{
		ResponseDTO: bresult.Success(ctx, "编辑成功"),
	}, nil
}
