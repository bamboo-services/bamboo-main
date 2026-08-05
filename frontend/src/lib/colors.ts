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

import type { LinkColor } from '@/api/types'

/** 内置炫彩颜色的保留 ID（与后端 constants.BuiltinFancyColorID 对齐） */
export const FANCY_COLOR_ID = 1n

/**
 * 竹绿炫彩渐变「竹影流光」：中深叶绿底色 + 顶部晨光晕，
 * 全部取自 styles.css 竹绿 token。纯渐变、不画竹——竹子作为卡背衬竹
 * 由 BambooArt 呈现，本函数仅供头像 / 色块 / 圆点等小元素使用。
 */
export function fancyGradient(): string {
  return [
    // 晨光晕：顶部偏右（用 leaf-muted 保证头像文字对比）
    'radial-gradient(140% 100% at 78% 0%, var(--leaf-muted) 0%, transparent 55%)',
    // 竹绿底色：上柔下深
    'linear-gradient(170deg, var(--leaf-muted) 0%, var(--leaf-deep) 58%)',
  ].join(', ')
}

/**
 * 判断颜色是否为炫彩。
 * 命中条件：保留 ID（color_id 直连兜底，炫彩为内置虚拟记录不落库）。
 */
export function isFancyColor(color: LinkColor | null | undefined): boolean {
  return color?.id === FANCY_COLOR_ID
}

/**
 * 取友链颜色主色/背景。
 * 炫彩返回竹绿渐变；普通颜色返回主色；未设置回退默认竹绿。
 */
export function accentOf(color: LinkColor | null | undefined): string {
  if (isFancyColor(color)) return fancyGradient()
  return color?.primary_color ?? 'var(--leaf-deep)'
}
