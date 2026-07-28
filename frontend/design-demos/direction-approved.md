# 方向确认 · Sidebar 竹林水墨美化

> huashu-design Gate 文件。本任务属于 Fallback「唯一豁免」中的 **已选定方向后的迭代**，
> 不重新过三方向门，记录豁免依据与用户原话。

## 豁免依据

- **方向来源**：竹林水墨设计语言已由用户在上一次 dashboard 优化会话拍板并落地
  （提交 `d6c0306`，标杆 `frontend/src/routes/_admin/admin/dashboard.tsx`）。
- **用户本次原话**：「现在请你开始对 Sidebar 也这么进行处理，现在 Sidebar 保持 out 的风格，
  然后对内部内容进行格式美学调整」——即在既定水墨语言内迭代，外部结构（floating /
  collapsible / 三段式）保持不变，仅美化内部内容。
- **约束来源**：`bamboo-ink-design` skill 的「竹林水墨 spec」作为硬约束（brand-spec）注入。

## 设计决策（在既定语言内）

- 标杆对齐：`dashboard.tsx`（斜墨条 / 衬线 / mono meta / 竹节线 / 墨韵竹叶）+ `styles.css` token。
- Header：墨晕 logo 底 + 衬线站名 + mono version + 竹节线。
- 分组标签：去 icon（反「每标题配 icon」slop），改斜墨条 + 衬线签名。
- 菜单项：激活态左侧斜墨条 + leaf-deep 文字图标 + 淡绿底；收起态 tooltip。
- Footer：竹节线 + leaf-deep 头像 + 衬线名 / mono email + 极淡墨韵竹叶点缀。
- 反 slop：无紫渐变、无 emoji、无玻璃拟态、不渲染印章；token 全部复用 styles.css，不自创颜色。

## 验证

- `cd frontend && npx tsc --noEmit` 零错误。
- 浏览器展开 / 收起两态零 console error。

---

# 方向确认 · About 三页竹林水墨重塑

> 本任务同样属于 Fallback「唯一豁免」中的 **已选定方向后的迭代**，
> 不重新过三方向门。记录豁免依据、用户原话与设计稿交付路径。

## 豁免依据

- **方向来源**：竹林水墨设计语言已由用户在 dashboard 优化会话拍板并落地
  （提交 `d6c0306`，标杆 `frontend/src/routes/_admin/admin/dashboard.tsx`）。
- **用户本次原话**：「现在，你有一个新的 UI 改进方案，参考一下 AGENT.md 以及
  huashu-design 的内容，主要改进一下主页扩展的 about 的三个界面，现在这三个界面
  还非常的朴素。」——即在既定水墨语言内重塑 `/about` 下 me / friends / sponsor 三页。
- **工作流选择原话**：「使用 HTML 涉及稿然后迁移，我自己打开看就行，不需要截图」
  ——先出 HTML 设计稿给用户自行打开查看，定稿后再迁移为 React 代码。

## 设计决策（在既定语言内）

- 标杆对齐：`dashboard.tsx`（斜墨条 CardHead / inkCard 宣纸块 / 竹节线 / 墨韵竹叶
  / 笔刷下划线 / enso 缺口圆空状态 / 衬线水印字）+ `styles.css` token。
- **Phase 0 共享原语**：新建 `src/lib/motion.ts`（抽取 `enter()` + `MotionDivProps`，
  原先在 6 个文件重复）、`src/components/ink-wash.tsx`（抽取 inkCard / CardHead /
  BambooRule / BambooArt / 新增 BrushUnderline + EnsoEmpty），dashboard 机械重构导入。
- **me「名士名帖」**：左名帖卡（头像 + 衬线名 + mono 站名 + 笔刷下划线 + 晨光墨晕 +
  墨韵竹叶）+ 右自介卷（CardHead + 衬线斜体导语 + 竹节线 + Markdown + 水印字「介」）。
- **friends「竹林群贤」**：分组 CardHead（斜墨条 + 衬线 + mono 友人数）+ 宣纸友链卡网格
  （hover 顶部墨线扫入 + 左侧友人主色细条 + 衬线名 hover→leaf-deep）+ 竹节线分组 +
  enso 空状态。
- **sponsor「感恩录」**：渠道 CardHead + 渠道宣纸卡网格 + 竹节线 + 记录大宣纸卡
  （mono 序号 + 衬线昵称 + mono 渠道·日期 + 留言斜体 + 衬线 leaf-deep 金额 + 水印字「谢」）+
  enso 空状态。
- **壳润色**：标题「关于小站」改衬线 + 居中笔刷下划线，副标题衬线斜体，`enter()` 改用共享模块。
- 反 slop：无紫渐变、无 emoji、无玻璃拟态、不渲染印章；token 全部复用 styles.css，
  不自创颜色；占位数据明确标注。

## 设计稿交付

三份自包含 HTML（双击即可在浏览器打开，含内联 token / keyframes / BambooArt SVG /
Google Fonts 衬线，占位数据明确标注）：

- `design-demos/about/me.html`
- `design-demos/about/friends.html`
- `design-demos/about/sponsor.html`

## 待办

- [x] Phase 0：抽取共享原语 + dashboard 机械重构（tsc 零错误）
- [x] Phase 1：三份 HTML 设计稿
- [ ] Phase 2：用户打开设计稿确认（等待）
- [ ] Phase 3：迁移 React（me/friends/sponsor.tsx + route.tsx 壳润色）
- [ ] Phase 4：tsc 零错误 + 浏览器验证

