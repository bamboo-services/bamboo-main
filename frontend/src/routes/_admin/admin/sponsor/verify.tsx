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

import { useCallback, useEffect, useState } from 'react'
import { Link, createFileRoute } from '@tanstack/react-router'
import { motion, useReducedMotion } from 'motion/react'
import { Check, ChevronRight, Coins, Mail, X } from 'lucide-react'
import type {
  SnowflakeID,
  SponsorRecordAdmin,
  UpdateRecordRequest,
} from '@/api/types'
import { Button } from '@/components/ui/button'
import { Checkbox } from '@/components/ui/checkbox'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { Select } from '@/components/ui/select'
import { Skeleton } from '@/components/ui/skeleton'
import { Textarea } from '@/components/ui/textarea'
import {
  BambooRule,
  CardHead,
  EnsoEmpty,
  InkBadge,
  PageHead,
  inkCard,
  inkTableWrap,
} from '@/components/ink-wash'
import { enter } from '@/lib/motion'
import { cn } from '@/lib/utils'
import {
  useAdminRecords,
  useAllChannels,
  useUpdateRecord,
  useUpdateSponsorStatus,
} from '@/hooks/use-sponsors'

export const Route = createFileRoute('/_admin/admin/sponsor/verify')({
  component: SponsorVerifyPage,
})

/** 布局形变统一弹簧：列表收窄与详情弹出共用，节奏一致 */
const morphSpring = { type: 'spring', stiffness: 320, damping: 32 } as const

/** 金额（分）→ 元字符串 */
function formatAmount(amount: number): string {
  return (amount / 100).toFixed(2)
}

/** 赞助记录状态徽章映射 */
function sponsorStatus(record: SponsorRecordAdmin): {
  label: string
  tone: 'leaf' | 'pending' | 'danger'
} {
  switch (record.status) {
    case 0:
      return { label: '待审核', tone: 'pending' }
    case 1:
      return { label: '已通过', tone: 'leaf' }
    case 2:
      return { label: '已拒绝', tone: 'danger' }
    default:
      return { label: '未知', tone: 'danger' }
  }
}

/**
 * 列表条目 —— 富/紧凑两种密度共用同一组件。
 * 选中其他项时通过 layout 动画原地收紧，保证「未选中 → 选中」是同一界面的自然演化。
 */
function VerifyItem({
  record,
  compact,
  active,
  onSelect,
}: {
  record: SponsorRecordAdmin
  compact: boolean
  active: boolean
  onSelect: (id: SnowflakeID) => void
}) {
  const status = sponsorStatus(record)
  return (
    <motion.button
      layout
      type="button"
      onClick={() => onSelect(record.id)}
      transition={morphSpring}
      className={cn(
        'group w-full cursor-pointer overflow-hidden rounded-lg border bg-card text-left transition-colors duration-200',
        active
          ? 'border-leaf-deep/60 bg-leaf-deep/6 shadow-[0_10px_24px_-18px_oklch(0.55_0.12_155/0.5)]'
          : 'border-border hover:border-leaf-muted hover:shadow-[0_10px_24px_-20px_oklch(0.32_0.06_155/0.35)]',
        compact ? 'p-2.5' : 'p-4',
      )}
    >
      {/* 常驻头部：首字头像 + 昵称 + 元信息 */}
      <div className="flex items-center gap-3">
        <div
          className={cn(
            'flex shrink-0 items-center justify-center rounded-lg bg-muted font-serif font-semibold text-text-secondary transition-all duration-300',
            compact ? 'size-8 text-xs' : 'size-11 text-sm',
          )}
        >
          {record.nickname.charAt(0)}
        </div>
        <div className="min-w-0 flex-1">
          <div className="flex items-center justify-between gap-2">
            <h3
              className={cn(
                'truncate font-serif font-semibold text-text-primary transition-all duration-300',
                compact ? 'text-sm' : 'text-base',
              )}
            >
              {record.nickname}
            </h3>
            {compact ? (
              active && (
                <span
                  className="size-2 shrink-0 rounded-full bg-leaf-deep"
                  aria-hidden="true"
                />
              )
            ) : (
              <InkBadge tone={status.tone}>{status.label}</InkBadge>
            )}
          </div>
          <div className="mt-0.5 flex items-center justify-between gap-2 font-mono text-xs text-text-secondary">
            <span className="truncate">
              {record.channel?.name ?? '未选择渠道'}
            </span>
            <span className="shrink-0 tabular-nums">
              ¥{formatAmount(record.amount)}
            </span>
          </div>
        </div>
      </div>

      {/* 富内容：紧凑模式下平滑收起（grid-rows 0fr/1fr 高度动画） */}
      <div
        className={cn(
          'grid transition-all duration-300 ease-out',
          compact ? 'grid-rows-[0fr] opacity-0' : 'grid-rows-[1fr] opacity-100',
        )}
      >
        <div className="overflow-hidden">
          <p
            className="mt-3 truncate text-sm text-text-secondary"
            title={record.message || undefined}
          >
            {record.message || '暂无留言'}
          </p>
          <div className="mt-3 flex items-center justify-between gap-2 border-t border-border/60 pt-3 text-xs">
            <span className="flex min-w-0 items-center gap-1.5 text-text-secondary">
              <Mail className="size-3.5 shrink-0" />
              <span className="truncate">{record.email ?? '未提供'}</span>
            </span>
            <span className="flex shrink-0 items-center gap-0.5 font-medium text-leaf-deep">
              审核
              <ChevronRight className="size-3.5 transition-transform duration-200 group-hover:translate-x-0.5" />
            </span>
          </div>
        </div>
      </div>
    </motion.button>
  )
}

