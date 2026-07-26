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
import { ArrowLeft, Check, Save, Unlink } from 'lucide-react'
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
import { Checkbox } from '@/components/ui/checkbox'
import { mockColors, mockLinks, mockLocations } from '@/data/mock/links'

export const Route = createFileRoute('/_admin/admin/link/$id/edit')({
  component: LinkEditPage,
})

function LinkEditPage() {
  const { id } = Route.useParams()
  const navigate = useNavigate()

  const link = mockLinks.find((l) => l.id === Number(id))

  const [formData, setFormData] = useState({
    siteName: link?.siteName ?? '',
    siteUrl: link?.siteUrl ?? '',
    siteLogo: link?.siteLogo ?? '',
    siteDescription: link?.siteDescription ?? '',
    webmasterEmail: link?.webmasterEmail ?? '',
    location: link?.location ?? null,
    color: link?.color ?? null,
    hasAdv: link?.hasAdv ?? false,
  })

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault()
    navigate({ to: '/admin/link' })
  }

  if (!link) {
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
            正在编辑「{link.siteName}」的信息
          </p>
        </div>
      </div>

      <Card>
        <CardHeader>
          <CardTitle>友链信息</CardTitle>
          <CardDescription>修改友链的基本信息，带 * 为必填项</CardDescription>
        </CardHeader>
        <CardContent>
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
                  value={formData.siteName}
                  onChange={(e) =>
                    setFormData({ ...formData, siteName: e.target.value })
                  }
                />
              </div>
              <div className="space-y-2">
                <Label htmlFor="siteUrl">
                  站点地址 <span className="text-destructive">*</span>
                </Label>
                <Input
                  id="siteUrl"
                  placeholder="https://example.com"
                  value={formData.siteUrl}
                  onChange={(e) =>
                    setFormData({ ...formData, siteUrl: e.target.value })
                  }
                />
              </div>
              <div className="space-y-2">
                <Label htmlFor="siteLogo">站点 Logo</Label>
                <Input
                  id="siteLogo"
                  placeholder="https://example.com/logo.png"
                  value={formData.siteLogo}
                  onChange={(e) =>
                    setFormData({ ...formData, siteLogo: e.target.value })
                  }
                />
              </div>
              <div className="space-y-2">
                <Label htmlFor="webmasterEmail">
                  站长邮箱 <span className="text-destructive">*</span>
                </Label>
                <Input
                  id="webmasterEmail"
                  type="email"
                  placeholder="admin@example.com"
                  value={formData.webmasterEmail}
                  onChange={(e) =>
                    setFormData({ ...formData, webmasterEmail: e.target.value })
                  }
                />
              </div>
            </div>

            {/* 站点描述 */}
            <div className="space-y-2">
              <Label htmlFor="siteDescription">站点描述</Label>
              <textarea
                id="siteDescription"
                className="min-h-[110px] w-full resize-y rounded-md border border-input bg-background px-3 py-2 text-sm shadow-xs transition-colors placeholder:text-muted-foreground focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring"
                placeholder="介绍一下这个站点吧…"
                value={formData.siteDescription}
                onChange={(e) =>
                  setFormData({ ...formData, siteDescription: e.target.value })
                }
              />
            </div>

            {/* 位置分类：药丸按钮 */}
            <div className="space-y-2">
              <Label>位置分类</Label>
              <div className="flex flex-wrap gap-2">
                {mockLocations.map((location) => (
                  <button
                    key={location.id}
                    type="button"
                    onClick={() =>
                      setFormData({ ...formData, location: location.id })
                    }
                    className={`cursor-pointer rounded-full px-4 py-1.5 text-sm font-medium transition-colors duration-200 ${
                      formData.location === location.id
                        ? 'bg-primary text-primary-foreground shadow-sm'
                        : 'bg-muted text-muted-foreground hover:bg-muted/80 hover:text-foreground'
                    }`}
                  >
                    {location.name}
                  </button>
                ))}
              </div>
            </div>

            {/* 颜色分类：可视化色块 */}
            <div className="space-y-2">
              <Label>颜色分类</Label>
              <div className="flex flex-wrap gap-3">
                {mockColors.map((color) => {
                  const selected = formData.color === color.id
                  return (
                    <button
                      key={color.id}
                      type="button"
                      title={color.name}
                      onClick={() =>
                        setFormData({ ...formData, color: color.id })
                      }
                      className={`flex size-9 cursor-pointer items-center justify-center rounded-full ring-offset-2 ring-offset-background transition-all duration-200 ${
                        selected
                          ? 'scale-110 ring-2 ring-ring'
                          : 'hover:scale-105'
                      }`}
                      style={{ backgroundColor: color.color }}
                    >
                      {selected && (
                        <Check className="size-4 text-white" strokeWidth={3} />
                      )}
                    </button>
                  )
                })}
              </div>
            </div>

            {/* 包含广告 */}
            <div className="flex items-center gap-2 rounded-lg border border-border/60 bg-muted/30 p-3">
              <Checkbox
                id="hasAdv"
                checked={formData.hasAdv}
                onCheckedChange={(checked) =>
                  setFormData({ ...formData, hasAdv: checked as boolean })
                }
              />
              <Label htmlFor="hasAdv" className="cursor-pointer font-normal">
                该站点包含广告内容
              </Label>
            </div>

            {/* 操作按钮 */}
            <div className="flex justify-end gap-3 border-t border-border/60 pt-6">
              <Link to="/admin/link">
                <Button variant="outline" className="cursor-pointer">
                  取消
                </Button>
              </Link>
              <Button type="submit" className="cursor-pointer">
                <Save className="mr-2 size-4" />
                保存
              </Button>
            </div>
          </form>
        </CardContent>
      </Card>
    </div>
  )
}
