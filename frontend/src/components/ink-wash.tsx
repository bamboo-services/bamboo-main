// --------------------------------------------------------------------------------
// Copyright (c) 2016-NOW(至今) 筱锋
// Author: 筱锋「xiao_lfeng」(https://www.x-lf.com)
// --------------------------------------------------------------------------------
// 许可证声明：版权所有 (c) 2016-2026 筱锋。保留所有权利。
// 有关MIT许可证的更多信息，请查看项目根目录下的LICENSE文件或访问：
// https://opensource.org/licenses/MIT
// --------------------------------------------------------------------------------

import { Link } from '@tanstack/react-router'
import { motion } from 'motion/react'
import { ArrowLeft, ChevronRight } from 'lucide-react'
import type { ReactNode } from 'react'
import type { LinkFriend } from '@/api/types'
import { CountUp } from '@/components/dashboard/count-up'
import { Skeleton } from '@/components/ui/skeleton'
import { enter } from '@/lib/motion'
import { cn } from '@/lib/utils'

/**
 * 竹林水墨 · 共享视觉原语。
 *
 * 以 `dashboard.tsx` 为标杆沉淀的东方签名元素，供后台与公开页复用。
 * 色彩一律走 `styles.css` 的 CSS 变量（leaf-* / text-* / border），
 * 动画（ink-sway 等）定义在 `styles.css`，reduced-motion 时静止。
 */

/**
 * 宣纸卡片块：hover 浮起 + 墨边软阴影。自带 `group` 供子元素联动。
 *
 * 过渡刻意收窄为 translate/border-color/box-shadow 三项 hover 属性，
 * 不用 transition-all——本类会直接挂在带 motion 入场动画的 section 上
 * （dashboard 的友链构成/系统状态/最近申请三块），若 transition-all 把
 * transform/opacity 也纳入 CSS 过渡，会与 motion 入场动画叠加出「二次上滑」。
 * Tailwind v4 的 translate-y 走原生 translate 属性，与 motion 的 transform
 * 分属不同属性，收窄后 hover 浮起依旧平滑且互不干扰。
 */
export const inkCard =
  'group relative overflow-hidden rounded-lg border border-border bg-card p-5 transition-[translate,border-color,box-shadow] duration-300 hover:-translate-y-0.5 hover:border-leaf-muted hover:shadow-[0_14px_30px_-22px_oklch(0.32_0.06_155/0.4)]'

/** 卡片标题：斜墨条（leaf-deep）+ 衬线标题 + 右侧 mono uppercase meta */
export function CardHead({ title, meta }: { title: string; meta?: string }) {
  return (
    <div className="mb-3 flex items-baseline justify-between gap-2">
      <h3 className="flex items-center gap-2.5 font-serif text-base font-semibold text-text-primary">
        <span className="h-[3px] w-4 -skew-x-12 rounded-sm bg-leaf-deep" />
        {title}
      </h3>
      {meta && (
        <span className="font-mono text-[11px] uppercase tracking-widest text-text-secondary">
          {meta}
        </span>
      )}
    </div>
  )
}

/** 竹节分隔线：横线 + 圆/宽节点交替，用于分大章节 */
export function BambooRule({ reduced, delay }: { reduced: boolean; delay: number }) {
  return (
    <motion.div
      {...enter(reduced, delay, {
        initial: { opacity: 0 },
        animate: { opacity: 1 },
        transition: { duration: 0.4 },
      })}
      className="my-5 flex items-center gap-4"
      aria-hidden
    >
      <span className="h-px flex-1 bg-border" />
      <span className="size-2.5 shrink-0 rounded-[4px] bg-leaf-muted opacity-85" />
      <span className="h-px flex-1 bg-border" />
      <span className="h-1.5 w-4 shrink-0 rounded bg-leaf-muted opacity-85" />
      <span className="h-px flex-1 bg-border" />
      <span className="size-2.5 shrink-0 rounded-[4px] bg-leaf-muted opacity-85" />
      <span className="h-px flex-1 bg-border" />
    </motion.div>
  )
}

