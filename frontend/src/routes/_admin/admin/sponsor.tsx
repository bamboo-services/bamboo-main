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

import { useEffect, useState } from 'react'
import { createFileRoute } from '@tanstack/react-router'
import {
  Heart,
  MoreHorizontal,
  Pencil,
  Plus,
  Search,
  SearchX,
  Sprout,
  Trash2,
} from 'lucide-react'
import type {
  SnowflakeID,
  SponsorChannelAdmin,
  SponsorRecordAdmin,
} from '@/api/types'
import { Badge } from '@/components/ui/badge'
import { Button } from '@/components/ui/button'
import { Card, CardContent } from '@/components/ui/card'
import { Checkbox } from '@/components/ui/checkbox'
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
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { Pagination } from '@/components/ui/pagination'
import { Skeleton } from '@/components/ui/skeleton'
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from '@/components/ui/table'
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs'
import { Textarea } from '@/components/ui/textarea'
import { cn } from '@/lib/utils'
import {
  useAdminChannels,
  useAdminRecords,
  useAllChannels,
  useCreateChannel,
  useCreateRecord,
  useDeleteChannel,
  useDeleteRecord,
  useUpdateChannel,
  useUpdateChannelStatus,
  useUpdateRecord,
} from '@/hooks/use-sponsors'

export const Route = createFileRoute('/_admin/admin/sponsor')({
  component: SponsorPage,
})

const PAGE_SIZE = 10

/** 原生 select 样式（与 Input 视觉对齐） */
const selectClass =
  'border-input dark:bg-input/30 h-9 w-full cursor-pointer rounded-md border bg-transparent px-3 py-1 text-sm shadow-xs transition-[color,box-shadow] outline-none focus-visible:border-ring focus-visible:ring-ring/50 focus-visible:ring-[3px] disabled:cursor-not-allowed disabled:opacity-50'

/** 金额（分）转元展示 */
function formatYuan(cents: number) {
  return `¥${(cents / 100).toLocaleString('zh-CN', {
    minimumFractionDigits: 2,
    maximumFractionDigits: 2,
  })}`
}

/** ISO 时间格式化为本地日期 */
function formatDate(iso: string) {
  return new Date(iso).toLocaleDateString('zh-CN')
}

/** 渠道图标：加载失败或为空时回退为首字色块 */
function ChannelIcon({
  name,
  icon,
  className,
}: {
  name: string
  icon: string | null
  className?: string
}) {
  const [failed, setFailed] = useState(false)
  if (failed || !icon) {
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
        src={icon}
        alt={name}
        className="size-full object-cover"
        loading="lazy"
        onError={() => setFailed(true)}
      />
    </div>
  )
}

/** 启用/禁用开关（button + role=switch 实现） */
function StatusSwitch({
  checked,
  disabled,
  onToggle,
}: {
  checked: boolean
  disabled?: boolean
  onToggle: () => void
}) {
  return (
    <button
      type="button"
      role="switch"
      aria-checked={checked}
      disabled={disabled}
      onClick={onToggle}
      className={cn(
        'relative inline-flex h-5 w-9 shrink-0 cursor-pointer items-center rounded-full transition-colors duration-200 disabled:cursor-not-allowed disabled:opacity-50',
        checked ? 'bg-primary' : 'bg-muted-foreground/25',
      )}
    >
      <span
        className={cn(
          'inline-block size-4 rounded-full bg-background shadow-sm transition-transform duration-200',
          checked ? 'translate-x-[18px]' : 'translate-x-0.5',
        )}
      />
    </button>
  )
}

