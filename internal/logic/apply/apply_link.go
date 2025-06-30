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

package apply

import (
	"bamboo-main/internal/consts"
	"bamboo-main/internal/dao"
	"bamboo-main/internal/model/dto/base"
	"bamboo-main/internal/service"
	"context"
	"fmt"
	"github.com/XiaoLFeng/bamboo-utils/berror"
	"github.com/XiaoLFeng/bamboo-utils/blog"
	"github.com/gogf/gf/v2/database/gdb"
	"github.com/google/uuid"
	"time"
)

// NewApplyLink
//
// 创建一个新的友情链接申请并存储到数据库中。
// 如果存储过程中出现错误，返回对应的错误码。
// 此方法还会生成唯一的UUID用于标识新申请，并清理全局缓存数据。
func (s *sApply) NewApplyLink(ctx context.Context, newApplyLink base.LinkFriendDTO) *berror.ErrorCode {
	blog.ServiceInfo(ctx, "NewApplyLink", "申请新的友情链接 %s", newApplyLink.LinkName)

	// 生成新的 UUID
	if newApplyLink.LinkUuid == "" {
		newApplyLink.LinkUuid = uuid.New().String()
	}

	// 存入数据库中
	_, daoErr := dao.LinkContext.Ctx(ctx).Cache(gdb.CacheOption{
		Duration: 12 * time.Hour,
		Name:     fmt.Sprintf(consts.LinkContextByUuidRedisKey, newApplyLink.LinkUuid),
	}).OmitEmpty().Insert(&newApplyLink)
	if daoErr != nil {
		blog.ServiceError(ctx, "NewApplyLink", "插入新的友情链接申请失败", daoErr)
		return berror.ErrorAddData(&berror.ErrDatabaseError, "插入新的友情链接申请失败: "+daoErr.Error())
	}

	// 清理全局缓存
	errorCode := service.Friend().DeleteGlobalImpactCache(ctx)
	if errorCode != nil {
		return errorCode
	}

	return nil
}
