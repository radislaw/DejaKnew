import { defineStore } from 'pinia'
import { ref } from 'vue'
import { createDeck, deleteDeck, listDecks } from '../api/deck.api'
import type { CreateDeckInput, Deck } from './types'

export const useDeckStore = defineStore('deck', () => {
  const decks = ref<Deck[]>([])
  const loading = ref(false)
  const error = ref<string | null>(null)

  async function fetchDecks(): Promise<void> {
    loading.value = true
    error.value = null
    try {
      decks.value = await listDecks()
    } catch (e) {
      error.value = e instanceof Error ? e.message : 'Failed to load decks'
    } finally {
      loading.value = false
    }
  }

  async function addDeck(input: CreateDeckInput): Promise<void> {
    const deck = await createDeck(input)
    decks.value = [deck, ...decks.value]
  }

  async function removeDeck(id: string): Promise<void> {
    await deleteDeck(id)
    decks.value = decks.value.filter((d) => d.id !== id)
  }

  return { decks, loading, error, fetchDecks, addDeck, removeDeck }
})
