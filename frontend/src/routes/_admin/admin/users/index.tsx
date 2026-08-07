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

import { useState } from 'react'
import { createFileRoute } from '@tanstack/react-router'
import { motion, useReducedMotion } from 'motion/react'
import { Search } from 'lucide-react'
import type { UserInfo } from '@/api/types'
import { Avatar, AvatarFallback, AvatarImage } from '@/components/ui/avatar'
import { Input } from '@/components/ui/input'
import { Pagination } from '@/components/ui/pagination'
import { Select } from '@/components/ui/select'
import { Skeleton } from '@/components/ui/skeleton'
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from '@/components/ui/table'
import {
  BambooRule,
  EnsoEmpty,
  InkBadge,
  InkSwitch,
  PageHead,
  inkTableHeadRow,
  inkTableRow,
  inkTableWrap,
  inkTd,
  inkTh,
} from '@/components/ink-wash'
import { useAdminUsers, useUpdateUserStatus } from '@/hooks/use-users'
import { formatDateTime } from '@/lib/datetime'
import { enter } from '@/lib/motion'
import { cn } from '@/lib/utils'

export const Route = createFileRoute('/_admin/admin/users/')({
  component: UsersPage,
})

const PAGE_SIZE = 10

/** 状态列文案：管理员恒显「管理员」；切换中显示「更新中」；否则按启用/禁用 */
function statusLabel(user: UserInfo, pending: boolean): string {
  if (user.is_admin) return '管理员'
  if (pending) return '更新中'
  return user.status === 1 ? '启用' : '禁用'
}

