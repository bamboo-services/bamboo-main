// --------------------------------------------------------------------------------
// Copyright (c) 2016-NOW(至今) 筱锋
// Author: 筱锋「xiao_lfeng」(https://www.x-lf.com)
// --------------------------------------------------------------------------------
// 许可证声明：版权所有 (c) 2016-2026 筱锋。保留所有权利。
// 有关MIT许可证的更多信息，请查看项目根目录下的LICENSE文件或访问：
// https://opensource.org/licenses/MIT
// --------------------------------------------------------------------------------

import { memo } from 'react'
import Markdown from 'react-markdown'
import remarkGfm from 'remark-gfm'

/**
 * Markdown 渲染器 —— 竹林水墨高级排版，不依赖 @tailwindcss/typography。
 * - 标题：衬线 display；正文：17px / 1.95 行高，限宽由父级 max-w-prose 控制
 * - 列表：斜墨条项目符（ul，CSS ::before）；有序列表保留衬线数字
 * - 引用块：渲染为「拉风引语」——衬线斜体大字 + 左墨条 + 悬浮「（CSS）
 * - 首段首字下沉、链接 leaf-deep 下划线，均由 styles.css 的 .markdown-view 承载
 */
export const MarkdownView = memo(function MarkdownView({
  content,
  className = '',
}: {
  content: string
  className?: string
}) {
  return (
    <div className={`markdown-view ${className}`}>
      <Markdown
        remarkPlugins={[remarkGfm]}
        components={{
          h1: ({ children }) => (
            <h1 className="mb-4 mt-9 font-serif text-3xl font-bold text-text-primary first:mt-0">
              {children}
            </h1>
          ),
          h2: ({ children }) => (
            <h2 className="mb-3 mt-8 font-serif text-2xl font-semibold text-text-primary">
              {children}
            </h2>
          ),
          h3: ({ children }) => (
            <h3 className="mb-3 mt-7 font-serif text-[1.3rem] font-semibold text-text-primary">
              {children}
            </h3>
          ),
          p: ({ children }) => (
            <p className="mb-4 text-[17px] leading-[1.95] text-text-secondary last:mb-0">
              {children}
            </p>
          ),
          ul: ({ children }) => (
            <ul className="my-5 list-none space-y-2.5 last:mb-0">{children}</ul>
          ),
          ol: ({ children }) => (
            <ol className="my-5 ml-6 list-decimal space-y-2.5 last:mb-0">
              {children}
            </ol>
          ),
          li: ({ children }) => (
            <li className="text-[16px] leading-[1.75] text-text-secondary">
              {children}
            </li>
          ),
          a: ({ children, href }) => (
            <a
              href={href}
              target="_blank"
              rel="noopener noreferrer"
              className="font-medium text-leaf-deep underline decoration-leaf-muted underline-offset-[3px] transition-colors hover:decoration-leaf-deep"
            >
              {children}
            </a>
          ),
          blockquote: ({ children }) => (
            <blockquote className="relative my-11 ml-auto max-w-[34rem] border-l-2 border-leaf-deep py-1 pl-7">
              {children}
            </blockquote>
          ),
          code: ({ className: cls, children }) => {
            const isBlock =
              typeof cls === 'string' && cls.startsWith('language-')
            if (isBlock) {
              return (
                <code className="block overflow-x-auto rounded-md bg-text-primary/5 p-3 font-mono text-sm text-text-primary">
                  {children}
                </code>
              )
            }
            return (
              <code className="rounded bg-leaf-light/30 px-1.5 py-0.5 font-mono text-sm text-text-primary">
                {children}
              </code>
            )
          },
          pre: ({ children }) => (
            <pre className="mb-4 last:mb-0">{children}</pre>
          ),
          hr: () => (
            <hr className="my-9 border-none border-t border-leaf-muted/50" />
          ),
          table: ({ children }) => (
            <div className="mb-4 overflow-x-auto last:mb-0">
              <table className="w-full border-collapse text-sm text-text-secondary">
                {children}
              </table>
            </div>
          ),
          th: ({ children }) => (
            <th className="border border-leaf-muted/40 px-3 py-1.5 text-left font-serif font-semibold text-text-primary">
              {children}
            </th>
          ),
          td: ({ children }) => (
            <td className="border border-leaf-muted/40 px-3 py-1.5">
              {children}
            </td>
          ),
        }}
      >
        {content}
      </Markdown>
    </div>
  )
})
