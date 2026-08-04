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

import {
  DndContext,
  DragOverlay,
  KeyboardSensor,
  PointerSensor,
  closestCorners,
  useDroppable,
  useSensor,
  useSensors,
} from '@dnd-kit/core'
import {
  SortableContext,
  arrayMove,
  rectSortingStrategy,
  sortableKeyboardCoordinates,
  useSortable,
  verticalListSortingStrategy,
} from '@dnd-kit/sortable'
import { CSS } from '@dnd-kit/utilities'
import { useBlocker } from '@tanstack/react-router'
import { GripVertical } from 'lucide-react'
import { useReducedMotion } from 'motion/react'
import { useCallback, useEffect, useMemo, useRef, useState } from 'react'
import { toast } from 'sonner'
import type { DragEndEvent, DragOverEvent, DragStartEvent } from '@dnd-kit/core'
import type { LinkFriend } from '@/api/types'
import type { FriendGroupSection } from '@/lib/friend-groups'
import { EnsoEmpty, InkBadge } from '@/components/ink-wash'
import { Skeleton } from '@/components/ui/skeleton'
import { domainOf } from '@/components/about/friend-card-shared'
import { CHAPTERS, LINK_LEVEL, groupLinksByGroup } from '@/lib/friend-groups'
import { accentOf, isFancyColor } from '@/lib/colors'
import { Button } from '@/components/ui/button'
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from '@/components/ui/dialog'
import { cn } from '@/lib/utils'
import { SiteAvatar } from '@/components/link/site-avatar'
import { usePublicLinks, useSortLinks } from '@/hooks/use-links'
import { useSortGroups } from '@/hooks/use-groups'

/** 章节 id（章节拖拽 sortable + 顶层 SortableContext）：'group:<id>' / 'group:none' */
const containerId = (section: FriendGroupSection) => `group:${section.groupId}`
/** 章节内栅格 droppable id（接收链接跨组拖入）：'container:<id>' */
const dropId = (section: FriendGroupSection) => `container:${section.groupId}`
/** 章节序序列化（groupIds × 链序），用于无意义拖拽抑制比对 */
const serializeSections = (sections: Array<FriendGroupSection>) =>
  JSON.stringify(
    sections.map((s) => [s.groupId, s.links.map((l) => l.id.toString())]),
  )

/** 级别徽章：好友/高级/广告（一般级别不显示），tone 走 ink-wash 既有 badgeTones */
function LevelBadge({ level }: { level: number }) {
  if (level === LINK_LEVEL.close) return <InkBadge tone="leaf">好友</InkBadge>
  if (level === LINK_LEVEL.premium) return <InkBadge tone="leaf">高级</InkBadge>
  if (level === LINK_LEVEL.ad) return <InkBadge tone="neutral">广告</InkBadge>
  return null
}

/** 左侧主色墨条：炫彩走 ink-fancy 流光，普通颜色取 accentOf 主色 */
function AccentBar({
  link,
  className,
}: {
  link: LinkFriend
  className?: string
}) {
  const fancy = isFancyColor(link.color_f_key)
  return (
    <span
      aria-hidden
      className={cn(
        'absolute left-0 rounded-r-full',
        fancy ? 'ink-fancy' : '',
        className,
      )}
      style={fancy ? undefined : { background: accentOf(link.color_f_key) }}
    />
  )
}

/**
 * 友链排位卡（纯视觉，h-full 填满栅格单元）：级别感知，与公开页逐格对应——
 * 高级 2×2 名帖式大卡；好友/广告 1×1 紧凑竖排；一般 1×1 横排富式（高度为好友一半）。
 * 复刻公开页 Bento 栅格层级，保证排位预览不失真。
 */
