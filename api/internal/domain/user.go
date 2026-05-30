package domain

import "time"

// User is a registered account. PasswordHash is empty for OAuth users.
type User struct {
	ID           string    `gorm:"type:uuid;primaryKey" json:"id"`
	Email        string    `gorm:"uniqueIndex;not null" json:"email"`
	PasswordHash string    `json:"-"`
	Provider     string    `gorm:"not null;default:'local'" json:"provider"`
	ProviderID   string    `json:"-"`
	Locale       string    `gorm:"not null;default:'en'" json:"locale"`
	CreatedAt    time.Time `json:"createdAt"`
}