/** 笔刷下划线：手绘 SVG path（leaf-deep），用于标题 / 名帖收尾 */
export function BrushUnderline({ className = '' }: { className?: string }) {
  return (
    <svg className={`block h-1.5 w-21 ${className}`} viewBox="0 0 84 6" aria-hidden>
      <path d="M0 3 C 20 -0.5 44 -1.2 84 0.5 C 50 5 20 6 0 3 Z" fill="var(--leaf-deep)" />
    </svg>
  )
}

/** enso 缺口圆 + 竹叶：空状态签名图标 */
export function EnsoIcon({ className = '' }: { className?: string }) {
  return (
    <svg width="48" height="48" viewBox="0 0 48 48" aria-hidden className={`shrink-0 ${className}`}>
      <path
        d="M27 7 A 17 17 0 1 0 41 31"
        fill="none"
        stroke="var(--leaf-muted)"
        strokeWidth="2"
        strokeLinecap="round"
      />
      <path
        d="M16 22 C 19 19 25 18 31 21 C 27 23 20 23 16 22 Z"
        fill="var(--leaf-muted)"
        opacity="0.55"
      />
    </svg>
  )
}

/** 空状态：enso 缺口圆 + 衬线标题 + mono 提示 + 可选动作区（children） */
export function EnsoEmpty({
  title,
  hint,
  children,
}: {
  title: string
  hint?: string
  children?: ReactNode
}) {
  return (
    <div className="flex items-center gap-4 py-4">
      <EnsoIcon />
      <div>
        <p className="font-serif text-[15px] text-text-primary">{title}</p>
        {hint && <p className="font-mono text-xs text-text-secondary">{hint}</p>}
      </div>
      {children}
    </div>
  )
}

/** 墨韵竹叶：三竿竹 SVG，墨色 opacity 0.045–0.1，ink-sway 极慢摇曳 */
export function BambooArt({ className }: { className?: string }) {
  return (
    <svg
      className={className}
      viewBox="0 0 580 432"
      preserveAspectRatio="xMaxYMax meet"
      aria-hidden="true"
    >
      <defs>
        <path id="bzleaf" d="M0 0 C12 -7 34 -11 58 -3 C36 4 12 5 0 0 Z" />
        <g id="spray">
          <path
            d="M0 0 Q26 -9 56 -7"
            fill="none"
            stroke="currentColor"
            strokeWidth="2"
            strokeLinecap="round"
          />
          <use href="#bzleaf" transform="translate(12 -3) rotate(-44) scale(1.08)" />
          <use href="#bzleaf" transform="translate(28 -6) rotate(-14)" />
          <use href="#bzleaf" transform="translate(44 -7) rotate(14) scale(.95)" />
          <use href="#bzleaf" transform="translate(56 -7) rotate(42) scale(.72)" />
        </g>
      </defs>

      {/* 竹一 */}
      <g className="ink-sway" opacity="0.1">
        <path
          d="M392 432 C388 366 393 314 390 258 C387 202 392 142 389 86 C388 56 390 30 388 8"
          fill="none"
          stroke="currentColor"
          strokeWidth="6.5"
          strokeLinecap="round"
        />
        <ellipse cx="390" cy="350" rx="7.5" ry="2.8" fill="none" stroke="currentColor" strokeWidth="2" />
        <ellipse cx="389.5" cy="258" rx="7.5" ry="2.8" fill="none" stroke="currentColor" strokeWidth="2" />
        <ellipse cx="390" cy="166" rx="7.5" ry="2.8" fill="none" stroke="currentColor" strokeWidth="2" />
        <ellipse cx="390" cy="82" rx="6.5" ry="2.4" fill="none" stroke="currentColor" strokeWidth="2" />
        <use href="#spray" transform="translate(391 318) rotate(-16)" />
        <use href="#spray" transform="translate(389 226) scale(-1 1) rotate(-22)" />
        <use href="#spray" transform="translate(390 134) rotate(-28) scale(1.12)" />
        <use href="#spray" transform="translate(388 50) rotate(-54) scale(1.2)" />
        <use href="#bzleaf" transform="translate(304 404) rotate(26) scale(.78)" />
      </g>

      {/* 竹二 */}
      <g className="ink-sway ink-sway-slow" opacity="0.065">
        <path
          d="M498 432 C495 380 499 338 497 294 C495 248 498 202 496 156 C495 132 497 106 496 84"
          fill="none"
          stroke="currentColor"
          strokeWidth="4.6"
          strokeLinecap="round"
        />
        <ellipse cx="497" cy="356" rx="5.5" ry="2" fill="none" stroke="currentColor" strokeWidth="1.6" />
        <ellipse cx="497" cy="266" rx="5.5" ry="2" fill="none" stroke="currentColor" strokeWidth="1.6" />
        <ellipse cx="497" cy="176" rx="5.5" ry="2" fill="none" stroke="currentColor" strokeWidth="1.6" />
        <use href="#spray" transform="translate(497 326) rotate(-12) scale(.88)" />
        <use href="#spray" transform="translate(496 236) scale(-1 1) rotate(-24) scale(.95)" />
        <use href="#spray" transform="translate(496 138) rotate(-42) scale(1.05)" />
      </g>

      {/* 竹三（更淡更远） */}
      <g className="ink-sway ink-sway-fast" opacity="0.045">
        <path
          d="M300 432 C298 400 300 372 299 340 C298 308 300 280 299 250"
          fill="none"
          stroke="currentColor"
          strokeWidth="3.4"
          strokeLinecap="round"
        />
        <use href="#spray" transform="translate(299 372) rotate(-10) scale(.8)" />
        <use href="#spray" transform="translate(299 300) scale(-1 1) rotate(-18) scale(.82)" />
      </g>
    </svg>
  )
}

