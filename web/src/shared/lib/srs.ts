export interface SrsState {
  interval: number
  ease: number
  reps: number
  lapses: number
  nextReview: number
  lastReview: number
}

export type CardStatus = 'new' | 'learning' | 'due' | 'learned' | 'difficult'

// Pure client-side derivation of a card's status from its SRS state.
// Next-interval math (calculateNextReview) is authoritative on the backend.
export function getCardStatus(srs: SrsState | undefined, now = Date.now()): CardStatus {
  if (!srs || srs.reps === 0) return 'new'
  if (srs.nextReview <= now) return 'due'
  if (srs.lapses >= 2) return 'difficult'
  if (srs.reps >= 5) return 'learned'
  return 'learning'
}
