<script setup lang="ts">
import { useRouter } from 'vue-router'
import { useI18n } from 'vue-i18n'
import { useAuthStore } from '@/features/auth'
import { LocaleSwitcher } from '@/shared/ui'

const { t } = useI18n()
const router = useRouter()
const auth = useAuthStore()
</script>

<template>
  <div class="min-h-screen bg-[#070d1c] px-4 py-6 text-slate-100">
    <header class="mx-auto mb-8 flex max-w-4xl items-center justify-between">
      <button
        class="text-xl font-semibold tracking-tight hover:text-slate-300 transition-colors"
        @click="router.push({ name: 'home' })"
      >
        DéjàKnew
      </button>
      <div class="flex items-center gap-3">
        <button
          class="text-sm text-slate-400 hover:text-slate-200 transition-colors"
          @click="router.push({ name: 'settings' })"
        >{{ auth.user?.email }}</button>
        <button
          class="rounded-lg px-3 py-1.5 text-sm text-slate-400 hover:text-slate-200"
          @click="auth.logout(); router.push('/login')"
        >
          {{ t('auth.signOut') }}
        </button>
        <LocaleSwitcher />
      </div>
    </header>
    <slot />
  </div>
</template>
