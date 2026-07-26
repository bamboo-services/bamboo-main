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
import { Link, createFileRoute } from '@tanstack/react-router'
import { motion } from 'motion/react'
import {
  ArrowLeft,
  CalendarDays,
  Check,
  ChevronRight,
  ExternalLink,
  Globe,
  Inbox,
  Mail,
  MapPin,
  X,
} from 'lucide-react'
import { Button } from '@/components/ui/button'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { Badge } from '@/components/ui/badge'
import { mockLinks } from '@/data/mock/links'
import type { LinkItem } from '@/data/mock/links'
import { cn } from '@/lib/utils'

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
  url: string
  className?: string
}) {
  const [failed, setFailed] = useState(false)
  if (failed) {
    return (
      <div
        className={cn(
          'flex shrink-0 items-center justify-center rounded-lg bg-amber-500/10 font-semibold text-amber-600',
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

/** 待审核徽章 */
function PendingBadge() {
  return (
    <Badge className="bg-amber-500/15 text-amber-700 hover:bg-amber-500/15">
      待审核
    </Badge>
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
  link: LinkItem
  compact: boolean
  active: boolean
  onSelect: (id: number) => void
}) {
  return (
    <motion.button
      layout
      type="button"
      onClick={() => onSelect(link.id)}
      transition={morphSpring}
      className={cn(
        'group w-full cursor-pointer overflow-hidden rounded-xl border bg-card text-left transition-colors duration-200',
        active
          ? 'border-amber-500/60 bg-amber-500/10 shadow-sm'
          : 'border-border/70 hover:border-amber-500/40 hover:shadow-sm',
        compact ? 'p-2.5' : 'p-4',
      )}
    >
      {/* 常驻头部：头像 + 名称 + 元信息 */}
      <div className="flex items-center gap-3">
        <SiteAvatar
          name={link.siteName}
          url={link.siteLogo}
          className={cn(
            'transition-all duration-300',
            compact ? 'size-8 text-xs' : 'size-11 text-sm',
          )}
        />
        <div className="min-w-0 flex-1">
          <div className="flex items-center justify-between gap-2">
            <h3
              className={cn(
                'truncate font-semibold transition-all duration-300',
                compact ? 'text-sm' : 'text-base',
              )}
            >
              {link.siteName}
            </h3>
            {compact ? (
              active && (
                <span
                  className="size-2 shrink-0 rounded-full bg-amber-500"
                  aria-hidden="true"
                />
              )
            ) : (
              <PendingBadge />
            )}
          </div>
          <div className="mt-0.5 flex items-center justify-between gap-2 text-xs text-muted-foreground">
            <span className="truncate">{link.locationName}</span>
            <span className="shrink-0 tabular-nums">{link.createdAt}</span>
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
          <p className="mt-3 line-clamp-2 text-sm leading-relaxed text-muted-foreground">
            {link.siteDescription}
          </p>
          <div className="mt-3 flex items-center justify-between gap-2 border-t border-border/40 pt-3 text-xs">
            <span className="flex min-w-0 items-center gap-1.5 text-muted-foreground">
              <Mail className="size-3.5 shrink-0" />
              <span className="truncate">{link.webmasterEmail}</span>
            </span>
            <span className="flex shrink-0 items-center gap-0.5 font-medium text-amber-600">
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
function VerifyDetail({ link }: { link: LinkItem }) {
  return (
    <div className="space-y-4">
      {/* 申请横幅：琥珀色呼应待审核状态 */}
      <Card className="relative overflow-hidden">
        <span
          className="absolute inset-y-0 left-0 z-10 w-1.5 bg-amber-500"
          aria-hidden="true"
        />
        <div
          className="pointer-events-none absolute inset-0"
          style={{
            background:
              'linear-gradient(120deg, rgb(245 158 11 / 0.14) 0%, rgb(245 158 11 / 0.05) 40%, transparent 70%)',
          }}
          aria-hidden="true"
        />
        <CardContent className="relative flex flex-wrap items-center gap-5 p-6 pl-9">
          <SiteAvatar
            name={link.siteName}
            url={link.siteLogo}
            className="size-20 rounded-2xl text-3xl"
          />
          <div className="min-w-0 flex-1">
            <div className="flex flex-wrap items-center gap-2.5">
              <h2 className="text-2xl font-bold tracking-tight">
                {link.siteName}
              </h2>
              <PendingBadge />
            </div>
            <a
              href={link.siteUrl}
              target="_blank"
              rel="noopener noreferrer"
              className="mt-2 inline-flex items-center gap-1.5 text-sm text-muted-foreground transition-colors hover:text-primary"
            >
              <Globe className="size-4" />
              {link.siteUrl}
              <ExternalLink className="size-3.5" />
            </a>
          </div>
        </CardContent>
      </Card>

      {/* 描述 + 申请信息 */}
      <div className="grid gap-4 xl:grid-cols-3">
        <Card className="xl:col-span-2">
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
            <CardTitle className="text-base">申请信息</CardTitle>
          </CardHeader>
          <CardContent className="space-y-3 pt-1">
            <div className="flex items-center gap-2.5 text-sm">
              <Mail className="size-4 shrink-0 text-muted-foreground" />
              <span className="truncate">{link.webmasterEmail}</span>
            </div>
            <div className="flex items-center gap-2.5 text-sm">
              <MapPin className="size-4 shrink-0 text-muted-foreground" />
              <span>{link.locationName}</span>
            </div>
            <div className="flex items-center gap-2.5 text-sm">
              <CalendarDays className="size-4 shrink-0 text-muted-foreground" />
              <span className="tabular-nums">{link.createdAt}</span>
            </div>
          </CardContent>
        </Card>
      </div>

      {/* 审核操作 */}
      <Card>
        <CardContent className="flex flex-wrap items-center justify-between gap-3 p-4">
          <p className="text-sm text-muted-foreground">
            审核通过后该站点将展示在友链页面
          </p>
          <div className="flex gap-2">
            <Button
              variant="outline"
              className="cursor-pointer text-destructive hover:bg-destructive/10 hover:text-destructive"
            >
              <X className="mr-2 size-4" />
              拒绝
            </Button>
            <Button className="cursor-pointer bg-green-600 hover:bg-green-700">
              <Check className="mr-2 size-4" />
              通过
            </Button>
          </div>
        </CardContent>
      </Card>
    </div>
  )
}

function LinkVerifyPage() {
  const pendingLinks = mockLinks.filter((link) => link.status === 'pending')
  const [selectedId, setSelectedId] = useState<number | null>(null)
  const selected = pendingLinks.find((l) => l.id === selectedId) ?? null
  const hasSelection = selected !== null

  /** 点击已选中项则取消选中 */
  const toggle = (id: number) =>
    setSelectedId((prev) => (prev === id ? null : id))

  return (
    <div className="space-y-5">
      {/* 页头 */}
      <div className="flex items-center gap-4">
        <Link to="/admin/link">
          <Button variant="ghost" size="icon" className="cursor-pointer">
            <ArrowLeft className="size-4" />
          </Button>
        </Link>
        <div>
          <h1 className="text-3xl font-bold tracking-tight">友链审核</h1>
          <p className="mt-1 flex items-center gap-2 text-muted-foreground">
            待审核的友链申请
            {pendingLinks.length > 0 && (
              <Badge variant="secondary">{pendingLinks.length}</Badge>
            )}
          </p>
        </div>
      </div>

      {pendingLinks.length === 0 ? (
        /* 空状态 */
        <Card className="border-dashed">
          <CardContent className="flex flex-col items-center justify-center gap-3 py-20 text-center">
            <div className="flex size-14 items-center justify-center rounded-full bg-muted">
              <Inbox className="size-6 text-muted-foreground" />
            </div>
            <div>
              <p className="font-medium">暂无待审核的友链申请</p>
              <p className="mt-1 text-sm text-muted-foreground">
                新的友链申请会出现在这里
              </p>
            </div>
            <Link to="/admin/link">
              <Button variant="outline" size="sm" className="cursor-pointer">
                返回列表
              </Button>
            </Link>
          </CardContent>
        </Card>
      ) : (
        /*
         * 主从布局：列表始终完整存在。
         * 未选中 → 列表占满（多列网格，富密度）；
         * 选中后 → 列表原地收紧为窄列（2），详情面板从右侧弹出（8）。
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
                key={link.id}
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
              <VerifyDetail link={selected} />
            </motion.div>
          )}
        </div>
      )}
    </div>
  )
}
