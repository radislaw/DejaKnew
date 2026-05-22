package domain

import "time"

// Deck is a named collection of cards owned by a single user.
type Deck struct {
	ID          string    `gorm:"type:uuid;primaryKey" json:"id"`
	OwnerID     string    `gorm:"type:uuid;index;not null" json:"ownerId"`
	Title       string    `gorm:"not null" json:"title"`
	Description string    `gorm:"not null;default:''" json:"description"`
	Tags        []string  `gorm:"serializer:json" json:"tags"`
	IsPublic    bool      `gorm:"not null;default:false" json:"isPublic"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}
