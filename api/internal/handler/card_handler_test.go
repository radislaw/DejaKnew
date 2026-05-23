package handler

import (
	"net/http"
	"strings"
	"testing"

	"flashcards/api/internal/domain"
)

// fakeCardRepo is an in-memory service.CardRepository for handler tests.
type fakeCardRepo struct {
	cards map[string]*domain.Card
}

func newFakeCardRepo() *fakeCardRepo {
	return &fakeCardRepo{cards: map[string]*domain.Card{}}
}

func (f *fakeCardRepo) Create(c *domain.Card) error {
	f.cards[c.ID] = c
	return nil
}

func (f *fakeCardRepo) ListByDeck(deckID string) ([]domain.Card, error) {
	var out []domain.Card
	for _, c := range f.cards {
		if c.DeckID == deckID {
			out = append(out, *c)
		}
	}
	return out, nil
}

func (f *fakeCardRepo) FindByID(id string) (*domain.Card, error) {
	if c, ok := f.cards[id]; ok {
		return c, nil
	}
	return nil, domain.ErrNotFound
}

func (f *fakeCardRepo) Update(c *domain.Card) error {
	f.cards[c.ID] = c
	return nil
}

func (f *fakeCardRepo) Delete(id string) error {
	delete(f.cards, id)
	return nil
}

func TestCreateCardReturns201(t *testing.T) {
	r := testRouter()
	token := registerUser(t, r, "c@b.com")

	created := doJSON(t, r, "POST", "/api/v1/decks", `{"title":"Deck"}`, token)
	if created.Code != http.StatusCreated {
		t.Fatalf("deck create failed: %d %s", created.Code, created.Body.String())
	}

	var deck struct{ ID string }
	decodeJSON(t, created, &deck)

	w := doJSON(t, r, "POST", "/api/v1/decks/"+deck.ID+"/cards",
		`{"front":"What is Go?","back":"A compiled language"}`, token)
	if w.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d: %s", w.Code, w.Body.String())
	}
}

func TestCreateCardRejectsMissingFields(t *testing.T) {
	r := testRouter()
	token := registerUser(t, r, "d@b.com")

	created := doJSON(t, r, "POST", "/api/v1/decks", `{"title":"Deck"}`, token)
	var deck struct{ ID string }
	decodeJSON(t, created, &deck)

	w := doJSON(t, r, "POST", "/api/v1/decks/"+deck.ID+"/cards", `{"front":"Q"}`, token)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400 for missing back, got %d", w.Code)
	}
}

func TestCreateCardForbiddenForOtherUser(t *testing.T) {
	r := testRouter()
	tokenA := registerUser(t, r, "e@b.com")
	tokenB := registerUser(t, r, "f@b.com")

	created := doJSON(t, r, "POST", "/api/v1/decks", `{"title":"A's deck"}`, tokenA)
	var deck struct{ ID string }
	decodeJSON(t, created, &deck)

	w := doJSON(t, r, "POST", "/api/v1/decks/"+deck.ID+"/cards",
		`{"front":"Q","back":"A"}`, tokenB)
	if w.Code != http.StatusForbidden {
		t.Fatalf("expected 403, got %d", w.Code)
	}
}

func TestGetDeckIncludesCards(t *testing.T) {
	r := testRouter()
	token := registerUser(t, r, "g@b.com")

	created := doJSON(t, r, "POST", "/api/v1/decks", `{"title":"Deck"}`, token)
	var deck struct{ ID string }
	decodeJSON(t, created, &deck)

	doJSON(t, r, "POST", "/api/v1/decks/"+deck.ID+"/cards",
		`{"front":"Q","back":"A"}`, token)

	w := doJSON(t, r, "GET", "/api/v1/decks/"+deck.ID, "", token)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
	if !strings.Contains(w.Body.String(), "cards") {
		t.Fatalf("expected cards field in response, got %s", w.Body.String())
	}
}

