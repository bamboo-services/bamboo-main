// ================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// You can delete these comments if you wish manually maintain this interface file.
// ================================================================================

package service

import (
	"bamboo-main/internal/model/dto"
	"bamboo-main/internal/model/dto/base"
	"bamboo-main/internal/model/entity"
	"context"

	"github.com/XiaoLFeng/bamboo-utils/berror"
)

type (
	IGroup interface {
		// Create
		//
		// 创建友链分组；
		// 创建时候需要提供分组名称、描述和排序值；
		// 创建将会匹配缓存，删除缓存中的分组列表；
		// 如果创建成功，则返回新创建的分组实体；
		// 如果创建失败，则返回错误码。
		Create(ctx context.Context, name string, description string, order int) (*entity.LinkGroup, *berror.ErrorCode)
		// GetOneByUUID
		//
		// 获取单个分组信息;
		// 通过分组 GroupUUID 获取分组信息；
		// 如果获取成功，则返回分组实体；
		// 如果获取失败，则返回错误码。
		GetOneByUUID(ctx context.Context, linkUUID string) (*entity.LinkGroup, *berror.ErrorCode)
		// GetOneByName
		//
		// 获取单个分组信息;
		// 通过分组 GroupName 获取分组信息；
		// 如果获取成功，则返回分组实体；
		// 如果获取失败，则返回错误码。
		GetOneByName(ctx context.Context, name string) (*entity.LinkGroup, *berror.ErrorCode)
		// GetAllList
		//
		// 获取所有友链分组列表；
		// 如果不存在任何友链将会返回空列表（不会产生 err）
		// 如果获取成功，则返回友链分组列表；
		// 如果获取失败，则返回错误码。
		GetAllList(ctx context.Context) ([]*entity.LinkGroup, *berror.ErrorCode)
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
		GetList(ctx context.Context, search string, page int, size int) (*dto.Page[base.LinkGroupDTO], *berror.ErrorCode)
		// Update
		//
		// 更新友链分组；
		// 通过分组 UUID 更新友链分组信息；
		// 如果更新成功，则返回 nil；如果更新失败，则返回错误码。
		//
		// 注意：更新会删除有关内容的所有缓存（如 consts.LinkGroupListRedisKey 等）
		// 若修改了 name 则会删除 consts.LinkGroupByNameRedisKey 缓存
		Update(ctx context.Context, uuid string, name string, description string, order int) *berror.ErrorCode
		// Delete
		//
		// 删除友链分组；
		// 通过分组 UUID 删除友链分组信息；
		// 如果删除成功，则返回 nil；如果删除失败，则返回错误码。
		//
		// 注意：删除会删除有关内容的所有缓存（如 consts.LinkGroupListRedisKey 等）
		Delete(ctx context.Context, uuid string) *berror.ErrorCode
	}
)

var (
	localGroup IGroup
)

func Group() IGroup {
	if localGroup == nil {
		panic("implement not found for interface IGroup, forgot register?")
	}
	return localGroup
}

func RegisterGroup(i IGroup) {
	localGroup = i
}
