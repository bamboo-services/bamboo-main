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

import { createFileRoute } from '@tanstack/react-router'
import { Plus } from 'lucide-react'
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import { Badge } from '@/components/ui/badge'

export const Route = createFileRoute('/_admin/admin/sponsor')({
  component: SponsorPage,
})

const mockSponsors = [
  {
    id: 1,
    name: '张三',
    amount: 50,
    method: '微信',
    message: '支持一下！',
    createdAt: '2024-12-10',
  },
  {
    id: 2,
    name: '李四',
    amount: 100,
    method: '支付宝',
    message: '加油！',
    createdAt: '2024-12-08',
  },
  {
    id: 3,
    name: '匿名',
    amount: 20,
    method: '微信',
    message: '',
    createdAt: '2024-12-05',
  },
]

function SponsorPage() {
  const totalAmount = mockSponsors.reduce((sum, s) => sum + s.amount, 0)

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-3xl font-bold tracking-tight">赞助管理</h1>
          <p className="text-muted-foreground">管理所有赞助记录</p>
        </div>
        <Button>
          <Plus className="mr-2 h-4 w-4" />
          添加赞助
        </Button>
      </div>

      {/* 统计卡片 */}
      <div className="grid gap-4 md:grid-cols-3">
        <Card>
          <CardHeader className="pb-2">
            <CardTitle className="text-sm font-medium">赞助总数</CardTitle>
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold">{mockSponsors.length}</div>
          </CardContent>
        </Card>
        <Card>
          <CardHeader className="pb-2">
            <CardTitle className="text-sm font-medium">赞助总额</CardTitle>
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold">¥{totalAmount}</div>
          </CardContent>
        </Card>
        <Card>
          <CardHeader className="pb-2">
            <CardTitle className="text-sm font-medium">本月赞助</CardTitle>
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold">¥{totalAmount}</div>
          </CardContent>
        </Card>
      </div>

      {/* 赞助列表 */}
      <Card>
        <CardHeader>
          <CardTitle>赞助记录</CardTitle>
          <CardDescription>所有赞助记录列表</CardDescription>
        </CardHeader>
        <CardContent>
          <div className="relative w-full overflow-auto">
            <table className="w-full caption-bottom text-sm">
              <thead className="[&_tr]:border-b">
                <tr className="border-b transition-colors hover:bg-muted/50">
                  <th className="h-12 px-4 text-left align-middle font-medium">
                    赞助者
                  </th>
                  <th className="h-12 px-4 text-left align-middle font-medium">
                    金额
                  </th>
                  <th className="h-12 px-4 text-left align-middle font-medium">
                    方式
                  </th>
                  <th className="h-12 px-4 text-left align-middle font-medium">
                    留言
                  </th>
                  <th className="h-12 px-4 text-left align-middle font-medium">
                    时间
                  </th>
                </tr>
              </thead>
              <tbody className="[&_tr:last-child]:border-0">
                {mockSponsors.map((sponsor) => (
                  <tr
                    key={sponsor.id}
                    className="border-b transition-colors hover:bg-muted/50"
                  >
                    <td className="p-4 align-middle">{sponsor.name}</td>
                    <td className="p-4 align-middle font-medium text-green-600">
                      ¥{sponsor.amount}
                    </td>
                    <td className="p-4 align-middle">
                      <Badge variant="secondary">{sponsor.method}</Badge>
                    </td>
                    <td className="p-4 align-middle text-muted-foreground">
                      {sponsor.message || '-'}
                    </td>
                    <td className="p-4 align-middle text-muted-foreground">
                      {sponsor.createdAt}
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        </CardContent>
      </Card>
    </div>
  )
}
