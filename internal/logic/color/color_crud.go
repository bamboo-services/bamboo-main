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

package color

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
// 创建友链颜色；
// 创建时候需要提供颜色名称、颜色值和描述；
// 创建将会匹配缓存，删除缓存中的颜色列表；
// 如果创建成功，则返回新创建的颜色实体；
// 如果创建失败，则返回错误码。
func (s *sColor) Create(ctx context.Context, name, value, description string) (*entity.LinkColor, *berror.ErrorCode) {
	blog.ServiceInfo(ctx, "Create", "创建友链颜色 %s", name)

	var newEntity = &entity.LinkColor{
		ColorUuid:      uuid.New().String(),
		ColorName:      name,
		ColorValue:     value,
		ColorDesc:      description,
		ColorCreatedAt: gtime.Now(),
		ColorUpdatedAt: gtime.Now(),
	}
	_, daoErr := dao.LinkColor.Ctx(ctx).Cache(gdb.CacheOption{
		Duration: -1,
		Name:     fmt.Sprintf(consts.LinkColorByUuidRedisKey, newEntity.ColorUuid),
	}).OmitEmpty().Insert(&newEntity)
	if daoErr != nil {
		blog.ServiceError(ctx, "Create", "创建友链颜色失败，错误：%v", daoErr)
		return nil, berror.ErrorAddData(&berror.ErrDatabaseError, daoErr.Error())
	}
	// 删除其余缓存
	_, cacheErr := g.Redis().Del(ctx, fmt.Sprintf(consts.LinkColorByUuidRedisKey, newEntity.ColorUuid))
	if cacheErr != nil {
		blog.ServiceError(ctx, "Create", "删除友链颜色缓存失败，错误：%v", cacheErr)
		return nil, berror.ErrorAddData(&berror.ErrDatabaseError, cacheErr.Error())
	}
	_, cacheErr = g.Redis().Del(ctx, consts.LinkColorByNameRedisKey, newEntity.ColorName)
	if cacheErr != nil {
		blog.ServiceError(ctx, "Create", "删除友链颜色名称缓存失败，错误：%v", cacheErr)
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
// 获取单个颜色信息;
// 通过颜色 ColorUUID 获取颜色信息；
// 如果获取成功，则返回颜色实体；
// 如果获取失败，则返回错误码。
func (s *sColor) GetOneByUUID(ctx context.Context, linkUUID string) (*entity.LinkColor, *berror.ErrorCode) {
	blog.ServiceInfo(ctx, "GetOneByUUID", "获取友链颜色 %s", linkUUID)
	// 检查输入是否是 UUID
	parseUUID, uuidErr := uuid.Parse(linkUUID)
	if uuidErr != nil {
		blog.ServiceError(ctx, "GetOneByUUID", "无效的 UUID 格式，错误：%v", uuidErr)
		return nil, berror.ErrorAddData(&berror.ErrInvalidParameters, "无效的 UUID 格式")
	}
	var colorEntity *entity.LinkColor
	daoErr := dao.LinkColor.Ctx(ctx).Cache(gdb.CacheOption{
		Duration: 7 * 24 * time.Hour,
		Name:     fmt.Sprintf(consts.LinkColorByUuidRedisKey, linkUUID),
	}).Where(&do.LinkColor{ColorUuid: parseUUID.String()}).OmitEmpty().Scan(&colorEntity)
	if daoErr != nil {
		blog.ServiceError(ctx, "GetOneByUUID", "获取友链颜色失败，错误：%v", daoErr)
		return nil, berror.ErrorAddData(&berror.ErrDatabaseError, daoErr.Error())
	}
	if colorEntity == nil {
		blog.ServiceNotice(ctx, "GetOneByUUID", "友链颜色不存在，UUID：%s", linkUUID)
		return nil, berror.ErrorAddData(&berror.ErrNotFound, "友链颜色不存在")
	}
	return colorEntity, nil
}

// GetOneByName
//
// 获取单个颜色信息;
// 通过颜色 ColorName 获取颜色信息；
// 如果获取成功，则返回颜色实体；
// 如果获取失败，则返回错误码。
func (s *sColor) GetOneByName(ctx context.Context, name string) (*entity.LinkColor, *berror.ErrorCode) {
	blog.ServiceInfo(ctx, "GetOneByName", "获取友链颜色 %s", name)
	var colorEntity *entity.LinkColor
	daoErr := dao.LinkColor.Ctx(ctx).Cache(gdb.CacheOption{
		Duration: 7 * 24 * time.Hour,
		Name:     fmt.Sprintf(consts.LinkColorByNameRedisKey, name),
	}).Where(&do.LinkColor{ColorName: name}).OmitEmptyWhere().Scan(&colorEntity)
	if daoErr != nil {
		blog.ServiceError(ctx, "GetOneByName", "获取友链颜色失败，错误：%v", daoErr)
		return nil, berror.ErrorAddData(&berror.ErrDatabaseError, daoErr.Error())
	}
	if colorEntity == nil {
		blog.ServiceNotice(ctx, "GetOneByName", "友链颜色不存在，名称：%s", name)
		return nil, berror.ErrorAddData(&berror.ErrNotFound, "友链颜色不存在")
	}
	return colorEntity, nil
}

// GetAllList
//
// 获取所有友链颜色列表；
// 如果不存在任何友链颜色将会返回空列表（不会产生 err）
// 如果获取成功，则返回友链颜色列表；
// 如果获取失败，则返回错误码。
func (s *sColor) GetAllList(ctx context.Context) ([]*entity.LinkColor, *berror.ErrorCode) {
	blog.ServiceInfo(ctx, "GetAllList", "获取所有友链颜色列表")
	var colorList []*entity.LinkColor
	daoErr := dao.LinkColor.Ctx(ctx).Cache(gdb.CacheOption{
		Duration: 7 * 24 * time.Hour,
		Name:     consts.LinkColorListRedisKey,
	}).Scan(&colorList)
	if daoErr != nil {
		blog.ServiceError(ctx, "GetAllList", "获取友链颜色列表失败，错误：%v", daoErr)
		return nil, berror.ErrorAddData(&berror.ErrDatabaseError, daoErr.Error())
	}
	return colorList, nil
}

// GetList
//
// 获取友链颜色的分页列表数据；
// 根据搜索条件检索颜色，支持按颜色名称和描述模糊搜索；
// search 参数支持空字符串，表示不进行搜索过滤；
// page 参数指定当前页码，从 1 开始；
// size 参数指定每页显示记录数；
// 方法返回带有分页信息的友链颜色 DTO 列表；
// 如果不存在任何符合条件的颜色，将返回空列表（不会产生错误）；
// 如果查询成功，返回分页数据；如果查询失败，返回错误码。
func (s *sColor) GetList(ctx context.Context, search string, page, size int) (*dto.Page[base.LinkColorDTO], *berror.ErrorCode) {
	blog.ServiceInfo(ctx, "GetList", "获取友链颜色列表，搜索条件：%s", utility.GetOrDefaultString(&search, "<NoSearchData>"))
	var colorList []*entity.LinkColor
	var whereBuilder *gdb.WhereBuilder
	if search != "" {
		whereBuilder = dao.LinkColor.Ctx(ctx).Builder().
			WhereOrLike(dao.LinkColor.Columns().ColorName, "%"+search+"%").
			WhereOrLike(dao.LinkColor.Columns().ColorDesc, "%"+search+"%")
	}
	daoErr := dao.LinkColor.Ctx(ctx).Cache(gdb.CacheOption{
		Duration: 7 * 24 * time.Hour,
		Name:     fmt.Sprintf(consts.LinkColorListHasSearchRedisKey, page, size, utility.GetOrDefaultString(&search, "<NoSearchData>")),
	}).Where(whereBuilder).Scan(&colorList)
	if daoErr != nil {
		blog.ServiceError(ctx, "GetList", "获取友链颜色列表失败，错误：%v", daoErr)
		return nil, berror.ErrorAddData(&berror.ErrDatabaseError, daoErr.Error())
	}
	return utility.MakePageToTarget(colorList, page, size, base.LinkColorDTO{}), nil
}

// Update
//
// 更新友链颜色；
// 通过颜色 UUID 更新友链颜色信息；
// 如果更新成功，则返回 nil；如果更新失败，则返回错误码。
//
// 注意：更新会删除有关内容的所有缓存（如 consts.LinkColorListRedisKey 等）
// 若修改了 name 则会删除 consts.LinkColorByNameRedisKey 缓存
func (s *sColor) Update(ctx context.Context, uuid, name, value, description string) *berror.ErrorCode {
	blog.ServiceInfo(ctx, "Update", "更新友链颜色 %s", uuid)

	// 查询原来的数据
	originalEntity, errorCode := s.GetOneByUUID(ctx, uuid)
	if errorCode != nil {
		return errorCode
	}

	// 更新数据
	var newEntity = &entity.LinkColor{
		ColorUuid:      originalEntity.ColorUuid,
		ColorName:      utility.GetOrDefaultString(&name, originalEntity.ColorName),
		ColorValue:     utility.GetOrDefaultString(&value, originalEntity.ColorValue),
		ColorDesc:      utility.GetOrDefaultString(&description, originalEntity.ColorDesc),
		ColorCreatedAt: originalEntity.ColorCreatedAt,
		ColorUpdatedAt: gtime.Now(),
	}
	_, daoErr := dao.LinkColor.Ctx(ctx).Cache(gdb.CacheOption{
		Duration: -1,
		Name:     fmt.Sprintf(consts.LinkColorByUuidRedisKey, uuid),
	}).OmitEmpty().WherePri(originalEntity.ColorUuid).OmitEmpty().Update(&newEntity)
	if daoErr != nil {
		blog.ServiceError(ctx, "Update", "更新友链颜色失败，错误：%v", daoErr)
		return berror.ErrorAddData(&berror.ErrDatabaseError, daoErr.Error())
	}
	// 删除缓存
	if originalEntity.ColorName != name {
		_, daoErr = g.Redis().Del(ctx, fmt.Sprintf(consts.LinkColorByNameRedisKey, originalEntity.ColorName))
		if daoErr != nil {
			blog.ServiceError(ctx, "Update", "删除友链颜色名称缓存失败，错误：%v", daoErr)
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
// 删除友链颜色；
// 通过颜色 UUID 删除友链颜色信息；
// 如果删除成功，则返回 nil；如果删除失败，则返回错误码。
//
// 注意：删除会删除有关内容的所有缓存（如 consts.LinkColorListRedisKey 等）
func (s *sColor) Delete(ctx context.Context, uuid string) *berror.ErrorCode {
	blog.ServiceInfo(ctx, "Delete", "删除友链颜色 %s", uuid)

	// 查询原来的数据
	originalEntity, errorCode := s.GetOneByUUID(ctx, uuid)
	if errorCode != nil {
		return errorCode
	}

	// 删除数据
	_, daoErr := dao.LinkColor.Ctx(ctx).Cache(gdb.CacheOption{
		Duration: -1,
		Name:     fmt.Sprintf(consts.LinkColorByUuidRedisKey, uuid),
	}).WherePri(originalEntity.ColorUuid).Delete()
	if daoErr != nil {
		blog.ServiceError(ctx, "Delete", "删除友链颜色失败，错误：%v", daoErr)
		return berror.ErrorAddData(&berror.ErrDatabaseError, daoErr.Error())
	}
	_, daoErr = g.Redis().Del(ctx, fmt.Sprintf(consts.LinkColorByNameRedisKey, originalEntity.ColorName))
	if daoErr != nil {
		blog.ServiceError(ctx, "Delete", "删除友链颜色名称缓存失败，错误：%v", daoErr)
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
// 主要是删除友链颜色列表缓存和友链颜色列表搜索缓存；
// 如果删除成功，则返回 nil；
// 如果删除失败，则返回错误码。
func deleteGlobalImpactCache(ctx context.Context) *berror.ErrorCode {
	_, cacheErr := g.Redis().Del(ctx, consts.SelectCache+consts.LinkColorListRedisKey)
	if cacheErr != nil {
		blog.ServiceError(ctx, "deleteGlobalImpactCache", "删除友链颜色列表缓存失败，错误：%v", cacheErr)
		return berror.ErrorAddData(&berror.ErrDatabaseError, cacheErr.Error())
	}
	patternKey := consts.SelectCache + strings.ReplaceAll(consts.LinkColorListHasSearchRedisKey, "%d:%d:%s", "*")
	keys, cacheErr := g.Redis().Keys(ctx, patternKey)
	g.Log().Debugf(ctx, "deleteGlobalImpactCache 获取友链颜色列表搜索缓存，匹配键：%s, 匹配到的键数量：%d", patternKey, len(keys))
	if cacheErr != nil {
		blog.ServiceError(ctx, "deleteGlobalImpactCache", "获取友链颜色列表搜索缓存失败，错误：%v", cacheErr)
		return berror.ErrorAddData(&berror.ErrDatabaseError, cacheErr.Error())
	}
	if keys != nil && len(keys) > 0 {
		delCount, cacheErr := g.Redis().Del(ctx, keys...)
		if cacheErr != nil {
			blog.ServiceError(ctx, "deleteGlobalImpactCache", "删除友链颜色列表搜索缓存失败，错误：%v", cacheErr)
			return berror.ErrorAddData(&berror.ErrDatabaseError, cacheErr.Error())
		}
		blog.ServiceDebug(ctx, "deleteGlobalImpactCache", "成功删除 %s 缓存数量：%d", patternKey, delCount)
	}
	return nil
}
