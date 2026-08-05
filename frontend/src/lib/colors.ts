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

/** 高级配色类型值（与后端 constants.ColorTypePremium 对齐） */
export const PREMIUM_COLOR_TYPE = 1

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
 * 判断颜色是否为高级配色（type=1，三色渐变渲染）。
 */
export function isPremiumColor(color: LinkColor | null | undefined): boolean {
  return color?.type === PREMIUM_COLOR_TYPE
}

/**
 * 高级配色三色混合渐变：副色浅色起点 → 主色中段 → 悬停色深色强调。
 * 副/悬停色缺失时逐级退化（副色缺失→双色渐变；全缺→主色单色），
 * 保证即便数据不完整也能渲染，不抛出异常。
 */
export function premiumGradient(color: LinkColor | null | undefined): string {
  const primary = color?.primary_color
  if (!primary) return fancyGradient()
  const sub = color?.sub_color
  const hover = color?.hover_color
  if (sub && hover) {
    return `linear-gradient(160deg, ${sub} 0%, ${primary} 52%, ${hover} 100%)`
  }
  if (sub) {
    return `linear-gradient(160deg, ${sub} 0%, ${primary} 100%)`
  }
  return primary
}

/**
 * 取友链颜色主色/背景。
 * 炫彩返回竹绿渐变；高级配色返回三色混合渐变；普通颜色返回主色；未设置回退默认竹绿。
 */
export function accentOf(color: LinkColor | null | undefined): string {
  if (isFancyColor(color)) return fancyGradient()
  if (isPremiumColor(color)) return premiumGradient(color)
  return color?.primary_color ?? 'var(--leaf-deep)'
}
