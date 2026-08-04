// --------------------------------------------------------------------------------
// Copyright (c) 2016-NOW(至今) 筱锋
// Author: 筱锋「xiao_lfeng」(https://www.x-lf.com)
// --------------------------------------------------------------------------------
// 许可证声明：版权所有 (c) 2016-2026 筱锋。保留所有权利。
// 有关MIT许可证的更多信息，请查看项目根目录下的LICENSE文件或访问：
// https://opensource.org/licenses/MIT
// --------------------------------------------------------------------------------

import { useEffect, useState } from 'react'
import {
  Link,
  Outlet,
  createFileRoute,
  redirect,
  useLocation,
} from '@tanstack/react-router'
import { motion, useReducedMotion } from 'motion/react'
import { ChevronsUpDown, LayoutDashboard, LogOut } from 'lucide-react'
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from '@/components/ui/dropdown-menu'
import { useAuth } from '@/hooks/use-auth'
import { getToken } from '@/lib/auth'
import { isAdmin, roleLabel as resolveRoleLabel } from '@/lib/role'
import { cn } from '@/lib/utils'
import { enter } from '@/lib/motion'

export const Route = createFileRoute('/_user')({
  // 路由守卫：无登录令牌时回跳登录页，并携带来源路径以便登录后返回
  beforeLoad: ({ location }) => {
    if (!getToken()) {
      throw redirect({
        to: '/auth/login',
        search: { redirect: location.pathname },
      })
    }
  },
  component: UserLayout,
})

/** 顶部导航项（居中展示，激活态用笔刷下划线） */
interface NavItem {
  to: string
  label: string
  /** 精确匹配之外的前缀匹配（子路由高亮，如 /user/links/* → 我的友链） */
  matchPrefix?: string
}

const NAV_ITEMS: Array<NavItem> = [
  { to: '/user/dashboard', label: '看板' },
  { to: '/user/links', label: '我的友链', matchPrefix: '/user/links' },
  { to: '/user/account', label: '账号设置' },
]

/** 竹叶小标（导航品牌用） */
function BambooLeafMark() {
  return (
    <svg
      width="20"
      height="20"
      viewBox="0 0 48 32"
      fill="var(--leaf-deep)"
      aria-hidden
    >
      <path d="M2 30C10 18 26 6 46 2c-3 12-16 24-44 28z" />
    </svg>
  )
}

/** 日期时间格式化为 2026-07-28 19:30；无效/缺省返回 — */
function formatDateTime(iso?: string | null): string {
  if (!iso) return '—'
  const d = new Date(iso)
  if (Number.isNaN(d.getTime())) return '—'
  const pad = (n: number) => n.toString().padStart(2, '0')
  return `${d.getFullYear()}-${pad(d.getMonth() + 1)}-${pad(d.getDate())} ${pad(d.getHours())}:${pad(d.getMinutes())}`
}

/** 状态信息行（label + mono 数值） */
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

/**
 * 用户中心顶导头像下拉：用户名片 + 状态行 + 退出登录。
 * 由 user-sidebar 收编而来，结构身份从「侧边栏 footer」迁为「顶导右侧下拉」。
 */
function UserAvatarMenu() {
  const { user, signOut } = useAuth()
  const [menuOpen, setMenuOpen] = useState(false)
  const displayName = user?.nickname || user?.username || '用户'
  const userInitial = (user?.username ?? '?').charAt(0).toUpperCase()
  const roleLabel = resolveRoleLabel(user?.role)

  const handleLogout = async () => {
    await signOut()
    window.location.href = '/auth/login'
  }

  return (
    <DropdownMenu open={menuOpen} onOpenChange={setMenuOpen}>
      <DropdownMenuTrigger asChild>
        <button
          type="button"
          className="flex cursor-pointer items-center gap-1.5 rounded-full border-2 border-leaf-light/60 bg-leaf-deep/8 p-0 pl-0 transition-colors hover:border-leaf-muted hover:bg-leaf-deep/12"
          aria-label="用户菜单"
        >
          <span className="grid size-8 place-items-center rounded-full bg-leaf-deep/15 font-serif text-sm font-semibold text-leaf-deep">
            {userInitial}
          </span>
          <ChevronsUpDown className="mr-1 size-3.5 text-text-secondary/60" />
        </button>
      </DropdownMenuTrigger>
      <DropdownMenuContent
        align="end"
        sideOffset={10}
        className="min-w-64 rounded-xl border-border bg-card p-0 shadow-[0_18px_44px_-18px_oklch(0.32_0.06_155/0.32)]"
      >
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
                  {roleLabel}
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
        {isAdmin(user) && (
          <>
            <DropdownMenuSeparator className="bg-border" />
            <div className="p-1.5">
              <Link
                to="/admin/dashboard"
                onClick={() => setMenuOpen(false)}
                className="relative flex items-center gap-2.5 rounded-sm px-3 py-2 font-serif text-[13px] font-medium text-text-primary transition-colors duration-200 before:absolute before:top-1/2 before:left-0 before:h-3.5 before:w-[3px] before:-translate-y-1/2 before:rounded-sm before:bg-leaf-deep before:opacity-0 before:transition-opacity hover:bg-leaf-light/40 hover:text-leaf-deep hover:before:opacity-100"
              >
                <LayoutDashboard className="size-4 text-text-secondary transition-colors duration-200" />
                切换到管理员端
              </Link>
            </div>
          </>
        )}
        <DropdownMenuSeparator className="bg-border" />
        <div className="p-1.5">
          <DropdownMenuItem variant="destructive" onClick={handleLogout}>
            <LogOut className="size-4" />
            退出登录
          </DropdownMenuItem>
        </div>
      </DropdownMenuContent>
    </DropdownMenu>
  )
}

