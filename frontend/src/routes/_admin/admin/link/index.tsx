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

import { useMemo, useState } from 'react'
import { Link, createFileRoute, useNavigate } from '@tanstack/react-router'
import {
  
  
  
  flexRender,
  getCoreRowModel,
  useReactTable
} from '@tanstack/react-table'
import {
  CheckCircle2,
  ExternalLink,
  List,
  MoreHorizontal,
  Pencil,
  Plus,
  Search,
  SearchX,
  Table2,
  Trash2,
} from 'lucide-react'
import type {ColumnDef, OnChangeFn, PaginationState} from '@tanstack/react-table';
import type { LinkFriend, SnowflakeID } from '@/api/types'
import { Button } from '@/components/ui/button'
import { Card, CardContent } from '@/components/ui/card'
import { Input } from '@/components/ui/input'
import { Badge } from '@/components/ui/badge'
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
import { cn } from '@/lib/utils'
import { useAdminLinks, useDeleteLink } from '@/hooks/use-links'
import { useAllGroups } from '@/hooks/use-groups'
import { useDashboardStats } from '@/hooks/use-dashboard'

export const Route = createFileRoute('/_admin/admin/link/')({
  component: LinkListPage,
})

type ViewMode = 'list' | 'table'

const PAGE_SIZE = 9

/** 友链状态徽章：0=待审核 1=已通过 2=已拒绝 */
function StatusBadge({ link }: { link: LinkFriend }) {
  if (link.is_failure === 1) {
    return <Badge variant="destructive">已失效</Badge>
  }
  switch (link.status) {
    case 1:
      return (
        <Badge className="bg-green-500/15 text-green-700 hover:bg-green-500/15">
          已通过
        </Badge>
      )
    case 2:
      return <Badge variant="destructive">已拒绝</Badge>
    default:
      return (
        <Badge className="bg-amber-500/15 text-amber-700 hover:bg-amber-500/15">
          待审核
        </Badge>
      )
  }
}

/** 头像加载失败时回退为首字色块 */
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
          'flex shrink-0 items-center justify-center rounded-lg bg-primary/10 font-semibold text-primary',
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
        className="size-full object-cover"
        loading="lazy"
        onError={() => setFailed(true)}
      />
    </div>
  )
}

