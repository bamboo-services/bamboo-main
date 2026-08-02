// --------------------------------------------------------------------------------
// Copyright (c) 2016-NOW(至今) 筱锋
// Author: 筱锋「xiao_lfeng」(https://www.x-lf.com)
// --------------------------------------------------------------------------------
// 许可证声明：版权所有 (c) 2016-2026 筱锋。保留所有权利。
// 有关MIT许可证的更多信息，请查看项目根目录下的LICENSE文件或访问：
// https://opensource.org/licenses/MIT
// --------------------------------------------------------------------------------

import { useEffect, useState } from 'react'
import { useInView } from 'react-intersection-observer'
import { cn } from '@/lib/utils'

/**
 * 懒加载图片（IntersectionObserver）——友链站点截图专用。
 *
 * 元素进入视口（默认提前 400px 预加载）才开始请求图片，避免整页友链截图一次性拉取；
 * 加载完成前显示墨晕渐变占位，加载完成后克制淡入。加载失败时停留占位，不重复请求。
 */
export function LazyImage({
  src,
  alt,
  className,
  rootMargin = '400px',
}: {
  src: string
  alt: string
  className?: string
  /** 提前加载距离（px）：元素进入视口前 rootMargin 距离即开始加载 */
  rootMargin?: string
}) {
  const [loaded, setLoaded] = useState(false)
  const [failed, setFailed] = useState(false)
  const { ref, inView } = useInView({ triggerOnce: true, rootMargin })

  // 图片地址变化时重置加载状态（如友链切换或重新截图后 URL 更新）
  useEffect(() => {
    setLoaded(false)
    setFailed(false)
  }, [src])

  const ready = inView && loaded && !failed

  return (
    <div ref={ref} className={cn('relative overflow-hidden', className)}>
      {/* 加载前 / 加载中 / 加载失败占位：墨晕渐变示意 */}
      <div
        aria-hidden={ready}
        className={cn(
          'absolute inset-0 flex items-center justify-center bg-gradient-to-br from-leaf-light/30 via-card to-leaf-muted/25 transition-opacity duration-500',
          ready && 'pointer-events-none opacity-0',
        )}
      >
        <span className="font-mono text-xs text-text-secondary">
          站点截图 · 占位
        </span>
      </div>

      {/* 进入视口后才请求图片，加载完成淡入 */}
      {inView && (
        <img
          src={src}
          alt={alt}
          loading="lazy"
          onLoad={() => setLoaded(true)}
          onError={() => setFailed(true)}
          className={cn(
            'h-full w-full object-cover transition-opacity duration-500',
            ready ? 'opacity-100' : 'opacity-0',
          )}
        />
      )}
    </div>
  )
}
