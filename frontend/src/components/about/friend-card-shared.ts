// --------------------------------------------------------------------------------
// Copyright (c) 2016-NOW(至今) 筱锋
// Author: 筱锋「xiao_lfeng」(https://www.x-lf.com)
// --------------------------------------------------------------------------------
// 许可证声明：版权所有 (c) 2016-2026 筱锋。保留所有权利。
// 有关MIT许可证的更多信息，请查看项目根目录下的LICENSE文件或访问：
// https://opensource.org/licenses/MIT
// --------------------------------------------------------------------------------

import { useRef } from 'react'
import type { MouseEvent } from 'react'
import type { LinkFriend } from '@/api/types'

/**
 * 友链卡共享管线：props 类型 + 「点击触发沉浸引导」hook + 域名提取。
 * 四张级别卡（premium/close/regular/ad）各自独立成组件，仅共用这套管线，
 * 渲染逻辑互不耦合。
 */

/** 四级友链卡统一的 props */
export interface FriendCardProps {
  link: LinkFriend
  /** 点击回调：传入 link 与卡片中心坐标，由父级触发 Interlude */
  onOpen: (link: LinkFriend, origin: { x: number; y: number }) => void
}

/** 从 URL 提取域名用于展示 */
export function domainOf(url: string): string {
  return url.replace(/^https?:\/\//, '').replace(/\/$/, '')
}

/**
 * 点击不直跳：拦截默认导航，取卡片几何中心作为 Interlude 的 clip-path 扩散起点。
 * 返回锚点 ref 与点击处理器。
 */
export function useFriendOpen(
  link: LinkFriend,
  onOpen: FriendCardProps['onOpen'],
) {
  const ref = useRef<HTMLAnchorElement>(null)

  const handleClick = (e: MouseEvent<HTMLAnchorElement>) => {
    e.preventDefault()
    const rect = ref.current?.getBoundingClientRect()
    const origin = rect
      ? { x: rect.left + rect.width / 2, y: rect.top + rect.height / 2 }
      : { x: window.innerWidth / 2, y: window.innerHeight / 2 }
    onOpen(link, origin)
  }

  return { ref, handleClick }
}
