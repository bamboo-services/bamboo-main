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
  UpdateChannelRequest,
  UpdateRecordRequest,
} from '@/api/types'
import {
  createChannel,
  createRecord,
  deleteChannel,
  deleteRecord,
  getAllChannels,
  listAdminChannels,
  listAdminRecords,
  updateChannel,
  updateChannelStatus,
  updateRecord,
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
    },
    onError: (err: Error) => toast.error(err.message || '赞助记录删除失败'),
  })
}