function SponsorPage() {
  // 统计仅取总数：轻量查询（page_size=1），与面板内的分页查询互不影响
  const recordsTotal = useAdminRecords({ page: 1, page_size: 1 })
  const allChannels = useAllChannels()

  const recordCount = recordsTotal.data?.pagination.total ?? 0
  const channelCount = allChannels.data?.length ?? 0

  return (
    <div className="space-y-6">
      {/* 页头 */}
      <div>
        <h1 className="text-3xl font-bold tracking-tight">赞助管理</h1>
        <p className="mt-1 text-muted-foreground">管理所有赞助记录与赞助渠道</p>
      </div>

      {/* 统计卡片 */}
      <div className="grid gap-4 sm:grid-cols-2">
        <Card className="transition-colors duration-200 hover:border-rose-300/60">
          <CardContent className="flex items-center gap-4 p-5">
            <div className="flex size-11 shrink-0 items-center justify-center rounded-xl bg-rose-500/10">
              <Heart className="size-5 text-rose-500" />
            </div>
            <div>
              {recordsTotal.isLoading ? (
                <Skeleton className="mb-1 h-7 w-16" />
              ) : (
                <div className="text-2xl font-bold tabular-nums">
                  {recordCount}
                </div>
              )}
              <p className="text-sm text-muted-foreground">赞助记录数</p>
            </div>
          </CardContent>
        </Card>
        <Card className="transition-colors duration-200 hover:border-emerald-300/60">
          <CardContent className="flex items-center gap-4 p-5">
            <div className="flex size-11 shrink-0 items-center justify-center rounded-xl bg-emerald-500/10">
              <Sprout className="size-5 text-emerald-500" />
            </div>
            <div>
              {allChannels.isLoading ? (
                <Skeleton className="mb-1 h-7 w-16" />
              ) : (
                <div className="text-2xl font-bold tabular-nums">
                  {channelCount}
                </div>
              )}
              <p className="text-sm text-muted-foreground">赞助渠道数</p>
            </div>
          </CardContent>
        </Card>
      </div>

      {/* 双标签：赞助记录 + 赞助渠道 */}
      <Tabs defaultValue="records">
        <TabsList>
          <TabsTrigger value="records">赞助记录</TabsTrigger>
          <TabsTrigger value="channels">赞助渠道</TabsTrigger>
        </TabsList>
        <TabsContent value="records" className="mt-4">
          <RecordsPanel />
        </TabsContent>
        <TabsContent value="channels" className="mt-4">
          <ChannelsPanel />
        </TabsContent>
      </Tabs>
    </div>
  )
}

// ---------------------------------------------------------------------------
// 赞助记录
// ---------------------------------------------------------------------------

