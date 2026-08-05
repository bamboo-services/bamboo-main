// --------------------------------------------------------------------------------
// Copyright (c) 2016-NOW(至今) 筱锋
// Author: 筱锋「xiao_lfeng」(https://www.x-lf.com)
// --------------------------------------------------------------------------------
// 许可证声明：版权所有 (c) 2016-2026 筱锋。保留所有权利。
// 有关MIT许可证的更多信息，请查看项目根目录下的LICENSE文件或访问：
// https://opensource.org/licenses/MIT
// --------------------------------------------------------------------------------

import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query'
import { toast } from 'sonner'
import type {
  ChannelPageParams,
  CreateChannelRequest,
  CreateRecordRequest,
  RecordPageParams,
  SnowflakeID,
  SponsorApplyRequest,
  SponsorStatusRequest,
  SponsorUserParams,
  SponsorUserUpdateRequest,
  UpdateChannelRequest,
  UpdateRecordRequest,
} from '@/api/types'
import type { GetRecordsParams } from '@/api/sponsor'
import {
  applySponsor,
  createChannel,
  createRecord,
  deleteChannel,
  deleteRecord,
  getAllChannels,
  getMySponsor,
  getPublicChannels,
  getPublicRecords,
  listAdminChannels,
  listAdminRecords,
  listMySponsors,
  updateChannel,
  updateChannelStatus,
  updateMySponsor,
  updateRecord,
  updateSponsorStatus,
} from '@/api/sponsor'

/** 赞助渠道 queryKey 工厂 */
export const channelKeys = {
  all: ['admin', 'channels'] as const,
  lists: () => [...channelKeys.all, 'list'] as const,
  list: (params: ChannelPageParams) =>
    [...channelKeys.lists(), params] as const,
  allList: () => [...channelKeys.all, 'all'] as const,
}

/** 赞助记录 queryKey 工厂（channel_id 转字符串入 key） */
export const recordKeys = {
  all: ['admin', 'records'] as const,
  lists: () => [...recordKeys.all, 'list'] as const,
  list: (params: RecordPageParams) =>
    [
      ...recordKeys.lists(),
      { ...params, channel_id: params.channel_id?.toString() },
    ] as const,
}

/** 管理端赞助渠道分页列表 */
export function useAdminChannels(params: ChannelPageParams = {}) {
  return useQuery({
    queryKey: channelKeys.list(params),
    queryFn: () => listAdminChannels(params),
  })
}

/** 赞助渠道全量列表（供记录表单选择器使用） */
export function useAllChannels() {
  return useQuery({
    queryKey: channelKeys.allList(),
    queryFn: () => getAllChannels(),
  })
}

/** 添加赞助渠道 */
export function useCreateChannel() {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: (req: CreateChannelRequest) => createChannel(req),
    onSuccess: () => {
      toast.success('赞助渠道添加成功')
      void qc.invalidateQueries({ queryKey: channelKeys.all })
    },
    onError: (err: Error) => toast.error(err.message || '赞助渠道添加失败'),
  })
}

/** 更新赞助渠道 */
export function useUpdateChannel() {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: ({ id, req }: { id: SnowflakeID; req: UpdateChannelRequest }) =>
      updateChannel(id, req),
    onSuccess: () => {
      toast.success('赞助渠道更新成功')
      void qc.invalidateQueries({ queryKey: channelKeys.all })
    },
    onError: (err: Error) => toast.error(err.message || '赞助渠道更新失败'),
  })
}

/** 切换赞助渠道启用/禁用 */
export function useUpdateChannelStatus() {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: ({ id, status }: { id: SnowflakeID; status: boolean }) =>
      updateChannelStatus(id, status),
    onSuccess: () => {
      toast.success('渠道状态已更新')
      void qc.invalidateQueries({ queryKey: channelKeys.all })
    },
    onError: (err: Error) => toast.error(err.message || '渠道状态更新失败'),
  })
}

/** 删除赞助渠道 */
export function useDeleteChannel() {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: (id: SnowflakeID) => deleteChannel(id),
    onSuccess: () => {
      toast.success('赞助渠道删除成功')
      void qc.invalidateQueries({ queryKey: channelKeys.all })
    },
    onError: (err: Error) => toast.error(err.message || '赞助渠道删除失败'),
  })
}

/** 管理端赞助记录分页列表 */
export function useAdminRecords(params: RecordPageParams = {}) {
  return useQuery({
    queryKey: recordKeys.list(params),
    queryFn: () => listAdminRecords(params),
  })
}

/** 添加赞助记录 */
export function useCreateRecord() {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: (req: CreateRecordRequest) => createRecord(req),
    onSuccess: () => {
      toast.success('赞助记录添加成功')
      void qc.invalidateQueries({ queryKey: recordKeys.all })
      void qc.invalidateQueries({ queryKey: channelKeys.all })
      void qc.invalidateQueries({ queryKey: publicSponsorKeys.records() })
    },
    onError: (err: Error) => toast.error(err.message || '赞助记录添加失败'),
  })
}

