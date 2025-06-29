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

package link_friend

import (
	"bamboo-main/internal/model/dto/base"
	"bamboo-main/internal/service"
	"context"
	"github.com/XiaoLFeng/bamboo-utils/berror"
	"github.com/XiaoLFeng/bamboo-utils/blog"
	"github.com/XiaoLFeng/bamboo-utils/bresult"
	"github.com/gogf/gf/v2/util/gconv"

	"bamboo-main/api/link_friend/v1"
)

// LinkFriendGet
//
// 获取友情链接详情的控制器方法；
// 接收一个 LinkFriendGetReq 请求参数，返回一个 LinkFriendGetRes 响应结果。
// 通过友链的UUID获取详细信息，包括名称、URL、头像、描述等属性。
func (c *ControllerV1) LinkFriendGet(ctx context.Context, req *v1.LinkFriendGetReq) (res *v1.LinkFriendGetRes, err error) {
	// 日志记录
	blog.ControllerInfo(ctx, "LinkFriendGet", "获取友情链接详情 %s", req.LinkUUID)

	// 调用逻辑层处理
	iFriend := service.Friend()
	friendEntity, errorCode := iFriend.GetOneByUUID(ctx, req.LinkUUID)
	if errorCode != nil {
		return nil, errorCode
	}

	// 数据结构转换
	var friendDTO base.LinkFriendDTO
	if err := gconv.Struct(friendEntity, &friendDTO); err != nil {
		blog.ControllerError(ctx, "LinkFriendGet", "转换友情链接实体到DTO失败: %v", err)
		return nil, berror.ErrorAddData(&berror.ErrInternalServer, "获取友情链接详情失败")
	}

	return &v1.LinkFriendGetRes{
		ResponseDTO: bresult.SuccessHasData(ctx, "获取成功", friendDTO),
	}, nil
}
