package service

import (
	"math"
	"time"

	"flashcards/api/internal/domain"

	"github.com/google/uuid"
)

// StudyProgressRepository is the persistence dependency of StudyService.
type StudyProgressRepository interface {
	GetDueQueue(userID, deckID string) ([]domain.Card, error)
	GetProgress(userID, cardID string) (*domain.StudyProgress, error)
	UpsertProgress(sp *domain.StudyProgress) error
	CreateReviewLog(entry *domain.ReviewLog) error
}

type StudyService struct {
	progress StudyProgressRepository
	cards    CardRepository
	decks    deckOwnerChecker
}

func NewStudyService(progress StudyProgressRepository, cards CardRepository, decks deckOwnerChecker) *StudyService {
	return &StudyService{progress: progress, cards: cards, decks: decks}
}

// GetQueue returns cards in a deck that are new or due for review today.
func (s *StudyService) GetQueue(userID, deckID string) ([]domain.Card, error) {
	deck, err := s.decks.FindByID(deckID)
	if err != nil {
		return nil, err
	}
	if deck.OwnerID != userID {
		return nil, domain.ErrForbidden
	}
	return s.progress.GetDueQueue(userID, deckID)
}

// SubmitReview applies the SRS algorithm for one review event and persists it.
func (s *StudyService) SubmitReview(userID, cardID string, quality int) (*domain.StudyProgress, error) {
	card, err := s.cards.FindByID(cardID)
	if err != nil {
		return nil, err
	}
	deck, err := s.decks.FindByID(card.DeckID)
	if err != nil {
		return nil, err
	}
	if deck.OwnerID != userID {
		return nil, domain.ErrForbidden
	}

	sp, err := s.progress.GetProgress(userID, cardID)
	if err != nil {
		return nil, err
	}
	if sp == nil {
		sp = &domain.StudyProgress{
			ID:       uuid.NewString(),
			UserID:   userID,
			CardID:   cardID,
			Interval: 0,
			Ease:     2.5,
			Reps:     0,
			Lapses:   0,
		}
	}

	applyReview(sp, quality)

	if err := s.progress.UpsertProgress(sp); err != nil {
		return nil, err
	}

	entry := &domain.ReviewLog{
		ID:         uuid.NewString(),
		UserID:     userID,
		CardID:     cardID,
		Quality:    quality,
		ReviewedAt: time.Now(),
	}
	if err := s.progress.CreateReviewLog(entry); err != nil {
		return nil, err
	}

	return sp, nil
}

// applyReview mutates sp in place using the SM2-like algorithm from the prototype.
func applyReview(sp *domain.StudyProgress, quality int) {
	now := time.Now()

	if quality == 0 {
		sp.Reps = 0
		sp.Interval = 0
		sp.Lapses++
		sp.Ease = math.Max(1.3, sp.Ease-0.2)
	} else {
		sp.Reps++
		delta := float64(3 - quality)
		sp.Ease = math.Max(1.3, sp.Ease+0.1-delta*(0.08+delta*0.02))
		switch sp.Reps {
		case 1:
			sp.Interval = 1
		case 2:
			if quality == 1 {
				sp.Interval = 1
			} else {
				sp.Interval = 3
			}
		default:
			var modifier float64
			switch quality {
			case 1:
				modifier = 0.6
			case 3:
				modifier = 1.3
			default:
				modifier = 1.0
			}
			sp.Interval = int(math.Round(float64(sp.Interval) * sp.Ease * modifier))
		}
	}

	nextReview := now.Add(time.Duration(sp.Interval) * 24 * time.Hour)
	sp.NextReview = &nextReview
	sp.LastReview = &now
}
