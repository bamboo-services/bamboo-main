# 设计方向确认 · 赞助页 CTA 按钮重设计

**日期:** 2026-08-05
**分支:** worktree-sponsor-cta-redesign

## 背景

本次为「赞助前端重设计」任务：为赞助页 `/about/sponsor` 的 CTA 按钮
（原为普通 outline Button「已赞助 · 申请展示」）设计手工定制按钮，
对齐友链页 `NameCardButton`（名帖母题）的定制水准。

视觉语言已定（竹林水墨 spec，`bamboo-ink-design` skill 硬约束），
故三方向为**同一语言内的母题差异化诠释**，非风格发散。

## 三方向初稿

三版 HTML 初稿（含默认态 + hover 态预览）存放于：
- `~/.claude/jobs/7006afb4/tmp/design-demos/方向A-随喜签.html`
- `~/.claude/jobs/7006afb4/tmp/design-demos/方向B-题名入册.html`
- `~/.claude/jobs/7006afb4/tmp/design-demos/方向C-谢帖回执.html`

| 方向 | 母题 | 结构签名 |
|---|---|---|
| A · 随喜签 | 随喜结缘（铜钱） | 方孔圆钱右上 + 竖排「随喜」签条 + 右下折角 |
| B · 题名入册 | 感恩账册 | 深绿实底签条 + 笔刷下划线 + 谢字水印 |
| C · 谢帖回执 | 与拜帖对偶 | 折角卡同构 + 竖排「谢帖」签条 |

## 用户选择

> 赞助页「已赞助·申请展示」CTA 按钮选哪个母题方向？
> **「A · 随喜签（铜钱）」** —— 用户原话选择项：
> 「方孔圆钱 + 竖排『随喜』签条，母题最独特、识别度最高，紧扣『随喜之门』章节词汇。」

用户补充要求（原话）：
> 「但是摆放位置要求跟 friend 一样，在右侧才对」

即：按钮摆放位置与友链页一致，放在 hero **右侧**（原 CTA 在左栏 intro 下方，
改为 friends 同款双栏布局——左栏开场文案、右栏「随喜结缘」竖排题跋 + 随喜签按钮）。

## 落地

- 新增 `SuixiSignButton` 组件（`frontend/src/routes/about/sponsor.tsx`）
- hero 改双栏 grid（`lg:grid-cols-12`），按钮移至右栏，配竖排「随喜结缘」题跋
- 移除原 outline Button 与 `Button` import
