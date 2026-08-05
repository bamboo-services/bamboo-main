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

import { useEffect, useState } from 'react'
import { Link, createFileRoute } from '@tanstack/react-router'
import { motion, useReducedMotion } from 'motion/react'
import { MoreHorizontal, Pencil, Plus, Search, Trash2 } from 'lucide-react'
import type {
  SnowflakeID,
  SponsorChannelAdmin,
  SponsorRecordAdmin,
} from '@/api/types'
import { Button } from '@/components/ui/button'
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
import { Select } from '@/components/ui/select'
import { Skeleton } from '@/components/ui/skeleton'
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from '@/components/ui/table'
import { Textarea } from '@/components/ui/textarea'
import {
  BambooRule,
  EnsoEmpty,
  InkBadge,
  InkSwitch,
  PageHead,
  inkTableHeadRow,
  inkTableRow,
  inkTableWrap,
  inkTd,
  inkTh,
} from '@/components/ink-wash'
import { enter } from '@/lib/motion'
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
        src={icon}
        alt={name}
        className="size-full object-cover"
        loading="lazy"
        onError={() => setFailed(true)}
      />
    </div>
  )
}

function SponsorPage() {
  const reduced = useReducedMotion() ?? false
  // 统计仅取总数：轻量查询（page_size=1），与面板内的分页查询互不影响
  const recordsTotal = useAdminRecords({ page: 1, page_size: 1 })
  const recordCount = recordsTotal.data?.pagination.total ?? 0

  return (
    <div className="mx-auto max-w-6xl space-y-6">
      <PageHead
        kicker="SPONSOR · 赞助"
        title="赞助管理"
        sub="管理所有赞助记录与赞助渠道，记录每一份支持。"
        actions={
          <div className="flex items-center gap-4">
            <span className="font-mono text-xs text-text-secondary tabular-nums">
              共{' '}
              <span className="font-semibold text-leaf-deep">{recordCount}</span>{' '}
              条记录
            </span>
            <Link to="/admin/sponsor/verify">
              <Button variant="outline" className="cursor-pointer">
                赞助审核
                {(recordsTotal.data?.pending_count ?? 0) > 0 && (
                  <InkBadge tone="pending" className="ml-2 px-2">
                    {recordsTotal.data?.pending_count}
                  </InkBadge>
                )}
              </Button>
            </Link>
          </div>
        }
      />

      <BambooRule reduced={reduced} delay={0.12} />

      {/* 记录面板：渠道卡带内联 + 记录表格，一屏同览 */}
      <RecordsPanel reduced={reduced} />
    </div>
  )
}

// ---------------------------------------------------------------------------
// 赞助记录
// ---------------------------------------------------------------------------

/** 赞助记录审核状态徽章 */
function recordStatusBadge(status: number) {
  switch (status) {
    case 0:
      return <InkBadge tone="pending">待审核</InkBadge>
    case 1:
      return <InkBadge tone="leaf">已通过</InkBadge>
    case 2:
      return <InkBadge tone="danger">已拒绝</InkBadge>
    default:
      return <InkBadge tone="neutral">未知</InkBadge>
  }
}

