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
  ColorListParams,
  ColorPageParams,
  CreateColorRequest,
  SnowflakeID,
  UpdateColorRequest,
} from '@/api/types'
import {
  createColor,
  deleteColor,
  getAllColors,
  listColors,
  updateColor,
  updateColorStatus,
} from '@/api/color'

/** 友链颜色 queryKey 工厂 */
export const colorKeys = {
  all: ['admin', 'colors'] as const,
  lists: () => [...colorKeys.all, 'list'] as const,
  list: (params: ColorPageParams) => [...colorKeys.lists(), params] as const,
  allList: () => [...colorKeys.all, 'all'] as const,
}

/** 友链颜色分页列表 */
export function useColors(params: ColorPageParams = {}) {
  return useQuery({
    queryKey: colorKeys.list(params),
    queryFn: () => listColors(params),
  })
}

/** 友链颜色全量列表（供选择器使用） */
export function useAllColors(params: ColorListParams = {}) {
  return useQuery({
    queryKey: [...colorKeys.allList(), params],
    queryFn: () => getAllColors(params),
  })
}

/** 添加友链颜色 */
export function useCreateColor() {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: (req: CreateColorRequest) => createColor(req),
    onSuccess: () => {
      toast.success('颜色添加成功')
      void qc.invalidateQueries({ queryKey: colorKeys.all })
    },
    onError: (err: Error) => toast.error(err.message || '颜色添加失败'),
  })
}

/** 更新友链颜色 */
export function useUpdateColor() {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: ({ id, req }: { id: SnowflakeID; req: UpdateColorRequest }) =>
      updateColor(id, req),
    onSuccess: () => {
      toast.success('颜色更新成功')
      void qc.invalidateQueries({ queryKey: colorKeys.all })
    },
    onError: (err: Error) => toast.error(err.message || '颜色更新失败'),
  })
}

/** 切换友链颜色启用/禁用 */
export function useUpdateColorStatus() {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: ({ id, status }: { id: SnowflakeID; status: boolean }) =>
      updateColorStatus(id, status),
    onSuccess: () => {
      toast.success('颜色状态已更新')
      void qc.invalidateQueries({ queryKey: colorKeys.all })
    },
    onError: (err: Error) => toast.error(err.message || '颜色状态更新失败'),
  })
}

/** 删除友链颜色 */
export function useDeleteColor() {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: ({ id, force }: { id: SnowflakeID; force?: boolean }) =>
      deleteColor(id, force),
    onSuccess: () => {
      toast.success('颜色删除成功')
      void qc.invalidateQueries({ queryKey: colorKeys.all })
    },
    onError: (err: Error) => toast.error(err.message || '颜色删除失败'),
  })
}
