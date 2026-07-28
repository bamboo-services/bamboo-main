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

import { Link, createFileRoute } from '@tanstack/react-router'
import { motion, useReducedMotion } from 'motion/react'
import { ArrowRight, Inbox } from 'lucide-react'
import type { ComponentProps } from 'react'
import type { DonutSegment } from '@/components/dashboard/donut-chart'
import { Avatar, AvatarFallback, AvatarImage } from '@/components/ui/avatar'
import { Skeleton } from '@/components/ui/skeleton'
import { useSidebar } from '@/components/ui/sidebar'
import { useDashboardStats, useHealth } from '@/hooks/use-dashboard'
import { getStoredUser } from '@/lib/auth'
import { CountUp } from '@/components/dashboard/count-up'
import { DonutChart } from '@/components/dashboard/donut-chart'
import { RadialGauge } from '@/components/dashboard/radial-gauge'

export const Route = createFileRoute('/_admin/admin/dashboard')({
  component: DashboardPage,
})

type MotionDivProps = ComponentProps<typeof motion.div>

/** 入场：淡入 + 轻微 y 位移，reduced 时退化为快速淡入 */
function enter(reduced: boolean, delay: number, full: MotionDivProps): MotionDivProps {
  if (!reduced) {
    return { ...full, transition: { ...full.transition, delay } }
  }
  return {
    initial: { opacity: 0 },
    animate: { opacity: 1 },
    transition: { duration: 0.25, delay: delay * 0.08 },
  }
}

/** 相对时间（最近申请索引用） */
function relativeTime(iso: string): string {
  const diff = Date.now() - new Date(iso).getTime()
  const minutes = Math.floor(diff / 60_000)
  if (minutes < 1) return '刚刚'
  if (minutes < 60) return `${minutes} 分钟前`
  const hours = Math.floor(minutes / 60)
  if (hours < 24) return `${hours} 小时前`
  const days = Math.floor(hours / 24)
  return `${days} 天前`
}

/** 按时段问候语 */
function greeting(): string {
  const hour = new Date().getHours()
  if (hour < 6) return '夜深了'
  if (hour < 12) return '早上好'
  if (hour < 14) return '中午好'
  if (hour < 18) return '下午好'
  return '晚上好'
}

/** 日期行：2026-07-28 周二 */
function dateLine(): string {
  const d = new Date()
  const pad = (n: number) => n.toString().padStart(2, '0')
  const week = ['周日', '周一', '周二', '周三', '周四', '周五', '周六'][d.getDay()]
  return `${d.getFullYear()}-${pad(d.getMonth() + 1)}-${pad(d.getDate())} ${week}`
}

/** 从 cpu_usage 字符串解析 0~1 比例；失败返回 null */
function parseCpuFraction(raw: string | undefined): number | null {
  if (!raw) return null
  const n = parseFloat(raw)
  return Number.isFinite(n) ? Math.min(1, Math.max(0, n / 100)) : null
}

/** 卡片基类：宣纸块，hover 浮起 + 墨边 */
const inkCard =
  'group relative overflow-hidden rounded-lg border border-border bg-card p-5 transition-all duration-300 hover:-translate-y-0.5 hover:border-leaf-muted hover:shadow-[0_14px_30px_-22px_oklch(0.32_0.06_155/0.4)]'

