// --------------------------------------------------------------------------------
// Copyright (c) 2016-NOW(至今) 筱锋
// Author: 筱锋「xiao_lfeng」(https://www.x-lf.com)
// --------------------------------------------------------------------------------
// 许可证声明：版权所有 (c) 2016-2026 筱锋。保留所有权利。
// 有关MIT许可证的更多信息，请查看项目根目录下的LICENSE文件或访问：
// https://opensource.org/licenses/MIT
// --------------------------------------------------------------------------------

import { request } from './client'
import type { FriendPublicResponse, LinkFriend } from './types'

/**
 * 获取公开友链列表。
 * - 后端只返回已通过审核且未失效的友链
 * - 可选 group_id 过滤
 */
export function getPublicLinks(groupId?: number): Promise<Array<LinkFriend>> {
  return request<FriendPublicResponse>({
    method: 'GET',
    url: '/links',
    params: groupId != null ? { group_id: groupId } : undefined,
  }).then((res) => res.links)
}
