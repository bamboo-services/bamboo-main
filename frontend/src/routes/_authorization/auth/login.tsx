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

import { createFileRoute, Link, useNavigate } from '@tanstack/react-router'
import { useState } from 'react'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { Checkbox } from '@/components/ui/checkbox'
import { siteInfo } from '@/data/mock/site-info'
import favicon from '@/assets/images/favicon.png'
import authBackground from '@/assets/images/auth_background.jpg'
import { User, Key, Send } from 'lucide-react'

export const Route = createFileRoute('/_authorization/auth/login')({
  component: LoginPage,
})

function LoginPage() {
  const navigate = useNavigate()
  const [formData, setFormData] = useState({
    username: '',
    password: '',
    remember: false,
  })

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault()
    // 静态界面，直接跳转到管理后台
    navigate({ to: '/admin/dashboard' })
  }

  return (
    <section className="bg-white min-h-dvh">
      <div className="lg:grid lg:min-h-screen lg:grid-cols-12">
        {/* 左侧背景图区域 */}
        <section className="relative flex h-32 items-end bg-gray-900 lg:col-span-5 lg:h-full xl:col-span-6">
          <img
            alt="Background"
            src={authBackground}
            className="absolute inset-0 h-full w-full object-cover opacity-80"
          />
          <div className="hidden lg:relative lg:block lg:p-12">
            <div className="flex">
              <Link to="/">
                <img
                  src={favicon}
                  alt="Logo"
                  className="rounded-3xl h-32 w-32 shadow-lg hover:scale-105 transition"
                />
              </Link>
            </div>
            <h2 className="mt-6 text-2xl font-bold text-white sm:text-3xl md:text-4xl">
              {siteInfo.site.siteName}
            </h2>
            <p className="mt-4 leading-relaxed text-white/90">
              {siteInfo.blogger.description}
            </p>
          </div>
        </section>

        {/* 右侧登录表单区域 */}
        <main className="flex items-center justify-center px-8 py-8 sm:px-12 lg:col-span-7 lg:px-16 lg:py-12 xl:col-span-6">
          <div className="max-w-xl lg:max-w-3xl w-full">
            {/* 移动端 Logo */}
            <div className="relative -mt-16 block lg:hidden">
              <Link
                to="/"
                className="inline-flex items-center justify-center rounded-full bg-white p-2 shadow-lg"
              >
                <img src={favicon} alt="Logo" className="rounded-2xl h-12" draggable={false} />
              </Link>
              <h1 className="mt-2 text-2xl font-bold text-gray-900 sm:text-3xl md:text-4xl">
                {siteInfo.site.siteName}
              </h1>
              <p className="mt-4 leading-relaxed text-gray-500">
                {siteInfo.blogger.description}
              </p>
            </div>

            {/* 登录表单 */}
            <form className="mt-8 grid gap-6" onSubmit={handleSubmit}>
              <div className="grid justify-center mb-6">
                <h2 className="text-4xl font-bold text-center">用户登录</h2>
              </div>

              {/* 用户名 */}
              <div className="grid gap-2">
                <Label htmlFor="username" className="flex gap-1">
                  <span>用户名</span>
                  <span className="text-xs text-red-500">*</span>
                </Label>
                <div className="relative">
                  <User className="absolute left-3 top-1/2 -translate-y-1/2 h-4 w-4 text-muted-foreground" />
                  <Input
                    id="username"
                    type="text"
                    placeholder="请输入用户名"
                    className="pl-10"
                    value={formData.username}
                    onChange={(e) => setFormData({ ...formData, username: e.target.value })}
                  />
                </div>
              </div>

              {/* 密码 */}
              <div className="grid gap-2">
                <Label htmlFor="password" className="flex gap-1">
                  <span>密码</span>
                  <span className="text-xs text-red-500">*</span>
                </Label>
                <div className="relative">
                  <Key className="absolute left-3 top-1/2 -translate-y-1/2 h-4 w-4 text-muted-foreground" />
                  <Input
                    id="password"
                    type="password"
                    placeholder="请输入密码"
                    className="pl-10"
                    value={formData.password}
                    onChange={(e) => setFormData({ ...formData, password: e.target.value })}
                  />
                </div>
              </div>

              {/* 记住登录 */}
              <div className="flex items-center gap-2">
                <Checkbox
                  id="remember"
                  checked={formData.remember}
                  onCheckedChange={(checked) =>
                    setFormData({ ...formData, remember: checked as boolean })
                  }
                />
                <Label htmlFor="remember" className="text-sm font-normal cursor-pointer">
                  记住登录
                </Label>
              </div>

              {/* 登录按钮 */}
              <div className="grid justify-center mt-2">
                <Button type="submit" size="lg" className="px-12">
                  <Send className="mr-2 h-4 w-4" />
                  登录
                </Button>
              </div>
            </form>
          </div>
        </main>
      </div>
    </section>
  )
}
