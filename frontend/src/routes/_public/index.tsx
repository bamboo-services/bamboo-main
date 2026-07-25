/*
 * --------------------------------------------------------------------------------
 * Copyright (c) 2016-NOW(至今) 筱锋
 * Author: 筱锋「xiao_lfeng」(https://www.x-lf.com)
 * --------------------------------------------------------------------------------
 * 许可证声明：版权所有 (c) 2016-2026 筱锋。保留所有权利。
 * 有关MIT许可证的更多信息，请查看项目根目录下的LICENSE文件或访问：
 * https://opensource.org/licenses/MIT
 * --------------------------------------------------------------------------------
 */

import { Link, createFileRoute } from '@tanstack/react-router'
import { motion, useReducedMotion } from 'motion/react'
import { useMemo, useState } from 'react'
import type { ComponentProps } from 'react'
import { BambooLogo } from '@/assets/svg/bamboo-logo'
import { FallingLeaves } from '@/components/decorative/falling-leaves'
import { siteInfo } from '@/data/mock/site-info'
import myAvatar from '@/assets/images/my_avatar.png'
import defaultBackground from '@/assets/images/default-background.webp'

export const Route = createFileRoute('/_public/')({
  component: HomePage,
})

type MotionDivProps = ComponentProps<typeof motion.div>

/**
 * 生成入场动画 props。
 * 无障碍：用户偏好减少动态时，退化为快速淡入（时序整体压缩）
 */
function enter(
  reduced: boolean,
  delay: number,
  full: MotionDivProps,
): MotionDivProps {
  if (!reduced) {
    return { ...full, transition: { ...full.transition, delay } }
  }
  return {
    initial: { opacity: 0 },
    animate: { opacity: 1 },
    transition: { duration: 0.3, delay: delay * 0.08 },
  }
}

function rand(min: number, max: number) {
  return min + Math.random() * (max - min)
}

/** 把角度归一化到 (-π, π]（与落叶物理同源） */
function wrapAngle(angle: number): number {
  const tau = Math.PI * 2
  let next = angle
  while (next > Math.PI) next -= tau
  while (next < -Math.PI) next += tau
  return next
}

interface GuidePath {
  /** 各关键帧水平偏移（相对头像中心，vw） */
  x: Array<string>
  /** 各关键帧垂直偏移（相对头像中心，vh） */
  y: Array<string>
  /** 各关键帧角度（deg，已 unwrap 保证最短弧插值） */
  rotate: Array<number>
  /** 归一化时间（0 → 1） */
  times: Array<number>
  /** 全程时长（s） */
  duration: number
}

/**
 * 引导叶路径 —— 与背景落叶同源的空气动力学积分，
 * 但给一个大初速度：从侧上方高速滑入，阻力让速度自然衰减，
 * 末段落入「静止空气层」，被一阵微风恰好送到头像中心。
 *
 * 与自由落叶的差异：
 * 1. 滑入段速度高，二次方阻力会瞬间吃掉动能（滑距只有对数级），
 *    故改用 Stokes 线性阻力（层流）—— 速度指数衰减，
 *    视觉上就是「快速滑入，慢慢变慢」
 * 2. 末段（后 40%）重力渐隐、阻力渐强（刹车），
 *    同时叠加与距离成正比的引导风：速度随距离自然减小，
 *    渐近收敛到落点，无加速冲击、轻轻落定
 * 3. 旋转沿用风向标对齐 + 颤振 + 阻尼，末段收敛到自然倾斜角
 */
