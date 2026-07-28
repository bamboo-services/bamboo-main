// --------------------------------------------------------------------------------
// Copyright (c) 2016-NOW(至今) 筱锋
// Author: 筱锋「xiao_lfeng」(https://www.x-lf.com)
// --------------------------------------------------------------------------------
// 许可证声明：版权所有 (c) 2016-2026 筱锋。保留所有权利。
// 有关MIT许可证的更多信息，请查看项目根目录下的LICENSE文件或访问：
// https://opensource.org/licenses/MIT
// --------------------------------------------------------------------------------

import { Link, Outlet, createFileRoute, useLocation } from '@tanstack/react-router'
import { motion, useReducedMotion } from 'motion/react'
import { useEffect, useState } from 'react'
import { cn } from '@/lib/utils'
import { enter } from '@/lib/motion'
import { FallingLeaves } from '@/components/decorative/falling-leaves'
import defaultBackground from '@/assets/images/default-background.webp'

export const Route = createFileRoute('/about')({
  component: AboutLayout,
})

const NAV_ITEMS = [
  { to: '/about/me', label: '关于我' },
  { to: '/about/friends', label: '友链' },
  { to: '/about/sponsor', label: '赞助' },
] as const

/** 竹叶小标（导航品牌用） */
function BambooLeafMark() {
  return (
    <svg width="20" height="20" viewBox="0 0 48 32" fill="var(--leaf-deep)" aria-hidden>
      <path d="M2 30C10 18 26 6 46 2c-3 12-16 24-44 28z" />
    </svg>
  )
}

function AboutLayout() {
  const reduced = useReducedMotion() ?? false
  const thisYear = new Date().getFullYear()
  const pathname = useLocation({ select: (loc) => loc.pathname })
  const [scrolled, setScrolled] = useState(false)

  useEffect(() => {
    const onScroll = () => setScrolled(window.scrollY > 40)
    onScroll()
    addEventListener('scroll', onScroll, { passive: true })
    return () => removeEventListener('scroll', onScroll)
  }, [])

  return (
    <div className="relative min-h-dvh overflow-hidden bg-background">
      {/* 背景图：钉住视口，不随滚动（适配展示屏幕即可）→ 缓慢进入模糊态 */}
      <motion.div
        initial={reduced ? false : { opacity: 0, filter: 'blur(0px)' }}
        animate={{ opacity: 1, filter: 'blur(8px)' }}
        transition={{ duration: 1.5, ease: 'easeOut' }}
        className="fixed inset-0 z-0 bg-cover bg-center bg-no-repeat"
        style={{ backgroundImage: `url(${defaultBackground})` }}
      />
      <div
        className="fixed inset-0 z-0"
        style={{
          background:
            'linear-gradient(135deg, var(--overlay-from) 0%, var(--overlay-via) 50%, var(--overlay-to) 100%)',
        }}
      />
      <FallingLeaves />

      {/* 顶部导航：sticky，滚动后着墨（实底 + 底部墨线） */}
      <nav
        className={cn(
          'fixed inset-x-0 top-0 z-50 border-b transition-colors duration-400',
          scrolled
            ? 'border-border bg-background/94 backdrop-blur-[2px]'
            : 'border-transparent',
        )}
      >
        <div className="mx-auto flex w-full max-w-6xl items-center justify-between px-6 py-5 md:px-10">
          <Link to="/" className="flex items-center gap-2.5 transition-opacity hover:opacity-70">
            <BambooLeafMark />
            <span className="font-serif text-base font-semibold tracking-wide text-text-primary">
              关于小站
            </span>
          </Link>
          <div className="flex items-center gap-7 md:gap-9">
            {NAV_ITEMS.map((item) => {
              const active = pathname === item.to
              return (
                <Link
                  key={item.to}
                  to={item.to}
                  className={cn(
                    'relative font-mono text-[11px] uppercase tracking-[0.28em] transition-colors',
                    'after:absolute after:-bottom-[7px] after:left-0 after:h-0.5 after:w-full after:origin-left after:bg-leaf-deep after:transition-transform after:duration-300 after:ease-[cubic-bezier(0.22,0.9,0.3,1)]',
                    active
                      ? 'text-text-primary after:scale-x-100'
                      : 'text-text-secondary after:scale-x-0 hover:text-text-primary hover:after:scale-x-100',
                  )}
                >
                  {item.label}
                </Link>
              )
            })}
          </div>
        </div>
      </nav>

      {/* 子路由内容：每页自带开场 hero 与宽度约束 */}
      <motion.div
        key={pathname}
        {...enter(reduced, 0, {
          initial: { opacity: 0 },
          animate: { opacity: 1 },
          transition: { duration: reduced ? 0.15 : 0.3 },
        })}
        className="relative z-10"
      >
        <Outlet />
      </motion.div>

      {/* 页脚 */}
      <footer className="relative z-10 mx-auto w-full max-w-6xl px-6 py-14 md:px-10">
        <div className="flex flex-col items-center gap-4 text-center">
          <BambooLeafMark />
          <p className="font-mono text-[11px] tracking-[0.2em] text-text-secondary">
            Copyright (C) 2016-{thisYear} 筱锋 xiao_lfeng · All Rights Reserved
          </p>
        </div>
      </footer>
    </div>
  )
}
