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

import { useMemo, useState } from 'react'
import { motion, useReducedMotion } from 'motion/react'

interface LeafConfig {
  /** 起始水平基准位置（vw） */
  baseX: number
  /** 叶片尺寸（px） */
  size: number
  /** 首次下落延迟（s），错开形成持续落叶 */
  delay: number
  /** 不透明度 —— 远景更淡 */
  opacity: number
  /** 叶片颜色（CSS 变量） */
  color: string
  /** 移动端是否隐藏（避免小屏过密） */
  hideOnMobile?: boolean
}

/** 八片叶子：大小 / 深浅各不相同，营造远近景深 */
const LEAVES: Array<LeafConfig> = [
  { baseX: 6, size: 24, delay: 0, opacity: 0.5, color: 'var(--leaf-deep)' },
  {
    baseX: 18,
    size: 17,
    delay: 4.5,
    opacity: 0.35,
    color: 'var(--leaf-muted)',
    hideOnMobile: true,
  },
  { baseX: 32, size: 28, delay: 7, opacity: 0.55, color: 'var(--leaf-deep)' },
  { baseX: 47, size: 19, delay: 2, opacity: 0.4, color: 'var(--leaf-light)' },
  {
    baseX: 60,
    size: 23,
    delay: 9,
    opacity: 0.5,
    color: 'var(--leaf-muted)',
    hideOnMobile: true,
  },
  { baseX: 72, size: 15, delay: 5.5, opacity: 0.3, color: 'var(--leaf-light)' },
  { baseX: 84, size: 26, delay: 1, opacity: 0.5, color: 'var(--leaf-deep)' },
  {
    baseX: 93,
    size: 18,
    delay: 10,
    opacity: 0.4,
    color: 'var(--leaf-muted)',
    hideOnMobile: true,
  },
]

/** 物理积分步长（s） */
const DT = 0.05
/** 轨迹采样间隔（s）—— 每 0.4s 采一个关键帧 */
const SAMPLE_INTERVAL = 0.4
/** 阻力 / 升力系数（视觉单位调谐） */
const FORCE_COEFF = 0.5
/** 重力加速度（vh/s²） */
const GRAVITY = 20

function rand(min: number, max: number) {
  return min + Math.random() * (max - min)
}

/** 把角度归一化到 (-π, π] */
function wrapAngle(angle: number): number {
  const tau = Math.PI * 2
  let next = angle
  while (next > Math.PI) next -= tau
  while (next < -Math.PI) next += tau
  return next
}

interface TrajectoryPoint {
  x: number
  y: number
  rotate: number
}

interface WindPath {
  /** 路径点（含起点） */
  points: Array<TrajectoryPoint>
  /** 归一化时间（0 → 1，严格递增） */
  times: Array<number>
  /** 本圈实际时长（s） */
  duration: number
}

/**
 * 落叶物理积分 —— 基于真实空气动力学模型
 *
 * 模型（参考 moruku36/falling-leaves-simulation）：
 * 1. 风场 = 基础风（连续低频正弦）+ 阵风 + 涡流（位置相关）
 * 2. 阻力：F_drag = ½ × Cd × A × v² × relVelocityDir（二次方阻力）
 * 3. 升力：F_lift = ½ × Cl × A × v² × sin(攻角 × 2) × normalDir
 *    叶片像机翼，攻角决定升力大小与方向，产生滑翔
 * 4. 旋转：风向标对齐力矩 + 颤振力矩 + 阻尼
 *    叶面趋向垂直于气流（像风向标），叠加正弦颤振，有阻尼不乱转
 * 5. 前向欧拉积分，dt=0.05s，每 0.4s 采一个关键帧
 *
 * 输出：交给 motion 用 ease: 'linear' 播放
 * （物理积分本身已是真实速度变化，关键帧足够密，无需额外缓冲）
 */
