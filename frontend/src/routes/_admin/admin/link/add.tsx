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

import { createFileRoute, useNavigate } from '@tanstack/react-router'
import { LinkForm } from '@/components/link-form'
import { PageHead } from '@/components/ink-wash'
import { useCreateLink } from '@/hooks/use-links'

export const Route = createFileRoute('/_admin/admin/link/add')({
  component: LinkAddPage,
})

/**
 * 添加友链：方案 C「信息卡群」。
 * 页面只给 PageHead（题跋开场），表单卡群由 LinkForm 自行编排。
 */
function LinkAddPage() {
  const navigate = useNavigate()
  const createLink = useCreateLink()

  return (
    <div className="mx-auto max-w-6xl space-y-6">
      <PageHead
        backTo="/admin/link"
        backLabel="返回友链管理"
        kicker="LINKS · 友链"
        title="添加友链"
        sub="添加一个新的友情链接，带 * 为必填项。"
      />
      <LinkForm
        submitting={createLink.isPending}
        onSubmit={(req) =>
          createLink.mutate(req, {
            onSuccess: () => navigate({ to: '/admin/link' }),
          })
        }
      />
    </div>
  )
}
