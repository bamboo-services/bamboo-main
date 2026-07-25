/*
 * --------------------------------------------------------------------------------
 * Copyright (c) 2016-NOW 筱锋
 * Author: 筱锋(https://www.x-lf.com)
 * License: MIT
 * --------------------------------------------------------------------------------
 */

import { createFileRoute, Link } from '@tanstack/react-router'
import { motion } from 'motion/react'
import { BambooLogo } from '@/assets/svg/bamboo-logo'
import { siteInfo } from '@/data/mock/site-info'
import myAvatar from '@/assets/images/my_avatar.png'
import defaultBackground from '@/assets/images/default-background.webp'

export const Route = createFileRoute('/_public/')({
  component: HomePage,
})

function HomePage() {
  const thisYear = new Date().getFullYear()

  return (
    <div className="overflow-hidden">
      <motion.div
        initial={{ opacity: 0, scale: 3 }}
        animate={{ opacity: 1, scale: 1 }}
        transition={{ type: 'spring', stiffness: 70, damping: 26 }}
        className="h-dvh grid"
        style={{
          backgroundImage: `url(${defaultBackground})`,
          backgroundSize: 'cover',
          backgroundPosition: 'center',
          backgroundRepeat: 'no-repeat',
        }}
      >
        <div className="grid grid-cols-4 gap-8 justify-center items-center px-8 lg:px-32">
          {/* 桌面端头像 */}
          <motion.div
            initial={{ opacity: 0, x: -40 }}
            animate={{ opacity: 1, x: 0 }}
            transition={{ type: 'spring', stiffness: 100, damping: 26, delay: 0.2 }}
            className="hidden lg:flex col-span-1 items-center"
          >
            <img
              alt="UserAvatar"
              className="rounded-full size-auto object-cover lg:h-48 xl:h-64"
              src={myAvatar}
              draggable={false}
            />
          </motion.div>

          {/* 内容区 */}
          <div className="col-span-4 lg:col-span-3">
            <div className="text-center grid gap-3">
              {/* 标题 */}
              <motion.div
                initial={{ opacity: 0, x: 40 }}
                animate={{ opacity: 1, x: 0 }}
                transition={{ type: 'spring', stiffness: 100, damping: 26, delay: 0.2 }}
                className="flex justify-center"
              >
                <h1
                  className="bg-gradient-to-r from-green-300 via-blue-500 to-purple-600 bg-clip-text text-4xl font-extrabold text-transparent sm:text-5xl flex gap-3 items-center"
                  style={{ textShadow: '1px 1px 4px rgba(38,164,192,0.32)' }}
                >
                  <BambooLogo size={52} />
                  <span>{siteInfo.site.siteName}</span>
                </h1>
              </motion.div>

              {/* 移动端头像 */}
              <motion.div
                initial={{ opacity: 0, x: 40 }}
                animate={{ opacity: 1, x: 0 }}
                transition={{ type: 'spring', stiffness: 100, damping: 26, delay: 0.7 }}
                className="lg:hidden items-center flex justify-center"
              >
                <img
                  alt="UserAvatar"
                  className="rounded-xl w-auto h-32"
                  src={myAvatar}
                  draggable={false}
                />
              </motion.div>

              {/* 描述 */}
              <motion.p
                initial={{ opacity: 0, x: 40 }}
                animate={{ opacity: 1, x: 0 }}
                transition={{ type: 'spring', stiffness: 100, damping: 26, delay: 1.2 }}
                className="lg:text-xl/relaxed"
              >
                {siteInfo.blogger.description}
              </motion.p>

              {/* 按钮组 */}
              <motion.div
                initial={{ opacity: 0, x: 40 }}
                animate={{ opacity: 1, x: 0 }}
                transition={{ type: 'spring', stiffness: 100, damping: 26, delay: 1.7 }}
                className="flex flex-wrap justify-center gap-4 pt-8 md:pt-6 lg:pt-4"
              >
                <a
                  href="https://blog.x-lf.com"
                  target="_blank"
                  rel="noopener noreferrer"
                  className="transition block rounded bg-blue-500 px-12 py-3 text-sm font-medium text-white hover:bg-blue-600 focus:outline-none focus:ring sm:w-auto shadow-xl shadow-blue-500/50"
                >
                  去我的博客吧
                </a>
                <Link
                  to="/auth/login"
                  className="transition block rounded bg-green-500 px-12 py-3 text-sm font-medium text-white hover:bg-green-600 focus:outline-none focus:ring sm:w-auto shadow-xl shadow-green-500/55"
                >
                  了解我的更多
                </Link>
              </motion.div>
            </div>
          </div>
        </div>

        {/* 桌面端底部 */}
        <footer className="hidden md:flex absolute inset-x-0 bottom-0 justify-between items-end p-3 text-gray-500">
          <motion.div
            initial={{ opacity: 0, x: -40 }}
            animate={{ opacity: 1, x: 0 }}
            transition={{ type: 'spring', stiffness: 100, damping: 26, delay: 2.2 }}
            className="grid"
          >
            <Link to="/admin/dashboard" className="hover:text-gray-700 transition">
              账户登录
            </Link>
            <span>Copyright (C) 2016-{thisYear} 筱锋xiao_lfeng. All Rights Reserved.</span>
          </motion.div>
          <motion.div
            initial={{ opacity: 0, x: 40 }}
            animate={{ opacity: 1, x: 0 }}
            transition={{ type: 'spring', stiffness: 100, damping: 26, delay: 2.2 }}
            className="grid text-end"
          >
            <a
              href="https://beian.miit.gov.cn/#/Integrated/index"
              target="_blank"
              rel="noopener noreferrer"
              className="hover:text-gray-700 transition"
            >
              粤ICP备 2022014822 号
            </a>
            <a
              href="https://beian.mps.gov.cn/#/query/webSearch"
              target="_blank"
              rel="noopener noreferrer"
              className="hover:text-gray-700 transition"
            >
              粤公网安备 44030702003207 号
            </a>
          </motion.div>
        </footer>

        {/* 移动端底部 */}
        <footer className="grid md:hidden text-gray-500 text-center pt-8 absolute inset-x-0 bottom-0 pb-3">
          <Link to="/auth/login" className="hover:text-gray-700 transition">
            账户登录
          </Link>
          <a
            href="https://beian.miit.gov.cn/#/Integrated/index"
            target="_blank"
            rel="noopener noreferrer"
            className="hover:text-gray-700 transition"
          >
            粤ICP备 2022014822 号
          </a>
        </footer>
      </motion.div>
    </div>
  )
}
