package service

import (
	"flashcards/api/internal/domain"

	"github.com/google/uuid"
)

// CardRepository is the persistence dependency of CardService.
type CardRepository interface {
	Create(card *domain.Card) error
	ListByDeck(deckID string) ([]domain.Card, error)
	FindByID(id string) (*domain.Card, error)
	Update(card *domain.Card) error
	Delete(id string) error
}

// deckOwnerChecker is the minimal deck dependency needed to verify ownership.
type deckOwnerChecker interface {
	FindByID(id string) (*domain.Deck, error)
}

type CardService struct {
	cards CardRepository
	decks deckOwnerChecker
}

func NewCardService(cards CardRepository, decks deckOwnerChecker) *CardService {
	return &CardService{cards: cards, decks: decks}
}

// CardInput carries the editable fields of a Basic card.
type CardInput struct {
	Front string
	Back  string
}

func (s *CardService) Create(ownerID, deckID string, in CardInput) (*domain.Card, error) {
	deck, err := s.decks.FindByID(deckID)
	if err != nil {
		return nil, err
	}
	if deck.OwnerID != ownerID {
		return nil, domain.ErrForbidden
	}
	card := &domain.Card{
		ID:     uuid.NewString(),
		DeckID: deckID,
		Type:   "basic",
		Content: domain.CardContent{
			Front: in.Front,
			Back:  in.Back,
		},
	}
	if err := s.cards.Create(card); err != nil {
		return nil, err
	}
	return card, nil
}

func (s *CardService) ListByDeck(deckID string) ([]domain.Card, error) {
	return s.cards.ListByDeck(deckID)
}

func (s *CardService) Update(ownerID, cardID string, in CardInput) (*domain.Card, error) {
	card, err := s.owned(ownerID, cardID)
	if err != nil {
		return nil, err
	}
	card.Content.Front = in.Front
	card.Content.Back = in.Back
	if err := s.cards.Update(card); err != nil {
		return nil, err
	}
	return card, nil
}

func (s *CardService) Delete(ownerID, cardID string) error {
	if _, err := s.owned(ownerID, cardID); err != nil {
		return err
	}
	return s.cards.Delete(cardID)
}

// owned loads a card and verifies ownerID owns its parent deck.
func (s *CardService) owned(ownerID, cardID string) (*domain.Card, error) {
	card, err := s.cards.FindByID(cardID)
	if err != nil {
		return nil, err
	}
	deck, err := s.decks.FindByID(card.DeckID)
	if err != nil {
		return nil, err
	}
	if deck.OwnerID != ownerID {
		return nil, domain.ErrForbidden
	}
	return card, nil
}
