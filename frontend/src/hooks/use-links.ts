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
  ApplyLinkRequest,
  CreateLinkRequest,
  LinkListParams,
  SnowflakeID,
  UpdateLinkFailRequest,
  UpdateLinkRequest,
  UpdateLinkStatusRequest,
  UpdateProfileRequest,
  UpdateUserLinkRequest,
  UserLinkParams,
} from '@/api/types'
import {
  applyLink,
  createLink,
  deleteLink,
  getAdminLink,
  getMyLink,
  listAdminLinks,
  listMyLinks,
  requestTakedown,
  updateLink,
  updateLinkFail,
  updateLinkStatus,
  updateMyLink,
  updateProfile,
} from '@/api/link'
import { AUTH_USER_QUERY_KEY } from '@/hooks/use-auth'

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

// ---------------------------------------------------------------------------
// 用户自助友链（与 admin 查询缓存隔离）
// ---------------------------------------------------------------------------

/** 用户友链 queryKey 工厂（id 一律 toString，避免 bigint 进入哈希） */
export const myLinkKeys = {
  all: ['user', 'links'] as const,
  lists: () => [...myLinkKeys.all, 'list'] as const,
  list: (params: UserLinkParams) => [...myLinkKeys.lists(), { ...params }] as const,
  details: () => [...myLinkKeys.all, 'detail'] as const,
  detail: (id: SnowflakeID) => [...myLinkKeys.details(), id.toString()] as const,
}

/** 我的友链分页列表 */
export function useMyLinks(params: UserLinkParams = {}) {
  return useQuery({
    queryKey: myLinkKeys.list(params),
    queryFn: () => listMyLinks(params),
  })
}

/** 我的友链详情 */
export function useMyLink(id: SnowflakeID) {
  return useQuery({
    queryKey: myLinkKeys.detail(id),
    queryFn: () => getMyLink(id),
  })
}

/** 访客自助申请友链（成功后失效我的友链缓存，便于登录态下即时可见） */
export function useApplyLink() {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: (req: ApplyLinkRequest) => applyLink(req),
    onSuccess: () => {
      void qc.invalidateQueries({ queryKey: myLinkKeys.all })
    },
    onError: (err: Error) => toast.error(err.message || '友链申请失败'),
  })
}

/** 更新我的友链 */
export function useUpdateMyLink() {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: ({
      id,
      req,
    }: {
      id: SnowflakeID
      req: UpdateUserLinkRequest
    }) => updateMyLink(id, req),
    onSuccess: () => {
      toast.success('友链更新成功')
      void qc.invalidateQueries({ queryKey: myLinkKeys.all })
    },
    onError: (err: Error) => toast.error(err.message || '友链更新失败'),
  })
}

/** 申请下架我的友链 */
export function useRequestTakedown() {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: (id: SnowflakeID) => requestTakedown(id),
    onSuccess: () => {
      toast.success('下架申请已提交，请等待管理员审核')
      void qc.invalidateQueries({ queryKey: myLinkKeys.all })
    },
    onError: (err: Error) => toast.error(err.message || '下架申请失败'),
  })
}

/** 更新用户资料（成功后失效当前用户缓存） */
export function useUpdateProfile() {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: (req: UpdateProfileRequest) => updateProfile(req),
    onSuccess: () => {
      toast.success('资料更新成功')
      void qc.invalidateQueries({ queryKey: AUTH_USER_QUERY_KEY })
    },
    onError: (err: Error) => toast.error(err.message || '资料更新失败'),
  })
}
