/*
 * --------------------------------------------------------------------------------
 * Copyright (c) 2016-NOW(至今) 筱锋
 * Author: 筱锋「xiao_lfeng」(https://www.x-lf.com)
 * --------------------------------------------------------------------------------
 * 许可证声明：版权所有 (c) 2016-2026 筱锋。保留所有权利。
 * 有关MIT许可证的更多信息，请查看项目根目录下的LICENSE文件或访问：
 * https://opensource.org/licenses/MIT
 * --------------------------------------------------------------------------------
 */

interface BambooLogoProps {
  size?: number
  className?: string
}

export function BambooLogo({ size = 32, className }: BambooLogoProps) {
  return (
    <svg
      className={className}
      viewBox="0 0 1024 1024"
      xmlns="http://www.w3.org/2000/svg"
      width={size}
      height={size}
    >
      {/* 竹节三档墨绿：亮/中/深 → chart-1 / chart-3 / leaf-deep，与水墨 token 同源 */}
      <path d="M482 297l-68.1 39.4L482 376l68.4-39.6z" fill="var(--chart-1)" />
      <path
        d="M482.3 375.9l-0.2 552.3-68.3-39.3 0.3-552.5z"
        fill="var(--chart-3)"
      />
      <path
        d="M482 375.9l0.2 552.3 68.4-39.3-0.4-552.5z"
        fill="var(--leaf-deep)"
      />
      <path
        d="M482 98.2l-68.1 39.3 68.1 39.6 68.4-39.6z"
        fill="var(--chart-1)"
      />
      <path
        d="M482.3 177l-0.2 146-68.3-39.4 0.3-146.1z"
        fill="var(--chart-3)"
      />
      <path d="M482 177l0.2 146 68.4-39.4-0.4-146.1z" fill="var(--leaf-deep)" />
      <path
        d="M693.8 191l-68.2 39.4 68.2 39.6 68.3-39.6z"
        fill="var(--chart-1)"
      />
      <path
        d="M694.1 269.9l-0.2 552.3-68.4-39.3 0.4-552.5z"
        fill="var(--chart-3)"
      />
      <path
        d="M693.8 269.9l0.2 552.3 68.3-39.3-0.3-552.5z"
        fill="var(--leaf-deep)"
      />
      <path d="M766 199l18 35 39.1 1.7L861 110.9z" fill="var(--chart-3)" />
      <path d="M861 110.9L784 234l39.1 1.7z" fill="var(--leaf-deep)" />
      <path
        d="M796.4 250.7l-7.3 38.7 30 25.3 106.4-75.3z"
        fill="var(--chart-3)"
      />
      <path d="M925.5 239.4l-136.4 50 30 25.3z" fill="var(--leaf-deep)" />
      <path
        d="M322.7 482.1l-68.2 39.4 68.2 39.5 68.3-39.5z"
        fill="var(--chart-1)"
      />
      <path
        d="M323 561l-0.2 336.2-68.4-39.4 0.4-336.3z"
        fill="var(--chart-3)"
      />
      <path
        d="M322.7 561l0.2 336.2 68.3-39.4-0.3-336.3z"
        fill="var(--leaf-deep)"
      />
      <path d="M223.9 398.3l2 39.4-33.1 21-95.4-88.9z" fill="var(--chart-3)" />
      <path d="M97.4 369.8l128.5 67.9-33.1 21z" fill="var(--leaf-deep)" />
      <path
        d="M223.5 458.3l25.7 29.9-13.3 36.9-129.8-11.8z"
        fill="var(--chart-3)"
      />
      <path d="M106.1 513.3l143.1-25.1-13.3 36.9z" fill="var(--leaf-deep)" />
    </svg>
  )
}
