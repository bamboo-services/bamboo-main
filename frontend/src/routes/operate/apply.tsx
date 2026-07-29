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
import { useState } from 'react'
import type { ApplyLinkRequest, BloggerInfoResponse } from '@/api/types'
import { UserLinkForm } from '@/components/user-link-form'
import { BrushUnderline, EnsoIcon, InkGlow } from '@/components/ink-wash'
import { Button } from '@/components/ui/button'
import { enter } from '@/lib/motion'
import { useApplyLink } from '@/hooks/use-links'
import { useBloggerInfo } from '@/hooks/use-site-info'
import { getToken } from '@/lib/auth'

export const Route = createFileRoute('/operate/apply')({
  component: ApplyPage,
})

/**
 * 自助管理 · 友链申请页（信笺版）。
 *
 * 结构身份：友链 = 结交 = 书信往来。开场是亲密的邀请小字（非 about 的大字名士帖），
 * 表单收进一张「信笺」——晨光墨晕 + 抬头 + 行格字段 + 落款递交；
 * 提交完成态以 enso 缺口圆签名「申请已寄出」。
 */
function ApplyPage() {
  const reduced = useReducedMotion() ?? false
  const applyMutation = useApplyLink()
  const [submitted, setSubmitted] = useState(false)
  const isAuthed = getToken() != null
  const { data: blogger } = useBloggerInfo()

  const handleSubmit = (req: ApplyLinkRequest) => {
    applyMutation.mutate(req, {
      onSuccess: () => setSubmitted(true),
    })
  }

  return (
    <main className="mx-auto w-full max-w-5xl flex-1 px-6 pb-20 md:px-10">
      {submitted ? (
        <SentState
          reduced={reduced}
          isAuthed={isAuthed}
          onReset={() => setSubmitted(false)}
        />
      ) : (
        <>
          {/* ═══════════ 开场 · 邀请（亲密小字，居中） ═══════════ */}
          <motion.section
            {...enter(reduced, 0.08, {
              initial: { opacity: 0, y: 16 },
              animate: { opacity: 1, y: 0 },
              transition: { duration: 0.7, ease: [0.22, 0.9, 0.3, 1] as const },
            })}
            className="mb-10 pt-4 text-center"
          >
            <p className="flex items-center justify-center gap-3 font-mono text-[11px] uppercase tracking-[0.35em] text-text-secondary">
              <span className="h-px w-10 bg-leaf-deep" aria-hidden />
              自助管理
              <span className="h-px w-10 bg-leaf-deep" aria-hidden />
            </p>
            <h1 className="mt-5 font-serif text-[clamp(1.75rem,4vw,2.5rem)] font-bold leading-[1.15] tracking-[0.02em] text-text-primary">
              让竹林<span className="text-leaf-deep">认识你</span>的小站
            </h1>
            <div className="mt-4 flex justify-center">
              <BrushUnderline />
            </div>
            <p className="mx-auto mt-4 max-w-md font-serif text-[15px] italic leading-relaxed text-text-secondary">
              递交一座你的小站，若竹林也认得你，便将它种在友链之畔，共沐晨光。
            </p>
          </motion.section>

          {/* ═══════════ 信笺 · 申请表单 + 博主信息（桌面双栏） ═══════════ */}
          <div className="grid items-start gap-6 lg:grid-cols-[minmax(0,1fr)_19rem]">
            {/* 信笺表单 */}
            <motion.section
              {...enter(reduced, 0.2, {
                initial: { opacity: 0, y: 20 },
                animate: { opacity: 1, y: 0 },
                transition: { duration: 0.6, ease: 'easeOut' },
              })}
              className="relative overflow-hidden rounded-lg border border-border bg-card p-6 shadow-[0_18px_50px_-28px_oklch(0.32_0.06_155/0.3)] md:p-9"
            >
              <InkGlow />

              {/* 信笺抬头 */}
              <div className="relative mb-6 flex items-center justify-between gap-2 border-b border-border pb-4">
                <h2 className="flex items-center gap-2.5 font-serif text-lg font-semibold text-text-primary">
                  <span className="h-[3px] w-5 -skew-x-12 rounded-sm bg-leaf-deep" aria-hidden />
                  友链申请
                </h2>
                <span className="font-mono text-[11px] uppercase tracking-[0.15em] text-text-secondary">
                  Application
                </span>
              </div>

              {/* 行格字段（与编辑页共用 UserLinkForm，字段对齐 ApplyLinkRequest） */}
              <div className="relative">
                <UserLinkForm
                  submitting={applyMutation.isPending}
                  submitLabel="递交申请"
                  onSubmit={handleSubmit}
                />
              </div>

              {/* 落款提示 */}
              <p className="relative mt-5 border-t border-border/60 pt-4 font-serif text-[13px] italic leading-relaxed text-text-secondary">
                提交后将进入待审核状态，审核结果将通过站长邮箱通知你。注册或登录时使用相同邮箱，便可在「我的友链」中管理它。
              </p>
            </motion.section>

            {/* 博主站点信息 · 交换友链提示 */}
            <ExchangeInfoCard blogger={blogger} reduced={reduced} />
          </div>
        </>
      )}
    </main>
  )
}

