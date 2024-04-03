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