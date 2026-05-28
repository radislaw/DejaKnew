<script setup lang="ts">
import { ref } from 'vue'
import { useI18n } from 'vue-i18n'
import { setLocale } from '@/shared/api'
import type { Locale } from '@/shared/config/i18n'

const { locale } = useI18n()
const open = ref(false)

const options: { value: Locale; label: string; flag: string }[] = [
  { value: 'ru', label: 'Русский', flag: '🇷🇺' },
  { value: 'en', label: 'English', flag: '🇬🇧' },
]

function select(lang: Locale) {
  locale.value = lang
  localStorage.setItem('locale', lang)
  setLocale(lang)
  open.value = false
}

function current() {
  return options.find(o => o.value === locale.value) ?? options[0]
}
</script>

<template>
  <div class="relative" @keydown.escape="open = false">
    <button
      class="flex items-center gap-1.5 rounded-md px-2.5 py-1.5 text-sm text-slate-400 hover:text-slate-200 hover:bg-white/5 transition-colors"
      @click="open = !open"
    >
      <span>{{ current().flag }}</span>
      <span class="font-medium uppercase tracking-wide text-xs">{{ current().value }}</span>
      <svg
        class="size-3 transition-transform"
        :class="open ? 'rotate-180' : ''"
        viewBox="0 0 12 12"
        fill="none"
        xmlns="http://www.w3.org/2000/svg"
      >
        <path d="M2 4l4 4 4-4" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round"/>
      </svg>
    </button>

    <Transition
      enter-active-class="transition duration-100 ease-out"
      enter-from-class="opacity-0 scale-95"
      enter-to-class="opacity-100 scale-100"
      leave-active-class="transition duration-75 ease-in"
      leave-from-class="opacity-100 scale-100"
      leave-to-class="opacity-0 scale-95"
    >
      <div
        v-if="open"
        class="absolute right-0 top-full mt-1 min-w-[130px] rounded-lg border border-white/10 bg-[#0f1829] py-1 shadow-xl z-50"
      >
        <button
          v-for="opt in options"
          :key="opt.value"
          class="flex w-full items-center gap-2.5 px-3 py-2 text-sm transition-colors"
          :class="opt.value === locale
            ? 'text-sky-400 bg-sky-500/10'
            : 'text-slate-300 hover:bg-white/5 hover:text-slate-100'"
          @click="select(opt.value)"
        >
          <span>{{ opt.flag }}</span>
          <span>{{ opt.label }}</span>
          <svg
            v-if="opt.value === locale"
            class="ml-auto size-3.5 text-sky-400"
            viewBox="0 0 14 14"
            fill="none"
          >
            <path d="M2 7l3.5 3.5L12 3" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round"/>
          </svg>
        </button>
      </div>
    </Transition>

    <div
      v-if="open"
      class="fixed inset-0 z-40"
      @click="open = false"
    />
  </div>
</template>
