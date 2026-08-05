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
import { Heart } from 'lucide-react'
import type { SponsorApplyRequest } from '@/api/types'
import { SponsorApplyForm } from '@/components/sponsor-apply-form'
import { BrushUnderline, EnsoIcon, InkGlow } from '@/components/ink-wash'
import { Button } from '@/components/ui/button'
import { enter } from '@/lib/motion'
import { useApplySponsor, usePublicChannels } from '@/hooks/use-sponsors'
import { getToken } from '@/lib/auth'

export const Route = createFileRoute('/operate/sponsor')({
  component: SponsorApplyPage,
})

/**
 * 自助管理 · 赞助展示申请页（信笺版）。
 *
 * 结构身份：赞助 = 心意的回礼 = 书信往来。开场是亲密的邀请小字（非 about 的大字感恩帖），
 * 表单收进一张「信笺」——晨光墨晕 + 抬头 + 行格字段 + 落款递交；
 * 提交完成态以 enso 缺口圆签名「申请已寄出」。
 */
function SponsorApplyPage() {
  const reduced = useReducedMotion() ?? false
  const applyMutation = useApplySponsor()
  const [submitted, setSubmitted] = useState(false)
  const isAuthed = getToken() != null

  const handleSubmit = (req: SponsorApplyRequest) => {
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
              让竹林<span className="text-leaf-deep">记住</span>你的心意
            </h1>
            <div className="mt-4 flex justify-center">
              <BrushUnderline />
            </div>
            <p className="mx-auto mt-4 max-w-md font-serif text-[15px] italic leading-relaxed text-text-secondary">
              递上一份你的赞助，若经核实入册，便将它记入感恩账册，与竹林共沐晨光。
            </p>
          </motion.section>

          {/* ═══════════ 信笺 · 申请表单 + 渠道说明（桌面双栏） ═══════════ */}
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
                  <span
                    className="h-[3px] w-5 -skew-x-12 rounded-sm bg-leaf-deep"
                    aria-hidden
                  />
                  赞助申请
                </h2>
                <span className="font-mono text-[11px] uppercase tracking-[0.15em] text-text-secondary">
                  Sponsor
                </span>
              </div>

              {/* 行格字段（SponsorApplyForm，金额元→分、渠道选择器、邮箱锁定） */}
              <div className="relative">
                <SponsorApplyForm
                  submitting={applyMutation.isPending}
                  submitLabel="递交申请"
                  onSubmit={handleSubmit}
                />
              </div>

              {/* 落款提示 */}
              <p className="relative mt-5 border-t border-border/60 pt-4 font-serif text-[13px] italic leading-relaxed text-text-secondary">
                提交后将进入待审核状态，审核结果将通过联系邮箱通知你。注册或登录时使用相同邮箱，便可在「我的赞助」中管理它。
              </p>
            </motion.section>

            {/* 渠道说明卡 */}
            <ChannelInfoCard reduced={reduced} />
          </div>
        </>
      )}
    </main>
  )
}

/**
 * 赞助渠道说明卡：展示当前可选的赞助渠道与申请说明。
 * 渠道数据取自公开接口 GET /sponsors/channels（后台「赞助渠道」面板维护）。
 */
function ChannelInfoCard({ reduced }: { reduced: boolean }) {
  const { data: channels } = usePublicChannels()

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
          <span
            className="h-[3px] w-4 -skew-x-12 rounded-sm bg-leaf-deep"
            aria-hidden
          />
          赞助渠道
        </h2>
        <span className="font-mono text-[11px] uppercase tracking-[0.15em] text-text-secondary">
          Channels
        </span>
      </div>

      {/* 提示 */}
      <div className="relative mb-4 rounded-md border border-leaf-muted/40 bg-leaf-light/20 px-3.5 py-3">
        <p className="font-serif text-[13px] italic leading-relaxed text-text-secondary">
          通过以下渠道实际赞助后，再递交本申请。经管理员核实后，您的支持将展示在「赞助」感恩账册中。
        </p>
      </div>

      {/* 渠道列表 */}
      <dl className="relative">
        {(channels?.length ? channels : []).map((ch, i) => (
          <div
            key={ch.id.toString()}
            className={`flex items-center justify-between gap-3 py-2.5 ${
              i < (channels?.length ?? 0) - 1 ? 'border-b border-border/60' : ''
            }`}
          >
            <dt className="font-serif text-sm text-text-primary">{ch.name}</dt>
            <dd className="flex items-center gap-1 font-mono text-[11px] text-text-secondary">
              <Heart className="size-3 text-leaf-deep" aria-hidden />
              {ch.sponsor_count} 次随喜
            </dd>
          </div>
        ))}
        {!channels?.length && (
          <p className="py-2 text-sm text-text-secondary">暂无渠道</p>
        )}
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
        你的赞助申请已进入待审核状态，审核通过后会展示在赞助页，
        审核结果将通过联系邮箱通知你。
      </p>
      <div className="relative mt-7 flex flex-wrap items-center justify-center gap-3">
        <Button variant="outline" className="cursor-pointer" onClick={onReset}>
          再递交一份
        </Button>
        {isAuthed ? (
          <Link to="/user/sponsors">
            <Button className="cursor-pointer">前往我的赞助</Button>
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
