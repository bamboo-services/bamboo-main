// --------------------------------------------------------------------------------
// Copyright (c) 2016-NOW(至今) 筱锋
// Author: 筱锋「xiao_lfeng」(https://www.x-lf.com)
// --------------------------------------------------------------------------------
// 许可证声明：版权所有 (c) 2016-2026 筱锋。保留所有权利。
// 有关MIT许可证的更多信息，请查看项目根目录下的LICENSE文件或访问：
// https://opensource.org/licenses/MIT
// --------------------------------------------------------------------------------

/**
 * contenteditable DOM → Markdown 序列化纯函数。
 *
 * 供 Markdown 编辑器（所见即所得模式）把编辑后的 DOM 转回 Markdown 字符串，
 * 与 `renderMarkdownToHtml`（Markdown → HTML）互为逆操作，覆盖
 * MarkdownView 渲染的全部元素：段落/标题/粗斜删/行内代码/链接/图片/
 * 有序无序列表（含嵌套）/引用/横线/代码块/软硬换行/GFM 表格。
 *
 * 仅依赖标准 DOM API，无 React 依赖，可在 jsdom 环境下单测。
 */

// ---------- 基础工具 ----------

/** 折叠空白：连续空白（含换行、NBSP）收敛为单个空格 */
function collapseWhitespace(text: string): string {
  return text.replace(/\s+/g, ' ')
}

