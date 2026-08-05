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
import {
  Check,
  ChevronRight,
  ExternalLink,
  Globe,
  Mail,
  MapPin,
  Palette,
  X,
} from 'lucide-react'
import type { LinkFriend, SnowflakeID, UpdateLinkRequest } from '@/api/types'
import { Button } from '@/components/ui/button'
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
  linkStatus,
} from '@/components/ink-wash'
import { enter } from '@/lib/motion'
import { AccentBar } from '@/components/link/accent-bar'
import { cn } from '@/lib/utils'
import {
  useAdminLinks,
  useApproveEditRequest,
  useRejectEditRequest,
  useUpdateLink,
  useUpdateLinkStatus,
} from '@/hooks/use-links'
import { useAllGroups } from '@/hooks/use-groups'
import { useAllColors } from '@/hooks/use-colors'

export const Route = createFileRoute('/_admin/admin/link/verify')({
  component: LinkVerifyPage,
})

/** 布局形变统一弹簧：列表收窄与详情弹出共用，节奏一致 */
const morphSpring = { type: 'spring', stiffness: 320, damping: 32 } as const

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
        className="size-full object-cover"
        loading="lazy"
        onError={() => setFailed(true)}
      />
    </div>
  )
}

/**
 * 列表条目 —— 富/紧凑两种密度共用同一组件。
 * 选中其他项时通过 layout 动画原地收紧，而非替换成另一个组件，
 * 保证「未选中 → 选中」是同一个界面的自然演化。
 */
