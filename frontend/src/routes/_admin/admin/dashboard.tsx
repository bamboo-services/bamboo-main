/*
 * --------------------------------------------------------------------------------
 * Copyright (c) 2016-NOW 筱锋
 * Author: 筱锋(https://www.x-lf.com)
 * License: MIT
 * --------------------------------------------------------------------------------
 */

import { createFileRoute } from '@tanstack/react-router'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { Link as LinkIcon, Users, CheckCircle, Clock } from 'lucide-react'

export const Route = createFileRoute('/_admin/admin/dashboard')({
  component: DashboardPage,
})

function DashboardPage() {
  const stats = [
    { title: '友链总数', value: '24', icon: LinkIcon, description: '已收录的友链数量' },
    { title: '待审核', value: '3', icon: Clock, description: '等待审核的申请' },
    { title: '已通过', value: '20', icon: CheckCircle, description: '审核通过的友链' },
    { title: '访问量', value: '1,234', icon: Users, description: '本月访问量' },
  ]

  return (
    <div className="space-y-6">
      <div>
        <h1 className="text-3xl font-bold tracking-tight">看板</h1>
        <p className="text-muted-foreground">欢迎回来，这里是系统概览。</p>
      </div>

      <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-4">
        {stats.map((stat) => (
          <Card key={stat.title}>
            <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
              <CardTitle className="text-sm font-medium">{stat.title}</CardTitle>
              <stat.icon className="h-4 w-4 text-muted-foreground" />
            </CardHeader>
            <CardContent>
              <div className="text-2xl font-bold">{stat.value}</div>
              <p className="text-xs text-muted-foreground">{stat.description}</p>
            </CardContent>
          </Card>
        ))}
      </div>

      <div className="grid gap-4 md:grid-cols-2">
        <Card>
          <CardHeader>
            <CardTitle>最近申请</CardTitle>
            <CardDescription>最近收到的友链申请</CardDescription>
          </CardHeader>
          <CardContent>
            <div className="space-y-4">
              <div className="flex items-center gap-4">
                <div className="h-10 w-10 rounded-full bg-muted" />
                <div className="flex-1">
                  <p className="text-sm font-medium">示例站点</p>
                  <p className="text-xs text-muted-foreground">2 小时前</p>
                </div>
              </div>
              <div className="flex items-center gap-4">
                <div className="h-10 w-10 rounded-full bg-muted" />
                <div className="flex-1">
                  <p className="text-sm font-medium">技术博客</p>
                  <p className="text-xs text-muted-foreground">5 小时前</p>
                </div>
              </div>
              <div className="flex items-center gap-4">
                <div className="h-10 w-10 rounded-full bg-muted" />
                <div className="flex-1">
                  <p className="text-sm font-medium">个人主页</p>
                  <p className="text-xs text-muted-foreground">1 天前</p>
                </div>
              </div>
            </div>
          </CardContent>
        </Card>

        <Card>
          <CardHeader>
            <CardTitle>系统状态</CardTitle>
            <CardDescription>当前系统运行状态</CardDescription>
          </CardHeader>
          <CardContent>
            <div className="space-y-4">
              <div className="flex items-center justify-between">
                <span className="text-sm">API 服务</span>
                <span className="flex items-center gap-1 text-sm text-green-500">
                  <span className="h-2 w-2 rounded-full bg-green-500" />
                  正常
                </span>
              </div>
              <div className="flex items-center justify-between">
                <span className="text-sm">数据库</span>
                <span className="flex items-center gap-1 text-sm text-green-500">
                  <span className="h-2 w-2 rounded-full bg-green-500" />
                  正常
                </span>
              </div>
              <div className="flex items-center justify-between">
                <span className="text-sm">缓存服务</span>
                <span className="flex items-center gap-1 text-sm text-green-500">
                  <span className="h-2 w-2 rounded-full bg-green-500" />
                  正常
                </span>
              </div>
            </div>
          </CardContent>
        </Card>
      </div>
    </div>
  )
}
