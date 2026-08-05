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

import { useMemo, useState } from 'react'
import { Link, createFileRoute, useNavigate } from '@tanstack/react-router'
import { motion, useReducedMotion } from 'motion/react'
import {
  flexRender,
  getCoreRowModel,
  useReactTable,
} from '@tanstack/react-table'
import {
  AlertTriangle,
  Camera,
  CheckCircle2,
  ExternalLink,
  List,
  MoreHorizontal,
  Pencil,
  Plus,
  Search,
  Table2,
  Trash2,
} from 'lucide-react'
import type {
  ColumnDef,
  OnChangeFn,
  PaginationState,
} from '@tanstack/react-table'
import type { LinkFriend, SnowflakeID } from '@/api/types'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Skeleton } from '@/components/ui/skeleton'
import { Pagination } from '@/components/ui/pagination'
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from '@/components/ui/dialog'
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from '@/components/ui/dropdown-menu'
import {
  BambooArt,
  BambooRule,
  EnsoEmpty,
  InkBadge,
  PageHead,
  inkCard,
  inkTableHeadRow,
  inkTableRow,
  inkTableWrap,
  inkTd,
  inkTh,
  linkStatus,
} from '@/components/ink-wash'
import { enter } from '@/lib/motion'
import { accentOf, isFancyColor, isPremiumColor } from '@/lib/colors'
import { AccentBar } from '@/components/link/accent-bar'
import { cn } from '@/lib/utils'
import {
  useAdminLinks,
  useDeleteLink,
  useReScreenshotLink,
} from '@/hooks/use-links'
import { useAllGroups } from '@/hooks/use-groups'
import { SiteAvatar } from '@/components/link/site-avatar'

export const Route = createFileRoute('/_admin/admin/link/')({
  component: LinkListPage,
})

type ViewMode = 'list' | 'table'

const PAGE_SIZE = 9

/** 行内操作菜单：编辑 + 重新截图（仅已通过）+ 删除（删除触发确认弹窗） */
function RowMenu({
  link,
  onConfirmDelete,
}: {
  link: LinkFriend
  onConfirmDelete: (link: LinkFriend) => void
}) {
  const reScreenshot = useReScreenshotLink()

  return (
    <DropdownMenu>
      <DropdownMenuTrigger asChild>
        <Button
          variant="ghost"
          size="icon"
          className="size-8 cursor-pointer opacity-50 transition-opacity hover:opacity-100 group-hover:opacity-100"
          onClick={(e) => e.stopPropagation()}
        >
          <MoreHorizontal className="size-4" />
        </Button>
      </DropdownMenuTrigger>
      <DropdownMenuContent align="end">
        <DropdownMenuItem asChild className="cursor-pointer">
          <Link
            to="/admin/link/$id/edit"
            params={{ id: link.id.toString() }}
            onClick={(e) => e.stopPropagation()}
          >
            <Pencil className="mr-2 size-4" />
            编辑
          </Link>
        </DropdownMenuItem>
        {link.status === 1 && (
          <>
            <DropdownMenuItem
              className="cursor-pointer"
              disabled={reScreenshot.isPending}
              onClick={(e) => {
                e.stopPropagation()
                reScreenshot.mutate(link.id)
              }}
            >
              <Camera className="mr-2 size-4" />
              重新截图
            </DropdownMenuItem>
          </>
        )}
        <DropdownMenuSeparator />
        <DropdownMenuItem
          className="cursor-pointer text-destructive focus:text-destructive"
          onClick={(e) => {
            e.stopPropagation()
            onConfirmDelete(link)
          }}
        >
          <Trash2 className="mr-2 size-4" />
          删除
        </DropdownMenuItem>
      </DropdownMenuContent>
    </DropdownMenu>
  )
}

/** 各列的响应式可见性 / 对齐 */
const columnClass: Record<string, string> = {
  group: 'hidden sm:table-cell',
  color: 'hidden md:table-cell',
  updatedAt: 'hidden lg:table-cell',
  actions: 'text-right',
}

