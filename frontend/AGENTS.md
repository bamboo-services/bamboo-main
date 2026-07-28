# FRONTEND 知识库

## 概述

友情链接管理的后台前端，独立于 Go 后端的子项目。基于 React 19 + TanStack Router（file-based routing）

- TanStack Query/Table + Tailwind v4 + shadcn/ui，由 Vite 7 驱动，包管理器固定为 pnpm。
  通过 `@` 路径别名引用 `src/`；路由树由 `@tanstack/router-plugin` 在开发/构建时自动生成
  `routeTree.gen.ts`，禁止手改。

## 目录结构

```text
frontend/
|- index.html                 # 入口 HTML（挂载点 #app）
|- vite.config.ts            # Vite 配置：tanstackRouter + react + tailwind 插件 + @ 别名
|- package.json              # pnpm 管理；scripts: dev/build/preview/test/lint/format/check
|- components.json            # shadcn/ui CLI 配置
|- tsconfig.json             # TS 配置
|- eslint.config.js          # ESLint（@tanstack/eslint-config）
|- prettier.config.js        # Prettier
|- public/                   # 静态资源
`- src/
   |- main.tsx               # 入口：RouterProvider + createRouter(routeTree) + reportWebVitals
   |- routeTree.gen.ts       # 自动生成，禁止手改
   |- styles.css            # Tailwind v4 入口
   |- routes/                # file-based routing
   |  |- __root.tsx          # 根路由：Outlet + TanStack Devtools
   |  |- _public/            # 公开页面路由组
   |  |- _admin/             # 管理后台布局路由组（SidebarProvider + Breadcrumb）
   |  |  `- admin/            #   dashboard / link/ / sponsor / setting 子页
   |  `- _authorization/      # 登录授权路由组（auth/login）
   |- components/
   |  |- ui/                 # shadcn/ui 组件（button/card/dialog/sidebar 等）
   |  `- layout/             # 布局组件（admin-sidebar）
   |- hooks/                 # 自定义 hooks（use-mobile）
   |- lib/utils.ts           # cn() class 合并助手
   |- data/                  # mock 数据（demo-table-data、mock/links、mock/site-info）
   `- assets/                # svg（bamboo-logo）、images
```

## 导航指南

| 任务           | 位置                                      | 说明                                             |
| -------------- | ----------------------------------------- | ------------------------------------------------ |
| 新增页面       | `src/routes/`                             | 按文件路径自动生成路由；放对位置即生效           |
| 路由分组/布局  | `src/routes/_*` 目录                      | `_admin`/`_public`/`_authorization` 为布局路由组 |
| 后台布局壳     | `src/routes/_admin/route.tsx`             | SidebarProvider + Breadcrumb + Outlet            |
| shadcn/ui 组件 | `src/components/ui/`                      | 由 `components.json` + shadcn CLI 维护           |
| 侧边栏导航     | `src/components/layout/admin-sidebar.tsx` | 管理后台侧边栏定义                               |
| 样式入口       | `src/styles.css`                          | Tailwind v4 `@import "tailwindcss"`              |
| 别名配置       | `vite.config.ts`、`tsconfig.json`         | `@` 指向 `./src`                                 |
| 路由树生成     | `src/routeTree.gen.ts`                    | `@tanstack/router-plugin` 自动生成               |

## 约定

- 包管理器固定为 **pnpm**（`pnpm-lock.yaml` 为真相源），不用 bun/npm/yarn
- 路由采用 file-based routing，新页面靠文件路径自动注册；`routeTree.gen.ts` 由插件生成，禁止手改
- `@` 别名指向 `src/`，导入内部模块一律走别名
- UI 基础组件统一走 shadcn/ui（`src/components/ui/`），用 CLI 增删而非手写
- 样式用 Tailwind v4（CSS-first，`@theme` 在 `styles.css`）；不写独立 CSS 文件
- 文件头部必须有版权横幅（TS/TSX/CSS/HTML），生成产物（`routeTree.gen.ts`）豁免
- 开发服务器固定 3000 端口（`vite --port 3000`）

## 反模式

- 手改 `routeTree.gen.ts`——下次生成会覆盖
- 在 `src/components/ui/` 手写非 shadcn 约定的组件——应放 `src/components/` 其他子目录
- 用 bun/npm/yarn 替代 pnpm 安装依赖——锁文件不一致
- 在路由文件中写与展示无关的副作用——应下沉到 hooks 或 lib
- 跳过 `cn()` 直接拼接 className——破坏样式合并约定

## 设计语言 · 竹林水墨

后台前端统一视觉语言为「竹林水墨 · 东方品牌签名」（淡绿宣纸 + 墨色衬线 + 东方签名元素）。
**设计语言与动画架构的完整 spec 沉淀在 [`ARCHITECTURE.md`](./ARCHITECTURE.md)**，包括：

- 色彩 token / 字体 / 布局与组件模式 / 品牌签名元素 / 反 AI slop 硬红线
- 动画分层模型（路由级 / 区块级 / 组件级 / 环境级）与四条核心铁律
- 入场编排 `enter()` 助手、CSS 环境动画、数据可视化组件动画约束

> 涉及任何视觉设计 / 动画 / 配色 / 新页面美化时，**直接读 [`ARCHITECTURE.md`](./ARCHITECTURE.md)**，
> 不在本文件重复内容。新视觉设计走 `huashu-design` skill 时，以 ARCHITECTURE.md 的 spec 为硬约束
> （约束器见 `.agents/skills/bamboo-ink-design/`）。

## 调试路径

1. 路由不生效/404：确认页面文件放在 `src/routes/` 对应路径下，保存后查看 `routeTree.gen.ts` 是否更新
2. 别名失效：核对 `vite.config.ts` 的 `resolve.alias` 与 `tsconfig.json` 的 `paths`
3. 样式丢失：检查 `src/styles.css` 是否被 `main.tsx` 导入，Tailwind v4 是否经 `@tailwindcss/vite` 插件接入
4. 依赖异常：用 `pnpm install` 重装，确认 `pnpm-lock.yaml` 与 `package.json` 一致
5. shadcn 组件缺失：用 `pnpm dlx shadcn@latest add <component>` 补齐（参照 `components.json`）
