<script setup lang="ts">
import { ref } from 'vue'
import { RouterLink, useRouter } from 'vue-router'
import { AuthForm, useAuthStore } from '@/features/auth'

const auth = useAuthStore()
const router = useRouter()
const loading = ref(false)
const error = ref<string | null>(null)

async function onSubmit(email: string, password: string) {
  loading.value = true
  error.value = null
  try {
    await auth.login(email, password)
    router.push('/')
  } catch (e) {
    error.value = e instanceof Error ? e.message : 'Sign in failed'
  } finally {
    loading.value = false
  }
}
</script>

<template>
  <main class="flex min-h-screen items-center justify-center bg-[#070d1c] px-4">
    <div class="w-full max-w-sm">
      <h1 class="mb-6 text-2xl font-semibold text-slate-100">Sign in</h1>
      <AuthForm submit-label="Sign in" :loading="loading" :error="error" @submit="onSubmit" />
      <p class="mt-4 text-sm text-slate-400">
        No account?
        <RouterLink to="/register" class="text-sky-400 hover:underline">Register</RouterLink>
      </p>
    </div>
  </main>
</template>
