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

import { Fragment, useEffect, useRef, useState } from 'react'
import {
  Link,
  Outlet,
  createFileRoute,
  redirect,
  useRouterState,
} from '@tanstack/react-router'
import { motion } from 'motion/react'
import {
  SidebarInset,
  SidebarProvider,
  SidebarTrigger,
} from '@/components/ui/sidebar'
import { AdminSidebar } from '@/components/layout/admin-sidebar'
import { getToken } from '@/lib/auth'
import {
  Breadcrumb,
  BreadcrumbItem,
  BreadcrumbLink,
  BreadcrumbList,
  BreadcrumbPage,
  BreadcrumbSeparator,
} from '@/components/ui/breadcrumb'

export const Route = createFileRoute('/_admin')({
  // 路由守卫：无登录令牌时回跳登录页，并携带来源路径以便登录后返回
  beforeLoad: ({ location }) => {
    if (!getToken()) {
      throw redirect({
        to: '/auth/login',
        search: { redirect: location.pathname },
      })
    }
  },
  component: AdminLayout,
})

/** 面包屑片段：label 为展示文案，to 存在则可点击跳转 */
interface Crumb {
  label: string
  to?: string
}

/**
 * 根据当前路径解析面包屑层级。
 * 先匹配具名子路由（add/verify），再匹配动态 $id 路由，避免误命中。
 */
function getCrumbs(pathname: string): Array<Crumb> {
  const path = pathname.replace(/\/+$/, '') || '/'
  const crumbs: Array<Crumb> = [{ label: '管理后台', to: '/admin/dashboard' }]

  if (path === '/admin' || path === '/admin/dashboard') {
    crumbs.push({ label: '仪表盘' })
  } else if (path === '/admin/link') {
    crumbs.push({ label: '友链管理' })
  } else if (path === '/admin/link/add') {
    crumbs.push({ label: '友链管理', to: '/admin/link' }, { label: '添加友链' })
  } else if (path === '/admin/link/verify') {
    crumbs.push({ label: '友链管理', to: '/admin/link' }, { label: '友链审核' })
  } else if (/^\/admin\/link\/[^/]+\/edit$/.test(path)) {
    crumbs.push({ label: '友链管理', to: '/admin/link' }, { label: '编辑友链' })
  } else if (/^\/admin\/link\/[^/]+$/.test(path)) {
    crumbs.push({ label: '友链管理', to: '/admin/link' }, { label: '友链详情' })
  } else if (path === '/admin/sponsor') {
    crumbs.push({ label: '赞助管理' })
  } else if (path === '/admin/setting') {
    crumbs.push({ label: '系统设置' })
  } else {
    crumbs.push({ label: '未知页面' })
  }

  return crumbs
}

/** 动态面包屑：仅当前页节点随路由渐显，前缀保持稳定不重复动画 */
function AdminBreadcrumb({ pathname }: { pathname: string }) {
  const crumbs = getCrumbs(pathname)
  const prefix = crumbs.slice(0, -1)
  const current = crumbs[crumbs.length - 1]

  return (
    <Breadcrumb>
      <BreadcrumbList>
        {prefix.map((crumb, i) => (
          <Fragment key={`${crumb.label}-${i.toString()}`}>
            {i > 0 && <BreadcrumbSeparator />}
            <BreadcrumbItem>
              {crumb.to ? (
                <BreadcrumbLink asChild>
                  <Link to={crumb.to}>{crumb.label}</Link>
                </BreadcrumbLink>
              ) : (
                <BreadcrumbPage>{crumb.label}</BreadcrumbPage>
              )}
            </BreadcrumbItem>
          </Fragment>
        ))}
        {prefix.length > 0 && <BreadcrumbSeparator />}
        <BreadcrumbItem>
          {/* 只有当前页节点动画，key 用 label 保证仅文案变化时才重放 */}
          <motion.span
            key={current.label}
            initial={{ opacity: 0, y: -4 }}
            animate={{ opacity: 1, y: 0 }}
            transition={{ duration: 0.2, ease: 'easeOut' }}
            className="inline-flex items-center"
          >
            <BreadcrumbPage>{current.label}</BreadcrumbPage>
          </motion.span>
        </BreadcrumbItem>
      </BreadcrumbList>
    </Breadcrumb>
  )
}

function AdminLayout() {
  const pathname = useRouterState({ select: (s) => s.location.pathname })
  const [isScrolled, setIsScrolled] = useState(false)
  const scrollRef = useRef<HTMLDivElement>(null)

  useEffect(() => {
    const el = scrollRef.current
    if (!el) return
    const handleScroll = () => {
      setIsScrolled(el.scrollTop > 0)
    }
    el.addEventListener('scroll', handleScroll, { passive: true })
    return () => el.removeEventListener('scroll', handleScroll)
  }, [])

  return (
    <SidebarProvider>
      <AdminSidebar />
      <SidebarInset className="flex flex-col">
        {/* 滚动感知 header：置顶透明，滑动时浮起 */}
        <motion.header
          className="sticky top-0 z-10 flex h-16 shrink-0 items-center gap-2 px-4"
          animate={{
            backgroundColor: isScrolled
              ? 'oklch(0.96 0.03 110 / 0.75)'
              : 'oklch(0.96 0.03 110 / 0)',
            boxShadow: isScrolled
              ? '0 1px 3px 0 rgb(0 0 0 / 0.08)'
              : '0 0 0 0 rgb(0 0 0 / 0)',
            backdropFilter: isScrolled ? 'blur(12px)' : 'blur(0px)',
          }}
          transition={{ duration: 0.25, ease: 'easeOut' }}
        >
          <SidebarTrigger className="-ml-1" />
          <AdminBreadcrumb pathname={pathname} />
        </motion.header>

        {/* 可滚动内容区 */}
        <div ref={scrollRef} className="flex-1 overflow-y-auto">
          <main className="p-4 md:p-6">
            <motion.div
              key={pathname}
              initial={{ opacity: 0, y: 12 }}
              animate={{ opacity: 1, y: 0 }}
              transition={{ duration: 0.25, ease: 'easeOut' }}
              className="h-full"
            >
              <Outlet />
            </motion.div>
          </main>
        </div>
      </SidebarInset>
    </SidebarProvider>
  )
}
