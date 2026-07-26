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
import { siteInfo } from '@/data/mock/site-info'

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

  const isActive = (url: string) => {
    if (url === '/admin/link') {
      return pathname === url || pathname.startsWith('/admin/link/')
    }
    return pathname === url
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
                  <span className="truncate font-semibold">
                    {siteInfo.site.siteName}
                  </span>
                  <span className="truncate text-xs text-muted-foreground">
                    v{siteInfo.site.version}
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
          <SidebarMenuItem>
            <SidebarMenuButton asChild>
              <Link to="/">
                <LogOut className="h-4 w-4" />
                <span>退出登录</span>
              </Link>
            </SidebarMenuButton>
          </SidebarMenuItem>
        </SidebarMenu>
      </SidebarFooter>
    </Sidebar>
  )
}
