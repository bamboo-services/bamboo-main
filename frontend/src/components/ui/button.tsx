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

import * as React from 'react'
import { Slot } from '@radix-ui/react-slot'
import { cva } from 'class-variance-authority'
import type { VariantProps } from 'class-variance-authority'

import { cn } from '@/lib/utils'

const buttonVariants = cva(
  "inline-flex items-center justify-center gap-2 whitespace-nowrap rounded-md text-sm font-medium transition-all disabled:pointer-events-none disabled:opacity-50 [&_svg]:pointer-events-none [&_svg:not([class*='size-'])]:size-4 shrink-0 [&_svg]:shrink-0 outline-none focus-visible:border-ring focus-visible:ring-ring/50 focus-visible:ring-[3px] aria-invalid:ring-destructive/20 dark:aria-invalid:ring-destructive/40 aria-invalid:border-destructive",
  {
    variants: {
      variant: {
        // 主操作：墨色实底 + 墨边软阴影（hover 晕开 / active 微沉），唯一 accent 落地处
        default:
          'bg-primary text-primary-foreground shadow-[0_2px_10px_-3px_oklch(0.32_0.06_155/0.4)] hover:bg-primary-hover hover:shadow-[0_6px_16px_-6px_oklch(0.32_0.06_155/0.5)] active:translate-y-px',
        destructive:
          'bg-destructive text-destructive-foreground shadow-[0_2px_10px_-3px_oklch(0.577_0.245_27.325/0.4)] hover:bg-destructive/90 active:translate-y-px focus-visible:ring-destructive/20 dark:focus-visible:ring-destructive/40 dark:bg-destructive/60',
        // 宣纸卡底 + 墨边，hover 淡墨染边，与 inkCard 同源
        outline:
          'border border-border bg-card text-text-primary shadow-xs hover:bg-muted hover:border-leaf-muted active:translate-y-px dark:bg-input/30 dark:border-input dark:hover:bg-input/50',
        secondary:
          'bg-secondary text-secondary-foreground hover:bg-secondary/80 active:translate-y-px',
        ghost:
          'text-text-secondary hover:bg-muted hover:text-text-primary dark:hover:bg-accent/50',
        // 文字链：leaf-deep 墨色 + 笔刷感下划线（hover 时墨色描实）
        link: 'text-leaf-deep underline-offset-4 decoration-leaf-deep/40 hover:underline hover:decoration-leaf-deep',
      },
      size: {
        default: 'h-9 px-4 py-2 has-[>svg]:px-3',
        sm: 'h-8 rounded-md gap-1.5 px-3 has-[>svg]:px-2.5',
        lg: 'h-10 rounded-md px-6 has-[>svg]:px-4',
        icon: 'size-9',
        'icon-sm': 'size-8',
        'icon-lg': 'size-10',
      },
    },
    defaultVariants: {
      variant: 'default',
      size: 'default',
    },
  },
)

function Button({
  className,
  variant = 'default',
  size = 'default',
  asChild = false,
  ...props
}: React.ComponentProps<'button'> &
  VariantProps<typeof buttonVariants> & {
    asChild?: boolean
  }) {
  const Comp = asChild ? Slot : 'button'

  return (
    <Comp
      data-slot="button"
      data-variant={variant}
      data-size={size}
      className={cn(buttonVariants({ variant, size, className }))}
      {...props}
    />
  )
}

export { Button, buttonVariants }
