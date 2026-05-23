import { apiRequest } from '@/shared/api'
import type { CreateDeckInput, Deck } from '../model/types'

export function listDecks(): Promise<Deck[]> {
  return apiRequest<Deck[]>('/decks')
}

export function createDeck(input: CreateDeckInput): Promise<Deck> {
  return apiRequest<Deck>('/decks', { method: 'POST', body: input })
}

export function deleteDeck(id: string): Promise<void> {
  return apiRequest<void>(`/decks/${id}`, { method: 'DELETE' })
}
