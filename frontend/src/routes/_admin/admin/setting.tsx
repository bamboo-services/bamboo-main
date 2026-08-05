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
import { Contact, Globe, Lock, Save, UserRound } from 'lucide-react'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { Skeleton } from '@/components/ui/skeleton'
import { Textarea } from '@/components/ui/textarea'
import { InkNavRow, PageHead } from '@/components/ink-wash'
import { enter } from '@/lib/motion'
import {
  useAbout,
  useBloggerInfo,
  useSiteInfo,
  useUpdateAbout,
  useUpdateBloggerInfo,
  useUpdateSiteInfo,
} from '@/hooks/use-site-info'
import {
  useBuiltinInvalidGroup,
  useUpdateBuiltinInvalidGroup,
} from '@/hooks/use-groups'

export const Route = createFileRoute('/_admin/admin/setting')({
  component: SettingPage,
})

type SettingSection = 'site' | 'about' | 'blogger' | 'invalid'

/**
 * 系统设置：主从面板 · 扁平化。
 * 不再套宣纸卡——左菜单与右编辑区直接落在宣纸上，
 * 靠晨光墨晕、斜墨条标题、留白与落款墨线立结构，轻、透气。
 */
function SettingPage() {
  const reduced = useReducedMotion() ?? false
  const [section, setSection] = useState<SettingSection>('site')

  return (
    <div className="mx-auto max-w-6xl space-y-6">
      <PageHead
        kicker="SETTINGS · 设置"
        title="系统设置"
        sub="管理站点信息与自我介绍。"
      />

      <div className="grid items-start gap-8 lg:grid-cols-[240px_1fr]">
        {/* 左侧菜单：扁平，无卡 */}
        <motion.nav
          {...enter(reduced, 0.12, {
            initial: { opacity: 0, y: 10 },
            animate: { opacity: 1, y: 0 },
            transition: { duration: 0.4, ease: 'easeOut' },
          })}
        >
          <p className="px-3.5 font-mono text-[10px] uppercase tracking-[0.18em] text-text-secondary">
            设置项 · PANELS
          </p>
          <div className="mt-2.5 flex flex-col gap-1">
            <InkNavRow
              icon={<Globe className="size-4" />}
              title="站点信息"
              desc="站名 · 描述 · 主页介绍"
              active={section === 'site'}
              onClick={() => setSection('site')}
            />
            <InkNavRow
              icon={<UserRound className="size-4" />}
              title="自我介绍"
              desc="Markdown 个人档案"
              active={section === 'about'}
              onClick={() => setSection('about')}
            />
            <InkNavRow
              icon={<Contact className="size-4" />}
              title="博主信息"
              desc="交换友链 · 站点资料"
              active={section === 'blogger'}
              onClick={() => setSection('blogger')}
            />
            <InkNavRow
              icon={<Lock className="size-4" />}
              title="已失效分组"
              desc="失效友链自动归集"
              active={section === 'invalid'}
              onClick={() => setSection('invalid')}
            />
          </div>
          <p className="mt-6 px-3.5 font-serif text-xs italic leading-relaxed text-text-secondary/80">
            「 一处落笔，全站生辉 」
          </p>
        </motion.nav>

        {/* 右侧编辑区：扁平，无卡 */}
        {section === 'site' ? (
          <SiteInfoPanel reduced={reduced} />
        ) : section === 'about' ? (
          <AboutPanel reduced={reduced} />
        ) : section === 'blogger' ? (
          <BloggerInfoPanel reduced={reduced} />
        ) : (
          <InvalidGroupPanel reduced={reduced} />
        )}
      </div>
    </div>
  )
}

/** 扁平面板头：晨光墨晕（氛围层）+ mono kicker + 斜墨条衬线大标题 */
function PanelHeader({ kicker, title }: { kicker: string; title: string }) {
  return (
    <>
      <div
        aria-hidden
        className="pointer-events-none absolute -top-8 inset-x-0 h-48"
        style={{
          background:
            'radial-gradient(ellipse 48% 58% at 38% 46%, oklch(0.88 0.1 105 / 0.15) 0%, oklch(0.88 0.1 105 / 0.08) 45%, transparent 75%)',
        }}
      />
      <header className="relative">
        <p className="font-mono text-[11px] uppercase tracking-[0.2em] text-leaf-deep">
          {kicker}
        </p>
        <h3 className="mt-2 flex items-center gap-3 font-serif text-2xl font-bold tracking-tight text-text-primary">
          <span className="h-1 w-5 -skew-x-12 rounded-sm bg-leaf-deep" />
          {title}
        </h3>
      </header>
    </>
  )
}

