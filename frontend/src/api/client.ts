// --------------------------------------------------------------------------------
// Copyright (c) 2016-NOW(至今) 筱锋
// Author: 筱锋「xiao_lfeng」(https://www.x-lf.com)
// --------------------------------------------------------------------------------
// 许可证声明：版权所有 (c) 2016-2026 筱锋。保留所有权利。
// 有关MIT许可证的更多信息，请查看项目根目录下的LICENSE文件或访问：
// https://opensource.org/licenses/MIT
// --------------------------------------------------------------------------------

import axios from 'axios'
import type {AxiosRequestConfig} from 'axios';
import type { BaseResponse } from './types'

/**
 * 统一 axios 实例。
 * - baseURL 优先取 VITE_API_BASE_URL，缺省走 Vite dev proxy（/api → 后端 5555）
 * - timeout 15s，匹配后端大部分公开接口
 */
const client = axios.create({
  baseURL: import.meta.env.VITE_API_BASE_URL ?? '/api/v1',
  timeout: 15000,
  headers: { 'Content-Type': 'application/json' },
})

/**
 * 请求辅助函数：拆 BaseResponse 包装。
 * - 成功（2xx）：返回 body.data，类型由调用方标注
 * - 业务错误（code 非 2xx）：抛 Error(message)
 * - 网络错误：抛 Error(后端 message || 浏览器 message)
 *
 * 写法上用泛型约束请求/响应，避免每个调用方重复拆包。
 */
export async function request<T>(config: AxiosRequestConfig): Promise<T> {
  try {
    const res = await client.request<BaseResponse<T>>(config)
    const body = res.data
    if (body.code < 200 || body.code >= 300) {
      throw new Error(body.message || body.error_message || '请求失败')
    }
    return body.data
  } catch (err) {
    if (axios.isAxiosError(err)) {
      const body = err.response?.data as BaseResponse<unknown> | undefined
      const msg = body?.message || err.message || '网络错误'
      throw new Error(msg)
    }
    throw err
  }
}

export { client }
