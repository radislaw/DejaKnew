package repository

import (
	"errors"

	"flashcards/api/internal/domain"

	"gorm.io/gorm"
)

// DeckRepo is the GORM-backed implementation of service.DeckRepository.
type DeckRepo struct {
	db *gorm.DB
}

func NewDeckRepo(db *gorm.DB) *DeckRepo {
	return &DeckRepo{db: db}
}

func (r *DeckRepo) Create(deck *domain.Deck) error {
	return r.db.Create(deck).Error
}

func (r *DeckRepo) ListByOwner(ownerID string) ([]domain.Deck, error) {
	var decks []domain.Deck
	err := r.db.Where("owner_id = ?", ownerID).Order("created_at desc").Find(&decks).Error
	if err != nil {
		return nil, err
	}
	return decks, nil
}

func (r *DeckRepo) FindByID(id string) (*domain.Deck, error) {
	var deck domain.Deck
	err := r.db.First(&deck, "id = ?", id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, domain.ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return &deck, nil
}

func (r *DeckRepo) Update(deck *domain.Deck) error {
	return r.db.Save(deck).Error
}

func (r *DeckRepo) Delete(id string) error {
	return r.db.Delete(&domain.Deck{}, "id = ?", id).Error
}
