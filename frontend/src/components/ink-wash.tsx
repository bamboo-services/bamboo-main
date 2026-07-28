// --------------------------------------------------------------------------------
// Copyright (c) 2016-NOW(至今) 筱锋
// Author: 筱锋「xiao_lfeng」(https://www.x-lf.com)
// --------------------------------------------------------------------------------
// 许可证声明：版权所有 (c) 2016-2026 筱锋。保留所有权利。
// 有关MIT许可证的更多信息，请查看项目根目录下的LICENSE文件或访问：
// https://opensource.org/licenses/MIT
// --------------------------------------------------------------------------------

import { motion } from 'motion/react'
import type { ReactNode } from 'react'
import { enter } from '@/lib/motion'

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
