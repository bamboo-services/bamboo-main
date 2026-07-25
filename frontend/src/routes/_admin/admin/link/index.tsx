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

import { createFileRoute, Link } from '@tanstack/react-router'
import { Button } from '@/components/ui/button'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { Input } from '@/components/ui/input'
import { Badge } from '@/components/ui/badge'
import { mockLinks, mockLocations } from '@/data/mock/links'
import { Plus, Search, CheckCircle, ExternalLink, MoreHorizontal } from 'lucide-react'
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger,
} from '@/components/ui/dropdown-menu'

export const Route = createFileRoute('/_admin/admin/link/')({
  component: LinkListPage,
})

function LinkListPage() {
  const approvedLinks = mockLinks.filter((link) => link.status === 'approved')
  const pendingLinks = mockLinks.filter((link) => link.status === 'pending')

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-3xl font-bold tracking-tight">友链管理</h1>
          <p className="text-muted-foreground">管理所有友情链接</p>
        </div>
        <div className="flex gap-2">
          <Link to="/admin/link/verify">
            <Button variant="outline">
              <CheckCircle className="mr-2 h-4 w-4" />
              友链审核
              {pendingLinks.length > 0 && (
                <Badge variant="destructive" className="ml-2">
                  {pendingLinks.length}
                </Badge>
              )}
            </Button>
          </Link>
          <Link to="/admin/link/add">
            <Button>
              <Plus className="mr-2 h-4 w-4" />
              添加友链
            </Button>
          </Link>
        </div>
      </div>

      {/* 筛选区域 */}
      <div className="flex flex-col gap-4 sm:flex-row">
        <div className="relative flex-1">
          <Search className="absolute left-3 top-1/2 h-4 w-4 -translate-y-1/2 text-muted-foreground" />
          <Input placeholder="搜索友链..." className="pl-10" />
        </div>
        <select className="rounded-md border border-input bg-background px-3 py-2 text-sm">
          <option value="">全部位置</option>
          {mockLocations.map((location) => (
            <option key={location.id} value={location.id}>
              {location.name}
            </option>
          ))}
        </select>
      </div>

      {/* 统计卡片 */}
      <div className="grid gap-4 md:grid-cols-3">
        <Card>
          <CardHeader className="pb-2">
            <CardTitle className="text-sm font-medium">总计</CardTitle>
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold">{mockLinks.length}</div>
          </CardContent>
        </Card>
        <Card>
          <CardHeader className="pb-2">
            <CardTitle className="text-sm font-medium">已通过</CardTitle>
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold text-green-600">{approvedLinks.length}</div>
          </CardContent>
        </Card>
        <Card>
          <CardHeader className="pb-2">
            <CardTitle className="text-sm font-medium">待审核</CardTitle>
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold text-yellow-600">{pendingLinks.length}</div>
          </CardContent>
        </Card>
      </div>

      {/* 友链列表 */}
      <div className="grid gap-4 md:grid-cols-2">
        {approvedLinks.map((link) => (
          <Card key={link.id}>
            <CardHeader className="flex flex-row items-start justify-between space-y-0">
              <div className="flex items-center gap-3">
                <div className="h-12 w-12 rounded-lg bg-muted flex items-center justify-center overflow-hidden">
                  <img
                    src={link.siteLogo}
                    alt={link.siteName}
                    className="h-full w-full object-cover"
                    onError={(e) => {
                      e.currentTarget.style.display = 'none'
                    }}
                  />
                </div>
                <div>
                  <CardTitle className="text-base">{link.siteName}</CardTitle>
                  <CardDescription className="flex items-center gap-1">
                    <a
                      href={link.siteUrl}
                      target="_blank"
                      rel="noopener noreferrer"
                      className="hover:underline flex items-center gap-1"
                    >
                      {link.siteUrl}
                      <ExternalLink className="h-3 w-3" />
                    </a>
                  </CardDescription>
                </div>
              </div>
              <DropdownMenu>
                <DropdownMenuTrigger asChild>
                  <Button variant="ghost" size="icon">
                    <MoreHorizontal className="h-4 w-4" />
                  </Button>
                </DropdownMenuTrigger>
                <DropdownMenuContent align="end">
                  <DropdownMenuItem asChild>
                    <Link to="/admin/link/$id/edit" params={{ id: String(link.id) }}>
                      编辑
                    </Link>
                  </DropdownMenuItem>
                  <DropdownMenuItem className="text-destructive">删除</DropdownMenuItem>
                </DropdownMenuContent>
              </DropdownMenu>
            </CardHeader>
            <CardContent>
              <p className="text-sm text-muted-foreground line-clamp-2">{link.siteDescription}</p>
              <div className="mt-3 flex items-center gap-2">
                <Badge variant="secondary">{link.locationName}</Badge>
                <Badge
                  variant="outline"
                  style={{ borderColor: mockLinks.find((l) => l.id === link.id)?.colorName }}
                >
                  {link.colorName}
                </Badge>
                {link.ableConnect ? (
                  <Badge variant="default" className="bg-green-500">
                    可访问
                  </Badge>
                ) : (
                  <Badge variant="destructive">不可访问</Badge>
                )}
              </div>
            </CardContent>
          </Card>
        ))}
      </div>
    </div>
  )
}
