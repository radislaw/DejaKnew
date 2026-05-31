package repository

import (
	"errors"

	"flashcards/api/internal/domain"

	"gorm.io/gorm"
)

// UserRepo is the GORM-backed implementation of service.UserRepository.
type UserRepo struct {
	db *gorm.DB
}

func NewUserRepo(db *gorm.DB) *UserRepo {
	return &UserRepo{db: db}
}

func (r *UserRepo) Create(user *domain.User) error {
	return r.db.Create(user).Error
}

func (r *UserRepo) FindByEmail(email string) (*domain.User, error) {
	var user domain.User
	err := r.db.Where("email = ?", email).First(&user).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, domain.ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepo) FindByID(id string) (*domain.User, error) {
	var user domain.User
	err := r.db.First(&user, "id = ?", id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, domain.ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepo) FindByProvider(provider, providerID string) (*domain.User, error) {
	var user domain.User
	err := r.db.Where("provider = ? AND provider_id = ?", provider, providerID).First(&user).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, domain.ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepo) Update(user *domain.User) error {
	return r.db.Save(user).Error
}

func (r *UserRepo) Delete(id string) error {
	return r.db.Delete(&domain.User{}, "id = ?", id).Error
}
