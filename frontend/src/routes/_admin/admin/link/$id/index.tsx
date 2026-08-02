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
import { Link, createFileRoute, useNavigate } from '@tanstack/react-router'
import { motion, useReducedMotion } from 'motion/react'
import {
  ArrowLeft,
  Camera,
  Check,
  Copy,
  ExternalLink,
  Globe,
  Mail,
  MapPin,
  MessageSquareText,
  Pencil,
  RefreshCcw,
  ShieldCheck,
  Trash2,
  Unlink,
} from 'lucide-react'
import type { LinkFriend } from '@/api/types'
import { Button } from '@/components/ui/button'
import { Avatar, AvatarFallback, AvatarImage } from '@/components/ui/avatar'
import { Skeleton } from '@/components/ui/skeleton'
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from '@/components/ui/dialog'
import {
  BambooArt,
  CardHead,
  InkBadge,
  inkCard,
  linkStatus,
} from '@/components/ink-wash'
import { enter } from '@/lib/motion'
import { accentOf, isFancyColor } from '@/lib/colors'
import { cn } from '@/lib/utils'
import { useAdminLink, useDeleteLink } from '@/hooks/use-links'

export const Route = createFileRoute('/_admin/admin/link/$id/')({
  component: LinkDetailPage,
})

/** 复制按钮：成功后短暂显示对勾反馈 */
function CopyButton({ text, label }: { text: string; label: string }) {
  const [copied, setCopied] = useState(false)
  const onCopy = () => {
    void navigator.clipboard.writeText(text).then(() => {
      setCopied(true)
      window.setTimeout(() => setCopied(false), 1500)
    })
  }
  return (
    <button
      type="button"
      onClick={onCopy}
      aria-label={label}
      title={label}
      className="inline-flex size-6 shrink-0 cursor-pointer items-center justify-center rounded-md text-text-secondary transition-colors duration-150 hover:bg-muted hover:text-text-primary"
    >
      {copied ? (
        <Check className="size-3.5 text-leaf-deep" />
      ) : (
        <Copy className="size-3.5" />
      )}
    </button>
  )
}

/** 信息行：左标签右值 */
function InfoRow({
  icon,
  label,
  children,
}: {
  icon: React.ReactNode
  label: string
  children: React.ReactNode
}) {
  return (
    <div className="flex items-center justify-between gap-3 border-b border-border/60 py-3 last:border-0">
      <span className="flex shrink-0 items-center gap-2 text-sm text-text-secondary">
        {icon}
        {label}
      </span>
      <span className="text-right text-sm font-medium text-text-primary">
        {children}
      </span>
    </div>
  )
}

/** 友链预览：1:1 还原公开页的展示卡片 */
function FriendPreview({ link, accent }: { link: LinkFriend; accent: string }) {
  return (
    <a
      href={link.url}
      target="_blank"
      rel="noopener noreferrer"
      className="group relative flex gap-3 overflow-hidden rounded-lg border border-border bg-card p-4 shadow-sm transition-[translate,border-color,box-shadow] duration-200 hover:-translate-y-0.5 hover:border-leaf-muted hover:shadow-[0_14px_30px_-22px_oklch(0.32_0.06_155/0.4)]"
    >
      <span
        className="absolute inset-y-0 left-0 w-1"
        style={{ background: accent }}
        aria-hidden="true"
      />
      <Avatar className="size-12 shrink-0 rounded-full">
        <AvatarImage src={link.avatar ?? undefined} alt={link.name} />
        <AvatarFallback className="font-serif">
          {link.name.slice(0, 1)}
        </AvatarFallback>
      </Avatar>
      <div className="min-w-0 flex-1">
        <div className="flex items-center justify-between gap-2">
          <h4 className="truncate font-serif font-semibold text-text-primary group-hover:text-leaf-deep">
            {link.name}
          </h4>
          <ExternalLink className="size-4 shrink-0 text-text-secondary opacity-0 transition-opacity duration-200 group-hover:opacity-100" />
        </div>
        <p className="mt-1 line-clamp-2 text-sm text-text-secondary">
          {link.description || '这个站点很神秘，没有留下描述。'}
        </p>
      </div>
    </a>
  )
}

