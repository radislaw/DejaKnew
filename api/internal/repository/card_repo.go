package repository

import (
	"errors"

	"flashcards/api/internal/domain"

	"gorm.io/gorm"
)

// CardRepo is the GORM-backed implementation of service.CardRepository.
type CardRepo struct {
	db *gorm.DB
}

func NewCardRepo(db *gorm.DB) *CardRepo {
	return &CardRepo{db: db}
}

func (r *CardRepo) Create(card *domain.Card) error {
	return r.db.Create(card).Error
}

func (r *CardRepo) ListByDeck(deckID string) ([]domain.Card, error) {
	var cards []domain.Card
	err := r.db.Where("deck_id = ?", deckID).Order("position asc, created_at asc").Find(&cards).Error
	if err != nil {
		return nil, err
	}
	return cards, nil
}

func (r *CardRepo) FindByID(id string) (*domain.Card, error) {
	var card domain.Card
	err := r.db.First(&card, "id = ?", id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, domain.ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return &card, nil
}

func (r *CardRepo) Update(card *domain.Card) error {
	return r.db.Save(card).Error
}

func (r *CardRepo) Delete(id string) error {
	return r.db.Delete(&domain.Card{}, "id = ?", id).Error
}

func (r *CardRepo) CountAll() (int64, error) {
	var n int64
	err := r.db.Model(&domain.Card{}).Count(&n).Error
	return n, err
}
