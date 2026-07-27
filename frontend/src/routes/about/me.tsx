// --------------------------------------------------------------------------------
// Copyright (c) 2016-NOW(至今) 筱锋
// Author: 筱锋「xiao_lfeng」(https://www.x-lf.com)
// --------------------------------------------------------------------------------
// 许可证声明：版权所有 (c) 2016-2026 筱锋。保留所有权利。
// 有关MIT许可证的更多信息，请查看项目根目录下的LICENSE文件或访问：
// https://opensource.org/licenses/MIT
// --------------------------------------------------------------------------------

import { createFileRoute } from '@tanstack/react-router'
import { useQuery } from '@tanstack/react-query'
import { motion, useReducedMotion } from 'motion/react'
import { ExternalLink, Mail } from 'lucide-react'
import type { ComponentProps } from 'react'
import { getAbout, getSiteInfo } from '@/api/info'
import { MarkdownView } from '@/components/markdown'
import { Avatar, AvatarFallback, AvatarImage } from '@/components/ui/avatar'
import { Separator } from '@/components/ui/separator'
import { Skeleton } from '@/components/ui/skeleton'
import { siteConfig } from '@/lib/site'
import myAvatar from '@/assets/images/my_avatar.png'

export const Route = createFileRoute('/about/me')({
  component: AboutMePage,
})

type MotionDivProps = ComponentProps<typeof motion.div>

function enter(
  reduced: boolean,
  delay: number,
  full: MotionDivProps,
): MotionDivProps {
  if (!reduced) return { ...full, transition: { ...full.transition, delay } }
  return {
    initial: { opacity: 0 },
    animate: { opacity: 1 },
    transition: { duration: 0.3, delay: delay * 0.08 },
  }
}

function AboutMePage() {
  const reduced = useReducedMotion() ?? false
  const { data: site } = useQuery({
    queryKey: ['public', 'site'],
    queryFn: getSiteInfo,
  })
  const { data: about, isLoading } = useQuery({
    queryKey: ['public', 'about'],
    queryFn: getAbout,
  })

  const nick = siteConfig.blogger.nick
  const desc = siteConfig.blogger.description
  const email = siteConfig.blogger.email
  const displayName = site?.site_name ?? siteConfig.defaultName
  const introduction = site?.introduction ?? ''

  return (
    <motion.div
      {...enter(reduced, 0.2, {
        initial: { opacity: 0, y: 16 },
        animate: { opacity: 1, y: 0 },
        transition: { duration: 0.5, ease: 'easeOut' },
      })}
      className="mx-auto max-w-4xl"
    >
      <div className="grid gap-8 md:grid-cols-[auto_1fr]">
        {/* 左：博主名片 */}
        <aside className="flex flex-col items-center gap-4 md:items-start">
          <Avatar className="size-28 rounded-full ring-4 ring-ring-glow md:size-32">
            <AvatarImage src={myAvatar} alt={nick} />
            <AvatarFallback>{nick.slice(0, 1)}</AvatarFallback>
          </Avatar>
          <div className="text-center md:text-left">
            <p className="text-lg font-semibold text-text-primary">{nick}</p>
            <p className="text-sm text-text-secondary">{displayName}</p>
          </div>
          <div className="flex flex-col gap-2 text-sm">
            <a
              href="https://blog.x-lf.com"
              target="_blank"
              rel="noopener noreferrer"
              className="inline-flex items-center gap-1.5 text-text-secondary transition-colors hover:text-primary"
            >
              <ExternalLink className="size-4" />
              blog.x-lf.com
            </a>
            <a
              href={`mailto:${email}`}
              className="inline-flex items-center gap-1.5 text-text-secondary transition-colors hover:text-primary"
            >
              <Mail className="size-4" />
              {email}
            </a>
          </div>
        </aside>

        {/* 右：自我介绍 */}
        <div className="min-w-0">
          <p className="mb-4 text-base leading-relaxed text-text-secondary">
            {desc}
          </p>
          <Separator className="mb-6 bg-leaf-muted/40" />
          {isLoading ? (
            <div className="space-y-3">
              {Array.from({ length: 4 }).map((_, i) => (
                <Skeleton key={i} className="h-5 w-full" />
              ))}
              <Skeleton className="h-5 w-2/3" />
            </div>
          ) : about?.content ? (
            <MarkdownView content={about.content} />
          ) : (
            <p className="text-sm text-text-secondary">
              {introduction || '博主尚未填写自我介绍。'}
            </p>
          )}
        </div>
      </div>
    </motion.div>
  )
}
