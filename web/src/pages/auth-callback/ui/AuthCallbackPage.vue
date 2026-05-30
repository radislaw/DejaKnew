<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useAuthStore } from '@/features/auth'

const route = useRoute()
const router = useRouter()
const auth = useAuthStore()
const errorMsg = ref<string | null>(null)

onMounted(async () => {
  const token = route.query.token as string | undefined
  const error = route.query.error as string | undefined

  if (error) {
    errorMsg.value = error
    setTimeout(() => router.replace({ name: 'login' }), 2500)
    return
  }

  if (token) {
    try {
      await auth.loginWithToken(token)
      router.replace('/')
    } catch {
      router.replace({ name: 'login' })
    }
  } else {
    router.replace({ name: 'login' })
  }
})
</script>

<template>
  <main class="flex min-h-screen items-center justify-center bg-[#070d1c]">
    <p v-if="errorMsg" class="text-rose-400">{{ errorMsg }}</p>
    <p v-else class="text-slate-400">...</p>
  </main>
</template>
