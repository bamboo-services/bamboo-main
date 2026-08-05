// --------------------------------------------------------------------------------
// Copyright (c) 2016-NOW(至今) 筱锋
// Author: 筱锋「xiao_lfeng」(https://www.x-lf.com)
// --------------------------------------------------------------------------------
// 许可证声明：版权所有 (c) 2016-2026 筱锋。保留所有权利。
// 有关MIT许可证的更多信息，请查看项目根目录下的LICENSE文件或访问：
// https://opensource.org/licenses/MIT
// --------------------------------------------------------------------------------

import { Link, createFileRoute } from '@tanstack/react-router'
import { useQuery } from '@tanstack/react-query'
import { motion, useReducedMotion } from 'motion/react'
import type { SponsorChannel, SponsorRecord } from '@/api/types'
import type { MotionDivProps } from '@/lib/motion'
import { getPublicChannels, getPublicRecords } from '@/api/sponsor'
import { Avatar, AvatarFallback, AvatarImage } from '@/components/ui/avatar'
import { Button } from '@/components/ui/button'
import { Skeleton } from '@/components/ui/skeleton'
import { BambooArt, EnsoEmpty } from '@/components/ink-wash'
import { enter } from '@/lib/motion'

export const Route = createFileRoute('/about/sponsor')({
  component: SponsorPage,
})

/** 展卷 reveal：滚动至视口内淡入上移（reduced-motion 时直接呈现） */
function scrollReveal(reduced: boolean): MotionDivProps {
  if (reduced) return {}
  return {
    initial: { opacity: 0, y: 30 },
    whileInView: { opacity: 1, y: 0 },
    viewport: { once: true, amount: 0.1 },
    transition: { duration: 0.9, ease: [0.22, 0.9, 0.3, 1] as const },
  }
}

/** 日期格式化 */
function formatDate(iso: string | null): string {
  if (!iso) return '—'
  return new Date(iso).toLocaleDateString('zh-CN', {
    year: 'numeric',
    month: '2-digit',
    day: '2-digit',
  })
}

/** 赞助金额（分）→「整数.小数」拆分，供衬线大字渲染 */
function splitAmount(cents: number): { int: string; dec: string } {
  const [int, dec] = (cents / 100).toFixed(2).split('.')
  return { int, dec }
}

