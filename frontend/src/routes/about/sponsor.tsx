// --------------------------------------------------------------------------------
// Copyright (c) 2016-NOW(至今) 筱锋
// Author: 筱锋「xiao_lfeng」(https://www.x-lf.com)
// --------------------------------------------------------------------------------
// 许可证声明：版权所有 (c) 2016-2026 筱锋。保留所有权利。
// 有关MIT许可证的更多信息，请查看项目根目录下的LICENSE文件或访问：
// https://opensource.org/licenses/MIT
// --------------------------------------------------------------------------------

import { createFileRoute } from '@tanstack/react-router'
import { useQuery } from '@tanstack/react-query'
import { motion, useReducedMotion } from 'motion/react'
import { Heart, Sprout } from 'lucide-react'
import type { ComponentProps } from 'react'
import type { SponsorChannel, SponsorRecord } from '@/api/types'
import { getPublicChannels, getPublicRecords } from '@/api/sponsor'
import { Avatar, AvatarFallback, AvatarImage } from '@/components/ui/avatar'
import { Badge } from '@/components/ui/badge'
import { Skeleton } from '@/components/ui/skeleton'

export const Route = createFileRoute('/about/sponsor')({
  component: SponsorPage,
})

type MotionDivProps = ComponentProps<typeof motion.div>

function enter(
  reduced: boolean,
  delay: number,
  full: MotionDivProps,
): MotionDivProps {
  if (!reduced) return { ...full, transition: { ...full.transition, delay } }
  return {
    initial: { opacity: 0 },
    animate: { opacity: 1 },
    transition: { duration: 0.3, delay: delay * 0.08 },
  }
}

/** 赞助金额（分）→ 元 */
function formatAmount(cents: number): string {
  return `¥${(cents / 100).toFixed(2)}`
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
    <motion.div
      {...enter(reduced, 0.2, {
        initial: { opacity: 0, y: 16 },
        animate: { opacity: 1, y: 0 },
        transition: { duration: 0.5, ease: 'easeOut' },
      })}
      className="space-y-10"
    >
      {/* 渠道 */}
      <section>
        <h3 className="mb-4 flex items-center gap-2 text-lg font-semibold text-text-primary">
          <Heart className="size-5 text-primary" />
          赞助渠道
        </h3>
        {channelsQuery.isLoading ? (
          <div className="grid gap-4 sm:grid-cols-2 lg:grid-cols-3">
            {Array.from({ length: 3 }).map((_, i) => (
              <Skeleton key={i} className="h-24 w-full rounded-xl" />
            ))}
          </div>
        ) : channels.length > 0 ? (
          <div className="grid gap-4 sm:grid-cols-2 lg:grid-cols-3">
            {channels.map((ch) => (
              <ChannelCard key={ch.id} channel={ch} />
            ))}
          </div>
        ) : (
          <p className="text-sm text-text-secondary">暂未开放赞助渠道。</p>
        )}
      </section>

      {/* 记录 */}
      <section>
        <h3 className="mb-4 flex items-center gap-2 text-lg font-semibold text-text-primary">
          <Sprout className="size-5 text-primary" />
          赞助记录
          {total > 0 && (
            <Badge variant="secondary" className="ml-1">
              共 {total} 位
            </Badge>
          )}
        </h3>
        {recordsQuery.isLoading ? (
          <div className="space-y-3">
            {Array.from({ length: 5 }).map((_, i) => (
              <Skeleton key={i} className="h-16 w-full rounded-xl" />
            ))}
          </div>
        ) : records.length > 0 ? (
          <div className="space-y-3">
            {records.map((rec) => (
              <RecordRow key={rec.id} record={rec} />
            ))}
          </div>
        ) : (
          <p className="text-sm text-text-secondary">
            还没有赞助记录，期待你的支持～
          </p>
        )}
      </section>
    </motion.div>
  )
}

function ChannelCard({ channel }: { channel: SponsorChannel }) {
  return (
    <div className="flex items-center gap-3 rounded-xl border border-leaf-muted/40 bg-card/80 p-4 shadow-sm transition-colors hover:border-primary/40">
      <Avatar className="size-12 rounded-full">
        <AvatarImage src={channel.icon ?? undefined} alt={channel.name} />
        <AvatarFallback className="bg-leaf-light/40 text-text-primary">
          {channel.name.slice(0, 1)}
        </AvatarFallback>
      </Avatar>
      <div className="flex-1">
        <p className="font-medium text-text-primary">{channel.name}</p>
        <p className="text-sm text-text-secondary">
          已收到 {channel.sponsor_count} 次赞助
        </p>
      </div>
      <Badge variant="default">{channel.status ? '开放' : '关闭'}</Badge>
    </div>
  )
}

function RecordRow({ record }: { record: SponsorRecord }) {
  return (
    <div className="flex items-center gap-3 rounded-xl border border-leaf-muted/40 bg-card/80 p-4 shadow-sm">
      <div className="flex size-10 shrink-0 items-center justify-center rounded-full bg-leaf-light/40 text-text-primary">
        <Heart className="size-5" />
      </div>
      <div className="min-w-0 flex-1">
        <div className="flex items-baseline justify-between gap-2">
          <span className="truncate font-medium text-text-primary">
            {record.nickname}
          </span>
          <span className="shrink-0 font-semibold text-primary">
            {formatAmount(record.amount)}
          </span>
        </div>
        <div className="mt-0.5 flex items-center gap-2 text-xs text-text-secondary">
          {record.channel && <span>{record.channel.name}</span>}
          {record.channel && <span>·</span>}
          <span>{formatDate(record.sponsor_at)}</span>
        </div>
        {record.message && (
          <p className="mt-1 truncate text-sm text-text-secondary">
            “{record.message}”
          </p>
        )}
      </div>
    </div>
  )
}
