// --------------------------------------------------------------------------------
// Copyright (c) 2016-NOW(至今) 筱锋
// Author: 筱锋「xiao_lfeng」(https://www.x-lf.com)
// --------------------------------------------------------------------------------
// 许可证声明：版权所有 (c) 2016-2026 筱锋。保留所有权利。
// 有关MIT许可证的更多信息，请查看项目根目录下的LICENSE文件或访问：
// https://opensource.org/licenses/MIT
// --------------------------------------------------------------------------------

import type { UserInfo } from '@/api/types'
import {
  Avatar,
  AvatarFallback,
  AvatarImage,
} from '@/components/ui/avatar'

/**
 * 当前用户头像：有 avatar URL 时渲染图片，缺失/加载失败回退首字色块。
 * 供管理侧栏 footer、用户中心顶导、公开页账户卡等身份展示位统一复用。
 *
 * - className：Avatar 容器样式（size / rounded / ring 等）
 * - fallbackClassName：首字回退块的样式（bg / text），默认沿用 ui/avatar 的 leaf-light 底
 */
export function UserAvatar({
  user,
  className,
  fallbackClassName,
}: {
  user?: UserInfo | null
  className?: string
  fallbackClassName?: string
}) {
  const displayName = user?.nickname || user?.username || '用户'
  const initial = (user?.username ?? '?').charAt(0).toUpperCase()

  return (
    <Avatar className={className}>
      <AvatarImage src={user?.avatar ?? undefined} alt={displayName} />
      <AvatarFallback className={fallbackClassName}>
        {initial}
      </AvatarFallback>
    </Avatar>
  )
}
