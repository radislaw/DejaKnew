import { createApp } from 'vue'
import { createPinia } from 'pinia'
import App from './App.vue'
import { router } from './router'
import { useAuthStore } from '@/features/auth'
import { i18n } from '@/shared/config'
import { setLocale } from '@/shared/api'
import './styles/index.css'

const app = createApp(App)
app.use(createPinia())
app.use(router)
app.use(i18n)

const detectedLocale = (i18n.global.locale as unknown as { value: string }).value
setLocale(detectedLocale)

// Restore the session (token is read from localStorage in the store;
// this fetches the user behind it) before mounting.
useAuthStore().restore()

app.mount('#app')
