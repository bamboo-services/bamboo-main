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

import { ChevronLeft, ChevronRight, MoreHorizontal } from 'lucide-react'
import { cn } from '@/lib/utils'

interface PaginationProps {
  /** 当前页码（从 0 开始，与 TanStack Table 的 pageIndex 对齐） */
  pageIndex: number
  /** 总页数 */
  pageCount: number
  /** 页码变化回调（返回 0 基页码） */
  onPageChange: (pageIndex: number) => void
  className?: string
}

type PageItem = number | 'ellipsis-l' | 'ellipsis-r'

/**
 * 计算需要展示的页码序列（1 基，含省略号）。
 * 总页数 ≤ 7 时全部展示；否则固定首尾，围绕当前页取窗口，缺口用省略号补齐。
 */
function buildPageItems(current: number, total: number): Array<PageItem> {
  const cur = current + 1
  if (total <= 7) {
    return Array.from({ length: total }, (_, i) => i + 1)
  }
  const items: Array<PageItem> = [1]
  const left = Math.max(2, cur - 1)
  const right = Math.min(total - 1, cur + 1)
  if (left > 2) items.push('ellipsis-l')
  for (let p = left; p <= right; p++) items.push(p)
  if (right < total - 1) items.push('ellipsis-r')
  items.push(total)
  return items
}

/** 紧凑页码导航：小尺寸按钮 + 页码 + 省略号，当前页高亮 */
export function Pagination({
  pageIndex,
  pageCount,
  onPageChange,
  className,
}: PaginationProps) {
  const items = buildPageItems(pageIndex, pageCount)

  const navBtn =
    'inline-flex size-7 cursor-pointer items-center justify-center rounded-md text-muted-foreground transition-colors duration-150 hover:bg-muted hover:text-foreground disabled:pointer-events-none disabled:opacity-40'

  return (
    <nav
      className={cn('flex items-center gap-1', className)}
      aria-label="分页导航"
    >
      <button
        type="button"
        className={navBtn}
        onClick={() => onPageChange(pageIndex - 1)}
        disabled={pageIndex <= 0}
        aria-label="上一页"
      >
        <ChevronLeft className="size-3.5" />
      </button>

      {items.map((item) =>
        typeof item === 'number' ? (
          <button
            key={item}
            type="button"
            onClick={() => onPageChange(item - 1)}
            aria-current={item - 1 === pageIndex ? 'page' : undefined}
            aria-label={`第 ${item} 页`}
            className={cn(
              'inline-flex size-7 cursor-pointer items-center justify-center rounded-md text-xs font-medium tabular-nums transition-colors duration-150',
              item - 1 === pageIndex
                ? 'bg-primary text-primary-foreground shadow-[0_2px_8px_-2px_oklch(0.32_0.06_155/0.45)]'
                : 'text-muted-foreground hover:bg-muted hover:text-foreground',
            )}
          >
            {item}
          </button>
        ) : (
          <span
            key={item}
            className="inline-flex size-7 items-center justify-center text-muted-foreground"
            aria-hidden="true"
          >
            <MoreHorizontal className="size-3.5" />
          </span>
        ),
      )}

      <button
        type="button"
        className={navBtn}
        onClick={() => onPageChange(pageIndex + 1)}
        disabled={pageIndex >= pageCount - 1}
        aria-label="下一页"
      >
        <ChevronRight className="size-3.5" />
      </button>
    </nav>
  )
}
