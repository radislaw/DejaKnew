import { describe, it, expect } from 'vitest'
import { getCardStatus, type SrsState } from './srs'

const srs = (over: Partial<SrsState> = {}): SrsState => ({
  interval: 1,
  ease: 2.5,
  reps: 1,
  lapses: 0,
  nextReview: Date.now() + 86_400_000,
  lastReview: Date.now(),
  ...over,
})

describe('getCardStatus', () => {
  it('returns "new" when there is no srs state', () => {
    expect(getCardStatus(undefined)).toBe('new')
  })

  it('returns "new" when reps is 0', () => {
    expect(getCardStatus(srs({ reps: 0 }))).toBe('new')
  })

  it('returns "due" when nextReview is in the past', () => {
    expect(getCardStatus(srs({ nextReview: Date.now() - 1000 }))).toBe('due')
  })

  it('returns "difficult" after 2 or more lapses', () => {
    expect(getCardStatus(srs({ lapses: 2 }))).toBe('difficult')
  })

  it('returns "learned" after 5 or more reps', () => {
    expect(getCardStatus(srs({ reps: 5 }))).toBe('learned')
  })
})
