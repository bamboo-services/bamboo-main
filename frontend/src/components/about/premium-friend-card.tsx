// --------------------------------------------------------------------------------
// Copyright (c) 2016-NOW(至今) 筱锋
// Author: 筱锋「xiao_lfeng」(https://www.x-lf.com)
// --------------------------------------------------------------------------------
// 许可证声明：版权所有 (c) 2016-2026 筱锋。保留所有权利。
// 有关MIT许可证的更多信息，请查看项目根目录下的LICENSE文件或访问：
// https://opensource.org/licenses/MIT
// --------------------------------------------------------------------------------

import { ArrowRight, Eye } from 'lucide-react'
import { domainOf, useFriendOpen } from './friend-card-shared'
import { LazyImage } from './lazy-image'
import type { FriendCardProps } from './friend-card-shared'
import { Avatar, AvatarFallback, AvatarImage } from '@/components/ui/avatar'
import {
  accentHoverOf,
  accentOf,
  fancyGradient,
  isFancyColor,
  isPremiumColor,
} from '@/lib/colors'
import { cn } from '@/lib/utils'
import { BambooArt } from '@/components/ink-wash'

/**
 * 高级友链卡（2×2）—— Bento 栅格中的特写位。
 * 默认：头像 + 名/址 + 描述，水平垂直居中（名帖式）。
 * hover：左侧墨条加宽延伸全高，站点截图面板自底部滑入覆盖（浏览器框）。
 * 点击：触发 Interlude 沉浸引导（截图背景）；截图未生成或加载失败时回退占位。
 */
export function PremiumFriendCard({ link, onOpen }: FriendCardProps) {
  const { ref, handleClick } = useFriendOpen(link, onOpen)
  const accent = accentOf(link.color_f_key)
  const fancy = isFancyColor(link.color_f_key)
  const premium = isPremiumColor(link.color_f_key)
  const hoverAccent = accentHoverOf(link.color_f_key)

  return (
    <a
      ref={ref}
      href={link.url}
      onClick={handleClick}
      className="group isolate relative col-span-2 row-span-4 flex flex-col overflow-hidden rounded-lg border border-border bg-card transition-all duration-500 hover:-translate-y-1 hover:border-leaf-muted hover:shadow-[0_24px_50px_-30px_oklch(0.32_0.06_155/0.45)]"
      style={
        fancy
          ? {
              background:
                'radial-gradient(130% 100% at 88% 0%, oklch(0.88 0.1 105 / 0.18), transparent 58%), var(--card)',
            }
          : undefined
      }
    >
      {/* 炫彩卡背衬竹：右下角墨竹（大卡，竹幅稍宽），hover 略深 */}
      {fancy && (
        <BambooArt className="pointer-events-none absolute -z-10 bottom-0 right-[-22px] top-0 h-full w-[210px] text-text-primary opacity-80 transition-opacity duration-500 group-hover:opacity-100" />
      )}
      {/* 左侧墨条：炫彩竹影流光 / 高级色三色渐变 / 普通色主色，hover 加宽延伸并切悬停色 */}
      <span
        className={cn(
          'absolute inset-y-5 left-0 z-20 rounded-r-full transition-all duration-500 group-hover:inset-y-0 group-hover:w-[6px]',
          fancy ? 'w-[4px] ink-fancy' : 'w-[3px]',
          !fancy &&
            !premium &&
            'bg-[var(--ink-accent)] group-hover:bg-[var(--ink-accent-hover)]',
        )}
        style={
          fancy
            ? undefined
            : premium
              ? { background: accent }
              : {
                  ['--ink-accent' as string]: accent,
                  ['--ink-accent-hover' as string]: hoverAccent ?? accent,
                }
        }
      />

      {/* 默认内容：水平垂直居中，hover 淡出 */}
      <div className="relative z-10 flex h-full flex-col items-center justify-center p-6 text-center transition-opacity duration-300 group-hover:opacity-0">
        <Avatar className="size-16 rounded-full ring-2 ring-ring-glow">
          <AvatarImage
            src={link.avatar ?? undefined}
            alt={link.name}
            loading="lazy"
          />
          <AvatarFallback
            className={cn(
              'font-serif text-2xl font-bold',
              fancy ? 'text-card' : 'bg-leaf-light/50 text-leaf-deep',
            )}
            style={fancy ? { background: fancyGradient() } : undefined}
          >
            {link.name.slice(0, 1)}
          </AvatarFallback>
        </Avatar>
        <h3 className="mt-4 font-serif text-xl font-bold text-text-primary">
          {link.name}
        </h3>
        <p className="mt-1 font-mono text-[11px] uppercase tracking-[0.18em] text-text-secondary">
          {domainOf(link.url)}
        </p>
        <p className="mt-4 max-w-xs text-[14px] leading-[1.8] text-text-secondary">
          {link.description ?? '这个站点很神秘，没有留下描述。'}
        </p>
        <p className="mt-5 flex items-center gap-2 font-mono text-[10px] uppercase tracking-[0.3em] text-leaf-deep/70">
          <Eye className="size-3.5" />
          悬停查看站点截图
        </p>
      </div>

      {/* 站点截图面板：hover 自底部滑入覆盖 */}
      <div className="absolute inset-0 z-10 flex translate-y-[101%] flex-col transition-transform duration-500 ease-[cubic-bezier(0.22,0.9,0.3,1)] group-hover:translate-y-0">
        {/* 浏览器框 */}
        <div className="flex items-center gap-1.5 border-b border-border bg-card px-3.5 py-2">
          <span className="size-2.5 rounded-full bg-leaf-muted/60" />
          <span className="size-2.5 rounded-full bg-leaf-muted/40" />
          <span className="size-2.5 rounded-full bg-leaf-muted/25" />
          <span className="ml-2.5 truncate font-mono text-[10px] text-text-secondary">
            {link.url}
          </span>
        </div>
        {/* 预览区：真实站点截图（滚动接近时懒加载，未生成时墨晕占位） */}
        <div className="relative flex-1 overflow-hidden bg-card">
          {link.screenshot_url ? (
            <LazyImage
              src={link.screenshot_url}
              alt={`${link.name} 站点截图`}
              className="h-full w-full"
            />
          ) : (
            <div className="flex h-full items-center justify-center bg-gradient-to-br from-leaf-light/30 via-card to-leaf-muted/25">
              <span className="font-mono text-xs text-text-secondary">
                站点截图 · 占位
              </span>
            </div>
          )}
          {/* 底部名号叠层（渐变保可读） */}
          <div className="absolute inset-x-0 bottom-0 bg-gradient-to-t from-card via-card/85 to-transparent px-5 pb-4 pt-10">
            <div className="flex items-end justify-between gap-3">
              <div>
                <h3 className="font-serif text-lg font-bold text-text-primary">
                  {link.name}
                </h3>
                <p className="font-mono text-[10px] uppercase tracking-[0.18em] text-text-secondary">
                  {domainOf(link.url)}
                </p>
              </div>
              <span className="inline-flex items-center gap-1.5 font-mono text-[10px] uppercase tracking-[0.3em] text-leaf-deep">
                访问
                <ArrowRight className="size-3.5 transition-transform duration-300 group-hover:translate-x-0.5" />
              </span>
            </div>
          </div>
        </div>
      </div>
    </a>
  )
}
