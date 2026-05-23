export interface Deck {
  id: string
  ownerId: string
  title: string
  description: string
  tags: string[]
  isPublic: boolean
  createdAt: string
  updatedAt: string
}

export interface CreateDeckInput {
  title: string
  description?: string
  tags?: string[]
  isPublic?: boolean
}

// Inline card shape to avoid cross-entity imports (FSD @x — v2 task)
export interface DeckDetailCard {
  id: string
  deckId: string
  type: string
  content: { front: string; back: string; frontImage: string | null; backImage: string | null }
  position: number
  createdAt: string
  updatedAt: string
}

export interface DeckDetail extends Deck {
  cards: DeckDetailCard[]
}
