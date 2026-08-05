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
import { useEffect, useState } from 'react'
import { cn } from '@/lib/utils'
import { enter } from '@/lib/motion'
import { AccountHoverCard } from '@/components/layout/account-hover-card'
import { FallingLeaves } from '@/components/decorative/falling-leaves'
import defaultBackground from '@/assets/images/default-background.webp'

export const Route = createFileRoute('/operate')({
  component: OperateLayout,
})

/** 自助申请切换：友链 / 赞助（同一 operate 路由组内互切，样式对齐 about 导航） */
const NAV_ITEMS = [
  { to: '/operate/apply', label: '友链申请' },
  { to: '/operate/sponsor', label: '赞助申请' },
] as const

/** 竹叶小标（导航品牌用） */
function BambooLeafMark() {
  return (
    <svg
      width="20"
      height="20"
      viewBox="0 0 48 32"
      fill="var(--leaf-deep)"
      aria-hidden
    >
      <path d="M2 30C10 18 26 6 46 2c-3 12-16 24-44 28z" />
    </svg>
  )
}

/**
 * 自助管理布局：信笺式单栏。
 *
 * 结构身份区别于 about（全屏模糊背景 + 落叶 + 大字叙事）：
 * operate 是访客的「递交一封信」，纯宣纸底、更窄聚焦、功能导向，
 * 顶导右侧复用公开页账户入口 AccountHoverCard（未登录墨韵胶囊 / 已登录状态卡）。
 */
function OperateLayout() {
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
    <div className="relative flex min-h-dvh flex-col overflow-hidden bg-background">
      {/* 背景图：钉住视口、不随滚动（公开页签名，与 about/首页同源）→ 缓慢进入模糊态 */}
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

      {/* 顶部导航：访客视角。与 about 的差异在导航内容与信笺式内容区，而非剥掉公开页背景 */}
      <nav
        className={cn(
          'fixed inset-x-0 top-0 z-50 border-b transition-colors duration-400',
          scrolled
            ? 'border-border bg-background/94 backdrop-blur-[2px]'
            : 'border-transparent',
        )}
      >
        <div className="mx-auto flex w-full max-w-5xl items-center justify-between px-6 py-4 md:px-10">
          <Link
            to="/"
            className="flex items-center gap-2.5 transition-opacity hover:opacity-70"
          >
            <BambooLeafMark />
            <span className="font-serif text-base font-semibold tracking-wide text-text-primary">
              自助管理
            </span>
          </Link>
          <div className="flex items-center gap-5 md:gap-6">
            {/* 申请切换：友链 / 赞助（参考 about 导航下划线样式） */}
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
            <AccountHoverCard />
          </div>
        </div>
      </nav>

      {/* 子路由内容：信笺式居中单栏 */}
      <motion.div
        {...enter(reduced, 0, {
          initial: { opacity: 0 },
          animate: { opacity: 1 },
          transition: { duration: reduced ? 0.15 : 0.3 },
        })}
        className="relative z-10 flex flex-1 flex-col pt-20"
      >
        <Outlet />
      </motion.div>

      {/* 页脚 */}
      <footer className="relative z-10 mx-auto w-full max-w-5xl px-6 py-12 md:px-10">
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
