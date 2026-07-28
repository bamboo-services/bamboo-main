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

import { createFileRoute } from '@tanstack/react-router'
import { useEffect, useState } from 'react'
import { motion, useReducedMotion } from 'motion/react'
import { Save } from 'lucide-react'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { Skeleton } from '@/components/ui/skeleton'
import { Textarea } from '@/components/ui/textarea'
import {
  BambooRule,
  CardHead,
  PageHead,
  inkCard,
} from '@/components/ink-wash'
import { enter } from '@/lib/motion'
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
  const reduced = useReducedMotion() ?? false

  return (
    <div className="mx-auto max-w-3xl space-y-6">
      <PageHead
        kicker="SETTINGS · 设置"
        title="系统设置"
        sub="管理站点信息与自我介绍。"
      />

      <BambooRule reduced={reduced} delay={0.12} />

      <div className="grid gap-6">
        <SiteInfoCard reduced={reduced} />
        <AboutCard reduced={reduced} />
      </div>
    </div>
  )
}

/** 站点设置卡片：站名 / 描述 / 主页介绍 */
function SiteInfoCard({ reduced }: { reduced: boolean }) {
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
    <motion.section
      {...enter(reduced, 0.18, {
        initial: { opacity: 0, y: 10 },
        animate: { opacity: 1, y: 0 },
        transition: { duration: 0.4, ease: 'easeOut' },
      })}
      className={inkCard}
    >
      <CardHead title="站点设置" meta="SITE INFO" />
      {isLoading ? (
        <div className="space-y-4">
          <Skeleton className="h-9 w-full" />
          <Skeleton className="h-9 w-full" />
          <Skeleton className="h-20 w-full" />
        </div>
      ) : (
        <div className="space-y-4">
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
        </div>
      )}
    </motion.section>
  )
}

/** 自我介绍卡片：Markdown 内容 */
function AboutCard({ reduced }: { reduced: boolean }) {
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
    <motion.section
      {...enter(reduced, 0.24, {
        initial: { opacity: 0, y: 10 },
        animate: { opacity: 1, y: 0 },
        transition: { duration: 0.4, ease: 'easeOut' },
      })}
      className={inkCard}
    >
      <CardHead title="自我介绍" meta="ABOUT · MARKDOWN" />
      {isLoading ? (
        <Skeleton className="h-40 w-full" />
      ) : (
        <div className="space-y-4">
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
        </div>
      )}
    </motion.section>
  )
}
