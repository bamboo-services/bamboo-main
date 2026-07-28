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
import { Link } from '@tanstack/react-router'
import { Check, Save } from 'lucide-react'
import type { FormEvent } from 'react'
import type {
  CreateLinkRequest,
  LinkColor,
  LinkFriend,
  SnowflakeID,
} from '@/api/types'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { Skeleton } from '@/components/ui/skeleton'
import { Textarea } from '@/components/ui/textarea'
import { InkPill } from '@/components/ink-wash'
import { cn } from '@/lib/utils'
import { useAllGroups } from '@/hooks/use-groups'
import { useAllColors } from '@/hooks/use-colors'

/** 表单状态（group/color 为 bigint 雪花 ID，未选为 null） */
interface LinkFormState {
  siteName: string
  siteUrl: string
  siteLogo: string
  webmasterEmail: string
  siteDescription: string
  groupId: SnowflakeID | null
  colorId: SnowflakeID | null
  order: string
  applyRemark: string
  level: number
}

/** 友链级别选项（与后端 pkg/constants LinkLevel 枚举对齐） */
const LEVEL_OPTIONS = [
  { value: 0, label: '一般' },
  { value: 1, label: '好友' },
  { value: 2, label: '高级' },
  { value: 3, label: '广告' },
] as const

/** 色块背景：type=1（炫彩）且有副色时渲染渐变，否则取主色 */
function colorBackground(color: LinkColor): string {
  if (color.type === 1 && color.sub_color) {
    return `linear-gradient(135deg, ${color.primary_color ?? 'var(--leaf-deep)'}, ${color.sub_color})`
  }
  return color.primary_color ?? 'var(--leaf-deep)'
}

/**
 * 友链信息表单：添加 / 编辑共用。
 * `initial` 提供时按友链详情预填（编辑），否则为空表单（添加）。
 * 提交时组装 CreateLinkRequest（与 UpdateLinkRequest 字段同构）交回调用方。
 */
