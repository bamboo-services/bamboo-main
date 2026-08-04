// --------------------------------------------------------------------------------
// Copyright (c) 2016-NOW(至今) 筱锋
// Author: 筱锋「xiao_lfeng」(https://www.x-lf.com)
// --------------------------------------------------------------------------------
// 许可证声明：版权所有 (c) 2016-2026 筱锋。保留所有权利。
// 有关MIT许可证的更多信息，请查看项目根目录下的LICENSE文件或访问：
// https://opensource.org/licenses/MIT
// --------------------------------------------------------------------------------

import { motion, useReducedMotion } from 'motion/react'
import { useEffect, useState } from 'react'

/**
 * 沉浸式跳转引导层（Interlude）——「随墨入站」。
 *
 * 点击友链卡不直接跳转：先从卡片中心以 clip-path circle() 扩散至全屏，
 * 内容（头像/站名/域名/进度线）错峰渐显，动画完毕才 window.open 真实跳转，
 * 随后引导层收拢。高级友链（premium）以站点截图为全屏背景并附浏览器地址栏，
 * 营造「进入对方小站」的沉浸感（截图未生成时回退墨晕渐变示意）。
 */

/** Interlude 展示所需的数据（由友链卡在点击时传入） */
export interface InterludeData {
  /** 站点名称 */
  name: string
  /** 目标 URL（动画完毕后 window.open） */
  url: string
  /** 头像字符（站名首字） */
  avatarChar: string
  /** 是否高级友链（截图背景 + 浏览器地址栏） */
  premium: boolean
  /** 高级友链站点截图 URL（可为空，空时回退墨晕渐变背景） */
  screenshotUrl?: string | null
  /** 扩散起点（卡片中心视口坐标） */
  origin: { x: number; y: number }
}

