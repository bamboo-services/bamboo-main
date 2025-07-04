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
	IFriend interface {
		// AddFriend
		//
		// 添加一个新的友链记录至数据库，同时检查是否存在重复的友链，若重复则返回对应错误。
		// 该函数会自动生成友链的唯一标识符（UUID），并设置友链的状态为启用（1）和创建、更新时间为当前时间。
		// 如果添加成功，则会删除相关的全局缓存，以确保数据的一致性和最新性。
		AddFriend(ctx context.Context, friend base.LinkFriendDTO) *berror.ErrorCode
		// GetOneByUUID
		//
		// 根据传入的 UUID 查询并返回友链信息，不存在或发生错误时返回对应的错误码。
		// 该函数会检查传入的 UUID 是否有效，并使用缓存机制来提高查询效率。
		// 如果查询成功，则返回友链实体；如果查询失败或未找到友链信息，则返回错误码。
		GetOneByUUID(ctx context.Context, linkUUID string) (*entity.LinkContext, *berror.ErrorCode)
		// GetOneByURL
		//
		// 根据传入的 URL 查询并返回友链信息，如果未找到或发生错误，则返回对应的错误码。
		// 该函数会解析传入的 URL，提取出域名部分，并使用缓存机制来提高查询效率。
		GetOneByURL(ctx context.Context, linkURL string) (*entity.LinkContext, *berror.ErrorCode)
		// Update
		//
		// 更新友链信息。如果友链不存在或更新操作失败，返回对应的错误码。更新成功后会删除相关的全局缓存。
		// 该函数会检查友链是否存在，并更新其信息，包括名称、URL、头像、RSS、描述、邮箱、分组、颜色、排序和审核备注等字段。
		// 如果更新成功，则会删除相关的全局缓存，以确保数据的一致性和最新性。
		// 如果更新失败，则返回错误码。
		Update(ctx context.Context, editFriendEntity base.LinkFriendDTO) *berror.ErrorCode
		// Delete
		//
		// 删除友情链接。根据传入的 UUID 删除友链信息。
		// 如果 UUID 格式无效或友链不存在，将返回对应的错误码。
		// 执行成功后会清理相关的全局缓存以确保数据一致性。
		Delete(ctx context.Context, linkUUID string) *berror.ErrorCode
		// GetPage
		//
		// 获取友链分页数据。根据传入的搜索内容、页码和每页数量进行查询。
		// 如果 isAll 为 false，则只查询未失败的友链；如果 isAll 为 true，则查询所有启用的友链。
		// 返回一个包含友链分页数据的 DTO 对象，如果查询失败则返回对应的错误码。
		GetPage(ctx context.Context, search string, page int, size int, isAll bool) (*dto.Page[base.LinkFriendDTO], *berror.ErrorCode)
		// DeleteGlobalImpactCache
		//
		// 删除友链相关的全局缓存；
		// 该函数会删除友链分组列表缓存和友链分组列表搜索缓存。
		DeleteGlobalImpactCache(ctx context.Context) *berror.ErrorCode
		// UpdateStatus
		//
		// 更新友链状态，根据传入的 UUID 和状态值进行操作。
		// 如果 UUID 格式无效或状态值非法，将返回对应的错误码。
		// 执行成功后会清理相关的全局缓存以确保数据一致性。
		UpdateStatus(ctx context.Context, linkUUID string, status int) *berror.ErrorCode
		// UpdateFailStatus
		//
		// 更新友链失败状态和原因，根据传入的 UUID、失败状态和原因进行操作。
		// 如果 UUID 格式无效，将返回对应的错误码。
		// 执行成功后会清理相关的全局缓存以确保数据一致性。
		UpdateFailStatus(ctx context.Context, linkUUID string, isFail bool, reason string) *berror.ErrorCode
	}
)

var (
	localFriend IFriend
)

func Friend() IFriend {
	if localFriend == nil {
		panic("implement not found for interface IFriend, forgot register?")
	}
	return localFriend
}

func RegisterFriend(i IFriend) {
	localFriend = i
}