function FriendRankCard({ link }: { link: LinkFriend }) {
  if (link.level === LINK_LEVEL.premium) {
    return (
      <div className="group relative flex h-full flex-col items-center justify-center overflow-hidden rounded-lg border border-border bg-card p-4 text-center transition-[translate,box-shadow,border-color] duration-300 hover:-translate-y-0.5 hover:border-leaf-muted hover:shadow-[0_18px_40px_-28px_oklch(0.32_0.06_155/0.45)]">
        <AccentBar
          link={link}
          className="inset-y-4 w-[2px] transition-all duration-300 group-hover:inset-y-0 group-hover:w-[4px]"
        />
        <SiteAvatar
          name={link.name}
          url={link.avatar}
          className="size-14 rounded-full ring-2 ring-ring-glow"
        />
        <p className="mt-2.5 w-full truncate font-serif text-[15px] font-bold leading-tight text-text-primary">
          {link.name}
        </p>
        <p className="mt-1 w-full truncate font-mono text-[10px] uppercase tracking-[0.15em] text-text-secondary">
          {domainOf(link.url)}
        </p>
        {link.description && (
          <p className="mt-2 line-clamp-2 text-xs leading-relaxed text-text-secondary">
            {link.description}
          </p>
        )}
        <div className="mt-2">
          <LevelBadge level={link.level} />
        </div>
      </div>
    )
  }

  // 好友：1×2 紧凑竖排（对齐公开页好友竖排卡，高度为一般的两倍）
  if (link.level === LINK_LEVEL.close) {
    return (
      <div className="group relative flex h-full flex-col items-center justify-center gap-1 overflow-hidden rounded-lg border border-border bg-card px-3 text-center transition-[translate,box-shadow,border-color] duration-300 hover:-translate-y-0.5 hover:border-leaf-muted hover:shadow-[0_12px_26px_-22px_oklch(0.32_0.06_155/0.4)]">
        <AccentBar
          link={link}
          className="inset-y-2.5 w-[2px] transition-all duration-300 group-hover:inset-y-0 group-hover:w-[3px]"
        />
        <SiteAvatar
          name={link.name}
          url={link.avatar}
          className="size-9 rounded-full ring-1 ring-ring-glow"
        />
        <p className="w-full truncate font-serif text-sm font-bold leading-tight text-text-primary">
          {link.name}
        </p>
        <LevelBadge level={link.level} />
      </div>
    )
  }

  // 广告：1×1 极简竖排（与一般同高，56px；公开页广告同为竖排带标识）
  if (link.level === LINK_LEVEL.ad) {
    return (
      <div className="group relative flex h-full flex-col items-center justify-center gap-0.5 overflow-hidden rounded-lg border border-border bg-card px-2 text-center transition-[translate,box-shadow,border-color] duration-300 hover:-translate-y-0.5 hover:border-leaf-muted hover:shadow-[0_12px_26px_-22px_oklch(0.32_0.06_155/0.4)]">
        <AccentBar
          link={link}
          className="inset-y-2 w-[2px] transition-all duration-300 group-hover:inset-y-0 group-hover:w-[3px]"
        />
        <SiteAvatar
          name={link.name}
          url={link.avatar}
          className="size-6 rounded-full ring-1 ring-ring-glow"
        />
        <p className="w-full truncate font-serif text-[10px] font-bold leading-tight text-text-primary">
          {link.name}
        </p>
      </div>
    )
  }

  // 一般：1×1 横排富式（对齐公开页横排卡），row-span-1 高度为好友一半
  return (
    <div className="group relative flex h-full items-center gap-2 overflow-hidden rounded-lg border border-border bg-card px-2.5 transition-[translate,box-shadow,border-color] duration-300 hover:-translate-y-0.5 hover:border-leaf-muted hover:shadow-[0_12px_26px_-22px_oklch(0.32_0.06_155/0.4)]">
      <AccentBar
        link={link}
        className="inset-y-2 w-[2px] transition-all duration-300 group-hover:inset-y-0 group-hover:w-[3px]"
      />
      <SiteAvatar
        name={link.name}
        url={link.avatar}
        className="size-7 rounded-full ring-1 ring-ring-glow"
      />
      <div className="min-w-0 flex-1">
        <p className="truncate font-serif text-xs font-bold leading-tight text-text-primary">
          {link.name}
        </p>
        <p className="mt-0.5 truncate font-mono text-[9px] uppercase tracking-[0.12em] text-text-secondary">
          {domainOf(link.url)}
        </p>
      </div>
    </div>
  )
}

