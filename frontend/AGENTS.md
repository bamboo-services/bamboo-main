# FRONTEND 知识库

## 概述

友情链接管理的前端子项目，独立于 Go 后端。基于 **React 19.2 + TanStack Router（file-based routing）+ TanStack Query/Table + Tailwind v4 + shadcn/ui**，由 **Vite 8** 驱动，包管理器固定为 **pnpm**。构建产物输出到 `../resources/frontend/dist`，由 Go 后端 `go:embed` 内嵌实现单二进制部署。通过 `@` 路径别名引用 `src/`；路由树由 `@tanstack/router-plugin` 在开发/构建时自动生成 `routeTree.gen.ts`，禁止手改。

## 目录结构

```text
frontend/
|- index.html                 # 入口 HTML（Google Fonts：Source Serif 4 + Noto Serif SC）
|- vite.config.ts            # Vite 配置：tanstackRouter + react + tailwind + @ 别名；outDir=../resources/frontend/dist；dev proxy /api、/screenshots → 5555
|- package.json              # pnpm 管理；scripts: dev/build/preview/test/lint/format/check
|- components.json            # shadcn/ui CLI 配置
|- tsconfig.json             # TS 配置（strict，@/* → ./src/*）
|- eslint.config.js          # ESLint（@tanstack/eslint-config）
|- prettier.config.js        # Prettier
|- ARCHITECTURE.md           # 设计语言与动画架构单一真相源
|- public/                   # 静态资源（favicon 等）
`- src/
   |- main.tsx               # 入口：createRouter + QueryClient + RouterProvider
   |- routeTree.gen.ts       # 自动生成，禁止手改
   |- styles.css            # Tailwind v4 入口 + 全部色彩 token（oklch）+ keyframes
   |- routes/                # file-based routing（_ 前缀 = 布局组）
   |  |- __root.tsx          # 根路由：Outlet + Toaster + devtools
   |  |- _public/            # 公开站布局组（首页友链墙）
   |  |- _admin/             # 管理后台布局组（admin-sidebar + 面包屑 + beforeLoad 鉴权）
   |  |- _authorization/     # 鉴权布局组（auth/login、register、callback）
   |  |- _user/              # 用户中心布局组（dashboard、account、links、sponsors）
   |  |- about/              # 关于页组（me 名士帖、friends 友链墙、sponsor 赞助墙）
   |  `- operate/            # 自助申请组（apply 友链申请、sponsor 赞助申请）
   |- components/
   |  |- ui/                 # shadcn/ui 组件（avatar/badge/dialog/sidebar/table 等 21 个）
   |  |- ink-wash.tsx        # 【核心视觉原语库】inkCard/CardHead/BambooRule/BambooArt 等
   |  |- layout/             # 布局（admin-sidebar、account-hover-card）
   |  |- about/              # 友链卡片系列（regular/close/premium/ad + interlude + lazy-image）
   |  |- link/               # 友链展示（accent-bar、site-avatar、ranking-board）
   |  |- dashboard/          # 数据可视化（count-up、donut-chart、radial-gauge）
   |  |- decorative/         # 环境装饰（falling-leaves）
   |  |- markdown.tsx        # MarkdownView 渲染（react-markdown + remark-gfm）
   |  |- markdown-editor.tsx # MarkdownEditor 编辑器
   |  |- link-form.tsx       # 管理端友链表单
   |  |- sponsor-apply-form.tsx # 访客赞助申请表单
   |  `- user-link-form.tsx  # 访客友链自助申请表单
   |- api/                   # axios 请求层：client + auth/color/dashboard/group/info/link/public/sponsor/types
   |- hooks/                 # TanStack Query hooks：use-auth/colors/dashboard/groups/links/site-info/sponsors/mobile
   |- lib/                   # 纯逻辑：auth(会话)/colors(颜色视觉)/datetime/friend-groups/locations/motion/role/site/utils
   |- assets/                # svg（bamboo-logo）、images（登录背景/默认背景）
   `- reportWebVitals.ts     # Web Vitals 性能上报
```

## 导航指南

