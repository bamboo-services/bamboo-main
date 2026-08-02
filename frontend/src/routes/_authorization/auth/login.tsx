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

import { Link, createFileRoute, redirect } from '@tanstack/react-router'
import { motion, useReducedMotion } from 'motion/react'
import { useState } from 'react'
import { Eye, EyeOff, KeyRound, LockKeyhole, User } from 'lucide-react'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { Checkbox } from '@/components/ui/checkbox'
import { BackLink, BambooArt, BrushUnderline } from '@/components/ink-wash'
import { BambooLogo } from '@/assets/svg/bamboo-logo'
import { siteConfig } from '@/lib/site'
import { SSO_OAUTH_LOGIN_URL, login } from '@/api/auth'
import { getStoredUser, getToken, setSession } from '@/lib/auth'
import { ROLE_ADMIN } from '@/lib/role'

/** 登录页 search 参数：redirect 为登录成功后的回跳路径 */
interface LoginSearch {
  redirect?: string
}

/** 仅信任同源内部路径（以单个 / 开头），避免开放重定向。
 *  无合法回跳目标时按角色分流：管理员去管理后台，其他去用户中心。 */
function resolveSafeRedirect(redirect?: string, role?: string): string {
  if (redirect && redirect.startsWith('/') && !redirect.startsWith('//')) {
    return redirect
  }
  return role === ROLE_ADMIN ? '/admin/dashboard' : '/user/dashboard'
}

export const Route = createFileRoute('/_authorization/auth/login')({
  validateSearch: (search: Record<string, unknown>): LoginSearch => ({
    redirect: typeof search.redirect === 'string' ? search.redirect : undefined,
  }),
  // 反向守卫：已登录用户按角色跳走，避免重复看到登录界面
  beforeLoad: ({ search }) => {
    if (getToken()) {
      throw redirect({
        to: resolveSafeRedirect(search.redirect, getStoredUser()?.role),
      })
    }
  },
  component: LoginPage,
})

