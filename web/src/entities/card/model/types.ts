export interface CardContent {
  front: string
  back: string
  frontImage: string | null
  backImage: string | null
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
}
