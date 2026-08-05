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
import { AnimatePresence, motion, useReducedMotion } from 'motion/react'
import { useEffect, useRef, useState } from 'react'
import { BambooLogo } from '@/assets/svg/bamboo-logo'
import { BackLink, BambooRule, EnsoIcon } from '@/components/ink-wash'
import { enter } from '@/lib/motion'
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
  const reduced = useReducedMotion() ?? false
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
    <div className="relative flex min-h-dvh flex-col items-center justify-center gap-6 overflow-hidden bg-background p-6">
      {/* 晨光墨晕：左上主光 + 右下辅光（与登录页左栏同源） */}
      <div
        aria-hidden
        className="pointer-events-none absolute inset-0"
        style={{
          background:
            'radial-gradient(560px 220px at 82% 0%, oklch(0.88 0.1 105 / 0.20), transparent 72%), radial-gradient(420px 260px at 10% 100%, oklch(0.88 0.1 105 / 0.08), transparent 70%)',
        }}
      />
      {/* 衬线水印大字：贴右下外缘，极淡 */}
      <span
        aria-hidden
        className="pointer-events-none absolute -bottom-8 right-[4%] select-none font-serif text-[200px] font-black leading-none text-text-primary opacity-[0.04]"
      >
        归
      </span>
      {/* 竹节竖线收边：左右分界线 */}
      <span
        aria-hidden
        className="absolute bottom-[12%] left-[5%] top-[12%] w-px"
        style={{
          background:
            'linear-gradient(to bottom, transparent, var(--leaf-muted), transparent)',
          opacity: 0.5,
        }}
      />
      <span
        aria-hidden
        className="absolute bottom-[12%] right-[5%] top-[12%] w-px"
        style={{
          background:
            'linear-gradient(to bottom, transparent, var(--leaf-muted), transparent)',
          opacity: 0.5,
        }}
      />
      {/* 顶部行：Logo（左）+ mono 角标（右上，仅桌面） */}
      <div className="relative z-10 flex w-full max-w-[420px] items-center justify-between">
        <Link
          to="/"
          className="flex items-center gap-2.5 transition-opacity hover:opacity-80"
        >
          <BambooLogo size={30} />
          <span className="font-serif text-lg font-semibold text-text-primary">
            {siteConfig.defaultName}
          </span>
        </Link>
        <span className="hidden font-mono text-[11px] uppercase tracking-[0.18em] text-text-secondary/70 lg:block">
          SSO · 竹林
        </span>
      </div>

      {/* 加载 / 错误两态：居中画心卡片，AnimatePresence 克制切换 */}
      <AnimatePresence mode="wait">
        {notice ? (
          <motion.div
            key={`notice-${notice.kind}`}
            {...enter(reduced, 0, {
              initial: { opacity: 0, y: 8 },
              animate: { opacity: 1, y: 0 },
              transition: { duration: 0.35, ease: 'easeOut' },
            })}
            className="relative z-10 w-full max-w-[420px]"
          >
            <div className="relative flex flex-col items-center gap-4 rounded-[2px] bg-[oklch(0.955_0.02_112)] px-9 py-10 text-center shadow-[0_0_0_1px_oklch(0.8_0.08_130/0.32),0_26px_54px_-30px_oklch(0.32_0.06_155/0.4)]">
              <EnsoIcon className="mt-1 size-12" />
              <h1 className="font-serif text-xl font-semibold tracking-wide text-text-primary">
                {notice.kind === 'denied' ? '授权未完成' : '登录失败'}
              </h1>
              <p
                className={`text-sm leading-relaxed ${
                  notice.kind === 'denied'
                    ? 'text-text-secondary'
                    : 'text-destructive'
                }`}
              >
                {notice.text}
              </p>
              <BambooRule reduced={reduced} delay={0} />
              <BackLink to="/auth/login" label="返回登录" className="mb-1" />
            </div>
          </motion.div>
        ) : (
          <motion.div
            key="loading"
            {...enter(reduced, 0, {
              initial: { opacity: 0, y: 8 },
              animate: { opacity: 1, y: 0 },
              transition: { duration: 0.35, ease: 'easeOut' },
            })}
            className="relative z-10 w-full max-w-[420px]"
          >
            <div className="relative flex flex-col items-center gap-4 rounded-[2px] bg-[oklch(0.955_0.02_112)] px-9 py-10 text-center shadow-[0_0_0_1px_oklch(0.8_0.08_130/0.32),0_26px_54px_-30px_oklch(0.32_0.06_155/0.4)]">
              {/* 一竿竹：完整复用 BambooArt 竹一数据，墨色淡雅、极慢摇曳 */}
              <div className="relative h-[150px] w-[150px]">
                <svg
                  className="text-text-primary absolute bottom-0 left-1/2 h-full -translate-x-1/2"
                  viewBox="258 0 260 432"
                  preserveAspectRatio="xMidYMax meet"
                  aria-hidden="true"
                >
                  <defs>
                    <path
                      id="bzleaf-cb"
                      d="M0 0 C12 -7 34 -11 58 -3 C36 4 12 5 0 0 Z"
                    />
                    <g id="spray-cb">
                      <path
                        d="M0 0 Q26 -9 56 -7"
                        fill="none"
                        stroke="currentColor"
                        strokeWidth="2"
                        strokeLinecap="round"
                      />
                      <use
                        href="#bzleaf-cb"
                        transform="translate(12 -3) rotate(-44) scale(1.08)"
                      />
                      <use
                        href="#bzleaf-cb"
                        transform="translate(28 -6) rotate(-14)"
                      />
                      <use
                        href="#bzleaf-cb"
                        transform="translate(44 -7) rotate(14) scale(.95)"
                      />
                      <use
                        href="#bzleaf-cb"
                        transform="translate(56 -7) rotate(42) scale(.72)"
                      />
                    </g>
                  </defs>
                  <g className="ink-sway" opacity="0.6">
                    <path
                      d="M392 432 C388 366 393 314 390 258 C387 202 392 142 389 86 C388 56 390 30 388 8"
                      fill="none"
                      stroke="currentColor"
                      strokeWidth="6.5"
                      strokeLinecap="round"
                    />
                    <ellipse
                      cx="390"
                      cy="350"
                      rx="7.5"
                      ry="2.8"
                      fill="none"
                      stroke="currentColor"
                      strokeWidth="2"
                    />
                    <ellipse
                      cx="389.5"
                      cy="258"
                      rx="7.5"
                      ry="2.8"
                      fill="none"
                      stroke="currentColor"
                      strokeWidth="2"
                    />
                    <ellipse
                      cx="390"
                      cy="166"
                      rx="7.5"
                      ry="2.8"
                      fill="none"
                      stroke="currentColor"
                      strokeWidth="2"
                    />
                    <ellipse
                      cx="390"
                      cy="82"
                      rx="6.5"
                      ry="2.4"
                      fill="none"
                      stroke="currentColor"
                      strokeWidth="2"
                    />
                    <use
                      href="#spray-cb"
                      transform="translate(391 318) rotate(-16)"
                    />
                    <use
                      href="#spray-cb"
                      transform="translate(389 226) scale(-1 1) rotate(-22)"
                    />
                    <use
                      href="#spray-cb"
                      transform="translate(390 134) rotate(-28) scale(1.12)"
                    />
                    <use
                      href="#spray-cb"
                      transform="translate(388 50) rotate(-54) scale(1.2)"
                    />
                    <use
                      href="#bzleaf-cb"
                      transform="translate(304 404) rotate(26) scale(.78)"
                    />
                  </g>
                </svg>
              </div>
              <h1 className="font-serif text-xl font-semibold tracking-wide text-text-primary">
                正在完成 SSO 登录
              </h1>
              <p className="font-serif text-[13px] tracking-[0.35em] text-text-secondary">
                以链相归 · 节节而生
              </p>
              <BambooRule reduced={reduced} delay={0} />
              <div className="flex items-center gap-2">
                <span
                  className="ink-pulse size-2 shrink-0 rounded-full bg-leaf-deep"
                  aria-hidden
                />
                <p className="font-mono text-xs tracking-[0.18em] text-text-secondary">
                  正在换取授权 · AUTHENTICATING
                </p>
              </div>
            </div>
          </motion.div>
        )}
      </AnimatePresence>
    </div>
  )
}