/** 站点信息面板 */
function SiteInfoPanel({ reduced }: { reduced: boolean }) {
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
      className="relative"
    >
      <PanelHeader kicker="SITE INFO" title="站点信息" />

      {isLoading ? (
        <div className="relative mt-8 grid gap-x-6 gap-y-5 md:grid-cols-2">
          <Skeleton className="h-9 w-full" />
          <Skeleton className="h-9 w-full" />
          <Skeleton className="h-20 w-full md:col-span-2" />
        </div>
      ) : (
        <>
          <div className="relative mt-8 grid gap-x-6 gap-y-5 md:grid-cols-2">
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
            <div className="space-y-2 md:col-span-2">
              <Label htmlFor="introduction">主页介绍</Label>
              <Textarea
                id="introduction"
                className="min-h-[96px]"
                value={introduction}
                onChange={(e) => setIntroduction(e.target.value)}
              />
            </div>
          </div>

          {/* 落款：墨线收束 + 衬线导语 + 保存 */}
          <div className="relative mt-8 flex flex-wrap items-center justify-between gap-3 border-t border-border/60 pt-5">
            <p className="font-serif text-sm italic text-text-secondary">
              站点信息将展示于公开页头部。
            </p>
            <Button onClick={handleSave} disabled={updateSite.isPending}>
              <Save className="mr-2 h-4 w-4" />
              {updateSite.isPending ? '保存中…' : '保存站点信息'}
            </Button>
          </div>
        </>
      )}
    </motion.section>
  )
}

/** 自我介绍面板 */
function AboutPanel({ reduced }: { reduced: boolean }) {
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
      {...enter(reduced, 0.18, {
        initial: { opacity: 0, y: 10 },
        animate: { opacity: 1, y: 0 },
        transition: { duration: 0.4, ease: 'easeOut' },
      })}
      className="relative"
    >
      <PanelHeader kicker="ABOUT · MARKDOWN" title="自我介绍" />

      {isLoading ? (
        <Skeleton className="relative mt-8 h-64 w-full" />
      ) : (
        <>
          <div className="relative mt-8 space-y-2">
            <Label htmlFor="aboutContent">介绍内容</Label>
            <Textarea
              id="aboutContent"
              className="min-h-[280px] font-mono text-sm"
              value={content}
              onChange={(e) => setContent(e.target.value)}
            />
          </div>

          <div className="relative mt-8 flex flex-wrap items-center justify-between gap-3 border-t border-border/60 pt-5">
            <p className="font-serif text-sm italic text-text-secondary">
              支持 Markdown，保存后即时生效。
            </p>
            <Button onClick={handleSave} disabled={updateAbout.isPending}>
              <Save className="mr-2 h-4 w-4" />
              {updateAbout.isPending ? '保存中…' : '保存自我介绍'}
            </Button>
          </div>
        </>
      )}
    </motion.section>
  )
}