/** 可拖拽友链卡：栅格单元（高级占 2×2）+ useSortable 变换 */
function SortableRankCard({
  link,
  reduced,
}: {
  link: LinkFriend
  reduced: boolean
}) {
  const {
    attributes,
    listeners,
    setNodeRef,
    transform,
    transition,
    isDragging,
  } = useSortable({
    id: link.id.toString(),
    data: { type: 'link' },
  })
  const premium = link.level === LINK_LEVEL.premium
  const close = link.level === LINK_LEVEL.close

  return (
    <div
      ref={setNodeRef}
      style={{
        transform: CSS.Transform.toString(transform),
        transition: reduced ? undefined : transition,
      }}
      {...attributes}
      {...listeners}
      className={cn(
        'cursor-grab focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-leaf-deep/50 active:cursor-grabbing',
        premium && 'col-span-2 row-span-4',
        close && 'row-span-2',
        isDragging && 'opacity-40',
      )}
    >
      <FriendRankCard link={link} />
    </div>
  )
}

/**
 * 章节容器：整章可拖拽（useSortable，握把在章节头）+ 栅格 droppable 接收链接拖入。
 * 未分组（'none'）章节置底且不可拖。
 */
function SectionBlock({
  section,
  index,
  reduced,
  children,
}: {
  section: FriendGroupSection
  index: number
  reduced: boolean
  children: React.ReactNode
}) {
  const ungrouped = section.groupId === 'none'
  const {
    attributes,
    listeners,
    setNodeRef,
    transform,
    transition,
    isDragging,
  } = useSortable({
    id: containerId(section),
    data: { type: 'group' },
    disabled: ungrouped,
  })
  const { setNodeRef: setDropRef } = useDroppable({
    id: dropId(section),
    data: { type: 'container', groupId: section.groupId },
  })

  return (
    <section
      ref={setNodeRef}
      style={{
        transform: CSS.Transform.toString(transform),
        transition: reduced ? undefined : transition,
      }}
      {...attributes}
      className={cn(isDragging && 'opacity-40')}
    >
      <div
        {...listeners}
        className={cn(
          'group mb-3 flex items-center gap-2.5',
          !ungrouped && 'cursor-grab active:cursor-grabbing',
        )}
      >
        <GripVertical
          className={cn(
            'size-4 shrink-0 text-text-secondary transition-opacity',
            ungrouped ? 'opacity-0' : 'opacity-0 group-hover:opacity-100',
          )}
          aria-hidden
        />
        <span className="font-mono text-[11px] uppercase tracking-[0.3em] text-leaf-deep">
          {CHAPTERS[index] ?? '其'}
        </span>
        <span className="truncate font-serif text-xl font-bold text-text-primary">
          {section.name}
        </span>
        <span className="ml-auto shrink-0 font-mono text-[11px] uppercase tracking-[0.2em] text-text-secondary">
          {section.links.length} 位
        </span>
      </div>
      <SortableContext
        items={section.links.map((l) => l.id.toString())}
        strategy={rectSortingStrategy}
      >
        <div
          ref={setDropRef}
          className="grid grid-flow-dense grid-cols-2 gap-3 [grid-auto-rows:56px] sm:grid-cols-3 lg:grid-cols-4"
        >
          {children}
        </div>
      </SortableContext>
    </section>
  )
}

/**
 * 排位管理看板：数据源与访客页同源（['public','links'] + getPublicLinks），
 * 章节/链序由 groupLinksByGroup 派生（与公开页逐字同算法），预览即所见。
 *
 * 拖拽机制：组内重排 / 跨分组移动 / 章节拖拽统一拍平为全局有序载荷，600ms 防抖
 * 后经 PATCH /admin/links/sort 持久化（乐观更新 + 失败回滚至最近成功快照）。
 * 章节序经全局链序「首次出场」映射公开页，与分组排序值无关。
 */
