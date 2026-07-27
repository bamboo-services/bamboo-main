// --------------------------------------------------------------------------------
// Copyright (c) 2016-NOW(至今) 筱锋
// Author: 筱锋「xiao_lfeng」(https://www.x-lf.com)
// --------------------------------------------------------------------------------
// 许可证声明：版权所有 (c) 2016-2026 筱锋。保留所有权利。
// 有关MIT许可证的更多信息，请查看项目根目录下的LICENSE文件或访问：
// https://opensource.org/licenses/MIT
// --------------------------------------------------------------------------------

import axios from 'axios'
import type { AxiosRequestConfig } from 'axios'
import type { BaseResponse } from './types'
import { clearSession, getToken } from '@/lib/auth'

/** 匹配 id 或 *_id 字段名（这些字段均为雪花 ID） */
const ID_KEY_PATTERN = /^(id|.*_id)$/
/** 匹配纯数字字符串 */
const NUMERIC_PATTERN = /^\d+$/

/**
 * JSON reviver：将 id / *_id 的纯数字字符串还原为 bigint。
 * 后端 SnowflakeID 以字符串传输，此处还原无损，规避 JS 安全整数精度问题。
 */
function reviveBigInt(key: string, value: unknown): unknown {
  if (
    typeof value === 'string' &&
    ID_KEY_PATTERN.test(key) &&
    NUMERIC_PATTERN.test(value)
  ) {
    return BigInt(value)
  }
  return value
}

/**
 * 统一 axios 实例。
 * - baseURL 优先取 VITE_API_BASE_URL，缺省走 Vite dev proxy（/api → 后端 5555）
 * - timeout 15s，匹配后端大部分公开接口
 * - transformResponse 将雪花 ID 字符串还原为 bigint
 * - transformRequest 将 bigint 序列化为字符串（后端反序列化兼容字符串）
 */
const client = axios.create({
  baseURL: import.meta.env.VITE_API_BASE_URL ?? '/api/v1',
  timeout: 15000,
  headers: { 'Content-Type': 'application/json' },
  transformResponse: [
    (data: unknown) => {
      if (typeof data !== 'string') return data
      try {
        return JSON.parse(data, reviveBigInt)
      } catch {
        return data
      }
    },
  ],
  transformRequest: [
    (data: unknown) => {
      if (data === undefined || data === null) return data
      return JSON.stringify(data, (_key, value) =>
        typeof value === 'bigint' ? value.toString() : value,
      )
    },
  ],
})

// 请求拦截器：存在登录令牌时自动注入 Authorization（调用方已显式指定则不覆盖）
client.interceptors.request.use((config) => {
  const token = getToken()
  if (token && !config.headers.Authorization) {
    config.headers.Authorization = `Bearer ${token}`
  }
  return config
})

// 响应拦截器：401 视为登录态失效，清空会话并回跳登录页（避免在登录页反复重定向）
client.interceptors.response.use(
  (res) => res,
  (error) => {
    if (axios.isAxiosError(error) && error.response?.status === 401) {
      clearSession()
      if (!window.location.pathname.startsWith('/auth/login')) {
        const redirect = encodeURIComponent(
          window.location.pathname + window.location.search,
        )
        window.location.href = `/auth/login?redirect=${redirect}`
      }
    }
    return Promise.reject(error)
  },
)

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
