// --------------------------------------------------------------------------------
// Copyright (c) 2016-NOW(至今) 筱锋
// Author: 筱锋「xiao_lfeng」(https://www.x-lf.com)
// --------------------------------------------------------------------------------
// 许可证声明：版权所有 (c) 2016-2026 筱锋。保留所有权利。
// 有关MIT许可证的更多信息，请查看项目根目录下的 LICENSE文件或访问：
// https://opensource.org/licenses/MIT
// --------------------------------------------------------------------------------

import type { ReactNode } from 'react'

interface RadialGaugeProps {
  /** 0 ~ 1 之间的进度比例 */
  value: number
  size?: number
  thickness?: number
  /** CSS 颜色字符串（如 var(--chart-3)） */
  color?: string
  /** 圆心内容（百分比 / 数值） */
  center?: ReactNode
}

/**
 * 270° 仪表盘。轨道 + 数值弧，缺口朝下。
 * pathLength=1 归一化；可视角 = 0.75（270°）。value 弧占比 = value * 0.75。
 *
 * 弧线在数据到达即直接画到位，不再做自生长扫掠动画——避免与外层 section
 * 的入场动画叠加，产生「二次加载」观感（见 dashboard.tsx 的 enter() 编排）。
 */
export function RadialGauge({
  value,
  size = 160,
  thickness = 14,
  color = 'var(--chart-3)',
  center,
}: RadialGaugeProps) {
  const clamped = Math.min(1, Math.max(0, value))
  const r = (size - thickness) / 2
  const cx = size / 2
  const cy = size / 2
  // rotate(135) 让弧从 7:30 起点顺时针扫 270°，缺口落在底部
  const rotate = 135

  // 留 0.002 间隙避免与轨道端点重叠
  const shown = Math.max(clamped * 0.75 - 0.002, 0)
  const dash = `${shown} ${1 - shown}`

  return (
    <div className="relative inline-flex" style={{ width: size, height: size }}>
      <svg width={size} height={size} viewBox={`0 0 ${size} ${size}`}>
        <g transform={`rotate(${rotate} ${cx} ${cy})`}>
          {/* 轨道：270° 暗轨 */}
          <circle
            cx={cx}
            cy={cy}
            r={r}
            fill="none"
            stroke="var(--muted)"
            strokeWidth={thickness}
            strokeLinecap="round"
            pathLength={1}
            strokeDasharray="0.75 0.25"
            opacity={0.45}
          />
          {/* 数值弧：直接画到 value*0.75 */}
          <circle
            cx={cx}
            cy={cy}
            r={r}
            fill="none"
            stroke={color}
            strokeWidth={thickness}
            strokeLinecap="round"
            pathLength={1}
            strokeDasharray={dash}
          />
        </g>
      </svg>
      {center && (
        <div className="pointer-events-none absolute inset-0 flex flex-col items-center justify-center">
          {center}
        </div>
      )}
    </div>
  )
}