function simulateGuidePath(): GuidePath {
  const dir = Math.random() < 0.5 ? 1 : -1
  const duration = 1.55
  const dt = 0.02
  const sampleEvery = 0.078

  // 起点：头像侧上方、屏幕之外
  let px = dir * rand(28, 40)
  let py = -rand(58, 66)
  // 大初速度：朝头像方向高速滑入
  let vx = -dir * rand(75, 90)
  let vy = rand(80, 92)
  const dragK = rand(1.7, 1.9)
  const gravity = rand(34, 40)

  // 旋转状态（沿用背景落叶模型）
  let rot = rand(-Math.PI, Math.PI)
  let angVel = rand(-2, 2)
  const alignStrength = rand(4, 6)
  const angularDrag = rand(2.2, 3)
  const flutterStrength = rand(1.2, 2)
  const flutterRate = rand(5, 7.5)
  const flutterPhase = rand(0, Math.PI * 2)

  // 落定参数：刹车阻力 + 引导风 + 自然倾斜角
  const settleAngle = (rand(-26, -12) * Math.PI) / 180
  const settleStart = duration * 0.6
  const settleFade = duration * 0.15
  const brakeK = 6.5
  const guideWind = 18

  const rawPoints: Array<{ x: number; y: number; rotate: number }> = [
    { x: px, y: py, rotate: rot },
  ]
  const rawTimes: Array<number> = [0]
  let t = 0
  let lastSample = 0

  const steps = Math.round(duration / dt)
  for (let i = 0; i < steps; i++) {
    t += dt
    // 末段权重：0 → 1（重力/阵风渐隐，刹车与引导风渐强）
    const k = Math.min(1, Math.max(0, (t - settleStart) / settleFade))

    const gustX = 1.2 * Math.sin(t * 0.9 + flutterPhase)
    const gustY = 0.8 * Math.cos(t * 0.7 + flutterPhase * 1.7)

    const dragNow = dragK + (brakeK - dragK) * k
    const fx = -dragNow * vx + gustX * (1 - k) - px * guideWind * k
    const fy =
      -dragNow * vy + gravity * (1 - k) + gustY * (1 - k) - py * guideWind * k

    vx += fx * dt
    vy += fy * dt
    px += vx * dt
    py += vy * dt

    // 旋转：风向标对齐（叶面趋向垂直于气流）+ 颤振 + 阻尼
    const speed = Math.hypot(vx, vy)
    if (speed > 0.5) {
      const flowAngle = Math.atan2(vy, vx)
      const broadside = flowAngle + Math.PI / 2
      const alignTorque = wrapAngle(broadside - rot) * alignStrength
      const flutterTorque =
        Math.sin(t * flutterRate + flutterPhase) * flutterStrength * (1 - k)
      angVel += (alignTorque + flutterTorque - angVel * angularDrag) * dt
    } else {
      angVel += -angVel * angularDrag * dt
    }
    // 末段把角度收敛到落定倾斜
    if (k > 0) {
      angVel += wrapAngle(settleAngle - rot) * 20 * k * dt
    }
    rot = wrapAngle(rot + angVel * dt)

    if (t - lastSample >= sampleEvery) {
      rawPoints.push({ x: px, y: py, rotate: rot })
      rawTimes.push(t)
      lastSample = t
    }
  }
  rawPoints.push({ x: px, y: py, rotate: rot })
  rawTimes.push(duration)

  // 角度 unwrap：保证相邻关键帧走最短弧，插值不反转
  const unwrapped: Array<number> = []
  for (let i = 0; i < rawPoints.length; i++) {
    let r = (rawPoints[i].rotate * 180) / Math.PI
    if (i > 0) {
      while (r - unwrapped[i - 1] > 180) r -= 360
      while (r - unwrapped[i - 1] < -180) r += 360
    }
    unwrapped.push(r)
  }

  return {
    x: rawPoints.map((p) => `${p.x.toFixed(2)}vw`),
    y: rawPoints.map((p) => `${p.y.toFixed(2)}vh`),
    rotate: unwrapped,
    times: rawTimes.map((tt) => tt / duration),
    duration,
  }
}

