/*
 * --------------------------------------------------------------------------------
 * Copyright (c) 2016-NOW 筱锋
 * Author: 筱锋(https://www.x-lf.com)
 * License: MIT
 * --------------------------------------------------------------------------------
 */

import { Outlet, createFileRoute, redirect } from '@tanstack/react-router'

export const Route = createFileRoute('/_admin/admin')({
  beforeLoad: ({ location }) => {
    if (location.pathname === '/admin' || location.pathname === '/admin/') {
      throw redirect({ to: '/admin/dashboard' })
    }
  },
  component: AdminRouteLayout,
})

function AdminRouteLayout() {
  return <Outlet />
}
