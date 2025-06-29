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

package group

import (
	"bamboo-main/internal/consts"
	"bamboo-main/internal/dao"
	"bamboo-main/internal/model/do"
	"bamboo-main/internal/model/dto"
	"bamboo-main/internal/model/dto/base"
	"bamboo-main/internal/model/entity"
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

// Create
//
// 创建友链分组；
// 创建时候需要提供分组名称、描述和排序值；
// 创建将会匹配缓存，删除缓存中的分组列表；
// 如果创建成功，则返回新创建的分组实体；
// 如果创建失败，则返回错误码。
func (s *sGroup) Create(ctx context.Context, name, description string, order int) (*entity.LinkGroup, *berror.ErrorCode) {
	blog.ServiceInfo(ctx, "Create", "创建友链分组 %s", name)

	var newEntity = &entity.LinkGroup{
		GroupUuid:      uuid.New().String(),
		GroupName:      name,
		GroupDesc:      description,
		GroupOrder:     order,
		GroupCreatedAt: gtime.Now(),
		GroupUpdatedAt: gtime.Now(),
	}
	_, daoErr := dao.LinkGroup.Ctx(ctx).Cache(gdb.CacheOption{
		Duration: -1,
		Name:     fmt.Sprintf(consts.LinkGroupByUuidRedisKey, newEntity.GroupUuid),
	}).OmitEmpty().Insert(&newEntity)
	if daoErr != nil {
		blog.ServiceError(ctx, "Create", "创建友链分组失败，错误：%v", daoErr)
		return nil, berror.ErrorAddData(&berror.ErrDatabaseError, daoErr.Error())
	}
	// 删除其余缓存
	_, cacheErr := g.Redis().Del(ctx, fmt.Sprintf(consts.LinkGroupByUuidRedisKey, newEntity.GroupUuid))
	if cacheErr != nil {
		blog.ServiceError(ctx, "Create", "删除友链分组缓存失败，错误：%v", cacheErr)
		return nil, berror.ErrorAddData(&berror.ErrDatabaseError, cacheErr.Error())
	}
	_, cacheErr = g.Redis().Del(ctx, consts.LinkGroupByNameRedisKey, newEntity.GroupName)
	if cacheErr != nil {
		blog.ServiceError(ctx, "Create", "删除友链分组名称缓存失败，错误：%v", cacheErr)
		return nil, berror.ErrorAddData(&berror.ErrDatabaseError, cacheErr.Error())
	}
	// 删除全局影响的缓存
	errorCode := deleteGlobalImpactCache(ctx)
	if errorCode != nil {
		return nil, errorCode
	}
	return newEntity, nil
}

// GetOneByUUID
//
// 获取单个分组信息;
// 通过分组 GroupUUID 获取分组信息；
// 如果获取成功，则返回分组实体；
// 如果获取失败，则返回错误码。
func (s *sGroup) GetOneByUUID(ctx context.Context, uuid string) (*entity.LinkGroup, *berror.ErrorCode) {
	blog.ServiceInfo(ctx, "GetOneByUUID", "获取友链分组 %s", uuid)
	var groupEntity *entity.LinkGroup
	daoErr := dao.LinkGroup.Ctx(ctx).Cache(gdb.CacheOption{
		Duration: 7 * 24 * time.Hour,
		Name:     fmt.Sprintf(consts.LinkGroupByUuidRedisKey, uuid),
	}).Where(&do.LinkGroup{GroupUuid: uuid}).OmitEmptyWhere().Scan(&groupEntity)
	if daoErr != nil {
		blog.ServiceError(ctx, "GetOneByUUID", "获取友链分组失败，错误：%v", daoErr)
		return nil, berror.ErrorAddData(&berror.ErrDatabaseError, daoErr.Error())
	}
	if groupEntity == nil {
		blog.ServiceNotice(ctx, "GetOneByUUID", "友链分组不存在，UUID：%s", uuid)
		return nil, berror.ErrorAddData(&berror.ErrNotFound, "友链分组不存在")
	}
	return groupEntity, nil
}

// GetOneByName
//
// 获取单个分组信息;
// 通过分组 GroupName 获取分组信息；
// 如果获取成功，则返回分组实体；
// 如果获取失败，则返回错误码。
func (s *sGroup) GetOneByName(ctx context.Context, name string) (*entity.LinkGroup, *berror.ErrorCode) {
	blog.ServiceInfo(ctx, "GetOneByName", "获取友链分组 %s", name)
	var groupEntity *entity.LinkGroup
	daoErr := dao.LinkGroup.Ctx(ctx).Cache(gdb.CacheOption{
		Duration: 7 * 24 * time.Hour,
		Name:     fmt.Sprintf(consts.LinkGroupByNameRedisKey, name),
	}).Where(&do.LinkGroup{GroupName: name}).OmitEmptyWhere().Scan(&groupEntity)
	if daoErr != nil {
		blog.ServiceError(ctx, "GetOneByName", "获取友链分组失败，错误：%v", daoErr)
		return nil, berror.ErrorAddData(&berror.ErrDatabaseError, daoErr.Error())
	}
	if groupEntity == nil {
		blog.ServiceNotice(ctx, "GetOneByName", "友链分组不存在，名称：%s", name)
		return nil, berror.ErrorAddData(&berror.ErrNotFound, "友链分组不存在")
	}
	return groupEntity, nil
}

// GetAllList
//
// 获取所有友链分组列表；
// 如果不存在任何友链将会返回空列表（不会产生 err）
// 如果获取成功，则返回友链分组列表；
// 如果获取失败，则返回错误码。
func (s *sGroup) GetAllList(ctx context.Context) ([]*entity.LinkGroup, *berror.ErrorCode) {
	blog.ServiceInfo(ctx, "GetAllList", "获取所有友链分组列表")
	var groupList []*entity.LinkGroup
	daoErr := dao.LinkGroup.Ctx(ctx).Cache(gdb.CacheOption{
		Duration: 7 * 24 * time.Hour,
		Name:     consts.LinkGroupListRedisKey,
	}).Scan(&groupList)
	if daoErr != nil {
		blog.ServiceError(ctx, "GetAllList", "获取友链分组列表失败，错误：%v", daoErr)
		return nil, berror.ErrorAddData(&berror.ErrDatabaseError, daoErr.Error())
	}
	return groupList, nil
}

// GetList
//
// 获取友链分组的分页列表数据；
// 根据搜索条件检索分组，支持按分组名称和描述模糊搜索；
// search 参数支持空字符串，表示不进行搜索过滤；
// page 参数指定当前页码，从 1 开始；
// size 参数指定每页显示记录数；
// 方法返回带有分页信息的友链分组 DTO 列表；
// 如果不存在任何符合条件的分组，将返回空列表（不会产生错误）；
// 如果查询成功，返回分页数据；如果查询失败，返回错误码。
func (s *sGroup) GetList(ctx context.Context, search string, page, size int) (*dto.Page[base.LinkGroupDTO], *berror.ErrorCode) {
	blog.ServiceInfo(ctx, "GetList", "获取友链分组列表，搜索条件：%s", utility.GetOrDefaultString(&search, "<NoSearchData>"))
	var groupList []*entity.LinkGroup
	var whereBuilder *gdb.WhereBuilder
	if search != "" {
		whereBuilder = dao.LinkGroup.Ctx(ctx).Builder().
			WhereOrLike(dao.LinkGroup.Columns().GroupName, "%"+search+"%").
			WhereOrLike(dao.LinkGroup.Columns().GroupDesc, "%"+search+"%")
	}
	daoErr := dao.LinkGroup.Ctx(ctx).Cache(gdb.CacheOption{
		Duration: 7 * 24 * time.Hour,
		Name:     fmt.Sprintf(consts.LinkGroupListHasSearchRedisKey, page, size, utility.GetOrDefaultString(&search, "<NoSearchData>")),
	}).Where(whereBuilder).OrderDesc(dao.LinkGroup.Columns().GroupOrder).Scan(&groupList)
	if daoErr != nil {
		blog.ServiceError(ctx, "GetList", "获取友链分组列表失败，错误：%v", daoErr)
		return nil, berror.ErrorAddData(&berror.ErrDatabaseError, daoErr.Error())
	}
	return utility.MakePageToTarget(groupList, page, size, base.LinkGroupDTO{}), nil
}

// Update
//
// 更新友链分组；
// 通过分组 UUID 更新友链分组信息；
// 如果更新成功，则返回 nil；如果更新失败，则返回错误码。
//
// 注意：更新会删除有关内容的所有缓存（如 consts.LinkGroupListRedisKey 等）
// 若修改了 name 则会删除 consts.LinkGroupByNameRedisKey 缓存
func (s *sGroup) Update(ctx context.Context, uuid, name, description string, order int) *berror.ErrorCode {
	blog.ServiceInfo(ctx, "Update", "更新友链分组 %s", uuid)

	// 查询原来的数据
	originalEntity, errorCode := s.GetOneByUUID(ctx, uuid)
	if errorCode != nil {
		return errorCode
	}

	// 更新数据
	var newEntity = &entity.LinkGroup{
		GroupUuid:      originalEntity.GroupUuid,
		GroupName:      utility.GetOrDefaultString(&name, originalEntity.GroupName),
		GroupDesc:      utility.GetOrDefaultString(&description, originalEntity.GroupDesc),
		GroupOrder:     utility.GetOrDefault(&order, originalEntity.GroupOrder),
		GroupCreatedAt: originalEntity.GroupCreatedAt,
		GroupUpdatedAt: gtime.Now(),
	}
	_, daoErr := dao.LinkGroup.Ctx(ctx).Cache(gdb.CacheOption{
		Duration: -1,
		Name:     fmt.Sprintf(consts.LinkGroupByUuidRedisKey, uuid),
	}).OmitEmpty().WherePri(originalEntity.GroupUuid).OmitEmpty().Update(&newEntity)
	if daoErr != nil {
		blog.ServiceError(ctx, "Update", "更新友链分组失败，错误：%v", daoErr)
		return berror.ErrorAddData(&berror.ErrDatabaseError, daoErr.Error())
	}
	// 删除缓存
	if originalEntity.GroupName != name {
		_, daoErr = g.Redis().Del(ctx, fmt.Sprintf(consts.LinkGroupByNameRedisKey, originalEntity.GroupName))
		if daoErr != nil {
			blog.ServiceError(ctx, "Update", "删除友链分组名称缓存失败，错误：%v", daoErr)
			return berror.ErrorAddData(&berror.ErrDatabaseError, daoErr.Error())
		}
	}
	errorCode = deleteGlobalImpactCache(ctx)
	if errorCode != nil {
		return errorCode
	}
	return nil
}

// Delete
//
// 删除友链分组；
// 通过分组 UUID 删除友链分组信息；
// 如果删除成功，则返回 nil；如果删除失败，则返回错误码。
//
// 注意：删除会删除有关内容的所有缓存（如 consts.LinkGroupListRedisKey 等）
func (s *sGroup) Delete(ctx context.Context, uuid string) *berror.ErrorCode {
	blog.ServiceInfo(ctx, "Delete", "删除友链分组 %s", uuid)

	// 查询原来的数据
	originalEntity, errorCode := s.GetOneByUUID(ctx, uuid)
	if errorCode != nil {
		return errorCode
	}

	// 删除数据
	_, daoErr := dao.LinkGroup.Ctx(ctx).Cache(gdb.CacheOption{
		Duration: -1,
		Name:     fmt.Sprintf(consts.LinkGroupByUuidRedisKey, uuid),
	}).WherePri(originalEntity.GroupUuid).Delete()
	if daoErr != nil {
		blog.ServiceError(ctx, "Delete", "删除友链分组失败，错误：%v", daoErr)
		return berror.ErrorAddData(&berror.ErrDatabaseError, daoErr.Error())
	}
	_, daoErr = g.Redis().Del(ctx, fmt.Sprintf(consts.LinkGroupByNameRedisKey, originalEntity.GroupName))
	if daoErr != nil {
		blog.ServiceError(ctx, "Delete", "删除友链分组名称缓存失败，错误：%v", daoErr)
		return berror.ErrorAddData(&berror.ErrDatabaseError, daoErr.Error())
	}

	// 删除缓存
	errorCode = deleteGlobalImpactCache(ctx)
	if errorCode != nil {
		return errorCode
	}
	return nil
}

// deleteGlobalImpactCache
//
// 删除全局影响的缓存；
// 主要是删除友链分组列表缓存和友链分组列表搜索缓存；
// 如果删除成功，则返回 nil；
// 如果删除失败，则返回错误码。
func deleteGlobalImpactCache(ctx context.Context) *berror.ErrorCode {
	_, cacheErr := g.Redis().Del(ctx, consts.SelectCache+consts.LinkGroupListRedisKey)
	if cacheErr != nil {
		blog.ServiceError(ctx, "deleteGlobalImpactCache", "删除友链分组列表缓存失败，错误：%v", cacheErr)
		return berror.ErrorAddData(&berror.ErrDatabaseError, cacheErr.Error())
	}
	patternKey := consts.SelectCache + strings.ReplaceAll(consts.LinkGroupListHasSearchRedisKey, "%d:%d:%s", "*")
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
