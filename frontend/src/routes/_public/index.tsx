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
import { useQuery } from '@tanstack/react-query'
import {
  AnimatePresence,
  motion,
  useMotionValue,
  useReducedMotion,
  useSpring,
  useTransform,
} from 'motion/react'
import { useEffect, useMemo, useRef, useState } from 'react'
import type { ComponentProps } from 'react'
import { BambooLogo } from '@/assets/svg/bamboo-logo'
import { FallingLeaves } from '@/components/decorative/falling-leaves'
import { BambooArt } from '@/components/ink-wash'
import { getSiteInfo } from '@/api/info'
import { getToken } from '@/lib/auth'
import { cn } from '@/lib/utils'
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

/* ------------------------------------------------------------------ */
/* 竹林景深：指针 → 归一化坐标，弹簧平滑后驱动各焦平面位移               */
/* ------------------------------------------------------------------ */

/**
 * 多焦平面视差 rig —— 模拟「相机对焦人物、环视竹林」的纵深。
 *
 * 指针归一化到 (-1..1)，经低刚度弹簧平滑（避免生硬跟手），再按「离相机越近
 * 位移越大」分配到各层：远景背景反向微移、近景边缘竹反向大幅位移，焦点正文
 * 保持锁定不动。reduced-motion 时不挂监听，各 motion value 恒为 0（静止）。
 */
function useGroveDepth(reduced: boolean) {
  const pointerX = useMotionValue(0)
  const pointerY = useMotionValue(0)

  useEffect(() => {
    if (reduced) return
    const onMove = (event: PointerEvent) => {
      pointerX.set((event.clientX / window.innerWidth) * 2 - 1)
      pointerY.set((event.clientY / window.innerHeight) * 2 - 1)
    }
    window.addEventListener('pointermove', onMove, { passive: true })
    return () => window.removeEventListener('pointermove', onMove)
  }, [reduced, pointerX, pointerY])

  const swayX = useSpring(pointerX, { stiffness: 55, damping: 18, mass: 0.9 })
  const swayY = useSpring(pointerY, { stiffness: 55, damping: 18, mass: 0.9 })

  return {
    /** 远景 · 背景图：随指针反向微移，位移最小 */
    bgX: useTransform(swayX, (v) => v * -8),
    bgY: useTransform(swayY, (v) => v * -5),
    /** 近景 · 边缘竹：随指针反向大幅位移，位移最大 */
    bambooX: useTransform(swayX, (v) => v * -22),
    bambooY: useTransform(swayY, (v) => v * -13),
  }
}

/* ------------------------------------------------------------------ */
/* Loading 界面：数据未就绪时展示的全屏呼吸动画                          */
/* ------------------------------------------------------------------ */

function LoadingScreen() {
  return (
    <motion.div
      key="loading"
      className="fixed inset-0 z-50 flex flex-col items-center justify-center gap-6 bg-background"
      initial={{ opacity: 1 }}
      exit={{ opacity: 0, scale: 1.02 }}
      transition={{ duration: 0.5, ease: 'easeInOut' }}
    >
      <motion.div
        animate={{ scale: [1, 1.08, 1], opacity: [0.85, 1, 0.85] }}
        transition={{
          duration: 2,
          repeat: Number.POSITIVE_INFINITY,
          ease: 'easeInOut',
        }}
      >
        <BambooLogo size={72} />
      </motion.div>
      <motion.p
        className="text-sm tracking-widest text-text-secondary"
        animate={{ opacity: [0.5, 1, 0.5] }}
        transition={{
          duration: 2,
          repeat: Number.POSITIVE_INFINITY,
          ease: 'easeInOut',
        }}
      >
        正在加载…
      </motion.p>
    </motion.div>
  )
}

/* ------------------------------------------------------------------ */
/* 打字机段落：逐字敲出 + 晨雾对焦，敲完后才放行后续元素入场             */
/* ------------------------------------------------------------------ */

interface TypewriterTextProps {
  text: string
  /** 开始敲击前的延迟（ms），用于配合整体入场时序 */
  delay?: number
  /** 临近敲完（约剩 500ms）时回调，父级据此提前放行按钮入场 */
  onNearEnd?: () => void
  /** 全部字符落墨完毕时回调，父级据此判定「中央内容加载完成」 */
  onDone?: () => void
  /** 段落排版类（默认正文 sans；信札体可传衬线大字） */
  className?: string
}

