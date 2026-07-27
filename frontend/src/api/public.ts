// --------------------------------------------------------------------------------
// Copyright (c) 2016-NOW(至今) 筱锋
// Author: 筱锋「xiao_lfeng」(https://www.x-lf.com)
// --------------------------------------------------------------------------------
// 许可证声明：版权所有 (c) 2016-2026 筱锋。保留所有权利。
// 有关MIT许可证的更多信息，请查看项目根目录下的LICENSE文件或访问：
// https://opensource.org/licenses/MIT
// --------------------------------------------------------------------------------

import { request } from './client'
import type { HealthResponse } from './types'

/** 系统健康检查（公开接口，返回运行时长/内存/CPU/协程数等运行时指标） */
export function getHealth(): Promise<HealthResponse> {
  return request<HealthResponse>({
    method: 'GET',
    url: '/public/health',
  })
}
