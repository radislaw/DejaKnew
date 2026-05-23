package domain

import "time"

// StudyProgress holds per-user per-card SRS state.
type StudyProgress struct {
	ID         string     `gorm:"type:uuid;primaryKey" json:"id"`
	UserID     string     `gorm:"type:uuid;uniqueIndex:idx_progress_user_card;not null" json:"userId"`
	CardID     string     `gorm:"type:uuid;uniqueIndex:idx_progress_user_card;not null" json:"cardId"`
	Interval   int        `gorm:"not null;default:0" json:"interval"`
	Ease       float64    `gorm:"not null;default:2.5" json:"ease"`
	Reps       int        `gorm:"not null;default:0" json:"reps"`
	Lapses     int        `gorm:"not null;default:0" json:"lapses"`
	NextReview *time.Time `json:"nextReview"`
	LastReview *time.Time `json:"lastReview"`
}

// ReviewLog records each individual review event for progress graphs.
type ReviewLog struct {
	ID         string    `gorm:"type:uuid;primaryKey" json:"id"`
	UserID     string    `gorm:"type:uuid;index;not null" json:"userId"`
	CardID     string    `gorm:"type:uuid;index;not null" json:"cardId"`
	Quality    int       `gorm:"not null" json:"quality"`
	ReviewedAt time.Time `gorm:"not null;default:now()" json:"reviewedAt"`
}
