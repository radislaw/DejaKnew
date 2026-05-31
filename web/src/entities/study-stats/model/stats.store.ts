import { ref } from 'vue'
import { defineStore } from 'pinia'
import { fetchStats } from '../api/stats.api'
import type { StudyStats } from './types'

export const useStatsStore = defineStore('study-stats', () => {
  const stats = ref<StudyStats | null>(null)
  const loading = ref(false)
  const error = ref<string | null>(null)

  async function fetchStudyStats() {
    loading.value = true
    error.value = null
    try {
      stats.value = await fetchStats()
    } catch (e) {
      error.value = e instanceof Error ? e.message : 'Failed to load stats'
    } finally {
      loading.value = false
    }
  }

  return { stats, loading, error, fetchStudyStats }
})