/** 赞助记录面板：搜索 / 渠道筛选 / 服务端分页表格 / 增删改 */
function RecordsPanel() {
  const [page, setPage] = useState(1)
  const [keyword, setKeyword] = useState('')
  const [channelFilter, setChannelFilter] = useState<SnowflakeID | null>(null)
  const [formTarget, setFormTarget] = useState<
    SponsorRecordAdmin | 'new' | null
  >(null)
  const [deleteTarget, setDeleteTarget] = useState<SponsorRecordAdmin | null>(
    null,
  )

  const allChannels = useAllChannels()
  const recordsQuery = useAdminRecords({
    page,
    page_size: PAGE_SIZE,
    nickname: keyword.trim() || undefined,
    channel_id: channelFilter ?? undefined,
    order_by: 'created_at',
    order: 'desc',
  })
  const deleteRecord = useDeleteRecord()

  const records = recordsQuery.data?.data ?? []
  const total = recordsQuery.data?.pagination.total ?? 0
  const totalPages = recordsQuery.data?.pagination.total_pages ?? 1
  const channels = allChannels.data ?? []

  const handleDelete = () => {
    if (!deleteTarget) return
    deleteRecord.mutate(deleteTarget.id, {
      onSuccess: () => setDeleteTarget(null),
    })
  }

  return (
    <div className="space-y-4">
      {/* 工具栏：搜索 + 渠道筛选 + 添加 */}
      <div className="flex flex-wrap items-center gap-3">
        <div className="relative w-full max-w-xs">
          <Search className="absolute left-3 top-1/2 size-4 -translate-y-1/2 text-muted-foreground" />
          <Input
            placeholder="搜索赞助者昵称…"
            className="pl-9"
            value={keyword}
            onChange={(e) => {
              setKeyword(e.target.value)
              setPage(1)
            }}
          />
        </div>
        <select
          aria-label="按渠道筛选"
          className={cn(selectClass, 'w-40')}
          value={channelFilter?.toString() ?? ''}
          onChange={(e) => {
            const value = e.target.value
            setChannelFilter(value ? BigInt(value) : null)
            setPage(1)
          }}
        >
          <option value="">全部渠道</option>
          {channels.map((channel) => (
            <option key={channel.id.toString()} value={channel.id.toString()}>
              {channel.name}
            </option>
          ))}
        </select>
        <div className="flex-1" />
        <Button className="cursor-pointer" onClick={() => setFormTarget('new')}>
          <Plus className="mr-2 size-4" />
          添加赞助
        </Button>
      </div>

      {/* 记录表格 */}
      <div className="overflow-hidden rounded-xl border border-border/70 bg-card">
        <Table>
          <TableHeader>
            <TableRow className="bg-muted/40 hover:bg-muted/40">
              <TableHead className="px-4">赞助者</TableHead>
              <TableHead className="px-4">金额</TableHead>
              <TableHead className="px-4">渠道</TableHead>
              <TableHead className="hidden px-4 md:table-cell">留言</TableHead>
              <TableHead className="px-4">状态</TableHead>
              <TableHead className="hidden px-4 lg:table-cell">时间</TableHead>
              <TableHead className="px-4 text-right">操作</TableHead>
            </TableRow>
          </TableHeader>
          <TableBody>
            {recordsQuery.isLoading ? (
              Array.from({ length: 5 }).map((_, i) => (
                <TableRow key={i}>
                  <TableCell colSpan={7} className="px-4 py-3">
                    <Skeleton className="h-6 w-full" />
                  </TableCell>
                </TableRow>
              ))
            ) : records.length === 0 ? (
              <TableRow className="hover:bg-transparent">
                <TableCell colSpan={7} className="px-4 py-14">
                  <div className="flex flex-col items-center gap-3 text-center">
                    <div className="flex size-12 items-center justify-center rounded-full bg-muted">
                      <SearchX className="size-5 text-muted-foreground" />
                    </div>
                    <div>
                      <p className="font-medium">没有找到赞助记录</p>
                      <p className="mt-1 text-sm text-muted-foreground">
                        试试调整搜索条件，或添加一条新的赞助记录
                      </p>
                    </div>
                  </div>
                </TableCell>
              </TableRow>
            ) : (
              records.map((record) => (
                <TableRow key={record.id.toString()} className="group">
                  <TableCell className="px-4 py-3">
                    <span className="font-medium">{record.nickname}</span>
                  </TableCell>
                  <TableCell className="px-4 py-3 font-medium tabular-nums text-green-600">
                    {formatYuan(record.amount)}
                  </TableCell>
                  <TableCell className="px-4 py-3">
                    <Badge variant="secondary">
                      {record.channel?.name ?? '未分类'}
                    </Badge>
                  </TableCell>
                  <TableCell className="hidden max-w-[220px] px-4 py-3 md:table-cell">
                    <span
                      className="block truncate text-muted-foreground"
                      title={record.message ?? undefined}
                    >
                      {record.message || '—'}
                    </span>
                  </TableCell>
                  <TableCell className="px-4 py-3">
                    <div className="flex flex-wrap gap-1">
                      {record.is_hidden ? (
                        <Badge variant="destructive">隐藏</Badge>
                      ) : (
                        <Badge className="bg-green-500/15 text-green-700 hover:bg-green-500/15">
                          公开
                        </Badge>
                      )}
                      {record.is_anonymous && (
                        <Badge variant="outline">匿名</Badge>
                      )}
                    </div>
                  </TableCell>
                  <TableCell className="hidden px-4 py-3 lg:table-cell">
                    <span className="tabular-nums text-muted-foreground">
                      {formatDate(record.sponsor_at ?? record.created_at)}
                    </span>
                  </TableCell>
                  <TableCell className="px-4 py-3 text-right">
                    <DropdownMenu>
                      <DropdownMenuTrigger asChild>
                        <Button
                          variant="ghost"
                          size="icon"
                          className="size-8 cursor-pointer opacity-50 transition-opacity hover:opacity-100 group-hover:opacity-100"
                        >
                          <MoreHorizontal className="size-4" />
                        </Button>
                      </DropdownMenuTrigger>
                      <DropdownMenuContent align="end">
                        <DropdownMenuItem
                          className="cursor-pointer"
                          onClick={() => setFormTarget(record)}
                        >
                          <Pencil className="mr-2 size-4" />
                          编辑
                        </DropdownMenuItem>
                        <DropdownMenuSeparator />
                        <DropdownMenuItem
                          className="cursor-pointer text-destructive focus:text-destructive"
                          onClick={() => setDeleteTarget(record)}
                        >
                          <Trash2 className="mr-2 size-4" />
                          删除
                        </DropdownMenuItem>
                      </DropdownMenuContent>
                    </DropdownMenu>
                  </TableCell>
                </TableRow>
              ))
            )}
          </TableBody>
        </Table>

        {/* 分页 */}
        <div className="flex flex-wrap items-center justify-between gap-3 border-t border-border/70 px-4 py-3">
          <span className="text-sm text-muted-foreground">
            共 {total} 条 · 第 {page} / {Math.max(totalPages, 1)} 页
          </span>
          <Pagination
            pageIndex={page - 1}
            pageCount={Math.max(totalPages, 1)}
            onPageChange={(i) => setPage(i + 1)}
          />
        </div>
      </div>

      {/* 新增 / 编辑弹窗 */}
      <RecordFormDialog
        target={formTarget}
        onClose={() => setFormTarget(null)}
      />

      {/* 删除确认弹窗 */}
      <Dialog
        open={deleteTarget !== null}
        onOpenChange={(open) => {
          if (!open) setDeleteTarget(null)
        }}
      >
        <DialogContent>
          <DialogHeader>
            <DialogTitle>删除赞助记录</DialogTitle>
            <DialogDescription>
              确定要删除「{deleteTarget?.nickname}
              」的赞助记录吗？此操作不可撤销。
            </DialogDescription>
          </DialogHeader>
          <DialogFooter>
            <Button
              variant="outline"
              onClick={() => setDeleteTarget(null)}
              disabled={deleteRecord.isPending}
            >
              取消
            </Button>
            <Button
              variant="destructive"
              onClick={handleDelete}
              disabled={deleteRecord.isPending}
            >
              {deleteRecord.isPending ? '删除中…' : '确认删除'}
            </Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>
    </div>
  )
}

