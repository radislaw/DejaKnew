<script setup lang="ts">
import { onMounted, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useI18n } from 'vue-i18n'
import { useStudyStore } from '@/entities/study-progress'
import { CardReview } from '@/features/card-review'

const { t } = useI18n()
const route = useRoute()
const router = useRouter()
const deckId = route.params.id as string

const study = useStudyStore()

onMounted(() => study.loadQueue(deckId))

watch(() => route.params.id, (id) => {
  study.reset()
  study.loadQueue(id as string)
})

async function onRate(quality: 0 | 1 | 2 | 3) {
  try {
    await study.review(quality)
  } catch {
    // error surfaced via store if needed
  }
}
</script>

<template>
  <div class="min-h-screen bg-[#070d1c] px-4 py-6 text-slate-100">
    <header class="mx-auto mb-8 flex max-w-xl items-center gap-4">
      <button
        class="text-sm text-slate-400 hover:text-slate-200"
        @click="router.push({ name: 'deck', params: { id: deckId } })"
      >
        {{ t('common.back') }}
      </button>
      <div class="flex-1">
        <h1 class="text-lg font-semibold tracking-tight">{{ t('study.title') }}</h1>
      </div>
      <span class="text-sm text-slate-500">{{ study.remaining }} {{ t('study.remaining') }}</span>
    </header>

    <main class="mx-auto max-w-xl">
      <p v-if="study.loading" class="text-center text-sm text-slate-400">{{ t('common.loading') }}</p>

      <p v-else-if="study.error" class="text-center text-sm text-red-400">{{ study.error }}</p>

      <div
        v-else-if="study.isDone"
        class="flex flex-col items-center gap-4 py-20 text-center"
      >
        <p class="text-2xl font-semibold text-emerald-400">{{ t('study.done') }}</p>
        <p class="text-sm text-slate-400">{{ t('study.doneMessage') }}</p>
        <button
          class="mt-4 rounded-lg bg-sky-600 px-6 py-2 text-sm font-medium text-white hover:bg-sky-500"
          @click="router.push({ name: 'deck', params: { id: deckId } })"
        >
          {{ t('study.backToDeck') }}
        </button>
      </div>

      <CardReview
        v-else-if="study.currentCard"
        :key="study.currentCard.id"
        :card="study.currentCard"
        @rate="onRate"
      />
    </main>
  </div>
</template>
