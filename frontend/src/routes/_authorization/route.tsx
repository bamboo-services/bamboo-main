/*
 * --------------------------------------------------------------------------------
 * Copyright (c) 2016-NOW 筱锋
 * Author: 筱锋(https://www.x-lf.com)
 * License: MIT
 * --------------------------------------------------------------------------------
 */

import { Outlet, createFileRoute } from '@tanstack/react-router'

export const Route = createFileRoute('/_authorization')({
  component: AuthorizationLayout,
})

function AuthorizationLayout() {
  return <Outlet />
}
