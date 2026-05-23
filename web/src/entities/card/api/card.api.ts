import { apiRequest } from '@/shared/api'
import type { Card, CreateCardInput } from '../model/types'

export function createCard(deckId: string, input: CreateCardInput): Promise<Card> {
  return apiRequest<Card>(`/decks/${deckId}/cards`, { method: 'POST', body: input })
}

export function updateCard(id: string, input: CreateCardInput): Promise<Card> {
  return apiRequest<Card>(`/cards/${id}`, { method: 'PATCH', body: input })
}

export function deleteCard(id: string): Promise<void> {
  return apiRequest<void>(`/cards/${id}`, { method: 'DELETE' })
}
