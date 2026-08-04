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

import { Avatar, AvatarFallback, AvatarImage } from '@/components/ui/avatar'
import { cn } from '@/lib/utils'

/**
 * 站点头像：加载失败或缺失时回退为首字色块（宣纸底 + 墨色衬线字）。
 * 复用 shadcn Avatar 原语（与 about 友链卡同套机制），onError 回退由 Radix 承担。
 */
export function SiteAvatar({
  name,
  url,
  className,
}: {
  name: string
  url: string | null
  className?: string
}) {
  return (
    <Avatar
      className={cn(
        'shrink-0 rounded-lg bg-muted',
        url && 'ring-1 ring-border/60',
        className,
      )}
    >
      <AvatarImage src={url ?? undefined} alt={name} loading="lazy" />
      <AvatarFallback className="rounded-lg bg-muted font-serif font-semibold text-text-secondary">
        {name.charAt(0)}
      </AvatarFallback>
    </Avatar>
  )
}
