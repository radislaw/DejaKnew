import { apiRequest } from '@/shared/api/client'
import type { StudyStats } from '../model/types'

export function fetchStats(): Promise<StudyStats> {
  return apiRequest<StudyStats>('/stats')
}
