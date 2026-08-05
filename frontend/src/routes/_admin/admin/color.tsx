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
import {
  Pencil,
  Plus,
  Power,
  PowerOff,
  Sparkles,
  Trash2,
} from 'lucide-react'
import type { LinkColor } from '@/api/types'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { Skeleton } from '@/components/ui/skeleton'
import { Checkbox } from '@/components/ui/checkbox'
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
import {
  useColors,
  useCreateColor,
  useDeleteColor,
  useUpdateColor,
  useUpdateColorStatus,
} from '@/hooks/use-colors'
import { fancyGradient, isFancyColor } from '@/lib/colors'

export const Route = createFileRoute('/_admin/admin/color')({
  component: ColorPage,
})

const PAGE_SIZE = 10

/** 新建时的默认色值 */
const DEFAULT_PRIMARY = '#0ea5e9'
const DEFAULT_SUB = '#22d3ee'
const DEFAULT_HOVER = '#0284c7'

/** input[type=color] 仅接受 #rrggbb，非法值回退默认色 */
function normalizeHex(value: string | null, fallback: string): string {
  return value && /^#[0-9a-fA-F]{6}$/.test(value)
    ? value.toLowerCase()
    : fallback
}

/** 颜色选择器：原生取色器 + hex 展示 */
function ColorField({
  id,
  label,
  value,
  onChange,
}: {
  id: string
  label: string
  value: string
  onChange: (value: string) => void
}) {
  return (
    <div className="space-y-2">
      <Label htmlFor={id}>{label}</Label>
      <div className="flex items-center gap-2.5 rounded-md border border-input bg-background px-2.5 py-1.5 transition-colors duration-150 focus-within:border-leaf-deep focus-within:ring-leaf-deep/30 focus-within:ring-[3px]">
        <input
          id={id}
          type="color"
          value={value}
          onChange={(e) => onChange(e.target.value)}
          className="size-7 shrink-0 cursor-pointer appearance-none rounded-md border-0 bg-transparent p-0 [&::-webkit-color-swatch-wrapper]:p-0 [&::-webkit-color-swatch]:rounded-md [&::-webkit-color-swatch]:border-0"
        />
        <code className="font-mono text-xs uppercase text-text-secondary">
          {value}
        </code>
      </div>
    </div>
  )
}

