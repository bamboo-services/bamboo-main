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

import { useEffect, useState } from 'react'
import { useReducedMotion } from 'motion/react'
import type { CSSProperties } from 'react'

interface CountUpProps {
  value: number
  /** 动画时长（秒），默认 1s */
  duration?: number
  className?: string
  /** 透传给数字 span 的内联样式（如衬线字体） */
  style?: CSSProperties
}

/**
 * 数字渐进动画。基于 requestAnimationFrame + easeOutCubic。
 * 无障碍：用户偏好减少动态时，直接跳到终值。
 */
export function CountUp({
  value,
  duration = 1,
  className,
  style,
}: CountUpProps) {
  const reduced = useReducedMotion() ?? false
  const [display, setDisplay] = useState(reduced ? value : 0)

  useEffect(() => {
    if (reduced) {
      setDisplay(value)
      return
    }
    let raf = 0
    const start = performance.now()
    const tick = (now: number) => {
      const t = Math.min(1, (now - start) / (duration * 1000))
      const eased = 1 - Math.pow(1 - t, 3)
      setDisplay(Math.round(value * eased))
      if (t < 1) raf = requestAnimationFrame(tick)
    }
    raf = requestAnimationFrame(tick)
    return () => cancelAnimationFrame(raf)
  }, [value, duration, reduced])

  return (
    <span className={className} style={style}>
      {display.toLocaleString()}
    </span>
  )
}
