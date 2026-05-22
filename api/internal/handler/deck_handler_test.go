package handler

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"flashcards/api/internal/domain"
)

// fakeDeckRepo is an in-memory service.DeckRepository for handler tests.
type fakeDeckRepo struct {
	decks map[string]*domain.Deck
}

func newFakeDeckRepo() *fakeDeckRepo {
	return &fakeDeckRepo{decks: map[string]*domain.Deck{}}
}

func (f *fakeDeckRepo) Create(d *domain.Deck) error {
	f.decks[d.ID] = d
	return nil
}

func (f *fakeDeckRepo) ListByOwner(ownerID string) ([]domain.Deck, error) {
	var out []domain.Deck
	for _, d := range f.decks {
		if d.OwnerID == ownerID {
			out = append(out, *d)
		}
	}
	return out, nil
}

func (f *fakeDeckRepo) FindByID(id string) (*domain.Deck, error) {
	if d, ok := f.decks[id]; ok {
		return d, nil
	}
	return nil, domain.ErrNotFound
}

func (f *fakeDeckRepo) Update(d *domain.Deck) error {
	f.decks[d.ID] = d
	return nil
}

func (f *fakeDeckRepo) Delete(id string) error {
	delete(f.decks, id)
	return nil
}

func deckIDFrom(t *testing.T, w *httptest.ResponseRecorder) string {
	t.Helper()
	var resp struct {
		ID string `json:"id"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("decode deck: %v", err)
	}
	return resp.ID
}

func TestCreateDeckReturns201(t *testing.T) {
	r := testRouter()
	token := registerUser(t, r, "a@b.com")
	w := doJSON(t, r, "POST", "/api/v1/decks", `{"title":"Go basics","tags":["go"]}`, token)
	if w.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d: %s", w.Code, w.Body.String())
	}
}

func TestCreateDeckRejectsEmptyTitle(t *testing.T) {
	r := testRouter()
	token := registerUser(t, r, "a@b.com")
	w := doJSON(t, r, "POST", "/api/v1/decks", `{"title":""}`, token)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400 for empty title, got %d", w.Code)
	}
}

func TestDecksRequireAuth(t *testing.T) {
	r := testRouter()
	w := doJSON(t, r, "GET", "/api/v1/decks", "", "")
	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401 without token, got %d", w.Code)
	}
}

func TestListReturnsCreatedDeck(t *testing.T) {
	r := testRouter()
	token := registerUser(t, r, "a@b.com")
	doJSON(t, r, "POST", "/api/v1/decks", `{"title":"Deck A"}`, token)
	w := doJSON(t, r, "GET", "/api/v1/decks", "", token)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
	if !strings.Contains(w.Body.String(), "Deck A") {
		t.Fatalf("expected list to contain Deck A, got %s", w.Body.String())
	}
}

func TestGetMissingDeckReturns404(t *testing.T) {
	r := testRouter()
	token := registerUser(t, r, "a@b.com")
	w := doJSON(t, r, "GET", "/api/v1/decks/no-such-id", "", token)
	if w.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", w.Code)
	}
}

func TestGetOtherUsersDeckReturns403(t *testing.T) {
	r := testRouter()
	tokenA := registerUser(t, r, "a@b.com")
	tokenB := registerUser(t, r, "b@b.com")
	created := doJSON(t, r, "POST", "/api/v1/decks", `{"title":"A's deck"}`, tokenA)
	id := deckIDFrom(t, created)

	w := doJSON(t, r, "GET", "/api/v1/decks/"+id, "", tokenB)
	if w.Code != http.StatusForbidden {
		t.Fatalf("expected 403 for another user's deck, got %d", w.Code)
	}
}

func TestUpdateAndDeleteDeck(t *testing.T) {
	r := testRouter()
	token := registerUser(t, r, "a@b.com")
	created := doJSON(t, r, "POST", "/api/v1/decks", `{"title":"Old"}`, token)
	id := deckIDFrom(t, created)

	upd := doJSON(t, r, "PATCH", "/api/v1/decks/"+id, `{"title":"New"}`, token)
	if upd.Code != http.StatusOK {
		t.Fatalf("expected 200 from update, got %d: %s", upd.Code, upd.Body.String())
	}

	del := doJSON(t, r, "DELETE", "/api/v1/decks/"+id, "", token)
	if del.Code != http.StatusNoContent {
		t.Fatalf("expected 204 from delete, got %d", del.Code)
	}

	get := doJSON(t, r, "GET", "/api/v1/decks/"+id, "", token)
	if get.Code != http.StatusNotFound {
		t.Fatalf("expected 404 after delete, got %d", get.Code)
	}
}