function HomePage() {
  const thisYear = new Date().getFullYear()
  const reduced = useReducedMotion() ?? false
  const siteName = siteInfo.site.siteName
  const [guideDone, setGuideDone] = useState(false)
  const guidePath = useMemo(() => simulateGuidePath(), [])

  return (
    <div className="relative min-h-dvh overflow-hidden bg-background">
      {/* 背景图：晨雾中显现（模糊放大 → 清晰） */}
      <motion.div
        initial={
          reduced ? false : { opacity: 0, scale: 1.06, filter: 'blur(14px)' }
        }
        animate={{ opacity: 1, scale: 1, filter: 'blur(0px)' }}
        transition={{ duration: 1.3, ease: [0.25, 0.6, 0.35, 1] }}
        className="absolute inset-0 bg-cover bg-center bg-no-repeat"
        style={{ backgroundImage: `url(${defaultBackground})` }}
      />
      {/* 浅黄绿渐变遮罩：随背景稍后淡入 */}
      <motion.div
        initial={reduced ? false : { opacity: 0 }}
        animate={{ opacity: 1 }}
        transition={{
          duration: 1.5,
          delay: reduced ? 0 : 0.15,
          ease: 'easeOut',
        }}
        className="absolute inset-0"
        style={{
          background:
            'linear-gradient(135deg, var(--overlay-from) 0%, var(--overlay-via) 50%, var(--overlay-to) 100%)',
        }}
      />
      {/* 移动端额外提亮，让背景图隐约可见 */}
      <motion.div
        initial={reduced ? false : { opacity: 0 }}
        animate={{ opacity: 1 }}
        transition={{
          duration: 1.2,
          delay: reduced ? 0 : 0.3,
          ease: 'easeOut',
        }}
        className="absolute inset-0 md:hidden"
        style={{
          background:
            'radial-gradient(ellipse at center, transparent 0%, var(--overlay-from) 70%)',
        }}
      />

      {/* 落叶飘零装饰 */}
      <FallingLeaves />

      {/* 主内容区（开场叙事全在子元素上，容器保持静态） */}
      <div className="relative z-10 flex min-h-dvh flex-col items-center justify-center gap-10 px-6 py-12 lg:flex-row lg:gap-12 xl:gap-16">
        {/* 头像区域 */}
        <div className="relative shrink-0">
          {/* 引导叶：高速滑入 → 慢慢变慢 → 落定化作头像（相对头像定位，响应式天然适配） */}
          {!reduced && !guideDone && (
            <div className="pointer-events-none absolute inset-0 z-20 flex items-center justify-center">
              <motion.div
                initial={{
                  x: guidePath.x[0],
                  y: guidePath.y[0],
                  rotate: guidePath.rotate[0],
                }}
                animate={{
                  x: guidePath.x,
                  y: guidePath.y,
                  rotate: guidePath.rotate,
                }}
                transition={{
                  duration: guidePath.duration,
                  delay: 0.25,
                  times: guidePath.times,
                  ease: 'linear',
                }}
                onAnimationComplete={() => setGuideDone(true)}
              >
                {/* 显现 / 缩放层：滑入时淡入并缓缓逼近，落定时缩小消失 */}
                <motion.div
                  initial={{ opacity: 0, scale: 0.7 }}
                  animate={{
                    opacity: [0, 1, 1, 0],
                    scale: [0.7, 0.85, 1, 0.35],
                  }}
                  transition={{
                    duration: guidePath.duration,
                    delay: 0.25,
                    times: [0, 0.05, 0.9, 1],
                    ease: 'linear',
                  }}
                >
                  <svg
                    viewBox="0 0 48 32"
                    width={52}
                    height={34}
                    fill="var(--leaf-deep)"
                    aria-hidden="true"
                  >
                    <path d="M2 30C10 18 26 6 46 2c-3 12-16 24-44 28z" />
                    <path
                      d="M8 26C16 16 28 8 40 5"
                      stroke="oklch(0.96 0.03 110 / 0.5)"
                      strokeWidth="0.9"
                      fill="none"
                    />
                  </svg>
                </motion.div>
              </motion.div>
            </div>
          )}

          {/* 涟漪光环：头像落定后扩散一次，像叶子点了一下水面 */}
          {!reduced && (
            <motion.span
              className="pointer-events-none absolute inset-0 rounded-full border-2 border-primary/50"
              initial={{ scale: 1, opacity: 0 }}
              animate={{ scale: [1, 1.35, 1.9], opacity: [0, 0.55, 0] }}
              transition={{ delay: 2.2, duration: 1.1, ease: 'easeOut' }}
            />
          )}

          {/* 头像绽放：引导叶落定的瞬间弹出来 */}
          <motion.div
            {...enter(reduced, 1.68, {
              initial: { scale: 0, rotate: -20, opacity: 0 },
              animate: { scale: 1, rotate: 0, opacity: 1 },
              transition: { type: 'spring', stiffness: 260, damping: 16 },
            })}
          >
            <div className="relative">
              <img
                alt="UserAvatar"
                className="size-36 rounded-full object-cover ring-4 ring-ring-glow shadow-xl shadow-primary/20 lg:size-48 xl:size-56"
                src={myAvatar}
                draggable={false}
              />
              <motion.div
                className="absolute -bottom-2 -right-2 size-8 rounded-full bg-leaf-light/60 blur-md"
                {...enter(reduced, 2.35, {
                  initial: { opacity: 0 },
                  animate: { opacity: 1 },
                  transition: { duration: 0.6 },
                })}
              />
            </div>
          </motion.div>
        </div>

        {/* 内容区域 */}
        <div className="flex max-w-xl flex-col items-center gap-5 text-center">
          {/* 标题：Logo 像竹笋弹出 + 站名逐字从晨雾中浮现 */}
          <div className="flex items-center gap-3">
            <motion.span
              className="text-primary"
              {...enter(reduced, 1.88, {
                initial: { opacity: 0, scale: 0, rotate: -160 },
                animate: { opacity: 1, scale: 1, rotate: 0 },
                transition: { type: 'spring', stiffness: 200, damping: 14 },
              })}
            >
              <BambooLogo size={44} />
            </motion.span>
            <h1 className="text-3xl font-bold tracking-tight text-text-primary sm:text-4xl lg:text-5xl">
              {siteName.split('').map((char, i) => (
                <motion.span
                  key={`${char}-${i.toString()}`}
                  className="inline-block"
                  {...enter(reduced, 2.02 + i * 0.07, {
                    initial: {
                      opacity: 0,
                      y: 16,
                      rotate: 8,
                      filter: 'blur(6px)',
                    },
                    animate: {
                      opacity: 1,
                      y: 0,
                      rotate: 0,
                      filter: 'blur(0px)',
                    },
                    transition: { type: 'spring', stiffness: 180, damping: 18 },
                  })}
                >
                  {char === ' ' ? ' ' : char}
                </motion.span>
              ))}
            </h1>
          </div>

          {/* 渐变分隔线：从中心向两侧生长 */}
          <motion.div
            {...enter(reduced, 2.45, {
              initial: { scaleX: 0 },
              animate: { scaleX: 1 },
              transition: { type: 'spring', stiffness: 60, damping: 22 },
            })}
            className="h-px w-24 origin-center bg-gradient-to-r from-leaf-light via-leaf-muted to-leaf-light"
          />

          {/* 描述：晨雾对焦，从模糊到清晰 */}
          <motion.p
            {...enter(reduced, 2.55, {
              initial: { opacity: 0, y: 8, filter: 'blur(8px)' },
              animate: { opacity: 1, y: 0, filter: 'blur(0px)' },
              transition: { duration: 0.9, ease: 'easeOut' },
            })}
            className="max-w-2xl text-base leading-relaxed text-text-secondary lg:text-lg"
          >
            {siteInfo.blogger.description}
          </motion.p>

          {/* 按钮组：春笋破土，Q 弹错峰 */}
          <div className="flex flex-wrap justify-center gap-4 pt-4 lg:pt-6">
            <motion.div
              {...enter(reduced, 2.7, {
                initial: { opacity: 0, y: 26, scale: 0.92 },
                animate: { opacity: 1, y: 0, scale: 1 },
                transition: { type: 'spring', stiffness: 240, damping: 14 },
              })}
            >
              <a
                href="https://blog.x-lf.com"
                target="_blank"
                rel="noopener noreferrer"
                className="inline-flex items-center justify-center rounded-lg bg-primary px-8 py-3 text-sm font-medium text-primary-foreground shadow-lg shadow-primary/25 transition-colors hover:bg-primary-hover focus:outline-none focus:ring-2 focus:ring-ring focus:ring-offset-2"
              >
                去我的博客吧
              </a>
            </motion.div>
            <motion.div
              {...enter(reduced, 2.82, {
                initial: { opacity: 0, y: 26, scale: 0.92 },
                animate: { opacity: 1, y: 0, scale: 1 },
                transition: { type: 'spring', stiffness: 240, damping: 14 },
              })}
            >
              <Link
                to="/about"
                className="inline-flex items-center justify-center rounded-lg border border-primary/30 bg-white/80 px-8 py-3 text-sm font-medium text-primary shadow-sm transition-colors hover:bg-primary/10 focus:outline-none focus:ring-2 focus:ring-ring focus:ring-offset-2"
              >
                了解我的更多
              </Link>
            </motion.div>
          </div>
        </div>
      </div>

      {/* 桌面端页脚 */}
      <motion.footer
        {...enter(reduced, 3.0, {
          initial: { opacity: 0 },
          animate: { opacity: 0.7 },
          transition: { duration: 0.6, ease: 'easeOut' },
        })}
        className="absolute inset-x-0 bottom-0 z-10 hidden items-end justify-between p-4 text-sm text-text-secondary md:flex"
      >
        <div className="grid gap-1">
          <Link
            to="/auth/login"
            className="transition-colors hover:text-primary"
          >
            账户登录
          </Link>
          <span>
            Copyright (C) 2016-{thisYear} 筱锋xiao_lfeng. All Rights Reserved.
          </span>
        </div>
        <div className="grid gap-1 text-end">
          <a
            href="https://beian.miit.gov.cn/#/Integrated/index"
            target="_blank"
            rel="noopener noreferrer"
            className="transition-colors hover:text-primary"
          >
            粤ICP备 2022014822 号
          </a>
          <a
            href="https://beian.mps.gov.cn/#/query/webSearch"
            target="_blank"
            rel="noopener noreferrer"
            className="transition-colors hover:text-primary"
          >
            粤公网安备 44030702003207 号
          </a>
        </div>
      </motion.footer>

      {/* 移动端页脚 */}
      <motion.footer
        {...enter(reduced, 3.0, {
          initial: { opacity: 0 },
          animate: { opacity: 0.7 },
          transition: { duration: 0.6, ease: 'easeOut' },
        })}
        className="absolute inset-x-0 bottom-0 z-10 grid gap-1 pb-4 text-center text-sm text-text-secondary md:hidden"
      >
        <Link to="/auth/login" className="transition-colors hover:text-primary">
          账户登录
        </Link>
        <a
          href="https://beian.miit.gov.cn/#/Integrated/index"
          target="_blank"
          rel="noopener noreferrer"
          className="transition-colors hover:text-primary"
        >
          粤ICP备 2022014822 号
        </a>
      </motion.footer>
    </div>
  )
}
