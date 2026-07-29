// --------------------------------------------------------------------------------
// Copyright (c) 2016-NOW(至今) 筱锋
// Author: 筱锋「xiao_lfeng」(https://www.x-lf.com)
// --------------------------------------------------------------------------------
// 许可证声明：版权所有 (c) 2016-2026 筱锋。保留所有权利。
// 有关MIT许可证的更多信息，请查看项目根目录下的LICENSE文件或访问：
// https://opensource.org/licenses/MIT
// --------------------------------------------------------------------------------

import { request } from './client'
import type {
  LoginRequest,
  LoginResponse,
  OAuthToken,
  PasswordChangeRequest,
  RegisterCodeRequest,
  RegisterRequest,
  RegisterResponse,
  UserInfoResponse,
} from './types'

/**
 * SSO 授权入口（完整路径，浏览器整页跳转用）。
 * 后端会 302 重定向到 SSO 提供商授权页，回调地址由 SSO_REDIRECT_URI 决定。
 */
export const SSO_OAUTH_LOGIN_URL = '/api/v1/sso/oauth/login'

/** 账号密码登录（POST /api/v1/auth/login） */
export function login(req: LoginRequest): Promise<LoginResponse> {
  return request<LoginResponse>({
    method: 'POST',
    url: '/auth/login',
    data: req,
  })
}

/** 用户注册（POST /api/v1/auth/register，需先经 sendRegisterCode 获取邮箱验证码） */
export function register(req: RegisterRequest): Promise<RegisterResponse> {
  return request<RegisterResponse>({
    method: 'POST',
    url: '/auth/register',
    data: req,
  })
}

/** 发送注册邮箱验证码（POST /api/v1/auth/register/code，60s 内不可重复发送） */
export function sendRegisterCode(email: string): Promise<void> {
  const payload: RegisterCodeRequest = { email }
  return request<void>({
    method: 'POST',
    url: '/auth/register/code',
    data: payload,
  })
}

/**
 * SSO OAuth 回调换取令牌（GET /api/v1/sso/oauth/callback）。
 * 浏览器从 SSO 回到 /auth/callback 后，用 query 中的 code + state 调用此接口。
 */
export function oauthCallback(
  code: string,
  state: string,
): Promise<OAuthToken> {
  return request<OAuthToken>({
    method: 'GET',
    url: '/sso/oauth/callback',
    params: { code, state },
  })
}

/**
 * 用 SSO 访问令牌换取本地会话（POST /api/v1/auth/oauth/login）。
 * 携带 Authorization: Bearer <accessToken>，由 client 拦截器注入。
 */
export function oauthLogin(accessToken: string): Promise<LoginResponse> {
  return request<LoginResponse>({
    method: 'POST',
    url: '/auth/oauth/login',
    headers: { Authorization: `Bearer ${accessToken}` },
  })
}

/** 获取当前登录用户信息（GET /api/v1/auth/user，需鉴权） */
export function getCurrentUser(): Promise<UserInfoResponse> {
  return request<UserInfoResponse>({ method: 'GET', url: '/auth/user' })
}

/** 修改密码（PUT /api/v1/auth/password/change，需鉴权） */
export function changePassword(req: PasswordChangeRequest): Promise<void> {
  return request<void>({
    method: 'PUT',
    url: '/auth/password/change',
    data: req,
  })
}

/**
 * 登出（PATCH /api/v1/auth/logout，需鉴权）。
 * 可显式传入令牌：调用方往往在登出时已清空本地会话，拦截器取不到令牌，
 * 故由调用方把捕获到的令牌经此显式携带，确保服务端登出成功。
 */
export function logout(token?: string): Promise<void> {
  return request<void>({
    method: 'PATCH',
    url: '/auth/logout',
    headers: token ? { Authorization: `Bearer ${token}` } : undefined,
  })
}
