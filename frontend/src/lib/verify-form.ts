/*
 * --------------------------------------------------------------------------------
 * Copyright (c) 2016-NOW(至今) 筱锋
 * Author: 筱锋「xiao_lfeng」(https://www.x-lf.com)
 * --------------------------------------------------------------------------------
 * 许可证声明：版权所有 (c) 2016-2026 筱锋。保留所有权利。
 * 有关MIT许可证的更多信息，请查看项目根目录下的 LICENSE 文件或访问：
 * https://opensource.org/licenses/MIT
 * --------------------------------------------------------------------------------
 */

import type { LinkFriend, SnowflakeID, UpdateLinkRequest } from '@/api/types'

/** 审核详情面板的可编辑表单状态 */
export interface VerifyFormState {
  siteName: string
  siteUrl: string
  siteLogo: string
  siteRss: string
  webmasterEmail: string
  siteDescription: string
  groupId: SnowflakeID | null
  colorId: SnowflakeID | null
  applyRemark: string
}

/** 从 LinkFriend 初始化编辑表单（期望值优先，回退当前正式值） */
export function initVerifyForm(link: LinkFriend): VerifyFormState {
  return {
    siteName: link.name,
    siteUrl: link.url,
    siteLogo: link.avatar ?? '',
    siteRss: link.rss ?? '',
    webmasterEmail: link.email ?? '',
    siteDescription: link.description ?? '',
    groupId: link.expected_group_id ?? link.group_id,
    colorId: link.expected_color_id ?? link.color_id,
    applyRemark: link.apply_remark ?? '',
  }
}

/**
 * 将编辑表单转为 UpdateLinkRequest。
 * `includeLocation` 默认 true（审核通过时预填的期望值随表单落正式位置/颜色）；
 * 拒绝审核时传 false，避免把期望值持久化为被拒友链的正式位置/颜色。
 */
export function verifyFormToUpdateReq(
  form: VerifyFormState,
  opts: { includeLocation?: boolean } = {},
): UpdateLinkRequest {
  const req: UpdateLinkRequest = {
    link_name: form.siteName.trim(),
    link_url: form.siteUrl.trim(),
    link_avatar: form.siteLogo.trim() || undefined,
    link_rss: form.siteRss.trim() || undefined,
    link_email: form.webmasterEmail.trim() || undefined,
    link_desc: form.siteDescription.trim() || undefined,
    link_apply_remark: form.applyRemark.trim() || undefined,
  }
  if (opts.includeLocation !== false) {
    req.link_group_id = form.groupId ?? null
    req.link_color_id = form.colorId ?? null
  }
  return req
}
