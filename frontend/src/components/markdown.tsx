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
 * Markdown 渲染器 —— 套用项目主题变量，不依赖 @tailwindcss/typography。
 * - 标题/正文：text-text-primary
 * - 链接：text-primary + 下划线 hover
 * - 代码/引用/分隔线：用 leaf-* 变量做淡绿点缀
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
            <h1 className="mb-4 mt-6 text-2xl font-bold text-text-primary first:mt-0">
              {children}
            </h1>
          ),
          h2: ({ children }) => (
            <h2 className="mb-3 mt-5 text-xl font-semibold text-text-primary">
              {children}
            </h2>
          ),
          h3: ({ children }) => (
            <h3 className="mb-2 mt-4 text-lg font-semibold text-text-primary">
              {children}
            </h3>
          ),
          p: ({ children }) => (
            <p className="mb-3 leading-7 text-text-secondary last:mb-0">
              {children}
            </p>
          ),
          ul: ({ children }) => (
            <ul className="mb-3 ml-5 list-disc space-y-1 text-text-secondary last:mb-0">
              {children}
            </ul>
          ),
          ol: ({ children }) => (
            <ol className="mb-3 ml-5 list-decimal space-y-1 text-text-secondary last:mb-0">
              {children}
            </ol>
          ),
          li: ({ children }) => <li className="leading-7">{children}</li>,
          a: ({ children, href }) => (
            <a
              href={href}
              target="_blank"
              rel="noopener noreferrer"
              className="font-medium text-primary underline-offset-2 hover:underline"
            >
              {children}
            </a>
          ),
          blockquote: ({ children }) => (
            <blockquote className="mb-3 border-l-2 border-leaf-muted bg-leaf-light/20 px-4 py-2 text-text-secondary italic last:mb-0">
              {children}
            </blockquote>
          ),
          code: ({ className: cls, children }) => {
            const isBlock = typeof cls === 'string' && cls.startsWith('language-')
            if (isBlock) {
              return (
                <code className="block overflow-x-auto rounded-md bg-text-primary/5 p-3 text-sm text-text-primary">
                  {children}
                </code>
              )
            }
            return (
              <code className="rounded bg-leaf-light/30 px-1.5 py-0.5 text-sm text-text-primary">
                {children}
              </code>
            )
          },
          pre: ({ children }) => (
            <pre className="mb-3 last:mb-0">{children}</pre>
          ),
          hr: () => (
            <hr className="my-4 border-none border-t border-leaf-muted/50" />
          ),
          table: ({ children }) => (
            <div className="mb-3 overflow-x-auto last:mb-0">
              <table className="w-full border-collapse text-sm text-text-secondary">
                {children}
              </table>
            </div>
          ),
          th: ({ children }) => (
            <th className="border border-leaf-muted/40 px-3 py-1.5 text-left font-semibold text-text-primary">
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
