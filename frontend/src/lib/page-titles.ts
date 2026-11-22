// --------------------------------------------------------------------------------
// Copyright (c) 2016-NOW(至今) 筱锋
// Author: 筱锋「xiao_lfeng」(https://www.x-lf.com)
// --------------------------------------------------------------------------------
// 许可证声明：版权所有 (c) 2016-2026 筱锋。保留所有权利。
// 有关MIT许可证的更多信息，请查看项目根目录下的LICENSE文件或访问：
// https://opensource.org/licenses/MIT
// --------------------------------------------------------------------------------

/**
 * 页面标题（界面简介）映射。
 * 非主页面的标题统一呈现为「站点名 - 界面简介」；
 * 主页（`/`）不在此表——其标题使用站点描述（见 useSitePageTitle）。
 */

/** 静态路径（无动态段）→ 界面简介 */
const STATIC_TITLES: Record<string, string> = {
  // 关于
  '/about/me': '关于我',
  '/about/friends': '好友',
  '/about/sponsor': '赞助',
  // 自助操作
  '/operate/apply': '申请友链',
  '/operate/sponsor': '赞助',
  // 认证
  '/auth/login': '登录',
  '/auth/register': '注册',
  '/auth/callback': '登录回调',
  // 用户中心
  '/user/dashboard': '我的面板',
  '/user/account': '账号设置',
  '/user/links': '我的友链',
  '/user/sponsors': '我的赞助',
  // 管理端
  '/admin/dashboard': '仪表盘',
  '/admin/link': '友链管理',
  '/admin/link/add': '添加友链',
  '/admin/link/anomaly': '异常管理',
  '/admin/link/verify': '友链审核',
  '/admin/sponsor': '赞助管理',
  '/admin/sponsor/verify': '赞助审核',
  '/admin/color': '颜色管理',
  '/admin/location': '位置管理',
  '/admin/ranking': '排位管理',
  '/admin/setting': '系统设置',
  '/admin/users': '用户管理',
}

/** 动态路径规则：path 模板段支持 `:param` 通配（数组序即优先级） */
const DYNAMIC_RULES: ReadonlyArray<{ title: string; path: string }> = [
  { title: '编辑友链', path: '/admin/link/:id/edit' },
  { title: '友链详情', path: '/admin/link/:id' },
  { title: '编辑友链', path: '/user/links/:id/edit' },
  { title: '调整位置', path: '/user/links/:id/edit-location' },
]

/** 路径模板匹配：按 `/` 分段比对，`:param` 段匹配任意一段 */
function matchPageTitle(pathname: string, template: string): boolean {
  const actual = pathname.split('/').filter(Boolean)
  const expected = template.split('/').filter(Boolean)
  if (actual.length !== expected.length) return false
  return expected.every((seg, i) => seg.startsWith(':') || seg === actual[i])
}

/** 取当前路径的界面简介；主页或无匹配返回 undefined */
export function getPageTitle(pathname: string): string | undefined {
  const staticTitle = STATIC_TITLES[pathname]
  if (staticTitle) return staticTitle
  return DYNAMIC_RULES.find((rule) => matchPageTitle(pathname, rule.path))?.title
}
