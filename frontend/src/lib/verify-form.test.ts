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

import { describe, expect, it } from 'vitest'
import type { LinkFriend } from '@/api/types'
import { initVerifyForm, verifyFormToUpdateReq } from './verify-form'

/** 构造最小 LinkFriend 测试样本，overrides 合并覆盖默认值。
 *  可空字段显式置 null，贴合后端 JSON 序列化的真实响应（null 而非字段缺失）。 */
function makeLink(overrides: Partial<LinkFriend> = {}): LinkFriend {
  return {
    id: 1n,
    name: '测试站点',
    url: 'https://example.com',
    avatar: null,
    rss: null,
    description: null,
    email: null,
    user_id: null,
    group_id: null,
    color_id: null,
    expected_group_id: null,
    expected_color_id: null,
    sort_order: 0,
    status: 0,
    is_failure: 0,
    level: 0,
    fail_reason: null,
    apply_remark: null,
    review_remark: null,
    screenshot_url: null,
    screenshot_at: null,
    updated_at: '2026-08-10T00:00:00Z',
    ...overrides,
  }
}

describe('initVerifyForm', () => {
  it('期望值优先预填：新申请存在期望位置/颜色时采用期望值', () => {
    const form = initVerifyForm(
      makeLink({ expected_group_id: 10n, expected_color_id: 20n }),
    )
    expect(form.groupId).toBe(10n)
    expect(form.colorId).toBe(20n)
  })

  it('期望值为空回退正式值：历史待审数据无期望值时回退 group_id/color_id', () => {
    const form = initVerifyForm(makeLink({ group_id: 5n, color_id: 6n }))
    expect(form.groupId).toBe(5n)
    expect(form.colorId).toBe(6n)
  })

  it('期望值与正式值均未提供时为 null', () => {
    const form = initVerifyForm(makeLink())
    expect(form.groupId).toBeNull()
    expect(form.colorId).toBeNull()
  })
})

describe('verifyFormToUpdateReq', () => {
  const form = {
    siteName: ' 测试站点 ',
    siteUrl: ' https://example.com ',
    siteLogo: '',
    siteRss: '',
    webmasterEmail: 'a@b.c',
    siteDescription: ' 站点描述 ',
    groupId: 10n,
    colorId: 20n,
    applyRemark: ' 申请备注 ',
  }

  it('默认携带位置/颜色（审核通过时预填的期望值随表单落正式值）', () => {
    const req = verifyFormToUpdateReq(form)
    expect(req.link_name).toBe('测试站点')
    expect(req.link_group_id).toBe(10n)
    expect(req.link_color_id).toBe(20n)
    expect(req.link_apply_remark).toBe('申请备注')
  })

  it('includeLocation=false 省略位置/颜色键（拒绝时不落正式位置/颜色）', () => {
    const req = verifyFormToUpdateReq(form, { includeLocation: false })
    expect(req.link_group_id).toBeUndefined()
    expect(req.link_color_id).toBeUndefined()
  })

  it('groupId/colorId 为 null 时清空正式位置/颜色', () => {
    const req = verifyFormToUpdateReq({ ...form, groupId: null, colorId: null })
    expect(req.link_group_id).toBeNull()
    expect(req.link_color_id).toBeNull()
  })
})
