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
import { ArrowLeft, LoaderCircle } from 'lucide-react'
import { BambooLogo } from '@/assets/svg/bamboo-logo'
import { oauthCallback, oauthLogin } from '@/api/auth'
import { setSession } from '@/lib/auth'
import { siteConfig } from '@/lib/site'

export const Route = createFileRoute('/_authorization/auth/callback')({
  component: AuthCallbackPage,
})

/**
 * SSO 授权回调页（路径 /auth/callback，对应 SSO_REDIRECT_URI）。
 * SSO 提供商把浏览器重定向回此页并携带 code + state，这里换取访问令牌、
 * 再换取本地会话后进入管理后台；任一步失败则提示并提供返回登录页入口。
 */
function AuthCallbackPage() {
  const [error, setError] = useState<string | null>(null)
  const cancelledRef = useRef(false)

  useEffect(() => {
    const params = new URLSearchParams(window.location.search)
    const code = params.get('code')
    const state = params.get('state')

    if (!code || !state) {
      setError('缺少授权参数，无法完成登录')
      return
    }

    void (async () => {
      try {
        const token = await oauthCallback(code, state)
        const res = await oauthLogin(token.access_token)
        if (cancelledRef.current) return
        setSession(res.token, res.user, true)
        window.location.href = '/admin/dashboard'
      } catch (err) {
        if (cancelledRef.current) return
        setError(err instanceof Error ? err.message : 'SSO 登录失败，请重试')
      }
    })()

    return () => {
      cancelledRef.current = true
    }
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

      {error ? (
        <div className="flex flex-col items-center gap-4 text-center">
          <p className="text-sm text-destructive">{error}</p>
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
