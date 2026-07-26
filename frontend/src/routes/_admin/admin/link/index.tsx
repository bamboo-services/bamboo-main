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
  type ColumnDef,
  type PaginationState,
  type SortingState,
  flexRender,
  getCoreRowModel,
  getPaginationRowModel,
  getSortedRowModel,
  useReactTable,
} from '@tanstack/react-table'
import {
  CheckCircle2,
  ChevronDown,
  ChevronUp,
  ChevronsUpDown,
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
import { Button } from '@/components/ui/button'
import { Card, CardContent } from '@/components/ui/card'
import { Input } from '@/components/ui/input'
import { Badge } from '@/components/ui/badge'
import { Pagination } from '@/components/ui/pagination'
import { mockColors, mockLinks, mockLocations } from '@/data/mock/links'
import type { LinkItem } from '@/data/mock/links'
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from '@/components/ui/dropdown-menu'
import { cn } from '@/lib/utils'

export const Route = createFileRoute('/_admin/admin/link/')({
  component: LinkListPage,
})

type ViewMode = 'list' | 'table'

function colorHex(id: number): string {
  return mockColors.find((c) => c.id === id)?.color ?? '#6366f1'
}

/** 头像加载失败时回退为首字色块，尺寸由 className 控制 */
function SiteAvatar({
  name,
  url,
  className,
}: {
  name: string
  url: string
  className?: string
}) {
  const [failed, setFailed] = useState(false)
  if (failed) {
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

/** 行内操作菜单：stopPropagation 避免触发整行点击 */
function RowMenu({ link }: { link: LinkItem }) {
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
          <Link to="/admin/link/$id/edit" params={{ id: String(link.id) }}>
            <Pencil className="mr-2 size-4" />
            编辑
          </Link>
        </DropdownMenuItem>
        <DropdownMenuSeparator />
        <DropdownMenuItem className="cursor-pointer text-destructive focus:text-destructive">
          <Trash2 className="mr-2 size-4" />
          删除
        </DropdownMenuItem>
      </DropdownMenuContent>
    </DropdownMenu>
  )
}

/** 可访问状态徽章 */
function StatusBadge({ link }: { link: LinkItem }) {
  return link.ableConnect ? (
    <Badge className="bg-green-500/15 text-green-700 hover:bg-green-500/15">
      可访问
    </Badge>
  ) : (
    <Badge variant="destructive">不可访问</Badge>
  )
}

/** 各列的响应式可见性 / 对齐（TanStack Table 无头，样式由渲染层控制） */
const columnClass: Record<string, string> = {
  locationName: 'hidden sm:table-cell',
  color: 'hidden md:table-cell',
  updatedAt: 'hidden lg:table-cell',
  actions: 'text-right',
}

/** 表格视图：基于 TanStack Table 的无头表格 */
function LinkTable({
  links,
  onOpenDetail,
}: {
  links: Array<LinkItem>
  onOpenDetail: (id: number) => void
}) {
  const columns = useMemo<ColumnDef<LinkItem>[]>(
    () => [
      {
        id: 'site',
        accessorFn: (row) => row.siteName,
        header: '站点',
        cell: ({ row }) => {
          const link = row.original
          return (
            <div className="flex items-center gap-3">
              <SiteAvatar
                name={link.siteName}
                url={link.siteLogo}
                className="size-8 text-xs"
              />
              <div className="min-w-0">
                <div className="truncate font-medium">{link.siteName}</div>
                <div className="flex items-center gap-1 text-xs text-muted-foreground">
                  <span className="truncate">{link.siteUrl}</span>
                  <ExternalLink className="size-3 shrink-0" />
                </div>
              </div>
            </div>
          )
        },
      },
      {
        accessorKey: 'locationName',
        header: '位置',
        cell: ({ row }) => (
          <Badge variant="secondary">{row.original.locationName}</Badge>
        ),
      },
      {
        id: 'color',
        header: '颜色',
        cell: ({ row }) => {
          const link = row.original
          return (
            <span className="flex items-center gap-1.5 text-muted-foreground">
              <span
                className="size-2.5 rounded-full"
                style={{ backgroundColor: colorHex(link.color) }}
                aria-hidden="true"
              />
              {link.colorName}
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
        accessorKey: 'updatedAt',
        header: '更新时间',
        cell: ({ getValue }) => (
          <span className="tabular-nums text-muted-foreground">
            {getValue<string>()}
          </span>
        ),
      },
      {
        id: 'actions',
        enableSorting: false,
        header: () => <span className="sr-only">操作</span>,
        cell: ({ row }) => <RowMenu link={row.original} />,
      },
    ],
    [],
  )

  const [sorting, setSorting] = useState<SortingState>([])
  const [pagination, setPagination] = useState<PaginationState>({
    pageIndex: 0,
    pageSize: 8,
  })

  const table = useReactTable({
    data: links,
    columns,
    state: { sorting, pagination },
    onSortingChange: setSorting,
    onPaginationChange: setPagination,
    getCoreRowModel: getCoreRowModel(),
    getSortedRowModel: getSortedRowModel(),
    getPaginationRowModel: getPaginationRowModel(),
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
                  {header.isPlaceholder ? null : header.column.getCanSort() ? (
                    <button
                      type="button"
                      onClick={header.column.getToggleSortingHandler()}
                      className="inline-flex cursor-pointer items-center gap-1 transition-colors duration-150 hover:text-foreground"
                    >
                      {flexRender(
                        header.column.columnDef.header,
                        header.getContext(),
                      )}
                      {{
                        asc: <ChevronUp className="size-3.5" />,
                        desc: <ChevronDown className="size-3.5" />,
                      }[header.column.getIsSorted() as string] ?? (
                        <ChevronsUpDown className="size-3.5 opacity-40" />
                      )}
                    </button>
                  ) : (
                    flexRender(
                      header.column.columnDef.header,
                      header.getContext(),
                    )
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

      {/* 分页控制栏 */}
      <div className="flex flex-wrap items-center justify-between gap-3 border-t border-border/70 px-4 py-3">
        <span className="text-sm text-muted-foreground">共 {links.length} 条</span>
        <Pagination
          pageIndex={table.getState().pagination.pageIndex}
          pageCount={Math.max(table.getPageCount(), 1)}
          onPageChange={(i) => table.setPageIndex(i)}
        />
      </div>
    </div>
  )
}

function LinkListPage() {
  const navigate = useNavigate()
  const [keyword, setKeyword] = useState('')
  const [locationFilter, setLocationFilter] = useState<number | null>(null)
  const [viewMode, setViewMode] = useState<ViewMode>('list')

  const approvedLinks = useMemo(
    () => mockLinks.filter((l) => l.status === 'approved'),
    [],
  )
  const pendingCount = useMemo(
    () => mockLinks.filter((l) => l.status === 'pending').length,
    [],
  )

  const filteredLinks = useMemo(() => {
    const kw = keyword.trim().toLowerCase()
    return approvedLinks.filter((link) => {
      const matchKeyword =
        kw === '' ||
        link.siteName.toLowerCase().includes(kw) ||
        link.siteUrl.toLowerCase().includes(kw) ||
        link.siteDescription.toLowerCase().includes(kw)
      const matchLocation =
        locationFilter === null || link.location === locationFilter
      return matchKeyword && matchLocation
    })
  }, [approvedLinks, keyword, locationFilter])

  const openDetail = (id: number) =>
    navigate({ to: '/admin/link/$id', params: { id: String(id) } })

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
              {pendingCount > 0 && (
                <Badge variant="destructive" className="ml-2">
                  {pendingCount}
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
              {mockLinks.length}
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
              {approvedLinks.length}
            </span>
            <span className="text-muted-foreground">已通过</span>
          </span>
          <span className="flex items-center gap-1.5">
            <span
              className="size-2 rounded-full"
              style={{ backgroundColor: '#f59e0b' }}
              aria-hidden="true"
            />
            <span className="font-semibold tabular-nums">{pendingCount}</span>
            <span className="text-muted-foreground">待审核</span>
          </span>
        </div>

        {/* 视图切换 */}
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

      {/* 筛选区：搜索 + 分组药丸 */}
      <div className="flex flex-wrap items-center gap-3">
        <div className="relative w-full max-w-xs">
          <Search className="absolute left-3 top-1/2 size-4 -translate-y-1/2 text-muted-foreground" />
          <Input
            placeholder="搜索站点…"
            className="pl-9"
            value={keyword}
            onChange={(e) => setKeyword(e.target.value)}
          />
        </div>
        <div className="flex flex-wrap gap-1.5">
          <button
            type="button"
            onClick={() => setLocationFilter(null)}
            className={cn(
              'cursor-pointer rounded-full px-3.5 py-1.5 text-sm font-medium transition-colors duration-200',
              locationFilter === null
                ? 'bg-primary text-primary-foreground shadow-sm'
                : 'bg-muted text-muted-foreground hover:text-foreground',
            )}
          >
            全部
          </button>
          {mockLocations.map((location) => (
            <button
              key={location.id}
              type="button"
              onClick={() => setLocationFilter(location.id)}
              className={cn(
                'cursor-pointer rounded-full px-3.5 py-1.5 text-sm font-medium transition-colors duration-200',
                locationFilter === location.id
                  ? 'bg-primary text-primary-foreground shadow-sm'
                  : 'bg-muted text-muted-foreground hover:text-foreground',
              )}
            >
              {location.name}
            </button>
          ))}
        </div>
      </div>

      {/* 内容区 */}
      {filteredLinks.length === 0 ? (
        <Card className="border-dashed">
          <CardContent className="flex flex-col items-center justify-center gap-3 py-14 text-center">
            <div className="flex size-12 items-center justify-center rounded-full bg-muted">
              <SearchX className="size-5 text-muted-foreground" />
            </div>
            <div>
              <p className="font-medium">没有找到匹配的友链</p>
              <p className="mt-1 text-sm text-muted-foreground">
                试试调整搜索关键词或分组筛选
              </p>
            </div>
            <Button
              variant="outline"
              size="sm"
              className="cursor-pointer"
              onClick={() => {
                setKeyword('')
                setLocationFilter(null)
              }}
            >
              清除筛选
            </Button>
          </CardContent>
        </Card>
      ) : viewMode === 'list' ? (
        /* 列表视图：响应式卡片网格（1/2/3 列自适应），整卡可点击进详情 */
        <div className="grid grid-cols-1 gap-3 sm:grid-cols-2 xl:grid-cols-3">
          {filteredLinks.map((link) => (
            <div
              key={link.id}
              role="button"
              tabIndex={0}
              onClick={() => openDetail(link.id)}
              onKeyDown={(e) => {
                if (e.key === 'Enter' || e.key === ' ') openDetail(link.id)
              }}
              className="group relative flex w-full cursor-pointer flex-col overflow-hidden rounded-xl border border-border/70 bg-card text-left transition-all duration-200 hover:-translate-y-0.5 hover:border-primary/40 hover:shadow-md"
            >
              {/* 顶部色彩条：呼应站点颜色 */}
              <span
                className="h-1 w-full"
                style={{ backgroundColor: colorHex(link.color) }}
                aria-hidden="true"
              />
              <div className="flex flex-1 flex-col p-4">
                {/* 头部：头像 + 站点信息 + 菜单 */}
                <div className="flex items-start justify-between gap-2">
                  <div className="flex min-w-0 items-center gap-3">
                    <SiteAvatar
                      name={link.siteName}
                      url={link.siteLogo}
                      className="size-11 text-sm"
                    />
                    <div className="min-w-0">
                      <h3 className="truncate font-semibold">
                        {link.siteName}
                      </h3>
                      <a
                        href={link.siteUrl}
                        target="_blank"
                        rel="noopener noreferrer"
                        onClick={(e) => e.stopPropagation()}
                        className="mt-0.5 inline-flex max-w-full items-center gap-1 text-xs text-muted-foreground transition-colors hover:text-primary"
                      >
                        <span className="truncate">{link.siteUrl}</span>
                        <ExternalLink className="size-3 shrink-0" />
                      </a>
                    </div>
                  </div>
                  <RowMenu link={link} />
                </div>

                {/* 描述 */}
                <p className="mt-3 line-clamp-2 flex-1 text-sm leading-relaxed text-muted-foreground">
                  {link.siteDescription}
                </p>

                {/* 徽章行：位置 + 颜色 + 状态 */}
                <div className="mt-3 flex flex-wrap items-center gap-1.5">
                  <Badge variant="secondary">{link.locationName}</Badge>
                  <Badge variant="outline" className="gap-1.5">
                    <span
                      className="size-2 rounded-full"
                      style={{ backgroundColor: colorHex(link.color) }}
                      aria-hidden="true"
                    />
                    {link.colorName}
                  </Badge>
                  <StatusBadge link={link} />
                </div>
              </div>
            </div>
          ))}
        </div>
      ) : (
        /* 表格视图：TanStack Table，整行可点击进详情 */
        <LinkTable links={filteredLinks} onOpenDetail={openDetail} />
      )}
    </div>
  )
}
