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

import {  useEffect, useState } from 'react'
import { Link, createFileRoute, useNavigate } from '@tanstack/react-router'
import { ArrowLeft, Check, Save, Unlink } from 'lucide-react'
import type {FormEvent} from 'react';
import type { LinkColor, SnowflakeID, UpdateLinkRequest } from '@/api/types'
import { Button } from '@/components/ui/button'
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from '@/components/ui/card'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { Skeleton } from '@/components/ui/skeleton'
import { Textarea } from '@/components/ui/textarea'
import { cn } from '@/lib/utils'
import { useAdminLink, useUpdateLink } from '@/hooks/use-links'
import { useAllGroups } from '@/hooks/use-groups'
import { useAllColors } from '@/hooks/use-colors'

export const Route = createFileRoute('/_admin/admin/link/$id/edit')({
  component: LinkEditPage,
})

/** 表单状态（group/color 为 bigint 雪花 ID，未选为 null） */
interface LinkFormState {
  siteName: string
  siteUrl: string
  siteLogo: string
  webmasterEmail: string
  siteDescription: string
  groupId: SnowflakeID | null
  colorId: SnowflakeID | null
  order: string
  applyRemark: string
}

/** 色块背景：type=1（炫彩）且有副色时渲染渐变，否则取主色 */
function colorBackground(color: LinkColor): string {
  if (color.type === 1 && color.sub_color) {
    return `linear-gradient(135deg, ${color.primary_color ?? '#6366f1'}, ${color.sub_color})`
  }
  return color.primary_color ?? '#6366f1'
}