/** 博主信息面板：交换友链场景的博主站点资料 */
function BloggerInfoPanel({ reduced }: { reduced: boolean }) {
  const { data, isLoading } = useBloggerInfo()
  const updateBlogger = useUpdateBloggerInfo()

  const [siteName, setSiteName] = useState('')
  const [siteDescription, setSiteDescription] = useState('')
  const [siteUrl, setSiteUrl] = useState('')
  const [siteImage, setSiteImage] = useState('')
  const [rss, setRss] = useState('')
  const [email, setEmail] = useState('')

  useEffect(() => {
    if (data) {
      setSiteName(data.site_name)
      setSiteDescription(data.site_description)
      setSiteUrl(data.site_url)
      setSiteImage(data.site_image)
      setRss(data.rss)
      setEmail(data.email)
    }
  }, [data])

  const handleSave = () => {
    updateBlogger.mutate({
      site_name: siteName,
      site_description: siteDescription,
      site_url: siteUrl,
      site_image: siteImage,
      rss,
      email,
    })
  }

  return (
    <motion.section
      {...enter(reduced, 0.18, {
        initial: { opacity: 0, y: 10 },
        animate: { opacity: 1, y: 0 },
        transition: { duration: 0.4, ease: 'easeOut' },
      })}
      className="relative"
    >
      <PanelHeader kicker="BLOGGER · EXCHANGE" title="博主信息" />

      {isLoading ? (
        <div className="relative mt-8 grid gap-x-6 gap-y-5 md:grid-cols-2">
          {Array.from({ length: 6 }).map((_, i) => (
            <Skeleton key={i} className="h-9 w-full" />
          ))}
        </div>
      ) : (
        <>
          <div className="relative mt-8 grid gap-x-6 gap-y-5 md:grid-cols-2">
            <div className="space-y-2">
              <Label htmlFor="bloggerSiteName">站点名字</Label>
              <Input
                id="bloggerSiteName"
                value={siteName}
                onChange={(e) => setSiteName(e.target.value)}
              />
            </div>
            <div className="space-y-2">
              <Label htmlFor="bloggerEmail">站长邮箱</Label>
              <Input
                id="bloggerEmail"
                value={email}
                onChange={(e) => setEmail(e.target.value)}
              />
            </div>
            <div className="space-y-2 md:col-span-2">
              <Label htmlFor="bloggerSiteDescription">站点描述</Label>
              <Input
                id="bloggerSiteDescription"
                value={siteDescription}
                onChange={(e) => setSiteDescription(e.target.value)}
              />
            </div>
            <div className="space-y-2">
              <Label htmlFor="bloggerSiteUrl">站点地址</Label>
              <Input
                id="bloggerSiteUrl"
                value={siteUrl}
                onChange={(e) => setSiteUrl(e.target.value)}
              />
            </div>
            <div className="space-y-2">
              <Label htmlFor="bloggerRss">站点订阅</Label>
              <Input
                id="bloggerRss"
                value={rss}
                onChange={(e) => setRss(e.target.value)}
              />
            </div>
            <div className="space-y-2 md:col-span-2">
              <Label htmlFor="bloggerSiteImage">站点图片</Label>
              <Input
                id="bloggerSiteImage"
                value={siteImage}
                onChange={(e) => setSiteImage(e.target.value)}
              />
            </div>
          </div>

          <div className="relative mt-8 flex flex-wrap items-center justify-between gap-3 border-t border-border/60 pt-5">
            <p className="font-serif text-sm italic text-text-secondary">
              访客申请友链前需先在自站添加博主友链，此处即供其复制的站点资料。
            </p>
            <Button onClick={handleSave} disabled={updateBlogger.isPending}>
              <Save className="mr-2 h-4 w-4" />
              {updateBlogger.isPending ? '保存中…' : '保存博主信息'}
            </Button>
          </div>
        </>
      )}
    </motion.section>
  )
}

/** 内置「已失效」分组面板：失效友链自动归集分组，名称/描述经 bm_system 热修改 */
function InvalidGroupPanel({ reduced }: { reduced: boolean }) {
  const { data, isLoading } = useBuiltinInvalidGroup()
  const updateBuiltinInvalid = useUpdateBuiltinInvalidGroup()

  const [invalidName, setInvalidName] = useState('')
  const [invalidDesc, setInvalidDesc] = useState('')

  useEffect(() => {
    if (data) {
      setInvalidName(data.name)
      setInvalidDesc(data.description ?? '')
    }
  }, [data])

  const handleSave = () => {
    updateBuiltinInvalid.mutate({
      name: invalidName.trim(),
      description: invalidDesc.trim(),
    })
  }

  return (
    <motion.section
      {...enter(reduced, 0.18, {
        initial: { opacity: 0, y: 10 },
        animate: { opacity: 1, y: 0 },
        transition: { duration: 0.4, ease: 'easeOut' },
      })}
      className="relative"
    >
      <PanelHeader kicker="INVALID GROUP · 已失效" title="已失效分组" />

      {isLoading ? (
        <div className="relative mt-8 grid gap-x-6 gap-y-5 md:grid-cols-2">
          <Skeleton className="h-9 w-full" />
          <Skeleton className="h-9 w-full" />
        </div>
      ) : (
        <>
          <div className="relative mt-8 grid gap-x-6 gap-y-5 md:grid-cols-2">
            <div className="space-y-2">
              <Label htmlFor="invalidGroupName">分组名称</Label>
              <Input
                id="invalidGroupName"
                value={invalidName}
                onChange={(e) => setInvalidName(e.target.value)}
              />
            </div>
            <div className="space-y-2">
              <Label htmlFor="invalidGroupDesc">分组描述</Label>
              <Input
                id="invalidGroupDesc"
                value={invalidDesc}
                onChange={(e) => setInvalidDesc(e.target.value)}
              />
            </div>
          </div>

          <div className="relative mt-8 flex flex-wrap items-center justify-between gap-3 border-t border-border/60 pt-5">
            <p className="font-serif text-sm italic text-text-secondary">
              失效友链自动归入该分组，名称与描述经 bm_system 热修改，公开「已失效」章节同步生效。
            </p>
            <Button
              onClick={handleSave}
              disabled={!invalidName.trim() || updateBuiltinInvalid.isPending}
            >
              <Save className="mr-2 h-4 w-4" />
              {updateBuiltinInvalid.isPending ? '保存中…' : '保存分组配置'}
            </Button>
          </div>
        </>
      )}
    </motion.section>
  )
}
