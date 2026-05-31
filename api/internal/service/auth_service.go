package service

import (
	"errors"

	"flashcards/api/internal/auth"
	"flashcards/api/internal/domain"

	"github.com/google/uuid"
)

// UserRepository is the persistence dependency of AuthService. It is declared
// in the consumer (this package) and satisfied by repository.UserRepo —
// this is the canonical repository-interface pattern for the project.
type UserRepository interface {
	Create(user *domain.User) error
	FindByEmail(email string) (*domain.User, error)
	FindByID(id string) (*domain.User, error)
	FindByProvider(provider, providerID string) (*domain.User, error)
	Update(user *domain.User) error
	Delete(id string) error
}

type AuthService struct {
	users  UserRepository
	tokens *auth.TokenManager
}

func NewAuthService(users UserRepository, tokens *auth.TokenManager) *AuthService {
	return &AuthService{users: users, tokens: tokens}
}

// Register creates a new account and returns a signed token.
func (s *AuthService) Register(email, password, locale string) (string, *domain.User, error) {
	existing, err := s.users.FindByEmail(email)
	if err != nil && !errors.Is(err, domain.ErrNotFound) {
		return "", nil, err
	}
	if existing != nil {
		return "", nil, domain.ErrUserExists
	}

	hash, err := auth.HashPassword(password)
	if err != nil {
		return "", nil, err
	}

	if locale != "ru" {
		locale = "en"
	}
	user := &domain.User{
		ID:           uuid.NewString(),
		Email:        email,
		PasswordHash: hash,
		Locale:       locale,
	}
	if err := s.users.Create(user); err != nil {
		return "", nil, err
	}

	token, err := s.tokens.Issue(user.ID)
	if err != nil {
		return "", nil, err
	}
	return token, user, nil
}

// Login verifies credentials and returns a signed token.
func (s *AuthService) Login(email, password string) (string, *domain.User, error) {
	user, err := s.users.FindByEmail(email)
	if errors.Is(err, domain.ErrNotFound) {
		return "", nil, domain.ErrInvalidCredentials
	}
	if err != nil {
		return "", nil, err
	}
	if !auth.CheckPassword(user.PasswordHash, password) {
		return "", nil, domain.ErrInvalidCredentials
	}

	token, err := s.tokens.Issue(user.ID)
	if err != nil {
		return "", nil, err
	}
	return token, user, nil
}

// CurrentUser loads the user behind an authenticated request.
func (s *AuthService) CurrentUser(id string) (*domain.User, error) {
	return s.users.FindByID(id)
}

// UpdateProfile changes email and/or locale for a user.
// Email uniqueness is enforced; OAuth users may not change their email.
func (s *AuthService) UpdateProfile(id, email, locale string) (*domain.User, error) {
	user, err := s.users.FindByID(id)
	if err != nil {
		return nil, err
	}

	if locale != "ru" {
		locale = "en"
	}
	user.Locale = locale

	if email != "" && email != user.Email {
		existing, err := s.users.FindByEmail(email)
		if err != nil && !errors.Is(err, domain.ErrNotFound) {
			return nil, err
		}
		if existing != nil {
			return nil, domain.ErrUserExists
		}
		user.Email = email
	}

	if err := s.users.Update(user); err != nil {
		return nil, err
	}
	return user, nil
}

// ChangePassword verifies the current password and replaces it. Forbidden for OAuth users.
func (s *AuthService) ChangePassword(id, current, newPass string) error {
	user, err := s.users.FindByID(id)
	if err != nil {
		return err
	}
	if user.PasswordHash == "" {
		return domain.ErrForbidden
	}
	if !auth.CheckPassword(user.PasswordHash, current) {
		return domain.ErrInvalidCredentials
	}
	hash, err := auth.HashPassword(newPass)
	if err != nil {
		return err
	}
	user.PasswordHash = hash
	return s.users.Update(user)
}

// DeleteAccount permanently removes the user record.
func (s *AuthService) DeleteAccount(id string) error {
	return s.users.Delete(id)
}

// LoginWithOAuth finds or creates a user from an OAuth provider and issues a JWT.
func (s *AuthService) LoginWithOAuth(provider, providerID, email, locale string) (string, *domain.User, error) {
	user, err := s.users.FindByProvider(provider, providerID)
	if err != nil && !errors.Is(err, domain.ErrNotFound) {
		return "", nil, err
	}
	if user != nil {
		token, err := s.tokens.Issue(user.ID)
		return token, user, err
	}

	if locale != "ru" {
		locale = "en"
	}
	user = &domain.User{
		ID:         uuid.NewString(),
		Email:      email,
		Provider:   provider,
		ProviderID: providerID,
		Locale:     locale,
	}
	if err := s.users.Create(user); err != nil {
		return "", nil, err
	}
	token, err := s.tokens.Issue(user.ID)
	return token, user, err
}
