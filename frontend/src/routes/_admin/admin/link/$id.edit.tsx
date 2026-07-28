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

import { Link, createFileRoute, useNavigate } from '@tanstack/react-router'
import { motion, useReducedMotion } from 'motion/react'
import { Unlink } from 'lucide-react'
import { Button } from '@/components/ui/button'
import { Skeleton } from '@/components/ui/skeleton'
import { LinkForm } from '@/components/link-form'
import { CardHead, PageHead, inkCard } from '@/components/ink-wash'
import { enter } from '@/lib/motion'
import { useAdminLink, useUpdateLink } from '@/hooks/use-links'

export const Route = createFileRoute('/_admin/admin/link/$id/edit')({
  component: LinkEditPage,
})

function LinkEditPage() {
  const reduced = useReducedMotion() ?? false
  const { id } = Route.useParams()
  const navigate = useNavigate()
  const linkId = BigInt(id)

  const { data: link, isLoading, isError } = useAdminLink(linkId)
  const updateLink = useUpdateLink()

  if (isError) {
    return (
      <div className="mx-auto flex max-w-md flex-col items-center justify-center gap-4 py-24 text-center">
        <Unlink className="size-10 text-text-secondary opacity-40" />
        <div>
          <h2 className="font-serif text-xl font-semibold text-text-primary">
            友链不存在
          </h2>
          <p className="mt-1 text-sm text-text-secondary">
            该友链可能已被删除或链接有误
          </p>
        </div>
        <Link to="/admin/link">
          <Button className="cursor-pointer">返回列表</Button>
        </Link>
      </div>
    )
  }

  return (
    <div className="mx-auto max-w-3xl space-y-6">
      <PageHead
        backTo="/admin/link"
        backLabel="返回友链管理"
        kicker="LINKS · 友链"
        title="编辑友链"
        sub={
          link ? `正在编辑「${link.name}」的信息` : '正在加载友链信息…'
        }
      />

      {isLoading ? (
        <div className={`${inkCard} space-y-6`}>
          <CardHead title="友链信息" />
          <div className="space-y-6">
            <div className="grid gap-4 md:grid-cols-2">
              {Array.from({ length: 4 }).map((_, i) => (
                <Skeleton key={i} className="h-9 w-full" />
              ))}
            </div>
            <Skeleton className="h-28 w-full" />
            <div className="grid gap-4 md:grid-cols-2">
              <Skeleton className="h-9 w-full" />
              <Skeleton className="h-9 w-full" />
            </div>
            <div className="flex gap-2">
              {Array.from({ length: 5 }).map((_, i) => (
                <Skeleton key={i} className="h-8 w-20 rounded-full" />
              ))}
            </div>
            <div className="flex gap-3">
              {Array.from({ length: 6 }).map((_, i) => (
                <Skeleton key={i} className="size-9 rounded-full" />
              ))}
            </div>
          </div>
        </div>
      ) : (
        <motion.section
          {...enter(reduced, 0.12, {
            initial: { opacity: 0, y: 10 },
            animate: { opacity: 1, y: 0 },
            transition: { duration: 0.4, ease: 'easeOut' },
          })}
          className={inkCard}
        >
          <CardHead title="友链信息" meta="EDIT LINK" />
          <LinkForm
            initial={link}
            submitting={updateLink.isPending}
            onSubmit={(req) =>
              updateLink.mutate(
                { id: linkId, req },
                { onSuccess: () => navigate({ to: '/admin/link' }) },
              )
            }
          />
        </motion.section>
      )}
    </div>
  )
}
