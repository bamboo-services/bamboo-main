// --------------------------------------------------------------------------------
// Copyright (c) 2016-NOW(至今) 筱锋
// Author: 筱锋「xiao_lfeng」(https://www.x-lf.com)
// --------------------------------------------------------------------------------
// 许可证声明：版权所有 (c) 2016-2026 筱锋。保留所有权利。
// 有关MIT许可证的更多信息，请查看项目根目录下的LICENSE文件或访问：
// https://opensource.org/licenses/MIT
// --------------------------------------------------------------------------------

import { useMutation, useQuery } from '@tanstack/react-query'
import { getCurrentUser, logout } from '@/api/auth'
import { clearSession, getStoredUser, getToken } from '@/lib/auth'

/** 当前用户查询 key */
export const AUTH_USER_QUERY_KEY = ['auth', 'user']

/**
 * 当前登录用户。
 * - 有令牌时向后端拉取最新用户信息（/auth/user）；
 * - 拉取期间用本地缓存的用户信息兜底展示，避免闪烁；
 * - 无令牌时不发起请求。
 */
export function useAuth() {
  const token = getToken()

  const { data, isLoading } = useQuery({
    queryKey: AUTH_USER_QUERY_KEY,
    queryFn: async () => (await getCurrentUser()).user,
    enabled: !!token,
    staleTime: 5 * 60 * 1000,
  })

  const signOutMutation = useMutation({
    // 先捕获令牌并立即清空本地会话（守卫即刻生效），再携带该令牌
    // best-effort 通知服务端登出；服务端失败不影响本地登出。调用方随后
    // 整页重载，查询缓存随之重建，故无需在此手动失效。
    mutationFn: async () => {
      const currentToken = getToken() ?? undefined
      clearSession()
      if (currentToken) {
        await logout(currentToken).catch(() => undefined)
      }
    },
  })

  return {
    user: data ?? getStoredUser(),
    isLoading: !!token && isLoading,
    isAuthenticated: !!token,
    signOut: () => signOutMutation.mutateAsync(),
    isSigningOut: signOutMutation.isPending,
  }
}
