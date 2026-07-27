// --------------------------------------------------------------------------------
// Copyright (c) 2016-NOW(至今) 筱锋
// Author: 筱锋「xiao_lfeng」(https://www.x-lf.com)
// --------------------------------------------------------------------------------
// 许可证声明：版权所有 (c) 2016-2026 筱锋。保留所有权利。
// 有关MIT许可证的更多信息，请查看项目根目录下的LICENSE文件或访问：
// https://opensource.org/licenses/MIT
// --------------------------------------------------------------------------------

import type { UserInfo } from '@/api/types'

const TOKEN_KEY = 'bamboo_auth_token'
const USER_KEY = 'bamboo_auth_user'

/** 匹配 id 或 *_id 字段名（雪花 ID），用于本地存储的 bigint 往返 */
const ID_KEY_PATTERN = /^(id|.*_id)$/
const NUMERIC_PATTERN = /^\d+$/

/** bigint 安全序列化：将 bigint 转为字符串（与 client.ts 请求侧一致） */
function stringifyUser(user: UserInfo): string {
  return JSON.stringify(user, (_key, value) =>
    typeof value === 'bigint' ? value.toString() : value,
  )
}

/** bigint 安全反序列化：将 id / *_id 的纯数字字符串还原为 bigint */
function parseUser(raw: string): UserInfo | null {
  return JSON.parse(raw, (key, value) => {
    if (
      typeof value === 'string' &&
      ID_KEY_PATTERN.test(key) &&
      NUMERIC_PATTERN.test(value)
    ) {
      return BigInt(value)
    }
    return value
  }) as UserInfo
}

/**
 * 写入登录会话。
 * - remember=true：localStorage（关闭浏览器后保留）
 * - remember=false：sessionStorage（关闭浏览器即清除）
 * 写入前先清空两种存储，避免残留旧会话造成读取歧义。
 * 用户对象的雪花 ID 以字符串存储，读取时还原为 bigint，保证类型一致。
 */
export function setSession(
  token: string,
  user: UserInfo,
  remember: boolean,
): void {
  clearSession()
  const storage = remember ? window.localStorage : window.sessionStorage
  storage.setItem(TOKEN_KEY, token)
  storage.setItem(USER_KEY, stringifyUser(user))
}

/** 读取当前访问令牌；两种存储都查，优先 localStorage。 */
export function getToken(): string | null {
  return (
    window.localStorage.getItem(TOKEN_KEY) ??
    window.sessionStorage.getItem(TOKEN_KEY)
  )
}

/** 读取本地缓存的用户信息（可能滞后于服务端，仅作展示兜底）。 */
export function getStoredUser(): UserInfo | null {
  const raw =
    window.localStorage.getItem(USER_KEY) ??
    window.sessionStorage.getItem(USER_KEY)
  if (!raw) return null
  try {
    return parseUser(raw)
  } catch {
    return null
  }
}

/** 清空登录会话（两种存储都清）。 */
export function clearSession(): void {
  window.localStorage.removeItem(TOKEN_KEY)
  window.localStorage.removeItem(USER_KEY)
  window.sessionStorage.removeItem(TOKEN_KEY)
  window.sessionStorage.removeItem(USER_KEY)
}
