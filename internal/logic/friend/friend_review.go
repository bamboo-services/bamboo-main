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

package friend

import (
	"bamboo-main/internal/consts"
	"bamboo-main/internal/dao"
	"bamboo-main/internal/model/do"
	"bamboo-main/internal/model/entity"
	"context"
	"fmt"
	"github.com/XiaoLFeng/bamboo-utils/berror"
	"github.com/XiaoLFeng/bamboo-utils/blog"
	"github.com/gogf/gf/v2/database/gdb"
	"github.com/google/uuid"
	"strconv"
	"time"
)

// UpdateStatus
//
// 更新友链状态，根据传入的 UUID 和状态值进行操作。
// 如果 UUID 格式无效或状态值非法，将返回对应的错误码。
// 执行成功后会清理相关的全局缓存以确保数据一致性。
func (s *sFriend) UpdateStatus(ctx context.Context, linkUUID string, status int) *berror.ErrorCode {
	blog.ServiceInfo(ctx, "UpdateStatus", "更新友链状态")
	// 检查 UUID 是否有效
	parseUUID, uuidErr := uuid.Parse(linkUUID)
	if uuidErr != nil {
		blog.ServiceError(ctx, "UpdateStatus", "UUID 格式错误", uuidErr)
		return berror.ErrorAddData(&berror.ErrInvalidParameters, "UUID 格式错误: "+uuidErr.Error())
	}
	// 检查状态值
	if status != 1 && status != 2 {
		blog.ServiceError(ctx, "UpdateStatus", "状态值错误", "Status: "+strconv.Itoa(status))
		return berror.ErrorAddData(&berror.ErrInvalidParameters, "状态值错误")
	}
	// 检查是否存在该友链
	var hasExistFriend *entity.LinkContext
	daoErr := dao.LinkContext.Ctx(ctx).Cache(gdb.CacheOption{
		Duration: 7 * 24 * time.Hour,
		Name:     fmt.Sprintf(consts.LinkContextByUuidRedisKey, parseUUID.String()),
	}).WherePri(parseUUID.String()).Scan(&hasExistFriend)
	if daoErr != nil {
		blog.ServiceError(ctx, "UpdateStatus", "查询友链失败", daoErr)
		return berror.ErrorAddData(&berror.ErrDatabaseError, "查询友链失败: "+daoErr.Error())
	}

	// 执行更新操作
	_, daoErr = dao.LinkContext.Ctx(ctx).Cache(gdb.CacheOption{
		Duration: -1,
		Name:     fmt.Sprintf(consts.LinkContextByUuidRedisKey, parseUUID.String()),
	}).WherePri(hasExistFriend.LinkUuid).OmitEmpty().Data(&do.LinkContext{LinkStatus: status}).Update()
	if daoErr != nil {
		blog.ServiceError(ctx, "UpdateStatus", "更新友链状态失败", daoErr)
		return berror.ErrorAddData(&berror.ErrDatabaseError, "更新友链状态失败: "+daoErr.Error())
	}

	// 删除相关的全局缓存
	errorCode := s.DeleteGlobalImpactCache(ctx)
	if errorCode != nil {
		return errorCode
	}
	return nil
}

// UpdateFailStatus
//
// 更新友链失败状态和原因，根据传入的 UUID、失败状态和原因进行操作。
// 如果 UUID 格式无效，将返回对应的错误码。
// 执行成功后会清理相关的全局缓存以确保数据一致性。
func (s *sFriend) UpdateFailStatus(ctx context.Context, linkUUID string, isFail bool, reason string) *berror.ErrorCode {
	blog.ServiceInfo(ctx, "UpdateFailStatus", "更新友链失败状态")
	// 检查 UUID 是否有效
	parseUUID, uuidErr := uuid.Parse(linkUUID)
	if uuidErr != nil {
		blog.ServiceError(ctx, "UpdateFailStatus", "UUID 格式错误", uuidErr)
		return berror.ErrorAddData(&berror.ErrInvalidParameters, "UUID 格式错误: "+uuidErr.Error())
	}
	// 执行查询操作，检查友链是否存在
	var hasExistFriend *entity.LinkContext
	daoErr := dao.LinkContext.Ctx(ctx).Cache(gdb.CacheOption{
		Duration: 7 * 24 * time.Hour,
		Name:     fmt.Sprintf(consts.LinkContextByUuidRedisKey, parseUUID.String()),
	}).WherePri(parseUUID.String()).Scan(&hasExistFriend)
	if daoErr != nil {
		blog.ServiceError(ctx, "UpdateFailStatus", "查询友链失败", daoErr)
		return berror.ErrorAddData(&berror.ErrDatabaseError, "查询友链失败: "+daoErr.Error())
	}

	// 执行更新操作
	_, daoErr = dao.LinkContext.Ctx(ctx).Cache(gdb.CacheOption{
		Duration: -1,
		Name:     fmt.Sprintf(consts.LinkContextByUuidRedisKey, parseUUID.String()),
	}).WherePri(parseUUID.String()).OmitEmpty().Data(&do.LinkContext{
		LinkFail:       isFail,
		LinkFailReason: reason,
	}).Update()
	if daoErr != nil {
		blog.ServiceError(ctx, "UpdateFailStatus", "更新友链失败状态失败", daoErr)
		return berror.ErrorAddData(&berror.ErrDatabaseError, "更新友链失败状态失败: "+daoErr.Error())
	}

	// 删除相关的全局缓存
	errorCode := s.DeleteGlobalImpactCache(ctx)
	if errorCode != nil {
		return errorCode
	}
	return nil
}