function SponsorPage() {
  const reduced = useReducedMotion() ?? false
  const channelsQuery = useQuery({
    queryKey: ['public', 'channels'],
    queryFn: getPublicChannels,
  })
  const recordsQuery = useQuery({
    queryKey: ['public', 'records', 1],
    queryFn: () => getPublicRecords({ page: 1, pageSize: 20 }),
  })

  const channels = channelsQuery.data ?? []
  const recordsPage = recordsQuery.data
  const records = recordsPage?.data ?? []
  const total = recordsPage?.pagination.total ?? 0

  return (
    <>
      {/* ═══════════ 开场 · 感恩帖（巨「谢」为锚，装饰破出容器） ═══════════ */}
      <section className="relative flex min-h-[80vh] items-center overflow-hidden">
        {/* 巨型「谢」字水印（视觉锚点） */}
        <span
          className="pointer-events-none absolute -right-8 top-1/2 hidden -translate-y-1/2 select-none font-serif text-[26rem] font-black leading-none text-text-primary opacity-[0.045] md:block"
          aria-hidden
        >
          谢
        </span>
        {/* 墨韵竹叶（左下，水平镜像） */}
        <BambooArt className="pointer-events-none absolute -left-16 bottom-0 h-[70%] w-[480px] -scale-x-100 text-text-primary" />

        <div className="relative z-10 mx-auto w-full max-w-6xl px-6 md:px-10 pb-20 pt-36">
          <motion.p
            {...enter(reduced, 0.1, {
              initial: { opacity: 0, y: 16 },
              animate: { opacity: 1, y: 0 },
              transition: { duration: 0.6, ease: 'easeOut' },
            })}
            className="flex items-center gap-4 font-mono text-[11px] uppercase tracking-[0.35em] text-text-secondary"
          >
            <span className="h-px w-10 bg-leaf-deep" />
            赞助 · 感恩录
          </motion.p>

          <motion.h1
            {...enter(reduced, 0.2, {
              initial: { opacity: 0, y: 24 },
              animate: { opacity: 1, y: 0 },
              transition: { duration: 0.7, ease: [0.22, 0.9, 0.3, 1] as const },
            })}
            className="mt-7 font-serif text-[clamp(3.4rem,9.5vw,6.5rem)] font-bold leading-[1.05] tracking-[0.03em] text-text-primary"
          >
            感恩<span className="text-leaf-deep">有你</span>
          </motion.h1>

          <motion.div
            {...enter(reduced, 0.32, {
              initial: { opacity: 0, scaleX: 0 },
              animate: { opacity: 1, scaleX: 1 },
              transition: { duration: 0.7, ease: [0.22, 0.9, 0.3, 1] as const },
            })}
            className="mt-5 origin-left"
          >
            <svg
              className="block h-3 w-40 md:w-52"
              viewBox="0 0 224 12"
              aria-hidden
            >
              <path
                d="M2 7 C 48 1 118 0 222 3 C 150 11 60 12 2 7 Z"
                fill="var(--leaf-deep)"
              />
            </svg>
          </motion.div>

          <motion.p
            {...enter(reduced, 0.42, {
              initial: { opacity: 0, y: 16 },
              animate: { opacity: 1, y: 0 },
              transition: { duration: 0.7, ease: 'easeOut' },
            })}
            className="mt-7 max-w-xl font-serif text-lg italic leading-relaxed text-text-secondary md:text-xl"
          >
            每一份支持，都是竹林里的一缕清风。它让小站得以安静地生长，也让我知道有人在听。
          </motion.p>

          {/* 申请展示入口：实际赞助后递交申请，经核实后入册感恩账册 */}
          <motion.div
            {...enter(reduced, 0.52, {
              initial: { opacity: 0, y: 16 },
              animate: { opacity: 1, y: 0 },
              transition: { duration: 0.7, ease: 'easeOut' },
            })}
            className="mt-10"
          >
            <Link to="/operate/sponsor">
              <Button
                variant="outline"
                className="cursor-pointer rounded-md px-7 py-5 font-serif text-[15px] tracking-wide"
              >
                已赞助 · 申请展示
              </Button>
            </Link>
          </motion.div>
        </div>
      </section>

      {/* ═══════════ 壹 · 渠道（雅致供奉） ═══════════ */}
      <section className="py-16 md:py-20">
        <div className="mx-auto w-full max-w-6xl px-6 md:px-10">
          <motion.div
            {...scrollReveal(reduced)}
            className="mb-8 flex items-end justify-between gap-4"
          >
            <div>
              <p className="font-mono text-[11px] uppercase tracking-[0.35em] text-leaf-deep">
                壹 · 渠道
              </p>
              <h2 className="mt-2.5 font-serif text-2xl font-bold text-text-primary md:text-3xl">
                随喜之门
              </h2>
            </div>
            <span className="font-mono text-[11px] uppercase tracking-[0.25em] text-text-secondary">
              Channels · {channels.length}
            </span>
          </motion.div>

          {channelsQuery.isLoading ? (
            <div className="grid gap-6 sm:grid-cols-2">
              {Array.from({ length: 2 }).map((_, i) => (
                <Skeleton key={i} className="h-28 rounded-lg" />
              ))}
            </div>
          ) : channels.length > 0 ? (
            <motion.div
              {...scrollReveal(reduced)}
              className="grid gap-6 sm:grid-cols-2"
            >
              {channels.map((ch) => (
                <ChannelCard key={ch.id.toString()} channel={ch} />
              ))}
            </motion.div>
          ) : (
            <p className="text-[15px] text-text-secondary">
              暂未开放赞助渠道。
            </p>
          )}
        </div>
      </section>

      {/* ═══════════ 贰 · 记录（感恩账册 · 水印「录」全出血） ═══════════ */}
      <section className="relative overflow-hidden border-y border-border/70">
        <span
          className="pointer-events-none absolute -bottom-12 right-[6%] select-none font-serif text-[170px] font-black leading-none text-text-primary opacity-[0.04] md:text-[220px]"
          aria-hidden
        >
          录
        </span>

        <div className="relative mx-auto w-full max-w-6xl px-6 md:px-10 py-20 md:py-24">
          <motion.div
            {...scrollReveal(reduced)}
            className="mb-4 flex items-end justify-between gap-4"
          >
            <div>
              <p className="font-mono text-[11px] uppercase tracking-[0.35em] text-leaf-deep">
                贰 · 记录
              </p>
              <h2 className="mt-2.5 font-serif text-2xl font-bold text-text-primary md:text-3xl">
                感恩账册
              </h2>
            </div>
            <span className="font-mono text-[11px] uppercase tracking-[0.25em] text-text-secondary">
              共 {total} 位
            </span>
          </motion.div>

          <motion.p
            {...scrollReveal(reduced)}
            className="mb-10 max-w-lg text-[15px] leading-relaxed text-text-secondary"
          >
            账册所载，不分多寡，皆是心意。排名按时间倒序，新近的心意在前。
          </motion.p>

          {recordsQuery.isLoading ? (
            <div className="space-y-3">
              {Array.from({ length: 4 }).map((_, i) => (
                <Skeleton key={i} className="h-20 rounded-lg" />
              ))}
            </div>
          ) : records.length > 0 ? (
            <motion.div
              {...scrollReveal(reduced)}
              className="divide-y divide-border border-y border-border"
            >
              {records.map((rec, i) => (
                <RecordRow key={rec.id.toString()} record={rec} index={i} />
              ))}
            </motion.div>
          ) : (
            <EnsoEmpty
              title="还没有赞助记录"
              hint="期待你的第一份支持，温暖这片竹林"
            />
          )}
        </div>
      </section>
    </>
  )
}

