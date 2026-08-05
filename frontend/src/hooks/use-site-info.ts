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
  UpdateApplySiteRequest,
  UpdateArchiveRequest,
  UpdateBloggerRequest,
  UpdateColorModeRequest,
  UpdateSiteRequest,
} from '@/api/types'
import {
  getApplySiteInfo,
  getArchive,
  getBloggerInfo,
  getColorMode,
  getSiteInfo,
  updateApplySiteInfo,
  updateArchive,
  updateBloggerInfo,
  updateColorMode,
  updateSiteInfo,
} from '@/api/info'

/** 站点信息 queryKey 工厂 */
export const siteInfoKeys = {
  all: ['site-info'] as const,
  site: () => [...siteInfoKeys.all, 'site'] as const,
  archive: () => [...siteInfoKeys.all, 'archive'] as const,
  applySite: () => [...siteInfoKeys.all, 'apply-site'] as const,
  blogger: () => [...siteInfoKeys.all, 'blogger'] as const,
  colorMode: () => [...siteInfoKeys.all, 'color-mode'] as const,
}

/** 站点信息（站名、主页介绍） */
export function useSiteInfo() {
  return useQuery({
    queryKey: siteInfoKeys.site(),
    queryFn: getSiteInfo,
  })
}

/** 站点档案（站点描述 + 自我介绍，均 Markdown） */
export function useArchive() {
  return useQuery({
    queryKey: siteInfoKeys.archive(),
    queryFn: getArchive,
  })
}

/** 申请站点展示（站名、描述、地址、图片、订阅、邮箱）— 供 operate/apply 交换友链场景 */
export function useApplySiteInfo() {
  return useQuery({
    queryKey: siteInfoKeys.applySite(),
    queryFn: getApplySiteInfo,
  })
}

/** 博主信息（昵称、简介、博客链接、头像）— 供「关于我」名士帖展示 */
export function useBloggerInfo() {
  return useQuery({
    queryKey: siteInfoKeys.blogger(),
    queryFn: getBloggerInfo,
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

/** 更新站点档案（站点描述与自我介绍一次保存） */
export function useUpdateArchive() {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: (req: UpdateArchiveRequest) => updateArchive(req),
    onSuccess: () => {
      toast.success('站点档案已保存')
      void qc.invalidateQueries({ queryKey: siteInfoKeys.archive() })
    },
    onError: (err: Error) => toast.error(err.message || '站点档案保存失败'),
  })
}

/** 更新申请站点展示 */
export function useUpdateApplySiteInfo() {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: (req: UpdateApplySiteRequest) => updateApplySiteInfo(req),
    onSuccess: () => {
      toast.success('申请站点展示已保存')
      void qc.invalidateQueries({ queryKey: siteInfoKeys.applySite() })
    },
    onError: (err: Error) => toast.error(err.message || '申请站点展示保存失败'),
  })
}

/** 更新博主信息 */
export function useUpdateBloggerInfo() {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: (req: UpdateBloggerRequest) => updateBloggerInfo(req),
    onSuccess: () => {
      toast.success('博主信息已保存')
      void qc.invalidateQueries({ queryKey: siteInfoKeys.blogger() })
    },
    onError: (err: Error) => toast.error(err.message || '博主信息保存失败'),
  })
}

/** 高级配色模式（normal=普通, premium=高级） */
export function useColorMode() {
  return useQuery({
    queryKey: siteInfoKeys.colorMode(),
    queryFn: getColorMode,
  })
}

/** 更新高级配色模式 */
export function useUpdateColorMode() {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: (req: UpdateColorModeRequest) => updateColorMode(req),
    onSuccess: () => {
      toast.success('高级配色模式已更新')
      void qc.invalidateQueries({ queryKey: siteInfoKeys.colorMode() })
      // 开关影响颜色选择器可见范围，联动刷新公开与管理端颜色列表
      void qc.invalidateQueries({ queryKey: ['public', 'colors'] })
      void qc.invalidateQueries({ queryKey: ['admin', 'colors'] })
    },
    onError: (err: Error) => toast.error(err.message || '高级配色模式更新失败'),
  })
}
