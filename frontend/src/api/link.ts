// --------------------------------------------------------------------------------
// Copyright (c) 2016-NOW(至今) 筱锋
// Author: 筱锋「xiao_lfeng」(https://www.x-lf.com)
// --------------------------------------------------------------------------------
// 许可证声明：版权所有 (c) 2016-2026 筱锋。保留所有权利。
// 有关MIT许可证的更多信息，请查看项目根目录下的LICENSE文件或访问：
// https://opensource.org/licenses/MIT
// --------------------------------------------------------------------------------

import { request } from './client'
import type {
  CreateLinkRequest,
  FriendPublicResponse,
  LinkFriend,
  LinkListParams,
  PaginationResponse,
  SnowflakeID,
  UpdateLinkFailRequest,
  UpdateLinkRequest,
  UpdateLinkStatusRequest,
} from './types'

/**
 * 获取公开友链列表。
 * - 后端只返回已通过审核且未失效的友链
 * - 可选 group_id 过滤
 */
export function getPublicLinks(
  groupId?: SnowflakeID,
): Promise<Array<LinkFriend>> {
  return request<FriendPublicResponse>({
    method: 'GET',
    url: '/links',
    params: groupId != null ? { group_id: groupId.toString() } : undefined,
  }).then((res) => res.links)
}

// ---------------------------------------------------------------------------
// 管理端接口（需 admin 角色）
// ---------------------------------------------------------------------------

/** 管理端友链分页列表（支持搜索/状态/分组筛选与排序） */
export function listAdminLinks(
  params: LinkListParams = {},
): Promise<PaginationResponse<LinkFriend>> {
  return request<PaginationResponse<LinkFriend>>({
    method: 'GET',
    url: '/admin/links',
    params: {
      page: params.page ?? 1,
      page_size: params.page_size ?? 10,
      link_name: params.link_name || undefined,
      link_status: params.link_status,
      link_fail: params.link_fail,
      link_group_id: params.link_group_id?.toString(),
      sort_by: params.sort_by,
      sort_order: params.sort_order,
    },
  })
}

/** 管理端友链详情 */
export function getAdminLink(id: SnowflakeID): Promise<LinkFriend> {
  return request<LinkFriend>({
    method: 'GET',
    url: `/admin/links/${id.toString()}`,
  })
}

/** 添加友链 */
export function createLink(req: CreateLinkRequest): Promise<LinkFriend> {
  return request<LinkFriend>({
    method: 'POST',
    url: '/admin/links',
    data: req,
  })
}

/** 更新友链 */
export function updateLink(
  id: SnowflakeID,
  req: UpdateLinkRequest,
): Promise<LinkFriend> {
  return request<LinkFriend>({
    method: 'PUT',
    url: `/admin/links/${id.toString()}`,
    data: req,
  })
}

/** 删除友链 */
export function deleteLink(id: SnowflakeID): Promise<void> {
  return request<void>({
    method: 'DELETE',
    url: `/admin/links/${id.toString()}`,
  })
}

/** 更新友链审核状态（通过/拒绝） */
export function updateLinkStatus(
  id: SnowflakeID,
  req: UpdateLinkStatusRequest,
): Promise<void> {
  return request<void>({
    method: 'PUT',
    url: `/admin/links/${id.toString()}/status`,
    data: req,
  })
}

/** 更新友链失效状态 */
export function updateLinkFail(
  id: SnowflakeID,
  req: UpdateLinkFailRequest,
): Promise<void> {
  return request<void>({
    method: 'PUT',
    url: `/admin/links/${id.toString()}/fail`,
    data: req,
  })
}
