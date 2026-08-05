// --------------------------------------------------------------------------------
// Copyright (c) 2016-NOW(至今) 筱锋
// Author: 筱锋「xiao_lfeng」(https://www.x-lf.com)
// --------------------------------------------------------------------------------
// 许可证声明：版权所有 (c) 2016-2026 筱锋。保留所有权利。
// 有关MIT许可证的更多信息，请查看项目根目录下的LICENSE文件或访问：
// https://opensource.org/licenses/MIT
// --------------------------------------------------------------------------------

import { Link, createFileRoute } from '@tanstack/react-router'
import { motion, useReducedMotion } from 'motion/react'
import { ArrowRight, Heart, Plus, Settings } from 'lucide-react'
import type { InkBadgeTone } from '@/components/ink-wash'
import { useAuth } from '@/hooks/use-auth'
import { useMyLinks } from '@/hooks/use-links'
import {
  BrushUnderline,
  CardHead,
  InkBadge,
  inkCard,
} from '@/components/ink-wash'
import { Skeleton } from '@/components/ui/skeleton'
import { enter } from '@/lib/motion'

export const Route = createFileRoute('/_user/user/dashboard')({
  component: UserDashboardPage,
})

/** 按时段问候语 */
function greeting(): string {
  const hour = new Date().getHours()
  if (hour < 6) return '夜深了'
  if (hour < 12) return '早上好'
  if (hour < 14) return '中午好'
  if (hour < 18) return '下午好'
  return '晚上好'
}

/** 状态行配置：色点 + 标签 + 徽章 + 计数（色点一律走 styles.css token） */
interface StatusRow {
  label: string
  badge: string
  tone: InkBadgeTone
  dot: string
  count: number
}

