package service

import (
	"errors"
	"testing"

	"flashcards/api/internal/domain"
)

// fakeDeckRepo is an in-memory DeckRepository for unit tests.
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

func TestCreateAndGetDeck(t *testing.T) {
	s := NewDeckService(newFakeDeckRepo())
	deck, err := s.Create("owner-1", DeckInput{Title: "Go basics"})
	if err != nil {
		t.Fatalf("create failed: %v", err)
	}
	got, err := s.Get("owner-1", deck.ID)
	if err != nil {
		t.Fatalf("get failed: %v", err)
	}
	if got.Title != "Go basics" {
		t.Fatalf("unexpected title: %q", got.Title)
	}
}

func TestGetRejectsOtherUsersDeck(t *testing.T) {
	s := NewDeckService(newFakeDeckRepo())
	deck, _ := s.Create("owner-1", DeckInput{Title: "Private"})
	_, err := s.Get("owner-2", deck.ID)
	if !errors.Is(err, domain.ErrForbidden) {
		t.Fatalf("expected ErrForbidden, got %v", err)
	}
}

func TestGetMissingDeckReturnsNotFound(t *testing.T) {
	s := NewDeckService(newFakeDeckRepo())
	_, err := s.Get("owner-1", "no-such-id")
	if !errors.Is(err, domain.ErrNotFound) {
		t.Fatalf("expected ErrNotFound, got %v", err)
	}
}

func TestListReturnsOnlyOwnDecks(t *testing.T) {
	s := NewDeckService(newFakeDeckRepo())
	s.Create("owner-1", DeckInput{Title: "A"})
	s.Create("owner-1", DeckInput{Title: "B"})
	s.Create("owner-2", DeckInput{Title: "C"})
	decks, err := s.List("owner-1")
	if err != nil {
		t.Fatalf("list failed: %v", err)
	}
	if len(decks) != 2 {
		t.Fatalf("expected 2 decks for owner-1, got %d", len(decks))
	}
}

func TestUpdateChangesFields(t *testing.T) {
	s := NewDeckService(newFakeDeckRepo())
	deck, _ := s.Create("owner-1", DeckInput{Title: "Old"})
	updated, err := s.Update("owner-1", deck.ID, DeckInput{Title: "New"})
	if err != nil {
		t.Fatalf("update failed: %v", err)
	}
	if updated.Title != "New" {
		t.Fatalf("expected title New, got %q", updated.Title)
	}
}

func TestUpdateRejectsOtherUsersDeck(t *testing.T) {
	s := NewDeckService(newFakeDeckRepo())
	deck, _ := s.Create("owner-1", DeckInput{Title: "X"})
	_, err := s.Update("owner-2", deck.ID, DeckInput{Title: "hacked"})
	if !errors.Is(err, domain.ErrForbidden) {
		t.Fatalf("expected ErrForbidden, got %v", err)
	}
}

func TestDeleteRemovesDeck(t *testing.T) {
	s := NewDeckService(newFakeDeckRepo())
	deck, _ := s.Create("owner-1", DeckInput{Title: "Temp"})
	if err := s.Delete("owner-1", deck.ID); err != nil {
		t.Fatalf("delete failed: %v", err)
	}
	if _, err := s.Get("owner-1", deck.ID); !errors.Is(err, domain.ErrNotFound) {
		t.Fatalf("expected ErrNotFound after delete, got %v", err)
	}
}
