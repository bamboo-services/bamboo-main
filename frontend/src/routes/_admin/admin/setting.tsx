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
import { Contact, Globe, Lock, Save, Send, UserRound } from 'lucide-react'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { Skeleton } from '@/components/ui/skeleton'
import { Textarea } from '@/components/ui/textarea'
import { InkNavRow, PageHead } from '@/components/ink-wash'
import { MarkdownEditor } from '@/components/markdown-editor'
import { enter } from '@/lib/motion'
import {
  useApplySiteInfo,
  useArchive,
  useBloggerInfo,
  useSiteInfo,
  useUpdateApplySiteInfo,
  useUpdateArchive,
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

type SettingSection = 'site' | 'archive' | 'applySite' | 'blogger' | 'invalid'

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
        {/* 左侧菜单：sticky 悬浮，滚动时吸附于 admin header 之下不滚走 */}
        <motion.nav
          {...enter(reduced, 0.12, {
            initial: { opacity: 0, y: 10 },
            animate: { opacity: 1, y: 0 },
            transition: { duration: 0.4, ease: 'easeOut' },
          })}
          className="lg:sticky lg:top-20"
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
              title="站点档案"
              desc="站点描述 · 自我介绍"
              active={section === 'archive'}
              onClick={() => setSection('archive')}
            />
            <InkNavRow
              icon={<Send className="size-4" />}
              title="申请展示"
              desc="交换友链 · 站点资料"
              active={section === 'applySite'}
              onClick={() => setSection('applySite')}
            />
            <InkNavRow
              icon={<Contact className="size-4" />}
              title="博主信息"
              desc="名士帖 · 个人展示"
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
        ) : section === 'archive' ? (
          <SiteArchivePanel reduced={reduced} />
        ) : section === 'applySite' ? (
          <ApplySitePanel reduced={reduced} />
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
  const [introduction, setIntroduction] = useState('')

  useEffect(() => {
    if (data) {
      setSiteName(data.site_name)
      setIntroduction(data.introduction)
    }
  }, [data])

  const handleSave = () => {
    updateSite.mutate({
      site_name: siteName,
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
              站名与主页介绍展示于公开页；站点描述已移入「站点档案」。
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

/** 站点档案面板：站点描述 + 自我介绍（均 Markdown，各自独立保存） */
function SiteArchivePanel({ reduced }: { reduced: boolean }) {
  const { data, isLoading } = useArchive()
  const updateArchive = useUpdateArchive()

  const [siteDesc, setSiteDesc] = useState('')
  const [aboutContent, setAboutContent] = useState('')

  useEffect(() => {
    if (data) {
      setSiteDesc(data.site_description)
      setAboutContent(data.about)
    }
  }, [data])

  const handleSave = () => {
    updateArchive.mutate({ site_description: siteDesc, about: aboutContent })
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
      <PanelHeader kicker="ARCHIVE · MARKDOWN" title="站点档案" />

      {isLoading ? (
        <div className="relative mt-8 space-y-8">
          <Skeleton className="h-64 w-full" />
          <Skeleton className="h-64 w-full" />
        </div>
      ) : (
        <>
          {/* 站点描述：Markdown 编辑器，展示于「贰 · 本站」 */}
          <div className="relative mt-8">
            <div className="mb-2 flex items-center justify-between">
              <Label htmlFor="archiveSiteDesc">站点描述</Label>
              <span className="font-mono text-[11px] tracking-wide text-text-secondary">
                · 展示于「贰 · 本站」
              </span>
            </div>
            <MarkdownEditor
              id="archiveSiteDesc"
              value={siteDesc}
              onChange={setSiteDesc}
              minHeight={200}
              placeholder="用 Markdown 书写本站介绍，支持加粗、横线等…"
              maxLength={5000}
            />
          </div>

          {/* 自我介绍：Markdown 编辑器，展示于「壹 · 自序」 */}
          <div className="relative mt-10 border-t border-border/60 pt-8">
            <div className="mb-2 flex items-center justify-between">
              <Label htmlFor="archiveAbout">自我介绍</Label>
              <span className="font-mono text-[11px] tracking-wide text-text-secondary">
                · 展示于「壹 · 自序」
              </span>
            </div>
            <MarkdownEditor
              id="archiveAbout"
              value={aboutContent}
              onChange={setAboutContent}
              minHeight={280}
              placeholder="用 Markdown 书写个人档案…"
              maxLength={10000}
            />
          </div>

          {/* 统一保存：站点描述与自我介绍一次提交 */}
          <div className="relative mt-8 flex flex-wrap items-center justify-between gap-3 border-t border-border/60 pt-5">
            <p className="font-serif text-sm italic text-text-secondary">
              站点描述与自我介绍统一保存，即时生效。
            </p>
            <Button onClick={handleSave} disabled={updateArchive.isPending}>
              <Save className="mr-2 h-4 w-4" />
              {updateArchive.isPending ? '保存中…' : '保存站点档案'}
            </Button>
          </div>
        </>
      )}
    </motion.section>
  )
}

/** 申请站点展示面板：operate/apply 交换友链场景的博主站点资料 */
function ApplySitePanel({ reduced }: { reduced: boolean }) {
  const { data, isLoading } = useApplySiteInfo()
  const updateApplySite = useUpdateApplySiteInfo()

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
    updateApplySite.mutate({
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
      <PanelHeader kicker="APPLY SITE · 交换友链" title="申请站点展示" />

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
              <Label htmlFor="applySiteName">站点名字</Label>
              <Input
                id="applySiteName"
                value={siteName}
                onChange={(e) => setSiteName(e.target.value)}
              />
            </div>
            <div className="space-y-2">
              <Label htmlFor="applySiteEmail">站长邮箱</Label>
              <Input
                id="applySiteEmail"
                value={email}
                onChange={(e) => setEmail(e.target.value)}
              />
            </div>
            <div className="space-y-2 md:col-span-2">
              <Label htmlFor="applySiteDescription">站点描述</Label>
              <Input
                id="applySiteDescription"
                value={siteDescription}
                onChange={(e) => setSiteDescription(e.target.value)}
              />
            </div>
            <div className="space-y-2">
              <Label htmlFor="applySiteUrl">站点地址</Label>
              <Input
                id="applySiteUrl"
                value={siteUrl}
                onChange={(e) => setSiteUrl(e.target.value)}
              />
            </div>
            <div className="space-y-2">
              <Label htmlFor="applySiteRss">站点订阅</Label>
              <Input
                id="applySiteRss"
                value={rss}
                onChange={(e) => setRss(e.target.value)}
              />
            </div>
            <div className="space-y-2 md:col-span-2">
              <Label htmlFor="applySiteImage">站点图片</Label>
              <Input
                id="applySiteImage"
                value={siteImage}
                onChange={(e) => setSiteImage(e.target.value)}
              />
            </div>
          </div>

          <div className="relative mt-8 flex flex-wrap items-center justify-between gap-3 border-t border-border/60 pt-5">
            <p className="font-serif text-sm italic text-text-secondary">
              访客申请友链前需先在自站添加博主友链，此处即供其复制的站点资料，展示于
              operate/apply 申请页。
            </p>
            <Button onClick={handleSave} disabled={updateApplySite.isPending}>
              <Save className="mr-2 h-4 w-4" />
              {updateApplySite.isPending ? '保存中…' : '保存申请展示'}
            </Button>
          </div>
        </>
      )}
    </motion.section>
  )
}

/** 博主信息面板：「关于我」名士帖个人展示 */
function BloggerInfoPanel({ reduced }: { reduced: boolean }) {
  const { data, isLoading } = useBloggerInfo()
  const updateBlogger = useUpdateBloggerInfo()

  const [nick, setNick] = useState('')
  const [description, setDescription] = useState('')
  const [blogUrl, setBlogUrl] = useState('')
  const [avatar, setAvatar] = useState('')

  useEffect(() => {
    if (data) {
      setNick(data.nick)
      setDescription(data.description)
      setBlogUrl(data.blog_url)
      setAvatar(data.avatar)
    }
  }, [data])

  const handleSave = () => {
    updateBlogger.mutate({
      nick,
      description,
      blog_url: blogUrl,
      avatar,
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
      <PanelHeader kicker="BLOGGER · 名士帖" title="博主信息" />

      {isLoading ? (
        <div className="relative mt-8 grid gap-x-6 gap-y-5 md:grid-cols-2">
          {Array.from({ length: 4 }).map((_, i) => (
            <Skeleton key={i} className="h-9 w-full" />
          ))}
        </div>
      ) : (
        <>
          <div className="relative mt-8 grid gap-x-6 gap-y-5 md:grid-cols-2">
            <div className="space-y-2">
              <Label htmlFor="bloggerNick">昵称</Label>
              <Input
                id="bloggerNick"
                value={nick}
                onChange={(e) => setNick(e.target.value)}
              />
            </div>
            <div className="space-y-2">
              <Label htmlFor="bloggerAvatar">头像地址</Label>
              <Input
                id="bloggerAvatar"
                value={avatar}
                onChange={(e) => setAvatar(e.target.value)}
              />
            </div>
            <div className="space-y-2 md:col-span-2">
              <Label htmlFor="bloggerDescription">个人简介</Label>
              <Input
                id="bloggerDescription"
                value={description}
                onChange={(e) => setDescription(e.target.value)}
              />
            </div>
            <div className="space-y-2 md:col-span-2">
              <Label htmlFor="bloggerBlogUrl">博客链接</Label>
              <Input
                id="bloggerBlogUrl"
                value={blogUrl}
                onChange={(e) => setBlogUrl(e.target.value)}
              />
            </div>
          </div>

          <div className="relative mt-8 flex flex-wrap items-center justify-between gap-3 border-t border-border/60 pt-5">
            <p className="font-serif text-sm italic text-text-secondary">
              昵称、个人简介、博客链接与头像展示于「关于我」名士帖。
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
              失效友链自动归入该分组，名称与描述经 bm_system
              热修改，公开「已失效」章节同步生效。
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
