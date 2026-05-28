<script setup lang="ts">
import {ref} from 'vue'
import {RouterLink, useRouter} from 'vue-router'
import {useI18n} from 'vue-i18n'
import {AuthForm, useAuthStore} from '@/features/auth'

const {t} = useI18n()
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
    error.value = e instanceof Error ? e.message : t('auth.signInFailed')
  } finally {
    loading.value = false
  }
}
</script>

<template>
  <main class="flex min-h-screen items-center justify-center bg-[#070d1c] px-4">
    <div class="w-full max-w-sm">
      <h1 class="mb-6 text-2xl font-semibold text-slate-100">{{ t('auth.signInTitle') }}</h1>
      <AuthForm :submit-label="t('auth.signInButton')" :loading="loading" :error="error" @submit="onSubmit"/>
      <p class="mt-4 text-sm text-slate-400">
        {{ t('auth.noAccount') }}
        <RouterLink to="/register" class="text-sky-400 hover:underline">{{ t('auth.registerButton') }}</RouterLink>
      </p>
    </div>
  </main>
</template>