/** 更新赞助记录 */
export function useUpdateRecord() {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: ({ id, req }: { id: SnowflakeID; req: UpdateRecordRequest }) =>
      updateRecord(id, req),
    onSuccess: () => {
      toast.success('赞助记录更新成功')
      void qc.invalidateQueries({ queryKey: recordKeys.all })
      void qc.invalidateQueries({ queryKey: publicSponsorKeys.records() })
    },
    onError: (err: Error) => toast.error(err.message || '赞助记录更新失败'),
  })
}

/** 删除赞助记录 */
export function useDeleteRecord() {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: (id: SnowflakeID) => deleteRecord(id),
    onSuccess: () => {
      toast.success('赞助记录删除成功')
      void qc.invalidateQueries({ queryKey: recordKeys.all })
      void qc.invalidateQueries({ queryKey: channelKeys.all })
      void qc.invalidateQueries({ queryKey: publicSponsorKeys.records() })
    },
    onError: (err: Error) => toast.error(err.message || '赞助记录删除失败'),
  })
}

// ---------------------------------------------------------------------------
// 公开接口（前台展示 / 申请表单选择器）
// ---------------------------------------------------------------------------

/** 公开赞助 queryKey 工厂（与 admin 查询缓存隔离） */
export const publicSponsorKeys = {
  channels: ['public', 'sponsors', 'channels'] as const,
  records: (page = 1) => ['public', 'sponsors', 'records', page] as const,
}

/** 公开赞助渠道列表（供申请表单选择器使用） */
export function usePublicChannels() {
  return useQuery({
    queryKey: publicSponsorKeys.channels,
    queryFn: getPublicChannels,
    staleTime: 10 * 60 * 1000,
  })
}

/** 公开赞助记录分页列表（前台赞助墙） */
export function usePublicRecords(params: GetRecordsParams = {}) {
  return useQuery({
    queryKey: publicSponsorKeys.records(params.page ?? 1),
    queryFn: () => getPublicRecords(params),
  })
}

// ---------------------------------------------------------------------------
// 用户自助赞助（与 admin 查询缓存隔离）
// ---------------------------------------------------------------------------

/** 我的赞助 queryKey 工厂（id 一律 toString，避免 bigint 进入哈希） */
export const mySponsorKeys = {
  all: ['user', 'sponsors'] as const,
  lists: () => [...mySponsorKeys.all, 'list'] as const,
  list: (params: SponsorUserParams) =>
    [...mySponsorKeys.lists(), { ...params }] as const,
  details: () => [...mySponsorKeys.all, 'detail'] as const,
  detail: (id: SnowflakeID) =>
    [...mySponsorKeys.details(), id.toString()] as const,
}

/** 我的赞助记录分页列表 */
export function useMySponsors(params: SponsorUserParams = {}) {
  return useQuery({
    queryKey: mySponsorKeys.list(params),
    queryFn: () => listMySponsors(params),
  })
}

/** 我的赞助记录详情 */
export function useMySponsor(id: SnowflakeID) {
  return useQuery({
    queryKey: mySponsorKeys.detail(id),
    queryFn: () => getMySponsor(id),
  })
}

/** 访客自助申请赞助展示（成功后失效我的赞助与公开账册缓存） */
export function useApplySponsor() {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: (req: SponsorApplyRequest) => applySponsor(req),
    onSuccess: () => {
      void qc.invalidateQueries({ queryKey: mySponsorKeys.all })
      void qc.invalidateQueries({ queryKey: publicSponsorKeys.records() })
    },
    onError: (err: Error) => toast.error(err.message || '赞助申请失败'),
  })
}

/** 更新我的赞助记录 */
export function useUpdateMySponsor() {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: ({
      id,
      req,
    }: {
      id: SnowflakeID
      req: SponsorUserUpdateRequest
    }) => updateMySponsor(id, req),
    onSuccess: () => {
      toast.success('赞助记录更新成功')
      void qc.invalidateQueries({ queryKey: mySponsorKeys.all })
    },
    onError: (err: Error) => toast.error(err.message || '赞助记录更新失败'),
  })
}

/** 审核赞助记录（成功后失效管理端 / 公开 / 我的赞助三端缓存） */
export function useUpdateSponsorStatus() {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: ({ id, req }: { id: SnowflakeID; req: SponsorStatusRequest }) =>
      updateSponsorStatus(id, req),
    onSuccess: () => {
      toast.success('审核状态已更新')
      void qc.invalidateQueries({ queryKey: recordKeys.all })
      void qc.invalidateQueries({ queryKey: publicSponsorKeys.records() })
      void qc.invalidateQueries({ queryKey: mySponsorKeys.all })
    },
    onError: (err: Error) => toast.error(err.message || '审核状态更新失败'),
  })
}