// 落墨节奏：标点/换行后停顿较长，普通字带确定性微抖动，避免等速机械感
const PUNCT_PAUSE = new Set([
  '，',
  '。',
  '、',
  '；',
  '：',
  '！',
  '？',
  ',',
  '.',
  ';',
  ':',
  '!',
  '?',
  '…',
  '—',
])

function charDelay(ch: string, index: number): number {
  if (ch === '\n') return 70
  if (PUNCT_PAUSE.has(ch)) return 55
  // 确定性抖动 [-2, 3]ms：快节奏下进一步收窄，落笔爽利而不凌乱
  const jitter = (((index * 2654435761) >>> 0) % 6) - 2
  return 12 + jitter
}

/**
 * 落墨收锋段落（逐字时间动画）。每个字预先排版占位（布局全程固定），打字机只控制
 * 「何时触发该字落墨」：未落墨的字以 blur + 透明 + 微沉降隐于纸面，触发后经 CSS
 * transition 自然收锋为清晰墨字——模糊是真高斯衰减、字是真写出来，单字自洽不穿帮。
 * 节奏非等速：标点与换行后停顿、普通字带确定性微抖动，模拟书写呼吸。
 * 临近落墨完成（约剩 500ms）回调 onNearEnd，父级据此提前放行按钮入场。
 * reduced-motion 时全文直接呈现、无动画。
 */
function TypewriterText({
  text,
  delay = 0,
  onNearEnd,
  onDone,
  className = 'max-w-2xl text-base leading-relaxed text-text-secondary lg:text-lg',
}: TypewriterTextProps) {
  const reduced = useReducedMotion() ?? false
  const chars = useMemo(() => Array.from(text), [text])
  const [count, setCount] = useState(reduced ? chars.length : 0)
  const nearEndRef = useRef(onNearEnd)
  const nearEndSent = useRef(false)
  nearEndRef.current = onNearEnd
  const doneRef = useRef(onDone)
  doneRef.current = onDone

  // 每个字落墨后到全部落墨完成的剩余时间，用于精确触发 onNearEnd（剩 ~500ms）
  const remaining = useMemo(() => {
    const rem = new Array<number>(chars.length).fill(0)
    for (let i = chars.length - 2; i >= 0; i--) {
      rem[i] = rem[i + 1] + charDelay(chars[i] ?? '', i)
    }
    return rem
  }, [chars])

  useEffect(() => {
    nearEndSent.current = false
    if (reduced || chars.length === 0) {
      setCount(chars.length)
      nearEndSent.current = true
      nearEndRef.current?.()
      doneRef.current?.()
      return
    }
    setCount(0)
    let k = -1
    let timer: ReturnType<typeof setTimeout>
    const step = () => {
      k += 1
      setCount(k + 1)
      if (!nearEndSent.current && (remaining[k] ?? 0) <= 500) {
        nearEndSent.current = true
        nearEndRef.current?.()
      }
      if (k >= chars.length - 1) {
        doneRef.current?.()
        return
      }
      timer = setTimeout(step, charDelay(chars[k] ?? '', k))
    }
    timer = setTimeout(step, delay)
    return () => clearTimeout(timer)
  }, [chars, reduced, delay, remaining])

  return (
    <p aria-label={text} className={`whitespace-pre-wrap ${className}`}>
      {chars.map((ch, i) => {
        if (ch === '\n') return <br key={`br-${i.toString()}`} aria-hidden />
        const shown = reduced || i < count
        return (
          <span
            key={`ch-${i.toString()}`}
            aria-hidden
            style={{
              display: 'inline-block',
              whiteSpace: 'pre',
              opacity: shown ? 1 : 0,
              filter: shown ? 'blur(0px)' : 'blur(8px)',
              transform: shown
                ? 'translateY(0) scale(1)'
                : 'translateY(0.14em) scale(1.08)',
              transition:
                'opacity 210ms ease-out, filter 240ms ease-out, transform 230ms cubic-bezier(0.22, 1, 0.36, 1)',
            }}
          >
            {ch}
          </span>
        )
      })}
    </p>
  )
}

/* ------------------------------------------------------------------ */
/* 主页内容：数据就绪后挂载，入场动画自然播放                            */
/* ------------------------------------------------------------------ */

/** 页脚站内导航（与 about/operate 顶导同源，公开路由） */
const HOME_NAV = [
  { to: '/about/me', label: '关于我' },
  { to: '/about/friends', label: '友链' },
  { to: '/about/sponsor', label: '赞助' },
  { to: '/operate/apply', label: '申请友链' },
] as const