/** 审核详情面板的可编辑表单状态 */
interface VerifyFormState {
  nickname: string
  amountYuan: string
  channelId: SnowflakeID | null
  message: string
  sponsorAt: string
  redirectUrl: string
  isAnonymous: boolean
  isHidden: boolean
  applyRemark: string
}

/** 从 SponsorRecordAdmin 初始化编辑表单 */
function initVerifyForm(record: SponsorRecordAdmin): VerifyFormState {
  return {
    nickname: record.nickname,
    amountYuan: formatAmount(record.amount),
    channelId: record.channel_id,
    message: record.message ?? '',
    sponsorAt: record.sponsor_at ? record.sponsor_at.slice(0, 10) : '',
    redirectUrl: record.redirect_url ?? '',
    isAnonymous: record.is_anonymous,
    isHidden: record.is_hidden,
    applyRemark: record.apply_remark ?? '',
  }
}

/** 将编辑表单转为 UpdateRecordRequest */
function verifyFormToUpdateReq(form: VerifyFormState): UpdateRecordRequest {
  return {
    nickname: form.nickname.trim(),
    amount: Math.round(Number(form.amountYuan) * 100),
    channel_id: form.channelId ?? null,
    message: form.message.trim() || undefined,
    sponsor_at: form.sponsorAt
      ? new Date(`${form.sponsorAt}T12:00:00`).toISOString()
      : undefined,
    redirect_url: form.redirectUrl.trim() || undefined,
    is_anonymous: form.isAnonymous,
    is_hidden: form.isHidden,
  }
}

