package service

import (
	"flashcards/api/internal/domain"

	"github.com/google/uuid"
)

// DeckRepository is the persistence dependency of DeckService.
type DeckRepository interface {
	Create(deck *domain.Deck) error
	ListByOwner(ownerID string) ([]domain.Deck, error)
	FindByID(id string) (*domain.Deck, error)
	Update(deck *domain.Deck) error
	Delete(id string) error
}

type DeckService struct {
	decks DeckRepository
}

func NewDeckService(decks DeckRepository) *DeckService {
	return &DeckService{decks: decks}
}

// DeckInput carries the editable fields of a deck.
type DeckInput struct {
	Title       string
	Description string
	Tags        []string
	IsPublic    bool
}

func (s *DeckService) Create(ownerID string, in DeckInput) (*domain.Deck, error) {
	deck := &domain.Deck{
		ID:          uuid.NewString(),
		OwnerID:     ownerID,
		Title:       in.Title,
		Description: in.Description,
		Tags:        in.Tags,
		IsPublic:    in.IsPublic,
	}
	if err := s.decks.Create(deck); err != nil {
		return nil, err
	}
	return deck, nil
}

func (s *DeckService) List(ownerID string) ([]domain.Deck, error) {
	return s.decks.ListByOwner(ownerID)
}

func (s *DeckService) Get(ownerID, deckID string) (*domain.Deck, error) {
	return s.owned(ownerID, deckID)
}

func (s *DeckService) Update(ownerID, deckID string, in DeckInput) (*domain.Deck, error) {
	deck, err := s.owned(ownerID, deckID)
	if err != nil {
		return nil, err
	}
	deck.Title = in.Title
	deck.Description = in.Description
	deck.Tags = in.Tags
	deck.IsPublic = in.IsPublic
	if err := s.decks.Update(deck); err != nil {
		return nil, err
	}
	return deck, nil
}

func (s *DeckService) Delete(ownerID, deckID string) error {
	if _, err := s.owned(ownerID, deckID); err != nil {
		return err
	}
	return s.decks.Delete(deckID)
}

// owned loads a deck and confirms it belongs to ownerID. It returns
// ErrNotFound if the deck is missing and ErrForbidden if it belongs to
// someone else — the single ownership gate reused by Get/Update/Delete.
func (s *DeckService) owned(ownerID, deckID string) (*domain.Deck, error) {
	deck, err := s.decks.FindByID(deckID)
	if err != nil {
		return nil, err
	}
	if deck.OwnerID != ownerID {
		return nil, domain.ErrForbidden
	}
	return deck, nil
}
