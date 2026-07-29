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
import { Unlink } from 'lucide-react'
import { Button } from '@/components/ui/button'
import { Skeleton } from '@/components/ui/skeleton'
import { LinkForm } from '@/components/link-form'
import { PageHead, inkCard } from '@/components/ink-wash'
import { useAdminLink, useUpdateLink } from '@/hooks/use-links'

export const Route = createFileRoute('/_admin/admin/link/$id/edit')({
  component: LinkEditPage,
})

/**
 * 编辑友链：方案 C「信息卡群」，与添加页共用 LinkForm。
 * 加载时展示卡群骨架，加载完成预填表单。
 */
function LinkEditPage() {
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
    <div className="mx-auto max-w-6xl space-y-6">
      <PageHead
        backTo="/admin/link"
        backLabel="返回友链管理"
        kicker="LINKS · 友链"
        title="编辑友链"
        sub={link ? `正在编辑「${link.name}」的信息` : '正在加载友链信息…'}
      />

      {isLoading ? (
        /* 卡群骨架：站点信息 / 分类与级别两列 + 排序与备注跨行 */
        <div className="grid items-start gap-4 lg:grid-cols-2">
          <div className={`${inkCard} space-y-4`}>
            <Skeleton className="h-5 w-24" />
            <div className="grid gap-4 md:grid-cols-2">
              {Array.from({ length: 4 }).map((_, i) => (
                <Skeleton key={i} className="h-9 w-full" />
              ))}
            </div>
            <Skeleton className="h-24 w-full" />
          </div>
          <div className={`${inkCard} space-y-4`}>
            <Skeleton className="h-5 w-24" />
            <Skeleton className="h-8 w-full rounded-full" />
            <Skeleton className="h-8 w-full rounded-full" />
            <div className="flex gap-3">
              {Array.from({ length: 5 }).map((_, i) => (
                <Skeleton key={i} className="size-9 rounded-full" />
              ))}
            </div>
            <Skeleton className="h-9 w-full" />
          </div>
          <div className={`${inkCard} space-y-4 lg:col-span-2`}>
            <Skeleton className="h-5 w-24" />
            <Skeleton className="h-20 w-full" />
          </div>
        </div>
      ) : (
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
      )}
    </div>
  )
}