/** 表格视图：服务端分页的无头表格 */
function LinkTable({
  links,
  pagination,
  pageCount,
  onPaginationChange,
  onOpenDetail,
  onConfirmDelete,
}: {
  links: Array<LinkFriend>
  pagination: PaginationState
  pageCount: number
  onPaginationChange: OnChangeFn<PaginationState>
  onOpenDetail: (id: SnowflakeID) => void
  onConfirmDelete: (link: LinkFriend) => void
}) {
  const columns = useMemo<Array<ColumnDef<LinkFriend>>>(
    () => [
      {
        id: 'site',
        accessorFn: (row) => row.name,
        header: '站点',
        cell: ({ row }) => {
          const link = row.original
          return (
            <div className="flex items-center gap-3">
              <SiteAvatar
                name={link.name}
                url={link.avatar}
                className="size-8 text-xs"
              />
              <div className="min-w-0">
                <div className="truncate font-serif font-semibold text-text-primary">
                  {link.name}
                </div>
                <div className="flex items-center gap-1 font-mono text-xs text-text-secondary">
                  <span className="truncate">{link.url}</span>
                  <ExternalLink className="size-3 shrink-0" />
                </div>
              </div>
            </div>
          )
        },
      },
      {
        id: 'group',
        header: '位置',
        cell: ({ row }) => (
          <InkBadge tone="neutral">
            {row.original.group_f_key?.name ?? '未分组'}
          </InkBadge>
        ),
      },
      {
        id: 'color',
        header: '颜色',
        cell: ({ row }) => {
          const color = row.original.color_f_key
          return (
            <span className="flex items-center gap-1.5 text-text-secondary">
              <span
                className="size-2.5 rounded-full"
                style={{ background: accentOf(color) }}
                aria-hidden="true"
              />
              {color?.name ?? '默认'}
              {isFancyColor(color) ? (
                <InkBadge tone="leaf">炫彩</InkBadge>
              ) : isPremiumColor(color) ? (
                <InkBadge tone="pending">高级</InkBadge>
              ) : color ? (
                <InkBadge tone="neutral">普通</InkBadge>
              ) : null}
            </span>
          )
        },
      },
      {
        id: 'status',
        header: '状态',
        cell: ({ row }) => {
          const s = linkStatus(row.original)
          return <InkBadge tone={s.tone}>{s.label}</InkBadge>
        },
      },
      {
        accessorKey: 'updated_at',
        id: 'updatedAt',
        header: '更新时间',
        cell: ({ row }) => (
          <span className="font-mono tabular-nums text-text-secondary">
            {new Date(row.original.updated_at).toLocaleDateString('zh-CN')}
          </span>
        ),
      },
      {
        id: 'actions',
        header: () => <span className="sr-only">操作</span>,
        cell: ({ row }) => (
          <RowMenu link={row.original} onConfirmDelete={onConfirmDelete} />
        ),
      },
    ],
    [onConfirmDelete],
  )

  const table = useReactTable({
    data: links,
    columns,
    state: { pagination },
    onPaginationChange,
    getCoreRowModel: getCoreRowModel(),
    manualPagination: true,
    pageCount,
  })

  return (
    <div className={inkTableWrap}>
      <table className="w-full text-sm">
        <thead>
          {table.getHeaderGroups().map((headerGroup) => (
            <tr key={headerGroup.id} className={inkTableHeadRow}>
              {headerGroup.headers.map((header) => (
                <th
                  key={header.id}
                  className={cn(inkTh, columnClass[header.column.id])}
                >
                  {flexRender(
                    header.column.columnDef.header,
                    header.getContext(),
                  )}
                </th>
              ))}
            </tr>
          ))}
        </thead>
        <tbody>
          {table.getRowModel().rows.map((row) => (
            <tr
              key={row.id}
              onClick={() => onOpenDetail(row.original.id)}
              className={cn('group cursor-pointer', inkTableRow)}
            >
              {row.getVisibleCells().map((cell) => (
                <td
                  key={cell.id}
                  className={cn(inkTd, columnClass[cell.column.id])}
                >
                  {flexRender(cell.column.columnDef.cell, cell.getContext())}
                </td>
              ))}
            </tr>
          ))}
        </tbody>
      </table>

      <div className="flex flex-wrap items-center justify-between gap-3 border-t border-border px-4 py-3">
        <span className="font-mono text-xs text-text-secondary">
          第 {pagination.pageIndex + 1} / {Math.max(pageCount, 1)} 页
        </span>
        <Pagination
          pageIndex={pagination.pageIndex}
          pageCount={Math.max(pageCount, 1)}
          onPageChange={(i) =>
            onPaginationChange((old) => ({ ...old, pageIndex: i }))
          }
        />
      </div>
    </div>
  )
}

