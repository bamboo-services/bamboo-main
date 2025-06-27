package link_group

import (
	"bamboo-main/internal/service"
	"context"
	"github.com/XiaoLFeng/bamboo-utils/blog"
	"github.com/XiaoLFeng/bamboo-utils/bresult"

	"bamboo-main/api/link_group/v1"
)

// LinkGroupDelete
//
// 删除链接分组的控制器方法；接收删除请求参数，执行分组删除操作并返回处理结果。
func (c *ControllerV1) LinkGroupDelete(ctx context.Context, req *v1.LinkGroupDeleteReq) (res *v1.LinkGroupDeleteRes, err error) {
	blog.ControllerInfo(ctx, "LinkGroupDelete", "删除链接分组")

	// 调用逻辑层处理
	iGroup := service.Group()
	errorCode := iGroup.Delete(ctx, req.GroupUuid)
	if errorCode != nil {
		return nil, errorCode
	}

	return &v1.LinkGroupDeleteRes{
		ResponseDTO: bresult.Success(ctx, "删除成功"),
	}, nil
}
