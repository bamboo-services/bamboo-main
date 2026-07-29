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
import { Save } from 'lucide-react'
import type { FormEvent } from 'react'
import type { ApplyLinkRequest, LinkFriend } from '@/api/types'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { Textarea } from '@/components/ui/textarea'

/** 用户友链表单状态 */
interface UserLinkFormState {
  siteName: string
  siteUrl: string
  siteLogo: string
  webmasterEmail: string
  siteDescription: string
  applyRemark: string
}

/**
 * 用户友链信息表单：访客申请 / 用户编辑共用。
 *
 * 仅含站点基础信息字段（名称/地址/Logo/邮箱/描述/备注），分组/颜色/级别/排序
 * 等管理员专属字段不在此开放。`initial` 提供时按友链详情预填（编辑），否则为空表单（申请）。
 * 提交时组装 ApplyLinkRequest（字段与 UpdateUserLinkRequest 同构且为其子集）交回调用方。
 */
export function UserLinkForm({
  initial,
  submitting,
  submitLabel = '保存',
  onSubmit,
}: {
  initial?: LinkFriend | null
  submitting: boolean
  submitLabel?: string
  onSubmit: (req: ApplyLinkRequest) => void
}) {
  const [form, setForm] = useState<UserLinkFormState>({
    siteName: '',
    siteUrl: '',
    siteLogo: '',
    webmasterEmail: '',
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
        webmasterEmail: initial.email ?? '',
        siteDescription: initial.description ?? '',
        applyRemark: initial.apply_remark ?? '',
      })
    }
  }, [initial])

  const handleSubmit = (e: FormEvent) => {
    e.preventDefault()
    onSubmit({
      link_name: form.siteName.trim(),
      link_url: form.siteUrl.trim(),
      link_avatar: form.siteLogo.trim() || undefined,
      link_email: form.webmasterEmail.trim(),
      link_desc: form.siteDescription.trim() || undefined,
      link_apply_remark: form.applyRemark.trim() || undefined,
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
          <Label htmlFor="webmasterEmail">
            站长邮箱 <span className="text-destructive">*</span>
          </Label>
          <Input
            id="webmasterEmail"
            type="email"
            placeholder="admin@example.com"
            required
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

      {/* 申请备注 */}
      <div className="space-y-2">
        <Label htmlFor="applyRemark">申请备注</Label>
        <Input
          id="applyRemark"
          placeholder="选填，例如申请来源说明"
          value={form.applyRemark}
          onChange={(e) => setForm({ ...form, applyRemark: e.target.value })}
        />
      </div>

      {/* 提交按钮 */}
      <div className="flex justify-end border-t border-border/60 pt-6">
        <Button type="submit" className="cursor-pointer" disabled={submitting}>
          <Save className="mr-2 size-4" />
          {submitting ? '提交中…' : submitLabel}
        </Button>
      </div>
    </form>
  )
}
