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

import { useState } from 'react'
import { Link, createFileRoute } from '@tanstack/react-router'
import { motion, useReducedMotion } from 'motion/react'
import {
  AlertTriangle,
  Camera,
  CheckCircle2,
  ExternalLink,
  Pencil,
} from 'lucide-react'
import { Button } from '@/components/ui/button'
import { Skeleton } from '@/components/ui/skeleton'
import { Pagination } from '@/components/ui/pagination'
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from '@/components/ui/table'
import {
  BambooRule,
  EnsoEmpty,
  InkBadge,
  PageHead,
  inkTableHeadRow,
  inkTableRow,
  inkTableWrap,
  inkTd,
  inkTh,
  linkStatus,
} from '@/components/ink-wash'
import { enter } from '@/lib/motion'
import { cn } from '@/lib/utils'
import {
  useAdminLinks,
  useReScreenshotLink,
  useUpdateLinkFail,
} from '@/hooks/use-links'

export const Route = createFileRoute('/_admin/admin/link/anomaly')({
  component: LinkAnomalyPage,
})

const PAGE_SIZE = 10

/** 头像加载失败时回退为首字色块，尺寸由 className 控制 */
function SiteAvatar({
  name,
  url,
  className,
}: {
  name: string
  url: string | null
  className?: string
}) {
  const [failed, setFailed] = useState(false)
  if (failed || !url) {
    return (
      <div
        className={cn(
          'flex shrink-0 items-center justify-center rounded-lg bg-muted font-serif font-semibold text-text-secondary',
          className,
        )}
      >
        {name.charAt(0)}
      </div>
    )
  }
  return (
    <div
      className={cn(
        'flex shrink-0 items-center justify-center overflow-hidden rounded-lg bg-muted ring-1 ring-border/60',
        className,
      )}
    >
      <img
        src={url}
        alt={name}
        loading="lazy"
        onError={() => setFailed(true)}
        className="size-full object-cover"
      />
    </div>
  )
}

/**
 * 异常管理：除「已通过 / 待审核」外的友链工作台
 * （已拒绝、下架待审核、已下架、已失效），支持恢复为正常与重新截图。
 */
