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

import { Toaster as Sonner } from 'sonner'
import type { ComponentProps } from 'react'

type ToasterProps = ComponentProps<typeof Sonner>

/** 全局 toast 容器：宣纸卡片底 + 衬线标题 + 类型左缘墨条（leaf-deep / destructive），样式对齐竹林水墨主题变量 */
function Toaster({ ...props }: ToasterProps) {
  return (
    <Sonner
      position="top-center"
      closeButton
      className="toaster group"
      toastOptions={{
        classNames: {
          toast:
            'group toast group-[.toaster]:bg-card group-[.toaster]:text-foreground group-[.toaster]:border group-[.toaster]:border-border group-[.toaster]:shadow-lg group-[.toaster]:rounded-lg',
          title: 'group-[.toast]:font-serif group-[.toast]:font-semibold',
          description: 'group-[.toast]:text-muted-foreground',
          actionButton:
            'group-[.toast]:bg-primary group-[.toast]:text-primary-foreground group-[.toast]:rounded-md',
          cancelButton:
            'group-[.toast]:bg-muted group-[.toast]:text-muted-foreground group-[.toast]:rounded-md',
          success:
            'group-[.toaster]:border-l-2 group-[.toaster]:border-l-leaf-deep',
          error:
            'group-[.toaster]:border-l-2 group-[.toaster]:border-l-destructive',
          warning:
            'group-[.toaster]:border-l-2 group-[.toaster]:border-l-leaf-muted',
          info: 'group-[.toaster]:border-l-2 group-[.toaster]:border-l-leaf-deep',
        },
      }}
      {...props}
    />
  )
}

export { Toaster }
