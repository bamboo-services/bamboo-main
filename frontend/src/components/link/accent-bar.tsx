/*
 * --------------------------------------------------------------------------------
 * Copyright (c) 2016-NOW(至今) 筱锋
 * Author: 筱锋「xiao_lfeng」(https://www.x-lf.com)
 * --------------------------------------------------------------------------------
 * 许可证声明：版权所有 (c) 2016-2026 筱锋。保留所有权利。
 * 有关MIT许可证的更多信息，请查看项目根目录下的 LICENSE 文件或访问：
 * https://opensource.org/licenses/MIT
 * --------------------------------------------------------------------------------
 */

import type { LinkColor } from '@/api/types'
import {
  accentHoverOf,
  accentOf,
  isFancyColor,
  isPremiumColor,
} from '@/lib/colors'
import { cn } from '@/lib/utils'

/**
 * 站点色墨条（左缘竖条）—— 颜色分级体系的统一渲染：
 * - 炫彩：ink-fancy 竹影流光
 * - 高级：三色混合渐变
 * - 普通：主色 + hover 悬停色（需父级挂 `group` 类触发切换）
 *
 * CSS 变量定义在自身，仅供本墨条消费；若需联动子元素（如文字 hover 变色），
 * 请参照 about 友链卡将变量提升到卡片根元素。
 */
export function AccentBar({
  color,
  className,
}: {
  color: LinkColor | null | undefined
  className?: string
}) {
  const accent = accentOf(color)
  const fancy = isFancyColor(color)
  const premium = isPremiumColor(color)
  const hoverAccent = accentHoverOf(color)

  return (
    <span
      aria-hidden
      className={cn(
        'absolute left-0 rounded-r-full',
        fancy && 'ink-fancy',
        !fancy &&
          !premium &&
          'bg-[var(--ink-accent)] group-hover:bg-[var(--ink-accent-hover)]',
        !fancy && premium && '[background:var(--ink-accent)]',
        className,
      )}
      style={
        !fancy && !premium
          ? {
              ['--ink-accent' as string]: accent,
              ['--ink-accent-hover' as string]: hoverAccent ?? accent,
            }
          : undefined
      }
    />
  )
}
