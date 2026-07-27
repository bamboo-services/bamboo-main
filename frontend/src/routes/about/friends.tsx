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
import { ExternalLink, Sprout } from 'lucide-react'
import type { ComponentProps } from 'react'
import type { LinkFriend } from '@/api/types'
import { getPublicLinks } from '@/api/link'
import { Avatar, AvatarFallback, AvatarImage } from '@/components/ui/avatar'
import { Badge } from '@/components/ui/badge'
import { Skeleton } from '@/components/ui/skeleton'

export const Route = createFileRoute('/about/friends')({
  component: FriendsPage,
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

/** 友链按分组聚合（分组名取自后端嵌套的 group_f_key） */
function groupLinksByGroup(links: Array<LinkFriend>): Array<{
  groupId: string
  name: string
  links: Array<LinkFriend>
}> {
  const map = new Map<string, { name: string; links: Array<LinkFriend> }>()
  for (const link of links) {
    const key = link.group_id != null ? link.group_id.toString() : 'none'
    const entry = map.get(key) ?? {
      name: link.group_f_key?.name ?? '未分组',
      links: [],
    }
    entry.links.push(link)
    map.set(key, entry)
  }
  return Array.from(map.entries())
    .map(([groupId, entry]) => ({
      groupId,
      name: entry.name,
      links: entry.links.sort((a, b) => a.sort_order - b.sort_order),
    }))
    .sort((a, b) => {
      if (a.groupId === 'none') return 1
      if (b.groupId === 'none') return -1
      return 0
    })
}

/** 取友链主色（后端嵌套的 color_f_key.primary_color） */
function colorOf(link: LinkFriend): string | null {
  return link.color_f_key?.primary_color ?? null
}

function FriendsPage() {
  const reduced = useReducedMotion() ?? false
  const {
    data: links,
    isLoading,
    error,
  } = useQuery({
    queryKey: ['public', 'links'],
    queryFn: () => getPublicLinks(),
  })

  if (isLoading) {
    return (
      <div className="grid gap-4 sm:grid-cols-2 lg:grid-cols-3">
        {Array.from({ length: 6 }).map((_, i) => (
          <Skeleton key={i} className="h-32 w-full rounded-xl" />
        ))}
      </div>
    )
  }
  if (error || !links || links.length === 0) {
    return (
      <div className="py-16 text-center text-text-secondary">
        <Sprout className="mx-auto mb-3 size-10 text-leaf-muted" />
        <p>还没有友链，快来申请第一个吧～</p>
      </div>
    )
  }

  const groups = groupLinksByGroup(links)

  return (
    <motion.div
      {...enter(reduced, 0.2, {
        initial: { opacity: 0, y: 16 },
        animate: { opacity: 1, y: 0 },
        transition: { duration: 0.5, ease: 'easeOut' },
      })}
      className="space-y-8"
    >
      {groups.map((group) => (
        <section key={`g-${group.groupId}`}>
          <h3 className="mb-4 flex items-center gap-2 text-lg font-semibold text-text-primary">
            <span className="size-2 rounded-full bg-primary" />
            {group.name}
            <Badge variant="secondary" className="ml-1">
              {group.links.length}
            </Badge>
          </h3>
          <div className="grid gap-4 sm:grid-cols-2 lg:grid-cols-3">
            {group.links.map((link) => (
              <FriendCard key={link.id.toString()} link={link} />
            ))}
          </div>
        </section>
      ))}
    </motion.div>
  )
}

function FriendCard({ link }: { link: LinkFriend }) {
  const accent = colorOf(link)
  return (
    <a
      href={link.url}
      target="_blank"
      rel="noopener noreferrer"
      className="group relative flex gap-3 overflow-hidden rounded-xl border border-leaf-muted/40 bg-card/80 p-4 shadow-sm transition-all hover:-translate-y-0.5 hover:border-primary/40 hover:shadow-md"
    >
      {/* 颜色色条 */}
      {accent && (
        <span
          className="absolute inset-y-0 left-0 w-1"
          style={{ backgroundColor: accent }}
        />
      )}
      <Avatar className="size-12 shrink-0 rounded-full">
        <AvatarImage src={link.avatar ?? undefined} alt={link.name} />
        <AvatarFallback className="bg-leaf-light/40 text-text-primary">
          {link.name.slice(0, 1)}
        </AvatarFallback>
      </Avatar>
      <div className="min-w-0 flex-1">
        <div className="flex items-center justify-between gap-2">
          <h4 className="truncate font-medium text-text-primary group-hover:text-primary">
            {link.name}
          </h4>
          <ExternalLink className="size-4 shrink-0 text-text-secondary opacity-0 transition-opacity group-hover:opacity-100" />
        </div>
        <p className="mt-1 line-clamp-2 text-sm text-text-secondary">
          {link.description ?? '这个站点很神秘，没有留下描述。'}
        </p>
      </div>
    </a>
  )
}
