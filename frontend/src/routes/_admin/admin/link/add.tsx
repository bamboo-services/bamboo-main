/*
 * --------------------------------------------------------------------------------
 * Copyright (c) 2016-NOW 筱锋
 * Author: 筱锋(https://www.x-lf.com)
 * License: MIT
 * --------------------------------------------------------------------------------
 */

import { createFileRoute, Link, useNavigate } from '@tanstack/react-router'
import { Button } from '@/components/ui/button'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { Checkbox } from '@/components/ui/checkbox'
import { mockLocations, mockColors } from '@/data/mock/links'
import { ArrowLeft, Save } from 'lucide-react'
import { useState } from 'react'

export const Route = createFileRoute('/_admin/admin/link/add')({
  component: LinkAddPage,
})

function LinkAddPage() {
  const navigate = useNavigate()
  const [formData, setFormData] = useState({
    siteName: '',
    siteUrl: '',
    siteLogo: '',
    siteDescription: '',
    webmasterEmail: '',
    location: '',
    color: '',
    hasAdv: false,
  })

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault()
    // 静态界面，直接返回列表
    navigate({ to: '/admin/link' })
  }

  return (
    <div className="space-y-6">
      <div className="flex items-center gap-4">
        <Link to="/admin/link">
          <Button variant="ghost" size="icon">
            <ArrowLeft className="h-4 w-4" />
          </Button>
        </Link>
        <div>
          <h1 className="text-3xl font-bold tracking-tight">添加友链</h1>
          <p className="text-muted-foreground">添加一个新的友情链接</p>
        </div>
      </div>

      <Card>
        <CardHeader>
          <CardTitle>友链信息</CardTitle>
          <CardDescription>填写友链的基本信息</CardDescription>
        </CardHeader>
        <CardContent>
          <form onSubmit={handleSubmit} className="space-y-6">
            <div className="grid gap-4 md:grid-cols-2">
              <div className="space-y-2">
                <Label htmlFor="siteName">
                  站点名称 <span className="text-red-500">*</span>
                </Label>
                <Input
                  id="siteName"
                  placeholder="请输入站点名称"
                  value={formData.siteName}
                  onChange={(e) => setFormData({ ...formData, siteName: e.target.value })}
                />
              </div>
              <div className="space-y-2">
                <Label htmlFor="siteUrl">
                  站点地址 <span className="text-red-500">*</span>
                </Label>
                <Input
                  id="siteUrl"
                  placeholder="https://example.com"
                  value={formData.siteUrl}
                  onChange={(e) => setFormData({ ...formData, siteUrl: e.target.value })}
                />
              </div>
              <div className="space-y-2">
                <Label htmlFor="siteLogo">站点 Logo</Label>
                <Input
                  id="siteLogo"
                  placeholder="https://example.com/logo.png"
                  value={formData.siteLogo}
                  onChange={(e) => setFormData({ ...formData, siteLogo: e.target.value })}
                />
              </div>
              <div className="space-y-2">
                <Label htmlFor="webmasterEmail">
                  站长邮箱 <span className="text-red-500">*</span>
                </Label>
                <Input
                  id="webmasterEmail"
                  type="email"
                  placeholder="admin@example.com"
                  value={formData.webmasterEmail}
                  onChange={(e) => setFormData({ ...formData, webmasterEmail: e.target.value })}
                />
              </div>
              <div className="space-y-2">
                <Label htmlFor="location">位置分类</Label>
                <select
                  id="location"
                  className="w-full rounded-md border border-input bg-background px-3 py-2 text-sm"
                  value={formData.location}
                  onChange={(e) => setFormData({ ...formData, location: e.target.value })}
                >
                  <option value="">请选择位置</option>
                  {mockLocations.map((location) => (
                    <option key={location.id} value={location.id}>
                      {location.name}
                    </option>
                  ))}
                </select>
              </div>
              <div className="space-y-2">
                <Label htmlFor="color">颜色分类</Label>
                <select
                  id="color"
                  className="w-full rounded-md border border-input bg-background px-3 py-2 text-sm"
                  value={formData.color}
                  onChange={(e) => setFormData({ ...formData, color: e.target.value })}
                >
                  <option value="">请选择颜色</option>
                  {mockColors.map((color) => (
                    <option key={color.id} value={color.id}>
                      {color.name}
                    </option>
                  ))}
                </select>
              </div>
            </div>

            <div className="space-y-2">
              <Label htmlFor="siteDescription">站点描述</Label>
              <textarea
                id="siteDescription"
                className="w-full min-h-[100px] rounded-md border border-input bg-background px-3 py-2 text-sm"
                placeholder="请输入站点描述"
                value={formData.siteDescription}
                onChange={(e) => setFormData({ ...formData, siteDescription: e.target.value })}
              />
            </div>

            <div className="flex items-center gap-2">
              <Checkbox
                id="hasAdv"
                checked={formData.hasAdv}
                onCheckedChange={(checked) => setFormData({ ...formData, hasAdv: checked as boolean })}
              />
              <Label htmlFor="hasAdv" className="font-normal">
                包含广告
              </Label>
            </div>

            <div className="flex justify-end gap-4">
              <Link to="/admin/link">
                <Button variant="outline">取消</Button>
              </Link>
              <Button type="submit">
                <Save className="mr-2 h-4 w-4" />
                保存
              </Button>
            </div>
          </form>
        </CardContent>
      </Card>
    </div>
  )
}
