import { apiRequest } from '@/shared/api'
import type { QueueCard, StudyProgress } from '../model/types'

export function getStudyQueue(deckId: string): Promise<QueueCard[]> {
  return apiRequest<QueueCard[]>(`/study/${deckId}/queue`)
}

export function submitReview(cardId: string, quality: 0 | 1 | 2 | 3): Promise<StudyProgress> {
  return apiRequest<StudyProgress>(`/study/${cardId}/review`, {
    method: 'POST',
    body: { quality },
  })
}
