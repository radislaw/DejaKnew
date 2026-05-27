<script setup lang="ts">
import { ref } from 'vue'
import { useI18n } from 'vue-i18n'

const { t } = useI18n()

defineProps<{
  submitLabel: string
  loading?: boolean
  error?: string | null
}>()

const emit = defineEmits<{
  submit: [email: string, password: string]
}>()

const email = ref('')
const password = ref('')
</script>

<template>
  <form class="flex flex-col gap-4" @submit.prevent="emit('submit', email, password)">
    <input
      v-model="email"
      type="email"
      required
      :placeholder="t('auth.emailPlaceholder')"
      class="rounded-lg border border-slate-700 bg-slate-900 px-4 py-2.5 text-slate-200 placeholder:text-slate-500 outline-none focus:border-sky-400"
    />
    <input
      v-model="password"
      type="password"
      required
      minlength="8"
      :placeholder="t('auth.passwordPlaceholder')"
      class="rounded-lg border border-slate-700 bg-slate-900 px-4 py-2.5 text-slate-200 placeholder:text-slate-500 outline-none focus:border-sky-400"
    />
    <p v-if="error" class="text-sm text-rose-400">{{ error }}</p>
    <button
      type="submit"
      :disabled="loading"
      class="rounded-lg bg-sky-400 px-4 py-2.5 font-medium text-slate-950 transition hover:bg-sky-300 disabled:opacity-50"
    >
      {{ loading ? '…' : submitLabel }}
    </button>
  </form>
</template>
