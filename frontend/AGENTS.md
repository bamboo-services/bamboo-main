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

后台前端统一视觉语言为「竹林水墨 · 东方品牌签名」：以 dashboard 为标杆，淡绿宣纸底 + 墨色文字 +
衬线 display + 单一 leaf-deep accent，配墨韵竹叶 / 晨光墨晕 / 竹节线 / 竖排题跋等东方签名元素。
**任何前端新页面 / 组件 / 视觉设计都必须沿用这套语言**，禁止引入紫渐变 / emoji 图标 / 玻璃拟态等
通用 SaaS 审美。涉及新视觉设计时，走 `huashu-design` skill，但必须以本节 spec 为硬约束
（详见 `.agents/skills/bamboo-ink-design/`）。

### 色彩 token（`styles.css`，勿自创颜色）

| token | oklch | 用途 |
| --- | --- | --- |
| `--background` | `oklch(0.975 0.016 110)` | 宣纸页面底 |
| `--card` | `oklch(0.99 0.005 110)` | 卡片块底 |
| `--text-primary` / 墨色 | `oklch(0.32 0.06 155)` | 主文字、浓墨 |
| `--text-secondary` / 淡墨 | `oklch(0.5 0.05 130)` | 次级文字 |
| `--leaf-light` | `oklch(0.88 0.1 105)` | 浅叶绿（晨光墨晕） |
| `--leaf-muted` | `oklch(0.8 0.08 130)` | 柔叶绿（竹节线） |
| `--leaf-deep` | `oklch(0.55 0.12 155)` | **深叶绿，唯一 accent** |
| `--seal` | `oklch(0.55 0.18 25)` | 朱砂红（备用，当前未启用） |
| `--chart-1/4/5` | 见 styles.css | donut 分段（已通过/待审核/其他） |

### 字体

- **display（标题 / 大数字）**：衬线 `--font-serif`（`Source Serif 4` + `Noto Serif SC`，经
  `index.html` Google Fonts 引入，系统衬线 fallback）。Tailwind 用 `font-serif`。
- **正文 / label**：系统 sans（`-apple-system` 栈）。
- **数据 / meta**：`font-mono` + `tabular-nums`（数字严格对齐）。

### 布局与组件模式（标杆见 `dashboard.tsx`）

- **宣纸卡片块** `inkCard`：`rounded-lg border border-border bg-card p-5` + hover 浮起 + 墨边软阴影。
- **卡片标题** `CardHead`：斜墨条（`-skew-x-12` leaf-deep）+ 衬线标题 + 右侧 mono uppercase meta。
- **竹节分隔线** `BambooRule`：横线 + 圆/宽节点交替，分大章节用。
- **双模式 hero**：sidebar 展开 → 顶部横排大问候；sidebar 收起（`useSidebar().open === false`）→
  左侧竖排卷轴（`writing-mode: vertical-rl`）+ 右侧数据区，grid `grid-cols-[280px_1fr]`。

### 品牌签名元素

- **墨韵竹叶** `BambooArt`：内联 SVG 三竿竹，墨色 opacity 0.045–0.1，`ink-sway` 极慢摇曳
  （`transform-box: fill-box` 底部锚点）。
- **晨光墨晕**：单色淡绿径向渐变（仅 leaf-light），非多色光斑。
- **竖排题跋**：`[writing-mode:vertical-rl] [text-orientation:upright]` 大号衬线问候 + 导语。
- **笔刷下划线**：手绘 SVG path（leaf-deep）。
- **enso 缺口圆**：空状态图标（圆环留缺口 + 竹叶）。
- **衬线水印字**：KPI 卡片背景大字（opacity 5%）。

### 动画（`styles.css`，reduced-motion 全部静止）

- `ink-sway` / `ink-sway-slow` / `ink-sway-fast`：竹叶摇曳。
- `ink-pulse`：状态点呼吸。
- 入场用 `motion/react` 的淡入 + 轻微位移，克制不过度。

### 反 AI slop（硬红线）

- 禁紫色 / 靛蓝 / 粉色渐变、渐变填充标题。
- 禁 emoji 当图标（用内联 SVG 线条图标或纯文字）。
- 禁玻璃拟态（backdrop-blur）、aurora 多色光斑背景。
- 禁近黑背景 + 霓虹 accent（本项目是宣纸浅色主题）。
- 禁每个标题配圆角 icon 块。
- 禁编造假数据 / 装饰性 stats。
- **印章元素已弃用**（用户决策），勿再加朱砂印章。

### 参考实现

- 标杆页面：`src/routes/_admin/admin/dashboard.tsx`（水墨叙事 + 双模式 hero）。
- 共享组件：`src/components/dashboard/`（count-up / donut-chart / radial-gauge）。
- token 与动画：`src/styles.css`。

## 调试路径

1. 路由不生效/404：确认页面文件放在 `src/routes/` 对应路径下，保存后查看 `routeTree.gen.ts` 是否更新
2. 别名失效：核对 `vite.config.ts` 的 `resolve.alias` 与 `tsconfig.json` 的 `paths`
3. 样式丢失：检查 `src/styles.css` 是否被 `main.tsx` 导入，Tailwind v4 是否经 `@tailwindcss/vite` 插件接入
4. 依赖异常：用 `pnpm install` 重装，确认 `pnpm-lock.yaml` 与 `package.json` 一致
5. shadcn 组件缺失：用 `pnpm dlx shadcn@latest add <component>` 补齐（参照 `components.json`）
