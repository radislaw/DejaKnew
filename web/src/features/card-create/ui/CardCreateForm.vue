<script setup lang="ts">
import { ref } from 'vue'
import type { CreateCardInput } from '@/entities/card'

const emit = defineEmits<{
  submit: [input: CreateCardInput]
  cancel: []
}>()

const front = ref('')
const back = ref('')
const topic = ref('')
const level = ref('')
const sourceUrl = ref('')
const tagsRaw = ref('')
const showOptional = ref(false)

function onSubmit() {
  const f = front.value.trim()
  const b = back.value.trim()
  if (!f || !b) return

  const input: CreateCardInput = { front: f, back: b }
  if (topic.value.trim()) input.topic = topic.value.trim()
  if (level.value.trim()) input.level = level.value.trim()
  if (sourceUrl.value.trim()) input.sourceUrl = sourceUrl.value.trim()
  const tags = tagsRaw.value.split(',').map(t => t.trim()).filter(Boolean)
  if (tags.length) input.tags = tags

  emit('submit', input)
  front.value = ''
  back.value = ''
  topic.value = ''
  level.value = ''
  sourceUrl.value = ''
  tagsRaw.value = ''
  showOptional.value = false
}
</script>

<template>
  <form
    class="flex flex-col gap-3 rounded-xl border border-slate-700 bg-slate-800/60 p-4"
    @submit.prevent="onSubmit"
  >
    <h2 class="text-sm font-medium text-slate-300">New card</h2>
    <textarea
      v-model="front"
      placeholder="Front"
      rows="2"
      required
      class="resize-none rounded-lg border border-slate-600 bg-slate-900 px-3 py-2 text-sm text-slate-100 placeholder-slate-500 outline-none focus:border-sky-500"
    />
    <textarea
      v-model="back"
      placeholder="Back (markdown supported)"
      rows="3"
      required
      class="resize-none rounded-lg border border-slate-600 bg-slate-900 px-3 py-2 text-sm text-slate-100 placeholder-slate-500 outline-none focus:border-sky-500"
    />

    <button
      type="button"
      class="self-start text-xs text-slate-500 hover:text-slate-300 transition-colors"
      @click="showOptional = !showOptional"
    >
      {{ showOptional ? '▲ Less' : '▼ More options' }}
    </button>

    <template v-if="showOptional">
      <div class="grid grid-cols-2 gap-2">
        <input
          v-model="topic"
          placeholder="Topic (e.g. Hoisting)"
          class="rounded-lg border border-slate-600 bg-slate-900 px-3 py-2 text-sm text-slate-100 placeholder-slate-500 outline-none focus:border-sky-500"
        />
        <input
          v-model="level"
          placeholder="Level (e.g. Уровень 1)"
          class="rounded-lg border border-slate-600 bg-slate-900 px-3 py-2 text-sm text-slate-100 placeholder-slate-500 outline-none focus:border-sky-500"
        />
      </div>
      <input
        v-model="sourceUrl"
        placeholder="Source URL (optional)"
        type="url"
        class="rounded-lg border border-slate-600 bg-slate-900 px-3 py-2 text-sm text-slate-100 placeholder-slate-500 outline-none focus:border-sky-500"
      />
      <input
        v-model="tagsRaw"
        placeholder="Tags, comma-separated (e.g. js, hoisting, let)"
        class="rounded-lg border border-slate-600 bg-slate-900 px-3 py-2 text-sm text-slate-100 placeholder-slate-500 outline-none focus:border-sky-500"
      />
    </template>

    <div class="flex gap-2">
      <button
        type="submit"
        class="rounded-lg bg-sky-600 px-4 py-1.5 text-sm font-medium text-white hover:bg-sky-500"
      >
        Add card
      </button>
      <button
        type="button"
        class="rounded-lg px-4 py-1.5 text-sm text-slate-400 hover:text-slate-200"
        @click="emit('cancel')"
      >
        Cancel
      </button>
    </div>
  </form>
</template>