export function RankingBoard() {
  const reduced = useReducedMotion() ?? false
  const { data: links, isLoading } = usePublicLinks()
  const sortLinks = useSortLinks()
  const sortGroups = useSortGroups()

  const [sections, setSections] = useState<Array<FriendGroupSection>>([])
  const [activeId, setActiveId] = useState<string | null>(null)
  const [saving, setSaving] = useState(false)
  const [dragActive, setDragActive] = useState(false)
  const [dirty, setDirty] = useState(false)

  const sectionsRef = useRef<Array<FriendGroupSection>>([])
  const snapshotRef = useRef<Array<FriendGroupSection>>([])
  const activeIdRef = useRef<string | null>(null)
  const savingRef = useRef(false)
  /** 拖拽开始时的序列化快照，用于无意义拖拽抑制 */
  const dragStartRef = useRef<string>('')
  /** 拖拽开始时的 sections 引用快照，用于 onDragCancel 回滚 onDragOver 的跨容器乐观移动 */
  const dragStartSectionsRef = useRef<Array<FriendGroupSection>>([])

  useEffect(() => {
    sectionsRef.current = sections
    activeIdRef.current = activeId
    savingRef.current = saving
  }, [sections, activeId, saving])

  // 仅在公开接口数据刷新时同步本地镜像（同源同缓存）；
  // 拖拽中/保存中跳过，避免用过期数据覆盖乐观状态（守卫走 ref，不进入依赖触发重同步）
  useEffect(() => {
    if (activeIdRef.current || savingRef.current) return
    const next = groupLinksByGroup(links ?? [])
    setSections(next)
    snapshotRef.current = next
    setDirty(false)
  }, [links])

  const sensors = useSensors(
    useSensor(PointerSensor, { activationConstraint: { distance: 6 } }),
    useSensor(KeyboardSensor, {
      coordinateGetter: sortableKeyboardCoordinates,
    }),
  )

  /** 由 id 定位所在章节（章节 id / droppable id → 自身；友链 id → 所属章节） */
  const findSection = useCallback(
    (id: string): FriendGroupSection | undefined => {
      if (id.startsWith('group:')) {
        return sectionsRef.current.find((s) => containerId(s) === id)
      }
      if (id.startsWith('container:')) {
        return sectionsRef.current.find((s) => dropId(s) === id)
      }
      return sectionsRef.current.find((s) =>
        s.links.some((l) => l.id.toString() === id),
      )
    },
    [],
  )

  /** 拍平 sections → 全局有序载荷（章节序 × 组内序，未分组置底） */
  const flush = useCallback(() => {
    const current = sectionsRef.current
    const items = current.flatMap((s) =>
      s.links.map((l) => ({
        id: l.id,
        group_id: s.groupId === 'none' ? null : BigInt(s.groupId),
      })),
    )
    if (items.length === 0) return
    // 章节序随分组排序值一并持久化（数字越小权重越大，与位置管理同源）
    const groupIds = current
      .filter((s) => s.groupId !== 'none')
      .map((s) => BigInt(s.groupId))
    setSaving(true)
    const done = () => {
      snapshotRef.current = sectionsRef.current
      setSaving(false)
      setDirty(false)
      toast.success('排位保存成功')
    }
    const rollback = (err?: Error) => {
      setSections(snapshotRef.current)
      setSaving(false)
      setDirty(false)
      toast.error(err?.message || '排位保存失败，已还原')
    }
    sortLinks.mutate(items, {
      onSuccess: () => {
        if (groupIds.length === 0) {
          done()
          return
        }
        sortGroups.mutate(groupIds, { onSuccess: done, onError: rollback })
      },
      onError: rollback,
    })
  }, [sortLinks, sortGroups])

  const onDragStart = useCallback(({ active }: DragStartEvent) => {
    setActiveId(String(active.id))
    setDragActive(true)
    dragStartRef.current = serializeSections(sectionsRef.current)
    dragStartSectionsRef.current = sectionsRef.current
  }, [])

  /** 跨容器拖动链接：实时把被拖卡片从源章节移到目标章节（乐观）；章节拖拽不在此处理 */
  const onDragOver = useCallback(
    ({ active, over }: DragOverEvent) => {
      if (!over || active.data.current?.type !== 'link') return
      const from = findSection(String(active.id))
      const to = findSection(String(over.id))
      if (!from || !to || from.groupId === to.groupId) return

      setSections((prev) => {
        const src = prev.find((s) => s.groupId === from.groupId)
        const dst = prev.find((s) => s.groupId === to.groupId)
        if (!src || !dst) return prev
        const srcLinks = [...src.links]
        const dstLinks = [...dst.links]
        const activeIndex = srcLinks.findIndex(
          (l) => l.id.toString() === String(active.id),
        )
        if (activeIndex < 0) return prev
        const [moved] = srcLinks.splice(activeIndex, 1)
        const overIndex = dstLinks.findIndex(
          (l) => l.id.toString() === String(over.id),
        )
        if (overIndex >= 0) dstLinks.splice(overIndex, 0, moved)
        else dstLinks.push(moved)
        return prev.map((s) => {
          if (s.groupId === src.groupId) return { ...s, links: srcLinks }
          if (s.groupId === dst.groupId) return { ...s, links: dstLinks }
          return s
        })
      })
    },
    [findSection],
  )

  /** 拖拽结束：链接定序 / 章节重排，标记未保存（无意义拖拽抑制） */
  const onDragEnd = useCallback(
    ({ active, over }: DragEndEvent) => {
      setActiveId(null)
      setDragActive(false)
      if (!over) return
      const type = active.data.current?.type
      // 先基于当前态计算拖拽后的 next（跨容器移动已在 onDragOver 落态），
      // 再与拖拽开始快照比较——避免 setSections 未提交导致 sectionsRef 仍为旧值的时序误判
      let next = sectionsRef.current

      if (type === 'group') {
        // 章节拖拽：整章移动（未分组置底，不作落点）
        const target = findSection(String(over.id))
        if (target && target.groupId !== 'none') {
          const fromIndex = next.findIndex(
            (s) => containerId(s) === String(active.id),
          )
          const toIndex = next.findIndex((s) => s.groupId === target.groupId)
          if (fromIndex >= 0 && toIndex >= 0 && fromIndex !== toIndex) {
            next = arrayMove(next, fromIndex, toIndex)
          }
        }
      } else {
        // 链接拖拽：同容器内定序（跨容器已在 onDragOver 完成移动）
        const section = findSection(String(active.id))
        const overSection = findSection(String(over.id))
        if (section && overSection && section.groupId === overSection.groupId) {
          const oldIndex = section.links.findIndex(
            (l) => l.id.toString() === String(active.id),
          )
          const newIndex = section.links.findIndex(
            (l) => l.id.toString() === String(over.id),
          )
          if (oldIndex >= 0 && newIndex >= 0 && oldIndex !== newIndex) {
            next = next.map((s) =>
              s.groupId === section.groupId
                ? { ...s, links: arrayMove(s.links, oldIndex, newIndex) }
                : s,
            )
          }
        }
      }

      // 无意义拖拽抑制：序列化未变则零请求
      const after = serializeSections(next)
      if (after !== dragStartRef.current) {
        setSections(next)
        setDirty(true)
      }
    },
    [findSection],
  )

  /** 拖拽取消（Esc / pointercancel / touchcancel）：复位状态并回滚跨容器乐观移动 */
  const onDragCancel = useCallback(() => {
    setActiveId(null)
    setDragActive(false)
    setSections(dragStartSectionsRef.current)
  }, [])

  /** 未保存更改时拦截 SPA 内路由导航与页面关闭（useBlocker 含 beforeunload 兜底） */
  const blocker = useBlocker({
    shouldBlockFn: () => dirty,
    enableBeforeUnload: dirty,
    withResolver: true,
  })

  const activeLink = useMemo(() => {
    if (!activeId) return null
    for (const s of sections) {
      const found = s.links.find((l) => l.id.toString() === activeId)
      if (found) return found
    }
    return null
  }, [activeId, sections])

  const total = useMemo(
    () => sections.reduce((sum, s) => sum + s.links.length, 0),
    [sections],
  )

  return (
    <div className="space-y-8">
      {/* 状态条：操作提示 + 总数 + 未保存指示 + 保存按钮 */}
      <div className="flex flex-wrap items-center justify-between gap-3">
        <p className="font-mono text-[11px] uppercase tracking-[0.22em] text-text-secondary">
          拖拽卡片调整排位 · 跨章拖动即改分组 · 点击保存持久化
        </p>
        <div className="flex items-center gap-2.5">
          {dirty && <InkBadge tone="leaf">有未保存的更改</InkBadge>}
          <span className="font-mono text-[11px] uppercase tracking-[0.2em] text-text-secondary">
            共 {total} 位
          </span>
          <Button size="sm" onClick={flush} disabled={!dirty || saving}>
            {saving ? '保存中…' : '保存'}
          </Button>
        </div>
      </div>

      {/* 看板主体 */}
      {isLoading ? (
        <div className="grid grid-flow-dense grid-cols-2 gap-3 [grid-auto-rows:56px] sm:grid-cols-3 lg:grid-cols-4">
          <Skeleton className="col-span-2 row-span-4 rounded-lg" />
          {Array.from({ length: 5 }).map((_, i) => (
            <Skeleton
              key={i}
              className={cn('rounded-lg', i % 2 === 0 && 'row-span-2')}
            />
          ))}
        </div>
      ) : sections.length === 0 ? (
        <EnsoEmpty title="暂无已通过的友链" hint="审核通过后在此管理排位" />
      ) : (
        <DndContext
          sensors={sensors}
          collisionDetection={closestCorners}
          onDragStart={onDragStart}
          onDragOver={onDragOver}
          onDragEnd={onDragEnd}
          onDragCancel={onDragCancel}
        >
          <SortableContext
            items={sections.map(containerId)}
            strategy={verticalListSortingStrategy}
          >
            <div
              className={cn('space-y-10', dragActive && 'pointer-events-none')}
            >
              {sections.map((section, gi) => (
                <SectionBlock
                  key={section.groupId}
                  section={section}
                  index={gi}
                  reduced={reduced}
                >
                  {section.links.map((link) => (
                    <SortableRankCard
                      key={link.id.toString()}
                      link={link}
                      reduced={reduced}
                    />
                  ))}
                </SectionBlock>
              ))}
            </div>
          </SortableContext>
          <DragOverlay
            dropAnimation={
              reduced ? null : { duration: 200, easing: 'ease-out' }
            }
            className="pointer-events-none"
          >
            {activeLink ? <FriendRankCard link={activeLink} /> : null}
          </DragOverlay>
        </DndContext>
      )}

      {/* 未保存更改的离开确认：放弃则 proceed 放行导航，继续编辑则 reset 留在当前页 */}
      {blocker.status === 'blocked' && (
        <Dialog
          open
          onOpenChange={(open) => {
            if (!open) blocker.reset()
          }}
        >
          <DialogContent>
            <DialogHeader>
              <DialogTitle>放弃未保存的更改？</DialogTitle>
              <DialogDescription>
                当前排位有未保存的调整，离开此页将丢失这些更改。
              </DialogDescription>
            </DialogHeader>
            <DialogFooter>
              <Button variant="outline" onClick={() => blocker.reset()}>
                继续编辑
              </Button>
              <Button onClick={() => blocker.proceed()}>放弃更改并离开</Button>
            </DialogFooter>
          </DialogContent>
        </Dialog>
      )}
    </div>
  )
}
