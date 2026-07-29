// --------------------------------------------------------------------------------
// Copyright (c) 2016-NOW(至今) 筱锋
// Author: 筱锋「xiao_lfeng」(https://www.x-lf.com)
// --------------------------------------------------------------------------------
// 许可证声明：版权所有 (c) 2016-2026 筱锋。保留所有权利。
// 有关MIT许可证的更多信息，请查看项目根目录下的LICENSE文件或访问：
// https://opensource.org/licenses/MIT
// --------------------------------------------------------------------------------

import {
  Link,
  createFileRoute,
  useNavigate,
} from '@tanstack/react-router'
import { motion, useReducedMotion } from 'motion/react'
import type { ApplyLinkRequest } from '@/api/types'
import { useMyLink, useUpdateMyLink } from '@/hooks/use-links'
import { UserLinkForm } from '@/components/user-link-form'
import { EnsoEmpty, PageHead, inkCard } from '@/components/ink-wash'
import { Skeleton } from '@/components/ui/skeleton'
import { Button } from '@/components/ui/button'
import { enter } from '@/lib/motion'

export const Route = createFileRoute('/_user/user/links/$id/edit')({
  component: EditLinkPage,
})

function EditLinkPage() {
  const reduced = useReducedMotion() ?? false
  const { id } = Route.useParams()
  const navigate = useNavigate()
  const linkId = BigInt(id)
  const { data: link, isLoading, error } = useMyLink(linkId)
  const updateMutation = useUpdateMyLink()

  const handleSubmit = (req: ApplyLinkRequest) => {
    updateMutation.mutate(
      { id: linkId, req },
      { onSuccess: () => void navigate({ to: '/user/links' }) },
    )
  }

  return (
    <div className="mx-auto max-w-3xl">
      <PageHead
        kicker="edit link"
        title="编辑友链"
        backTo="/user/links"
        backLabel="返回我的友链"
        sub="修改你的友链信息，保存后重新进入审核队列。"
      />

      {/* ═══════════ 表单 / 骨架 / 错误态 ═══════════ */}
      <div className="mt-8">
        {isLoading ? (
          <Skeleton className="h-96 rounded-lg" />
        ) : error || !link ? (
          <div className={`${inkCard} p-8`}>
            <EnsoEmpty
              title="友链不存在或无权访问"
              hint="它可能已被删除，或不属于当前账号"
            >
              <Link to="/user/links" className="ml-auto">
                <Button variant="outline" className="cursor-pointer">
                  返回我的友链
                </Button>
              </Link>
            </EnsoEmpty>
          </div>
        ) : (
          <motion.div
            {...enter(reduced, 0.15, {
              initial: { opacity: 0, y: 20 },
              animate: { opacity: 1, y: 0 },
              transition: { duration: 0.6, ease: 'easeOut' },
            })}
            className={`${inkCard} p-6 md:p-8`}
          >
            <UserLinkForm
              initial={link}
              submitting={updateMutation.isPending}
              submitLabel="保存修改"
              onSubmit={handleSubmit}
            />
          </motion.div>
        )}
      </div>
    </div>
  )
}
