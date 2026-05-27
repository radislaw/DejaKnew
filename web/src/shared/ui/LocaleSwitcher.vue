<script setup lang="ts">
import { useI18n } from 'vue-i18n'
import { setLocale } from '@/shared/api'
import type { Locale } from '@/shared/config/i18n'

const { locale } = useI18n()

function switchLocale() {
  const next: Locale = locale.value === 'ru' ? 'en' : 'ru'
  locale.value = next
  localStorage.setItem('locale', next)
  setLocale(next)
}
</script>

<template>
  <button
    class="rounded-md px-2 py-1 text-xs font-medium text-slate-400 hover:text-slate-200 transition-colors"
    @click="switchLocale"
  >
    {{ locale === 'ru' ? 'EN' : 'RU' }}
  </button>
</template>
