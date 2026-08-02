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

import { Link } from '@tanstack/react-router'
import { AnimatePresence, motion, useReducedMotion } from 'motion/react'
import { useState } from 'react'
import { HoverCard as HoverCardPrimitive } from 'radix-ui'
import {
  ChevronDown,
  CircleUserRound,
  LayoutDashboard,
  LogIn,
  LogOut,
} from 'lucide-react'
import type { LucideIcon } from 'lucide-react'
import { useAuth } from '@/hooks/use-auth'
import { formatDateTime } from '@/lib/datetime'
import { isAdmin, roleLabel } from '@/lib/role'

/** S 型墨韵缓动：与首页 header/footer 对焦凝定同节奏 */
const EASE_INK: [number, number, number, number] = [0.76, 0, 0.24, 1]

/**
 * 公开页导航右侧的账户入口。
 * - 未登录：墨韵胶囊「登录」按钮（LogIn 图标），悬停图标前倾；
 * - 已登录：人像图标 + 昵称，悬停展开状态卡片（同管理侧栏用户卡语言），
 *   入口按角色分发——管理员可进管理后台与用户中心（用户中心含个人信息
 *   修改，向下兼容），普通用户仅用户中心。
 */
export function AccountHoverCard() {
  const { user, isAuthenticated, signOut } = useAuth()
  const reduced = useReducedMotion() ?? false
  const [open, setOpen] = useState(false)

  // ───────── 未登录：墨韵胶囊登录按钮 ─────────
  if (!isAuthenticated) {
    return (
      <Link
        to="/auth/login"
        className="group flex items-center gap-2 rounded-full border border-leaf-muted/45 bg-leaf-light/20 px-4 py-2 transition-colors duration-300 hover:border-leaf-deep/40 hover:bg-leaf-light/40"
      >
        <LogIn className="size-3.5 text-leaf-deep transition-transform duration-300 ease-out group-hover:translate-x-0.5" />
        <span className="font-mono text-[11px] uppercase tracking-[0.28em] text-text-primary transition-colors duration-300 group-hover:text-leaf-deep">
          登录
        </span>
      </Link>
    )
  }

  const displayName = user?.nickname || user?.username || '—'
  const userInitial = (user?.username ?? '?').charAt(0).toUpperCase()

  // 角色分发：admin 双入口（向下兼容用户中心），普通用户仅用户中心
  const entries: Array<{ to: string; label: string; icon: LucideIcon }> = [
    ...(isAdmin(user)
      ? [{ to: '/admin/dashboard', label: '管理后台', icon: LayoutDashboard }]
      : []),
    { to: '/user/dashboard', label: '用户中心', icon: CircleUserRound },
  ]

  const handleLogout = async () => {
    setOpen(false)
    await signOut()
    window.location.href = '/auth/login'
  }

  return (
    <HoverCardPrimitive.Root
      open={open}
      onOpenChange={setOpen}
      openDelay={100}
      closeDelay={160}
    >
      <HoverCardPrimitive.Trigger asChild>
        <button
          type="button"
          className="group flex items-center gap-2 rounded-full border border-transparent py-1.5 pl-2.5 pr-3 transition-colors duration-300 hover:border-leaf-muted/50 hover:bg-leaf-light/30 data-[state=open]:border-leaf-muted/50 data-[state=open]:bg-leaf-light/35"
        >
          <CircleUserRound className="size-[18px] text-text-secondary transition-colors duration-300 group-hover:text-leaf-deep group-data-[state=open]:text-leaf-deep" />
          <span className="max-w-28 truncate font-serif text-[13px] font-semibold tracking-wide text-text-primary">
            {displayName}
          </span>
          <motion.span
            aria-hidden
            className="flex"
            animate={{ rotate: open ? 180 : 0 }}
            transition={
              reduced ? { duration: 0 } : { duration: 0.3, ease: EASE_INK }
            }
          >
            <ChevronDown className="size-3.5 text-text-secondary/60 transition-colors duration-300 group-hover:text-text-secondary" />
          </motion.span>
        </button>
      </HoverCardPrimitive.Trigger>

      <AnimatePresence>
        {open && (
          <HoverCardPrimitive.Content
            asChild
            forceMount
            side="bottom"
            align="end"
            sideOffset={12}
            collisionPadding={12}
          >
            <motion.div
              className="z-50 w-72 overflow-hidden rounded-xl border border-border bg-card shadow-[0_18px_44px_-18px_oklch(0.32_0.06_155/0.32)]"
              initial={
                reduced
                  ? { opacity: 0 }
                  : { opacity: 0, y: -8, scale: 0.98, filter: 'blur(8px)' }
              }
              animate={
                reduced
                  ? { opacity: 1, transition: { duration: 0.15 } }
                  : {
                      opacity: 1,
                      y: 0,
                      scale: 1,
                      filter: 'blur(0px)',
                      transition: { duration: 0.38, ease: EASE_INK },
                    }
              }
              exit={
                reduced
                  ? { opacity: 0, transition: { duration: 0.12 } }
                  : {
                      opacity: 0,
                      y: -6,
                      scale: 0.99,
                      filter: 'blur(6px)',
                      transition: { duration: 0.24, ease: EASE_INK },
                    }
              }
            >
              {/* 状态信息栏：晨光墨晕 + 首字头像 + 角色徽章 */}
              <div className="relative overflow-hidden px-4 pt-4 pb-3.5">
                <div
                  className="pointer-events-none absolute inset-0"
                  style={{
                    background:
                      'radial-gradient(220px 130px at 88% -10%, oklch(0.88 0.1 105 / 0.2), transparent 70%)',
                  }}
                />
                <div className="relative flex items-center gap-3">
                  <div className="flex size-11 shrink-0 items-center justify-center rounded-lg bg-leaf-deep/12 font-serif text-lg font-semibold text-leaf-deep ring-1 ring-leaf-deep/15">
                    {userInitial}
                  </div>
                  <div className="min-w-0 flex-1">
                    <div className="flex items-center gap-2">
                      <span className="truncate font-serif text-[15px] font-semibold text-text-primary">
                        {displayName}
                      </span>
                      <span className="shrink-0 rounded border border-leaf-muted/50 bg-leaf-light/30 px-1.5 py-px font-mono text-[10px] tracking-wide text-leaf-deep">
                        {roleLabel(user?.role)}
                      </span>
                    </div>
                    <p className="truncate font-mono text-[11px] text-text-secondary">
                      {user?.email ?? ''}
                    </p>
                  </div>
                </div>
                <dl className="relative mt-3.5">
                  <StatusRow
                    label="邮箱认证"
                    value={user?.email_verify ? '已认证' : '未认证'}
                  />
                  <StatusRow
                    label="上次登录"
                    value={formatDateTime(user?.last_login_at)}
                  />
                  <StatusRow
                    label="注册于"
                    value={formatDateTime(user?.created_at)}
                  />
                </dl>
              </div>

              {/* 角色感知入口：悬停墨条签名与管理侧栏同源 */}
              <div className="border-t border-border p-1.5">
                {entries.map((entry) => (
                  <Link
                    key={entry.to}
                    to={entry.to}
                    onClick={() => setOpen(false)}
                    className="relative flex items-center gap-2.5 rounded-sm px-3 py-2 font-serif text-[13px] font-medium text-text-primary transition-colors duration-200 before:absolute before:top-1/2 before:left-0 before:h-3.5 before:w-[3px] before:-translate-y-1/2 before:rounded-sm before:bg-leaf-deep before:opacity-0 before:transition-opacity hover:bg-leaf-light/40 hover:text-leaf-deep hover:before:opacity-100"
                  >
                    <entry.icon className="size-4 text-text-secondary transition-colors duration-200" />
                    {entry.label}
                  </Link>
                ))}
              </div>

              <div className="border-t border-border p-1.5">
                <button
                  type="button"
                  onClick={handleLogout}
                  className="relative flex w-full items-center gap-2.5 rounded-sm px-3 py-2 font-serif text-[13px] font-medium text-destructive transition-colors duration-200 before:absolute before:top-1/2 before:left-0 before:h-3.5 before:w-[3px] before:-translate-y-1/2 before:rounded-sm before:bg-destructive before:opacity-0 before:transition-opacity hover:bg-destructive/10 hover:before:opacity-100"
                >
                  <LogOut className="size-4" />
                  退出登录
                </button>
              </div>
            </motion.div>
          </HoverCardPrimitive.Content>
        )}
      </AnimatePresence>
    </HoverCardPrimitive.Root>
  )
}

/* ───────── StatusRow：状态信息行（label + mono 数值） ───────── */

function StatusRow({ label, value }: { label: string; value: string }) {
  return (
    <div className="flex items-baseline justify-between gap-3 border-b border-border py-1.5 last:border-b-0">
      <dt className="text-[11px] text-text-secondary">{label}</dt>
      <dd className="truncate font-mono text-[11px] font-medium tabular-nums text-text-primary">
        {value}
      </dd>
    </div>
  )
}
