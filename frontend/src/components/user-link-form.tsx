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

import { useEffect, useMemo, useState } from 'react'
import { MapPin, Palette, Save } from 'lucide-react'
import type { FormEvent } from 'react'
import type { ApplyLinkRequest, LinkFriend } from '@/api/types'
import { useAuth } from '@/hooks/use-auth'
import { usePublicColors, usePublicGroups } from '@/hooks/use-links'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { Select } from '@/components/ui/select'
import { Textarea } from '@/components/ui/textarea'

/** 用户友链表单状态 */
interface UserLinkFormState {
  siteName: string
  siteUrl: string
  siteLogo: string
  siteRss: string
  webmasterEmail: string
  groupId: string
  colorId: string
  siteDescription: string
  applyRemark: string
}

/**
 * 用户友链信息表单：访客申请 / 用户编辑共用。
 *
 * 仅含站点基础信息字段（名称/地址/Logo/邮箱/描述），分组/颜色/级别/排序
 * 等管理员专属字段不在此开放。`initial` 提供时按友链详情预填（编辑），否则为空表单（申请）。
 *
 * `mode` 区分两种语义：
 * - `apply`（默认）：开放「展示位置/展示颜色」选择与「申请备注」输入（备注仅博主审核可见）。
 * - `edit`：展示位置/颜色改为只读展示（不可编辑），且不提供申请备注输入；
 *   提交载荷不携带位置/颜色/备注字段，仅更新站点基础信息。
 *
 * 提交时组装 ApplyLinkRequest 交回调用方（edit 模式为其子集 UpdateUserLinkRequest 所需字段）。
 */