function LoginPage() {
  const reduced = useReducedMotion() ?? false
  const { redirect: redirectTarget } = Route.useSearch()
  const [showPassword, setShowPassword] = useState(false)
  const [error, setError] = useState<string | null>(null)
  const [loading, setLoading] = useState(false)
  const [formData, setFormData] = useState({
    username: '',
    password: '',
    remember: false,
  })

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    setError(null)
    setLoading(true)
    try {
      const res = await login({
        username: formData.username,
        password: formData.password,
      })
      setSession(res.token, res.user, formData.remember)
      // 登录后整页跳转，确保应用以干净状态重新装载；按角色分流落地页
      window.location.href = resolveSafeRedirect(redirectTarget, res.user.role)
    } catch (err) {
      setError(err instanceof Error ? err.message : '登录失败，请稍后重试')
      setLoading(false)
    }
  }

  const handleSsoLogin = () => {
    window.location.href = SSO_OAUTH_LOGIN_URL
  }

  return (
    <div className="grid min-h-dvh bg-background lg:grid-cols-2">
      {/* ───────── 左栏：表单（视口左，装饰朝左外边缘） ───────── */}
      <div className="relative flex flex-col gap-6 overflow-hidden p-6 md:p-10">
        {/* 晨光墨晕：单色淡绿径向，重心偏左上 */}
        <div
          aria-hidden
          className="pointer-events-none absolute inset-0"
          style={{
            background:
              'radial-gradient(560px 220px at 18% 0%, oklch(0.88 0.1 105 / 0.20), transparent 72%), radial-gradient(420px 260px at 90% 100%, oklch(0.88 0.1 105 / 0.08), transparent 70%)',
          }}
        />
        {/* 墨韵竹叶水印：左栏朝左外，极淡 */}
        <BambooArt
          mirror
          className="pointer-events-none absolute -bottom-10 -left-[60px] h-full w-[420px] text-text-primary opacity-50"
        />

        {/* 顶部 Logo */}
        <div className="relative z-10 flex justify-center md:justify-start">
          <Link
            to="/"
            className="flex items-center gap-2.5 transition-opacity hover:opacity-80"
          >
            <BambooLogo size={30} />
            <span className="font-serif text-lg font-semibold text-text-primary">
              {siteConfig.defaultName}
            </span>
          </Link>
        </div>

        {/* 居中窄表单 */}
        <div className="relative z-10 flex flex-1 items-center justify-center">
          <motion.div
            initial={reduced ? false : { opacity: 0, y: 12 }}
            animate={{ opacity: 1, y: 0 }}
            transition={{ duration: 0.5, ease: 'easeOut' }}
            className="w-full max-w-[340px]"
          >
            <form className="flex flex-col" onSubmit={handleSubmit}>
              {/* 返回链接：标题上方的逃生口 */}
              <BackLink to="/" label="返回首页" className="mb-5" />

              {/* 标题组：衬线 + 笔刷下划线 */}
              <h1 className="font-serif text-[32px] font-bold leading-tight tracking-[0.01em] text-text-primary">
                欢迎回来
              </h1>
              <BrushUnderline className="mb-2.5 mt-2.5" />
              <p className="text-sm leading-relaxed text-text-secondary">
                登录以管理你的友情链接
              </p>

              <div className="mt-7 grid gap-[18px]">
                {/* 用户名 */}
                <div className="grid gap-2">
                  <Label htmlFor="username">用户名</Label>
                  <div className="relative">
                    <User className="absolute left-3 top-1/2 size-4 -translate-y-1/2 text-text-secondary/70" />
                    <Input
                      id="username"
                      type="text"
                      placeholder="请输入用户名"
                      autoComplete="username"
                      className="pl-9"
                      value={formData.username}
                      onChange={(e) =>
                        setFormData({ ...formData, username: e.target.value })
                      }
                      required
                    />
                  </div>
                </div>

                {/* 密码（可见性切换） */}
                <div className="grid gap-2">
                  <Label htmlFor="password">密码</Label>
                  <div className="relative">
                    <LockKeyhole className="absolute left-3 top-1/2 size-4 -translate-y-1/2 text-text-secondary/70" />
                    <Input
                      id="password"
                      type={showPassword ? 'text' : 'password'}
                      placeholder="请输入密码"
                      autoComplete="current-password"
                      className="pl-9 pr-9"
                      value={formData.password}
                      onChange={(e) =>
                        setFormData({ ...formData, password: e.target.value })
                      }
                      required
                    />
                    <button
                      type="button"
                      onClick={() => setShowPassword((v) => !v)}
                      className="absolute right-3 top-1/2 -translate-y-1/2 text-text-secondary transition-colors hover:text-text-primary"
                      aria-label={showPassword ? '隐藏密码' : '显示密码'}
                    >
                      {showPassword ? (
                        <EyeOff className="size-4" />
                      ) : (
                        <Eye className="size-4" />
                      )}
                    </button>
                  </div>
                </div>

                {/* 记住登录 */}
                <div className="flex items-center gap-2">
                  <Checkbox
                    id="remember"
                    checked={formData.remember}
                    onCheckedChange={(checked) =>
                      setFormData({ ...formData, remember: checked === true })
                    }
                  />
                  <Label
                    htmlFor="remember"
                    className="cursor-pointer text-sm font-normal text-text-secondary"
                  >
                    记住登录状态
                  </Label>
                </div>

                {/* 错误提示 */}
                {error && (
                  <p role="alert" className="text-sm text-destructive">
                    {error}
                  </p>
                )}

                {/* 登录按钮（全宽） */}
                <Button type="submit" className="w-full" disabled={loading}>
                  {loading ? '登录中…' : '登 录'}
                </Button>

                {/* 竹节分隔：或 */}
                <div className="flex items-center gap-3">
                  <span className="h-px flex-1 bg-border" />
                  <span className="size-1.5 shrink-0 rounded-[3px] bg-leaf-muted opacity-85" />
                  <span className="text-xs text-text-secondary">或</span>
                  <span className="size-1.5 shrink-0 rounded-[3px] bg-leaf-muted opacity-85" />
                  <span className="h-px flex-1 bg-border" />
                </div>

                {/* SSO 辅助登录 */}
                <Button
                  type="button"
                  variant="outline"
                  className="w-full"
                  onClick={handleSsoLogin}
                  disabled={loading}
                >
                  <KeyRound className="size-4" />
                  使用 SSO 账号登录
                </Button>
              </div>

              {/* 底部：仅认证流程内切换 */}
              <div className="mt-6 flex justify-center text-sm">
                <Link
                  to="/auth/register"
                  className="font-medium text-leaf-deep transition-opacity hover:opacity-75"
                >
                  注册账号
                </Link>
              </div>
            </form>
          </motion.div>
        </div>

        {/* 底部版权 */}
        <p className="relative z-10 text-center font-mono text-[11px] tracking-[0.1em] text-text-secondary/70">
          BAMBOO · FRIENDSHIP LINKS
        </p>
      </div>

      {/* ───────── 右栏：水墨叙事面板 ───────── */}
      <motion.div
        initial={reduced ? false : { opacity: 0 }}
        animate={{ opacity: 1 }}
        transition={{ duration: 0.8, ease: 'easeOut' }}
        className="relative hidden overflow-hidden lg:block"
        style={{
          background:
            'radial-gradient(700px 400px at 75% 12%, oklch(0.88 0.1 105 / 0.26), transparent 68%), radial-gradient(500px 380px at 20% 90%, oklch(0.88 0.1 105 / 0.12), transparent 70%), oklch(0.955 0.02 112)',
        }}
      >
        {/* 竹节竖线收边：贴左分界线 */}
        <span
          aria-hidden
          className="absolute bottom-[8%] left-0 top-[8%] w-px"
          style={{
            background:
              'linear-gradient(to bottom, transparent, var(--leaf-muted), transparent)',
            opacity: 0.5,
          }}
        />
        {/* 衬线水印大字：贴右外边缘 */}
        <span
          aria-hidden
          className="pointer-events-none absolute right-[100px] top-10 select-none font-serif text-[180px] font-black leading-none text-text-primary opacity-[0.04]"
        >
          竹
        </span>
        {/* 右上角 mono 标记 */}
        <span className="absolute right-12 top-10 font-mono text-[11px] uppercase tracking-[0.18em] text-text-secondary/70">
          Sign in · 竹林
        </span>
        {/* 墨韵竹叶：右栏朝右外 */}
        <BambooArt className="pointer-events-none absolute -right-10 top-0 h-full w-[560px] text-text-primary" />

        {/* 竖排题跋：贴左分界线内侧 */}
        <div className="absolute left-[72px] top-1/2 -translate-y-1/2 [writing-mode:vertical-rl] [text-orientation:upright]">
          <div className="flex">
            <p className="font-serif text-[52px] font-bold leading-[1.15] tracking-[0.28em] text-text-primary">
              竹林清晨
            </p>
            <p className="mr-[18px] font-serif text-[34px] font-semibold tracking-[0.3em] text-leaf-deep">
              归
            </p>
          </div>
          <p className="mr-4 font-serif text-[15px] leading-[2.4] tracking-[0.4em] text-text-secondary">
            万物生长，节节而成
          </p>
        </div>

        {/* 底部引言 */}
        <div className="absolute inset-x-0 bottom-14 pl-[72px] pr-[60px]">
          <p className="font-serif text-[17px] italic leading-relaxed tracking-[0.02em] text-text-primary/85">
            「{siteConfig.blogger.description}」
          </p>
          <p className="mt-2.5 font-mono text-[11px] uppercase tracking-[0.12em] text-text-secondary">
            {siteConfig.defaultName} · 友情链接管理
          </p>
        </div>
      </motion.div>
    </div>
  )
}
