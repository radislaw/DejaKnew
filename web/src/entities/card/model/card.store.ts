import { defineStore } from 'pinia'
import { ref } from 'vue'
import { createCard, deleteCard, updateCard } from '../api/card.api'
import type { Card, CreateCardInput } from './types'

export const useCardStore = defineStore('card', () => {
  const cards = ref<Card[]>([])
  const loading = ref(false)
  const error = ref<string | null>(null)

  function setCards(c: Card[]): void {
    cards.value = c
  }

  async function addCard(deckId: string, input: CreateCardInput): Promise<void> {
    const card = await createCard(deckId, input)
    cards.value = [...cards.value, card]
  }

  async function updateCardItem(cardId: string, input: CreateCardInput): Promise<void> {
    const updated = await updateCard(cardId, input)
    cards.value = cards.value.map((c) => (c.id === cardId ? updated : c))
  }

  async function removeCard(cardId: string): Promise<void> {
    await deleteCard(cardId)
    cards.value = cards.value.filter((c) => c.id !== cardId)
  }

  return { cards, loading, error, setCards, addCard, updateCardItem, removeCard }
})
