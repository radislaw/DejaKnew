export interface SrsState {
  interval: number
  ease: number
  reps: number
  lapses: number
  nextReview: number
  lastReview: number
}

export type Quality = 0 | 1 | 2 | 3

export type CardStatus = 'new' | 'learning' | 'due' | 'learned' | 'difficult'

const DAY_MS = 24 * 60 * 60 * 1000

export function calculateNextReview(
  prev: Partial<SrsState> | undefined,
  quality: Quality,
): SrsState {
  const now = Date.now()
  let interval = prev?.interval ?? 0
  let ease = prev?.ease ?? 2.5
  let reps = prev?.reps ?? 0
  let lapses = prev?.lapses ?? 0

  if (quality === 0) {
    reps = 0
    interval = 0
    lapses += 1
    ease = Math.max(1.3, ease - 0.2)
  } else {
    reps += 1
    ease = Math.max(1.3, ease + (0.1 - (3 - quality) * (0.08 + (3 - quality) * 0.02)))
    if (reps === 1) interval = 1
    else if (reps === 2) interval = quality === 1 ? 1 : 3
    else interval = Math.round(interval * ease * (quality === 1 ? 0.6 : quality === 3 ? 1.3 : 1))
  }

  return {
    interval,
    ease,
    reps,
    lapses,
    nextReview: now + interval * DAY_MS,
    lastReview: now,
  }
}

export function getCardStatus(srs: SrsState | undefined, now = Date.now()): CardStatus {
  if (!srs || srs.reps === 0) return 'new'
  if (srs.nextReview <= now) return 'due'
  if (srs.lapses >= 2) return 'difficult'
  if (srs.reps >= 5) return 'learned'
  return 'learning'
}
