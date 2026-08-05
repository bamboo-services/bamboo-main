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

import type { LinkGroup } from '@/api/types'

/** 内置「首页」分组保留 ID（与后端 constants.BuiltinGroupHomepageID 对齐） */
export const HOMEPAGE_GROUP_ID = 1n

/** 内置「友链页」分组保留 ID（与后端 constants.BuiltinGroupFriendsID 对齐） */
export const FRIENDS_GROUP_ID = 2n

/** 内置分组保留 ID 集合 */
export const BUILTIN_GROUP_IDS: ReadonlySet<bigint> = new Set([
  HOMEPAGE_GROUP_ID,
  FRIENDS_GROUP_ID,
])

/** 内置分组名称（后端注入虚拟记录时固定） */
export const BUILTIN_GROUP_NAMES: ReadonlyMap<bigint, string> = new Map([
  [HOMEPAGE_GROUP_ID, '首页'],
  [FRIENDS_GROUP_ID, '友链页'],
])

/** 内置分组优先级（数字越小越靠前；非内置返回 null），供聚合排序保证恒置顶 */
export function builtinPriority(id: bigint | string): number | null {
  const bid = typeof id === 'string' ? BigInt(id) : id
  if (bid === HOMEPAGE_GROUP_ID) return 0
  if (bid === FRIENDS_GROUP_ID) return 1
  return null
}

/** 判断分组 ID 是否为内置保留 ID */
export function isBuiltinGroupId(
  id: bigint | string | null | undefined,
): boolean {
  if (id == null) return false
  const bid = typeof id === 'string' ? BigInt(id) : id
  return BUILTIN_GROUP_IDS.has(bid)
}

/** 判断分组对象是否为内置分组 */
export function isBuiltinGroup(
  group: LinkGroup | null | undefined,
): boolean {
  return isBuiltinGroupId(group?.id)
}
