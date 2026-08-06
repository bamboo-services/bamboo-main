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

import { useEffect, useState } from 'react'
import { createFileRoute } from '@tanstack/react-router'
import { motion, useReducedMotion } from 'motion/react'
import { Lock, Pencil, Plus, Power, PowerOff, Trash2 } from 'lucide-react'
import type { LinkGroup } from '@/api/types'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { Textarea } from '@/components/ui/textarea'
import { Skeleton } from '@/components/ui/skeleton'
import { Checkbox } from '@/components/ui/checkbox'
import { Select } from '@/components/ui/select'
import { Pagination } from '@/components/ui/pagination'
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from '@/components/ui/table'
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from '@/components/ui/dialog'
import {
  BambooRule,
  EnsoEmpty,
  InkBadge,
  PageHead,
  inkTableHeadRow,
  inkTableRow,
  inkTableWrap,
  inkTd,
  inkTh,
} from '@/components/ink-wash'
import { enter } from '@/lib/motion'
import { cn } from '@/lib/utils'
import { isBuiltinGroup } from '@/lib/locations'
import {
  useAllGroups,
  useCreateGroup,
  useDeleteGroup,
  useGroups,
  useUpdateGroup,
  useUpdateGroupStatus,
} from '@/hooks/use-groups'

export const Route = createFileRoute('/_admin/admin/location')({
  component: LocationPage,
})

const PAGE_SIZE = 10

