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
import { useCallback, useState } from 'react'
import type { LinkFriend } from '@/api/types'
import type { InterludeData } from '@/components/about/interlude'
import type { FriendCardProps } from '@/components/about/friend-card-shared'
import type { MotionDivProps } from '@/lib/motion'
import { BambooArt, EnsoEmpty } from '@/components/ink-wash'
import { Skeleton } from '@/components/ui/skeleton'
import { Button } from '@/components/ui/button'
import { AdFriendCard } from '@/components/about/ad-friend-card'
import { CloseFriendCard } from '@/components/about/close-friend-card'
import { PremiumFriendCard } from '@/components/about/premium-friend-card'
import { RegularFriendCard } from '@/components/about/regular-friend-card'
import { Interlude } from '@/components/about/interlude'
import { enter } from '@/lib/motion'
import { CHAPTERS, LINK_LEVEL, groupLinksByGroup } from '@/lib/friend-groups'
import { usePublicLinks } from '@/hooks/use-links'

export const Route = createFileRoute('/about/friends')({
  component: FriendsPage,
})

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

/**
 * 名帖款行动入口：折角拜帖卡片 + 竖排「拜帖」签条 + 双行文案。
 *
 * 母题：「群贤毕至」典出兰亭雅集——文人入集先递名帖，故申请友链
 * 不是通用 CTA，而是「递上一张名帖」的动作。
 *
 * 交互叙事：名帖默认微微倾斜，如搁在案头静候；hover 时被「拾起」——
 * 摆正、抬起、折角翻开，即双手呈上之态。色彩字体一律走 styles.css token。
 */
function NameCardButton() {
  return (
    <Link
      to="/operate/apply"
      className="group/mingtie relative inline-flex rotate-[2.5deg] items-stretch overflow-hidden rounded-sm border border-leaf-deep/35 bg-card shadow-[0_3px_10px_-3px_oklch(0.32_0.06_155/0.22)] transition-[translate,rotate,border-color,box-shadow] duration-300 hover:-translate-y-1.5 hover:rotate-0 hover:border-leaf-deep hover:shadow-[0_18px_36px_-14px_oklch(0.32_0.06_155/0.38)]"
    >
      {/* 折角翻片：纸的背面，hover 时翻开一角 */}
      <span
        aria-hidden
        className="absolute right-0 top-0 size-[22px] transition-[width,height] duration-300 group-hover/mingtie:size-[26px]"
        style={{
          background:
            'linear-gradient(225deg, oklch(0.9 0.05 120) 0%, oklch(0.82 0.07 130) 48%, oklch(0.94 0.03 115) 52%, oklch(0.97 0.02 110) 100%)',
          clipPath: 'polygon(0 0, 100% 0, 100% 100%)',
        }}
      />
      {/* 竖排「拜帖」签条 */}
      <span
        aria-hidden
        className="flex w-[34px] shrink-0 items-center justify-center border-r border-leaf-deep/20 bg-leaf-deep/8 py-2 [writing-mode:vertical-rl] [text-orientation:upright] font-serif text-[13px] font-bold tracking-[0.3em] text-leaf-deep"
      >
        拜帖
      </span>
      {/* 正文：衬线主标 + mono 副标 */}
      <span className="px-4 py-3 text-left">
        <span className="block font-serif text-[17px] font-bold leading-tight tracking-[0.08em] text-text-primary transition-colors group-hover/mingtie:text-leaf-deep">
          递上名帖
        </span>
        <span className="mt-0.5 block font-mono text-[10px] uppercase tracking-[0.22em] text-text-secondary">
          申请友链 · Apply
        </span>
      </span>
    </Link>
  )
}

function FriendsPage() {
  const reduced = useReducedMotion() ?? false
  const { data: links, isLoading, error } = usePublicLinks()
  const [interlude, setInterlude] = useState<InterludeData | null>(null)

  /** 点击友链卡 → 组装 Interlude 数据（高级友链用截图背景） */
  const handleOpen = useCallback(
    (link: LinkFriend, origin: { x: number; y: number }) => {
      setInterlude({
        name: link.name,
        url: link.url,
        avatarChar: link.name.slice(0, 1),
        premium: link.level === LINK_LEVEL.premium,
        screenshotUrl: link.screenshot_url,
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
          <div className="grid grid-cols-1 items-center gap-12 lg:grid-cols-12">
            {/* ── 左栏 · 开场邀请 ── */}
            <div className="lg:col-span-7">
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
                  transition: {
                    duration: 0.7,
                    ease: [0.22, 0.9, 0.3, 1] as const,
                  },
                })}
                className="mt-7 font-serif text-[clamp(3.4rem,9.5vw,6.5rem)] font-bold leading-[1.05] tracking-[0.03em] text-text-primary"
              >
                群贤<span className="text-leaf-deep">毕至</span>
              </motion.h1>

              <motion.div
                {...enter(reduced, 0.32, {
                  initial: { opacity: 0, scaleX: 0 },
                  animate: { opacity: 1, scaleX: 1 },
                  transition: {
                    duration: 0.7,
                    ease: [0.22, 0.9, 0.3, 1] as const,
                  },
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
                少长咸集，各美其美。轻点任意一位朋友，随墨而入其小站。
              </motion.p>
            </div>

            {/* ── 右栏 · 虚位以待 + 名帖（与左侧墨竹对望，案头呈帖） ── */}
            <div className="flex items-center justify-start gap-7 lg:col-span-5 lg:justify-end">
              {/* 竖排题跋：扫榻以待之意，desktop 专属 */}
              <motion.span
                {...enter(reduced, 0.5, {
                  initial: { opacity: 0, x: 12 },
                  animate: { opacity: 1, x: 0 },
                  transition: { duration: 0.7, ease: 'easeOut' },
                })}
                aria-hidden
                className="hidden select-none font-serif text-sm tracking-[0.55em] text-text-secondary/80 [writing-mode:vertical-rl] [text-orientation:upright] lg:block"
              >
                虚位以待
              </motion.span>

              {/* 名帖：搁在案头微倾，hover 被拾起呈上 */}
              <motion.div
                {...enter(reduced, 0.58, {
                  initial: { opacity: 0, y: 20, rotate: 6, scale: 0.94 },
                  animate: { opacity: 1, y: 0, rotate: 0, scale: 1 },
                  transition: { type: 'spring', stiffness: 220, damping: 15 },
                })}
              >
                <NameCardButton />
              </motion.div>
            </div>
          </div>
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
            <EnsoEmpty title="还没有友链" hint="快来申请第一个，让竹林热闹起来">
              <Link to="/operate/apply" className="ml-auto">
                <Button className="cursor-pointer">申请友链</Button>
              </Link>
            </EnsoEmpty>
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
