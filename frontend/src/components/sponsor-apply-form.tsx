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
import { Coins, Save } from 'lucide-react'
import type { FormEvent } from 'react'
import type { SponsorApplyRequest } from '@/api/types'
import { useAuth } from '@/hooks/use-auth'
import { usePublicChannels } from '@/hooks/use-sponsors'
import { Button } from '@/components/ui/button'
import { Checkbox } from '@/components/ui/checkbox'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { Select } from '@/components/ui/select'
import { Textarea } from '@/components/ui/textarea'

/** 赞助展示申请表单状态 */
interface SponsorApplyFormState {
  nickname: string
  amountYuan: string
  channelId: string
  message: string
  sponsorAt: string
  contactEmail: string
  redirectUrl: string
  isAnonymous: boolean
  applyRemark: string
}

/**
 * 赞助展示申请表单：访客申请 / 用户编辑共用。
 *
 * 仅含展示类基础字段（昵称/金额/渠道/留言/时间/邮箱/跳转/匿名/备注），
 * 金额以「元」录入、提交时转为「分」；is_hidden/sort_order 等管理员专属字段不在此开放。
 * 联系邮箱必填：已登录时锁定为账户邮箱只读（归属确认与审核结果通知用）。
 */
export function SponsorApplyForm({
  submitting,
  submitLabel = '递交申请',
  onSubmit,
}: {
  submitting: boolean
  submitLabel?: string
  onSubmit: (req: SponsorApplyRequest) => void
}) {
  const { user, isAuthenticated } = useAuth()
  const lockedEmail = isAuthenticated ? (user?.email ?? '') : ''
  const { data: channels } = usePublicChannels()

  const [form, setForm] = useState<SponsorApplyFormState>({
    nickname: '',
    amountYuan: '',
    channelId: '',
    message: '',
    sponsorAt: '',
    contactEmail: '',
    redirectUrl: '',
    isAnonymous: false,
    applyRemark: '',
  })

  const handleSubmit = (e: FormEvent) => {
    e.preventDefault()
    onSubmit({
      nickname: form.nickname.trim(),
      amount: Math.round(Number(form.amountYuan) * 100),
      channel_id: form.channelId ? BigInt(form.channelId) : undefined,
      message: form.message.trim() || undefined,
      sponsor_at: form.sponsorAt
        ? new Date(`${form.sponsorAt}T12:00:00`).toISOString()
        : undefined,
      email: (lockedEmail || form.contactEmail).trim(),
      redirect_url: form.redirectUrl.trim() || undefined,
      is_anonymous: form.isAnonymous || undefined,
      apply_remark: form.applyRemark.trim() || undefined,
    })
  }

  return (
    <form onSubmit={handleSubmit} className="space-y-5">
      {/* 行格字段：昵称|金额 · 渠道|时间 · 跳转全宽收尾（连续行格，不设分区标题） */}
      <div className="grid gap-4 md:grid-cols-2">
        <div className="space-y-2">
          <Label htmlFor="nickname">
            昵称 <span className="text-destructive">*</span>
          </Label>
          <Input
            id="nickname"
            placeholder="希望展示的署名"
            required
            value={form.nickname}
            onChange={(e) => setForm({ ...form, nickname: e.target.value })}
          />
        </div>
        <div className="space-y-2">
          <Label htmlFor="amountYuan" className="flex items-center gap-1.5">
            <Coins className="size-3.5 text-text-secondary" />
            赞助金额（元） <span className="text-destructive">*</span>
          </Label>
          <Input
            id="amountYuan"
            type="number"
            min="0.01"
            step="0.01"
            placeholder="6.66"
            required
            value={form.amountYuan}
            onChange={(e) => setForm({ ...form, amountYuan: e.target.value })}
          />
        </div>
        <div className="space-y-2">
          <Label htmlFor="channelId">赞助渠道</Label>
          <Select
            id="channelId"
            value={form.channelId}
            onChange={(e) => setForm({ ...form, channelId: e.target.value })}
          >
            <option value="">请选择渠道</option>
            {channels?.map((ch) => (
              <option key={ch.id.toString()} value={ch.id.toString()}>
                {ch.name}
              </option>
            ))}
          </Select>
        </div>
        <div className="space-y-2">
          <Label htmlFor="sponsorAt">赞助时间</Label>
          <Input
            id="sponsorAt"
            type="date"
            value={form.sponsorAt}
            onChange={(e) => setForm({ ...form, sponsorAt: e.target.value })}
          />
        </div>
        <div className="space-y-2 md:col-span-2">
          <Label htmlFor="contactEmail">
            联系邮箱 <span className="text-destructive">*</span>
          </Label>
          <Input
            id="contactEmail"
            type="email"
            placeholder="admin@example.com"
            required
            readOnly={isAuthenticated}
            value={isAuthenticated ? lockedEmail : form.contactEmail}
            onChange={(e) => setForm({ ...form, contactEmail: e.target.value })}
            className={
              isAuthenticated ? 'cursor-not-allowed bg-muted/30 opacity-70' : ''
            }
          />
          {isAuthenticated ? (
            <p className="text-[11px] text-text-secondary">
              已登录，邮箱自动填入且不可修改
            </p>
          ) : (
            <p className="text-[11px] text-text-secondary">
              用于确认归属与接收审核结果通知
            </p>
          )}
        </div>
      </div>

      {/* 留言 */}
      <div className="space-y-2">
        <Label htmlFor="message">留言</Label>
        <Textarea
          id="message"
          className="min-h-[90px]"
          placeholder="想对博主说的话…"
          value={form.message}
          onChange={(e) => setForm({ ...form, message: e.target.value })}
        />
      </div>

      {/* 跳转链接 */}
      <div className="space-y-2">
        <Label htmlFor="redirectUrl">跳转链接</Label>
        <Input
          id="redirectUrl"
          type="url"
          placeholder="https://example.com"
          value={form.redirectUrl}
          onChange={(e) => setForm({ ...form, redirectUrl: e.target.value })}
        />
      </div>

      {/* 匿名展示 */}
      <div className="flex items-start gap-2.5">
        <Checkbox
          id="isAnonymous"
          checked={form.isAnonymous}
          onCheckedChange={(checked) =>
            setForm({ ...form, isAnonymous: checked === true })
          }
        />
        <div className="grid gap-1">
          <Label htmlFor="isAnonymous" className="cursor-pointer">
            匿名展示
          </Label>
          <p className="text-[11px] text-text-secondary">
            勾选后前台将以「匿名用户」展示，不显示跳转链接
          </p>
        </div>
      </div>

      {/* 申请备注 */}
      <div className="space-y-2">
        <Label htmlFor="applyRemark">申请备注</Label>
        <Input
          id="applyRemark"
          placeholder="选填，例如付款单号或留言说明"
          value={form.applyRemark}
          onChange={(e) => setForm({ ...form, applyRemark: e.target.value })}
        />
      </div>

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
