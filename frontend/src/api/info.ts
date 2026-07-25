// --------------------------------------------------------------------------------
// Copyright (c) 2016-NOW(至今) 筱锋
// Author: 筱锋「xiao_lfeng」(https://www.x-lf.com)
// --------------------------------------------------------------------------------
// 许可证声明：版权所有 (c) 2016-2026 筱锋。保留所有权利。
// 有关MIT许可证的更多信息，请查看项目根目录下的LICENSE文件或访问：
// https://opensource.org/licenses/MIT
// --------------------------------------------------------------------------------

import { request } from './client'
import type { AboutResponse, SiteInfoResponse } from './types'

/** 获取站点信息（站名、描述、主页介绍） */
export function getSiteInfo(): Promise<SiteInfoResponse> {
  return request<SiteInfoResponse>({ method: 'GET', url: '/info/site' })
}

/** 获取 Markdown 格式的自我介绍 */
export function getAbout(): Promise<AboutResponse> {
  return request<AboutResponse>({ method: 'GET', url: '/info/about' })
}
