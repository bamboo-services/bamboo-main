// --------------------------------------------------------------------------------
// Copyright (c) 2016-NOW(至今) 筱锋
// Author: 筱锋「xiao_lfeng」(https://www.x-lf.com)
// --------------------------------------------------------------------------------
// 许可证声明：版权所有 (c) 2016-2026 筱锋。保留所有权利。
// 有关MIT许可证的更多信息，请查看项目根目录下的LICENSE文件或访问：
// https://opensource.org/licenses/MIT
// --------------------------------------------------------------------------------

import { createElement } from 'react'
import { renderToStaticMarkup } from 'react-dom/server'
import { MarkdownView } from '@/components/markdown'

/**
 * 把 Markdown 渲染为 HTML 内容字符串（剥掉外层 `.markdown-view` 壳）。
 *
 * 复用 `MarkdownView` 作为唯一渲染真相源，保证 Markdown 编辑器的预览
 * 与公开页排版像素级一致。`react-dom/server` 的 `renderToStaticMarkup`
 * 为纯静态渲染，无副作用，可安全用于 contenteditable 的 innerHTML。
 */
export function renderMarkdownToHtml(value: string): string {
  const wrapper = document.createElement('div')
  wrapper.innerHTML = renderToStaticMarkup(
    createElement(MarkdownView, { content: value }),
  )
  return wrapper.firstElementChild?.innerHTML ?? ''
}
