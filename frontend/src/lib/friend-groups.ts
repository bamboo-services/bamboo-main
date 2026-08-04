/*
 * --------------------------------------------------------------------------------
 * Copyright (c) 2016-NOW(至今) 筱锋
 * Author: 筱锋「xiao_lfeng」(https://www.x-lf.com)
 * --------------------------------------------------------------------------------
 * 许可证声明：版权所有 (c) 2016-2026 筱锋。保留所有权利。
 * 有关MIT许可证的更多信息，请查看项目根目录下的LICENSE文件或访问：
 * https://opensource.org/licenses/MIT
 * --------------------------------------------------------------------------------
 */

import type { LinkFriend } from '@/api/types'

/** 友链级别枚举（与后端 pkg/constants LinkLevel 对齐：0 一般 / 1 好友 / 2 高级 / 3 广告） */
export const LINK_LEVEL = { regular: 0, close: 1, premium: 2, ad: 3 } as const

/** 章节序号（壹 贰 叁 …） */
export const CHAPTERS = [
  '壹',
  '贰',
  '叁',
  '肆',
  '伍',
  '陆',
  '柒',
  '捌',
  '玖',
  '拾',
] as const

/** 友链分组章节（groupLinksByGroup 的聚合单元；groupId 'none' 表示未分组） */
export interface FriendGroupSection {
  groupId: string
  name: string
  links: Array<LinkFriend>
}

/**
 * 友链按分组聚合（分组名取自后端嵌套的 group_f_key）
 *
 * 章节序遵循分组排序值「数字越小权重越大」（group_f_key.sort_order ASC，与位置管理一致）；
 * 未分组以最大值置底；章内按友链排序值升序。
 */
export function groupLinksByGroup(
  links: Array<LinkFriend>,
): Array<FriendGroupSection> {
  const map = new Map<
    string,
    { name: string; order: number; links: Array<LinkFriend> }
  >()
  for (const link of links) {
    const key = link.group_id != null ? link.group_id.toString() : 'none'
    const entry = map.get(key) ?? {
      name: link.group_f_key?.name ?? '未分组',
      order: link.group_f_key?.sort_order ?? Number.MAX_SAFE_INTEGER,
      links: [],
    }
    entry.links.push(link)
    map.set(key, entry)
  }
  return Array.from(map.entries())
    .sort((a, b) => a[1].order - b[1].order)
    .map(([groupId, entry]) => ({
      groupId,
      name: entry.name,
      links: entry.links.sort((a, b) => a.sort_order - b.sort_order),
    }))
}
