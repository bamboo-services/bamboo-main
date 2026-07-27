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
import type { UpdateAboutRequest, UpdateSiteRequest } from '@/api/types'
import { getAbout, getSiteInfo, updateAbout, updateSiteInfo } from '@/api/info'

/** 站点信息 queryKey 工厂 */
export const siteInfoKeys = {
  all: ['site-info'] as const,
  site: () => [...siteInfoKeys.all, 'site'] as const,
  about: () => [...siteInfoKeys.all, 'about'] as const,
}

/** 站点信息（站名、描述、主页介绍） */
export function useSiteInfo() {
  return useQuery({
    queryKey: siteInfoKeys.site(),
    queryFn: getSiteInfo,
  })
}

/** 自我介绍（Markdown） */
export function useAbout() {
  return useQuery({
    queryKey: siteInfoKeys.about(),
    queryFn: getAbout,
  })
}

/** 更新站点信息 */
export function useUpdateSiteInfo() {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: (req: UpdateSiteRequest) => updateSiteInfo(req),
    onSuccess: () => {
      toast.success('站点信息已保存')
      void qc.invalidateQueries({ queryKey: siteInfoKeys.site() })
    },
    onError: (err: Error) => toast.error(err.message || '站点信息保存失败'),
  })
}

/** 更新自我介绍 */
export function useUpdateAbout() {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: (req: UpdateAboutRequest) => updateAbout(req),
    onSuccess: () => {
      toast.success('自我介绍已保存')
      void qc.invalidateQueries({ queryKey: siteInfoKeys.about() })
    },
    onError: (err: Error) => toast.error(err.message || '自我介绍保存失败'),
  })
}
