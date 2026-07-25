/*
 * --------------------------------------------------------------------------------
 * Copyright (c) 2016-NOW 筱锋
 * Author: 筱锋(https://www.x-lf.com)
 * License: MIT
 * --------------------------------------------------------------------------------
 */

import { createFileRoute, Link } from '@tanstack/react-router'
import { Button } from '@/components/ui/button'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { Badge } from '@/components/ui/badge'
import { mockLinks } from '@/data/mock/links'
import { ArrowLeft, Check, X, ExternalLink } from 'lucide-react'

export const Route = createFileRoute('/_admin/admin/link/verify')({
  component: LinkVerifyPage,
})

function LinkVerifyPage() {
  const pendingLinks = mockLinks.filter((link) => link.status === 'pending')

  return (
    <div className="space-y-6">
      <div className="flex items-center gap-4">
        <Link to="/admin/link">
          <Button variant="ghost" size="icon">
            <ArrowLeft className="h-4 w-4" />
          </Button>
        </Link>
        <div>
          <h1 className="text-3xl font-bold tracking-tight">友链审核</h1>
          <p className="text-muted-foreground">
            待审核的友链申请 ({pendingLinks.length})
          </p>
        </div>
      </div>

      {pendingLinks.length === 0 ? (
        <Card>
          <CardContent className="flex flex-col items-center justify-center py-12">
            <p className="text-muted-foreground">暂无待审核的友链申请</p>
            <Link to="/admin/link" className="mt-4">
              <Button variant="outline">返回列表</Button>
            </Link>
          </CardContent>
        </Card>
      ) : (
        <div className="space-y-4">
          {pendingLinks.map((link) => (
            <Card key={link.id}>
              <CardHeader>
                <div className="flex items-start justify-between">
                  <div className="flex items-center gap-4">
                    <div className="h-16 w-16 rounded-lg bg-muted flex items-center justify-center overflow-hidden">
                      <img
                        src={link.siteLogo}
                        alt={link.siteName}
                        className="h-full w-full object-cover"
                        onError={(e) => {
                          e.currentTarget.style.display = 'none'
                        }}
                      />
                    </div>
                    <div>
                      <CardTitle>{link.siteName}</CardTitle>
                      <CardDescription className="flex items-center gap-1 mt-1">
                        <a
                          href={link.siteUrl}
                          target="_blank"
                          rel="noopener noreferrer"
                          className="hover:underline flex items-center gap-1"
                        >
                          {link.siteUrl}
                          <ExternalLink className="h-3 w-3" />
                        </a>
                      </CardDescription>
                    </div>
                  </div>
                  <Badge variant="secondary">待审核</Badge>
                </div>
              </CardHeader>
              <CardContent className="space-y-4">
                <div>
                  <p className="text-sm font-medium mb-1">站点描述</p>
                  <p className="text-sm text-muted-foreground">{link.siteDescription}</p>
                </div>
                <div className="grid gap-4 md:grid-cols-3">
                  <div>
                    <p className="text-sm font-medium mb-1">站长邮箱</p>
                    <p className="text-sm text-muted-foreground">{link.webmasterEmail}</p>
                  </div>
                  <div>
                    <p className="text-sm font-medium mb-1">申请位置</p>
                    <p className="text-sm text-muted-foreground">{link.locationName}</p>
                  </div>
                  <div>
                    <p className="text-sm font-medium mb-1">申请时间</p>
                    <p className="text-sm text-muted-foreground">{link.createdAt}</p>
                  </div>
                </div>
                <div className="flex justify-end gap-2 pt-4 border-t">
                  <Button variant="outline" className="text-destructive hover:text-destructive">
                    <X className="mr-2 h-4 w-4" />
                    拒绝
                  </Button>
                  <Button>
                    <Check className="mr-2 h-4 w-4" />
                    通过
                  </Button>
                </div>
              </CardContent>
            </Card>
          ))}
        </div>
      )}
    </div>
  )
}
