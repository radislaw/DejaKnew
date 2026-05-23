package service

import (
	"errors"
	"testing"

	"flashcards/api/internal/domain"
)

// fakeCardRepo is an in-memory CardRepository for unit tests.
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

func newCardService() (*CardService, *fakeCardRepo, *fakeDeckRepo) {
	cardRepo := newFakeCardRepo()
	deckRepo := newFakeDeckRepo()
	svc := NewCardService(cardRepo, deckRepo)
	return svc, cardRepo, deckRepo
}

func addDeck(repo *fakeDeckRepo, ownerID, deckID string) {
	repo.decks[deckID] = &domain.Deck{ID: deckID, OwnerID: ownerID}
}

func TestCreateCardSuccess(t *testing.T) {
	svc, _, deckRepo := newCardService()
	addDeck(deckRepo, "user-1", "deck-1")

	card, err := svc.Create("user-1", "deck-1", CardInput{Front: "Q", Back: "A"})
	if err != nil {
		t.Fatalf("create failed: %v", err)
	}
	if card.Content.Front != "Q" || card.Content.Back != "A" {
		t.Fatalf("unexpected content: %+v", card.Content)
	}
	if card.DeckID != "deck-1" {
		t.Fatalf("expected deckId deck-1, got %q", card.DeckID)
	}
}

func TestCreateCardForbiddenForOtherUser(t *testing.T) {
	svc, _, deckRepo := newCardService()
	addDeck(deckRepo, "user-1", "deck-1")

	_, err := svc.Create("user-2", "deck-1", CardInput{Front: "Q", Back: "A"})
	if !errors.Is(err, domain.ErrForbidden) {
		t.Fatalf("expected ErrForbidden, got %v", err)
	}
}

func TestCreateCardDeckNotFound(t *testing.T) {
	svc, _, _ := newCardService()

	_, err := svc.Create("user-1", "no-such-deck", CardInput{Front: "Q", Back: "A"})
	if !errors.Is(err, domain.ErrNotFound) {
		t.Fatalf("expected ErrNotFound, got %v", err)
	}
}

func TestUpdateCardSuccess(t *testing.T) {
	svc, _, deckRepo := newCardService()
	addDeck(deckRepo, "user-1", "deck-1")

	card, _ := svc.Create("user-1", "deck-1", CardInput{Front: "Old Q", Back: "Old A"})
	updated, err := svc.Update("user-1", card.ID, CardInput{Front: "New Q", Back: "New A"})
	if err != nil {
		t.Fatalf("update failed: %v", err)
	}
	if updated.Content.Front != "New Q" {
		t.Fatalf("expected New Q, got %q", updated.Content.Front)
	}
}

func TestUpdateCardForbiddenForOtherUser(t *testing.T) {
	svc, _, deckRepo := newCardService()
	addDeck(deckRepo, "user-1", "deck-1")

	card, _ := svc.Create("user-1", "deck-1", CardInput{Front: "Q", Back: "A"})
	_, err := svc.Update("user-2", card.ID, CardInput{Front: "X", Back: "Y"})
	if !errors.Is(err, domain.ErrForbidden) {
		t.Fatalf("expected ErrForbidden, got %v", err)
	}
}

func TestDeleteCardSuccess(t *testing.T) {
	svc, _, deckRepo := newCardService()
	addDeck(deckRepo, "user-1", "deck-1")

	card, _ := svc.Create("user-1", "deck-1", CardInput{Front: "Q", Back: "A"})
	if err := svc.Delete("user-1", card.ID); err != nil {
		t.Fatalf("delete failed: %v", err)
	}
	_, err := svc.cards.FindByID(card.ID)
	if !errors.Is(err, domain.ErrNotFound) {
		t.Fatalf("expected ErrNotFound after delete, got %v", err)
	}
}

func TestDeleteCardNotFound(t *testing.T) {
	svc, _, _ := newCardService()
	err := svc.Delete("user-1", "no-such-card")
	if !errors.Is(err, domain.ErrNotFound) {
		t.Fatalf("expected ErrNotFound, got %v", err)
	}
}

func TestListByDeckReturnsOwnCards(t *testing.T) {
	svc, _, deckRepo := newCardService()
	addDeck(deckRepo, "user-1", "deck-1")
	addDeck(deckRepo, "user-1", "deck-2")

	svc.Create("user-1", "deck-1", CardInput{Front: "A", Back: "1"})
	svc.Create("user-1", "deck-1", CardInput{Front: "B", Back: "2"})
	svc.Create("user-1", "deck-2", CardInput{Front: "C", Back: "3"})

	cards, err := svc.ListByDeck("deck-1")
	if err != nil {
		t.Fatalf("list failed: %v", err)
	}
	if len(cards) != 2 {
		t.Fatalf("expected 2 cards for deck-1, got %d", len(cards))
	}
}
