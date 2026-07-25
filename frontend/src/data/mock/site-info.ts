/*
 * --------------------------------------------------------------------------------
 * Copyright (c) 2016-NOW 筱锋
 * Author: 筱锋(https://www.x-lf.com)
 * License: MIT
 * --------------------------------------------------------------------------------
 */

export interface SiteInfo {
  site: {
    siteName: string
    author: string
    version: string
    description: string
    keywords: string
  }
  blogger: {
    name: string
    nick: string
    email: string
    description: string
  }
}

export const siteInfo: SiteInfo = {
  site: {
    siteName: 'Bamboo',
    author: '筱锋',
    version: '2.0.0',
    description: '友链管理系统',
    keywords: '友链,博客,管理',
  },
  blogger: {
    name: 'xiao_lfeng',
    nick: '筱锋',
    email: 'gm@x-lf.cn',
    description: '一个热爱技术的开发者，喜欢折腾各种有趣的东西。',
  },
}
