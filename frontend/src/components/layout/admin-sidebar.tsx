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

import { Link, useLocation } from '@tanstack/react-router'
import { motion, useReducedMotion } from 'motion/react'
import {
  ChevronsUpDown,
  Heart,
  LayoutDashboard,
  Link as LinkIcon,
  LogOut,
  MapPin,
  Palette,
  Settings,
} from 'lucide-react'
import {
  Sidebar,
  SidebarContent,
  SidebarFooter,
  SidebarGroup,
  SidebarGroupContent,
  SidebarGroupLabel,
  SidebarHeader,
  SidebarMenu,
  SidebarMenuButton,
  SidebarMenuItem,
} from '@/components/ui/sidebar'
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from '@/components/ui/dropdown-menu'
import { BambooLogo } from '@/assets/svg/bamboo-logo'
import { siteConfig } from '@/lib/site'
import { useAuth } from '@/hooks/use-auth'
import { useSiteInfo } from '@/hooks/use-site-info'
import { cn } from '@/lib/utils'

/** 角色中文映射（与后端 pkg/constants 的 Role* 常量对应） */
const roleLabels: Record<string, string> = {
  admin: '管理员',
  moderator: '协作者',
  user: '用户',
}

/** 日期时间格式化为 2026-07-28 19:30；无效/缺省返回 — */
function formatDateTime(iso?: string | null): string {
  if (!iso) return '—'
  const d = new Date(iso)
  if (Number.isNaN(d.getTime())) return '—'
  const pad = (n: number) => n.toString().padStart(2, '0')
  return `${d.getFullYear()}-${pad(d.getMonth() + 1)}-${pad(d.getDate())} ${pad(d.getHours())}:${pad(d.getMinutes())}`
}

const menuGroups = [
  {
    label: '首页',
    items: [{ title: '看板', url: '/admin/dashboard', icon: LayoutDashboard }],
  },
  {
    label: '友链',
    items: [
      { title: '友链管理', url: '/admin/link', icon: LinkIcon },
      { title: '位置管理', url: '/admin/location', icon: MapPin },
      { title: '颜色管理', url: '/admin/color', icon: Palette },
    ],
  },
  {
    label: '其他',
    items: [
      { title: '赞助', url: '/admin/sponsor', icon: Heart },
      { title: '设置', url: '/admin/setting', icon: Settings },
    ],
  },
]

