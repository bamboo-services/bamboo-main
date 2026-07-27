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

import { Link, createFileRoute } from '@tanstack/react-router'
import {
  Activity,
  CheckCircle,
  Clock,
  Cpu,
  Link as LinkIcon,
  MemoryStick,
  Timer,
} from 'lucide-react'
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from '@/components/ui/card'
import { Avatar, AvatarFallback, AvatarImage } from '@/components/ui/avatar'
import { Skeleton } from '@/components/ui/skeleton'
import { useDashboardStats, useHealth } from '@/hooks/use-dashboard'

export const Route = createFileRoute('/_admin/admin/dashboard')({
  component: DashboardPage,
})

/** 相对时间（用于最近申请） */
function relativeTime(iso: string): string {
  const diff = Date.now() - new Date(iso).getTime()
  const minutes = Math.floor(diff / 60_000)
  if (minutes < 1) return '刚刚'
  if (minutes < 60) return `${minutes} 分钟前`
  const hours = Math.floor(minutes / 60)
  if (hours < 24) return `${hours} 小时前`
  const days = Math.floor(hours / 24)
  return `${days} 天前`
}

function DashboardPage() {
  const statsQuery = useDashboardStats()
  const healthQuery = useHealth()
  const stats = statsQuery.data

  const statCards = [
    {
      title: '友链总数',
      value: stats?.total_links,
      icon: LinkIcon,
      description: '已收录的友链数量',
    },
    {
      title: '待审核',
      value: stats?.pending_links,
      icon: Clock,
      description: '等待审核的申请',
    },
    {
      title: '已通过',
      value: stats?.approved_links,
      icon: CheckCircle,
      description: '审核通过的友链',
    },
  ]

  const runtime = healthQuery.data?.runtime

  return (
    <div className="space-y-6">
      <div>
        <h1 className="text-3xl font-bold tracking-tight">看板</h1>
        <p className="text-muted-foreground">欢迎回来，这里是系统概览。</p>
      </div>

      <div className="grid gap-4 md:grid-cols-3">
        {statCards.map((stat) => (
          <Card key={stat.title}>
            <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
              <CardTitle className="text-sm font-medium">
                {stat.title}
              </CardTitle>
              <stat.icon className="h-4 w-4 text-muted-foreground" />
            </CardHeader>
            <CardContent>
              {statsQuery.isLoading ? (
                <Skeleton className="h-8 w-16" />
              ) : (
                <div className="text-2xl font-bold">{stat.value ?? 0}</div>
              )}
              <p className="text-xs text-muted-foreground">
                {stat.description}
              </p>
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
            {statsQuery.isLoading ? (
              <div className="space-y-4">
                {Array.from({ length: 3 }).map((_, i) => (
                  <Skeleton key={i} className="h-10 w-full" />
                ))}
              </div>
            ) : stats && stats.recent_applications.length > 0 ? (
              <div className="space-y-4">
                {stats.recent_applications.map((item) => (
                  <Link
                    key={item.id.toString()}
                    to="/admin/link/verify"
                    className="flex items-center gap-4 rounded-lg p-1 transition-colors hover:bg-muted/50"
                  >
                    <Avatar className="h-10 w-10">
                      <AvatarImage
                        src={item.avatar || undefined}
                        alt={item.name}
                      />
                      <AvatarFallback>{item.name.slice(0, 1)}</AvatarFallback>
                    </Avatar>
                    <div className="flex-1">
                      <p className="text-sm font-medium">{item.name}</p>
                      <p className="text-xs text-muted-foreground">
                        {relativeTime(item.created_at)}
                      </p>
                    </div>
                  </Link>
                ))}
              </div>
            ) : (
              <p className="py-4 text-center text-sm text-muted-foreground">
                暂无待审核申请
              </p>
            )}
          </CardContent>
        </Card>

        <Card>
          <CardHeader>
            <CardTitle className="flex items-center gap-2">
              <Activity className="h-4 w-4" />
              系统状态
            </CardTitle>
            <CardDescription>当前系统运行时指标</CardDescription>
          </CardHeader>
          <CardContent>
            {healthQuery.isLoading ? (
              <div className="space-y-4">
                {Array.from({ length: 4 }).map((_, i) => (
                  <Skeleton key={i} className="h-5 w-full" />
                ))}
              </div>
            ) : runtime ? (
              <div className="space-y-4">
                <RuntimeRow
                  icon={Timer}
                  label="运行时长"
                  value={runtime.uptime}
                />
                <RuntimeRow
                  icon={MemoryStick}
                  label="内存使用"
                  value={runtime.memory_usage}
                />
                <RuntimeRow
                  icon={Cpu}
                  label="CPU 使用率"
                  value={runtime.cpu_usage}
                />
                <RuntimeRow
                  icon={Activity}
                  label="协程数量"
                  value={String(runtime.goroutines)}
                />
              </div>
            ) : (
              <p className="py-4 text-center text-sm text-muted-foreground">
                暂无运行时数据
              </p>
            )}
          </CardContent>
        </Card>
      </div>
    </div>
  )
}

function RuntimeRow({
  icon: Icon,
  label,
  value,
}: {
  icon: typeof Timer
  label: string
  value: string
}) {
  return (
    <div className="flex items-center justify-between">
      <span className="flex items-center gap-2 text-sm text-muted-foreground">
        <Icon className="h-4 w-4" />
        {label}
      </span>
      <span className="text-sm font-medium">{value}</span>
    </div>
  )
}