/* ───────── 页面级原语：供 dashboard 之外的后台各页收编进同一套语言 ───────── */

/**
 * 页头：mono kicker + 衬线大标题 + 笔刷下划线 + 描述 + 右侧动作区。
 * 取代旧版 `text-3xl font-bold tracking-tight` 通用页头，是全后台页面的统一开场。
 * `backTo` 提供时渲染返回链接（箭头悬停微移）。
 */
export function PageHead({
  kicker,
  title,
  sub,
  backTo,
  backLabel,
  actions,
}: {
  kicker?: string
  title: string
  sub?: string
  backTo?: string
  backLabel?: string
  actions?: ReactNode
}) {
  return (
    <div className="flex flex-wrap items-end justify-between gap-4">
      <div className="min-w-0">
        {backTo && (
          <Link
            to={backTo}
            className="group mb-2.5 inline-flex items-center gap-1.5 font-mono text-xs text-text-secondary transition-colors duration-150 hover:text-leaf-deep"
          >
            <ArrowLeft className="size-3.5 transition-transform duration-200 group-hover:-translate-x-0.5" />
            {backLabel ?? '返回'}
          </Link>
        )}
        {kicker && (
          <p className="font-mono text-[11px] uppercase tracking-[0.2em] text-leaf-deep">
            {kicker}
          </p>
        )}
        <h1 className="font-serif text-4xl font-bold leading-tight tracking-tight text-text-primary">
          {title}
        </h1>
        <BrushUnderline className="mt-2.5" />
        {sub && (
          <p className="mt-3 max-w-xl text-sm leading-relaxed text-text-secondary">
            {sub}
          </p>
        )}
      </div>
      {actions && (
        <div className="flex shrink-0 flex-wrap items-center gap-2">{actions}</div>
      )}
    </div>
  )
}

/**
 * KPI 统计块：宣纸卡片 + 顶部扫入墨线 + 衬线水印字 + CountUp 大数字。
 * 与 dashboard 的 KPI 三块同源，供列表统计条 / 赞助统计等复用。
 */
