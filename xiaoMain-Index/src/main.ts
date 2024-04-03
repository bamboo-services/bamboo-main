import "@/assets/tailwind.css"

import { createApp } from 'vue'
import Index from '@/Index.vue'
import router from '@/router'

const app = createApp(Index)

app
    .use(router)
    .mount('#app')
