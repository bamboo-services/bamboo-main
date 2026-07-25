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
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { Save } from 'lucide-react'
import { siteInfo } from '@/data/mock/site-info'

export const Route = createFileRoute('/_admin/admin/setting')({
  component: SettingPage,
})

function SettingPage() {
  return (
    <div className="space-y-6">
      <div>
        <h1 className="text-3xl font-bold tracking-tight">系统设置</h1>
        <p className="text-muted-foreground">管理系统配置</p>
      </div>

      <div className="grid gap-6">
        {/* 站点设置 */}
        <Card>
          <CardHeader>
            <CardTitle>站点设置</CardTitle>
            <CardDescription>配置站点的基本信息</CardDescription>
          </CardHeader>
          <CardContent className="space-y-4">
            <div className="grid gap-4 md:grid-cols-2">
              <div className="space-y-2">
                <Label htmlFor="siteName">站点名称</Label>
                <Input id="siteName" defaultValue={siteInfo.site.siteName} />
              </div>
              <div className="space-y-2">
                <Label htmlFor="siteAuthor">站点作者</Label>
                <Input id="siteAuthor" defaultValue={siteInfo.site.author} />
              </div>
            </div>
            <div className="space-y-2">
              <Label htmlFor="siteDescription">站点描述</Label>
              <textarea
                id="siteDescription"
                className="w-full min-h-[80px] rounded-md border border-input bg-background px-3 py-2 text-sm"
                defaultValue={siteInfo.site.description}
              />
            </div>
            <div className="space-y-2">
              <Label htmlFor="siteKeywords">站点关键词</Label>
              <Input id="siteKeywords" defaultValue={siteInfo.site.keywords} />
            </div>
          </CardContent>
        </Card>

        {/* 博主设置 */}
        <Card>
          <CardHeader>
            <CardTitle>博主设置</CardTitle>
            <CardDescription>配置博主的个人信息</CardDescription>
          </CardHeader>
          <CardContent className="space-y-4">
            <div className="grid gap-4 md:grid-cols-2">
              <div className="space-y-2">
                <Label htmlFor="bloggerName">博主名称</Label>
                <Input id="bloggerName" defaultValue={siteInfo.blogger.name} />
              </div>
              <div className="space-y-2">
                <Label htmlFor="bloggerNick">博主昵称</Label>
                <Input id="bloggerNick" defaultValue={siteInfo.blogger.nick} />
              </div>
            </div>
            <div className="space-y-2">
              <Label htmlFor="bloggerEmail">博主邮箱</Label>
              <Input id="bloggerEmail" type="email" defaultValue={siteInfo.blogger.email} />
            </div>
            <div className="space-y-2">
              <Label htmlFor="bloggerDescription">博主简介</Label>
              <textarea
                id="bloggerDescription"
                className="w-full min-h-[80px] rounded-md border border-input bg-background px-3 py-2 text-sm"
                defaultValue={siteInfo.blogger.description}
              />
            </div>
          </CardContent>
        </Card>

        <div className="flex justify-end">
          <Button>
            <Save className="mr-2 h-4 w-4" />
            保存设置
          </Button>
        </div>
      </div>
    </div>
  )
}