/**
 * 博主站点信息卡：交换友链前需先在自站添加博主友链。
 *
 * 六个字段（站点名字/描述/地址/图片/订阅/邮箱）全部取自后端
 * GET /info/blogger（admin「博主信息」设置面板维护），不再写死前端。
 */
function ExchangeInfoCard({
  blogger,
  reduced,
}: {
  blogger?: BloggerInfoResponse
  reduced: boolean
}) {
  const rows = [
    { label: '站点名字', value: blogger?.site_name ?? '—' },
    { label: '站点描述', value: blogger?.site_description ?? '—' },
    { label: '站点地址', value: blogger?.site_url ?? '—' },
    { label: '站点图片', value: blogger?.site_image ?? '—' },
    { label: '站点订阅', value: blogger?.rss ?? '—' },
    { label: '站长邮箱', value: blogger?.email ?? '—' },
  ]

  return (
    <motion.aside
      {...enter(reduced, 0.32, {
        initial: { opacity: 0, y: 20 },
        animate: { opacity: 1, y: 0 },
        transition: { duration: 0.6, ease: 'easeOut' },
      })}
      className="relative overflow-hidden rounded-lg border border-border bg-card p-5 shadow-[0_18px_50px_-28px_oklch(0.32_0.06_155/0.3)]"
    >
      <InkGlow />

      {/* 抬头 */}
      <div className="relative mb-4 flex items-center justify-between gap-2 border-b border-border pb-3">
        <h2 className="flex items-center gap-2.5 font-serif text-base font-semibold text-text-primary">
          <span className="h-[3px] w-4 -skew-x-12 rounded-sm bg-leaf-deep" aria-hidden />
          添加博主友链
        </h2>
        <span className="font-mono text-[11px] uppercase tracking-[0.15em] text-text-secondary">
          Exchange
        </span>
      </div>

      {/* 提示：先加博主友链再申请 */}
      <div className="relative mb-4 rounded-md border border-leaf-muted/40 bg-leaf-light/20 px-3.5 py-3">
        <p className="font-serif text-[13px] italic leading-relaxed text-text-secondary">
          友链审核通过后，会在友链页面展示您的站点信息。添加友链之前，请先在您的站点添加博主友链，添加后再申请，避免审核时未看到友链导致审核失败。
        </p>
      </div>

      {/* 博主站点信息行 */}
      <dl className="relative">
        {rows.map((row, i) => (
          <div
            key={row.label}
            className={`py-2.5 ${i < rows.length - 1 ? 'border-b border-border/60' : ''}`}
          >
            <dt className="font-mono text-[11px] uppercase tracking-[0.12em] text-text-secondary">
              {row.label}
            </dt>
            <dd className="mt-1 break-all font-mono text-xs leading-relaxed text-text-primary">
              {row.value}
            </dd>
          </div>
        ))}
      </dl>
    </motion.aside>
  )
}

/** 提交完成态：enso 缺口圆 +「申请已寄出」 */
function SentState({
  reduced,
  isAuthed,
  onReset,
}: {
  reduced: boolean
  isAuthed: boolean
  onReset: () => void
}) {
  return (
    <motion.section
      {...enter(reduced, 0.1, {
        initial: { opacity: 0, y: 16 },
        animate: { opacity: 1, y: 0 },
        transition: { duration: 0.6, ease: 'easeOut' },
      })}
      className="relative mt-6 overflow-hidden rounded-lg border border-border bg-card px-6 py-16 text-center shadow-[0_18px_50px_-28px_oklch(0.32_0.06_155/0.3)]"
    >
      <InkGlow />
      <EnsoIcon className="relative mx-auto size-16" />
      <h2 className="relative mt-5 font-serif text-2xl font-semibold text-text-primary">
        申请已寄出
      </h2>
      <p className="relative mx-auto mt-2 max-w-sm text-sm leading-relaxed text-text-secondary">
        你的友链申请已进入待审核状态，审核通过后会展示在友链页，
        审核结果将通过站长邮箱通知你。
      </p>
      <div className="relative mt-7 flex flex-wrap items-center justify-center gap-3">
        <Button variant="outline" className="cursor-pointer" onClick={onReset}>
          再递交一座
        </Button>
        {isAuthed ? (
          <Link to="/user/links">
            <Button className="cursor-pointer">前往我的友链</Button>
          </Link>
        ) : (
          <Link to="/auth/register">
            <Button className="cursor-pointer">注册账号以跟踪审核状态</Button>
          </Link>
        )}
      </div>
    </motion.section>
  )
}
