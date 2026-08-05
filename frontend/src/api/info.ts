// --------------------------------------------------------------------------------
// Copyright (c) 2016-NOW(至今) 筱锋
// Author: 筱锋「xiao_lfeng」(https://www.x-lf.com)
// --------------------------------------------------------------------------------
// 许可证声明：版权所有 (c) 2016-2026 筱锋。保留所有权利。
// 有关MIT许可证的更多信息，请查看项目根目录下的LICENSE文件或访问：
// https://opensource.org/licenses/MIT
// --------------------------------------------------------------------------------

import { request } from './client'
import type {
  ApplySiteInfoResponse,
  ArchiveResponse,
  BloggerInfoResponse,
  ColorModeResponse,
  SiteInfoResponse,
  UpdateApplySiteRequest,
  UpdateArchiveRequest,
  UpdateBloggerRequest,
  UpdateColorModeRequest,
  UpdateSiteRequest,
} from './types'

/** 获取站点信息（站名、主页介绍） */
export function getSiteInfo(): Promise<SiteInfoResponse> {
  return request<SiteInfoResponse>({ method: 'GET', url: '/info/site' })
}

/** 获取站点档案（站点描述 + 自我介绍，均 Markdown） */
export function getArchive(): Promise<ArchiveResponse> {
  return request<ArchiveResponse>({ method: 'GET', url: '/info/archive' })
}

/** 获取申请站点展示（站名、描述、地址、图片、订阅、邮箱）— 供 operate/apply 交换友链场景 */
export function getApplySiteInfo(): Promise<ApplySiteInfoResponse> {
  return request<ApplySiteInfoResponse>({
    method: 'GET',
    url: '/info/apply-site',
  })
}

/** 获取博主信息（昵称、简介、博客链接、头像）— 供「关于我」名士帖展示 */
export function getBloggerInfo(): Promise<BloggerInfoResponse> {
  return request<BloggerInfoResponse>({ method: 'GET', url: '/info/blogger' })
}

/** 更新站点信息（管理端） */
export function updateSiteInfo(
  req: UpdateSiteRequest,
): Promise<SiteInfoResponse> {
  return request<SiteInfoResponse>({
    method: 'PUT',
    url: '/info/admin/site',
    data: req,
  })
}

/** 更新站点档案（管理端，站点描述与自我介绍一次保存） */
export function updateArchive(
  req: UpdateArchiveRequest,
): Promise<ArchiveResponse> {
  return request<ArchiveResponse>({
    method: 'PUT',
    url: '/info/admin/archive',
    data: req,
  })
}

/** 更新申请站点展示（管理端） */
export function updateApplySiteInfo(
  req: UpdateApplySiteRequest,
): Promise<ApplySiteInfoResponse> {
  return request<ApplySiteInfoResponse>({
    method: 'PUT',
    url: '/info/admin/apply-site',
    data: req,
  })
}

/** 更新博主信息（管理端） */
export function updateBloggerInfo(
  req: UpdateBloggerRequest,
): Promise<BloggerInfoResponse> {
  return request<BloggerInfoResponse>({
    method: 'PUT',
    url: '/info/admin/blogger',
    data: req,
  })
}

/** 获取高级配色模式（normal=普通, premium=高级） */
export function getColorMode(): Promise<ColorModeResponse> {
  return request<ColorModeResponse>({ method: 'GET', url: '/info/color-mode' })
}

/** 更新高级配色模式（管理端） */
export function updateColorMode(
  req: UpdateColorModeRequest,
): Promise<ColorModeResponse> {
  return request<ColorModeResponse>({
    method: 'PUT',
    url: '/info/admin/color-mode',
    data: req,
  })
}