func TestDeleteCardReturns204(t *testing.T) {
	r := testRouter()
	token := registerUser(t, r, "h@b.com")

	created := doJSON(t, r, "POST", "/api/v1/decks", `{"title":"Deck"}`, token)
	var deck struct{ ID string }
	decodeJSON(t, created, &deck)

	cardResp := doJSON(t, r, "POST", "/api/v1/decks/"+deck.ID+"/cards",
		`{"front":"Q","back":"A"}`, token)
	var card struct{ ID string }
	decodeJSON(t, cardResp, &card)

	w := doJSON(t, r, "DELETE", "/api/v1/cards/"+card.ID, "", token)
	if w.Code != http.StatusNoContent {
		t.Fatalf("expected 204, got %d: %s", w.Code, w.Body.String())
	}
}

func TestExportDeckCSV(t *testing.T) {
	r := testRouter()
	token := registerUser(t, r, "exp-csv@b.com")

	created := doJSON(t, r, "POST", "/api/v1/decks", `{"title":"Go"}`, token)
	var deck struct{ ID string }
	decodeJSON(t, created, &deck)

	doJSON(t, r, "POST", "/api/v1/decks/"+deck.ID+"/cards",
		`{"front":"Q","back":"A"}`, token)

	w := doJSON(t, r, "GET", "/api/v1/decks/"+deck.ID+"/export?format=csv", "", token)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
	}
	if ct := w.Header().Get("Content-Type"); !strings.Contains(ct, "text/csv") {
		t.Errorf("unexpected Content-Type: %s", ct)
	}
	body := w.Body.String()
	if !strings.Contains(body, "front") || !strings.Contains(body, "Q") {
		t.Errorf("CSV missing expected content: %s", body)
	}
}

func TestExportDeckJSON(t *testing.T) {
	r := testRouter()
	token := registerUser(t, r, "exp-json@b.com")

	created := doJSON(t, r, "POST", "/api/v1/decks", `{"title":"Go"}`, token)
	var deck struct{ ID string }
	decodeJSON(t, created, &deck)

	doJSON(t, r, "POST", "/api/v1/decks/"+deck.ID+"/cards",
		`{"front":"Q","back":"A"}`, token)

	w := doJSON(t, r, "GET", "/api/v1/decks/"+deck.ID+"/export?format=json", "", token)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
	}
	if ct := w.Header().Get("Content-Type"); !strings.Contains(ct, "application/json") {
		t.Errorf("unexpected Content-Type: %s", ct)
	}
	body := w.Body.String()
	if !strings.Contains(body, `"name"`) || !strings.Contains(body, `"Q"`) {
		t.Errorf("JSON missing expected content: %s", body)
	}
}

func TestExportDeckAnkiCSV(t *testing.T) {
	r := testRouter()
	token := registerUser(t, r, "exp-anki@b.com")

	created := doJSON(t, r, "POST", "/api/v1/decks", `{"title":"Go"}`, token)
	var deck struct{ ID string }
	decodeJSON(t, created, &deck)

	doJSON(t, r, "POST", "/api/v1/decks/"+deck.ID+"/cards",
		`{"front":"Question","back":"Answer"}`, token)

	w := doJSON(t, r, "GET", "/api/v1/decks/"+deck.ID+"/export?format=anki-csv", "", token)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
	}
	body := w.Body.String()
	if !strings.Contains(body, "Question\tAnswer") {
		t.Errorf("Anki CSV missing tab-separated content: %s", body)
	}
}

func TestExportDeckBadFormat(t *testing.T) {
	r := testRouter()
	token := registerUser(t, r, "exp-bad@b.com")

	created := doJSON(t, r, "POST", "/api/v1/decks", `{"title":"Go"}`, token)
	var deck struct{ ID string }
	decodeJSON(t, created, &deck)

	w := doJSON(t, r, "GET", "/api/v1/decks/"+deck.ID+"/export?format=xlsx", "", token)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w.Code)
	}
}

func TestExportAll(t *testing.T) {
	r := testRouter()
	token := registerUser(t, r, "exp-all@b.com")

	d1 := doJSON(t, r, "POST", "/api/v1/decks", `{"title":"Deck1"}`, token)
	var deck1 struct{ ID string }
	decodeJSON(t, d1, &deck1)
	doJSON(t, r, "POST", "/api/v1/decks/"+deck1.ID+"/cards",
		`{"front":"F1","back":"B1"}`, token)

	w := doJSON(t, r, "GET", "/api/v1/export?format=csv", "", token)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
	}
	body := w.Body.String()
	if !strings.Contains(body, "F1") {
		t.Errorf("export missing card content: %s", body)
	}
}
