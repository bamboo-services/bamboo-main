// --------------------------------------------------------------------------------
// Copyright (c) 2016-NOW(至今) 筱锋
// Author: 筱锋「xiao_lfeng」(https://www.x-lf.com)
// --------------------------------------------------------------------------------
// 许可证声明：版权所有 (c) 2016-2026 筱锋。保留所有权利。
// 有关MIT许可证的更多信息，请查看项目根目录下的LICENSE文件或访问：
// https://opensource.org/licenses/MIT
// --------------------------------------------------------------------------------

import { Link, createFileRoute } from '@tanstack/react-router'
import { motion, useReducedMotion } from 'motion/react'
import { useEffect, useState } from 'react'
import { Heart, Pencil } from 'lucide-react'
import type { SponsorRecordAdmin, SponsorUserUpdateRequest } from '@/api/types'
import { useMySponsors, useUpdateMySponsor } from '@/hooks/use-sponsors'
import { EnsoEmpty, InkBadge, PageHead, inkCard } from '@/components/ink-wash'
import { Skeleton } from '@/components/ui/skeleton'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { Textarea } from '@/components/ui/textarea'
import { Checkbox } from '@/components/ui/checkbox'
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from '@/components/ui/dialog'
import { enter } from '@/lib/motion'

export const Route = createFileRoute('/_user/user/sponsors/')({
  component: MySponsorsPage,
})

/** 赞助记录状态徽章 */
function sponsorStatus(record: SponsorRecordAdmin) {
  switch (record.status) {
    case 0:
      return { label: '待审核', tone: 'pending' as const }
    case 1:
      return { label: '已通过', tone: 'leaf' as const }
    case 2:
      return { label: '已拒绝', tone: 'danger' as const }
    default:
      return { label: '未知', tone: 'neutral' as const }
  }
}

/** 金额（分）→ 元字符串 */
function formatAmount(amount: number): string {
  return (amount / 100).toFixed(2)
}

/** 内联编辑表单状态（金额与渠道不可由用户修改） */
interface EditFormState {
  nickname: string
  message: string
  sponsorAt: string
  redirectUrl: string
  isAnonymous: boolean
  applyRemark: string
}

/** 从赞助记录初始化编辑表单 */
function initEditForm(record: SponsorRecordAdmin): EditFormState {
  return {
    nickname: record.nickname,
    message: record.message ?? '',
    sponsorAt: record.sponsor_at ? record.sponsor_at.slice(0, 10) : '',
    redirectUrl: record.redirect_url ?? '',
    isAnonymous: record.is_anonymous,
    applyRemark: record.apply_remark ?? '',
  }
}