/** 渠道卡（雅致供奉） */
function ChannelCard({ channel }: { channel: SponsorChannel }) {
  return (
    <div className="group relative overflow-hidden rounded-lg border border-border bg-card p-8 transition-all duration-500 hover:-translate-y-1 hover:border-leaf-muted hover:shadow-[0_24px_50px_-30px_oklch(0.32_0.06_155/0.45)]">
      <span className="absolute left-8 right-8 top-[-1px] h-0.5 origin-left scale-x-0 rounded-full bg-leaf-deep transition-transform duration-500 group-hover:scale-x-100" />
      <div className="flex items-center gap-5">
        <Avatar className="size-16 shrink-0 rounded-full ring-2 ring-ring-glow">
          <AvatarImage src={channel.icon ?? undefined} alt={channel.name} />
          <AvatarFallback className="bg-leaf-light/50 font-serif text-2xl font-bold text-leaf-deep">
            {channel.name.slice(0, 1)}
          </AvatarFallback>
        </Avatar>
        <div className="min-w-0">
          <h3 className="font-serif text-xl font-bold text-text-primary">
            {channel.name}
          </h3>
          <p className="mt-1 font-mono text-[11px] uppercase tracking-[0.18em] text-text-secondary">
            已收到 {channel.sponsor_count} 次随喜
          </p>
        </div>
        <span className="ml-auto shrink-0 rounded border border-leaf-muted/60 bg-leaf-light/20 px-2.5 py-1 font-mono text-[10px] uppercase tracking-[0.25em] text-leaf-deep">
          {channel.status ? '开放' : '关闭'}
        </span>
      </div>
    </div>
  )
}

/** 账册行：mono 序号 + 头像 + 衬线昵称 + 渠道·日期 + 斜体留言 + 衬线大字金额 */
function RecordRow({
  record,
  index,
}: {
  record: SponsorRecord
  index: number
}) {
  const { int, dec } = splitAmount(record.amount)
  return (
    <div className="group relative flex items-center gap-5 px-4 py-6 transition-colors duration-300 hover:bg-muted/45 md:gap-7 md:px-6">
      {/* 左侧墨线：hover 自上而下扫入 */}
      <span className="absolute inset-y-0 left-0 w-0.5 origin-top scale-y-0 bg-leaf-deep transition-transform duration-300 group-hover:scale-y-100" />

      <span className="w-8 shrink-0 font-mono text-[11px] tabular-nums text-text-secondary">
        {(index + 1).toString().padStart(2, '0')}
      </span>

      <div className="flex size-12 shrink-0 items-center justify-center rounded-full bg-leaf-light/40 font-serif text-lg font-semibold text-leaf-deep">
        {record.nickname.slice(0, 1)}
      </div>

      <div className="min-w-0 flex-1">
        <p className="truncate font-serif text-lg font-semibold text-text-primary">
          {record.nickname}
        </p>
        <p className="mt-0.5 truncate font-mono text-[11px] text-text-secondary">
          {record.channel?.name ?? '—'} · {formatDate(record.sponsor_at)}
        </p>
        {record.message && (
          <p className="mt-1.5 truncate text-sm italic text-text-secondary">
            「{record.message}」
          </p>
        )}
      </div>

      <span className="shrink-0 font-serif text-2xl font-semibold tabular-nums text-leaf-deep md:text-3xl">
        ¥{int}
        <span className="text-lg">.{dec}</span>
      </span>
    </div>
  )
}
