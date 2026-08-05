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
import { ExternalLink, Mail } from 'lucide-react'
import type { MotionDivProps } from '@/lib/motion'
import { getArchive, getSiteInfo } from '@/api/info'
import { useApplySiteInfo, useBloggerInfo } from '@/hooks/use-site-info'
import { MarkdownView } from '@/components/markdown'
import { Avatar, AvatarFallback, AvatarImage } from '@/components/ui/avatar'
import { Skeleton } from '@/components/ui/skeleton'
import { BambooArt, BrushUnderline } from '@/components/ink-wash'
import { enter } from '@/lib/motion'
import { siteConfig } from '@/lib/site'
import myAvatar from '@/assets/images/my_avatar.png'

export const Route = createFileRoute('/about/me')({
  component: AboutMePage,
})

/** 展卷 reveal：滚动至视口内淡入上移（reduced-motion 时直接呈现） */
function scrollReveal(reduced: boolean): MotionDivProps {
  if (reduced) return {}
  return {
    initial: { opacity: 0, y: 30 },
    whileInView: { opacity: 1, y: 0 },
    viewport: { once: true, amount: 0.12 },
    transition: { duration: 0.9, ease: [0.22, 0.9, 0.3, 1] as const },
  }
}

/** 竹叶点角装饰（头像旁） */
function LeafAccent() {
  return (
    <svg
      className="absolute -bottom-2 -left-3 h-9 w-14 rotate-[-18deg]"
      viewBox="0 0 48 32"
      fill="var(--leaf-deep)"
      aria-hidden
    >
      <path d="M2 30C10 18 26 6 46 2c-3 12-16 24-44 28z" />
    </svg>
  )
}

