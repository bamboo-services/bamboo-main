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

import * as React from 'react'
import { ChevronDown } from 'lucide-react'

import { cn } from '@/lib/utils'

/**
 * 原生 select 的字段纸化包装：与 Input / Textarea 共用同一套「字段纸」token
 * （field-surface 竹青纸底 + field-line 柔叶边 + leaf-deep 焦点光环），
 * 替换项目中散落的旧 `border-input bg-transparent` 内联样式，保证视觉一致。
 * 原生 option 面板保留系统渲染（macOS 原生下拉），仅本体换肤。
 */
function Select({
  className,
  children,
  ...props
}: React.ComponentProps<'select'>) {
  return (
    <div className="relative">
      <select
        data-slot="select"
        className={cn(
          'h-9 w-full cursor-pointer appearance-none rounded-md border border-field-line bg-field-surface pl-3 pr-8 text-base text-text-primary transition-[color,box-shadow,border-color] outline-none focus-visible:border-leaf-deep focus-visible:ring-ring/50 focus-visible:ring-[3px] disabled:cursor-not-allowed disabled:opacity-50 md:text-sm',
          className,
        )}
        {...props}
      >
        {children}
      </select>
      <ChevronDown className="pointer-events-none absolute right-2.5 top-1/2 size-4 -translate-y-1/2 text-text-secondary" />
    </div>
  )
}

export { Select }
