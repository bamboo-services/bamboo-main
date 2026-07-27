// --------------------------------------------------------------------------------
// Copyright (c) 2016-NOW(至今) 筱锋
// Author: 筱锋「xiao_lfeng」(https://www.x-lf.com)
// --------------------------------------------------------------------------------
// 许可证声明：版权所有 (c) 2016-2026 筱锋。保留所有权利。
// 有关MIT许可证的更多信息，请查看项目根目录下的LICENSE文件或访问：
// https://opensource.org/licenses/MIT
// --------------------------------------------------------------------------------

import { useQuery } from '@tanstack/react-query'
import { getDashboardStats } from '@/api/dashboard'
import { getHealth } from '@/api/public'

/** 仪表盘统计 queryKey */
export const dashboardKeys = {
  stats: () => ['admin', 'dashboard', 'stats'] as const,
}

/** 健康检查 queryKey */
export const healthKeys = {
  status: () => ['health'] as const,
}

/** 仪表盘统计（友链计数 + 最近申请） */
export function useDashboardStats() {
  return useQuery({
    queryKey: dashboardKeys.stats(),
    queryFn: getDashboardStats,
  })
}

/** 系统健康检查（运行时指标） */
export function useHealth() {
  return useQuery({
    queryKey: healthKeys.status(),
    queryFn: getHealth,
    refetchInterval: 30_000, // 每 30s 刷新一次运行时指标
  })
}
