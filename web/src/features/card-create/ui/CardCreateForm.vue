<script setup lang="ts">
import { ref } from 'vue'

const emit = defineEmits<{
  submit: [front: string, back: string]
  cancel: []
}>()

const front = ref('')
const back = ref('')

function onSubmit() {
  const f = front.value.trim()
  const b = back.value.trim()
  if (!f || !b) return
  emit('submit', f, b)
  front.value = ''
  back.value = ''
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
      placeholder="Back"
      rows="2"
      required
      class="resize-none rounded-lg border border-slate-600 bg-slate-900 px-3 py-2 text-sm text-slate-100 placeholder-slate-500 outline-none focus:border-sky-500"
    />
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
