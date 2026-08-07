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
  SnowflakeID,
  UserListParams,
  UserStatusRequest,
} from '@/api/types'
import { listAdminUsers, updateUserStatus } from '@/api/user'

/** 用户 queryKey 工厂（id 一律 toString，避免 bigint 进入哈希） */
export const userKeys = {
  all: ['admin', 'users'] as const,
  lists: () => [...userKeys.all, 'list'] as const,
  list: (params: UserListParams) =>
    [...userKeys.lists(), { ...params }] as const,
}

/** 管理端用户分页列表 */
export function useAdminUsers(params: UserListParams = {}) {
  return useQuery({
    queryKey: userKeys.list(params),
    queryFn: () => listAdminUsers(params),
  })
}

/** 更新用户启用/禁用状态（成功后失效用户列表缓存） */
export function useUpdateUserStatus() {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: ({
      id,
      req,
    }: {
      id: SnowflakeID
      req: UserStatusRequest
    }) => updateUserStatus(id, req),
    onSuccess: () => {
      toast.success('用户状态已更新')
      void qc.invalidateQueries({ queryKey: userKeys.all })
    },
    onError: (err: Error) => toast.error(err.message || '用户状态更新失败'),
  })
}
