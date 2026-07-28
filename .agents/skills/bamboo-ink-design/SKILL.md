---
name: bamboo-ink-design
description: >-
  bamboo-main 前端「竹林水墨」设计语言约束器。任何涉及本项目 frontend/ 的视觉设计、
  UI 改版、新页面/组件美化、配色选型、设计探索的任务都必须使用此技能——即使用户只说
  "做个好看的页面"、"调整一下美学"、"这个界面太丑"、"设计一版 XX"、"出几个方案"。
  本技能不自己做设计，而是引导并约束 huashu-design skill：调用 huashu-design 时，
  必须把本技能提供的「竹林水墨」spec 作为硬约束（brand-spec）喂给它，确保产出沿用
  项目既定的淡绿宣纸 + 墨色衬线 + 东方签名语言，而不是泛化的通用 SaaS 审美。
---

# Bamboo Ink Design · 竹林水墨设计约束

本项目（bamboo-main 友情链接管理后台）的前端有**既定且已落地的视觉语言**：「竹林水墨 · 东方品牌签名」，
标杆实现是 `frontend/src/routes/_admin/admin/dashboard.tsx`。

**本技能的唯一职责**：当你需要做前端视觉设计时，**调用 `huashu-design` skill 去执行设计**，
但必须把下面的「竹林水墨 spec」作为**硬约束（相当于 huashu-design 的 brand-spec / 设计 spec）**
传给它，让它在既定语言内发挥，而不是凭空发明一套新风格。

> 为什么需要这层约束：huashu-design 默认会探索各种风格方向（报刊/数据密集/叙事/瑞士粗野……）。
> 但本项目的设计方向**已经由用户拍板**（水墨叙事 + 双模式 hero），不需要再发散探索风格，
> 只需要在既定语言内把具体页面做扎实。本技能就是把"已拍板的方向"固化成约束。

---

## 执行流程

1. **识别任务**：用户要前端视觉设计（新页面 / 改版 / 美化 / 出方案）。
2. **调用 huashu-design**：通过 Skill 工具加载 `huashu-design`。
3. **注入约束**：在 huashu-design 的「设计 spec / brand-spec」环节，把本技能下方
   「竹林水墨 spec」整段作为**不可违背的既定约束**喂给它。明确告知 huashu-design：
   - 风格方向**已定**（水墨叙事），**跳过三方向风格探索门**（属于 huashu-design
     「唯一豁免」中的「已选定方向后的迭代」），直接在该语言内执行。
   - 色彩 / 字体 / 布局 / 品牌签名 / 反 slop 全部按本 spec，不许自创。
4. **标杆参照**：让 huashu-design 先读 `frontend/src/routes/_admin/admin/dashboard.tsx`
   与 `frontend/src/styles.css`，以现有实现为对齐基准。
5. **产出落地**：设计稿确认后，迁移为 React + Tailwind v4 代码，复用 `src/components/dashboard/`
   下的共享组件（count-up / donut-chart / radial-gauge），token 一律走 `styles.css` 的 CSS 变量。

---

## 竹林水墨 spec（喂给 huashu-design 的硬约束）

### 气质定调

竹林清晨的雅致、克制、有书卷气。淡绿宣纸底 + 墨色文字 + 衬线 display + 单一 leaf-deep accent，
配墨韵竹叶 / 晨光墨晕 / 竹节线 / 竖排题跋等东方签名。**像一份精心排印的晨报 / 书画题跋，
不是通用 SaaS dashboard。**

### 色彩 token（`frontend/src/styles.css`，禁止自创颜色）

```css
--background: oklch(0.975 0.016 110);  /* 宣纸页面底 */
--card:       oklch(0.99 0.005 110);   /* 卡片块底 */
--text-primary: oklch(0.32 0.06 155);  /* 墨色，主文字/浓墨 */
--text-secondary: oklch(0.5 0.05 130); /* 淡墨，次级文字 */
--leaf-light: oklch(0.88 0.1 105);     /* 浅叶绿（晨光墨晕） */
--leaf-muted: oklch(0.8 0.08 130);     /* 柔叶绿（竹节线） */
--leaf-deep:  oklch(0.55 0.12 155);    /* 深叶绿，唯一 accent */
--seal:       oklch(0.55 0.18 25);     /* 朱砂红，备用，当前未启用 */
--chart-1/4/5: 见 styles.css           /* donut 分段 */
```

### 字体

- **display（标题 / 大数字）**：衬线 `font-serif`（`Source Serif 4` + `Noto Serif SC`，
  经 `index.html` Google Fonts 引入，系统衬线 fallback）。
- **正文 / label**：系统 sans。
- **数据 / meta**：`font-mono` + `tabular-nums`（数字严格对齐）。

### 布局与组件模式

- **宣纸卡片块**：`rounded-lg border border-border bg-card p-5` + hover 浮起 + 墨边软阴影。
- **卡片标题**：斜墨条（`-skew-x-12` leaf-deep）+ 衬线标题 + 右侧 mono uppercase meta。
- **竹节分隔线**：横线 + 圆/宽节点交替，分大章节。
- **双模式 hero**：sidebar 展开 → 顶部横排大问候；sidebar 收起（`useSidebar().open === false`）→
  左侧竖排卷轴（`writing-mode: vertical-rl`）+ 右侧数据区，grid `grid-cols-[280px_1fr]`。

### 品牌签名元素

- **墨韵竹叶**：内联 SVG 三竿竹，墨色 opacity 0.045–0.1，`ink-sway` 极慢摇曳
  （`transform-box: fill-box` 底部锚点）。
- **晨光墨晕**：单色淡绿径向渐变（仅 leaf-light），非多色光斑。
- **竖排题跋**：`[writing-mode:vertical-rl] [text-orientation:upright]` 大号衬线问候 + 导语。
- **笔刷下划线**：手绘 SVG path（leaf-deep）。
- **enso 缺口圆**：空状态图标（圆环留缺口 + 竹叶）。
- **衬线水印字**：卡片背景大字（opacity 5%）。

### 动画（`styles.css`，reduced-motion 全部静止）

- `ink-sway` / `ink-sway-slow` / `ink-sway-fast`：竹叶摇曳。
- `ink-pulse`：状态点呼吸。
- 入场用 `motion/react` 淡入 + 轻微位移，克制不过度。

### 反 AI slop（硬红线）

- 禁紫色 / 靛蓝 / 粉色渐变、渐变填充标题。
- 禁 emoji 当图标（用内联 SVG 线条图标或纯文字）。
- 禁玻璃拟态（backdrop-blur）、aurora 多色光斑背景。
- 禁近黑背景 + 霓虹 accent（本项目是宣纸浅色主题）。
- 禁每个标题配圆角 icon 块。
- 禁编造假数据 / 装饰性 stats。
- **印章元素已弃用**（用户决策），勿再加朱砂印章。

---

## 标杆与参考

- **标杆页面**：`frontend/src/routes/_admin/admin/dashboard.tsx`（水墨叙事 + 双模式 hero）。
- **共享组件**：`frontend/src/components/dashboard/`（count-up / donut-chart / radial-gauge）。
- **token 与动画**：`frontend/src/styles.css`。
- **完整设计文档**：`frontend/AGENTS.md` 的「设计语言 · 竹林水墨」章节。

---

## 一句话原则

> 设计方向已定，别再发散。调用 huashu-design，把「竹林水墨 spec」当 brand-spec 喂进去，
> 在淡绿宣纸 + 墨色衬线 + 东方签名的既定语言内，把页面做扎实。
