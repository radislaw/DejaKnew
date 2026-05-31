package repository

import (
	"errors"
	"time"

	"flashcards/api/internal/domain"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// GetUserStats computes streak and review counts from ReviewLog.
func (r *StudyProgressRepo) GetUserStats(userID string) (*domain.StudyStats, error) {
	var totalReviews int64
	if err := r.db.Model(&domain.ReviewLog{}).Where("user_id = ?", userID).Count(&totalReviews).Error; err != nil {
		return nil, err
	}

	var cardsStudied int64
	if err := r.db.Model(&domain.ReviewLog{}).Where("user_id = ?", userID).
		Distinct("card_id").Count(&cardsStudied).Error; err != nil {
		return nil, err
	}

	var todayReviews int64
	if err := r.db.Model(&domain.ReviewLog{}).
		Where("user_id = ? AND DATE(reviewed_at) = CURRENT_DATE", userID).
		Count(&todayReviews).Error; err != nil {
		return nil, err
	}

	// Distinct review dates ordered newest-first for streak calculation.
	var dates []time.Time
	if err := r.db.Raw(
		"SELECT DATE(reviewed_at) FROM review_logs WHERE user_id = ? GROUP BY DATE(reviewed_at) ORDER BY DATE(reviewed_at) DESC",
		userID,
	).Scan(&dates).Error; err != nil {
		return nil, err
	}

	current, longest := calculateStreaks(dates)

	lastDate := ""
	if len(dates) > 0 {
		lastDate = dates[0].Format("2006-01-02")
	}

	return &domain.StudyStats{
		Streak: domain.StreakStats{Current: current, Longest: longest, LastDate: lastDate},
		Today:  domain.DayStats{Reviews: int(todayReviews)},
		Total:  domain.TotalStats{Reviews: int(totalReviews), CardsStudied: int(cardsStudied)},
	}, nil
}

// calculateStreaks derives current and longest streaks from a DESC-ordered list of distinct review dates.
func calculateStreaks(dates []time.Time) (current, longest int) {
	if len(dates) == 0 {
		return 0, 0
	}

	today := time.Now().UTC().Truncate(24 * time.Hour)
	yesterday := today.AddDate(0, 0, -1)

	// Current streak: walk back from the most recent date only if it's today or yesterday.
	if dates[0].Equal(today) || dates[0].Equal(yesterday) {
		current = 1
		for i := 1; i < len(dates); i++ {
			if dates[i].Equal(dates[i-1].AddDate(0, 0, -1)) {
				current++
			} else {
				break
			}
		}
	}

	// Longest streak: scan all consecutive groups.
	run := 1
	longest = 1
	for i := 1; i < len(dates); i++ {
		if dates[i].Equal(dates[i-1].AddDate(0, 0, -1)) {
			run++
			if run > longest {
				longest = run
			}
		} else {
			run = 1
		}
	}

	return current, longest
}

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
