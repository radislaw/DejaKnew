<script setup lang="ts">
import { ref } from 'vue'
import type { QueueCard } from '@/entities/study-progress'

defineProps<{ card: QueueCard }>()

const emit = defineEmits<{
  rate: [quality: 0 | 1 | 2 | 3]
}>()

const flipped = ref(false)

function flip() {
  flipped.value = true
}
</script>

<template>
  <div class="flex flex-col items-center gap-6">
    <div
      class="relative h-64 w-full max-w-xl cursor-pointer select-none rounded-2xl border border-slate-700 bg-slate-800/60 p-8"
      style="perspective: 1200px"
      @click="!flipped && flip()"
    >
      <div
        class="absolute inset-0 flex items-center justify-center rounded-2xl p-8 transition-all duration-500"
        :class="flipped ? 'opacity-0 scale-95 pointer-events-none' : 'opacity-100'"
      >
        <p class="text-center text-xl font-medium text-slate-100">{{ card.content.front }}</p>
      </div>

      <div
        class="absolute inset-0 flex items-center justify-center rounded-2xl p-8 transition-all duration-500"
        :class="flipped ? 'opacity-100' : 'opacity-0 scale-95 pointer-events-none'"
      >
        <p class="text-center text-xl text-slate-200">{{ card.content.back }}</p>
      </div>

      <span
        v-if="!flipped"
        class="absolute bottom-3 left-1/2 -translate-x-1/2 text-xs text-slate-500"
      >tap to reveal</span>
    </div>

    <div v-if="flipped" class="grid w-full max-w-xl grid-cols-4 gap-2">
      <button
        class="rounded-xl border border-slate-600 bg-slate-800 py-3 text-sm font-medium text-red-400 hover:bg-red-900/30"
        @click="emit('rate', 0)"
      >
        Again
      </button>
      <button
        class="rounded-xl border border-slate-600 bg-slate-800 py-3 text-sm font-medium text-orange-400 hover:bg-orange-900/30"
        @click="emit('rate', 1)"
      >
        Hard
      </button>
      <button
        class="rounded-xl border border-slate-600 bg-slate-800 py-3 text-sm font-medium text-sky-400 hover:bg-sky-900/30"
        @click="emit('rate', 2)"
      >
        Good
      </button>
      <button
        class="rounded-xl border border-slate-600 bg-slate-800 py-3 text-sm font-medium text-emerald-400 hover:bg-emerald-900/30"
        @click="emit('rate', 3)"
      >
        Easy
      </button>
    </div>
  </div>
</template>