/** 赞助记录新增 / 编辑复用弹窗（target 为 'new' 时新增，为记录对象时编辑） */
function RecordFormDialog({
  target,
  onClose,
}: {
  target: SponsorRecordAdmin | 'new' | null
  onClose: () => void
}) {
  const allChannels = useAllChannels()
  const createRecord = useCreateRecord()
  const updateRecord = useUpdateRecord()

  const [nickname, setNickname] = useState('')
  const [amountYuan, setAmountYuan] = useState('')
  const [channelId, setChannelId] = useState('')
  const [message, setMessage] = useState('')
  const [redirectUrl, setRedirectUrl] = useState('')
  const [sponsorAt, setSponsorAt] = useState('')
  const [sortOrder, setSortOrder] = useState('0')
  const [isAnonymous, setIsAnonymous] = useState(false)
  const [isHidden, setIsHidden] = useState(false)

  // 弹窗目标变化时初始化表单
  useEffect(() => {
    if (target === null) return
    if (target === 'new') {
      setNickname('')
      setAmountYuan('')
      setChannelId('')
      setMessage('')
      setRedirectUrl('')
      setSponsorAt('')
      setSortOrder('0')
      setIsAnonymous(false)
      setIsHidden(false)
    } else {
      setNickname(target.nickname)
      setAmountYuan((target.amount / 100).toFixed(2))
      setChannelId(target.channel_id?.toString() ?? '')
      setMessage(target.message ?? '')
      setRedirectUrl(target.redirect_url ?? '')
      setSponsorAt(target.sponsor_at ? target.sponsor_at.slice(0, 10) : '')
      setSortOrder(String(target.sort_order))
      setIsAnonymous(target.is_anonymous)
      setIsHidden(target.is_hidden)
    }
  }, [target])

  const channels = allChannels.data ?? []
  const isEdit = target !== null && target !== 'new'
  const pending = createRecord.isPending || updateRecord.isPending

  const amountInvalid =
    amountYuan.trim() === '' ||
    Number.isNaN(Number(amountYuan)) ||
    Number(amountYuan) < 0
  const submitDisabled = pending || nickname.trim() === '' || amountInvalid

  const handleSubmit = () => {
    if (!target || submitDisabled) return
    const req = {
      nickname: nickname.trim(),
      // 元 → 分
      amount: Math.round(Number(amountYuan) * 100),
      channel_id: channelId ? BigInt(channelId) : undefined,
      message: message.trim() || undefined,
      redirect_url: redirectUrl.trim() || undefined,
      sponsor_at: sponsorAt ? `${sponsorAt}T00:00:00Z` : undefined,
      sort_order: Number(sortOrder) || 0,
      is_anonymous: isAnonymous,
      is_hidden: isHidden,
    }
    if (target === 'new') {
      createRecord.mutate(req, { onSuccess: () => onClose() })
    } else {
      updateRecord.mutate(
        { id: target.id, req },
        { onSuccess: () => onClose() },
      )
    }
  }

  return (
    <Dialog
      open={target !== null}
      onOpenChange={(open) => {
        if (!open) onClose()
      }}
    >
      <DialogContent className="sm:max-w-xl">
        <DialogHeader>
          <DialogTitle>{isEdit ? '编辑赞助记录' : '添加赞助记录'}</DialogTitle>
          <DialogDescription>
            {isEdit ? '修改该条赞助记录的信息' : '手动录入一条赞助记录'}
          </DialogDescription>
        </DialogHeader>

        <div className="grid gap-4">
          <div className="grid gap-4 sm:grid-cols-2">
            <div className="space-y-2">
              <Label htmlFor="recordNickname">昵称 *</Label>
              <Input
                id="recordNickname"
                value={nickname}
                onChange={(e) => setNickname(e.target.value)}
                placeholder="赞助者昵称"
              />
            </div>
            <div className="space-y-2">
              <Label htmlFor="recordAmount">金额（元）*</Label>
              <Input
                id="recordAmount"
                type="number"
                min="0"
                step="0.01"
                value={amountYuan}
                onChange={(e) => setAmountYuan(e.target.value)}
                placeholder="0.00"
              />
            </div>
          </div>

          <div className="grid gap-4 sm:grid-cols-2">
            <div className="space-y-2">
              <Label htmlFor="recordChannel">赞助渠道</Label>
              <select
                id="recordChannel"
                className={selectClass}
                value={channelId}
                onChange={(e) => setChannelId(e.target.value)}
              >
                <option value="">未分类</option>
                {channels.map((channel) => (
                  <option
                    key={channel.id.toString()}
                    value={channel.id.toString()}
                  >
                    {channel.name}
                  </option>
                ))}
              </select>
            </div>
            <div className="space-y-2">
              <Label htmlFor="recordSponsorAt">赞助时间</Label>
              <Input
                id="recordSponsorAt"
                type="date"
                value={sponsorAt}
                onChange={(e) => setSponsorAt(e.target.value)}
              />
            </div>
          </div>

          <div className="space-y-2">
            <Label htmlFor="recordMessage">留言</Label>
            <Textarea
              id="recordMessage"
              className="min-h-[72px]"
              value={message}
              onChange={(e) => setMessage(e.target.value)}
              placeholder="赞助者留下的话（可选）"
            />
          </div>

          <div className="space-y-2">
            <Label htmlFor="recordRedirectUrl">跳转链接</Label>
            <Input
              id="recordRedirectUrl"
              value={redirectUrl}
              onChange={(e) => setRedirectUrl(e.target.value)}
              placeholder="https://（可选）"
            />
          </div>

          <div className="grid gap-4 sm:grid-cols-2">
            <div className="space-y-2">
              <Label htmlFor="recordSortOrder">排序</Label>
              <Input
                id="recordSortOrder"
                type="number"
                value={sortOrder}
                onChange={(e) => setSortOrder(e.target.value)}
              />
            </div>
            <div className="flex items-end gap-5 pb-1.5">
              <div className="flex items-center gap-2">
                <Checkbox
                  id="recordAnonymous"
                  checked={isAnonymous}
                  onCheckedChange={(v) => setIsAnonymous(v === true)}
                />
                <Label
                  htmlFor="recordAnonymous"
                  className="cursor-pointer font-normal"
                >
                  匿名
                </Label>
              </div>
              <div className="flex items-center gap-2">
                <Checkbox
                  id="recordHidden"
                  checked={isHidden}
                  onCheckedChange={(v) => setIsHidden(v === true)}
                />
                <Label
                  htmlFor="recordHidden"
                  className="cursor-pointer font-normal"
                >
                  隐藏
                </Label>
              </div>
            </div>
          </div>
        </div>

        <DialogFooter>
          <Button variant="outline" onClick={onClose} disabled={pending}>
            取消
          </Button>
          <Button
            className="cursor-pointer"
            onClick={handleSubmit}
            disabled={submitDisabled}
          >
            {pending ? '保存中…' : isEdit ? '保存修改' : '确认添加'}
          </Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  )
}

