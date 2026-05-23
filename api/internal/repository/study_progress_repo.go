package repository

import (
	"errors"
	"time"

	"flashcards/api/internal/domain"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// StudyProgressRepo is the GORM-backed implementation of service.StudyProgressRepository.
type StudyProgressRepo struct {
	db *gorm.DB
}

func NewStudyProgressRepo(db *gorm.DB) *StudyProgressRepo {
	return &StudyProgressRepo{db: db}
}

// GetDueQueue returns cards in a deck that are new or due for review.
func (r *StudyProgressRepo) GetDueQueue(userID, deckID string) ([]domain.Card, error) {
	var cards []domain.Card
	now := time.Now()
	err := r.db.
		Where("deck_id = ?", deckID).
		Where(
			"id NOT IN (SELECT card_id FROM study_progresses WHERE user_id = ? AND next_review > ?)",
			userID, now,
		).
		Order("position asc").
		Find(&cards).Error
	if err != nil {
		return nil, err
	}
	return cards, nil
}

// GetProgress returns the SRS state for a user+card pair, or nil if not started.
func (r *StudyProgressRepo) GetProgress(userID, cardID string) (*domain.StudyProgress, error) {
	var sp domain.StudyProgress
	err := r.db.First(&sp, "user_id = ? AND card_id = ?", userID, cardID).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &sp, nil
}

// UpsertProgress creates or updates the SRS state for a user+card pair.
func (r *StudyProgressRepo) UpsertProgress(sp *domain.StudyProgress) error {
	return r.db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "user_id"}, {Name: "card_id"}},
		DoUpdates: clause.AssignmentColumns([]string{"interval", "ease", "reps", "lapses", "next_review", "last_review"}),
	}).Create(sp).Error
}

// CreateReviewLog appends one review event to the log.
func (r *StudyProgressRepo) CreateReviewLog(entry *domain.ReviewLog) error {
	return r.db.Create(entry).Error
}