function ColorPage() {
  const reduced = useReducedMotion() ?? false
  const [pageIndex, setPageIndex] = useState(0)

  // 新建 / 编辑弹窗（共用一个受控表单）
  const [formOpen, setFormOpen] = useState(false)
  const [editing, setEditing] = useState<LinkColor | null>(null)
  const [name, setName] = useState('')
  const [primary, setPrimary] = useState(DEFAULT_PRIMARY)
  const [sub, setSub] = useState(DEFAULT_SUB)
  const [hover, setHover] = useState(DEFAULT_HOVER)
  const [order, setOrder] = useState('0')

  // 删除确认弹窗
  const [deleteTarget, setDeleteTarget] = useState<LinkColor | null>(null)
  const [forceDelete, setForceDelete] = useState(false)

  const colorsQuery = useColors({
    page: pageIndex + 1,
    page_size: PAGE_SIZE,
    order_by: 'sort_order',
    order: 'asc',
  })
  const createColor = useCreateColor()
  const updateColor = useUpdateColor()
  const updateStatus = useUpdateColorStatus()
  const deleteColor = useDeleteColor()

  const colors = colorsQuery.data?.data ?? []
  const total = colorsQuery.data?.pagination.total ?? 0
  const totalPages = colorsQuery.data?.pagination.total_pages ?? 1

  const submitting = createColor.isPending || updateColor.isPending

  // 弹窗打开时按编辑对象预填 / 重置表单
  useEffect(() => {
    if (formOpen) {
      setName(editing?.name ?? '')
      setPrimary(normalizeHex(editing?.primary_color ?? null, DEFAULT_PRIMARY))
      setSub(normalizeHex(editing?.sub_color ?? null, DEFAULT_SUB))
      setHover(normalizeHex(editing?.hover_color ?? null, DEFAULT_HOVER))
      setOrder(editing ? String(editing.sort_order) : '0')
    }
  }, [formOpen, editing])

  const openCreate = () => {
    setEditing(null)
    setFormOpen(true)
  }

  const openEdit = (color: LinkColor) => {
    setEditing(color)
    setFormOpen(true)
  }

  const handleSubmit = () => {
    const parsedOrder = Number.parseInt(order, 10)
    const colorOrder = Number.isNaN(parsedOrder) ? 0 : parsedOrder
    const close = () => setFormOpen(false)
    // 颜色均需配置主色、副色与悬停色；炫彩为内置颜色，无需创建
    const colorFields = {
      primary_color: primary,
      sub_color: sub,
      hover_color: hover,
    }

    if (editing) {
      updateColor.mutate(
        {
          id: editing.id,
          req: {
            color_name: name.trim(),
            ...colorFields,
            color_order: colorOrder,
          },
        },
        { onSuccess: close },
      )
    } else {
      createColor.mutate(
        {
          color_name: name.trim(),
          ...colorFields,
          color_order: colorOrder,
        },
        { onSuccess: close },
      )
    }
  }

  const handleDelete = () => {
    if (!deleteTarget) return
    deleteColor.mutate(
      { id: deleteTarget.id, force: forceDelete },
      {
        onSuccess: () => {
          setDeleteTarget(null)
          setForceDelete(false)
        },
      },
    )
  }

  return (
    <div className="mx-auto max-w-6xl space-y-6">
      <PageHead
        kicker="COLORS · 颜色"
        title="颜色管理"
        sub="管理友链的展示颜色，配置普通配色；炫彩为系统内置颜色。"
        actions={
          <Button className="cursor-pointer" onClick={openCreate}>
            <Plus className="mr-2 size-4" />
            新建颜色
          </Button>
        }
      />

      <BambooRule reduced={reduced} delay={0.12} />

      {/* 颜色表格 */}
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
              <TableHead className={inkTh}>预览</TableHead>
              <TableHead className={inkTh}>名称</TableHead>
              <TableHead className={inkTh}>类型</TableHead>
              <TableHead className={inkTh}>排序</TableHead>
              <TableHead className={inkTh}>状态</TableHead>
              <TableHead className={cn(inkTh, 'text-right')}>操作</TableHead>
            </TableRow>
          </TableHeader>
          <TableBody>
            {colorsQuery.isLoading ? (
              Array.from({ length: 5 }).map((_, i) => (
                <TableRow key={i}>
                  <TableCell className={inkTd}>
                    <Skeleton className="size-7 rounded-md" />
                  </TableCell>
                  <TableCell className={inkTd}>
                    <Skeleton className="h-4 w-24" />
                  </TableCell>
                  <TableCell className={inkTd}>
                    <Skeleton className="h-5 w-12 rounded-full" />
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
            ) : colors.length === 0 ? (
              <TableRow className="hover:bg-transparent">
                <TableCell colSpan={6} className="px-4 py-10">
                  <EnsoEmpty
                    title="暂无颜色"
                    hint="创建一种颜色，让友链卡片更有个性"
                  >
                    <Button
                      variant="outline"
                      size="sm"
                      className="ml-auto cursor-pointer"
                      onClick={openCreate}
                    >
                      <Plus className="mr-2 size-4" />
                      新建颜色
                    </Button>
                  </EnsoEmpty>
                </TableCell>
              </TableRow>
            ) : (
              colors.map((color) => (
                <TableRow key={color.id.toString()} className={inkTableRow}>
                  <TableCell className={inkTd}>
                    {isFancyColor(color) ? (
                      <span
                        className="block size-7 rounded-md ring-1 ring-inset ring-border/60"
                        style={{ background: fancyGradient() }}
                        title="炫彩"
                        aria-hidden="true"
                      />
                    ) : (
                      <div className="flex items-center gap-1.5">
                        <span
                          className="block size-7 rounded-md ring-1 ring-inset ring-border/60"
                          style={{
                            backgroundColor:
                              color.primary_color ?? 'var(--leaf-deep)',
                          }}
                          title={`主色 ${color.primary_color ?? '未设置'}`}
                          aria-hidden="true"
                        />
                        <div className="flex flex-col gap-1">
                          {color.sub_color && (
                            <span
                              className="block size-3 rounded-full ring-1 ring-inset ring-border/60"
                              style={{ backgroundColor: color.sub_color }}
                              title={`副色 ${color.sub_color}`}
                              aria-hidden="true"
                            />
                          )}
                          {color.hover_color && (
                            <span
                              className="block size-3 rounded-full ring-1 ring-inset ring-border/60"
                              style={{ backgroundColor: color.hover_color }}
                              title={`悬停色 ${color.hover_color}`}
                              aria-hidden="true"
                            />
                          )}
                        </div>
                      </div>
                    )}
                  </TableCell>
                  <TableCell className={inkTd}>
                    <span className="font-serif font-semibold text-text-primary">
                      {color.name}
                    </span>
                    {!isFancyColor(color) && color.primary_color && (
                      <div className="font-mono text-xs uppercase text-text-secondary">
                        {color.primary_color}
                      </div>
                    )}
                  </TableCell>
                  <TableCell className={inkTd}>
                    {isFancyColor(color) ? (
                      <InkBadge tone="leaf">
                        <Sparkles className="size-3" />
                        炫彩
                      </InkBadge>
                    ) : (
                      <InkBadge tone="neutral">普通</InkBadge>
                    )}
                  </TableCell>
                  <TableCell className={inkTd}>
                    <span className="font-mono tabular-nums text-text-secondary">
                      {color.sort_order}
                    </span>
                  </TableCell>
                  <TableCell className={inkTd}>
                    {color.status ? (
                      <InkBadge tone="leaf">启用</InkBadge>
                    ) : (
                      <InkBadge tone="neutral">禁用</InkBadge>
                    )}
                  </TableCell>
                  <TableCell className={inkTd}>
                    <div className="flex items-center justify-end gap-1">
                      <Button
                        variant="ghost"
                        size="icon"
                        title="编辑"
                        className="size-8 cursor-pointer"
                        onClick={() => openEdit(color)}
                      >
                        <Pencil className="size-4" />
                      </Button>
                      <Button
                        variant="ghost"
                        size="icon"
                        title={color.status ? '禁用' : '启用'}
                        className="size-8 cursor-pointer"
                        disabled={updateStatus.isPending}
                        onClick={() =>
                          updateStatus.mutate({
                            id: color.id,
                            status: !color.status,
                          })
                        }
                      >
                        {color.status ? (
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
                        onClick={() => setDeleteTarget(color)}
                      >
                        <Trash2 className="size-4" />
                      </Button>
                    </div>
                  </TableCell>
                </TableRow>
              ))
            )}
          </TableBody>
        </Table>

        {!colorsQuery.isLoading && colors.length > 0 && (
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

      {/* 新建 / 编辑颜色弹窗 */}
      <Dialog
        open={formOpen}
        onOpenChange={(open) => {
          if (!open) setFormOpen(false)
        }}
      >
        <DialogContent>
          <DialogHeader>
            <DialogTitle>{editing ? '编辑颜色' : '新建颜色'}</DialogTitle>
            <DialogDescription>
              {editing
                ? '修改颜色配置后点击保存生效。'
                : '创建一种友链颜色，配置普通配色。'}
            </DialogDescription>
          </DialogHeader>
          <div className="space-y-4">
            <div className="space-y-2">
              <Label htmlFor="colorName">颜色名称</Label>
              <Input
                id="colorName"
                placeholder="如：天蓝"
                value={name}
                onChange={(e) => setName(e.target.value)}
              />
            </div>
            <div className="space-y-2">
              <div className="flex items-center gap-2 rounded-md border border-input bg-background px-3 py-2">
                <span className="text-sm font-medium text-text-primary">
                  普通配色
                </span>
                <span className="ml-auto text-xs text-text-secondary">
                  炫彩为内置颜色，无需创建
                </span>
              </div>
              <p className="text-xs text-text-secondary">
                需配置主色、副色与悬停色。
              </p>
            </div>
            <div className="grid gap-4 sm:grid-cols-3">
              <ColorField
                id="primaryColor"
                label="主色"
                value={primary}
                onChange={setPrimary}
              />
              <ColorField
                id="subColor"
                label="副色"
                value={sub}
                onChange={setSub}
              />
              <ColorField
                id="hoverColor"
                label="悬停色"
                value={hover}
                onChange={setHover}
              />
            </div>
            <div className="space-y-2">
              <Label htmlFor="colorOrder">排序</Label>
              <Input
                id="colorOrder"
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
          }
        }}
      >
        <DialogContent>
          <DialogHeader>
            <DialogTitle>删除颜色</DialogTitle>
            <DialogDescription>
              确定要删除颜色「{deleteTarget?.name}」吗？此操作不可撤销。
            </DialogDescription>
          </DialogHeader>
          <label className="flex cursor-pointer items-center gap-2.5 rounded-md border border-border bg-muted/30 px-3 py-2.5 text-sm transition-colors duration-150 hover:bg-muted/50">
            <Checkbox
              checked={forceDelete}
              onCheckedChange={(checked) => setForceDelete(checked === true)}
            />
            <span>强制删除（使用该颜色的友链将恢复默认颜色）</span>
          </label>
          <DialogFooter>
            <Button
              variant="outline"
              onClick={() => setDeleteTarget(null)}
              disabled={deleteColor.isPending}
            >
              取消
            </Button>
            <Button
              variant="destructive"
              onClick={handleDelete}
              disabled={deleteColor.isPending}
            >
              {deleteColor.isPending ? '删除中…' : '确认删除'}
            </Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>
    </div>
  )
}