/** 竹叶点角装饰（头像旁，与 about 名士帖同源） */
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

interface SealAvatarProps {
  reduced: boolean
  guideDone: boolean
  onGuideDone: () => void
  guidePath: GuidePath
}

/** 印章头像：引导叶滑入 → 墨晕涟漪 → 头像绽放（ring-glow 光环 + 竹叶点角） */
function SealAvatar({
  reduced,
  guideDone,
  onGuideDone,
  guidePath,
}: SealAvatarProps) {
  return (
    <div className="relative">
      {/* 引导叶：高速滑入 → 慢慢变慢 → 落定化作头像（相对头像定位） */}
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
              delay: 0.6,
              times: guidePath.times,
              ease: 'linear',
            }}
            onAnimationComplete={onGuideDone}
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
                delay: 0.6,
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

      {/* 墨晕涟漪：头像落定后淡绿墨晕扩散一次 */}
      {!reduced && (
        <motion.span
          className="pointer-events-none absolute inset-0 rounded-full"
          initial={{ scale: 0.9, opacity: 0 }}
          animate={{ scale: [0.9, 1.6, 2.1], opacity: [0, 0.5, 0] }}
          transition={{ delay: 2.35, duration: 1.3, ease: 'easeOut' }}
          style={{
            background:
              'radial-gradient(circle, oklch(0.88 0.1 105 / 0.4), transparent 68%)',
          }}
        />
      )}

      {/* 头像绽放：引导叶落定的瞬间弹出来 */}
      <motion.div
        {...enter(reduced, 2.1, {
          initial: { scale: 0, rotate: -20, opacity: 0 },
          animate: { scale: 1, rotate: 0, opacity: 1 },
          transition: { type: 'spring', stiffness: 260, damping: 16 },
        })}
        className="relative"
      >
        {/* 主体背光：柔和晨光光晕，随绽放一同弹出，把头像从竹雾中托出 */}
        <span
          aria-hidden
          className="pointer-events-none absolute left-1/2 top-1/2 size-[260%] -translate-x-1/2 -translate-y-1/2 rounded-full"
          style={{
            background:
              'radial-gradient(circle, oklch(0.93 0.06 110 / 0.5) 0%, oklch(0.93 0.06 110 / 0.22) 42%, transparent 68%)',
          }}
        />
        <img
          alt="UserAvatar"
          className="relative size-[clamp(5rem,16vh,9rem)] rounded-full object-cover shadow-xl shadow-primary/15 ring-4 ring-ring-glow lg:size-56"
          src={myAvatar}
          draggable={false}
        />
        <LeafAccent />
      </motion.div>
    </div>
  )
}

interface HomeContentProps {
  siteName: string
  description: string
}

