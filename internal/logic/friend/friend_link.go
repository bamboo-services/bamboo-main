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
	"bamboo-main/internal/model/dto/base"
	"bamboo-main/internal/model/entity"
	"bamboo-main/pkg/cerror"
	"bamboo-main/pkg/utility"
	"context"
	"fmt"
	"github.com/XiaoLFeng/bamboo-utils/berror"
	"github.com/XiaoLFeng/bamboo-utils/blog"
	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/google/uuid"
	"strings"
	"time"
)

// AddFriend
//
// 添加一个新的友链记录至数据库，同时检查是否存在重复的友链，若重复则返回对应错误。
// 该函数会自动生成友链的唯一标识符（UUID），并设置友链的状态为启用（1）和创建、更新时间为当前时间。
// 如果添加成功，则会删除相关的全局缓存，以确保数据的一致性和最新性。
func (s *sFriend) AddFriend(ctx context.Context, friend base.LinkFriendDTO) *berror.ErrorCode {
	blog.ServiceInfo(ctx, "AddFriend", "添加友链")

	// 添加缺失的内容
	friend.LinkUuid = uuid.New().String()
	friend.LinkStatus = 1
	friend.LinkCreatedAt = gtime.Now()
	friend.LinkUpdatedAt = gtime.Now()

	// 检查是否已经存在相同的友链
	getDomain := utility.GetBaseDomain(friend.LinkUrl)
	var hasExistFriend *entity.LinkContext
	daoErr := dao.LinkContext.Ctx(ctx).WhereLike("link_url", "%"+getDomain+"%").Scan(&hasExistFriend)
	if daoErr != nil {
		blog.ServiceError(ctx, "AddFriend", "查询友链失败", daoErr)
		return berror.ErrorAddData(&berror.ErrDatabaseError, "查询友链失败: "+daoErr.Error())
	}
	if hasExistFriend != nil {
		blog.ServiceNotice(ctx, "AddFriend", "友链已存在", "Domain: "+getDomain)
		return cerror.ErrDomainExisted
	}

	// 执行插入操作
	_, daoErr = dao.LinkContext.Ctx(ctx).Cache(gdb.CacheOption{
		Duration: 7 * 24 * time.Hour,
		Name:     fmt.Sprintf(consts.LinkContextByUuidRedisKey, friend.LinkUuid),
	}).OmitEmpty().Insert(&friend)
	if daoErr != nil {
		blog.ServiceError(ctx, "AddFriend", "插入友链失败", daoErr)
		return berror.ErrorAddData(&berror.ErrDatabaseError, "插入友链失败: "+daoErr.Error())
	}

	// 删除关于友链的全局缓存
	errorCode := deleteGlobalImpactCache(ctx)
	if errorCode != nil {
		return errorCode
	}

	return nil
}

// GetOneByUUID
//
// 根据传入的 UUID 查询并返回友链信息，不存在或发生错误时返回对应的错误码。
// 该函数会检查传入的 UUID 是否有效，并使用缓存机制来提高查询效率。
// 如果查询成功，则返回友链实体；如果查询失败或未找到友链信息，则返回错误码。
func (s *sFriend) GetOneByUUID(ctx context.Context, linkUUID string) (*entity.LinkContext, *berror.ErrorCode) {
	blog.ServiceInfo(ctx, "GetOneByUUID", "获取友链信息 UUID: %s", linkUUID)
	// 检查 UUID 是否有效
	parseUUID, uuidErr := uuid.Parse(linkUUID)
	if uuidErr != nil {
		blog.ServiceError(ctx, "GetOneByUUID", "无效的友链 UUID: %s, 错误: %v", linkUUID, uuidErr)
		return nil, berror.ErrorAddData(&berror.ErrInvalidParameters, "无效的友链 UUID: "+uuidErr.Error())
	}
	// 查询友链信息
	var linkEntity *entity.LinkContext
	daoErr := dao.LinkContext.Ctx(ctx).Cache(gdb.CacheOption{
		Duration: 7 * 24 * time.Hour,
		Name:     fmt.Sprintf(consts.LinkContextByUuidRedisKey, parseUUID.String()),
	}).Where(&do.LinkContext{LinkUuid: parseUUID.String()}).Scan(&linkEntity)
	if daoErr != nil {
		blog.ServiceError(ctx, "GetOneByUUID", "查询友链信息失败，错误：%v", daoErr)
		return nil, berror.ErrorAddData(&berror.ErrDatabaseError, daoErr.Error())
	}
	if linkEntity == nil {
		blog.ServiceNotice(ctx, "GetOneByUUID", "未找到友链信息，UUID: %s", linkUUID)
		return nil, berror.ErrorAddData(&berror.ErrNotFound, "未找到友链信息")
	}
	return linkEntity, nil
}