function LinkDetailPage() {
  const reduced = useReducedMotion() ?? false
  const { id } = Route.useParams()
  const navigate = useNavigate()
  const [deleteOpen, setDeleteOpen] = useState(false)
  const [shotFailed, setShotFailed] = useState(false)

  const linkQuery = useAdminLink(BigInt(id))
  const deleteLink = useDeleteLink()
  const link = linkQuery.data

  const handleDelete = () => {
    if (!link) return
    deleteLink.mutate(link.id, {
      onSuccess: () => {
        setDeleteOpen(false)
        navigate({ to: '/admin/link' })
      },
    })
  }

  if (linkQuery.isLoading) {
    return (
      <div className="mx-auto max-w-6xl space-y-5">
        <Skeleton className="h-5 w-40" />
        <Skeleton className="h-40 w-full rounded-lg" />
        <div className="grid gap-4 lg:grid-cols-3">
          <div className="space-y-4 lg:col-span-2">
            <Skeleton className="h-32 w-full rounded-lg" />
            <Skeleton className="h-28 w-full rounded-lg" />
          </div>
          <Skeleton className="h-64 w-full rounded-lg" />
        </div>
      </div>
    )
  }

  if (!link) {
    return (
      <div className="mx-auto flex min-h-[60vh] max-w-6xl flex-col items-center justify-center gap-4 text-center">
        <Unlink className="size-10 text-text-secondary opacity-40" />
        <div>
          <h2 className="font-serif text-xl font-semibold text-text-primary">
            友链不存在
          </h2>
          <p className="mt-1 text-sm text-text-secondary">
            该友链可能已被删除或链接有误
          </p>
        </div>
        <Link to="/admin/link">
          <Button className="cursor-pointer">返回列表</Button>
        </Link>
      </div>
    )
  }

  const accent = accentOf(link.color_f_key)
  const fancy = isFancyColor(link.color_f_key)
  const s = linkStatus(link)

  return (
    <div className="mx-auto max-w-6xl space-y-5">
      {/* 返回：轻量文字链接，箭头悬停微移 */}
      <Link
        to="/admin/link"
        className="group inline-flex cursor-pointer items-center gap-1.5 font-mono text-xs text-text-secondary transition-colors duration-150 hover:text-leaf-deep"
      >
        <ArrowLeft className="size-3.5 transition-transform duration-150 group-hover:-translate-x-0.5" />
        返回友链管理
      </Link>

      {/* 站点档案头：宣纸卡 + 左侧站点色墨条 + 墨韵竹叶水印 */}
      <motion.section
        {...enter(reduced, 0.08, {
          initial: { opacity: 0, y: 12 },
          animate: { opacity: 1, y: 0 },
          transition: { duration: 0.5, ease: 'easeOut' },
        })}
        className={`${inkCard} relative overflow-hidden p-0`}
      >
        <span
          className={cn(
            'absolute inset-y-0 left-0 z-10 w-1.5',
            fancy && 'ink-fancy',
          )}
          style={{ background: accent }}
          aria-hidden="true"
        />
        {/* 晨光墨晕 + 墨韵竹叶 */}
        <div
          className="pointer-events-none absolute inset-0"
          style={{
            background:
              'radial-gradient(620px 300px at 88% 0%, oklch(0.88 0.1 105 / 0.20), transparent 70%)',
          }}
          aria-hidden="true"
        />
        <BambooArt className="pointer-events-none absolute top-0 right-[-40px] h-full w-[420px] text-text-primary opacity-70" />

        <div className="relative z-10 flex flex-wrap items-center gap-5 p-6 pl-9">
          <Avatar className="size-24 shrink-0 rounded-2xl shadow-lg ring-4 ring-card">
            <AvatarImage src={link.avatar ?? undefined} alt={link.name} />
            <AvatarFallback className="rounded-2xl bg-muted font-serif text-3xl font-semibold text-text-secondary">
              {link.name.charAt(0)}
            </AvatarFallback>
          </Avatar>

          <div className="min-w-0 flex-1">
            <div className="flex flex-wrap items-center gap-2.5">
              <h1 className="truncate font-serif text-2xl font-bold tracking-tight text-text-primary lg:text-3xl">
                {link.name}
              </h1>
              <InkBadge tone={s.tone}>{s.label}</InkBadge>
            </div>
            <div className="mt-2 flex flex-wrap items-center gap-x-4 gap-y-1.5 text-sm text-text-secondary">
              <span className="flex min-w-0 items-center gap-1.5">
                <Globe className="size-4 shrink-0" />
                <a
                  href={link.url}
                  target="_blank"
                  rel="noopener noreferrer"
                  className="truncate font-mono transition-colors duration-150 hover:text-leaf-deep"
                >
                  {link.url}
                </a>
                <CopyButton text={link.url} label="复制站点地址" />
              </span>
              <span className="flex items-center gap-1.5">
                <MapPin className="size-4" />
                {link.group_f_key?.name ?? '未分组'}
              </span>
              <span className="flex items-center gap-1.5">
                <span
                  className="size-2.5 rounded-full"
                  style={{ background: accent }}
                  aria-hidden="true"
                />
                {link.color_f_key?.name ?? '默认'}
              </span>
            </div>
          </div>

          <div className="flex shrink-0 flex-wrap gap-2">
            <a
              href={link.url}
              target="_blank"
              rel="noopener noreferrer"
              className="inline-flex h-9 cursor-pointer items-center justify-center gap-2 rounded-md bg-leaf-deep px-4 text-sm font-medium text-card transition-colors duration-150 hover:bg-leaf-deep/90"
            >
              <ExternalLink className="size-4" />
              访问站点
            </a>
            <Button
              variant="outline"
              className="cursor-pointer"
              onClick={() =>
                navigate({
                  to: '/admin/link/$id/edit',
                  params: { id: link.id.toString() },
                })
              }
            >
              <Pencil className="mr-2 size-4" />
              编辑
            </Button>
            <Button
              variant="outline"
              className="cursor-pointer text-destructive hover:bg-destructive/10 hover:text-destructive"
              onClick={() => setDeleteOpen(true)}
            >
              <Trash2 className="mr-2 size-4" />
              删除
            </Button>
          </div>
        </div>
      </motion.section>

      {/* 内容网格：描述 + 预览（2 份）｜ 基本信息（1 份） */}
      <div className="grid gap-4 lg:grid-cols-3">
        <div className="space-y-4 lg:col-span-2">
          <motion.section
            {...enter(reduced, 0.16, {
              initial: { opacity: 0, y: 10 },
              animate: { opacity: 1, y: 0 },
              transition: { duration: 0.4, ease: 'easeOut' },
            })}
            className={inkCard}
          >
            <CardHead title="站点描述" />
            <p className="leading-relaxed text-text-secondary">
              {link.description || '暂无描述'}
            </p>
          </motion.section>

          <motion.section
            {...enter(reduced, 0.22, {
              initial: { opacity: 0, y: 10 },
              animate: { opacity: 1, y: 0 },
              transition: { duration: 0.4, ease: 'easeOut' },
            })}
            className={inkCard}
          >
            <CardHead title="友链预览" meta="访客在公开页看到的卡片样式" />
            <FriendPreview link={link} accent={accent} />
          </motion.section>

          <motion.section
            {...enter(reduced, 0.25, {
              initial: { opacity: 0, y: 10 },
              animate: { opacity: 1, y: 0 },
              transition: { duration: 0.4, ease: 'easeOut' },
            })}
            className={inkCard}
          >
            <CardHead
              title="站点截图"
              meta={
                link.screenshot_at
                  ? `更新于 ${new Date(link.screenshot_at).toLocaleString('zh-CN')}`
                  : '每日 0 点自动更新'
              }
            />
            {link.screenshot_url && !shotFailed ? (
              <img
                src={link.screenshot_url}
                alt={`${link.name} 站点截图`}
                loading="lazy"
                onError={() => setShotFailed(true)}
                className="block w-full rounded-lg border border-border/70 bg-card shadow-sm"
              />
            ) : (
              <div className="flex h-44 flex-col items-center justify-center gap-2 rounded-lg border border-dashed border-border/70 bg-gradient-to-br from-leaf-light/15 via-card to-leaf-muted/10 text-center">
                <Camera className="size-7 text-text-secondary opacity-40" />
                <p className="font-mono text-xs text-text-secondary">
                  截图尚未生成
                </p>
                <p className="text-xs text-text-secondary/80">
                  审核通过后自动生成，此后每日 0 点更新
                </p>
              </div>
            )}
          </motion.section>

          {(link.apply_remark || link.review_remark) && (
            <motion.section
              {...enter(reduced, 0.28, {
                initial: { opacity: 0, y: 10 },
                animate: { opacity: 1, y: 0 },
                transition: { duration: 0.4, ease: 'easeOut' },
              })}
              className={inkCard}
            >
              <CardHead title="备注信息" />
              <div className="space-y-3">
                {link.apply_remark && (
                  <div className="flex items-start gap-2.5 text-sm">
                    <MessageSquareText className="mt-0.5 size-4 shrink-0 text-text-secondary" />
                    <div className="min-w-0">
                      <p className="font-mono text-[11px] uppercase tracking-widest text-text-secondary">
                        申请备注
                      </p>
                      <p className="mt-0.5 leading-relaxed text-text-primary">
                        {link.apply_remark}
                      </p>
                    </div>
                  </div>
                )}
                {link.review_remark && (
                  <div className="flex items-start gap-2.5 text-sm">
                    <ShieldCheck className="mt-0.5 size-4 shrink-0 text-text-secondary" />
                    <div className="min-w-0">
                      <p className="font-mono text-[11px] uppercase tracking-widest text-text-secondary">
                        审核备注
                      </p>
                      <p className="mt-0.5 leading-relaxed text-text-primary">
                        {link.review_remark}
                      </p>
                    </div>
                  </div>
                )}
              </div>
            </motion.section>
          )}
        </div>

        <motion.section
          {...enter(reduced, 0.2, {
            initial: { opacity: 0, y: 10 },
            animate: { opacity: 1, y: 0 },
            transition: { duration: 0.4, ease: 'easeOut' },
          })}
          className={`${inkCard} h-fit`}
        >
          <CardHead title="基本信息" meta="INFO" />
          <InfoRow icon={<Mail className="size-4" />} label="站长邮箱">
            {link.email ? (
              <span className="flex items-center gap-1">
                {link.email}
                <CopyButton text={link.email} label="复制站长邮箱" />
              </span>
            ) : (
              <span className="text-text-secondary">—</span>
            )}
          </InfoRow>
          <InfoRow icon={<MapPin className="size-4" />} label="所在位置">
            {link.group_f_key?.name ?? '未分组'}
          </InfoRow>
          <InfoRow icon={<Globe className="size-4" />} label="审核状态">
            <InkBadge tone={s.tone}>{s.label}</InkBadge>
          </InfoRow>
          <InfoRow icon={<RefreshCcw className="size-4" />} label="更新时间">
            <span className="font-mono tabular-nums">
              {new Date(link.updated_at).toLocaleString('zh-CN')}
            </span>
          </InfoRow>
        </motion.section>
      </div>

      {/* 删除确认弹窗 */}
      <Dialog
        open={deleteOpen}
        onOpenChange={(open) => {
          if (!open) setDeleteOpen(false)
        }}
      >
        <DialogContent>
          <DialogHeader>
            <DialogTitle>删除友链</DialogTitle>
            <DialogDescription>
              确定要删除友链「{link.name}」吗？此操作不可撤销。
            </DialogDescription>
          </DialogHeader>
          <DialogFooter>
            <Button
              variant="outline"
              onClick={() => setDeleteOpen(false)}
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
