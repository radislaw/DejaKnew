import { createApp } from 'vue'
import { createPinia } from 'pinia'
import App from './App.vue'
import { router } from './router'
import { useAuthStore } from '@/features/auth'
import './styles/index.css'

const app = createApp(App)
app.use(createPinia())
app.use(router)

// Restore the session (token is read from localStorage in the store;
// this fetches the user behind it) before mounting.
useAuthStore().restore()

app.mount('#app')