function LinkEditPage() {
  const { id } = Route.useParams()
  const navigate = useNavigate()
  const linkId = BigInt(id)

  const { data: link, isLoading, isError } = useAdminLink(linkId)
  const updateLink = useUpdateLink()
  const groupsQuery = useAllGroups()
  const colorsQuery = useAllColors()

  const groups = groupsQuery.data ?? []
  const colors = colorsQuery.data ?? []

  const [form, setForm] = useState<LinkFormState>({
    siteName: '',
    siteUrl: '',
    siteLogo: '',
    webmasterEmail: '',
    siteDescription: '',
    groupId: null,
    colorId: null,
    order: '',
    applyRemark: '',
  })

  // 详情加载完成后预填表单
  useEffect(() => {
    if (link) {
      setForm({
        siteName: link.name,
        siteUrl: link.url,
        siteLogo: link.avatar ?? '',
        webmasterEmail: link.email ?? '',
        siteDescription: link.description ?? '',
        groupId: link.group_id,
        colorId: link.color_id,
        order: String(link.sort_order),
        applyRemark: link.apply_remark ?? '',
      })
    }
  }, [link])

  const handleSubmit = (e: FormEvent) => {
    e.preventDefault()
    const order = Number.parseInt(form.order, 10)
    const req: UpdateLinkRequest = {
      link_name: form.siteName.trim(),
      link_url: form.siteUrl.trim(),
      link_avatar: form.siteLogo.trim() || undefined,
      link_email: form.webmasterEmail.trim() || undefined,
      link_desc: form.siteDescription.trim() || undefined,
      link_group_id: form.groupId ?? undefined,
      link_color_id: form.colorId ?? undefined,
      link_order: Number.isNaN(order) ? undefined : order,
      link_apply_remark: form.applyRemark.trim() || undefined,
    }
    updateLink.mutate(
      { id: linkId, req },
      { onSuccess: () => navigate({ to: '/admin/link' }) },
    )
  }

  if (isError) {
    return (
      <div className="mx-auto flex max-w-md flex-col items-center justify-center gap-4 py-24 text-center">
        <div className="flex size-16 items-center justify-center rounded-full bg-muted">
          <Unlink className="size-7 text-muted-foreground" />
        </div>
        <div>
          <h2 className="text-xl font-semibold">友链不存在</h2>
          <p className="mt-1 text-sm text-muted-foreground">
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
      {/* 页头 */}
      <div className="flex items-center gap-4">
        <Link to="/admin/link">
          <Button variant="ghost" size="icon" className="cursor-pointer">
            <ArrowLeft className="size-4" />
          </Button>
        </Link>
        <div>
          <h1 className="text-3xl font-bold tracking-tight">编辑友链</h1>
          <p className="mt-1 text-muted-foreground">
            {link ? `正在编辑「${link.name}」的信息` : '正在加载友链信息…'}
          </p>
        </div>
      </div>

      <Card>
        <CardHeader>
          <CardTitle>友链信息</CardTitle>
          <CardDescription>修改友链的基本信息，带 * 为必填项</CardDescription>
        </CardHeader>
        <CardContent>
          {isLoading ? (
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
          ) : (
            <form onSubmit={handleSubmit} className="space-y-6">
              {/* 基本信息 */}
              <div className="grid gap-4 md:grid-cols-2">
                <div className="space-y-2">
                  <Label htmlFor="siteName">
                    站点名称 <span className="text-destructive">*</span>
                  </Label>
                  <Input
                    id="siteName"
                    placeholder="请输入站点名称"
                    required
                    value={form.siteName}
                    onChange={(e) =>
                      setForm({ ...form, siteName: e.target.value })
                    }
                  />
                </div>
                <div className="space-y-2">
                  <Label htmlFor="siteUrl">
                    站点地址 <span className="text-destructive">*</span>
                  </Label>
                  <Input
                    id="siteUrl"
                    type="url"
                    placeholder="https://example.com"
                    required
                    value={form.siteUrl}
                    onChange={(e) =>
                      setForm({ ...form, siteUrl: e.target.value })
                    }
                  />
                </div>
                <div className="space-y-2">
                  <Label htmlFor="siteLogo">站点 Logo</Label>
                  <Input
                    id="siteLogo"
                    placeholder="https://example.com/logo.png"
                    value={form.siteLogo}
                    onChange={(e) =>
                      setForm({ ...form, siteLogo: e.target.value })
                    }
                  />
                </div>
                <div className="space-y-2">
                  <Label htmlFor="webmasterEmail">站长邮箱</Label>
                  <Input
                    id="webmasterEmail"
                    type="email"
                    placeholder="admin@example.com"
                    value={form.webmasterEmail}
                    onChange={(e) =>
                      setForm({ ...form, webmasterEmail: e.target.value })
                    }
                  />
                </div>
              </div>

              {/* 站点描述 */}
              <div className="space-y-2">
                <Label htmlFor="siteDescription">站点描述</Label>
                <Textarea
                  id="siteDescription"
                  className="min-h-[110px]"
                  placeholder="介绍一下这个站点吧…"
                  value={form.siteDescription}
                  onChange={(e) =>
                    setForm({ ...form, siteDescription: e.target.value })
                  }
                />
              </div>

              {/* 排序 + 申请备注 */}
              <div className="grid gap-4 md:grid-cols-2">
                <div className="space-y-2">
                  <Label htmlFor="linkOrder">排序</Label>
                  <Input
                    id="linkOrder"
                    type="number"
                    step={1}
                    placeholder="数值越小越靠前，默认 0"
                    value={form.order}
                    onChange={(e) =>
                      setForm({ ...form, order: e.target.value })
                    }
                  />
                </div>
                <div className="space-y-2">
                  <Label htmlFor="applyRemark">申请备注</Label>
                  <Input
                    id="applyRemark"
                    placeholder="选填，例如申请来源说明"
                    value={form.applyRemark}
                    onChange={(e) =>
                      setForm({ ...form, applyRemark: e.target.value })
                    }
                  />
                </div>
              </div>

              {/* 位置分类：药丸按钮 */}
              <div className="space-y-2">
                <div className="flex items-baseline justify-between">
                  <Label>位置分类</Label>
                  <span className="text-xs text-muted-foreground">
                    再次点击可取消选择
                  </span>
                </div>
                {groupsQuery.isLoading ? (
                  <div className="flex gap-2">
                    {Array.from({ length: 4 }).map((_, i) => (
                      <Skeleton key={i} className="h-8 w-20 rounded-full" />
                    ))}
                  </div>
                ) : (
                  <div className="flex flex-wrap gap-2">
                    <button
                      type="button"
                      onClick={() => setForm({ ...form, groupId: null })}
                      className={cn(
                        'cursor-pointer rounded-full px-4 py-1.5 text-sm font-medium transition-colors duration-200',
                        form.groupId === null
                          ? 'bg-primary text-primary-foreground shadow-sm'
                          : 'bg-muted text-muted-foreground hover:bg-muted/80 hover:text-foreground',
                      )}
                    >
                      未分组
                    </button>
                    {groups.map((group) => (
                      <button
                        key={group.id.toString()}
                        type="button"
                        title={group.description ?? undefined}
                        onClick={() =>
                          setForm({
                            ...form,
                            groupId:
                              form.groupId === group.id ? null : group.id,
                          })
                        }
                        className={cn(
                          'cursor-pointer rounded-full px-4 py-1.5 text-sm font-medium transition-colors duration-200',
                          form.groupId === group.id
                            ? 'bg-primary text-primary-foreground shadow-sm'
                            : 'bg-muted text-muted-foreground hover:bg-muted/80 hover:text-foreground',
                        )}
                      >
                        {group.name}
                      </button>
                    ))}
                    {groups.length === 0 && (
                      <span className="self-center text-xs text-muted-foreground">
                        暂无分组，可先前往分组管理创建
                      </span>
                    )}
                  </div>
                )}
              </div>

              {/* 颜色分类：可视化色块 */}
              <div className="space-y-2">
                <div className="flex items-baseline justify-between">
                  <Label>颜色分类</Label>
                  <span className="text-xs text-muted-foreground">
                    再次点击可恢复默认
                  </span>
                </div>
                {colorsQuery.isLoading ? (
                  <div className="flex gap-3">
                    {Array.from({ length: 6 }).map((_, i) => (
                      <Skeleton key={i} className="size-9 rounded-full" />
                    ))}
                  </div>
                ) : (
                  <div className="flex flex-wrap items-center gap-3">
                    <button
                      type="button"
                      title="默认颜色"
                      onClick={() => setForm({ ...form, colorId: null })}
                      className={cn(
                        'flex size-9 cursor-pointer items-center justify-center rounded-full border-2 border-dashed border-border bg-muted/40 transition-all duration-200',
                        form.colorId === null
                          ? 'scale-110 ring-2 ring-ring ring-offset-2 ring-offset-background'
                          : 'hover:scale-105',
                      )}
                    >
                      {form.colorId === null && (
                        <Check
                          className="size-4 text-muted-foreground"
                          strokeWidth={3}
                        />
                      )}
                    </button>
                    {colors.map((color) => {
                      const selected = form.colorId === color.id
                      return (
                        <button
                          key={color.id.toString()}
                          type="button"
                          title={color.name}
                          onClick={() =>
                            setForm({
                              ...form,
                              colorId: selected ? null : color.id,
                            })
                          }
                          className={cn(
                            'flex size-9 cursor-pointer items-center justify-center rounded-full ring-offset-2 ring-offset-background transition-all duration-200',
                            selected
                              ? 'scale-110 ring-2 ring-ring'
                              : 'hover:scale-105',
                          )}
                          style={{ background: colorBackground(color) }}
                        >
                          {selected && (
                            <Check
                              className="size-4 text-white"
                              strokeWidth={3}
                            />
                          )}
                        </button>
                      )
                    })}
                    {colors.length === 0 && (
                      <span className="text-xs text-muted-foreground">
                        暂无颜色，将使用默认颜色
                      </span>
                    )}
                  </div>
                )}
              </div>

              {/* 操作按钮 */}
              <div className="flex justify-end gap-3 border-t border-border/60 pt-6">
                <Link to="/admin/link">
                  <Button variant="outline" className="cursor-pointer">
                    取消
                  </Button>
                </Link>
                <Button
                  type="submit"
                  className="cursor-pointer"
                  disabled={updateLink.isPending}
                >
                  <Save className="mr-2 size-4" />
                  {updateLink.isPending ? '保存中…' : '保存'}
                </Button>
              </div>
            </form>
          )}
        </CardContent>
      </Card>
    </div>
  )
}