function LinkAnomalyPage() {
  const reduced = useReducedMotion() ?? false
  const [pageIndex, setPageIndex] = useState(0)

  const linksQuery = useAdminLinks({
    page: pageIndex + 1,
    page_size: PAGE_SIZE,
    link_anomaly: true,
    sort_by: 'updated_at',
    sort_order: 'desc',
  })
  const restoreFail = useUpdateLinkFail()
  const reshot = useReScreenshotLink()

  const links = linksQuery.data?.data ?? []
  const total = linksQuery.data?.pagination.total ?? 0
  const totalPages = linksQuery.data?.pagination.total_pages ?? 1

  return (
    <div className="mx-auto max-w-6xl space-y-6">
      <PageHead
        kicker="ANOMALY · 异常"
        title="异常管理"
        sub="已拒绝、下架待审核、已下架与已失效友链集中处理。"
      />
      <BambooRule reduced={reduced} delay={0.12} />

      <motion.section
        {...enter(reduced, 0.18, {
          initial: { opacity: 0, y: 10 },
          animate: { opacity: 1, y: 0 },
          transition: { duration: 0.4, ease: 'easeOut' },
        })}
        className={inkTableWrap}
      >
        <Table>
          <TableHeader>
            <TableRow className={cn(inkTableHeadRow, 'hover:bg-muted/30')}>
              <TableHead className={inkTh}>站点</TableHead>
              <TableHead className={inkTh}>失效原因</TableHead>
              <TableHead className={inkTh}>状态</TableHead>
              <TableHead className={inkTh}>更新时间</TableHead>
              <TableHead className={cn(inkTh, 'text-right')}>操作</TableHead>
            </TableRow>
          </TableHeader>
          <TableBody>
            {linksQuery.isLoading ? (
              Array.from({ length: 4 }).map((_, i) => (
                <TableRow key={i}>
                  <TableCell className={inkTd}>
                    <Skeleton className="h-9 w-56" />
                  </TableCell>
                  <TableCell className={inkTd}>
                    <Skeleton className="h-4 w-32" />
                  </TableCell>
                  <TableCell className={inkTd}>
                    <Skeleton className="h-5 w-14 rounded-full" />
                  </TableCell>
                  <TableCell className={inkTd}>
                    <Skeleton className="h-4 w-28" />
                  </TableCell>
                  <TableCell className={cn(inkTd, 'text-right')}>
                    <div className="flex justify-end gap-1">
                      <Skeleton className="size-8 rounded-md" />
                      <Skeleton className="size-8 rounded-md" />
                      <Skeleton className="size-8 rounded-md" />
                    </div>
                  </TableCell>
                </TableRow>
              ))
            ) : links.length === 0 ? (
              <TableRow className="hover:bg-transparent">
                <TableCell colSpan={5} className="px-4 py-10">
                  <EnsoEmpty
                    title="暂无异常友链"
                    hint="失效的友链会出现在这里，方便你排查处理"
                  >
                    <AlertTriangle className="size-4 text-text-secondary" />
                  </EnsoEmpty>
                </TableCell>
              </TableRow>
            ) : (
              links.map((link) => {
                const s = linkStatus(link)
                return (
                  <TableRow key={link.id.toString()} className={inkTableRow}>
                    <TableCell className={inkTd}>
                      <div className="flex items-center gap-2.5">
                        <SiteAvatar
                          name={link.name}
                          url={link.avatar}
                          className="size-9 text-xs"
                        />
                        <div className="min-w-0">
                          <div className="truncate font-serif font-semibold text-text-primary">
                            {link.name}
                          </div>
                          <a
                            href={link.url}
                            target="_blank"
                            rel="noopener noreferrer"
                            className="mt-0.5 inline-flex max-w-full items-center gap-1 font-mono text-[11px] text-text-secondary transition-colors hover:text-leaf-deep"
                          >
                            <span className="truncate">{link.url}</span>
                            <ExternalLink className="size-3 shrink-0" />
                          </a>
                        </div>
                      </div>
                    </TableCell>
                    <TableCell className={inkTd}>
                      <span className="text-sm leading-snug text-text-secondary">
                        {link.fail_reason || '—'}
                      </span>
                    </TableCell>
                    <TableCell className={inkTd}>
                      <InkBadge tone={s.tone}>{s.label}</InkBadge>
                    </TableCell>
                    <TableCell className={inkTd}>
                      <span className="font-mono text-xs tabular-nums text-text-secondary">
                        {new Date(link.updated_at).toLocaleString('zh-CN')}
                      </span>
                    </TableCell>
                    <TableCell className={cn(inkTd, 'text-right')}>
                      <div className="flex items-center justify-end gap-1">
                        <Button
                          variant="ghost"
                          size="sm"
                          className="cursor-pointer text-leaf-deep"
                          disabled={restoreFail.isPending}
                          onClick={() =>
                            restoreFail.mutate({
                              id: link.id,
                              req: { link_fail: 0 },
                            })
                          }
                        >
                          <CheckCircle2 className="size-4" />
                          恢复
                        </Button>
                        <Button
                          variant="ghost"
                          size="sm"
                          className="cursor-pointer"
                          disabled={reshot.isPending}
                          onClick={() => reshot.mutate(link.id)}
                        >
                          <Camera className="size-4" />
                          重截图
                        </Button>
                        <Link
                          to="/admin/link/$id"
                          params={{ id: link.id.toString() }}
                        >
                          <Button
                            variant="ghost"
                            size="sm"
                            className="cursor-pointer"
                          >
                            <Pencil className="size-4" />
                            详情
                          </Button>
                        </Link>
                      </div>
                    </TableCell>
                  </TableRow>
                )
              })
            )}
          </TableBody>
        </Table>
      </motion.section>

      {!linksQuery.isLoading && links.length > 0 && (
        <div className="flex flex-wrap items-center justify-between gap-3">
          <span className="font-mono text-xs text-text-secondary">
            第 {pageIndex + 1} / {Math.max(totalPages, 1)} 页 · 共 {total} 条
          </span>
          <Pagination
            pageIndex={pageIndex}
            pageCount={Math.max(totalPages, 1)}
            onPageChange={setPageIndex}
          />
        </div>
      )}
    </div>
  )
}
