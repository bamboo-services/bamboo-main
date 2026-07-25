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