| 任务 | 位置 | 说明 |
|---|---|---|
| 新增页面 | `src/routes/` | 按文件路径自动生成路由；放对位置即生效 |
| 路由分组/布局 | `src/routes/_*` 目录 | `_admin`/`_user`/`_authorization`/`_public` 为布局路由组 |
| 后台布局壳 | `src/routes/_admin/route.tsx` | AdminSidebar + Breadcrumb + Outlet，beforeLoad 校验 admin |
| 用户中心壳 | `src/routes/_user/route.tsx` | 顶导 + 头像下拉，beforeLoad 校验登录 |
| 请求层 | `src/api/` | `client.ts` 统一 axios 封装，各领域 api 文件 |
| 数据层 hooks | `src/hooks/` | TanStack Query 封装，queryKey 工厂统一管理缓存 |
| 纯逻辑 | `src/lib/` | 颜色视觉/分级/会话/日期等无副作用逻辑 |
| 竹林水墨原语 | `src/components/ink-wash.tsx` | 共享视觉组件单一来源 |
| shadcn/ui 组件 | `src/components/ui/` | 由 `components.json` + shadcn CLI 维护 |
| 侧边栏导航 | `src/components/layout/admin-sidebar.tsx` | 管理后台侧边栏定义 |
| 样式入口 | `src/styles.css` | Tailwind v4 + 色彩 token + keyframes |
| 别名配置 | `vite.config.ts`、`tsconfig.json` | `@` 指向 `./src` |
| 路由树生成 | `src/routeTree.gen.ts` | `@tanstack/router-plugin` 自动生成 |

## 约定

- 包管理器固定为 **pnpm**（`pnpm-lock.yaml` 为真相源），不用 bun/npm/yarn
- 路由采用 file-based routing，新页面靠文件路径自动注册；`routeTree.gen.ts` 由插件生成，禁止手改
- `@` 别名指向 `src/`，导入内部模块一律走别名
- UI 基础组件统一走 shadcn/ui（`src/components/ui/`），用 CLI 增删而非手写
- 数据请求一律经 `src/api/*` + `src/hooks/*`，不在组件内裸写 fetch；queryKey 走各 hook 的工厂函数保证一致性
- 样式用 Tailwind v4（CSS-first，`@theme`/token 在 `styles.css`）；不写独立 CSS 文件
- 纯展示常量（颜色视觉、友链分级、内置分组）放 `src/lib/`，不在组件内散落 magic number
- 文件头部必须有版权横幅（TS/TSX/CSS/HTML），生成产物（`routeTree.gen.ts`）豁免
- 开发服务器固定 3000 端口（`vite --port 3000`），API 经 Vite proxy 转发到 5555

## 反模式

- 手改 `routeTree.gen.ts`——下次生成会覆盖
- 在 `src/components/ui/` 手写非 shadcn 约定的组件——应放 `src/components/` 其他子目录
- 用 bun/npm/yarn 替代 pnpm 安装依赖——锁文件不一致
- 在路由文件中写与展示无关的副作用——应下沉到 hooks 或 lib
- 跳过 `cn()` 直接拼接 className——破坏样式合并约定
- 自创颜色 token 或渐变——色彩必须取自 `styles.css` 的竹林水墨 palette
- 在组件内裸写 fetch/axios——应走 `src/api/` 请求层

## 设计语言 · 竹林水墨

前端统一视觉语言为「竹林水墨 · 东方品牌签名」（淡绿宣纸 + 墨色衬线 + 东方签名元素）。
**设计语言与动画架构的完整 spec 沉淀在 [`ARCHITECTURE.md`](./ARCHITECTURE.md)**，包括：

- 色彩 token（oklch，禁止自创） / 字体 / 布局与组件模式 / 品牌签名元素 / 反 AI slop 硬红线
- 动画分层模型（路由级 opacity / 区块级 opacity+y / 组件级静态 / 环境级 keyframes）与四条核心铁律
- 入场编排 `enter()` 助手、`reduced-motion` 全降级约定、数据可视化组件静态约束

> 涉及任何视觉设计 / 动画 / 配色 / 新页面美化时，**直接读 [`ARCHITECTURE.md`](./ARCHITECTURE.md)**，
> 不在本文件重复内容。新视觉设计走 `huashu-design` skill 时，以 ARCHITECTURE.md 的 spec 为硬约束
> （约束器见 `.agents/skills/bamboo-ink-design/`）。

## 调试路径

1. 路由不生效/404：确认页面文件放在 `src/routes/` 对应路径下，保存后查看 `routeTree.gen.ts` 是否更新
2. 别名失效：核对 `vite.config.ts` 的 `resolve.alias` 与 `tsconfig.json` 的 `paths`
3. 样式丢失：检查 `src/styles.css` 是否被 `main.tsx` 导入，Tailwind v4 是否经 `@tailwindcss/vite` 插件接入
4. 依赖异常：用 `pnpm install` 重装，确认 `pnpm-lock.yaml` 与 `package.json` 一致
5. shadcn 组件缺失：用 `pnpm dlx shadcn@latest add <component>` 补齐（参照 `components.json`）
6. 接口请求异常：dev 环境核对 Vite proxy（/api → 5555），生产核对内嵌产物是否刷新（`make build-frontend`）
7. 页面在单二进制下 404：确认已执行 `make build-frontend` 产出 `resources/frontend/dist`（构建产物不入库）
