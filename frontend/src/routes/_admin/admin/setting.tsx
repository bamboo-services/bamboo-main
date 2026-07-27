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
import { useEffect, useState } from 'react'
import { Save } from 'lucide-react'
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { Skeleton } from '@/components/ui/skeleton'
import { Textarea } from '@/components/ui/textarea'
import {
  useAbout,
  useSiteInfo,
  useUpdateAbout,
  useUpdateSiteInfo,
} from '@/hooks/use-site-info'

export const Route = createFileRoute('/_admin/admin/setting')({
  component: SettingPage,
})

function SettingPage() {
  return (
    <div className="space-y-6">
      <div>
        <h1 className="text-3xl font-bold tracking-tight">系统设置</h1>
        <p className="text-muted-foreground">管理站点信息与自我介绍</p>
      </div>

      <div className="grid gap-6">
        <SiteInfoCard />
        <AboutCard />
      </div>
    </div>
  )
}

/** 站点设置卡片：站名 / 描述 / 主页介绍 */
function SiteInfoCard() {
  const { data, isLoading } = useSiteInfo()
  const updateSite = useUpdateSiteInfo()

  const [siteName, setSiteName] = useState('')
  const [siteDescription, setSiteDescription] = useState('')
  const [introduction, setIntroduction] = useState('')

  useEffect(() => {
    if (data) {
      setSiteName(data.site_name)
      setSiteDescription(data.site_description)
      setIntroduction(data.introduction)
    }
  }, [data])

  const handleSave = () => {
    updateSite.mutate({
      site_name: siteName,
      site_description: siteDescription,
      introduction,
    })
  }

  return (
    <Card>
      <CardHeader>
        <CardTitle>站点设置</CardTitle>
        <CardDescription>配置站点的名称、描述与主页介绍</CardDescription>
      </CardHeader>
      <CardContent className="space-y-4">
        {isLoading ? (
          <div className="space-y-4">
            <Skeleton className="h-9 w-full" />
            <Skeleton className="h-9 w-full" />
            <Skeleton className="h-20 w-full" />
          </div>
        ) : (
          <>
            <div className="space-y-2">
              <Label htmlFor="siteName">站点名称</Label>
              <Input
                id="siteName"
                value={siteName}
                onChange={(e) => setSiteName(e.target.value)}
              />
            </div>
            <div className="space-y-2">
              <Label htmlFor="siteDescription">站点描述</Label>
              <Input
                id="siteDescription"
                value={siteDescription}
                onChange={(e) => setSiteDescription(e.target.value)}
              />
            </div>
            <div className="space-y-2">
              <Label htmlFor="introduction">主页介绍</Label>
              <Textarea
                id="introduction"
                className="min-h-[80px]"
                value={introduction}
                onChange={(e) => setIntroduction(e.target.value)}
              />
            </div>
            <div className="flex justify-end">
              <Button onClick={handleSave} disabled={updateSite.isPending}>
                <Save className="mr-2 h-4 w-4" />
                {updateSite.isPending ? '保存中…' : '保存站点设置'}
              </Button>
            </div>
          </>
        )}
      </CardContent>
    </Card>
  )
}

/** 自我介绍卡片：Markdown 内容 */
function AboutCard() {
  const { data, isLoading } = useAbout()
  const updateAbout = useUpdateAbout()

  const [content, setContent] = useState('')

  useEffect(() => {
    if (data) setContent(data.content)
  }, [data])

  const handleSave = () => {
    updateAbout.mutate({ content })
  }

  return (
    <Card>
      <CardHeader>
        <CardTitle>自我介绍</CardTitle>
        <CardDescription>Markdown 格式，展示在「关于我」页面</CardDescription>
      </CardHeader>
      <CardContent className="space-y-4">
        {isLoading ? (
          <Skeleton className="h-40 w-full" />
        ) : (
          <>
            <div className="space-y-2">
              <Label htmlFor="aboutContent">介绍内容</Label>
              <Textarea
                id="aboutContent"
                className="min-h-[240px] font-mono text-sm"
                value={content}
                onChange={(e) => setContent(e.target.value)}
              />
            </div>
            <div className="flex justify-end">
              <Button onClick={handleSave} disabled={updateAbout.isPending}>
                <Save className="mr-2 h-4 w-4" />
                {updateAbout.isPending ? '保存中…' : '保存自我介绍'}
              </Button>
            </div>
          </>
        )}
      </CardContent>
    </Card>
  )
}
