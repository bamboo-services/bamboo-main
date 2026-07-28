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
import { useCallback, useState } from 'react'
import type { LinkFriend } from '@/api/types'
import { getPublicLinks } from '@/api/link'
import { BambooArt, EnsoEmpty } from '@/components/ink-wash'
import { Skeleton } from '@/components/ui/skeleton'
import { AdFriendCard } from '@/components/about/ad-friend-card'
import { CloseFriendCard } from '@/components/about/close-friend-card'
import { PremiumFriendCard } from '@/components/about/premium-friend-card'
import { RegularFriendCard } from '@/components/about/regular-friend-card'
import { Interlude } from '@/components/about/interlude'
import type { InterludeData } from '@/components/about/interlude'
import type { FriendCardProps } from '@/components/about/friend-card-shared'
import { enter } from '@/lib/motion'
import type { MotionDivProps } from '@/lib/motion'

export const Route = createFileRoute('/about/friends')({
  component: FriendsPage,
})

/** 友链级别枚举（与后端 pkg/constants LinkLevel 对齐：0 一般 / 1 好友 / 2 高级 / 3 广告） */
const LINK_LEVEL = { regular: 0, close: 1, premium: 2, ad: 3 } as const

/** 章节序号（壹 贰 叁 …） */
const CHAPTERS = ['壹', '贰', '叁', '肆', '伍', '陆', '柒', '捌', '玖', '拾'] as const

/** 展卷 reveal：滚动至视口内淡入上移（reduced-motion 时直接呈现） */
function scrollReveal(reduced: boolean): MotionDivProps {
  if (reduced) return {}
  return {
    initial: { opacity: 0, y: 30 },
    whileInView: { opacity: 1, y: 0 },
    viewport: { once: true, amount: 0.08 },
    transition: { duration: 0.9, ease: [0.22, 0.9, 0.3, 1] as const },
  }
}

/** 友链按分组聚合（分组名取自后端嵌套的 group_f_key） */
function groupLinksByGroup(links: Array<LinkFriend>): Array<{
  groupId: string
  name: string
  links: Array<LinkFriend>
}> {
  const map = new Map<string, { name: string; links: Array<LinkFriend> }>()
  for (const link of links) {
    const key = link.group_id != null ? link.group_id.toString() : 'none'
    const entry = map.get(key) ?? {
      name: link.group_f_key?.name ?? '未分组',
      links: [],
    }
    entry.links.push(link)
    map.set(key, entry)
  }
  return Array.from(map.entries())
    .map(([groupId, entry]) => ({
      groupId,
      name: entry.name,
      links: entry.links.sort((a, b) => a.sort_order - b.sort_order),
    }))
    .sort((a, b) => {
      if (a.groupId === 'none') return 1
      if (b.groupId === 'none') return -1
      return 0
    })
}

/** 按级别选用对应的卡片组件（级别来自后端 level 字段） */
function FriendCardSwitch({ link, onOpen }: FriendCardProps) {
  switch (link.level) {
    case LINK_LEVEL.premium:
      return <PremiumFriendCard link={link} onOpen={onOpen} />
    case LINK_LEVEL.close:
      return <CloseFriendCard link={link} onOpen={onOpen} />
    case LINK_LEVEL.ad:
      return <AdFriendCard link={link} onOpen={onOpen} />
    default:
      return <RegularFriendCard link={link} onOpen={onOpen} />
  }
}