// Update
//
// 更新友链信息。如果友链不存在或更新操作失败，返回对应的错误码。更新成功后会删除相关的全局缓存。
// 该函数会检查友链是否存在，并更新其信息，包括名称、URL、头像、RSS、描述、邮箱、分组、颜色、排序和审核备注等字段。
// 如果更新成功，则会删除相关的全局缓存，以确保数据的一致性和最新性。
// 如果更新失败，则返回错误码。
func (s *sFriend) Update(ctx context.Context, editFriendEntity base.LinkFriendDTO) *berror.ErrorCode {
	blog.ServiceInfo(ctx, "Update", "更新友链信息，UUID: %s", editFriendEntity.LinkUuid)

	// 检查友链是否存在
	getLinkEntity, errorCode := s.GetOneByUUID(ctx, editFriendEntity.LinkUuid)
	if errorCode != nil {
		return errorCode
	}

	// 更新友链信息
	editFriendEntity.LinkUpdatedAt = gtime.Now()
	_, daoErr := dao.LinkContext.Ctx(ctx).Cache(gdb.CacheOption{
		Duration: -1,
		Name:     fmt.Sprintf(consts.LinkContextByUuidRedisKey, editFriendEntity.LinkUuid),
	}).WherePri(getLinkEntity.LinkUuid).OmitEmpty().Update(&editFriendEntity)
	if daoErr != nil {
		blog.ServiceError(ctx, "Update", "更新友链信息失败，错误：%v", daoErr)
		return berror.ErrorAddData(&berror.ErrDatabaseError, daoErr.Error())
	}

	// 删除相关的全局缓存
	errorCode = deleteGlobalImpactCache(ctx)
	if errorCode != nil {
		return errorCode
	}

	return nil
}

// Delete
//
// 删除友情链接。根据传入的 UUID 删除友链信息。
// 如果 UUID 格式无效或友链不存在，将返回对应的错误码。
// 执行成功后会清理相关的全局缓存以确保数据一致性。
func (s *sFriend) Delete(ctx context.Context, linkUUID string) *berror.ErrorCode {
	blog.ServiceInfo(ctx, "Delete", "删除友情链接 %s", linkUUID)

	// 检查 UUID 是否有效
	parseUUID, uuidErr := uuid.Parse(linkUUID)
	if uuidErr != nil {
		blog.ServiceError(ctx, "Delete", "UUID 格式错误", uuidErr)
		return berror.ErrorAddData(&berror.ErrInvalidParameters, "UUID 格式错误: "+uuidErr.Error())
	}

	// 执行删除操作
	_, daoErr := dao.LinkContext.Ctx(ctx).Cache(gdb.CacheOption{
		Duration: -1,
		Name:     fmt.Sprintf(consts.LinkContextByUuidRedisKey, parseUUID.String()),
	}).WherePri(parseUUID.String()).Delete()
	if daoErr != nil {
		blog.ServiceError(ctx, "Delete", "删除友情链接失败", daoErr)
		return berror.ErrorAddData(&berror.ErrDatabaseError, "删除友情链接失败: "+daoErr.Error())
	}

	// 删除缓存
	_, cacheErr := g.Redis().Del(ctx, fmt.Sprintf(consts.LinkContextByUuidRedisKey, parseUUID.String()))
	if cacheErr != nil {
		blog.ServiceError(ctx, "Delete", "删除友情链接缓存失败", cacheErr)
		return berror.ErrorAddData(&berror.ErrDatabaseError, cacheErr.Error())
	}

	// 删除相关的全局缓存
	errorCode := deleteGlobalImpactCache(ctx)
	if errorCode != nil {
		return errorCode
	}
	return nil
}

// deleteGlobalImpactCache
//
// 删除友链相关的全局缓存；
// 该函数会删除友链分组列表缓存和友链分组列表搜索缓存。
func deleteGlobalImpactCache(ctx context.Context) *berror.ErrorCode {
	_, cacheErr := g.Redis().Del(ctx, consts.SelectCache+consts.LinkContextListRedisKey)
	if cacheErr != nil {
		blog.ServiceError(ctx, "deleteGlobalImpactCache", "删除友链分组列表缓存失败，错误：%v", cacheErr)
		return berror.ErrorAddData(&berror.ErrDatabaseError, cacheErr.Error())
	}
	patternKey := consts.SelectCache + strings.ReplaceAll(consts.LinkContextListSearchRedisKey, "%d:%d:%s", "*")
	keys, cacheErr := g.Redis().Keys(ctx, patternKey)
	g.Log().Debugf(ctx, "deleteGlobalImpactCache 获取友链分组列表搜索缓存，匹配键：%s, 匹配到的键数量：%d", patternKey, len(keys))
	if cacheErr != nil {
		blog.ServiceError(ctx, "deleteGlobalImpactCache", "获取友链分组列表搜索缓存失败，错误：%v", cacheErr)
		return berror.ErrorAddData(&berror.ErrDatabaseError, cacheErr.Error())
	}
	if keys != nil && len(keys) > 0 {
		delCount, cacheErr := g.Redis().Del(ctx, keys...)
		if cacheErr != nil {
			blog.ServiceError(ctx, "deleteGlobalImpactCache", "删除友链分组列表搜索缓存失败，错误：%v", cacheErr)
			return berror.ErrorAddData(&berror.ErrDatabaseError, cacheErr.Error())
		}
		blog.ServiceDebug(ctx, "deleteGlobalImpactCache", "成功删除 %s 缓存数量：%d", patternKey, delCount)
	}
	return nil
}