/** 段落行首转义：避免「# / > / - / 数字.」开头被渲染成标题/列表/引用 */
function escapeBlockStart(text: string): string {
  return text.replace(/^(#{1,6}\s|>\s|[-+*]\s|\d+\.\s)/, (m) => `\\${m}`)
}

/** 行内代码围栏：内容含反引号时自动升级为双反引号围栏 */
function codeFence(text: string): string {
  if (!text.includes('`')) return `\`${text}\``
  if (!text.includes('``')) return `\`\`${text}\`\``
  return text
}

// ---------- 行内序列化 ----------

/** 递归序列化元素的所有子节点到 out */
function serializeInlineChildren(el: Element, out: string[]): void {
  for (const child of Array.from(el.childNodes)) {
    serializeInline(child, out)
  }
}

/** 递归序列化元素的子节点并拼接为字符串 */
function inlineRaw(el: Element): string {
  const parts: string[] = []
  serializeInlineChildren(el, parts)
  return parts.join('')
}

/** 序列化单个行内节点 */
function serializeInline(node: Node, out: string[]): void {
  if (node.nodeType === Node.TEXT_NODE) {
    const text = collapseWhitespace(node.textContent ?? '')
    if (text) out.push(text)
    return
  }
  if (node.nodeType === Node.COMMENT_NODE) return
  if (node.nodeType !== Node.ELEMENT_NODE) return

  const el = node as HTMLElement
  switch (el.tagName.toLowerCase()) {
    case 'br':
      // 行尾两空格 + 换行 = Markdown 硬换行（remark-gfm 渲染为 <br>）
      out.push('  \n')
      return
    case 'strong':
    case 'b':
      out.push(`**${inlineRaw(el)}**`)
      return
    case 'em':
    case 'i':
      out.push(`*${inlineRaw(el)}*`)
      return
    case 'del':
    case 's':
    case 'strike':
      out.push(`~~${inlineRaw(el)}~~`)
      return
    case 'code':
      out.push(codeFence(el.textContent ?? ''))
      return
    case 'a': {
      const href = el.getAttribute('href') ?? ''
      const text = inlineRaw(el).replace(/]/g, '\\]')
      out.push(`[${text}](${href})`)
      return
    }
    case 'img': {
      const src = el.getAttribute('src') ?? ''
      const alt = el.getAttribute('alt') ?? ''
      out.push(`![${alt}](${src})`)
      return
    }
    case 'span': {
      // styled-span 语义探测：加粗/斜体/删除线的内联样式保留语义
      const style = el.style
      const isBold = /\bbold\b|600|700|800|900/.test(style.fontWeight)
      const isItalic = style.fontStyle === 'italic'
      const isStrike = `${style.textDecorationLine} ${style.textDecoration}`.includes(
        'line-through',
      )
      if (isBold || isItalic || isStrike) {
        let wrapped = inlineRaw(el)
        if (isStrike) wrapped = `~~${wrapped}~~`
        if (isItalic) wrapped = `*${wrapped}*`
        if (isBold) wrapped = `**${wrapped}**`
        out.push(wrapped)
      } else {
        // 纯样式壳（FONT/U/MARK 等）：解包透传子节点
        serializeInlineChildren(el, out)
      }
      return
    }
    default:
      serializeInlineChildren(el, out)
  }
}

// ---------- 块级序列化 ----------

/** 序列化单个块节点，返回其 Markdown（多行块内部用换行，块间由入口拼接） */
function serializeBlock(node: Node): string {
  if (node.nodeType === Node.TEXT_NODE) {
    const text = collapseWhitespace(node.textContent ?? '')
    return text ? escapeBlockStart(text.trim()) : ''
  }
  if (node.nodeType === Node.COMMENT_NODE) return ''
  if (node.nodeType !== Node.ELEMENT_NODE) return ''

  const el = node as HTMLElement
  switch (el.tagName.toLowerCase()) {
    case 'p':
    case 'div': {
      const text = inlineRaw(el).trim()
      return text ? escapeBlockStart(text) : ''
    }
    case 'h1':
    case 'h2':
    case 'h3': {
      const level = Number(el.tagName[1])
      const text = inlineRaw(el).trim()
      return text ? `${'#'.repeat(level)} ${text}` : ''
    }
    case 'ul':
    case 'ol':
      return serializeList(el, 0)
    case 'blockquote': {
      const inner: string[] = []
      for (const child of Array.from(el.childNodes)) {
        const block = serializeBlock(child)
        if (block) inner.push(block)
      }
      if (inner.length === 0) return ''
      // 每行加「> 」前缀；嵌套引用会因 inner 已含「> 」而自然累积为「>> 」
      return inner
        .join('\n\n')
        .split('\n')
        .map((line) => `> ${line}`)
        .join('\n')
    }
    case 'pre': {
      const codeEl = el.querySelector('code')
      const langMatch = codeEl?.className?.match(/language-([\w-]+)/)
      const lang = langMatch ? langMatch[1] : ''
      const content = (el.textContent ?? '').replace(/^\n+|\n+$/g, '')
      const fence = content.includes('```') ? '````' : '```'
      return `${fence}${lang}\n${content}\n${fence}`
    }
    case 'hr':
      return '---'
    case 'table':
      return serializeTable(el)
    case 'img': {
      const src = el.getAttribute('src') ?? ''
      const alt = el.getAttribute('alt') ?? ''
      return `![${alt}](${src})`
    }
    default: {
      // 容器（section/li 兜底）：拼接子块
      const inner: string[] = []
      for (const child of Array.from(el.childNodes)) {
        const block = serializeBlock(child)
        if (block) inner.push(block)
      }
      return inner.join('\n\n')
    }
  }
}

/** 有序/无序列表（含嵌套）序列化为多行 Markdown */
function serializeList(list: HTMLElement, depth: number): string {
  const ordered = list.tagName.toLowerCase() === 'ol'
  const indent = '  '.repeat(depth)
  const lines: string[] = []
  let index = 1

  for (const li of Array.from(list.children)) {
    if (li.tagName.toLowerCase() !== 'li') continue

    // 行内内容与嵌套列表分离
    const inlineParts: string[] = []
    const nested: HTMLElement[] = []
    for (const child of Array.from(li.childNodes)) {
      if (child.nodeType === Node.ELEMENT_NODE) {
        const tag = (child as HTMLElement).tagName.toLowerCase()
        if (tag === 'ul' || tag === 'ol') {
          nested.push(child as HTMLElement)
          continue
        }
      }
      serializeInline(child, inlineParts)
    }

    const text = inlineParts.join('').trim()
    const marker = ordered ? `${index}. ` : '- '
    lines.push(`${indent}${marker}${text}`)
    for (const nestedList of nested) {
      const nestedStr = serializeList(nestedList, depth + 1)
      if (nestedStr) lines.push(nestedStr)
    }
    index++
  }

  return lines.join('\n')
}

/** GFM 表格序列化：表头行 + 分隔行 + 数据行 */
function serializeTable(table: HTMLElement): string {
  const rows: string[][] = []
  for (const tr of Array.from(table.querySelectorAll('tr'))) {
    const cells: string[] = []
    for (const cell of Array.from(tr.children)) {
      const tag = cell.tagName.toLowerCase()
      if (tag !== 'th' && tag !== 'td') continue
      const text = inlineRaw(cell).replace(/\|/g, '\\|').replace(/\n/g, '<br>')
      cells.push(text)
    }
    if (cells.length > 0) rows.push(cells)
  }
  if (rows.length === 0) return ''

  const fmtRow = (cells: string[]): string => `| ${cells.join(' | ')} |`
  const separator = rows[0].map(() => '---')
  return [
    fmtRow(rows[0]),
    fmtRow(separator),
    ...rows.slice(1).map(fmtRow),
  ].join('\n')
}

// ---------- 入口 ----------

/**
 * 将 contenteditable 根节点序列化为 Markdown 字符串。
 * 块间以空行（\n\n）分隔；返回结果末尾不含多余换行。
 */
export function domToMarkdown(root: HTMLElement): string {
  const blocks: string[] = []
  for (const child of Array.from(root.childNodes)) {
    const block = serializeBlock(child)
    if (block) blocks.push(block)
  }
  return blocks.join('\n\n')
}
