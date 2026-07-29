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
 * 角色常量，对齐后端 `pkg/constants/context.go` 的 Role* 常量。
 * - admin：管理员，可访问 /admin 管理后台
 * - moderator：协作者，按普通用户处理（暂无专属界面）
 * - user：普通用户，登录后进入 /user 用户中心
 */
export const ROLE_ADMIN = 'admin'
export const ROLE_MODERATOR = 'moderator'
export const ROLE_USER = 'user'

/** 角色中文映射（与 admin-sidebar 的 roleLabels 同源，收编于此统一复用） */
const ROLE_LABELS: Record<string, string> = {
  [ROLE_ADMIN]: '管理员',
  [ROLE_MODERATOR]: '协作者',
  [ROLE_USER]: '用户',
}

/**
 * 判断用户是否为管理员。接受显式传入的 user，缺省时读取本地缓存的当前用户，
 * 便于在 beforeLoad 等同步上下文中无 hook 时调用。
 */
export function isAdmin(user?: UserInfo | null): boolean {
  const target = user ?? getStoredUser()
  return target?.role === ROLE_ADMIN
}

/**
 * 判断用户是否持有指定角色之一。缺省 user 时回退读取本地缓存用户。
 */
export function hasRole(user: UserInfo | null | undefined, ...roles: string[]): boolean {
  const target = user ?? getStoredUser()
  if (!target) return false
  return roles.includes(target.role)
}

/**
 * 取角色中文标签；未知角色回退为原始 role 字符串，再回退为占位横线。
 */
export function roleLabel(role?: string | null): string {
  if (!role) return '—'
  return ROLE_LABELS[role] ?? role
}