/** 选中后右侧的审核详情面板（可编辑赞助信息 + 审核操作） */
function VerifyDetail({
  record,
  form,
  onFormChange,
  remark,
  onRemarkChange,
  onApprove,
  onReject,
  isPending,
}: {
  record: SponsorRecordAdmin
  form: VerifyFormState
  onFormChange: (patch: Partial<VerifyFormState>) => void
  remark: string
  onRemarkChange: (value: string) => void
  onApprove: () => void
  onReject: () => void
  isPending: boolean
}) {
  const status = sponsorStatus(record)
  const channels = useAllChannels().data ?? []

  return (
    <div className="space-y-4">
      {/* 申请横幅 */}
      <div className={`${inkCard} relative overflow-hidden p-0`}>
        <div
          className="pointer-events-none absolute inset-0"
          style={{
            background:
              'radial-gradient(520px 240px at 88% 0%, oklch(0.88 0.1 105 / 0.22), transparent 70%)',
          }}
          aria-hidden="true"
        />
        <div className="relative flex flex-wrap items-center gap-5 p-6">
          <div className="flex size-16 shrink-0 items-center justify-center rounded-2xl bg-muted font-serif text-3xl text-text-secondary">
            {form.nickname || record.nickname}
          </div>
          <div className="min-w-0 flex-1">
            <div className="flex flex-wrap items-center gap-2.5">
              <h2 className="font-serif text-2xl font-bold tracking-tight text-text-primary">
                {form.nickname || record.nickname}
              </h2>
              <InkBadge tone={status.tone}>{status.label}</InkBadge>
            </div>
            <p className="mt-2 flex items-center gap-1.5 font-mono text-sm text-text-secondary">
              <Coins className="size-4" />¥
              {form.amountYuan || formatAmount(record.amount)}
              {form.channelId || record.channel
                ? ` · ${record.channel?.name ?? '渠道'}`
                : ''}
            </p>
          </div>
        </div>
      </div>

      {/* 可编辑赞助信息 */}
      <div className={inkCard}>
        <CardHead title="赞助信息" meta="EDITABLE" />
        <div className="grid gap-4 md:grid-cols-2">
          <div className="space-y-2">
            <Label htmlFor="v-nickname">
              昵称 <span className="text-destructive">*</span>
            </Label>
            <Input
              id="v-nickname"
              value={form.nickname}
              onChange={(e) => onFormChange({ nickname: e.target.value })}
              disabled={isPending}
            />
          </div>
          <div className="space-y-2">
            <Label htmlFor="v-amount">
              赞助金额（元） <span className="text-destructive">*</span>
            </Label>
            <Input
              id="v-amount"
              type="number"
              min="0.01"
              step="0.01"
              value={form.amountYuan}
              onChange={(e) => onFormChange({ amountYuan: e.target.value })}
              disabled={isPending}
            />
          </div>
          <div className="space-y-2">
            <Label htmlFor="v-channelId">赞助渠道</Label>
            <Select
              id="v-channelId"
              value={form.channelId?.toString() ?? ''}
              onChange={(e) =>
                onFormChange({
                  channelId: e.target.value ? BigInt(e.target.value) : null,
                })
              }
              disabled={isPending}
            >
              <option value="">未选择渠道</option>
              {channels.map((ch) => (
                <option key={ch.id.toString()} value={ch.id.toString()}>
                  {ch.name}
                </option>
              ))}
            </Select>
          </div>
          <div className="space-y-2">
            <Label htmlFor="v-sponsorAt">赞助时间</Label>
            <Input
              id="v-sponsorAt"
              type="date"
              value={form.sponsorAt}
              onChange={(e) => onFormChange({ sponsorAt: e.target.value })}
              disabled={isPending}
            />
          </div>
          <div className="space-y-2 md:col-span-2">
            <Label htmlFor="v-redirectUrl">跳转链接</Label>
            <Input
              id="v-redirectUrl"
              type="url"
              value={form.redirectUrl}
              onChange={(e) => onFormChange({ redirectUrl: e.target.value })}
              disabled={isPending}
            />
          </div>
        </div>
        <div className="mt-4 flex flex-wrap gap-6">
          <div className="flex items-center gap-2">
            <Checkbox
              id="v-anonymous"
              checked={form.isAnonymous}
              onCheckedChange={(checked) =>
                onFormChange({ isAnonymous: checked === true })
              }
              disabled={isPending}
            />
            <Label htmlFor="v-anonymous" className="cursor-pointer">
              匿名展示
            </Label>
          </div>
          <div className="flex items-center gap-2">
            <Checkbox
              id="v-hidden"
              checked={form.isHidden}
              onCheckedChange={(checked) =>
                onFormChange({ isHidden: checked === true })
              }
              disabled={isPending}
            />
            <Label htmlFor="v-hidden" className="cursor-pointer">
              前台隐藏
            </Label>
          </div>
        </div>
        <div className="mt-4 space-y-2">
          <Label htmlFor="v-message">留言</Label>
          <Textarea
            id="v-message"
            className="min-h-[70px]"
            value={form.message}
            onChange={(e) => onFormChange({ message: e.target.value })}
            disabled={isPending}
          />
        </div>
        <div className="mt-4 space-y-2">
          <Label htmlFor="v-remark">申请备注</Label>
          <Textarea
            id="v-remark"
            className="min-h-[50px]"
            value={form.applyRemark}
            onChange={(e) => onFormChange({ applyRemark: e.target.value })}
            disabled={isPending}
          />
        </div>
      </div>

      {/* 审核操作 */}
      <div className={inkCard}>
        <CardHead title="审核操作" meta="REVIEW" />
        <div className="space-y-2">
          <label
            htmlFor="review-remark"
            className="font-mono text-[11px] uppercase tracking-widest text-text-secondary"
          >
            审核备注（可选）
          </label>
          <Textarea
            id="review-remark"
            placeholder="填写审核说明，将反馈给申请人…"
            value={remark}
            onChange={(e) => onRemarkChange(e.target.value)}
            rows={3}
            disabled={isPending}
          />
        </div>
        <div className="mt-4 flex flex-wrap items-center justify-between gap-3">
          <p className="text-sm text-text-secondary">
            审核通过后该条记录将展示在赞助页，上方修改会同步保存
          </p>
          <div className="flex gap-2">
            <Button
              variant="outline"
              className="cursor-pointer text-destructive hover:bg-destructive/10 hover:text-destructive"
              onClick={onReject}
              disabled={isPending}
            >
              <X className="mr-2 size-4" />
              {isPending ? '处理中…' : '拒绝'}
            </Button>
            <Button
              className="cursor-pointer"
              onClick={onApprove}
              disabled={isPending}
            >
              <Check className="mr-2 size-4" />
              {isPending ? '处理中…' : '通过'}
            </Button>
          </div>
        </div>
      </div>
    </div>
  )
}

