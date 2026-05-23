<script setup lang="ts">
import { ref } from 'vue'
import { importCards } from '@/entities/card/api/card.api'
import type { Card } from '@/entities/card'

const props = defineProps<{ deckId: string }>()

const emit = defineEmits<{
  imported: [cards: Card[]]
  cancel: []
}>()

const format = ref<'json' | 'csv'>('json')
const fileText = ref('')
const fileName = ref('')
const preview = ref<number | null>(null)
const error = ref<string | null>(null)
const loading = ref(false)

function onFileChange(e: Event) {
  const input = e.target as HTMLInputElement
  const file = input.files?.[0]
  if (!file) return

  fileName.value = file.name
  error.value = null
  preview.value = null

  if (file.name.endsWith('.csv')) format.value = 'csv'
  else if (file.name.endsWith('.json')) format.value = 'json'

  const reader = new FileReader()
  reader.onload = (ev) => {
    fileText.value = (ev.target?.result as string) ?? ''
    previewCount()
  }
  reader.readAsText(file)
}

function previewCount() {
  if (!fileText.value) return
  try {
    if (format.value === 'json') {
      const parsed = JSON.parse(fileText.value)
      preview.value = Array.isArray(parsed) ? parsed.filter((r) => r.front && r.back).length : 0
    } else {
      const lines = fileText.value.split('\n').filter((l) => l.trim())
      preview.value = Math.max(0, lines.length - 1)
    }
    error.value = null
  } catch {
    error.value = 'Failed to parse file'
    preview.value = null
  }
}

async function onSubmit() {
  if (!fileText.value) return
  error.value = null
  loading.value = true
  try {
    const result = await importCards(props.deckId, format.value, fileText.value)
    emit('imported', result.cards)
  } catch (e) {
    error.value = e instanceof Error ? e.message : 'Import failed'
  } finally {
    loading.value = false
  }
}
</script>

<template>
  <div class="rounded-xl border border-slate-700 bg-slate-800/60 p-4 flex flex-col gap-3">
    <div class="flex items-center justify-between">
      <span class="text-sm font-medium text-slate-200">Import cards</span>
      <div class="flex gap-1 rounded-lg bg-slate-900 p-0.5">
        <button
          v-for="f in ['json', 'csv'] as const"
          :key="f"
          class="rounded-md px-3 py-1 text-xs font-medium transition-colors"
          :class="format === f ? 'bg-sky-600 text-white' : 'text-slate-400 hover:text-slate-200'"
          @click="format = f; previewCount()"
        >
          {{ f.toUpperCase() }}
        </button>
      </div>
    </div>

    <label
      class="flex cursor-pointer flex-col items-center gap-2 rounded-lg border border-dashed border-slate-600 px-4 py-6 text-sm text-slate-400 hover:border-slate-500 hover:text-slate-300 transition-colors"
    >
      <span class="text-2xl">📂</span>
      <span>{{ fileName || 'Choose file (.json or .csv)' }}</span>
      <input type="file" accept=".json,.csv" class="sr-only" @change="onFileChange" />
    </label>

    <div v-if="preview !== null" class="text-xs text-slate-400">
      Preview: <span class="text-slate-200 font-medium">{{ preview }} card{{ preview !== 1 ? 's' : '' }}</span> found
    </div>

    <p v-if="error" class="text-xs text-red-400">{{ error }}</p>

    <div class="text-xs text-slate-600">
      <template v-if="format === 'csv'">
        Columns: <code>front, back, topic, level, sourceUrl, tags</code> — tags separated by <code>|</code>
      </template>
      <template v-else>
        JSON array: <code>[{"{"}"front":"...","back":"..."{"}"}, ...]</code>
      </template>
    </div>

    <div class="flex gap-2">
      <button
        :disabled="!fileText || loading || preview === 0"
        class="rounded-lg bg-emerald-600 px-4 py-1.5 text-sm font-medium text-white hover:bg-emerald-500 disabled:opacity-40 disabled:cursor-not-allowed"
        @click="onSubmit"
      >
        {{ loading ? 'Importing…' : 'Import' }}
      </button>
      <button
        class="rounded-lg px-3 py-1.5 text-sm text-slate-400 hover:text-slate-200"
        @click="emit('cancel')"
      >
        Cancel
      </button>
    </div>
  </div>
</template>