/** 行内操作菜单：编辑 + 删除（删除触发确认弹窗） */
function RowMenu({
  link,
  onConfirmDelete,
}: {
  link: LinkFriend
  onConfirmDelete: (link: LinkFriend) => void
}) {
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
          <Link to="/admin/link/$id/edit" params={{ id: link.id.toString() }}>
            <Pencil className="mr-2 size-4" />
            编辑
          </Link>
        </DropdownMenuItem>
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
                <div className="truncate font-medium">{link.name}</div>
                <div className="flex items-center gap-1 text-xs text-muted-foreground">
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
          <Badge variant="secondary">
            {row.original.group_f_key?.name ?? '未分组'}
          </Badge>
        ),
      },
      {
        id: 'color',
        header: '颜色',
        cell: ({ row }) => {
          const color = row.original.color_f_key
          return (
            <span className="flex items-center gap-1.5 text-muted-foreground">
              <span
                className="size-2.5 rounded-full"
                style={{ backgroundColor: color?.primary_color ?? '#6366f1' }}
                aria-hidden="true"
              />
              {color?.name ?? '默认'}
            </span>
          )
        },
      },
      {
        id: 'status',
        header: '状态',
        cell: ({ row }) => <StatusBadge link={row.original} />,
      },
      {
        accessorKey: 'updated_at',
        id: 'updatedAt',
        header: '更新时间',
        cell: ({ row }) => (
          <span className="tabular-nums text-muted-foreground">
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
    <div className="overflow-hidden rounded-xl border border-border/70 bg-card">
      <table className="w-full text-sm">
        <thead>
          {table.getHeaderGroups().map((headerGroup) => (
            <tr
              key={headerGroup.id}
              className="border-b border-border/70 bg-muted/40 text-left text-xs text-muted-foreground"
            >
              {headerGroup.headers.map((header) => (
                <th
                  key={header.id}
                  className={cn(
                    'px-4 py-2.5 font-medium',
                    columnClass[header.column.id],
                  )}
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
              className="group cursor-pointer border-b border-border/40 transition-colors duration-150 last:border-0 hover:bg-muted/40"
            >
              {row.getVisibleCells().map((cell) => (
                <td
                  key={cell.id}
                  className={cn('px-4 py-2.5', columnClass[cell.column.id])}
                >
                  {flexRender(cell.column.columnDef.cell, cell.getContext())}
                </td>
              ))}
            </tr>
          ))}
        </tbody>
      </table>

      <div className="flex flex-wrap items-center justify-between gap-3 border-t border-border/70 px-4 py-3">
        <span className="text-sm text-muted-foreground">
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
  const navigate = useNavigate()
  const [keyword, setKeyword] = useState('')
  const [statusFilter, setStatusFilter] = useState<number | null>(null)
  const [groupFilter, setGroupFilter] = useState<SnowflakeID | null>(null)
  const [viewMode, setViewMode] = useState<ViewMode>('list')
  const [pagination, setPagination] = useState<PaginationState>({
    pageIndex: 0,
    pageSize: PAGE_SIZE,
  })
  const [deleteTarget, setDeleteTarget] = useState<LinkFriend | null>(null)

  const deleteLink = useDeleteLink()
  const groupsQuery = useAllGroups()
  const statsQuery = useDashboardStats()

  const linksQuery = useAdminLinks({
    page: pagination.pageIndex + 1,
    page_size: pagination.pageSize,
    link_name: keyword.trim() || undefined,
    link_status: statusFilter ?? undefined,
    link_group_id: groupFilter ?? undefined,
    sort_by: 'created_at',
    sort_order: 'desc',
  })

  const links = linksQuery.data?.data ?? []
  const total = linksQuery.data?.pagination.total ?? 0
  const totalPages = linksQuery.data?.pagination.total_pages ?? 1
  const stats = statsQuery.data
  const groups = groupsQuery.data ?? []

  const openDetail = (id: SnowflakeID) =>
    navigate({ to: '/admin/link/$id', params: { id: id.toString() } })

  const handleDelete = () => {
    if (!deleteTarget) return
    deleteLink.mutate(deleteTarget.id, {
      onSuccess: () => setDeleteTarget(null),
    })
  }

  const statusPills = [
    { label: '全部', value: null },
    { label: '待审核', value: 0 },
    { label: '已通过', value: 1 },
    { label: '已拒绝', value: 2 },
  ]

  return (
    <div className="space-y-5">
      {/* 页头 */}
      <div className="flex flex-wrap items-center justify-between gap-4">
        <div>
          <h1 className="text-3xl font-bold tracking-tight">友链管理</h1>
          <p className="mt-1 text-muted-foreground">管理所有友情链接</p>
        </div>
        <div className="flex gap-2">
          <Link to="/admin/link/verify">
            <Button variant="outline" className="cursor-pointer">
              <CheckCircle2 className="mr-2 size-4" />
              友链审核
              {(stats?.pending_links ?? 0) > 0 && (
                <Badge variant="destructive" className="ml-2">
                  {stats?.pending_links}
                </Badge>
              )}
            </Button>
          </Link>
          <Link to="/admin/link/add">
            <Button className="cursor-pointer">
              <Plus className="mr-2 size-4" />
              添加友链
            </Button>
          </Link>
        </div>
      </div>

      {/* 紧凑统计条 + 视图切换 */}
      <div className="flex flex-wrap items-center justify-between gap-3">
        <div className="flex items-center gap-5 text-sm">
          <span className="flex items-center gap-1.5">
            <span
              className="size-2 rounded-full"
              style={{ backgroundColor: '#0891b2' }}
              aria-hidden="true"
            />
            <span className="font-semibold tabular-nums">
              {stats?.total_links ?? 0}
            </span>
            <span className="text-muted-foreground">总收录</span>
          </span>
          <span className="flex items-center gap-1.5">
            <span
              className="size-2 rounded-full"
              style={{ backgroundColor: '#22c55e' }}
              aria-hidden="true"
            />
            <span className="font-semibold tabular-nums">
              {stats?.approved_links ?? 0}
            </span>
            <span className="text-muted-foreground">已通过</span>
          </span>
          <span className="flex items-center gap-1.5">
            <span
              className="size-2 rounded-full"
              style={{ backgroundColor: '#f59e0b' }}
              aria-hidden="true"
            />
            <span className="font-semibold tabular-nums">
              {stats?.pending_links ?? 0}
            </span>
            <span className="text-muted-foreground">待审核</span>
          </span>
        </div>

        <div className="flex rounded-lg border border-input bg-background p-0.5">
          <button
            type="button"
            title="列表视图"
            onClick={() => setViewMode('list')}
            className={cn(
              'cursor-pointer rounded-md p-1.5 transition-colors duration-200',
              viewMode === 'list'
                ? 'bg-muted text-foreground shadow-xs'
                : 'text-muted-foreground hover:text-foreground',
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
                ? 'bg-muted text-foreground shadow-xs'
                : 'text-muted-foreground hover:text-foreground',
            )}
          >
            <Table2 className="size-4" />
          </button>
        </div>
      </div>

      {/* 筛选区：搜索 + 状态 + 分组 */}
      <div className="flex flex-wrap items-center gap-3">
        <div className="relative w-full max-w-xs">
          <Search className="absolute left-3 top-1/2 size-4 -translate-y-1/2 text-muted-foreground" />
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
        <div className="flex flex-wrap gap-1.5">
          {statusPills.map((pill) => (
            <button
              key={pill.label}
              type="button"
              onClick={() => {
                setStatusFilter(pill.value)
                setPagination((old) => ({ ...old, pageIndex: 0 }))
              }}
              className={cn(
                'cursor-pointer rounded-full px-3.5 py-1.5 text-sm font-medium transition-colors duration-200',
                statusFilter === pill.value
                  ? 'bg-primary text-primary-foreground shadow-sm'
                  : 'bg-muted text-muted-foreground hover:text-foreground',
              )}
            >
              {pill.label}
            </button>
          ))}
        </div>
      </div>

      {groups.length > 0 && (
        <div className="flex flex-wrap gap-1.5">
          <button
            type="button"
            onClick={() => {
              setGroupFilter(null)
              setPagination((old) => ({ ...old, pageIndex: 0 }))
            }}
            className={cn(
              'cursor-pointer rounded-full px-3.5 py-1.5 text-sm font-medium transition-colors duration-200',
              groupFilter === null
                ? 'bg-foreground text-background shadow-sm'
                : 'bg-muted text-muted-foreground hover:text-foreground',
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
                'cursor-pointer rounded-full px-3.5 py-1.5 text-sm font-medium transition-colors duration-200',
                groupFilter === group.id
                  ? 'bg-foreground text-background shadow-sm'
                  : 'bg-muted text-muted-foreground hover:text-foreground',
              )}
            >
              {group.name}
            </button>
          ))}
        </div>
      )}

      {/* 内容区 */}
      {linksQuery.isLoading ? (
        <div className="grid grid-cols-1 gap-3 sm:grid-cols-2 xl:grid-cols-3">
          {Array.from({ length: 6 }).map((_, i) => (
            <Skeleton key={i} className="h-44 w-full rounded-xl" />
          ))}
        </div>
      ) : links.length === 0 ? (
        <Card className="border-dashed">
          <CardContent className="flex flex-col items-center justify-center gap-3 py-14 text-center">
            <div className="flex size-12 items-center justify-center rounded-full bg-muted">
              <SearchX className="size-5 text-muted-foreground" />
            </div>
            <div>
              <p className="font-medium">没有找到匹配的友链</p>
              <p className="mt-1 text-sm text-muted-foreground">
                试试调整搜索关键词或筛选条件
              </p>
            </div>
            <Button
              variant="outline"
              size="sm"
              className="cursor-pointer"
              onClick={() => {
                setKeyword('')
                setStatusFilter(null)
                setGroupFilter(null)
              }}
            >
              清除筛选
            </Button>
          </CardContent>
        </Card>
      ) : viewMode === 'list' ? (
        <div className="grid grid-cols-1 gap-3 sm:grid-cols-2 xl:grid-cols-3">
          {links.map((link) => (
            <div
              key={link.id.toString()}
              role="button"
              tabIndex={0}
              onClick={() => openDetail(link.id)}
              onKeyDown={(e) => {
                if (e.key === 'Enter' || e.key === ' ') openDetail(link.id)
              }}
              className="group relative flex w-full cursor-pointer flex-col overflow-hidden rounded-xl border border-border/70 bg-card text-left transition-all duration-200 hover:-translate-y-0.5 hover:border-primary/40 hover:shadow-md"
            >
              <span
                className="h-1 w-full"
                style={{
                  backgroundColor: link.color_f_key?.primary_color ?? '#6366f1',
                }}
                aria-hidden="true"
              />
              <div className="flex flex-1 flex-col p-4">
                <div className="flex items-start justify-between gap-2">
                  <div className="flex min-w-0 items-center gap-3">
                    <SiteAvatar
                      name={link.name}
                      url={link.avatar}
                      className="size-11 text-sm"
                    />
                    <div className="min-w-0">
                      <h3 className="truncate font-semibold">{link.name}</h3>
                      <a
                        href={link.url}
                        target="_blank"
                        rel="noopener noreferrer"
                        onClick={(e) => e.stopPropagation()}
                        className="mt-0.5 inline-flex max-w-full items-center gap-1 text-xs text-muted-foreground transition-colors hover:text-primary"
                      >
                        <span className="truncate">{link.url}</span>
                        <ExternalLink className="size-3 shrink-0" />
                      </a>
                    </div>
                  </div>
                  <RowMenu link={link} onConfirmDelete={setDeleteTarget} />
                </div>

                <p className="mt-3 line-clamp-2 flex-1 text-sm leading-relaxed text-muted-foreground">
                  {link.description ?? '这个站点没有留下描述。'}
                </p>

                <div className="mt-3 flex flex-wrap items-center gap-1.5">
                  <Badge variant="secondary">
                    {link.group_f_key?.name ?? '未分组'}
                  </Badge>
                  <Badge variant="outline" className="gap-1.5">
                    <span
                      className="size-2 rounded-full"
                      style={{
                        backgroundColor:
                          link.color_f_key?.primary_color ?? '#6366f1',
                      }}
                      aria-hidden="true"
                    />
                    {link.color_f_key?.name ?? '默认'}
                  </Badge>
                  <StatusBadge link={link} />
                </div>
              </div>
            </div>
          ))}
        </div>
      ) : (
        <LinkTable
          links={links}
          pagination={pagination}
          pageCount={totalPages}
          onPaginationChange={setPagination}
          onOpenDetail={openDetail}
          onConfirmDelete={setDeleteTarget}
        />
      )}

      {/* 列表视图分页 */}
      {viewMode === 'list' && links.length > 0 && (
        <div className="flex items-center justify-between">
          <span className="text-sm text-muted-foreground">共 {total} 条</span>
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
