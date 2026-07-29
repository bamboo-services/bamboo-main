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

import { Link, createFileRoute } from '@tanstack/react-router'
import { useEffect, useRef, useState } from 'react'
import { ArrowLeft, Ban, CircleAlert, LoaderCircle } from 'lucide-react'
import { BambooLogo } from '@/assets/svg/bamboo-logo'
import { oauthCallback, oauthLogin } from '@/api/auth'
import { setSession } from '@/lib/auth'
import { ROLE_ADMIN } from '@/lib/role'
import { siteConfig } from '@/lib/site'

export const Route = createFileRoute('/_authorization/auth/callback')({
  component: AuthCallbackPage,
})

/**
 * SSO 授权回调页（路径 /auth/callback，对应 SSO_REDIRECT_URI）。
 * SSO 提供商把浏览器重定向回此页并携带 code + state，这里换取访问令牌、
 * 再换取本地会话后进入管理后台；任一步失败则提示并提供返回登录页入口。
 */

/** 回调页提示：denied 为用户主动取消授权（中性提示），error 为授权或登录失败 */
interface CallbackNotice {
  kind: 'denied' | 'error'
  text: string
}

/** 将 SSO 回调的 error 参数转译为面向用户的友好提示 */
function resolveOAuthNotice(
  error: string,
  description: string | null,
): CallbackNotice {
  if (error === 'access_denied') {
    return { kind: 'denied', text: '已取消 SSO 授权，未登录账号' }
  }
  return {
    kind: 'error',
    text: description
      ? `SSO 授权失败：${description}`
      : `SSO 授权失败（${error}）`,
  }
}

function AuthCallbackPage() {
  const [notice, setNotice] = useState<CallbackNotice | null>(null)
  // React StrictMode 开发模式会「挂载→卸载→再挂载」双重调用 effect。
  // oauthCallback 消费一次性 code+state，重复执行必然因 state 已消费而失败，
  // 故用 ref 幂等守卫确保整个令牌换取流程只触发一次（ref 在双重挂载间保留）。
  const startedRef = useRef(false)

  useEffect(() => {
    if (startedRef.current) return
    startedRef.current = true

    const params = new URLSearchParams(window.location.search)
    const oauthError = params.get('error')
    const code = params.get('code')
    const state = params.get('state')

    // SSO 提供商在用户拒绝授权或授权异常时，会携带 error 参数回调本页，
    // 此时不会有 code，需优先识别并给出友好提示，而非继续走换取流程。
    if (oauthError) {
      setNotice(resolveOAuthNotice(oauthError, params.get('error_description')))
      return
    }

    if (!code || !state) {
      setNotice({ kind: 'error', text: '缺少授权参数，无法完成登录' })
      return
    }

    void (async () => {
      try {
        const token = await oauthCallback(code, state)
        const res = await oauthLogin(token.access_token)
        setSession(res.token, res.user, true)
        // 按角色分流落地页：管理员去管理后台，其他去用户中心
        window.location.href =
          res.user.role === ROLE_ADMIN ? '/admin/dashboard' : '/user/dashboard'
      } catch (err) {
        setNotice({
          kind: 'error',
          text: err instanceof Error ? err.message : 'SSO 登录失败，请重试',
        })
      }
    })()
  }, [])

  return (
    <div className="flex min-h-dvh flex-col items-center justify-center gap-6 bg-background p-6">
      <Link
        to="/"
        className="flex items-center gap-2.5 font-medium transition-opacity hover:opacity-80"
      >
        <BambooLogo size={30} />
        <span className="text-lg font-semibold text-text-primary">
          {siteConfig.defaultName}
        </span>
      </Link>

      {notice ? (
        <div className="flex flex-col items-center gap-4 text-center">
          <div
            className={
              notice.kind === 'denied'
                ? 'flex items-center gap-2 text-muted-foreground'
                : 'flex items-center gap-2 text-destructive'
            }
          >
            {notice.kind === 'denied' ? (
              <Ban className="size-4" />
            ) : (
              <CircleAlert className="size-4" />
            )}
            <p className="text-sm">{notice.text}</p>
          </div>
          <Link
            to="/auth/login"
            className="inline-flex items-center gap-1.5 text-sm text-muted-foreground transition-colors hover:text-primary"
          >
            <ArrowLeft className="size-3.5" />
            返回登录
          </Link>
        </div>
      ) : (
        <div className="flex flex-col items-center gap-3 text-muted-foreground">
          <LoaderCircle className="size-6 animate-spin" />
          <p className="text-sm">正在完成 SSO 登录…</p>
        </div>
      )}
    </div>
  )
}
