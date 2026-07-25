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
  PaginationResponse,
  SponsorChannel,
  SponsorRecord,
} from './types'

/** 获取公开赞助渠道列表（按 sort_order 升序） */
export function getPublicChannels(): Promise<Array<SponsorChannel>> {
  return request<{ channels?: Array<SponsorChannel> } | Array<SponsorChannel>>({
    method: 'GET',
    url: '/sponsors/channels',
  }).then((res) => {
    // 兼容后端直接返回数组或包裹 { channels } 两种形态
    if (Array.isArray(res)) return res
    return res.channels ?? []
  })
}

export interface GetRecordsParams {
  page?: number
  pageSize?: number
  channelId?: number
  orderBy?: 'amount' | 'sponsor_at' | 'sort_order'
  order?: 'asc' | 'desc'
}

/** 获取公开赞助记录分页列表 */
export function getPublicRecords(
  params: GetRecordsParams = {},
): Promise<PaginationResponse<SponsorRecord>> {
  return request<PaginationResponse<SponsorRecord>>({
    method: 'GET',
    url: '/sponsors/records',
    params: {
      page: params.page ?? 1,
      page_size: params.pageSize ?? 20,
      channel_id: params.channelId,
      order_by: params.orderBy,
      order: params.order,
    },
  })
}
