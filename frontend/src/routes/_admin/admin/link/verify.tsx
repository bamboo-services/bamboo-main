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
import {
  Check,
  ChevronRight,
  ExternalLink,
  Globe,
  Mail,
  MapPin,
  RefreshCcw,
  X,
} from 'lucide-react'
import type { LinkFriend, SnowflakeID } from '@/api/types'
import { Button } from '@/components/ui/button'
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
import { cn } from '@/lib/utils'
import { useAdminLinks, useUpdateLinkStatus } from '@/hooks/use-links'

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
          <p className="mt-3 line-clamp-2 text-sm leading-relaxed text-text-secondary">
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

/** 选中后右侧的审核详情面板 */
function VerifyDetail({
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
  const isTakedown = link.status === 3
  const status = linkStatus(link)
  return (
    <div className="space-y-4">
      {/* 申请横幅：左侧墨条 + 晨光墨晕呼应待审核状态 */}
      <div className={`${inkCard} relative overflow-hidden p-0`}>
        <span
          className="absolute inset-y-0 left-0 z-10 w-1.5 bg-leaf-deep"
          aria-hidden="true"
        />
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
            name={link.name}
            url={link.avatar}
            className="size-20 rounded-2xl text-3xl"
          />
          <div className="min-w-0 flex-1">
            <div className="flex flex-wrap items-center gap-2.5">
              <h2 className="font-serif text-2xl font-bold tracking-tight text-text-primary">
                {link.name}
              </h2>
              <InkBadge tone={status.tone}>{status.label}</InkBadge>
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

      {/* 描述 + 申请信息 */}
      <div className="grid gap-4 xl:grid-cols-3">
        <div className={`${inkCard} xl:col-span-2`}>
          <CardHead title="站点描述" />
          <p className="leading-relaxed text-text-secondary">
            {link.description || '暂无描述'}
          </p>
          {link.apply_remark && (
            <div className="mt-4 border-t border-border/60 pt-3">
              <p className="font-mono text-[11px] uppercase tracking-widest text-text-secondary">
                申请备注
              </p>
              <p className="mt-1 leading-relaxed text-text-primary">
                {link.apply_remark}
              </p>
            </div>
          )}
        </div>

        <div className={inkCard}>
          <CardHead title="申请信息" />
          <div className="space-y-3">
            <div className="flex items-center gap-2.5 text-sm text-text-primary">
              <Mail className="size-4 shrink-0 text-text-secondary" />
              <span className="truncate">{link.email ?? '未提供'}</span>
            </div>
            <div className="flex items-center gap-2.5 text-sm text-text-primary">
              <MapPin className="size-4 shrink-0 text-text-secondary" />
              <span>{link.group_f_key?.name ?? '未分组'}</span>
            </div>
            <div className="flex items-center gap-2.5 text-sm text-text-primary">
              <RefreshCcw className="size-4 shrink-0 text-text-secondary" />
              <span className="font-mono tabular-nums">
                {new Date(link.updated_at).toLocaleString('zh-CN')}
              </span>
            </div>
          </div>
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
            {isTakedown
              ? '批准下架后该站点将从友链页面移除'
              : '审核通过后该站点将展示在友链页面'}
          </p>
          <div className="flex gap-2">
            <Button
              variant="outline"
              className="cursor-pointer text-destructive hover:bg-destructive/10 hover:text-destructive"
              onClick={onReject}
              disabled={isPending}
            >
              <X className="mr-2 size-4" />
              {isPending ? '处理中…' : isTakedown ? '驳回下架' : '拒绝'}
            </Button>
            <Button
              className="cursor-pointer"
              onClick={onApprove}
              disabled={isPending}
            >
              <Check className="mr-2 size-4" />
              {isPending ? '处理中…' : isTakedown ? '批准下架' : '通过'}
            </Button>
          </div>
        </div>
      </div>
    </div>
  )
}

function LinkVerifyPage() {
  const reduced = useReducedMotion() ?? false
  // 同时拉取「待审核(0)」与「下架待审核(3)」两类待处理友链
  const pendingQuery = useAdminLinks({ link_status: 0, page: 1, page_size: 50 })
  const takedownQuery = useAdminLinks({ link_status: 3, page: 1, page_size: 50 })
  const isLoading = pendingQuery.isLoading || takedownQuery.isLoading
  const pendingLinks = [
    ...(pendingQuery.data?.data ?? []),
    ...(takedownQuery.data?.data ?? []),
  ]
  const updateStatus = useUpdateLinkStatus()

  const [selectedId, setSelectedId] = useState<SnowflakeID | null>(null)
  const [remark, setRemark] = useState('')

  const selected = pendingLinks.find((l) => l.id === selectedId) ?? null
  const hasSelection = selected !== null

  /** 点击已选中项则取消选中 */
  const toggle = (id: SnowflakeID) =>
    setSelectedId((prev) => (prev === id ? null : id))

  // 主操作：新申请 → 通过(1)；下架申请 → 批准下架(4)
  const handleApprove = () => {
    if (!selected) return
    const targetStatus = selected.status === 3 ? 4 : 1
    updateStatus.mutate(
      {
        id: selected.id,
        req: {
          link_status: targetStatus,
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
  }

  // 次操作：新申请 → 拒绝(2)；下架申请 → 驳回下架（恢复已通过 1）
  const handleReject = () => {
    if (!selected) return
    const targetStatus = selected.status === 3 ? 1 : 2
    updateStatus.mutate(
      {
        id: selected.id,
        req: {
          link_status: targetStatus,
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
  }

  // 列表刷新后若当前选中项已不在列表中，则清空选中
  useEffect(() => {
    if (selectedId && !pendingLinks.some((l) => l.id === selectedId)) {
      setSelectedId(null)
      setRemark('')
    }
  }, [selectedId, pendingLinks])

  return (
    <div className="mx-auto max-w-6xl space-y-6">
      <PageHead
        backTo="/admin/link"
        backLabel="返回友链管理"
        kicker="REVIEW · 审核"
        title="友链审核"
        sub="逐条审阅新提交的友链申请，通过后即可展示在友链页面。"
        actions={
          pendingLinks.length > 0 ? (
            <InkBadge tone="pending" className="px-2.5 py-1 text-sm">
              {pendingLinks.length} 项待处理
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
      ) : pendingLinks.length === 0 ? (
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
              title="暂无待审核的友链申请"
              hint="新的友链申请会出现在这里"
            >
              <Link to="/admin/link" className="ml-auto">
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
                remark={remark}
                onRemarkChange={setRemark}
                onApprove={handleApprove}
                onReject={handleReject}
                isPending={updateStatus.isPending}
              />
            </motion.div>
          )}
        </div>
      )}
    </div>
  )
}
