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

export interface LinkItem {
  id: number
  siteName: string
  siteUrl: string
  siteLogo: string
  siteDescription: string
  webmasterEmail: string
  location: number
  locationName: string
  color: number
  colorName: string
  status: 'pending' | 'approved' | 'rejected'
  hasAdv: boolean
  ableConnect: boolean
  createdAt: string
  updatedAt: string
}

export interface LocationItem {
  id: number
  name: string
  description: string
}

export interface ColorItem {
  id: number
  name: string
  color: string
}

export const mockLocations: LocationItem[] = [
  { id: 1, name: '技术博客', description: '分享技术文章的博客' },
  { id: 2, name: '个人博客', description: '记录生活的个人博客' },
  { id: 3, name: '资源站点', description: '提供资源下载的站点' },
]

export const mockColors: ColorItem[] = [
  { id: 1, name: '默认', color: '#6366f1' },
  { id: 2, name: '绿色', color: '#22c55e' },
  { id: 3, name: '蓝色', color: '#3b82f6' },
  { id: 4, name: '紫色', color: '#a855f7' },
]

export const mockLinks: LinkItem[] = [
  {
    id: 1,
    siteName: '筱锋的博客',
    siteUrl: 'https://blog.x-lf.com',
    siteLogo: 'https://blog.x-lf.com/avatar.png',
    siteDescription: '记录技术与生活的个人博客',
    webmasterEmail: 'gm@x-lf.cn',
    location: 1,
    locationName: '技术博客',
    color: 2,
    colorName: '绿色',
    status: 'approved',
    hasAdv: false,
    ableConnect: true,
    createdAt: '2024-01-15',
    updatedAt: '2024-12-01',
  },
  {
    id: 2,
    siteName: '示例站点',
    siteUrl: 'https://example.com',
    siteLogo: 'https://example.com/logo.png',
    siteDescription: '这是一个示例站点的描述',
    webmasterEmail: 'admin@example.com',
    location: 2,
    locationName: '个人博客',
    color: 1,
    colorName: '默认',
    status: 'approved',
    hasAdv: false,
    ableConnect: true,
    createdAt: '2024-02-20',
    updatedAt: '2024-11-15',
  },
  {
    id: 3,
    siteName: '待审核站点',
    siteUrl: 'https://pending.example.com',
    siteLogo: 'https://pending.example.com/logo.png',
    siteDescription: '这是一个等待审核的站点',
    webmasterEmail: 'pending@example.com',
    location: 1,
    locationName: '技术博客',
    color: 3,
    colorName: '蓝色',
    status: 'pending',
    hasAdv: false,
    ableConnect: true,
    createdAt: '2024-12-10',
    updatedAt: '2024-12-10',
  },
  {
    id: 4,
    siteName: '资源分享站',
    siteUrl: 'https://resources.example.com',
    siteLogo: 'https://resources.example.com/logo.png',
    siteDescription: '分享各种有用资源的站点',
    webmasterEmail: 'resource@example.com',
    location: 3,
    locationName: '资源站点',
    color: 4,
    colorName: '紫色',
    status: 'approved',
    hasAdv: true,
    ableConnect: true,
    createdAt: '2024-03-05',
    updatedAt: '2024-10-20',
  },
]
