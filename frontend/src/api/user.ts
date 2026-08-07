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
  AdminUserListResponse,
  SnowflakeID,
  UserListParams,
  UserStatusRequest,
} from './types'

// ---------------------------------------------------------------------------
// 管理端用户接口（需 admin 角色）
// ---------------------------------------------------------------------------

/** 管理端用户分页列表（支持关键词搜索/状态筛选与排序） */
export function listAdminUsers(
  params: UserListParams = {},
): Promise<AdminUserListResponse> {
  return request<AdminUserListResponse>({
    method: 'GET',
    url: '/admin/users',
    params: {
      page: params.page ?? 1,
      page_size: params.page_size ?? 10,
      keyword: params.keyword || undefined,
      status: params.status,
      sort_by: params.sort_by,
      sort_order: params.sort_order,
    },
  })
}

/** 更新用户启用/禁用状态（PATCH /api/v1/admin/users/:id/status） */
export function updateUserStatus(
  id: SnowflakeID,
  req: UserStatusRequest,
): Promise<void> {
  return request<void>({
    method: 'PATCH',
    url: `/admin/users/${id.toString()}/status`,
    data: req,
  })
}
