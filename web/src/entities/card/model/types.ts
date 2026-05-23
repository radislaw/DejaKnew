export interface CardContent {
  front: string
  back: string
  frontImage: string | null
  backImage: string | null
  topic?: string
  level?: string
  sourceUrl?: string
  tags?: string[]
}

export interface Card {
  id: string
  deckId: string
  type: string
  content: CardContent
  position: number
  createdAt: string
  updatedAt: string
}

export interface CreateCardInput {
  front: string
  back: string
  topic?: string
  level?: string
  sourceUrl?: string
  tags?: string[]
}
