// --------------------------------------------------------------------------------
// Copyright (c) 2016-NOW(至今) 筱锋
// Author: 筱锋「xiao_lfeng」(https://www.x-lf.com)
// --------------------------------------------------------------------------------
// 许可证声明：版权所有 (c) 2016-2026 筱锋。保留所有权利。
// 有关MIT许可证的更多信息，请查看项目根目录下的LICENSE文件或访问：
// https://opensource.org/licenses/MIT
// --------------------------------------------------------------------------------

import {
  Link,
  Outlet,
  createFileRoute,
  useLocation,
} from '@tanstack/react-router'
import { motion, useReducedMotion } from 'motion/react'
import { ArrowLeft, Heart, Sprout, User } from 'lucide-react'
import { useEffect, useRef } from 'react'
import type { ComponentProps } from 'react'
import { FallingLeaves } from '@/components/decorative/falling-leaves'
import defaultBackground from '@/assets/images/default-background.webp'

export const Route = createFileRoute('/about')({
  component: AboutLayout,
})

type MotionDivProps = ComponentProps<typeof motion.div>

/** 入场动画：reduced-motion 退化为快速淡入 */
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

const NAV_ITEMS = [
  { to: '/about/me', label: '关于我', icon: User },
  { to: '/about/friends', label: '友链', icon: Sprout },
  { to: '/about/sponsor', label: '赞助', icon: Heart },
] as const

/** 子路由载入过渡：淡入上移（离开无动画，瞬间切换） */
const pageTransition = {
  initial: { opacity: 0, y: 12 },
  animate: { opacity: 1, y: 0 },
}

function AboutLayout() {
  const reduced = useReducedMotion() ?? false
  const thisYear = new Date().getFullYear()
  const pathname = useLocation({ select: (loc) => loc.pathname })

  // 区分「初次载入」（内容有顺序 delay）与「路由切换」（立即载入）
  const isFirstRender = useRef(true)
  useEffect(() => {
    isFirstRender.current = false
  }, [])

  return (
    <div className="relative flex min-h-dvh flex-col overflow-hidden bg-background">
      {/* 背景图：清晰 → 缓慢进入模糊态（1.5s，比首页的晨雾动画更轻） */}
      <motion.div
        initial={reduced ? false : { opacity: 0, filter: 'blur(0px)' }}
        animate={{ opacity: 1, filter: 'blur(8px)' }}
        transition={{ duration: 1.5, ease: 'easeOut' }}
        className="absolute inset-0 bg-cover bg-center bg-no-repeat"
        style={{ backgroundImage: `url(${defaultBackground})` }}
      />
      <div
        className="absolute inset-0"
        style={{
          background:
            'linear-gradient(135deg, var(--overlay-from) 0%, var(--overlay-via) 50%, var(--overlay-to) 100%)',
        }}
      />
      <FallingLeaves />

      {/* 主内容：约束最大宽度，超大屏不无限拉宽；移动端底部留白给 Dock 栏 */}
      <div className="relative z-10 mx-auto w-full max-w-4xl px-4 pb-28 pt-12 md:px-8 md:pb-20 md:pt-20">
        <motion.div
          {...enter(reduced, 0.15, {
            initial: { opacity: 0, y: 16 },
            animate: { opacity: 1, y: 0 },
            transition: { duration: 0.5, ease: 'easeOut' },
          })}
          className="mb-8 text-center"
        >
          <h1 className="text-3xl font-bold tracking-tight text-text-primary md:text-4xl">
            关于小站
          </h1>
          <p className="mt-2 text-sm text-text-secondary md:text-base">
            这里记录着我、朋友们，以及支持过我的人
          </p>
        </motion.div>

        {/* 顶部路由导航：仅桌面端显示（顺序第二步，淡入上移） */}
        <motion.nav
          {...enter(reduced, 0.3, {
            initial: { opacity: 0, y: 12 },
            animate: { opacity: 1, y: 0 },
            transition: { duration: 0.4, ease: 'easeOut' },
          })}
          className="mx-auto mb-8 hidden w-full max-w-md grid-cols-3 gap-1 rounded-lg bg-muted/50 p-1 md:grid"
        >
          {NAV_ITEMS.map((item) => (
            <Link
              key={item.to}
              to={item.to}
              activeProps={{
                className: 'bg-background text-text-primary shadow-sm',
              }}
              inactiveProps={{
                className: 'text-muted-foreground hover:text-text-primary',
              }}
              className="flex items-center justify-center gap-1.5 rounded-md px-3 py-2 text-sm font-medium transition-colors"
            >
              <item.icon className="size-4" />
              {item.label}
            </Link>
          ))}
        </motion.nav>

        {/* 子路由内容：离开无动画（瞬间切换）。
            初次载入时 delay 0.45 让标题→导航先出来，后续路由切换立即载入 */}
        <motion.div
          key={pathname}
          initial={reduced ? { opacity: 0 } : pageTransition.initial}
          animate={pageTransition.animate}
          transition={{
            duration: reduced ? 0.15 : 0.3,
            delay: isFirstRender.current ? 0.45 : 0,
            ease: 'easeOut',
          }}
        >
          <Outlet />
        </motion.div>
      </div>

      {/* 移动端 Dock 栏：iOS 风格毛玻璃底部导航（顺序：内容之后弹入） */}
      <motion.nav
        {...enter(reduced, 0.5, {
          initial: { opacity: 0, y: 24 },
          animate: { opacity: 1, y: 0 },
          transition: { type: 'spring', stiffness: 220, damping: 20 },
        })}
        className="fixed inset-x-4 bottom-4 z-40 md:hidden"
      >
        <div className="mx-auto flex max-w-sm items-center justify-around rounded-3xl bg-white/70 px-2 py-2 shadow-lg shadow-primary/20 ring-1 ring-black/5 backdrop-blur-md">
          {NAV_ITEMS.map((item) => (
            <Link
              key={item.to}
              to={item.to}
              activeProps={{
                className: 'bg-primary/15 text-primary',
              }}
              inactiveProps={{
                className: 'text-text-secondary',
              }}
              className="flex flex-col items-center gap-0.5 rounded-2xl px-5 py-1.5 transition-colors"
            >
              <item.icon className="size-5" />
              <span className="text-xs font-medium">{item.label}</span>
            </Link>
          ))}
        </div>
      </motion.nav>

      {/* 返回首页：右下角悬浮按钮（顺序收尾；移动端抬高避开 Dock 栏） */}
      <motion.div
        {...enter(reduced, 0.6, {
          initial: { opacity: 0, scale: 0.8 },
          animate: { opacity: 1, scale: 1 },
          transition: { type: 'spring', stiffness: 260, damping: 20 },
        })}
        className="fixed bottom-28 right-6 z-40 md:bottom-6 md:right-20"
      >
        <Link
          to="/"
          className="flex size-12 items-center justify-center rounded-full bg-primary text-primary-foreground shadow-lg shadow-primary/30 transition-transform hover:scale-110 focus:outline-none focus:ring-2 focus:ring-ring focus:ring-offset-2"
          aria-label="返回首页"
          title="返回首页"
        >
          <ArrowLeft className="size-5" />
        </Link>
      </motion.div>

      {/* 页脚：仅版权，不加额外背景。mt-auto 让它内容不足时贴底 */
      }
      <footer className="relative z-10 mt-auto px-4 py-6 text-center text-sm text-text-secondary md:px-8">
        <p>Copyright (C) 2016-{thisYear} 筱锋 xiao_lfeng. All Rights Reserved.</p>
      </footer>
    </div>
  )
}
