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
import { BackLink, BambooArt, BrushUnderline } from '@/components/ink-wash'
import { BambooLogo } from '@/assets/svg/bamboo-logo'
import { siteConfig } from '@/lib/site'
import { getStoredUser, getToken, setSession } from '@/lib/auth'
import { useRegister, useSendRegisterCode } from '@/hooks/use-auth'

/** 注册页 search 参数：redirect 为注册成功后的回跳路径 */
interface RegisterSearch {
  redirect?: string
}

/** 仅信任同源内部路径（以单个 / 开头），避免开放重定向。
 *  无合法回跳目标时按管理员身份分流：管理员去管理后台，其他去用户中心。 */
function resolveSafeRedirect(target?: string, isAdminUser?: boolean): string {
  if (target && target.startsWith('/') && !target.startsWith('//')) {
    return target
  }
  return isAdminUser ? '/admin/dashboard' : '/user/dashboard'
}

/** 简单邮箱格式校验（发送验证码前置校验，最终以后端为准） */
function isValidEmail(email: string): boolean {
  return /^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(email)
}

export const Route = createFileRoute('/_authorization/auth/register')({
  validateSearch: (search: Record<string, unknown>): RegisterSearch => ({
    redirect: typeof search.redirect === 'string' ? search.redirect : undefined,
  }),
  // 反向守卫：已登录用户按管理员身份跳走，避免重复看到注册界面
  beforeLoad: ({ search }) => {
    if (getToken()) {
      throw redirect({
        to: resolveSafeRedirect(search.redirect, getStoredUser()?.is_admin),
      })
    }
  },
  component: RegisterPage,
})

