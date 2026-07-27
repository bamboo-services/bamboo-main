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
  CreateGroupRequest,
  GroupListParams,
  GroupPageParams,
  LinkGroup,
  PaginationResponse,
  SnowflakeID,
  UpdateGroupRequest,
} from './types'

/** 友链分组分页列表 */
export function listGroups(
  params: GroupPageParams = {},
): Promise<PaginationResponse<LinkGroup>> {
  return request<PaginationResponse<LinkGroup>>({
    method: 'GET',
    url: '/admin/groups',
    params: {
      page: params.page ?? 1,
      page_size: params.page_size ?? 10,
      status: params.status,
      name: params.name,
      order_by: params.order_by,
      order: params.order,
    },
  })
}

/** 友链分组全量列表（不分页，供选择器使用） */
export function getAllGroups(
  params: GroupListParams = {},
): Promise<Array<LinkGroup>> {
  return request<Array<LinkGroup>>({
    method: 'GET',
    url: '/admin/groups/all',
    params: {
      status: params.status,
      name: params.name,
      with_links: params.with_links,
      only_enabled: params.only_enabled,
      order_by: params.order_by ?? 'sort_order',
      order: params.order ?? 'asc',
    },
  })
}

/** 添加友链分组 */
export function createGroup(req: CreateGroupRequest): Promise<LinkGroup> {
  return request<LinkGroup>({
    method: 'POST',
    url: '/admin/groups',
    data: req,
  })
}

/** 更新友链分组 */
export function updateGroup(
  id: SnowflakeID,
  req: UpdateGroupRequest,
): Promise<LinkGroup> {
  return request<LinkGroup>({
    method: 'PUT',
    url: `/admin/groups/${id.toString()}`,
    data: req,
  })
}

/** 切换友链分组启用/禁用状态 */
export function updateGroupStatus(
  id: SnowflakeID,
  status: boolean,
): Promise<{ status: boolean }> {
  return request<{ status: boolean }>({
    method: 'PATCH',
    url: `/admin/groups/${id.toString()}/status`,
    data: { status },
  })
}

/** 删除友链分组（force=true 时强制删除并清空关联） */
export function deleteGroup(id: SnowflakeID, force = false): Promise<void> {
  return request<void>({
    method: 'DELETE',
    url: `/admin/groups/${id.toString()}`,
    params: { force },
  })
}