function simulateLeafPath(baseX: number): WindPath {
  // 物理参数（每片叶子每圈都重新随机，轨迹永不重复）
  const mass = rand(1, 2.5)
  const dragCoeff = rand(0.85, 1.25)
  const liftCoeff = rand(0.18, 0.48)
  const angularDrag = rand(1.8, 2.8)
  const alignStrength = rand(3.2, 5.6)
  const flutterStrength = rand(0.6, 1.4)
  const flutterRate = rand(3.5, 6.2)
  const turbOffset = rand(0, Math.PI * 2)

  // 风向：朝屏幕中心或对侧漂移（避免飞出屏幕太久）
  const windDir = baseX < 50 ? 1 : -1
  const windBase = rand(0.8, 2.2) * windDir

  // 初始状态
  let x = baseX + rand(-3, 3)
  let y = -8
  let vx = rand(-1, 1) + windBase * 0.4
  let vy = rand(1.5, 4)
  let rot = rand(-Math.PI, Math.PI)
  let angVel = rand(-1.0, 1.0)

  const points: Array<TrajectoryPoint> = [
    { x, y, rotate: (rot * 180) / Math.PI },
  ]
  const rawTimes: Array<number> = [0]
  let t = 0
  let lastSample = 0

  // 积分直到落地或超时
  while (y < 108 && t < 30) {
    t += DT

    // 风场：基础风 + 阵风（低频正弦）+ 涡流（位置相关）
    const windX =
      windBase +
      0.8 * Math.sin(t * 0.45 + turbOffset) +
      0.4 * Math.sin(y * 0.12 + t * 1.1 + turbOffset)
    const windY = 0.6 * Math.cos(t * 0.33 + turbOffset * 1.7)

    // 相对风速
    const relX = windX - vx
    const relY = windY - vy
    const relSpeed = Math.max(Math.hypot(relX, relY), 0.001)

    // 阻力（二次方，方向与相对速度相反）
    const dragMag = 0.5 * dragCoeff * FORCE_COEFF * relSpeed * relSpeed
    const dragX = (dragMag * relX) / relSpeed
    const dragY = (dragMag * relY) / relSpeed

    // 升力（叶片像机翼，攻角决定升力大小与方向）
    const flowAngle = Math.atan2(relY, relX)
    const angleOfAttack = wrapAngle(flowAngle - rot)
    const liftMag =
      0.5 *
      liftCoeff *
      FORCE_COEFF *
      relSpeed *
      relSpeed *
      Math.sin(angleOfAttack * 2)
    // 法向（垂直于相对风速，指向旋转后的左侧）
    const normalX = -relY / relSpeed
    const normalY = relX / relSpeed
    const liftX = normalX * liftMag
    const liftY = normalY * liftMag

    // 加速度
    const ax = (dragX + liftX) / mass
    const ay = (dragY + liftY + GRAVITY * mass) / mass

    // 前向欧拉积分
    vx += ax * DT
    vy += ay * DT
    x += vx * DT
    y += vy * DT

    // 旋转：风向标对齐（叶面趋向垂直于气流）+ 颤振 + 阻尼
    const broadsideAngle = flowAngle + Math.PI / 2
    const alignTorque = wrapAngle(broadsideAngle - rot) * alignStrength
    const flutterTorque =
      Math.sin(t * flutterRate + turbOffset) * flutterStrength
    const dampingTorque = -angVel * angularDrag
    angVel += (alignTorque + flutterTorque + dampingTorque) * DT
    rot = wrapAngle(rot + angVel * DT)

    // 采样
    if (t - lastSample >= SAMPLE_INTERVAL) {
      points.push({ x, y, rotate: (rot * 180) / Math.PI })
      rawTimes.push(t)
      lastSample = t
    }
  }

  // 确保最后一个点（落地时）
  if (rawTimes[rawTimes.length - 1] < t) {
    points.push({ x, y, rotate: (rot * 180) / Math.PI })
    rawTimes.push(t)
  }

  // 归一化 times（严格递增 0 → 1）
  const totalT = rawTimes[rawTimes.length - 1]
  const times = rawTimes.map((tt) => tt / totalT)

  return { points, times, duration: totalT }
}

/** 单片竹叶：细长披针形 + 叶脉 */
function LeafShape({ size, color }: { size: number; color: string }) {
  return (
    <svg
      viewBox="0 0 48 32"
      width={size}
      height={Math.round(size * 0.66)}
      fill={color}
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
  )
}

/** 单片落叶：每圈落完重新积分生成新路径，轨迹永不重复 */
function FallingLeaf({ leaf }: { leaf: LeafConfig }) {
  const [cycle, setCycle] = useState(0)

  const path = useMemo(() => simulateLeafPath(leaf.baseX), [cycle, leaf.baseX])

  return (
    <motion.div
      key={cycle}
      className={`absolute top-0 ${leaf.hideOnMobile ? 'hidden md:block' : ''}`}
      style={{ opacity: leaf.opacity }}
      initial={{
        x: `${path.points[0].x}vw`,
        y: `${path.points[0].y}vh`,
        rotate: path.points[0].rotate,
      }}
      animate={{
        x: path.points.map((p) => `${p.x}vw`),
        y: path.points.map((p) => `${p.y}vh`),
        rotate: path.points.map((p) => p.rotate),
      }}
      transition={{
        duration: path.duration,
        delay: cycle === 0 ? leaf.delay : 0,
        times: path.times,
        ease: 'linear',
      }}
      onAnimationComplete={() => setCycle((c) => c + 1)}
    >
      <LeafShape size={leaf.size} color={leaf.color} />
    </motion.div>
  )
}

/**
 * 落叶飘零 —— 真实空气动力学物理积分驱动，沉在内容层之下
 * 尊重 prefers-reduced-motion：用户偏好减少动画时不渲染
 */
export function FallingLeaves() {
  const shouldReduceMotion = useReducedMotion()

  if (shouldReduceMotion) {
    return null
  }

  return (
    <div
      className="pointer-events-none absolute inset-0 z-[5] overflow-hidden"
      aria-hidden="true"
    >
      {LEAVES.map((leaf) => (
        <FallingLeaf key={leaf.baseX} leaf={leaf} />
      ))}
    </div>
  )
}
