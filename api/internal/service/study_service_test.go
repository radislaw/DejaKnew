package service

import (
	"errors"
	"testing"
	"time"

	"flashcards/api/internal/domain"
)

// fakeStudyProgressRepo is an in-memory StudyProgressRepository for unit tests.
type fakeStudyProgressRepo struct {
	progress   map[string]*domain.StudyProgress // key: userID+":"+cardID
	reviewLogs []*domain.ReviewLog
	cards      map[string]*domain.Card // deckID -> cards for queue simulation
}

func newFakeStudyProgressRepo() *fakeStudyProgressRepo {
	return &fakeStudyProgressRepo{
		progress: map[string]*domain.StudyProgress{},
		cards:    map[string]*domain.Card{},
	}
}

func (f *fakeStudyProgressRepo) key(userID, cardID string) string {
	return userID + ":" + cardID
}

func (f *fakeStudyProgressRepo) GetDueQueue(userID, deckID string) ([]domain.Card, error) {
	now := time.Now()
	var out []domain.Card
	for _, c := range f.cards {
		if c.DeckID != deckID {
			continue
		}
		sp := f.progress[f.key(userID, c.ID)]
		if sp == nil || sp.NextReview == nil || !sp.NextReview.After(now) {
			out = append(out, *c)
		}
	}
	return out, nil
}

func (f *fakeStudyProgressRepo) GetProgress(userID, cardID string) (*domain.StudyProgress, error) {
	sp := f.progress[f.key(userID, cardID)]
	return sp, nil
}

func (f *fakeStudyProgressRepo) UpsertProgress(sp *domain.StudyProgress) error {
	f.progress[f.key(sp.UserID, sp.CardID)] = sp
	return nil
}

func (f *fakeStudyProgressRepo) CreateReviewLog(entry *domain.ReviewLog) error {
	f.reviewLogs = append(f.reviewLogs, entry)
	return nil
}

func newStudyService() (*StudyService, *fakeStudyProgressRepo, *fakeCardRepo, *fakeDeckRepo) {
	progressRepo := newFakeStudyProgressRepo()
	cardRepo := newFakeCardRepo()
	deckRepo := newFakeDeckRepo()
	svc := NewStudyService(progressRepo, cardRepo, deckRepo)
	return svc, progressRepo, cardRepo, deckRepo
}

func addCard(cardRepo *fakeCardRepo, progressRepo *fakeStudyProgressRepo, cardID, deckID string) {
	c := &domain.Card{ID: cardID, DeckID: deckID}
	cardRepo.cards[cardID] = c
	progressRepo.cards[cardID] = c
}

func TestGetQueueForbiddenForNonOwner(t *testing.T) {
	svc, progressRepo, cardRepo, deckRepo := newStudyService()
	addDeck(deckRepo, "user-1", "deck-1")
	addCard(cardRepo, progressRepo, "card-1", "deck-1")

	_, err := svc.GetQueue("user-2", "deck-1")
	if !errors.Is(err, domain.ErrForbidden) {
		t.Fatalf("expected ErrForbidden, got %v", err)
	}
}

func TestGetQueueReturnsDueCards(t *testing.T) {
	svc, progressRepo, cardRepo, deckRepo := newStudyService()
	addDeck(deckRepo, "user-1", "deck-1")
	addCard(cardRepo, progressRepo, "card-1", "deck-1")
	addCard(cardRepo, progressRepo, "card-2", "deck-1")

	// Mark card-2 as reviewed and not yet due
	future := time.Now().Add(48 * time.Hour)
	progressRepo.progress["user-1:card-2"] = &domain.StudyProgress{
		UserID: "user-1", CardID: "card-2", NextReview: &future,
	}

	queue, err := svc.GetQueue("user-1", "deck-1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(queue) != 1 || queue[0].ID != "card-1" {
		t.Fatalf("expected only card-1 in queue, got %+v", queue)
	}
}

func TestSubmitReviewFail(t *testing.T) {
	svc, progressRepo, cardRepo, deckRepo := newStudyService()
	addDeck(deckRepo, "user-1", "deck-1")
	addCard(cardRepo, progressRepo, "card-1", "deck-1")

	sp, err := svc.SubmitReview("user-1", "card-1", 0)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if sp.Reps != 0 || sp.Lapses != 1 || sp.Interval != 0 {
		t.Fatalf("expected reps=0 lapses=1 interval=0, got %+v", sp)
	}
	if len(progressRepo.reviewLogs) != 1 {
		t.Fatalf("expected 1 review log entry, got %d", len(progressRepo.reviewLogs))
	}
}

func TestSubmitReviewGoodProgression(t *testing.T) {
	svc, progressRepo, cardRepo, deckRepo := newStudyService()
	addDeck(deckRepo, "user-1", "deck-1")
	addCard(cardRepo, progressRepo, "card-1", "deck-1")

	sp, _ := svc.SubmitReview("user-1", "card-1", 2) // reps=1, interval=1
	if sp.Reps != 1 || sp.Interval != 1 {
		t.Fatalf("first review: expected reps=1 interval=1, got %+v", sp)
	}

	sp, _ = svc.SubmitReview("user-1", "card-1", 2) // reps=2, interval=3
	if sp.Reps != 2 || sp.Interval != 3 {
		t.Fatalf("second review: expected reps=2 interval=3, got %+v", sp)
	}

	sp, _ = svc.SubmitReview("user-1", "card-1", 3) // reps=3, interval=round(3 * ease * 1.3)
	if sp.Reps != 3 || sp.Interval <= 3 {
		t.Fatalf("third review: expected reps=3 interval>3, got %+v", sp)
	}
}

func TestSubmitReviewForbiddenForNonOwner(t *testing.T) {
	svc, progressRepo, cardRepo, deckRepo := newStudyService()
	addDeck(deckRepo, "user-1", "deck-1")
	addCard(cardRepo, progressRepo, "card-1", "deck-1")

	_, err := svc.SubmitReview("user-2", "card-1", 2)
	if !errors.Is(err, domain.ErrForbidden) {
		t.Fatalf("expected ErrForbidden, got %v", err)
	}
}

func TestApplyReviewEaseFloor(t *testing.T) {
	sp := &domain.StudyProgress{Ease: 1.4, Reps: 5, Interval: 10}
	// multiple fails should not drop ease below 1.3
	for i := 0; i < 5; i++ {
		applyReview(sp, 0)
	}
	if sp.Ease < 1.3 {
		t.Fatalf("ease dropped below floor: %f", sp.Ease)
	}
}
