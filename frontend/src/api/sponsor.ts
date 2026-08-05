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
  AdminRecordListResponse,
  ChannelPageParams,
  CreateChannelRequest,
  CreateRecordRequest,
  PaginationResponse,
  RecordPageParams,
  SnowflakeID,
  SponsorApplyRequest,
  SponsorChannel,
  SponsorChannelAdmin,
  SponsorRecord,
  SponsorRecordAdmin,
  SponsorStatusRequest,
  SponsorUserParams,
  SponsorUserUpdateRequest,
  UpdateChannelRequest,
  UpdateRecordRequest,
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
  channelId?: SnowflakeID
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
      channel_id: params.channelId?.toString(),
      order_by: params.orderBy,
      order: params.order,
    },
  })
}

// ---------------------------------------------------------------------------
// 管理端接口（需 admin 角色）
// ---------------------------------------------------------------------------

/** 管理端赞助渠道分页列表 */
export function listAdminChannels(
  params: ChannelPageParams = {},
): Promise<PaginationResponse<SponsorChannelAdmin>> {
  return request<PaginationResponse<SponsorChannelAdmin>>({
    method: 'GET',
    url: '/admin/sponsors/channels',
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

/** 管理端赞助渠道全量列表（不分页，供记录表单选择器使用） */
export function getAllChannels(): Promise<Array<SponsorChannel>> {
  return request<Array<SponsorChannel>>({
    method: 'GET',
    url: '/admin/sponsors/channels/all',
    params: { order_by: 'sort_order', order: 'asc' },
  })
}

/** 添加赞助渠道 */
export function createChannel(
  req: CreateChannelRequest,
): Promise<SponsorChannelAdmin> {
  return request<SponsorChannelAdmin>({
    method: 'POST',
    url: '/admin/sponsors/channels',
    data: req,
  })
}

/** 更新赞助渠道 */
export function updateChannel(
  id: SnowflakeID,
  req: UpdateChannelRequest,
): Promise<SponsorChannelAdmin> {
  return request<SponsorChannelAdmin>({
    method: 'PUT',
    url: `/admin/sponsors/channels/${id.toString()}`,
    data: req,
  })
}

/** 切换赞助渠道启用/禁用状态 */
export function updateChannelStatus(
  id: SnowflakeID,
  status: boolean,
): Promise<{ status: boolean }> {
  return request<{ status: boolean }>({
    method: 'PATCH',
    url: `/admin/sponsors/channels/${id.toString()}/status`,
    data: { status },
  })
}

/** 删除赞助渠道 */
export function deleteChannel(id: SnowflakeID): Promise<void> {
  return request<void>({
    method: 'DELETE',
    url: `/admin/sponsors/channels/${id.toString()}`,
  })
}

/** 管理端赞助记录分页列表 */
export function listAdminRecords(
  params: RecordPageParams = {},
): Promise<AdminRecordListResponse> {
  return request<AdminRecordListResponse>({
    method: 'GET',
    url: '/admin/sponsors/records',
    params: {
      page: params.page ?? 1,
      page_size: params.page_size ?? 10,
      channel_id: params.channel_id?.toString(),
      nickname: params.nickname,
      is_anonymous: params.is_anonymous,
      is_hidden: params.is_hidden,
      status: params.status,
      order_by: params.order_by,
      order: params.order,
    },
  })
}

/** 添加赞助记录 */
export function createRecord(
  req: CreateRecordRequest,
): Promise<SponsorRecordAdmin> {
  return request<SponsorRecordAdmin>({
    method: 'POST',
    url: '/admin/sponsors/records',
    data: req,
  })
}

/** 更新赞助记录 */
export function updateRecord(
  id: SnowflakeID,
  req: UpdateRecordRequest,
): Promise<SponsorRecordAdmin> {
  return request<SponsorRecordAdmin>({
    method: 'PUT',
    url: `/admin/sponsors/records/${id.toString()}`,
    data: req,
  })
}

/** 删除赞助记录 */
export function deleteRecord(id: SnowflakeID): Promise<void> {
  return request<void>({
    method: 'DELETE',
    url: `/admin/sponsors/records/${id.toString()}`,
  })
}

/** 审核赞助记录（PUT /api/v1/admin/sponsors/records/:id/status） */
export function updateSponsorStatus(
  id: SnowflakeID,
  req: SponsorStatusRequest,
): Promise<void> {
  return request<void>({
    method: 'PUT',
    url: `/admin/sponsors/records/${id.toString()}/status`,
    data: req,
  })
}

// ---------------------------------------------------------------------------
// 申请 / 用户自助接口
// ---------------------------------------------------------------------------

/** 访客自助申请赞助展示（公开，POST /api/v1/sponsors/apply） */
export function applySponsor(
  req: SponsorApplyRequest,
): Promise<SponsorRecordAdmin> {
  return request<SponsorRecordAdmin>({
    method: 'POST',
    url: '/sponsors/apply',
    data: req,
  })
}

/** 我的赞助记录分页列表（GET /api/v1/user/sponsors） */
export function listMySponsors(
  params: SponsorUserParams = {},
): Promise<PaginationResponse<SponsorRecordAdmin>> {
  return request<PaginationResponse<SponsorRecordAdmin>>({
    method: 'GET',
    url: '/user/sponsors',
    params: {
      page: params.page ?? 1,
      page_size: params.page_size ?? 10,
      sponsor_status: params.sponsor_status,
    },
  })
}

/** 我的赞助记录详情（GET /api/v1/user/sponsors/:id） */
export function getMySponsor(id: SnowflakeID): Promise<SponsorRecordAdmin> {
  return request<SponsorRecordAdmin>({
    method: 'GET',
    url: `/user/sponsors/${id.toString()}`,
  })
}

/** 更新我的赞助记录（PUT /api/v1/user/sponsors/:id） */
export function updateMySponsor(
  id: SnowflakeID,
  req: SponsorUserUpdateRequest,
): Promise<SponsorRecordAdmin> {
  return request<SponsorRecordAdmin>({
    method: 'PUT',
    url: `/user/sponsors/${id.toString()}`,
    data: req,
  })
}