function UserDashboardPage() {
  const reduced = useReducedMotion() ?? false
  const { user } = useAuth()
  const { data, isLoading } = useMyLinks({ page: 1, page_size: 100 })
  const links = data?.data ?? []

  // 本地聚合：待审核（含下架待审核）/ 已展示 / 被拒 / 已下架
  const pending = links.filter((l) => l.status === 0 || l.status === 3).length
  const live = links.filter((l) => l.status === 1).length
  const rejected = links.filter((l) => l.status === 2).length
  const offline = links.filter((l) => l.status === 4).length

  const displayName = user?.nickname || user?.username

  const statusRows: Array<StatusRow> = [
    {
      label: '待审核',
      badge: '审核中',
      tone: 'pending',
      dot: 'var(--leaf-light)',
      count: pending,
    },
    {
      label: '已展示',
      badge: '公开展示',
      tone: 'leaf',
      dot: 'var(--chart-1)',
      count: live,
    },
    {
      label: '被拒',
      badge: '未通过',
      tone: 'danger',
      dot: 'oklch(0.577 0.245 27.325 / 0.5)',
      count: rejected,
    },
    {
      label: '已下架',
      badge: '不可见',
      tone: 'neutral',
      dot: 'var(--leaf-muted)',
      count: offline,
    },
  ]

  return (
    <div className="mx-auto max-w-4xl">
      {/* ═══════════ 轻盈 hero：kicker + 问候 + 笔刷 + 导语 ═══════════ */}
      <motion.section
        {...enter(reduced, 0.05, {
          initial: { opacity: 0, y: 12 },
          animate: { opacity: 1, y: 0 },
          transition: { duration: 0.6, ease: [0.22, 0.9, 0.3, 1] as const },
        })}
        className="mb-8"
      >
        <p className="flex items-center gap-4 font-mono text-[11px] uppercase tracking-[0.35em] text-text-secondary">
          <span className="h-px w-10 bg-leaf-deep" />
          我的友链 · 各美其美
        </p>
        <h1 className="mt-6 font-serif text-[clamp(2rem,5vw,3rem)] font-bold leading-[1.08] tracking-[0.02em] text-text-primary">
          {greeting()}，
          <span className="text-leaf-deep">{displayName ?? '朋友'}</span>
        </h1>
        <div className="mt-4">
          <BrushUnderline />
        </div>
        <p className="mt-4 max-w-xl font-serif text-base italic leading-relaxed text-text-secondary">
          竹林清晨，万物生长。你的友链如竹节，一节一节，缓缓而成。
        </p>
      </motion.section>

      {/* ═══════════ 友链状态概览 ═══════════ */}
      <motion.section
        {...enter(reduced, 0.15, {
          initial: { opacity: 0, y: 12 },
          animate: { opacity: 1, y: 0 },
          transition: { duration: 0.5, ease: 'easeOut' },
        })}
        className={`${inkCard} mb-6`}
      >
        <CardHead
          title="友链状态"
          meta={isLoading ? '—' : `TOTAL ${links.length}`}
        />
        {isLoading ? (
          <div className="space-y-3">
            <Skeleton className="h-12 w-full" />
            <Skeleton className="h-12 w-full" />
          </div>
        ) : (
          <div>
            {statusRows.map((row, i) => (
              <div
                key={row.label}
                className={`flex items-center justify-between py-3.5 ${
                  i < statusRows.length - 1 ? 'border-b border-border' : ''
                }`}
              >
                <div className="flex items-center gap-3">
                  <span
                    className="size-2 rounded-[2px]"
                    style={{ backgroundColor: row.dot }}
                    aria-hidden
                  />
                  <span className="text-sm text-text-primary">{row.label}</span>
                  <InkBadge tone={row.tone}>{row.badge}</InkBadge>
                </div>
                <span className="font-serif text-2xl font-semibold tabular-nums text-text-primary">
                  {row.count}
                </span>
              </div>
            ))}
          </div>
        )}
      </motion.section>

      {/* ═══════════ 快捷入口 ═══════════ */}
      <motion.section
        {...enter(reduced, 0.25, {
          initial: { opacity: 0, y: 12 },
          animate: { opacity: 1, y: 0 },
          transition: { duration: 0.5, ease: 'easeOut' },
        })}
        className="grid gap-4 sm:grid-cols-2 xl:grid-cols-3"
      >
        <Link
          to="/operate/apply"
          className={`${inkCard} group flex items-center gap-4 p-5`}
        >
          <span className="grid size-11 shrink-0 place-items-center rounded-lg bg-leaf-deep/12 text-leaf-deep ring-1 ring-leaf-deep/15 transition-transform duration-300 group-hover:scale-105">
            <Plus className="size-5" />
          </span>
          <span className="min-w-0 flex-1">
            <span className="block font-serif text-[15px] font-semibold text-text-primary">
              申请友链
            </span>
            <span className="mt-0.5 block text-xs text-text-secondary">
              提交新的小站，让竹林认识你
            </span>
          </span>
          <ArrowRight className="size-4 shrink-0 text-text-secondary transition-transform duration-300 group-hover:translate-x-0.5" />
        </Link>
        <Link
          to="/user/account"
          className={`${inkCard} group flex items-center gap-4 p-5`}
        >
          <span className="grid size-11 shrink-0 place-items-center rounded-lg bg-leaf-deep/12 text-leaf-deep ring-1 ring-leaf-deep/15 transition-transform duration-300 group-hover:scale-105">
            <Settings className="size-5" />
          </span>
          <span className="min-w-0 flex-1">
            <span className="block font-serif text-[15px] font-semibold text-text-primary">
              账号设置
            </span>
            <span className="mt-0.5 block text-xs text-text-secondary">
              管理昵称、头像与密码
            </span>
          </span>
          <ArrowRight className="size-4 shrink-0 text-text-secondary transition-transform duration-300 group-hover:translate-x-0.5" />
        </Link>
        <Link
          to="/operate/sponsor"
          className={`${inkCard} group flex items-center gap-4 p-5`}
        >
          <span className="grid size-11 shrink-0 place-items-center rounded-lg bg-leaf-deep/12 text-leaf-deep ring-1 ring-leaf-deep/15 transition-transform duration-300 group-hover:scale-105">
            <Heart className="size-5" />
          </span>
          <span className="min-w-0 flex-1">
            <span className="block font-serif text-[15px] font-semibold text-text-primary">
              申请赞助展示
            </span>
            <span className="mt-0.5 block text-xs text-text-secondary">
              递交赞助心意，经核实后入册
            </span>
          </span>
          <ArrowRight className="size-4 shrink-0 text-text-secondary transition-transform duration-300 group-hover:translate-x-0.5" />
        </Link>
      </motion.section>
    </div>
  )
}