export function LinkForm({
  initial,
  submitting,
  onSubmit,
}: {
  initial?: LinkFriend | null
  submitting: boolean
  onSubmit: (req: CreateLinkRequest) => void
}) {
  const groupsQuery = useAllGroups()
  const colorsQuery = useAllColors()

  const groups = groupsQuery.data ?? []
  const colors = colorsQuery.data ?? []

  const [form, setForm] = useState<LinkFormState>({
    siteName: '',
    siteUrl: '',
    siteLogo: '',
    webmasterEmail: '',
    siteDescription: '',
    groupId: null,
    colorId: null,
    order: '',
    applyRemark: '',
    level: 0,
  })

  // 详情加载完成后预填表单（编辑场景）
  useEffect(() => {
    if (initial) {
      setForm({
        siteName: initial.name,
        siteUrl: initial.url,
        siteLogo: initial.avatar ?? '',
        webmasterEmail: initial.email ?? '',
        siteDescription: initial.description ?? '',
        groupId: initial.group_id,
        colorId: initial.color_id,
        order: String(initial.sort_order),
        applyRemark: initial.apply_remark ?? '',
        level: initial.level,
      })
    }
  }, [initial])

  const handleSubmit = (e: FormEvent) => {
    e.preventDefault()
    const order = Number.parseInt(form.order, 10)
    onSubmit({
      link_name: form.siteName.trim(),
      link_url: form.siteUrl.trim(),
      link_avatar: form.siteLogo.trim() || undefined,
      link_email: form.webmasterEmail.trim() || undefined,
      link_desc: form.siteDescription.trim() || undefined,
      link_group_id: form.groupId ?? undefined,
      link_color_id: form.colorId ?? undefined,
      link_order: Number.isNaN(order) ? undefined : order,
      link_apply_remark: form.applyRemark.trim() || undefined,
      link_level: form.level,
    })
  }

  return (
    <form onSubmit={handleSubmit} className="space-y-6">
      {/* 基本信息 */}
      <div className="grid gap-4 md:grid-cols-2">
        <div className="space-y-2">
          <Label htmlFor="siteName">
            站点名称 <span className="text-destructive">*</span>
          </Label>
          <Input
            id="siteName"
            placeholder="请输入站点名称"
            required
            value={form.siteName}
            onChange={(e) => setForm({ ...form, siteName: e.target.value })}
          />
        </div>
        <div className="space-y-2">
          <Label htmlFor="siteUrl">
            站点地址 <span className="text-destructive">*</span>
          </Label>
          <Input
            id="siteUrl"
            type="url"
            placeholder="https://example.com"
            required
            value={form.siteUrl}
            onChange={(e) => setForm({ ...form, siteUrl: e.target.value })}
          />
        </div>
        <div className="space-y-2">
          <Label htmlFor="siteLogo">站点 Logo</Label>
          <Input
            id="siteLogo"
            placeholder="https://example.com/logo.png"
            value={form.siteLogo}
            onChange={(e) => setForm({ ...form, siteLogo: e.target.value })}
          />
        </div>
        <div className="space-y-2">
          <Label htmlFor="webmasterEmail">站长邮箱</Label>
          <Input
            id="webmasterEmail"
            type="email"
            placeholder="admin@example.com"
            value={form.webmasterEmail}
            onChange={(e) =>
              setForm({ ...form, webmasterEmail: e.target.value })
            }
          />
        </div>
      </div>

      {/* 站点描述 */}
      <div className="space-y-2">
        <Label htmlFor="siteDescription">站点描述</Label>
        <Textarea
          id="siteDescription"
          className="min-h-[110px]"
          placeholder="介绍一下这个站点吧…"
          value={form.siteDescription}
          onChange={(e) =>
            setForm({ ...form, siteDescription: e.target.value })
          }
        />
      </div>

      {/* 排序 + 申请备注 */}
      <div className="grid gap-4 md:grid-cols-2">
        <div className="space-y-2">
          <Label htmlFor="linkOrder">排序</Label>
          <Input
            id="linkOrder"
            type="number"
            step={1}
            placeholder="数值越小越靠前，默认 0"
            value={form.order}
            onChange={(e) => setForm({ ...form, order: e.target.value })}
          />
        </div>
        <div className="space-y-2">
          <Label htmlFor="applyRemark">申请备注</Label>
          <Input
            id="applyRemark"
            placeholder="选填，例如申请来源说明"
            value={form.applyRemark}
            onChange={(e) => setForm({ ...form, applyRemark: e.target.value })}
          />
        </div>
      </div>

      {/* 友链级别：药丸按钮 */}
      <div className="space-y-2">
        <div className="flex items-baseline justify-between">
          <Label>友链级别</Label>
          <span className="font-mono text-xs text-text-secondary">
            决定前台展示样式
          </span>
        </div>
        <div className="flex flex-wrap gap-2">
          {LEVEL_OPTIONS.map((opt) => (
            <InkPill
              key={opt.value}
              active={form.level === opt.value}
              onClick={() => setForm({ ...form, level: opt.value })}
            >
              {opt.label}
            </InkPill>
          ))}
        </div>
      </div>

      {/* 位置分类：药丸按钮 */}
      <div className="space-y-2">
        <div className="flex items-baseline justify-between">
          <Label>位置分类</Label>
          <span className="font-mono text-xs text-text-secondary">
            再次点击可取消选择
          </span>
        </div>
        {groupsQuery.isLoading ? (
          <div className="flex gap-2">
            {Array.from({ length: 4 }).map((_, i) => (
              <Skeleton key={i} className="h-8 w-20 rounded-full" />
            ))}
          </div>
        ) : (
          <div className="flex flex-wrap gap-2">
            <InkPill
              active={form.groupId === null}
              onClick={() => setForm({ ...form, groupId: null })}
            >
              未分组
            </InkPill>
            {groups.map((group) => (
              <InkPill
                key={group.id.toString()}
                title={group.description ?? undefined}
                active={form.groupId === group.id}
                onClick={() =>
                  setForm({
                    ...form,
                    groupId: form.groupId === group.id ? null : group.id,
                  })
                }
              >
                {group.name}
              </InkPill>
            ))}
            {groups.length === 0 && (
              <span className="self-center text-xs text-text-secondary">
                暂无分组，可先前往分组管理创建
              </span>
            )}
          </div>
        )}
      </div>

      {/* 颜色分类：可视化色块 */}
      <div className="space-y-2">
        <div className="flex items-baseline justify-between">
          <Label>颜色分类</Label>
          <span className="font-mono text-xs text-text-secondary">
            再次点击可恢复默认
          </span>
        </div>
        {colorsQuery.isLoading ? (
          <div className="flex gap-3">
            {Array.from({ length: 6 }).map((_, i) => (
              <Skeleton key={i} className="size-9 rounded-full" />
            ))}
          </div>
        ) : (
          <div className="flex flex-wrap items-center gap-3">
            <button
              type="button"
              title="默认颜色"
              onClick={() => setForm({ ...form, colorId: null })}
              className={cn(
                'flex size-9 cursor-pointer items-center justify-center rounded-full border-2 border-dashed border-border bg-muted/40 transition-all duration-200',
                form.colorId === null
                  ? 'scale-110 ring-2 ring-leaf-deep ring-offset-2 ring-offset-background'
                  : 'hover:scale-105',
              )}
            >
              {form.colorId === null && (
                <Check
                  className="size-4 text-text-secondary"
                  strokeWidth={3}
                />
              )}
            </button>
            {colors.map((color) => {
              const selected = form.colorId === color.id
              return (
                <button
                  key={color.id.toString()}
                  type="button"
                  title={color.name}
                  onClick={() =>
                    setForm({
                      ...form,
                      colorId: selected ? null : color.id,
                    })
                  }
                  className={cn(
                    'flex size-9 cursor-pointer items-center justify-center rounded-full ring-offset-2 ring-offset-background transition-all duration-200',
                    selected
                      ? 'scale-110 ring-2 ring-leaf-deep'
                      : 'hover:scale-105',
                  )}
                  style={{ background: colorBackground(color) }}
                >
                  {selected && (
                    <Check className="size-4 text-card" strokeWidth={3} />
                  )}
                </button>
              )
            })}
            {colors.length === 0 && (
              <span className="text-xs text-text-secondary">
                暂无颜色，将使用默认颜色
              </span>
            )}
          </div>
        )}
      </div>

      {/* 操作按钮 */}
      <div className="flex justify-end gap-3 border-t border-border/60 pt-6">
        <Link to="/admin/link">
          <Button variant="outline" className="cursor-pointer">
            取消
          </Button>
        </Link>
        <Button type="submit" className="cursor-pointer" disabled={submitting}>
          <Save className="mr-2 size-4" />
          {submitting ? '保存中…' : '保存'}
        </Button>
      </div>
    </form>
  )
}