function SponsorVerifyPage() {
  const reduced = useReducedMotion() ?? false
  // 赞助审核仅展示待审核(0)
  const pendingQuery = useAdminRecords({ status: 0, page: 1, page_size: 50 })
  const isLoading = pendingQuery.isLoading
  const pendingRecords = pendingQuery.data?.data ?? []
  const updateStatus = useUpdateSponsorStatus()
  const updateRecord = useUpdateRecord()

  const [selectedId, setSelectedId] = useState<SnowflakeID | null>(null)
  const [remark, setRemark] = useState('')
  const [editForm, setEditForm] = useState<VerifyFormState | null>(null)

  const selected = pendingRecords.find((r) => r.id === selectedId) ?? null
  const hasSelection = selected !== null

  /** 选中项变化时初始化/重置编辑表单 */
  useEffect(() => {
    if (selected) {
      setEditForm(initVerifyForm(selected))
    } else {
      setEditForm(null)
    }
  }, [selectedId])

  const patchForm = useCallback((patch: Partial<VerifyFormState>) => {
    setEditForm((prev) => (prev ? { ...prev, ...patch } : prev))
  }, [])

  /** 点击已选中项则取消选中 */
  const toggle = (id: SnowflakeID) => {
    setSelectedId((prev) => (prev === id ? null : id))
    setRemark('')
  }

  // 主操作：先保存编辑 → 再改状态
  const handleApprove = () => {
    if (!selected || !editForm) return
    const updateReq = verifyFormToUpdateReq(editForm)

    updateRecord.mutate(
      { id: selected.id, req: updateReq },
      {
        onSuccess: () => {
          updateStatus.mutate(
            {
              id: selected.id,
              req: {
                sponsor_status: 1,
                sponsor_review_remark: remark.trim() || undefined,
              },
            },
            {
              onSuccess: () => {
                setSelectedId(null)
                setRemark('')
              },
            },
          )
        },
      },
    )
  }

  // 次操作：先保存编辑 → 再改状态
  const handleReject = () => {
    if (!selected || !editForm) return
    const updateReq = verifyFormToUpdateReq(editForm)

    updateRecord.mutate(
      { id: selected.id, req: updateReq },
      {
        onSuccess: () => {
          updateStatus.mutate(
            {
              id: selected.id,
              req: {
                sponsor_status: 2,
                sponsor_review_remark: remark.trim() || undefined,
              },
            },
            {
              onSuccess: () => {
                setSelectedId(null)
                setRemark('')
              },
            },
          )
        },
      },
    )
  }

  // 列表刷新后若当前选中项已不在列表中，则清空选中
  useEffect(() => {
    if (selectedId && !pendingRecords.some((r) => r.id === selectedId)) {
      setSelectedId(null)
      setRemark('')
    }
  }, [selectedId, pendingRecords])

  return (
    <div className="mx-auto max-w-6xl space-y-6">
      <PageHead
        backTo="/admin/sponsor"
        backLabel="返回赞助管理"
        kicker="REVIEW · 审核"
        title="赞助审核"
        sub="逐条审阅新提交的赞助展示申请，通过后即可展示在赞助页。"
        actions={
          pendingRecords.length > 0 ? (
            <InkBadge tone="pending" className="px-2.5 py-1 text-sm">
              {pendingRecords.length} 项待处理
            </InkBadge>
          ) : undefined
        }
      />

      <BambooRule reduced={reduced} delay={0.12} />

      {isLoading ? (
        <div className="grid flex-1 gap-3 sm:grid-cols-2 xl:grid-cols-3">
          {Array.from({ length: 6 }).map((_, i) => (
            <Skeleton key={i} className="h-40 w-full rounded-lg" />
          ))}
        </div>
      ) : pendingRecords.length === 0 ? (
        /* 空状态 */
        <motion.div
          {...enter(reduced, 0.2, {
            initial: { opacity: 0, y: 10 },
            animate: { opacity: 1, y: 0 },
            transition: { duration: 0.4, ease: 'easeOut' },
          })}
          className={inkTableWrap}
        >
          <div className="py-8">
            <EnsoEmpty
              title="暂无待审核的赞助申请"
              hint="新的赞助展示申请会出现在这里"
            >
              <Link to="/admin/sponsor" className="ml-auto">
                <Button variant="outline" size="sm" className="cursor-pointer">
                  返回列表
                </Button>
              </Link>
            </EnsoEmpty>
          </div>
        </motion.div>
      ) : (
        /*
         * 主从布局：列表始终完整存在。
         * 未选中 → 列表占满（多列网格，富密度）；
         * 选中后 → 列表原地收紧为窄列，详情面板从右侧弹出。
         */
        <div className="flex flex-col gap-4 lg:flex-row lg:items-start">
          <motion.div
            layout
            transition={morphSpring}
            className={cn(
              'min-w-0',
              hasSelection
                ? 'flex flex-col gap-2 lg:w-60 lg:shrink-0'
                : 'grid flex-1 gap-3 sm:grid-cols-2 xl:grid-cols-3',
            )}
          >
            {pendingRecords.map((record) => (
              <VerifyItem
                key={record.id.toString()}
                record={record}
                compact={hasSelection}
                active={record.id === selectedId}
                onSelect={toggle}
              />
            ))}
          </motion.div>

          {/* 详情面板：慢慢弹出 */}
          {selected && (
            <motion.div
              key="detail"
              initial={{ opacity: 0, x: 48, scale: 0.98 }}
              animate={{ opacity: 1, x: 0, scale: 1 }}
              transition={morphSpring}
              className="min-w-0 flex-1"
            >
              <VerifyDetail
                record={selected}
                form={editForm ?? initVerifyForm(selected)}
                onFormChange={patchForm}
                remark={remark}
                onRemarkChange={setRemark}
                onApprove={handleApprove}
                onReject={handleReject}
                isPending={updateStatus.isPending || updateRecord.isPending}
              />
            </motion.div>
          )}
        </div>
      )}
    </div>
  )
}
