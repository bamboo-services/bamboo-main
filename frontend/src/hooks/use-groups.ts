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
  CreateGroupRequest,
  GroupListParams,
  GroupPageParams,
  SnowflakeID,
  UpdateBuiltinInvalidGroupRequest,
  UpdateGroupRequest,
} from '@/api/types'
import {
  createGroup,
  deleteGroup,
  getAllGroups,
  getBuiltinInvalidGroup,
  listGroups,
  sortGroups,
  updateBuiltinInvalidGroup,
  updateGroup,
  updateGroupStatus,
} from '@/api/group'
import { failedLinkKeys } from '@/hooks/use-links'

/** 友链分组 queryKey 工厂 */
export const groupKeys = {
  all: ['admin', 'groups'] as const,
  lists: () => [...groupKeys.all, 'list'] as const,
  list: (params: GroupPageParams) => [...groupKeys.lists(), params] as const,
  allList: () => [...groupKeys.all, 'all'] as const,
}

/** 友链分组分页列表 */
export function useGroups(params: GroupPageParams = {}) {
  return useQuery({
    queryKey: groupKeys.list(params),
    queryFn: () => listGroups(params),
  })
}

/** 友链分组全量列表（供选择器使用） */
export function useAllGroups(params: GroupListParams = {}) {
  return useQuery({
    queryKey: [...groupKeys.allList(), params],
    queryFn: () => getAllGroups(params),
  })
}

/** 添加友链分组 */
export function useCreateGroup() {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: (req: CreateGroupRequest) => createGroup(req),
    onSuccess: () => {
      toast.success('分组添加成功')
      void qc.invalidateQueries({ queryKey: groupKeys.all })
    },
    onError: (err: Error) => toast.error(err.message || '分组添加失败'),
  })
}

/** 更新友链分组 */
export function useUpdateGroup() {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: ({ id, req }: { id: SnowflakeID; req: UpdateGroupRequest }) =>
      updateGroup(id, req),
    onSuccess: () => {
      toast.success('分组更新成功')
      void qc.invalidateQueries({ queryKey: groupKeys.all })
    },
    onError: (err: Error) => toast.error(err.message || '分组更新失败'),
  })
}

/** 切换友链分组启用/禁用 */
export function useUpdateGroupStatus() {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: ({ id, status }: { id: SnowflakeID; status: boolean }) =>
      updateGroupStatus(id, status),
    onSuccess: () => {
      toast.success('分组状态已更新')
      void qc.invalidateQueries({ queryKey: groupKeys.all })
    },
    onError: (err: Error) => toast.error(err.message || '分组状态更新失败'),
  })
}

/** 删除友链分组（force=true 移至未分组；targetGroupId 迁移到指定分组，二者互斥） */
export function useDeleteGroup() {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: ({
      id,
      force,
      targetGroupId,
    }: {
      id: SnowflakeID
      force?: boolean
      targetGroupId?: SnowflakeID
    }) =>
      deleteGroup(id, {
        ...(force ? { force: true } : {}),
        ...(targetGroupId ? { target_group_id: targetGroupId } : {}),
      }),
    onSuccess: () => {
      toast.success('分组删除成功')
      void qc.invalidateQueries({ queryKey: groupKeys.all })
    },
    onError: (err: Error) => toast.error(err.message || '分组删除失败'),
  })
}

/** 批量更新分组排序（排位看板章节拖拽联动；成功提示由友链排序 mutation 统一承担） */
export function useSortGroups() {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: (groupIds: Array<SnowflakeID>) => sortGroups(groupIds),
    onSuccess: () => {
      void qc.invalidateQueries({ queryKey: groupKeys.all })
    },
    onError: (err: Error) => toast.error(err.message || '分组排序保存失败'),
  })
}

/** 内置「已失效」分组配置 queryKey */
const builtinInvalidKeys = {
  detail: ['admin', 'groups', 'builtin', 'invalid'] as const,
}

/** 获取内置「已失效」分组配置（GET /api/v1/info/builtin-invalid-group） */
export function useBuiltinInvalidGroup() {
  return useQuery({
    queryKey: builtinInvalidKeys.detail,
    queryFn: getBuiltinInvalidGroup,
  })
}

/** 更新内置「已失效」分组配置（经 bm_system 热修改，同步刷新公开「已失效」章节） */
export function useUpdateBuiltinInvalidGroup() {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: (req: UpdateBuiltinInvalidGroupRequest) =>
      updateBuiltinInvalidGroup(req),
    onSuccess: () => {
      toast.success('已失效分组配置已更新')
      void qc.invalidateQueries({ queryKey: builtinInvalidKeys.detail })
      // 分组名称变更需同步公开「已失效」章节标题
      void qc.invalidateQueries({ queryKey: failedLinkKeys.all })
    },
    onError: (err: Error) => toast.error(err.message || '已失效分组配置更新失败'),
  })
}