/** 站点 Logo：认证页两栏镜像时，桌面/移动的视口左上承载者不同，故抽成片段复用 */
function SiteLogo() {
  return (
    <Link
      to="/"
      className="flex items-center gap-2.5 transition-opacity hover:opacity-80"
    >
      <BambooLogo size={30} />
      <span className="font-serif text-lg font-semibold text-text-primary">
        {siteConfig.defaultName}
      </span>
    </Link>
  )
}

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
      // 注册后整页跳转，确保应用以干净状态重新装载；按管理员身份分流落地页
      window.location.href = resolveSafeRedirect(
        redirectTarget,
        res.user.is_admin,
      )
    } catch (err) {
      setError(err instanceof Error ? err.message : '注册失败，请稍后重试')
    }
  }

  return (
    <div className="grid min-h-dvh bg-background lg:grid-cols-2">
      {/* ───────── 左栏：水墨叙事面板（镜像，与登录页互为镜面） ───────── */}
      <motion.div
        initial={reduced ? false : { opacity: 0 }}
        animate={{ opacity: 1 }}
        transition={{ duration: 0.8, ease: 'easeOut' }}
        className="relative hidden overflow-hidden lg:block"
        style={{
          background:
            'radial-gradient(700px 400px at 25% 12%, oklch(0.88 0.1 105 / 0.26), transparent 68%), radial-gradient(500px 380px at 80% 90%, oklch(0.88 0.1 105 / 0.12), transparent 70%), oklch(0.955 0.02 112)',
        }}
      >
        {/* 桌面 Logo：承载视口左上（表单在右栏，故 Logo 落在水墨栏） */}
        <div className="absolute left-12 top-10 z-10">
          <SiteLogo />
        </div>
        {/* 衬线水印大字：贴左外边缘 */}
        <span
          aria-hidden
          className="pointer-events-none absolute left-[100px] top-10 select-none font-serif text-[180px] font-black leading-none text-text-primary opacity-[0.04]"
        >
          林
        </span>
        {/* 墨韵竹叶：左栏朝左外（镜像） */}
        <BambooArt
          mirror
          className="pointer-events-none absolute -left-10 top-0 h-full w-[560px] text-text-primary"
        />
        {/* 竹节竖线收边：贴右分界线 */}
        <span
          aria-hidden
          className="absolute bottom-[8%] right-0 top-[8%] w-px"
          style={{
            background:
              'linear-gradient(to bottom, transparent, var(--leaf-muted), transparent)',
            opacity: 0.5,
          }}
        />

        {/* 竖排题跋：贴右分界线内侧 */}
        <div className="absolute right-[72px] top-1/2 -translate-y-1/2 [writing-mode:vertical-rl] [text-orientation:upright]">
          <div className="flex">
            <p className="font-serif text-[52px] font-bold leading-[1.15] tracking-[0.28em] text-text-primary">
              竹林新友
            </p>
            <p className="mr-[18px] font-serif text-[34px] font-semibold tracking-[0.3em] text-leaf-deep">
              遇
            </p>
          </div>
          <p className="mr-4 font-serif text-[15px] leading-[2.4] tracking-[0.4em] text-text-secondary">
            以文会友，以链相连
          </p>
        </div>

        {/* 底部引言 */}
        <div className="absolute inset-x-0 bottom-14 pl-[60px] pr-[72px]">
          <p className="font-serif text-[17px] italic leading-relaxed tracking-[0.02em] text-text-primary/85">
            「{siteConfig.blogger.description}」
          </p>
          <p className="mt-2.5 font-mono text-[11px] uppercase tracking-[0.12em] text-text-secondary">
            {siteConfig.defaultName} · 友情链接管理
          </p>
        </div>
      </motion.div>

      {/* ───────── 右栏：表单（视口右，装饰朝右外边缘） ───────── */}
      <div className="relative flex flex-col gap-6 overflow-hidden p-6 md:p-10">
        {/* 晨光墨晕：镜像，重心偏右上 */}
        <div
          aria-hidden
          className="pointer-events-none absolute inset-0"
          style={{
            background:
              'radial-gradient(560px 220px at 82% 0%, oklch(0.88 0.1 105 / 0.20), transparent 72%), radial-gradient(420px 260px at 10% 100%, oklch(0.88 0.1 105 / 0.08), transparent 70%)',
          }}
        />
        {/* 墨韵竹叶水印：右栏朝右外，极淡 */}
        <BambooArt className="pointer-events-none absolute -bottom-10 -right-[60px] h-full w-[420px] text-text-primary opacity-50" />
        {/* 角标：承载视口右上（仅桌面） */}
        <span className="absolute right-12 top-10 hidden font-mono text-[11px] uppercase tracking-[0.18em] text-text-secondary/70 lg:block">
          Register · 竹林
        </span>

        {/* 顶部 Logo：仅移动（桌面 Logo 在左栏水墨面板） */}
        <div className="relative z-10 flex justify-center md:justify-start lg:hidden">
          <SiteLogo />
        </div>

        {/* 居中窄表单 */}
        <div className="relative z-10 flex flex-1 items-center justify-center">
          <motion.div
            initial={reduced ? false : { opacity: 0, y: 12 }}
            animate={{ opacity: 1, y: 0 }}
            transition={{ duration: 0.5, ease: 'easeOut' }}
            className="w-full max-w-[360px]"
          >
            <form className="flex flex-col" onSubmit={handleSubmit}>
              {/* 返回链接：标题上方的逃生口 */}
              <BackLink to="/" label="返回首页" className="mb-5" />

              {/* 标题组：衬线 + 笔刷下划线 */}
              <h1 className="font-serif text-[32px] font-bold leading-tight tracking-[0.01em] text-text-primary">
                注册账号
              </h1>
              <BrushUnderline className="mb-2.5 mt-2.5" />
              <p className="text-sm leading-relaxed text-text-secondary">
                注册后即可申请并管理你的友情链接
              </p>

              <div className="mt-6 grid gap-4">
                {/* 邮箱 + 发送验证码 */}
                <div className="grid gap-2">
                  <Label htmlFor="email">邮箱</Label>
                  <div className="flex gap-2">
                    <div className="relative flex-1">
                      <Mail className="absolute left-3 top-1/2 size-4 -translate-y-1/2 text-text-secondary/70" />
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
                    <ShieldCheck className="absolute left-3 top-1/2 size-4 -translate-y-1/2 text-text-secondary/70" />
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

                {/* 昵称（可选） */}
                <div className="grid gap-2">
                  <Label htmlFor="nickname">
                    昵称 <span className="text-text-secondary">（可选）</span>
                  </Label>
                  <div className="relative">
                    <AtSign className="absolute left-3 top-1/2 size-4 -translate-y-1/2 text-text-secondary/70" />
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
                    <LockKeyhole className="absolute left-3 top-1/2 size-4 -translate-y-1/2 text-text-secondary/70" />
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
                  {registerMutation.isPending ? '注册中…' : '注 册'}
                </Button>
              </div>

              {/* 底部：仅认证流程内切换 */}
              <div className="mt-6 flex justify-center text-sm">
                <Link
                  to="/auth/login"
                  className="font-medium text-leaf-deep transition-opacity hover:opacity-75"
                >
                  已有账号？去登录
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
    </div>
  )
}
