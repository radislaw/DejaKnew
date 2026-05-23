// Inline card shape to avoid cross-entity imports (FSD @x — v2 task)
export interface QueueCard {
  id: string
  deckId: string
  content: { front: string; back: string; frontImage: string | null; backImage: string | null; topic?: string; level?: string; sourceUrl?: string; tags?: string[] }
  position: number
}

export interface StudyProgress {
  id: string
  userId: string
  cardId: string
  interval: number
  ease: number
  reps: number
  lapses: number
  nextReview: string | null
  lastReview: string | null
}
