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
import { motion, useReducedMotion } from 'motion/react'
import { useState } from 'react'
import {
  ArrowLeft,
  Eye,
  EyeOff,
  KeyRound,
  LockKeyhole,
  User,
} from 'lucide-react'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { Checkbox } from '@/components/ui/checkbox'
import { BambooLogo } from '@/assets/svg/bamboo-logo'
import { siteConfig } from '@/lib/site'
import defaultBackground from '@/assets/images/default-background.webp'
import { SSO_OAUTH_LOGIN_URL, login } from '@/api/auth'
import { setSession } from '@/lib/auth'

/** 登录页 search 参数：redirect 为登录成功后的回跳路径 */
interface LoginSearch {
  redirect?: string
}

export const Route = createFileRoute('/_authorization/auth/login')({
  validateSearch: (search: Record<string, unknown>): LoginSearch => ({
    redirect: typeof search.redirect === 'string' ? search.redirect : undefined,
  }),
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

  // 仅信任同源内部路径（以单个 / 开头），避免开放重定向
  const safeRedirect =
    redirectTarget &&
    redirectTarget.startsWith('/') &&
    !redirectTarget.startsWith('//')
      ? redirectTarget
      : '/admin/dashboard'

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
      // 登录后整页跳转，确保应用以干净状态重新装载
      window.location.href = safeRedirect
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
      {/* 左侧：表单栏 */}
      <div className="flex flex-col gap-6 p-6 md:p-10">
        {/* 顶部 Logo（shadcn login-04 风格） */}
        <div className="flex justify-center md:justify-start">
          <Link
            to="/"
            className="flex items-center gap-2.5 font-medium transition-opacity hover:opacity-80"
          >
            <BambooLogo size={30} />
            <span className="text-lg font-semibold text-text-primary">
              {siteConfig.defaultName}
            </span>
          </Link>
        </div>

        {/* 居中窄表单 */}
        <div className="flex flex-1 items-center justify-center">
          <motion.div
            initial={reduced ? false : { opacity: 0, y: 10 }}
            animate={{ opacity: 1, y: 0 }}
            transition={{ duration: 0.45, ease: 'easeOut' }}
            className="w-full max-w-xs"
          >
            <form className="flex flex-col gap-6" onSubmit={handleSubmit}>
              {/* 标题组 */}
              <div className="flex flex-col gap-2 text-center">
                <h1 className="text-2xl font-bold text-text-primary">
                  欢迎回来
                </h1>
                <p className="text-sm text-muted-foreground">
                  登录以管理你的友情链接
                </p>
              </div>

              <div className="grid gap-4">
                {/* 用户名 */}
                <div className="grid gap-2">
                  <Label htmlFor="username">用户名</Label>
                  <div className="relative">
                    <User className="absolute left-3 top-1/2 size-4 -translate-y-1/2 text-muted-foreground" />
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
                    <LockKeyhole className="absolute left-3 top-1/2 size-4 -translate-y-1/2 text-muted-foreground" />
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
                      className="absolute right-3 top-1/2 -translate-y-1/2 text-muted-foreground transition-colors hover:text-foreground"
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
                    className="cursor-pointer text-sm font-normal text-muted-foreground"
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
                  {loading ? '登录中…' : '登录'}
                </Button>

                {/* SSO 辅助登录 */}
                <div className="relative flex items-center justify-center">
                  <span className="absolute inset-x-0 top-1/2 h-px -translate-y-1/2 bg-border" />
                  <span className="relative bg-background px-2 text-xs text-muted-foreground">
                    或
                  </span>
                </div>
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

              {/* 返回首页 */}
              <div className="text-center text-sm">
                <Link
                  to="/"
                  className="inline-flex items-center gap-1.5 text-muted-foreground transition-colors hover:text-primary"
                >
                  <ArrowLeft className="size-3.5" />
                  返回首页
                </Link>
              </div>
            </form>
          </motion.div>
        </div>
      </div>

      {/* 右侧：图片栏（移动端隐藏，沿用首页竹林清晨图保持色调统一） */}
      <motion.div
        initial={reduced ? false : { opacity: 0 }}
        animate={{ opacity: 1 }}
        transition={{ duration: 0.8, ease: 'easeOut' }}
        className="relative hidden lg:block"
      >
        <img
          alt="Background"
          src={defaultBackground}
          className="absolute inset-0 h-full w-full object-cover"
        />
        {/* 主题绿调遮罩：底部压深以便承载引言文字 */}
        <div
          className="absolute inset-0"
          style={{
            background:
              'linear-gradient(to top, oklch(0.35 0.08 155 / 0.82) 0%, oklch(0.5 0.08 155 / 0.25) 45%, oklch(0.96 0.03 110 / 0.2) 100%)',
          }}
        />
        {/* 底部引言 */}
        <div className="absolute inset-x-0 bottom-0 p-10">
          <p className="text-lg font-medium leading-relaxed text-white">
            {siteConfig.blogger.description}
          </p>
          <p className="mt-2 text-sm text-white/70">
            {siteConfig.defaultName} · 友情链接管理
          </p>
        </div>
      </motion.div>
    </div>
  )
}
