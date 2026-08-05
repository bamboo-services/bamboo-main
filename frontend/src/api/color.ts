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
  ColorListParams,
  ColorPageParams,
  CreateColorRequest,
  LinkColor,
  PaginationResponse,
  SnowflakeID,
  UpdateColorRequest,
} from './types'

/** 友链颜色分页列表 */
export function listColors(
  params: ColorPageParams = {},
): Promise<PaginationResponse<LinkColor>> {
  return request<PaginationResponse<LinkColor>>({
    method: 'GET',
    url: '/admin/colors',
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

/** 友链颜色全量列表（不分页，供选择器使用） */
export function getAllColors(
  params: ColorListParams = {},
): Promise<Array<LinkColor>> {
  return request<Array<LinkColor>>({
    method: 'GET',
    url: '/admin/colors/all',
    params: {
      status: params.status,
      name: params.name,
      only_enabled: params.only_enabled,
      order_by: params.order_by ?? 'sort_order',
      order: params.order ?? 'asc',
    },
  })
}

/** 添加友链颜色 */
export function createColor(req: CreateColorRequest): Promise<LinkColor> {
  return request<LinkColor>({
    method: 'POST',
    url: '/admin/colors',
    data: req,
  })
}

/** 更新友链颜色 */
export function updateColor(
  id: SnowflakeID,
  req: UpdateColorRequest,
): Promise<LinkColor> {
  return request<LinkColor>({
    method: 'PUT',
    url: `/admin/colors/${id.toString()}`,
    data: req,
  })
}

/** 切换友链颜色启用/禁用状态 */
export function updateColorStatus(
  id: SnowflakeID,
  status: boolean,
): Promise<{ status: boolean }> {
  return request<{ status: boolean }>({
    method: 'PATCH',
    url: `/admin/colors/${id.toString()}/status`,
    data: { status },
  })
}

/** 删除友链颜色（force=true 时强制删除并清空关联） */
export function deleteColor(id: SnowflakeID, force = false): Promise<void> {
  return request<void>({
    method: 'DELETE',
    url: `/admin/colors/${id.toString()}`,
    params: { force },
  })
}