export function InkStat({
  label,
  value,
  hint,
  watermark,
  loading,
}: {
  label: string
  value: number
  hint?: string
  watermark?: string
  loading?: boolean
}) {
  return (
    <div className={`${inkCard} p-5`}>
      <span className="absolute left-5 right-5 top-[-1px] h-0.5 origin-left scale-x-0 rounded-full bg-leaf-deep transition-transform duration-400 group-hover:scale-x-100" />
      {watermark && (
        <span
          className="pointer-events-none absolute -bottom-4 right-1.5 font-serif text-[80px] font-black leading-none text-text-primary opacity-5"
          aria-hidden
        >
          {watermark}
        </span>
      )}
      <p className="font-mono text-[11px] uppercase tracking-[0.1em] text-text-secondary">
        {label}
      </p>
      {loading ? (
        <Skeleton className="mt-2.5 h-10 w-14" />
      ) : (
        <CountUp
          value={value}
          className="mt-2.5 block font-serif text-[42px] font-semibold leading-none tabular-nums text-text-primary"
        />
      )}
      {hint && <p className="mt-2.5 text-xs text-text-secondary">{hint}</p>}
    </div>
  )
}

/** 筛选药丸：选中 leaf-deep 实底，未选 muted。供状态 / 分组 / 视图切换复用 */
export function InkPill({
  active,
  onClick,
  title,
  children,
}: {
  active: boolean
  onClick: () => void
  title?: string
  children: ReactNode
}) {
  return (
    <button
      type="button"
      title={title}
      onClick={onClick}
      className={cn(
        'cursor-pointer rounded-full px-3.5 py-1.5 text-sm font-medium transition-colors duration-200',
        active
          ? 'bg-leaf-deep text-card shadow-sm'
          : 'bg-muted text-text-secondary hover:bg-muted/70 hover:text-text-primary',
      )}
    >
      {children}
    </button>
  )
}

/** 表格壳与单元格排版常量：衬线/mono 排印 + 宣纸行 hover，配合 shadcn Table 使用 */
export const inkTableWrap =
  'overflow-hidden rounded-lg border border-border bg-card'
export const inkTableHeadRow = 'border-b border-border bg-muted/30'
export const inkTh =
  'px-4 py-2.5 font-mono text-[11px] font-medium uppercase tracking-widest text-text-secondary'
export const inkTd = 'px-4 py-3 text-sm'
export const inkTableRow =
  'border-b border-border/60 transition-colors duration-150 last:border-0 hover:bg-muted/30'

/** 状态开关：leaf-deep 通断，role=switch */
export function InkSwitch({
  checked,
  disabled,
  onToggle,
}: {
  checked: boolean
  disabled?: boolean
  onToggle: () => void
}) {
  return (
    <button
      type="button"
      role="switch"
      aria-checked={checked}
      disabled={disabled}
      onClick={onToggle}
      className={cn(
        'relative inline-flex h-5 w-9 shrink-0 cursor-pointer items-center rounded-full transition-colors duration-200 disabled:cursor-not-allowed disabled:opacity-50',
        checked ? 'bg-leaf-deep' : 'bg-border',
      )}
    >
      <span
        className={cn(
          'inline-block size-4 rounded-full bg-card shadow-sm transition-transform duration-200',
          checked ? 'translate-x-[18px]' : 'translate-x-0.5',
        )}
      />
    </button>
  )
}

/** 徽章色调：一律走 styles.css 既有 token，与 donut 分段同源，禁自创颜色 */
const badgeTones = {
  leaf: 'border-chart-1/30 bg-chart-1/12 text-leaf-deep',
  pending: 'border-chart-4/50 bg-chart-4/20 text-leaf-deep',
  danger: 'border-destructive/25 bg-destructive/10 text-destructive',
  neutral: 'border-border bg-muted/60 text-text-secondary',
} as const

export type InkBadgeTone = keyof typeof badgeTones

