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
import { Link, createFileRoute, useNavigate } from '@tanstack/react-router'
import { motion, useReducedMotion } from 'motion/react'
import { Save } from 'lucide-react'
import type { FormEvent } from 'react'
import type { EditApplyRequest } from '@/api/types'
import {
  useEditApply,
  useMyLink,
  usePublicColors,
  usePublicGroups,
} from '@/hooks/use-links'
import { EnsoEmpty, PageHead, inkCard } from '@/components/ink-wash'
import { Button } from '@/components/ui/button'
import { Label } from '@/components/ui/label'
import { Select } from '@/components/ui/select'
import { Textarea } from '@/components/ui/textarea'
import { Skeleton } from '@/components/ui/skeleton'
import { enter } from '@/lib/motion'

export const Route = createFileRoute('/_user/user/links/$id/edit-location')({
  component: EditLocationPage,
})

function EditLocationPage() {
  const reduced = useReducedMotion() ?? false
  const { id } = Route.useParams()
  const navigate = useNavigate()
  const linkId = BigInt(id)
  const { data: link, isLoading, error } = useMyLink(linkId)
  const { data: groups } = usePublicGroups()
  const { data: colors } = usePublicColors()
  const editApply = useEditApply()

  const [groupId, setGroupId] = useState('')
  const [colorId, setColorId] = useState('')
  const [remark, setRemark] = useState('')

  // 至少提供一个预期变更才可提交
  const canSubmit = groupId !== '' || colorId !== ''

  const handleSubmit = (e: FormEvent) => {
    e.preventDefault()
    if (!canSubmit || editApply.isPending) return
    const req: EditApplyRequest = {
      // 空值即「保持原值」（undefined → 后端省略），仅发送实际变更项
      ...(groupId !== '' ? { link_group_id: BigInt(groupId) } : {}),
      ...(colorId !== '' ? { link_color_id: BigInt(colorId) } : {}),
      ...(remark.trim() ? { link_apply_remark: remark.trim() } : {}),
    }
    editApply.mutate(
      { id: linkId, req },
      { onSuccess: () => void navigate({ to: '/user/links' }) },
    )
  }

  return (
    <div className="mx-auto max-w-3xl">
      <PageHead
        kicker="edit location · 申请修改"
        title="申请修改展示位置/颜色"
        backTo="/user/links"
        backLabel="返回我的友链"
        sub="选择预期的展示位置与展示颜色，提交后进入待审核队列；审核通过前友链仍按当前设置展示。"
      />

      {/* ═══════════ 表单 / 骨架 / 守卫态 ═══════════ */}
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
        ) : link.status !== 1 || link.is_failure === 1 ? (
          <div className={`${inkCard} p-8`}>
            <EnsoEmpty
              title="仅已通过的友链可申请修改"
              hint="友链通过审核后才能申请调整展示位置与颜色"
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
            <form onSubmit={handleSubmit} className="space-y-5">
              {/* 当前 vs 预期：位置 | 颜色 */}
              <div className="grid gap-4 md:grid-cols-2">
                <div className="space-y-2">
                  <Label>当前展示位置</Label>
                  <div className="rounded-md border border-border/70 bg-muted/40 px-3 py-2 text-sm text-text-secondary">
                    {link.group_f_key?.name ??
                      (link.group_id ? `位置 #${link.group_id}` : '未分组')}
                  </div>
                </div>
                <div className="space-y-2">
                  <Label>当前展示颜色</Label>
                  <div className="rounded-md border border-border/70 bg-muted/40 px-3 py-2 text-sm text-text-secondary">
                    {link.color_f_key?.name ??
                      (link.color_id ? `颜色 #${link.color_id}` : '默认颜色')}
                  </div>
                </div>
                <div className="space-y-2">
                  <Label htmlFor="groupId">预期展示位置</Label>
                  <Select
                    id="groupId"
                    value={groupId}
                    onChange={(e) => setGroupId(e.target.value)}
                  >
                    <option value="">保持原位置</option>
                    {groups?.map((g) => (
                      <option key={g.id.toString()} value={g.id.toString()}>
                        {g.name}
                      </option>
                    ))}
                  </Select>
                </div>
                <div className="space-y-2">
                  <Label htmlFor="colorId">预期展示颜色</Label>
                  <Select
                    id="colorId"
                    value={colorId}
                    onChange={(e) => setColorId(e.target.value)}
                  >
                    <option value="">保持原颜色</option>
                    {colors?.map((c) => (
                      <option key={c.id.toString()} value={c.id.toString()}>
                        {c.name}
                      </option>
                    ))}
                  </Select>
                </div>
              </div>

              {/* 申请备注（仅博主审核可见） */}
              <div className="space-y-2">
                <Label htmlFor="remark">申请备注</Label>
                <Textarea
                  id="remark"
                  className="min-h-[80px]"
                  placeholder="选填，例如说明调整原因…"
                  value={remark}
                  onChange={(e) => setRemark(e.target.value)}
                />
                <p className="text-[11px] text-text-secondary">
                  备注仅供博主审核时查看，不会公开展示
                </p>
              </div>

              {/* 提交按钮 */}
              <div className="flex justify-end border-t border-border/60 pt-5">
                <Button
                  type="submit"
                  className="cursor-pointer"
                  disabled={editApply.isPending || !canSubmit}
                >
                  <Save className="mr-2 size-4" />
                  {editApply.isPending ? '提交中…' : '提交申请'}
                </Button>
              </div>
            </form>
          </motion.div>
        )}
      </div>
    </div>
  )
}
