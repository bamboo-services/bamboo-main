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

/** 取友链主色（后端嵌套的 color_f_key.primary_color），无则为 leaf-deep */
function colorOf(link: FriendCardProps['link']): string {
  return link.color_f_key?.primary_color ?? 'var(--leaf-deep)'
}

/**
 * 一般友链卡（1×1）—— 紧凑横排。
 * 友人主色左墨条 2px，hover 加宽至 4px 并延伸全高；衬线半粗名 hover→leaf-deep。
 * 点击：触发 Interlude 沉浸引导。
 */
export function RegularFriendCard({ link, onOpen }: FriendCardProps) {
  const { ref, handleClick } = useFriendOpen(link, onOpen)
  const accent = colorOf(link)

  return (
    <a
      ref={ref}
      href={link.url}
      onClick={handleClick}
      className="group relative flex items-center gap-3 overflow-hidden rounded-lg border border-border bg-card p-4 transition-all duration-300 hover:-translate-y-0.5 hover:border-leaf-muted hover:shadow-[0_14px_30px_-22px_oklch(0.32_0.06_155/0.4)]"
    >
      {/* 左侧墨条：友人主色，hover 加宽并延伸全高 */}
      <span
        className="absolute inset-y-3 left-0 w-[2px] rounded-r-full transition-all duration-500 group-hover:inset-y-0 group-hover:w-[4px]"
        style={{ backgroundColor: accent }}
      />

      <Avatar className="size-10 shrink-0 rounded-full">
        <AvatarImage src={link.avatar ?? undefined} alt={link.name} loading="lazy" />
        <AvatarFallback className="bg-leaf-light/30 font-serif text-sm font-semibold text-leaf-deep">
          {link.name.slice(0, 1)}
        </AvatarFallback>
      </Avatar>

      <div className="min-w-0 flex-1">
        <h3 className="truncate font-serif text-sm font-semibold text-text-primary transition-colors group-hover:text-leaf-deep">
          {link.name}
        </h3>
        <p className="mt-0.5 truncate text-xs leading-snug text-text-secondary">
          {link.description ?? '这个站点很神秘，没有留下描述。'}
        </p>
      </div>
    </a>
  )
}