export function UserLinkForm({
  initial,
  submitting,
  submitLabel = '保存',
  mode = 'apply',
  onSubmit,
}: {
  initial?: LinkFriend | null
  submitting: boolean
  submitLabel?: string
  mode?: 'apply' | 'edit'
  onSubmit: (req: ApplyLinkRequest) => void
}) {
  const { user, isAuthenticated } = useAuth()
  const lockedEmail = isAuthenticated ? (user?.email ?? '') : ''
  const { data: groups } = usePublicGroups()
  const { data: publicColors } = usePublicColors()

  // 颜色下拉兜底：服务端已按高级配色开关过滤；编辑场景下若当前选中颜色不在
  // 可见列表（如普通模式下既有高级色），追加进下拉选项，避免选中态丢失被误清。
  const colors = useMemo(() => {
    const list = publicColors ?? []
    const selected = initial?.color_f_key
    if (selected && !list.some((c) => c.id === selected.id)) {
      return [...list, selected]
    }
    return list
  }, [publicColors, initial?.color_f_key])

  const [form, setForm] = useState<UserLinkFormState>({
    siteName: '',
    siteUrl: '',
    siteLogo: '',
    siteRss: '',
    webmasterEmail: '',
    groupId: '',
    colorId: '',
    siteDescription: '',
    applyRemark: '',
  })

  // 详情加载完成后预填表单（编辑场景）
  useEffect(() => {
    if (initial) {
      setForm({
        siteName: initial.name,
        siteUrl: initial.url,
        siteLogo: initial.avatar ?? '',
        siteRss: initial.rss ?? '',
        webmasterEmail: lockedEmail || (initial.email ?? ''),
        groupId: initial.group_id?.toString() ?? '',
        colorId: initial.color_id?.toString() ?? '',
        siteDescription: initial.description ?? '',
        applyRemark: initial.apply_remark ?? '',
      })
    }
  }, [initial, lockedEmail])

  const handleSubmit = (e: FormEvent) => {
    e.preventDefault()
    onSubmit({
      link_name: form.siteName.trim(),
      link_url: form.siteUrl.trim(),
      link_avatar: form.siteLogo.trim() || undefined,
      link_rss: form.siteRss.trim() || undefined,
      link_email: (lockedEmail || form.webmasterEmail).trim(),
      // edit 模式：展示位置/颜色不可编辑、无申请备注，载荷仅站点基础信息
      ...(mode === 'edit'
        ? {}
        : {
            link_group_id: form.groupId ? BigInt(form.groupId) : undefined,
            link_color_id: form.colorId ? BigInt(form.colorId) : undefined,
            link_apply_remark: form.applyRemark.trim() || undefined,
          }),
      link_desc: form.siteDescription.trim() || undefined,
    })
  }

  return (
    <form onSubmit={handleSubmit} className="space-y-5">
      {/* 行格字段：名称|地址 · Logo|RSS · 位置|颜色 · 邮箱全宽收尾（连续行格，不设分区标题） */}
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
          <Label htmlFor="siteRss">订阅地址</Label>
          <Input
            id="siteRss"
            type="url"
            placeholder="https://example.com/atom.xml"
            value={form.siteRss}
            onChange={(e) => setForm({ ...form, siteRss: e.target.value })}
          />
        </div>
        <div className="space-y-2">
          <Label htmlFor="groupId" className="flex items-center gap-1.5">
            <MapPin className="size-3.5 text-text-secondary" />
            展示位置
          </Label>
          {mode === 'edit' ? (
            <div className="rounded-md border border-border/70 bg-muted/40 px-3 py-2 text-sm text-text-secondary">
              {initial?.group_f_key?.name ??
                (initial?.group_id ? `位置 #${initial.group_id}` : '未分组')}
            </div>
          ) : (
            <Select
              id="groupId"
              value={form.groupId}
              onChange={(e) => setForm({ ...form, groupId: e.target.value })}
            >
              <option value="">请选择位置</option>
              {groups?.map((g) => (
                <option key={g.id.toString()} value={g.id.toString()}>
                  {g.name}
                </option>
              ))}
            </Select>
          )}
        </div>
        <div className="space-y-2">
          <Label htmlFor="colorId" className="flex items-center gap-1.5">
            <Palette className="size-3.5 text-text-secondary" />
            展示颜色
          </Label>
          {mode === 'edit' ? (
            <div className="rounded-md border border-border/70 bg-muted/40 px-3 py-2 text-sm text-text-secondary">
              {initial?.color_f_key?.name ??
                (initial?.color_id ? `颜色 #${initial.color_id}` : '默认颜色')}
            </div>
          ) : (
            <Select
              id="colorId"
              value={form.colorId}
              onChange={(e) => setForm({ ...form, colorId: e.target.value })}
            >
              <option value="">请选择颜色</option>
              {colors?.map((c) => (
                <option key={c.id.toString()} value={c.id.toString()}>
                  {c.name}
                </option>
              ))}
            </Select>
          )}
        </div>
        <div className="space-y-2 md:col-span-2">
          <Label htmlFor="webmasterEmail">
            站长邮箱 <span className="text-destructive">*</span>
          </Label>
          <Input
            id="webmasterEmail"
            type="email"
            placeholder="admin@example.com"
            required
            readOnly={isAuthenticated}
            value={isAuthenticated ? lockedEmail : form.webmasterEmail}
            onChange={(e) =>
              setForm({ ...form, webmasterEmail: e.target.value })
            }
            className={
              isAuthenticated ? 'cursor-not-allowed bg-muted/30 opacity-70' : ''
            }
          />
          {isAuthenticated && (
            <p className="text-[11px] text-text-secondary">
              已登录，邮箱自动填入且不可修改
            </p>
          )}
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

      {/* 申请备注（仅博主审核可见；编辑基础信息时不提供） */}
      {mode === 'edit' ? null : (
        <div className="space-y-2">
          <Label htmlFor="applyRemark">申请备注</Label>
          <Input
            id="applyRemark"
            placeholder="选填，例如申请来源说明"
            value={form.applyRemark}
            onChange={(e) => setForm({ ...form, applyRemark: e.target.value })}
          />
          <p className="text-[11px] text-text-secondary">
            备注仅供博主审核时查看，不会公开展示
          </p>
        </div>
      )}

      {/* 提交按钮 */}
      <div className="flex justify-end border-t border-border/60 pt-5">
        <Button type="submit" className="cursor-pointer" disabled={submitting}>
          <Save className="mr-2 size-4" />
          {submitting ? '提交中…' : submitLabel}
        </Button>
      </div>
    </form>
  )
}
