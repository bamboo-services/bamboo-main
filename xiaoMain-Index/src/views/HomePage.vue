<!--
  - -------------------------------------------------------------------------------
  - Copyright (c) 2016-NOW(至今) 筱锋
  - Author: 筱锋(https://www.x-lf.com)
  -
  - 本文件包含 XiaoMain 的源代码，该项目的所有源代码均遵循MIT开源许可证协议。
  - -------------------------------------------------------------------------------
  - 许可证声明：
  -
  - 版权所有 (c) 2016-2024 筱锋。保留所有权利。
  -
  - 本软件是“按原样”提供的，没有任何形式的明示或暗示的保证，包括但不限于
  - 对适销性、特定用途的适用性和非侵权性的暗示保证。在任何情况下，
  - 作者或版权持有人均不承担因软件或软件的使用或其他交易而产生的、
  - 由此引起的或以任何方式与此软件有关的任何索赔、损害或其他责任。
  -
  - 使用本软件即表示您了解此声明并同意其条款。
  -
  - 有关MIT许可证的更多信息，请查看项目根目录下的LICENSE文件或访问：
  - https://opensource.org/licenses/MIT
  - -------------------------------------------------------------------------------
  - 免责声明：
  -
  - 使用本软件的风险由用户自担。作者或版权持有人在法律允许的最大范围内，
  - 对因使用本软件内容而导致的任何直接或间接的损失不承担任何责任。
  - -------------------------------------------------------------------------------
  -->

<template>
  <section>
    <div class="mx-auto max-w-screen-xl px-4 py-32 lg:flex lg:h-screen lg:items-center">
      <div class="mx-auto max-w-xxl text-center">
        <div class="grid grid-cols-4 gap-3">
          <div class="col-span-1 hidden md:block">
            <div class="flex items-center">
              <img alt="UserAvatar" class="rounded-full w-auto h-auto xl:h-64" src="https://api.x-lf.cn/avatar/?uid=1">
            </div>
          </div>
          <div class="col-span-3 h-full flex items-center">
            <div>
              <transition-group>
                <h1 class="bg-gradient-to-r from-green-300 via-blue-500 to-purple-600 bg-clip-text text-3xl font-extrabold text-transparent sm:text-5xl"
                    style="text-shadow: 1px 1px 2px rgba(38,164,192,0.32)">
                  凌中的锋雨
                </h1>
                <p class="mt-4 sm:text-xl/relaxed">
                  愿你的人生如璀璨星辰，勇敢梦想，坚韧追寻，发现真我之光。在每个转角，以坚定意志和无尽创意，拥抱挑战，种下希望。遇见激励之人，共成长；面对风雨，保持乐观，让每一步都踏出意义深远的足迹。在追梦途中，听从内心之声，珍惜遇见，让生活不仅是冒险，更是自我发现的诗篇。
                </p>

                <div class="mt-8 flex flex-wrap justify-center gap-4">
                  <a class="block w-full rounded bg-blue-500 px-12 py-3 text-sm font-medium text-white hover:bg-blue-600 focus:outline-none focus:ring sm:w-auto shadow-xl shadow-blue-500/50"
                     @click="() => window.location.replace('https://blog.x-lf.com')"
                  >
                    去我的博客吧
                  </a>
                  <a class="block w-full rounded bg-green-500 px-12 py-3 text-sm font-medium text-white hover:bg-green-600 focus:outline-none focus:ring sm:w-auto shadow-xl shadow-green-500/55"
                     @click="() => router.push({name: 'AboutMe', replace: true})">
                    了解我的更多
                  </a>
                </div>
              </transition-group>
            </div>
          </div>
        </div>
      </div>
    </div>
    <footer class="absolute inset-x-0 bottom-0 flex justify-between items-end p-3 text-gray-500">
      <div class="grid">
        <span @click="() => router.push({name: 'AccountLogin', replace: true})">账户登录</span>
        <span>Copyright (C) 2016-2024 筱锋xiao_lfeng. All Rights Reserved.</span>
      </div>
      <div class="grid text-end">
        <span><a href="">粤ICP备 2022014822 号</a></span>
        <span><a href="">粤公网安备 44030702003207 号</a></span>
      </div>
    </footer>
  </section>
</template>

<style>
body {
  background: url("@/assets/default-background.webp") no-repeat center;
  width: 100%
}
</style>

<script lang="ts" setup>
import router from "@/router";
</script>