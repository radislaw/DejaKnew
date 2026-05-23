package domain

import "time"

// CardContent holds the v1 Basic card fields stored as JSONB.
type CardContent struct {
	Front      string  `json:"front"`
	Back       string  `json:"back"`
	FrontImage *string `json:"frontImage"`
	BackImage  *string `json:"backImage"`
}

// Card is a single flashcard belonging to a Deck.
type Card struct {
	ID        string      `gorm:"type:uuid;primaryKey" json:"id"`
	DeckID    string      `gorm:"type:uuid;index;not null" json:"deckId"`
	Type      string      `gorm:"not null;default:'basic'" json:"type"`
	Content   CardContent `gorm:"serializer:json" json:"content"`
	Position  int         `gorm:"not null;default:0" json:"position"`
	CreatedAt time.Time   `json:"createdAt"`
	UpdatedAt time.Time   `json:"updatedAt"`
}