/** 赞助记录面板：搜索 / 渠道筛选 / 审核状态筛选 / 服务端分页表格 / 增删改 */
function RecordsPanel({ reduced }: { reduced: boolean }) {
  const [page, setPage] = useState(1)
  const [keyword, setKeyword] = useState('')
  const [channelFilter, setChannelFilter] = useState<SnowflakeID | null>(null)
  const [statusFilter, setStatusFilter] = useState<number | null>(null)
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
    status: statusFilter ?? undefined,
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
      {/* 工具 + 渠道卡带：一张宣纸卡，记录操作与渠道配置同屏 */}
      <motion.div
        {...enter(reduced, 0.18, {
          initial: { opacity: 0, y: 10 },
          animate: { opacity: 1, y: 0 },
          transition: { duration: 0.4, ease: 'easeOut' },
        })}
        className="rounded-lg border border-border bg-card/60 p-5"
      >
        {/* 工具栏：搜索 + 渠道筛选 + 状态筛选 + 添加赞助 */}
        <div className="flex flex-wrap items-center gap-3">
          <div className="relative w-full max-w-xs">
            <Search className="absolute left-3 top-1/2 size-4 -translate-y-1/2 text-text-secondary" />
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
          <Select
            aria-label="按渠道筛选"
            className="w-40"
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
          </Select>
          <Select
            aria-label="按审核状态筛选"
            className="w-36"
            value={statusFilter?.toString() ?? ''}
            onChange={(e) => {
              const value = e.target.value
              setStatusFilter(value ? Number(value) : null)
              setPage(1)
            }}
          >
            <option value="">全部状态</option>
            <option value="0">待审核</option>
            <option value="1">已通过</option>
            <option value="2">已拒绝</option>
          </Select>
          <div className="flex-1" />
          <Button className="cursor-pointer" onClick={() => setFormTarget('new')}>
            <Plus className="mr-2 size-4" />
            添加赞助
          </Button>
        </div>

        {/* 渠道卡带：渠道配置内联，紧随搜索/筛选之后 */}
        <div className="mt-5 border-t border-border/60 pt-5">
          <ChannelStrip reduced={reduced} />
        </div>
      </motion.div>

      {/* 记录表格 */}
      <motion.div
        {...enter(reduced, 0.2, {
          initial: { opacity: 0 },
          animate: { opacity: 1 },
          transition: { duration: 0.4, ease: 'easeOut' },
        })}
        className={inkTableWrap}
      >
        <Table>
          <TableHeader>
            <TableRow className={cn(inkTableHeadRow, 'hover:bg-muted/30')}>
              <TableHead className={inkTh}>赞助者</TableHead>
              <TableHead className={inkTh}>金额</TableHead>
              <TableHead className={inkTh}>渠道</TableHead>
              <TableHead className={cn(inkTh, 'hidden md:table-cell')}>
                留言
              </TableHead>
              <TableHead className={inkTh}>状态</TableHead>
              <TableHead className={cn(inkTh, 'hidden lg:table-cell')}>
                时间
              </TableHead>
              <TableHead className={cn(inkTh, 'text-right')}>操作</TableHead>
            </TableRow>
          </TableHeader>
          <TableBody>
            {recordsQuery.isLoading ? (
              Array.from({ length: 5 }).map((_, i) => (
                <TableRow key={i}>
                  <TableCell colSpan={7} className={inkTd}>
                    <Skeleton className="h-6 w-full" />
                  </TableCell>
                </TableRow>
              ))
            ) : records.length === 0 ? (
              <TableRow className="hover:bg-transparent">
                <TableCell colSpan={7} className="px-4 py-10">
                  <EnsoEmpty
                    title="没有找到赞助记录"
                    hint="试试调整搜索条件，或添加一条新的赞助记录"
                  />
                </TableCell>
              </TableRow>
            ) : (
              records.map((record) => (
                <TableRow
                  key={record.id.toString()}
                  className={cn('group', inkTableRow)}
                >
                  <TableCell className={inkTd}>
                    <span className="font-serif font-semibold text-text-primary">
                      {record.nickname}
                    </span>
                  </TableCell>
                  <TableCell
                    className={cn(
                      inkTd,
                      'font-mono font-semibold tabular-nums text-leaf-deep',
                    )}
                  >
                    {formatYuan(record.amount)}
                  </TableCell>
                  <TableCell className={inkTd}>
                    <InkBadge tone="neutral">
                      {record.channel?.name ?? '未分类'}
                    </InkBadge>
                  </TableCell>
                  <TableCell
                    className={cn(inkTd, 'hidden max-w-[220px] md:table-cell')}
                  >
                    <span
                      className="block truncate text-text-secondary"
                      title={record.message ?? undefined}
                    >
                      {record.message || '—'}
                    </span>
                  </TableCell>
                  <TableCell className={inkTd}>
                    <div className="flex flex-wrap gap-1">
                      {recordStatusBadge(record.status)}
                      {record.is_hidden ? (
                        <InkBadge tone="danger">隐藏</InkBadge>
                      ) : (
                        <InkBadge tone="neutral">公开</InkBadge>
                      )}
                      {record.is_anonymous && (
                        <InkBadge tone="neutral">匿名</InkBadge>
                      )}
                    </div>
                  </TableCell>
                  <TableCell className={cn(inkTd, 'hidden lg:table-cell')}>
                    <span className="font-mono tabular-nums text-text-secondary">
                      {formatDate(record.sponsor_at ?? record.created_at)}
                    </span>
                  </TableCell>
                  <TableCell className={cn(inkTd, 'text-right')}>
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
        <div className="flex flex-wrap items-center justify-between gap-3 border-t border-border px-4 py-3">
          <span className="font-mono text-xs text-text-secondary">
            共 {total} 条 · 第 {page} / {Math.max(totalPages, 1)} 页
          </span>
          <Pagination
            pageIndex={page - 1}
            pageCount={Math.max(totalPages, 1)}
            onPageChange={(i) => setPage(i + 1)}
          />
        </div>
      </motion.div>

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
              <Select
                id="recordChannel"
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
              </Select>
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

/** 渠道卡带：全量渠道卡片横排，内联在记录工具栏之后；开关 / 编辑 / 删除 */
function ChannelStrip({ reduced }: { reduced: boolean }) {
  // 渠道全量（page_size 上限 100）：取完整字段供编辑回填；赞助渠道数量远低于上限
  const channelsQuery = useAdminChannels({
    page: 1,
    page_size: 100,
    order_by: 'sort_order',
    order: 'asc',
  })
  const deleteChannel = useDeleteChannel()
  const updateStatus = useUpdateChannelStatus()

  const [formTarget, setFormTarget] = useState<
    SponsorChannelAdmin | 'new' | null
  >(null)
  const [deleteTarget, setDeleteTarget] = useState<SponsorChannelAdmin | null>(
    null,
  )

  const channels = channelsQuery.data?.data ?? []

  const handleDelete = () => {
    if (!deleteTarget) return
    deleteChannel.mutate(deleteTarget.id, {
      onSuccess: () => setDeleteTarget(null),
    })
  }

  return (
    <div>
      <div className="flex items-center justify-between gap-3">
        <p className="flex items-center gap-2 font-serif font-semibold text-text-primary">
          <span
            aria-hidden
            className="h-3.5 w-1 -skew-x-12 rounded-sm bg-leaf-deep"
          />
          赞助渠道
          <span className="font-mono text-xs font-normal text-text-secondary tabular-nums">
            {channels.length} 个
          </span>
        </p>
        <Button
          variant="outline"
          size="sm"
          className="cursor-pointer"
          onClick={() => setFormTarget('new')}
        >
          <Plus className="mr-1.5 size-3.5" />
          添加渠道
        </Button>
      </div>

      <motion.div
        {...enter(reduced, 0.22, {
          initial: { opacity: 0, y: 8 },
          animate: { opacity: 1, y: 0 },
          transition: { duration: 0.4, ease: 'easeOut' },
        })}
        className="mt-3 flex flex-wrap gap-3"
      >
        {channelsQuery.isLoading ? (
          Array.from({ length: 3 }).map((_, i) => (
            <Skeleton key={i} className="h-[76px] w-44 rounded-lg" />
          ))
        ) : channels.length === 0 ? (
          <p className="w-full rounded-lg border border-dashed border-border px-4 py-6 text-center text-sm text-text-secondary">
            还没有配置赞助渠道，点击右上角「添加渠道」创建第一个收款渠道。
          </p>
        ) : (
          channels.map((channel) => (
            <ChannelCard
              key={channel.id.toString()}
              channel={channel}
              statusPending={
                updateStatus.isPending &&
                updateStatus.variables.id === channel.id
              }
              onToggle={() =>
                updateStatus.mutate({ id: channel.id, status: !channel.status })
              }
              onEdit={() => setFormTarget(channel)}
              onDelete={() => setDeleteTarget(channel)}
            />
          ))
        )}
      </motion.div>

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

/** 单个渠道卡片：图标 + 名称 + 启用开关 + 赞助次数 + 悬停操作菜单 */
function ChannelCard({
  channel,
  statusPending,
  onToggle,
  onEdit,
  onDelete,
}: {
  channel: SponsorChannelAdmin
  statusPending: boolean
  onToggle: () => void
  onEdit: () => void
  onDelete: () => void
}) {
  return (
    <div className="group relative w-44 rounded-lg border border-border bg-card p-3 transition-shadow duration-200 hover:shadow-md">
      <div className="flex items-center gap-2.5">
        <ChannelIcon
          name={channel.name}
          icon={channel.icon}
          className="size-8 text-[10px]"
        />
        <span className="min-w-0 flex-1 truncate font-serif font-semibold text-text-primary">
          {channel.name}
        </span>
        <DropdownMenu>
          <DropdownMenuTrigger asChild>
            <Button
              variant="ghost"
              size="icon"
              className="size-7 cursor-pointer opacity-0 transition-opacity group-hover:opacity-100 focus-visible:opacity-100"
            >
              <MoreHorizontal className="size-4" />
            </Button>
          </DropdownMenuTrigger>
          <DropdownMenuContent align="end">
            <DropdownMenuItem className="cursor-pointer" onClick={onEdit}>
              <Pencil className="mr-2 size-4" />
              编辑
            </DropdownMenuItem>
            <DropdownMenuSeparator />
            <DropdownMenuItem
              className="cursor-pointer text-destructive focus:text-destructive"
              onClick={onDelete}
            >
              <Trash2 className="mr-2 size-4" />
              删除
            </DropdownMenuItem>
          </DropdownMenuContent>
        </DropdownMenu>
      </div>
      <div className="mt-2.5 flex items-center justify-between gap-2">
        <div className="flex items-center gap-1.5">
          <InkSwitch
            checked={channel.status}
            disabled={statusPending}
            onToggle={onToggle}
          />
          <span className="text-xs text-text-secondary">
            {statusPending ? '更新中' : channel.status ? '启用' : '禁用'}
          </span>
        </div>
        <span className="font-mono text-xs text-text-secondary tabular-nums">
          {channel.sponsor_count} 次
        </span>
      </div>
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
