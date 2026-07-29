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
import { useEffect, useState } from 'react'
import {
  ArrowLeft,
  AtSign,
  Eye,
  EyeOff,
  LockKeyhole,
  Mail,
  ShieldCheck,
  User,
} from 'lucide-react'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { BambooLogo } from '@/assets/svg/bamboo-logo'
import { siteConfig } from '@/lib/site'
import defaultBackground from '@/assets/images/default-background.webp'
import { getStoredUser, getToken, setSession } from '@/lib/auth'
import { ROLE_ADMIN } from '@/lib/role'
import { useRegister, useSendRegisterCode } from '@/hooks/use-auth'

/** 注册页 search 参数：redirect 为注册成功后的回跳路径 */
interface RegisterSearch {
  redirect?: string
}

/** 仅信任同源内部路径（以单个 / 开头），避免开放重定向。
 *  无合法回跳目标时按角色分流：管理员去管理后台，其他去用户中心。 */
function resolveSafeRedirect(target?: string, role?: string): string {
  if (target && target.startsWith('/') && !target.startsWith('//')) {
    return target
  }
  return role === ROLE_ADMIN ? '/admin/dashboard' : '/user/dashboard'
}

/** 简单邮箱格式校验（发送验证码前置校验，最终以后端为准） */
function isValidEmail(email: string): boolean {
  return /^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(email)
}

export const Route = createFileRoute('/_authorization/auth/register')({
  validateSearch: (search: Record<string, unknown>): RegisterSearch => ({
    redirect: typeof search.redirect === 'string' ? search.redirect : undefined,
  }),
  // 反向守卫：已登录用户按角色跳走，避免重复看到注册界面
  beforeLoad: ({ search }) => {
    if (getToken()) {
      throw redirect({
        to: resolveSafeRedirect(search.redirect, getStoredUser()?.role),
      })
    }
  },
  component: RegisterPage,
})

function RegisterPage() {
  const reduced = useReducedMotion() ?? false
  const { redirect: redirectTarget } = Route.useSearch()
  const [showPassword, setShowPassword] = useState(false)
  const [error, setError] = useState<string | null>(null)
  const [countdown, setCountdown] = useState(0)
  const [formData, setFormData] = useState({
    email: '',
    code: '',
    username: '',
    nickname: '',
    password: '',
  })

  const sendCode = useSendRegisterCode()
  const registerMutation = useRegister()

  // 验证码发送倒计时（每秒递减，归零后允许重发）
  useEffect(() => {
    if (countdown <= 0) return
    const timer = setTimeout(() => setCountdown((c) => c - 1), 1000)
    return () => clearTimeout(timer)
  }, [countdown])

  const handleSendCode = () => {
    if (countdown > 0 || sendCode.isPending) return
    if (!isValidEmail(formData.email)) {
      setError('请输入正确的邮箱地址')
      return
    }
    setError(null)
    sendCode.mutate(formData.email, {
      onSuccess: () => setCountdown(60),
    })
  }

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    setError(null)
    try {
      const res = await registerMutation.mutateAsync({
        username: formData.username,
        email: formData.email,
        nickname: formData.nickname.trim() || undefined,
        password: formData.password,
        code: formData.code,
      })
      setSession(res.token, res.user, true)
      // 注册后整页跳转，确保应用以干净状态重新装载；按角色分流落地页
      window.location.href = resolveSafeRedirect(redirectTarget, res.user.role)
    } catch (err) {
      setError(err instanceof Error ? err.message : '注册失败，请稍后重试')
    }
  }

  return (
    <div className="grid min-h-dvh bg-background lg:grid-cols-2">
      {/* 左侧：表单栏 */}
      <div className="flex flex-col gap-6 p-6 md:p-10">
        {/* 顶部 Logo */}
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
                  注册账号
                </h1>
                <p className="text-sm text-muted-foreground">
                  注册后即可申请并管理你的友情链接
                </p>
              </div>

              <div className="grid gap-4">
                {/* 邮箱 + 发送验证码 */}
                <div className="grid gap-2">
                  <Label htmlFor="email">邮箱</Label>
                  <div className="flex gap-2">
                    <div className="relative flex-1">
                      <Mail className="absolute left-3 top-1/2 size-4 -translate-y-1/2 text-muted-foreground" />
                      <Input
                        id="email"
                        type="email"
                        placeholder="请输入邮箱"
                        autoComplete="email"
                        className="pl-9"
                        value={formData.email}
                        onChange={(e) =>
                          setFormData({ ...formData, email: e.target.value })
                        }
                        required
                      />
                    </div>
                    <Button
                      type="button"
                      variant="outline"
                      className="shrink-0 cursor-pointer"
                      disabled={countdown > 0 || sendCode.isPending}
                      onClick={handleSendCode}
                    >
                      {sendCode.isPending
                        ? '发送中…'
                        : countdown > 0
                          ? `${countdown}s`
                          : '发送验证码'}
                    </Button>
                  </div>
                </div>

                {/* 验证码 */}
                <div className="grid gap-2">
                  <Label htmlFor="code">验证码</Label>
                  <div className="relative">
                    <ShieldCheck className="absolute left-3 top-1/2 size-4 -translate-y-1/2 text-muted-foreground" />
                    <Input
                      id="code"
                      type="text"
                      inputMode="numeric"
                      maxLength={6}
                      placeholder="请输入 6 位验证码"
                      className="pl-9"
                      value={formData.code}
                      onChange={(e) =>
                        setFormData({ ...formData, code: e.target.value })
                      }
                      required
                    />
                  </div>
                </div>

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

                {/* 昵称（可选） */}
                <div className="grid gap-2">
                  <Label htmlFor="nickname">
                    昵称 <span className="text-muted-foreground">（可选）</span>
                  </Label>
                  <div className="relative">
                    <AtSign className="absolute left-3 top-1/2 size-4 -translate-y-1/2 text-muted-foreground" />
                    <Input
                      id="nickname"
                      type="text"
                      placeholder="给自己起个昵称吧"
                      className="pl-9"
                      value={formData.nickname}
                      onChange={(e) =>
                        setFormData({ ...formData, nickname: e.target.value })
                      }
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
                      placeholder="至少 6 位密码"
                      autoComplete="new-password"
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

                {/* 错误提示 */}
                {error && (
                  <p role="alert" className="text-sm text-destructive">
                    {error}
                  </p>
                )}

                {/* 注册按钮（全宽） */}
                <Button
                  type="submit"
                  className="w-full"
                  disabled={registerMutation.isPending}
                >
                  {registerMutation.isPending ? '注册中…' : '注册'}
                </Button>
              </div>

              {/* 已有账号 / 返回首页 */}
              <div className="flex items-center justify-center gap-4 text-sm">
                <Link
                  to="/auth/login"
                  className="text-primary transition-colors hover:opacity-80"
                >
                  已有账号？去登录
                </Link>
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
