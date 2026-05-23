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
