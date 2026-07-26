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

import { useState } from 'react'
import { Link, createFileRoute, useNavigate } from '@tanstack/react-router'
import {
  ArrowLeft,
  CalendarDays,
  Check,
  Clock3,
  Copy,
  ExternalLink,
  Globe,
  Mail,
  MapPin,
  Pencil,
  RefreshCcw,
  Trash2,
  Unlink,
} from 'lucide-react'
import { Button } from '@/components/ui/button'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { Badge } from '@/components/ui/badge'
import { Avatar, AvatarFallback, AvatarImage } from '@/components/ui/avatar'
import { mockColors, mockLinks } from '@/data/mock/links'
import type { LinkItem } from '@/data/mock/links'
import { cn } from '@/lib/utils'

export const Route = createFileRoute('/_admin/admin/link/$id')({
  component: LinkDetailPage,
})

function colorHex(id: number): string {
  return mockColors.find((c) => c.id === id)?.color ?? '#6366f1'
}

function colorName(id: number): string {
  return mockColors.find((c) => c.id === id)?.name ?? '默认'
}

/** 计算收录天数 */
function daysSince(dateStr: string): number {
  const diff = Date.now() - new Date(dateStr).getTime()
  return Math.max(0, Math.floor(diff / 86_400_000))
}

/** 复制按钮：成功后短暂显示对勾反馈 */
function CopyButton({ text, label }: { text: string; label: string }) {
  const [copied, setCopied] = useState(false)
  const onCopy = () => {
    void navigator.clipboard?.writeText(text).then(() => {
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
      className="inline-flex size-6 shrink-0 cursor-pointer items-center justify-center rounded-md text-muted-foreground transition-colors duration-150 hover:bg-muted hover:text-foreground"
    >
      {copied ? (
        <Check className="size-3.5 text-green-600" />
      ) : (
        <Copy className="size-3.5" />
      )}
    </button>
  )
}

/** 可访问状态徽章：可访问时带脉冲呼吸点 */
function StatusPill({ link }: { link: LinkItem }) {
  return link.ableConnect ? (
    <Badge className="bg-green-500/15 text-green-700 hover:bg-green-500/15">
      <span className="relative mr-1.5 flex size-2">
        <span className="absolute inline-flex size-full animate-ping rounded-full bg-green-500 opacity-60 motion-reduce:animate-none" />
        <span className="relative inline-flex size-2 rounded-full bg-green-500" />
      </span>
      可访问
    </Badge>
  ) : (
    <Badge variant="destructive">
      <span className="mr-1.5 size-2 rounded-full bg-current opacity-70" />
      不可访问
    </Badge>
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
    <div className="flex items-center justify-between gap-3 border-b border-border/40 py-3 last:border-0">
      <span className="flex shrink-0 items-center gap-2 text-sm text-muted-foreground">
        {icon}
        {label}
      </span>
      <span className="text-right text-sm font-medium">{children}</span>
    </div>
  )
}

/** 友链预览：1:1 还原公开页的展示卡片 */
function FriendPreview({ link }: { link: LinkItem }) {
  const accent = colorHex(link.color)
  return (
    <a
      href={link.siteUrl}
      target="_blank"
      rel="noopener noreferrer"
      className="group relative flex gap-3 overflow-hidden rounded-xl border border-leaf-muted/40 bg-card/80 p-4 shadow-sm transition-all duration-200 hover:-translate-y-0.5 hover:border-primary/40 hover:shadow-md"
    >
      <span
        className="absolute inset-y-0 left-0 w-1"
        style={{ backgroundColor: accent }}
        aria-hidden="true"
      />
      <Avatar className="size-12 shrink-0 rounded-full">
        <AvatarImage src={link.siteLogo} alt={link.siteName} />
        <AvatarFallback className="bg-leaf-light/40 text-text-primary">
          {link.siteName.slice(0, 1)}
        </AvatarFallback>
      </Avatar>
      <div className="min-w-0 flex-1">
        <div className="flex items-center justify-between gap-2">
          <h4 className="truncate font-medium text-text-primary group-hover:text-primary">
            {link.siteName}
          </h4>
          <ExternalLink className="size-4 shrink-0 text-text-secondary opacity-0 transition-opacity duration-200 group-hover:opacity-100" />
        </div>
        <p className="mt-1 line-clamp-2 text-sm text-text-secondary">
          {link.siteDescription || '这个站点很神秘，没有留下描述。'}
        </p>
      </div>
    </a>
  )
}

function LinkDetailPage() {
  const { id } = Route.useParams()
  const navigate = useNavigate()
  const link = mockLinks.find((l) => l.id === Number(id))

  if (!link) {
    return (
      <div className="flex min-h-[60vh] flex-col items-center justify-center gap-4 text-center">
        <div className="flex size-16 items-center justify-center rounded-full bg-muted">
          <Unlink className="size-7 text-muted-foreground" />
        </div>
        <div>
          <h2 className="text-xl font-semibold">友链不存在</h2>
          <p className="mt-1 text-sm text-muted-foreground">
            该友链可能已被删除或链接有误
          </p>
        </div>
        <Link to="/admin/link">
          <Button className="cursor-pointer">返回列表</Button>
        </Link>
      </div>
    )
  }

  const accent = colorHex(link.color)

  return (
    <div className="space-y-5">
      {/* 返回：轻量文字链接，箭头悬停微移 */}
      <Link
        to="/admin/link"
        className="group inline-flex cursor-pointer items-center gap-1.5 text-sm text-muted-foreground transition-colors duration-150 hover:text-foreground"
      >
        <ArrowLeft className="size-4 transition-transform duration-150 group-hover:-translate-x-0.5" />
        返回友链管理
      </Link>

      {/* 站点档案头：整卡沉浸渐变背景（absolute inset-0 天然贴满，无置顶问题） */}
      <Card className="relative overflow-hidden">
        {/* 多层站点色渐变：左上方主光晕 + 右下方副光晕 + 斜向底色 */}
        <div
          className="absolute inset-0"
          style={{
            background: `radial-gradient(ellipse 90% 130% at 18% 0%, ${accent}3d 0%, transparent 55%), radial-gradient(ellipse 70% 100% at 88% 100%, ${accent}26 0%, transparent 50%), linear-gradient(120deg, ${accent}1c 0%, ${accent}0d 50%, transparent 80%)`,
          }}
          aria-hidden="true"
        />
        {/* 柔光光斑，增加层次 */}
        <div
          className="pointer-events-none absolute -right-14 -top-24 size-64 rounded-full opacity-30 blur-3xl"
          style={{ backgroundColor: accent }}
          aria-hidden="true"
        />

        <CardContent className="relative flex flex-wrap items-center gap-5">
          <Avatar className="size-24 shrink-0 rounded-2xl shadow-xl ring-4 ring-white/90">
            <AvatarImage src={link.siteLogo} alt={link.siteName} />
            <AvatarFallback className="rounded-2xl bg-white/70 text-3xl font-semibold text-primary">
              {link.siteName.charAt(0)}
            </AvatarFallback>
          </Avatar>

          <div className="min-w-0 flex-1">
            <div className="flex flex-wrap items-center gap-2.5">
              <h1 className="truncate text-2xl font-bold tracking-tight lg:text-3xl">
                {link.siteName}
              </h1>
              <StatusPill link={link} />
              {link.hasAdv && (
                <Badge className="bg-amber-500/15 text-amber-700 hover:bg-amber-500/15">
                  含广告
                </Badge>
              )}
            </div>
            <div className="mt-2 flex flex-wrap items-center gap-x-4 gap-y-1.5 text-sm text-muted-foreground">
              <span className="flex min-w-0 items-center gap-1.5">
                <Globe className="size-4 shrink-0" />
                <a
                  href={link.siteUrl}
                  target="_blank"
                  rel="noopener noreferrer"
                  className="truncate transition-colors duration-150 hover:text-primary"
                >
                  {link.siteUrl}
                </a>
                <CopyButton text={link.siteUrl} label="复制站点地址" />
              </span>
              <span className="flex items-center gap-1.5">
                <MapPin className="size-4" />
                {link.locationName}
              </span>
              <span className="flex items-center gap-1.5">
                <span
                  className="size-2.5 rounded-full"
                  style={{ backgroundColor: accent }}
                  aria-hidden="true"
                />
                {colorName(link.color)}
              </span>
            </div>
          </div>

          <div className="flex shrink-0 flex-wrap gap-2">
            <a
              href={link.siteUrl}
              target="_blank"
              rel="noopener noreferrer"
              className="inline-flex h-9 cursor-pointer items-center justify-center gap-2 rounded-lg px-4 text-sm font-medium text-white transition duration-150 hover:brightness-95"
              style={{
                backgroundColor: accent,
                boxShadow: `0 4px 14px ${accent}4d`,
              }}
            >
              <ExternalLink className="size-4" />
              访问站点
            </a>
            <Button
              variant="outline"
              className="cursor-pointer bg-white/80 backdrop-blur-sm"
              onClick={() =>
                navigate({
                  to: '/admin/link/$id/edit',
                  params: { id: String(link.id) },
                })
              }
            >
              <Pencil className="mr-2 size-4" />
              编辑
            </Button>
            <Button
              variant="outline"
              className="cursor-pointer bg-white/80 text-destructive backdrop-blur-sm hover:bg-destructive/10 hover:text-destructive"
            >
              <Trash2 className="mr-2 size-4" />
              删除
            </Button>
          </div>
        </CardContent>
      </Card>

      {/* 内容网格：描述 + 预览（2 份）｜ 基本信息（1 份） */}
      <div className="grid gap-4 lg:grid-cols-3">
        <div className="space-y-4 lg:col-span-2">
          <Card>
            <CardHeader className="pb-3">
              <CardTitle className="text-base">站点描述</CardTitle>
            </CardHeader>
            <CardContent>
              <p className="leading-relaxed text-muted-foreground">
                {link.siteDescription || '暂无描述'}
              </p>
            </CardContent>
          </Card>

          <Card>
            <CardHeader className="pb-3">
              <CardTitle className="flex flex-wrap items-baseline gap-2 text-base">
                友链预览
                <span className="text-xs font-normal text-muted-foreground">
                  访客在公开页看到的卡片样式
                </span>
              </CardTitle>
            </CardHeader>
            <CardContent>
              <FriendPreview link={link} />
            </CardContent>
          </Card>
        </div>

        <Card className="h-fit">
          <CardHeader className="pb-3">
            <CardTitle className="text-base">基本信息</CardTitle>
          </CardHeader>
          <CardContent className="pt-1">
            <InfoRow icon={<Mail className="size-4" />} label="站长邮箱">
              <span className="flex items-center gap-1">
                {link.webmasterEmail}
                <CopyButton text={link.webmasterEmail} label="复制站长邮箱" />
              </span>
            </InfoRow>
            <InfoRow icon={<Clock3 className="size-4" />} label="收录时长">
              <span className="font-semibold text-primary tabular-nums">
                {daysSince(link.createdAt)} 天
              </span>
            </InfoRow>
            <InfoRow
              icon={<CalendarDays className="size-4" />}
              label="创建时间"
            >
              <span className="tabular-nums">{link.createdAt}</span>
            </InfoRow>
            <InfoRow icon={<RefreshCcw className="size-4" />} label="更新时间">
              <span className="tabular-nums">{link.updatedAt}</span>
            </InfoRow>
            <InfoRow icon={<Globe className="size-4" />} label="审核状态">
              <span
                className={cn(
                  link.status === 'approved' && 'text-green-600',
                  link.status === 'pending' && 'text-amber-600',
                  link.status === 'rejected' && 'text-destructive',
                )}
              >
                {link.status === 'approved'
                  ? '已通过'
                  : link.status === 'pending'
                    ? '待审核'
                    : '已拒绝'}
              </span>
            </InfoRow>
          </CardContent>
        </Card>
      </div>
    </div>
  )
}
