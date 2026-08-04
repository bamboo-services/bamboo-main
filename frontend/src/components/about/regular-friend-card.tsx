// --------------------------------------------------------------------------------
// Copyright (c) 2016-NOW(至今) 筱锋
// Author: 筱锋「xiao_lfeng」(https://www.x-lf.com)
// --------------------------------------------------------------------------------
// 许可证声明：版权所有 (c) 2016-2026 筱锋。保留所有权利。
// 有关MIT许可证的更多信息，请查看项目根目录下的LICENSE文件或访问：
// https://opensource.org/licenses/MIT
// --------------------------------------------------------------------------------

import { useFriendOpen } from './friend-card-shared'
import type { FriendCardProps } from './friend-card-shared'
import { Avatar, AvatarFallback, AvatarImage } from '@/components/ui/avatar'
import { accentOf, fancyGradient, isFancyColor } from '@/lib/colors'
import { cn } from '@/lib/utils'
import { BambooArt } from '@/components/ink-wash'

/**
 * 一般友链卡（1×1）—— 横排富式。
 * ring-glow 头像 + 衬线粗体名，左墨条 3px，hover 加宽至 5px 并延伸全高。
 * 点击：触发 Interlude 沉浸引导。
 */
export function RegularFriendCard({ link, onOpen }: FriendCardProps) {
  const { ref, handleClick } = useFriendOpen(link, onOpen)
  const accent = accentOf(link.color_f_key)
  const fancy = isFancyColor(link.color_f_key)

  return (
    <a
      ref={ref}
      href={link.url}
      onClick={handleClick}
      className="group isolate relative flex items-center gap-3 overflow-hidden rounded-lg border border-border bg-card p-4 transition-all duration-300 hover:-translate-y-0.5 hover:border-leaf-muted hover:shadow-[0_14px_30px_-22px_oklch(0.32_0.06_155/0.4)]"
      style={
        fancy
          ? {
              background:
                'radial-gradient(130% 100% at 88% 0%, oklch(0.88 0.1 105 / 0.16), transparent 58%), var(--card)',
            }
          : undefined
      }
    >
      {/* 炫彩卡背衬竹：右下角墨竹（主页同款 BambooArt），hover 略深 */}
      {fancy && (
        <BambooArt className="pointer-events-none absolute -z-10 bottom-0 right-[-14px] top-0 h-full w-[150px] text-text-primary opacity-80 transition-opacity duration-500 group-hover:opacity-100" />
      )}
      {/* 左侧墨条：友链主色（炫彩为竹影流光），hover 加宽并延伸全高 */}
      <span
        className={cn(
          'absolute inset-y-3 left-0 rounded-r-full transition-all duration-500 group-hover:inset-y-0 group-hover:w-[5px]',
          fancy ? 'w-[4px] ink-fancy' : 'w-[3px]',
        )}
        style={fancy ? undefined : { background: accent }}
      />

      <Avatar className="size-11 shrink-0 rounded-full ring-1 ring-ring-glow">
        <AvatarImage
          src={link.avatar ?? undefined}
          alt={link.name}
          loading="lazy"
        />
        <AvatarFallback
          className={cn(
            'font-serif text-base font-semibold',
            fancy ? 'text-card' : 'bg-leaf-light/40 text-leaf-deep',
          )}
          style={fancy ? { background: fancyGradient() } : undefined}
        >
          {link.name.slice(0, 1)}
        </AvatarFallback>
      </Avatar>

      <div className="min-w-0 flex-1">
        <h3 className="truncate font-serif text-[15px] font-bold text-text-primary transition-colors group-hover:text-leaf-deep">
          {link.name}
        </h3>
        <p className="mt-0.5 truncate text-xs leading-snug text-text-secondary">
          {link.description ?? '这个站点很神秘，没有留下描述。'}
        </p>
      </div>
    </a>
  )
}