function UsersPage() {
  const reduced = useReducedMotion() ?? false
  const [page, setPage] = useState(1)
  const [keyword, setKeyword] = useState('')
  const [statusFilter, setStatusFilter] = useState<number | null>(null)

  const usersQuery = useAdminUsers({
    page,
    page_size: PAGE_SIZE,
    keyword: keyword.trim() || undefined,
    status: statusFilter ?? undefined,
    sort_by: 'created_at',
    sort_order: 'desc',
  })
  const updateStatus = useUpdateUserStatus()

  const users = usersQuery.data?.data ?? []
  const total = usersQuery.data?.pagination.total ?? 0
  const totalPages = usersQuery.data?.pagination.total_pages ?? 1

  return (
    <div className="mx-auto max-w-6xl space-y-6">
      <PageHead
        kicker="USERS · 用户"
        title="用户管理"
        sub="管理平台注册用户，启用或禁用其登录访问。系统管理员不受影响。"
        actions={
          <span className="font-mono text-xs text-text-secondary tabular-nums">
            共{' '}
            <span className="font-semibold text-leaf-deep">{total}</span> 位用户
          </span>
        }
      />

      <BambooRule reduced={reduced} delay={0.12} />

      {/* 工具栏：关键词搜索 + 状态筛选 */}
      <motion.div
        {...enter(reduced, 0.18, {
          initial: { opacity: 0, y: 10 },
          animate: { opacity: 1, y: 0 },
          transition: { duration: 0.4, ease: 'easeOut' },
        })}
        className="rounded-lg border border-border bg-card/60 p-5"
      >
        <div className="flex flex-wrap items-center gap-3">
          <div className="relative w-full max-w-xs">
            <Search className="absolute left-3 top-1/2 size-4 -translate-y-1/2 text-text-secondary" />
            <Input
              placeholder="搜索用户名 / 邮箱 / 昵称…"
              className="pl-9"
              value={keyword}
              onChange={(e) => {
                setKeyword(e.target.value)
                setPage(1)
              }}
            />
          </div>
          <Select
            aria-label="按状态筛选"
            className="w-36"
            value={statusFilter?.toString() ?? ''}
            onChange={(e) => {
              const value = e.target.value
              setStatusFilter(value ? Number(value) : null)
              setPage(1)
            }}
          >
            <option value="">全部状态</option>
            <option value="1">启用</option>
            <option value="0">禁用</option>
          </Select>
        </div>
      </motion.div>

      {/* 用户表格 */}
      <motion.div
        {...enter(reduced, 0.2, {
          initial: { opacity: 0 },
          animate: { opacity: 1 },
          transition: { duration: 0.4, ease: 'easeOut' },
        })}
        className={inkTableWrap}
      >
        <Table>
          <TableHeader>
            <TableRow className={cn(inkTableHeadRow, 'hover:bg-muted/30')}>
              <TableHead className={inkTh}>用户</TableHead>
              <TableHead className={inkTh}>昵称</TableHead>
              <TableHead className={inkTh}>邮箱认证</TableHead>
              <TableHead className={inkTh}>状态</TableHead>
              <TableHead className={cn(inkTh, 'hidden lg:table-cell')}>
                最后登录
              </TableHead>
              <TableHead className={cn(inkTh, 'hidden lg:table-cell')}>
                注册时间
              </TableHead>
            </TableRow>
          </TableHeader>
          <TableBody>
            {usersQuery.isLoading ? (
              Array.from({ length: 5 }).map((_, i) => (
                <TableRow key={i}>
                  <TableCell colSpan={6} className={inkTd}>
                    <Skeleton className="h-6 w-full" />
                  </TableCell>
                </TableRow>
              ))
            ) : users.length === 0 ? (
              <TableRow className="hover:bg-transparent">
                <TableCell colSpan={6} className="px-4 py-10">
                  <EnsoEmpty
                    title="没有找到用户"
                    hint="试试调整搜索条件或状态筛选"
                  />
                </TableCell>
              </TableRow>
            ) : (
              users.map((user) => {
                const isPending =
                  updateStatus.isPending &&
                  updateStatus.variables?.id.toString() === user.id.toString()
                return (
                  <TableRow
                    key={user.id.toString()}
                    className={cn('group', inkTableRow)}
                  >
                    {/* 用户：头像 + 用户名/管理员徽章 + 邮箱 */}
                    <TableCell className={inkTd}>
                      <div className="flex items-center gap-3">
                        <Avatar className="size-9 shrink-0 rounded-lg">
                          <AvatarImage
                            src={user.avatar ?? undefined}
                            alt={user.username}
                            loading="lazy"
                          />
                          <AvatarFallback className="bg-leaf-light/50 font-serif font-semibold text-leaf-deep">
                            {user.username.charAt(0)}
                          </AvatarFallback>
                        </Avatar>
                        <div className="min-w-0">
                          <div className="flex items-center gap-2">
                            <span className="truncate font-serif font-semibold text-text-primary">
                              {user.username}
                            </span>
                            {user.is_admin && (
                              <InkBadge tone="leaf">管理员</InkBadge>
                            )}
                          </div>
                          <p className="truncate font-mono text-xs text-text-secondary">
                            {user.email}
                          </p>
                        </div>
                      </div>
                    </TableCell>

                    {/* 昵称 */}
                    <TableCell className={cn(inkTd, 'text-text-secondary')}>
                      {user.nickname || '—'}
                    </TableCell>

                    {/* 邮箱认证 */}
                    <TableCell className={inkTd}>
                      <InkBadge
                        tone={user.email_verify ? 'leaf' : 'neutral'}
                      >
                        {user.email_verify ? '已认证' : '未认证'}
                      </InkBadge>
                    </TableCell>

                    {/* 状态：InkSwitch 直接切换；系统管理员行禁用开关 */}
                    <TableCell className={inkTd}>
                      <div className="flex items-center gap-2.5">
                        <InkSwitch
                          checked={user.status === 1}
                          disabled={user.is_admin || isPending}
                          onToggle={() =>
                            updateStatus.mutate({
                              id: user.id,
                              req: { status: user.status === 1 ? 0 : 1 },
                            })
                          }
                        />
                        <span className="font-mono text-xs text-text-secondary">
                          {statusLabel(user, isPending)}
                        </span>
                      </div>
                    </TableCell>

                    {/* 最后登录 / 注册时间 */}
                    <TableCell
                      className={cn(inkTd, 'hidden lg:table-cell')}
                    >
                      <span className="font-mono tabular-nums text-text-secondary">
                        {formatDateTime(user.last_login_at)}
                      </span>
                    </TableCell>
                    <TableCell
                      className={cn(inkTd, 'hidden lg:table-cell')}
                    >
                      <span className="font-mono tabular-nums text-text-secondary">
                        {formatDateTime(user.created_at)}
                      </span>
                    </TableCell>
                  </TableRow>
                )
              })
            )}
          </TableBody>
        </Table>

        {/* 分页 */}
        <div className="flex flex-wrap items-center justify-between gap-3 border-t border-border px-4 py-3">
          <span className="font-mono text-xs text-text-secondary">
            共 {total} 位 · 第 {page} / {Math.max(totalPages, 1)} 页
          </span>
          <Pagination
            pageIndex={page - 1}
            pageCount={Math.max(totalPages, 1)}
            onPageChange={(i) => setPage(i + 1)}
          />
        </div>
      </motion.div>
    </div>
  )
}