function FriendsPage() {
  const reduced = useReducedMotion() ?? false
  const {
    data: links,
    isLoading,
    error,
  } = useQuery({
    queryKey: ['public', 'links'],
    queryFn: () => getPublicLinks(),
  })
  const [interlude, setInterlude] = useState<InterludeData | null>(null)

  /** 点击友链卡 → 组装 Interlude 数据（高级友链用截图背景） */
  const handleOpen = useCallback(
    (link: LinkFriend, origin: { x: number; y: number }) => {
      setInterlude({
        name: link.name,
        url: link.url,
        avatarChar: link.name.slice(0, 1),
        premium: link.level === LINK_LEVEL.premium,
        origin,
      })
    },
    [],
  )

  const closeInterlude = useCallback(() => setInterlude(null), [])

  return (
    <>
      {/* ═══════════ 开场 · 群贤帖（装饰破出容器） ═══════════ */}
      <section className="relative flex min-h-[72vh] items-center overflow-hidden">
        {/* 墨韵竹叶（左侧全出血，水平镜像） */}
        <BambooArt className="pointer-events-none absolute -left-20 top-0 h-full w-[560px] -scale-x-100 text-text-primary md:w-[680px]" />

        <div className="relative z-10 mx-auto w-full max-w-6xl px-6 md:px-10 pb-16 pt-36">
          <motion.p
            {...enter(reduced, 0.1, {
              initial: { opacity: 0, y: 16 },
              animate: { opacity: 1, y: 0 },
              transition: { duration: 0.6, ease: 'easeOut' },
            })}
            className="flex items-center gap-4 font-mono text-[11px] uppercase tracking-[0.35em] text-text-secondary"
          >
            <span className="h-px w-10 bg-leaf-deep" />
            竹林友链 · 以文会友
          </motion.p>

          <motion.h1
            {...enter(reduced, 0.2, {
              initial: { opacity: 0, y: 24 },
              animate: { opacity: 1, y: 0 },
              transition: { duration: 0.7, ease: [0.22, 0.9, 0.3, 1] as const },
            })}
            className="mt-7 font-serif text-[clamp(3.4rem,9.5vw,6.5rem)] font-bold leading-[1.05] tracking-[0.03em] text-text-primary"
          >
            群贤<span className="text-leaf-deep">毕至</span>
          </motion.h1>

          <motion.div
            {...enter(reduced, 0.32, {
              initial: { opacity: 0, scaleX: 0 },
              animate: { opacity: 1, scaleX: 1 },
              transition: { duration: 0.7, ease: [0.22, 0.9, 0.3, 1] as const },
            })}
            className="mt-5 origin-left"
          >
            <svg className="block h-3 w-40 md:w-52" viewBox="0 0 224 12" aria-hidden>
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
            少长咸集，各美其美。轻点任意一位朋友，随墨而入其小站。
          </motion.p>
        </div>
      </section>

      {/* ═══════════ 分组章节 · Bento 混排栅格 ═══════════ */}
      {isLoading ? (
        <div className="mx-auto w-full max-w-6xl px-6 md:px-10">
          <div className="grid grid-flow-dense grid-cols-2 gap-4 pb-20 [grid-auto-rows:168px] sm:grid-cols-3 sm:[grid-auto-rows:158px] md:grid-cols-4 md:[grid-auto-rows:150px]">
            <Skeleton className="col-span-2 row-span-2 rounded-lg" />
            {Array.from({ length: 5 }).map((_, i) => (
              <Skeleton key={i} className="rounded-lg" />
            ))}
          </div>
        </div>
      ) : error || !links || links.length === 0 ? (
        <div className="mx-auto w-full max-w-6xl px-6 md:px-10">
          <div className="pb-24 pt-8">
            <EnsoEmpty
              title="还没有友链"
              hint="快来申请第一个，让竹林热闹起来"
            />
          </div>
        </div>
      ) : (
        groupLinksByGroup(links).map((group, gi) => (
          <section key={`g-${group.groupId}`} className="py-16 md:py-20">
            <div className="mx-auto w-full max-w-6xl px-6 md:px-10">
              <motion.div
                {...scrollReveal(reduced)}
                className="mb-8 flex items-end justify-between gap-4"
              >
                <div>
                  <p className="font-mono text-[11px] uppercase tracking-[0.35em] text-leaf-deep">
                    {CHAPTERS[gi] ?? '其'}
                  </p>
                  <h2 className="mt-2.5 font-serif text-2xl font-bold text-text-primary md:text-3xl">
                    {group.name}
                  </h2>
                </div>
                <span className="font-mono text-[11px] uppercase tracking-[0.25em] text-text-secondary">
                  {group.links.length} 位
                </span>
              </motion.div>

              {/* Bento 栅格：1×1 基础，高级 2×2，dense 流自动填补空隙 */}
              <motion.div
                {...scrollReveal(reduced)}
                className="grid grid-flow-dense grid-cols-2 gap-4 [grid-auto-rows:168px] sm:grid-cols-3 sm:[grid-auto-rows:158px] md:grid-cols-4 md:[grid-auto-rows:150px]"
              >
                {group.links.map((link) => (
                  <FriendCardSwitch
                    key={link.id.toString()}
                    link={link}
                    onOpen={handleOpen}
                  />
                ))}
              </motion.div>
            </div>
          </section>
        ))
      )}

      {/* 沉浸式跳转引导层 */}
      <Interlude data={interlude} onDone={closeInterlude} />
    </>
  )
}