function AboutMePage() {
  const reduced = useReducedMotion() ?? false
  const { data: site } = useQuery({
    queryKey: ['public', 'site'],
    queryFn: getSiteInfo,
  })
  const { data: archive, isLoading } = useQuery({
    queryKey: ['public', 'archive'],
    queryFn: getArchive,
  })
  const { data: blogger } = useBloggerInfo()
  const { data: applySite } = useApplySiteInfo()

  const nick = blogger?.nick ?? siteConfig.blogger.nick
  const desc = blogger?.description ?? siteConfig.blogger.description
  const email = applySite?.email ?? siteConfig.blogger.email
  const blogUrl = blogger?.blog_url ?? 'https://blog.x-lf.com'
  const blogHost = blogUrl.replace(/^https?:\/\//, '')
  const avatar = blogger?.avatar || myAvatar
  const displayName = site?.site_name ?? siteConfig.defaultName
  const siteDesc = archive?.site_description ?? ''
  const introduction = site?.introduction ?? ''

  return (
    <>
      {/* ═══════════ 开场 · 名士帖（全幅展卷，装饰破出容器） ═══════════ */}
      <section className="relative flex min-h-[94vh] items-center overflow-hidden">
        {/* 墨韵竹叶（右侧全出血） */}
        <BambooArt className="pointer-events-none absolute -right-16 top-0 h-full w-[560px] text-text-primary md:w-[760px]" />

        <div className="relative z-10 mx-auto w-full max-w-6xl px-6 md:px-10">
          <div className="grid w-full grid-cols-12 items-center gap-10 pb-24 pt-36">
            {/* 左 7：名号 */}
            <div className="col-span-12 lg:col-span-7">
              <motion.p
                {...enter(reduced, 0.1, {
                  initial: { opacity: 0, y: 16 },
                  animate: { opacity: 1, y: 0 },
                  transition: { duration: 0.6, ease: 'easeOut' },
                })}
                className="flex items-center gap-4 font-mono text-[11px] uppercase tracking-[0.35em] text-text-secondary"
              >
                <span className="h-px w-10 bg-leaf-deep" />
                {displayName} · 站主
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
                className="mt-8 font-serif text-[clamp(4.5rem,13vw,8.5rem)] font-bold leading-[1.02] tracking-[0.02em] text-text-primary"
              >
                {nick.slice(0, 1)}
                <span className="text-leaf-deep">{nick.slice(1)}</span>
              </motion.h1>

              {/* 大写意笔刷 */}
              <motion.div
                {...enter(reduced, 0.32, {
                  initial: { opacity: 0, scaleX: 0 },
                  animate: { opacity: 1, scaleX: 1 },
                  transition: {
                    duration: 0.7,
                    ease: [0.22, 0.9, 0.3, 1] as const,
                  },
                })}
                className="mt-6 origin-left"
              >
                <svg
                  className="block h-3 w-44 md:w-56"
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
                className="mt-8 max-w-lg font-serif text-xl italic leading-relaxed tracking-[0.01em] text-text-secondary md:text-2xl"
              >
                {desc}
              </motion.p>

              {/* 联系方式：克制的行内链接 */}
              <motion.div
                {...enter(reduced, 0.52, {
                  initial: { opacity: 0, y: 16 },
                  animate: { opacity: 1, y: 0 },
                  transition: { duration: 0.7, ease: 'easeOut' },
                })}
                className="mt-10 flex flex-wrap items-center gap-x-8 gap-y-3"
              >
                <a
                  href={blogUrl}
                  target="_blank"
                  rel="noopener noreferrer"
                  className="group inline-flex items-baseline gap-2"
                >
                  <span className="font-mono text-[11px] uppercase tracking-[0.25em] text-text-secondary transition-colors group-hover:text-leaf-deep">
                    博客
                  </span>
                  <span className="flex items-center gap-1 font-serif text-base text-text-primary underline decoration-leaf-muted/60 decoration-1 underline-offset-4 transition-colors group-hover:text-leaf-deep group-hover:decoration-leaf-deep">
                    {blogHost}
                    <ExternalLink className="size-3.5" />
                  </span>
                </a>
                <a
                  href={`mailto:${email}`}
                  className="group inline-flex items-baseline gap-2"
                >
                  <span className="font-mono text-[11px] uppercase tracking-[0.25em] text-text-secondary transition-colors group-hover:text-leaf-deep">
                    邮箱
                  </span>
                  <span className="flex items-center gap-1 font-serif text-base text-text-primary underline decoration-leaf-muted/60 decoration-1 underline-offset-4 transition-colors group-hover:text-leaf-deep group-hover:decoration-leaf-deep">
                    {email}
                    <Mail className="size-3.5" />
                  </span>
                </a>
              </motion.div>
            </div>

            {/* 右 5：头像 + 竖排题跋 */}
            <div className="col-span-12 flex justify-center lg:col-span-5 lg:justify-end">
              <motion.div
                {...enter(reduced, 0.4, {
                  initial: { opacity: 0, scale: 0.92 },
                  animate: { opacity: 1, scale: 1 },
                  transition: {
                    duration: 0.8,
                    ease: [0.22, 0.9, 0.3, 1] as const,
                  },
                })}
                className="flex items-center gap-7"
              >
                <p className="hidden font-serif text-sm tracking-[0.55em] text-text-secondary [text-orientation:upright] [writing-mode:vertical-rl] sm:block">
                  名士帖
                </p>
                <div className="relative">
                  <Avatar className="size-48 rounded-full ring-4 ring-ring-glow md:size-60">
                    <AvatarImage src={avatar} alt={nick} />
                    <AvatarFallback className="bg-leaf-light/50 font-serif text-6xl font-bold text-leaf-deep">
                      {nick.slice(0, 1)}
                    </AvatarFallback>
                  </Avatar>
                  <LeafAccent />
                </div>
              </motion.div>
            </div>
          </div>
        </div>

        {/* 滚动提示 */}
        <div className="absolute bottom-8 left-1/2 z-10 flex -translate-x-1/2 flex-col items-center gap-3">
          <span className="font-mono text-[10px] uppercase tracking-[0.4em] text-text-secondary">
            向下展卷
          </span>
          <motion.span
            className="block h-10 w-px bg-leaf-deep/70"
            animate={reduced ? undefined : { scaleY: [0, 1, 1, 0] }}
            transition={{
              duration: 2.2,
              repeat: Number.POSITIVE_INFINITY,
              ease: [0.65, 0, 0.35, 1] as const,
            }}
            style={{ transformOrigin: 'top' }}
          />
        </div>
      </section>

      {/* ═══════════ 壹 · 自序（动态 Markdown） ═══════════ */}
      <section className="py-24 md:py-32">
        <div className="mx-auto w-full max-w-6xl px-6 md:px-10">
          <div className="grid grid-cols-12 gap-10">
            <aside className="col-span-12 md:col-span-3">
              <motion.div
                {...scrollReveal(reduced)}
                className="md:sticky md:top-28"
              >
                <p className="font-mono text-[11px] uppercase tracking-[0.35em] text-leaf-deep">
                  壹 · 自序
                </p>
                <h2 className="mt-4 font-serif text-3xl font-bold leading-snug text-text-primary md:text-4xl">
                  自我介绍
                </h2>
                <BrushUnderline className="mt-5 w-20" />
              </motion.div>
            </aside>

            <div className="col-span-12 md:col-span-9">
              <motion.div {...scrollReveal(reduced)} className="max-w-prose">
                {isLoading ? (
                  <div className="space-y-4">
                    {Array.from({ length: 4 }).map((_, i) => (
                      <Skeleton key={i} className="h-5 w-full" />
                    ))}
                    <Skeleton className="h-5 w-2/3" />
                  </div>
                ) : archive?.about ? (
                  <MarkdownView content={archive.about} />
                ) : (
                  <p className="text-[17px] leading-[1.95] text-text-secondary">
                    {introduction || '博主尚未填写自我介绍。'}
                  </p>
                )}
              </motion.div>
            </div>
          </div>
        </div>
      </section>

      {/* ═══════════ 贰 · 本站（水印「站」全出血） ═══════════ */}
      <section className="relative overflow-hidden border-y border-border/70">
        <span
          className="pointer-events-none absolute -bottom-10 right-[6%] select-none font-serif text-[180px] font-black leading-none text-text-primary opacity-[0.04] md:text-[240px]"
          aria-hidden
        >
          站
        </span>

        <div className="relative mx-auto w-full max-w-6xl px-6 md:px-10">
          <div className="grid grid-cols-12 gap-10 py-24 md:py-28">
            <aside className="col-span-12 md:col-span-3">
              <motion.div {...scrollReveal(reduced)}>
                <p className="font-mono text-[11px] uppercase tracking-[0.35em] text-leaf-deep">
                  贰 · 本站
                </p>
                <h2 className="mt-4 font-serif text-3xl font-bold leading-snug text-text-primary md:text-4xl">
                  关于本站
                </h2>
                <BrushUnderline className="mt-5 w-20" />
              </motion.div>
            </aside>

            <div className="col-span-12 md:col-span-9">
              <motion.div {...scrollReveal(reduced)} className="max-w-prose">
                {siteDesc ? (
                  <MarkdownView content={siteDesc} />
                ) : (
                  <p className="text-[17px] leading-[1.95] text-text-secondary">
                    博主还没有写下关于本站的介绍。
                  </p>
                )}
              </motion.div>
            </div>
          </div>
        </div>
      </section>
    </>
  )
}
