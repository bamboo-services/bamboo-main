// --------------------------------------------------------------------------------
// Copyright (c) 2016-NOW(至今) 筱锋
// Author: 筱锋「xiao_lfeng」(https://www.x-lf.com)
// --------------------------------------------------------------------------------
// 许可证声明：版权所有 (c) 2016-2026 筱锋。保留所有权利。
// 有关MIT许可证的更多信息，请查看项目根目录下的LICENSE文件或访问：
// https://opensource.org/licenses/MIT
// --------------------------------------------------------------------------------

/**
 * 后端统一响应包装（见 bamboo-base-go BaseResponse）
 * - code: HTTP 风格状态码，成功为 200
 * - data: 业务数据，失败时省略
 */
export interface BaseResponse<T> {
  context: string
  output: string
  code: number
  message: string
  error_message?: string
  overhead?: number
  data: T
}

/** 站点信息（GET /api/v1/info/site） */
export interface SiteInfoResponse {
  site_name: string
  site_description: string
  introduction: string
  updated_at: string
}

/** 自我介绍（GET /api/v1/info/about） */
export interface AboutResponse {
  content: string
  updated_at: string
}

/** 友链公开条目（GET /api/v1/links） */
export interface LinkFriend {
  id: number
  name: string
  url: string
  avatar: string | null
  rss: string | null
  description: string | null
  email: string | null
  group_id: number | null
  color_id: number | null
  sort_order: number
  status: number
  is_failure: number
  fail_reason: string | null
  apply_remark: string | null
  review_remark: string | null
  created_at: string
  updated_at: string
}

export interface FriendPublicResponse {
  links: Array<LinkFriend>
}

/** 赞助渠道（GET /api/v1/sponsors/channels） */
export interface SponsorChannel {
  id: number
  name: string
  icon: string | null
  sort_order: number
  status: boolean
  sponsor_count: number
}

/** 赞助记录中的渠道简表 */
export interface SponsorChannelSimple {
  id: number
  name: string
  icon: string | null
}

/** 赞助记录公开条目（GET /api/v1/sponsors/records） */
export interface SponsorRecord {
  id: number
  nickname: string
  redirect_url: string | null
  amount: number
  message: string | null
  sponsor_at: string | null
  channel: SponsorChannelSimple | null
}

/** 分页信息 */
export interface PaginationInfo {
  page: number
  page_size: number
  total: number
  total_pages: number
  has_next: boolean
  has_prev: boolean
}

/** 分页响应（data + pagination） */
export interface PaginationResponse<T> {
  data: Array<T>
  pagination: PaginationInfo
}
