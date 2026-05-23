import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { getStudyQueue, submitReview } from '../api/study.api'
import type { QueueCard, StudyProgress } from './types'

export const useStudyStore = defineStore('study', () => {
  const queue = ref<QueueCard[]>([])
  const currentIndex = ref(0)
  const lastProgress = ref<StudyProgress | null>(null)
  const loading = ref(false)
  const error = ref<string | null>(null)

  const currentCard = computed<QueueCard | null>(() => queue.value[currentIndex.value] ?? null)
  const isDone = computed(() => currentIndex.value >= queue.value.length)
  const remaining = computed(() => queue.value.length - currentIndex.value)

  async function loadQueue(deckId: string): Promise<void> {
    loading.value = true
    error.value = null
    try {
      queue.value = await getStudyQueue(deckId)
      currentIndex.value = 0
      lastProgress.value = null
    } catch (e) {
      error.value = e instanceof Error ? e.message : 'Failed to load queue'
    } finally {
      loading.value = false
    }
  }

  async function review(quality: 0 | 1 | 2 | 3): Promise<void> {
    const card = currentCard.value
    if (!card) return
    lastProgress.value = await submitReview(card.id, quality)
    currentIndex.value += 1
  }

  function reset(): void {
    queue.value = []
    currentIndex.value = 0
    lastProgress.value = null
    error.value = null
  }

  return { queue, currentIndex, lastProgress, loading, error, currentCard, isDone, remaining, loadQueue, review, reset }
})
