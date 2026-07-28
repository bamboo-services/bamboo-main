# FRONTEND 设计与动画架构

> 本文档是前端「设计语言」与「动画架构」的**单一真相源**。
> `AGENTS.md` 仅作指引，不在此重复内容；任何视觉 / 动画决策以本文为准。

---

## 一、设计语言 · 竹林水墨

后台前端统一视觉语言为「竹林水墨 · 东方品牌签名」：以 dashboard 为标杆，淡绿宣纸底 + 墨色文字 +
衬线 display + 单一 leaf-deep accent，配墨韵竹叶 / 晨光墨晕 / 竹节线 / 竖排题跋等东方签名元素。
**任何前端新页面 / 组件 / 视觉设计都必须沿用这套语言**，禁止引入紫渐变 / emoji 图标 / 玻璃拟态等
通用 SaaS 审美。涉及新视觉设计时，走 `huashu-design` skill，但必须以本节 spec 为硬约束
（详见 `.agents/skills/bamboo-ink-design/`）。

### 1.1 气质定调

竹林清晨的雅致、克制、有书卷气。淡绿宣纸底 + 墨色文字 + 衬线 display + 单一 leaf-deep accent，
配墨韵竹叶 / 晨光墨晕 / 竹节线 / 竖排题跋等东方签名。**像一份精心排印的晨报 / 书画题跋，
不是通用 SaaS dashboard。**

### 1.2 色彩 token（`src/styles.css`，禁止自创颜色）

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

### 1.3 字体

- **display（标题 / 大数字）**：衬线 `--font-serif`（`Source Serif 4` + `Noto Serif SC`，经
  `index.html` Google Fonts 引入，系统衬线 fallback）。Tailwind 用 `font-serif`。
- **正文 / label**：系统 sans（`-apple-system` 栈）。
- **数据 / meta**：`font-mono` + `tabular-nums`（数字严格对齐）。

### 1.4 布局与组件模式（标杆见 `dashboard.tsx`）

- **宣纸卡片块** `inkCard`：`rounded-lg border border-border bg-card p-5` + hover 浮起 + 墨边软阴影。
- **卡片标题** `CardHead`：斜墨条（`-skew-x-12` leaf-deep）+ 衬线标题 + 右侧 mono uppercase meta。
- **竹节分隔线** `BambooRule`：横线 + 圆/宽节点交替，分大章节用。
- **双模式 hero**：sidebar 展开 → 顶部横排大问候；sidebar 收起（`useSidebar().open === false`）→
  左侧竖排卷轴（`writing-mode: vertical-rl`）+ 右侧数据区，grid `grid-cols-[280px_1fr]`。

### 1.5 品牌签名元素

- **墨韵竹叶** `BambooArt`：内联 SVG 三竿竹，墨色 opacity 0.045–0.1，`ink-sway` 极慢摇曳
  （`transform-box: fill-box` 底部锚点）。
- **晨光墨晕**：单色淡绿径向渐变（仅 leaf-light），非多色光斑。
- **竖排题跋**：`[writing-mode:vertical-rl] [text-orientation:upright]` 大号衬线问候 + 导语。
- **笔刷下划线**：手绘 SVG path（leaf-deep）。
- **enso 缺口圆**：空状态图标（圆环留缺口 + 竹叶）。
- **衬线水印字**：KPI 卡片背景大字（opacity 5%）。

### 1.6 反 AI slop（硬红线）

- 禁紫色 / 靛蓝 / 粉色渐变、渐变填充标题。
- 禁 emoji 当图标（用内联 SVG 线条图标或纯文字）。
- 禁玻璃拟态（backdrop-blur）、aurora 多色光斑背景。
- 禁近黑背景 + 霓虹 accent（本项目是宣纸浅色主题）。
- 禁每个标题配圆角 icon 块。
- 禁编造假数据 / 装饰性 stats。
- **印章元素已弃用**（用户决策），勿再加朱砂印章。

---

## 二、动画架构

前端动画分三层，各司其职。**核心铁律：跨层不得叠加动画同一 CSS 属性**——历史 bug「二次上滑」
正是路由级 `y` 位移与区块级 `y` 位移叠加、以及 `transition-all` 把 motion 驱动的 `transform`/`opacity`
卷入 CSS 过渡所致。下文每一层都给出规避规则。

### 2.1 分层模型

| 层级 | 载体 | 职责 | 允许的属性 | 关键约束 |
| --- | --- | --- | --- | --- |
| **路由级** | `src/routes/_admin/route.tsx` 的 `<motion.div key={pathname}>` | 页面切换整体淡入 | **仅 `opacity`** | 禁 `y` 位移，否则与各页区块级上滑嵌套叠加 |
| **区块级** | 各页 `motion.section` + `enter()` 助手 | 卡片错峰入场 | `opacity` + `y` | 入场动画只允许出现在本层 |
| **组件级** | `src/components/dashboard/*` | 数据可视化呈现 | `strokeDasharray`（静态）/ raf 计数 | **禁自生长动画**，遵守单一入场原则 |
| **环境级** | `src/styles.css` CSS keyframes | 持续氛围 | `transform: rotate` / `box-shadow` | reduced-motion 全静止，不参与入场 |

### 2.2 入场编排（区块级）

统一走 `src/lib/motion.ts` 的 `enter(reduced, delay, full)` 助手，禁止各页各自重复定义：

```ts
// 正常：原样返回 full 动画 + 叠加 delay
// reduced-motion：退化为纯 opacity 快速淡入（delay × 0.08）
export function enter(reduced: boolean, delay: number, full: MotionDivProps): MotionDivProps
```