function MySponsorsPage() {
  const reduced = useReducedMotion() ?? false
  const { data, isLoading, error } = useMySponsors({ page: 1, page_size: 50 })
  const updateSponsor = useUpdateMySponsor()
  const [editTarget, setEditTarget] = useState<SponsorRecordAdmin | null>(null)
  const [editForm, setEditForm] = useState<EditFormState | null>(null)

  const records = data?.data ?? []

  // 打开编辑弹窗时预填表单
  useEffect(() => {
    if (editTarget) {
      setEditForm(initEditForm(editTarget))
    } else {
      setEditForm(null)
    }
  }, [editTarget])

  const handleSave = () => {
    if (!editTarget || !editForm) return
    const req: SponsorUserUpdateRequest = {
      nickname: editForm.nickname.trim() || undefined,
      message: editForm.message.trim() || undefined,
      sponsor_at: editForm.sponsorAt
        ? new Date(`${editForm.sponsorAt}T12:00:00`).toISOString()
        : undefined,
      redirect_url: editForm.redirectUrl.trim() || undefined,
      is_anonymous: editForm.isAnonymous || undefined,
      apply_remark: editForm.applyRemark.trim() || undefined,
    }
    updateSponsor.mutate(
      { id: editTarget.id, req },
      { onSuccess: () => setEditTarget(null) },
    )
  }

  return (
    <div className="mx-auto max-w-4xl">
      <PageHead
        kicker="my sponsor · 心意的回礼"
        title="我的赞助"
        sub="查看你的赞助支持与审核状态，可随时编辑展示信息。"
        actions={
          <Link to="/operate/sponsor">
            <Button className="cursor-pointer">申请赞助展示</Button>
          </Link>
        }
      />

      {/* ═══════════ 列表 / 骨架 / 空态 ═══════════ */}
      <div className="mt-8">
        {isLoading ? (
          <div className="space-y-4">
            {Array.from({ length: 3 }).map((_, i) => (
              <Skeleton key={i} className="h-28 rounded-lg" />
            ))}
          </div>
        ) : error || records.length === 0 ? (
          <div className={`${inkCard} p-8`}>
            <EnsoEmpty
              title="还没有赞助记录"
              hint="实际赞助后递交申请，经核实后展示在赞助页"
            >
              <Link to="/operate/sponsor" className="ml-auto">
                <Button className="cursor-pointer">去申请赞助</Button>
              </Link>
            </EnsoEmpty>
          </div>
        ) : (
          <div className="space-y-4">
            {records.map((record) => {
              const status = sponsorStatus(record)
              return (
                <motion.div
                  key={record.id.toString()}
                  {...enter(reduced, 0, {
                    initial: { opacity: 0, y: 12 },
                    animate: { opacity: 1, y: 0 },
                    transition: { duration: 0.4, ease: 'easeOut' },
                  })}
                  className={`${inkCard} p-5`}
                >
                  <div className="flex items-start justify-between gap-4">
                    <div className="min-w-0">
                      <div className="flex flex-wrap items-center gap-2.5">
                        <h3 className="font-serif text-lg font-semibold text-text-primary">
                          {record.nickname}
                        </h3>
                        <InkBadge tone={status.tone}>{status.label}</InkBadge>
                      </div>
                      <p className="mt-1 flex items-center gap-1.5 font-mono text-sm text-text-secondary">
                        <Heart className="size-3.5 text-leaf-deep" />¥
                        {formatAmount(record.amount)}
                        {record.channel
                          ? ` · ${record.channel.name}`
                          : ' · 未选择渠道'}
                      </p>
                      {record.message && (
                        <p className="mt-2 line-clamp-2 text-sm text-text-secondary">
                          {record.message}
                        </p>
                      )}
                      {record.apply_remark && (
                        <p className="mt-1 text-xs text-text-secondary">
                          申请备注：{record.apply_remark}
                        </p>
                      )}
                      {record.review_remark && (
                        <p className="mt-1 text-xs text-text-secondary">
                          审核备注：{record.review_remark}
                        </p>
                      )}
                    </div>
                  </div>

                  {/* 操作区 */}
                  <div className="mt-4 flex flex-wrap items-center gap-2 border-t border-border/60 pt-4">
                    <Button
                      variant="outline"
                      size="sm"
                      className="cursor-pointer"
                      onClick={() => setEditTarget(record)}
                    >
                      <Pencil className="size-4" />
                      编辑
                    </Button>
                  </div>
                </motion.div>
              )
            })}
          </div>
        )}
      </div>

      {/* 编辑 Dialog（金额与渠道不可修改） */}
      <Dialog
        open={editTarget != null}
        onOpenChange={(open) => !open && setEditTarget(null)}
      >
        <DialogContent>
          <DialogHeader>
            <DialogTitle>编辑赞助记录</DialogTitle>
            <DialogDescription>
              金额与渠道需与实际支付核验，不可修改；修改后重新进入审核队列。
            </DialogDescription>
          </DialogHeader>
          {editTarget && editForm && (
            <div className="grid gap-4">
              {/* 只读摘要 */}
              <div className="flex items-center gap-2 rounded-md border border-border/60 bg-muted/30 px-3.5 py-2.5 text-sm">
                <span className="font-serif font-semibold text-text-primary">
                  {editTarget.nickname}
                </span>
                <span className="font-mono text-leaf-deep">
                  ¥{formatAmount(editTarget.amount)}
                </span>
                {editTarget.channel && (
                  <span className="text-text-secondary">
                    · {editTarget.channel.name}
                  </span>
                )}
              </div>

              <div className="space-y-2">
                <Label htmlFor="e-nickname">
                  昵称 <span className="text-destructive">*</span>
                </Label>
                <Input
                  id="e-nickname"
                  value={editForm.nickname}
                  onChange={(e) =>
                    setEditForm({ ...editForm, nickname: e.target.value })
                  }
                />
              </div>
              <div className="grid gap-4 md:grid-cols-2">
                <div className="space-y-2">
                  <Label htmlFor="e-sponsorAt">赞助时间</Label>
                  <Input
                    id="e-sponsorAt"
                    type="date"
                    value={editForm.sponsorAt}
                    onChange={(e) =>
                      setEditForm({ ...editForm, sponsorAt: e.target.value })
                    }
                  />
                </div>
                <div className="space-y-2">
                  <Label htmlFor="e-redirectUrl">跳转链接</Label>
                  <Input
                    id="e-redirectUrl"
                    type="url"
                    value={editForm.redirectUrl}
                    onChange={(e) =>
                      setEditForm({ ...editForm, redirectUrl: e.target.value })
                    }
                  />
                </div>
              </div>
              <div className="space-y-2">
                <Label htmlFor="e-message">留言</Label>
                <Textarea
                  id="e-message"
                  className="min-h-[70px]"
                  value={editForm.message}
                  onChange={(e) =>
                    setEditForm({ ...editForm, message: e.target.value })
                  }
                />
              </div>
              <div className="space-y-2">
                <Label htmlFor="e-remark">申请备注</Label>
                <Input
                  id="e-remark"
                  value={editForm.applyRemark}
                  onChange={(e) =>
                    setEditForm({ ...editForm, applyRemark: e.target.value })
                  }
                />
              </div>
              <div className="flex items-center gap-2">
                <Checkbox
                  id="e-anonymous"
                  checked={editForm.isAnonymous}
                  onCheckedChange={(checked) =>
                    setEditForm({ ...editForm, isAnonymous: checked === true })
                  }
                />
                <Label htmlFor="e-anonymous" className="cursor-pointer">
                  匿名展示
                </Label>
              </div>
            </div>
          )}
          <DialogFooter>
            <Button
              variant="outline"
              className="cursor-pointer"
              onClick={() => setEditTarget(null)}
            >
              取消
            </Button>
            <Button
              className="cursor-pointer"
              onClick={handleSave}
              disabled={updateSponsor.isPending}
            >
              {updateSponsor.isPending ? '保存中…' : '保存'}
            </Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>
    </div>
  )
}
