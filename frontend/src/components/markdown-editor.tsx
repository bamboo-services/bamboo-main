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
  useEffect,
  useImperativeHandle,
  useRef,
  useState,
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
import {
  Tabs,
  TabsList,
  TabsTrigger,
} from '@/components/ui/tabs'
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger,
} from '@/components/ui/dropdown-menu'
import { cn } from '@/lib/utils'
import { renderMarkdownToHtml } from '@/lib/markdown-html'
import { domToMarkdown } from '@/lib/markdown-serializer'

/** 编辑器实例句柄：父组件保存前调用 flush 兜底同步最新内容 */
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
  /** 软限制：仅展示字符计数，不强制截断（源码模式可超限，保存由后端校验） */
  maxLength?: number
  disabled?: boolean
  className?: string
}

type Mode = 'wysiwyg' | 'source'

/**
 * 竹林水墨 Markdown 所见即所得编辑器。
 *
 * - 「预览」模式：contenteditable 渲染与公开页一致的排版，编辑即预览；
 * - 「源码」模式：纯 Markdown 源码 textarea，可手工修正；
 * - 工具栏：加粗/斜体/删除线/行内代码/链接/图片/列表/横线/块格式；
 * - 渲染单一真相源为 `MarkdownView`（经 renderMarkdownToHtml），
 *   DOM 回写经 domToMarkdown 序列化，双向保真。
 *
 * 光标铁律：聚焦（isFocused）期间绝不动 innerHTML，避免光标丢失。
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
  const [mode, setMode] = useState<Mode>('wysiwyg')
  const [draft, setDraft] = useState(value)

  const editorRef = useRef<HTMLDivElement>(null)
  const valueRef = useRef(value)
  const lastRenderedMd = useRef(value)
  const isFocused = useRef(false)
  const isComposing = useRef(false)

  // 同步最新外部 value，供 flush 比较
  useEffect(() => {
    valueRef.current = value
  }, [value])

  /** 序列化当前 DOM 并回调 onChange（内容有变时） */
  const flush = useCallback(() => {
    const el = editorRef.current
    if (!el) return
    const md = domToMarkdown(el)
    if (md !== valueRef.current) onChange(md)
  }, [onChange])

  // 暴露 flush 供父组件保存前兜底
  useImperativeHandle(ref, () => ({ flush }), [flush])

  /** 写入 innerHTML 的唯一入口：聚焦期间跳过 */
  const syncFromValue = useCallback(() => {
    const el = editorRef.current
    if (!el || isFocused.current) return
    if (value === lastRenderedMd.current) return
    el.innerHTML = renderMarkdownToHtml(value)
    lastRenderedMd.current = value
  }, [value])

  // 挂载 / 外部 value 变化 / 模式切换时同步
  useEffect(() => {
    syncFromValue()
  }, [value, mode, syncFromValue])

  /** 输入即时序列化（IME 组合期跳过），保证保存永远拿到最新值 */
  const handleInput = useCallback(() => {
    if (isComposing.current) return
    flush()
  }, [flush])

  /** 统一执行 execCommand（已废弃但全浏览器兼容，自带 undo 栈） */
  const exec = useCallback(
    (command: string, execValue?: string) => {
      editorRef.current?.focus()
      document.execCommand(command, false, execValue)
      handleInput()
    },
    [handleInput],
  )

  const setBlock = (tag: string) => exec('formatBlock', tag)

  const insertInlineCode = () =>
    exec('insertHTML', '<code>代码</code>')

  const insertCodeBlock = () =>
    exec('insertHTML', '<pre><code>代码</code></pre>')

  const insertLink = () => {
    const url = window.prompt('链接地址', 'https://')
    if (url) exec('createLink', url)
  }

  const insertImage = () => {
    const url = window.prompt('图片地址', 'https://')
    const alt = window.prompt('图片描述', '') ?? ''
    if (url) exec('insertHTML', `<img src="${url}" alt="${alt}" />`)
  }

  const handleModeChange = (next: string) => {
    if (next === mode) return
    if (next === 'source') {
      // 切源码前把 DOM 序列化写入 draft
      const el = editorRef.current
      if (el) {
        const md = domToMarkdown(el)
        setDraft(md)
        valueRef.current = md
        if (md !== value) onChange(md)
      } else {
        setDraft(valueRef.current)
      }
    }
    setMode(next as Mode)
  }

  const handleKeyDown = (e: React.KeyboardEvent) => {
    const mod = e.metaKey || e.ctrlKey
    if (mod && e.key.toLowerCase() === 'b') {
      e.preventDefault()
      exec('bold')
    } else if (mod && e.key.toLowerCase() === 'i') {
      e.preventDefault()
      exec('italic')
    } else if (mod && e.key.toLowerCase() === 'k') {
      e.preventDefault()
      insertLink()
    } else if (mod && e.shiftKey && e.key.toLowerCase() === 'x') {
      e.preventDefault()
      exec('strikeThrough')
    }
  }

  const handlePaste = (e: React.ClipboardEvent) => {
    // 只落纯文本，杜绝外站脏 HTML
    e.preventDefault()
    const text = e.clipboardData.getData('text/plain')
    document.execCommand('insertText', false, text)
    handleInput()
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
              <DropdownMenuItem onSelect={() => setBlock('P')}>
                正文
              </DropdownMenuItem>
              <DropdownMenuItem onSelect={() => setBlock('H1')}>
                标题一
              </DropdownMenuItem>
              <DropdownMenuItem onSelect={() => setBlock('H2')}>
                标题二
              </DropdownMenuItem>
              <DropdownMenuItem onSelect={() => setBlock('H3')}>
                标题三
              </DropdownMenuItem>
              <DropdownMenuItem onSelect={() => setBlock('BLOCKQUOTE')}>
                引用
              </DropdownMenuItem>
              <DropdownMenuItem onSelect={insertCodeBlock}>
                代码块
              </DropdownMenuItem>
            </DropdownMenuContent>
          </DropdownMenu>

          <ToolButton title="加粗 (⌘B)" onClick={() => exec('bold')}>
            <Bold className="size-4" />
          </ToolButton>
          <ToolButton title="斜体 (⌘I)" onClick={() => exec('italic')}>
            <Italic className="size-4" />
          </ToolButton>
          <ToolButton
            title="删除线 (⌘⇧X)"
            onClick={() => exec('strikeThrough')}
          >
            <Strikethrough className="size-4" />
          </ToolButton>

          <div className="mx-1 h-4 w-px bg-border/60" aria-hidden />

          <ToolButton title="行内代码" onClick={insertInlineCode}>
            <Code className="size-4" />
          </ToolButton>
          <ToolButton title="链接 (⌘K)" onClick={insertLink}>
            <Link className="size-4" />
          </ToolButton>
          <ToolButton title="图片" onClick={insertImage}>
            <Image className="size-4" />
          </ToolButton>

          <div className="mx-1 h-4 w-px bg-border/60" aria-hidden />

          <ToolButton title="无序列表" onClick={() => exec('insertUnorderedList')}>
            <List className="size-4" />
          </ToolButton>
          <ToolButton title="有序列表" onClick={() => exec('insertOrderedList')}>
            <ListOrdered className="size-4" />
          </ToolButton>
          <ToolButton title="横线" onClick={() => exec('insertHorizontalRule')}>
            <Minus className="size-4" />
          </ToolButton>
        </div>

        {/* 模式切换：预览 / 源码 */}
        <Tabs value={mode} onValueChange={handleModeChange} className="shrink-0">
          <TabsList variant="line" className="h-7 gap-1">
            <TabsTrigger value="wysiwyg" className="text-xs">
              预览
            </TabsTrigger>
            <TabsTrigger value="source" className="text-xs">
              源码
            </TabsTrigger>
          </TabsList>
        </Tabs>
      </div>

      {/* ───────── 编辑区：预览（contenteditable）/ 源码（textarea） ───────── */}
      {mode === 'wysiwyg' ? (
        <div
          id={id}
          ref={editorRef}
          className="markdown-view cursor-text px-4 py-3"
          contentEditable={!disabled}
          suppressContentEditableWarning
          spellCheck={false}
          data-placeholder={placeholder}
          style={{ minHeight }}
          onInput={handleInput}
          onFocus={() => {
            isFocused.current = true
          }}
          onBlur={() => {
            isFocused.current = false
            flush()
          }}
          onPaste={handlePaste}
          onKeyDown={handleKeyDown}
          onCompositionStart={() => {
            isComposing.current = true
          }}
          onCompositionEnd={() => {
            isComposing.current = false
            handleInput()
          }}
        />
      ) : (
        <Textarea
          id={id}
          value={draft}
          onChange={(e) => {
            setDraft(e.target.value)
            onChange(e.target.value)
          }}
          className="min-h-[320px] resize-y rounded-none border-0 bg-transparent font-mono text-sm shadow-none focus-visible:ring-0"
          style={{ minHeight }}
          placeholder={placeholder}
          spellCheck={false}
        />
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

/** 工具栏按钮：mousedown 阻止抢焦点，保持 contenteditable 选区 */
function ToolButton({
  title,
  onClick,
  children,
}: {
  title: string
  onClick: () => void
  children: React.ReactNode
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