// ---------------------------------------------------------------------------
// 赞助渠道
// ---------------------------------------------------------------------------

/** 赞助渠道面板：搜索 / 服务端分页表格 / 增删改 / 状态切换 */
function ChannelsPanel() {
  const [page, setPage] = useState(1)
  const [keyword, setKeyword] = useState('')
  const [formTarget, setFormTarget] = useState<
    SponsorChannelAdmin | 'new' | null
  >(null)
  const [deleteTarget, setDeleteTarget] = useState<SponsorChannelAdmin | null>(
    null,
  )

  const channelsQuery = useAdminChannels({
    page,
    page_size: PAGE_SIZE,
    name: keyword.trim() || undefined,
    order_by: 'sort_order',
    order: 'asc',
  })
  const deleteChannel = useDeleteChannel()
  const updateStatus = useUpdateChannelStatus()

  const channels = channelsQuery.data?.data ?? []
  const total = channelsQuery.data?.pagination.total ?? 0
  const totalPages = channelsQuery.data?.pagination.total_pages ?? 1

  const handleDelete = () => {
    if (!deleteTarget) return
    deleteChannel.mutate(deleteTarget.id, {
      onSuccess: () => setDeleteTarget(null),
    })
  }

  return (
    <div className="space-y-4">
      {/* 工具栏：搜索 + 添加 */}
      <div className="flex flex-wrap items-center gap-3">
        <div className="relative w-full max-w-xs">
          <Search className="absolute left-3 top-1/2 size-4 -translate-y-1/2 text-muted-foreground" />
          <Input
            placeholder="搜索渠道名称…"
            className="pl-9"
            value={keyword}
            onChange={(e) => {
              setKeyword(e.target.value)
              setPage(1)
            }}
          />
        </div>
        <div className="flex-1" />
        <Button className="cursor-pointer" onClick={() => setFormTarget('new')}>
          <Plus className="mr-2 size-4" />
          添加渠道
        </Button>
      </div>

      {/* 渠道表格 */}
      <div className="overflow-hidden rounded-xl border border-border/70 bg-card">
        <Table>
          <TableHeader>
            <TableRow className="bg-muted/40 hover:bg-muted/40">
              <TableHead className="px-4">渠道</TableHead>
              <TableHead className="hidden px-4 md:table-cell">描述</TableHead>
              <TableHead className="px-4">排序</TableHead>
              <TableHead className="px-4">赞助次数</TableHead>
              <TableHead className="px-4">状态</TableHead>
              <TableHead className="px-4 text-right">操作</TableHead>
            </TableRow>
          </TableHeader>
          <TableBody>
            {channelsQuery.isLoading ? (
              Array.from({ length: 5 }).map((_, i) => (
                <TableRow key={i}>
                  <TableCell colSpan={6} className="px-4 py-3">
                    <Skeleton className="h-6 w-full" />
                  </TableCell>
                </TableRow>
              ))
            ) : channels.length === 0 ? (
              <TableRow className="hover:bg-transparent">
                <TableCell colSpan={6} className="px-4 py-14">
                  <div className="flex flex-col items-center gap-3 text-center">
                    <div className="flex size-12 items-center justify-center rounded-full bg-muted">
                      <SearchX className="size-5 text-muted-foreground" />
                    </div>
                    <div>
                      <p className="font-medium">没有找到赞助渠道</p>
                      <p className="mt-1 text-sm text-muted-foreground">
                        试试调整搜索关键词，或添加一个新的赞助渠道
                      </p>
                    </div>
                  </div>
                </TableCell>
              </TableRow>
            ) : (
              channels.map((channel) => {
                const statusPending =
                  updateStatus.isPending &&
                  updateStatus.variables.id === channel.id
                return (
                  <TableRow key={channel.id.toString()} className="group">
                    <TableCell className="px-4 py-3">
                      <div className="flex items-center gap-3">
                        <ChannelIcon
                          name={channel.name}
                          icon={channel.icon}
                          className="size-8 text-xs"
                        />
                        <span className="font-medium">{channel.name}</span>
                      </div>
                    </TableCell>
                    <TableCell className="hidden max-w-[260px] px-4 py-3 md:table-cell">
                      <span
                        className="block truncate text-muted-foreground"
                        title={channel.description ?? undefined}
                      >
                        {channel.description || '—'}
                      </span>
                    </TableCell>
                    <TableCell className="px-4 py-3 tabular-nums text-muted-foreground">
                      {channel.sort_order}
                    </TableCell>
                    <TableCell className="px-4 py-3 font-medium tabular-nums">
                      {channel.sponsor_count}
                    </TableCell>
                    <TableCell className="px-4 py-3">
                      <div className="flex items-center gap-2">
                        <StatusSwitch
                          checked={channel.status}
                          disabled={statusPending}
                          onToggle={() =>
                            updateStatus.mutate({
                              id: channel.id,
                              status: !channel.status,
                            })
                          }
                        />
                        <span
                          className={cn(
                            'text-sm',
                            channel.status
                              ? 'text-foreground'
                              : 'text-muted-foreground',
                          )}
                        >
                          {statusPending
                            ? '更新中…'
                            : channel.status
                              ? '启用'
                              : '禁用'}
                        </span>
                      </div>
                    </TableCell>
                    <TableCell className="px-4 py-3 text-right">
                      <DropdownMenu>
                        <DropdownMenuTrigger asChild>
                          <Button
                            variant="ghost"
                            size="icon"
                            className="size-8 cursor-pointer opacity-50 transition-opacity hover:opacity-100 group-hover:opacity-100"
                          >
                            <MoreHorizontal className="size-4" />
                          </Button>
                        </DropdownMenuTrigger>
                        <DropdownMenuContent align="end">
                          <DropdownMenuItem
                            className="cursor-pointer"
                            onClick={() => setFormTarget(channel)}
                          >
                            <Pencil className="mr-2 size-4" />
                            编辑
                          </DropdownMenuItem>
                          <DropdownMenuSeparator />
                          <DropdownMenuItem
                            className="cursor-pointer text-destructive focus:text-destructive"
                            onClick={() => setDeleteTarget(channel)}
                          >
                            <Trash2 className="mr-2 size-4" />
                            删除
                          </DropdownMenuItem>
                        </DropdownMenuContent>
                      </DropdownMenu>
                    </TableCell>
                  </TableRow>
                )
              })
            )}
          </TableBody>
        </Table>

        {/* 分页 */}
        <div className="flex flex-wrap items-center justify-between gap-3 border-t border-border/70 px-4 py-3">
          <span className="text-sm text-muted-foreground">
            共 {total} 条 · 第 {page} / {Math.max(totalPages, 1)} 页
          </span>
          <Pagination
            pageIndex={page - 1}
            pageCount={Math.max(totalPages, 1)}
            onPageChange={(i) => setPage(i + 1)}
          />
        </div>
      </div>

      {/* 新增 / 编辑弹窗 */}
      <ChannelFormDialog
        target={formTarget}
        onClose={() => setFormTarget(null)}
      />

      {/* 删除确认弹窗 */}
      <Dialog
        open={deleteTarget !== null}
        onOpenChange={(open) => {
          if (!open) setDeleteTarget(null)
        }}
      >
        <DialogContent>
          <DialogHeader>
            <DialogTitle>删除赞助渠道</DialogTitle>
            <DialogDescription>
              确定要删除赞助渠道「{deleteTarget?.name}」吗？此操作不可撤销。
            </DialogDescription>
          </DialogHeader>
          <DialogFooter>
            <Button
              variant="outline"
              onClick={() => setDeleteTarget(null)}
              disabled={deleteChannel.isPending}
            >
              取消
            </Button>
            <Button
              variant="destructive"
              onClick={handleDelete}
              disabled={deleteChannel.isPending}
            >
              {deleteChannel.isPending ? '删除中…' : '确认删除'}
            </Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>
    </div>
  )
}