export function AdminSidebar() {
  const { pathname } = useLocation()
  const { user, signOut } = useAuth()
  const { data: site } = useSiteInfo()
  const reduced = useReducedMotion() ?? false
  const siteName = site?.site_name || siteConfig.defaultName

  const isActive = (url: string) => {
    if (url === '/admin/link') {
      return pathname === url || pathname.startsWith('/admin/link/')
    }
    return pathname === url
  }

  const displayName = user?.nickname || user?.username || '管理员'
  const userInitial = (user?.username ?? '?').charAt(0).toUpperCase()
  const roleLabel = roleLabels[user?.role ?? ''] ?? user?.role ?? '—'

  const handleLogout = async () => {
    await signOut()
    window.location.href = '/auth/login'
  }

  return (
    <Sidebar variant="floating" collapsible="icon">
      {/* ───────── Header：墨晕 logo + 衬线站名 + mono 版本 + 竹节线 ───────── */}
      <SidebarHeader>
        <SidebarMenu>
          <SidebarMenuItem>
            <SidebarMenuButton size="lg" asChild>
              <Link to="/">
                {/* 晨光墨晕底圈，承托竹 logo */}
                <div className="flex aspect-square size-9 items-center justify-center rounded-lg bg-gradient-to-br from-leaf-light/70 via-leaf-light/25 to-transparent ring-1 ring-leaf-muted/40">
                  <BambooLogo size={28} />
                </div>
                <div className="grid flex-1 text-left leading-tight">
                  <span className="truncate font-serif text-[15px] font-semibold tracking-wide text-text-primary">
                    {siteName}
                  </span>
                  <span className="truncate font-mono text-[10px] uppercase tracking-[0.16em] text-text-secondary">
                    v{siteConfig.version} · {siteConfig.tagline}
                  </span>
                </div>
              </Link>
            </SidebarMenuButton>
          </SidebarMenuItem>
        </SidebarMenu>
        <BambooRule className="mt-1" />
      </SidebarHeader>

      {/* ───────── Content：斜墨条分组签名 + 激活墨条菜单 ───────── */}
      <SidebarContent>
        {menuGroups.map((group) => (
          <SidebarGroup key={group.label}>
            {/* 分组签名：斜墨条 + 衬线（收起态自动隐藏） */}
            <SidebarGroupLabel className="gap-2 font-serif text-xs font-semibold tracking-wide text-text-secondary">
              <span className="h-[3px] w-3 -skew-x-12 rounded-sm bg-leaf-muted" />
              {group.label}
            </SidebarGroupLabel>
            <SidebarGroupContent>
              <SidebarMenu>
                {group.items.map((item) => {
                  const active = isActive(item.url)
                  return (
                    <SidebarMenuItem key={item.title}>
                      {/* 激活墨条：layoutId 跨菜单项滑动；置于按钮外侧，避开按钮 overflow-hidden 裁切 */}
                      {active && (
                        <motion.span
                          layoutId="sidebar-ink-bar"
                          transition={
                            reduced
                              ? { duration: 0 }
                              : { type: 'spring', stiffness: 380, damping: 30 }
                          }
                          className="absolute top-2 left-1 z-10 h-4 w-[3px] group-data-[collapsible=icon]:hidden"
                          aria-hidden
                        >
                          <span className="block h-full w-full -skew-x-12 rounded-sm bg-leaf-deep" />
                        </motion.span>
                      )}
                      <SidebarMenuButton
                        asChild
                        isActive={active}
                        tooltip={item.title}
                        className={cn(
                          // 激活才 pl-4 给墨条让位 + 墨色；未激活保持默认 p-2 原状
                          'transition-[width,height,padding,background-color,color] duration-300',
                          'data-[active=true]:pl-4 data-[active=true]:bg-leaf-light/45 data-[active=true]:font-semibold data-[active=true]:text-leaf-deep',
                        )}
                      >
                        <Link to={item.url}>
                          <item.icon className="h-4 w-4" />
                          <span>{item.title}</span>
                        </Link>
                      </SidebarMenuButton>
                    </SidebarMenuItem>
                  )
                })}
              </SidebarMenu>
            </SidebarGroupContent>
          </SidebarGroup>
        ))}
      </SidebarContent>

      {/* ───────── Footer：竹节线 + 用户状态菜单（点击展开）+ 墨韵竹叶点缀 ───────── */}
      <SidebarFooter className="relative">
        {/* 极淡墨韵竹叶，贴右下摇曳（收起态隐藏） */}
        <SidebarBambooArt className="pointer-events-none absolute right-1 bottom-1 h-24 w-16 text-leaf-deep opacity-[0.06] group-data-[collapsible=icon]:hidden" />
        <BambooRule className="mb-1" />
        <SidebarMenu>
          {/* 点击用户卡片展开状态信息栏 + 退出登录 */}
          <SidebarMenuItem>
            <DropdownMenu>
              <DropdownMenuTrigger asChild>
                <SidebarMenuButton size="lg" className="transition-colors">
                  <div className="flex aspect-square size-8 items-center justify-center rounded-lg bg-leaf-deep/12 font-serif text-sm font-semibold text-leaf-deep ring-1 ring-leaf-deep/15">
                    {userInitial}
                  </div>
                  <div className="grid flex-1 text-left leading-tight">
                    <span className="truncate font-serif text-sm font-semibold text-text-primary">
                      {displayName}
                    </span>
                    <span className="truncate font-mono text-[10px] tracking-wide text-text-secondary">
                      {user?.email ?? ''}
                    </span>
                  </div>
                  <ChevronsUpDown className="ml-auto size-4 shrink-0 text-text-secondary/60" />
                </SidebarMenuButton>
              </DropdownMenuTrigger>
              <DropdownMenuContent
                side="top"
                align="start"
                sideOffset={8}
                className="w-(--radix-dropdown-menu-trigger-width) min-w-60 rounded-xl border-border bg-card p-0 shadow-[0_18px_44px_-18px_oklch(0.32_0.06_155/0.32)]"
              >
                {/* 状态信息栏 */}
                <div className="relative overflow-hidden px-4 pt-4 pb-3.5">
                  {/* 晨光墨晕 */}
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
                <DropdownMenuSeparator className="bg-border" />
                <div className="p-1.5">
                  <DropdownMenuItem variant="destructive" onClick={handleLogout}>
                    <LogOut className="size-4" />
                    退出登录
                  </DropdownMenuItem>
                </div>
              </DropdownMenuContent>
            </DropdownMenu>
          </SidebarMenuItem>
        </SidebarMenu>
      </SidebarFooter>
    </Sidebar>
  )
}

/* ───────── BambooRule：竹节分隔线（横线 + 节点，收起态隐藏） ───────── */

function BambooRule({ className }: { className?: string }) {
  return (
    <div
      className={cn(
        'flex items-center gap-2 px-3 group-data-[collapsible=icon]:hidden',
        className,
      )}
      aria-hidden
    >
      <span className="h-px flex-1 bg-border" />
      <span className="size-1.5 shrink-0 rounded-[3px] bg-leaf-muted opacity-80" />
      <span className="h-px flex-1 bg-border" />
      <span className="h-1 w-3 shrink-0 rounded bg-leaf-muted opacity-80" />
      <span className="h-px flex-1 bg-border" />
    </div>
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

/* ───────── SidebarBambooArt：单竿墨韵竹叶（footer 点缀，ink-sway 摇曳） ───────── */

function SidebarBambooArt({ className }: { className?: string }) {
  return (
    <svg
      className={className}
      viewBox="0 0 120 200"
      preserveAspectRatio="xMaxYMax meet"
      aria-hidden="true"
    >
      <defs>
        <path id="sbleaf" d="M0 0 C8 -5 22 -7 38 -2 C24 3 8 3 0 0 Z" />
      </defs>
      <g className="ink-sway">
        <path
          d="M82 200 C80 160 83 120 81 80 C80 50 82 22 81 4"
          fill="none"
          stroke="currentColor"
          strokeWidth="3"
          strokeLinecap="round"
        />
        <ellipse cx="81" cy="150" rx="4" ry="1.5" fill="none" stroke="currentColor" strokeWidth="1.2" />
        <ellipse cx="81" cy="96" rx="4" ry="1.5" fill="none" stroke="currentColor" strokeWidth="1.2" />
        <use href="#sbleaf" transform="translate(82 158) rotate(-22)" />
        <use href="#sbleaf" transform="translate(81 104) scale(-1 1) rotate(-26)" />
        <use href="#sbleaf" transform="translate(81 52) rotate(-32) scale(1.05)" />
      </g>
    </svg>
  )
}
