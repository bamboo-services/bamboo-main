// --------------------------------------------------------------------------------
// Copyright (c) 2016-NOW(至今) 筱锋
// Author: 筱锋「xiao_lfeng」(https://www.x-lf.com)
// --------------------------------------------------------------------------------
// 许可证声明：版权所有 (c) 2016-2026 筱锋。保留所有权利。
// 有关MIT许可证的更多信息，请查看项目根目录下的LICENSE文件或访问：
// https://opensource.org/licenses/MIT
// --------------------------------------------------------------------------------

import { useEffect } from 'react'
import { useRouterState } from '@tanstack/react-router'
import { useArchive, useSiteInfo } from '@/hooks/use-site-info'
import { getPageTitle } from '@/lib/page-titles'

/** 站点信息加载完成前的默认标题（与 index.html 首屏 <title> 保持一致） */
const TITLE_LOADING = 'Loading......'

/**
 * 站点驱动的页面标题。
 * 非主页拼接为「站名 - 界面简介」；主页（无界面简介）拼接为「站名 - 站点描述」；
 * 依赖未就绪（含请求失败）期间保持默认 Loading 文案。
 */
export function useSitePageTitle() {
  const { data: site } = useSiteInfo()
  const { data: archive } = useArchive()
  const pageTitle = useRouterState({
    select: (state) => getPageTitle(state.location.pathname),
  })

  useEffect(() => {
    const siteName = site?.site_name
    if (!siteName) {
      document.title = TITLE_LOADING
      return
    }
    if (pageTitle) {
      document.title = `${siteName} - ${pageTitle}`
      return
    }
    // 主页：站名 - 站点描述
    document.title = archive?.site_description
      ? `${siteName} - ${archive.site_description}`
      : TITLE_LOADING
  }, [site, archive, pageTitle])
}
