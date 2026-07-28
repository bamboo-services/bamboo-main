// --------------------------------------------------------------------------------
// Copyright (c) 2016-NOW(至今) 筱锋
// Author: 筱锋「xiao_lfeng」(https://www.x-lf.com)
// --------------------------------------------------------------------------------
// 许可证声明：版权所有 (c) 2016-2026 筱锋。保留所有权利。
// 有关MIT许可证的更多信息，请查看项目根目录下的LICENSE文件或访问：
// https://opensource.org/licenses/MIT
// --------------------------------------------------------------------------------

import type { ComponentProps } from 'react'
import type { motion } from 'motion/react'

/** motion.div 的 props 类型别名，供 enter() 与各页面复用 */
export type MotionDivProps = ComponentProps<typeof motion.div>

/**
 * 入场动画助手：无障碍降级。
 *
 * 用户偏好减少动态时（reduced-motion），退化为快速淡入，
 * 时序整体压缩（每 0.08s 一拍）。否则原样返回 full 动画并叠加 delay。
 *
 * 此函数原先在 dashboard / 首页 / about 壳与三个子页面各自重复定义，
 * 现统一抽取于此，供所有页面复用，保证无障碍行为一致。
 */
export function enter(
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
