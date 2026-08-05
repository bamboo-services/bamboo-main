// --------------------------------------------------------------------------------
// Copyright (c) 2016-NOW(至今) 筱锋
// Author: 筱锋「xiao_lfeng」(https://www.x-lf.com)
// --------------------------------------------------------------------------------
// 许可证声明：版权所有 (c) 2016-2026 筱锋。保留所有权利。
// 有关MIT许可证的更多信息，请查看项目根目录下的LICENSE文件或访问：
// https://opensource.org/licenses/MIT
// --------------------------------------------------------------------------------

import { createFileRoute } from '@tanstack/react-router'
import { motion, useReducedMotion } from 'motion/react'
import { useEffect, useState } from 'react'
import { toast } from 'sonner'
import type { FormEvent } from 'react'
import { useAuth } from '@/hooks/use-auth'
import { useUpdateProfile } from '@/hooks/use-links'
import { changePassword } from '@/api/auth'
import { CardHead, InkBadge, PageHead, inkCard } from '@/components/ink-wash'
import { Skeleton } from '@/components/ui/skeleton'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { enter } from '@/lib/motion'

export const Route = createFileRoute('/_user/user/account')({
  component: AccountPage,
})

function AccountPage() {
  const reduced = useReducedMotion() ?? false
  const { user, isLoading } = useAuth()
  const updateProfile = useUpdateProfile()

  // 资料表单（用户加载完成后预填）
  const [profile, setProfile] = useState({ username: '', nickname: '', avatar: '' })
  useEffect(() => {
    if (user) {
      setProfile({
        username: user.username ?? '',
        nickname: user.nickname ?? '',
        avatar: user.avatar ?? '',
      })
    }
  }, [user])

  const handleProfileSubmit = (e: FormEvent) => {
    e.preventDefault()
    updateProfile.mutate({
      username: profile.username.trim() || undefined,
      nickname: profile.nickname.trim() || undefined,
      avatar: profile.avatar.trim() || undefined,
    })
  }

  // 修改密码表单
  const [pwd, setPwd] = useState({ old: '', next: '', confirm: '' })
  const [pwdLoading, setPwdLoading] = useState(false)

  const handlePasswordSubmit = async (e: FormEvent) => {
    e.preventDefault()
    if (pwd.next !== pwd.confirm) {
      toast.error('两次输入的新密码不一致')
      return
    }
    setPwdLoading(true)
    try {
      await changePassword({
        old_password: pwd.old,
        new_password: pwd.next,
      })
      toast.success('密码修改成功')
      setPwd({ old: '', next: '', confirm: '' })
    } catch (err) {
      toast.error(err instanceof Error ? err.message : '密码修改失败')
    } finally {
      setPwdLoading(false)
    }
  }

  return (
    <div className="mx-auto max-w-3xl">
      <PageHead
        kicker="account · profile"
        title="账号设置"
        sub="管理你的个人资料与登录密码。"
      />

      {/* ═══════════ 资料 / 密码 两区块 ═══════════ */}
      <div className="mt-8 flex flex-col gap-6">
        {/* 个人资料 */}
        <motion.div
          {...enter(reduced, 0.15, {
            initial: { opacity: 0, y: 20 },
            animate: { opacity: 1, y: 0 },
            transition: { duration: 0.6, ease: 'easeOut' },
          })}
          className={`${inkCard} p-6 md:p-8`}
        >
          <CardHead title="个人资料" meta="profile" />
          {isLoading || !user ? (
            <div className="space-y-3">
              <Skeleton className="h-9 w-full" />
              <Skeleton className="h-9 w-full" />
            </div>
          ) : (
            <form onSubmit={handleProfileSubmit} className="space-y-5">
              {/* 邮箱（只读 + 验证状态） */}
              <div className="space-y-2">
                <Label>邮箱</Label>
                <div className="flex items-center gap-3">
                  <Input value={user.email} disabled className="bg-muted/40" />
                  {user.email_verify ? (
                    <InkBadge tone="leaf" className="shrink-0">
                      已验证
                    </InkBadge>
                  ) : (
                    <InkBadge tone="pending" className="shrink-0">
                      未验证
                    </InkBadge>
                  )}
                </div>
                {!user.email_verify && (
                  <p className="text-xs text-text-secondary">
                    该邮箱尚未验证（SSO 账号可在首次验证后获得标识）。
                  </p>
                )}
              </div>

              <div className="space-y-2">
                <Label htmlFor="username">用户名</Label>
                <Input
                  id="username"
                  placeholder="登录使用的用户名"
                  value={profile.username}
                  onChange={(e) =>
                    setProfile({ ...profile, username: e.target.value })
                  }
                />
                <p className="text-xs text-text-secondary">
                  用户名是登录凭证，修改后请使用新用户名登录。
                </p>
              </div>

              <div className="space-y-2">
                <Label htmlFor="nickname">昵称</Label>
                <Input
                  id="nickname"
                  placeholder="给自己起个昵称吧"
                  value={profile.nickname}
                  onChange={(e) =>
                    setProfile({ ...profile, nickname: e.target.value })
                  }
                />
              </div>

              <div className="space-y-2">
                <Label htmlFor="avatar">头像 URL</Label>
                <Input
                  id="avatar"
                  placeholder="https://example.com/avatar.jpg"
                  value={profile.avatar}
                  onChange={(e) =>
                    setProfile({ ...profile, avatar: e.target.value })
                  }
                />
              </div>

              <div className="flex justify-end border-t border-border/60 pt-5">
                <Button
                  type="submit"
                  className="cursor-pointer"
                  disabled={updateProfile.isPending}
                >
                  {updateProfile.isPending ? '保存中…' : '保存资料'}
                </Button>
              </div>
            </form>
          )}
        </motion.div>

        {/* 修改密码 */}
        <motion.div
          {...enter(reduced, 0.25, {
            initial: { opacity: 0, y: 20 },
            animate: { opacity: 1, y: 0 },
            transition: { duration: 0.6, ease: 'easeOut' },
          })}
          className={`${inkCard} p-6 md:p-8`}
        >
          <CardHead title="修改密码" meta="security" />
          <form onSubmit={handlePasswordSubmit} className="space-y-5">
            <div className="space-y-2">
              <Label htmlFor="oldPwd">当前密码</Label>
              <Input
                id="oldPwd"
                type="password"
                autoComplete="current-password"
                placeholder="请输入当前密码"
                value={pwd.old}
                onChange={(e) => setPwd({ ...pwd, old: e.target.value })}
                required
              />
            </div>
            <div className="grid gap-4 md:grid-cols-2">
              <div className="space-y-2">
                <Label htmlFor="newPwd">新密码</Label>
                <Input
                  id="newPwd"
                  type="password"
                  autoComplete="new-password"
                  placeholder="至少 6 位"
                  value={pwd.next}
                  onChange={(e) => setPwd({ ...pwd, next: e.target.value })}
                  required
                />
              </div>
              <div className="space-y-2">
                <Label htmlFor="confirmPwd">确认新密码</Label>
                <Input
                  id="confirmPwd"
                  type="password"
                  autoComplete="new-password"
                  placeholder="再次输入新密码"
                  value={pwd.confirm}
                  onChange={(e) => setPwd({ ...pwd, confirm: e.target.value })}
                  required
                />
              </div>
            </div>
            <div className="flex justify-end border-t border-border/60 pt-5">
              <Button
                type="submit"
                className="cursor-pointer"
                disabled={pwdLoading}
              >
                {pwdLoading ? '修改中…' : '修改密码'}
              </Button>
            </div>
          </form>
        </motion.div>
      </div>
    </div>
  )
}
