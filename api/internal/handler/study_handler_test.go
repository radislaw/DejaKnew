package handler

import (
	"net/http"
	"testing"
	"time"

	"flashcards/api/internal/domain"
)

// fakeStudyProgressRepo is an in-memory StudyProgressRepository for handler tests.
type fakeStudyProgressRepo struct {
	progress   map[string]*domain.StudyProgress
	reviewLogs []*domain.ReviewLog
	cardIndex  map[string]*domain.Card // cardID -> card, for queue filtering
}

func newFakeStudyProgressRepo() *fakeStudyProgressRepo {
	return &fakeStudyProgressRepo{
		progress:  map[string]*domain.StudyProgress{},
		cardIndex: map[string]*domain.Card{},
	}
}

func (f *fakeStudyProgressRepo) progressKey(userID, cardID string) string {
	return userID + ":" + cardID
}

func (f *fakeStudyProgressRepo) GetDueQueue(userID, deckID string) ([]domain.Card, error) {
	now := time.Now()
	var out []domain.Card
	for _, c := range f.cardIndex {
		if c.DeckID != deckID {
			continue
		}
		sp := f.progress[f.progressKey(userID, c.ID)]
		if sp == nil || sp.NextReview == nil || !sp.NextReview.After(now) {
			out = append(out, *c)
		}
	}
	return out, nil
}

func (f *fakeStudyProgressRepo) GetProgress(userID, cardID string) (*domain.StudyProgress, error) {
	return f.progress[f.progressKey(userID, cardID)], nil
}

func (f *fakeStudyProgressRepo) UpsertProgress(sp *domain.StudyProgress) error {
	f.progress[f.progressKey(sp.UserID, sp.CardID)] = sp
	return nil
}

func (f *fakeStudyProgressRepo) CreateReviewLog(entry *domain.ReviewLog) error {
	f.reviewLogs = append(f.reviewLogs, entry)
	return nil
}

func TestGetQueueReturns200(t *testing.T) {
	r := testRouter()
	token := registerUser(t, r, "study1@b.com")

	deckResp := doJSON(t, r, "POST", "/api/v1/decks", `{"title":"SRS Deck"}`, token)
	var deck struct{ ID string }
	decodeJSON(t, deckResp, &deck)

	doJSON(t, r, "POST", "/api/v1/decks/"+deck.ID+"/cards",
		`{"front":"Q","back":"A"}`, token)

	w := doJSON(t, r, "GET", "/api/v1/study/"+deck.ID+"/queue", "", token)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
	}
}

func TestGetQueueForbiddenForNonOwner(t *testing.T) {
	r := testRouter()
	tokenA := registerUser(t, r, "study2a@b.com")
	tokenB := registerUser(t, r, "study2b@b.com")

	deckResp := doJSON(t, r, "POST", "/api/v1/decks", `{"title":"Deck"}`, tokenA)
	var deck struct{ ID string }
	decodeJSON(t, deckResp, &deck)

	w := doJSON(t, r, "GET", "/api/v1/study/"+deck.ID+"/queue", "", tokenB)
	if w.Code != http.StatusForbidden {
		t.Fatalf("expected 403, got %d: %s", w.Code, w.Body.String())
	}
}

func TestSubmitReviewReturns200(t *testing.T) {
	r := testRouter()
	token := registerUser(t, r, "study3@b.com")

	deckResp := doJSON(t, r, "POST", "/api/v1/decks", `{"title":"Deck"}`, token)
	var deck struct{ ID string }
	decodeJSON(t, deckResp, &deck)

	cardResp := doJSON(t, r, "POST", "/api/v1/decks/"+deck.ID+"/cards",
		`{"front":"Q","back":"A"}`, token)
	var card struct{ ID string }
	decodeJSON(t, cardResp, &card)

	w := doJSON(t, r, "POST", "/api/v1/study/"+card.ID+"/review", `{"quality":2}`, token)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
	}
}

func TestSubmitReviewRejectsBadQuality(t *testing.T) {
	r := testRouter()
	token := registerUser(t, r, "study4@b.com")

	deckResp := doJSON(t, r, "POST", "/api/v1/decks", `{"title":"Deck"}`, token)
	var deck struct{ ID string }
	decodeJSON(t, deckResp, &deck)

	cardResp := doJSON(t, r, "POST", "/api/v1/decks/"+deck.ID+"/cards",
		`{"front":"Q","back":"A"}`, token)
	var card struct{ ID string }
	decodeJSON(t, cardResp, &card)

	w := doJSON(t, r, "POST", "/api/v1/study/"+card.ID+"/review", `{"quality":5}`, token)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400 for quality=5, got %d: %s", w.Code, w.Body.String())
	}
}