---

## v3 迭代 · 高级感重塑（竹林三帖）

- **用户原话**：「重新设计，我觉得还是欠缺了不够高级，基于当前的整体样式风格整体重新设计 UI/UX」
  ——属同一竹林水墨方向内的质感升级迭代，不重新过三方向门。
- **级别数据**：用户确认「当前只确认前端样式，后续再补充」——FriendCard 分级组件先做前端，
  level 数据来源（后端字段）后续再定。
- **v3 诊断**：v2 仍停留在「卡片网格」范式——均匀圆角卡、安全对称布局、温和字号、
  水墨仅作点缀、无叙事节奏，故显「不够高级」。
- **v3 设计语言（竹林三帖 · 手卷式编辑布局）**：
  - 巨型衬线（clamp 至 7-8.5rem）× 微型 mono（11px / 0.35em 字距）极端对比；
  - 不对称构图：全出血竹叶、12 栅格 5/7 与 3/9 分栏、元素交叠、大量留白；
  - 水墨作结构：大竹叶章节背景、墨晕分章带、巨型水印字（介/站/谢/录）、大写意笔刷；
  - 展卷叙事：IntersectionObserver 逐段 reveal + sticky 着墨导航 + 滚动提示线；
  - 分级友链：壹 高级（全幅特写+浏览器框截图）/ 贰 好友（3:9 分栏+竖分隔）/
    叁 一般（编辑式名录，hover 墨线自左扫入）/ 肆 广告（居中+「推广」标）；
  - 赞助账册：mono 序号 + 衬线昵称 + 斜体留言 + 衬线大字金额（leaf-deep）。
- **反 slop**：无紫渐变、无 emoji、无玻璃拟态（导航用实底非 backdrop-blur）、不渲染印章、
  token 全复用 styles.css。
- **交付**：`design-demos/about/{me,friends,sponsor}.html`（v3，双击即开）。

---

## 落地实施 · 全栈对接（v4.1 定稿迁移）

- **用户原话**：「允许进入计划模式，开始实际做计划然后对接进入实际的项目中，现在内容
  以及足够成熟了」+「要求后端字段多一个叫做 level 方案记录在案，数据库采用枚举数字
  存储，映射由 go 映射信息」+「friend-card 四个模式需要各自建一个 components 写入
  组件库中」。
- **级别数据源（最终拍板）**：新增后端 `level` 字段，DB 枚举数字存储 + Go 常量映射，
  前台按 `link.level` 选用四级卡片组件。

### 后端 level 字段（DB 枚举数字 + Go 映射）

| 值 | Go 常量 | 含义 | 前台卡片 |
|---|---|---|---|
| `0` | `LinkLevelRegular` | 一般 | `RegularFriendCard`（1×1 紧凑） |
| `1` | `LinkLevelClose` | 好友 | `CloseFriendCard`（1×1 富式） |
| `2` | `LinkLevelPremium` | 高级 | `PremiumFriendCard`（2×2 特写） |
| `3` | `LinkLevelAd` | 广告 | `AdFriendCard`（1×1 居中 + 推广标） |

- **常量**：`pkg/constants/context.go`（紧随 `LinkFailBroken`，仿既有分组惯例，未新建文件）。
- **实体**：`internal/entity/link_context.go` → `LinkFriend.Level int`
  （`gorm:"type:int;default:0;comment:友链级别（0: 一般, 1: 好友, 2: 高级, 3: 广告）"`，
  仿 Status/IsFailure 惯例；响应 DTO 内嵌 entity 自动透出 `level`）。
- **DTO**：`api/link/link.go` → `FriendAddRequest.LinkLevel int` /
  `FriendUpdateRequest.LinkLevel *int`（均 `binding:"omitempty,oneof=0 1 2 3"`）。
- **logic**：`internal/logic/link.go` → Add 直取 `req.LinkLevel`；Update 走
  `if req.LinkLevel != nil { link.Level = *req.LinkLevel }` 指针判空。
- **DB 迁移**：GORM AutoMigrate 自动加列（default 0，存量友链全为「一般」）。
- **swagger**：`make swag` 已重新生成。

### 前端组件库（`src/components/about/`，named export）

四级友链卡**各自独立成组件**（用户要求，非单组件 + level prop）：
`premium-friend-card.tsx` / `close-friend-card.tsx` / `regular-friend-card.tsx` /
`ad-friend-card.tsx`；共享 `friend-card-shared.ts`（`FriendCardProps` / `domainOf()` /
`useFriendOpen()`）；`interlude.tsx` 提供 clip-path 沉浸式跳转过渡。
`friends.tsx` 按 `link.level` switch 选用组件，Bento 栅格 `grid-flow-dense` 自动填充。

### 管理端级别可配置

`link/add.tsx` + `$id.edit.tsx` 增加「友链级别」药丸按钮组（一般/好友/高级/广告），
提交 `link_level`；编辑页从 `link.level` 预填。

### 验证

- 后端：`go build ./... && go vet ./...` 零错误；`make swag` 成功。
- 前端：`cd frontend && npx tsc --noEmit` 零错误。
- 浏览器：`pnpm dev` 逐页核对（见下）。
