// --------------------------------------------------------------------------------
// Copyright (c) 2016-NOW(至今) 筱锋
// Author: 筱锋「xiao_lfeng」(https://www.x-lf.com)
// --------------------------------------------------------------------------------
// 许可证声明：版权所有 (c) 2016-2026 筱锋。保留所有权利。
// 有关MIT许可证的更多信息，请查看项目根目录下的LICENSE文件或访问：
// https://opensource.org/licenses/MIT
// --------------------------------------------------------------------------------

import {
  forwardRef,
  useCallback,
  useImperativeHandle,
  useRef,
  useState,
  type KeyboardEvent,
  type ReactNode,
} from 'react'
import {
  Bold,
  Code,
  Heading1,
  Image,
  Italic,
  Link,
  List,
  ListOrdered,
  Minus,
  Strikethrough,
} from 'lucide-react'
import { Button } from '@/components/ui/button'
import { Textarea } from '@/components/ui/textarea'
import { Tabs, TabsList, TabsTrigger } from '@/components/ui/tabs'
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger,
} from '@/components/ui/dropdown-menu'
import { MarkdownView } from '@/components/markdown'
import { cn } from '@/lib/utils'

/** 编辑器实例句柄：源码模式受控同步，flush 保留为空操作以兼容父组件调用 */
export interface MarkdownEditorHandle {
  flush: () => void
}

export interface MarkdownEditorProps {
  id?: string
  value: string
  onChange: (value: string) => void
  /** 编辑区最小高度（px），默认 320 */
  minHeight?: number
  placeholder?: string
  /** 软限制：仅展示字符计数，不强制截断（保存由后端校验） */
  maxLength?: number
  disabled?: boolean
  className?: string
}

type Mode = 'source' | 'preview'

/**
 * 竹林水墨 Markdown 编辑器。
 *
 * - 「源码」模式（默认）：Markdown 源码 textarea，受控编辑，内容永不丢失；
 * - 「预览」模式：复用 `MarkdownView` 只读渲染，所见即公开页排版；
 * - 工具栏：加粗/斜体/删除线/行内代码/链接/图片/列表/横线/块格式，
 *   全部基于 textarea 选区包裹，纯文本编辑、无 DOM 双向同步风险。
 */
export const MarkdownEditor = forwardRef<
  MarkdownEditorHandle,
  MarkdownEditorProps
