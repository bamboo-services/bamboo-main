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
  CreateLinkRequest,
  LinkListParams,
  SnowflakeID,
  UpdateLinkFailRequest,
  UpdateLinkRequest,
  UpdateLinkStatusRequest,
} from '@/api/types'
import {
  createLink,
  deleteLink,
  getAdminLink,
  listAdminLinks,
  updateLink,
  updateLinkFail,
  updateLinkStatus,
} from '@/api/link'

/** 友链 queryKey 工厂（id 一律 toString，避免 bigint 进入哈希） */
export const linkKeys = {
  all: ['admin', 'links'] as const,
  lists: () => [...linkKeys.all, 'list'] as const,
  list: (params: LinkListParams) =>
    [
      ...linkKeys.lists(),
      { ...params, link_group_id: params.link_group_id?.toString() },
    ] as const,
  details: () => [...linkKeys.all, 'detail'] as const,
  detail: (id: SnowflakeID) => [...linkKeys.details(), id.toString()] as const,
}

/** 管理端友链分页列表 */
export function useAdminLinks(params: LinkListParams = {}) {
  return useQuery({
    queryKey: linkKeys.list(params),
    queryFn: () => listAdminLinks(params),
  })
}

/** 管理端友链详情 */
export function useAdminLink(id: SnowflakeID) {
  return useQuery({
    queryKey: linkKeys.detail(id),
    queryFn: () => getAdminLink(id),
  })
}

/** 添加友链 */
export function useCreateLink() {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: (req: CreateLinkRequest) => createLink(req),
    onSuccess: () => {
      toast.success('友链添加成功')
      void qc.invalidateQueries({ queryKey: linkKeys.all })
    },
    onError: (err: Error) => toast.error(err.message || '友链添加失败'),
  })
}

/** 更新友链 */
export function useUpdateLink() {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: ({ id, req }: { id: SnowflakeID; req: UpdateLinkRequest }) =>
      updateLink(id, req),
    onSuccess: () => {
      toast.success('友链更新成功')
      void qc.invalidateQueries({ queryKey: linkKeys.all })
    },
    onError: (err: Error) => toast.error(err.message || '友链更新失败'),
  })
}

/** 删除友链 */
export function useDeleteLink() {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: (id: SnowflakeID) => deleteLink(id),
    onSuccess: () => {
      toast.success('友链删除成功')
      void qc.invalidateQueries({ queryKey: linkKeys.all })
    },
    onError: (err: Error) => toast.error(err.message || '友链删除失败'),
  })
}

/** 审核友链（通过/拒绝） */
export function useUpdateLinkStatus() {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: ({
      id,
      req,
    }: {
      id: SnowflakeID
      req: UpdateLinkStatusRequest
    }) => updateLinkStatus(id, req),
    onSuccess: () => {
      toast.success('审核状态已更新')
      void qc.invalidateQueries({ queryKey: linkKeys.all })
    },
    onError: (err: Error) => toast.error(err.message || '审核状态更新失败'),
  })
}

/** 更新友链失效状态 */
export function useUpdateLinkFail() {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: ({
      id,
      req,
    }: {
      id: SnowflakeID
      req: UpdateLinkFailRequest
    }) => updateLinkFail(id, req),
    onSuccess: () => {
      toast.success('失效状态已更新')
      void qc.invalidateQueries({ queryKey: linkKeys.all })
    },
    onError: (err: Error) => toast.error(err.message || '失效状态更新失败'),
  })
}