function HomeContent({ siteName, description }: HomeContentProps) {
  const thisYear = new Date().getFullYear()
  const reduced = useReducedMotion() ?? false
  const [guideDone, setGuideDone] = useState(reduced)
  const [typeStarted, setTypeStarted] = useState(reduced)
  const [typeDone, setTypeDone] = useState(reduced)
  const [scrolled, setScrolled] = useState(false)
  const guidePath = useMemo(() => simulateGuidePath(), [])
  const isAuthed = getToken() != null
  const { bgX, bgY, bambooX, bambooY } = useGroveDepth(reduced)

  // 中央内容加载完成 = 信札落墨完毕 ∧ 引导叶落定（头像绽放）。
  // header / footer 此前保持隐身，待内容就绪后才入场收束画面。
  const contentDone = typeDone && guideDone

  useEffect(() => {
    const onScroll = () => setScrolled(window.scrollY > 40)
    onScroll()
    addEventListener('scroll', onScroll, { passive: true })
    return () => removeEventListener('scroll', onScroll)
  }, [])

  // 主/次 CTA：占位层与入场层共用，保证高度一致、布局不跳动
  const primaryCta = (
    <a
      href="https://blog.x-lf.com"
      target="_blank"
      rel="noopener noreferrer"
      className="inline-flex items-center justify-center rounded-lg bg-primary px-8 py-3 font-serif text-sm font-medium text-primary-foreground shadow-md shadow-primary/20 transition-[translate,background-color,box-shadow] duration-300 hover:-translate-y-0.5 hover:bg-primary-hover hover:shadow-lg hover:shadow-primary/25 focus:outline-none focus:ring-2 focus:ring-ring focus:ring-offset-2"
    >
      去我的博客吧
    </a>
  )
  const secondaryCta = (
    <Link
      to="/about"
      className="inline-flex items-center justify-center rounded-lg border border-leaf-deep/40 bg-transparent px-8 py-3 font-serif text-sm font-medium text-leaf-deep transition-[translate,background-color] duration-300 hover:-translate-y-0.5 hover:bg-leaf-deep/8 focus:outline-none focus:ring-2 focus:ring-ring focus:ring-offset-2"
    >
      了解我的更多
    </Link>
  )

  return (
    <motion.div
      key="content"
      className="relative flex h-dvh flex-col overflow-hidden bg-background"
      initial={{ opacity: 0 }}
      animate={{ opacity: 1 }}
      transition={{ duration: 0.4, ease: 'easeOut' }}
    >
      {/* 背景图 · 远景：晨雾中显现后常驻轻微焦外柔化（远景虚化），
          scale 1.04 预留视差位移量防露边 */}
      <motion.div
        initial={
          reduced ? false : { opacity: 0, scale: 1.08, filter: 'blur(16px)' }
        }
        animate={{ opacity: 1, scale: 1.04, filter: 'blur(2px)' }}
        transition={{ duration: 1.3, ease: [0.25, 0.6, 0.35, 1] }}
        className="absolute inset-0 bg-cover bg-center bg-no-repeat"
        style={{ backgroundImage: `url(${defaultBackground})`, x: bgX, y: bgY }}
      />
      {/* 宣纸遮罩：与 about/operate 同源全局 token，通透不压背景 */}
      <div
        className="absolute inset-0"
        style={{
          background:
            'linear-gradient(135deg, var(--overlay-from) 0%, var(--overlay-via) 50%, var(--overlay-to) 100%)',
        }}
      />
      {/* 晨光墨晕：单色淡绿径向（左上主光 + 右下辅光），轻盈透气 */}
      <motion.div
        initial={reduced ? false : { opacity: 0 }}
        animate={{ opacity: 1 }}
        transition={{
          duration: 1.6,
          delay: reduced ? 0 : 0.2,
          ease: 'easeOut',
        }}
        className="pointer-events-none absolute inset-0"
        style={{
          background:
            'radial-gradient(880px 520px at 14% 8%, oklch(0.88 0.1 105 / 0.16), transparent 70%), radial-gradient(760px 460px at 88% 92%, oklch(0.88 0.1 105 / 0.1), transparent 72%)',
        }}
      />
      {/* 流雾 · 上：中层空气透视，极慢横移，暗示竹林深处的湿气 */}
      <motion.div
        aria-hidden
        className="pointer-events-none absolute -left-[15%] -right-[15%] top-[13%] h-40"
        style={{
          background:
            'linear-gradient(90deg, transparent 4%, oklch(0.975 0.015 110 / 0.42) 32%, oklch(0.975 0.015 110 / 0.55) 52%, oklch(0.975 0.015 110 / 0.42) 72%, transparent 96%)',
          filter: 'blur(26px)',
        }}
        animate={reduced ? undefined : { x: ['0%', '-4%', '0%'] }}
        transition={{
          duration: 30,
          repeat: Number.POSITIVE_INFINITY,
          ease: 'easeInOut',
        }}
      />
      {/* 流雾 · 下：贴地晨雾，更浓更慢，与上层反向漂移 */}
      <motion.div
        aria-hidden
        className="pointer-events-none absolute -left-[15%] -right-[15%] bottom-[16%] h-52"
        style={{
          background:
            'linear-gradient(90deg, transparent 6%, oklch(0.97 0.02 110 / 0.5) 38%, oklch(0.97 0.02 110 / 0.62) 55%, oklch(0.97 0.02 110 / 0.5) 74%, transparent 94%)',
          filter: 'blur(30px)',
        }}
        animate={reduced ? undefined : { x: ['0%', '5%', '0%'] }}
        transition={{
          duration: 38,
          repeat: Number.POSITIVE_INFINITY,
          ease: 'easeInOut',
        }}
      />

      {/* 墨韵竹叶 · 近景：左右缘贴边，视差位移最大 + 焦外微虚化（离相机最近） */}
      <motion.div
        aria-hidden
        className="pointer-events-none absolute inset-0 [filter:blur(1.5px)]"
        style={{ x: bambooX, y: bambooY }}
      >
        <BambooArt
          className="absolute bottom-0 left-[-80px] top-0 hidden h-full w-[460px] text-text-primary opacity-90 lg:block"
          mirror
        />
        <BambooArt className="absolute bottom-0 right-[-80px] top-0 hidden h-full w-[460px] text-text-primary opacity-90 xl:block" />
      </motion.div>

      {/* 宣纸晕影：四周微沉、中央透亮，把视线收向名号（单色绿调，非暗角黑） */}
      <div
        aria-hidden
        className="pointer-events-none absolute inset-0"
        style={{
          background:
            'radial-gradient(115% 92% at 50% 44%, transparent 52%, oklch(0.88 0.06 120 / 0.22) 100%)',
        }}
      />

      {/* 落叶飘零装饰 */}
      <FallingLeaves />

      {/* 顶部导航：与 about/operate 同语法（滚动着墨） */}
      <motion.nav
        initial={false}
        animate={
          contentDone
            ? { opacity: 1, y: 0, filter: 'blur(0px)' }
            : { opacity: 0, y: -16, filter: 'blur(10px)' }
        }
        transition={
          reduced
            ? { duration: 0.3 }
            : {
                duration: 0.85,
                delay: 0.2,
                ease: [0.76, 0, 0.24, 1] as const,
              }
        }
        style={{ visibility: contentDone ? 'visible' : 'hidden' }}
        className={cn(
          'fixed inset-x-0 top-0 z-50 border-b transition-colors duration-400',
          scrolled
            ? 'border-border bg-background/94 backdrop-blur-[2px]'
            : 'border-transparent',
        )}
      >
        <div className="mx-auto flex w-full max-w-6xl items-center justify-between px-6 py-5 md:px-10">
          <Link
            to="/"
            className="flex items-center gap-2.5 transition-opacity hover:opacity-70"
          >
            <BambooLogo size={26} />
            <span className="font-serif text-base font-semibold tracking-wide text-text-primary">
              {siteName}
            </span>
          </Link>
          <div className="flex items-center gap-7 md:gap-9">
            {HOME_NAV.slice(0, 3).map((item) => (
              <Link
                key={item.to}
                to={item.to}
                className="relative font-mono text-[11px] uppercase tracking-[0.28em] text-text-secondary transition-colors after:absolute after:-bottom-[7px] after:left-0 after:h-0.5 after:w-full after:origin-left after:scale-x-0 after:bg-leaf-deep after:transition-transform after:duration-300 after:ease-[cubic-bezier(0.22,0.9,0.3,1)] hover:text-text-primary hover:after:scale-x-100"
              >
                {item.label}
              </Link>
            ))}
          </div>
        </div>
      </motion.nav>

      {/* 主内容：名士帖式 hero（12 栏，左 7 名号 / 右 5 头像），与 about/me 同语法 */}
      <main className="relative z-10 flex min-h-0 flex-1 flex-col pt-16 lg:pt-20">
        <section className="relative flex flex-1 items-center overflow-hidden">
          <div className="mx-auto w-full max-w-6xl px-6 md:px-10">
            <div className="grid w-full grid-cols-12 items-center gap-y-4 gap-x-0 pb-6 pt-6 lg:gap-x-10 lg:gap-y-10 lg:pb-12 lg:pt-10">
              {/* 左 7：名号 */}
              <div className="col-span-12 min-w-0 lg:col-span-7">
                {/* 双色巨型衬线站名：首字墨色 + 其余 leaf-deep，晨雾逐字浮现 */}
                <h1 className="font-serif font-bold leading-[1.05] tracking-[0.02em]">
                  {siteName.split('').map((char, i) => (
                    <motion.span
                      key={`${char}-${i.toString()}`}
                      className={`inline-block text-[clamp(3rem,9vw,6rem)] ${i === 0 ? 'text-text-primary' : 'text-leaf-deep'}`}
                      {...enter(reduced, 0.2 + i * 0.08, {
                        initial: {
                          opacity: 0,
                          y: 22,
                          rotate: 6,
                          filter: 'blur(6px)',
                        },
                        animate: {
                          opacity: 1,
                          y: 0,
                          rotate: 0,
                          filter: 'blur(0px)',
                        },
                        transition: {
                          duration: 0.7,
                          ease: [0.22, 0.9, 0.3, 1] as const,
                        },
                      })}
                    >
                      {char === ' ' ? '\u00A0' : char}
                    </motion.span>
                  ))}
                </h1>

                {/* 大写意笔刷 */}
                <motion.div
                  {...enter(reduced, 0.55, {
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

                {/* 信札正文：衬线打字机，逐字落墨收锋 */}
                <div className="mt-6 lg:mt-8">
                  <TypewriterText
                    text={description}
                    delay={reduced ? 0 : 1300}
                    onNearEnd={() => setTypeStarted(true)}
                    onDone={() => setTypeDone(true)}
                    className="max-w-lg text-left font-serif text-base leading-normal text-text-secondary md:text-xl md:leading-loose"
                  />
                </div>

                {/* CTA：始终挂载占位（布局不跳动），打字机临近结束时弹入 */}
                <div className="mt-6 flex flex-wrap items-center gap-4 lg:mt-10">
                  <motion.div
                    initial={false}
                    animate={
                      typeStarted
                        ? { opacity: 1, y: 0, scale: 1 }
                        : { opacity: 0, y: 26, scale: 0.92 }
                    }
                    transition={{ type: 'spring', stiffness: 300, damping: 18 }}
                    style={{ visibility: typeStarted ? 'visible' : 'hidden' }}
                  >
                    {primaryCta}
                  </motion.div>
                  <motion.div
                    initial={false}
                    animate={
                      typeStarted
                        ? { opacity: 1, y: 0, scale: 1 }
                        : { opacity: 0, y: 26, scale: 0.92 }
                    }
                    transition={{
                      type: 'spring',
                      stiffness: 300,
                      damping: 18,
                      delay: typeStarted ? 0.08 : 0,
                    }}
                  >
                    {secondaryCta}
                  </motion.div>
                </div>
              </div>

              {/* 右 5：印章头像（引导叶落定后绽放） */}
              <div className="col-span-12 flex min-w-0 justify-center lg:col-span-5 lg:justify-end">
                <SealAvatar
                  reduced={reduced}
                  guideDone={guideDone}
                  onGuideDone={() => setGuideDone(true)}
                  guidePath={guidePath}
                />
              </div>
            </div>
          </div>
        </section>
      </main>

      {/* 页脚：文档流内、居中（与 about/operate 同语法） */}
      <motion.footer
        initial={false}
        animate={
          contentDone
            ? { opacity: 1, y: 0, filter: 'blur(0px)' }
            : { opacity: 0, y: 16, filter: 'blur(10px)' }
        }
        transition={
          reduced
            ? { duration: 0.3 }
            : {
                duration: 0.85,
                delay: 0.2,
                ease: [0.76, 0, 0.24, 1] as const,
              }
        }
        style={{ visibility: contentDone ? 'visible' : 'hidden' }}
        className="relative z-10 shrink-0 border-t border-border/50"
      >
        <div className="mx-auto flex w-full max-w-6xl flex-col items-center gap-2 px-6 py-5 text-center md:px-10">
          <div className="flex flex-wrap items-center justify-center gap-x-5 gap-y-1">
            <a
              href="https://beian.miit.gov.cn/#/Integrated/index"
              target="_blank"
              rel="noopener noreferrer"
              className="font-mono text-[11px] tracking-[0.15em] text-text-secondary transition-colors hover:text-leaf-deep"
            >
              粤ICP备 2022014822 号
            </a>
            <a
              href="https://beian.mps.gov.cn/#/query/webSearch"
              target="_blank"
              rel="noopener noreferrer"
              className="font-mono text-[11px] tracking-[0.15em] text-text-secondary transition-colors hover:text-leaf-deep"
            >
              粤公网安备 44030702003207 号
            </a>
            <span
              className="hidden text-text-secondary/40 sm:inline"
              aria-hidden
            >
              ·
            </span>
            <Link
              to={isAuthed ? '/user/dashboard' : '/auth/login'}
              className="font-mono text-[11px] tracking-[0.15em] text-text-secondary/60 transition-colors hover:text-leaf-deep"
            >
              {isAuthed ? '用户中心' : '管理'}
            </Link>
          </div>
          <p className="font-mono text-[11px] tracking-[0.2em] text-text-secondary">
            Copyright (C) 2016-{thisYear} 筱锋 xiao_lfeng · All Rights Reserved
          </p>
        </div>
      </motion.footer>
    </motion.div>
  )
}

/* ------------------------------------------------------------------ */
/* 页面入口：数据加载 → loading / 内容 切换                             */
/* ------------------------------------------------------------------ */

function HomePage() {
  const { data, isLoading } = useQuery({
    queryKey: ['public', 'site'],
    queryFn: getSiteInfo,
  })

  const ready = !isLoading && data != null

  return (
    <AnimatePresence mode="wait">
      {ready ? (
        <HomeContent
          siteName={data.site_name}
          description={data.introduction}
        />
      ) : (
        <LoadingScreen />
      )}
    </AnimatePresence>
  )
}
