// --------------------------------------------------------------------------------
// Copyright (c) 2016-NOW(至今) 筱锋
// Author: 筱锋「xiao_lfeng」(https://www.x-lf.com)
// --------------------------------------------------------------------------------
// 许可证声明：版权所有 (c) 2016-2026 筱锋。保留所有权利。
// 有关MIT许可证的更多信息，请查看项目根目录下的LICENSE文件或访问：
// https://opensource.org/licenses/MIT
// --------------------------------------------------------------------------------

import type { UserInfo } from '@/api/types'
import { getStoredUser } from '@/lib/auth'

/**
 * 判断用户是否为系统唯一管理员。接受显式传入的 user，缺省时读取本地缓存的当前用户，
 * 便于在 beforeLoad 等同步上下文中无 hook 时调用。
 *
 * 管理员身份由后端计算字段 `is_admin` 承载（对齐 bm_system.system.admin.id 判定）。
 */
export function isAdmin(user?: UserInfo | null): boolean {
  const target = user ?? getStoredUser()
  return target?.is_admin === true
}

/**
 * 取当前用户身份徽章文案：管理员 → 管理员，其余 → 用户。
 * 缺省 user 时回退读取本地缓存用户。
 */
export function adminLabel(user?: UserInfo | null): string {
  return isAdmin(user) ? '管理员' : '用户'
}
