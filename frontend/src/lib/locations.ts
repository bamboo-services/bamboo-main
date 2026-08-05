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

/** 内置「已失效」分组保留 ID（与后端 constants.BuiltinGroupInvalidID 对齐） */
export const INVALID_GROUP_ID = 1n

/** 内置分组保留 ID 集合（现仅「已失效」） */
export const BUILTIN_GROUP_IDS: ReadonlySet<bigint> = new Set([INVALID_GROUP_ID])

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
