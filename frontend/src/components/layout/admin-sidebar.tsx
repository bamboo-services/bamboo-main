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

import { Link, useLocation } from '@tanstack/react-router'
import {
  Heart,
  Home,
  LayoutDashboard,
  Link as LinkIcon,
  LogOut,
  MapPin,
  MoreHorizontal,
  Palette,
  Settings,
  Users,
} from 'lucide-react'
import {
  Sidebar,
  SidebarContent,
  SidebarFooter,
  SidebarGroup,
  SidebarGroupContent,
  SidebarGroupLabel,
  SidebarHeader,
  SidebarMenu,
  SidebarMenuButton,
  SidebarMenuItem,
} from '@/components/ui/sidebar'
import { BambooLogo } from '@/assets/svg/bamboo-logo'
import { siteConfig } from '@/lib/site'
import { useAuth } from '@/hooks/use-auth'
import { useSiteInfo } from '@/hooks/use-site-info'

const menuGroups = [
  {
    label: '首页',
    icon: Home,
    items: [{ title: '看板', url: '/admin/dashboard', icon: LayoutDashboard }],
  },
  {
    label: '友链',
    icon: Users,
    items: [
      { title: '友链管理', url: '/admin/link', icon: LinkIcon },
      { title: '位置管理', url: '/admin/location', icon: MapPin },
      { title: '颜色管理', url: '/admin/color', icon: Palette },
    ],
  },
  {
    label: '其他',
    icon: MoreHorizontal,
    items: [
      { title: '赞助', url: '/admin/sponsor', icon: Heart },
      { title: '设置', url: '/admin/setting', icon: Settings },
    ],
  },
]

export function AdminSidebar() {
  const { pathname } = useLocation()
  const { user, signOut } = useAuth()
  const { data: site } = useSiteInfo()
  const siteName = site?.site_name || siteConfig.defaultName

  const isActive = (url: string) => {
    if (url === '/admin/link') {
      return pathname === url || pathname.startsWith('/admin/link/')
    }
    return pathname === url
  }

  const displayName = user?.nickname || user?.username || '管理员'
  const userInitial = (user?.username ?? '?').charAt(0).toUpperCase()

  const handleLogout = async () => {
    await signOut()
    window.location.href = '/auth/login'
  }

  return (
    <Sidebar variant="floating" collapsible="icon">
      <SidebarHeader>
        <SidebarMenu>
          <SidebarMenuItem>
            <SidebarMenuButton size="lg" asChild>
              <Link to="/">
                <div className="flex aspect-square size-8 items-center justify-center rounded-lg">
                  <BambooLogo size={32} />
                </div>
                <div className="grid flex-1 text-left text-sm leading-tight">
                  <span className="truncate font-semibold">{siteName}</span>
                  <span className="truncate text-xs text-muted-foreground">
                    v{siteConfig.version}
                  </span>
                </div>
              </Link>
            </SidebarMenuButton>
          </SidebarMenuItem>
        </SidebarMenu>
      </SidebarHeader>

      <SidebarContent>
        {menuGroups.map((group) => (
          <SidebarGroup key={group.label}>
            <SidebarGroupLabel className="flex items-center gap-1">
              <group.icon className="h-4 w-4" />
              <span>{group.label}</span>
            </SidebarGroupLabel>
            <SidebarGroupContent>
              <SidebarMenu>
                {group.items.map((item) => (
                  <SidebarMenuItem key={item.title}>
                    <SidebarMenuButton asChild isActive={isActive(item.url)}>
                      <Link to={item.url}>
                        <item.icon className="h-4 w-4" />
                        <span>{item.title}</span>
                      </Link>
                    </SidebarMenuButton>
                  </SidebarMenuItem>
                ))}
              </SidebarMenu>
            </SidebarGroupContent>
          </SidebarGroup>
        ))}
      </SidebarContent>

      <SidebarFooter>
        <SidebarMenu>
          {/* 当前登录用户 */}
          <SidebarMenuItem>
            <SidebarMenuButton size="lg">
              <div className="flex aspect-square size-8 items-center justify-center rounded-lg bg-primary/10 text-sm font-semibold text-primary">
                {userInitial}
              </div>
              <div className="grid flex-1 text-left text-sm leading-tight">
                <span className="truncate font-semibold">{displayName}</span>
                <span className="truncate text-xs text-muted-foreground">
                  {user?.email ?? ''}
                </span>
              </div>
            </SidebarMenuButton>
          </SidebarMenuItem>
          {/* 退出登录 */}
          <SidebarMenuItem>
            <SidebarMenuButton onClick={handleLogout}>
              <LogOut className="h-4 w-4" />
              <span>退出登录</span>
            </SidebarMenuButton>
          </SidebarMenuItem>
        </SidebarMenu>
      </SidebarFooter>
    </Sidebar>
  )
}