/** 卡片标题：斜墨条 + 衬线标题 + 右侧 mono meta */
function CardHead({ title, meta }: { title: string; meta?: string }) {
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

function DashboardPage() {
  const reduced = useReducedMotion() ?? false
  const { open } = useSidebar()
  const statsQuery = useDashboardStats()
  const healthQuery = useHealth()
  const stats = statsQuery.data
  const runtime = healthQuery.data?.runtime
  const system = healthQuery.data?.system

  const user = getStoredUser()
  const displayName = user?.nickname || user?.username

  const total = stats?.total_links ?? 0
  const approved = stats?.approved_links ?? 0
  const pending = stats?.pending_links ?? 0
  const other = Math.max(0, total - approved - pending)

  const kpiItems = [
    { label: '友链总数', value: total, hint: '已收录 · Total', watermark: '林' },
    { label: '待审核', value: pending, hint: '待处理 · Awaiting', watermark: '候' },
    { label: '已通过', value: approved, hint: '已上线 · Live', watermark: '通' },
  ]

  const donutSegments: Array<DonutSegment> = [
    { value: approved, color: 'var(--chart-1)', label: '已通过' },
    { value: pending, color: 'var(--chart-4)', label: '待审核' },
    { value: other, color: 'var(--chart-5)', label: '其他' },
  ]

  const cpuFrac = parseCpuFraction(runtime?.cpu_usage)
  const isHealthy = healthQuery.isSuccess && Boolean(runtime)

  const heroProps = {
    reduced,
    healthy: isHealthy,
    loading: healthQuery.isLoading,
    displayName,
    dateStr: dateLine(),
    version: system?.version,
    environment: system?.environment,
    goVersion: system?.go_version,
    platform: system?.platform,
  }

  return (
    <div
      className={
        open
          ? 'mx-auto max-w-6xl'
          : 'mx-auto grid max-w-6xl gap-6 lg:grid-cols-[280px_1fr]'
      }
    >
      {/* ───────── Hero：展开=顶部横排 / 收起=左侧竖排卷轴 ───────── */}
      {open ? (
        <>
          <HorizontalHero {...heroProps} />
          {/* 竹节分隔线 */}
          <BambooRule reduced={reduced} delay={0.18} />
        </>
      ) : (
        <VerticalColophon {...heroProps} />
      )}

      {/* ───────── 数据区：展开时在 hero 下方，收起时在卷轴右侧 ───────── */}
      <div>
        {/* KPI 三块 */}
        <motion.section
          {...enter(reduced, 0.22, {
            initial: { opacity: 0, y: 10 },
            animate: { opacity: 1, y: 0 },
            transition: { duration: 0.4, ease: 'easeOut' },
          })}
          className="grid grid-cols-3 gap-4"
        >
          {kpiItems.map((item) => (
            <div key={item.label} className={`${inkCard} p-5`}>
              {/* 顶部墨线，hover 时从左扫入 */}
              <span className="absolute left-5 right-5 top-[-1px] h-0.5 origin-left scale-x-0 rounded-full bg-leaf-deep transition-transform duration-400 group-hover:scale-x-100" />
              {/* 衬线水印字 */}
              <span
                className="pointer-events-none absolute -bottom-4 right-1.5 font-serif text-[96px] font-black leading-none text-text-primary opacity-5"
                aria-hidden
              >
                {item.watermark}
              </span>
              <p className="font-mono text-[11px] uppercase tracking-[0.1em] text-text-secondary">
                {item.label}
              </p>
              {statsQuery.isLoading ? (
                <Skeleton className="mt-2.5 h-12 w-16" />
              ) : (
                <CountUp
                  value={item.value}
                  className="mt-2.5 block font-serif text-[54px] font-semibold leading-none tabular-nums text-text-primary"
                />
              )}
              <p className="mt-2.5 text-xs text-text-secondary">{item.hint}</p>
            </div>
          ))}
        </motion.section>

        {/* 中部两列：donut + 系统状态 */}
        <div className="mt-4 grid gap-4 lg:grid-cols-12">
          {/* 友链构成 donut（5/12） */}
          <motion.section
            {...enter(reduced, 0.3, {
              initial: { opacity: 0, y: 10 },
              animate: { opacity: 1, y: 0 },
              transition: { duration: 0.4, ease: 'easeOut' },
            })}
            className={`${inkCard} flex flex-col lg:col-span-5`}
          >
            <CardHead title="友链构成" meta={`TOTAL ${total}`} />
            {statsQuery.isLoading ? (
              <Skeleton className="mx-auto my-6 h-44 w-44 rounded-full" />
            ) : total > 0 ? (
              <div className="mt-2 flex flex-wrap items-center justify-center gap-6">
                <DonutChart
                  segments={donutSegments}
                  size={170}
                  thickness={16}
                  center={
                    <>
                      <CountUp
                        value={total}
                        className="font-serif text-[38px] font-semibold leading-none tabular-nums text-text-primary"
                      />
                      <span className="font-mono text-[10px] uppercase tracking-[0.14em] text-text-secondary">
                        Total Links
                      </span>
                    </>
                  }
                />
                <ul className="flex min-w-[180px] flex-1 flex-col">
                  {donutSegments.map((seg) => (
                    <li
                      key={seg.label}
                      className="flex items-center gap-2.5 border-b border-border px-1.5 py-2.5 last:border-b-0 transition-colors hover:bg-muted/50"
                    >
                      <span
                        className="size-2.5 shrink-0 rounded-[2px]"
                        style={{ backgroundColor: seg.color }}
                      />
                      <span className="text-sm text-text-primary">{seg.label}</span>
                      <span className="ml-auto font-mono text-sm font-medium tabular-nums text-text-primary">
                        {statsQuery.isLoading ? '—' : seg.value}
                      </span>
                      <span className="w-10 text-right font-mono text-[11px] tabular-nums text-text-secondary">
                        {total > 0 ? `${Math.round((seg.value / total) * 100)}%` : '0%'}
                      </span>
                    </li>
                  ))}
                </ul>
              </div>
            ) : (
              <div className="flex h-44 flex-col items-center justify-center gap-2 text-text-secondary">
                <Inbox className="size-7 opacity-40" />
                <p className="text-sm">暂无友链数据</p>
              </div>
            )}
          </motion.section>

          {/* 系统状态（7/12） */}
          <motion.section
            {...enter(reduced, 0.36, {
              initial: { opacity: 0, y: 10 },
              animate: { opacity: 1, y: 0 },
              transition: { duration: 0.4, ease: 'easeOut' },
            })}
            className={`${inkCard} flex flex-col lg:col-span-7`}
          >
            <CardHead title="系统状态" meta="SYSTEM · METRICS" />
            {healthQuery.isLoading ? (
              <div className="flex flex-1 items-center justify-center py-6">
                <Skeleton className="h-36 w-36 rounded-full" />
              </div>
            ) : runtime ? (
              <>
                <div className="mt-2 grid flex-1 grid-cols-1 gap-5 sm:grid-cols-[200px_1fr]">
                  {/* CPU 仪表 */}
                  <div className="flex flex-col items-center gap-2 pt-1">
                    {cpuFrac != null ? (
                      <RadialGauge
                        value={cpuFrac}
                        size={140}
                        thickness={12}
                        color="var(--leaf-deep)"
                        center={
                          <>
                            <span className="font-mono text-xl font-semibold tabular-nums text-text-primary">
                              {runtime.cpu_usage}
                            </span>
                            <span className="font-mono text-[11px] uppercase tracking-widest text-text-secondary">
                              CPU 占用
                            </span>
                          </>
                        }
                      />
                    ) : (
                      <div className="flex size-[140px] flex-col items-center justify-center gap-1">
                        <Inbox className="size-5 text-text-secondary" />
                        <span className="font-mono text-sm text-text-primary">
                          {runtime.cpu_usage || '—'}
                        </span>
                      </div>
                    )}
                  </div>

                  {/* 指标行 */}
                  <dl className="flex flex-col">
                    <RuntimeRow label="运行时长" value={runtime.uptime} />
                    <RuntimeRow label="内存占用" value={runtime.memory_usage} />
                    <RuntimeRow label="协程数" value={String(runtime.goroutines)} />
                  </dl>
                </div>

                {/* 运行环境 chips */}
                {system && (
                  <div className="mt-4 flex flex-wrap items-center gap-2 border-t border-border pt-3">
                    <span className="font-mono text-[11px] uppercase tracking-widest text-text-secondary">
                      Runtime
                    </span>
                    {[
                      system.version,
                      system.environment,
                      system.go_version,
                      system.platform,
                    ].map((item, i) => (
                      <span
                        key={`${item}-${i.toString()}`}
                        className="rounded border border-border bg-muted/40 px-2 py-0.5 font-mono text-[11px] text-text-primary transition-colors hover:border-leaf-muted"
                      >
                        {item}
                      </span>
                    ))}
                  </div>
                )}
              </>
            ) : (
              <div className="flex flex-1 items-center justify-center gap-2 py-8 text-text-secondary">
                <Inbox className="size-7 opacity-40" />
                <p className="text-sm">暂无运行时数据</p>
              </div>
            )}
          </motion.section>
        </div>

        {/* 最近申请 */}
        <motion.section
          {...enter(reduced, 0.44, {
            initial: { opacity: 0, y: 10 },
            animate: { opacity: 1, y: 0 },
            transition: { duration: 0.4, ease: 'easeOut' },
          })}
          className={`${inkCard} mt-4`}
        >
          <CardHead title="最近申请" meta={pending > 0 ? `${pending} 项待处理` : '0 项待处理'} />
          {statsQuery.isLoading ? (
            <div className="space-y-2">
              {Array.from({ length: 3 }).map((_, i) => (
                <Skeleton key={i} className="h-10 w-full" />
              ))}
            </div>
          ) : stats && stats.recent_applications.length > 0 ? (
            <ul className="divide-y divide-border">
              {stats.recent_applications.map((item, i) => (
                <li key={item.id.toString()}>
                  <Link
                    to="/admin/link/verify"
                    className="flex items-center gap-3 py-2.5 transition-colors hover:bg-muted/40"
                  >
                    <span className="w-6 shrink-0 font-mono text-[11px] text-text-secondary">
                      {(i + 1).toString().padStart(2, '0')}
                    </span>
                    <Avatar className="size-8 rounded">
                      <AvatarImage src={item.avatar || undefined} alt={item.name} />
                      <AvatarFallback className="rounded text-xs">
                        {item.name.slice(0, 1)}
                      </AvatarFallback>
                    </Avatar>
                    <div className="min-w-0 flex-1">
                      <p className="truncate font-serif text-base font-semibold text-text-primary">
                        {item.name}
                      </p>
                      <p className="truncate font-mono text-xs text-text-secondary">
                        {item.url}
                      </p>
                    </div>
                    <div className="flex shrink-0 items-center gap-3">
                      <span className="font-mono text-xs text-text-secondary">
                        {relativeTime(item.created_at)}
                      </span>
                      <span className="text-xs font-medium text-leaf-deep">待审核</span>
                    </div>
                  </Link>
                </li>
              ))}
            </ul>
          ) : (
            <div className="flex items-center gap-4 py-4">
              {/* enso 缺口圆 + 竹叶空状态 */}
              <svg width="48" height="48" viewBox="0 0 48 48" aria-hidden className="shrink-0">
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
              <div>
                <p className="font-serif text-[15px] text-text-primary">暂无待审核申请</p>
                <p className="font-mono text-xs text-text-secondary">
                  新的友链申请提交后将在此处进入审核队列
                </p>
              </div>
              <Link
                to="/admin/link/verify"
                className="ml-auto inline-flex items-center gap-1.5 text-sm font-medium text-leaf-deep transition-opacity hover:opacity-70"
              >
                前往审核
                <ArrowRight className="size-3.5 transition-transform duration-300 group-hover:translate-x-0.5" />
              </Link>
            </div>
          )}
        </motion.section>
      </div>
    </div>
  )
}

/* ───────── Hero props ───────── */

interface HeroProps {
  reduced: boolean
  healthy: boolean
  loading: boolean
  displayName?: string
  dateStr: string
  version?: string
  environment?: string
  goVersion?: string
  platform?: string
}

/** 横排 hero（sidebar 展开，内容区窄）：顶部大问候 + 笔刷 + 导语 + meta */
function HorizontalHero(props: HeroProps) {
  const metaItems = [
    props.environment ? `环境 · ${props.environment}` : null,
    props.version ? `版本 · ${props.version}` : null,
    props.goVersion ? `${props.goVersion}` : null,
    props.platform ? `${props.platform}` : null,
  ].filter(Boolean) as string[]

  return (
    <motion.section
      initial={{ opacity: 0, y: 12 }}
      animate={{ opacity: 1, y: 0 }}
      transition={
        props.reduced
          ? { duration: 0.25 }
          : { duration: 0.6, ease: [0.22, 0.9, 0.3, 1] as const }
      }
      className="relative flex min-h-[300px] items-center overflow-hidden rounded-lg border border-border bg-background"
    >
      {/* 晨光墨晕：单色淡绿径向渐变 */}
      <div
        className="pointer-events-none absolute inset-0"
        style={{
          background:
            'radial-gradient(640px 360px at 83% 4%, oklch(0.88 0.1 105 / 0.22), transparent 70%), radial-gradient(520px 320px at 6% 98%, oklch(0.88 0.1 105 / 0.10), transparent 72%)',
        }}
      />
      {/* 墨韵竹叶背景 */}
      <BambooArt className="pointer-events-none absolute top-0 right-[-30px] h-full w-[580px] text-text-primary" />

      <div className="relative z-10 w-full px-6 py-12 md:px-10">
        <div className="mb-5 flex items-center gap-3.5">
          <StatusDot healthy={props.healthy} loading={props.loading} />
          <span className="font-mono text-xs text-text-secondary">系统运行中</span>
          <span className="text-leaf-muted">·</span>
          <span className="font-mono text-xs text-text-secondary">{props.dateStr}</span>
        </div>

        <h1 className="font-serif text-[64px] font-bold leading-[1.12] tracking-[0.005em] text-text-primary">
          {greeting()}
          {props.displayName ? `，${props.displayName}` : ''}
        </h1>

        {/* 笔刷下划线 */}
        <svg className="my-5 block h-1.5 w-21" viewBox="0 0 84 6" aria-hidden>
          <path
            d="M0 3 C 20 -0.5 44 -1.2 84 0.5 C 50 5 20 6 0 3 Z"
            fill="var(--leaf-deep)"
          />
        </svg>

        <p className="max-w-[520px] font-serif text-xl italic leading-relaxed tracking-[0.015em] text-text-secondary">
          竹林清晨，万物生长。今日友链如竹节，一节一节，缓缓而成。
        </p>

        <div className="mt-7 flex flex-wrap items-center gap-4">
          {metaItems.map((m, i) => (
            <span key={m} className="flex items-center gap-4">
              <span className="font-mono text-xs text-text-secondary">{m}</span>
              {i < metaItems.length - 1 && <span className="text-leaf-muted">·</span>}
            </span>
          ))}
        </div>
      </div>
    </motion.section>
  )
}

/** 竖排卷轴（sidebar 收起，内容区宽）：左侧竖排题跋，数据区在右 */
function VerticalColophon(props: HeroProps) {
  return (
    <motion.aside
      initial={{ opacity: 0, x: -12 }}
      animate={{ opacity: 1, x: 0 }}
      transition={
        props.reduced
          ? { duration: 0.25 }
          : { duration: 0.6, ease: [0.22, 0.9, 0.3, 1] as const }
      }
      className="relative flex min-h-[420px] flex-col overflow-hidden rounded-lg border border-border bg-background"
    >
      {/* 晨光墨晕 */}
      <div
        className="pointer-events-none absolute inset-0"
        style={{
          background:
            'radial-gradient(420px 300px at 80% 8%, oklch(0.88 0.1 105 / 0.20), transparent 70%), radial-gradient(380px 280px at 10% 95%, oklch(0.88 0.1 105 / 0.10), transparent 72%)',
        }}
      />
      {/* 墨韵竹叶背景（更窄，贴右） */}
      <BambooArt className="pointer-events-none absolute top-0 right-[-70px] h-full w-[380px] text-text-primary opacity-70" />

      {/* 顶部状态 */}
      <div className="relative z-10 flex items-center gap-2.5 px-6 pt-6">
        <StatusDot healthy={props.healthy} loading={props.loading} />
        <span className="font-mono text-xs text-text-secondary">系统运行中</span>
        <span className="text-leaf-muted">·</span>
        <span className="font-mono text-xs text-text-secondary">{props.dateStr}</span>
      </div>

      {/* 竖排题跋 */}
      <div className="relative z-10 flex flex-1 items-center justify-center px-6 py-8">
        <div className="[writing-mode:vertical-rl] [text-orientation:upright]">
          <p className="flex flex-col items-start gap-4 font-serif font-bold leading-[1.05] tracking-[0.3em] text-text-primary">
            <span className="text-[56px]">{greeting()}</span>
            {props.displayName && (
              <span className="text-[40px] text-leaf-deep">{props.displayName}</span>
            )}
          </p>
          <p className="mt-8 max-h-[260px] font-serif text-base leading-[2.2] tracking-[0.4em] text-text-secondary">
            竹林清晨，万物生长。
          </p>
        </div>
      </div>

      {/* 底部 meta */}
      <div className="relative z-10 flex flex-col gap-1 px-6 pb-6">
        {props.environment && (
          <span className="font-mono text-[11px] text-text-secondary">
            环境 · {props.environment}
          </span>
        )}
        {props.version && (
          <span className="font-mono text-[11px] text-text-secondary">
            版本 · {props.version}
          </span>
        )}
        {(props.goVersion || props.platform) && (
          <span className="font-mono text-[11px] text-text-secondary">
            {props.goVersion}
            {props.goVersion && props.platform ? ' / ' : ''}
            {props.platform}
          </span>
        )}
      </div>
    </motion.aside>
  )
}

/** 状态点：绿点呼吸（reduced 时静止） */
function StatusDot({ healthy, loading }: { healthy: boolean; loading: boolean }) {
  const color = loading
    ? 'bg-text-secondary/60'
    : healthy
      ? 'bg-chart-1'
      : 'bg-destructive/70'
  return (
    <span
      className={`ink-pulse size-2 shrink-0 rounded-full ${color}`}
      aria-hidden
    />
  )
}

/* ───────── BambooArt：墨韵竹叶 SVG（hero 背景） ───────── */

function BambooArt({ className }: { className?: string }) {
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

/* ───────── BambooRule：竹节分隔线 ───────── */

function BambooRule({ reduced, delay }: { reduced: boolean; delay: number }) {
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

/* ───────── RuntimeRow：运行时指标行 ───────── */

function RuntimeRow({ label, value }: { label: string; value: string }) {
  return (
    <div className="flex items-baseline justify-between gap-3 border-b border-border py-2.5 last:border-b-0">
      <dt className="text-xs text-text-secondary">{label}</dt>
      <dd className="truncate font-mono text-sm font-medium tabular-nums text-text-primary">
        {value}
      </dd>
    </div>
  )
}
