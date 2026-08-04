/*
 * --------------------------------------------------------------------------------
 * Copyright (c) 2016-NOW(至今) 筱锋
 * Author: 筱锋「xiao_lfeng」(https://www.x-lf.com)
 * --------------------------------------------------------------------------------
 * 许可证声明：版权所有 (c) 2016-2026 筱锋。保留所有权利。
 * 有关MIT许可证的更多信息，请查看项目根目录下的 LICENSE 文件或访问：
 * https://opensource.org/licenses/MIT
 * --------------------------------------------------------------------------------
 */

import { createFileRoute } from '@tanstack/react-router'
import { motion, useReducedMotion } from 'motion/react'
import { RankingBoard } from '@/components/link/ranking-board'
import { BambooRule, PageHead } from '@/components/ink-wash'
import { enter } from '@/lib/motion'

export const Route = createFileRoute('/_admin/admin/ranking')({
  component: RankingPage,
})

function RankingPage() {
  const reduced = useReducedMotion() ?? false

  return (
    <div className="mx-auto max-w-6xl space-y-6">
      <PageHead
        kicker="RANKING · 排位"
        title="排位管理"
        backTo="/admin/link"
        backLabel="友链管理"
        sub="拖拽卡片调整友链排位，版式与访客页逐格同构——高级友链以 2×2 大卡呈现，其余为 1×1 小卡。组内重排、跨组移动、章节拖拽松手后自动保存，所见即所得。"
      />

      <BambooRule reduced={reduced} delay={0.12} />

      <motion.section
        {...enter(reduced, 0.18, {
          initial: { opacity: 0, y: 10 },
          animate: { opacity: 1, y: 0 },
          transition: { duration: 0.4, ease: 'easeOut' },
        })}
      >
        <RankingBoard />
      </motion.section>
    </div>
  )
}
