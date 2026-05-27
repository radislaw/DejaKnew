<script setup lang="ts">
import { ref, computed } from 'vue'
import { useI18n } from 'vue-i18n'
import { marked } from 'marked'
import type { QueueCard } from '@/entities/study-progress'

const { t } = useI18n()

const props = defineProps<{ card: QueueCard }>()

const emit = defineEmits<{
  rate: [quality: 0 | 1 | 2 | 3]
}>()

const flipped = ref(false)

function flip() {
  flipped.value = true
}

const backHtml = computed(() => {
  const result = marked.parse(props.card.content.back)
  return typeof result === 'string' ? result : ''
})

const tags = computed(() => props.card.content.tags ?? [])
</script>

<template>
  <div class="flex flex-col items-center gap-6">
    <div
      class="relative w-full max-w-xl cursor-pointer select-none rounded-2xl border border-slate-700 bg-slate-800/60"
      style="perspective: 1200px; min-height: 16rem"
      @click="!flipped && flip()"
    >
      <!-- Front -->
      <div
        class="absolute inset-0 flex flex-col items-center justify-center rounded-2xl p-8 transition-all duration-500"
        :class="flipped ? 'opacity-0 scale-95 pointer-events-none' : 'opacity-100'"
      >
        <div v-if="card.content.level || card.content.topic" class="mb-3 flex gap-2 self-start">
          <span
            v-if="card.content.level"
            class="rounded-full bg-sky-900/60 px-2.5 py-0.5 text-xs font-medium text-sky-300"
          >{{ card.content.level }}</span>
          <span
            v-if="card.content.topic"
            class="rounded-full bg-slate-700/80 px-2.5 py-0.5 text-xs text-slate-400"
          >{{ card.content.topic }}</span>
        </div>
        <p class="text-center text-xl font-medium text-slate-100">{{ card.content.front }}</p>
        <span class="mt-6 text-xs text-slate-500">{{ t('cardReview.tapToReveal') }}</span>
      </div>

      <!-- Back -->
      <div
        class="absolute inset-0 flex flex-col rounded-2xl p-8 transition-all duration-500 overflow-auto"
        :class="flipped ? 'opacity-100' : 'opacity-0 scale-95 pointer-events-none'"
      >
        <div
          class="prose prose-invert prose-sm max-w-none flex-1"
          v-html="backHtml"
        />

        <div v-if="tags.length" class="mt-4 flex flex-wrap gap-1.5">
          <span
            v-for="tag in tags"
            :key="tag"
            class="rounded-full bg-slate-700/60 px-2 py-0.5 text-xs text-slate-400"
          >#{{ tag }}</span>
        </div>

        <a
          v-if="card.content.sourceUrl"
          :href="card.content.sourceUrl"
          target="_blank"
          rel="noopener noreferrer"
          class="mt-3 self-start text-xs text-sky-400 hover:text-sky-300 hover:underline"
          @click.stop
        >
          {{ t('cardReview.source') }}
        </a>
      </div>
    </div>

    <div v-if="flipped" class="grid w-full max-w-xl grid-cols-4 gap-2">
      <button
        class="rounded-xl border border-slate-600 bg-slate-800 py-3 text-sm font-medium text-red-400 hover:bg-red-900/30"
        @click="emit('rate', 0)"
      >
        {{ t('cardReview.again') }}
      </button>
      <button
        class="rounded-xl border border-slate-600 bg-slate-800 py-3 text-sm font-medium text-orange-400 hover:bg-orange-900/30"
        @click="emit('rate', 1)"
      >
        {{ t('cardReview.hard') }}
      </button>
      <button
        class="rounded-xl border border-slate-600 bg-slate-800 py-3 text-sm font-medium text-sky-400 hover:bg-sky-900/30"
        @click="emit('rate', 2)"
      >
        {{ t('cardReview.good') }}
      </button>
      <button
        class="rounded-xl border border-slate-600 bg-slate-800 py-3 text-sm font-medium text-emerald-400 hover:bg-emerald-900/30"
        @click="emit('rate', 3)"
      >
        {{ t('cardReview.easy') }}
      </button>
    </div>
  </div>
</template>
