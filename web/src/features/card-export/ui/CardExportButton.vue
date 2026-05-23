<script setup lang="ts">
import { ref } from 'vue'
import { exportCards, type ExportFormat } from '@/entities/card'

const props = defineProps<{ deckId?: string }>()

const formats: { value: ExportFormat; label: string }[] = [
  { value: 'csv', label: 'CSV' },
  { value: 'json', label: 'JSON' },
  { value: 'anki-csv', label: 'Anki CSV' },
]

const loading = ref(false)
const open = ref(false)
const error = ref<string | null>(null)

async function download(format: ExportFormat) {
  open.value = false
  loading.value = true
  error.value = null
  try {
    const { blob, filename } = await exportCards(format, props.deckId)
    const url = URL.createObjectURL(blob)
    const a = document.createElement('a')
    a.href = url
    a.download = filename
    a.click()
    URL.revokeObjectURL(url)
  } catch (e) {
    error.value = e instanceof Error ? e.message : 'Export failed'
  } finally {
    loading.value = false
  }
}
</script>

<template>
  <div class="relative">
    <button
      class="rounded-lg border border-slate-600 px-4 py-1.5 text-sm font-medium text-slate-300 hover:border-slate-500 hover:text-slate-100 disabled:opacity-50"
      :disabled="loading"
      @click="open = !open"
    >
      {{ loading ? 'Exporting…' : 'Export' }}
    </button>
    <div
      v-if="open"
      class="absolute right-0 top-full z-10 mt-1 rounded-lg border border-slate-700 bg-slate-800 py-1 shadow-xl"
    >
      <button
        v-for="fmt in formats"
        :key="fmt.value"
        class="block w-full px-4 py-2 text-left text-sm text-slate-200 hover:bg-slate-700"
        @click="download(fmt.value)"
      >
        {{ fmt.label }}
      </button>
    </div>
    <p v-if="error" class="mt-1 text-xs text-red-400">{{ error }}</p>
  </div>
</template>