function LocationPage() {
  const reduced = useReducedMotion() ?? false
  const [pageIndex, setPageIndex] = useState(0)

  // 新建 / 编辑弹窗（共用一个受控表单）
  const [formOpen, setFormOpen] = useState(false)
  const [editing, setEditing] = useState<LinkGroup | null>(null)
  const [name, setName] = useState('')
  const [desc, setDesc] = useState('')
  const [order, setOrder] = useState('0')

  // 删除确认弹窗
  const [deleteTarget, setDeleteTarget] = useState<LinkGroup | null>(null)
  const [forceDelete, setForceDelete] = useState(false)
  const [migrateTargetID, setMigrateTargetID] = useState('')

  const groupsQuery = useGroups({
    page: pageIndex + 1,
    page_size: PAGE_SIZE,
    order_by: 'sort_order',
    order: 'asc',
  })
  const allGroupsQuery = useAllGroups()
  const createGroup = useCreateGroup()
  const updateGroup = useUpdateGroup()
  const updateStatus = useUpdateGroupStatus()
  const deleteGroup = useDeleteGroup()

  const groups = groupsQuery.data?.data ?? []
  const total = groupsQuery.data?.pagination.total ?? 0
  const totalPages = groupsQuery.data?.pagination.total_pages ?? 1

  const submitting = createGroup.isPending || updateGroup.isPending

  // 可迁移的目标分组：排除内置「已失效」分组与当前被删分组自身
  const migrateCandidates = (allGroupsQuery.data ?? []).filter(
    (g) => !isBuiltinGroup(g) && g.id.toString() !== deleteTarget?.id.toString(),
  )

  // 弹窗打开时按编辑对象预填 / 重置表单
  useEffect(() => {
    if (formOpen) {
      setName(editing?.name ?? '')
      setDesc(editing?.description ?? '')
      setOrder(editing ? String(editing.sort_order) : '0')
    }
  }, [formOpen, editing])

  const openCreate = () => {
    setEditing(null)
    setFormOpen(true)
  }

  const openEdit = (group: LinkGroup) => {
    setEditing(group)
    setFormOpen(true)
  }

  const handleSubmit = () => {
    const parsedOrder = Number.parseInt(order, 10)
    const groupOrder = Number.isNaN(parsedOrder) ? 0 : parsedOrder
    const close = () => setFormOpen(false)

    if (editing) {
      updateGroup.mutate(
        {
          id: editing.id,
          req: {
            group_name: name.trim(),
            group_desc: desc.trim(),
            group_order: groupOrder,
          },
        },
        { onSuccess: close },
      )
    } else {
      createGroup.mutate(
        {
          group_name: name.trim(),
          group_desc: desc.trim() || undefined,
          group_order: groupOrder,
        },
        { onSuccess: close },
      )
    }
  }

  const handleDelete = () => {
    if (!deleteTarget) return
    deleteGroup.mutate(
      {
        id: deleteTarget.id,
        ...(migrateTargetID
          ? { targetGroupId: BigInt(migrateTargetID) }
          : { force: forceDelete }),
      },
      {
        onSuccess: () => {
          setDeleteTarget(null)
          setForceDelete(false)
          setMigrateTargetID('')
        },
      },
    )
  }

  return (
    <div className="mx-auto max-w-6xl space-y-6">
      <PageHead
        kicker="GROUPS · 位置"
        title="位置管理"
        sub="管理友链的展示位置（分组），控制友链在页面上的归类方式。"
        actions={
          <Button className="cursor-pointer" onClick={openCreate}>
            <Plus className="mr-2 size-4" />
            新建分组
          </Button>
        }
      />

      <BambooRule reduced={reduced} delay={0.12} />

      {/* 分组表格 */}
      <motion.section
        {...enter(reduced, 0.18, {
          initial: { opacity: 0, y: 10 },
          animate: { opacity: 1, y: 0 },
          transition: { duration: 0.4, ease: 'easeOut' },
        })}
        className={inkTableWrap}
      >
        <Table>
          <TableHeader>
            <TableRow className={cn(inkTableHeadRow, 'hover:bg-muted/30')}>
              <TableHead className={inkTh}>名称</TableHead>
              <TableHead className={cn(inkTh, 'hidden md:table-cell')}>
                描述
              </TableHead>
              <TableHead className={inkTh}>排序</TableHead>
              <TableHead className={inkTh}>状态</TableHead>
              <TableHead className={cn(inkTh, 'text-right')}>操作</TableHead>
            </TableRow>
          </TableHeader>
          <TableBody>
            {groupsQuery.isLoading ? (
              Array.from({ length: 5 }).map((_, i) => (
                <TableRow key={i}>
                  <TableCell className={inkTd}>
                    <Skeleton className="h-4 w-24" />
                  </TableCell>
                  <TableCell className={cn(inkTd, 'hidden md:table-cell')}>
                    <Skeleton className="h-4 w-48" />
                  </TableCell>
                  <TableCell className={inkTd}>
                    <Skeleton className="h-4 w-8" />
                  </TableCell>
                  <TableCell className={inkTd}>
                    <Skeleton className="h-5 w-12 rounded-full" />
                  </TableCell>
                  <TableCell className={inkTd}>
                    <div className="flex justify-end gap-1">
                      <Skeleton className="size-8 rounded-md" />
                      <Skeleton className="size-8 rounded-md" />
                      <Skeleton className="size-8 rounded-md" />
                    </div>
                  </TableCell>
                </TableRow>
              ))
            ) : groups.length === 0 ? (
              <TableRow className="hover:bg-transparent">
                <TableCell colSpan={5} className="px-4 py-10">
                  <EnsoEmpty
                    title="暂无分组"
                    hint="创建一个分组，让友链展示更有条理"
                  >
                    <Button
                      variant="outline"
                      size="sm"
                      className="ml-auto cursor-pointer"
                      onClick={openCreate}
                    >
                      <Plus className="mr-2 size-4" />
                      新建分组
                    </Button>
                  </EnsoEmpty>
                </TableCell>
              </TableRow>
            ) : (
              groups.map((group) => (
                <TableRow key={group.id.toString()} className={inkTableRow}>
                  <TableCell className={inkTd}>
                    <span className="flex items-center gap-2">
                      <span className="font-serif font-semibold text-text-primary">
                        {group.name}
                      </span>
                      {isBuiltinGroup(group) && (
                        <InkBadge tone="leaf">
                          <Lock className="size-3" />
                          内置
                        </InkBadge>
                      )}
                    </span>
                  </TableCell>
                  <TableCell className={cn(inkTd, 'hidden md:table-cell')}>
                    <span className="block max-w-xs truncate text-text-secondary">
                      {group.description ?? '—'}
                    </span>
                  </TableCell>
                  <TableCell className={inkTd}>
                    <span className="font-mono tabular-nums text-text-secondary">
                      {group.sort_order}
                    </span>
                  </TableCell>
                  <TableCell className={inkTd}>
                    {isBuiltinGroup(group) || group.status ? (
                      <InkBadge tone="leaf">启用</InkBadge>
                    ) : (
                      <InkBadge tone="neutral">禁用</InkBadge>
                    )}
                  </TableCell>
                  <TableCell className={inkTd}>
                    {isBuiltinGroup(group) ? (
                      <span className="font-mono text-xs text-text-secondary">
                        内置位置
                      </span>
                    ) : (
                      <div className="flex items-center justify-end gap-1">
                        <Button
                          variant="ghost"
                          size="icon"
                          title="编辑"
                          className="size-8 cursor-pointer"
                          onClick={() => openEdit(group)}
                        >
                          <Pencil className="size-4" />
                        </Button>
                        <Button
                          variant="ghost"
                          size="icon"
                          title={group.status ? '禁用' : '启用'}
                          className="size-8 cursor-pointer"
                          disabled={updateStatus.isPending}
                          onClick={() =>
                            updateStatus.mutate({
                              id: group.id,
                              status: !group.status,
                            })
                          }
                        >
                          {group.status ? (
                            <PowerOff className="size-4" />
                          ) : (
                            <Power className="size-4" />
                          )}
                        </Button>
                        <Button
                          variant="ghost"
                          size="icon"
                          title="删除"
                          className="size-8 cursor-pointer text-destructive hover:text-destructive"
                          onClick={() => setDeleteTarget(group)}
                        >
                          <Trash2 className="size-4" />
                        </Button>
                      </div>
                    )}
                  </TableCell>
                </TableRow>
              ))
            )}
          </TableBody>
        </Table>

        {!groupsQuery.isLoading && groups.length > 0 && (
          <div className="flex flex-wrap items-center justify-between gap-3 border-t border-border px-4 py-3">
            <span className="font-mono text-xs text-text-secondary">
              第 {pageIndex + 1} / {Math.max(totalPages, 1)} 页 · 共 {total} 条
            </span>
            <Pagination
              pageIndex={pageIndex}
              pageCount={Math.max(totalPages, 1)}
              onPageChange={setPageIndex}
            />
          </div>
        )}
      </motion.section>

      {/* 新建 / 编辑分组弹窗 */}
      <Dialog
        open={formOpen}
        onOpenChange={(open) => {
          if (!open) setFormOpen(false)
        }}
      >
        <DialogContent>
          <DialogHeader>
            <DialogTitle>{editing ? '编辑分组' : '新建分组'}</DialogTitle>
            <DialogDescription>
              {editing
                ? '修改分组信息后点击保存生效。'
                : '创建一个友链分组，友链可按位置归类展示。'}
            </DialogDescription>
          </DialogHeader>
          <div className="space-y-4">
            <div className="space-y-2">
              <Label htmlFor="groupName">分组名称</Label>
              <Input
                id="groupName"
                placeholder="如：技术博客"
                value={name}
                onChange={(e) => setName(e.target.value)}
              />
            </div>
            <div className="space-y-2">
              <Label htmlFor="groupDesc">分组描述</Label>
              <Textarea
                id="groupDesc"
                className="min-h-[80px]"
                placeholder="可选，简要说明该分组的用途"
                value={desc}
                onChange={(e) => setDesc(e.target.value)}
              />
            </div>
            <div className="space-y-2">
              <Label htmlFor="groupOrder">排序</Label>
              <Input
                id="groupOrder"
                type="number"
                value={order}
                onChange={(e) => setOrder(e.target.value)}
              />
              <p className="text-xs text-text-secondary">
                数值越小越靠前，默认为 0
              </p>
            </div>
          </div>
          <DialogFooter>
            <Button
              variant="outline"
              onClick={() => setFormOpen(false)}
              disabled={submitting}
            >
              取消
            </Button>
            <Button
              onClick={handleSubmit}
              disabled={!name.trim() || submitting}
            >
              {submitting ? '保存中…' : '保存'}
            </Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>

      {/* 删除确认弹窗 */}
      <Dialog
        open={deleteTarget !== null}
        onOpenChange={(open) => {
          if (!open) {
            setDeleteTarget(null)
            setForceDelete(false)
            setMigrateTargetID('')
          }
        }}
      >
        <DialogContent>
          <DialogHeader>
            <DialogTitle>删除分组</DialogTitle>
            <DialogDescription>
              确定要删除分组「{deleteTarget?.name}」吗？请先决定分组下友链的去向。
            </DialogDescription>
          </DialogHeader>
          <div className="space-y-4">
            <label className="flex cursor-pointer items-center gap-2.5 rounded-md border border-border bg-muted/30 px-3 py-2.5 text-sm transition-colors duration-150 hover:bg-muted/50">
              <Checkbox
                checked={forceDelete}
                onCheckedChange={(checked) => {
                  setForceDelete(checked === true)
                  if (checked) setMigrateTargetID('')
                }}
              />
              <span>移至未分组（该分组下的友链将移出分组）</span>
            </label>
            <div className="space-y-2">
              <Label htmlFor="migrateGroup">迁移到其他分组</Label>
              <Select
                id="migrateGroup"
                value={migrateTargetID}
                onChange={(e) => {
                  setMigrateTargetID(e.target.value)
                  if (e.target.value) setForceDelete(false)
                }}
              >
                <option value="">不迁移</option>
                {migrateCandidates.map((g) => (
                  <option
                    key={g.id.toString()}
                    value={g.id.toString()}
                    disabled={!g.status}
                  >
                    {g.name}
                    {!g.status ? '（已禁用）' : ''}
                  </option>
                ))}
              </Select>
              <p className="text-xs text-text-secondary">
                选择目标分组后，分组下友链将迁移过去再删除当前分组；禁用分组不可作为迁移目标。
              </p>
            </div>
          </div>
          <DialogFooter>
            <Button
              variant="outline"
              onClick={() => setDeleteTarget(null)}
              disabled={deleteGroup.isPending}
            >
              取消
            </Button>
            <Button
              variant="destructive"
              onClick={handleDelete}
              disabled={deleteGroup.isPending}
            >
              {deleteGroup.isPending ? '删除中…' : '确认删除'}
            </Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>
    </div>
  )
}
