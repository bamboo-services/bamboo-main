// --------------------------------------------------------------------------------
// Copyright (c) 2016-NOW(至今) 筱锋
// Author: 筱锋「xiao_lfeng」(https://www.x-lf.com)
// --------------------------------------------------------------------------------
// 许可证声明：版权所有 (c) 2016-2026 筱锋。保留所有权利。
// 有关MIT许可证的更多信息，请查看项目根目录下的LICENSE文件或访问：
// https://opensource.org/licenses/MIT
// --------------------------------------------------------------------------------

import { Avatar, AvatarFallback, AvatarImage } from '@/components/ui/avatar'
import { useFriendOpen } from './friend-card-shared'
import type { FriendCardProps } from './friend-card-shared'

/**
 * 广告友链卡（1×1）—— 居中竖排，带「推广」标识。
 * leaf-muted 左墨条 2px（弱化，区别于有机友链），hover 加宽至 4px 并延伸全高。
 * 点击：触发 Interlude 沉浸引导。
 */
export function AdFriendCard({ link, onOpen }: FriendCardProps) {
  const { ref, handleClick } = useFriendOpen(link, onOpen)

  return (
    <a
      ref={ref}
      href={link.url}
      onClick={handleClick}
      className="group relative flex flex-col items-center justify-center overflow-hidden rounded-lg border border-border bg-card p-4 text-center transition-all duration-300 hover:-translate-y-0.5 hover:border-leaf-muted hover:shadow-[0_14px_30px_-22px_oklch(0.32_0.06_155/0.4)]"
    >
      {/* 「推广」标识 */}
      <span className="absolute right-2.5 top-2.5 z-20 rounded border border-leaf-muted/60 bg-leaf-light/20 px-1.5 py-0.5 font-mono text-[9px] uppercase tracking-[0.2em] text-leaf-deep">
        推广
      </span>

      {/* 左侧墨条：leaf-muted（弱化），hover 加宽并延伸全高 */}
      <span className="absolute inset-y-3 left-0 w-[2px] rounded-r-full bg-leaf-muted transition-all duration-500 group-hover:inset-y-0 group-hover:w-[4px]" />

      <Avatar className="size-11 rounded-full ring-1 ring-ring-glow">
        <AvatarImage src={link.avatar ?? undefined} alt={link.name} loading="lazy" />
        <AvatarFallback className="bg-leaf-light/50 font-serif text-base font-bold text-leaf-deep">
          {link.name.slice(0, 1)}
        </AvatarFallback>
      </Avatar>

      <h3 className="mt-2.5 truncate font-serif text-sm font-bold text-text-primary transition-colors group-hover:text-leaf-deep">
        {link.name}
      </h3>
      <p className="mt-0.5 truncate text-[11px] leading-snug text-text-secondary">
        {link.description ?? '这个站点很神秘，没有留下描述。'}
      </p>
    </a>
  )
}
