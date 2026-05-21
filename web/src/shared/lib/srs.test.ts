import { describe, it, expect } from 'vitest'
import { calculateNextReview, getCardStatus } from './srs'

describe('calculateNextReview', () => {
  it('first successful review sets interval to 1 day', () => {
    const s = calculateNextReview(undefined, 2)
    expect(s.reps).toBe(1)
    expect(s.interval).toBe(1)
  })

  it('quality 0 resets reps and counts a lapse', () => {
    const prev = calculateNextReview(undefined, 2)
    const s = calculateNextReview(prev, 0)
    expect(s.reps).toBe(0)
    expect(s.lapses).toBe(1)
    expect(s.interval).toBe(0)
  })

  it('ease never drops below 1.3', () => {
    let s = calculateNextReview(undefined, 0)
    for (let i = 0; i < 20; i++) s = calculateNextReview(s, 0)
    expect(s.ease).toBeGreaterThanOrEqual(1.3)
  })
})

describe('getCardStatus', () => {
  it('returns "new" when there is no srs state', () => {
    expect(getCardStatus(undefined)).toBe('new')
  })

  it('returns "due" when nextReview is in the past', () => {
    const s = calculateNextReview(undefined, 2)
    expect(getCardStatus({ ...s, nextReview: Date.now() - 1000 })).toBe('due')
  })
})