function VerifyItem({
  link,
  compact,
  active,
  onSelect,
}: {
  link: LinkFriend
  compact: boolean
  active: boolean
  onSelect: (id: SnowflakeID) => void
}) {
  const status = linkStatus(link)
  return (
    <motion.button
      layout
      type="button"
      onClick={() => onSelect(link.id)}
      transition={morphSpring}
      className={cn(
        'group w-full cursor-pointer overflow-hidden rounded-lg border bg-card text-left transition-colors duration-200',
        active
          ? 'border-leaf-deep/60 bg-leaf-deep/6 shadow-[0_10px_24px_-18px_oklch(0.55_0.12_155/0.5)]'
          : 'border-border hover:border-leaf-muted hover:shadow-[0_10px_24px_-20px_oklch(0.32_0.06_155/0.35)]',
        compact ? 'p-2.5' : 'p-4',
      )}
    >
      {/* 常驻头部：头像 + 名称 + 元信息 */}
      <div className="flex items-center gap-3">
        <SiteAvatar
          name={link.name}
          url={link.avatar}
          className={cn(
            'transition-all duration-300',
            compact ? 'size-8 text-xs' : 'size-11 text-sm',
          )}
        />
        <div className="min-w-0 flex-1">
          <div className="flex items-center justify-between gap-2">
            <h3
              className={cn(
                'truncate font-serif font-semibold text-text-primary transition-all duration-300',
                compact ? 'text-sm' : 'text-base',
              )}
            >
              {link.name}
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
              {link.group_f_key?.name ?? '未分组'}
            </span>
            <span className="shrink-0 tabular-nums">
              {new Date(link.updated_at).toLocaleDateString('zh-CN')}
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
            title={link.description || undefined}
          >
            {link.description || '暂无描述'}
          </p>
          <div className="mt-3 flex items-center justify-between gap-2 border-t border-border/60 pt-3 text-xs">
            <span className="flex min-w-0 items-center gap-1.5 text-text-secondary">
              <Mail className="size-3.5 shrink-0" />
              <span className="truncate">{link.email ?? '未提供'}</span>
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
  siteName: string
  siteUrl: string
  siteLogo: string
  siteRss: string
  webmasterEmail: string
  siteDescription: string
  groupId: SnowflakeID | null
  colorId: SnowflakeID | null
  applyRemark: string
}

/** 从 LinkFriend 初始化编辑表单 */
function initVerifyForm(link: LinkFriend): VerifyFormState {
  return {
    siteName: link.name,
    siteUrl: link.url,
    siteLogo: link.avatar ?? '',
    siteRss: link.rss ?? '',
    webmasterEmail: link.email ?? '',
    siteDescription: link.description ?? '',
    groupId: link.group_id,
    colorId: link.color_id,
    applyRemark: link.apply_remark ?? '',
  }
}

/** 将编辑表单转为 UpdateLinkRequest */
function verifyFormToUpdateReq(form: VerifyFormState): UpdateLinkRequest {
  return {
    link_name: form.siteName.trim(),
    link_url: form.siteUrl.trim(),
    link_avatar: form.siteLogo.trim() || undefined,
    link_rss: form.siteRss.trim() || undefined,
    link_email: form.webmasterEmail.trim() || undefined,
    link_desc: form.siteDescription.trim() || undefined,
    link_group_id: form.groupId ?? null,
    link_color_id: form.colorId ?? null,
    link_apply_remark: form.applyRemark.trim() || undefined,
  }
}

/** 选中后右侧的审核详情面板（可编辑站点信息 + 审核操作） */
function VerifyDetail({
  link,
  form,
  onFormChange,
  remark,
  onRemarkChange,
  onApprove,
  onReject,
  isPending,
}: {
  link: LinkFriend
  form: VerifyFormState
  onFormChange: (patch: Partial<VerifyFormState>) => void
  remark: string
  onRemarkChange: (value: string) => void
  onApprove: () => void
  onReject: () => void
  isPending: boolean
}) {
  const status = linkStatus(link)
  const groups = useAllGroups().data ?? []
  const colors = useAllColors().data ?? []

  return (
    <div className="space-y-4">
      {/* 申请横幅：左侧墨条 + 晨光墨晕呼应待审核状态 */}
      <div className={`${inkCard} group relative overflow-hidden p-0`}>
        <AccentBar color={link.color_f_key} className="inset-y-0 z-10 w-1.5" />
        <div
          className="pointer-events-none absolute inset-0"
          style={{
            background:
              'radial-gradient(520px 240px at 88% 0%, oklch(0.88 0.1 105 / 0.22), transparent 70%)',
          }}
          aria-hidden="true"
        />
        <div className="relative flex flex-wrap items-center gap-5 p-6 pl-9">
          <SiteAvatar
            name={form.siteName || link.name}
            url={form.siteLogo || link.avatar}
            className="size-20 rounded-2xl text-3xl"
          />
          <div className="min-w-0 flex-1">
            <div className="flex flex-wrap items-center gap-2.5">
              <h2 className="font-serif text-2xl font-bold tracking-tight text-text-primary">
                {form.siteName || link.name}
              </h2>
              <InkBadge tone={status.tone}>{status.label}</InkBadge>
            </div>
            <a
              href={form.siteUrl || link.url}
              target="_blank"
              rel="noopener noreferrer"
              className="mt-2 inline-flex items-center gap-1.5 font-mono text-sm text-text-secondary transition-colors hover:text-leaf-deep"
            >
              <Globe className="size-4" />
              {form.siteUrl || link.url}
              <ExternalLink className="size-3.5" />
            </a>
          </div>
        </div>
      </div>

      {/* 可编辑站点信息 */}
      <div className={inkCard}>
        <CardHead title="站点信息" meta="EDITABLE" />
        <div className="grid gap-4 md:grid-cols-2">
          <div className="space-y-2">
            <Label htmlFor="v-siteName">
              站点名称 <span className="text-destructive">*</span>
            </Label>
            <Input
              id="v-siteName"
              value={form.siteName}
              onChange={(e) => onFormChange({ siteName: e.target.value })}
              disabled={isPending}
            />
          </div>
          <div className="space-y-2">
            <Label htmlFor="v-siteUrl">
              站点地址 <span className="text-destructive">*</span>
            </Label>
            <Input
              id="v-siteUrl"
              type="url"
              value={form.siteUrl}
              onChange={(e) => onFormChange({ siteUrl: e.target.value })}
              disabled={isPending}
            />
          </div>
          <div className="space-y-2">
            <Label htmlFor="v-siteLogo">站点 Logo</Label>
            <Input
              id="v-siteLogo"
              value={form.siteLogo}
              onChange={(e) => onFormChange({ siteLogo: e.target.value })}
              disabled={isPending}
            />
          </div>
          <div className="space-y-2">
            <Label htmlFor="v-siteRss">订阅地址</Label>
            <Input
              id="v-siteRss"
              type="url"
              value={form.siteRss}
              onChange={(e) => onFormChange({ siteRss: e.target.value })}
              disabled={isPending}
            />
          </div>
          <div className="space-y-2">
            <Label htmlFor="v-email">站长邮箱</Label>
            <Input
              id="v-email"
              type="email"
              value={form.webmasterEmail}
              onChange={(e) => onFormChange({ webmasterEmail: e.target.value })}
              disabled={isPending}
            />
          </div>
          <div className="space-y-2">
            <Label htmlFor="v-groupId" className="flex items-center gap-1.5">
              <MapPin className="size-3.5 text-text-secondary" />
              展示位置
            </Label>
            <Select
              id="v-groupId"
              value={form.groupId?.toString() ?? ''}
              onChange={(e) =>
                onFormChange({
                  groupId: e.target.value ? BigInt(e.target.value) : null,
                })
              }
              disabled={isPending}
            >
              <option value="">未分组</option>
              {groups.map((g) => (
                <option key={g.id.toString()} value={g.id.toString()}>
                  {g.name}
                </option>
              ))}
            </Select>
          </div>
          <div className="space-y-2">
            <Label htmlFor="v-colorId" className="flex items-center gap-1.5">
              <Palette className="size-3.5 text-text-secondary" />
              展示颜色
            </Label>
            <Select
              id="v-colorId"
              value={form.colorId?.toString() ?? ''}
              onChange={(e) =>
                onFormChange({
                  colorId: e.target.value ? BigInt(e.target.value) : null,
                })
              }
              disabled={isPending}
            >
              <option value="">默认颜色</option>
              {colors.map((c) => (
                <option key={c.id.toString()} value={c.id.toString()}>
                  {c.name}
                </option>
              ))}
            </Select>
          </div>
        </div>
        <div className="mt-4 space-y-2">
          <Label htmlFor="v-desc">站点描述</Label>
          <Textarea
            id="v-desc"
            className="min-h-[80px]"
            value={form.siteDescription}
            onChange={(e) => onFormChange({ siteDescription: e.target.value })}
            disabled={isPending}
          />
        </div>
        <div className="mt-4 space-y-2">
          <Label htmlFor="v-remark">申请备注</Label>
          <Textarea
            id="v-remark"
            className="min-h-[60px]"
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
            审核通过后该站点将展示在友链页面，上方修改会同步保存
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

/** 修改申请列表条目（edit 段） */
function EditVerifyItem({
  link,
  active,
  onSelect,
}: {
  link: LinkFriend
  active: boolean
  onSelect: (id: SnowflakeID) => void
}) {
  const status = linkStatus(link)
  return (
    <motion.button
      layout
      type="button"
      onClick={() => onSelect(link.id)}
      transition={morphSpring}
      className={cn(
        'group w-full cursor-pointer overflow-hidden rounded-lg border bg-card p-4 text-left transition-colors duration-200',
        active
          ? 'border-leaf-deep/60 bg-leaf-deep/6 shadow-[0_10px_24px_-18px_oklch(0.55_0.12_155/0.5)]'
          : 'border-border hover:border-leaf-muted hover:shadow-[0_10px_24px_-20px_oklch(0.32_0.06_155/0.35)]',
      )}
    >
      <div className="flex items-center gap-3">
        <SiteAvatar
          name={link.name}
          url={link.avatar}
          className="size-10 text-sm"
        />
        <div className="min-w-0 flex-1">
          <div className="flex items-center justify-between gap-2">
            <h3 className="truncate font-serif text-base font-semibold text-text-primary">
              {link.name}
            </h3>
            <InkBadge tone={status.tone}>{status.label}</InkBadge>
          </div>
          <p className="mt-1 truncate text-sm text-text-secondary">
            {link.apply_remark || '暂无申请备注'}
          </p>
        </div>
      </div>
    </motion.button>
  )
}

/** 修改申请详情面板：当前 vs 预期对比 + 通过/拒绝 */
function EditVerifyDetail({
  link,
  remark,
  onRemarkChange,
  onApprove,
  onReject,
  isPending,
}: {
  link: LinkFriend
  remark: string
  onRemarkChange: (value: string) => void
  onApprove: () => void
  onReject: () => void
  isPending: boolean
}) {
  const groups = useAllGroups().data ?? []
  const colors = useAllColors().data ?? []
  const groupName = (id: SnowflakeID | null) =>
    id == null
      ? '未分组'
      : (groups.find((g) => g.id === id)?.name ?? `位置 #${id}`)
  const colorName = (id: SnowflakeID | null) =>
    id == null
      ? '默认颜色'
      : (colors.find((c) => c.id === id)?.name ?? `颜色 #${id}`)

  return (
    <div className="space-y-4">
      {/* 站点横幅 */}
      <div className={`${inkCard} relative overflow-hidden p-0`}>
        <div className="relative flex flex-wrap items-center gap-5 p-6">
          <SiteAvatar
            name={link.name}
            url={link.avatar}
            className="size-16 rounded-2xl text-2xl"
          />
          <div className="min-w-0 flex-1">
            <div className="flex flex-wrap items-center gap-2.5">
              <h2 className="font-serif text-2xl font-bold tracking-tight text-text-primary">
                {link.name}
              </h2>
              <InkBadge tone="pending">修改审核中</InkBadge>
            </div>
            <a
              href={link.url}
              target="_blank"
              rel="noopener noreferrer"
              className="mt-2 inline-flex items-center gap-1.5 font-mono text-sm text-text-secondary transition-colors hover:text-leaf-deep"
            >
              <Globe className="size-4" />
              {link.url}
              <ExternalLink className="size-3.5" />
            </a>
          </div>
        </div>
      </div>

      {/* 当前 vs 预期 */}
      <div className={inkCard}>
        <CardHead title="展示调整对比" meta="CURRENT → EXPECTED" />
        <div className="grid gap-4 md:grid-cols-2">
          <div className="space-y-2">
            <Label>当前展示位置</Label>
            <div className="rounded-md border border-border/70 bg-muted/40 px-3 py-2 text-sm text-text-secondary">
              {link.group_f_key?.name ??
                (link.group_id ? `位置 #${link.group_id}` : '未分组')}
            </div>
          </div>
          <div className="space-y-2">
            <Label>预期展示位置</Label>
            <div className="rounded-md border border-leaf-deep/40 bg-leaf-deep/6 px-3 py-2 text-sm font-medium text-leaf-deep">
              {groupName(link.expected_group_id)}
            </div>
          </div>
          <div className="space-y-2">
            <Label>当前展示颜色</Label>
            <div className="rounded-md border border-border/70 bg-muted/40 px-3 py-2 text-sm text-text-secondary">
              {link.color_f_key?.name ??
                (link.color_id ? `颜色 #${link.color_id}` : '默认颜色')}
            </div>
          </div>
          <div className="space-y-2">
            <Label>预期展示颜色</Label>
            <div className="rounded-md border border-leaf-deep/40 bg-leaf-deep/6 px-3 py-2 text-sm font-medium text-leaf-deep">
              {colorName(link.expected_color_id)}
            </div>
          </div>
        </div>
        <div className="mt-4 space-y-2">
          <Label>申请备注</Label>
          <div className="rounded-md border border-border/70 bg-muted/40 px-3 py-2 text-sm text-text-secondary">
            {link.apply_remark || '无'}
          </div>
        </div>
      </div>

      {/* 审核操作 */}
      <div className={inkCard}>
        <CardHead title="审核操作" meta="REVIEW" />
        <div className="space-y-2">
          <label
            htmlFor="edit-review-remark"
            className="font-mono text-[11px] uppercase tracking-widest text-text-secondary"
          >
            审核备注（可选）
          </label>
          <Textarea
            id="edit-review-remark"
            placeholder="填写审核说明，将反馈给申请人…"
            value={remark}
            onChange={(e) => onRemarkChange(e.target.value)}
            rows={3}
            disabled={isPending}
          />
        </div>
        <div className="mt-4 flex flex-wrap items-center justify-between gap-3">
          <p className="text-sm text-text-secondary">
            通过后将应用预期位置/颜色，拒绝则保持当前设置
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

function LinkVerifyPage() {
  const reduced = useReducedMotion() ?? false

  // 双段审核：新申请(0) / 修改位置·颜色(5)，其余状态分流至异常管理
  const [tab, setTab] = useState<'new' | 'edit'>('new')
  const newQuery = useAdminLinks({ link_status: 0, page: 1, page_size: 50 })
  const editQuery = useAdminLinks({ link_status: 5, page: 1, page_size: 50 })
  const isLoading = tab === 'new' ? newQuery.isLoading : editQuery.isLoading
  const pendingLinks = newQuery.data?.data ?? []
  const editPendingLinks = editQuery.data?.data ?? []
  const updateStatus = useUpdateLinkStatus()
  const updateLink = useUpdateLink()
  const approveEdit = useApproveEditRequest()
  const rejectEdit = useRejectEditRequest()

  // new 段选中态
  const [selectedId, setSelectedId] = useState<SnowflakeID | null>(null)
  const [remark, setRemark] = useState('')
  const [editForm, setEditForm] = useState<VerifyFormState | null>(null)

  // edit 段选中态
  const [editSelectedId, setEditSelectedId] = useState<SnowflakeID | null>(null)
  const [editRemark, setEditRemark] = useState('')

  const selected = pendingLinks.find((l) => l.id === selectedId) ?? null
  const hasSelection = selected !== null
  const editSelected =
    editPendingLinks.find((l) => l.id === editSelectedId) ?? null

  /** 切换审核段并清空各自选中态 */
  const switchTab = (next: 'new' | 'edit') => {
    setTab(next)
    setSelectedId(null)
    setEditSelectedId(null)
    setRemark('')
    setEditRemark('')
  }

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

  /** edit 段：点击已选中项则取消选中 */
  const toggleEdit = (id: SnowflakeID) => {
    setEditSelectedId((prev) => (prev === id ? null : id))
    setEditRemark('')
  }

  /** edit 段：通过修改申请（应用预期位置/颜色） */
  const handleApproveEdit = () => {
    if (!editSelected) return
    approveEdit.mutate(
      {
        id: editSelected.id,
        req: { link_review_remark: editRemark.trim() || undefined },
      },
      { onSuccess: () => setEditSelectedId(null) },
    )
  }

  /** edit 段：拒绝修改申请（保持原位置/颜色） */
  const handleRejectEdit = () => {
    if (!editSelected) return
    rejectEdit.mutate(
      {
        id: editSelected.id,
        req: { link_review_remark: editRemark.trim() || undefined },
      },
      { onSuccess: () => setEditSelectedId(null) },
    )
  }

  // 主操作：先保存编辑 → 再改状态
  const handleApprove = () => {
    if (!selected || !editForm) return
    const updateReq = verifyFormToUpdateReq(editForm)

    // 先保存站点信息修改
    updateLink.mutate(
      { id: selected.id, req: updateReq },
      {
        onSuccess: () => {
          // 保存成功后再改审核状态
          updateStatus.mutate(
            {
              id: selected.id,
              req: {
                link_status: 1,
                link_review_remark: remark.trim() || undefined,
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

    updateLink.mutate(
      { id: selected.id, req: updateReq },
      {
        onSuccess: () => {
          updateStatus.mutate(
            {
              id: selected.id,
              req: {
                link_status: 2,
                link_review_remark: remark.trim() || undefined,
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
    if (selectedId && !pendingLinks.some((l) => l.id === selectedId)) {
      setSelectedId(null)
      setRemark('')
    }
  }, [selectedId, pendingLinks])

  useEffect(() => {
    if (
      editSelectedId &&
      !editPendingLinks.some((l) => l.id === editSelectedId)
    ) {
      setEditSelectedId(null)
      setEditRemark('')
    }
  }, [editSelectedId, editPendingLinks])

  const currentList = tab === 'new' ? pendingLinks : editPendingLinks

  return (
    <div className="mx-auto max-w-6xl space-y-6">
      <PageHead
        backTo="/admin/link"
        backLabel="返回友链管理"
        kicker="REVIEW · 审核"
        title="友链审核"
        sub="逐条审阅友链申请与修改申请，通过后即可生效。"
        actions={
          currentList.length > 0 ? (
            <InkBadge tone="pending" className="px-2.5 py-1 text-sm">
              {currentList.length} 项待处理
            </InkBadge>
          ) : undefined
        }
      />

      <BambooRule reduced={reduced} delay={0.12} />

      {/* 双段切换：新申请 / 修改位置·颜色 */}
      <div className="flex items-center gap-1 rounded-lg border border-border bg-card p-1">
        <button
          type="button"
          onClick={() => switchTab('new')}
          className={cn(
            'flex-1 cursor-pointer rounded-md px-3 py-1.5 text-sm font-medium transition-colors',
            tab === 'new'
              ? 'bg-leaf-deep/10 text-leaf-deep'
              : 'text-text-secondary hover:bg-muted/60 hover:text-text-primary',
          )}
        >
          新申请
          <span className="ml-1.5 text-xs text-text-secondary">
            ({newQuery.data?.data.length ?? 0})
          </span>
        </button>
        <button
          type="button"
          onClick={() => switchTab('edit')}
          className={cn(
            'flex-1 cursor-pointer rounded-md px-3 py-1.5 text-sm font-medium transition-colors',
            tab === 'edit'
              ? 'bg-leaf-deep/10 text-leaf-deep'
              : 'text-text-secondary hover:bg-muted/60 hover:text-text-primary',
          )}
        >
          修改位置·颜色
          <span className="ml-1.5 text-xs text-text-secondary">
            ({editQuery.data?.data.length ?? 0})
          </span>
        </button>
      </div>

      {isLoading ? (
        <div className="grid flex-1 gap-3 sm:grid-cols-2 xl:grid-cols-3">
          {Array.from({ length: 6 }).map((_, i) => (
            <Skeleton key={i} className="h-40 w-full rounded-lg" />
          ))}
        </div>
      ) : currentList.length === 0 ? (
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
              title={
                tab === 'new' ? '暂无待审核的友链申请' : '暂无待审核的修改申请'
              }
              hint={
                tab === 'new'
                  ? '新的友链申请会出现在这里'
                  : '用户提交的展示位置/颜色修改申请会出现在这里'
              }
            >
              <Link to="/admin/link" className="ml-auto">
                <Button variant="outline" size="sm" className="cursor-pointer">
                  返回列表
                </Button>
              </Link>
            </EnsoEmpty>
          </div>
        </motion.div>
      ) : tab === 'new' ? (
        /*
         * new 段主从布局：列表始终完整存在。
         * 未选中 → 列表占满（多列网格，富密度）；
         * 选中后 → 列表原地收紧为窄列，详情面板从右侧弹出。
         * 全程同一批条目、同一个界面，仅密度与占比变化。
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
            {pendingLinks.map((link) => (
              <VerifyItem
                key={link.id.toString()}
                link={link}
                compact={hasSelection}
                active={link.id === selectedId}
                onSelect={toggle}
              />
            ))}
          </motion.div>

          {/* 详情面板：慢慢弹出（仅入场动画，列表让位由 layout 完成） */}
          {selected && (
            <motion.div
              key="detail"
              initial={{ opacity: 0, x: 48, scale: 0.98 }}
              animate={{ opacity: 1, x: 0, scale: 1 }}
              transition={morphSpring}
              className="min-w-0 flex-1"
            >
              <VerifyDetail
                link={selected}
                form={editForm ?? initVerifyForm(selected)}
                onFormChange={patchForm}
                remark={remark}
                onRemarkChange={setRemark}
                onApprove={handleApprove}
                onReject={handleReject}
                isPending={updateStatus.isPending || updateLink.isPending}
              />
            </motion.div>
          )}
        </div>
      ) : (
        /* edit 段主从布局 */
        <div className="flex flex-col gap-4 lg:flex-row lg:items-start">
          <motion.div
            layout
            transition={morphSpring}
            className={cn(
              'min-w-0',
              editSelected
                ? 'flex flex-col gap-2 lg:w-60 lg:shrink-0'
                : 'grid flex-1 gap-3 sm:grid-cols-2 xl:grid-cols-3',
            )}
          >
            {editPendingLinks.map((link) => (
              <EditVerifyItem
                key={link.id.toString()}
                link={link}
                active={link.id === editSelectedId}
                onSelect={toggleEdit}
              />
            ))}
          </motion.div>

          {editSelected && (
            <motion.div
              key="edit-detail"
              initial={{ opacity: 0, x: 48, scale: 0.98 }}
              animate={{ opacity: 1, x: 0, scale: 1 }}
              transition={morphSpring}
              className="min-w-0 flex-1"
            >
              <EditVerifyDetail
                link={editSelected}
                remark={editRemark}
                onRemarkChange={setEditRemark}
                onApprove={handleApproveEdit}
                onReject={handleRejectEdit}
                isPending={approveEdit.isPending || rejectEdit.isPending}
              />
            </motion.div>
          )}
        </div>
      )}
    </div>
  )
}