>(function MarkdownEditor(
  {
    id,
    value,
    onChange,
    minHeight = 320,
    placeholder,
    maxLength,
    disabled,
    className,
  },
  ref,
) {
  const [mode, setMode] = useState<Mode>('source')
  const taRef = useRef<HTMLTextAreaElement>(null)

  // 源码模式 value 已受控同步，flush 无需额外操作
  useImperativeHandle(ref, () => ({ flush: () => {} }), [])

  /** 下一帧恢复 textarea 焦点与光标位置（等 React 更新 value） */
  const restoreSelection = useCallback((pos: number) => {
    requestAnimationFrame(() => {
      const el = taRef.current
      if (!el) return
      el.focus()
      el.setSelectionRange(pos, pos)
    })
  }, [])

  /** 选区包裹：在选区前后插入标记，未选时以 hint 占位 */
  const wrapSelection = useCallback(
    (before: string, after: string, hint?: string) => {
      const el = taRef.current
      if (!el) return
      const start = el.selectionStart
      const end = el.selectionEnd
      const insert = value.slice(start, end) || hint || ''
      onChange(value.slice(0, start) + before + insert + after + value.slice(end))
      restoreSelection(start + before.length + insert.length)
    },
    [value, onChange, restoreSelection],
  )

  /** 行首前缀：光标所在行加前缀，已有则移除（toggle） */
  const toggleBlockPrefix = useCallback(
    (prefix: string) => {
      const el = taRef.current
      if (!el) return
      const start = el.selectionStart
      const lineStart = value.lastIndexOf('\n', start - 1) + 1
      const lineEndRel = value.indexOf('\n', start)
      const lineEnd = lineEndRel === -1 ? value.length : lineEndRel
      const line = value.slice(lineStart, lineEnd)
      let next: string
      let pos: number
      if (line.startsWith(prefix)) {
        next = value.slice(0, lineStart) + line.slice(prefix.length) + value.slice(lineEnd)
        pos = lineStart
      } else {
        next = value.slice(0, lineStart) + prefix + line + value.slice(lineEnd)
        pos = lineStart + prefix.length
      }
      onChange(next)
      restoreSelection(pos)
    },
    [value, onChange, restoreSelection],
  )

  /** 正文：移除光标所在行的标题/列表/引用标记 */
  const clearBlockPrefix = useCallback(() => {
    const el = taRef.current
    if (!el) return
    const start = el.selectionStart
    const lineStart = value.lastIndexOf('\n', start - 1) + 1
    const lineEndRel = value.indexOf('\n', start)
    const lineEnd = lineEndRel === -1 ? value.length : lineEndRel
    const line = value.slice(lineStart, lineEnd)
    const next =
      value.slice(0, lineStart) +
      line.replace(/^(#{1,6}\s|>\s|[-+*]\s|\d+\.\s)/, '') +
      value.slice(lineEnd)
    onChange(next)
    restoreSelection(lineStart)
  }, [value, onChange, restoreSelection])

  /** 插入横线 */
  const insertHorizontalRule = useCallback(() => {
    const el = taRef.current
    if (!el) return
    const start = el.selectionStart
    const end = el.selectionEnd
    const insert = '\n\n---\n\n'
    onChange(value.slice(0, start) + insert + value.slice(end))
    restoreSelection(start + insert.length)
  }, [value, onChange, restoreSelection])

  const insertLink = useCallback(() => {
    const url = window.prompt('链接地址', 'https://')
    if (url) wrapSelection('[', `](${url})`, '链接文字')
  }, [wrapSelection])

  const insertImage = useCallback(() => {
    const url = window.prompt('图片地址', 'https://')
    if (url) wrapSelection('![', `](${url})`, '图片描述')
  }, [wrapSelection])

  const handleKeyDown = (e: KeyboardEvent<HTMLTextAreaElement>) => {
    const mod = e.metaKey || e.ctrlKey
    if (mod && e.key.toLowerCase() === 'b') {
      e.preventDefault()
      wrapSelection('**', '**', '加粗文字')
    } else if (mod && e.key.toLowerCase() === 'i') {
      e.preventDefault()
      wrapSelection('*', '*', '斜体文字')
    } else if (mod && e.key.toLowerCase() === 'k') {
      e.preventDefault()
      insertLink()
    } else if (mod && e.shiftKey && e.key.toLowerCase() === 'x') {
      e.preventDefault()
      wrapSelection('~~', '~~', '删除文字')
    }
  }

  return (
    <div
      className={cn(
        'markdown-editor overflow-hidden rounded-lg border border-field-line bg-card transition-[box-shadow,border-color]',
        'focus-within:border-leaf-deep focus-within:ring-[3px] focus-within:ring-ring/50',
        disabled && 'pointer-events-none opacity-50',
        className,
      )}
    >
      {/* ───────── 工具条：格式按钮 + 模式切换 ───────── */}
      <div className="flex flex-wrap items-center justify-between gap-2 border-b border-border/60 bg-muted/40 px-2 py-1.5">
        <div className="flex flex-wrap items-center gap-0.5">
          {/* 块格式：正文/标题/引用/代码块 */}
          <DropdownMenu>
            <DropdownMenuTrigger asChild>
              <Button
                variant="ghost"
                size="icon-sm"
                title="块格式"
                className="cursor-pointer"
                onMouseDown={(e) => e.preventDefault()}
              >
                <Heading1 className="size-4" />
              </Button>
            </DropdownMenuTrigger>
            <DropdownMenuContent align="start" className="min-w-32">
              <DropdownMenuItem onSelect={clearBlockPrefix}>
                正文
              </DropdownMenuItem>
              <DropdownMenuItem onSelect={() => toggleBlockPrefix('# ')}>
                标题一
              </DropdownMenuItem>
              <DropdownMenuItem onSelect={() => toggleBlockPrefix('## ')}>
                标题二
              </DropdownMenuItem>
              <DropdownMenuItem onSelect={() => toggleBlockPrefix('### ')}>
                标题三
              </DropdownMenuItem>
              <DropdownMenuItem onSelect={() => toggleBlockPrefix('> ')}>
                引用
              </DropdownMenuItem>
              <DropdownMenuItem
                onSelect={() => wrapSelection('```\n', '\n```', '代码')}
              >
                代码块
              </DropdownMenuItem>
            </DropdownMenuContent>
          </DropdownMenu>

          <ToolButton
            title="加粗 (⌘B)"
            onClick={() => wrapSelection('**', '**', '加粗文字')}
          >
            <Bold className="size-4" />
          </ToolButton>
          <ToolButton
            title="斜体 (⌘I)"
            onClick={() => wrapSelection('*', '*', '斜体文字')}
          >
            <Italic className="size-4" />
          </ToolButton>
          <ToolButton
            title="删除线 (⌘⇧X)"
            onClick={() => wrapSelection('~~', '~~', '删除文字')}
          >
            <Strikethrough className="size-4" />
          </ToolButton>

          <div className="mx-1 h-4 w-px bg-border/60" aria-hidden />

          <ToolButton
            title="行内代码"
            onClick={() => wrapSelection('`', '`', '代码')}
          >
            <Code className="size-4" />
          </ToolButton>
          <ToolButton title="链接 (⌘K)" onClick={insertLink}>
            <Link className="size-4" />
          </ToolButton>
          <ToolButton title="图片" onClick={insertImage}>
            <Image className="size-4" />
          </ToolButton>

          <div className="mx-1 h-4 w-px bg-border/60" aria-hidden />

          <ToolButton
            title="无序列表"
            onClick={() => toggleBlockPrefix('- ')}
          >
            <List className="size-4" />
          </ToolButton>
          <ToolButton
            title="有序列表"
            onClick={() => toggleBlockPrefix('1. ')}
          >
            <ListOrdered className="size-4" />
          </ToolButton>
          <ToolButton title="横线" onClick={insertHorizontalRule}>
            <Minus className="size-4" />
          </ToolButton>
        </div>

        {/* 模式切换：源码（默认）/ 预览（只读） */}
        <Tabs
          value={mode}
          onValueChange={(v) => setMode(v as Mode)}
          className="shrink-0"
        >
          <TabsList variant="line" className="h-7 gap-1">
            <TabsTrigger value="source" className="text-xs">
              源码
            </TabsTrigger>
            <TabsTrigger value="preview" className="text-xs">
              预览
            </TabsTrigger>
          </TabsList>
        </Tabs>
      </div>

      {/* ───────── 编辑区：源码（textarea）/ 预览（只读 MarkdownView） ───────── */}
      {mode === 'source' ? (
        <Textarea
          id={id}
          ref={taRef}
          value={value}
          onChange={(e) => onChange(e.target.value)}
          className="resize-y rounded-none border-0 bg-transparent font-mono text-sm shadow-none focus-visible:ring-0"
          style={{ minHeight }}
          placeholder={placeholder}
          spellCheck={false}
          disabled={disabled}
          onKeyDown={handleKeyDown}
        />
      ) : (
        <div
          className="max-h-[480px] overflow-y-auto px-4 py-3"
          style={{ minHeight }}
        >
          <MarkdownView content={value} />
        </div>
      )}

      {/* ───────── 字数统计 ───────── */}
      {maxLength != null && (
        <div className="border-t border-border/60 px-3 py-1 text-right font-mono text-[11px] text-text-secondary">
          {value.length} / {maxLength}
        </div>
      )}
    </div>
  )
})

/** 工具栏按钮：mousedown 阻止抢焦点，保持 textarea 选区 */
function ToolButton({
  title,
  onClick,
  children,
}: {
  title: string
  onClick: () => void
  children: ReactNode
}) {
  return (
    <Button
      variant="ghost"
      size="icon-sm"
      title={title}
      className="cursor-pointer text-text-secondary hover:bg-leaf-light/40 hover:text-leaf-deep"
      onMouseDown={(e) => e.preventDefault()}
      onClick={onClick}
    >
      {children}
    </Button>
  )
}
