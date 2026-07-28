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

import {  useEffect } from 'react'
import { animate, motion, useMotionValue, useReducedMotion, useTransform } from 'motion/react'
import type {ReactNode} from 'react';

/** 环形图分段：value 为绝对数值，color 为 CSS 颜色字符串（如 var(--chart-1)） */
export interface DonutSegment {
  value: number
  color: string
  label: string
}

interface DonutChartProps {
  segments: Array<DonutSegment>
  size?: number
  thickness?: number
  /** 圆心内容（通常放总数 + 标签） */
  center?: ReactNode
}

/**
 * 手搓 SVG 环形图。pathLength=1 归一化，每段以 transform rotate 定位、
 * strokeDasharray 控制可见比例，MotionValue 驱动「从无到有」的生长动画。
 */
export function DonutChart({
  segments,
  size = 200,
  thickness = 22,
  center,
}: DonutChartProps) {
  const reduced = useReducedMotion() ?? false
  const total = segments.reduce((sum, s) => sum + s.value, 0)
  const r = (size - thickness) / 2
  const cx = size / 2
  const cy = size / 2

  // 累计起点角度（度），用于把每段旋转到正确位置
  let acc = 0

  return (
    <div className="relative inline-flex" style={{ width: size, height: size }}>
      <svg width={size} height={size} viewBox={`0 0 ${size} ${size}`}>
        {/* 轨道底圈 */}
        <circle
          cx={cx}
          cy={cy}
          r={r}
          fill="none"
          stroke="var(--muted)"
          strokeWidth={thickness}
          opacity={0.45}
        />
        {total > 0 &&
          segments.map((seg, i) => {
            const frac = seg.value / total
            const rotate = acc * 360 - 90 // -90 让起点落在 12 点钟
            acc += frac
            return (
              <DonutArc
                key={`${seg.label}-${i.toString()}`}
                frac={frac}
                color={seg.color}
                rotate={rotate}
                reduced={reduced}
                delay={0.15 + i * 0.12}
                thickness={thickness}
                r={r}
                cx={cx}
                cy={cy}
              />
            )
          })}
      </svg>
      {center && (
        <div className="pointer-events-none absolute inset-0 flex flex-col items-center justify-center">
          {center}
        </div>
      )}
    </div>
  )
}

/** 单段弧：MotionValue 驱动 dasharray 从 0 生长到目标占比 */
function DonutArc({
  frac,
  color,
  rotate,
  reduced,
  delay,
  thickness,
  r,
  cx,
  cy,
}: {
  frac: number
  color: string
  rotate: number
  reduced: boolean
  delay: number
  thickness: number
  r: number
  cx: number
  cy: number
}) {
  const f = useMotionValue(reduced ? frac : 0)
  useEffect(() => {
    if (reduced) {
      f.set(frac)
      return
    }
    const controls = animate(f, frac, {
      duration: 0.9,
      ease: 'easeOut',
      delay,
    })
    return () => controls.stop()
  }, [frac, delay, reduced, f])
  // 留 0.002 间隙避免相邻段重叠出现锯齿
  const dash = useTransform(f, (v) => {
    const shown = Math.max(v - 0.002, 0)
    return `${shown} ${1 - shown}`
  })
  return (
    <g transform={`rotate(${rotate} ${cx} ${cy})`}>
      <motion.circle
        cx={cx}
        cy={cy}
        r={r}
        fill="none"
        stroke={color}
        strokeWidth={thickness}
        pathLength={1}
        style={{ strokeDasharray: dash }}
      />
    </g>
  )
}