function UserLayout() {
  const reduced = useReducedMotion() ?? false
  const thisYear = new Date().getFullYear()
  const pathname = useLocation({ select: (loc) => loc.pathname })
  const [scrolled, setScrolled] = useState(false)

  useEffect(() => {
    const onScroll = () => setScrolled(window.scrollY > 40)
    onScroll()
    addEventListener('scroll', onScroll, { passive: true })
    return () => removeEventListener('scroll', onScroll)
  }, [])

  return (
    <div className="relative min-h-dvh bg-background">
      {/* 顶部导航：纯宣纸底（区别于 about 的模糊背景图），功能化个人门户 */}
      <nav
        className={cn(
          'fixed inset-x-0 top-0 z-50 border-b transition-colors duration-400',
          scrolled
            ? 'border-border bg-background/94 backdrop-blur-[2px]'
            : 'border-transparent',
        )}
      >
        <div className="mx-auto flex w-full max-w-5xl items-center px-6 py-4 md:px-10">
          {/* 左：品牌 */}
          <Link
            to="/"
            className="flex shrink-0 items-center gap-2.5 transition-opacity hover:opacity-70"
          >
            <BambooLeafMark />
            <span className="font-serif text-base font-semibold tracking-wide text-text-primary">
              用户中心
            </span>
          </Link>

          {/* 中：导航项（居中） */}
          <div className="flex flex-1 items-center justify-center gap-7 md:gap-9">
            {NAV_ITEMS.map((item) => {
              const active = item.matchPrefix
                ? pathname === item.to || pathname.startsWith(item.matchPrefix)
                : pathname === item.to
              return (
                <Link
                  key={item.to}
                  to={item.to}
                  className={cn(
                    'relative font-mono text-[11px] uppercase tracking-[0.28em] transition-colors',
                    'after:absolute after:-bottom-[7px] after:left-0 after:h-0.5 after:w-full after:origin-left after:bg-leaf-deep after:transition-transform after:duration-300 after:ease-[cubic-bezier(0.22,0.9,0.3,1)]',
                    active
                      ? 'text-text-primary after:scale-x-100'
                      : 'text-text-secondary after:scale-x-0 hover:text-text-primary hover:after:scale-x-100',
                  )}
                >
                  {item.label}
                </Link>
              )
            })}
          </div>

          {/* 右：头像下拉 */}
          <div className="flex shrink-0 items-center">
            <UserAvatarMenu />
          </div>
        </div>
      </nav>

      {/* 子路由内容：每页自带开场与宽度约束 */}
      <motion.div
        key={pathname}
        {...enter(reduced, 0, {
          initial: { opacity: 0 },
          animate: { opacity: 1 },
          transition: { duration: reduced ? 0.15 : 0.3 },
        })}
        className="relative z-10 pt-20"
      >
        <Outlet />
      </motion.div>

      {/* 页脚 */}
      <footer className="relative z-10 mx-auto w-full max-w-5xl px-6 py-14 md:px-10">
        <div className="flex flex-col items-center gap-4 text-center">
          <BambooLeafMark />
          <p className="font-mono text-[11px] tracking-[0.2em] text-text-secondary">
            Copyright (C) 2016-{thisYear} 筱锋 xiao_lfeng · All Rights Reserved
          </p>
        </div>
      </footer>
    </div>
  )
}
