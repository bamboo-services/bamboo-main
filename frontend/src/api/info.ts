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
  AboutResponse,
  BloggerInfoResponse,
  SiteInfoResponse,
  UpdateAboutRequest,
  UpdateBloggerRequest,
  UpdateSiteRequest,
} from './types'

/** 获取站点信息（站名、描述、主页介绍） */
export function getSiteInfo(): Promise<SiteInfoResponse> {
  return request<SiteInfoResponse>({ method: 'GET', url: '/info/site' })
}

/** 获取 Markdown 格式的自我介绍 */
export function getAbout(): Promise<AboutResponse> {
  return request<AboutResponse>({ method: 'GET', url: '/info/about' })
}

/** 获取博主信息（站名、描述、地址、图片、订阅、邮箱）— 供交换友链场景 */
export function getBloggerInfo(): Promise<BloggerInfoResponse> {
  return request<BloggerInfoResponse>({ method: 'GET', url: '/info/blogger' })
}

/** 更新站点信息（管理端） */
export function updateSiteInfo(
  req: UpdateSiteRequest,
): Promise<SiteInfoResponse> {
  return request<SiteInfoResponse>({
    method: 'PUT',
    url: '/admin/info/site',
    data: req,
  })
}

/** 更新自我介绍（管理端） */
export function updateAbout(req: UpdateAboutRequest): Promise<AboutResponse> {
  return request<AboutResponse>({
    method: 'PUT',
    url: '/admin/info/about',
    data: req,
  })
}

/** 更新博主信息（管理端） */
export function updateBloggerInfo(
  req: UpdateBloggerRequest,
): Promise<BloggerInfoResponse> {
  return request<BloggerInfoResponse>({
    method: 'PUT',
    url: '/admin/info/blogger',
    data: req,
  })
}