/** 从 URL 提取域名用于展示 */
function domainOf(url: string): string {
  return url.replace(/^https?:\/\//, '').replace(/\/$/, '')
}

export function Interlude({
  data,
  onDone,
}: {
  /** 为 null 时引导层关闭 */
  data: InterludeData | null
  /** 引导完成（已跳转并收拢）后的回调，用于清空 data */
  onDone: () => void
}) {
  const reduced = useReducedMotion() ?? false
  const [phase, setPhase] = useState<'expand' | 'collapse'>('expand')

  useEffect(() => {
    if (!data) return
    setPhase('expand')

    // 动画完毕 → 真实跳转 → 收拢 → 通知父级清空
    const navTimer = setTimeout(
      () => {
        window.open(data.url, '_blank', 'noopener')
        setPhase('collapse')
      },
      reduced ? 400 : 1650,
    )
    const doneTimer = setTimeout(onDone, reduced ? 500 : 2450)

    return () => {
      clearTimeout(navTimer)
      clearTimeout(doneTimer)
    }
  }, [data, reduced, onDone])

  if (!data) return null

  const { x, y } = data.origin
  const expanded = `circle(150% at ${x}px ${y}px)`
  const collapsed = `circle(0px at ${x}px ${y}px)`
  const domain = domainOf(data.url)

  return (
    <motion.div
      aria-hidden="true"
      className="fixed inset-0 z-[100] flex items-center justify-center"
      initial={{ clipPath: collapsed }}
      animate={{ clipPath: phase === 'expand' ? expanded : collapsed }}
      transition={{
        duration: reduced ? 0.01 : 0.75,
        ease: [0.22, 0.9, 0.3, 1],
      }}
    >
      {/* 背景层：常规 = 宣纸墨晕；高级 = 站点截图（暗化叠层保证可读，缺图回退墨晕） */}
      {data.premium ? (
        <div
          className="absolute inset-0 bg-cover bg-center"
          style={
            data.screenshotUrl
              ? {
                  backgroundImage: `linear-gradient(180deg, oklch(0.32 0.06 155 / 0.28), oklch(0.32 0.06 155 / 0.5)), url(${data.screenshotUrl})`,
                }
              : {
                  background:
                    'linear-gradient(180deg, oklch(0.32 0.06 155 / 0.28), oklch(0.32 0.06 155 / 0.5)), linear-gradient(135deg, oklch(0.88 0.1 105 / 0.55), oklch(0.99 0.005 110) 42%, oklch(0.8 0.08 130 / 0.5))',
                }
          }
        />
      ) : (
        <div
          className="absolute inset-0"
          style={{
            background:
              'radial-gradient(900px 620px at 50% 38%, oklch(0.88 0.1 105 / 0.4), transparent 70%), radial-gradient(700px 500px at 8% 100%, oklch(0.88 0.1 105 / 0.18), transparent 70%), oklch(0.975 0.016 110)',
          }}
        />
      )}

      {/* 浏览器地址栏（仅高级 · 截图模式） */}
      {data.premium && (
        <motion.div
          className="absolute inset-x-0 top-0 flex items-center gap-[7px] px-[22px] py-4 font-mono"
          initial={{ opacity: 0 }}
          animate={{ opacity: 1 }}
          transition={{ duration: 0.5, delay: reduced ? 0 : 0.55 }}
        >
          <span className="size-[11px] rounded-full bg-card/35" />
          <span className="size-[11px] rounded-full bg-card/35" />
          <span className="size-[11px] rounded-full bg-card/35" />
          <span className="ml-3 text-xs tracking-[0.08em] text-card/75">
            {data.url}
          </span>
        </motion.div>
      )}

      {/* 中央内容：头像 → 站名 → 域名 → 进度线 → 提示，错峰渐显 */}
      <div className="relative z-10 flex flex-col items-center px-6 text-center">
        <motion.div
          className={
            data.premium
              ? 'flex size-24 items-center justify-center rounded-full bg-leaf-light/55 font-serif text-[40px] font-bold text-leaf-deep shadow-[0_0_0_4px_oklch(0.99_0.005_110/0.35)]'
              : 'flex size-24 items-center justify-center rounded-full bg-leaf-light/55 font-serif text-[40px] font-bold text-leaf-deep shadow-[0_0_0_4px_oklch(0.72_0.13_155/0.3)]'
          }
          initial={{ opacity: 0, scale: 0.7 }}
          animate={{ opacity: 1, scale: 1 }}
          transition={{
            duration: reduced ? 0.01 : 0.5,
            delay: reduced ? 0 : 0.25,
            ease: [0.22, 0.9, 0.3, 1],
          }}
        >
          {data.avatarChar}
        </motion.div>

        <motion.h2
          className={
            data.premium
              ? 'mt-[26px] font-serif text-[clamp(2rem,6vw,2.75rem)] font-bold text-card'
              : 'mt-[26px] font-serif text-[clamp(2rem,6vw,2.75rem)] font-bold text-text-primary'
          }
          initial={{ opacity: 0, y: 18 }}
          animate={{ opacity: 1, y: 0 }}
          transition={{
            duration: reduced ? 0.01 : 0.5,
            delay: reduced ? 0 : 0.35,
            ease: [0.22, 0.9, 0.3, 1],
          }}
        >
          {data.name}
        </motion.h2>

        <motion.p
          className={
            data.premium
              ? 'mt-3 font-mono text-[13px] tracking-[0.18em] text-card/72'
              : 'mt-3 font-mono text-[13px] tracking-[0.18em] text-text-secondary'
          }
          initial={{ opacity: 0, y: 14 }}
          animate={{ opacity: 1, y: 0 }}
          transition={{
            duration: reduced ? 0.01 : 0.5,
            delay: reduced ? 0 : 0.45,
            ease: [0.22, 0.9, 0.3, 1],
          }}
        >
          {domain}
        </motion.p>

        {/* 进度线：走满即跳转 */}
        <motion.div
          className={
            data.premium
              ? 'mt-[34px] h-0.5 w-[190px] overflow-hidden rounded-[2px] bg-card/25'
              : 'mt-[34px] h-0.5 w-[190px] overflow-hidden rounded-[2px] bg-border'
          }
          initial={{ opacity: 0 }}
          animate={{ opacity: 1 }}
          transition={{
            duration: reduced ? 0.01 : 0.4,
            delay: reduced ? 0 : 0.55,
          }}
        >
          <motion.span
            className={
              data.premium
                ? 'block h-full rounded-[2px] bg-leaf-light'
                : 'block h-full rounded-[2px] bg-leaf-deep'
            }
            initial={{ width: '0%' }}
            animate={{ width: '100%' }}
            transition={{
              duration: reduced ? 0.3 : 0.9,
              delay: reduced ? 0 : 0.6,
              ease: [0.4, 0, 0.2, 1],
            }}
          />
        </motion.div>

        <motion.p
          className={
            data.premium
              ? 'mt-4 font-mono text-[11px] uppercase tracking-[0.4em] text-card/60'
              : 'mt-4 font-mono text-[11px] uppercase tracking-[0.4em] text-text-secondary'
          }
          initial={{ opacity: 0 }}
          animate={{ opacity: 1 }}
          transition={{
            duration: reduced ? 0.01 : 0.4,
            delay: reduced ? 0 : 0.65,
          }}
        >
          正在前往
        </motion.p>
      </div>
    </motion.div>
  )
}
