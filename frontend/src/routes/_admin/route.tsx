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

import { Fragment, useEffect, useState } from 'react'
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

  // 内容实际由 window 滚动：sidebar 布局仅 min-h-svh 保底、不约束高度，
  // 内容会把布局撑高（内部 div 并不可滚），故滚动感知必须监听 window。
  useEffect(() => {
    const handleScroll = () => {
      setIsScrolled(window.scrollY > 0)
    }
    handleScroll()
    window.addEventListener('scroll', handleScroll, { passive: true })
    return () => window.removeEventListener('scroll', handleScroll)
  }, [])

  return (
    <SidebarProvider>
      <AdminSidebar />
      <SidebarInset className="flex flex-col">
        {/* 滚动感知 header：置顶时透明无缝、满高；下滑后收缩为脱离顶缘的圆角宣纸浮条。
            z-30 压住内容区（hero 等为 z-10）；bg-background 实底确保浮条悬浮间隙不透出滚动内容 */}
        <header className="sticky top-0 z-30 flex h-16 shrink-0 items-center bg-background px-3">
          <motion.div
            className="flex w-full items-center gap-2 px-2"
            initial={false}
            animate={{
              height: isScrolled ? 48 : 64,
              borderRadius: isScrolled ? 14 : 0,
              backgroundColor: isScrolled
                ? 'oklch(0.985 0.012 110)'
                : 'oklch(0.975 0.016 110 / 0)',
              boxShadow: isScrolled
                ? '0 16px 40px -16px oklch(0.32 0.06 155 / 0.28), 0 0 0 1px oklch(0.9 0.03 120)'
                : '0 0 0 0 oklch(0.32 0.06 155 / 0), 0 0 0 0 oklch(0.9 0.03 120 / 0)',
            }}
            transition={{ duration: 0.32, ease: [0.32, 0.72, 0, 1] as const }}
          >
            <SidebarTrigger className="-ml-1" />
            <AdminBreadcrumb pathname={pathname} />
          </motion.div>
        </header>

        {/* 内容区：随 window 滚动 */}
        <div className="flex-1">
          <main className="p-4 md:p-6">
            {/* 路由级过渡只做 opacity 淡入，不做位移：
                各页面（如 dashboard）自带错峰入场编排，若此处再叠加 y 位移，
                会与页面内部 section 的上滑嵌套叠加，产生「二次上滑」观感。 */}
            <motion.div
              key={pathname}
              initial={{ opacity: 0 }}
              animate={{ opacity: 1 }}
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
