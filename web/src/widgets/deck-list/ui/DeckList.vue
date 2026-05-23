<script setup lang="ts">
import type { Deck } from '@/entities/deck'

defineProps<{ decks: Deck[] }>()
const emit = defineEmits<{ delete: [id: string] }>()
</script>

<template>
  <ul class="grid gap-3 sm:grid-cols-2 lg:grid-cols-3">
    <li
      v-for="deck in decks"
      :key="deck.id"
      class="group flex flex-col gap-2 rounded-xl border border-slate-700 bg-slate-800/60 p-4 transition hover:border-slate-500"
    >
      <div class="flex items-start justify-between gap-2">
        <span class="text-sm font-medium text-slate-100 leading-snug">{{ deck.title }}</span>
        <button
          class="shrink-0 text-slate-600 opacity-0 transition hover:text-red-400 group-hover:opacity-100"
          title="Delete deck"
          @click="emit('delete', deck.id)"
        >
          ✕
        </button>
      </div>
      <p v-if="deck.description" class="text-xs text-slate-400 line-clamp-2">
        {{ deck.description }}
      </p>
      <div v-if="deck.tags?.length" class="flex flex-wrap gap-1">
        <span
          v-for="tag in deck.tags"
          :key="tag"
          class="rounded-full bg-slate-700 px-2 py-0.5 text-xs text-slate-300"
        >
          {{ tag }}
        </span>
      </div>
    </li>
  </ul>
</template>