/** 赞助渠道新增 / 编辑复用弹窗（target 为 'new' 时新增，为渠道对象时编辑） */
function ChannelFormDialog({
  target,
  onClose,
}: {
  target: SponsorChannelAdmin | 'new' | null
  onClose: () => void
}) {
  const createChannel = useCreateChannel()
  const updateChannel = useUpdateChannel()

  const [name, setName] = useState('')
  const [icon, setIcon] = useState('')
  const [description, setDescription] = useState('')
  const [sortOrder, setSortOrder] = useState('0')

  // 弹窗目标变化时初始化表单
  useEffect(() => {
    if (target === null) return
    if (target === 'new') {
      setName('')
      setIcon('')
      setDescription('')
      setSortOrder('0')
    } else {
      setName(target.name)
      setIcon(target.icon ?? '')
      setDescription(target.description ?? '')
      setSortOrder(String(target.sort_order))
    }
  }, [target])

  const isEdit = target !== null && target !== 'new'
  const pending = createChannel.isPending || updateChannel.isPending
  const submitDisabled = pending || name.trim() === ''

  const handleSubmit = () => {
    if (!target || submitDisabled) return
    const req = {
      name: name.trim(),
      icon: icon.trim() || undefined,
      description: description.trim() || undefined,
      sort_order: Number(sortOrder) || 0,
    }
    if (target === 'new') {
      createChannel.mutate(req, { onSuccess: () => onClose() })
    } else {
      updateChannel.mutate(
        { id: target.id, req },
        { onSuccess: () => onClose() },
      )
    }
  }

  return (
    <Dialog
      open={target !== null}
      onOpenChange={(open) => {
        if (!open) onClose()
      }}
    >
      <DialogContent>
        <DialogHeader>
          <DialogTitle>{isEdit ? '编辑赞助渠道' : '添加赞助渠道'}</DialogTitle>
          <DialogDescription>
            {isEdit
              ? '修改该赞助渠道的信息'
              : '创建一个赞助渠道（如微信、支付宝）'}
          </DialogDescription>
        </DialogHeader>

        <div className="grid gap-4">
          <div className="space-y-2">
            <Label htmlFor="channelName">名称 *</Label>
            <Input
              id="channelName"
              value={name}
              onChange={(e) => setName(e.target.value)}
              placeholder="渠道名称，如：微信支付"
            />
          </div>

          <div className="space-y-2">
            <Label htmlFor="channelIcon">图标 URL</Label>
            <Input
              id="channelIcon"
              value={icon}
              onChange={(e) => setIcon(e.target.value)}
              placeholder="https://（可选）"
            />
          </div>

          <div className="space-y-2">
            <Label htmlFor="channelDescription">描述</Label>
            <Textarea
              id="channelDescription"
              className="min-h-[72px]"
              value={description}
              onChange={(e) => setDescription(e.target.value)}
              placeholder="渠道的简要说明（可选）"
            />
          </div>

          <div className="space-y-2">
            <Label htmlFor="channelSortOrder">排序</Label>
            <Input
              id="channelSortOrder"
              type="number"
              value={sortOrder}
              onChange={(e) => setSortOrder(e.target.value)}
            />
          </div>
        </div>

        <DialogFooter>
          <Button variant="outline" onClick={onClose} disabled={pending}>
            取消
          </Button>
          <Button
            className="cursor-pointer"
            onClick={handleSubmit}
            disabled={submitDisabled}
          >
            {pending ? '保存中…' : isEdit ? '保存修改' : '确认添加'}
          </Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  )
}