function LinkListPage() {
  const reduced = useReducedMotion() ?? false
  const navigate = useNavigate()
  const [keyword, setKeyword] = useState('')
  const [groupFilter, setGroupFilter] = useState<SnowflakeID | null>(null)
  const [viewMode, setViewMode] = useState<ViewMode>('list')
  const [pagination, setPagination] = useState<PaginationState>({
    pageIndex: 0,
    pageSize: PAGE_SIZE,
  })
  const [deleteTarget, setDeleteTarget] = useState<LinkFriend | null>(null)

  const deleteLink = useDeleteLink()
  const groupsQuery = useAllGroups()

  const linksQuery = useAdminLinks({
    page: pagination.pageIndex + 1,
    page_size: pagination.pageSize,
    link_name: keyword.trim() || undefined,
    // 友链管理仅展示已通过友链，其余状态分流至审核/异常管理
    link_status: 1,
    link_group_id: groupFilter ?? undefined,
    sort_by: 'created_at',
    sort_order: 'desc',
  })

  const links = linksQuery.data?.data ?? []
  const total = linksQuery.data?.pagination.total ?? 0
  const totalPages = linksQuery.data?.pagination.total_pages ?? 1
  const groups = groupsQuery.data ?? []

  const openDetail = (id: SnowflakeID) =>
    navigate({ to: '/admin/link/$id', params: { id: id.toString() } })

  const handleDelete = () => {
    if (!deleteTarget) return
    deleteLink.mutate(deleteTarget.id, {
      onSuccess: () => setDeleteTarget(null),
    })
  }

  return (
    <div className="mx-auto max-w-6xl space-y-6">
      <PageHead
        kicker="LINKS · 友链"
        title="友链管理"
        sub="管理所有友情链接，审核申请、维护站点信息。"
        actions={
          <>
            <Link to="/admin/link/anomaly">
              <Button variant="outline" className="cursor-pointer">
                <AlertTriangle className="mr-2 size-4" />
                异常管理
                {(linksQuery.data?.anomaly_count ?? 0) > 0 && (
                  <InkBadge tone="danger" className="ml-2">
                    {linksQuery.data?.anomaly_count}
                  </InkBadge>
                )}
              </Button>
            </Link>
            <Link to="/admin/link/verify">
              <Button variant="outline" className="cursor-pointer">
                <CheckCircle2 className="mr-2 size-4" />
                友链审核
                {(linksQuery.data?.pending_count ?? 0) +
                  (linksQuery.data?.edit_pending_count ?? 0) >
                  0 && (
                  <InkBadge tone="pending" className="ml-2">
                    {(linksQuery.data?.pending_count ?? 0) +
                      (linksQuery.data?.edit_pending_count ?? 0)}
                  </InkBadge>
                )}
              </Button>
            </Link>
            <Link to="/admin/link/add">
              <Button className="cursor-pointer">
                <Plus className="mr-2 size-4" />
                添加友链
              </Button>
            </Link>
          </>
        }
      />

      <BambooRule reduced={reduced} delay={0.12} />

      {/* 搜索 + 位置分组 + 视图切换（一行） */}
      <motion.div
        {...enter(reduced, 0.18, {
          initial: { opacity: 0, y: 10 },
          animate: { opacity: 1, y: 0 },
          transition: { duration: 0.4, ease: 'easeOut' },
        })}
        className="flex flex-wrap items-center justify-between gap-3"
      >
        <div className="relative w-full max-w-xs sm:w-72">
          <Search className="absolute left-3 top-1/2 size-4 -translate-y-1/2 text-text-secondary" />
          <Input
            placeholder="搜索站点…"
            className="pl-9"
            value={keyword}
            onChange={(e) => {
              setKeyword(e.target.value)
              setPagination((old) => ({ ...old, pageIndex: 0 }))
            }}
          />
        </div>
        {groups.length > 0 && (
          <div className="flex min-w-0 flex-1 items-center justify-start">
            <div className="flex max-w-full items-center gap-1 overflow-x-auto rounded-lg border border-border bg-background p-0.5">
              <button
                type="button"
                onClick={() => {
                  setGroupFilter(null)
                  setPagination((old) => ({ ...old, pageIndex: 0 }))
                }}
                className={cn(
                  'cursor-pointer whitespace-nowrap rounded-md px-3 py-1.5 text-sm font-medium transition-colors duration-200',
                  groupFilter === null
                    ? 'bg-leaf-deep/12 text-leaf-deep'
                    : 'text-text-secondary hover:text-text-primary',
                )}
              >
                全部位置
              </button>
              {groups.map((group) => (
                <button
                  key={group.id.toString()}
                  type="button"
                  onClick={() => {
                    setGroupFilter(group.id)
                    setPagination((old) => ({ ...old, pageIndex: 0 }))
                  }}
                  className={cn(
                    'cursor-pointer whitespace-nowrap rounded-md px-3 py-1.5 text-sm font-medium transition-colors duration-200',
                    groupFilter === group.id
                      ? 'bg-leaf-deep/12 text-leaf-deep'
                      : 'text-text-secondary hover:text-text-primary',
                  )}
                >
                  {group.name}
                </button>
              ))}
            </div>
          </div>
        )}
        <div className="flex rounded-lg border border-border bg-background p-0.5">
          <button
            type="button"
            title="列表视图"
            onClick={() => setViewMode('list')}
            className={cn(
              'cursor-pointer rounded-md p-1.5 transition-colors duration-200',
              viewMode === 'list'
                ? 'bg-leaf-deep/12 text-leaf-deep'
                : 'text-text-secondary hover:text-text-primary',
            )}
          >
            <List className="size-4" />
          </button>
          <button
            type="button"
            title="表格视图"
            onClick={() => setViewMode('table')}
            className={cn(
              'cursor-pointer rounded-md p-1.5 transition-colors duration-200',
              viewMode === 'table'
                ? 'bg-leaf-deep/12 text-leaf-deep'
                : 'text-text-secondary hover:text-text-primary',
            )}
          >
            <Table2 className="size-4" />
          </button>
        </div>
      </motion.div>

      {/* 内容区 */}
      {linksQuery.isLoading ? (
        <div className="grid grid-cols-1 gap-3 sm:grid-cols-2 xl:grid-cols-3">
          {Array.from({ length: 6 }).map((_, i) => (
            <Skeleton key={i} className="h-44 w-full rounded-lg" />
          ))}
        </div>
      ) : links.length === 0 ? (
        <motion.div
          {...enter(reduced, 0.3, {
            initial: { opacity: 0 },
            animate: { opacity: 1 },
            transition: { duration: 0.4, ease: 'easeOut' },
          })}
          className={inkTableWrap}
        >
          <div className="py-6">
            <EnsoEmpty
              title="没有找到匹配的友链"
              hint="试试调整搜索关键词或筛选条件"
            >
              <Button
                variant="outline"
                size="sm"
                className="ml-auto cursor-pointer"
                onClick={() => {
                  setKeyword('')
                  setGroupFilter(null)
                }}
              >
                清除筛选
              </Button>
            </EnsoEmpty>
          </div>
        </motion.div>
      ) : viewMode === 'list' ? (
        <motion.div
          {...enter(reduced, 0.3, {
            initial: { opacity: 0 },
            animate: { opacity: 1 },
            transition: { duration: 0.4, ease: 'easeOut' },
          })}
          className="grid grid-cols-1 gap-3 sm:grid-cols-2 xl:grid-cols-3"
        >
          {links.map((link, i) => {
            const fancy = isFancyColor(link.color_f_key)
            return (
              <motion.div
                key={link.id.toString()}
                {...enter(reduced, 0.3 + i * 0.05, {
                  initial: { opacity: 0, y: 12 },
                  animate: { opacity: 1, y: 0 },
                  transition: { duration: 0.4, ease: 'easeOut' },
                })}
                role="button"
                tabIndex={0}
                onClick={() => openDetail(link.id)}
                onKeyDown={(e) => {
                  if (e.key === 'Enter' || e.key === ' ') openDetail(link.id)
                }}
                className={cn(
                  inkCard,
                  'group isolate cursor-pointer text-left',
                )}
              >
                {/* 炫彩卡背衬竹：右下角墨竹，hover 略深 */}
                {fancy && (
                  <BambooArt className="pointer-events-none absolute -z-10 bottom-0 right-[-18px] top-0 h-full w-[180px] text-text-primary opacity-80 transition-opacity duration-500 group-hover:opacity-100" />
                )}
                {/* 左侧站点色墨条：普通主色 hover 悬停色 / 高级渐变 / 炫彩竹影 */}
                <AccentBar color={link.color_f_key} className="inset-y-0 w-1" />
                <div className="flex flex-1 flex-col">
                  <div className="flex items-start justify-between gap-2">
                    <div className="flex min-w-0 items-center gap-3">
                      <SiteAvatar
                        name={link.name}
                        url={link.avatar}
                        className="size-11 text-sm"
                      />
                      <div className="min-w-0">
                        <div className="flex items-center gap-2">
                          <h3 className="truncate font-serif text-base font-semibold text-text-primary">
                            {link.name}
                          </h3>
                          <InkBadge tone="neutral" className="shrink-0">
                            {link.group_f_key?.name ?? '未分组'}
                          </InkBadge>
                        </div>
                        <a
                          href={link.url}
                          target="_blank"
                          rel="noopener noreferrer"
                          onClick={(e) => e.stopPropagation()}
                          className="mt-0.5 inline-flex max-w-full items-center gap-1 font-mono text-xs text-text-secondary transition-colors hover:text-leaf-deep"
                        >
                          <span className="truncate">{link.url}</span>
                          <ExternalLink className="size-3 shrink-0" />
                        </a>
                      </div>
                    </div>
                    <RowMenu link={link} onConfirmDelete={setDeleteTarget} />
                  </div>
                  <p
                    className="mt-3 flex-1 truncate text-sm text-text-secondary"
                    title={link.description ?? undefined}
                  >
                    {link.description ?? '这个站点没有留下描述。'}
                  </p>
                </div>
              </motion.div>
            )
          })}
        </motion.div>
      ) : (
        <motion.div
          {...enter(reduced, 0.3, {
            initial: { opacity: 0 },
            animate: { opacity: 1 },
            transition: { duration: 0.4, ease: 'easeOut' },
          })}
        >
          <LinkTable
            links={links}
            pagination={pagination}
            pageCount={totalPages}
            onPaginationChange={setPagination}
            onOpenDetail={openDetail}
            onConfirmDelete={setDeleteTarget}
          />
        </motion.div>
      )}

      {/* 列表视图分页 */}
      {viewMode === 'list' && links.length > 0 && (
        <div className="flex items-center justify-between">
          <span className="font-mono text-xs text-text-secondary">
            共 {total} 条
          </span>
          <Pagination
            pageIndex={pagination.pageIndex}
            pageCount={Math.max(totalPages, 1)}
            onPageChange={(i) =>
              setPagination((old) => ({ ...old, pageIndex: i }))
            }
          />
        </div>
      )}

      {/* 删除确认弹窗 */}
      <Dialog
        open={deleteTarget !== null}
        onOpenChange={(open) => {
          if (!open) setDeleteTarget(null)
        }}
      >
        <DialogContent>
          <DialogHeader>
            <DialogTitle>删除友链</DialogTitle>
            <DialogDescription>
              确定要删除友链「{deleteTarget?.name}」吗？此操作不可撤销。
            </DialogDescription>
          </DialogHeader>
          <DialogFooter>
            <Button
              variant="outline"
              onClick={() => setDeleteTarget(null)}
              disabled={deleteLink.isPending}
            >
              取消
            </Button>
            <Button
              variant="destructive"
              onClick={handleDelete}
              disabled={deleteLink.isPending}
            >
              {deleteLink.isPending ? '删除中…' : '确认删除'}
            </Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>
    </div>
  )
}