/** 状态徽章：圆角方片 + 细边框 + 色调 token */
export function InkBadge({
  tone = 'neutral',
  className,
  children,
}: {
  tone?: InkBadgeTone
  className?: string
  children: ReactNode
}) {
  return (
    <span
      className={cn(
        'inline-flex items-center gap-1 rounded-[4px] border px-1.5 py-0.5 text-xs font-medium',
        badgeTones[tone],
        className,
      )}
    >
      {children}
    </span>
  )
}

/** 友链状态 → 徽章（is_failure=1 已失效；status 1=已通过 2=已拒绝 3=下架待审核 4=已下架 0=待审核） */
export function linkStatus(link: LinkFriend): {
  label: string
  tone: InkBadgeTone
} {
  if (link.is_failure === 1) return { label: '已失效', tone: 'danger' }
  switch (link.status) {
    case 1:
      return { label: '已通过', tone: 'leaf' }
    case 2:
      return { label: '已拒绝', tone: 'danger' }
    case 3:
      return { label: '下架待审核', tone: 'pending' }
    case 4:
      return { label: '已下架', tone: 'neutral' }
    default:
      return { label: '待审核', tone: 'pending' }
  }
}

/**
 * 晨光墨晕：单色淡绿径向，挂于 inkCard / 容器顶部，作「答卷纸」的晨光签名。
 * 与 dialog 浮层、详情 hero 的晨光同源（仅 leaf-light，非多色光斑）。
 * pointer-events-none 且淡度 0.18，覆于其上墨字标题/字段不影响可读，
 * 反如晨光落在题跋首行。容器须为 relative + overflow-hidden（inkCard 已满足）。
 */
export function InkGlow({ className = '' }: { className?: string }) {
  return (
    <span
      aria-hidden
      className={`pointer-events-none absolute inset-x-0 top-0 h-32 ${className}`}
      style={{
        background:
          'radial-gradient(520px 140px at 50% 0%, oklch(0.88 0.1 105 / 0.18), transparent 72%)',
      }}
    />
  )
}

/**
 * 主从导航行：斜墨条（选中态）+ 图标块 + 衬线标题 + 淡墨描述 + chevron。
 * 方案 C「主从面板」的左侧菜单签名，用于设置等页的内部导航。
 * 图标块为功能性导航图标（随选中态着色），非标题装饰。
 */
export function InkNavRow({
  icon,
  title,
  desc,
  active,
  onClick,
}: {
  icon: ReactNode
  title: string
  desc?: string
  active: boolean
  onClick: () => void
}) {
  return (
    <button
      type="button"
      onClick={onClick}
      aria-current={active || undefined}
      className={cn(
        'relative flex w-full cursor-pointer items-center gap-3 rounded-lg px-3.5 py-3 text-left transition-colors duration-150',
        active ? 'bg-leaf-deep/10' : 'hover:bg-muted',
      )}
    >
      <span
        aria-hidden
        className={cn(
          'absolute bottom-2.5 left-0 top-2.5 w-[3px] rounded-sm bg-leaf-deep transition-opacity duration-150',
          active ? 'opacity-100' : 'opacity-0',
        )}
      />
      <span
        className={cn(
          'grid size-9 shrink-0 place-items-center rounded-lg transition-colors duration-150',
          active ? 'bg-leaf-deep/15 text-leaf-deep' : 'bg-muted text-text-secondary',
        )}
      >
        {icon}
      </span>
      <span className="min-w-0 flex-1">
        <span className="block font-serif text-[15px] font-semibold leading-tight text-text-primary">
          {title}
        </span>
        {desc && (
          <span className="mt-0.5 block text-xs text-text-secondary">{desc}</span>
        )}
      </span>
      <ChevronRight
        className={cn(
          'size-4 shrink-0 transition-colors duration-150',
          active ? 'text-leaf-deep' : 'text-text-secondary/40',
        )}
      />
    </button>
  )
}
