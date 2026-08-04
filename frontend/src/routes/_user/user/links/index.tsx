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
import { useCallback, useState } from 'react'
import { ArchiveX, ExternalLink, Pencil } from 'lucide-react'
import type { LinkFriend } from '@/api/types'
import type { InterludeData } from '@/components/about/interlude'
import { useMyLinks, useRequestTakedown } from '@/hooks/use-links'
import {
  EnsoEmpty,
  InkBadge,
  PageHead,
  inkCard,
  linkStatus,
} from '@/components/ink-wash'
import { Skeleton } from '@/components/ui/skeleton'
import { Button } from '@/components/ui/button'
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from '@/components/ui/dialog'
import { Interlude } from '@/components/about/interlude'
import { enter } from '@/lib/motion'

export const Route = createFileRoute('/_user/user/links/')({
  component: MyLinksPage,
})

function MyLinksPage() {
  const reduced = useReducedMotion() ?? false
  const { data, isLoading, error } = useMyLinks({ page: 1, page_size: 50 })
  const takedown = useRequestTakedown()
  const [interlude, setInterlude] = useState<InterludeData | null>(null)
  const [takedownTarget, setTakedownTarget] = useState<LinkFriend | null>(null)

  const links = data?.data ?? []

  /** 点击访问站点 → 组装 Interlude 数据，随墨入站 */
  const handleOpen = useCallback(
    (link: LinkFriend, origin: { x: number; y: number }) => {
      setInterlude({
        name: link.name,
        url: link.url,
        avatarChar: link.name.slice(0, 1),
        premium: false,
        origin,
      })
    },
    [],
  )

  const closeInterlude = useCallback(() => setInterlude(null), [])

  const confirmTakedown = () => {
    if (!takedownTarget) return
    takedown.mutate(takedownTarget.id, {
      onSuccess: () => setTakedownTarget(null),
    })
  }

  return (
    <div className="mx-auto max-w-4xl">
      <PageHead
        kicker="my links · 各美其美"
        title="我的友链"
        sub="查看你申请的友链与审核状态，可随时编辑或申请下架。"
        actions={
          <Link to="/operate/apply">
            <Button className="cursor-pointer">申请友链</Button>
          </Link>
        }
      />

      {/* ═══════════ 列表 / 骨架 / 空态 ═══════════ */}
      <div className="mt-8">
        {isLoading ? (
          <div className="space-y-4">
            {Array.from({ length: 3 }).map((_, i) => (
              <Skeleton key={i} className="h-32 rounded-lg" />
            ))}
          </div>
        ) : error || links.length === 0 ? (
          <div className={`${inkCard} p-8`}>
            <EnsoEmpty title="还没有友链" hint="提交申请，让竹林认识你的小站">
              <Link to="/operate/apply" className="ml-auto">
                <Button className="cursor-pointer">去申请友链</Button>
              </Link>
            </EnsoEmpty>
          </div>
        ) : (
          <div className="space-y-4">
            {links.map((link) => {
              const status = linkStatus(link)
              return (
                <motion.div
                  key={link.id.toString()}
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
                          {link.name}
                        </h3>
                        <InkBadge tone={status.tone}>{status.label}</InkBadge>
                      </div>
                      <p className="mt-1 truncate font-mono text-xs text-text-secondary">
                        {link.url}
                      </p>
                      {link.description && (
                        <p className="mt-2 line-clamp-2 text-sm text-text-secondary">
                          {link.description}
                        </p>
                      )}
                      {link.review_remark && (
                        <p className="mt-2 text-xs text-text-secondary">
                          审核备注：{link.review_remark}
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
                      onClick={(e) =>
                        handleOpen(link, { x: e.clientX, y: e.clientY })
                      }
                    >
                      <ExternalLink className="size-4" />
                      访问站点
                    </Button>
                    <Link
                      to="/user/links/$id/edit"
                      params={{ id: link.id.toString() }}
                    >
                      <Button
                        variant="outline"
                        size="sm"
                        className="cursor-pointer"
                      >
                        <Pencil className="size-4" />
                        编辑
                      </Button>
                    </Link>
                    {link.status === 1 && (
                      <Button
                        variant="outline"
                        size="sm"
                        className="cursor-pointer"
                        onClick={() => setTakedownTarget(link)}
                      >
                        <ArchiveX className="size-4" />
                        申请下架
                      </Button>
                    )}
                  </div>
                </motion.div>
              )
            })}
          </div>
        )}
      </div>

      {/* 沉浸式跳转引导层 */}
      <Interlude data={interlude} onDone={closeInterlude} />

      {/* 下架确认 Dialog */}
      <Dialog
        open={takedownTarget != null}
        onOpenChange={(open) => !open && setTakedownTarget(null)}
      >
        <DialogContent>
          <DialogHeader>
            <DialogTitle>申请下架友链</DialogTitle>
            <DialogDescription>
              确定要申请下架「{takedownTarget?.name}
              」吗？提交后需等待管理员审核， 审核通过后该友链将不再公开展示。
            </DialogDescription>
          </DialogHeader>
          <DialogFooter>
            <Button
              variant="outline"
              className="cursor-pointer"
              onClick={() => setTakedownTarget(null)}
            >
              取消
            </Button>
            <Button
              className="cursor-pointer"
              onClick={confirmTakedown}
              disabled={takedown.isPending}
            >
              {takedown.isPending ? '提交中…' : '确认申请'}
            </Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>
    </div>
  )
}