- 错峰 delay 递增（如 dashboard 的 KPI/友链构成/系统状态/最近申请 = 0.22 / 0.3 / 0.36 / 0.44）。
- 单段 `duration: 0.4`，`ease: 'easeOut'`，位移幅度 `y: 10`（克制，不过度）。
- reduced-motion 时整体压缩为纯 opacity 淡入，无位移。

### 2.3 路由级过渡（`_admin/route.tsx`）

包裹 `<Outlet />` 的 `motion.div` **只做 `opacity` 淡入，不做 `y` 位移**。原因：路由级 `y` 与
区块级 `enter()` 的 `y` 是嵌套 transform，会叠加出「整页先上滑、卡片又上滑」的二次位移观感。
页面切换只留一段柔和淡入，位移交给区块级编排。

```tsx
<motion.div
  key={pathname}
  initial={{ opacity: 0 }}
  animate={{ opacity: 1 }}
  transition={{ duration: 0.25, ease: 'easeOut' }}
>
  <Outlet />
</motion.div>
```

### 2.4 CSS 环境动画（`styles.css`）

均为**持续循环**氛围动画，不参与入场；`prefers-reduced-motion: reduce` 下全部 `animation: none`：

- `ink-sway` / `ink-sway-slow` / `ink-sway-fast`：竹叶摇曳（`transform-box: fill-box` 底部锚点，
  13/17/11s，alternate）。
- `ink-pulse`：状态点呼吸（`box-shadow` 扩散，2.6s）。
- `ink-stamp`：朱砂落印（**已弃用**，保留定义但勿再启用）。

### 2.5 数据可视化组件（`src/components/dashboard/`）

| 组件 | 机制 | 动画 | 约束 |
| --- | --- | --- | --- |
| `count-up.tsx` | raf + easeOutCubic 数字滚动 | 非 motion，1s 计数 | reduced-motion 跳终值；属「数值更新」非入场 |
| `donut-chart.tsx` | 纯 SVG，`pathLength=1` + `strokeDasharray` | **无** | 数据到达即直接画到位，禁自生长 |
| `radial-gauge.tsx` | 纯 SVG，270° 弧 + `strokeDasharray` | **无** | 同上，禁自生长 |

> 历史教训：donut/gauge 曾用 `useMotionValue` + `animate()` 做弧线生长动画（0.9~1s），在数据到达后
> 才触发，与区块级入场动画叠加成「二次加载」。现统一移除——**入场只走区块级 `enter()` 一层**。

### 2.6 核心铁律

> 这些规则是用真实 bug 换来的，违反任一条都会重现「二次上滑 / 二次加载」。

#### 铁律 1 · 单一入场原则

每个视觉单元**只允许一段入场动画**。入场动画只出现在区块级 `enter()` 层；数据可视化组件
（donut/gauge 等）**不得自带生长/扫掠动画**，否则数据异步到达时会叠加出「第二段」入场。

- 反例：donut 弧线在 section 上滑后才开始生长 → 视觉上「卡片先滑进来，图表又画一遍」。
- 正解：图表数据到达即静态画到位，入场交给外层 section。

#### 铁律 2 · motion 与 CSS transition 不叠加同属性

一个元素若由 motion 驱动 `transform`/`opacity` 入场，**禁止**再挂 `transition-all`。`transition-all`
会把 motion 正在逐帧驱动的 `transform`/`opacity` 也卷入 CSS 过渡，叠出二次位移（motion 播一遍，
CSS 过渡再过渡一遍同一批属性变化）。

规避手法（利用 Tailwind v4 特性）：

- Tailwind v4 的 `translate-y-*` / `scale-*` / `rotate-*` 生成的是**原生 `translate`/`scale`/`rotate` 属性**，
  而 motion 的 `y`/`scale`/`rotate` 生成的是 **`transform: translateY()/scale()/rotate()`**——两者分属不同 CSS 属性。
- 因此卡片需要 hover 浮起时，把 CSS 过渡**收窄**到 hover 实际用到的属性即可，motion 的
  `transform`/`opacity` 入场自然不被卷入：

  ```ts
  // inkCard（ink-wash.tsx）—— 过渡只含 translate / border-color / box-shadow
  '... transition-[translate,border-color,box-shadow] duration-300 hover:-translate-y-0.5 ...'
  ```

- 判别口诀：**motion 动什么，CSS transition 就不过渡什么**。

#### 铁律 3 · 区块级入场与 hover 过渡分属不同元素更佳

当 KPI 三块的 `motion.section` 本身不带 `inkCard`（卡片是内层 `div.inkCard`），motion 动画与
`transition-all` 在不同元素，天然互不干扰。新页面布局优先采用这种「外层 motion / 内层 hover 过渡」
的分层结构；确需同一元素兼顾两者时，回到铁律 2 的收窄手法。

#### 铁律 4 · reduced-motion 全降级

- CSS keyframes：`@media (prefers-reduced-motion: reduce)` 下 `animation: none`。
- motion 入场：经 `enter()` 退化为纯 opacity 快速淡入（无位移、delay 压缩）。
- CountUp：直接跳终值。

### 2.7 标杆与参考

- **标杆页面**：`src/routes/_admin/admin/dashboard.tsx`（水墨叙事 + 双模式 hero + 错峰入场）。
- **入场助手**：`src/lib/motion.ts`（`enter()`）。
- **路由级过渡**：`src/routes/_admin/route.tsx`。
- **共享视觉原语**：`src/components/ink-wash.tsx`（`inkCard` / `CardHead` / `BambooRule` / `BambooArt` 等）。
- **数据可视化**：`src/components/dashboard/`（count-up / donut-chart / radial-gauge）。
- **token 与 CSS 动画**：`src/styles.css`。
