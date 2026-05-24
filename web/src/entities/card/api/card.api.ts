import { apiRequest, apiBlobRequest } from '@/shared/api'
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

export interface ImportResult {
  imported: number
  cards: Card[]
}

export function importCards(
  deckId: string,
  format: 'json' | 'csv',
  data: string,
): Promise<ImportResult> {
  return apiRequest<ImportResult>(`/decks/${deckId}/import`, {
    method: 'POST',
    body: { format, data },
  })
}

export type ExportFormat = 'csv' | 'json' | 'anki-csv'

export function exportCards(format: ExportFormat, deckId?: string): Promise<{ blob: Blob; filename: string }> {
  const path = deckId
    ? `/decks/${deckId}/export?format=${format}`
    : `/export?format=${format}`
  return apiBlobRequest(path)
}
